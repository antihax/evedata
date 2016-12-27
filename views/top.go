package views

import (
	"bytes"
	"evedata/appContext"
	evedata "evedata/server"
	"time"

	"evedata/templates"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"text/tabwriter"

	humanize "github.com/dustin/go-humanize"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/sessions"
)

func init() {
	evedata.AddRoute("top", "GET", "/top", topPage)
	evedata.AddRoute("top", "GET", "/X/topStatistics", topStatisticsTxt)
	go GenerateStatistics(evedata.GetContext())
	statisticsLast = make(map[string]int)
}

var (
	statisticsTxt  []byte
	statisticsLast map[string]int
)

func statisticsChange(t string, v int) (out string) {
	l := statisticsLast[t]
	n := v - l
	if n != 0 {
		out = fmt.Sprintf("(%+d)", n)
	}
	statisticsLast[t] = v
	return
}

// [TODO] Break this into a package and fix the stupid.
func GenerateStatistics(c *appContext.AppContext) {
	for c.Cache == nil {
	} // Stupid
	log.Printf("Start collecting statistics\n")
	red := c.Cache.Get()
	tick := time.NewTicker(time.Second * 5)

	for {
		w := bytes.NewBuffer(statisticsTxt)
		w.Reset()

		out := tabwriter.NewWriter(w, 40, 4, 2, ' ', tabwriter.AlignRight)

		kills, _ := redis.Int(red.Do("SCARD", "EVEDATA_knownKills"))
		fmt.Fprintf(out, "%s \tKnown Kills %s\n", humanize.Comma((int64)(kills)), statisticsChange("kills", kills))
		killq, _ := redis.Int(red.Do("SCARD", "EVEDATA_killQueue"))
		fmt.Fprintf(out, "%s \tKills in Queue %s\n", humanize.Comma((int64)(killq)), statisticsChange("killq", killq))

		err := out.Flush()
		statisticsTxt = w.Bytes()
		if err != nil {
			log.Printf("top error: %v\n", err)
		}
		<-tick.C
	}
}

func topPage(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 60*60)
	p := newPage(s, r, "EVEData.org backend statistics")
	templates.Templates = template.Must(template.ParseFiles("templates/top.html", templates.LayoutPath))

	if err := templates.Templates.ExecuteTemplate(w, "base", p); err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func topStatisticsTxt(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
	setCache(w, 5)
	w.Write(statisticsTxt)
	return http.StatusOK, nil
}
