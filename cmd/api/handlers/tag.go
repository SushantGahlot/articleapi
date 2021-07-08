package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sushantgahlot/articleapi/pkg/application"
)

type jsonResponse struct {
	Tag         string   `json:"tag"`
	Count       int      `json:"count"`
	Articles    []string `json:"articles"`
	RelatedTags []string `json:"related_tags"`
}

type articleTag struct {
	articleId string
}

func GetTagSummary(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		tag := p.ByName("tagName")
		dt := p.ByName("date")

		parsedDate, err := time.Parse("20060102", dt)

		if err != nil {
			http.Error(w, "Error parsing date", http.StatusBadRequest)
		}

		resp := jsonResponse{}
		resp.Tag = tag

		// Get all articles for the given date
		stmt := `
		SELECT article_id
		FROM article
		WHERE article_date = $1
		ORDER BY article_date DESC
		`
		var articleIDs []uuid.UUID
		var articleID uuid.UUID
		rows, err := app.DB.DBClient.Query(context.Background(), stmt, parsedDate)
		if err != nil {
			fmt.Println("Error", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			err = rows.Scan(&articleID)

			if err != nil {
				fmt.Println("Got error", err)
				continue
			}
			articleIDs = append(articleIDs, articleID)
		}

		// Get all article ids related to given tag for the given day
		stmt = `
		SELECT article_id 
		FROM article
		WHERE article_id = ANY(
			SELECT article_id
			FROM article_tag
			WHERE tag_id = (SELECT tag_id FROM tags WHERE tag=$1)
		)	
		AND article_date = $2
		`
		var articleTags []articleTag
		var articleTag articleTag
		rows, err = app.DB.DBClient.Query(context.Background(), stmt, tag, parsedDate)
		if err != nil {
			fmt.Println("Error getting related article ids", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			err = rows.Scan(&articleTag.articleId)

			if err != nil {
				fmt.Println("Got error", err)
				continue
			}
			articleTags = append(articleTags, articleTag)
		}

		// Update count
		resp.Count = len(articleTags)

		// Update response json with article IDs
		count := 0
		for _, arTag := range articleTags {
			if count == 10 {
				break
			}
			resp.Articles = append(resp.Articles, arTag.articleId)
			count++
		}

		// Get related tags
		stmt = `
		SELECT tag FROM tags WHERE tag_id = ANY (
			SELECT tag_id FROM article_tag WHERE article_id = ANY (
				SELECT article_id 
				FROM article
				WHERE article_id = ANY (
					SELECT article_id
					FROM article_tag
					WHERE tag_id = ANY (SELECT tag_id FROM article_tag WHERE article_id = ANY ($1))
				)
				AND article_date = ($2)
			)
		)
		`
		var relatedTags = make(map[string]struct{})
		var relatedTag string
		rows, err = app.DB.DBClient.Query(context.Background(), stmt, &articleIDs, parsedDate)
		if err != nil {
			fmt.Println("Error at last query", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			err = rows.Scan(&relatedTag)

			if err != nil {
				fmt.Println("Got error", err)
				continue
			}
			relatedTags[relatedTag] = struct{}{}
		}
		delete(relatedTags, tag)

		for k := range relatedTags {
			resp.RelatedTags = append(resp.RelatedTags, k)
		}

		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(resp)
		w.Write(response)
	}
}
