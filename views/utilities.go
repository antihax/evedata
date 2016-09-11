package views

import (
	"net/http"
	"strconv"
	"time"
)

func setCache(w http.ResponseWriter, cacheTime int) {
	w.Header().Set("Cache-Control", "max-age:"+strconv.Itoa(cacheTime)+", public")
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Expires", time.Now().UTC().Add(time.Second*time.Duration(cacheTime)).Format(http.TimeFormat))
}
