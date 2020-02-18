package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis" // Redis
	"github.com/streadway/amqp" // RabbitMQ
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

type Redis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	Connect  *redis.Client
}

func main() {
	time_wait := 15 * 60
	key := "fb_page_list"
	access_token := "EAACY2TDTj4kBAPFcGcHaPsF6V6ldbj3iqU72EanY6W142ImRSg51ZAJtGlUZCDIefhixs27fEruqpovqiSPQ7kTWaunZBgrYZAH221DmyNZAhn9m3x5Vw5cU0uL4grOyuNemGTJMx3YZAxnP54LZAgGBuila8WTU7nylZC3ZCSGRgYEfsn8bs9c64eVAdUhJFxueNWWsjlwQ33ZBWB7MmWrgaN"
	page_id := "100869428117671"

	LoopMessage(time_wait, key, access_token, page_id)
}

func LoopMessage(time_wait int, key string, access_token string, page_id string) {
	for {
		// Redis Connect(Client)
		redis_struct := Redis{
			Host:     "localhost",
			Port:     "32777", // docker port | 6378 default port
			Password: "",
			DB:       0,
		}
		redis_struct.Connect = redis_struct.ConnectRedis()
		client := redis_struct.Connect

		//create fb_page_list is map[string]interface for save page_id and access_token
		var fb_page_list = make(map[string]interface{})
		fb_page_list["page_id"] = page_id
		fb_page_list["access_token"] = access_token

		// insert data [key,value] in redis
		Insert_Key_Value(client, key, fb_page_list)

		// get access_token in redis
		value_list := Get_Key_Value(client, key)

		// change map[string]interface to json
		jsonSendMessage, _ := json.Marshal(value_list) // Data from Redis

		// send Message
		SendMessage(jsonSendMessage)

		// wait send message = time_wait (min)
		time.Sleep(time.Duration(time_wait) * time.Second)
	}
}

func Insert_Key_Value(client *redis.Client, key string, list_data map[string]interface{}) { // insert key field value [field value ...]
	err := client.HMSet(key, list_data).Err()
	fmt.Println("---------------------------------------------------------------")
	if err != nil {
		// panic(err)
		log.Println("Can't Insert Data in Redis")
	} else {
		log.Println("Insert Data in Redis Successfully")
	}
}

func Get_Key_Value(client *redis.Client, key string) map[string]string { // get key field value
	value_list, err := client.HGetAll(key).Result()
	if err != nil {
		// panic(err)
		log.Println("Can't Select Data in Redis")
	} else {
		log.Println("Get Value : ", value_list)
	}
	return value_list
}

func SendMessage(jsonSendMessage []byte) {
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
	err := rabbit.Channel.Publish(
		// "",               // exchange
		"ha_fbfeed",      // exchange
		rabbit.QueueName, // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/json", // type = json
			Body:        jsonSendMessage,
		})
	log.Printf(" [x] SentMessage %s", jsonSendMessage)
	failOnError(err, "Failed to publish a message")
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

func (redis_struct Redis) ConnectRedis() *redis.Client {
	addr := fmt.Sprintf("%s:%s", redis_struct.Host, redis_struct.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redis_struct.Password,
		DB:       redis_struct.DB,
	})
	return client
}
