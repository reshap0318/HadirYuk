package services

import (
	"context"
	"fmt"
	"time"

	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
)

func (s *Services) EmployeeDashboard(ctx context.Context) (*dtos.EmployeeDashboardResponse, error) {
	s.Logger.LogStart("EmployeeDashboard", "Fetching employee dashboard data")

	callerID := helpers.GetCallerID(ctx)
	if callerID == 0 {
		s.Logger.LogEndWithError("EmployeeDashboard", "Invalid token")
		return nil, helpers.ErrInvalidToken
	}

	now := time.Now()

	// Today status
	todayStatus, err := s.GetTodayStatus(ctx)
	if err != nil {
		s.Logger.LogWarn("EmployeeDashboard", "Failed to get today status: %v", err)
	}

	// Monthly stats
	monthlyStats, err := s.AttendanceMonthlyStats(ctx, now.Year(), int(now.Month()))
	if err != nil {
		s.Logger.LogWarn("EmployeeDashboard", "Failed to get monthly stats: %v", err)
	}

	// Upcoming shifts (next 7 days)
	upcomingShifts, err := s.getUpcomingShifts(ctx, callerID, now)
	if err != nil {
		s.Logger.LogWarn("EmployeeDashboard", "Failed to get upcoming shifts: %v", err)
	}

	response := &dtos.EmployeeDashboardResponse{
		TodayStatus:    todayStatus,
		MonthlyStats:   monthlyStats,
		UpcomingShifts: upcomingShifts,
	}

	s.Logger.LogEnd("EmployeeDashboard", "Dashboard data fetched for user %d", callerID)
	return response, nil
}

func (s *Services) getUpcomingShifts(ctx context.Context, userID uint, now time.Time) ([]dtos.UpcomingShiftDTO, error) {
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endDate := startDate.AddDate(0, 0, 7)

	assignments, err := s.repo.UserShiftAssignment.FindAllActiveForUserDateRange(nil, userID, startDate, endDate, "Shift")
	if err != nil {
		return nil, err
	}

	dayNames := map[time.Weekday]string{
		time.Monday:    "Senin",
		time.Tuesday:   "Selasa",
		time.Wednesday: "Rabu",
		time.Thursday:  "Kamis",
		time.Friday:    "Jumat",
		time.Saturday:  "Sabtu",
		time.Sunday:    "Minggu",
	}

	dateMap := make(map[string]bool)
	var result []dtos.UpcomingShiftDTO

	for day := startDate; day.Before(endDate); day = day.AddDate(0, 0, 1) {
		dateStr := day.Format("2006-01-02")
		for _, a := range assignments {
			if a.StartDate.After(day) || (a.EndDate != nil && a.EndDate.Before(day)) {
				continue
			}
			key := dateStr + "_" + fmt.Sprint(a.ShiftID)
			if dateMap[key] {
				continue
			}
			dateMap[key] = true
			result = append(result, dtos.UpcomingShiftDTO{
				Date:      dateStr,
				DayName:   dayNames[day.Weekday()],
				ShiftID:   a.Shift.ID,
				ShiftName: a.Shift.Name,
				StartTime: a.Shift.StartTime,
				EndTime:   a.Shift.EndTime,
				ColorCode: a.Shift.ColorCode,
			})
		}
	}

	return result, nil
}

func (s *Services) HRDashboard(ctx context.Context) (*dtos.HRDashboardResponse, error) {
	s.Logger.LogStart("HRDashboard", "Fetching HR dashboard data")

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Get total employees (exclude super admin)
	allUsers, err := s.repo.User.FindAllExceptSuperAdmin("Profile")
	if err != nil {
		s.Logger.LogEndWithError("HRDashboard", "Failed to fetch users: %v", err)
		return nil, err
	}
	totalEmployees := len(allUsers)

	// Get today's attendances
	todayAttendances, err := s.repo.Attendance.FindByDateAll(nil, today, "User", "Shift")
	if err != nil {
		s.Logger.LogWarn("HRDashboard", "Failed to fetch today attendances: %v", err)
	}

	// Count by status
	presentCount := 0
	lateCount := 0
	absentCount := 0
	totalOvertime := 0
	checkedInUserIDs := make(map[uint]bool)
	var recentActivity []dtos.RecentActivityDTO

	for _, att := range todayAttendances {
		checkedInUserIDs[att.UserID] = true
		switch att.Status {
		case "present":
			presentCount++
		case "late":
			lateCount++
		case "absent":
			absentCount++
		}
		totalOvertime += att.OvertimeMinutes

		// Build recent activity
		if att.TimeIn != nil {
			userName := ""
			if att.User.ID != 0 {
				userName = att.User.Name
			}
			shiftName := ""
			if att.Shift.ID != 0 {
				shiftName = att.Shift.Name
			}
			recentActivity = append(recentActivity, dtos.RecentActivityDTO{
				ID:        att.ID,
				UserID:    att.UserID,
				UserName:  userName,
				Avatar:    helpers.GetFileURL(att.User.Avatar),
				Action:    "checkin",
				ShiftName: shiftName,
				Time:      *att.TimeIn,
				Status:    att.Status,
			})
		}
		if att.TimeOut != nil {
			recentActivity = append(recentActivity, dtos.RecentActivityDTO{
				ID:        att.ID,
				UserID:    att.UserID,
				UserName:  att.User.Name,
				Avatar:    helpers.GetFileURL(att.User.Avatar),
				Action:    "checkout",
				ShiftName: att.Shift.Name,
				Time:      *att.TimeOut,
				Status:    att.Status,
			})
		}
	}

	// Sort recent activity by time descending (simple bubble sort)
	for i := 0; i < len(recentActivity); i++ {
		for j := i + 1; j < len(recentActivity); j++ {
			if recentActivity[j].Time.After(recentActivity[i].Time) {
				recentActivity[i], recentActivity[j] = recentActivity[j], recentActivity[i]
			}
		}
	}
	if len(recentActivity) > 10 {
		recentActivity = recentActivity[:10]
	}

	// Employees not yet checked in
	notYetCheckIn := 0
	employeeIDsWithShift := make(map[uint]bool)
	for _, u := range allUsers {
		assignments, _ := s.repo.UserShiftAssignment.FindAllActiveForUserDate(nil, u.ID, today)
		if len(assignments) > 0 {
			employeeIDsWithShift[u.ID] = true
		}
	}
	for userID := range employeeIDsWithShift {
		if !checkedInUserIDs[userID] {
			notYetCheckIn++
		}
	}

	// Department stats
	departmentMap := make(map[string]*dtos.DepartmentStatDTO)
	for _, u := range allUsers {
		dept := ""
		if u.Profile != nil {
			dept = u.Profile.Department
		}
		if dept == "" {
			dept = "Lainnya"
		}

		if _, exists := departmentMap[dept]; !exists {
			departmentMap[dept] = &dtos.DepartmentStatDTO{
				Department:    dept,
				TotalEmployees: 0,
			}
		}
		departmentMap[dept].TotalEmployees++

		if checkedInUserIDs[u.ID] {
			for _, att := range todayAttendances {
				if att.UserID == u.ID {
					switch att.Status {
					case "present":
						departmentMap[dept].Present++
					case "late":
						departmentMap[dept].Late++
					case "absent":
						departmentMap[dept].Absent++
					}
					break
				}
			}
		}
	}

	departmentStats := make([]dtos.DepartmentStatDTO, 0, len(departmentMap))
	for _, stat := range departmentMap {
		departmentStats = append(departmentStats, *stat)
	}

	s.Logger.LogEnd("HRDashboard", "Dashboard data: %d employees, %d present, %d late, %d absent", totalEmployees, presentCount, lateCount, absentCount)

	return &dtos.HRDashboardResponse{
		Date:            today.Format("2006-01-02"),
		TotalEmployees:  totalEmployees,
		Present:         presentCount,
		Late:            lateCount,
		Absent:          absentCount,
		NotYetCheckIn:   notYetCheckIn,
		TotalOvertime:   totalOvertime,
		DepartmentStats: departmentStats,
		RecentActivity:  recentActivity,
	}, nil
}

// AttendanceSchedule returns all employees' shift schedules for a date range.
func (s *Services) AttendanceSchedule(ctx context.Context, dateFrom, dateTo string) ([]dtos.ScheduleEmployeeDTO, error) {
	s.Logger.LogStart("AttendanceSchedule", "Fetching schedule from %s to %s", dateFrom, dateTo)

	callerID := helpers.GetCallerID(ctx)
	if callerID == 0 {
		s.Logger.LogEndWithError("AttendanceSchedule", "Invalid token")
		return nil, helpers.ErrInvalidToken
	}

	from, err := time.Parse("2006-01-02", dateFrom)
	if err != nil {
		s.Logger.LogEndWithError("AttendanceSchedule", "Invalid date_from: %v", err)
		return nil, &helpers.FieldError{Field: "date_from", Message: "Format tanggal tidak valid, gunakan YYYY-MM-DD"}
	}

	to, err := time.Parse("2006-01-02", dateTo)
	if err != nil {
		s.Logger.LogEndWithError("AttendanceSchedule", "Invalid date_to: %v", err)
		return nil, &helpers.FieldError{Field: "date_to", Message: "Format tanggal tidak valid, gunakan YYYY-MM-DD"}
	}

	// Get all active assignments within the date range
	assignments, err := s.repo.UserShiftAssignment.FindAllInDateRange(nil, from, to, "User", "Shift")
	if err != nil {
		s.Logger.LogEndWithError("AttendanceSchedule", "Failed to fetch assignments: %v", err)
		return nil, err
	}

	// Group by user
	userMap := make(map[uint]*dtos.ScheduleEmployeeDTO)
	for _, a := range assignments {
		if _, exists := userMap[a.UserID]; !exists {
			userName := ""
			avatar := ""
			if a.User.ID != 0 {
				userName = a.User.Name
				avatar = helpers.GetFileURL(a.User.Avatar)
			}
			userMap[a.UserID] = &dtos.ScheduleEmployeeDTO{
				UserID:   a.UserID,
				UserName: userName,
				Avatar:   avatar,
				Shifts:   []dtos.ScheduleShiftDTO{},
			}
		}

		// Expand assignments into individual dates
		start := a.StartDate
		if start.Before(from) {
			start = from
		}
		end := to
		if a.EndDate != nil && a.EndDate.Before(end) {
			end = *a.EndDate
		}

		for day := start; !day.After(end); day = day.AddDate(0, 0, 1) {
			shiftName := ""
			startTime := ""
			endTime := ""
			colorCode := ""
			if a.Shift.ID != 0 {
				shiftName = a.Shift.Name
				startTime = a.Shift.StartTime
				endTime = a.Shift.EndTime
				colorCode = a.Shift.ColorCode
			}
			userMap[a.UserID].Shifts = append(userMap[a.UserID].Shifts, dtos.ScheduleShiftDTO{
				Date:      day.Format("2006-01-02"),
				ShiftID:   a.ShiftID,
				ShiftName: shiftName,
				StartTime: startTime,
				EndTime:   endTime,
				ColorCode: colorCode,
			})
		}
	}

	result := make([]dtos.ScheduleEmployeeDTO, 0, len(userMap))
	for _, user := range userMap {
		result = append(result, *user)
	}

	s.Logger.LogEnd("AttendanceSchedule", "Found %d employees with schedules", len(result))
	return result, nil
}
