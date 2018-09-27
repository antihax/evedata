package marketwatch

import (
	"net/http"
	"strconv"
	"time"

	"github.com/antihax/goesi"
)

func getPages(r *http.Response) (int32, error) {
	// Decode the page into int32. Return if this fails as there were no extra pages.
	pagesInt, err := strconv.Atoi(r.Header.Get("x-pages"))
	if err != nil {
		return 0, err
	}
	pages := int32(pagesInt)
	return pages, err
}

func timeUntilCacheExpires(r *http.Response) time.Duration {
	duration := time.Until(goesi.CacheExpires(r))
	if duration < time.Second {
		duration = time.Second * 10
	} else {
		duration += time.Second * 15
	}
	return duration
}
