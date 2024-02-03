package main

import (
	"context"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool

// User model
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"fullname"`
	Password string `json:"password"`
	City     string `json:"city"`
	Age      int    `json:"age"`
}

func md5Hash(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}

func signupHandler(c *gin.Context) {
	var user User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if the username already exists
	var count int
	err = pool.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM users WHERE email = $1",
		user.Email).Scan(&count)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking email"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User already exists"})
		return
	}

	hashedPassword := md5Hash(user.Password)

	// Insert the new user into the database
	_, err = pool.Query(context.Background(),
		"INSERT INTO users (email, fullname, password, city, age) VALUES ($1, $2, $3, $4, $5)",
		user.Email,
		user.FullName,
		hashedPassword,
		user.City,
		user.Age)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}

func loginHandler(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password with MD5 (for educational purposes only)
	hashedPassword := md5Hash(user.Password)

	// Check if the username and hashed password match
	err := pool.QueryRow(
		context.Background(),
		"SELECT fullname FROM users WHERE email = $1 AND password = $2",
		user.Email,
		hashedPassword).Scan(&user.FullName)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"FullName": user.FullName})
}

func getUserHandler(c *gin.Context) {
	users := []User{}

	user_id := c.DefaultQuery("id", "1")
	query := fmt.Sprintf("SELECT id, email, fullname, coalesce(city, 'NULL'), coalesce(age, -1) FROM users WHERE id=%s", user_id)

	rows, err := pool.Query(context.Background(), query)
	defer rows.Close()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FullName,
			&user.City,
			&user.Age)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

func initDb() {
	create_table_query := `
	CREATE TABLE IF NOT EXISTS users (
		id       SERIAL PRIMARY KEY,
		email    VARCHAR(64) NOT NULL,
		fullname VARCHAR(32) NOT NULL,
		password VARCHAR(32) NOT NULL,
		city     VARCHAR(16),
		age      INT
	)
	`

	_, err := pool.Query(context.Background(), create_table_query)

	if err != nil {
		log.Println("Error in init SQL queries: ", err)
	}

	insert_rows_query := `
	INSERT INTO users (email, fullname, password, city, age) 
	VALUES
		('john.doe@example.com', 'John Doe', 'b8a76c56d41e570a6e73f55c232572e9', 'New York', 28),
		('alice.smith@example.com', 'Alice Smith', '3a6da9ad70dfe6bd6129a0858aaa1fd0', 'San Francisco', 35),
		('bob.jones@example.com', 'Bob Jones', '1543a45232df76aaec95af184e246c69', 'Los Angeles', 22),
		('sara.williams@example.com', 'Sara Williams', '1543a45232df76aaec95af184e246c69', NULL, 29),
		('michael.brown@example.com', 'Michael Brown', 'd0b8dfe012c6aad2be13e3430439f581', 'Chicago', 31),
		('emily.wang@example.com', 'Emily Wang', 'b85b5926d887f5dfa4782549a3e97793', NULL, 26),
		('david.nguyen@example.com', 'David Nguyen', '6060bcea977c78994a6587382fce4c4b', 'Houston', 33),
		('olivia.garcia@example.com', 'Olivia Garcia', 'f57b4b400b8ed7956556c24339c2d48f', 'Miami', 29),
		('ryan.miller@example.com', 'Ryan Miller', '9ed1707a05359d46a03a8ebe129c2964', 'Seattle', 38),
		('emma.davis@example.com', 'Emma Davis', 'ce28650b9ba722fbd4da8c0a4c2b8cb9', 'Denver', 27),
		('admin@example.com', 'Admin User', 'a9946a9d51be374db363c9492850c0a8', NULL, NULL)
	ON CONFLICT DO NOTHING
	`

	_, err = pool.Query(context.Background(), insert_rows_query)

	if err != nil {
		log.Println("Error in init SQL queries: ", err)
	}
}

func main() {
	pool, _ = pgxpool.Connect(context.Background(), "postgres://192.168.241.132/postgres?user=poc_update_user&password=testpassword")

	initDb()

	r := gin.Default()
	r.GET("/user", getUserHandler)
	r.POST("/signup", signupHandler)
	r.POST("/login", loginHandler)
	r.Run(":8000")

	defer pool.Close()
}
