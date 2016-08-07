package eveConsumer

import "log"

func (c *EveConsumer) contactSync() {
	rows, err := c.db.Query(
		`SELECT source, group_concat(destination)
			FROM contactSyncs GROUP BY source
		    HAVING max(nextSync) < UTC_TIMESTAMP()`)
	tx, err := c.db.Beginx()
	if err != nil {
		log.Printf("EVEConsumer: Failed starting transaction: %v", err)
		return
	}

	for rows.Next() {
		var (
			source int
			dest   string
		)

		err = rows.Scan(&source, &dest)
		//destinations := strings.Split(dest, ",")
		if err != nil {
			log.Printf("EVEConsumer: Failed Scanning Rows: %v", err)
			return
		}
		char, err := c.eve.GetCharacterInfo(source)
		if err != nil {
			log.Printf("EVEConsumer: Failed getting character info %v", err)
			continue
		}
		if char.AllianceID > 0 {

		}

	}
	err = tx.Commit()
}
