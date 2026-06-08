# Backend Implementation Plan - Face Photo (T-024)

## Overview & Scope

Implementasi fitur **Face Photo** untuk sistem absensi HadirYuk. Fitur ini mencakup:
1. **Face Photo Upload** — endpoint untuk admin mengupload foto wajah employee (selain endpoint `/me` yang sudah ada)
2. **Face Embedding Generation** — service untuk generate face embedding dari foto yang diupload
3. **Face Recognition Service** — service untuk mencocokkan wajah saat check-in/check-out
4. **Face Verification** — validasi bahwa wajah terdeteksi dalam foto sebelum disimpan

> **Catatan:** Endpoint `/api/me/face-photo` (self-upload) dan `/api/me/face-photo` (delete) **sudah ada** di `profile_handler.go` dan `profile_service.go`. Plan ini fokus pada:
> - Admin upload face photo untuk user lain (`POST /api/users/:id/face-photo`)
> - Face embedding generation service
> - Face recognition/matching service
> - GoCV integration

## Tech Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| Face Detection & Recognition | GoCV (OpenCV bindings for Go) | Face detection, landmark extraction, embedding generation |
| Alternative | face-recognition-go | Higher-level face recognition API (wrapper around dlib) |
| File Storage | Local (`storage/face-photos/`) | Store face photo files |
| Database | MySQL (`user_profiles.face_embedding`) | Store face embedding as text (JSON array of float64) |

## Database Changes

**No new tables or migrations needed.** Model `UserProfile` sudah memiliki field yang diperlukan:

| Field | Type | Status |
|-------|------|--------|
| `face_photo_url` | `varchar(500)` | ✅ Sudah ada |
| `face_embedding` | `text` | ✅ Sudah ada |

## Dependencies

### Go Dependencies (to be added to go.mod)

```bash
# Option A: GoCV (recommended for production)
go get gocv.io/x/gocv

# Option B: face-recognition-go (simpler API)
go get github.com/Kagami/go-face

# Note: Both require OpenCV/dlib native libraries installed on the system
# For development without native libs, implement a stub/mock service first
```

### System Dependencies

| Dependency | Purpose | Install |
|------------|---------|---------|
| OpenCV 4.x | Face detection (Haar Cascade / DNN) | `apt install libopencv-dev` (Linux) / `brew install opencv` (macOS) |
| dlib (optional) | Face landmark detection (if using face-recognition-go) | `apt install libdlib-dev` |

## Files to Create (6 files)

### 1. Service — Face Recognition Service

| File | Purpose |
|------|---------|
| `internal/services/face_service.go` | Face detection, embedding generation, face matching |

**Methods:**
```go
type FaceService struct {
    logger     *helpers.Logger
    threshold  float64  // default 0.85
    modelPath  string   // path to face detection model
}

// DetectFace validates that a face is detected in the image file
func (s *FaceService) DetectFace(imagePath string) (bool, error)

// GenerateEmbedding generates a face embedding vector from an image
func (s *FaceService) GenerateEmbedding(imagePath string) ([]float64, error)

// MatchFace compares two embeddings and returns similarity score
func (s *FaceService) MatchFace(embedding1, embedding2 []float64) (float64, bool)

// ProcessFacePhoto handles the full pipeline: detect → embed → return
func (s *FaceService) ProcessFacePhoto(imagePath string) ([]float64, error)
```

### 2. Service — Admin Face Photo Upload

| File | Purpose |
|------|---------|
| `internal/services/admin_face_service.go` | Admin upload face photo for any user |

**Methods:**
```go
// AdminUploadFacePhoto uploads face photo for a specific user (admin action)
func (s *Services) AdminUploadFacePhoto(ctx context.Context, userID uint, fileUUID string) (*dtos.UserDTO, error)

// AdminDeleteFacePhoto removes face photo for a specific user (admin action)
func (s *Services) AdminDeleteFacePhoto(ctx context.Context, userID uint) (*dtos.UserDTO, error)
```

### 3. DTO — Face Photo Request/Response

| File | Purpose |
|------|---------|
| `internal/dtos/face_dto.go` | Face photo DTOs (jika belum ada di master_dto.go) |

**DTOs:**
```go
type FacePhotoRequest struct {
    FacePhoto string `json:"face_photo" binding:"required"` // UUID from /api/upload
}

type FacePhotoResponse struct {
    PhotoURL       string  `json:"photo_url"`
    EmbeddingReady bool    `json:"embedding_ready"`
    Similarity     float64 `json:"similarity,omitempty"` // if matching
}

type FaceMatchRequest struct {
    PhotoBase64 string `json:"photo" binding:"required"` // base64 encoded image
}

type FaceMatchResponse struct {
    Matched    bool    `json:"matched"`
    Similarity float64 `json:"similarity"`
    UserID     uint    `json:"user_id,omitempty"`
}
```

### 4. Handler — Admin Face Photo Handler

| File | Purpose |
|------|---------|
| `internal/handlers/admin_face_handler.go` | Handlers for admin face photo operations |

**Handlers:**
```go
// AdminUploadFacePhoto handles POST /api/users/:id/face-photo
func (h *Handlers) AdminUploadFacePhoto(c *gin.Context)

// AdminDeleteFacePhoto handles DELETE /api/users/:id/face-photo
func (h *Handlers) AdminDeleteFacePhoto(c *gin.Context)

// FaceMatch handles POST /api/face/match (for attendance verification)
func (h *Handlers) FaceMatch(c *gin.Context)
```

### 5. Routes — Face Photo Routes

| File | Purpose |
|------|---------|
| `internal/routes/face_route.go` | Route definitions for face photo endpoints |

### 6. Helper — Face Image Processing

| File | Purpose |
|------|---------|
| `internal/helpers/face_helper.go` | Utility functions for face image processing (base64 decode, temp file handling) |

**Functions:**
```go
// Base64ToTempFile decodes base64 image to a temporary file
func Base64ToTempFile(base64Data string) (string, error)

// ValidateFaceImageFormat checks if image is valid format and size
func ValidateFaceImageFormat(filePath string) error
```

## Files to Modify (5 files)

| File | Change |
|------|--------|
| `internal/services/00_services.go` | Add `FaceService` field to `Services` struct, initialize in `NewServices` |
| `internal/handlers/00_handlers.go` | Add `FaceService` dependency to `Handlers` struct |
| `internal/routes/user_route.go` | Register face photo endpoints under `/api/users/:id` |
| `cmd/api/main.go` | Register face route group |
| `cmd/migration/main.go` | No change needed (fields already exist) |

## API Endpoints

### Admin Face Photo Management

| Method | Endpoint | Permission | Request | Response |
|--------|----------|------------|---------|----------|
| POST | `/api/users/:id/face-photo` | `user.upload-face` | `{ "face_photo": "uuid" }` | `{ "code": 200, "data": { "photo_url": "...", "embedding_ready": true } }` |
| DELETE | `/api/users/:id/face-photo` | `user.upload-face` | - | `{ "code": 200, "data": { ... } }` |

### Face Recognition (for attendance)

| Method | Endpoint | Permission | Request | Response |
|--------|----------|------------|---------|----------|
| POST | `/api/face/match` | baseline | `{ "photo": "base64" }` | `{ "code": 200, "data": { "matched": true, "similarity": 0.92, "user_id": 1 } }` |

> **Note:** `/api/face/match` digunakan oleh attendance check-in flow (T-027) untuk memverifikasi wajah karyawan.

## Step-by-Step Implementation Tasks

### Phase 1: Foundation (No GoCV Dependency)

| Step | Task | Description | Effort |
|------|------|-------------|--------|
| 1 | Create `face_dto.go` | Define FacePhotoRequest, FacePhotoResponse, FaceMatchRequest, FaceMatchResponse DTOs | 1h |
| 2 | Create `face_helper.go` | Base64 decode, temp file handling, image format validation | 1h |
| 3 | Create `face_service.go` (stub) | Implement stub service that returns mock embeddings for development | 2h |
| 4 | Create `admin_face_service.go` | Admin upload/delete face photo with embedding generation call | 2h |
| 5 | Create `admin_face_handler.go` | Handler for admin face photo endpoints | 1h |
| 6 | Create `face_route.go` | Route definitions | 0.5h |
| 7 | Modify `00_services.go` | Add FaceService to Services struct | 0.5h |
| 8 | Modify `00_handlers.go` | Add FaceService to Handlers struct | 0.5h |
| 9 | Modify `user_route.go` | Add face photo routes to user routes | 0.5h |
| 10 | Modify `main.go` | Register face route group | 0.5h |

### Phase 2: GoCV Integration (Production)

| Step | Task | Description | Effort |
|------|------|-------------|--------|
| 11 | Install OpenCV | Install OpenCV 4.x on development/production server | 1h |
| 12 | Add GoCV dependency | `go get gocv.io/x/gocv` | 0.5h |
| 13 | Implement `DetectFace` | Use GoCV Haar Cascade or DNN face detector | 2h |
| 14 | Implement `GenerateEmbedding` | Extract face embedding using GoCV DNN (FaceNet/ResNet) | 3h |
| 15 | Implement `MatchFace` | Cosine similarity comparison between embeddings | 1h |
| 16 | Implement `FaceMatch` handler | Full pipeline: base64 → temp file → detect → embed → match against all users | 2h |
| 17 | Add face model files | Download pre-trained face detection/recognition models | 0.5h |
| 18 | Add env config | `FACE_MODEL_PATH`, `FACE_THRESHOLD`, `FACE_PHOTO_MAX_SIZE` | 0.5h |

### Phase 3: Testing & Refinement

| Step | Task | Description | Effort |
|------|------|-------------|--------|
| 19 | Unit tests — face service | Test DetectFace, GenerateEmbedding, MatchFace with sample images | 2h |
| 20 | Integration tests — admin upload | Test full admin upload flow via API | 1h |
| 21 | Integration tests — face match | Test face matching with known photos | 1h |
| 22 | Error handling | Handle: no face detected, multiple faces, poor quality, wrong format | 1h |

## Testing Approach

### Unit Tests

| Test | Description |
|------|-------------|
| `TestFaceService_DetectFace` | Test face detection with images containing 0, 1, and multiple faces |
| `TestFaceService_GenerateEmbedding` | Test embedding generation returns consistent vector for same face |
| `TestFaceService_MatchFace` | Test cosine similarity: same face > threshold, different face < threshold |
| `TestFaceService_ProcessFacePhoto` | Test full pipeline integration |

### Integration Tests

| Test | Description |
|------|-------------|
| `TestAdminUploadFacePhoto` | Upload photo → verify file saved → verify embedding stored in DB |
| `TestAdminDeleteFacePhoto` | Delete photo → verify file removed → verify embedding cleared |
| `TestFaceMatchEndpoint` | POST base64 image → verify match result with correct user |

### Test Data

- Sample face photos: 3-5 photos per test user (front-facing, different lighting)
- Negative test images: no face, multiple faces, non-image files

## Environment Variables

```env
# Face Recognition
FACE_MODEL_PATH=./models/face_model.dat
FACE_THRESHOLD=0.85
FACE_PHOTO_MAX_SIZE=2097152  # 2MB
```

## Error Handling

| Scenario | Error Message | HTTP Status |
|----------|---------------|-------------|
| No face detected in photo | "Wajah tidak terdeteksi. Pastikan wajah terlihat jelas" | 422 |
| Multiple faces detected | "Hanya satu wajah yang diperbolehkan dalam foto" | 422 |
| Face embedding generation failed | "Gagal memproses foto wajah. Coba lagi" | 500 |
| Face match below threshold | "Wajah tidak dikenali. Pastikan pencahayaan cukup" | 401 |
| File size exceeds limit | "Ukuran file maksimal 2MB" | 422 |
| Invalid image format | "Format file harus JPG/PNG/WebP" | 422 |
| User not found | "User tidak ditemukan" | 404 |

## Security Considerations

1. **File validation:** Strict image format and size checks before processing
2. **Face detection mandatory:** Reject photos without detected faces
3. **Embedding storage:** Store as encrypted text in database (future enhancement)
4. **Rate limiting:** Apply rate limiting on `/api/face/match` to prevent brute-force
5. **Access control:** Admin endpoints require `user.upload-face` permission
