package tron

import (
	"testing"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"google.golang.org/grpc"
)

const test_TronGrpcNode string = "grpc.nile.trongrid.io:50051"
const test_TronGrpcNodeWrong string = "bar.com:50051"

func Test_NewTronClient(t *testing.T) {
	// connection to correct url
	c, err := NewTronClient(test_TronGrpcNode, grpc.WithInsecure())
	defer c.grpc.Stop()
	if err != nil {
		t.Fatalf("Can't connect to node url %s: %v", test_TronGrpcNode, err)
	}
}

func Test_SetTimeout(t *testing.T) {
	var c *TronClient
	err := c.SetTimeout(time.Second, grpc.WithInsecure())
	if err == nil {
		t.Fatalf("SetTimeout shouldn't use nil TronClient")
	}
	c = new(TronClient)
	c.url = test_TronGrpcNode
	err = c.SetTimeout(time.Second, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Can't connect to node url %s: %v", test_TronGrpcNode, err)
	}
}
func Test_keepConnect(t *testing.T) {
	var c *TronClient
	err := c.keepConnect()
	if err == nil {
		t.Fatal("keepConnect shouldn't use nil TronClient")
	}
	c = new(TronClient)
	c.url = test_TronGrpcNodeWrong
	err = c.keepConnect()
	if err == nil {
		t.Fatal("keepConnect shouldn't use unitialized TronClient")
	}
	c.grpc = client.NewGrpcClient(test_TronGrpcNodeWrong)
	err = c.grpc.Start(grpc.WithInsecure())
	defer c.grpc.Stop()
	if err != nil {
		t.Fatalf("Can't start the grpc client: %v", err)
	}
	err = c.keepConnect()
	if err == nil {
		t.Fatalf("keepConnect shouldn't connect to wrong node url")
	}
}