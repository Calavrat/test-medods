package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SingIn(c *gin.Context) {
	userID := c.Param("id")

	ts, err := h.service.Authorization.GenerateToken(userID)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	tokens := map[string]string{
		"access_token":  ts.AccesToken,
		"refresh_token": ts.RefreshToken,
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) Refresh(c *gin.Context) {

	mapTokens := make(map[string]string)
	// access_token
	header := c.GetHeader("Authorization")
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}
	mapTokens["access_token"] = headerParts[1]

	//refresh_token
	if err := c.BindJSON(&mapTokens); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//проверка токенов на связь
	userId, err := h.service.ParseTokens(&mapTokens)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	//генерация новых токенов
	ts, err := h.service.Authorization.GenerateToken(userId)
	if err != nil {
		newErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	tokens := map[string]string{
		"access_token":  ts.AccesToken,
		"refresh_token": ts.RefreshToken,
	}

	c.JSON(http.StatusOK, tokens)
}
