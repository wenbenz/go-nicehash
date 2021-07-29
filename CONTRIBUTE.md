# Welcome!
Thanks for contributing to the go-nicehash project! This document will describe how the project is structured and how to quickly add endpoints.

# Project Structure
`client.go` is the workhorse of this project. Its `Do` method accepts a request does all the magic for signing the request as described [here](https://www.nicehash.com/docs/). So all you have to do, dear contributor, is define the struct and populate the request with the URL, parameters, and body.

# How to add a new endpoint
1. Create a new branch to work under
2. Find/Create a file under the `client` directory, named after the operation(s).
3. If one doesn't exist in the file, add a constant for the URL format string. E.g.
```
REPORTS_LIST_ENDPOINT    = "/main/api/v2/reports/list"
```
4. If the target endpoint returns data, define a data type. E.g.
```
//ReportMetadata describes a NiceHash report
type ReportMetadata struct {
	Id      string    `json:"id"`
	Status  int       `json:"status"` // 0 = not ready; 1 = ready
	Name    string    `json:"name"`
	Created time.Time `json:"createdTs"`
	Updated time.Time `json:"updatedTs"`
}
```
5. Define a public `Client` method which populates a `http.Request` with the following information, then pass the request into the `Client.Do` method:
    - the method (e.g. "GET")
    - the URL defined in step 3 (and args if any); use `Client.getUrl` to format the url string
    - (optional) the body
    - (optional) query parameters using `Client.addQueryParams`
```
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

```

...and that's it. `Client.Do` will sign the request. Manually verify that the new method works as expected.

**Note:** if the response is type `text/csv`, `Client.Do` puts an `io.ReadCloser` in the destination.