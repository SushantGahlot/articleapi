package router

import (
	"github.com/julienschmidt/httprouter"
	"github.com/sushantgahlot/articleapi/cmd/api/handlers"
	"github.com/sushantgahlot/articleapi/pkg/application"
	"github.com/sushantgahlot/articleapi/pkg/middleware"
)

func GetHandlers(app *application.Application) *httprouter.Router {
	mux := httprouter.New()
	mux.GET("/tags/:tagName/:date", middleware.RequestLogger(handlers.GetTagSummary(app)))
	mux.POST("/articles", middleware.RequestLogger(handlers.PostArticle(app)))
	mux.GET("/articles/:id", middleware.RequestLogger(handlers.GetArticle(app)))

	return mux
}
