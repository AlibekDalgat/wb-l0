package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) getOrder(c *gin.Context) {
	id := c.Param("id")
	info, err := h.services.Order.GetOrder(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, info)
}
