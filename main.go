package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/segmentio/go-env"
)

// additionalParameters for the impersonation request.
type additionalParameters struct {
	ResponseType string `json:"response_type"`
	Scope        string `json:"scope"`
}

// impersonationInput request input.
type impersonationInput struct {
	Protocol             string               `json:"protocol"`
	ImpersonatorID       string               `json:"impersonator_id"`
	ClientID             string               `json:"client_id"`
	AdditionalParameters additionalParameters `json:"additionalParameters"`
}

// tokenInput request input.
type tokenInput struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
}

// tokenOutput request output.
type tokenOutput struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func main() {
	clientID := env.MustGet("AUTH0_CLIENT_ID")
	clientSecret := env.MustGet("AUTH0_CLIENT_SECRET")

	impersonatorID := flag.String("impersonator-id", "", "User ID of impersonator")
	appClientID := flag.String("client-id", "", "Client ID of the application")
	scope := flag.String("scope", "openid name user_id nickname email picture", "OAuth scope")

	flag.Parse()

	userID := flag.Arg(0)
	if userID == "" {
		log.Fatalf("<user-id> argument required")
	}

	token, err := getToken("apex-inc", clientID, clientSecret)
	if err != nil {
		log.Fatalf("error fetching token: %s", err)
	}

	link, err := getImpersionationLink("apex-inc", userID, *impersonatorID, *appClientID, token, *scope)
	if err != nil {
		log.Fatalf("error fetching link: %s", err)
	}

	fmt.Println(link)
}

// getToken returns a token which expires within 24 hours.
func getToken(account, clientID, clientSecret string) (string, error) {
	url := fmt.Sprintf("https://%s.auth0.com/oauth/token", account)

	body := &tokenInput{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
	}

	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	res, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	var v tokenOutput
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return "", err
	}

	return v.AccessToken, nil
}

// getImpersionationLink returns a link which can be used to authenticate as the user.
func getImpersionationLink(account, userID, impersonatorID, clientID, token, scope string) (string, error) {
	url := fmt.Sprintf("https://%s.auth0.com/users/%s/impersonate", account, userID)

	body := &impersonationInput{
		Protocol:       "oauth2",
		ImpersonatorID: impersonatorID,
		ClientID:       clientID,
		AdditionalParameters: additionalParameters{
			ResponseType: "token",
			Scope:        scope,
		},
	}

	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
