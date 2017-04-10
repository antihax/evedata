package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/antihax/evedata/fpGrowth"
)

type KnownAlts struct {
	CharacterID   int64  `db:"characterID" json:"id"`
	CharacterName string `db:"characterName" json:"name"`
	Frequency     int    `db:"frequency" json:"frequency"`
	Type          string `db:"type" json:"type"`
	Source        uint8  `db:"source" json:"source"`
}

// Obtain Character Associates by ID.
// [BENCHMARK] 0.000 sec / 0.000 sec
func GetCharacterKnownAssociates(id int64) ([]KnownAlts, error) {
	ref := []KnownAlts{}
	if err := database.Select(&ref, `
		SELECT 	associateID AS characterID,
				frequency,
				C.name AS characterName,
				IFNULL(source, 0) AS source
		FROM evedata.characterAssociations A
		INNER JOIN evedata.characters C ON A.associateID = C.characterID
		INNER JOIN evedata.characters M ON A.characterID = M.characterID
		WHERE A.characterID = ?
		AND (M.allianceID != C.allianceID OR C.allianceID = 0) AND M.corporationID != C.corporationID
		`, id); err != nil {
		return nil, err
	}

	for i := range ref {
		ref[i].Type = "character"
	}

	return ref, nil
}

func GetCorporationKnownAssociates(id int64) ([]KnownAlts, error) {
	ref := []KnownAlts{}
	if err := database.Select(&ref, `
		SELECT	associateID AS characterID,
				SUM(frequency) AS frequency,
				C.name AS characterName,
				IFNULL(source, 0) AS source
		FROM evedata.characterAssociations A
        INNER JOIN evedata.characters C ON A.associateID = C.characterID
        INNER JOIN evedata.characters M ON A.characterID = M.characterID
        WHERE A.characterID IN (SELECT characterID FROM evedata.characters WHERE corporationID = ?)
        AND (M.allianceID != C.allianceID OR C.allianceID = 0) AND M.corporationID != C.corporationID
        GROUP BY associateID
		`, id); err != nil {
		return nil, err
	}

	for i := range ref {
		ref[i].Type = "character"
	}

	return ref, nil
}

func GetAllianceKnownAssociates(id int64) ([]KnownAlts, error) {
	ref := []KnownAlts{}
	if err := database.Select(&ref, `
		SELECT	associateID AS characterID,
				SUM(frequency) AS frequency,
				C.name AS characterName,
				IFNULL(source, 0) AS source
		FROM evedata.characterAssociations A
        INNER JOIN evedata.characters C ON A.associateID = C.characterID
        INNER JOIN evedata.characters M ON A.characterID = M.characterID
        WHERE A.characterID IN (SELECT characterID FROM evedata.characters WHERE allianceID = ?)
        AND (M.allianceID != C.allianceID OR C.allianceID = 0) AND M.corporationID != C.corporationID
        GROUP BY associateID
		`, id); err != nil {
		return nil, err
	}

	for i := range ref {
		ref[i].Type = "character"
	}

	return ref, nil
}

func BuildRelationships() error {
	if err := buildKillmailRelationships(); err != nil {
		log.Println(err)
		return err
	}
	if err := buildCorpJoinRelationships(); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Find relationships between characters in killmails
func buildKillmailRelationships() error {
	log.Printf("Character Associations: Pull Database")
	rows, err := database.Query(`
        SELECT K.id, GROUP_CONCAT(characterID) 
        FROM evedata.killmailAttackers A
        INNER JOIN evedata.killmails K ON K.id = A.id
        WHERE killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 90 DAY)
        GROUP BY K.id
        HAVING count(*) > 1 AND count(*) < 11;
        `)
	if err != nil {
		return err
	}
	return findRelationships(rows, 1)
}

// Find relationships between characters from corp history
func buildCorpJoinRelationships() error {
	log.Printf("Character Associations: Pull Database")
	rows, err := database.Query(`
        SELECT 
            UNIX_TIMESTAMP(startDate)+FLOOR(1 + (RAND() * 86028157)), 
            GROUP_CONCAT(DISTINCT H.characterID)
        FROM evedata.corporationHistory H
        INNER JOIN evedata.characters C ON C.characterID = H.characterID
        WHERE H.corporationID > 1999999 AND C.corporationID > 1000001
        GROUP BY H.corporationID, DATE(startDate)
        HAVING count(*) > 2 AND count(*) < 11;
        `)
	if err != nil {
		return err
	}
	return findRelationships(rows, 2)
}

// Find relationships in the results
func findRelationships(rows *sql.Rows, associationType uint8) error {
	log.Printf("Character Associations: Build Transaction History")
	transactions := fpGrowth.ItemSet{}

	for rows.Next() {
		var (
			transactionID int
			items         string
		)

		// Build history of common transactions
		rows.Scan(&transactionID, &items)
		transactions[transactionID] = splitToInt(items)
	}
	rows.Close()

	log.Printf("Character Associations: Build fpTree")
	fp := fpGrowth.NewFPTree(transactions, 2)

	log.Printf("Character Associations: Growth")
	associations := fp.Growth()

	log.Printf("Character Associations: Build Values")
	var values []string
	for _, association := range associations {
		for _, char1 := range association.Items {
			for _, char2 := range association.Items {
				if char1 != char2 {
					values = append(values, fmt.Sprintf("(%d,%d,%d,UTC_TIMESTAMP(), %d)",
						char1, char2, association.Frequency, associationType))
				}
			}
		}
	}

	log.Printf("Character Associations: Update Database")
	for start := 0; start < len(values); start = start + 20000 {
		//	log.Printf("Character Associations: %.2f%%\n", (float32(start)/float32(len(values)))*100)
		end := min(start+20000, len(values))

		stmt := fmt.Sprintf(`
			INSERT INTO evedata.characterAssociations 
				(characterID, associateID, frequency, added, source) 
				VALUES %s 
        	ON DUPLICATE KEY UPDATE 
				added = VALUES(added),
				frequency = IF(frequency > VALUES(frequency), frequency, VALUES(frequency)),
				source = source | VALUES(source);
			`, strings.Join(values[start:end], ",\n"))

		tx, err := Begin()
		if err != nil {
			return err
		}
		_, err = tx.Exec(stmt)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = RetryTransaction(tx)
		if err != nil {
			return err
		}
	}
	log.Printf("Character Associations: Finished")
	return nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func splitToInt(list string) []int {
	a := strings.Split(list, ",")
	b := make([]int, len(a))
	for i, v := range a {
		b[i], _ = strconv.Atoi(v)
	}
	return b
}
