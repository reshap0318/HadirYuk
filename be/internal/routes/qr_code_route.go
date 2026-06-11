package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/reshap0318/hadirYuk/internal/handlers"
	"github.com/reshap0318/hadirYuk/internal/helpers"
	"github.com/reshap0318/hadirYuk/internal/middleware"
)

func RegisterQRCodeRoutes(r *gin.RouterGroup, h *handlers.Handlers, acc *helpers.Access) {
	qrCodes := r.Group("/qr-codes")
	{
		qrCodes.POST("/generate", middleware.RequirePermission(acc, "qrcode.generate"), h.QRCodeGenerate)
		qrCodes.GET("", middleware.RequirePermission(acc, "qrcode.view"), h.QRCodeGetAll)
		qrCodes.POST("/:id/revoke", middleware.RequirePermission(acc, "qrcode.revoke"), h.QRCodeRevoke)
	}
}
