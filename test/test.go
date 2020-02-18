package main

import (
	"fmt"

	// "encoding/json"

	"github.com/go-redis/redis"
)

type Author struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// json, err := json.Marshal(Author{Name: "Elliot", Age: 25})
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// author := &Author{
	// 	Name: "Elliot",
	// 	Age:  25,
	// }

	// err = client.Set("id1234", json, 0).Err()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// val, err := client.Get("id1234").Result()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(val)

	// err = client.HMSet("keyqq", author).Err()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// val, err := client.HGET()

	//Init a map[string]interface{}
	var m = make(map[string]interface{})
	m["id"] = "13124125124"
	m["content"] = "herro"

	hash, err := client.HMSet("i", m).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(hash)

	vars, err := client.HMGet("i", "id").Result()
	// fmt.Printf("%q %s\n", v, err)
	fmt.Println(vars, err)

	var fb = make(map[string]interface{})
	fb["page_id"] = "100869428117671"
	fb["access_token"] = "12123123123"
	key := "fb_page_list"

	hash_fb, err := client.HMSet(key, fb).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(hash_fb)

	var_fb, err := client.HMGet(key, "page_id", "access_token").Result()
	// fmt.Printf("%q %s\n", v, err)
	fmt.Println(var_fb, err)

	// fmt.Println(m)
	// v, err := client.Do("HMGet", "i").String()

	// err = client.HVals("i").Err()
	// err = client.HGet("i", "1)").Err()
	// fmt.Println(err)
	// if err != nil {
	// 	//If the counter is not set, set it to 0
	// 	// err := client.Set("i", "0", 0).Err() // set = update
	// 	// if err != nil {
	// 	// 	panic(err)
	// 	// }
	// }
	// fmt.Println(var)
	// } else {
	// 	//Else, increment it
	// 	counter, err := client.Incr("i").Result()
	// 	//Cast it to string
	// 	uid = strconv.FormatInt(index, 10)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	//Set new counter
	// 	err = client.Set("i", counter, 0).Err()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
}
