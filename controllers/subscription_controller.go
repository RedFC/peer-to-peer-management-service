package controllers

import (
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"

	"p2p-management-service/models"
	"p2p-management-service/services"
)

type SubscriptionController struct {
	Service *services.SubscriptionService
}

func NewSubscriptionController(service *services.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{Service: service}
}

// GetSubscriptions godoc
// @Summary      Get list of subscriptions
// @Description  Get all subscriptions with pagination
// @Tags         subscriptions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        pageSize  query     int  false  "Page size"    default(10)
// @Param        pageNo    query     int  false  "Page number"  default(1)
// @Success      200  {object}  models.PaginatedResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /subscriber [get]
func (c *SubscriptionController) GetSubscriptions(ctx iris.Context) {
	pageSize := ctx.URLParamDefault("pageSize", "10")
	pageNo := ctx.URLParamDefault("pageNo", "1")

	totalCount, err := c.Service.DB.GetTotalSubscriptionsCount(ctx)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(models.ErrorResponse{Message: "Failed to retrieve subscriptions count", Trace: err.Error()})
		return
	}

	subscriptions, err := c.Service.GetSubscriptions(pageNo, pageSize, totalCount)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(models.ErrorResponse{Message: "Failed to retrieve subscriptions", Trace: err.Error()})
		return
	}

	ctx.JSON(models.PaginatedResponse{
		Data:       subscriptions,
		PageSize:   pageSize,
		PageNo:     pageNo,
		TotalCount: totalCount,
		Message:    "Subscriptions retrieved successfully",
	})
}

// GetSubscriptionByID godoc
// @Summary      Get subscription by ID
// @Description  Get a subscription by its ID
// @Tags         subscriptions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Subscription ID"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /subscriber/{id} [get]
func (c *SubscriptionController) GetSubscriptionByID(ctx iris.Context) {
	idParam := ctx.Params().Get("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid subscription ID format", Trace: err.Error()})
		return
	}

	subscription, err := c.Service.GetSubscriptionByID(id)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(models.ErrorResponse{Message: "Subscription not found", Trace: err.Error()})
		return
	}

	ctx.JSON(models.SuccessResponse{Data: subscription, Message: "Subscription retrieved successfully"})
}

// CreateSubscription godoc
// @Summary      Create a new subscription
// @Description  Create a new subscription with the input payload
// @Tags         subscriptions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        subscription  body      models.CreateSubscriptionParams  true  "Create Subscription Request"
// @Success      201   {object}  models.SuccessResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      500   {object}  models.ErrorResponse
// @Router       /subscriber [post]
func (c *SubscriptionController) CreateSubscription(ctx iris.Context) {
	var createSubscriptionRequest models.CreateSubscriptionParams
	if err := ctx.ReadJSON(&createSubscriptionRequest); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid request", Trace: err.Error()})
		return
	}

	newSubscription, err := c.Service.CreateSubscription(createSubscriptionRequest)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(models.ErrorResponse{Message: "Failed to create subscription", Trace: err.Error()})
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(models.SuccessResponse{Data: newSubscription, Message: "Subscription created successfully"})
}

// UpdateSubscription godoc
// @Summary      Update an existing subscription
// @Description  Update an existing subscription by ID
// @Tags         subscriptions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      int                           true  "Subscription ID"
// @Param        subscription  body      models.UpdateSubscriptionParams  true  "Update Subscription Request"
// @Success      200   {object}  models.SuccessResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      404   {object}  models.ErrorResponse
// @Router       /subscriber/{id} [put]
func (c *SubscriptionController) UpdateSubscription(ctx iris.Context) {
	idParam := ctx.Params().Get("id")
	idInt, err := uuid.Parse(idParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid subscription ID format", Trace: err.Error()})
		return
	}

	var updateSubscriptionRequest models.UpdateSubscriptionParams
	if err := ctx.ReadJSON(&updateSubscriptionRequest); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid request", Trace: err.Error()})
		return
	}

	updatedSubscription, err := c.Service.UpdateSubscription(idInt, updateSubscriptionRequest)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(models.ErrorResponse{Message: "Subscription not found", Trace: err.Error()})
		return
	}

	ctx.JSON(models.SuccessResponse{Data: updatedSubscription, Message: "Subscription updated successfully"})
}

// DeleteSubscription godoc
// @Summary      Delete a subscription
// @Description  Delete a subscription by user and group identifiers
// @Tags         subscriptions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        group_id  path      int     true  "Group ID"
// @Param        user_id   path      string  true  "User ID (UUID)"
// @Success      204  "No Content"
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /subscriber/{group_id}/{user_id} [delete]
func (c *SubscriptionController) DeleteSubscription(ctx iris.Context) {
	userIDParam := ctx.Params().Get("user_id")
	groupIDParam := ctx.Params().Get("group_id")

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid user ID format", Trace: err.Error()})
		return
	}

	groupID, err := uuid.Parse(groupIDParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid group ID format", Trace: err.Error()})
		return
	}

	err = c.Service.DeleteSubscription(userID, groupID)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(models.ErrorResponse{Message: "Subscription not found", Trace: err.Error()})
		return
	}

	ctx.StatusCode(iris.StatusNoContent)
}

// RevokeSubscription godoc
// @Summary      Revoke a subscription
// @Description  Revoke an active subscription by user and group identifiers
// @Tags         subscriptions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        group_id  path      int     true  "Group ID"
// @Param        user_id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /subscriber/{group_id}/{user_id}/revoke [post]
func (c *SubscriptionController) RevokeSubscription(ctx iris.Context) {
	userIDParam := ctx.Params().Get("user_id")
	groupIDParam := ctx.Params().Get("group_id")

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid user ID format", Trace: err.Error()})
		return
	}

	groupID, err := uuid.Parse(groupIDParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid group ID format", Trace: err.Error()})
		return
	}

	err = c.Service.RevokeSubscription(userID, groupID)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(models.ErrorResponse{Message: "Subscription not found", Trace: err.Error()})
		return
	}

	ctx.JSON(models.SuccessResponse{Message: "Subscription revoked successfully"})
}

// SendEmailInvitation godoc
// @Summary      Resend invitation email
// @Description  Triggers a new invitation email to the subscriber for the provided group
// @Tags         subscriptions
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        group_id  path      int     true  "Group ID"
// @Param        user_id   path      string  true  "User ID (UUID)"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /subscriber/{group_id}/{user_id} [get]
func (c *SubscriptionController) SendEmailInvitation(ctx iris.Context) {
	userIDParam := ctx.Params().Get("user_id")
	groupIDParam := ctx.Params().Get("group_id")

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid user ID format", Trace: err.Error()})
		return
	}

	groupID, err := uuid.Parse(groupIDParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid group ID format", Trace: err.Error()})
		return
	}

	_, err = c.Service.ResendInvitation(userID, groupID)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(models.ErrorResponse{Message: "Subscription not found", Trace: err.Error()})
		return
	}

	ctx.JSON(models.SuccessResponse{Message: "Email sent successfully"})
}
