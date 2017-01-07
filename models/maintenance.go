package models

func MaintKillMails() error { // Broken into smaller chunks so we have a chance of it getting completed.
	// Delete stuff older than a year, we do not care...
	if err := retryExec(`
		DELETE A.* FROM evedata.killmailAttackers A 
            INNER JOIN evedata.killmails K ON A.id = K.id
            WHERE killTime < DATE_SUB(UTC_TIMESTAMP, INTERVAL 1 YEAR); 
            `); err != nil {
		return err
	}
	if err := retryExec(`
		DELETE A.* FROM evedata.killmailItems A 
        INNER JOIN evedata.killmails K ON A.id = K.id
        WHERE killTime < DATE_SUB(UTC_TIMESTAMP, INTERVAL 1 YEAR); 
            `); err != nil {
		return err
	}
	if err := retryExec(`
		DELETE FROM evedata.killmails
        WHERE killTime < DATE_SUB(UTC_TIMESTAMP, INTERVAL 1 YEAR);
            `); err != nil {
		return err
	}

	// Remove any invalid items
	if err := retryExec(`
        DELETE A.* FROM evedata.killmailAttackers A
        LEFT OUTER JOIN evedata.killmails K ON A.id = K.id
        WHERE K.id IS NULL;
            `); err != nil {
		return err
	}
	if err := retryExec(`
        DELETE A.* FROM evedata.killmailItems A
        LEFT OUTER JOIN evedata.killmails K ON A.id = K.id
        WHERE K.id IS NULL;
            `); err != nil {
		return err
	}

	// Prefill stats for known entities that may have no kills
	if err := retryExec(`
        INSERT IGNORE INTO evedata.entityKillStats (id)
	    (SELECT corporationID AS id FROM evedata.corporations); 
            `); err != nil {
		return err
	}
	if err := retryExec(`
        INSERT IGNORE INTO evedata.entityKillStats (id)
	    (SELECT allianceID AS id FROM evedata.alliances); 
            `); err != nil {
		return err
	}

	// Build entity stats
	if err := retryExec(`
        INSERT IGNORE INTO evedata.entityKillStats (id, losses)
            (SELECT 
                victimCorporationID AS id,
                COUNT(DISTINCT K.id) AS losses
            FROM evedata.killmails K
            WHERE K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 30 DAY)
            GROUP BY victimCorporationID
            ) ON DUPLICATE KEY UPDATE losses = values(losses);
            `); err != nil {
		return err
	}
	if err := retryExec(`
        INSERT IGNORE INTO evedata.entityKillStats (id, losses)
            (SELECT 
                victimAllianceID AS id,
                COUNT(DISTINCT K.id) AS losses
            FROM evedata.killmails K
            WHERE K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 30 DAY)
            GROUP BY victimAllianceID
            ) ON DUPLICATE KEY UPDATE losses = values(losses);
            `); err != nil {
		return err
	}

	if err := retryExec(`
        INSERT IGNORE INTO evedata.entityKillStats (id, kills)
            (SELECT 
                corporationID AS id,
                COUNT(DISTINCT K.id) AS kills
            FROM evedata.killmails K
            INNER JOIN evedata.killmailAttackers A ON A.id = K.id
            WHERE K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 30 DAY)
            GROUP BY A.corporationID
            ) ON DUPLICATE KEY UPDATE kills = values(kills);
            `); err != nil {
		return err
	}
	if err := retryExec(`
        INSERT IGNORE INTO evedata.entityKillStats (id, kills)
            (SELECT 
                allianceID AS id,
                COUNT(DISTINCT K.id) AS kills
            FROM evedata.killmails K
            INNER JOIN evedata.killmailAttackers A ON A.id = K.id
            WHERE K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 30 DAY)
            GROUP BY A.allianceID
            ) ON DUPLICATE KEY UPDATE kills = values(kills);
            `); err != nil {
		return err
	}

	// Update everyone efficiency
	if err := retryExec(`
        UPDATE evedata.entityKillStats SET efficiency = IF(losses+kills, (kills/(kills+losses)) , 1.0000);
            `); err != nil {
		return err
	}

	return nil
}

// Retry the exec until we get no error (deadlocks)
func retryExec(sql string, args ...interface{}) error {
	var err error
	for {
		_, err = database.Exec(sql, args...)
		if err != nil {
			break
		}
	}
	return err
}
