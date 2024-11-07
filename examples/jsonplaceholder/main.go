package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/zrbecker/fetch"
	"github.com/zrbecker/fetch/examples/jsonplaceholder/placeholder"
	"github.com/zrbecker/fetch/transports"
)

var baseURL = flag.String("base-url", "https://jsonplaceholder.typicode.com", "base url for jsonplaceholder api")
var debug = flag.Bool("debug", false, "should print debug request data")

func main() {
	flag.Parse()

	ctx := context.Background()

	opts := make([]fetch.FetcherOption, 0)
	if *debug {
		httpClient := &http.Client{
			Transport: transports.NewDebugTransport(http.DefaultTransport),
		}
		opts = append(opts, fetch.Client(httpClient))
	}

	placeholderClient := placeholder.NewClient(*baseURL, opts...)
	response, err := placeholderClient.CreatePost(ctx, placeholder.CreatePostParams{})
	if err != nil {
		log.Panic(err)
	}

	jsonBz, err := json.Marshal(response)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(jsonBz))
}
