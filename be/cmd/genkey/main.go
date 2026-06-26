package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/reshap0318/hadirYuk/internal/helpers"
)

func main() {
	force := flag.Bool("f", false, "Force overwrite existing keys without confirmation")
	flag.Parse()

	privateKeyPath := helpers.GetEnv("JWT_PRIVATE_KEY_PATH", "storage/keys/private.pem")
	publicKeyPath := helpers.GetEnv("JWT_PUBLIC_KEY_PATH", "storage/keys/public.pem")
	passphraseFile := helpers.GetEnv("JWT_PASSPHRASE_PATH", "storage/keys/passphrase")

	if _, err := os.Stat(privateKeyPath); err == nil && !*force {
		fmt.Println("⚠️  WARNING: Existing keys found!")
		fmt.Println("⚠️  This will overwrite existing keys!")
		fmt.Println("⚠️  All existing JWT tokens will be invalidated!")
		fmt.Print("❓ Continue? (y/n): ")

		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Println("❌ Operation cancelled")
			return
		}
		fmt.Println()
	}

	privateKey, publicKey, err := generateRSAKeys()
	if err != nil {
		log.Fatalf("Error generating keys: %v", err)
	}

	passphrase, err := helpers.GenerateRandomString(32)
	if err != nil {
		log.Fatalf("Error generating passphrase: %v", err)
	}

	if err := os.MkdirAll("storage/keys", 0700); err != nil {
		log.Fatalf("Error creating storage/keys directory: %v", err)
	}

	if err := saveEncryptedPrivateKey(privateKeyPath, privateKey, passphrase); err != nil {
		log.Fatalf("Error saving private key: %v", err)
	}

	if err := savePublicKey(publicKeyPath, publicKey); err != nil {
		log.Fatalf("Error saving public key: %v", err)
	}

	if err := os.WriteFile(passphraseFile, []byte(passphrase), 0600); err != nil {
		log.Fatalf("Error saving passphrase: %v", err)
	}

	fmt.Println("✅ Keys generated successfully!")
	fmt.Printf("📁 Private key : %s\n", privateKeyPath)
	fmt.Printf("📁 Public key  : %s\n", publicKeyPath)
	fmt.Printf("📁 Passphrase  : %s\n", passphraseFile)
	fmt.Println("⚠️  If you regenerate keys, all existing JWT tokens will be invalid!")
}

func generateRSAKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}
	return privateKey, &privateKey.PublicKey, nil
}

func saveEncryptedPrivateKey(path string, key *rsa.PrivateKey, passphrase string) error {
	der := x509.MarshalPKCS1PrivateKey(key)
	block, err := x509.EncryptPEMBlock(rand.Reader, "RSA PRIVATE KEY", der, []byte(passphrase), x509.PEMCipherAES256)
	if err != nil {
		return fmt.Errorf("failed to encrypt private key: %w", err)
	}
	if err := os.WriteFile(path, pem.EncodeToMemory(block), 0600); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}
	return nil
}

func savePublicKey(path string, key *rsa.PublicKey) error {
	pub, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}
	block := &pem.Block{Type: "PUBLIC KEY", Bytes: pub}
	if err := os.WriteFile(path, pem.EncodeToMemory(block), 0644); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}
	return nil
}
