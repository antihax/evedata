package eveapi

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type crestSimpleFrame struct {
	CacheUntil time.Time
}

type crestPagedFrame struct {
	crestSimpleFrame

	Next struct {
		HRef string
	}
	TotalCount int
	PageCount  int
	Page       int
}

func (c *crestSimpleFrame) getFrameInfo(r *http.Response) error {

	maxAge := strings.Split(r.Header.Get("Cache-Control"), "=")[1]
	iMaxAge, err := strconv.Atoi(maxAge)
	if err != nil {
		return err
	}

	date, err := time.Parse(time.RFC1123, r.Header.Get("Date"))
	if err != nil {
		return err
	}

	c.CacheUntil = date.Add(time.Duration(iMaxAge) * time.Second)

	return nil
}

func (c *crestPagedFrame) getFrameInfo(url string, r *http.Response) error {
	if err := c.crestSimpleFrame.getFrameInfo(r); err != nil {
		return err
	}
	page, err := getPageNumberFromURL(url)
	if err != nil {
		return err
	}
	c.Page = page
	return nil
}

func getPageNumberFromURL(s string) (int, error) {
	u, err := url.Parse(s)
	if err != nil {
		return 0, err
	}

	m, err := url.ParseQuery(u.RawQuery)
	if m != nil {
		if m["page"] != nil {
			i, err := strconv.Atoi(m["page"][0])
			if err != nil {
				return 0, err
			}
			return i, nil
		}
	}
	return 0, err
}
