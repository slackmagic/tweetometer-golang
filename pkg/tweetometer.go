package pkg

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var searchStream *twitter.Stream
var twitterDemux twitter.SwitchDemux

func init(){
	fmt.Println("INIT TWEETOMETER")
	OpenDB()
}

func CreateTwitterClient() *twitter.Client {
	
	config := oauth1.NewConfig(os.Getenv("ConsumerKey"), os.Getenv("ConsumerSecret"))
	token := oauth1.NewToken(os.Getenv("Token"), os.Getenv("TokenSecret"))

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
	//DisplayTweet(tweet)
	encodedTweet := Compress(EncodeToBytes(tweet))
	decodedTweet := DecodeToTweet(Decompress(encodedTweet))

	DisplayTweet(&decodedTweet)
}

func DisplayTweet(tweet *twitter.Tweet){
	fmt.Println("================================================")
	fmt.Println("[" + tweet.User.Name + "] @ " + tweet.CreatedAt)
	fmt.Println("-----------------------------")
	fmt.Println(tweet.Text)
}

func EncodeToBytes(p interface{})[]byte {
	
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("uncompressed size (bytes): ", len(buf.Bytes()))
	return buf.Bytes()
}

func DecodeToTweet(s []byte) twitter.Tweet {

	p := twitter.Tweet{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)
	if err != nil {
		log.Fatal(err)
	}
	return p
}
