package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	NICEHASH_API = "https://api2.nicehash.com"
)

//Client sends signed requests to the NiceHash API using the provided credentials
type Client struct {
	Credentials Credentials
	httpClient  http.Client
}

//Credentials stores the user's API key, secret, and org id used to sign requests.
type Credentials struct {
	Key    string
	Secret string
	OrgId  string
}

//NewClient creates a new client from the credentials struct.
func NewClient(cred Credentials) *Client {
	return &Client{
		Credentials: cred,
		httpClient:  http.Client{},
	}
}

//NewClientReadFrom returns a client created from reading credentials from a file.
//path is the system path to a file with 3 lines:
// 1. api key
// 2. api secret
// 3. org id
func NewClientReadFrom(path string) (*Client, error) {
	credentials, err := getCredentials(path)
	if err != nil {
		return nil, err
	}
	return &Client{
		Credentials: credentials,
		httpClient:  http.Client{},
	}, nil
}

//Do populates required request headers and signs the request using the client secret key.
//Decodes response into destination if destination is not nil.
//JSON responses get unmarshalled into destination;
//CSV responses returns the reader.
func (c *Client) Do(request *http.Request) (*http.Response, error) {
	c.populateHeaders(*request)
	if err := c.populateAuth(*request); err != nil {
		return nil, err
	}
	return c.httpClient.Do(request)
}

func getUrl(endpoint string, args ...interface{}) string {
	return NICEHASH_API + fmt.Sprintf(endpoint, args...)
}

func addQueryParams(request *http.Request, params map[string]string) {
	query := request.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	request.URL.RawQuery = query.Encode()
}

func getCredentials(path string) (Credentials, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return Credentials{}, err
	}
	fields := strings.Split(string(fileBytes), "\n")
	return Credentials{
		Key:    fields[0],
		Secret: fields[1],
		OrgId:  fields[2],
	}, err
}

//timeToLongString converts a timestamp to a string representation of the milliseconds from epoch
func timeToLongString(t time.Time) string {
	return strconv.FormatInt(t.Unix()*1000, 10)
}

func (c *Client) populateHeaders(request http.Request) {
	// make random nonce
	nonceBytes := make([]byte, 18) // becomes 36 byte hex
	rand.Read(nonceBytes)
	nonce := hex.EncodeToString(nonceBytes)

	request.Header.Add("X-Time", timeToLongString(time.Now()))
	request.Header.Add("X-Nonce", nonce)
	request.Header.Add("X-Organization-Id", c.Credentials.OrgId)
	request.Header.Add("X-Request-Id", nonce)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Accept", "text/csv")
}

func (c *Client) populateAuth(request http.Request) error {
	input := bytes.NewBuffer([]byte{})
	input.WriteString(c.Credentials.Key)
	input.WriteByte(0x00)
	input.WriteString(request.Header.Get("X-Time"))
	input.WriteByte(0x00)
	input.WriteString(request.Header.Get("X-Nonce"))
	input.WriteByte(0x00)
	input.WriteByte(0x00)
	input.WriteString(request.Header.Get("X-Organization-Id"))
	input.WriteByte(0x00)
	input.WriteByte(0x00)
	input.WriteString(request.Method)
	input.WriteByte(0x00)
	input.WriteString(request.URL.Path)
	input.WriteByte(0x00)
	input.WriteString(request.URL.RawQuery)
	if request.Body != nil {
		bodyReader, err := request.GetBody()
		if err != nil {
			return err
		}
		body, err := ioutil.ReadAll(bodyReader)
		if err != nil {
			return err
		}
		bodyStr := string(body)
		if bodyStr != "" {
			input.WriteByte(0x00)
			input.WriteString(bodyStr)
		}
	}

	inputString := input.String()

	// hash
	mac := hmac.New(sha256.New, []byte(c.Credentials.Secret))
	mac.Write([]byte(inputString))
	request.Header.Add("X-Auth", c.Credentials.Key+":"+hex.EncodeToString(mac.Sum(nil)))
	return nil
}
