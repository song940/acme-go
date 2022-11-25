package main

import (
	"log"
	"net/http"
)

type ACMEClientConfig struct {
	DirectoryUrl string
}

type ACMEClient struct {
	config *ACMEClientConfig
}

func NewClient(config *ACMEClientConfig) *ACMEClient {
	client := &ACMEClient{config}
	return client
}

func (c *ACMEClient) Request(url string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	client.Do(req)
}

func (client *ACMEClient) GetDirectory() {
	client.Request(client.config.DirectoryUrl)
}

func (client *ACMEClient) Register(contact []string, termsOfServiceAgreed bool) (account any, err error) {
	log.Println(contact, termsOfServiceAgreed)
	return
}

func (client *ACMEClient) CreateOrder(domains []string) (order any, err error) {
	return
}

func main() {
	config := &ACMEClientConfig{
		DirectoryUrl: "https://acme-staging-v02.api.letsencrypt.org/directory",
	}
	client := NewClient(config)

	account, err := client.Register([]string{"song940@gmail.com"}, false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(account)

	order, err := client.CreateOrder([]string{
		"lsong.one",
		"lsong.org",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(order)
}
