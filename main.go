package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/golang/glog"
    "golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/theQRL/walletd-rest-proxy/qrl_proto"
    "time"
	"google.golang.org/protobuf/encoding/protojson"
)


func main() {
	// Open a connection to the server.
	conn, err := grpc.Dial("testnet-1.automated.theqrl.org:19009", grpc.WithInsecure())
	if err != nil {
		glog.Fatalf("failed connecting to server: %s", err)
	}
	defer conn.Close()
 
	// Create a User service Client on the connection.
	client := qrl_proto.NewPublicAPIClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// GetStats
	response, err := client.GetStats(ctx, &qrl_proto.GetStatsReq{})
	if err != nil {
		glog.Fatalf("%v",err)
	}

	// Serialize protobuf to json
	b, err := protojson.Marshal(response)
	if err != nil {
		glog.Fatalf("%v",err)
	}

	router := gin.Default()

	// Simple group: mainnet
	mainnet := router.Group("/mainnet")
	{
		mainnet.GET("/GetStats", func(c *gin.Context) {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.String(http.StatusOK, "%v", string(b))
		})	
	}

	// Simple group: testnet
	testnet := router.Group("/testnet")
	{
		testnet.GET("/GetStats", func(c *gin.Context) {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.String(http.StatusOK, "%v", string(b))
		})	
	}

	router.Run(":8080")
}