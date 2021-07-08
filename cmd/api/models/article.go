package models

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/sushantgahlot/articleapi/pkg/application"
)

type Article struct {
	ID           uuid.UUID `json:"ID"`
	ArticleTitle string    `json:"articleTitle"`
	ArticleDate  time.Time `json:"articleDate"`
	ArticleBody  string    `json:"articleBody"`
	Tags         []string  `json:"tags"`
}

func (ar *Article) Insert(app *application.Application, ctx context.Context) error {
	today := time.Now()

	stmt := `
	INSERT INTO article (
		article_id,
		article_title,
		article_date,
		article_body
	)
	VALUES($1, $2, $3, $4)
	`

	commandTag, err := app.DB.DBClient.Exec(
		context.Background(),
		stmt,
		ar.ID,
		ar.ArticleTitle,
		today,
		ar.ArticleBody,
	)

	if err != nil || commandTag.RowsAffected() != 1 {
		return err
	}

	return nil
}
