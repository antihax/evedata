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
			var left, right []string

			// Remove old entries from host table
			red.Do("ZREMRANGEBYSCORE", "EVEDATA_HOST", 0, time.Now().UTC().Unix())

			w := bytes.NewBuffer(statisticsTxt)
			w.Reset()

			out := tabwriter.NewWriter(w, 20, 1, 1, ' ', tabwriter.AlignRight)

			// this will store the keys of each iteration
			var host []string

			if arr, err := redis.MultiBulk(red.Do("ZRANGEBYSCORE", "EVEDATA_HOST", 0, "inf")); err != nil {
				log.Println(err)
			} else {
				host, _ = redis.Strings(arr, nil)
			}
			sort.Strings(host)
			for i := 0; i < len(host); i++ {
				fmt.Fprintf(out, "%s\n", host[i])
			}

			// HTTP Statistics
			fmt.Fprintln(out)

			killq, _ := redis.Int(red.Do("SCARD", "EVEDATA_killQueue"))
			left = append(left, fmt.Sprintf("%s %s\tKill    Queue", statisticsChange("killq", killq), humanize.Comma((int64)(killq))))
			entityq, _ := redis.Int(red.Do("SCARD", "EVEDATA_entityQueue"))
			left = append(left, fmt.Sprintf("%s %s\tEntity  Queue", statisticsChange("entityq", entityq), humanize.Comma((int64)(entityq))))
			history, _ := redis.Int(red.Do("SCARD", "EVEDATA_marketHistory"))
			left = append(left, fmt.Sprintf("%s %s\tHist    Queue", statisticsChange("history", history), humanize.Comma((int64)(history))))
			orders, _ := redis.Int(red.Do("SCARD", "EVEDATA_marketOrders"))
			left = append(left, fmt.Sprintf("%s %s\tMarket  Queue", statisticsChange("orders", orders), humanize.Comma((int64)(orders))))
			regions, _ := redis.Int(red.Do("ZCARD", "EVEDATA_marketRegions"))
			left = append(left, fmt.Sprintf("%s %s\tRegion  Queue", statisticsChange("regions", regions), humanize.Comma((int64)(regions))))
			contacts, _ := redis.Int(red.Do("SCARD", "EVEDATA_contactSyncQueue"))
			left = append(left, fmt.Sprintf("%s %s\tSync    Queue", statisticsChange("contacts", contacts), humanize.Comma((int64)(contacts))))
			assets, _ := redis.Int(red.Do("SCARD", "EVEDATA_assetQueue"))
			left = append(left, fmt.Sprintf("%s %s\tAsset   Queue", statisticsChange("assets", assets), humanize.Comma((int64)(assets))))
			wars, _ := redis.Int(red.Do("SCARD", "EVEDATA_warQueue"))
			left = append(left, fmt.Sprintf("%s %s\tWars    Queue", statisticsChange("wars", wars), humanize.Comma((int64)(wars))))
			journal, _ := redis.Int(red.Do("SCARD", "EVEDATA_walletQueue"))
			left = append(left, fmt.Sprintf("%s %s\tJournal  Queue", statisticsChange("journal", journal), humanize.Comma((int64)(journal))))

			kills, _ := redis.Int(red.Do("SCARD", "EVEDATA_knownKills"))
			right = append(right, fmt.Sprintf("%s \tKnown Kills %s", humanize.Comma((int64)(kills)), statisticsChange("kills", kills)))

			finWars, _ := redis.Int(red.Do("SCARD", "EVEDATA_knownFinishedWars"))
			right = append(right, fmt.Sprintf("%s \tFinished Wars %s", humanize.Comma((int64)(finWars)), statisticsChange("finWars", finWars)))

			right = append(right, "")

			red.Do("ZREMRANGEBYSCORE", "EVEDATA_HTTPRequest", 0, time.Now().UTC().Add(time.Second*-30).Unix())
			httpreq, _ := redis.Int(red.Do("ZCARD", "EVEDATA_HTTPRequest"))
			right = append(right, fmt.Sprintf("%.1f rps \tHTTP Requests", (float64)(httpreq)/30))

			red.Do("ZREMRANGEBYSCORE", "EVEDATA_HTTPCount", 0, time.Now().UTC().Add(time.Second*-30).Unix())
			httpcount, _ := redis.Int(red.Do("ZCARD", "EVEDATA_HTTPCount"))
			right = append(right, fmt.Sprintf("%.1f rps \tHTTP API Calls  ", (float64)(httpcount)/30))

			red.Do("ZREMRANGEBYSCORE", "EVEDATA_HTTPErrorCount", 0, time.Now().UTC().Add(time.Second*-30).Unix())
			httperr, _ := redis.Int(red.Do("ZCARD", "EVEDATA_HTTPErrorCount"))
			right = append(right, fmt.Sprintf("%.1f eps \tHTTP Errors  ", (float64)(httperr)/30))

			l := 0
			ll := len(left)
			lr := len(right)
			if ll > lr {
				l = ll
			} else {
				l = lr
			}

			for i := 0; i < l; i++ {
				if ll > i && lr > i {
					fmt.Fprintf(out, "%s\t%s\n", left[i], right[i])
				} else if ll > i && lr <= i {
					fmt.Fprintf(out, "%s\t\n", left[i])
				} else {
					fmt.Fprintf(out, "\t\t%s\n", right[i])
				}
			}

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
