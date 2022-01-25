package server

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/tradersclub/TCTemplateBack/model"
	"github.com/tradersclub/TCUtils/logger"
	"github.com/tradersclub/TCUtils/mongodb/v2"
	"github.com/tradersclub/TCUtils/tcerr"
)

func createSqlConnection(url string) *sqlx.DB {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "mysql", url)
	if err != nil {
		logger.Fatal("createSqlConnection: ", err.Error())
	}

	return db
}

func createNoSqlConnection(url, database string, isReader bool) mongodb.Database {
	return mongodb.New(
		mongodb.Options{
			URI:      url,
			Database: database,
			IsReader: isReader,
		},
	)
}

func createHTTPErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		if err := c.JSON(tcerr.GetHTTPCode(err), model.Response{Err: err}); err != nil {
			logger.ErrorContext(c.Request().Context(), err)
		}
	}
}
