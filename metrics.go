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

var vkMessageNew = promauto.NewCounter(prometheus.CounterOpts{
	Name: "vk_message_new",
})
var vkWallReplyNew = promauto.NewCounter(prometheus.CounterOpts{
	Name: "vk_wall_reply_new",
})
var vkGroupLeave = promauto.NewCounter(prometheus.CounterOpts{
	Name: "vk_group_leave",
})
var vkGroupJoin = promauto.NewCounter(prometheus.CounterOpts{
	Name: "vk_group_join"/**/,
})

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

	switch requestType := parsedBody["type"]; requestType {
	case "confirmation":
		_, _ = fmt.Fprint(responseWriter, confirmationString)
		return
	case "message_new":
		vkMessageNew.Add(1)
	case "wall_reply_new":
		vkWallReplyNew.Add(1)
	case "group_leave":
		vkGroupLeave.Add(1)
	case "group_join":
		vkGroupJoin.Add(1)
	}

	_, _ = fmt.Fprint(responseWriter, "ok")
}

func main() {
	confirmationString = os.Getenv("CONFIRMATION_STR")
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/callback", callbackHandler)
	log.Printf("Starting metrics server on 2112...\nConfirmation string: %s...\n", confirmationString)
	log.Fatal(http.ListenAndServe(":2112", nil))
}
