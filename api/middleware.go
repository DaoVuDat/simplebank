package api

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"simple_bank/token"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authorizationHeader := c.Request().Header.Get(authorizationHeaderKey)

			if len(authorizationHeader) == 0 {
				err := errors.New("authorization header is not provided")
				return c.JSON(http.StatusUnauthorized, err)
			}

			fields := strings.Fields(authorizationHeader)
			if len(fields) < 2 {
				err := errors.New("invalid authorization header format")
				return c.JSON(http.StatusUnauthorized, err)
			}

			authorizationType := strings.ToLower(fields[0])
			if authorizationType != authorizationTypeBearer {
				err := fmt.Errorf("unsupported authorization type %s", authorizationType)
				return c.JSON(http.StatusUnauthorized, err)

			}

			accessToken := fields[1]
			payload, err := tokenMaker.VerifyToken(accessToken)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, err)
			}

			c.Set(authorizationPayloadKey, payload)
			return next(c)
		}
	}
}

func test(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		return next(c)
	}
}
