package routes

import (
	"p2p-management-service/config"
	"p2p-management-service/controllers"
	"p2p-management-service/db/generated"
	"p2p-management-service/middlewares"
	"p2p-management-service/services"
	"p2p-management-service/utils"

	"github.com/kataras/iris/v12"
)

func SetupRoutes(app *iris.Application, DB *generated.Queries, encryptor *utils.Encryptor) {

	api := app.Party("/api/v1")
	{

		api.Get("/health-check", func(ctx iris.Context) {
			ctx.JSON(iris.Map{"status": "ok"})
		})

		auth := api.Party("/auth")
		{
			authService := services.NewAuthService(DB, encryptor)
			authController := controllers.NewAuthController(authService)
			auth.Post("/login", authController.Login)
		}

		users := api.Party("/users", middlewares.AuthMiddleware, middlewares.SuperAdminMiddleware)
		{
			userService := services.NewUserService(DB, encryptor)
			relationService := services.NewRelationService(DB, encryptor)
			userController := controllers.NewUserController(userService, relationService)

			users.Get("/", userController.GetUsers)
			users.Post("/", userController.CreateUser)
			users.Get("/{id:string}", userController.GetUserByID)
			users.Put("/{id:string}", userController.UpdateUser)
			users.Delete("/{id:string}", userController.DeleteUser)
			users.Post("/{userId}/role/{roleId}", userController.AttachRole)
			users.Post("/{userId}/group/{groupId}", userController.AttachGroup)
		}

		roles := api.Party("/roles", middlewares.AuthMiddleware, middlewares.SuperAdminMiddleware)
		{
			roleService := services.NewRoleService(DB, encryptor)
			roleController := controllers.NewRoleController(roleService)

			roles.Get("/", roleController.GetRoles)
			roles.Post("/", roleController.CreateRole)
			roles.Get("/{id:string}", roleController.GetRoleByID)
			roles.Put("/{id:string}", roleController.UpdateRole)
			roles.Delete("/{id:string}", roleController.DeleteRole)
		}

		groups := api.Party("/groups", middlewares.AuthMiddleware, middlewares.SuperAdminMiddleware)
		{
			groupService := services.NewGroupService(DB, encryptor)
			groupController := controllers.NewGroupController(groupService)

			groups.Get("/", groupController.GetGroups)
			groups.Post("/", groupController.CreateGroup)
			groups.Get("/{id:string}", groupController.GetGroupByID)
			groups.Put("/{id:string}", groupController.UpdateGroup)
			groups.Delete("/{id:string}", groupController.DeleteGroup)
		}

		subscribers := api.Party("/subscriber", middlewares.AuthMiddleware, middlewares.SuperAdminMiddleware)
		{
			// choose email backend based on configuration
			var emailSvc services.EmailService
			switch config.AppConfig.EMAIL_BACKEND {
			case "local":
				emailSvc = services.NewLocalEmailService()
			case "ses":
				emailSvc = services.NewSESEmailService()
			case "postfix":
				emailSvc = services.NewPostfixEmailService()
			default:
				// default to nil (no emails)
				emailSvc = nil
			}

			subscriberService := services.NewSubscriptionService(DB, encryptor, emailSvc)
			subscriberController := controllers.NewSubscriptionController(subscriberService)

			subscribers.Get("/", subscriberController.GetSubscriptions)
			subscribers.Post("/", subscriberController.CreateSubscription)
			subscribers.Get("/{id:uuid}", subscriberController.GetSubscriptionByID)
			subscribers.Put("/{id:uuid}", subscriberController.UpdateSubscription)
			subscribers.Delete("/{group_id:uuid}/{user_id:uuid}", subscriberController.DeleteSubscription)
			subscribers.Post("/{group_id:uuid}/{user_id:uuid}/revoke", subscriberController.RevokeSubscription)
			subscribers.Get("/{group_id:uuid}/{user_id:uuid}/resend-email", subscriberController.SendEmailInvitation)
		}
	}
}
