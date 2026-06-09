package services

import (
	"context"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
)

// GetTodayStatus returns the attendance status for the current user today.
func (s *Services) GetTodayStatus(ctx context.Context) (*dtos.AttendanceStatusResponse, error) {
	s.Logger.LogStart("GetTodayStatus", "Fetching today's attendance status")

	userID := helpers.GetCallerID(ctx)
	if userID == 0 {
		s.Logger.LogEndWithError("GetTodayStatus", "Invalid token: caller ID not found")
		return nil, helpers.ErrInvalidToken
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	existing, err := s.repo.Attendance.FindByUserAndDate(nil, userID, today)
	if err != nil && err != gorm.ErrRecordNotFound {
		s.Logger.LogEndWithError("GetTodayStatus", "Failed to fetch attendance: %v", err)
		return nil, err
	}

	response := &dtos.AttendanceStatusResponse{
		HasCheckedIn: false,
	}

	if existing != nil && existing.TimeIn != nil {
		response.HasCheckedIn = true
		response.TimeIn = existing.TimeIn
		response.ShiftID = existing.ShiftID
		response.Status = existing.Status
		response.Distance = existing.DistanceMeters

		if existing.TimeOut != nil {
			response.HasCheckedOut = true
			response.TimeOut = existing.TimeOut
			response.Duration = existing.Duration
			s.Logger.LogEnd("GetTodayStatus", "User %d checked in at %v, checked out at %v", userID, *existing.TimeIn, *existing.TimeOut)
		} else {
			s.Logger.LogEnd("GetTodayStatus", "User %d checked in at %v, not yet checked out", userID, *existing.TimeIn)
		}
	} else {
		s.Logger.LogEnd("GetTodayStatus", "User %d has not checked in today", userID)
	}

	return response, nil
}

// NearestOffice finds the nearest active office location to the given coordinates.
func (s *Services) NearestOffice(ctx context.Context, req dtos.NearestOfficeRequest) (*dtos.NearestOfficeResponse, error) {
	s.Logger.LogStart("NearestOffice", "Finding nearest office for coordinates: %.6f, %.6f", req.Lat, req.Lng)

	locations, err := s.repo.OfficeLocation.FindByFieldMap(nil, map[string]interface{}{
		"is_active": true,
	})
	if err != nil {
		s.Logger.LogEndWithError("NearestOffice", "Failed to fetch office locations: %v", err)
		return nil, err
	}

	if len(locations) == 0 {
		s.Logger.LogEndWithError("NearestOffice", "No active office locations found")
		return nil, &helpers.CustomError{Message: "Tidak ada lokasi kantor yang aktif"}
	}

	var nearestOffice *models.OfficeLocation
	minDistance := math.MaxFloat64

	for i := range locations {
		dist := haversineDistance(req.Lat, req.Lng, locations[i].Latitude, locations[i].Longitude)
		if dist < minDistance {
			minDistance = dist
			nearestOffice = &locations[i]
		}
	}

	if nearestOffice == nil {
		s.Logger.LogEndWithError("NearestOffice", "No nearest office found")
		return nil, &helpers.CustomError{Message: "Tidak dapat menemukan kantor terdekat"}
	}

	response := dtos.NearestOfficeResponse{
		ID:           nearestOffice.ID,
		Name:         nearestOffice.Name,
		Latitude:     nearestOffice.Latitude,
		Longitude:    nearestOffice.Longitude,
		RadiusMeters: nearestOffice.RadiusMeters,
		Distance:     math.Round(minDistance*100) / 100,
	}

	s.Logger.LogEnd("NearestOffice", "Nearest office: %s, distance: %.2f meters", nearestOffice.Name, minDistance)
	return &response, nil
}

// AttendanceCheckIn handles the check-in process with geotagging, Haversine validation, and photo evidence.
func (s *Services) AttendanceCheckIn(ctx context.Context, req dtos.AttendanceCheckInRequest) (*dtos.AttendanceDTO, error) {
	s.Logger.LogStart("AttendanceCheckIn", "Processing check-in for user")

	userID := helpers.GetCallerID(ctx)
	if userID == 0 {
		s.Logger.LogEndWithError("AttendanceCheckIn", "Invalid token: caller ID not found")
		return nil, helpers.ErrInvalidToken
	}

	// Validate photo UUID - check file exists in tmp
	s.Logger.LogStep("AttendanceCheckIn", "Validating photo evidence UUID: %s", req.Image)
	fileMeta, err := helpers.GetFileMetadata(req.Image, "storage/tmp")
	if err != nil {
		s.Logger.LogEndWithError("AttendanceCheckIn", "Photo validation failed: %v", err)
		return nil, &helpers.CustomError{Message: "Foto bukti tidak valid atau tidak ditemukan"}
	}

	// Validate file extension
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[fileMeta.Extension] {
		s.Logger.LogEndWithError("AttendanceCheckIn", "Photo extension not allowed: %s", fileMeta.Extension)
		return nil, &helpers.CustomError{Message: "Format foto tidak didukung. Gunakan JPG, PNG, atau WEBP"}
	}

	// Validate file size (max 5MB)
	maxSizeMB := 5.0
	if fileMeta.SizeMB > maxSizeMB {
		s.Logger.LogEndWithError("AttendanceCheckIn", "Photo size too large: %.2fMB", fileMeta.SizeMB)
		return nil, &helpers.CustomError{Message: fmt.Sprintf("Ukuran foto terlalu besar (maks %.0fMB)", maxSizeMB)}
	}

	s.Logger.LogStep("AttendanceCheckIn", "Photo validated: %s, size: %.2fMB, ext: %s", req.Image, fileMeta.SizeMB, fileMeta.Extension)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Check for duplicate check-in on same date
	existing, err := s.repo.Attendance.FindByUserAndDate(nil, userID, today)
	if err != nil && err != gorm.ErrRecordNotFound {
		s.Logger.LogEndWithError("AttendanceCheckIn", "Failed to check existing attendance: %v", err)
		return nil, err
	}
	if existing != nil && existing.TimeIn != nil {
		s.Logger.LogEndWithError("AttendanceCheckIn", "User already checked in today")
		return nil, &helpers.CustomError{Message: "Anda sudah melakukan check-in hari ini"}
	}
	if err == gorm.ErrRecordNotFound {
		s.Logger.LogStep("AttendanceCheckIn", "No existing attendance found for today - proceeding with check-in")
	}

	// Find nearest active office location
	locations, err := s.repo.OfficeLocation.FindByFieldMap(nil, map[string]interface{}{
		"is_active": true,
	})
	if err != nil {
		s.Logger.LogEndWithError("AttendanceCheckIn", "Failed to fetch office locations: %v", err)
		return nil, err
	}

	if len(locations) == 0 {
		s.Logger.LogEndWithError("AttendanceCheckIn", "No active office locations found")
		return nil, &helpers.CustomError{Message: "Tidak ada lokasi kantor yang aktif"}
	}

	// Find nearest office using Haversine distance
	var nearestOffice *models.OfficeLocation
	minDistance := math.MaxFloat64

	for i := range locations {
		dist := haversineDistance(req.Lat, req.Lng, locations[i].Latitude, locations[i].Longitude)
		if dist < minDistance {
			minDistance = dist
			nearestOffice = &locations[i]
		}
	}

	s.Logger.LogStep("AttendanceCheckIn", "Nearest office: %s, distance: %.2f meters, radius: %d meters", nearestOffice.Name, minDistance, nearestOffice.RadiusMeters)

	// Validate distance against office radius
	if minDistance > float64(nearestOffice.RadiusMeters) {
		s.Logger.LogEndWithError("AttendanceCheckIn", "User outside office area: %.2f > %d", minDistance, nearestOffice.RadiusMeters)
		return nil, &helpers.CustomError{Message: "Anda berada di luar area kantor"}
	}

	// Check for active shift assignment
	userShift, err := s.repo.UserShiftAssignment.FindByUserID(nil, userID, "Shift")
	if err != nil || userShift == nil {
		s.Logger.LogEndWithError("AttendanceCheckIn", "User has no active shift assignment")
		return nil, &helpers.CustomError{Message: "Tidak ada jadwal shift kerja saat ini"}
	}

	// Validate today is within shift assignment date range
	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	if todayDate.Before(userShift.StartDate) {
		s.Logger.LogEndWithError("AttendanceCheckIn", "Check-in before shift start date: %s < %s", today.Format("2006-01-02"), userShift.StartDate.Format("2006-01-02"))
		return nil, &helpers.CustomError{Message: "Tidak ada jadwal shift kerja saat ini"}
	}
	if userShift.EndDate != nil && todayDate.After(*userShift.EndDate) {
		s.Logger.LogEndWithError("AttendanceCheckIn", "Check-in after shift end date: %s > %s", today.Format("2006-01-02"), userShift.EndDate.Format("2006-01-02"))
		return nil, &helpers.CustomError{Message: "Tidak ada jadwal shift kerja saat ini"}
	}

	// Validate check-in time is within shift window with flexi time
	shiftStart, parseStartErr := time.Parse("15:04", userShift.Shift.StartTime)
	shiftEnd, parseEndErr := time.Parse("15:04", userShift.Shift.EndTime)
	if parseStartErr != nil || parseEndErr != nil {
		s.Logger.LogEndWithError("AttendanceCheckIn", "Invalid shift time format")
		return nil, &helpers.CustomError{Message: "Tidak ada jadwal shift kerja saat ini"}
	}

	checkInTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, now.Location())
	shiftStartTime := time.Date(now.Year(), now.Month(), now.Day(), shiftStart.Hour(), shiftStart.Minute(), 0, 0, now.Location())
	shiftEndTime := time.Date(now.Year(), now.Month(), now.Day(), shiftEnd.Hour(), shiftEnd.Minute(), 0, 0, now.Location())

	// Apply flexi time: window starts at (start_time - flexi_minutes)
	flexiMinutes := userShift.Shift.FlexiMinutes
	flexiStart := shiftStartTime.Add(-time.Duration(flexiMinutes) * time.Minute)
	flexiThreshold := shiftStartTime.Add(time.Duration(flexiMinutes) * time.Minute)

	if checkInTime.Before(flexiStart) || checkInTime.After(shiftEndTime) {
		s.Logger.LogEndWithError("AttendanceCheckIn", "Check-in outside shift window: %s not in %s-%s", checkInTime.Format("15:04"), flexiStart.Format("15:04"), shiftEndTime.Format("15:04"))
		return nil, &helpers.CustomError{Message: "Tidak ada jadwal shift kerja saat ini"}
	}

	// Determine status: present if check-in <= start_time + flexi, late if after
	status := "present"
	if checkInTime.After(flexiThreshold) {
		status = "late"
	}

	// Move photo from tmp to attendance-evidence
	fotoPath, err := helpers.MoveFile(req.Image, "storage/tmp", "storage/attendance-evidence")
	if err != nil {
		s.Logger.LogEndWithError("AttendanceCheckIn", "Failed to move photo evidence: %v", err)
		return nil, &helpers.CustomError{Message: "Gagal menyimpan foto bukti"}
	}

	// Create attendance record
	attendance := &models.Attendance{
		UserID:         userID,
		ShiftID:        userShift.ShiftID,
		Date:           today,
		TimeIn:         &now,
		LatIn:          &req.Lat,
		LngIn:          &req.Lng,
		OfficeID:       nearestOffice.ID,
		Status:         status,
		DistanceMeters: &minDistance,
		ImageIn:        fotoPath,
	}

	var result *models.Attendance
	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		var err error
		result, err = s.repo.Attendance.Create(tx, attendance)
		if err != nil {
			return nil, err
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "success",
			Title:   "Check-in Berhasil",
			Message: fmt.Sprintf("Check-in berhasil di %s (%.0fm dari kantor)", nearestOffice.Name, minDistance),
			Data: map[string]interface{}{
				"id":              result.ID,
				"user_id":         userID,
				"status":          status,
				"distance_meters": minDistance,
				"office":          nearestOffice.Name,
			},
		})

		return result, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("AttendanceCheckIn", "Failed to create attendance: %v", err)
		return nil, err
	}

	result = res.(*models.Attendance)
	dto := dtos.ToAttendanceDTO(result)
	s.Logger.LogEnd("AttendanceCheckIn", "Check-in successful: user %d, status %s, distance %.2fm", userID, status, minDistance)
	return &dto, nil
}

// AttendanceCheckOut handles the check-out process with geotagging, photo evidence, and duration calculation.
func (s *Services) AttendanceCheckOut(ctx context.Context, req dtos.AttendanceCheckInRequest) (*dtos.AttendanceDTO, error) {
	s.Logger.LogStart("AttendanceCheckOut", "Processing check-out for user")

	userID := helpers.GetCallerID(ctx)
	if userID == 0 {
		s.Logger.LogEndWithError("AttendanceCheckOut", "Invalid token: caller ID not found")
		return nil, helpers.ErrInvalidToken
	}

	// Validate photo UUID - check file exists in tmp
	s.Logger.LogStep("AttendanceCheckOut", "Validating photo evidence UUID: %s", req.Image)
	fileMeta, err := helpers.GetFileMetadata(req.Image, "storage/tmp")
	if err != nil {
		s.Logger.LogEndWithError("AttendanceCheckOut", "Photo validation failed: %v", err)
		return nil, &helpers.CustomError{Message: "Foto bukti tidak valid atau tidak ditemukan"}
	}

	// Validate file extension
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[fileMeta.Extension] {
		s.Logger.LogEndWithError("AttendanceCheckOut", "Photo extension not allowed: %s", fileMeta.Extension)
		return nil, &helpers.CustomError{Message: "Format foto tidak didukung. Gunakan JPG, PNG, atau WEBP"}
	}

	// Validate file size (max 5MB)
	maxSizeMB := 5.0
	if fileMeta.SizeMB > maxSizeMB {
		s.Logger.LogEndWithError("AttendanceCheckOut", "Photo size too large: %.2fMB", fileMeta.SizeMB)
		return nil, &helpers.CustomError{Message: fmt.Sprintf("Ukuran foto terlalu besar (maks %.0fMB)", maxSizeMB)}
	}

	s.Logger.LogStep("AttendanceCheckOut", "Photo validated: %s, size: %.2fMB, ext: %s", req.Image, fileMeta.SizeMB, fileMeta.Extension)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Find existing attendance record for today
	existing, err := s.repo.Attendance.FindByUserAndDate(nil, userID, today)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.Logger.LogEndWithError("AttendanceCheckOut", "No check-in record found for today")
			return nil, &helpers.CustomError{Message: "Belum melakukan check-in hari ini"}
		}
		s.Logger.LogEndWithError("AttendanceCheckOut", "Failed to fetch attendance: %v", err)
		return nil, err
	}

	// Check if already checked out
	if existing.TimeOut != nil {
		s.Logger.LogEndWithError("AttendanceCheckOut", "User already checked out today")
		return nil, &helpers.CustomError{Message: "Anda sudah melakukan check-out hari ini"}
	}

	// Find nearest active office location
	locations, err := s.repo.OfficeLocation.FindByFieldMap(nil, map[string]interface{}{
		"is_active": true,
	})
	if err != nil {
		s.Logger.LogEndWithError("AttendanceCheckOut", "Failed to fetch office locations: %v", err)
		return nil, err
	}

	if len(locations) == 0 {
		s.Logger.LogEndWithError("AttendanceCheckOut", "No active office locations found")
		return nil, &helpers.CustomError{Message: "Tidak ada lokasi kantor yang aktif"}
	}

	// Find nearest office using Haversine distance
	var nearestOffice *models.OfficeLocation
	minDistance := math.MaxFloat64

	for i := range locations {
		dist := haversineDistance(req.Lat, req.Lng, locations[i].Latitude, locations[i].Longitude)
		if dist < minDistance {
			minDistance = dist
			nearestOffice = &locations[i]
		}
	}

	s.Logger.LogStep("AttendanceCheckOut", "Nearest office: %s, distance: %.2f meters", nearestOffice.Name, minDistance)

	// Validate distance against office radius
	if minDistance > float64(nearestOffice.RadiusMeters) {
		s.Logger.LogEndWithError("AttendanceCheckOut", "User outside office area: %.2f > %d", minDistance, nearestOffice.RadiusMeters)
		return nil, &helpers.CustomError{Message: "Anda berada di luar area kantor"}
	}

	// Move photo from tmp to attendance-evidence
	fotoPath, err := helpers.MoveFile(req.Image, "storage/tmp", "storage/attendance-evidence")
	if err != nil {
		s.Logger.LogEndWithError("AttendanceCheckOut", "Failed to move photo evidence: %v", err)
		return nil, &helpers.CustomError{Message: "Gagal menyimpan foto bukti"}
	}

	// Calculate duration
	duration := calculateDuration(*existing.TimeIn, now)

	// Update attendance record
	existing.TimeOut = &now
	existing.LatOut = &req.Lat
	existing.LngOut = &req.Lng
	existing.ImageOut = fotoPath
	existing.Duration = duration

	var result *models.Attendance
	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		var err error
		result, err = s.repo.Attendance.Update(tx, &models.Attendance{ID: existing.ID}, existing)
		if err != nil {
			return nil, err
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "info",
			Title:   "Check-out Berhasil",
			Message: fmt.Sprintf("Check-out berhasil di %s. Durasi kerja: %s", nearestOffice.Name, duration),
			Data: map[string]interface{}{
				"id":              result.ID,
				"user_id":         userID,
				"duration":        duration,
				"distance_meters": minDistance,
				"office":          nearestOffice.Name,
			},
		})

		return result, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("AttendanceCheckOut", "Failed to update attendance: %v", err)
		return nil, err
	}

	result = res.(*models.Attendance)
	dto := dtos.ToAttendanceDTO(result)
	s.Logger.LogEnd("AttendanceCheckOut", "Check-out successful: user %d, duration %s", userID, duration)
	return &dto, nil
}

// calculateDuration calculates the duration between two times and returns a human-readable string.
func calculateDuration(start, end time.Time) string {
	diff := end.Sub(start)
	hours := int(diff.Hours())
	minutes := int(diff.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

// haversineDistance calculates the distance between two points on Earth using the Haversine formula.
// Returns distance in meters.
func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371000 // meters

	dLat := degreesToRadians(lat2 - lat1)
	dLng := degreesToRadians(lng2 - lng1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// degreesToRadians converts degrees to radians.
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
