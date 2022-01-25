package server

import (
	"context"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tradersclub/TCUtils/cache"
	"github.com/tradersclub/TCUtils/logger"
	"github.com/tradersclub/TCUtils/mongodb/v2"
	"github.com/tradersclub/TCUtils/natstan"
	"github.com/tradersclub/TCUtils/validator"

	"github.com/labstack/echo-contrib/prometheus"
	emiddleware "github.com/labstack/echo/v4/middleware"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/tradersclub/TCTemplateBack/api"
	"github.com/tradersclub/TCTemplateBack/api/swagger"
	"github.com/tradersclub/TCTemplateBack/app"
	"github.com/tradersclub/TCTemplateBack/config"
	"github.com/tradersclub/TCTemplateBack/event"
	"github.com/tradersclub/TCTemplateBack/store"

	auth "github.com/tradersclub/TCAuth/middleware/echo"
)

// Server is a interface to define contract to server up
type Server interface {
	Start()
	Stop()
	ReloadConnections()
}

type server struct {
	startedAt time.Time

	echo *echo.Echo

	session auth.Middleware

	queueMessager natstan.QueueMessager

	dbNoSQLReader mongodb.Database
	dbNoSQLWriter mongodb.Database

	dbSQLReader *sqlx.DB
	dbSQLWriter *sqlx.DB

	cache cache.Cache

	store *store.Container
	app   *app.Container
}

// New is instance the server
func New() Server {
	return &server{
		startedAt: time.Now(),
	}
}

func (s *server) Start() {

	// ---- setup echo ----
	s.echo = echo.New()
	s.echo.Validator = validator.New()
	s.echo.Debug = config.ConfigGlobal.ENV != "prod"
	s.echo.HideBanner = true
	s.echo.HTTPErrorHandler = createHTTPErrorHandler()

	// ---- setup middlewares ----
	s.echo.Use(emiddleware.Logger())
	s.echo.Use(emiddleware.BodyLimit("2M"))
	s.echo.Use(emiddleware.Recover())
	s.echo.Use(emiddleware.RequestID())

	// ---- setup prometheus ----
	prometheus.NewPrometheus("TCTemplateBack", nil).Use(s.echo)

	// ---- start SQL connection ----
	s.dbSQLWriter = createSqlConnection(config.ConfigGlobal.Database.Writer.URL)
	s.dbSQLReader = createSqlConnection(config.ConfigGlobal.Database.Reader.URL)

	// ---- start NoSQL connection ----
	s.dbNoSQLWriter = createNoSqlConnection(config.ConfigGlobal.Database.WriterNoSql.URL, config.ConfigGlobal.Database.WriterNoSql.Scheme, false)
	s.dbNoSQLReader = createNoSqlConnection(config.ConfigGlobal.Database.ReaderNoSql.URL, config.ConfigGlobal.Database.ReaderNoSql.Scheme, true)

	// ---- start NATS connection ----
	s.queueMessager = natstan.NewQueueMessager(natstan.Options{
		URL: config.ConfigGlobal.Nats.URL,
	})

	// ---- start Memcached connection ----
	s.cache = cache.NewMemcache(config.ConfigGlobal.Cache)

	// ---- start TCAuth connection ----
	s.session = auth.NewMiddle(config.ConfigGlobal.Auth)

	// ---- setup Store ----
	s.store = store.New(store.Options{
		ReaderSQL:   s.dbSQLReader,
		WriterSQL:   s.dbSQLWriter,
		WriterNoSQL: s.dbNoSQLWriter,
		ReaderNoSQL: s.dbNoSQLReader,
	})

	// ---- setup App ----
	s.app = app.New(app.Options{
		Stores:        s.store,
		Cache:         s.cache,
		QueueMessager: s.queueMessager,
		Session:       s.session,
	})

	// ---- setup Event ----
	event.Register(event.Options{
		Apps:          s.app,
		QueueMessager: s.queueMessager,
	})

	// ---- setup Api ----
	api.Register(api.Options{
		Group: s.echo.Group(""),
		Apps:  s.app,
	})

	// ---- setup documentation ----
	if s.echo.Debug {
		swagger.Register(swagger.Options{
			Port:      config.ConfigGlobal.Server.Port,
			Group:     s.echo.Group("/swagger"),
			AccessKey: config.ConfigGlobal.Docs.Key,
		})
	}

	// ---- start server ----
	logger.Info("Start server PID: ", os.Getpid())
	if err := s.echo.Start(config.ConfigGlobal.Server.Port); err != nil {
		logger.Error("cannot starting server ", err.Error())
	}
}

func (s *server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if s.dbSQLReader != nil {
		if err := s.dbSQLReader.Close(); err != nil {
			logger.Error("cannot close database Reader SQL: ", err.Error())
		}
	}

	if s.dbSQLWriter != nil {
		if err := s.dbSQLWriter.Close(); err != nil {
			logger.Error("cannot close database Writer SQL: ", err.Error())
		}
	}

	if s.dbNoSQLReader != nil {
		if err := s.dbNoSQLReader.Client().Disconnect(ctx); err != nil {
			logger.Error("cannot close database Reader NoSQL: ", err.Error())
		}
	}

	if s.dbNoSQLWriter != nil {
		if err := s.dbNoSQLWriter.Client().Disconnect(ctx); err != nil {
			logger.Error("cannot close database Writer NoSQL: ", err.Error())
		}
	}

	if s.queueMessager != nil {
		s.queueMessager.Close()
	}

	if s.session != nil {
		if err := s.session.Close(); err != nil {
			logger.Error("cannot close TCAuth connection: ", err.Error())
		}
	}

	s.cache = nil

	if err := s.echo.Close(); err != nil {
		logger.Error("cannot close echo ", err.Error())
	}
}

// ReloadConnections all connections like DB, Nats, ...
func (s *server) ReloadConnections() {
	// Aqui o ideal é você realizar apenas o reload de conexões não do server inteiro como está no exemplo abaixo
	// Pois dando o Stop e Start, causara indisponibilidade.
	// Vejam só que sentido tem dar close no Echo e subir ele de novo, se tu só precisa fazer reconexão com o Banco de Dados
	// Então dexarei comentado abaixo um exemplo como se tu quisesse só dar reload no Auth e no Banco de dados.
	// if s.session != nil {
	// 	if err := s.session.Close(); err != nil {
	// 		logger.Error("cannot close TCAuth connection: ", err.Error())
	// 	}
	// }
	// if s.dbSQLReader != nil {
	// 	if err := s.dbSQLReader.Close(); err != nil {
	// 		logger.Error("cannot close database Reader SQL: ", err.Error())
	// 	}
	// }
	// if s.dbSQLWriter != nil {
	// 	if err := s.dbSQLWriter.Close(); err != nil {
	// 		logger.Error("cannot close database Writer SQL: ", err.Error())
	// 	}
	// }

	// // ---- start SQL connection ----
	// s.dbSQLWriter = createSqlConnection(config.ConfigGlobal.Database.Writer.URL)
	// s.dbSQLReader = createSqlConnection(config.ConfigGlobal.Database.Reader.URL)

	// // ---- start TCAuth connection ----
	// s.session = auth.NewMiddle(config.ConfigGlobal.Auth)

	s.Stop()
	s.Start()
}
