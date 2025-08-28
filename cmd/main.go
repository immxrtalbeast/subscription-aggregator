package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/immxrtalbeast/subscription-aggregator/cmd/docs"
	"github.com/immxrtalbeast/subscription-aggregator/internal/config"
	"github.com/immxrtalbeast/subscription-aggregator/internal/controller"
	"github.com/immxrtalbeast/subscription-aggregator/internal/lib/logger/sl"
	"github.com/immxrtalbeast/subscription-aggregator/internal/lib/logger/slogpretty"
	"github.com/immxrtalbeast/subscription-aggregator/internal/service/subscription"
	"github.com/immxrtalbeast/subscription-aggregator/internal/storage/psql"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	gPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Subscribe API
// @version 1.0
// @description API для управления подписками
// @host localhost:8080
// @BasePath /api/v1

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
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	{
		api.POST("/create", subscriptionController.AddSubcription)
		api.GET("/:id", subscriptionController.Subscription)
		api.GET("/all", subscriptionController.ListSubscription)
		api.PUT("/update", subscriptionController.UpdateSubscription)
		api.DELETE("/:id", subscriptionController.DeleteSubscription)
		api.GET("/total", subscriptionController.TotalCost)
	}
	addr := ":" + cfg.Port
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Info("starting server", "port", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server failed to start", sl.Err(err))
			panic("fatal")
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	log.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown:", sl.Err(err))
		panic("fatal")
	}
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Error("failed to close db: ", sl.Err(err))
		panic("fatal")
	}

	log.Info("server exiting")

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
