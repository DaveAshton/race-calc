package race

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

type RaceEntrant struct {
	EntrantId     int       `form:"entrant_id" json:"entrant_id,omitempty"`
	EntrantName   string    `form:"entrant_name" json:"entrant_name" binding:"required"`
	RaceId        int       `form:"race_id" json:"race_id" binding:"required"`
	BoatClass     string    `form:"boat_class" json:"boat_class" binding:"required"`
	Py            int       `form:"py" json:"py" binding:"required"`
	FinishTime    time.Time `form:"finish_time" json:"finish_time" binding:"required"`
	ElapsedSecs   int       `form:"elapsed_secs" json:"elapsed_secs" binding:"required"`
	CorrectedSecs int       `form:"corrected_secs" json:"corrected_secs" binding:"required"`
}
type RaceEntrants []RaceEntrant

const insertSql string = `INSERT INTO entrants (entrant_name, 
	race_id,
	boat_class,
	py,
	finish_time, 
	elapsed_seconds,
	corrected_seconds) VALUES ($1, $2, $3, $4, $5, $6, $7)`

func PostEntrants(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var entrants RaceEntrants

		if err := c.BindJSON(&entrants); err != nil {
			log.Print("Error parsing entrants data..", err)
			return
		}

		log.Printf("Parsed entrants OK: count: %q", len(entrants))
		if len(entrants) == 0 {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error inserting entrants, no entrants passed in json"))
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

const selectSql string = `SELECT 
	entrant_id, 
	entrant_name, 
	race_id, 
	boat_class, 
	py, 
	finish_time, 
	elapsed_seconds, 
	corrected_seconds 
	FROM entrants where race_id = $1`

func getQueryAsInt(c *gin.Context, queryParam string) (int, string, error) {
	var val = c.Query(queryParam)
	var ret, err = strconv.Atoi(val)
	if err != nil {
		log.Println(err)
		return 0, "", err
	}
	return ret, val, nil
}

func GetEntrants(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		const query = "race_id"
		_, raceId, err := getQueryAsInt(c, "race_id")

		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error parsing query string: %q value: %q", query, raceId))
			return
		}

		log.Printf("About to search for entrants with race_id = %q..", raceId)
		rows, err := db.Query(selectSql, raceId)
		if err != nil {
			c.String(http.StatusInternalServerError,
				fmt.Sprintf("Error reading entrants table: %q", err))
			return
		}

		defer rows.Close()

		var entrants RaceEntrants
		for rows.Next() {
			var entrant RaceEntrant
			if err := rows.Scan(
				&entrant.EntrantId,
				&entrant.EntrantName,
				&entrant.RaceId,
				&entrant.BoatClass,
				&entrant.Py,
				&entrant.FinishTime,
				&entrant.ElapsedSecs,
				&entrant.CorrectedSecs,
			); err != nil {
				c.String(http.StatusInternalServerError,
					fmt.Sprintf("Error scanning entrant: %q", err))
				return
			}
			entrants = append(entrants, entrant)
		}
		c.JSON(http.StatusOK, entrants)
	}
}
