package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis" // Redis
	"github.com/streadway/amqp"
)

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
	// Redis Client
	client := redis.NewClient(&redis.Options{
		// Addr:     "localhost:6379",
		Addr:     "localhost:32777", // docker
		Password: "",                // no password set
		DB:       0,                 // use default DB
	})
	// Rabbit Connect
	rabbit := Rabbit{
		Host:      "localhost",
		Port:      "32771",
		Username:  "guest",
		Password:  "guest",
		QueueName: "fbfeed.page.add",
	}
	rabbit.Connection = rabbit.Connect()
	rabbit.Channel, _ = rabbit.Connection.Channel() // _ = ไม่ใช้ตัวแปรนั้น
	count_time := 0

	// loop time wait 1 min
	for {
		// Redis
		key := "fb_page_list"
		access_token := "EAACY2TDTj4kBAPFcGcHaPsF6V6ldbj3iqU72EanY6W142ImRSg51ZAJtGlUZCDIefhixs27fEruqpovqiSPQ7kTWaunZBgrYZAH221DmyNZAhn9m3x5Vw5cU0uL4grOyuNemGTJMx3YZAxnP54LZAgGBuila8WTU7nylZC3ZCSGRgYEfsn8bs9c64eVAdUhJFxueNWWsjlwQ33ZBWB7MmWrgaN"
		// set access_token
		var fb = make(map[string]interface{})
		fb["page_id"] = "100869428117671"
		fb["access_token"] = access_token
		err := client.HMSet(key, fb).Err()
		if err != nil {
			panic(err)
		}
		// get access_token
		var_fb_check, err := client.HGetAll(key).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println(len(var_fb_check))
		fmt.Println(var_fb_check)

		// if len(var_fb_check) == 0 { // insert when no key value
		// 	var fb = make(map[string]interface{})
		// 	fb["page_id"] = "100869428117671"
		// 	fb["access_token"] = access_token
		// 	err = client.HMSet(key, fb).Err()
		// 	if err != nil {
		// 		panic(err)
		// 	}

		// } else if count_time%3 == 0 || count_time == 0 { // update
		// 	var fb = make(map[string]interface{})
		// 	fb["page_id"] = "100869428117671"
		// 	fb["access_token"] = access_token
		// 	err = client.HMSet(key, fb).Err()
		// 	if err != nil {
		// 		panic(err)
		// 	}

		// } else if len(var_fb_check) == 2 {
		// 	var fb = make(map[string]interface{})
		// 	fb["page_id"] = "100869428117671"
		// 	fb["access_token"] = access_token
		// 	err = client.HMSet(key, fb).Err()
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// }

		// else {
		// 	fmt.Println(var_fb_check["page_id"])
		// }

		// end
		// Rabbit
		// body := "Hello World!22224341323123"
		jsonString, _ := json.Marshal(var_fb_check) // Redis

		err = rabbit.Channel.Publish(
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
		fmt.Println("----------------------------------------")
		log.Printf(" [x] Sent %s", var_fb_check)
		failOnError(err, "Failed to publish a message")

		time.Sleep(900 * time.Second) // wait send message 1 min
		count_time += 1
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
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
