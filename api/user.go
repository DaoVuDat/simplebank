package api

import (
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"net/http"
	db "simple_bank/db/sqlc"
	"simple_bank/util"
)

type createUserRequest struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type userResponse struct {
	Username          string             `json:"username"`
	FullName          string             `json:"full_name"`
	Email             string             `json:"email"`
	PasswordChangedAt pgtype.Timestamptz `json:"password_changed_at"`
	CreatedAt         pgtype.Timestamptz `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
}

func (server *Server) createUser(c echo.Context) error {
	var req createUserRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	err = c.Validate(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, struct {
			StatusCode int
			Message    string
		}{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		FullName:       req.FullName,
		Email:          req.Email,
		HashedPassword: hashPassword,
	}

	user, err := server.store.CreateUser(c.Request().Context(), arg)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return echo.NewHTTPError(http.StatusForbidden, err)
			}
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	resp := newUserResponse(user)

	return c.JSON(http.StatusOK, resp)
}

type loginUserRequest struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
	User        userResponse
}

func (server *Server) loginUser(c echo.Context) error {
	var req loginUserRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	err = c.Validate(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, struct {
			StatusCode int
			Message    string
		}{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	user, err := server.store.GetUser(c.Request().Context(), req.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := loginUserResponse{
		User:        newUserResponse(user),
		AccessToken: accessToken,
	}

	return c.JSON(http.StatusOK, resp)

}
