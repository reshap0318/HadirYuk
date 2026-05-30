package dtos

// FacePhotoResponse represents the response after face photo upload.
type FacePhotoResponse struct {
	PhotoURL       string `json:"photo_url"`
	EmbeddingReady bool   `json:"embedding_ready"`
}

// FaceMatchRequest represents the request to match a face photo.
type FaceMatchRequest struct {
	PhotoBase64 string `json:"photo" validate:"required"` // base64 encoded image
}

// FaceMatchResponse represents the response after face matching.
type FaceMatchResponse struct {
	Matched    bool    `json:"matched"`
	Similarity float64 `json:"similarity"`
	UserID     uint    `json:"user_id,omitempty"`
}
