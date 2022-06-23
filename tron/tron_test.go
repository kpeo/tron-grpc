package tron

import (
	"testing"

	"google.golang.org/grpc"
)

const test_TronGrpcNodeCorrect string = "grpc.nile.trongrid.io:50051"
const test_TronGrpcNodeWrong string = "localhost:50054"

func Test_NewTronClient(t *testing.T) {
	// connection to correct url
	c, err := NewTronClient(test_TronGrpcNodeCorrect, grpc.WithInsecure())
	if err != nil {
		t.Errorf("Can't connect to correct node url %s: %v", test_TronGrpcNodeCorrect, err)
	}
	c.grpc.Stop()
	// connection to wrong url
	c, err = NewTronClient(test_TronGrpcNodeWrong, grpc.WithInsecure())
	if err == nil {
		t.Errorf("Successfuly connected to incorrect node url %s: %v", test_TronGrpcNodeWrong, err)
	}
	c.grpc.Stop();
}

func Test_SetTimeout(t *testing.T) {
	t.Errorf("Not implemented")
}
func Test_keepConnect(t *testing.T) {
	t.Errorf("Not implemented")
}