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

const attendanceBufferMinutes = 15

// GetTodayStatus returns the enriched attendance status for the current user today.
func (s *Services) GetTodayStatus(ctx context.Context) (*dtos.AttendanceTodayResponse, error) {
	s.Logger.LogStart("GetTodayStatus", "Fetching today's attendance status")

	userID := helpers.GetCallerID(ctx)
	if userID == 0 {
		s.Logger.LogEndWithError("GetTodayStatus", "Invalid token: caller ID not found")
		return nil, helpers.ErrInvalidToken
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	response := &dtos.AttendanceTodayResponse{
		Sessions:    []dtos.AttendanceSessionDTO{},
		TodaysShifts: []dtos.TodaysShiftDTO{},
	}

	// Check for active session across ALL dates (cross-day scenario)
	activeSession, err := s.repo.Attendance.FindActiveSessionByUserID(nil, userID, "Shift")
	if err != nil && err != gorm.ErrRecordNotFound {
		s.Logger.LogEndWithError("GetTodayStatus", "Failed to fetch active session: %v", err)
		return nil, err
	}

	// Get today's completed sessions
	todayAttendances, err := s.repo.Attendance.FindByUserAndDate(nil, userID, today, "Shift")
	if err != nil && err != gorm.ErrRecordNotFound {
		s.Logger.LogEndWithError("GetTodayStatus", "Failed to fetch today's attendances: %v", err)
		return nil, err
	}

	// Build sessions array from today's attendances
	for _, att := range todayAttendances {
		dto := dtos.ToAttendanceDTO(&att)
		response.Sessions = append(response.Sessions, dtos.AttendanceSessionDTO{
			AttendanceDTO: dto,
			ShiftName:     att.Shift.Name,
		})
	}

	// Get all shifts assigned today
	allShifts, err := s.repo.UserShiftAssignment.FindAllActiveForUserDate(nil, userID, today, "Shift")
	if err != nil {
		s.Logger.LogError("GetTodayStatus", "Failed to fetch today's shifts: %v", err)
	}

	// Build todays_shifts with session linkage
	for _, assignment := range allShifts {
		shiftDTO := dtos.TodaysShiftDTO{
			ID:        assignment.ShiftID,
			Name:      assignment.Shift.Name,
			StartTime: assignment.Shift.StartTime,
			EndTime:   assignment.Shift.EndTime,
			ColorCode: assignment.Shift.ColorCode,
			Status:    "not_started",
		}

		// Find linked session
		for _, session := range response.Sessions {
			if session.ShiftID == assignment.ShiftID {
				sessionCopy := session
				shiftDTO.Session = &sessionCopy
				if session.TimeOut != nil {
					shiftDTO.Status = "completed"
				} else {
					shiftDTO.Status = "active"
				}
				break
			}
		}

		response.TodaysShifts = append(response.TodaysShifts, shiftDTO)
	}

	// Determine current action
	if activeSession != nil {
		// There's an active session — user must check out first
		response.CurrentAction = dtos.CurrentActionDTO{
			Action: "checkout",
			Shift: &dtos.ShiftDTO{
				ID:        activeSession.ShiftID,
				Name:      activeSession.Shift.Name,
				StartTime: activeSession.Shift.StartTime,
				EndTime:   activeSession.Shift.EndTime,
				ColorCode: activeSession.Shift.ColorCode,
			},
		}

		// Check if cross-day (active session from previous date)
		activeSessionDate := time.Date(activeSession.Date.Year(), activeSession.Date.Month(), activeSession.Date.Day(), 0, 0, 0, 0, activeSession.Date.Location())
		if !activeSessionDate.Equal(today) {
			response.CurrentAction.CrossDaySession = &dtos.CrossDaySessionDTO{
				ID:        activeSession.ID,
				ShiftName: activeSession.Shift.Name,
				Date:      activeSession.Date.Format("2006-01-02"),
			}
		}

		s.Logger.LogEnd("GetTodayStatus", "User %d has active session, action: checkout", userID)
		return response, nil
	}

	// Check if any today session is still active (shouldn't happen, but be safe)
	for _, att := range todayAttendances {
		if att.TimeIn != nil && att.TimeOut == nil {
			response.CurrentAction = dtos.CurrentActionDTO{
				Action: "checkout",
				Shift: &dtos.ShiftDTO{
					ID:        att.ShiftID,
					Name:      att.Shift.Name,
					StartTime: att.Shift.StartTime,
					EndTime:   att.Shift.EndTime,
					ColorCode: att.Shift.ColorCode,
				},
			}
			s.Logger.LogEnd("GetTodayStatus", "User %d has active session today, action: checkout", userID)
			return response, nil
		}
	}

	// Check if there's an applicable shift for check-in
	applicableShift := s.findApplicableShift(userID, now)
	if applicableShift != nil {
		response.CurrentAction = dtos.CurrentActionDTO{
			Action: "checkin",
			Shift: &dtos.ShiftDTO{
				ID:        applicableShift.ShiftID,
				Name:      applicableShift.Shift.Name,
				StartTime: applicableShift.Shift.StartTime,
				EndTime:   applicableShift.Shift.EndTime,
				ColorCode: applicableShift.Shift.ColorCode,
			},
		}
	} else {
		// No applicable shift — check if all shifts for today are done
		allDone := true
		for _, ts := range response.TodaysShifts {
			if ts.Status != "completed" {
				allDone = false
				break
			}
		}
		if allDone && len(response.TodaysShifts) > 0 {
			response.CurrentAction = dtos.CurrentActionDTO{Action: "done"}
		} else {
			response.CurrentAction = dtos.CurrentActionDTO{Action: "done"}
		}
	}

	s.Logger.LogEnd("GetTodayStatus", "User %d has %d sessions, action: %s", userID, len(todayAttendances), response.CurrentAction.Action)
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

// AttendanceCheckIn handles the check-in process with auto-detect shift.
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

	// MS-05: Validate NO active session exists (check ALL dates, not just today)
	activeSession, err := s.repo.Attendance.FindActiveSessionByUserID(nil, userID, "Shift")
	if err != nil && err != gorm.ErrRecordNotFound {
		s.Logger.LogEndWithError("AttendanceCheckIn", "Failed to check active session: %v", err)
		return nil, err
	}
	if activeSession != nil {
		sessionDate := activeSession.Date.Format("02/01")
		s.Logger.LogEndWithError("AttendanceCheckIn", "User has active session from %s", sessionDate)
		return nil, &helpers.CustomError{
			Message: fmt.Sprintf("Anda memiliki sesi check-in yang belum di-checkout dari tanggal %s. Silakan check-out terlebih dahulu.", sessionDate),
		}
	}

	// MS-06: Auto-detect shift via FindApplicableShift
	applicableShift := s.findApplicableShift(userID, now)
	if applicableShift == nil {
		s.Logger.LogEndWithError("AttendanceCheckIn", "No applicable shift found for current time")
		return nil, &helpers.CustomError{Message: "Tidak ada shift yang tersedia untuk check-in saat ini"}
	}

	s.Logger.LogStep("AttendanceCheckIn", "Applicable shift detected: %s (%s-%s)", applicableShift.Shift.Name, applicableShift.Shift.StartTime, applicableShift.Shift.EndTime)

	// MS-08: Duplicate check per shift — check (user_id, date, shift_id)
	existing, err := s.repo.Attendance.FindByUserDateShift(nil, userID, today, applicableShift.ShiftID)
	if err != nil && err != gorm.ErrRecordNotFound {
		s.Logger.LogEndWithError("AttendanceCheckIn", "Failed to check existing attendance: %v", err)
		return nil, err
	}
	if existing != nil && existing.TimeIn != nil {
		s.Logger.LogEndWithError("AttendanceCheckIn", "User already checked in for this shift today")
		return nil, &helpers.CustomError{Message: fmt.Sprintf("Anda sudah melakukan check-in untuk shift %s hari ini", applicableShift.Shift.Name)}
	}

	// Find nearest active office location
	nearestOffice, minDistance, err := s.findNearestOffice(req.Lat, req.Lng)
	if err != nil {
		return nil, err
	}

	s.Logger.LogStep("AttendanceCheckIn", "Nearest office: %s, distance: %.2f meters, radius: %d meters", nearestOffice.Name, minDistance, nearestOffice.RadiusMeters)

	// Validate distance against office radius
	if minDistance > float64(nearestOffice.RadiusMeters) {
		s.Logger.LogEndWithError("AttendanceCheckIn", "User outside office area: %.2f > %d", minDistance, nearestOffice.RadiusMeters)
		return nil, &helpers.CustomError{Message: "Anda berada di luar area kantor"}
	}

	// MS-07: Determine status based on window logic
	shiftStart, _ := time.Parse("15:04", applicableShift.Shift.StartTime)
	shiftEnd, _ := time.Parse("15:04", applicableShift.Shift.EndTime)
	shiftStartTime := time.Date(now.Year(), now.Month(), now.Day(), shiftStart.Hour(), shiftStart.Minute(), 0, 0, now.Location())
	shiftEndTime := time.Date(now.Year(), now.Month(), now.Day(), shiftEnd.Hour(), shiftEnd.Minute(), 0, 0, now.Location())

	buffer := time.Duration(attendanceBufferMinutes) * time.Minute
	flexiStart := shiftStartTime.Add(-buffer)
	flexiThreshold := shiftStartTime.Add(buffer)

	checkInTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, now.Location())

	if checkInTime.Before(flexiStart) || checkInTime.After(shiftEndTime) {
		s.Logger.LogEndWithError("AttendanceCheckIn", "Check-in outside shift window: %s not in %s-%s", checkInTime.Format("15:04"), flexiStart.Format("15:04"), shiftEndTime.Format("15:04"))
		return nil, &helpers.CustomError{Message: "Waktu check-in di luar jendela shift yang diperbolehkan"}
	}

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
		ShiftID:        applicableShift.ShiftID,
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
			Message: fmt.Sprintf("Check-in berhasil di %s (%.0fm dari kantor) — Shift %s", nearestOffice.Name, minDistance, applicableShift.Shift.Name),
			Data: map[string]interface{}{
				"id":              result.ID,
				"user_id":         userID,
				"status":          status,
				"distance_meters": minDistance,
				"office":          nearestOffice.Name,
				"shift_id":        applicableShift.ShiftID,
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
	s.Logger.LogEnd("AttendanceCheckIn", "Check-in successful: user %d, shift %s, status %s, distance %.2fm", userID, applicableShift.Shift.Name, status, minDistance)
	return &dto, nil
}

// AttendanceCheckOut handles the check-out process with window validation and overtime calculation.
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

	// MS-09: Validate active session exists (TimeIn != nil AND TimeOut == nil)
	activeSession, err := s.repo.Attendance.FindActiveSessionByUserID(nil, userID, "Shift")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.Logger.LogEndWithError("AttendanceCheckOut", "No active session found")
			return nil, &helpers.CustomError{Message: "Belum melakukan check-in"}
		}
		s.Logger.LogEndWithError("AttendanceCheckOut", "Failed to fetch active session: %v", err)
		return nil, err
	}

	s.Logger.LogStep("AttendanceCheckOut", "Active session found: ID=%d, Shift=%s, Date=%v", activeSession.ID, activeSession.Shift.Name, activeSession.Date)

	// MS-10 & MS-11: Validate window and calculate overtime
	shiftEnd, _ := time.Parse("15:04", activeSession.Shift.EndTime)
	shiftEndTime := time.Date(activeSession.Date.Year(), activeSession.Date.Month(), activeSession.Date.Day(), shiftEnd.Hour(), shiftEnd.Minute(), 0, 0, activeSession.Date.Location())

	buffer := time.Duration(attendanceBufferMinutes) * time.Minute
	earliestCheckout := shiftEndTime.Add(-buffer)
	latestNormalCheckout := shiftEndTime.Add(buffer)

	if now.Before(earliestCheckout) {
		s.Logger.LogEndWithError("AttendanceCheckOut", "Check-out before window: %s < %s", now.Format("15:04"), earliestCheckout.Format("15:04"))
		return nil, &helpers.CustomError{Message: fmt.Sprintf("Belum waktunya check-out. Check-out dapat dilakukan setelah %s", earliestCheckout.Format("15:04"))}
	}

	// MS-11: Calculate overtime_minutes if check-out > shiftEnd + buffer
	overtimeMinutes := 0
	if now.After(latestNormalCheckout) {
		overtimeMinutes = int(now.Sub(latestNormalCheckout).Minutes())
		s.Logger.LogStep("AttendanceCheckOut", "Overtime detected: %d minutes (checkout %s > %s)", overtimeMinutes, now.Format("15:04"), latestNormalCheckout.Format("15:04"))
	}

	// Find nearest active office location
	nearestOffice, minDistance, err := s.findNearestOffice(req.Lat, req.Lng)
	if err != nil {
		return nil, err
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
	duration := calculateDuration(*activeSession.TimeIn, now)

	// Update attendance record
	activeSession.TimeOut = &now
	activeSession.LatOut = &req.Lat
	activeSession.LngOut = &req.Lng
	activeSession.ImageOut = fotoPath
	activeSession.Duration = duration
	activeSession.OvertimeMinutes = overtimeMinutes

	var result *models.Attendance
	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		var err error
		result, err = s.repo.Attendance.Update(tx, &models.Attendance{ID: activeSession.ID}, activeSession)
		if err != nil {
			return nil, err
		}

		msg := fmt.Sprintf("Check-out berhasil di %s. Durasi kerja: %s", nearestOffice.Name, duration)
		if overtimeMinutes > 0 {
			msg += fmt.Sprintf(" (Lembur: %d menit)", overtimeMinutes)
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "info",
			Title:   "Check-out Berhasil",
			Message: msg,
			Data: map[string]interface{}{
				"id":               result.ID,
				"user_id":          userID,
				"duration":         duration,
				"distance_meters":  minDistance,
				"office":           nearestOffice.Name,
				"overtime_minutes": overtimeMinutes,
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
	s.Logger.LogEnd("AttendanceCheckOut", "Check-out successful: user %d, duration %s, overtime %dm", userID, duration, overtimeMinutes)
	return &dto, nil
}

// findApplicableShift finds the shift whose window matches the current time.
// Returns the shift assignment whose [start-buffer, end] window contains `now`.
// If multiple shifts overlap, returns the one with the earliest end time.
func (s *Services) findApplicableShift(userID uint, now time.Time) *models.UserShiftAssignment {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	assignments, err := s.repo.UserShiftAssignment.FindAllActiveForUserDate(nil, userID, today, "Shift")
	if err != nil || len(assignments) == 0 {
		return nil
	}

	buffer := time.Duration(attendanceBufferMinutes) * time.Minute
	var bestMatch *models.UserShiftAssignment
	var bestEnd time.Time

	for i := range assignments {
		assignment := &assignments[i]
		shiftStart, parseErr := time.Parse("15:04", assignment.Shift.StartTime)
		shiftEnd, parseEndErr := time.Parse("15:04", assignment.Shift.EndTime)
		if parseErr != nil || parseEndErr != nil {
			continue
		}

		windowStart := time.Date(now.Year(), now.Month(), now.Day(), shiftStart.Hour(), shiftStart.Minute(), 0, 0, now.Location()).Add(-buffer)
		windowEnd := time.Date(now.Year(), now.Month(), now.Day(), shiftEnd.Hour(), shiftEnd.Minute(), 0, 0, now.Location())

		if now.After(windowStart) && !now.After(windowEnd) {
			// This shift's window matches — pick the one with earliest end time
			if bestMatch == nil || windowEnd.Before(bestEnd) {
				bestMatch = assignment
				bestEnd = windowEnd
			}
		}
	}

	return bestMatch
}

// findNearestOffice finds the nearest active office to the given coordinates.
func (s *Services) findNearestOffice(lat, lng float64) (*models.OfficeLocation, float64, error) {
	locations, err := s.repo.OfficeLocation.FindByFieldMap(nil, map[string]interface{}{
		"is_active": true,
	})
	if err != nil {
		s.Logger.LogError("findNearestOffice", "Failed to fetch office locations: %v", err)
		return nil, 0, err
	}

	if len(locations) == 0 {
		return nil, 0, &helpers.CustomError{Message: "Tidak ada lokasi kantor yang aktif"}
	}

	var nearestOffice *models.OfficeLocation
	minDistance := math.MaxFloat64

	for i := range locations {
		dist := haversineDistance(lat, lng, locations[i].Latitude, locations[i].Longitude)
		if dist < minDistance {
			minDistance = dist
			nearestOffice = &locations[i]
		}
	}

	if nearestOffice == nil {
		return nil, 0, &helpers.CustomError{Message: "Tidak dapat menemukan kantor terdekat"}
	}

	return nearestOffice, minDistance, nil
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
