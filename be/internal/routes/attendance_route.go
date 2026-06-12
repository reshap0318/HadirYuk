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
	}
}
