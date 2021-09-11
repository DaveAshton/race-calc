package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"race-calc/race"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
	"github.com/russross/blackfriday"
)

func repeatHandler(r int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var buffer bytes.Buffer
		for i := 0; i < r; i++ {
			buffer.WriteString("Hello from Go!\n")
		}
		c.String(http.StatusOK, buffer.String())
	}
}

func dbFunc(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS ticks (tick timestamp)"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error creating database table: %q", err))
			return
		}

		if _, err := db.Exec("INSERT INTO ticks VALUES (now())"); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error incrementing tick: %q", err))
			return
		}

		rows, err := db.Query("SELECT tick FROM ticks")
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error reading ticks: %q", err))
			return
		}

		defer rows.Close()
		for rows.Next() {
			var tick time.Time
			if err := rows.Scan(&tick); err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error scanning ticks: %q", err))
				return
			}
			c.String(http.StatusOK, fmt.Sprintf("Read from DB: %s\n", tick.String()))
		}
	}
}

// type RaceEntrant struct {
// 	Entrant       string `form:"entrant" json:"entrant" binding:"required"`
// 	BoatClass     string `form:"boat_class" json:"boat_class" binding:"required"`
// 	FinishTime    string `form:"finish_time" json:"finish_time" binding:"required"`
// 	ElapsedSecs   int    `form:"elapsed_secs" json:"elapsed_secs" binding:"required"`
// 	CorrectedSecs int    `form:"corrected_secs" json:"corrected_secs" binding:"required"`
// }
// type Race []RaceEntrant

// func postRace(db *sql.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var raceData Race

// 		if err := c.BindJSON(&raceData); err != nil {
// 			log.Print("Error parsing race data", err)
// 			return
// 		}

// 		log.Print("Parsed OK")
// 		c.IndentedJSON(http.StatusCreated, raceData)
// 	}
// }

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Print("$PORT must be set")
		port = "5000"
	}

	tStr := os.Getenv("REPEAT")
	repeat, err := strconv.Atoi(tStr)
	if err != nil {
		log.Printf("Error converting $REPEAT to an int: %q - Using default\n", err)
		repeat = 5
	}

	connString := os.Getenv("DATABASE_URL") + "sslmode=disable"
	log.Printf("DATABASE_URL: %q", os.Getenv("DATABASE_URL"))
	log.Printf("connection string: %q", connString)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/mark", func(c *gin.Context) {
		c.String(http.StatusOK, string(blackfriday.Run([]byte("**hi!**"))))
	})

	router.GET("/repeat", repeatHandler(repeat))

	router.GET("/db", dbFunc(db))

	router.POST("/race", race.PostRace(db))

	router.Run(":" + port)
}
