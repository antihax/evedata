package models

import (
	"net/http"
	"net/http/httputil"
)

func AddHTTPError(req *http.Request, res *http.Response) error {
	reqText, err := httputil.DumpRequest(req, true)
	if err != nil {
		return err
	}
	resText, err := httputil.DumpResponse(res, true)
	if err != nil {
		return err
	}
	_, err = database.Exec(`INSERT INTO evedata.httpErrors 
                    (url, status, request, response, time) 
                    VALUES(?,?,?,?,NOW());`,
		req.URL.String(), res.StatusCode, reqText, resText)
	return err
}
