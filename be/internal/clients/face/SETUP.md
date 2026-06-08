# GoCV Setup Guide for Windows

GoCV requires OpenCV 4.x to be installed on your system. Follow these steps to set it up:

## Option 1: Using MSYS2 (Recommended)

1. **Install MSYS2** from https://www.msys2.org/

2. **Open MSYS2 MinGW 64-bit terminal** and run:
   ```bash
   pacman -Syu
   pacman -S mingw-w64-x86_64-opencv
   ```

3. **Add MSYS2 to PATH**:
   Add `C:\msys64\mingw64\bin` to your system PATH environment variable.

4. **Set CGO environment variables**:
   ```powershell
   $env:CGO_CXXFLAGS="-std=c++11"
   $env:CGO_CPPFLAGS="-IC:/msys64/mingw64/include"
   $env:CGO_LDFLAGS="-LC:/msys64/mingw64/lib -lopencv_core4100 -lopencv_imgproc4100 -lopencv_imgcodecs4100 -lopencv_objdetect4100 -lopencv_dnn4100"
   ```
   Note: Replace `4100` with your OpenCV version (e.g., `4090` for 4.9.0).

5. **Build the project**:
   ```powershell
   cd be
   go build ./...
   ```

## Option 2: Using Pre-built OpenCV Binaries

1. **Download OpenCV** from https://opencv.org/releases/
   - Download the Windows pack (e.g., `opencv-4.9.0-windows.exe`)
   - Extract to `C:\opencv`

2. **Set environment variables**:
   ```powershell
   $env:CGO_CXXFLAGS="-std=c++11"
   $env:CGO_CPPFLAGS="-IC:/opencv/build/include"
   $env:CGO_LDFLAGS="-LC:/opencv/build/x64/mingw/lib -lopencv_core490 -lopencv_imgproc490 -lopencv_imgcodecs490 -lopencv_objdetect490 -lopencv_dnn490"
   $env:PATH="C:\opencv\build\x64\mingw\bin;$env:PATH"
   ```

3. **Build the project**:
   ```powershell
   cd be
   go build ./...
   ```

## Download Model Files

Run the PowerShell script to download required model files:
```powershell
cd be\internal\clients\face\models
.\download_models.ps1
```

This downloads:
- `haarcascade_frontalface_default.xml` - Face detection model
- `nn4.small2.v1.t7` - OpenFace FaceNet embedding model

## Verify Installation

Run the following to verify GoCV is working:
```powershell
cd be
go run -tags customenv gocv.io/x/gocv/cmd/version
```

## Troubleshooting

### "undefined: CalibFlag" errors
- OpenCV is not installed or not in PATH
- Verify OpenCV version matches CGO_LDFLAGS

### "cannot find -lopencv_core" errors
- OpenCV libraries not found
- Check CGO_LDFLAGS paths are correct

### "dll not found" at runtime
- Add OpenCV bin directory to PATH
- Copy OpenCV DLLs to executable directory

## Alternative: Development Mode (Without GoCV)

If you want to develop without installing OpenCV, the code includes a fallback stub implementation. Set the environment variable:
```powershell
$env:FACE_CLIENT_STUB="true"
```

This will use the mock embedding generator for development purposes.
