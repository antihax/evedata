package killmailstats

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
)

type statFunc func(age int) error
type entityStatFunc func() error

func funcName(f interface{}) string {
	p := reflect.ValueOf(f).Pointer()
	rf := runtime.FuncForPC(p)
	return rf.Name()
}

// Run the service
func (s *KillmailStats) runStats() {
	ages := [5]int{18250, 7, 14, 30, 90}
	stats := []statFunc{
		s.wars,
		s.ganks,
		s.lowsec,
		s.nullsec,
		s.highsec,
		s.wh,
		s.lowsecFW,
		s.highsecFW,
		s.total,
	}

	entityStats := []entityStatFunc{
		s.entity_highsec,
		s.entity_lowsec,
		s.entity_nullsec,
		s.entity_wh,
	}

	for _, age := range ages {
		for _, stat := range stats {
			fmt.Printf("Processing %d for %s\n", age, funcName(stat))
			err := stat(age)
			if err != nil {
				log.Println(err)
			}
		}
	}
	for _, stat := range entityStats {
		fmt.Printf("Processing  %s\n", funcName(stat))
		err := stat()
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *KillmailStats) lowsec(age int) error {
	return s.doSQL(`
		INSERT INTO evedata.killmailStatistics  (month, year, characterAge, lowsec)
		SELECT MONTH(K.killTime) AS month, YEAR(K.killTime) AS year, ? as characterAge, count(DISTINCT K.victimCharacterID)
		FROM evedata.killmailAttributes A
		INNER JOIN evedata.killmails K ON K.id = A.id
		INNER JOIN mapSolarSystems S ON S.solarSystemID = K.solarSystemID
			WHERE characterAge < ?
			AND ROUND(S.security, 1) < 0.5 
			AND ROUND(S.security, 1) > 0.0
			AND K.killTime > DATE_FORMAT(DATE_SUB(UTC_TIMESTAMP(), INTERVAL 60 DAY), '%Y-%m-01')
		GROUP BY YEAR(K.killTime), MONTH(K.killTime)
		ON DUPLICATE KEY UPDATE lowsec = VALUES(lowsec);
	`, age, age)
}

func (s *KillmailStats) lowsecFW(age int) error {
	return s.doSQL(`
		INSERT INTO evedata.killmailStatistics  (month, year, characterAge, lowsecFW)
		SELECT MONTH(K.killTime) AS month, YEAR(K.killTime) AS year, ? as characterAge, count(DISTINCT K.victimCharacterID)
		FROM evedata.killmailAttributes A
		INNER JOIN evedata.killmails K ON K.id = A.id
		INNER JOIN mapSolarSystems S ON S.solarSystemID = K.solarSystemID
			WHERE characterAge < ?
			AND ROUND(S.security, 1) < 0.5 
			AND ROUND(S.security, 1) > 0.0
			AND K.killTime > DATE_FORMAT(DATE_SUB(UTC_TIMESTAMP(), INTERVAL 60 DAY), '%Y-%m-01')
			AND K.factionID > 0
		GROUP BY YEAR(K.killTime), MONTH(K.killTime)
		ON DUPLICATE KEY UPDATE lowsecFW = VALUES(lowsecFW);
	`, age, age)
}

func (s *KillmailStats) highsec(age int) error {
	return s.doSQL(`
		INSERT INTO evedata.killmailStatistics  (month, year, characterAge, highsec)
		SELECT MONTH(K.killTime) AS month, YEAR(K.killTime) AS year, ? as characterAge, count(DISTINCT K.victimCharacterID)
		FROM evedata.killmailAttributes A
		INNER JOIN evedata.killmails K ON K.id = A.id
		INNER JOIN mapSolarSystems S ON S.solarSystemID = K.solarSystemID
			WHERE characterAge < ?
			AND ROUND(S.security, 1) >= 0.5
			AND K.killTime > DATE_FORMAT(DATE_SUB(UTC_TIMESTAMP(), INTERVAL 60 DAY), '%Y-%m-01')
		GROUP BY YEAR(K.killTime), MONTH(K.killTime)
		ON DUPLICATE KEY UPDATE highsec = VALUES(highsec);
	`, age, age)
}

func (s *KillmailStats) highsecFW(age int) error {
	return s.doSQL(`
		INSERT INTO evedata.killmailStatistics  (month, year, characterAge, highsecFW)
		SELECT MONTH(K.killTime) AS month, YEAR(K.killTime) AS year, ? as characterAge, count(DISTINCT K.victimCharacterID)
		FROM evedata.killmailAttributes A
		INNER JOIN evedata.killmails K ON K.id = A.id
		INNER JOIN mapSolarSystems S ON S.solarSystemID = K.solarSystemID
			WHERE characterAge < ?
			AND ROUND(S.security, 1) >= 0.5 
			AND K.factionID > 0
			AND K.killTime > DATE_FORMAT(DATE_SUB(UTC_TIMESTAMP(), INTERVAL 60 DAY), '%Y-%m-01')
		GROUP BY YEAR(K.killTime), MONTH(K.killTime)
		ON DUPLICATE KEY UPDATE highsecFW = VALUES(highsecFW);
	`, age, age)
}

func (s *KillmailStats) wars(age int) error {
	return s.doSQL(`
		INSERT INTO evedata.killmailStatistics  (month, year, characterAge, wars)
		SELECT MONTH(K.killTime) AS month, YEAR(K.killTime) AS year, ? as characterAge, count(DISTINCT K.victimCharacterID)
		FROM evedata.killmailAttributes A
		INNER JOIN evedata.killmails K ON K.id = A.id
		INNER JOIN mapSolarSystems S ON S.solarSystemID = K.solarSystemID
			WHERE characterAge < ?
			AND ROUND(S.security, 1) >= 0.5 AND K.warID > 0
			AND K.killTime > DATE_FORMAT(DATE_SUB(UTC_TIMESTAMP(), INTERVAL 60 DAY), '%Y-%m-01')
		GROUP BY YEAR(K.killTime), MONTH(K.killTime)
		ON DUPLICATE KEY UPDATE wars = VALUES(wars);
	`, age, age)
}

func (s *KillmailStats) ganks(age int) error {
	return s.doSQL(`
		INSERT INTO evedata.killmailStatistics  (month, year, characterAge, ganks)
		SELECT MONTH(K.killTime) AS month, YEAR(K.killTime) AS year, ? as characterAge, count(DISTINCT K.victimCharacterID)
		FROM evedata.killmailAttributes A
		INNER JOIN evedata.killmails K ON K.id = A.id
		INNER JOIN mapSolarSystems S ON S.solarSystemID = K.solarSystemID
			WHERE characterAge < ?
			AND ROUND(S.security, 1) >= 0.5 AND A.meanSecurity < -4.5
			AND K.killTime > DATE_FORMAT(DATE_SUB(UTC_TIMESTAMP(), INTERVAL 60 DAY), '%Y-%m-01')
		GROUP BY YEAR(K.killTime), MONTH(K.killTime)
		ON DUPLICATE KEY UPDATE ganks = VALUES(ganks);
	`, age, age)
}

func (s *KillmailStats) wh(age int) error {
	return s.doSQL(`
		INSERT INTO evedata.killmailStatistics  (month, year, characterAge, wh)
		SELECT MONTH(K.killTime) AS month, YEAR(K.killTime) AS year, ? as characterAge, count(DISTINCT K.victimCharacterID)
		FROM evedata.killmailAttributes A
		INNER JOIN evedata.killmails K ON K.id = A.id
		INNER JOIN mapSolarSystems S ON S.solarSystemID = K.solarSystemID
			WHERE characterAge < ?
			AND S.solarSystemID >= 31000000 AND S.solarSystemID < 32000000
			AND K.killTime > DATE_FORMAT(DATE_SUB(UTC_TIMESTAMP(), INTERVAL 60 DAY), '%Y-%m-01')
		GROUP BY YEAR(K.killTime), MONTH(K.killTime)       
		ON DUPLICATE KEY UPDATE wh = VALUES(wh);
	`, age, age)
}

func (s *KillmailStats) nullsec(age int) error {
	return s.doSQL(`
		INSERT INTO evedata.killmailStatistics  (month, year, characterAge, nullsec)
		SELECT MONTH(K.killTime) AS month, YEAR(K.killTime) AS year, ? as characterAge, count(DISTINCT K.victimCharacterID)
		FROM evedata.killmailAttributes A
		INNER JOIN evedata.killmails K ON K.id = A.id
		INNER JOIN mapSolarSystems S ON S.solarSystemID = K.solarSystemID
			WHERE characterAge < ?
			AND S.solarSystemID < 31000000 AND ROUND(S.security, 1) <= 0.0 
			AND K.killTime > DATE_FORMAT(DATE_SUB(UTC_TIMESTAMP(), INTERVAL 60 DAY), '%Y-%m-01')
		GROUP BY YEAR(K.killTime), MONTH(K.killTime)       
		ON DUPLICATE KEY UPDATE nullsec = VALUES(nullsec);
	`, age, age)
}

func (s *KillmailStats) total(age int) error {
	return s.doSQL(`
		INSERT INTO evedata.killmailStatistics  (month, year, characterAge, total)
		SELECT MONTH(K.killTime) AS month, YEAR(K.killTime) AS year, ? as characterAge, count(DISTINCT K.victimCharacterID)
		FROM evedata.killmailAttributes A
		INNER JOIN evedata.killmails K ON K.id = A.id
			WHERE characterAge < ?
			AND K.killTime > DATE_FORMAT(DATE_SUB(UTC_TIMESTAMP(), INTERVAL 60 DAY), '%Y-%m-01')
		GROUP BY YEAR(K.killTime), MONTH(K.killTime)       
		ON DUPLICATE KEY UPDATE total = VALUES(total);
	`, age, age)
}

func (s *KillmailStats) entity_highsec() error {
	return s.doSQL(`
		INSERT INTO evedata.killmailKillers (month, year, id, kills, area)
		SELECT MONTH(K.killTime) AS month, YEAR(K.killTime) AS year, E.ID, COUNT(DISTINCT K.id) AS kills, "highsec"
				FROM evedata.killmailAttackers A 
				INNER JOIN evedata.killmails K ON K.id = A.id
				INNER JOIN evedata.entities E ON E.id = IF(A.allianceID , A.allianceID, A.corporationID)
				INNER JOIN mapSolarSystems S ON S.solarSystemID = K.solarSystemID
				WHERE ROUND(S.security, 1) >= 0.5
				AND K.killTime > DATE_FORMAT(DATE_SUB(UTC_TIMESTAMP(), INTERVAL 60 DAY), '%Y-%m-01') 
				GROUP BY YEAR(K.killTime), MONTH(K.killTime), E.ID
				ON DUPLICATE KEY UPDATE kills=VALUES(kills);
	`)
}

func (s *KillmailStats) entity_lowsec() error {
	return s.doSQL(`
		INSERT INTO evedata.killmailKillers (month, year, id, kills, area)
		SELECT MONTH(K.killTime) AS month, YEAR(K.killTime) AS year, E.ID, COUNT(DISTINCT K.id) AS kills, "lowsec"
				FROM evedata.killmailAttackers A 
				INNER JOIN evedata.killmails K ON K.id = A.id
				INNER JOIN evedata.entities E ON E.id = IF(A.allianceID , A.allianceID, A.corporationID)
				INNER JOIN mapSolarSystems S ON S.solarSystemID = K.solarSystemID
					WHERE ROUND(S.security, 1) < 0.5 AND ROUND(S.security, 1) > 0.0 
					AND K.killTime > DATE_FORMAT(DATE_SUB(UTC_TIMESTAMP(), INTERVAL 60 DAY), '%Y-%m-01') 

					GROUP BY YEAR(K.killTime), MONTH(K.killTime), E.ID
				ON DUPLICATE KEY UPDATE kills=VALUES(kills);
	`)
}

func (s *KillmailStats) entity_nullsec() error {
	return s.doSQL(`
	INSERT INTO evedata.killmailKillers (month, year, id, kills, area)
		SELECT MONTH(K.killTime) AS month, YEAR(K.killTime) AS year, E.ID, COUNT(DISTINCT K.id) AS kills, "nullsec"
		FROM evedata.killmailAttackers A 
		INNER JOIN evedata.killmails K ON K.id = A.id
		INNER JOIN evedata.entities E ON E.id = IF(A.allianceID , A.allianceID, A.corporationID)
		INNER JOIN mapSolarSystems S ON S.solarSystemID = K.solarSystemID
			WHERE S.solarSystemID < 31000000 AND ROUND(S.security, 1) <= 0.0 
			AND K.killTime > DATE_FORMAT(DATE_SUB(UTC_TIMESTAMP(), INTERVAL 60 DAY), '%Y-%m-01') 
		GROUP BY YEAR(K.killTime), MONTH(K.killTime), E.ID
		ON DUPLICATE KEY UPDATE kills=VALUES(kills);
	`)
}

func (s *KillmailStats) entity_wh() error {
	return s.doSQL(`
	INSERT INTO evedata.killmailKillers (month, year, id, kills, area)
		SELECT MONTH(K.killTime) AS month, YEAR(K.killTime) AS year, E.ID, COUNT(DISTINCT K.id) AS kills, "wh"
		FROM evedata.killmailAttackers A 
		INNER JOIN evedata.killmails K ON K.id = A.id
		INNER JOIN evedata.entities E ON E.id = IF(A.allianceID , A.allianceID, A.corporationID)
		INNER JOIN mapSolarSystems S ON S.solarSystemID = K.solarSystemID
			WHERE S.solarSystemID >= 31000000 AND S.solarSystemID < 32000000
			AND K.killTime > DATE_FORMAT(DATE_SUB(UTC_TIMESTAMP(), INTERVAL 60 DAY), '%Y-%m-01') 
		GROUP BY YEAR(K.killTime), MONTH(K.killTime), E.ID
		ON DUPLICATE KEY UPDATE kills=VALUES(kills);
	`)
}
