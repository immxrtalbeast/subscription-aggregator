package main

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/immxrtalbeast/subscription-aggregator/internal/config"
	"github.com/immxrtalbeast/subscription-aggregator/internal/controller"
	"github.com/immxrtalbeast/subscription-aggregator/internal/lib/logger/sl"
	"github.com/immxrtalbeast/subscription-aggregator/internal/lib/logger/slogpretty"
	"github.com/immxrtalbeast/subscription-aggregator/internal/service/subscription"
	"github.com/immxrtalbeast/subscription-aggregator/internal/storage/psql"
	gPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting subscription service")
	if err := runMigrations(cfg); err != nil {
		log.Error("failed to run migrations", sl.Err(err))
		panic("fatal")
	}
	log.Info("Migrations applied successfully")
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode)

	db, err := gorm.Open(gPostgres.New(gPostgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	subscriptionRepository := psql.NewSubscriptionRepository(db)
	subscriptionInteractor := subscription.NewSubscriptionInteractor(log, subscriptionRepository)
	subscriptionController := controller.NewSubscriptionController(subscriptionInteractor)
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		api.POST("/create", subscriptionController.AddSubcription)
		api.GET("/:id", subscriptionController.Subscription)
		api.GET("/all", subscriptionController.ListSubscription)
		api.PUT("/update", subscriptionController.UpdateSubscription)
		api.DELETE("/:id", subscriptionController.DeleteSubscription)
		api.GET("/total", subscriptionController.TotalCost)
	}
	router.Run(":8080")
}

func runMigrations(cfg *config.Config) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(cfg.DB.User),
		url.QueryEscape(cfg.DB.Password),
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
		cfg.DB.SSLMode,
	)

	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
