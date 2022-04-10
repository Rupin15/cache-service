package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	//proto "RPC/proto"
	"github.com/gin-gonic/gin"
	redis "github.com/gomodule/redigo/redis"
)

type User struct {
	Username string `json:"name"`
	Class    string `json:"class"`
	Roll     string `json:"roll"`
}

func newPool() *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func main() {
	pool := newPool()
	conn := pool.Get()
	defer conn.Close()
	err := ping(conn)
	if err != nil {
		fmt.Println(err)
	}

	g := gin.Default()
	g.GET("/Getuser/:name", func(ctx *gin.Context) {
		var name string = ctx.Param("name")

		//req := &proto.Request{Name: string(name)}
		if value := getStruct(conn, name); value == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"Rupin:" + name: fmt.Sprint(name),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Sorry Failed to retrieve value"})
		}
	})

	g.GET("/Setuser/:name/:class/:roll", func(ctx *gin.Context) {
		var name string = ctx.Param("name")
		var class string = ctx.Param("class")
		var roll string = ctx.Param("roll")

		// 	err = set(client)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		//req := &proto.Request{Name: string(name)}
		if status := setStruct(conn, name, class, roll); status == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"Success": "Value has been set like this---- Rupin:" + name})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Sorry Failed to retrieve value"})
		}
	})

	if err := g.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}

func ping(c redis.Conn) error {
	// Send PING command to Redis
	pong, err := c.Do("PING")
	if err != nil {
		return err
	}

	// PING command returns a Redis "Simple String"
	// Use redis.String to convert the interface type to string
	s, err := redis.String(pong, err)
	if err != nil {
		return err
	}

	fmt.Printf("PING Response = %s\n", s)
	// Output: PONG

	return nil
}

func setStruct(c redis.Conn, name string, class string, roll string) error {

	const objectPrefix string = "Rupin:"

	usr := User{
		Username: name,
		Class:    class,
		Roll:     roll,
	}

	// serialize User object to JSON
	json, err := json.Marshal(usr)
	if err != nil {
		return err
	}

	// SET object
	_, err = c.Do("SET", objectPrefix+usr.Username, json)
	if err != nil {
		return err
	}

	return nil
}

func getStruct(c redis.Conn, name string) error {

	const objectPrefix string = "Rupin:"

	username := name
	s, err := redis.String(c.Do("GET", objectPrefix+username))
	if err == redis.ErrNil {
		fmt.Println("User does not exist")
	} else if err != nil {
		return err
	}

	usr := User{}
	err = json.Unmarshal([]byte(s), &usr)

	fmt.Printf("\n\n%+v\n\n", usr)

	return nil

}
