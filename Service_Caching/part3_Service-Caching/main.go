package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	//proto "RPC/proto"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func rClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return client
}

func main() {

	// creates a client
	client := rClient()
	//pool := client
	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err, pong)
	}

	// Using the SET command to set Key-value pair

	g := gin.Default()
	g.GET("/Getuser/:name", func(ctx *gin.Context) {
		var name string = ctx.Param("name")

		//req := &proto.Request{Name: string(name)}
		if value := get(client, name); value != "Sorry" {
			ctx.JSON(http.StatusOK, gin.H{
				"Rupin:" + name: fmt.Sprint(value),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Sorry Failed to retrieve value"})
		}
	})

	g.GET("/Setuser/:key/:value", func(ctx *gin.Context) {
		var key string = ctx.Param("key")
		var value string = ctx.Param("value")

		// 	err = set(client)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		//req := &proto.Request{Name: string(name)}
		if status := set(client, key, value); status == "OK" {
			ctx.JSON(http.StatusOK, gin.H{
				"Success": "Value has been set like this---- Rupin:" + key + "=" + value})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Sorry Failed to retrieve value"})
		}
	})

	if err := g.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

	// check connection status

	// Using the GET command to get values from keys

}

func get(client *redis.Client, toSearch string) string {

	nameVal, err := client.Get(fmt.Sprintf("Rupin:%s", toSearch)).Result()
	if err != nil {
		return "Sorry"
	}
	return nameVal
	// fmt.Println("name", nameVal)

	// countryVal, err := client.Get("Road").Result()
	// if err == redis.Nil {
	// 	fmt.Println("no value found")
	// } else if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Println("country", countryVal)
	// }

	// return nil
}

func set(client *redis.Client, key string, value string) string {
	var myname strings.Builder
	myname.WriteString("Rupin:")
	myname.WriteString(key)
	err := client.Set(myname.String(), value, 0).Err()
	if err != nil {
		return "Sorry"
	}
	return "OK"
}
