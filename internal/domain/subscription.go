package domain

import (
	"context"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ServiceName string     `gorm:"not null" json:"service_name"`
	Price       int        `gorm:"not null" json:"price"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	StartDate   MonthYear  `gorm:"not null" json:"start_date"`
	EndDate     *MonthYear `json:"end_date"`
}

type SubscriptionInteractor interface {
	AddSubscription(ctx context.Context, serviceName string, price int, userID uuid.UUID, startDate MonthYear, endDate *MonthYear) (uuid.UUID, error)
	Subscription(ctx context.Context, subscriptionID uuid.UUID) (*Subscription, error)
	DeleteSubscription(ctx context.Context, subscriptionID uuid.UUID) error
	UpdateSubscription(ctx context.Context, subscriptionID uuid.UUID, serviceName string, price int, userID uuid.UUID, startDate MonthYear, endDate *MonthYear) error
	ListSubscription(ctx context.Context, offset, limit int) ([]*Subscription, int64, error)
	TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate MonthYear) (int, error)
}

type SubscriptionRepository interface {
	SaveSubscription(ctx context.Context, subscription *Subscription) (uuid.UUID, error)
	Subscription(ctx context.Context, subscriptionID uuid.UUID) (*Subscription, error)
	DeleteSubscription(ctx context.Context, subscriptionID uuid.UUID) error
	UpdateSubscription(ctx context.Context, subscription *Subscription) error
	ListSubscription(ctx context.Context, offset, limit int) ([]*Subscription, error)
	TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate MonthYear) ([]Subscription, error)
	Count(ctx context.Context) (int64, error)
}

type AddSubcriptionRequest struct {
	ServiceName  string  `json:"service_name" binding:"required" example:"Yandex Plus"`
	Price        float64 `json:"price" binding:"required" example:"400"`
	UserIDRaw    string  `json:"user_id" binding:"required" example:"a19df875-4040-4fc3-84ad-003d013fcd89"`
	StartDateRaw string  `json:"start_date" binding:"required" example:"07-2025"`
	EndDateRaw   string  `json:"end_date" example:"07-2026"`
}

type UpdateSubcriptionRequest struct {
	SubscriptionIDRaw string `json:"id" binding:"required"`
	ServiceName       string `json:"service_name" binding:"required" example:"Yandex Plus"`
	Price             int    `json:"price" binding:"required" example:"400"`
	UserIDRaw         string `json:"user_id" binding:"required" example:"a19df875-4040-4fc3-84ad-003d013fcd89"`
	StartDateRaw      string `json:"start_date" binding:"required" example:"07-2025"`
	EndDateRaw        string `json:"end_date" example:"07-2026"`
}
