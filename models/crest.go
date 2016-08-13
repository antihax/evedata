package models

func AddCRESTRef(id int, ref string) error {
	_, err := database.Exec(`INSERT IGNORE INTO crestID (id, crestRef) VALUES(?,?);`, id, ref)
	if err != nil {

		return err
	}
	return nil
}
