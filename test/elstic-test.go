package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"

	"strings"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func main() {
	var (
		r map[string]interface{}
		// wg sync.WaitGroup
	)

	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:32774",
		},
		Username: "foo",
		Password: "bar",
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   150,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
			},
		},
	}

	es, _ := elasticsearch.NewClient(cfg)
	//

	// 1. Get cluster info
	//
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	// log.Println(res)

	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))

	// 2. Index documents concurrently
	//
	var b strings.Builder
	b.WriteString(`{"titles" : "`)
	b.WriteString(`"Test One"`)
	b.WriteString(`"}`)

	i := 0
	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      "test",
		DocumentID: strconv.Itoa(i + 1),
		Body:       strings.NewReader(b.String()),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err = req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

	log.Println(strings.Repeat("-", 37))

	// // 2. Index documents concurrently
	// //
	// for i, title := range []string{"Test One", "Test Two"} {
	// 	wg.Add(1)

	// 	go func(i int, title string) {
	// 		defer wg.Done()

	// 		// Build the request body.
	// 		var b strings.Builder
	// 		b.WriteString(`{"title" : "`)
	// 		b.WriteString(title)
	// 		b.WriteString(`"}`)

	// 		// Set up the request object.
	// 		req := esapi.IndexRequest{
	// 			Index:      "test",
	// 			DocumentID: strconv.Itoa(i + 1),
	// 			Body:       strings.NewReader(b.String()),
	// 			Refresh:    "true",
	// 		}

	// 		// Perform the request with the client.
	// 		res, err := req.Do(context.Background(), es)
	// 		if err != nil {
	// 			log.Fatalf("Error getting response: %s", err)
	// 		}
	// 		defer res.Body.Close()

	// 		if res.IsError() {
	// 			log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
	// 		} else {
	// 			// Deserialize the response into a map.
	// 			var r map[string]interface{}
	// 			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
	// 				log.Printf("Error parsing the response body: %s", err)
	// 			} else {
	// 				// Print the response status and indexed document version.
	// 				log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
	// 			}
	// 		}
	// 	}(i, title)
	// }
	// wg.Wait()

}
