package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	srv, err := getGmailService()
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail service: %v", err)
		return
	}

	uEnv := os.Getenv("USERS")
	users := strings.Split(uEnv, ",")
	for _, user := range users {
		r, err := srv.Users.Messages.List(user).Q("from: postbank \"Successful purchase with\"").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve messages: %v", err)
			return
		}

		if len(r.Messages) == 0 {
			fmt.Println("No messages found.")
			return
		}

		for _, m := range r.Messages {
			msg, _ := srv.Users.Messages.Get(user, m.Id).Do()
			html, _ := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)

			r, _ := regexp.Compile("amount.+\\.")
			info := r.FindString(string(html))
			tokens := strings.Fields(info[:len(info)-1])

			amount := tokens[1]

			date := tokens[len(tokens)-2]
			hour := tokens[len(tokens)-1]

			mFrom := strings.LastIndex(info, tokens[4])
			mTo := strings.Index(info, date)
			merchant := info[mFrom:mTo]

			fmt.Println(date, hour, merchant, amount)
		}
	}
}
