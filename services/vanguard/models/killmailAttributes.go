package models

import "fmt"

type KillmailAttributes struct {
	ID                   int32   `db:"id" json:"id"`
	TypeName             string  `db:"typeName" json:"typeName"`
	TypeID               string  `db:"typeID" json:"typeID"`
	RPS                  float64 `db:"rps" json:"rps"`
	DPS                  float64 `db:"dps" json:"dps"`
	EHP                  float64 `db:"ehp" json:"ehp"`
	Alpha                float64 `db:"alpha" json:"alpha"`
	ScanResolution       float64 `db:"scanResolution" json:"scanResolution"`
	SignatureRadiusNoMWD float64 `db:"signatureRadiusNoMWD" json:"signatureRadiusNoMWD"`
	Agility              float64 `db:"agility" json:"agility"`
	WarpSpeed            float64 `db:"warpSpeed" json:"warpSpeed"`
	Speed                float64 `db:"speed" json:"speed"`
	RemoteArmorRepair    float64 `db:"remoteArmorRepair" json:"remoteArmorRepair"`
	RemoteShieldRepair   float64 `db:"remoteShieldRepair" json:"remoteShieldRepair"`
	RemoteEnergyTransfer float64 `db:"remoteEnergyTransfer" json:"remoteEnergyTransfer"`
	EnergyNeutralization float64 `db:"energyNeutralization" json:"energyNeutralization"`
	SensorStrength       float64 `db:"sensorStrength" json:"sensorStrength"`
	CapacitorNoMWD       float64 `db:"capacitorNoMWD" json:"capacitorNoMWD"`
	CapacitorTimeNoMWD   float64 `db:"capacitorTimeNoMWD" json:"capacitorTimeNoMWD"`
}

func GetKillmailAttributes(groupID int64, value int64, points int64) ([]KillmailAttributes, error) {
	v := []KillmailAttributes{}

	otherFilters := ""
	if value > 0 {
		otherFilters += fmt.Sprintf(" AND fittedValue <= %d", value*1000000)
	}
	if value > 0 {
		otherFilters += fmt.Sprintf(" AND totalWarpScrambleStrength >= %d", points)
	}

	if err := database.Select(&v, `
		SELECT 	K.id, typeName, typeID, rps, dps, ehp, alpha, scanResolution, signatureRadiusNoMWD, agility, 
			warpSpeed, speed, remoteArmorRepair, remoteShieldRepair, remoteEnergyTransfer,
			energyNeutralization, sensorStrength, capacitorNoMWD, capacitorTimeNoMWD
		FROM evedata.killmails K
		INNER JOIN evedata.killmailAttributes A FORCE INDEX(ix_id_cpu_pg_ehp) ON K.id = A.id 
		INNER JOIN invTypes T ON T.typeID = K.shipType
		WHERE T.groupID = ? AND powerRemaining >= 0 AND CPURemaining >= 0 AND eHP > 0 AND (alpha > 0 OR rps > 0 OR remoteArmorRepair > 0 OR remoteShieldRepair > 0)
		AND K.killTime > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 3 MONTH)`+otherFilters+` ORDER BY killTime DESC LIMIT 10000;`, groupID); err != nil {
		return nil, err
	}
	return v, nil
}

type OffensiveGroups struct {
	GroupID   int32  `db:"groupID" json:"groupID"`
	GroupName string `db:"groupName" json:"groupName"`
}

func GetOffensiveShipGroupID() ([]OffensiveGroups, error) {
	v := []OffensiveGroups{}

	if err := database.Select(&v, `
		SELECT groupID, groupName FROM eve.invGroups 
		WHERE categoryID = 6 AND groupID NOT IN(29, 902, 31, 30, 547, 659, 1972, 513, 1202, 381, 513, 1022)
		ORDER BY groupName;`); err != nil {
		return nil, err
	}
	return v, nil
}

type KillmailStatistics struct {
	Month        int `db:"month" json:"month"`
	Year         int `db:"year" json:"year"`
	CharacterAge int `db:"characterAge" json:"characterAge"`
	Wars         int `db:"wars" json:"wars"`
	Ganks        int `db:"ganks" json:"ganks"`
	Lowsec       int `db:"lowsec" json:"lowsec"`
	Nullsec      int `db:"nullsec" json:"nullsec"`
	Highsec      int `db:"highsec" json:"highsec"`
	NPCKills     int `db:"npcKills" json:"npcKills"`
	WH           int `db:"wh" json:"wh"`
	LowsecFW     int `db:"lowsecFW" json:"lowsecFW"`
	HighsecFW    int `db:"highsecFW" json:"highsecFW"`
	Total        int `db:"total" json:"total"`
}

func GetKillmailStatistics() ([]KillmailStatistics, error) {
	v := []KillmailStatistics{}

	if err := database.Select(&v, `
		SELECT month, year, characterAge, wars, ganks, lowsec,
		nullsec, wh, lowsecFW, highsecFW, highsec, npcKills, total 
		FROM evedata.killmailStatistics;	
	`); err != nil {
		return nil, err
	}
	return v, nil
}

type KillmailAreaEntityStatistics struct {
	Area  string `db:"area" json:"area"`
	Name  string `db:"name" json:"name"`
	ID    int32  `db:"id" json:"id"`
	Kills int32  `db:"kills" json:"kills"`
}

func GetKillmailAreaEntityStatistics() ([]KillmailAreaEntityStatistics, error) {
	v := []KillmailAreaEntityStatistics{}

	if err := database.Select(&v, `
	SELECT id, sum(kills) AS kills, area, coalesce(A.name,C.name) as name
	FROM evedata.killmailKillers K
	LEFT OUTER JOIN evedata.alliances A ON A.allianceID = K.id
	LEFT OUTER JOIN evedata.corporations C ON C.corporationID = K.id
    WHERE year > YEAR(DATE_SUB(utc_timestamp(), INTERVAL 2 year))
	GROUP BY id, area
	`); err != nil {
		return nil, err
	}
	return v, nil
}

type KillmailCapVsSubcapStatistics struct {
	Month int    `db:"month" json:"month"`
	Year  int    `db:"year" json:"year"`
	Area  string `db:"area" json:"area"`
	Type  string `db:"type" json:"type"`
	Kills int32  `db:"kills" json:"kills"`
}

func GetKillmailCapVsSubcapStatistics() ([]KillmailCapVsSubcapStatistics, error) {
	v := []KillmailCapVsSubcapStatistics{}

	if err := database.Select(&v, `
	SELECT month, year, kills, area, type
	FROM evedata.killmailCapsVsSubcap K
	`); err != nil {
		return nil, err
	}
	return v, nil
}
