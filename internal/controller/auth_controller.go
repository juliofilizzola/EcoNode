package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"EcoNode/internal/dto"
	"EcoNode/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const bearerTokenType = "Bearer"

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var request dto.LoginRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "invalid request body",
			Details: validationDetails(err),
		})
		return
	}

	issued, err := c.authService.Login(ctx.Request.Context(), request.Username, request.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "invalid credentials"})
			return
		}

		log.Printf("login failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken: issued.Value,
		TokenType:   bearerTokenType,
		ExpiresIn:   int64(time.Until(issued.ExpiresAt).Seconds()),
	})
}

func validationDetails(err error) []string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return nil
	}

	details := make([]string, 0, len(validationErrors))
	for _, fieldError := range validationErrors {
		details = append(details, fmt.Sprintf("%s failed on '%s' validation", fieldError.Field(), fieldError.Tag()))
	}

	return details
}
