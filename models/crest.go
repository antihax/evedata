package models

import "strings"

func AddCRESTRef(id int64, ref string) error {
	var t string
	if strings.Contains(ref, "alliances") {
		t = "alliance"
	} else if strings.Contains(ref, "corporations") {
		t = "corporation"
	} else if strings.Contains(ref, "characters") {
		t = "character"
	}

	_, err := database.Exec(`INSERT IGNORE INTO crestID (id, crestRef, type) VALUES(?,?,?);`, id, ref, t)
	if err != nil {

		return err
	}
	return nil
}
