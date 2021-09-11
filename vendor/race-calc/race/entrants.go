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

type RaceEntrant struct {
	EntrantName   string    `form:"entrant_name" json:"entrant_name" binding:"required"`
	RaceId        int       `form:"race_id" json:"race_id" binding:"required"`
	BoatClass     string    `form:"boat_class" json:"boat_class" binding:"required"`
	Py            int       `form:"py" json:"py" binding:"required"`
	FinishTime    time.Time `form:"finish_time" json:"finish_time" binding:"required"`
	ElapsedSecs   int       `form:"elapsed_secs" json:"elapsed_secs" binding:"required"`
	CorrectedSecs int       `form:"corrected_secs" json:"corrected_secs" binding:"required"`
}
type RaceEntrants []RaceEntrant

const insertSql string = "INSERT INTO entrants (entrant_name, race_id, boat_class, py, finish_time, elapsed_seconds, corrected_seconds) VALUES ($1, $2, $3, $4, $5, $6, $7)"

func PostEntrants(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var entrants RaceEntrants

		if err := c.BindJSON(&entrants); err != nil {
			log.Print("Error parsing entrants data..", err)
			return
		}

		log.Printf("Parsed entrants OK: count: %q", len(entrants))
		if len(entrants) == 0 {
			return
		}

		if tx, txErr := db.Begin(); txErr == nil {
			for i := 0; i < len(entrants); i++ {
				var ent RaceEntrant = entrants[i]
				if _, err := db.Exec(
					insertSql,
					ent.EntrantName,
					ent.RaceId,
					ent.BoatClass,
					ent.Py,
					ent.FinishTime,
					ent.ElapsedSecs,
					ent.CorrectedSecs,
				); err != nil {
					c.String(http.StatusInternalServerError,
						fmt.Sprintf("Error inserting entrants  table: %q", err))
					tx.Rollback()
					return
				}
			}
			var commitErr = tx.Commit()
			if commitErr != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error committing transaction entrants  table: %q", commitErr))
			}
		}

		c.IndentedJSON(http.StatusCreated, len(entrants))
	}
}
