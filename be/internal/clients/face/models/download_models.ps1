# Download OpenCV models for face detection and embedding
# Run this script once to set up the required model files

$modelsDir = Join-Path $PSScriptRoot "models"
New-Item -ItemType Directory -Force -Path $modelsDir | Out-Null

Write-Host "Downloading Haar Cascade for face detection..."
$haarcascadeUrl = "https://raw.githubusercontent.com/opencv/opencv/master/data/haarcascades/haarcascade_frontalface_default.xml"
$haarcascadePath = Join-Path $modelsDir "haarcascade_frontalface_default.xml"
Invoke-WebRequest -Uri $haarcascadeUrl -OutFile $haarcascadePath
Write-Host "✓ Haar cascade downloaded"

Write-Host "Downloading OpenFace FaceNet model (Torch format)..."
# OpenFace model - produces 128-dim embeddings
$openfaceUrl = "https://storage.cmusatyalab.org/openface-models/nn4.small2.v1.t7"
$openfacePath = Join-Path $modelsDir "nn4.small2.v1.t7"
Invoke-WebRequest -Uri $openfaceUrl -OutFile $openfacePath
Write-Host "✓ OpenFace model downloaded"

Write-Host ""
Write-Host "All models downloaded successfully!"
Write-Host "Location: $modelsDir"
Write-Host ""
Write-Host "Files:"
Get-ChildItem $modelsDir | ForEach-Object { Write-Host "  - $($_.Name) ($([math]::Round($_.Length/1MB, 2)) MB)" }
