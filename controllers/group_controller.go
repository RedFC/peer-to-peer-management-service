package controllers

import (
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"

	"p2p-management-service/models"
	"p2p-management-service/services"
)

type GroupController struct {
	Service *services.GroupService
}

func NewGroupController(service *services.GroupService) *GroupController {
	return &GroupController{Service: service}
}

// GetGroups godoc
// @Summary      Get list of groups
// @Description  Get all groups
// @Tags         groups
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        pageSize  query     int  false  "Page size"    default(10)
// @Param        pageNo    query     int  false  "Page number"  default(1)
// @Success      200  {object}  models.PaginatedResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /groups [get]
func (c *GroupController) GetGroups(ctx iris.Context) {
	pageSize := ctx.URLParamDefault("pageSize", "10")
	pageNo := ctx.URLParamDefault("pageNo", "1")
	totalCount, err := c.Service.GetTotalGroupsCount()
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(models.ErrorResponse{Message: "Failed to retrieve groups count", Trace: err.Error()})
		return
	}
	groups, err := c.Service.GetGroups(pageNo, pageSize)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(models.ErrorResponse{Message: "Failed to retrieve groups", Trace: err.Error()})
		return
	}
	ctx.JSON(models.PaginatedResponse{
		Data:       groups,
		TotalCount: totalCount,
		PageNo:     pageNo,
		PageSize:   pageSize,
		Message:    "Groups retrieved successfully",
	})
}

// GetGroupByID godoc
// @Summary      Get group by ID
// @Description  Get a group by its ID
// @Tags         groups
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Group ID"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /groups/{id} [get]
func (c *GroupController) GetGroupByID(ctx iris.Context) {
	id, err := uuid.Parse(ctx.Params().Get("id"))
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Failed to retrieve group ID", Trace: err.Error()})
		return
	}

	group, err := c.Service.GetGroupByID(id)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(models.ErrorResponse{Message: "Group not found", Trace: err.Error()})
		return
	}

	ctx.JSON(models.SuccessResponse{
		Data:    group,
		Message: "Group retrieved successfully",
	})
}

// CreateGroup godoc
// @Summary      Create a new group
// @Description  Create a new group with the provided details
// @Tags         groups
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        group  body  models.CreateGroupParams  true  "Group details"
// @Success      201  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /groups [post]
func (c *GroupController) CreateGroup(ctx iris.Context) {
	var groupRequest models.CreateGroupParams
	if err := ctx.ReadJSON(&groupRequest); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid request", Trace: err.Error()})
		return
	}

	group, err := c.Service.CreateGroup(groupRequest)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(models.ErrorResponse{Message: "Failed to create group", Trace: err.Error()})
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(models.SuccessResponse{Data: group, Message: "Group created successfully"})
}

// UpdateGroup godoc
// @Summary      Update an existing group
// @Description  Update a group with the provided details
// @Tags         groups
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Group ID"
// @Param        group  body  models.UpdateGroupParams  true  "Group details"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /groups/{id} [put]
func (c *GroupController) UpdateGroup(ctx iris.Context) {
	id := ctx.Params().Get("id")
	groupId, err := uuid.Parse(id)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid group ID", Trace: err.Error()})
		return
	}

	var groupRequest models.UpdateGroupParams
	if err := ctx.ReadJSON(&groupRequest); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid request", Trace: err.Error()})
		return
	}

	group, err := c.Service.UpdateGroup(groupId, groupRequest)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(models.ErrorResponse{Message: "Group not found", Trace: err.Error()})
		return
	}

	ctx.JSON(models.SuccessResponse{Data: group, Message: "Group updated successfully"})
}

// DeleteGroup godoc
// @Summary      Delete a group
// @Description  Delete a group by its ID
// @Tags         groups
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Group ID"
// @Success      204  "No Content"
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /groups/{id} [delete]
func (c *GroupController) DeleteGroup(ctx iris.Context) {
	id := ctx.Params().Get("id")
	groupId, err := uuid.Parse(id)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid group ID", Trace: err.Error()})
		return
	}

	err = c.Service.DeleteGroup(groupId)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(models.ErrorResponse{Message: "Group not found", Trace: err.Error()})
		return
	}

	ctx.StatusCode(iris.StatusNoContent)
}
