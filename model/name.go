package model

import "log"

func QueryNameWithID(id int) ( name string, err error) {
	stm, err := db.Prepare("SELECT * FROM customer WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}

	defer stm.Close()

	var tempid int

	err = stm.QueryRow(id).Scan(&name, &tempid)

	if err != nil {
		log.Fatal(err)
	}

	return

} // QueryNameWithID()