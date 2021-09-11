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

// type RaceEntrant struct {
// 	Entrant       string `form:"entrant" json:"entrant" binding:"required"`
// 	BoatClass     string `form:"boat_class" json:"boat_class" binding:"required"`
// 	FinishTime    string `form:"finish_time" json:"finish_time" binding:"required"`
// 	ElapsedSecs   int    `form:"elapsed_secs" json:"elapsed_secs" binding:"required"`
// 	CorrectedSecs int    `form:"corrected_secs" json:"corrected_secs" binding:"required"`
// }
// type Race []RaceEntrant

type RaceDef struct {
	RaceName  string    `form:"race_name" json:"race_name" binding:"required"`
	StartTime time.Time `form:"start_time" json:"start_time" binding:"required"`
}

func PostRace(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var raceData RaceDef

		if err := c.BindJSON(&raceData); err != nil {
			log.Print("Error parsing race data..", err)
			return
		}

		log.Printf("Parsed Race OK: %q %q", raceData.RaceName, raceData.StartTime.String())

		if _, err := db.Exec("INSERT INTO race (race_name, start_time) VALUES ($1, $2)", raceData.RaceName, raceData.StartTime); err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error inserting into Race table: %q", err))
			return
		}
		c.IndentedJSON(http.StatusCreated, raceData)
	}
}
