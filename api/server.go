package api

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
	db "simple_bank/db/sqlc"
	"simple_bank/token"
	"simple_bank/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *echo.Echo
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func NewServer(config util.Config, store db.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{store: store, config: config, tokenMaker: tokenMaker}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := echo.New()
	router.Validator = &CustomValidator{validator: validator.New()}

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/", authMiddleware(server.tokenMaker))

	authRoutes.POST("accounts", server.createAccount)
	authRoutes.GET("accounts/:id", server.getAccount)
	authRoutes.GET("accounts", server.listAccount)

	authRoutes.POST("transfer", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) {
	server.router.Logger.Fatal(server.router.Start(address))
}
