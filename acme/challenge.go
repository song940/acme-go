package acme

import "encoding/json"

type ChallengeResponse struct {
	Type   string `json:"type"`
	Status string `json:"status"`
	URL    string `json:"url"`
	Token  string `json:"token"`
}

func (client *Client) GetChallenge(challengeUrl string) (resp *ChallengeResponse, err error) {
	_, data, err := client.get(challengeUrl)
	if err != nil {
		return
	}
	resp = &ChallengeResponse{}
	err = json.Unmarshal(data, &resp)
	return
}
