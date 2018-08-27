package artifice

import (
	"log"
	"time"
)

func init() {
	registerTrigger("alliancehistoryMaint", alliancehistoryMaint, time.NewTicker(time.Hour*2))
	registerTrigger("corphistoryMaint", corphistoryMaint, time.NewTicker(time.Hour*2))
	registerTrigger("marketMaint", marketMaint, time.NewTicker(time.Hour))
	registerTrigger("marketUpdate", marketUpdate, time.NewTicker(time.Hour*2))
	registerTrigger("discoveredAssetsMaint", discoveredAssetsMaint, time.NewTicker(time.Second*3620))
	registerTrigger("entityMaint", entityMaint, time.NewTicker(time.Second*3630*3))
	registerTrigger("killmailMaint", killmailMaint, time.NewTicker(time.Hour*3))
	registerTrigger("entityStatsMaint", entityStatsMaint, time.NewTicker(time.Hour*8))
	registerTrigger("contactSyncMaint", contactSyncMaint, time.NewTicker(time.Second*3615*6))
}

func corphistoryMaint(s *Artifice) error {
	if err := s.doSQL(`
		UPDATE evedata.corporationHistory A
		INNER JOIN (
			SELECT H.recordID, MIN(L.startDate) AS endDate
				FROM 	evedata.corporationHistory H
				INNER JOIN evedata.corporationHistory L ON H.characterID = L.characterID AND L.startDate > H.startDate
				WHERE   H.endDate IS NULL
				GROUP BY H.recordID) B ON A.recordID = B.recordID
			SET A.endDate = B.endDate`); err != nil {
		return err
	}
	return nil
}

func alliancehistoryMaint(s *Artifice) error {
	if err := s.doSQL(`
		UPDATE evedata.allianceHistory A
		INNER JOIN (SELECT H.recordID, MIN(L.startDate) AS endDate
			FROM evedata.allianceHistory H
			INNER JOIN evedata.allianceHistory L ON H.corporationID = L.corporationID AND L.startDate > H.startDate
			WHERE H.endDate IS NULL
			GROUP BY H.recordID) B ON A.recordID = B.recordID
		SET A.endDate = B.endDate`); err != nil {
		return err
	}

	return nil
}

func contactSyncMaint(s *Artifice) error {
	if err := s.doSQL(`
        DELETE S.* FROM evedata.contactSyncs S
        LEFT OUTER JOIN evedata.crestTokens T ON S.destination = T.tokenCharacterID
        WHERE tokenCharacterID IS NULL;`); err != nil {
		return err
	}
	if err := s.doSQL(`
        DELETE S.* FROM evedata.contactSyncs S
        LEFT OUTER JOIN evedata.crestTokens T ON S.source = T.tokenCharacterID
        WHERE tokenCharacterID IS NULL;`); err != nil {
		return err
	}

	return nil
}

func entityStatsMaint(s *Artifice) error {
	// Prefill stats for known entities that may have no kills
	if err := s.doSQL(`
			INSERT IGNORE INTO evedata.entityKillStats (id) (SELECT characterID AS id FROM evedata.characters);
				`); err != nil {
		return err
	}
	if err := s.doSQL(`
			INSERT IGNORE INTO evedata.entityKillStats (id) (SELECT corporationID AS id FROM evedata.corporations WHERE memberCount > 0);
				`); err != nil {
		return err
	}
	if err := s.doSQL(`
			INSERT IGNORE INTO evedata.entityKillStats (id) (SELECT allianceID AS id FROM evedata.alliances);
				`); err != nil {
		return err
	}
	if err := s.doSQL(`
		UPDATE evedata.entityKillStats SET kills=0, losses=0, capKills=0, efficiency = 0;
			`); err != nil {
		return err
	}

	// Build entity stats
	if err := s.doSQL(`
		INSERT INTO evedata.entityKillStats (id, losses)
			(SELECT 
				victimCharacterID AS id,
				COUNT(DISTINCT K.id) AS losses
			FROM evedata.killmails K
			WHERE K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 90 DAY)
			GROUP BY victimCharacterID
			) ON DUPLICATE KEY UPDATE losses = values(losses);
			`); err != nil {
		return err
	}

	// Build entity stats
	if err := s.doSQL(`
			INSERT INTO evedata.entityKillStats (id, losses)
				(SELECT 
					victimCorporationID AS id,
					COUNT(DISTINCT K.id) AS losses
				FROM evedata.killmails K
				WHERE K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 90 DAY)
				GROUP BY victimCorporationID
				) ON DUPLICATE KEY UPDATE losses = values(losses);
				`); err != nil {
		return err
	}

	if err := s.doSQL(`
			INSERT INTO evedata.entityKillStats (id, losses)
				(SELECT 
					victimAllianceID AS id,
					COUNT(DISTINCT K.id) AS losses
				FROM evedata.killmails K
				WHERE K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 90 DAY)
				GROUP BY victimAllianceID
				) ON DUPLICATE KEY UPDATE losses = values(losses);
				`); err != nil {
		return err
	}

	if err := s.doSQL(`
		INSERT INTO evedata.entityKillStats (id, kills)
			(SELECT 
				characterID AS id,
				COUNT(DISTINCT K.id) AS kills
			FROM evedata.killmails K
			INNER JOIN evedata.killmailAttackers A ON A.id = K.id
			WHERE K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 90 DAY)
			GROUP BY A.characterID
			) ON DUPLICATE KEY UPDATE kills = values(kills);
			`); err != nil {
		return err
	}

	if err := s.doSQL(`
			INSERT INTO evedata.entityKillStats (id, kills)
				(SELECT 
					corporationID AS id,
					COUNT(DISTINCT K.id) AS kills
				FROM evedata.killmails K
				INNER JOIN evedata.killmailAttackers A ON A.id = K.id
				WHERE K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 90 DAY)
				GROUP BY A.corporationID
				) ON DUPLICATE KEY UPDATE kills = values(kills);
				`); err != nil {
		return err
	}

	if err := s.doSQL(`
			INSERT INTO evedata.entityKillStats (id, kills)
				(SELECT 
					allianceID AS id,
					COUNT(DISTINCT K.id) AS kills
				FROM evedata.killmails K
				INNER JOIN evedata.killmailAttackers A ON A.id = K.id
				WHERE K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 90 DAY)
				GROUP BY A.allianceID
				) ON DUPLICATE KEY UPDATE kills = values(kills);
				`); err != nil {
		return err
	}

	if err := s.doSQL(`
		INSERT INTO evedata.entityKillStats (id, capKills) (
			SELECT C.characterID, COUNT(DISTINCT A.id) FROM evedata.characters C
			INNER JOIN evedata.killmailAttackers A ON C.characterID = A.characterID
			INNER JOIN evedata.killmails K ON A.id = K.id AND K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 90 DAY)
            INNER JOIN evedata.killmailAttackers A2 ON A.id = A2.id
				AND A2.shipType IN (SELECT typeID FROM eve.invTypes WHERE groupID IN (30, 547, 485, 1538, 659))
			GROUP BY characterID) ON DUPLICATE KEY UPDATE capKills = values(capKills);	
			`); err != nil {
		return err
	}

	if err := s.doSQL(`
		INSERT INTO evedata.entityKillStats (id, capKills) (
			SELECT C.corporationID, COUNT(DISTINCT A.id) FROM evedata.corporations C
			INNER JOIN evedata.killmailAttackers A ON C.corporationID = A.corporationID
			INNER JOIN evedata.killmails K ON A.id = K.id AND K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 90 DAY)
            INNER JOIN evedata.killmailAttackers A2 ON A.id = A2.id
				AND A2.shipType IN (SELECT typeID FROM eve.invTypes WHERE groupID IN (30, 547, 485, 1538, 659))
			GROUP BY corporationID) ON DUPLICATE KEY UPDATE capKills = values(capKills);	
			`); err != nil {
		return err
	}

	if err := s.doSQL(`
		INSERT INTO evedata.entityKillStats (id, capKills) (
			SELECT C.allianceID, COUNT(DISTINCT A.id) FROM evedata.alliances C
			INNER JOIN evedata.killmailAttackers A ON C.allianceID = A.allianceID
			INNER JOIN evedata.killmails K ON A.id = K.id AND K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 90 DAY)
            INNER JOIN evedata.killmailAttackers A2 ON A.id = A2.id
				AND A2.shipType IN (SELECT typeID FROM eve.invTypes WHERE groupID IN (30, 547, 485, 1538, 659))
			GROUP BY allianceID) ON DUPLICATE KEY UPDATE capKills = values(capKills);	
			`); err != nil {
		return err
	}

	// Update everyone efficiency
	if err := s.doSQL(`
			UPDATE evedata.entityKillStats SET efficiency = IF(losses+kills, (kills/(kills+losses)), 1.0000);
				`); err != nil {
		return err
	}
	return nil

}

func killmailMaint(s *Artifice) error { // Broken into smaller chunks so we have a chance of it getting completed.
	// Delete old killmails
	if err := s.RetryExecTillNoRows(`
		DELETE FROM evedata.killmails
		WHERE killTime < DATE_SUB(UTC_TIMESTAMP, INTERVAL 365 DAY) LIMIT 8000;
			`); err != nil {
		return err
	}

	// Remove any orphan attackers
	if err := s.RetryExecTillNoRows(`
		DELETE D.* FROM evedata.killmailAttackers D
		JOIN (select A.id FROM evedata.killmailAttackers A
			LEFT JOIN evedata.killmails K ON K.id = A.id
			WHERE K.id IS NULL LIMIT 500) S ON D.id = S.id;
				   `); err != nil {
		return err
	}

	// Remove any orphan killmails
	if err := s.RetryExecTillNoRows(`
			DELETE D.* FROM evedata.killmails D
			JOIN (select K.id FROM evedata.killmails K
				LEFT JOIN evedata.killmailAttackers A ON A.id = K.id
				WHERE A.id IS NULL LIMIT 1000) S ON D.id = S.id;
					   `); err != nil {
		return err
	}

	return nil
}

func discoveredAssetsMaint(s *Artifice) error {
	if err := s.doSQL(`
        INSERT INTO evedata.discoveredAssets 
            SELECT 
                A.corporationID, 
                C.allianceID, 
                typeID, 
                K.solarSystemID, 
                K.x, 
                K.y, 
                K.z, 
                evedata.closestCelestial(K.solarSystemID, K.x, K.y, K.z) AS locationID, 
                MAX(killTime) as lastSeen 
            FROM evedata.killmailAttackers A
            INNER JOIN invTypes T ON shipType = typeID
            INNER JOIN evedata.corporations C ON C.corporationID = A.corporationID
            INNER JOIN evedata.killmails K ON K.id = A.id
            INNER JOIN mapSolarSystems S ON S.solarSystemID = K.solarSystemID
			WHERE K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 120 DAY) AND 
				characterID = 0 AND groupID IN (
					365, 549, 1023, 1404, 1406, 1537, 1652, 1653, 1657, 2233,
					1321,1322,1327,1328,1329,1330,1331,1332,1333,1415,1429,1430,
					1441,1442,1535,1537,1546,1547,1548,1549,1551,1562,1613,1614,
					1615,1616,1617,1618,1619,1620,1621,1622,1629,1630,1631,1632,
					1633,1634,1635,1639,1640,1641,1642,1652,1653,1717,1719,1816,
					1819,1820,1821,1822,1823,1824,1825,1826,1827,1828,1829,1830,
					1831,1832,1833,1834,1835,1836,1837,1838,1839,1840,1841,1842,
					1843,1844,1845,1846,1847,1850,1851,1852,1853,1854,1855,1856,
					1857,1858,1859,1860,1861,1862,1863,1864,1865,1867,1868,1869,
					1870,1887,1912,1913,1914,1933,1934,1935,1936,1937,1938,1939,
					1941,1942,1943,1944,1945,1962,1966,1967,1968)
            GROUP BY A.corporationID, solarSystemID, typeID
        ON DUPLICATE KEY UPDATE lastSeen =  VALUES(lastSeen);
            `); err != nil {
		return err
	}

	if err := s.doSQL(`
        INSERT INTO evedata.discoveredAssets 
            SELECT 
                K.victimCorporationID AS corporationID, 
                C.allianceID, 
                typeID, 
                K.solarSystemID, 
                K.x, 
                K.y, 
                K.z, 
                evedata.closestCelestial(K.solarSystemID, K.x, K.y, K.z) AS locationID, 
                MAX(killTime) as lastSeen 
            FROM evedata.killmails K 
            INNER JOIN invTypes T ON K.shipType = typeID
            INNER JOIN evedata.corporations C ON C.corporationID = K.victimCorporationID
            INNER JOIN mapSolarSystems S ON S.solarSystemID = K.solarSystemID
			WHERE K.killTime > DATE_SUB(UTC_TIMESTAMP, INTERVAL 120 DAY) AND 
				victimCharacterID = 0 AND groupID IN (
					365, 549, 1023, 1404, 1406, 1537, 1652, 1653, 1657, 2233,
					1321,1322,1327,1328,1329,1330,1331,1332,1333,1415,1429,1430,
					1441,1442,1535,1537,1546,1547,1548,1549,1551,1562,1613,1614,
					1615,1616,1617,1618,1619,1620,1621,1622,1629,1630,1631,1632,
					1633,1634,1635,1639,1640,1641,1642,1652,1653,1717,1719,1816,
					1819,1820,1821,1822,1823,1824,1825,1826,1827,1828,1829,1830,
					1831,1832,1833,1834,1835,1836,1837,1838,1839,1840,1841,1842,
					1843,1844,1845,1846,1847,1850,1851,1852,1853,1854,1855,1856,
					1857,1858,1859,1860,1861,1862,1863,1864,1865,1867,1868,1869,
					1870,1887,1912,1913,1914,1933,1934,1935,1936,1937,1938,1939,
					1941,1942,1943,1944,1945,1962,1966,1967,1968)
            GROUP BY K.victimCorporationID, solarSystemID, typeID
        ON DUPLICATE KEY UPDATE lastSeen =  VALUES(lastSeen);`); err != nil {
		return err
	}

	if err := s.doSQL(`
		INSERT INTO evedata.discoveredAssets 
			SELECT corporationID, allianceID, typeID, solarSystemID, x, y, z,
			evedata.closestCelestial(solarSystemID, x, y, z) AS locationID,
			S.updated AS lastSeen
			FROM evedata.structures S
			INNER JOIN evedata.corporations C ON C.corporationID = S.ownerID
		ON DUPLICATE KEY UPDATE lastSeen = VALUES(lastSeen);`); err != nil {
		return err
	}

	if err := s.doSQL(`
		DELETE FROM evedata.discoveredAssets 
		WHERE lastSeen < DATE_SUB(UTC_TIMESTAMP(), INTERVAL 1 YEAR);`); err != nil {
		return err
	}
	return nil
}

func entityMaint(s *Artifice) error {
	if err := s.doSQL(`
        UPDATE evedata.alliances A SET memberCount = 
            IFNULL(
                    (SELECT sum(memberCount) AS memberCount FROM evedata.corporations  C
                    WHERE C.allianceID = A.allianceID
                    GROUP BY allianceID LIMIT 1),
                    0
            );
            `); err != nil {
		return err
	}
	return nil
}

type marketRegion struct {
	RegionID   int32  `db:"regionID"`
	RegionName string `db:"regionName"`
}

// Anywhere can now have a public market.
func getMarketRegions(s *Artifice) ([]marketRegion, error) {
	v := []marketRegion{}
	err := s.db.Select(&v, `
		SELECT 	regionID, regionName 
		FROM 	mapRegions 
		WHERE regionID < 11000000;
	`)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func marketUpdate(s *Artifice) error {
	regions, err := getMarketRegions(s)
	if err != nil {
		log.Println(err)
	}

	if err := s.doSQL(`
        DELETE FROM evedata.marketStations ORDER BY stationID;
             `); err != nil {
		log.Println(err)
	}

	if err := s.doSQL(`
        INSERT INTO evedata.marketStations SELECT stationName, M.stationID, Count(*) as Count
        FROM    evedata.market M
                INNER JOIN staStations S ON M.stationID = S.stationID
        WHERE   reported >= DATE_SUB(UTC_TIMESTAMP(), INTERVAL 7 DAY)
		GROUP BY M.stationID 
		ORDER BY stationID
        ON DUPLICATE KEY UPDATE stationID=stationID;
            `); err != nil {
		log.Println(err)
	}

	if err := s.doSQL(`
       UPDATE evedata.market_vol SET quantity = 0;
             `); err != nil {
		log.Println(err)
	}

	for _, region := range regions {
		if err := s.doSQL(`
        REPLACE INTO evedata.market_vol (
            SELECT count(*) as number,sum(quantity)/7 as quantity, regionID, itemID 
                FROM evedata.market_history 
                WHERE date > DATE_SUB(UTC_TIMESTAMP(),INTERVAL 7 DAY) 
                AND regionID = ?
                GROUP BY regionID, itemID);
            `, region.RegionID); err != nil {
			log.Println(err)
			return err
		}
	}

	if err := s.doSQL(`
		DELETE FROM evedata.jitaPrice ORDER BY itemID;
			  `); err != nil {
		log.Println(err)
		return err
	}

	if err := s.doSQL(`
		 INSERT IGNORE INTO evedata.jitaPrice (
		 SELECT S.typeID as itemID, buy, sell, high, low, mean, quantity FROM
			 (SELECT typeID, min(price) AS sell FROM evedata.market WHERE regionID = 10000002 AND bid = 0 GROUP BY typeID) S
			 INNER JOIN (SELECT typeID, max(price) AS buy FROM evedata.market WHERE regionID = 10000002 AND bid = 1 GROUP BY typeID) B ON S.typeID = B.typeID
			 LEFT OUTER JOIN (SELECT itemID, max(high) AS high, avg(mean) AS mean, min(low) AS low, sum(quantity) AS quantity FROM evedata.market_history WHERE regionID = 10000002 AND date > DATE_SUB(UTC_DATE(), INTERVAL 7 DAY) GROUP BY itemID) H on H.itemID = S.typeID
		 HAVING mean IS NOT NULL
		 ) ORDER BY itemID;
			 `); err != nil {
		log.Println(err)
		return err
	}

	if err := s.doSQL(`
		DELETE FROM evedata.iskPerLp ORDER BY typeID;
			  `); err != nil {
		log.Println(err)
		return err
	}

	if err := s.doSQL(`
		 INSERT IGNORE INTO evedata.iskPerLp (
		 SELECT
				 N.itemName,
				 S.typeID,
				 T.typeName,
				 MIN(lpCost) AS lpCost,
				 MIN(iskCost) AS iskCost,
				 ROUND(MIN(C.buy),0) AS JitaPrice,
				 ROUND(MIN(C.quantity),0) AS JitaVolume,
				 ROUND(COALESCE(MIN(P.price),0) + iskCost, 0)  AS itemCost,
				 ROUND(
						 (
								 ( MIN(S.quantity) * AVG(C.buy) ) -
								 ( COALESCE( MIN(P.price), 0) + iskCost )
						 )
						 / MIN(lpCost)
				 , 0) AS ISKperLP,
				 P.offerID
		 FROM evedata.lpOffers S
 
		 INNER JOIN invNames N ON S.corporationID = N.itemID
		 INNER JOIN invTypes T ON S.typeID = T.typeID
		 INNER JOIN evedata.jitaPrice C ON C.itemID = S.typeID
 
		 LEFT OUTER JOIN         (
								 SELECT offerID, sum(H.sell * L.quantity) AS price
								 FROM evedata.lpOfferRequirements L
								 INNER JOIN evedata.jitaPrice H ON H.itemID = L.typeID
								 GROUP BY offerID
						 ) AS P ON S.offerID = P.offerID
 
		 GROUP BY S.offerID, S.corporationID
		 HAVING ISKperLP > 0) ORDER BY typeID;
			 `); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func marketMaint(s *Artifice) error {
	// Deal with any possible orphaned orders
	if err := s.RetryExecTillNoRows(`
	DELETE FROM evedata.market WHERE DATE_ADD(issued, INTERVAL duration DAY) < utc_timestamp();
	            `); err != nil {
		log.Println(err)
	}

	return nil
}
