package services

import (
	"github.com/reshap0318/hadirYuk/internal/clients/face"
	"github.com/reshap0318/hadirYuk/internal/helpers"
)

// FaceService handles face detection, embedding generation, and face matching.
// It delegates actual face processing to FaceClient and handles business logic.
type FaceService struct {
	faceClient *face.FaceClient
	logger     *helpers.Logger
}

// NewFaceService creates a new FaceService instance.
func NewFaceService(faceClient *face.FaceClient, logger *helpers.Logger) *FaceService {
	return &FaceService{
		faceClient: faceClient,
		logger:     logger,
	}
}

// SetThreshold sets the similarity threshold for face matching.
func (s *FaceService) SetThreshold(threshold float64) {
	s.faceClient.SetThreshold(threshold)
}

// DetectFace validates that a face is detected in the image file.
func (s *FaceService) DetectFace(imagePath string) (bool, error) {
	s.logger.LogStart("FaceService.DetectFace", "Checking face in image: %s", imagePath)

	faceDetected, err := s.faceClient.DetectFace(imagePath)
	if err != nil {
		s.logger.LogEndWithError("FaceService.DetectFace", "Face detection failed: %v", err)
		return false, err
	}

	s.logger.LogEnd("FaceService.DetectFace", "Face detected")
	return faceDetected, nil
}

// GenerateEmbedding generates a face embedding vector from an image.
func (s *FaceService) GenerateEmbedding(imagePath string) ([]float64, error) {
	s.logger.LogStart("FaceService.GenerateEmbedding", "Generating embedding for: %s", imagePath)

	embedding, err := s.faceClient.GenerateEmbedding(imagePath)
	if err != nil {
		s.logger.LogEndWithError("FaceService.GenerateEmbedding", "Embedding generation failed: %v", err)
		return nil, err
	}

	s.logger.LogEnd("FaceService.GenerateEmbedding", "Generated 128-dim embedding")
	return embedding, nil
}

// MatchFace compares two embeddings and returns similarity score and match result.
func (s *FaceService) MatchFace(embedding1, embedding2 []float64) (float64, bool) {
	s.logger.LogStart("FaceService.MatchFace", "Comparing two embeddings")

	similarity, matched := s.faceClient.MatchFace(embedding1, embedding2)

	s.logger.LogEnd("FaceService.MatchFace", "Similarity: %.4f, Matched: %v", similarity, matched)
	return similarity, matched
}

// ProcessFacePhoto handles the full pipeline: detect → embed → return embedding.
func (s *FaceService) ProcessFacePhoto(imagePath string) ([]float64, error) {
	s.logger.LogStart("FaceService.ProcessFacePhoto", "Processing face photo: %s", imagePath)

	embedding, err := s.faceClient.ProcessFacePhoto(imagePath)
	if err != nil {
		s.logger.LogEndWithError("FaceService.ProcessFacePhoto", "Face photo processing failed: %v", err)
		return nil, err
	}

	s.logger.LogEnd("FaceService.ProcessFacePhoto", "Face photo processed successfully")
	return embedding, nil
}
