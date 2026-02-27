package controllers

import (
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"

	"p2p-management-service/models"
	"p2p-management-service/services"
	"p2p-management-service/utils"
)

type UserController struct {
	Service         *services.UserService
	RelationService *services.RelationService
}

func NewUserController(service *services.UserService, relationService *services.RelationService) *UserController {
	return &UserController{
		Service:         service,
		RelationService: relationService,
	}
}

// GetUsers godoc
// @Summary      Get list of users
// @Description  Get all users
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        pageSize  query     int  false  "Page size"  default(10)
// @Param        pageNo    query     int  false  "Page number"  default(1)
// @Param        roleId    query     string  false  "Role ID"
// @Param        groupId   query     string  false  "Group ID"
// @Success      200  {object}  models.PaginatedResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /users [get]
func (c *UserController) GetUsers(ctx iris.Context) {
	pageSize := ctx.URLParamDefault("pageSize", "10")
	pageNo := ctx.URLParamDefault("pageNo", "1")

	// NEW: Optional Filters
	roleIdParam := ctx.URLParam("roleId")   // ex: ?roleId=00
	groupIdParam := ctx.URLParam("groupId") // ex: ?groupId=10

	roleHashes := []string{utils.Hash(models.RoleSuperAdmin)}
	totalCount, err := c.Service.GetTotalUsersCount(roleHashes, roleIdParam, groupIdParam)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(models.ErrorResponse{Message: "Failed to retrieve users count", Trace: err.Error()})
		return
	}

	users, err := c.Service.GetUsers(pageNo, pageSize, totalCount, roleHashes, roleIdParam, groupIdParam)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(models.ErrorResponse{Message: "Failed to retrieve users", Trace: err.Error()})
		return
	}

	ctx.JSON(models.PaginatedResponse{
		Data:       users,
		TotalCount: totalCount,
		PageNo:     pageNo,
		PageSize:   pageSize,
		Message:    "Users retrieved successfully",
	})
}

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Get a user by their ID
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /users/{id} [get]
func (c *UserController) GetUserByID(ctx iris.Context) {
	idParam := ctx.Params().Get("id")
	uuidValue, err := uuid.Parse(idParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid user ID", Trace: err.Error()})
		return
	}

	user, err := c.Service.GetUserByID(uuidValue)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(models.ErrorResponse{Message: "User not found", Trace: err.Error()})
		return
	}

	ctx.JSON(models.SuccessResponse{Data: user, Message: "User retrieved successfully"})
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user with the input payload
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        user  body      models.CreateUserParams  true  "Create User Request"
// @Success      201   {object}  models.SuccessResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      500   {object}  models.ErrorResponse
// @Router       /users [post]
func (c *UserController) CreateUser(ctx iris.Context) {
	var createUserRequest models.CreateUserParams
	if err := ctx.ReadJSON(&createUserRequest); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid request", Trace: err.Error()})
		return
	}

	newUser, err := c.Service.CreateUser(createUserRequest)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(models.ErrorResponse{Message: "Failed to create user", Trace: err.Error()})
		return
	}

	_, err = c.RelationService.RolesToUsers(models.AttachRoleParams{
		User: *newUser,
		Role: models.RoleUser,
	})
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(models.ErrorResponse{Message: "Failed to assign default role", Trace: err.Error()})
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(models.SuccessResponse{Data: newUser, Message: "User created successfully"})
}

// UpdateUser godoc
// @Summary      Update an existing user
// @Description  Update an existing user by ID
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id    path      string                   true  "User ID"
// @Param        user  body      models.UpdateUserParams  true  "Update User Request"
// @Success      200   {object}  models.SuccessResponse
// @Failure      400   {object}  models.ErrorResponse
// @Failure      404   {object}  models.ErrorResponse
// @Router       /users/{id} [put]
func (c *UserController) UpdateUser(ctx iris.Context) {
	idParam := ctx.Params().Get("id")
	uuidValue, err := uuid.Parse(idParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid user ID", Trace: err.Error()})
		return
	}

	var updateUserRequest models.UpdateUserParams
	if err := ctx.ReadJSON(&updateUserRequest); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid request", Trace: err.Error()})
		return
	}

	updatedUser, err := c.Service.UpdateUser(uuidValue, updateUserRequest)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(models.ErrorResponse{Message: "User not found", Trace: err.Error()})
		return
	}

	ctx.JSON(models.SuccessResponse{Data: updatedUser, Message: "User updated successfully"})
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Delete a user by ID
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      204  "No Content"
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /users/{id} [delete]
func (c *UserController) DeleteUser(ctx iris.Context) {
	idParam := ctx.Params().Get("id")
	uuidValue, err := uuid.Parse(idParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid user ID", Trace: err.Error()})
		return
	}

	err = c.Service.DeleteUser(uuidValue)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(models.ErrorResponse{Message: "User not found", Trace: err.Error()})
		return
	}

	ctx.StatusCode(iris.StatusNoContent)
}

// AttachRole godoc
// @Summary      Attach role to a user
// @Description  Attach a role to a user by user and role identifiers
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        userId  path      string  true  "User ID"
// @Param        roleId  path      int     true  "Role ID"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /users/{userId}/role/{roleId} [post]
func (c *UserController) AttachRole(ctx iris.Context) {
	idParam := ctx.Params().Get("userId")
	uuidValue, err := uuid.Parse(idParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid user ID", Trace: err.Error()})
		return
	}

	roleIdParam := ctx.Params().Get("roleId")
	roleId, err := uuid.Parse(roleIdParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid role ID format", Trace: err.Error()})
		return
	}

	userData, err := c.Service.GetUserByID(uuidValue)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(models.ErrorResponse{Message: "User not found", Trace: err.Error()})
		return
	}

	_, err = c.RelationService.RolesToUsersByID(userData.ID, roleId)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(models.ErrorResponse{Message: "Failed to attach role", Trace: err.Error()})
		return
	}

	ctx.JSON(models.SuccessResponse{Message: "Role attached"})
}

// AttachGroup godoc
// @Summary      Attach group to a user
// @Description  Attach a group to a user by user and group identifiers
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        userId   path      string  true  "User ID"
// @Param        groupId  path      int     true  "Group ID"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Failure      500  {object}  models.ErrorResponse
// @Router       /users/{userId}/group/{groupId} [post]
func (c *UserController) AttachGroup(ctx iris.Context) {
	idParam := ctx.Params().Get("userId")
	uuidValue, err := uuid.Parse(idParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid user ID", Trace: err.Error()})
		return
	}

	groupIdParam := ctx.Params().Get("groupId")
	groupId, err := uuid.Parse(groupIdParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(models.ErrorResponse{Message: "Invalid group ID format", Trace: err.Error()})
		return
	}

	userData, err := c.Service.GetUserByID(uuidValue)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(models.ErrorResponse{Message: "User not found", Trace: err.Error()})
		return
	}

	_, err = c.RelationService.GroupToUsersByID(userData.ID, groupId)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(models.ErrorResponse{Message: "Failed to attach group", Trace: err.Error()})
		return
	}

	ctx.JSON(models.SuccessResponse{Message: "Group attached"})
}
