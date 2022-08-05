package app

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/zedisdog/sweetbean/database/seed"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Application interface {
	MigratorUp() error
	Seed() error
	Start()
	Stop(duration time.Duration)
}

var _ Application = (*App)(nil)

type App struct {
	Migrator    *migrate.Migrate
	HttpServers map[string]*http.Server
	Log         *zap.Logger
	RunWait     sync.WaitGroup
	Seeder      seed.GormSeedFunc
	DB          *gorm.DB
}

func (a App) MigratorUp() error {
	if a.Migrator != nil {
		return a.Migrator.Up()
	}
	return nil
}

func (a App) Seed() error {
	if a.Seeder != nil && a.DB != nil {
		return a.Seeder(a.DB)
	}
	return nil
}

func (a *App) Start() {
	for name, svr := range a.HttpServers {
		a.RunWait.Add(1)
		go func(name string, svr *http.Server) {
			if a.Log != nil {
				a.Log.Info("starting http server...\n", zap.String("name", name))
			}
			err := svr.ListenAndServe()
			if err != nil && a.Log != nil {
				a.Log.Warn("http server stopped", zap.String("name", name), zap.Error(err))
			}
			a.RunWait.Done()
		}(name, svr)
	}
}

func (a *App) Stop(duration time.Duration) {
	for name, svr := range a.HttpServers {
		go func(name string, svr *http.Server) {
			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()
			_ = svr.Shutdown(ctx)
		}(name, svr)
	}
}
