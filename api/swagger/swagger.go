package swagger

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	docs "github.com/tradersclub/TCTemplateBack/docs"
	"github.com/tradersclub/TCUtils/logger"
)

// Options struct de opções para a criação de uma instancia do swagger
type Options struct {
	Group     *echo.Group
	AccessKey string
	Port      string
}

// Register group item check
func Register(opts Options) {

	docs.SwaggerInfo.Title = "Swagger {Nome do projeto} API"
	docs.SwaggerInfo.Description = "Swagger com as rotas e modelos de uso de parâmetros da API do {Nome do Projeto}"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost" + opts.Port // Host of application
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	opts.Group.GET("/:key", func(c echo.Context) error {
		key := c.Param("key")
		if key != opts.AccessKey {
			return nil
		}
		return c.Redirect(http.StatusFound, "/swagger/"+key+"/index.html")
	})

	opts.Group.GET("/:key/*", func(c echo.Context) error {
		key := c.Param("key")

		if key != opts.AccessKey {
			return c.JSON(
				http.StatusUnauthorized,
				notAuthorized{http.StatusUnauthorized, "You are not wellcome here"},
			)
		}

		return echoSwagger.WrapHandler(c)
	})

	logger.Info("Swagger is initializing...")
}
