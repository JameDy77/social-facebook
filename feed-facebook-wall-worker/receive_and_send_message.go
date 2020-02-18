package main

import (
	"encoding/json"
	"fmt"
	"log"
	// "reflect"
	"strconv"
	"time"

	fb "github.com/huandu/facebook"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// send data form page post
type Social_Data_Facebook struct {
	Post_id      string    `json:"post_id"`
	Page_id      string    `json:"page_id"`
	Page_name    string    `json:"page_name"`
	Message      string    `json:"message"`
	Full_picture string    `json:"full_picture"`
	Shares       int       `json:"shares"`
	Story        string    `json:"story"`
	Status_type  string    `json:"status_type"`
	Created_time time.Time `json:"created_time"`
	Updated_time time.Time `json:"updated_time"`
}

// struct post for stream_facebook
type stream_facebook struct {
	Page_id   string `json:"page_id"`
	Page_name string `json:"page_name"`
	Data      []data_facebook
}

// struct data for postpage
type data_facebook struct {
	Post_Id, Message, Full_Picture, Story, Status_type string
	Shares                                             int
	Created_time, Updated_time                         time.Time
}

// struct rabbit
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
	// setup receive_message from send_message [page_id,access_token] => queuename = fbfeed.page.add
	rabbit_receive := Rabbit{
		Host:      "localhost",
		Port:      "32771",
		Username:  "guest",
		Password:  "guest",
		QueueName: "fbfeed.page.add",
	}
	rabbit_receive.Connection = rabbit_receive.Connect()            //	connect rabbit
	rabbit_receive.Channel, _ = rabbit_receive.Connection.Channel() // _ = ไม่ใช้ตัวแปรนั้น connect channel rabbit
	// setup channel
	msgs, err := rabbit_receive.Channel.Consume(
		rabbit_receive.QueueName, // queue
		"",                       // consumer
		true,                     // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,                      // args
	)
	failOnError(err, "Failed to register a consumer")
	// setup send_message to folder social-data-facebook => queuename = fbfeed.post.add
	rabbit_send := Rabbit{
		Host:      "localhost",
		Port:      "32771",
		Username:  "guest",
		Password:  "guest",
		QueueName: "fbfeed.post.add",
	}
	rabbit_send.Connection = rabbit_send.Connect()            // connect rabbit
	rabbit_send.Channel, _ = rabbit_send.Connection.Channel() // _ = ไม่ใช้ตัวแปรนั้น connect channel rabbit
	//
	hour_stream := 240.00
	Loop_Message(msgs, rabbit_send, hour_stream)

	// for data := range msgs {
	// 	fmt.Println("----------------------------------------")
	// 	log.Printf("[x] Receive")
	// 	json.Unmarshal([]byte(data.Body), &facebookMap)
	// 	log.Println("page_id : ", facebookMap["page_id"])
	// 	log.Println("access_token : ", facebookMap["access_token"])

	// 	// facebook
	// 	page_id := fmt.Sprintf("%v", facebookMap["page_id"])
	// 	res, _ := fb.Get("/"+page_id, fb.Params{
	// 		"fields":       "id,name,feed{message,shares,created_time,updated_time,status_type,story,full_picture}",
	// 		"access_token": facebookMap["access_token"],
	// 	})
	// 	if res["error"] != nil {
	// 		log.Println("Please new access_totken in fill !!!!!")
	// 	} else {
	// 		// page_name := fmt.Sprintf("%v", res["name"])
	// 		var stream_page []data_facebook
	// 		if value_feed, has := res["feed"]; has {
	// 			if value_data, has := value_feed.(map[string]interface{})["data"].([]interface{}); has {
	// 				for i := 0; i < len(value_data); i++ {
	// 					data := data_facebook{}
	// 					// social_data := Social_Data_Facebook{}
	// 					value := value_data[i].(map[string]interface{})["created_time"]
	// 					value_str := fmt.Sprintf("%v", value)
	// 					format := "2006-01-02T15:04:05+0000"
	// 					location, _ := time.LoadLocation("Asia/Bangkok")
	// 					createdAt, _ := time.Parse(format, value_str)
	// 					createdAt = createdAt.In(location)

	// 					expiresAt := time.Now()
	// 					diff := expiresAt.Sub(createdAt)

	// 					if diff.Hours() >= 0 && diff.Hours() <= 240 { // check data stream 24 hr

	// 						data.Created_time = createdAt
	// 						data.Post_Id = fmt.Sprintf("%v", value_data[i].(map[string]interface{})["id"])      // * change interface to string
	// 						data.Message = fmt.Sprintf("%v", value_data[i].(map[string]interface{})["message"]) // * change interface to string
	// 						shares := fmt.Sprintf("%v", value_data[i].(map[string]interface{})["shares"])       // * change interface to string
	// 						shares_count, _ := strconv.Atoi(shares)
	// 						data.Shares = shares_count
	// 						// updated_time , status_type , story , full_picture
	// 						updated_time_string := fmt.Sprintf("%v", value_data[i].(map[string]interface{})["updated_time"])
	// 						updated_time, _ := time.Parse(format, updated_time_string)
	// 						updated_time = createdAt.In(location)
	// 						data.Updated_time = updated_time
	// 						data.Stroy = fmt.Sprintf("%v", value_data[i].(map[string]interface{})["story"])
	// 						data.Status_type = fmt.Sprintf("%v", value_data[i].(map[string]interface{})["status_type"])
	// 						data.Full_Picture = fmt.Sprintf("%v", value_data[i].(map[string]interface{})["full_picture"])
	// 						// add data to stream_page
	// 						stream_page = append(stream_page, data)

	// 						// ---------------------------------------------------------------------
	// 						// add social_data to stream fbfeed.post.add
	// 						// social_data.Created_time = createdAt
	// 						// social_data.Post_Id = fmt.Sprintf("%v", value_data[i].(map[string]interface{})["id"])      // * change interface to string
	// 						// social_data.Message = fmt.Sprintf("%v", value_data[i].(map[string]interface{})["message"]) // * change interface to string
	// 						// shares := fmt.Sprintf("%v", value_data[i].(map[string]interface{})["shares"])              // * change interface to string
	// 						// shares_count, _ := strconv.Atoi(shares)
	// 						// social_data.Shares = shares_count
	// 						// // updated_time , status_type , story , full_picture
	// 						// updated_time_string := fmt.Sprintf("%v", value_data[i].(map[string]interface{})["updated_time"])
	// 						// updated_time, _ := time.Parse(format, updated_time_string)
	// 						// updated_time = createdAt.In(location)
	// 						// social_data.Updated_time = updated_time
	// 						// social_data.Stroy = fmt.Sprintf("%v", value_data[i].(map[string]interface{})["story"])
	// 						// social_data.Status_type = fmt.Sprintf("%v", value_data[i].(map[string]interface{})["status_type"])
	// 						// social_data.Full_Picture = fmt.Sprintf("%v", value_data[i].(map[string]interface{})["full_picture"])

	// 						// Publish data => queuename = fbfeed.post.add
	// 						jsonDataStream_facebook, _ := json.Marshal(data) // change stream_facebook struct to json
	// 						err := rabbit_send.Channel.Publish(
	// 							// "",               // exchange
	// 							"ha_fbfeed",           // exchange
	// 							rabbit_send.QueueName, // routing key
	// 							false,                 // mandatory
	// 							false,                 // immediate
	// 							amqp.Publishing{
	// 								ContentType: "application/json",
	// 								Body:        jsonDataStream_facebook,
	// 							})
	// 						failOnError(err, "Failed to publish a message")
	// 					}
	// 				}
	// 			}
	// 		}
	// 		// stream_facebook := stream_facebook{
	// 		// 	Page_id:   page_id,
	// 		// 	Page_name: page_name,
	// 		// 	Data:      stream_page,
	// 		// }
	// 		// // end
	// 		// jsonStream_facebook, _ := json.Marshal(stream_facebook) // change stream_facebook struct to json
	// 		// err := rabbit_send.Channel.Publish(
	// 		// 	// "",               // exchange
	// 		// 	"ha_fbfeed",           // exchange
	// 		// 	rabbit_send.QueueName, // routing key
	// 		// 	false,                 // mandatory
	// 		// 	false,                 // immediate
	// 		// 	amqp.Publishing{
	// 		// 		ContentType: "application/json",
	// 		// 		Body:        jsonStream_facebook,
	// 		// 	})
	// 		// failOnError(err, "Failed to publish a message")
	// 	}
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

// get fields from facebook page
func FacebookResult(page_id string, access_token string) fb.Result {
	res, _ := fb.Get("/"+page_id, fb.Params{
		"fields":       "id,name,feed{message,shares,created_time,updated_time,status_type,story,full_picture}",
		"access_token": access_token,
	})
	return res
}

func Loop_Message(msgs <-chan amqp.Delivery, rabbit_send Rabbit, hour_stream float64) {
	// create facebookMap for receive_message
	facebookMap := make(map[string]interface{})
	for data := range msgs {
		// printlog page_id and access_token
		fmt.Println("----------------------------------------")
		log.Printf("[x] Receive")

		// change json data.Body to struct facebookMap
		json.Unmarshal([]byte(data.Body), &facebookMap)
		log.Println("page_id : ", facebookMap["page_id"])
		log.Println("access_token : ", facebookMap["access_token"])

		// get value page_id and access_token from facebookMap
		page_id := fmt.Sprintf("%v", facebookMap["page_id"])
		access_token := fmt.Sprintf("%v", facebookMap["access_token"])

		// find post from page
		facebook_result := FacebookResult(page_id, access_token)

		// check error field when error aceess_token
		if facebook_result["error"] != nil {
			log.Println("Please new access_token in field !!!!!")
		} else {
			Check_Field_Stream(facebook_result, hour_stream, rabbit_send)
		}
	}
}

func Check_Field_Stream(facebook_result fb.Result, hour_stream float64, rabbit_send Rabbit) {
	// check value field feed
	page_id := fmt.Sprintf("%v", facebook_result["id"])
	page_name := fmt.Sprintf("%v", facebook_result["name"])
	if value_feed, has := facebook_result["feed"]; has {
		// check value field data
		if value_data, has := value_feed.(map[string]interface{})["data"].([]interface{}); has {
			// loop find value_data
			for i := 0; i < len(value_data); i++ {
				// check value created_time
				created_timeMap := value_data[i].(map[string]interface{})["created_time"]
				created_timeString := fmt.Sprintf("%v", created_timeMap)

				// change format facebook_datetime to golang_time
				format := "2006-01-02T15:04:05+0000"
				location, _ := time.LoadLocation("Asia/Bangkok")
				createdAt, _ := time.Parse(format, created_timeString)
				createdAt = createdAt.In(location)
				// create timeNow
				expiresAt := time.Now()
				// check timeNow and createdAt
				diff := expiresAt.Sub(createdAt)

				social_data := Social_Data_Facebook{}

				if diff.Hours() >= 0 && diff.Hours() <= hour_stream { // check data stream = hour_stream (hr)
					social_data.Created_time = createdAt
					social_data.Post_id = fmt.Sprintf("%v", value_data[i].(map[string]interface{})["id"])                           // * change interface to string
					social_data.Message = fmt.Sprintf("%v", value_data[i].(map[string]interface{})["message"])                      // * change interface to string
					shares := fmt.Sprintf("%v", value_data[i].(map[string]interface{})["shares"].(map[string]interface{})["count"]) // * change interface to string
					shares_count, _ := strconv.Atoi(shares)
					social_data.Shares = shares_count

					// updated_time , status_type , story , full_picture
					updated_time_string := fmt.Sprintf("%v", value_data[i].(map[string]interface{})["updated_time"])
					updated_time, _ := time.Parse(format, updated_time_string)
					updated_time = createdAt.In(location)
					social_data.Updated_time = updated_time
					social_data.Story = fmt.Sprintf("%v", value_data[i].(map[string]interface{})["story"])
					social_data.Status_type = fmt.Sprintf("%v", value_data[i].(map[string]interface{})["status_type"])
					social_data.Full_picture = fmt.Sprintf("%v", value_data[i].(map[string]interface{})["full_picture"])

					social_data.Page_id = page_id
					social_data.Page_name = page_name

					// Publish data => queuename = fbfeed.post.add
					Send_message(social_data, rabbit_send)
				}

			}
		}
	}
}

func Send_message(social_data Social_Data_Facebook, rabbit_send Rabbit) {
	jsonDataStream_facebook, _ := json.Marshal(social_data) // change stream_facebook struct to json
	err := rabbit_send.Channel.Publish(
		// "",               // exchange
		"ha_fbfeed",           // exchange
		rabbit_send.QueueName, // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonDataStream_facebook,
		})
	failOnError(err, "Failed to publish a message")
}
