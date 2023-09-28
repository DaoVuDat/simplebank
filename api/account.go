package api

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"net/http"
	db "simple_bank/db/sqlc"
	"simple_bank/token"
	"strconv"
)

type createAccountRequest struct {
	Currency string `json:"currency"`
}

func (c createAccountRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Currency, validation.Required, validation.In("USD", "EUR")))
}

func (server *Server) createAccount(c echo.Context) error {
	var req createAccountRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	err = req.Validate()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, struct {
			StatusCode int
			Message    string
		}{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	authPayload := c.Get(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(c.Request().Context(), arg)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.ForeignKeyViolation, pgerrcode.UniqueViolation:
				return echo.NewHTTPError(http.StatusForbidden, err)
			}
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, account)
}

func (server *Server) getAccount(c echo.Context) error {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	account, err := server.store.GetAccount(c.Request().Context(), int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	authPayload := c.Get(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("account doesn't belong to the authenticated user"))
	}

	return c.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageID   int32 `query:"page_id"`
	PageSize int32 `query:"page_size"`
}

func (l listAccountRequest) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.PageID, validation.Required, validation.Min(1)),
		validation.Field(&l.PageSize, validation.Required, validation.Min(5), validation.Max(10)),
	)
}

func (server *Server) listAccount(c echo.Context) error {
	var req listAccountRequest
	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	authPayload := c.Get(authorizationPayloadKey).(*token.Payload)

	arg := db.ListAccountParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccount(c.Request().Context(), arg)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, accounts)
}
