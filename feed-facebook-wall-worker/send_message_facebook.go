package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	fb "github.com/huandu/facebook"
	"github.com/streadway/amqp"
)

// struct post for stream_facebook
type stream_facebook struct {
	Page_id   string `json:"page_id"`
	Page_name string `json:"page_name"`
	Data      []data_facebook
}

// struct data for postpage
type data_facebook struct {
	Post_Id, Message string
	Created_time     time.Time
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
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
	// facebook
	page_id := "100869428117671"
	access_token := "EAACY2TDTj4kBABZBv3PZBbPegcrpYVoqxyYLeEn4AZAZAn9FEZBphb8JavhiE0ZChTDwYEaaMB7SGBCsFyTqIKRTW3w0B8Pivc1AlWsJJmFHjtQ3ZCJtp4xmv8jKq9KI7ECY2TXndRm3vJDHX1E6aNf6st3WaE5vdTrqtkHZA53km2iFwfA26EERZAdbVZCoZB7LFyioyykWLJQAgZDZD"
	res, _ := fb.Get("/"+page_id, fb.Params{
		// "fields":       "name,country_page_likes,feed.limit(2)",
		"fields":       "id,name,feed",
		"access_token": access_token,
	})
	page_name := fmt.Sprintf("%v", res["name"])
	var stream_page []data_facebook
	if value_feed, has := res["feed"]; has {
		if value_data, has := value_feed.(map[string]interface{})["data"].([]interface{}); has {
			for i := 0; i < len(value_data); i++ {
				data := data_facebook{}
				value := value_data[i].(map[string]interface{})["created_time"]
				value_str := fmt.Sprintf("%v", value)
				format := "2006-01-02T15:04:05+0000"
				location, _ := time.LoadLocation("Asia/Bangkok")
				createdAt, _ := time.Parse(format, value_str)
				createdAt = createdAt.In(location)

				expiresAt := time.Now()
				diff := expiresAt.Sub(createdAt)

				if diff.Hours() <= 24 { // check data stream 24 hr
					fmt.Println("----------------")
					data.Created_time = createdAt
					post_id := fmt.Sprintf("%v", value_data[i].(map[string]interface{})["id"]) // * change interface to string
					data.Post_Id = post_id
					message := fmt.Sprintf("%v", value_data[i].(map[string]interface{})["message"]) // * change interface to string
					data.Message = message
					stream_page = append(stream_page, data)
				}
			}
		}
	}

	stream_facebook := stream_facebook{
		Page_id:   page_id,
		Page_name: page_name,
		Data:      stream_page,
	}
	// fmt.Println(stream_facebook)

	// end
	// start rabbit
	rabbit := Rabbit{
		Host:      "localhost",
		Port:      "32771",
		Username:  "guest",
		Password:  "guest",
		QueueName: "fbfeed.post.add",
	}
	rabbit.Connection = rabbit.Connect()
	rabbit.Channel, _ = rabbit.Connection.Channel() // _ = ไม่ใช้ตัวแปรนั้น

	// body := "Hello World!22224341323123"
	jsonString, _ := json.Marshal(stream_facebook) // change stream_facebook struct to json

	err := rabbit.Channel.Publish(
		// "",               // exchange
		"ha_fbfeed",      // exchange
		rabbit.QueueName, // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			// ContentType: "text/plain",
			// Body:        []byte(body),
			ContentType: "application/json",
			Body:        jsonString,
		})
	// log.Printf(" [x] Sent %s", body)
	log.Printf(" [x] Sent %s", stream_facebook)
	failOnError(err, "Failed to publish a message")
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

func (rabbit Rabbit) Publish(ch *amqp.Channel) {
	ch, err := rabbit.Connection.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
}
