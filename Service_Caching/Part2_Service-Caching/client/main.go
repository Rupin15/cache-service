package main

import (
	proto "Part2_Service-Caching/proto"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewUserServiceClient(conn)

	g := gin.Default()
	g.GET("/user/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")

		req := &proto.Request{Name: string(name)}
		if response, err := client.GetUserByName(ctx, req); err == nil {
			ctx.JSON(http.StatusOK, gin.H{
				"Name":  fmt.Sprint(response.Name),
				"Roll":  fmt.Sprint(response.Roll),
				"Class": fmt.Sprint(response.Class),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	})

	if err := g.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
