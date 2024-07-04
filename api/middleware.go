package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rouclec/simplebank/token"
)

const (
	authHeaderKey  = "authorization"
	authTypeBearer = "Bearer "
	authPayloadKey = "auth_payload"
)

func extractTokenFromHeader(authorizationHeader string) (string, error) {
	if !strings.HasPrefix(authorizationHeader, authTypeBearer) {
		return "", errors.New("invalid token, not a bearer token") // Not a Bearer token
	}
	return strings.TrimSpace(strings.TrimPrefix(authorizationHeader, authTypeBearer)), nil // Extract and trim token
}

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var authToken string
		var err error
		authToken, err = ctx.Cookie("auth")

		if err != nil {
			authorizationHeader := ctx.GetHeader(authHeaderKey)
			if len(authorizationHeader) == 0 {
				err := errors.New("authorization header is not provided")
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
				return
			}
			authToken, err = extractTokenFromHeader(authorizationHeader)

			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
				return
			}
		}

		payload, err := tokenMaker.VerifyToken(authToken)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authPayloadKey, payload)
		ctx.Next()
	}
}
