package domain

import (
	"context"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ServiceName string    `gorm:"not null"`
	Price       int       `gorm:"not null"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	StartDate   MonthYear `gorm:"not null"`
	EndDate     *MonthYear
}

type SubscriptionInteractor interface {
	AddSubscription(ctx context.Context, serviceName string, price int, userID uuid.UUID, startDate MonthYear, endDate *MonthYear) (uuid.UUID, error)
	Subscription(ctx context.Context, subscriptionID uuid.UUID) (*Subscription, error)
	DeleteSubscription(ctx context.Context, subscriptionID uuid.UUID) error
	UpdateSubscription(ctx context.Context, subscriptionID uuid.UUID, serviceName string, price int, userID uuid.UUID, startDate MonthYear, endDate *MonthYear) error
	ListSubscription(ctx context.Context) ([]*Subscription, error)
	// TotalCost(ctx context.Context, userID uuid.UUID, serviceName string, startDate, endDate MonthYear) (int, error)
}

type SubscriptionRepository interface {
	SaveSubscription(ctx context.Context, subscription *Subscription) (uuid.UUID, error)
	Subscription(ctx context.Context, subscriptionID uuid.UUID) (*Subscription, error)
	DeleteSubscription(ctx context.Context, subscriptionID uuid.UUID) error
	UpdateSubscription(ctx context.Context, subscription *Subscription) error
	ListSubscription(ctx context.Context) ([]*Subscription, error)
	// TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate MonthYear) (int, error)
}
