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
	"encoding/hex"
)

var connMainnet *grpc.ClientConn
var errMainnet error
var connTestnet *grpc.ClientConn 
var errTestnet error
var clientMainnet qrl_proto.PublicAPIClient 
var clientTestnet qrl_proto.PublicAPIClient

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

	// Serialize protobuf to json
	GetStatsRespMainnetSerialized, err := protojson.Marshal(GetStatsRespMainnet)
	if err != nil {
		glog.Fatalf("%v",err)
	}
	return string(GetStatsRespMainnetSerialized)
}

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

	// Serialize protobuf to json
	GetStatsRespTestnetSerialized, err := protojson.Marshal(GetStatsRespTestnet)
	if err != nil {
		glog.Fatalf("%v",err)
	}
	return string(GetStatsRespTestnetSerialized)
}

func GetBalanceMainnet(qaddress string) string {
	if errMainnet != nil {
		glog.Fatalf("failed connecting to server: %s", errMainnet)
	}

	// Create a User service Client on the connection.
	ctxMainnet, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	binAddress, err := hex.DecodeString(qaddress[1:])

	// GetBalance Mainnet
	GetBalanceRespMainnet, err := clientMainnet.GetBalance(ctxMainnet, &qrl_proto.GetBalanceReq{Address:binAddress})
	if err != nil {
		glog.Fatalf("%v",err)
	}

	// Serialize protobuf to json
	GetBalanceRespMainnetSerialized, err := protojson.Marshal(GetBalanceRespMainnet)
	if err != nil {
		glog.Fatalf("%v",err)
	}
	return string(GetBalanceRespMainnetSerialized)
}

func GetBalanceTestnet(qaddress string) string {
	if errMainnet != nil {
		glog.Fatalf("failed connecting to server: %s", errMainnet)
	}

	// Create a User service Client on the connection.
	ctxTestnet, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	binAddress, err := hex.DecodeString(qaddress[1:])

	// GetBalance Testnet
	GetBalanceRespTestnet, err := clientTestnet.GetBalance(ctxTestnet, &qrl_proto.GetBalanceReq{Address:binAddress})
	if err != nil {
		glog.Fatalf("%v",err)
	}

	// Serialize protobuf to json
	GetBalanceRespTestnetSerialized, err := protojson.Marshal(GetBalanceRespTestnet)
	if err != nil {
		glog.Fatalf("%v",err)
	}
	return string(GetBalanceRespTestnetSerialized)
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
		mainnet.GET("/GetBalance/:qaddress", func(c *gin.Context) {
			qaddress := c.Param("qaddress")
			c.Writer.Header().Set("Content-Type", "application/json")
			c.String(http.StatusOK, "%v", GetBalanceMainnet(qaddress))
		})	
	}

	// Simple group: testnet
	testnet := router.Group("/testnet")
	{
		testnet.GET("/GetStats", func(c *gin.Context) {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.String(http.StatusOK, "%v", GetStatsTestnet())
		})
		testnet.GET("/GetBalance/:qaddress", func(c *gin.Context) {
			qaddress := c.Param("qaddress")
			c.Writer.Header().Set("Content-Type", "application/json")
			c.String(http.StatusOK, "%v", GetBalanceTestnet(qaddress))
		})	
	}

	router.Run(":8080")
}