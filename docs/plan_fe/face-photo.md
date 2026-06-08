# Frontend Implementation Plan - Face Photo (T-024)

## Overview & Scope

Implementasi fitur **Face Photo** pada frontend Vue 3 untuk sistem absensi HadirYuk. Fitur ini mencakup:
1. **Face Photo Capture Component** — komponen kamera untuk capture foto wajah langsung dari browser
2. **Face Photo Upload Component** — komponen upload foto wajah (reusable)
3. **Admin Face Photo Management** — upload face photo saat create/edit employee (User Management)
4. **Face Photo Integration** — integrasi dengan profile page (sudah ada, perlu enhancement)

> **Catatan:** Profile page (`/profile`) sudah memiliki face photo modal dengan upload file. Plan ini fokus pada:
> - Camera capture component (reusable)
> - Face photo upload di User Management (admin)
> - Face photo preview dan management yang lebih baik

## Tech Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| Camera Capture | `navigator.mediaDevices.getUserMedia` | Access device camera for face capture |
| Canvas API | HTML5 Canvas | Capture frame from video stream |
| File Upload | Existing `uploadFile` helper | Upload captured photo via `/api/upload` |
| UI Components | Existing UiButton, UiModal, FormFile | Reuse existing UI primitives |

## Database Changes

**None.** Frontend only interacts with existing API endpoints.

## Dependencies

### npm Packages (no new packages needed)

> Semua yang dibutuhkan sudah tersedia. Camera access menggunakan native browser API (`navigator.mediaDevices.getUserMedia`), tidak perlu library tambahan.

### Browser Requirements

| Requirement | Support |
|-------------|---------|
| HTTPS (required for camera) | All modern browsers |
| `getUserMedia` API | Chrome 53+, Firefox 36+, Safari 11+, Edge 79+ |
| Canvas API | All modern browsers |

## Files to Create (4 files)

### 1. Composable — Camera Capture

| File | Purpose |
|------|---------|
| `src/composables/useCamera.ts` | Composable untuk mengakses kamera, capture foto, dan manage stream |

**Interface:**
```typescript
interface UseCameraReturn {
  videoElement: Ref<HTMLVideoElement | null>
  stream: Ref<MediaStream | null>
  isCameraActive: Ref<boolean>
  error: Ref<string | null>
  
  startCamera: (facingMode?: 'user' | 'environment') => Promise<void>
  stopCamera: () => void
  capturePhoto: () => File | null
  switchCamera: () => void
}
```

**Features:**
- Start/stop camera stream
- Capture frame to File object (JPEG)
- Switch between front/back camera (mobile)
- Error handling (permission denied, no camera)
- Auto-cleanup on component unmount

### 2. Component — Face Photo Capture

| File | Purpose |
|------|---------|
| `src/components/utils/FacePhotoCapture.vue` | Reusable component untuk capture foto wajah dari kamera |

**Props:**
```typescript
interface Props {
  maxFileSize?: number        // default 2MB
  acceptedFormats?: string    // default 'image/jpeg,image/png,image/webp'
  showGuidelines?: boolean    // show face outline overlay
}
```

**Emits:**
```typescript
interface Emits {
  (e: 'captured', file: File): void
  (e: 'error', message: string): void
}
```

**Features:**
- Live camera preview with face guidelines overlay
- Capture button
- Retake button
- File size and format validation
- Loading state during capture

### 3. Component — Face Photo Manager

| File | Purpose |
|------|---------|
| `src/components/utils/FacePhotoManager.vue` | Reusable component untuk manage face photo (preview, upload, capture, remove) |

**Props:**
```typescript
interface Props {
  currentPhoto?: string | null   // URL of current face photo
  userId?: number                // User ID (for admin upload)
  mode?: 'self' | 'admin'        // self = /me endpoint, admin = /users/:id endpoint
}
```

**Emits:**
```typescript
interface Emits {
  (e: 'uploaded', photoUrl: string): void
  (e: 'removed'): void
  (e: 'error', message: string): void
}
```

**Features:**
- Display current face photo or placeholder
- Upload from file (FormFile)
- Capture from camera (FacePhotoCapture)
- Remove current photo with confirmation
- Tab/switch between upload and capture modes
- Progress indicator during upload

### 4. Page — User Face Photo (Admin)

| File | Purpose |
|------|---------|
| `src/pages/users/FacePhotoModal.vue` | Modal untuk admin upload face photo saat create/edit employee |

**Integration:**
- Dipanggil dari `users/FormModal.vue` atau sebagai standalone modal
- Menggunakan `FacePhotoManager` component dengan `mode="admin"`
- Calls API: `POST /api/users/:id/face-photo`

## Files to Modify (4 files)

| File | Change |
|------|--------|
| `src/stores/profile.ts` | No changes needed — sudah ada `uploadFacePhoto` dan `removeFacePhoto` |
| `src/stores/user.ts` | Add `uploadFacePhoto(userId, file)` dan `removeFacePhoto(userId)` methods |
| `src/pages/users/FormModal.vue` | Add face photo section with FacePhotoManager or link to FacePhotoModal |
| `src/pages/profile/IndexView.vue` | Enhance existing face photo modal to use FacePhotoCapture component |

## API Integration

### Existing Endpoints (already implemented)

| Method | Endpoint | Used By |
|--------|----------|---------|
| POST | `/api/upload` | Upload file → get UUID |
| PUT | `/api/me/face-photo` | Self-upload face photo |
| DELETE | `/api/me/face-photo` | Self-delete face photo |

### New Endpoints (to be used after backend implementation)

| Method | Endpoint | Used By |
|--------|----------|---------|
| POST | `/api/users/:id/face-photo` | Admin upload for employee |
| DELETE | `/api/users/:id/face-photo` | Admin delete for employee |
| POST | `/api/face/match` | Attendance face verification (T-027) |

## Step-by-Step Implementation Tasks

### Phase 1: Camera Composable & Capture Component

| Step | Task | Description | Effort |
|------|------|-------------|--------|
| 1 | Create `useCamera.ts` | Composable: start/stop camera, capture frame, error handling | 2h |
| 2 | Create `FacePhotoCapture.vue` | Component: camera preview, guidelines overlay, capture button | 2h |
| 3 | Test camera access | Test on desktop (webcam) and mobile (front camera) | 1h |

### Phase 2: Face Photo Manager Component

| Step | Task | Description | Effort |
|------|------|-------------|--------|
| 4 | Create `FacePhotoManager.vue` | Component: combine upload + capture + preview + remove | 2h |
| 5 | Integrate with profile page | Replace existing face photo modal content with FacePhotoManager | 1h |
| 6 | Add camera tab to profile modal | Add tab/switch between "Upload File" and "Capture Camera" | 1h |

### Phase 3: Admin Face Photo Management

| Step | Task | Description | Effort |
|------|------|-------------|--------|
| 7 | Extend `user.ts` store | Add `uploadFacePhoto(userId, file)` and `removeFacePhoto(userId)` methods | 1h |
| 8 | Create `FacePhotoModal.vue` | Modal for admin to manage employee face photo | 1.5h |
| 9 | Integrate with User Form | Add "Manage Face Photo" button in user create/edit form | 1h |
| 10 | Add permission check | Only show face photo management for users with `user.upload-face` | 0.5h |

### Phase 4: Polish & Testing

| Step | Task | Description | Effort |
|------|------|-------------|--------|
| 11 | Face guidelines overlay | Add SVG face outline overlay on camera preview for better positioning | 1h |
| 12 | Image quality validation | Client-side check: minimum resolution, face detection hint | 1h |
| 13 | Mobile responsiveness | Ensure camera component works well on mobile devices | 1h |
| 14 | Error states | Handle: no camera permission, camera in use, capture failed | 0.5h |

## UI/UX Design

### Face Photo Capture Component Layout

```
┌─────────────────────────────────┐
│  📷 Capture Foto Wajah          │
├─────────────────────────────────┤
│  ┌───────────────────────────┐  │
│  │                           │  │
│  │    [Camera Preview]       │  │
│  │    ┌─────────────────┐    │  │
│  │    │   Face Outline  │    │  │
│  │    │    (guideline)  │    │  │
│  │    └─────────────────┘    │  │
│  │                           │  │
│  └───────────────────────────┘  │
│                                 │
│  [📸 Capture]  [🔄 Switch Cam] │
│                                 │
│  💡 Tips:                       │
│  • Pastikan wajah terlihat jelas│
│  • Pencahayaan yang cukup       │
│  • Hadapkan wajah ke kamera     │
└─────────────────────────────────┘
```

### Face Photo Manager Layout

```
┌─────────────────────────────────┐
│  Kelola Foto Wajah              │
├─────────────────────────────────┤
│  [📁 Upload File] [📷 Kamera]  │ ← Tab switch
├─────────────────────────────────┤
│                                 │
│  Current Photo:                 │
│  ┌─────────────────┐            │
│  │   [Photo/Placeholder]       │
│  └─────────────────┘            │
│                                 │
│  If Upload tab:                 │
│  [FormFile component]           │
│                                 │
│  If Camera tab:                 │
│  [FacePhotoCapture component]   │
│                                 │
│  [🗑️ Hapus Foto] (if exists)   │
├─────────────────────────────────┤
│  [Batal]        [Simpan Foto]  │
└─────────────────────────────────┘
```

## Testing Approach

### Manual Testing Checklist

| Test | Description | Browser |
|------|-------------|---------|
| Camera permission | Request and grant camera permission | Chrome, Firefox, Safari |
| Camera denied | Handle permission denied gracefully | All |
| Capture photo | Capture and verify photo quality | All |
| Upload photo | Upload file and verify | All |
| Switch camera | Switch front/back on mobile | Mobile Chrome/Safari |
| Remove photo | Delete and verify removal | All |
| Admin upload | Admin uploads for employee | All |
| File validation | Reject oversized/wrong format files | All |

### Component Testing (if unit test framework available)

| Component | Test |
|-----------|------|
| `useCamera` | startCamera resolves, stopCamera cleans up stream, capturePhoto returns File |
| `FacePhotoCapture` | Renders camera preview, capture button emits 'captured', error state displays message |
| `FacePhotoManager` | Tab switching works, upload mode shows file input, capture mode shows camera |

## Browser Compatibility

| Feature | Chrome | Firefox | Safari | Edge |
|---------|--------|---------|--------|------|
| `getUserMedia` | ✅ 53+ | ✅ 36+ | ✅ 11+ | ✅ 79+ |
| Canvas `toBlob` | ✅ | ✅ | ✅ | ✅ |
| `MediaDevices.enumerateDevices` | ✅ | ✅ | ✅ | ✅ |
| Fallback message | "Browser tidak mendukung akses kamera. Silakan upload foto." |

## Security Considerations

1. **HTTPS required:** `getUserMedia` only works on secure contexts (HTTPS or localhost)
2. **Permission prompt:** Browser handles camera permission — no custom permission UI needed
3. **Stream cleanup:** Always stop camera stream on component unmount to release camera
4. **File validation:** Client-side validation before upload (format, size)
5. **No face data storage:** Captured photo is sent to server immediately, not stored locally
