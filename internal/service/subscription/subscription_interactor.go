package subscription

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/immxrtalbeast/subscription-aggregator/internal/domain"
	"github.com/immxrtalbeast/subscription-aggregator/internal/lib/logger/sl"
)

type SubscriptionInteractor struct {
	log      *slog.Logger
	subsRepo domain.SubscriptionRepository
}

func NewSubscriptionInteractor(log *slog.Logger, subsRepo domain.SubscriptionRepository) *SubscriptionInteractor {
	return &SubscriptionInteractor{log: log, subsRepo: subsRepo}
}

func (si *SubscriptionInteractor) AddSubscription(ctx context.Context, serviceName string, price int, userID uuid.UUID, startDate domain.MonthYear, endDate *domain.MonthYear) (uuid.UUID, error) {
	const op = "service.subscription.add"
	log := si.log.With(
		slog.String("op", op),
		slog.String("service_name", serviceName),
		slog.String("userID", userID.String()),
	)
	log.Info("adding subscription")
	subscription := &domain.Subscription{
		ServiceName: serviceName,
		Price:       price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	id, err := si.subsRepo.SaveSubscription(ctx, subscription)
	if err != nil {
		log.Error("failed to save subscription", sl.Err(err))
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("Subscription saved!")
	return id, nil
}

func (si *SubscriptionInteractor) Subscription(ctx context.Context, subscriptionID uuid.UUID) (*domain.Subscription, error) {
	const op = "service.subscription.get"
	log := si.log.With(
		slog.String("op", op),
		slog.String("id", subscriptionID.String()),
	)
	log.Info("getting subscription")
	subscription, err := si.subsRepo.Subscription(ctx, subscriptionID)
	if err != nil {
		log.Error("failed to get subscription", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("subscription provided")
	return subscription, nil
}

func (si *SubscriptionInteractor) DeleteSubscription(ctx context.Context, subscriptionID uuid.UUID) error {
	const op = "service.subscription.delete"
	log := si.log.With(
		slog.String("op", op),
		slog.String("id", subscriptionID.String()),
	)
	log.Info("deleting subscription")
	if err := si.subsRepo.DeleteSubscription(ctx, subscriptionID); err != nil {
		log.Error("failed to delete subscription", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("subscription deleted")
	return nil
}

func (si *SubscriptionInteractor) UpdateSubscription(ctx context.Context, subscriptionID uuid.UUID, serviceName string, price int, userID uuid.UUID, startDate domain.MonthYear, endDate *domain.MonthYear) error {
	const op = "service.subscription.update"
	log := si.log.With(
		slog.String("op", op),
		slog.String("subscription_id", subscriptionID.String()),
		slog.String("service_name", serviceName),
		slog.String("user_id", userID.String()),
	)
	log.Info("updating subscription")
	subscription := &domain.Subscription{
		ID:          subscriptionID,
		ServiceName: serviceName,
		Price:       price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}
	if err := si.subsRepo.UpdateSubscription(ctx, subscription); err != nil {
		log.Error("failed to update subscription")
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (si *SubscriptionInteractor) ListSubscription(ctx context.Context, offset, limit int) ([]*domain.Subscription, int64, error) {
	const op = "service.subscription.list"
	log := si.log.With(
		slog.String("op", op),
	)
	log.Info("getting list of subscriptions")
	list, err := si.subsRepo.ListSubscription(ctx, offset, limit)
	if err != nil {
		log.Error("failed to get list of subscription")
		return nil, 0, fmt.Errorf("%s: %w", op, err)
	}
	total, err := si.subsRepo.Count(ctx)
	if err != nil {
		log.Error("failed to count of subscription")
		return list, 0, err
	}

	log.Info("list provided")
	return list, total, nil
}

func (si *SubscriptionInteractor) TotalCost(ctx context.Context, userID *uuid.UUID, serviceName *string, startDate, endDate domain.MonthYear) (int, error) {
	const op = "service.subscription.totalCost"
	log := si.log.With(
		slog.String("op", op),
		slog.String("start_date", startDate.String()),
		slog.String("end_date", endDate.String()),
	)
	if endDate.IsBefore(startDate) {
		log.Error("start date cannot be after end date")
		return 0, errors.New("start date cannot be after end date")
	}

	subscriptions, err := si.subsRepo.TotalCost(ctx, userID, serviceName, startDate, endDate)
	if err != nil {
		log.Error("failed to get suscriptions", sl.Err(err))
		return 0, err
	}
	total := 0
	for _, sub := range subscriptions {
		months := si.calculateActiveMonths(sub.StartDate, sub.EndDate, startDate, endDate)
		total += sub.Price * months
	}
	return total, nil
}

func (si *SubscriptionInteractor) calculateActiveMonths(subStart domain.MonthYear, subEnd *domain.MonthYear, periodStart, periodEnd domain.MonthYear) int {
	startMonth := domain.MaxMonthYear(subStart, periodStart)

	var endMonth domain.MonthYear
	if subEnd == nil {
		endMonth = periodEnd
	} else {
		endMonth = domain.MinMonthYear(*subEnd, periodEnd)
	}

	if domain.CompareMonthYears(startMonth, endMonth) > 0 {
		return 0
	}

	return domain.MonthDifference(startMonth, endMonth) + 1
}
