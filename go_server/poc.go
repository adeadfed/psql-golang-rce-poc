package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool

type Phrase struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

func phraseHandler(c *gin.Context) {
	phrases := []Phrase{}

	phrase_id := c.DefaultQuery("id", "1")
	query := fmt.Sprintf("SELECT id, text FROM phrases WHERE id=%s", phrase_id)

	rows, err := pool.Query(context.Background(), query)
	defer rows.Close()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for rows.Next() {
		var phrase Phrase
		err := rows.Scan(&phrase.ID, &phrase.Text)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		phrases = append(phrases, phrase)
	}

	c.JSON(http.StatusOK, phrases)
}

func initDb() {
	create_table_query := `
	CREATE TABLE IF NOT EXISTS phrases (
		id serial PRIMARY KEY,
		text VARCHAR(255) NOT NULL
	)
	`

	_, err := pool.Query(context.Background(), create_table_query)

	if err != nil {
		log.Println("Error in init SQL queries: ", err)
	}

	insert_rows_query := `
	INSERT INTO phrases (text)
	VALUES
		('Hello, world!'),
		('A day in paradise.'),
		('Live, laugh, love.'),
		('Carpe diem.'),
		('The sun always shines.'),
		('Smile and the world smiles with you.'),
		('Happiness is a choice.'),
		('Keep calm and carry on.'),
		('Chase your dreams.'),
		('Life is beautiful.'),
		('Seize the moment.'),
		('Just breathe.'),
		('Adventure awaits.'),
		('Find your inner peace.'),
		('Dance in the rain')
	ON CONFLICT DO NOTHING
	`

	_, err = pool.Query(context.Background(), insert_rows_query)

	if err != nil {
		log.Println("Error in init SQL queries: ", err)
	}
}

func main() {
	pool, _ = pgxpool.Connect(context.Background(), "postgres://localhost/postgres?user=poc_user&password=poc_pass")

	initDb()

	r := gin.Default()
	r.GET("/phrases", phraseHandler)
	r.Run(":8000")

	defer pool.Close()
}
