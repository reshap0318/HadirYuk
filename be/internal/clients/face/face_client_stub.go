//go:build !gocv

package face

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"os"
)

// FaceClient handles face detection, embedding generation, and face matching.
// This is the stub implementation used when GoCV/OpenCV is not available.
// For production with real face recognition, build with: go build -tags gocv
type FaceClient struct {
	threshold float64
}

// NewFaceClient creates a new FaceClient instance with stub implementation.
func NewFaceClient() *FaceClient {
	return &FaceClient{
		threshold: 0.85,
	}
}

func (c *FaceClient) loadModels() error {
	return nil
}

func (c *FaceClient) findModelDir() string {
	return "."
}

func (c *FaceClient) SetThreshold(threshold float64) {
	c.threshold = threshold
}

func (c *FaceClient) GetThreshold() float64 {
	return c.threshold
}

func (c *FaceClient) IsReady() bool {
	return true
}

func (c *FaceClient) UseStub() bool {
	return true
}

func (c *FaceClient) InitError() error {
	return fmt.Errorf("GoCV not enabled. Build with -tags gocv for real face detection")
}

func (c *FaceClient) DetectFace(imagePath string) (bool, error) {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return false, fmt.Errorf("image file not found")
	}
	// Stub: assume face is detected if file exists
	return true, nil
}

func (c *FaceClient) GenerateEmbedding(imagePath string) ([]float64, error) {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("image file not found")
	}

	// Read file content to generate deterministic embedding
	// Using file content hash ensures same photo produces same embedding
	// regardless of storage location (tmp vs face-photos)
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file: %w", err)
	}

	// Generate deterministic 128-dim embedding from file content hash
	hash := sha256.Sum256(data)
	hexStr := hex.EncodeToString(hash[:])

	embedding := make([]float64, 128)
	for i := 0; i < 128; i++ {
		charIdx := i % len(hexStr)
		val := float64(hexStr[charIdx]) / 255.0 * 2.0 - 1.0
		embedding[i] = val
	}

	// Normalize
	norm := 0.0
	for _, v := range embedding {
		norm += v * v
	}
	norm = math.Sqrt(norm)
	if norm > 0 {
		for i := range embedding {
			embedding[i] /= norm
		}
	}

	return embedding, nil
}

func (c *FaceClient) MatchFace(embedding1, embedding2 []float64) (float64, bool) {
	if len(embedding1) != len(embedding2) {
		return 0, false
	}
	similarity := cosineSimilarity(embedding1, embedding2)
	return similarity, similarity >= c.threshold
}

func (c *FaceClient) ProcessFacePhoto(imagePath string) ([]float64, error) {
	faceDetected, err := c.DetectFace(imagePath)
	if err != nil {
		return nil, err
	}
	if !faceDetected {
		return nil, fmt.Errorf("no face detected in the image")
	}
	return c.GenerateEmbedding(imagePath)
}

func (c *FaceClient) Close() {
	// No resources to release in stub implementation
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	dotProduct, normA, normB := 0.0, 0.0, 0.0
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	normA, normB = math.Sqrt(normA), math.Sqrt(normB)
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (normA * normB)
}
