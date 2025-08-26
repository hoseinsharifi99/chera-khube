package app

import (
	"chera_khube/internal/database"
	"chera_khube/internal/helper"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application interface {
	Setup(ctx context.Context) error
}

type application struct {
	db     *gorm.DB
	config *helper.ServiceConfig
	ctx    context.Context
}

func NewApplication(config *helper.ServiceConfig) Application {
	return &application{config: config}
}

func (a *application) Setup(ctx context.Context) error {
	a.ctx = ctx
	defer a.handleCancelInterrupt(ctx)
	err := errors.Join(a.openDatabaseConnection())
	app := fx.New(
		fx.Provide(
			a.InitRouter,
			a.InitController,
			a.InitService,
			a.InitLogger,
			a.InitRepository,
		),
		fx.Invoke(func(router *gin.Engine) {
			addr := fmt.Sprintf("%s:%s", a.config.Http.Host, a.config.Http.Port)
			go func() {
				_ = router.Run(addr)
			}()
		}),

		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
	)
	app.Run()
	return err
}

func (a *application) handleCancelInterrupt(ctx context.Context) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	<-c
	fmt.Println("shutting down...")
	_, cancel := context.WithTimeout(ctx, time.Duration(a.config.System.ShutdownTimeout)*time.Second)
	defer cancel()
	err := database.Close(a.db)
	if err != nil {
		fmt.Printf("shut down error: %s\n", err.Error())
		return
	}
	fmt.Println("Fin")
}
