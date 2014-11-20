package main

import (
	"./secret"
	"encoding/json"
	"errors"
	mailGun "github.com/mailgun/mailgun-go"
	"net/http"
	"strconv"
)

type GithubPullRequest struct {
	Action      string `json:"action"`
	PullRequest struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		Url    string `json:"url"`
		User   struct {
			Login string `json:"login"`
		} `json:"user"`
	} `json:"pull_request"`
}

func main() {
	http.HandleFunc("/app/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		header := r.Header.Get("X-Github-Event")

		if header != "pull_request" {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(header)
			return
		}

		var pullRequest GithubPullRequest
		if err := json.NewDecoder(r.Body).Decode(&pullRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := handlePullRequest(pullRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(&pullRequest)
	})

	http.ListenAndServe(":8080", nil)
}

func handlePullRequest(pullRequest GithubPullRequest) error {
	switch pullRequest.Action {
	case "opened", "reopened":
		alertNewPullRequest(pullRequest)
		return nil
	}
	errorMsg := "Not yet handling Action: " + pullRequest.Action
	return errors.New(errorMsg)
}

func alertNewPullRequest(pullRequest GithubPullRequest) {
	gun := mailGun.NewMailgun(secret.Domain,
		secret.ApiKey,
		secret.PublicApiKey)

	message := gun.NewMessage(
		"Music Collaboratory <postmaster@sandbox81f2f8104dc74714a4a50f801c022e9c.mailgun.org>",
		"New Pull Request #"+strconv.Itoa(pullRequest.PullRequest.Number),
		"Pull Request opened by: "+pullRequest.PullRequest.User.Login+"\n\n"+
			"Title: "+pullRequest.PullRequest.Title+"\n"+
			"Please review at "+pullRequest.PullRequest.Url)

	for _, user := range secret.MailingList {
		message.AddRecipient(user)
	}

	gun.Send(message)
}
