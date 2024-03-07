package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

const mondayURL = "https://api.monday.com/v2"

var mondayToken = ""

func PostToMonday(data PISignUP) {
	req, err := http.NewRequest("POST", mondayURL,
		strings.NewReader(`{"query": "mutation {
				create_item (
					board_id: 1373816796,
					group_id: topics,
					column_values:
				)
			}"
		}`), // FIXME: prod values, don't hardcode this.
	)
	if err != nil {
		log.Println("Could not add new member" +
			data.Member.Infix + " " + data.Member.LastName + ", " + data.Member.FirstName +
			" to the members file.",
		)
	}

	req.Header.Add("Authorization", mondayToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Could not add new member" +
			data.Member.Infix + " " + data.Member.LastName + ", " + data.Member.FirstName +
			" to the members file.",
		)
		buf, _ := io.ReadAll(resp.Body)
		log.Println(string(buf))
	}
}
