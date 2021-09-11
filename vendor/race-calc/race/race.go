package race

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

type RaceDef struct {
	RaceName  string    `form:"race_name" json:"race_name" binding:"required"`
	StartTime time.Time `form:"start_time" json:"start_time" binding:"required"`
}

type LastId struct {
	Id int `form:"id" json:"id" binding:"required"`
}

func PostRace(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var raceData RaceDef

		if err := c.BindJSON(&raceData); err != nil {
			log.Print("Error parsing race data..", err)
			return
		}

		log.Printf("Parsed Race OK: %q %q", raceData.RaceName, raceData.StartTime.String())

		lastInsertId := 0
		if err := db.QueryRow(
			"INSERT INTO race (race_name, start_time) VALUES ($1, $2) returning race_id",
			raceData.RaceName,
			raceData.StartTime,
		).Scan(&lastInsertId); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error inserting into Race table: %q", err))
			return
		}
		lastId := LastId{Id: lastInsertId}
		c.IndentedJSON(http.StatusCreated, lastId)
	}
}
