package pkg

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var searchStream *twitter.Stream
var twitterDemux twitter.SwitchDemux

func CreateTwitterClient() *twitter.Client {
	config := oauth1.NewConfig(os.Getenv("ConsumerKey"), os.Getenv("ConsumerSecret"))
	token := oauth1.NewToken(os.Getenv("Token"), os.Getenv("TokenSecret"))

	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient)
}

func StartStream() {
	client := CreateTwitterClient()
	twitterDemux = twitter.NewSwitchDemux()

	twitterDemux.Tweet = Process

	filterParams := &twitter.StreamFilterParams{
		Track:         strings.Split(strings.ReplaceAll(os.Getenv("Track"), "\"", ""), "|"),
		Language:      strings.Split(strings.ReplaceAll(os.Getenv("Lang"), "\"", ""), "|"),
		StallWarnings: twitter.Bool(true),
	}

	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	searchStream = stream

	fmt.Println("Starting Stream...")
	go twitterDemux.HandleChan(searchStream.Messages)

}

func StopStream() {
	fmt.Println("Stopping Stream...")
	searchStream.Stop()
}

func Process(tweet *twitter.Tweet) {
	fmt.Println("================================================")
	fmt.Println("[" + tweet.User.Name + "] @ " + tweet.CreatedAt)
	fmt.Println("-----------------------------")
	fmt.Println(tweet.Text)
}
