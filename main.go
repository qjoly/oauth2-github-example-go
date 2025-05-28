package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var (
	clientID     = os.Getenv("GITHUB_CLIENT_ID")
	clientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
)

func main() {
	if clientID == "" || clientSecret == "" {
		log.Fatal("GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET must be set")
	}

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:8080/auth/github/callback",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<html><body>
			<h1>GitHub OAuth2 Example</h1>
			<p><a href="/login">Login with GitHub</a></p>
			</body></html>`)
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOnline)
		http.Redirect(w, r, url, http.StatusFound)
	})

	http.HandleFunc("/auth/github/callback", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != "state" {
			http.Error(w, "State invalid", http.StatusBadRequest)
			return
		}

		token, err := conf.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		client := conf.Client(r.Context(), token)
		resp, err := client.Get("https://api.github.com/user")
		if err != nil {
			http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var userInfo map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
			http.Error(w, "Failed to parse user info: "+err.Error(), http.StatusInternalServerError)
			return
		}

		emailsResp, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			http.Error(w, "Failed to get user emails: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer emailsResp.Body.Close()

		var emails []map[string]interface{}
		if err := json.NewDecoder(emailsResp.Body).Decode(&emails); err != nil {
			http.Error(w, "Failed to parse user emails: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, `<html><body>
			<h1>GitHub Authentication Successful</h1>
			<h2>User Info</h2>
			<pre>%s</pre>
			<h2>User Emails</h2>
			<pre>%s</pre>
			<h2>Token Info</h2>
			<pre>%s</pre>
			</body></html>`,
			prettyPrint(userInfo),
			prettyPrint(emails),
			prettyPrint(token))
	})

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func prettyPrint(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}
