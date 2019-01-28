package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var confirmationString string

var counterMap = make(map[string]prometheus.Counter)
var eventTypes = []string{
	"message_new",
	"message_reply",
	"message_edit",
	"message_allow",
	"message_deny",
	"photo_new",
	"photo_comment_new",
	"photo_comment_restore",
	"photo_comment_delete",
	"audio_new",
	"view_new",
	"video_comment_new",
	"video_comment_edit",
	"video_comment_restore",
	"video_comment_delete",
	"wall_post_new",
	"wall_repost",
	"wall_reply_new",
	"wall_reply_edit",
	"wall_reply_restore",
	"wall_reply_delete",
	"board_post_new",
	"board_post_edit",
	"board_post_restore",
	"board_post_delete",
	"market_comment_new",
	"market_comment_edit",
	"market_comment_restore",
	"market_comment_delete",
	"group_leave",
	"group_join",
	"user_block",
	"user_unblock",
	"poll_vote_new",
	"group_officers_edit",
	"group_change_settings",
	"group_change_photo",
}

func initializeHandlers() {
	for _, eventType := range eventTypes {
		counterMap[eventType] = promauto.NewCounter(prometheus.CounterOpts{
			Name: eventType,
		})
	}
}

func callbackHandler(responseWriter http.ResponseWriter, request *http.Request) {
	var body []byte
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Print(err)
		return
	}

	var parsedBody map[string]interface{}
	err = json.Unmarshal(body, &parsedBody)
	if err != nil {
		fmt.Print(err)
		return
	}

	requestType := parsedBody["type"].(string)
	if handler, ok := counterMap[requestType]; ok {
		handler.Add(1)
	}
	_, _ = fmt.Fprint(responseWriter, "ok")
}

func main() {
	initializeHandlers()
	confirmationString = os.Getenv("CONFIRMATION_STR")
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/callback", callbackHandler)
	log.Printf("Starting metrics server on 2112...\nConfirmation string: %s...\n", confirmationString)
	log.Fatal(http.ListenAndServe(":2112", nil))
}
