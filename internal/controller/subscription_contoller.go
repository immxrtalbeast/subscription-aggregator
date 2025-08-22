package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/immxrtalbeast/subscription-aggregator/internal/domain"
)

type SubscriptionController struct {
	subscriptionService domain.SubscriptionInteractor
}

func NewSubscriptionController(subscriptionService domain.SubscriptionInteractor) *SubscriptionController {
	return &SubscriptionController{subscriptionService: subscriptionService}
}

func (c *SubscriptionController) AddSubcription(ctx *gin.Context) {
	type AddSubcriptionRequest struct {
		ServiceName  string `json:"service_name" binding:"required"`
		Price        int    `json:"price" binding:"required"`
		UserIDRaw    string `json:"user_id" binding:"required"`
		StartDateRaw string `json:"start_date" binding:"required"`
		EndDateRaw   string `json:"end_date"`
	}
	var req AddSubcriptionRequest

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
	subscriptionID, err := c.subscriptionService.AddSubscription(ctx, req.ServiceName, req.Price, userID, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to get subscription",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"subscription": subscription,
	})
}

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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to delete subscription",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "subscription deleted successfully",
	})
}

func (c *SubscriptionController) UpdateSubscription(ctx *gin.Context) {
	type UpdateSubcriptionRequest struct {
		SubscriptionIDRaw string `json:"id" binding:"required"`
		ServiceName       string `json:"service_name" binding:"required"`
		Price             int    `json:"price" binding:"required"`
		UserIDRaw         string `json:"user_id" binding:"required"`
		StartDateRaw      string `json:"start_date" binding:"required"`
		EndDateRaw        string `json:"end_date"`
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to create subscription",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "subscription updated successfuly",
	})
}

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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "failed to get cost",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"total_sum": sum,
	})
}
