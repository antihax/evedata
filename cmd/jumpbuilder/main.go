package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/RyanCarrier/dijkstra"
	"github.com/antihax/evedata/internal/sqlhelper"
)

type system struct {
	SystemID  int
	Neighbors []int
	Security  float64
}

type pair struct {
	to   system
	from system
}

type path struct {
	to          int
	from        int
	jumps       int64
	secureJumps int64
}

func Round(x, unit float64) float64 {
	return float64(int64(x/unit+0.5)) * unit
}

func dbUpdater(p chan path) {
	db := sqlhelper.NewDatabase()
	for {
		e := <-p

		_, err := db.Exec(`
		INSERT INTO evedata.jumps (toSolarSystemID, fromSolarSystemID, jumps, securejumps)
		VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE jumps=VALUES(jumps), securejumps=VALUES(securejumps)`,
			e.to, e.from, e.jumps, e.secureJumps)
		fmt.Printf("%v %s\n", e, err)
	}
}

func processor(systems []system, in chan pair, out chan path) {
	fmt.Printf("build graphs\n")
	secureGraph := dijkstra.NewGraph()
	graph := dijkstra.NewGraph()

	fmt.Printf("build vertices\n")
	for _, s := range systems {
		if Round(s.Security, 0.1) >= 0.5 {
			secureGraph.AddVertex(s.SystemID)
		}
		graph.AddVertex(s.SystemID)
	}

	fmt.Printf("build arcs\n")
	for _, s := range systems {
		for _, n := range s.Neighbors {
			if Round(s.Security, 0.1) >= 0.5 {
				secureGraph.AddArc(s.SystemID, n, 1)
			}
			graph.AddArc(s.SystemID, n, 1)
		}
	}

	for {
		pair := <-in
		to := pair.to
		from := pair.from
		s := path{to: to.SystemID, from: from.SystemID, jumps: 9999, secureJumps: 9999}
		jumps, err := graph.Shortest(to.SystemID, from.SystemID)
		if err != nil {
			s.jumps = 9999
		} else {
			s.jumps = jumps.Distance
		}

		if Round(to.Security, 0.1) >= 0.5 && Round(from.Security, 0.1) >= 0.5 {
			jumps, err := graph.Shortest(to.SystemID, from.SystemID)
			if err != nil {
				s.secureJumps = 9999
			} else {
				s.secureJumps = jumps.Distance
			}
		}
		out <- s
	}
}

// Add any new refTypes into the database
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("evedata journal import: ")

	fmt.Printf("load data\n")
	raw, _ := ioutil.ReadFile("./jumpmap.json")
	systems := []system{}
	if err := json.Unmarshal(raw, &systems); err != nil {
		log.Panicln(err)
	}

	c := make(chan path, 1000)
	pairs := make(chan pair)
	go dbUpdater(c)

	for i := 0; i < 5; i++ {
		go processor(systems, pairs, c)
	}

	fmt.Printf("build paths\n")
	for _, from := range systems {
		for _, to := range systems {
			p := pair{to: to, from: from}
			pairs <- p
		}
	}
}
