package main

import (
	"context"
	"crypto/subtle"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/utking/etcd-ui/internal/controllers"
	"github.com/utking/etcd-ui/internal/helpers/http"
	"github.com/utking/etcd-ui/internal/helpers/utils"
	"github.com/utking/etcd-ui/static"
	"golang.org/x/sync/errgroup"
)

const (
	fileAllRWMode = 0o666
	gracefullWait = time.Second
	gzipLevel     = 5
)

func main() {
	if envErr := godotenv.Load(".env.dev"); envErr != nil {
		if envErr = godotenv.Load(".env"); envErr != nil {
			fmt.Println(envErr)
		}
	}

	accessLogPath := fmt.Sprintf(
		"%s/%s-access.log",
		utils.ReadEnvVar("LOG_DIR_PATH", "."),
		time.Now().Format(time.DateOnly),
	)

	app := echo.New()

	// Enable BasicAuth only if the UI username is not empty
	if utils.GetUIUsername() != "" {
		var (
			uiUsername = utils.GetUIUsername()
			uiPassword = utils.GetUIPassword()
		)

		app.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
			Validator: func(username, password string, _ echo.Context) (bool, error) {
				// Be careful to use constant time comparison to prevent timing attacks
				if subtle.ConstantTimeCompare([]byte(username), []byte(uiUsername)) == 1 &&
					subtle.ConstantTimeCompare([]byte(password), []byte(uiPassword)) == 1 {
					return true, nil
				}

				return false, nil
			},
		}))
	}

	app.HideBanner = true
	app.HTTPErrorHandler = controllers.HTTPErrorHandler

	if err := http.InitTemplates(app); err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(accessLogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, fileAllRWMode)
	if err != nil {
		panic(fmt.Sprintf("error opening file for access logs: %v", err))
	}

	defer f.Close()

	app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: f,
	}))
	app.Use(middleware.Recover())

	app.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: gzipLevel,
	}))

	app.Use(middleware.Secure())
	app.StaticFS("/", static.StaticFiles)
	app.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf",
	}))

	controllers.Setup(app)

	// Graceful shutdown
	ctx, cancel := signal.NotifyContext(
		context.Background(), os.Interrupt, syscall.SIGTERM,
	)

	defer cancel()

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return app.Start(fmt.Sprintf("%s:8080", utils.ReadEnvVar("HOST", "")))
	})

	g.Go(func() error {
		<-gCtx.Done()
		time.Sleep(gracefullWait)

		return app.Shutdown(gCtx)
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("exit with: %v\n", err)
	}
}
