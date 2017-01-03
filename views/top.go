package views

import (
	"bytes"
	"time"

	"github.com/antihax/evedata/appContext"
	evedata "github.com/antihax/evedata/server"

	"fmt"
	"html/template"
	"log"
	"net/http"
	"text/tabwriter"

	"github.com/antihax/evedata/templates"

	"sort"

	humanize "github.com/dustin/go-humanize"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/sessions"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
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

func statisticsLoadHostStats(r redis.Conn) {
	l, _ := load.Avg()
	i, _ := host.Info()
	m, _ := mem.VirtualMemory()
	cpuPercent, _ := cpu.Percent(0, false)
	cpus, _ := cpu.Counts(true)

	data := fmt.Sprintf("%s: Load: %.2f %.2f %.2f  CPU(%d cores): %.1f%%  Memory: %d/%d GiB  ", i.Hostname, l.Load1, l.Load5, l.Load15, cpus, cpuPercent[0], m.Used/1024/1024/1024, m.Total/1024/1024/1024)

	r.Do("ZADD", "EVEDATA_HOST", time.Now().UTC().Unix()+5, data)

	return
}

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
	defer red.Close()
	tick := time.NewTicker(time.Second * 5)

	for {
		statisticsLoadHostStats(red)
		if c.Conf.GenerateStats {
			// Remove old entries from host table
			red.Do("ZREMRANGEBYSCORE", "EVEDATA_HOST", 0, time.Now().UTC().Unix())

			w := bytes.NewBuffer(statisticsTxt)
			w.Reset()

			out := tabwriter.NewWriter(w, 40, 4, 2, ' ', tabwriter.AlignRight)

			// this will store the keys of each iteration
			var host []string

			if arr, err := redis.MultiBulk(red.Do("ZRANGEBYSCORE", "EVEDATA_HOST", 0, "inf")); err != nil {
				fmt.Println(err)
			} else {
				host, _ = redis.Strings(arr, nil)
			}
			sort.Strings(host)
			for i := 0; i < len(host); i++ {
				fmt.Fprintf(out, "%s\n", host[i])
			}

			fmt.Fprintln(out)

			kills, _ := redis.Int(red.Do("SCARD", "EVEDATA_knownKills"))
			fmt.Fprintf(out, "%s \tKnown Kills %s\n", humanize.Comma((int64)(kills)), statisticsChange("kills", kills))
			killq, _ := redis.Int(red.Do("SCARD", "EVEDATA_killQueue"))
			fmt.Fprintf(out, "%s \tKills in Queue %s\n", humanize.Comma((int64)(killq)), statisticsChange("killq", killq))
			entityq, _ := redis.Int(red.Do("SCARD", "EVEDATA_entityQueue"))
			fmt.Fprintf(out, "%s \tEntities in Queue %s\n", humanize.Comma((int64)(entityq)), statisticsChange("entityq", entityq))
			fmt.Fprintln(out)
			history, _ := redis.Int(red.Do("SCARD", "EVEDATA_marketHistory"))
			fmt.Fprintf(out, "%s \tMarket History in Queue %s\n", humanize.Comma((int64)(history)), statisticsChange("history", history))
			orders, _ := redis.Int(red.Do("SCARD", "EVEDATA_marketOrders"))
			fmt.Fprintf(out, "%s \tMarket Orders in Queue %s\n", humanize.Comma((int64)(orders)), statisticsChange("orders", orders))
			regions, _ := redis.Int(red.Do("ZCARD", "EVEDATA_marketRegions"))
			fmt.Fprintf(out, "%s \tMarket Regions %s\n", humanize.Comma((int64)(regions)), statisticsChange("regions", regions))
			contacts, _ := redis.Int(red.Do("SCARD", "EVEDATA_contactSyncQueue"))
			fmt.Fprintf(out, "%s \tContactSyncs in Queue %s\n", humanize.Comma((int64)(contacts)), statisticsChange("contacts", contacts))

			// HTTP Statistics
			fmt.Fprintln(out)

			httpreq, _ := redis.Int(red.Do("ZCARD", "EVEDATA_HTTPRequest"))
			red.Do("ZREMRANGEBYSCORE", "EVEDATA_HTTPRequest", 0, time.Now().UTC().Unix())
			fmt.Fprintf(out, "%s \tHTTP Requests (%d rps) %s\n", humanize.Comma((int64)(httpreq)), httpreq/5, statisticsChange("httpreq", httpreq))

			httpcount, _ := redis.Int(red.Do("ZCARD", "EVEDATA_HTTPCount"))
			red.Do("ZREMRANGEBYSCORE", "EVEDATA_HTTPCount", 0, time.Now().UTC().Unix())
			fmt.Fprintf(out, "%s \tHTTP API Calls (%d rps) %s\n", humanize.Comma((int64)(httpcount)), httpcount/5, statisticsChange("httpcount", httpcount))

			httperr, _ := redis.Int(red.Do("ZCARD", "EVEDATA_HTTPErrorCount"))
			red.Do("ZREMRANGEBYSCORE", "EVEDATA_HTTPErrorCount", 0, time.Now().UTC().Unix())
			fmt.Fprintf(out, "%s \tHTTP Errors %s\n", humanize.Comma((int64)(httperr)), statisticsChange("httperr", httperr))

			// Write out the stats
			err := out.Flush()
			statisticsTxt = w.Bytes()

			red.Do("SET", "EVEDATA_statistics", statisticsTxt)
			if err != nil {
				log.Printf("top error: %v\n", err)
			}
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
	setCache(w, 0)
	redisConn := c.Cache.Get()
	defer redisConn.Close()
	stats, _ := redis.Bytes(redisConn.Do("GET", "EVEDATA_statistics"))
	w.Write(stats)
	return http.StatusOK, nil
}
