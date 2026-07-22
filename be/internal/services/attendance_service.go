package services

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"gorm.io/gorm"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/models"
	"github.com/reshap0318/hadirYuk/internal/repositories"
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
	yesterday := today.AddDate(0, 0, -1)

	response := &dtos.AttendanceTodayResponse{
		Sessions:     []dtos.AttendanceSessionDTO{},
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

	// Also include cross-midnight sessions from yesterday (shift end < shift start)
	yesterdayAttendances, _ := s.repo.Attendance.FindByUserAndDate(nil, userID, yesterday, "Shift")
	for _, att := range yesterdayAttendances {
		endT, errE := time.Parse("15:04", att.Shift.EndTime)
		startT, errS := time.Parse("15:04", att.Shift.StartTime)
		if errE == nil && errS == nil && endT.Before(startT) {
			dto := dtos.ToAttendanceDTO(&att)
			response.Sessions = append(response.Sessions, dtos.AttendanceSessionDTO{
				AttendanceDTO: dto,
				ShiftName:     att.Shift.Name,
			})
		}
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
	match := s.findApplicableShift(userID, now)
	if match != nil {
		response.CurrentAction = dtos.CurrentActionDTO{
			Action: "checkin",
			Shift: &dtos.ShiftDTO{
				ID:        match.Assignment.ShiftID,
				Name:      match.Assignment.Shift.Name,
				StartTime: match.Assignment.Shift.StartTime,
				EndTime:   match.Assignment.Shift.EndTime,
				ColorCode: match.Assignment.Shift.ColorCode,
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
		if allDone || len(response.TodaysShifts) == 0 {
			response.CurrentAction = dtos.CurrentActionDTO{Action: "done"}
		} else {
			// Shifts exist but none are applicable right now (too early or between windows)
			response.CurrentAction = dtos.CurrentActionDTO{Action: "waiting"}
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

	// MS-06: Auto-detect shift via findApplicableShift (supports cross-midnight shifts)
	shiftResult := s.findApplicableShift(userID, now)
	if shiftResult == nil {
		s.Logger.LogEndWithError("AttendanceCheckIn", "No applicable shift found for current time")
		return nil, &helpers.CustomError{Message: "Tidak ada shift yang tersedia untuk check-in saat ini"}
	}
	applicableShift := shiftResult.Assignment
	shiftDate := shiftResult.ShiftDate

	s.Logger.LogStep("AttendanceCheckIn", "Applicable shift detected: %s (%s-%s), shiftDate: %s", applicableShift.Shift.Name, applicableShift.Shift.StartTime, applicableShift.Shift.EndTime, shiftDate.Format("2006-01-02"))

	// MS-08: Duplicate check using shiftDate (cross-midnight shifts belong to the day they started)
	existing, err := s.repo.Attendance.FindByUserDateShift(nil, userID, shiftDate, applicableShift.ShiftID)
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
	// shiftStartTime uses shiftDate so cross-midnight shifts (e.g., 22:00 on Day1) are handled correctly.
	shiftStart, _ := time.Parse("15:04", applicableShift.Shift.StartTime)
	shiftStartTime := time.Date(shiftDate.Year(), shiftDate.Month(), shiftDate.Day(), shiftStart.Hour(), shiftStart.Minute(), 0, 0, now.Location())
	shiftEndTime := shiftResult.WindowEnd

	buffer := time.Duration(attendanceBufferMinutes) * time.Minute
	flexiStart := shiftStartTime.Add(-buffer)
	flexiThreshold := shiftStartTime.Add(buffer)

	if now.Before(flexiStart) || now.After(shiftEndTime) {
		s.Logger.LogEndWithError("AttendanceCheckIn", "Check-in outside shift window: %s not in %s-%s", now.Format("15:04"), flexiStart.Format("15:04"), shiftEndTime.Format("15:04"))
		return nil, &helpers.CustomError{Message: "Waktu check-in di luar jendela shift yang diperbolehkan"}
	}

	status := "present"
	if now.After(flexiThreshold) {
		status = "late"
	}

	// Move photo from tmp to attendance-evidence
	fotoPath, err := helpers.MoveFile(req.Image, "storage/tmp", "storage/attendance-evidence")
	if err != nil {
		s.Logger.LogEndWithError("AttendanceCheckIn", "Failed to move photo evidence: %v", err)
		return nil, &helpers.CustomError{Message: "Gagal menyimpan foto bukti"}
	}

	// Create attendance record (Date = shiftDate, not today, for cross-midnight shifts)
	attendance := &models.Attendance{
		UserID:         userID,
		ShiftID:        applicableShift.ShiftID,
		Date:           shiftDate,
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
		if _, rollbackErr := helpers.MoveFile(req.Image, "storage/attendance-evidence", "storage/tmp"); rollbackErr != nil {
			s.Logger.LogEndWithError("AttendanceCheckIn", "Failed to rollback photo move: %v", rollbackErr)
		}
		return nil, err
	}

	result = res.(*models.Attendance)
	dto := dtos.ToAttendanceDTO(result)
	s.Logger.LogEnd("AttendanceCheckIn", "Check-in successful: user %d, shift %s, status %s, distance %.2fm", userID, applicableShift.Shift.Name, status, minDistance)
	return &dto, nil
}

// AttendanceQRCheckIn handles check-in via QR Code (no photo, no GPS).
func (s *Services) AttendanceQRCheckIn(ctx context.Context, req dtos.AttendanceQRCheckInRequest) (*dtos.AttendanceDTO, error) {
	s.Logger.LogStart("AttendanceQRCheckIn", "Processing QR check-in")

	userID := helpers.GetCallerID(ctx)
	if userID == 0 {
		s.Logger.LogEndWithError("AttendanceQRCheckIn", "Invalid token: caller ID not found")
		return nil, helpers.ErrInvalidToken
	}

	// Validate QR Code
	qrCode, err := s.QRCodeValidate(ctx, req.CodeValue)
	if err != nil {
		s.Logger.LogEndWithError("AttendanceQRCheckIn", "QR validation failed: %v", err)
		return nil, &helpers.CustomError{Message: "QR Code tidak valid"}
	}

	s.Logger.LogStep("AttendanceQRCheckIn", "QR code validated: office_id=%d", qrCode.OfficeID)

	now := time.Now()

	// Validate NO active session exists (check ALL dates)
	activeSession, err := s.repo.Attendance.FindActiveSessionByUserID(nil, userID, "Shift")
	if err != nil && err != gorm.ErrRecordNotFound {
		s.Logger.LogEndWithError("AttendanceQRCheckIn", "Failed to check active session: %v", err)
		return nil, err
	}
	if activeSession != nil {
		sessionDate := activeSession.Date.Format("02/01")
		s.Logger.LogEndWithError("AttendanceQRCheckIn", "User has active session from %s", sessionDate)
		return nil, &helpers.CustomError{
			Message: fmt.Sprintf("Anda memiliki sesi check-in yang belum di-checkout dari tanggal %s. Silakan check-out terlebih dahulu.", sessionDate),
		}
	}

	// Auto-detect shift (supports cross-midnight shifts)
	shiftResult := s.findApplicableShift(userID, now)
	if shiftResult == nil {
		s.Logger.LogEndWithError("AttendanceQRCheckIn", "No applicable shift found")
		return nil, &helpers.CustomError{Message: "Tidak ada shift yang tersedia untuk check-in saat ini"}
	}
	applicableShift := shiftResult.Assignment
	shiftDate := shiftResult.ShiftDate

	s.Logger.LogStep("AttendanceQRCheckIn", "Applicable shift: %s, shiftDate: %s", applicableShift.Shift.Name, shiftDate.Format("2006-01-02"))

	// Duplicate check using shiftDate
	existing, err := s.repo.Attendance.FindByUserDateShift(nil, userID, shiftDate, applicableShift.ShiftID)
	if err != nil && err != gorm.ErrRecordNotFound {
		s.Logger.LogEndWithError("AttendanceQRCheckIn", "Failed to check existing attendance: %v", err)
		return nil, err
	}
	if existing != nil && existing.TimeIn != nil {
		s.Logger.LogEndWithError("AttendanceQRCheckIn", "User already checked in for this shift")
		return nil, &helpers.CustomError{Message: fmt.Sprintf("Anda sudah melakukan check-in untuk shift %s hari ini", applicableShift.Shift.Name)}
	}

	// Window validation
	shiftStart, _ := time.Parse("15:04", applicableShift.Shift.StartTime)
	shiftStartTime := time.Date(shiftDate.Year(), shiftDate.Month(), shiftDate.Day(), shiftStart.Hour(), shiftStart.Minute(), 0, 0, now.Location())
	shiftEndTime := shiftResult.WindowEnd

	buffer := time.Duration(attendanceBufferMinutes) * time.Minute
	flexiStart := shiftStartTime.Add(-buffer)
	flexiThreshold := shiftStartTime.Add(buffer)

	if now.Before(flexiStart) || now.After(shiftEndTime) {
		s.Logger.LogEndWithError("AttendanceQRCheckIn", "Check-in outside window")
		return nil, &helpers.CustomError{Message: "Waktu check-in di luar jendela shift yang diperbolehkan"}
	}

	status := "present"
	if now.After(flexiThreshold) {
		status = "late"
	}

	// Create attendance record (Date = shiftDate for cross-midnight support)
	attendance := &models.Attendance{
		UserID:   userID,
		ShiftID:  applicableShift.ShiftID,
		Date:     shiftDate,
		TimeIn:   &now,
		OfficeID: qrCode.OfficeID,
		Status:   status,
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
			Title:   "Check-in Berhasil (QR Code)",
			Message: fmt.Sprintf("Check-in via QR Code berhasil — Shift %s", applicableShift.Shift.Name),
			Data: map[string]interface{}{
				"id":       result.ID,
				"user_id":  userID,
				"status":   status,
				"office":   qrCode.OfficeID,
				"shift_id": applicableShift.ShiftID,
			},
		})

		return result, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("AttendanceQRCheckIn", "Failed to create attendance: %v", err)
		return nil, err
	}

	result = res.(*models.Attendance)
	dto := dtos.ToAttendanceDTO(result)
	s.Logger.LogEnd("AttendanceQRCheckIn", "QR check-in successful: user %d, shift %s, status %s", userID, applicableShift.Shift.Name, status)
	return &dto, nil
}

// AttendanceQRCheckOut handles check-out via QR Code (no photo, no GPS).
func (s *Services) AttendanceQRCheckOut(ctx context.Context, req dtos.AttendanceQRCheckInRequest) (*dtos.AttendanceDTO, error) {
	s.Logger.LogStart("AttendanceQRCheckOut", "Processing QR check-out")

	userID := helpers.GetCallerID(ctx)
	if userID == 0 {
		s.Logger.LogEndWithError("AttendanceQRCheckOut", "Invalid token: caller ID not found")
		return nil, helpers.ErrInvalidToken
	}

	// Validate QR Code
	qrCode, err := s.QRCodeValidate(ctx, req.CodeValue)
	if err != nil {
		s.Logger.LogEndWithError("AttendanceQRCheckOut", "QR validation failed: %v", err)
		return nil, &helpers.CustomError{Message: "QR Code tidak valid"}
	}

	s.Logger.LogStep("AttendanceQRCheckOut", "QR code validated: office_id=%d", qrCode.OfficeID)

	now := time.Now()

	// Validate active session exists
	activeSession, err := s.repo.Attendance.FindActiveSessionByUserID(nil, userID, "Shift")
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			s.Logger.LogEndWithError("AttendanceQRCheckOut", "No active session found")
			return nil, &helpers.CustomError{Message: "Belum melakukan check-in"}
		}
		s.Logger.LogEndWithError("AttendanceQRCheckOut", "Failed to fetch active session: %v", err)
		return nil, err
	}

	s.Logger.LogStep("AttendanceQRCheckOut", "Active session found: ID=%d, Shift=%s", activeSession.ID, activeSession.Shift.Name)

	// Window validation
	shiftEnd, _ := time.Parse("15:04", activeSession.Shift.EndTime)
	shiftStart, _ := time.Parse("15:04", activeSession.Shift.StartTime)
	shiftEndDay := activeSession.Date
	if shiftEnd.Hour()*60+shiftEnd.Minute() < shiftStart.Hour()*60+shiftStart.Minute() {
		shiftEndDay = activeSession.Date.AddDate(0, 0, 1)
	}
	shiftEndTime := time.Date(shiftEndDay.Year(), shiftEndDay.Month(), shiftEndDay.Day(), shiftEnd.Hour(), shiftEnd.Minute(), 0, 0, now.Location())

	buffer := time.Duration(attendanceBufferMinutes) * time.Minute
	earliestCheckout := shiftEndTime.Add(-buffer)

	if now.Before(earliestCheckout) {
		s.Logger.LogEndWithError("AttendanceQRCheckOut", "Check-out before window")
		return nil, &helpers.CustomError{Message: fmt.Sprintf("Belum waktunya check-out. Check-out dapat dilakukan setelah %s", earliestCheckout.Format("15:04"))}
	}

	// Overtime dihitung langsung dari shiftEnd (buffer hanya untuk checkin)
	overtimeMinutes := 0
	if now.After(shiftEndTime) {
		overtimeMinutes = int(now.Sub(shiftEndTime).Minutes())
		s.Logger.LogStep("AttendanceQRCheckOut", "Overtime detected: %d minutes", overtimeMinutes)
	}

	// Duration dihitung dari max(checkIn, shiftStart) — early check-in tidak dihitung
	shiftStartTime := time.Date(activeSession.Date.Year(), activeSession.Date.Month(), activeSession.Date.Day(), shiftStart.Hour(), shiftStart.Minute(), 0, 0, now.Location())
	durationStart := *activeSession.TimeIn
	if durationStart.Before(shiftStartTime) {
		durationStart = shiftStartTime
	}
	duration := calculateDuration(durationStart, now)
	durationMinutes := int(now.Sub(durationStart).Minutes())

	// Update attendance record
	var result *models.Attendance
	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		var err error
		result, err = s.repo.Attendance.UpdateMap(tx, &models.Attendance{ID: activeSession.ID}, map[string]interface{}{
			"time_out":         &now,
			"duration":         duration,
			"duration_minutes": durationMinutes,
			"overtime_minutes": overtimeMinutes,
		})
		if err != nil {
			return nil, err
		}

		msg := fmt.Sprintf("Check-out via QR Code berhasil. Durasi kerja: %s", duration)
		if overtimeMinutes > 0 {
			msg += fmt.Sprintf(" (Lembur: %d menit)", overtimeMinutes)
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "info",
			Title:   "Check-out Berhasil (QR Code)",
			Message: msg,
			Data: map[string]interface{}{
				"id":               result.ID,
				"user_id":          userID,
				"duration":         duration,
				"overtime_minutes": overtimeMinutes,
			},
		})

		return result, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("AttendanceQRCheckOut", "Failed to update attendance: %v", err)
		return nil, err
	}

	result = res.(*models.Attendance)
	dto := dtos.ToAttendanceDTO(result)
	s.Logger.LogEnd("AttendanceQRCheckOut", "QR check-out successful: user %d, duration %s, overtime %dm", userID, duration, overtimeMinutes)
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
	shiftStart, _ := time.Parse("15:04", activeSession.Shift.StartTime)
	shiftEndDay := activeSession.Date
	if shiftEnd.Hour()*60+shiftEnd.Minute() < shiftStart.Hour()*60+shiftStart.Minute() {
		// Cross-midnight shift: end time is on the day after Date
		shiftEndDay = activeSession.Date.AddDate(0, 0, 1)
	}
	shiftEndTime := time.Date(shiftEndDay.Year(), shiftEndDay.Month(), shiftEndDay.Day(), shiftEnd.Hour(), shiftEnd.Minute(), 0, 0, now.Location())

	buffer := time.Duration(attendanceBufferMinutes) * time.Minute
	earliestCheckout := shiftEndTime.Add(-buffer)

	if now.Before(earliestCheckout) {
		s.Logger.LogEndWithError("AttendanceCheckOut", "Check-out before window: %s < %s", now.Format("15:04"), earliestCheckout.Format("15:04"))
		return nil, &helpers.CustomError{Message: fmt.Sprintf("Belum waktunya check-out. Check-out dapat dilakukan setelah %s", earliestCheckout.Format("15:04"))}
	}

	// Overtime dihitung langsung dari shiftEnd (buffer hanya untuk checkin)
	overtimeMinutes := 0
	if now.After(shiftEndTime) {
		overtimeMinutes = int(now.Sub(shiftEndTime).Minutes())
		s.Logger.LogStep("AttendanceCheckOut", "Overtime detected: %d minutes (checkout %s > shiftEnd %s)", overtimeMinutes, now.Format("15:04"), shiftEndTime.Format("15:04"))
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

	// Duration dihitung dari max(checkIn, shiftStart) — early check-in tidak dihitung
	shiftStartTime := time.Date(activeSession.Date.Year(), activeSession.Date.Month(), activeSession.Date.Day(), shiftStart.Hour(), shiftStart.Minute(), 0, 0, now.Location())
	durationStart := *activeSession.TimeIn
	if durationStart.Before(shiftStartTime) {
		durationStart = shiftStartTime
	}
	duration := calculateDuration(durationStart, now)
	durationMinutes := int(now.Sub(durationStart).Minutes())

	// Update attendance record
	var result *models.Attendance
	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		var err error
		result, err = s.repo.Attendance.UpdateMap(tx, &models.Attendance{ID: activeSession.ID}, map[string]interface{}{
			"time_out":         &now,
			"lat_out":          &req.Lat,
			"lng_out":          &req.Lng,
			"image_out":        fotoPath,
			"duration":         duration,
			"duration_minutes": durationMinutes,
			"overtime_minutes": overtimeMinutes,
		})
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
		if _, rollbackErr := helpers.MoveFile(req.Image, "storage/attendance-evidence", "storage/tmp"); rollbackErr != nil {
			s.Logger.LogEndWithError("AttendanceCheckOut", "Failed to rollback photo move: %v", rollbackErr)
		}
		return nil, err
	}

	result = res.(*models.Attendance)
	dto := dtos.ToAttendanceDTO(result)
	s.Logger.LogEnd("AttendanceCheckOut", "Check-out successful: user %d, duration %s, overtime %dm", userID, duration, overtimeMinutes)
	return &dto, nil
}

// shiftMatch holds the result of findApplicableShift.
type shiftMatch struct {
	Assignment *models.UserShiftAssignment
	ShiftDate  time.Time // calendar date the shift belongs to (may be yesterday for cross-midnight)
	WindowEnd  time.Time
}

// findApplicableShift finds the shift whose window matches the current time.
// Returns the shift assignment whose [start-buffer, end] window contains `now`.
// If multiple shifts overlap, returns the one with the earliest end time.
// Supports cross-midnight shifts (e.g., 22:00–06:00).
func (s *Services) findApplicableShift(userID uint, now time.Time) *shiftMatch {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)

	// Fetch assignments for today and yesterday — yesterday is needed for
	// cross-midnight shifts where check-in happens after midnight.
	assignments, err := s.repo.UserShiftAssignment.FindAllActiveForUserDate(nil, userID, today, "Shift")
	if err != nil {
		return nil
	}
	if yAssignments, err := s.repo.UserShiftAssignment.FindAllActiveForUserDate(nil, userID, yesterday, "Shift"); err == nil {
		assignments = append(assignments, yAssignments...)
	}

	// Deduplicate by ID (ongoing assignments appear in both today and yesterday queries)
	seen := make(map[uint]bool)
	var deduped []models.UserShiftAssignment
	for _, a := range assignments {
		if !seen[a.ID] {
			seen[a.ID] = true
			deduped = append(deduped, a)
		}
	}
	if len(deduped) == 0 {
		return nil
	}

	buffer := time.Duration(attendanceBufferMinutes) * time.Minute
	var best *shiftMatch

	for i := range deduped {
		assignment := &deduped[i]
		shiftStart, parseErr := time.Parse("15:04", assignment.Shift.StartTime)
		shiftEnd, parseEndErr := time.Parse("15:04", assignment.Shift.EndTime)
		if parseErr != nil || parseEndErr != nil {
			continue
		}

		var windowStart, windowEnd, shiftDate time.Time

		if shiftEnd.Before(shiftStart) {
			// Cross-midnight shift (e.g., 22:00–06:00)
			nowMins := now.Hour()*60 + now.Minute()
			endMins := shiftEnd.Hour()*60 + shiftEnd.Minute()

			if nowMins < endMins {
				// After-midnight portion — shift belongs to yesterday
				shiftDate = yesterday
				windowStart = time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), shiftStart.Hour(), shiftStart.Minute(), 0, 0, now.Location()).Add(-buffer)
				windowEnd = time.Date(now.Year(), now.Month(), now.Day(), shiftEnd.Hour(), shiftEnd.Minute(), 0, 0, now.Location())
			} else {
				// Before-midnight portion — shift starts today, ends tomorrow
				tomorrow := today.AddDate(0, 0, 1)
				shiftDate = today
				windowStart = time.Date(now.Year(), now.Month(), now.Day(), shiftStart.Hour(), shiftStart.Minute(), 0, 0, now.Location()).Add(-buffer)
				windowEnd = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), shiftEnd.Hour(), shiftEnd.Minute(), 0, 0, now.Location())
			}
		} else {
			// Normal shift
			shiftDate = today
			windowStart = time.Date(now.Year(), now.Month(), now.Day(), shiftStart.Hour(), shiftStart.Minute(), 0, 0, now.Location()).Add(-buffer)
			windowEnd = time.Date(now.Year(), now.Month(), now.Day(), shiftEnd.Hour(), shiftEnd.Minute(), 0, 0, now.Location())
		}

		if now.After(windowStart) && !now.After(windowEnd) {
			if best == nil || windowEnd.Before(best.WindowEnd) {
				best = &shiftMatch{
					Assignment: assignment,
					ShiftDate:  shiftDate,
					WindowEnd:  windowEnd,
				}
			}
		}
	}

	return best
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

// AttendanceHistory returns paginated attendance history for the current user.
func (s *Services) AttendanceHistory(ctx context.Context, req dtos.AttendanceHistoryRequest) (*repositories.PagedResult[dtos.AttendanceHistoryDTO], error) {
	s.Logger.LogStart("AttendanceHistory", "Fetching attendance history")

	userID := helpers.GetCallerID(ctx)
	if userID == 0 {
		s.Logger.LogEndWithError("AttendanceHistory", "Invalid token: caller ID not found")
		return nil, helpers.ErrInvalidToken
	}

	var dateFrom, dateTo *time.Time
	if req.DateFrom != nil {
		parsed, err := time.Parse("2006-01-02", *req.DateFrom)
		if err == nil {
			dateFrom = &parsed
		}
	}
	if req.DateTo != nil {
		parsed, err := time.Parse("2006-01-02", *req.DateTo)
		if err == nil {
			dateTo = &parsed
		}
	}

	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	result, err := s.repo.Attendance.FindHistory(nil, userID, dateFrom, dateTo, req.Status, page, pageSize, "User", "Shift", "Office")
	if err != nil {
		s.Logger.LogEndWithError("AttendanceHistory", "Failed to fetch history: %v", err)
		return nil, err
	}

	// Convert to DTOs with related data
	historyDTOs := make([]dtos.AttendanceHistoryDTO, len(result.Data))
	for i, att := range result.Data {
		dto := dtos.AttendanceHistoryDTO{
			AttendanceDTO: dtos.ToAttendanceDTO(&att),
		}
		if att.User.ID != 0 {
			dto.UserName = att.User.Name
		}
		if att.Shift.ID != 0 {
			dto.ShiftName = att.Shift.Name
		}
		if att.Office.ID != 0 {
			dto.OfficeName = att.Office.Name
		}
		historyDTOs[i] = dto
	}

	s.Logger.LogEnd("AttendanceHistory", "Found %d records", result.Total)

	return &repositories.PagedResult[dtos.AttendanceHistoryDTO]{
		Data:       historyDTOs,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// AttendanceHistoryAll returns paginated attendance history for all users (HR only).
func (s *Services) AttendanceHistoryAll(ctx context.Context, req dtos.AttendanceHistoryRequest) (*repositories.PagedResult[dtos.AttendanceHistoryDTO], error) {
	s.Logger.LogStart("AttendanceHistoryAll", "Fetching all attendance history")

	var dateFrom, dateTo *time.Time
	if req.DateFrom != nil {
		parsed, err := time.Parse("2006-01-02", *req.DateFrom)
		if err == nil {
			dateFrom = &parsed
		}
	}
	if req.DateTo != nil {
		parsed, err := time.Parse("2006-01-02", *req.DateTo)
		if err == nil {
			dateTo = &parsed
		}
	}

	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	result, err := s.repo.Attendance.FindHistoryAll(nil, dateFrom, dateTo, req.Status, page, pageSize, "User", "Shift", "Office")
	if err != nil {
		s.Logger.LogEndWithError("AttendanceHistoryAll", "Failed to fetch history: %v", err)
		return nil, err
	}

	// Convert to DTOs with related data
	historyDTOs := make([]dtos.AttendanceHistoryDTO, len(result.Data))
	for i, att := range result.Data {
		dto := dtos.AttendanceHistoryDTO{
			AttendanceDTO: dtos.ToAttendanceDTO(&att),
		}
		if att.User.ID != 0 {
			dto.UserName = att.User.Name
		}
		if att.Shift.ID != 0 {
			dto.ShiftName = att.Shift.Name
		}
		if att.Office.ID != 0 {
			dto.OfficeName = att.Office.Name
		}
		historyDTOs[i] = dto
	}

	s.Logger.LogEnd("AttendanceHistoryAll", "Found %d records", result.Total)

	return &repositories.PagedResult[dtos.AttendanceHistoryDTO]{
		Data:       historyDTOs,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// AttendanceMonthlyStats returns monthly attendance statistics for the current user.
func (s *Services) AttendanceMonthlyStats(ctx context.Context, year, month int) (*dtos.MonthlyStatsResponse, error) {
	s.Logger.LogStart("AttendanceMonthlyStats", "Fetching monthly stats for %d-%d", year, month)

	userID := helpers.GetCallerID(ctx)
	if userID == 0 {
		s.Logger.LogEndWithError("AttendanceMonthlyStats", "Invalid token: caller ID not found")
		return nil, helpers.ErrInvalidToken
	}

	present, late, absent, overtime, avgDurationMins, err := s.repo.Attendance.FindMonthlyStats(nil, userID, year, month)
	if err != nil {
		s.Logger.LogEndWithError("AttendanceMonthlyStats", "Failed to fetch stats: %v", err)
		return nil, err
	}

	avgDuration := "0h 0m"
	if avgDurationMins > 0 {
		h := int(avgDurationMins) / 60
		m := int(avgDurationMins) % 60
		avgDuration = fmt.Sprintf("%dh %dm", h, m)
	}

	response := &dtos.MonthlyStatsResponse{
		TotalPresent:  present,
		TotalLate:     late,
		TotalAbsent:   absent,
		TotalOvertime: overtime,
		AvgDuration:   avgDuration,
	}

	s.Logger.LogEnd("AttendanceMonthlyStats", "Stats: present=%d, late=%d, absent=%d, overtime=%d", present, late, absent, overtime)
	return response, nil
}

// AttendanceCorrect corrects an attendance record (HR only).
func (s *Services) AttendanceCorrect(ctx context.Context, id uint, req dtos.AttendanceCorrectRequest) (*dtos.AttendanceDTO, error) {
	s.Logger.LogStart("AttendanceCorrect", "Correcting attendance record %d", id)

	callerID := helpers.GetCallerID(ctx)
	if callerID == 0 {
		s.Logger.LogEndWithError("AttendanceCorrect", "Invalid token: caller ID not found")
		return nil, helpers.ErrInvalidToken
	}

	// Fetch existing record with shift
	existing, err := s.repo.Attendance.FindByID(nil, id, "Shift")
	if err != nil {
		s.Logger.LogEndWithError("AttendanceCorrect", "Attendance record not found: %v", err)
		return nil, helpers.ErrNotFound
	}

	if req.TimeIn == nil || req.TimeOut == nil {
		return nil, &helpers.CustomError{Message: "time_in dan time_out wajib diisi"}
	}

	// Recalculate status based on new time_in vs shift start + buffer
	// Gunakan existing.Date sebagai base (benar untuk cross-midnight check-in setelah tengah malam)
	buffer := time.Duration(attendanceBufferMinutes) * time.Minute
	shiftStart, _ := time.Parse("15:04", existing.Shift.StartTime)
	shiftStartTime := time.Date(existing.Date.Year(), existing.Date.Month(), existing.Date.Day(), shiftStart.Hour(), shiftStart.Minute(), 0, 0, req.TimeIn.Location())
	flexiThreshold := shiftStartTime.Add(buffer)

	status := "present"
	if req.TimeIn.After(flexiThreshold) {
		status = "late"
	}

	// Duration dari max(timeIn, shiftStart) — early check-in tidak dihitung
	// Gunakan existing.Date sebagai base shiftStart (benar untuk shift cross-midnight)
	shiftStartForDuration := time.Date(existing.Date.Year(), existing.Date.Month(), existing.Date.Day(), shiftStart.Hour(), shiftStart.Minute(), 0, 0, req.TimeIn.Location())
	durationStart := *req.TimeIn
	if durationStart.Before(shiftStartForDuration) {
		durationStart = shiftStartForDuration
	}
	duration := calculateDuration(durationStart, *req.TimeOut)
	durationMinutes := int(req.TimeOut.Sub(durationStart).Minutes())

	// Overtime dihitung langsung dari shiftEnd (buffer hanya untuk checkin)
	overtimeMinutes := 0
	shiftEnd, _ := time.Parse("15:04", existing.Shift.EndTime)
	shiftEndTime := time.Date(req.TimeOut.Year(), req.TimeOut.Month(), req.TimeOut.Day(), shiftEnd.Hour(), shiftEnd.Minute(), 0, 0, req.TimeOut.Location())
	if req.TimeOut.After(shiftEndTime) {
		overtimeMinutes = int(req.TimeOut.Sub(shiftEndTime).Minutes())
	}

	now := time.Now()
	var result *models.Attendance

	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		var err error
		result, err = s.repo.Attendance.UpdateMap(tx, &models.Attendance{ID: id}, map[string]interface{}{
			"time_in":           req.TimeIn,
			"time_out":          req.TimeOut,
			"status":            status,
			"duration":          duration,
			"duration_minutes":  durationMinutes,
			"overtime_minutes":  overtimeMinutes,
			"corrected_by":      callerID,
			"corrected_at":      now,
			"correction_reason": req.CorrectionReason,
		})
		if err != nil {
			return nil, err
		}

		_ = s.NotificationCreate(ctx, &NotificationCreateParams{
			Type:    "warning",
			Title:   "Absensi Dikoreksi",
			Message: fmt.Sprintf("Record absensi #%d telah dikoreksi oleh HR", id),
			Data: map[string]interface{}{
				"id":           id,
				"user_id":      result.UserID,
				"corrected_by": callerID,
			},
		})

		return result, nil
	})
	if err != nil {
		s.Logger.LogEndWithError("AttendanceCorrect", "Failed to correct attendance: %v", err)
		return nil, err
	}

	result = res.(*models.Attendance)
	dto := dtos.ToAttendanceDTO(result)
	s.Logger.LogEnd("AttendanceCorrect", "Attendance %d corrected by user %d, status=%s, duration=%s", id, callerID, status, duration)
	return &dto, nil
}

// AttendanceReport returns attendance report with filters (HR only).
func (s *Services) AttendanceReport(ctx context.Context, req dtos.AttendanceReportRequest) (*repositories.PagedResult[dtos.AttendanceHistoryDTO], error) {
	s.Logger.LogStart("AttendanceReport", "Fetching attendance report")

	var dateFrom, dateTo *time.Time
	if req.DateFrom != nil {
		parsed, err := time.Parse("2006-01-02", *req.DateFrom)
		if err == nil {
			dateFrom = &parsed
		}
	}
	if req.DateTo != nil {
		parsed, err := time.Parse("2006-01-02", *req.DateTo)
		if err == nil {
			dateTo = &parsed
		}
	}

	page := req.Page
	pageSize := req.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	// Use FindHistoryAll with optional user_id filter
	var result *repositories.PagedResult[models.Attendance]
	var err error

	if req.UserID != nil {
		result, err = s.repo.Attendance.FindHistory(nil, *req.UserID, dateFrom, dateTo, req.Status, page, pageSize, "User", "Shift", "Office")
	} else {
		result, err = s.repo.Attendance.FindHistoryAll(nil, dateFrom, dateTo, req.Status, page, pageSize, "User", "Shift", "Office")
	}

	if err != nil {
		s.Logger.LogEndWithError("AttendanceReport", "Failed to fetch report: %v", err)
		return nil, err
	}

	// Convert to DTOs with related data
	reportDTOs := make([]dtos.AttendanceHistoryDTO, len(result.Data))
	for i, att := range result.Data {
		dto := dtos.AttendanceHistoryDTO{
			AttendanceDTO: dtos.ToAttendanceDTO(&att),
		}
		if att.User.ID != 0 {
			dto.UserName = att.User.Name
		}
		if att.Shift.ID != 0 {
			dto.ShiftName = att.Shift.Name
		}
		if att.Office.ID != 0 {
			dto.OfficeName = att.Office.Name
		}
		reportDTOs[i] = dto
	}

	s.Logger.LogEnd("AttendanceReport", "Found %d records", result.Total)

	return &repositories.PagedResult[dtos.AttendanceHistoryDTO]{
		Data:       reportDTOs,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// AttendanceLateStats returns late attendance statistics (HR only).
func (s *Services) AttendanceLateStats(ctx context.Context, dateFrom, dateTo *string, userID *uint, page, pageSize int) (*dtos.LateStatsResponse, error) {
	s.Logger.LogStart("AttendanceLateStats", "Fetching late statistics")

	var df, dt *time.Time
	if dateFrom != nil {
		parsed, err := time.Parse("2006-01-02", *dateFrom)
		if err == nil {
			df = &parsed
		}
	}
	if dateTo != nil {
		parsed, err := time.Parse("2006-01-02", *dateTo)
		if err == nil {
			dt = &parsed
		}
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	result, err := s.repo.Attendance.FindLateStats(nil, df, dt, userID, page, pageSize, "User", "Shift", "Office")
	if err != nil {
		s.Logger.LogEndWithError("AttendanceLateStats", "Failed to fetch late stats: %v", err)
		return nil, err
	}

	// Calculate statistics
	totalLateDays := len(result.Data)
	totalLateMinutes := 0
	trendMap := make(map[string]*dtos.LateTrendDTO)
	details := make([]dtos.LateDetailDTO, len(result.Data))

	for i, att := range result.Data {
		// Calculate late minutes (difference between time_in and shift start)
		lateMinutes := 0
		if att.TimeIn != nil && att.Shift.StartTime != "" {
			shiftStart, _ := time.Parse("15:04", att.Shift.StartTime)
			// Use att.Date (shiftDate) as the base — correct for cross-midnight check-ins after midnight
			// where att.TimeIn.Date() would be the next day but shiftStart belongs to att.Date.
			shiftStartTime := time.Date(att.Date.Year(), att.Date.Month(), att.Date.Day(), shiftStart.Hour(), shiftStart.Minute(), 0, 0, att.TimeIn.Location())
			diff := att.TimeIn.Sub(shiftStartTime)
			if diff > 0 {
				lateMinutes = int(diff.Minutes())
			}
		}
		totalLateMinutes += lateMinutes

		// Build trend
		dateStr := att.Date.Format("2006-01-02")
		if trend, exists := trendMap[dateStr]; exists {
			trend.LateCount++
			trend.TotalMinutes += lateMinutes
		} else {
			trendMap[dateStr] = &dtos.LateTrendDTO{
				Date:         dateStr,
				LateCount:    1,
				TotalMinutes: lateMinutes,
			}
		}

		// Build detail
		userName := ""
		if att.User.ID != 0 {
			userName = att.User.Name
		}
		shiftName := ""
		if att.Shift.ID != 0 {
			shiftName = att.Shift.Name
		}
		officeName := ""
		if att.Office.ID != 0 {
			officeName = att.Office.Name
		}

		timeIn := time.Time{}
		if att.TimeIn != nil {
			timeIn = *att.TimeIn
		}

		details[i] = dtos.LateDetailDTO{
			ID:          att.ID,
			UserID:      att.UserID,
			UserName:    userName,
			Date:        att.Date,
			ShiftName:   shiftName,
			TimeIn:      timeIn,
			LateMinutes: lateMinutes,
			OfficeName:  officeName,
		}
	}

	// Convert trend map to sorted slice
	trend := make([]dtos.LateTrendDTO, 0, len(trendMap))
	for _, t := range trendMap {
		trend = append(trend, *t)
	}
	sort.Slice(trend, func(i, j int) bool {
		return trend[i].Date < trend[j].Date
	})

	avgLateMinutes := float64(0)
	if totalLateDays > 0 {
		avgLateMinutes = float64(totalLateMinutes) / float64(totalLateDays)
	}

	response := &dtos.LateStatsResponse{
		TotalLateDays:  totalLateDays,
		TotalRecords:   result.Total,
		AvgLateMinutes: avgLateMinutes,
		Trend:          trend,
		Details:        details,
	}

	s.Logger.LogEnd("AttendanceLateStats", "Found %d late records, avg %.1f minutes", totalLateDays, avgLateMinutes)
	return response, nil
}
