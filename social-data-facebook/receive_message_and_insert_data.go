package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"strings"
	"time"

	"github.com/streadway/amqp"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// struct insert data form Social_Data_Facebook
type Social_Data_Facebook struct {
	Post_id      string    `bson:"post_id" json:"post_id"`
	Page_id      string    `bson:"page_id" json:"page_id"`
	Page_name    string    `bson:"page_name" json:"page_name"`
	Message      string    `bson:"message" json:"message"`
	Full_picture string    `bson:"full_picture" json:"full_picture"`
	Shares       int       `bson:"shares" json:"shares"`
	Story        string    `bson:"story" json:"story"`
	Status_type  string    `bson:"status_type" json:"status_type"`
	Created_time time.Time `bson:"created_time" json:"created_time"`
	Updated_time time.Time `bson:"updated_time" json:"updated_time"`
	Created_at   time.Time
	Updated_at   time.Time `bson:"updated_at" json:"updated_at"`
}

// struct post form stream_facebook
type Stream_facebook struct {
	Page_id   string `bson:"page_id" json:"page_id"`
	Page_name string `bson:"page_name" json:"page_name"`
	Data      []Data_facebook
}

// struct data form postpage
type Data_facebook struct {
	Page_id      string    `bson:"page_id" json:"page_id"`
	Page_name    string    `bson:"page_name" json:"page_name"`
	Post_id      string    `bson:"post_id" json:"post_id"`
	Message      string    `bson:"message" json:"message"`
	Full_picture string    `bson:"full_picture" json:"full_picture"`
	Shares       int       `bson:"shares" json:"shares"`
	Story        string    `bson:"story" json:"story"`
	Status_type  string    `bson:"status_type" json:"status_type"`
	Created_time time.Time `bson:"created_time" json:"created_time"`
	Updated_time time.Time `bson:"updated_time" json:"updated_time"`
	Created_at   time.Time
	Updated_at   time.Time `bson:"updated_at" json:"updated_at"`
}

type Data_facebook_elastic struct {
	Post_id      string `bson:"post_id" json:"post_id"`
	Created_time string `bson:"created_time" json:"created_time"`
}

type Rabbit struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	QueueName  string `json:"queue_name"`
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func main() {
	rabbit := Rabbit{
		Host:      "localhost",
		Port:      "32771",
		Username:  "guest",
		Password:  "guest",
		QueueName: "fbfeed.post.add",
	}
	rabbit.Connection = rabbit.Connect()
	rabbit.Channel, _ = rabbit.Connection.Channel() // _ = ไม่ใช้ตัวแปรนั้น

	msgs, err := rabbit.Channel.Consume(
		rabbit.QueueName, // queue
		"",               // consumer
		true,             // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	failOnError(err, "Failed to register a consumer")
	Loop_Message(msgs)

	// facebook_post_Map := make(map[string]interface{})
	// facebook_post_Map := Stream_facebook{}
	// Connect mongodb
	// session, err := mgo.Dial("mongodb://127.0.0.1:32776/social-data-facebook")
	// if err != nil {
	// 	panic(err)
	// }
	// defer session.Close()
	// session.SetMode(mgo.Monotonic, true)
	// db := session.DB("social-data-facebook").C("steram-page")

	// Connect elasticsearch
	// var (
	// 	// r  map[string]interface{}
	// 	wg sync.WaitGroup
	// )

	// cfg := elasticsearch.Config{
	// 	Addresses: []string{
	// 		"http://localhost:32774",
	// 	},
	// 	Username: "foo",
	// 	Password: "bar",
	// 	Transport: &http.Transport{
	// 		MaxIdleConnsPerHost:   100,
	// 		ResponseHeaderTimeout: time.Second,
	// 		DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
	// 		TLSClientConfig: &tls.Config{
	// 			MinVersion: tls.VersionTLS11,
	// 		},
	// 	},
	// }

	// es, _ := elasticsearch.NewClient(cfg)

	// for data_stream_facebook := range msgs {
	// 	json.Unmarshal([]byte(data_stream_facebook.Body), &facebook_post_Map)
	// 	for index_message := range facebook_post_Map.Data {
	// 		data_facebook := Data_facebook{}
	// 		// fmt.Println(index_message)
	// 		// fmt.Println(facebook_post_Map.Data[index_message])
	// 		data_facebook.Page_id = facebook_post_Map.Page_id
	// 		data_facebook.Page_name = facebook_post_Map.Page_name
	// 		data_facebook.Post_id = facebook_post_Map.Data[index_message].Post_id
	// 		data_facebook.Message = facebook_post_Map.Data[index_message].Message
	// 		data_facebook.Full_picture = facebook_post_Map.Data[index_message].Full_picture
	// 		data_facebook.Shares = facebook_post_Map.Data[index_message].Shares
	// 		data_facebook.Story = facebook_post_Map.Data[index_message].Story
	// 		data_facebook.Status_type = facebook_post_Map.Data[index_message].Status_type
	// 		data_facebook.Created_time = facebook_post_Map.Data[index_message].Created_time
	// 		data_facebook.Updated_time = facebook_post_Map.Data[index_message].Updated_time
	// 		data_facebook.Created_at = time.Now()
	// 		data_facebook.Updated_at = time.Now()
	// 		// fmt.Println(data_facebook)

	// 		colQuerier := bson.M{"post_id": data_facebook.Post_id} //* [find post_id in Collection]
	// 		count, err := db.Find(colQuerier).Count()
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		// log.Println("count : ", count, "Post_id : ", data_facebook.Post_id)
	// 		if count > 0 {
	// 			log.Println("Post_id : ", data_facebook.Post_id, "Unique")
	// 		}
	// 		if count == 0 {
	// 			err = db.Insert(data_facebook)
	// 			if err != nil {
	// 				if mgo.IsDup(err) {
	// 					panic(err)
	// 				}
	// 				panic(err)
	// 			} else {
	// 				log.Println("Created Post_id : ", data_facebook.Post_id)
	// 				log.Println("Time : ", data_facebook.Created_time.String())

	// 				// elasticsearch
	// 				// wg.Add(1)

	// 				// // go func(i int data string) {
	// 				// defer wg.Done()

	// 				// Build the request body.
	// 				var b strings.Builder
	// 				b.WriteString(`{"post_id" : "`)
	// 				b.WriteString(data_facebook.Post_id)
	// 				b.WriteString(`","create_time" : "`)
	// 				b.WriteString(data_facebook.Created_time.String())
	// 				b.WriteString(`"}`)

	// 				elastic := Data_facebook_elastic{}
	// 				elastic.Post_id = data_facebook.Post_id
	// 				elastic.Created_time = data_facebook.Created_time.String()

	// 				// e, _ := json.Marshal(elastic)
	// 				// dec := json.NewDecoder(elastic)

	// 				// Set up the request object.
	// 				req := esapi.IndexRequest{
	// 					Index:      "facebook",
	// 					DocumentID: data_facebook.Post_id, //strconv.Itoa(0 + 1),
	// 					Body:       strings.NewReader(b.String()),
	// 					// strings.NewReader(b.String()),
	// 					Refresh: "true",
	// 				}

	// 				// Perform the request with the client.
	// 				res, err := req.Do(context.Background(), es)
	// 				if err != nil {
	// 					log.Fatalf("Error getting response: %s", err)
	// 				}
	// 				defer res.Body.Close()

	// 				if res.IsError() {
	// 					log.Printf("[%s] Error indexing document ID=%d", res.Status(), data_facebook.Post_id)
	// 				} else {
	// 					// Deserialize the response into a map.
	// 					var r map[string]interface{}
	// 					if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
	// 						log.Printf("Error parsing the response body: %s", err)
	// 					} else {
	// 						// Print the response status and indexed document version.
	// 						log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
	// 					}
	// 				}

	// 				// }
	// 			}
	// 			// wg.Wait()
	// 		}
	// 		log.Println(strings.Repeat("-", 50))

	//	}

	// }
}

func (rabbit Rabbit) Connect() *amqp.Connection {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s", rabbit.Username, rabbit.Password, rabbit.Host, rabbit.Port)
	//graylog2
	conn, err := amqp.Dial(url)
	if err != nil {
		errorMessage := fmt.Sprintf("Cannot connect :%s", err)
		panic(errorMessage)
	}
	return conn
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Loop_Message(msgs <-chan amqp.Delivery) {
	// facebookMap := make(map[string]interface{})
	Social_Data_FacebookMap := Social_Data_Facebook{}
	for data := range msgs {
		// printlog page_id and access_token
		fmt.Println("----------------------------------------")
		log.Printf("[x] Receive")

		// change json data.Body to struct facebookMap
		json.Unmarshal([]byte(data.Body), &Social_Data_FacebookMap)
		Social_Data_FacebookMap.Created_at = time.Now()
		Social_Data_FacebookMap.Updated_at = time.Now()
		// fmt.Println(Social_Data_FacebookMap.Page_id)
		// fmt.Println(Social_Data_FacebookMap.Page_name)
		// fmt.Println(Social_Data_FacebookMap.Post_id)
		// fmt.Println(Social_Data_FacebookMap.Message)
		// fmt.Println(Social_Data_FacebookMap.Full_picture)
		// fmt.Println(Social_Data_FacebookMap.Shares)
		// fmt.Println(Social_Data_FacebookMap.Story)
		// fmt.Println(Social_Data_FacebookMap.Status_type)
		// fmt.Println(Social_Data_FacebookMap.Created_time)
		// fmt.Println(Social_Data_FacebookMap.Updated_time)
		// fmt.Println(Social_Data_FacebookMap.Created_at)
		// fmt.Println(Social_Data_FacebookMap.Updated_at)

		cfg := elasticsearch.Config{
			Addresses: []string{
				"http://localhost:32774",
			},
			Username: "foo",
			Password: "bar",
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   100,
				ResponseHeaderTimeout: time.Second,
				DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS11,
				},
			},
		}
		es, _ := elasticsearch.NewClient(cfg)

		session := ConnectMongoDb()
		session.SetMode(mgo.Monotonic, true)
		defer session.Close()
		db := session.DB("social-data-facebook").C("steram-page")

		colQuerier := bson.M{"post_id": Social_Data_FacebookMap.Post_id} //* [find post_id in Collection]
		count, err := db.Find(colQuerier).Count()
		if err != nil {
			panic(err)
		}
		if count > 0 {
			log.Println("Post_id : ", Social_Data_FacebookMap.Post_id, "Unique")
		}
		if count == 0 {
			err = db.Insert(Social_Data_FacebookMap)
			if err != nil {
				if mgo.IsDup(err) {
					panic(err)
				}
				panic(err)
			} else {
				log.Println("Created Post_id : ", Social_Data_FacebookMap.Post_id)
				var b strings.Builder
				b.WriteString(`{"post_id" : "`)
				b.WriteString(Social_Data_FacebookMap.Post_id)
				b.WriteString(`","create_time" : "`)
				b.WriteString(Social_Data_FacebookMap.Created_time.String())
				b.WriteString(`"}`)

				elastic := Data_facebook_elastic{}
				elastic.Post_id = Social_Data_FacebookMap.Post_id
				elastic.Created_time = Social_Data_FacebookMap.Created_time.String()

				// Set up the request object.
				req := esapi.IndexRequest{
					Index:      "facebook",
					DocumentID: Social_Data_FacebookMap.Post_id, //strconv.Itoa(0 + 1),
					Body:       strings.NewReader(b.String()),
					Refresh:    "true",
				}

				// Perform the request with the client.
				res, err := req.Do(context.Background(), es)
				if err != nil {
					log.Fatalf("Error getting response: %s", err)
				}
				defer res.Body.Close()

				if res.IsError() {
					log.Printf("[%s] Error indexing document ID=%d", res.Status(), Social_Data_FacebookMap.Post_id)
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

			}
		}
		log.Println(strings.Repeat("-", 50))
	}
}

func ConnectMongoDb() *mgo.Session {
	session, err := mgo.Dial("mongodb://127.0.0.1:32776/social-data-facebook")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	return session
}

// func ConnectElastic() *elasticsearch.Config {
// 	cfg := elasticsearch.Config{
// 		Addresses: []string{
// 			"http://localhost:32774",
// 		},
// 		Username: "foo",
// 		Password: "bar",
// 		Transport: &http.Transport{
// 			MaxIdleConnsPerHost:   100,
// 			ResponseHeaderTimeout: time.Second,
// 			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
// 			TLSClientConfig: &tls.Config{
// 				MinVersion: tls.VersionTLS11,
// 			},
// 		},
// 	}

// 	return cfg
// }
