package main

import (
	"log"

	"github.com/song940/acme/acme"
)

const accountKey = `
-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgA8aaCjB1AlV2ndWt
y/1mBDxzEZdvXymT/aBCXw1E26KhRANCAASv8qR9xkSTsOHGBB8F1OEPYQ4gmst1
k3JMM1Bg/XKlyNfynRX+WfB6VtQaiPllh5qazOgX3xfOeNcQQIqzQeVU
-----END PRIVATE KEY-----
`

func main() {
	client, err := acme.NewClient(&acme.Config{
		DirectoryURL: "https://acme-staging-v02.api.letsencrypt.org/directory",
	})
	client.Account.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	accountUrl, accResp, err := client.Register(&acme.AccountRequest{
		TermsOfServiceAgreed: true,
		Contact: []string{
			"mailto: song940@gmail.com",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	client.Account.URL = accountUrl
	log.Println(accountUrl, accResp)
	orderUrl, orderResp, err := client.CreateOrder(&acme.OrderRequest{
		Identifiers: []acme.Identifier{
			{
				Type:  "dns",
				Value: "lsong.org",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(orderUrl, orderResp)
	order, err := client.GetOrder(orderUrl)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(order)
	for _, authUrl := range order.Authorizations {
		auth, err := client.GetAuthorization(authUrl)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(auth.Status)
		for _, ch := range auth.Challenges {
			log.Println(ch.Type, ch.Status, ch.URL, ch.Token)
		}
	}
}
