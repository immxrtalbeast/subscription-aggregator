package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/immxrtalbeast/subscription-aggregator/internal/domain"
	"github.com/immxrtalbeast/subscription-aggregator/internal/storage/psql"
)

type SubscriptionController struct {
	subscriptionService domain.SubscriptionInteractor
}

func NewSubscriptionController(subscriptionService domain.SubscriptionInteractor) *SubscriptionController {
	return &SubscriptionController{subscriptionService: subscriptionService}
}

// @Summary Создать подписку
// @Param   subscription body domain.AddSubcriptionRequest true "Данные подписки"
// @Success 200 {object} map[string]interface{}
// @Router /create [post]
func (c *SubscriptionController) AddSubcription(ctx *gin.Context) {

	var req domain.AddSubcriptionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	if req.Price <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "price should be > 0",
		})
		return
	}
	startDate, err := domain.ParseMonthYear(req.StartDateRaw)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid start date format",
			"details": err.Error(),
		})
		return
	}
	var endDate *domain.MonthYear
	if req.EndDateRaw != "" {
		parsed, err := domain.ParseMonthYear(req.EndDateRaw)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid start date format",
				"details": err.Error(),
			})
			return
		}
		endDate = &parsed
		if endDate.IsBefore(*&startDate) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "start date should be after end date",
			})
			return
		}
	}
	userID, err := uuid.Parse(req.UserIDRaw)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid user_id",
			"details": err.Error(),
		})
		return
	}
	subscriptionID, err := c.subscriptionService.AddSubscription(ctx, req.ServiceName, req.Price, userID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create subscription",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":         "subscription added successfuly",
		"subscription_id": subscriptionID,
	})
}

// @Summary Получить подписку
// @Param   id path int true "ID подписки"
// @Success 200 {object} domain.Subscription
// @Router /{id} [get]
func (c *SubscriptionController) Subscription(ctx *gin.Context) {
	subscriptionIDRaw := ctx.Param("id")
	subscriptionID, err := uuid.Parse(subscriptionIDRaw)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "couldn`t parse uuid",
			"details": err.Error(),
		})
		return
	}
	subscription, err := c.subscriptionService.Subscription(ctx, subscriptionID)
	if err != nil {
		if errors.Is(err, psql.ErrSubscriptNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": psql.ErrSubscriptNotFound.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to get subscription",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"subscription": subscription,
	})
}

// @Summary Удалить подписку
// @Param   id path int true "ID подписки"
// @Success 200
// @Router /{id} [delete]
func (c *SubscriptionController) DeleteSubscription(ctx *gin.Context) {
	subscriptionIDRaw := ctx.Param("id")
	subscriptionID, err := uuid.Parse(subscriptionIDRaw)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "couldn`t parse uuid",
			"details": err.Error(),
		})
		return
	}
	if err := c.subscriptionService.DeleteSubscription(ctx, subscriptionID); err != nil {
		if errors.Is(err, psql.ErrSubscriptNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": psql.ErrSubscriptNotFound,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to delete subscription",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "subscription deleted successfully",
	})
}

// @Summary Изменить подписку
// @Param   subscription body domain.UpdateSubcriptionRequest true "Данные"
// @Success 200 {object} map[string]interface{}
// @Router /update [put]
func (c *SubscriptionController) UpdateSubscription(ctx *gin.Context) {
	type UpdateSubcriptionRequest struct {
		SubscriptionIDRaw string `json:"id" binding:"required"`
		ServiceName       string `json:"service_name" binding:"required" example:"Yandex Plus"`
		Price             int    `json:"price" binding:"required" example:"400"`
		UserIDRaw         string `json:"user_id" binding:"required" example:"a19df875-4040-4fc3-84ad-003d013fcd89"`
		StartDateRaw      string `json:"start_date" binding:"required" example:"07-2025"`
		EndDateRaw        string `json:"end_date" example:"07-2026"`
	}
	var req UpdateSubcriptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	startDate, err := domain.ParseMonthYear(req.StartDateRaw)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid start date format",
			"details": err.Error(),
		})
		return
	}
	var endDate *domain.MonthYear
	if req.EndDateRaw != "" {
		parsed, err := domain.ParseMonthYear(req.EndDateRaw)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "invalid start date format",
				"details": err.Error(),
			})
			return
		}
		endDate = &parsed
	}
	userID, err := uuid.Parse(req.UserIDRaw)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid user_id",
			"details": err.Error(),
		})
		return
	}
	subscriptionID, err := uuid.Parse(req.SubscriptionIDRaw)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid subscription_id",
			"details": err.Error(),
		})
		return
	}
	if err = c.subscriptionService.UpdateSubscription(ctx, subscriptionID, req.ServiceName, req.Price, userID, startDate, endDate); err != nil {
		if errors.Is(err, psql.ErrSubscriptNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": psql.ErrSubscriptNotFound,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to create subscription",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "subscription updated successfuly",
	})
}

// @Summary Получить все подписки
// @Param param_name query string false "Описание параметра"
// @Success 200 {object} map[string]interface{}
// @Router /all [get]
func (c *SubscriptionController) ListSubscription(ctx *gin.Context) {
	subscriptions, err := c.subscriptionService.ListSubscription(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to get list of subscriptions",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"subscriptions": subscriptions,
	})
}

// @Summary Подсчет суммарной стоимости всех подписок за выбранный период с фильтрацией по id пользователя и названию подписки
// @Param   user_id      query string  false "ID пользователя"
// @Param   service_name query string  false "Название сервиса"
// @Param   start_date   query string  true  "Начальная дата (YYYY-MM-DD)"
// @Param   end_date     query string  true  "Конечная дата (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Router /total [get]
func (c *SubscriptionController) TotalCost(ctx *gin.Context) {
	var req struct {
		UserID      *string `form:"user_id"`
		ServiceName *string `form:"service_name"`
		StartDate   string  `form:"start_date" binding:"required"`
		EndDate     string  `form:"end_date" binding:"required"`
	}
	if err := ctx.BindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var userID *uuid.UUID
	if req.UserID != nil {
		id, err := uuid.Parse(*req.UserID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "failed to parse userID",
				"details": err.Error(),
			})
			return
		}
		userID = &id
	}
	startDate, err := domain.ParseMonthYear(req.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to parse start date",
			"details": err.Error(),
		})
		return
	}
	endDate, err := domain.ParseMonthYear(req.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to parse end date",
			"details": err.Error(),
		})
		return
	}
	sum, err := c.subscriptionService.TotalCost(ctx, userID, req.ServiceName, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to get cost",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"total_sum": sum,
	})
}
