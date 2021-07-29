package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientReadFrom(t *testing.T) {
	client, err := NewClientReadFrom("./testdata/in/testCredentials.txt")
	assert.Nil(t, err)
	credentials := client.Credentials
	assert.Equal(t, "aa69aa69-0e0e-0420-aa11-420aaaa42069", credentials.Key)
	assert.Equal(t, "a123456b-7cd8-9e01-23f4-a5678bc9012def34ab56-78c9-012e-3456-78fa9b0cde12", credentials.Secret)
	assert.Equal(t, "a1a11111-1a1a-290b-8c7d-ef6f54321098", credentials.OrgId)
}

func TestPopulateAuth(t *testing.T) {
	client := NewClient(Credentials{
		Key:    "4ebd366d-76f4-4400-a3b6-e51515d054d6",
		Secret: "fd8a1652-728b-42fe-82b8-f623e56da8850750f5bf-ce66-4ca7-8b84-93651abc723b",
		OrgId:  "da41b3bc-3d0b-4226-b7ea-aee73f94a518",
	})
	request, _ := http.NewRequest(
		"GET",
		"https://api2.nicehash.com/main/api/v2/hashpower/orderBook?algorithm=X16R&page=0&size=100",
		nil,
	)
	request.Header.Add("X-Time", "1543597115712")
	request.Header.Add("X-Nonce", "9675d0f8-1325-484b-9594-c9d6d3268890")
	request.Header.Add("X-Organization-Id", client.Credentials.OrgId)
	client.populateAuth(*request)
	assert.Equal(t, "4ebd366d-76f4-4400-a3b6-e51515d054d6:21e6a16f6eb34ac476d59f969f548b47fffe3fea318d9c99e77fc710d2fed798", request.Header.Get("X-Auth"))
}

func TestPopulateHeader(t *testing.T) {
	client := NewClient(Credentials{
		Key:    "4ebd366d-76f4-4400-a3b6-e51515d054d6",
		Secret: "fd8a1652-728b-42fe-82b8-f623e56da8850750f5bf-ce66-4ca7-8b84-93651abc723b",
		OrgId:  "da41b3bc-3d0b-4226-b7ea-aee73f94a518",
	})
	request, _ := http.NewRequest(
		"GET",
		"https://api2.nicehash.com/main/api/v2/hashpower/orderBook?algorithm=X16R&page=0&size=100",
		nil,
	)
	client.populateHeaders(*request)
	assert.Len(t, request.Header.Get("X-Time"), 13)
	assert.Len(t, request.Header.Get("X-Nonce"), 36)
}
