package router

import (
	"database/sql"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/iwanjunaid/basesvc/adapter/controller"
	_ "github.com/iwanjunaid/basesvc/docs"
	"github.com/iwanjunaid/basesvc/shared/logger"
	"os"
	"os/signal"

	"github.com/iwanjunaid/basesvc/registry"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Rest struct {
	port   string
	db     *sql.DB
	server *fiber.App
	controller *controller.AppController
}

func NewRest(port string, db *sql.DB) *Rest {
	r := &Rest{
		db:   db,
		port: port,
	}
	return r
}

func (r *Rest) Serve() {
	r.setup()
	if err := r.server.Listen(fmt.Sprintf(":%s", r.port)); err != nil {
		logger.WithFields(logger.Fields{"component": "api_command"}).Fatalf("failed running rest server, error: %s", err)
	}
}

func (r *Rest) setup() {
	r.initServer()
	r.initRouter()
}

// InitServer: initialize server instance
func (r *Rest) initServer(){
	app := fiber.New()
	app.Use(cors.New())
	app.Use("/swagger", swagger.Handler)
	app.Use(recover.New())

	// adding graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		logger.WithFields(logger.Fields{"component": "api_command"}).Infof("gracefully shutting down service...")
		_ = app.Shutdown()
	}()

	registry := registry.NewRegistry(r.db)
	r.controller = registry.NewAppController()
	r.server = app
}

// @title BaseSVC API
// @version 1.0
// @description This is a sample basesvc server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1
// InitRouter: initialize server router
func (r *Rest) initRouter() {
	v1 := r.server.Group("/v1")
	author := v1.Group("/author")
	author.Get("/", func(ctx *fiber.Ctx) error {
		return r.controller.Author.GetAuthors(ctx)
	})
	r.server.Get("/:id", func(ctx *fiber.Ctx) error {
		return r.controller.Author.GetAuthors(ctx)
	})
	r.server.Post("/", func(ctx *fiber.Ctx) error {
		return r.controller.Author.GetAuthors(ctx)
	})
	r.server.Put("/:id", func(ctx *fiber.Ctx) error {
		return r.controller.Author.GetAuthors(ctx)
	})
}
