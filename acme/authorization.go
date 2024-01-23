package acme

import "encoding/json"

func (client *Client) GetAuthorization(authorizationUrl string) (resp *AuthorizationResponse, err error) {
	_, data, err := client.get(authorizationUrl)
	if err != nil {
		return
	}
	resp = &AuthorizationResponse{}
	err = json.Unmarshal(data, &resp)
	return
}

func (client *Client) GetAuthorizations(authorizationUrls []string) (resp []*AuthorizationResponse, err error) {
	for _, authorizationUrl := range authorizationUrls {
		auth, err := client.GetAuthorization(authorizationUrl)
		if err != nil {
			return nil, err
		}
		resp = append(resp, auth)
	}
	return
}

func (client *Client) DeactivateAuthorization(authorizationUrl string) (err error) {
	deactivateReq := map[string]interface{}{
		"status": "deactivated",
	}
	_, _, err = client.post(authorizationUrl, deactivateReq)
	if err != nil {
		return
	}
	return
}

// func (client *Client) GetChallenge(challengeUrl string) (resp *ChallengeResponse, err error) {
// 	_, data, err := client.get(challengeUrl)
// 	if err != nil {
// 		return
// 	}
// 	resp = &ChallengeResponse{}
// 	err = json.Unmarshal(data, &resp)
// 	return
// }

func (client *Client) CompleteChallenge(challengeUrl string) (err error) {
	_, _, err = client.post(challengeUrl, nil)
	if err != nil {
		return
	}
	return
}
