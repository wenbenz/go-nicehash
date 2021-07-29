package client

import (
	"net/http"
	"strings"
	"time"
)

const (
	CREATE_REPORT_ENDPOINT   = "/main/api/v2/reports/add"
	REPORTS_LIST_ENDPOINT    = "/main/api/v2/reports/list"
	DOWNLOAD_REPORT_ENDPOINT = "/main/api/v2/reports/download/%s"
)

type ReportMetadata struct {
	Id      string    `json:"id"`
	Status  int       `json:"status"` // 0 = not ready; 1 = ready
	Name    string    `json:"name"`
	Created time.Time `json:"createdTs"`
	Updated time.Time `json:"updatedTs"`
}

//transactionType, crypto, fiat, aggregation string,
// start, end time.Time, timezone, timezoneOffset int
func (c *Client) CreateReport() error {
	requestBody := `{
		"transaction": "ALL",
		"currency": "BTC",
		"fiat": "CAD",
		"aggregation": "NONE",
		"dateFrom": "1609480861000",
		"dateTo":   "1627433670000",
		"timezoneOffset": "0",
		"timezoneValue": "0",
		"personal": true
	}`
	request, err := http.NewRequest("POST", getUrl(CREATE_REPORT_ENDPOINT), strings.NewReader(requestBody))
	if err != nil {
		return err
	}
	err = c.Do(request, nil)
	return err
}

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

func (c *Client) GetReport(id string) ([]byte, error) {
	request, err := http.NewRequest("GET", getUrl(DOWNLOAD_REPORT_ENDPOINT, id), nil)
	if err != nil {
		return nil, err
	}
	var b []byte
	if err = c.Do(request, &b); err == nil {
		return b, nil
	}
	return nil, err
}
