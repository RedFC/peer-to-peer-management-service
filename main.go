package main

import (
	"fmt"
	"log"
	"strings"

	"p2p-management-service/config"
	"p2p-management-service/db"
	"p2p-management-service/db/generated"
	_ "p2p-management-service/docs"
	"p2p-management-service/middlewares"
	"p2p-management-service/routes"
	"p2p-management-service/scripts"
	"p2p-management-service/utils"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           P2p Management Services API
// @version         1.0
// @description     API for P2p Management Services
// @host            localhost:8080
// @host            p2pmanagement.net
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	// init logger
	utils.InitLogger()

	// Create a new Iris application
	app := iris.New()

	// 🔹 Add global request/error logging middleware
	app.Use(middlewares.LoggerMiddleware)

	// 🔹 Panic recovery logging
	app.Use(middlewares.RecoveryMiddleware)

	// Load environment variables
	config.LoadConfig(".env")

	// Initialize the encryptor with a secret key from env
	secretKey := config.AppConfig.ENCRYPTION_SECRET
	if len(secretKey) != 32 {
		log.Fatal("ENCRYPTION_SECRET must be 32 bytes for AES-256")
	}
	encryptor := utils.NewEncryptor(secretKey)

	corsConfig := config.AppConfig.CORS
	crs := cors.New(cors.Options{
		AllowedOrigins:   splitAndTrim(corsConfig.AllowedOrigins),
		AllowedMethods:   splitAndTrim(corsConfig.AllowedMethods),
		AllowedHeaders:   splitAndTrim(corsConfig.AllowedHeaders),
		AllowCredentials: true,
		ExposedHeaders:   splitAndTrim(corsConfig.ExposedHeaders),
	})

	app.UseRouter(crs)

	// Connect to the database
	db.ConnectDatabase()
	queries := generated.New(db.Conn)

	// seed database
	scripts.Seed(db.Conn)

	// ✅ Convert net/http Swagger handler to iris.Handler
	app.Get("/swagger/{any:path}", iris.FromStd(httpSwagger.WrapHandler))

	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("Running on port %d",
			config.AppConfig.Port)
	})

	routes.SetupRoutes(app, queries, encryptor)

	addr := fmt.Sprintf(":%d", config.AppConfig.Port)
	app.Listen(addr)
}

func splitAndTrim(value string) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
