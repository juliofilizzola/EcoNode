package middleware

import (
	"net/http"
	"strings"

	"EcoNode/internal/dto"
	"EcoNode/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	bearerScheme        = "Bearer "
	claimsContextKey    = "claims"
)

func Authenticate(tokens service.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString, found := strings.CutPrefix(ctx.GetHeader(authorizationHeader), bearerScheme)
		if !found || tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "missing bearer token"})
			return
		}

		claims, err := tokens.Parse(tokenString)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Message: "invalid or expired token"})
			return
		}

		ctx.Set(claimsContextKey, claims)
		ctx.Next()
	}
}

func ClaimsFrom(ctx *gin.Context) (service.Claims, bool) {
	value, exists := ctx.Get(claimsContextKey)
	if !exists {
		return service.Claims{}, false
	}

	claims, ok := value.(service.Claims)

	return claims, ok
}
