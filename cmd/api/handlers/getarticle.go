package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sushantgahlot/articleapi/cmd/api/models"
	"github.com/sushantgahlot/articleapi/pkg/application"
)

func GetArticle(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer r.Body.Close()

		articleId := p.ByName("id")
		_, err := uuid.FromString(articleId)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		stmt := `
	SELECT *
	FROM article
	WHERE article_id = $1
	`

		article := models.Article{}

		err = app.DB.DBClient.QueryRow(context.Background(), stmt, articleId).Scan(&article.ID, &article.ArticleTitle, &article.ArticleDate, &article.ArticleBody)
		if err != nil {
			fmt.Println("Error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		stmt = `
	SELECT tag_id
	FROM article_tag
	WHERE article_id = $1
	`

		var tagIDs []uuid.UUID
		var tagID uuid.UUID
		rows, err := app.DB.DBClient.Query(context.Background(), stmt, articleId)
		if err != nil {
			fmt.Println("Error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			err = rows.Scan(&tagID)

			if err != nil {
				fmt.Println("Got error", err)
				continue
			}
			tagIDs = append(tagIDs, tagID)
		}

		stmt = `
	SELECT *
	FROM tags
	WHERE tag_id = ANY ($1)
	`

		var tagStructs []models.Tag
		var tagStruct models.Tag
		rows, err = app.DB.DBClient.Query(context.Background(), stmt, tagIDs)
		if err != nil {
			fmt.Println("Error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			err = rows.Scan(&tagStruct.ID, &tagStruct.Tag)

			if err != nil {
				fmt.Println("Got error", err)
				continue
			}
			tagStructs = append(tagStructs, tagStruct)
		}

		for _, tag := range tagStructs {
			article.Tags = append(article.Tags, tag.Tag)
		}

		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(article)
		w.Write(response)
	}
}
