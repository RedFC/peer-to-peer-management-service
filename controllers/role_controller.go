package controllers

import (
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"

	"p2p-management-service/models"
	"p2p-management-service/services"
)

type RoleController struct {
	Service *services.RoleService
}

func NewRoleController(service *services.RoleService) *RoleController {
	return &RoleController{Service: service}
}

// GetRoles godoc
// @Summary      Get list of roles
// @Description  Get all roles
// @Tags         roles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.SuccessResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /roles [get]
func (c *RoleController) GetRoles(ctx iris.Context) {
	roles, err := c.Service.GetRoles()
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(models.ErrorResponse{Message: "Failed to retrieve roles", Trace: err.Error()})
		return
	}
	ctx.JSON(models.SuccessResponse{Data: roles, Message: "Roles retrieved successfully"})
}

// GetRoleByID godoc
// @Summary      Get role by ID
// @Description  Get a role by its ID
// @Tags         roles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Role ID"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /roles/{id} [get]
func (c *RoleController) GetRoleByID(ctx iris.Context) {
	id, err := uuid.Parse(ctx.Params().Get("id"))
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid role ID", Trace: err.Error()})
		return
	}

	role, err := c.Service.GetRoleByID(id)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(models.ErrorResponse{Message: "Role not found", Trace: err.Error()})
		return
	}

	ctx.JSON(models.SuccessResponse{Data: role, Message: "Role retrieved successfully"})
}

// CreateRole godoc
// @Summary      Create a new role
// @Description  Create a new role with name and description
// @Tags         roles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        role body models.CreateRoleParams true "Role request"
// @Success      201  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /roles [post]
func (c *RoleController) CreateRole(ctx iris.Context) {
	var roleRequest models.CreateRoleParams
	if err := ctx.ReadJSON(&roleRequest); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid request", Trace: err.Error()})
		return
	}

	role, err := c.Service.CreateRole(roleRequest)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(models.ErrorResponse{Message: "Failed to create role", Trace: err.Error()})
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(models.SuccessResponse{Data: role, Message: "Role created successfully"})
}

// UpdateRole godoc
// @Summary      Update an existing role
// @Description  Update a role by its ID with new name and description
// @Tags         roles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path  int  true  "Role ID"
// @Param        role body models.UpdateRoleParams true "Role request"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /roles/{id} [put]
func (c *RoleController) UpdateRole(ctx iris.Context) {
	id := ctx.Params().Get("id")
	roleId, err := uuid.Parse(id)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid role ID", Trace: err.Error()})
		return
	}

	var roleRequest models.UpdateRoleParams
	if err := ctx.ReadJSON(&roleRequest); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid request", Trace: err.Error()})
		return
	}

	role, err := c.Service.UpdateRole(roleId, roleRequest)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(models.ErrorResponse{Message: "Role not found", Trace: err.Error()})
		return
	}

	ctx.JSON(models.SuccessResponse{Data: role, Message: "Role updated successfully"})
}

// DeleteRole godoc
// @Summary      Delete a role
// @Description  Delete a role by its ID
// @Tags         roles
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Role ID"
// @Success      204
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /roles/{id} [delete]
func (c *RoleController) DeleteRole(ctx iris.Context) {
	id := ctx.Params().Get("id")

	roleId, err := uuid.Parse(id)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid role ID", Trace: err.Error()})
		return
	}

	err = c.Service.DeleteRole(roleId)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(models.ErrorResponse{Message: "Role not found", Trace: err.Error()})
		return
	}

	ctx.StatusCode(iris.StatusNoContent)
}
