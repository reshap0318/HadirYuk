//go:build gocv

package face

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"gocv.io/x/gocv"
)

// FaceClient handles face detection, embedding generation, and face matching.
// Uses GoCV (OpenCV bindings) for real face detection and embedding extraction.
type FaceClient struct {
	threshold    float64
	classifier   gocv.CascadeClassifier
	net          gocv.Net
	modelsLoaded bool
	initError    error
}

// NewFaceClient creates a new FaceClient instance and attempts to load OpenCV models.
func NewFaceClient() *FaceClient {
	client := &FaceClient{
		threshold: 0.85,
	}

	if err := client.loadModels(); err != nil {
		client.initError = err
	}

	return client
}

func (c *FaceClient) loadModels() error {
	modelDir := c.findModelDir()

	haarcascadePath := filepath.Join(modelDir, "haarcascade_frontalface_default.xml")
	c.classifier = gocv.NewCascadeClassifier()
	if !c.classifier.Load(haarcascadePath) {
		c.classifier.Close()
		return fmt.Errorf("failed to load Haar cascade from %s", haarcascadePath)
	}

	facenetPath := filepath.Join(modelDir, "nn4.small2.v1.t7")
	c.net = gocv.ReadNetFromTorch(facenetPath)
	if c.net.Empty() {
		c.classifier.Close()
		return fmt.Errorf("failed to load FaceNet model from %s", facenetPath)
	}

	c.net.SetPreferableBackend(gocv.NetBackendDefault)
	c.net.SetPreferableTarget(gocv.NetTargetCPU)
	c.modelsLoaded = true
	return nil
}

func (c *FaceClient) findModelDir() string {
	if _, err := os.Stat("haarcascade_frontalface_default.xml"); err == nil {
		return "."
	}

	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		if _, err := os.Stat(filepath.Join(exeDir, "haarcascade_frontalface_default.xml")); err == nil {
			return exeDir
		}
		parentDir := filepath.Dir(exeDir)
		if _, err := os.Stat(filepath.Join(parentDir, "haarcascade_frontalface_default.xml")); err == nil {
			return parentDir
		}
	}

	if _, err := os.Stat("internal/clients/face/haarcascade_frontalface_default.xml"); err == nil {
		return "internal/clients/face"
	}

	return "."
}

func (c *FaceClient) SetThreshold(threshold float64) {
	c.threshold = threshold
}

func (c *FaceClient) GetThreshold() float64 {
	return c.threshold
}

func (c *FaceClient) IsReady() bool {
	return c.modelsLoaded
}

func (c *FaceClient) UseStub() bool {
	return false
}

func (c *FaceClient) InitError() error {
	return c.initError
}

func (c *FaceClient) DetectFace(imagePath string) (bool, error) {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return false, fmt.Errorf("image file not found: %s", imagePath)
	}

	img := gocv.IMRead(imagePath, gocv.IMReadColor)
	if img.Empty() {
		return false, fmt.Errorf("failed to read image file: %s", imagePath)
	}
	defer img.Close()

	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	rects := c.classifier.DetectMultiScale(gray)
	return len(rects) > 0, nil
}

func (c *FaceClient) GenerateEmbedding(imagePath string) ([]float64, error) {
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("image file not found: %s", imagePath)
	}

	img := gocv.IMRead(imagePath, gocv.IMReadColor)
	if img.Empty() {
		return nil, fmt.Errorf("failed to read image file: %s", imagePath)
	}
	defer img.Close()

	gray := gocv.NewMat()
	defer gray.Close()
	gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)

	rects := c.classifier.DetectMultiScale(gray)
	if len(rects) == 0 {
		return nil, fmt.Errorf("no face detected in the image")
	}

	faceRect := rects[0]
	faceImg := img.Region(faceRect)
	defer faceImg.Close()

	resized := gocv.NewMat()
	defer resized.Close()
	gocv.Resize(faceImg, &resized, gocv.NewSize(96, 96), 0, 0, gocv.InterpolationLinear)

	faceBlob := gocv.NewMat()
	defer faceBlob.Close()
	resized.ConvertTo(&faceBlob, gocv.MatTypeCV32F)
	faceBlob.ConvertTo(&faceBlob, gocv.MatTypeCV32F, 1.0/255.0, 0)

	rgb := gocv.NewMat()
	defer rgb.Close()
	gocv.CvtColor(faceBlob, &rgb, gocv.ColorBGRToRGB)

	blob := gocv.BlobFromImage(rgb, 1.0, gocv.NewSize(96, 96), gocv.NewScalar(0, 0, 0, 0), false, false)
	defer blob.Close()

	c.net.SetInput(blob, "")
	output := c.net.Forward("")
	defer output.Close()

	embedding := make([]float64, 128)
	for i := 0; i < 128; i++ {
		embedding[i] = float64(output.GetFloatAt(0, i))
	}

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
	if !c.classifier.Empty() {
		c.classifier.Close()
	}
	if !c.net.Empty() {
		c.net.Close()
	}
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
