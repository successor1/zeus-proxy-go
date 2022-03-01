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

var connMainnet *grpc.ClientConn
var errMainnet error
var connTestnet *grpc.ClientConn 
var errTestnet error
var clientMainnet qrl_proto.PublicAPIClient 
var clientTestnet qrl_proto.PublicAPIClient

func GetStatsTestnet() string{
	// Open a connection to the testnet server.
	if errTestnet != nil {
		glog.Fatalf("failed connecting to server: %s", errTestnet)
	}
 
	// Create a User service Client on the connection.
	ctxTestnet, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// GetStats Testnet
	GetStatsRespTestnet, err := clientTestnet.GetStats(ctxTestnet, &qrl_proto.GetStatsReq{})
	if err != nil {
		glog.Fatalf("%v",err)
	}

	// Serialize protobuf to json Testnet
	GetStatsRespTestnetSerialized, err := protojson.Marshal(GetStatsRespTestnet)
	if err != nil {
		glog.Fatalf("%v",err)
	}
	return string(GetStatsRespTestnetSerialized)
}

func GetStatsMainnet() string {
	if errMainnet != nil {
		glog.Fatalf("failed connecting to server: %s", errMainnet)
	}

	// Create a User service Client on the connection.
	ctxMainnet, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// GetStats Mainnet
	GetStatsRespMainnet, err := clientMainnet.GetStats(ctxMainnet, &qrl_proto.GetStatsReq{})
	if err != nil {
		glog.Fatalf("%v",err)
	}

	// Serialize protobuf to json Mainnet
	GetStatsRespMainnetSerialized, err := protojson.Marshal(GetStatsRespMainnet)
	if err != nil {
		glog.Fatalf("%v",err)
	}
	return string(GetStatsRespMainnetSerialized)
}

func main() {
	// Open a connection to the mainnet server.
	connMainnet, errMainnet = grpc.Dial("mainnet-1.automated.theqrl.org:19009", grpc.WithInsecure())
	connTestnet, errTestnet = grpc.Dial("testnet-1.automated.theqrl.org:19009", grpc.WithInsecure())
	clientMainnet = qrl_proto.NewPublicAPIClient(connMainnet)
	clientTestnet = qrl_proto.NewPublicAPIClient(connTestnet)

	defer connMainnet.Close()
	defer connTestnet.Close()

	router := gin.Default()

	// Simple group: mainnet
	mainnet := router.Group("/mainnet")
	{
		mainnet.GET("/GetStats", func(c *gin.Context) {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.String(http.StatusOK, "%v", GetStatsMainnet())
		})	
	}

	// Simple group: testnet
	testnet := router.Group("/testnet")
	{
		testnet.GET("/GetStats", func(c *gin.Context) {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.String(http.StatusOK, "%v", GetStatsTestnet())
		})	
	}

	router.Run(":8080")
}