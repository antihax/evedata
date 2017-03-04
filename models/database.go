package models

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	database      *sqlx.DB
	SQLTimeFormat = "2006-01-02 15:04:05"
)

// Set Database handle
func SetDatabase(DB *sqlx.DB) {
	database = DB
}

func SetupDatabase(driver string, spec string) (*sqlx.DB, error) {
	var err error

	// Build Connection Pool
	if database, err = sqlx.Connect(driver, spec); err != nil {
		return nil, err
	}

	// Check we can connect
	if err = database.Ping(); err != nil {
		return nil, err
	}

	// Put some finite limits to prevent opening too many connections
	database.SetConnMaxLifetime(time.Minute * 2)
	database.SetMaxIdleConns(100)

	SetDatabase(database)
	return database, nil
}

func DumpDatabase(file string, db string) (err error) {
	f, err := os.Create(file)
	defer f.Close()

	f.WriteString(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;\n\n", db))

	f.WriteString(fmt.Sprintf("USE %s;\n\n", db))

	tables, err := database.Query(`SELECT table_name
			FROM information_schema.TABLES WHERE table_schema = ?;`, db)
	if err != nil {
		return err
	}
	defer tables.Close()

	for tables.Next() {
		var table, create string
		err = tables.Scan(&table)
		if err != nil {
			return err
		}
		row := database.QueryRow(fmt.Sprintf(`SHOW CREATE TABLE %s.%s;`, db, table))
		err = row.Scan(&table, &create)
		if err != nil {
			return err
		}
		f.WriteString(fmt.Sprintf("%s;\n\n", create))
	}

	f.WriteString(`
		DELIMITER $$
		CREATE FUNCTION constellationIDBySolarSystem(system INT UNSIGNED) RETURNS int(10) unsigned
			DETERMINISTIC
		BEGIN
			DECLARE constellation int(10) unsigned;
			SELECT constellationID INTO constellation
				FROM eve.mapSolarSystems
				WHERE solarSystemID = system
				LIMIT 1;
			
		RETURN constellation;
		END$$
		DELIMITER ;
		`)

	f.WriteString(`
		DELIMITER $$
		CREATE FUNCTION closestCelestial(s INT UNSIGNED, x1 FLOAT, y1 FLOAT, z1 FLOAT) RETURNS int(10) unsigned
			DETERMINISTIC
		BEGIN
			DECLARE celestialID int(10) unsigned;
			SELECT itemID INTO celestialID
				FROM eve.mapDenormalize
				WHERE orbitID IS NOT NULL AND solarSystemID = s
				ORDER BY POW(( x1 - x), 2) + POW(( y1 - y), 2) + POW(( z1 - z), 2)
				LIMIT 1;
			
		RETURN celestialID;
		END$$
		DELIMITER ;
		`)

	f.WriteString(`DELIMITER $$
		CREATE FUNCTION regionIDBySolarSystem(system INT UNSIGNED) RETURNS int(10) unsigned
			DETERMINISTIC
		BEGIN
			DECLARE region int(10) unsigned;
			SELECT regionID INTO region
				FROM eve.mapSolarSystems
				WHERE solarSystemID = system
				LIMIT 1;
			
		RETURN region;
		END$$
		DELIMITER ;
		`)

	f.WriteString(`DELIMITER $$
		CREATE FUNCTION regionIDByStructureID(structure BIGINT UNSIGNED) RETURNS int(10) unsigned
			DETERMINISTIC
		BEGIN
			DECLARE region int(10) unsigned;
			SELECT regionID INTO region
				FROM eve.mapSolarSystems M
				INNER JOIN evedata.structures S ON S.solarSystemID = M.solarSystemID
				WHERE stationID = structure
				LIMIT 1;
			
		RETURN region;
		END$$
		DELIMITER ;
		`)

	f.WriteString(`DELIMITER $$
		CREATE FUNCTION raceByID(inRaceID int UNSIGNED) RETURNS VARCHAR(20) 
			DETERMINISTIC
		BEGIN
			DECLARE race VARCHAR(20) ;
			SELECT raceName INTO race
				FROM eve.chrRaces 
				WHERE raceID = inRaceID
				LIMIT 1;
			
		RETURN race;
		END$$
		DELIMITER ;
		`)

	return
}
