package main

import (
	"fmt"
	"time"

	fb "github.com/huandu/facebook"
)

type stream_facebook struct {
	Page_id   string `json:"page_id"`
	Page_name string `json:"page_name"`
	Data      []data_facebook
	// Message      string    `json:"message"`
	// Created_time time.Time `json:"created_time"`
}

type data_facebook struct {
	Post_Id, Message string
	Created_time     time.Time
}

func main() {
	page_id := "100869428117671"
	access_token := "EAACY2TDTj4kBACRczPeXuxOqOPCY3ZAoE9ZCZAUEWUaKj4jOtYoVr0upRuWTVN1gDOZBc6FZBrCPapZAYYWvzzshSseDw1iHRjG8IDLzrkSyjLNZA1sAcpbxFozARUZAHuqVBsJnKmZBA8bZAizddj16h0qjPSQkgMfaQ4GP6zoQZAQxs28GLS06iEc5pVFIBViiptVW6SgJ9NkDfLTps7ufonD"
	res, _ := fb.Get("/"+page_id, fb.Params{
		// "fields":       "name,country_page_likes,feed.limit(2)",
		"fields":       "id,name,feed",
		"access_token": access_token,
	})
	// fmt.Println(res)
	fmt.Println("Here is my Facebook first name:", res["name"])
	page_name := fmt.Sprintf("%v", res["name"])
	fmt.Println("Here is page_likes:", res["country_page_likes"])
	// fmt.Println(res["feed"])
	// data := res["feed"]
	// for key, value := range res {
	// 	fmt.Println("Key:", key, "Value:", value)
	// }
	// fmt.Println(res["feed"])

	// jsonString, _ := json.Marshal(res["feed"])
	// fmt.Println(string(jsonString))
	// [""].([]int)[index]

	// if value, has := mapstring["a"].([]int);has{
	// fmt.Println(value[3])
	// }
	// fmt.Println("----------------------------------")
	// fmt.Println(res["feed"].(map[string]interface{})["data"])
	// fmt.Println("----------------------------------")
	// data2 := res["feed"].(map[string]interface{})["data"].([]interface{})
	// fmt.Println(data2)

	// for key, value := range data2[0].(map[string]interface{}) {
	// 	fmt.Println(key, value)
	// }

	fmt.Println("----------------------__________________-")
	var data_all []data_facebook
	// var time_check string
	if value_feed, has := res["feed"]; has {
		if value_data, has := value_feed.(map[string]interface{})["data"].([]interface{}); has {
			for i := 0; i < len(value_data); i++ {
				data := data_facebook{}
				// fmt.Println(value_data[i].(map[string]interface{})["created_time"])
				value := value_data[i].(map[string]interface{})["created_time"]
				value_str := fmt.Sprintf("%v", value)
				format := "2006-01-02T15:04:05+0000"
				location, _ := time.LoadLocation("Asia/Bangkok")
				createdAt, _ := time.Parse(format, value_str)
				createdAt = createdAt.In(location)

				expiresAt := time.Now()
				diff := expiresAt.Sub(createdAt)

				if diff.Hours() <= 24 { // check stream 16 hr
					fmt.Println("----------------")
					// fmt.Println(key, " : ", createdAt, " : ", value)
					data.Created_time = createdAt
					post_id := fmt.Sprintf("%v", value_data[i].(map[string]interface{})["id"])
					data.Post_Id = post_id
					message := fmt.Sprintf("%v", value_data[i].(map[string]interface{})["message"])
					data.Message = message
					// check_created_time = true
					data_all = append(data_all, data)
				}

				// check_created_time := false
				// for key, value := range value_data[i].(map[string]interface{}) {

				// 	if key == "created_time" {
				// 		// fmt.Println("----------------")
				// 		// fmt.Println(key, " : ", value)

				// 		// setup a start and end time
				// 		value_str := fmt.Sprintf("%v", value)
				// 		format := "2006-01-02T15:04:05+0000"
				// 		location, _ := time.LoadLocation("Asia/Bangkok")
				// 		createdAt, _ := time.Parse(format, value_str)
				// 		createdAt = createdAt.In(location)

				// 		expiresAt := time.Now()
				// 		// fmt.Println("time : ", expiresAt)
				// 		// fmt.Println("Hour : ", expiresAt.Hour())

				// 		// fmt.Println(reflect.TypeOf(createdAt), createdAt)
				// 		// fmt.Println(reflect.TypeOf(expiresAt), expiresAt)

				// 		// get the diff
				// 		diff := expiresAt.Sub(createdAt)
				// 		// fmt.Printf("Lifespan is %+v\n", diff)
				// 		// fmt.Println(diff.Hours())
				// 		// fmt.Println(reflect.TypeOf(diff))

				// 		if diff.Hours() <= 24 { // check stream 16 hr
				// 			fmt.Println("----------------")
				// 			fmt.Println(key, " : ", createdAt, " : ", value)
				// 			data.Created_time = createdAt
				// 			check_created_time = true
				// 		}
				// 	}
				// 	if key == "message" {
				// 		if check_created_time == true {
				// 			fmt.Println("message : ", value)
				// 			data.Message = fmt.Sprintf("%v", value)
				// 		}
				// 	}
				// 	if key == "id" {
				// 		if check_created_time == true {
				// 			fmt.Println("id : ", value)
				// 			data.Id = fmt.Sprintf("%v", value)
				// 		}
				// 	}

				// }

			}

		}
	}
	fmt.Println("----------------------__________________-")
	// fmt.Println(data_all)

	stream_facebook := stream_facebook{
		Page_id:   page_id,
		Page_name: page_name,
		Data:      data_all,
	}

	fmt.Println(stream_facebook)
}
