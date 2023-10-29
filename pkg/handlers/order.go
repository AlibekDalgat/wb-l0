package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	wb_l0 "wb-l0"
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

func (h *Handler) AddOrder(data []byte) error {
	var dataJson wb_l0.Order
	err := json.Unmarshal(data, &dataJson)
	if err != nil {
		return errors.New("Некорректное содержание сообщения в канале")
	}
	if dataJson.IsValid() {
		return h.services.AddOrder(dataJson, data)
	} else {
		return errors.New("Невалидный json")
	}
}
