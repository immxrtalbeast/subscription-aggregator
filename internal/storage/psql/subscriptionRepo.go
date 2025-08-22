package psql

import (
	"context"

	"github.com/google/uuid"
	"github.com/immxrtalbeast/subscription-aggregator/internal/domain"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) SaveSubscription(ctx context.Context, subscription *domain.Subscription) (uuid.UUID, error) {
	result := r.db.WithContext(ctx).Create(&subscription)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}
	return subscription.ID, nil
}

func (r *SubscriptionRepository) Subscription(ctx context.Context, subscriptionID uuid.UUID) (*domain.Subscription, error) {
	var subscription *domain.Subscription
	err := r.db.Where("id = ?", subscriptionID).First(&subscription).Error
	return subscription, err
}

func (r *SubscriptionRepository) DeleteSubscription(ctx context.Context, subscriptionID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", subscriptionID).Delete(&domain.Subscription{}).Error
}

func (r *SubscriptionRepository) UpdateSubscription(ctx context.Context, subscription *domain.Subscription) error {
	result := r.db.WithContext(ctx).Model(&domain.Subscription{}).
		Where("id = ?", subscription.ID).
		Omit("id").
		Updates(&subscription)

	return result.Error
}
func (r *SubscriptionRepository) ListSubscription(ctx context.Context) ([]*domain.Subscription, error) {
	var subscriptions []*domain.Subscription
	err := r.db.WithContext(ctx).Model(&domain.Subscription{}).Scan(&subscriptions).Error
	return subscriptions, err
}
