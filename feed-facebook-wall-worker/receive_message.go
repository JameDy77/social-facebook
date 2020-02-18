package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

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
	rabbit := Rabbit{
		Host:      "localhost",
		Port:      "32771",
		Username:  "guest",
		Password:  "guest",
		QueueName: "fbfeed.page.add",
	}
	rabbit.Connection = rabbit.Connect()
	rabbit.Channel, _ = rabbit.Connection.Channel() // _ = ไม่ใช้ตัวแปรนั้น

	msgs, _ := rabbit.Channel.Consume(
		rabbit.QueueName, // queue
		"",               // consumer
		true,             // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	// fmt.Println(msgs)
	// failOnError(err, "Failed to register a consumer")
	personMap := make(map[string]interface{})
	for data := range msgs {
		fmt.Println("----------------------------------------")
		log.Printf(" [x] Receive")
		// fmt.Println(string(data.Body))
		json.Unmarshal([]byte(data.Body), &personMap)
		// fmt.Println(personMap)
		fmt.Println("page_id : ", personMap["page_id"])
		fmt.Println("access_token : ", personMap["access_token"])
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
