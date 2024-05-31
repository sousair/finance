package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/sousair/go-finance/internal/entities"
	"github.com/sousair/go-finance/internal/infra/cipher"
	"github.com/sousair/go-finance/internal/infra/database"
	"github.com/sousair/go-finance/internal/infra/financeinquiry"
	"github.com/sousair/go-finance/internal/infra/token"
	httpxhandlers "github.com/sousair/go-finance/internal/presentation/httpx/handlers"
	httpxmiddleware "github.com/sousair/go-finance/internal/presentation/httpx/middlewares"
	"github.com/sousair/go-finance/internal/usecases"
	"github.com/sousair/go-finance/internal/workers"
)

var (
	postgresConnectionURL = os.Getenv("POSTGRES_CONNECTION_URL")
	port                  = os.Getenv("PORT")
	bcryptCostStr         = os.Getenv("BCRYPT_COST")
	userJwtSecret         = os.Getenv("JWT_USER_SECRET")
)

func main() {
	ctx := context.Background()
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

	assetPriceRepo, err := database.NewRepository[entities.AssetPrice](db)
	if err != nil {
		panic(err)
	}

	userInputRepo, err := database.NewRepository[entities.UserInput](db)
	if err != nil {
		panic(err)
	}

	userAssetRepo, err := database.NewRepository[entities.UserAsset](db)
	if err != nil {
		panic(err)
	}

	bcryptCost, err := strconv.Atoi(bcryptCostStr)
	if err != nil {
		panic(err)
	}

	cipher := cipher.NewBcrypt(bcryptCost)
	userJwt := token.NewJWT[usecases.UserTokenPayload](userJwtSecret)
	finance := financeinquiry.NewBrazilianFinance()

	userAuthMiddleware := httpxmiddleware.NewUserAuthMiddleware(userJwt).Execute

	createUserUc := usecases.NewCreateUserUsecase(userRepo, cipher)
	userLoginUc := usecases.NewUserLoginUsecase(userRepo, cipher, userJwt)
	createAssetUc := usecases.NewCreateAssetUsecase(assetRepo)
	addUserAssetUc := usecases.NewAddUserAssetUsecase(userAssetRepo)
	createUserInputUc := usecases.NewCreateUserInputUsecase(userInputRepo, userAssetRepo, addUserAssetUc)
	updateStockAssetsUc := usecases.NewUpdateStockAssetsUsecase(assetRepo, assetPriceRepo, finance)

	validator := validator.New()

	createUserHandler := httpxhandlers.NewCreateUserHandler(validator, createUserUc)
	userLoginHandler := httpxhandlers.NewUserLoginHandler(validator, userLoginUc)
	createAssetHandler := httpxhandlers.NewCreateAssetHandler(validator, createAssetUc)
	createUserInputHandler := httpxhandlers.NewCreateUserInputHandler(validator, createUserInputUc)

	e := echo.New()

	e.POST("/users", createUserHandler.Handle)
	e.POST("/users/login", userLoginHandler.Handle)

	e.POST("/users/input", createUserInputHandler.Handle, userAuthMiddleware)

	e.POST("/assets", createAssetHandler.Handle, userAuthMiddleware)

	cronWorker := workers.NewCronWorker()

	cronWorker.Register(workers.CronJob{
		Interval: time.Minute * 5,
		Name:     "UpdateStockAssets",
		Fn:       updateStockAssetsUc.UpdateStockAssets,
	})

	go cronWorker.Start(ctx)

	e.Logger.Fatal(e.Start(fmt.Sprintf(": %s", port)))
}
