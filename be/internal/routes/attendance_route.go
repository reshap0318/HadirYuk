package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/reshap0318/hadirYuk/internal/handlers"
)

func RegisterAttendanceRoutes(r *gin.RouterGroup, handlers *handlers.Handlers) {
	attendance := r.Group("/attendance")
	{
		attendance.GET("/today-status", handlers.GetTodayStatus)
		attendance.POST("/checkin", handlers.AttendanceCheckIn)
		attendance.POST("/nearest-office", handlers.GetNearestOffice)
	}
}
