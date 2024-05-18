package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-playground/validator"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/sousair/go-finance/internal/entities"
	"github.com/sousair/go-finance/internal/infra/cipher"
	"github.com/sousair/go-finance/internal/infra/database"
	"github.com/sousair/go-finance/internal/infra/token"
	httpxhandlers "github.com/sousair/go-finance/internal/presentation/httpx/handlers"
	httpxmiddleware "github.com/sousair/go-finance/internal/presentation/httpx/middlewares"
	"github.com/sousair/go-finance/internal/usecases"
)

var (
	postgresConnectionURL = os.Getenv("POSTGRES_CONNECTION_URL")
	port                  = os.Getenv("PORT")
	bcryptCostStr         = os.Getenv("BCRYPT_COST")
	userJwtSecret         = os.Getenv("JWT_USER_SECRET")
)

func main() {
	db, err := database.NewPostgres(postgresConnectionURL)

	if err != nil {
		panic(err)
	}

	userRepo, err := database.NewRepository[entities.User](db)

	if err != nil {
		panic(err)
	}

	assetRepo, err := database.NewRepository[entities.Asset](db)

	if err != nil {
		panic(err)
	}

	bcryptCost, err := strconv.Atoi(bcryptCostStr)

	if err != nil {
		panic(err)
	}

	cipher := cipher.NewBcrypt(bcryptCost)
	userJwt := token.NewJWT[usecases.UserTokenPayload](userJwtSecret)

	userAuthMiddleware := httpxmiddleware.NewUserAuthMiddleware(userJwt).Execute

	createUserUc := usecases.NewCreateUserUsecase(userRepo, cipher)
	userLoginUc := usecases.NewUserLoginUsecase(userRepo, cipher, userJwt)
	createAssetUc := usecases.NewCreateAssetUsecase(assetRepo)

	createUserHandler := httpxhandlers.NewCreateUserHandler(validator.New(), createUserUc)
	userLoginHandler := httpxhandlers.NewUserLoginHandler(validator.New(), userLoginUc)
	createAssetHandler := httpxhandlers.NewCreateAssetHandler(validator.New(), createAssetUc)

	e := echo.New()

	e.POST("/users", createUserHandler.Handle)
	e.POST("/users/login", userLoginHandler.Handle)
	e.POST("/assets", createAssetHandler.Handle, userAuthMiddleware)

	e.Logger.Fatal(e.Start(fmt.Sprintf(": %s", port)))
}
