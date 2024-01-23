package acme

import (
	"encoding/json"
)

// Identifier object used in order and authorization objects
// See https://tools.ietf.org/html/rfc8555#section-7.1.4
type Identifier struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type OrderRequest struct {
	Identifiers []Identifier `json:"identifiers"`
}

type OrderResponse struct {
	Status         string       `json:"status"`
	Expires        string       `json:"expires"`
	Identifiers    []Identifier `json:"identifiers"`
	Authorizations []string     `json:"authorizations"`
	Finalize       string       `json:"finalize"`
}

type AuthorizationResponse struct {
	Identifier Identifier `json:"identifier"`
	Status     string     `json:"status"`
	Expires    string     `json:"expires"`
	Challenges []struct {
		Type   string `json:"type"`
		Status string `json:"status"`
		URL    string `json:"url"`
		Token  string `json:"token"`
	} `json:"challenges"`
}

func (client *Client) CreateOrder(request *OrderRequest) (url string, resp *OrderResponse, err error) {
	headers, data, err := client.post(client.Directory.NewOrder, request)
	if err != nil {
		return
	}
	url = headers.Get("Location")
	resp = &OrderResponse{}
	err = json.Unmarshal(data, &resp)
	return
}

func (client *Client) GetOrder(orderUrl string) (resp *OrderResponse, err error) {
	_, data, err := client.get(orderUrl)
	if err != nil {
		return
	}
	resp = &OrderResponse{}
	err = json.Unmarshal(data, &resp)
	return
}

func (client *Client) FinalizeOrder(orderUrl string, csr string) (resp *OrderResponse, err error) {
	finalizeReq := map[string]interface{}{
		"csr": csr,
	}
	_, data, err := client.post(orderUrl, finalizeReq)
	if err != nil {
		return
	}
	resp = &OrderResponse{}
	err = json.Unmarshal(data, &resp)
	return
}
