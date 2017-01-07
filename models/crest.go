package models

import "strings"

type CRESTRef struct {
	ID        int64  `db:"id" json:"id"`
	CrestRef  string `db:"crestRef" json:"crestRef"`
	CrestType string `db:"type" json:"type"`
}

func AddCRESTRef(id int64, ref string) error {
	var t string
	if strings.Contains(ref, "alliances") {
		t = "alliance"
	} else if strings.Contains(ref, "corporations") {
		t = "corporation"
	} else if strings.Contains(ref, "characters") {
		t = "character"
	}

	_, err := database.Exec(`INSERT IGNORE INTO evedata.crestID (id, crestRef, type) VALUES(?,?,?);`, id, ref, t)
	if err != nil {

		return err
	}
	return nil
}

// [BENCHMARK] 0.000 sec / 0.000 sec
func GetCRESTRef(id int64) (*CRESTRef, error) {
	ref := &CRESTRef{}
	if err := database.Select(&ref, `SELECT id FROM killmails WHERE id = ? LIMIT 1;`, id); err != nil {
		return nil, err
	}
	return ref, nil
}
