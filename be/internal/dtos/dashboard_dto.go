package dtos

import "time"

type EmployeeDashboardResponse struct {
	TodayStatus    *AttendanceTodayResponse `json:"today_status"`
	MonthlyStats   *MonthlyStatsResponse    `json:"monthly_stats"`
	UpcomingShifts []UpcomingShiftDTO       `json:"upcoming_shifts"`
}

type UpcomingShiftDTO struct {
	Date      string   `json:"date"`
	DayName   string   `json:"day_name"`
	ShiftID   uint     `json:"shift_id"`
	ShiftName string   `json:"shift_name"`
	StartTime string   `json:"start_time"`
	EndTime   string   `json:"end_time"`
	ColorCode string   `json:"color_code"`
}

type HRDashboardResponse struct {
	Date             string                   `json:"date"`
	TotalEmployees   int                      `json:"total_employees"`
	Present          int                      `json:"present"`
	Late             int                      `json:"late"`
	Absent           int                      `json:"absent"`
	NotYetCheckIn    int                      `json:"not_yet_check_in"`
	TotalOvertime    int                      `json:"total_overtime"`
	DepartmentStats  []DepartmentStatDTO      `json:"department_stats"`
	RecentActivity   []RecentActivityDTO      `json:"recent_activity"`
}

type DepartmentStatDTO struct {
	Department    string `json:"department"`
	TotalEmployees int   `json:"total_employees"`
	Present       int    `json:"present"`
	Late          int    `json:"late"`
	Absent        int    `json:"absent"`
}

type RecentActivityDTO struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	UserName   string    `json:"user_name"`
	Avatar     string    `json:"avatar"`
	Action     string    `json:"action"`
	ShiftName  string    `json:"shift_name"`
	Time       time.Time `json:"time"`
	Status     string    `json:"status"`
}

type ScheduleEmployeeDTO struct {
	UserID   uint              `json:"user_id"`
	UserName string            `json:"user_name"`
	Avatar   string            `json:"avatar"`
	Shifts   []ScheduleShiftDTO `json:"shifts"`
}

type ScheduleShiftDTO struct {
	Date      string `json:"date"`
	ShiftID   uint   `json:"shift_id"`
	ShiftName string `json:"shift_name"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	ColorCode string `json:"color_code"`
}
