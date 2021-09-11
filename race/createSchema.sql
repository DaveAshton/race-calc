CREATE TABLE IF NOT EXISTS race (
    race_id INT GENERATED ALWAYS AS IDENTITY,
    race_name VARCHAR NOT NULL,
    start_time DATE NULL,
    PRIMARY KEY(race_id)
);

CREATE TABLE IF NOT EXISTS entrants (
    entrant_id INT GENERATED ALWAYS AS IDENTITY,
    entrant_name VARCHAR NOT NULL,
    race_id INT NOT NULL,
    boat_class VARCHAR NOT NULL,
    py INT NOT NULL,
    finish_time DATE NULL,
    elapsed_seconds INT NOT NULL,
    corrected_seconds INT NOT NULL,

    PRIMARY KEY(entrant_id),
    CONSTRAINT fk_race
      FOREIGN KEY(race_id) 
	  REFERENCES race(race_id)
	  ON DELETE CASCADE
);


INSERT INTO race (race_name, start_time) VALUES
('Thaw','2019-10-12T07:20:50.52Z');