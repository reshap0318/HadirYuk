package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/handlers"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/middleware"
)

func RegisterAttendanceRoutes(r *gin.RouterGroup, handlers *handlers.Handlers, acc *helpers.Access) {
	attendance := r.Group("/attendance")
	{
		attendance.GET("/today", handlers.GetTodayStatus)
		attendance.POST("/checkin", middleware.RequirePermission(acc, "attendance.checkin"), handlers.AttendanceCheckIn)
		attendance.POST("/checkout", middleware.RequirePermission(acc, "attendance.checkout"), handlers.AttendanceCheckOut)
		attendance.POST("/checkin/qr", middleware.RequirePermission(acc, "attendance.checkin"), handlers.AttendanceQRCheckIn)
		attendance.POST("/checkout/qr", middleware.RequirePermission(acc, "attendance.checkout"), handlers.AttendanceQRCheckOut)
		attendance.POST("/nearest-office", handlers.GetNearestOffice)
		attendance.GET("/history", handlers.AttendanceHistory)
		attendance.GET("/stats", handlers.AttendanceMonthlyStats)
		attendance.PUT("/:id/correct", middleware.RequirePermission(acc, "attendance.correct"), handlers.AttendanceCorrect)
		attendance.GET("/late-statistics", middleware.RequirePermission(acc, "late-statistic.view"), handlers.AttendanceLateStats)
	}

	// Report routes
	reports := r.Group("/attendance/reports")
	{
		reports.GET("", middleware.RequirePermission(acc, "report.view"), handlers.AttendanceReport)
	}
}
