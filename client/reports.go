package client

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	CREATE_REPORT_ENDPOINT   = "/main/api/v2/reports/add"
	REPORTS_LIST_ENDPOINT    = "/main/api/v2/reports/list"
	DOWNLOAD_REPORT_ENDPOINT = "/main/api/v2/reports/download/%s"
	DELETE_REPORT_ENDPOINT   = "/main/api/v2/reports/delete/%s"
)

//ReportMetadata describes a NiceHash report
type ReportMetadata struct {
	Id      string    `json:"id"`
	Status  int       `json:"status"` // 0 = not ready; 1 = ready
	Name    string    `json:"name"`
	Created time.Time `json:"createdTs"`
	Updated time.Time `json:"updatedTs"`
}

//CreateReport creates a report with the provided parameters
//	transactionType: type of transactions (e.g. "ALL", "DEPOSIT", "WITHDRAWAL", "EXCHANGE", "HASHPOWER", "MINING", "OTHER")
//	crypto: cryptocurrency symbol (e.g. "BTC", "ETH")
//	fiat: fiat currency symbol (e.g. "USD", "CAD")
// 	aggregation: time aggregate ("NONE", "DAY", "MONTH", "QUARTER", "YEAR")
//	start: timestamp for earliest record
// 	end: timestamp for latest record
// 	timezone: integer code for timezone (e.g. 0 for GMT)
//  timezoneOffset: constant added to time. (e.g. setting timezone to "0" and timezone offset to "0" will produce a UTC timestamp)
func (c *Client) CreateReport(transactionType, currency, fiat, aggregation string,
	start, end time.Time, timezone, timezoneOffset string) error {
	requestBody := fmt.Sprintf(`{
		"transaction": "%s",
		"currency": "%s",
		"fiat": "%s",
		"aggregation": "%s",
		"dateFrom": "%s",
		"dateTo": "%s",
		"timezoneOffset": "%s",
		"timezoneValue": "%s",
		"personal": true
	}`, transactionType, currency, fiat, aggregation, timeToLongString(start), timeToLongString(end), timezoneOffset, timezone)
	request, err := http.NewRequest("POST", getUrl(CREATE_REPORT_ENDPOINT), strings.NewReader(requestBody))
	if err != nil {
		return err
	}
	err = c.Do(request, nil)
	return err
}

//GetReportsList returns a list of report metadata.
//Note that this is separate from the reports seen under "settings" in the NiceHash UI.
func (c *Client) GetReportsList() ([]ReportMetadata, error) {
	var reports []ReportMetadata
	var err error
	if request, err := http.NewRequest("GET", getUrl(REPORTS_LIST_ENDPOINT), nil); err == nil {
		if err = c.Do(request, &reports); err == nil {
			return reports, nil
		}
	}
	return nil, err
}

//GetReport returns the CSV bytes of a report
func (c *Client) GetReport(id string) ([]byte, error) {
	request, err := http.NewRequest("GET", getUrl(DOWNLOAD_REPORT_ENDPOINT, id), nil)
	if err != nil {
		return nil, err
	}
	var b io.ReadCloser
	if err = c.Do(request, &b); err == nil {
		return ioutil.ReadAll(b)
	}
	return nil, err
}

//DeleteReport deletes the report with the specified ID.
func (c *Client) DeleteReport(id string) error {
	request, err := http.NewRequest("DELETE", getUrl(DELETE_REPORT_ENDPOINT, id), nil)
	if err != nil {
		return err
	}
	if err = c.Do(request, nil); err == nil {
		return nil
	}
	return err
}
