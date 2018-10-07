package artifice

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/antihax/evedata/internal/fpgrowth"
)

func init() {
	registerTrigger("buildRelationships", buildRelationships, time.NewTicker(time.Second*86400))
}

func buildRelationships(s *Artifice) error {
	if err := s.cleanupRelationships(); err != nil {
		log.Println(err)
		return err
	}

	if err := s.buildKillmailRelationships(); err != nil {
		log.Println(err)
		return err
	}
	if err := s.buildCorpJoinRelationships(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Find relationships between characters in killmails
func (s *Artifice) cleanupRelationships() error {
	// Remove any orphan killmails
	return s.doSQL(`
		DELETE FROM evedata.characterAssociations 
		WHERE
			characterID = 0
			OR added < DATE_SUB(UTC_TIMESTAMP(),
			INTERVAL 6 MONTH)`)
}

// Find relationships between characters in killmails
func (s *Artifice) buildKillmailRelationships() error {
	rows, err := s.db.Query(`
        SELECT K.id, GROUP_CONCAT(characterID) 
        FROM evedata.killmailAttackers A
        INNER JOIN evedata.killmails K ON K.id = A.id
        WHERE killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 90 DAY) AND characterID > 0 
        GROUP BY K.id
        HAVING count(*) > 1 AND count(*) < 11;
        `)
	if err != nil {
		return err
	}
	defer rows.Close()

	return s.findRelationships(rows, 1)
}

// Find relationships between characters from corp history
func (s *Artifice) buildCorpJoinRelationships() error {
	rows, err := s.db.Query(`
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
	defer rows.Close()
	return s.findRelationships(rows, 2)
}

// Find relationships in the results
func (s *Artifice) findRelationships(rows *sql.Rows, associationType uint8) error {
	log.Printf("Character Associations: Build Transaction History")
	transactions := fpgrowth.ItemSet{}

	for rows.Next() {
		var (
			transactionID int
			items         string
		)

		// Build history of common transactions
		rows.Scan(&transactionID, &items)
		transactions[transactionID] = SplitToInt(items)
	}
	rows.Close()

	log.Printf("Character Associations: Build fpTree")
	fp := fpgrowth.NewFPTree(transactions, 2)

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

		err := s.doSQL(stmt)
		if err != nil {
			return err
		}
	}
	log.Printf("Character Associations: Finished")
	return nil
}

// SplitToInt splits a csv into []int
func SplitToInt(list string) []int {
	a := strings.Split(list, ",")
	b := make([]int, len(a))
	for i, v := range a {
		b[i], _ = strconv.Atoi(v)
	}
	return b
}
