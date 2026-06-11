package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/reshap0318/hadirYuk/internal/dtos"
	"github.com/reshap0318/hadirYuk/internal/helpers"
)

func (h *Handlers) QRCodeGenerate(c *gin.Context) {
	var req dtos.QRCodeGenerateRequest
	if err := c.BindJSON(&req); err != nil {
		helpers.BadRequest(c, "Invalid JSON payload")
		return
	}
	if err := h.Validate.Struct(&req); err != nil {
		helpers.ValidationResponse(c, h.getErrorsMap(err))
		return
	}

	result, err := h.svcs.QRCodeGenerate(c.Request.Context(), &req)
	if err != nil {
		if err == helpers.ErrNotFound {
			helpers.NotFound(c, "Office not found")
			return
		}
		helpers.InternalServerError(c, "Failed to generate QR code")
		return
	}

	helpers.Created(c, "QR code berhasil dibuat", result)
}

func (h *Handlers) QRCodeGetAll(c *gin.Context) {
	result, err := h.svcs.QRCodeGetAll(c.Request.Context())
	if err != nil {
		helpers.InternalServerError(c, "Failed to fetch QR codes")
		return
	}

	helpers.OK(c, "QR codes berhasil diambil", result)
}

func (h *Handlers) QRCodeRevoke(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		helpers.BadRequest(c, "Invalid QR code ID")
		return
	}

	result, err := h.svcs.QRCodeRevoke(c.Request.Context(), uint(id))
	if err != nil {
		if err == helpers.ErrNotFound {
			helpers.NotFound(c, "QR code not found")
			return
		}
		helpers.InternalServerError(c, "Failed to revoke QR code")
		return
	}

	helpers.OK(c, "QR code berhasil dicabut", result)
}
