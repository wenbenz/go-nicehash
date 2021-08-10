package client

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const ACCOUNT_ENDPOINT = "/main/api/v2/accounting/account2/%s"

type Account struct {
	Active         bool                  `json:"active"`
	Currency       string                `json:"currency"`
	TotalBalance   string                `json:"totalBalance"`
	Available      string                `json:"available"`
	Pending        string                `json:"pending"`
	PendingDetails AccountExtendedDetail `json:"pendingDetails"`
	BtcRate        float64               `json:"btcRate"`
}

type AccountExtendedDetail struct {
	Deposit         string `json:"deposit"`
	Withdrawal      string `json:"withdrawal"`
	Exchange        string `json:"exchange"`
	HashpowerOrders string `json:"hashpowerOrders"`
	UnpaidMining    string `json:"unpaidMining"`
}

func (c *Client) GetAccount(currency string, extended bool) (*Account, error) {
	request, err := http.NewRequest("GET", getUrl(ACCOUNT_ENDPOINT, "BTC"), nil)
	if err != nil {
		return nil, err
	}
	addQueryParams(request, map[string]string{
		"extendedResponse": strconv.FormatBool(extended),
	})
	response, err := c.Do(request)
	if err != nil {
		return nil, err
	}
	var account Account
	if err = NewUnexpectedResponseError(response.StatusCode, 200); err != nil {
		return nil, err
	}
	err = json.NewDecoder(response.Body).Decode(&account)
	return &account, err
}
