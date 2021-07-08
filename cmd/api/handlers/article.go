package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
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

func PostArticle(app *application.Application) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		defer r.Body.Close()

		article := models.Article{}
		err := json.NewDecoder(r.Body).Decode(&article)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if len(article.Tags) == 0 {
			http.Error(w, "Error. Tags can not be empty", http.StatusBadRequest)
			return
		}

		err = checkForDuplicates(app, &article, w)
		if err != nil {
			fmt.Println(err)
			return
		}

		validatedTags := validateTags(article.Tags)

		if len(validatedTags) == 0 {
			http.Error(w, "Error. Tags can not be empty", http.StatusBadRequest)
			return
		}

		article.ID, err = uuid.NewV4()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		err = article.Insert(app, r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		savedTags := getSavedTags(app, validatedTags)

		tagsToSave := tagsToInsert(savedTags, validatedTags)

		tagsInserted, err := saveTags(app, tagsToSave)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Slice of tag ids that need to be mapped with saved article
		var tagsToMap []uuid.UUID

		for _, tag := range savedTags {
			tagsToMap = append(tagsToMap, tag.ID)
		}

		for _, tag := range tagsInserted {
			tagsToMap = append(tagsToMap, tag.ID)
		}

		err = mapArticleToTags(app, &article, tagsToMap)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(article)
		w.Write(response)
	}
}

func mapArticleToTags(app *application.Application, ar *models.Article, tagIDs []uuid.UUID) error {
	stmt := `
		INSERT INTO article_tag (article_id, tag_id)
		VALUES %s
		ON CONFLICT DO NOTHING
	`
	argsPerRow := 2
	valueArgs := make([]interface{}, 0, len(tagIDs))

	for _, tag := range tagIDs {
		valueArgs = append(valueArgs, ar.ID, tag)
	}

	stmt = app.DB.GetBulkInsertQuery(stmt, argsPerRow, len(tagIDs))
	commandTag, err := app.DB.DBClient.Exec(context.Background(), stmt, valueArgs...)

	if err != nil || commandTag.RowsAffected() == 0 {
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Error. Nothing got inserted when bulk inserting tags")
		}
		return err
	}
	return nil
}

func saveTags(app *application.Application, tagsToSave []string) ([]models.Tag, error) {
	var tagStructs []models.Tag
	if len(tagsToSave) == 0 {
		return tagStructs, nil
	}
	stmt := `
		INSERT INTO tags (tag_id, tag)
		VALUES %s
	`
	argsPerRow := 2
	valueArgs := make([]interface{}, 0, len(tagsToSave))
	numQueries := 0

	for _, tag := range tagsToSave {
		tag_id, err := uuid.NewV4()
		if err != nil {
			fmt.Println("Could not insert tag", tag)
			continue
		}

		valueArgs = append(valueArgs, tag_id, tag)
		tagStructs = append(tagStructs, models.Tag{ID: tag_id})
		numQueries += 1
	}

	stmt = app.DB.GetBulkInsertQuery(stmt, argsPerRow, numQueries)
	commandTag, err := app.DB.DBClient.Exec(context.Background(), stmt, valueArgs...)

	if err != nil || commandTag.RowsAffected() == 0 {
		if err != nil {
			fmt.Println(err, "Error in bulk insert of tags.")
		} else {
			fmt.Println("Error. Nothing got inserted when bulk inserting tags")
		}
		return nil, err
	}

	return tagStructs, nil
}

func getSavedTags(app *application.Application, validatedTags []string) []models.Tag {
	stmt := `
	SELECT *
	FROM tags
	WHERE tag = ANY ($1)
	`

	var savedTags []models.Tag

	rows, err := app.DB.DBClient.Query(context.Background(), stmt, validatedTags)

	if err != nil {
		fmt.Println("Got error", err)
	}

	defer rows.Close()
	tag := models.Tag{}

	for rows.Next() {
		err = rows.Scan(&tag.ID, &tag.Tag)

		if err != nil {
			fmt.Println("Got error", err)
			continue
		}
		savedTags = append(savedTags, tag)
	}

	return savedTags
}

func tagsToInsert(savedTags []models.Tag, validatedTags []string) []string {
	savedTagsMap := make(map[string]bool)
	unsavedTags := make([]string, 0, len(validatedTags))

	// Populate map with saved tags
	for _, savedTag := range savedTags {
		savedTagsMap[savedTag.Tag] = true
	}

	// If a tag is not in saved tags map, append it to unsaved tags
	for _, tag := range validatedTags {
		if _, value := savedTagsMap[tag]; !value {
			unsavedTags = append(unsavedTags, tag)
		}
	}

	return unsavedTags
}

func validateTags(tags []string) []string {
	validatedTags := make([]string, 0, len(tags))
	for _, tag := range tags {
		if len(tag) != 0 {
			validatedTags = append(validatedTags, tag)
		}
	}

	if len(validatedTags) > 1 {
		uniqueTagsMap := make(map[string]bool)
		uniqueTags := make([]string, 0, len(validatedTags))
		for _, tag := range validatedTags {
			if _, value := uniqueTagsMap[tag]; !value {
				uniqueTagsMap[tag] = true
				uniqueTags = append(uniqueTags, tag)
			}
		}
		return uniqueTags
	}
	return validatedTags
}

func checkForDuplicates(app *application.Application, ar *models.Article, w http.ResponseWriter) error {
	var testString string
	stmt := `
	SELECT article_title
	FROM article
	WHERE article_title=$1
	LIMIT 1
	`

	err := app.DB.DBClient.QueryRow(context.Background(), stmt, ar.ArticleTitle).Scan(&testString)

	if err != nil && err != pgx.ErrNoRows {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return err
	}

	if len(testString) > 0 {
		http.Error(w, "Error. Article title already exists in the database", http.StatusBadRequest)
		return errors.New("duplicate article title")
	}

	stmt = `
	SELECT article_body
	FROM article
	WHERE article_body=$1
	LIMIT 1
	`

	err = app.DB.DBClient.QueryRow(context.Background(), stmt, ar.ArticleBody).Scan(&testString)

	if err != nil && err != pgx.ErrNoRows {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return errors.New("internal server error")
	}

	if len(testString) > 0 {
		http.Error(w, "Error. Article body already exists in the database", http.StatusBadRequest)
		return errors.New("duplicate article body")
	}

	return nil
}
