package pkg

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var searchStream *twitter.Stream
var twitterDemux twitter.SwitchDemux
var charSeparator string = "|"

func init() {
	fmt.Println("INIT TWEETOMETER")
	OpenDB()
}

func createTwitterClient() *twitter.Client {

	config := oauth1.NewConfig(os.Getenv("ConsumerKey"), os.Getenv("ConsumerSecret"))
	token := oauth1.NewToken(os.Getenv("Token"), os.Getenv("TokenSecret"))

	httpClient := config.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient)
}

func StartExtractionProcess() {
	client := createTwitterClient()
	twitterDemux = twitter.NewSwitchDemux()

	twitterDemux.Tweet = process

	filterParams := &twitter.StreamFilterParams{
		Track:         strings.Split(strings.ReplaceAll(os.Getenv("Track"), "\"", ""), charSeparator),
		Language:      strings.Split(strings.ReplaceAll(os.Getenv("Lang"), "\"", ""), charSeparator),
		StallWarnings: twitter.Bool(true),
	}

	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	searchStream = stream

	fmt.Println("Starting Stream...")
	twitterDemux.HandleChan(searchStream.Messages)

}

func StopExtractionProcess() {
	fmt.Println("Stopping Stream...")
	searchStream.Stop()
	CloseDB()
}

func process(tweet *twitter.Tweet) {
	displayTweet(tweet)
	encodedTweet := Compress(encodeToBytes(tweet))
	InsertData(createKeyFromTweet(tweet), encodedTweet)

}

func createKeyFromTweet(tweet *twitter.Tweet) []byte {
	createdAt, err := tweet.CreatedAtTime()
	if err != nil {
		log.Fatal(err)
	}

	key := createdAt.UTC().Format(time.RFC3339Nano)
	key += "@" + strconv.FormatInt(tweet.User.ID, 10)
	return []byte(key)
}

func displayTweet(tweet *twitter.Tweet) {
	fmt.Println("Insert [" + tweet.User.Name + "] @ " + tweet.CreatedAt)
	fmt.Println("-----------------------------")
}

func encodeToBytes(p interface{}) []byte {

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}

func decodeToTweet(s []byte) twitter.Tweet {

	p := twitter.Tweet{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal(err)
	}
	return p
}
