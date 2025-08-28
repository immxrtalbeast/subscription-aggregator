package psql

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/immxrtalbeast/subscription-aggregator/internal/domain"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

var ErrSubscriptNotFound = errors.New("Subscript not found")

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) SaveSubscription(ctx context.Context, subscription *domain.Subscription) (uuid.UUID, error) {
	result := r.db.WithContext(ctx).Create(&subscription)
	return subscription.ID, result.Error
}

func (r *SubscriptionRepository) Subscription(ctx context.Context, subscriptionID uuid.UUID) (*domain.Subscription, error) {
	var subscription *domain.Subscription
	err := r.db.Where("id = ?", subscriptionID).First(&subscription).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSubscriptNotFound
	}
	return subscription, err
}

func (r *SubscriptionRepository) DeleteSubscription(ctx context.Context, subscriptionID uuid.UUID) error {
	err := r.db.WithContext(ctx).Where("id = ?", subscriptionID).Delete(&domain.Subscription{}).Error
	if errors.Is(err, ErrSubscriptNotFound) {
		return ErrSubscriptNotFound
	}
	return err
}

func (r *SubscriptionRepository) UpdateSubscription(ctx context.Context, subscription *domain.Subscription) error {
	result := r.db.WithContext(ctx).Model(&domain.Subscription{}).
		Where("id = ?", subscription.ID).
		Omit("id").
		Updates(&subscription)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ErrSubscriptNotFound
	}
	return result.Error
}
func (r *SubscriptionRepository) ListSubscription(ctx context.Context, offset, limit int) ([]*domain.Subscription, error) {
	var subscriptions []*domain.Subscription
	err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Model(&domain.Subscription{}).Scan(&subscriptions).Error
	return subscriptions, err
}

func (r *SubscriptionRepository) TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate domain.MonthYear) ([]domain.Subscription, error) {
	var subscriptions []domain.Subscription

	query := r.db.WithContext(ctx).Model(&domain.Subscription{}).
		Where("start_date <= ?", endDate).
		Where("(end_date IS NULL OR end_date >= ?)", startDate)

	if userID != nil {
		query = query.Where("user_id = ?", userID)
	}
	if serviceName != nil {
		query = query.Where("service_name = ?", serviceName)
	}
	err := query.Scan(&subscriptions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to calculate total cost: %w", err)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	result := r.db.Model(&domain.Subscription{}).Count(&count)
	return count, result.Error
}
