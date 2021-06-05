package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	"github.com/joho/godotenv"
)

func init() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println(os.Getenv("NAME"))
}

func main() {

	config := oauth1.NewConfig(os.Getenv("ConsumerKey"), os.Getenv("ConsumerSecret"))
	token := oauth1.NewToken(os.Getenv("Token"), os.Getenv("TokenSecret"))

	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()

	demux.Tweet = func(tweet *twitter.Tweet) {
		fmt.Println("================================================")
		fmt.Println("[" + tweet.User.Name + "] @ " + tweet.CreatedAt)
		fmt.Println("-----------------------------")
		fmt.Println(tweet.Text)

	}

	//demux.Tweet = pkg.process(*twitter.tweet)

	fmt.Println("Starting Stream...")

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{"#zlan2021"},
		Language:      []string{"fr"},
		StallWarnings: twitter.Bool(true),
	}

	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits
	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}
