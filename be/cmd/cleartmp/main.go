package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/reshap0318/hadirYuk/internal/helpers"
)

func main() {
	dir := flag.String("dir", helpers.GetEnv("TMP_DIR", "storage/tmp"), "Directory to clean")
	age := flag.Int("age", 24, "Delete files older than this many hours")
	dryRun := flag.Bool("dry-run", false, "Print files that would be deleted without deleting")
	flag.Parse()

	logger, err := helpers.NewLoggerWithSuffix("storage/logs", "cleartmp")
	if err != nil {
		log.Fatalf("[cleartmp] failed to initialize logger: %v", err)
	}
	defer logger.Close()

	threshold := time.Now().Add(-time.Duration(*age) * time.Hour)

	logger.Printf("Scanning: %s (files older than %dh)", *dir, *age)
	if *dryRun {
		logger.Println("DRY RUN — no files will be deleted")
	}

	deleted, skipped, errors := 0, 0, 0

	err = filepath.WalkDir(*dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			logger.Printf("Error accessing %s: %v", path, err)
			errors++
			return nil
		}

		// skip directories
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			logger.Printf("Error reading info %s: %v", path, err)
			errors++
			return nil
		}

		if info.ModTime().After(threshold) {
			skipped++
			return nil
		}

		if *dryRun {
			logger.Printf("  [dry-run] would delete: %s (modified: %s)", path, info.ModTime().Format(time.DateTime))
			deleted++
			return nil
		}

		if err := helpers.DeleteFile(path); err != nil {
			logger.Printf("Failed to delete %s: %v", path, err)
			errors++
			return nil
		}

		logger.Printf("  deleted: %s", path)
		deleted++
		return nil
	})

	if err != nil {
		logger.Fatalf("Failed to walk directory %s: %v", *dir, err)
	}

	logger.Printf("Done — deleted: %d, skipped: %d, errors: %d", deleted, skipped, errors)
	if errors > 0 {
		os.Exit(1)
	}
}
