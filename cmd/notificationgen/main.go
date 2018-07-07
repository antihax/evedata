package main

import (
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"sort"

	"github.com/antihax/evedata/internal/sqlhelper"
)

// Add any new refTypes into the database
func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("evedata notificationgen: ")
	db := sqlhelper.NewDatabase()
	defer db.Close()

	log.Println("Getting Notifications")
	rows, err := db.Query(`
		SELECT DISTINCT type, text FROM evedata.notifications GROUP BY type, length(text)`)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	log.Println("Processing Notifications")
	structs := make(map[string]string)
	for rows.Next() {
		var notifType, notifString string
		err := rows.Scan(&notifType, &notifString)
		if err != nil {
			log.Fatalln(err)
		}
		cmd := exec.Command("nodejs", "yaml-to-go.js", notifType)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}

		io.WriteString(stdin, notifString)
		stdin.Close()

		out, err := cmd.Output()
		if err != nil {
			log.Fatal(err)
		}
		s := string(out)

		if _, ok := structs[notifType]; ok {
			if len(structs[notifType]) < len(s) {
				structs[notifType] = s
			}
		} else {
			structs[notifType] = s
		}
	}

	structures := "package notifications\n\n"

	keys := make([]string, 0)
	for k := range structs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, s := range keys {
		structures += structs[s] + "\n"
	}

	err = ioutil.WriteFile("notifications.txt", []byte(structures), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
