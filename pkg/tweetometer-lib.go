package pkg

import (
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
)

func process(tweet *twitter.Tweet) {
	fmt.Println("================================================")
	fmt.Println("[" + tweet.User.Name + "] @ " + tweet.CreatedAt)
	fmt.Println("-----------------------------")
	fmt.Println(tweet.Text)

}
