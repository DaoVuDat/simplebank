package api

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	db "simple_bank/db/sqlc"
	"simple_bank/token"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id"`
	ToAccountID   int64  `json:"to_account_id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
}

func (c transferRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.FromAccountID, validation.Required, validation.Min(1)),
		validation.Field(&c.ToAccountID, validation.Required, validation.Min(1)),
		validation.Field(&c.Amount, validation.Required, validation.Min(0)),
		validation.Field(&c.Currency, validation.Required, validation.In("USD", "EUR")))
}

func (server *Server) createTransfer(c echo.Context) error {
	var req transferRequest
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

	fromAccount, isValid, err := server.validAccount(c, req.FromAccountID, req.Currency)
	if !isValid {
		return err
	}

	authPayload := c.Get(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("account doesn't belong to the authenticated user"))
	}

	_, isValid, err = server.validAccount(c, req.ToAccountID, req.Currency)
	if !isValid {
		return err
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(c.Request().Context(), arg)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(c echo.Context, accountId int64, currency string) (db.Account, bool, error) {
	account, err := server.store.GetAccount(c.Request().Context(), accountId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return account, false, echo.NewHTTPError(http.StatusBadRequest, struct {
				StatusCode int
				Message    string
			}{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
			})
		}

		return account, false, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if account.Currency != currency {
		return account, false, echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency))
	}

	return account, true, nil
}
