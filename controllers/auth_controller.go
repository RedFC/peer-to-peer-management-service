package controllers

import (
	"p2p-management-service/models"
	"p2p-management-service/services"

	"github.com/kataras/iris/v12"
)

type AuthController struct {
	Service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{Service: service}
}

// Login godoc
// @Summary      User login
// @Description  User login with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login body models.LoginRequest true "Login request"
// @Success      200  {object}  models.SuccessResponse
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Router       /auth/login [post]
func (c *AuthController) Login(ctx iris.Context) {
	var loginRequest models.LoginRequest
	if err := ctx.ReadJSON(&loginRequest); err != nil {
		ctx.StatusCode(400)
		ctx.JSON(models.ErrorResponse{Message: "Invalid request", Trace: err.Error()})
		return
	}

	// Perform login logic here
	loginResponse, err := c.Service.Authenticate(loginRequest.Email, loginRequest.Password)
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(models.ErrorResponse{Message: "Invalid credentials", Trace: err.Error()})
		return
	}

	ctx.JSON(models.SuccessResponse{Data: loginResponse, Message: "Login successful"})
}
