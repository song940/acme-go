package acme

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"io"
	"net/http"
)

type H map[string]interface{}

type Directory struct {
	NewNonce    string `json:"newNonce"`    // url to new nonce endpoint
	NewAccount  string `json:"newAccount"`  // url to new account endpoint
	NewOrder    string `json:"newOrder"`    // url to new order endpoint
	NewAuthz    string `json:"newAuthz"`    // url to new authz endpoint
	RevokeCert  string `json:"revokeCert"`  // url to revoke cert endpoint
	KeyChange   string `json:"keyChange"`   // url to key change endpoint
	RenewalInfo string `json:"renewalInfo"` // url to renewal info endpoint

	// meta object containing directory metadata
	Meta struct {
		TermsOfService          string   `json:"termsOfService"`
		Website                 string   `json:"website"`
		CaaIdentities           []string `json:"caaIdentities"`
		ExternalAccountRequired bool     `json:"externalAccountRequired"`
	} `json:"meta"`
}

type Response struct {
	Status string `json:"status"`
	Type   string `json:"type"`
	Detail string `json:"detail"`
}

type RevokeCertRequest struct {
	Certificate string `json:"certificate"`
	Reason      int    `json:"reason"`
}

type Account struct {
	PrivateKey crypto.Signer
	URL        string
}

func (acc *Account) GenerateKey() (err error) {
	acc.PrivateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return
}

func (acc *Account) ImportKey(key string) (err error) {
	priv, err := x509.ParsePKCS8PrivateKey([]byte(key))
	if err != nil {
		return
	}
	acc.PrivateKey = priv.(crypto.Signer)
	return
}

func (acc *Account) ExportKey() (key string, err error) {
	data, err := x509.MarshalPKCS8PrivateKey(acc.PrivateKey)
	if err != nil {
		return
	}
	key = string(data)
	return
}

type Config struct {
	AccountKey   string `json:"accountKey"`
	AccountURL   string `json:"accountURL"`
	DirectoryURL string `json:"directoryURL"`
}

type Client struct {
	Config    *Config
	Account   *Account
	Directory *Directory
	nonce     string
}

func NewClient(config *Config) (client *Client, err error) {
	account := &Account{
		URL: config.AccountURL,
	}
	if config.AccountKey != "" {
		err = account.ImportKey(config.AccountKey)
	} else {
		err = account.GenerateKey()
	}
	if err != nil {
		return
	}
	client = &Client{Config: config, Account: account}
	client.Directory, err = client.GetDirectory()
	return
}

func (c *Client) request(method, url string, payload []byte) (res *http.Response, err error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return
	}
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/jose+json")
	}
	res, err = client.Do(req)
	if err != nil {
		return
	}
	c.nonce = res.Header.Get("Replay-Nonce")
	return
}

func (client *Client) get(url string) (headers http.Header, body []byte, err error) {
	res, err := client.request(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	headers = res.Header
	defer res.Body.Close()
	body, err = io.ReadAll(res.Body)
	return
}

func (client *Client) post(url string, payload interface{}) (headers http.Header, body []byte, err error) {
	data, err := client.buildSignedRequestData(url, payload)
	res, err := client.request(http.MethodPost, url, data)
	if err != nil {
		return
	}
	defer res.Body.Close()
	headers = res.Header
	body, err = io.ReadAll(res.Body)
	return
}

func (client *Client) GetDirectory() (directory *Directory, err error) {
	_, data, err := client.get(client.Config.DirectoryURL)
	if err != nil {
		return
	}
	directory = &Directory{}
	err = json.Unmarshal(data, directory)
	return
}

func (client *Client) getNonce() (string, error) {
	if client.nonce != "" {
		return client.nonce, nil
	}
	_, err := client.request(http.MethodHead, client.Directory.NewNonce, nil)
	return client.nonce, err
}
