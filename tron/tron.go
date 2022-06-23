package tron

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"google.golang.org/grpc"
)

// TronClient grpc client structure
type TronClient struct {
	url string
	grpc *client.GrpcClient
}

// NewClient creates a new grpc client
func NewTronClient(url string, opts ...grpc.DialOption) (*TronClient, error) {
	c := new(TronClient)
	c.url = url

	c.grpc = client.NewGrpcClient(url)

	err := c.grpc.Start(opts...) // grpc.WithInsecure()
	if err != nil {
		return nil, fmt.Errorf("can't start grpc client: %v", err)
	}
	return c, nil
}

// SetTimeout sets timeout for node connection
func (c *TronClient) SetTimeout(timeout time.Duration, opts ...grpc.DialOption) error {
	if c == nil {
		return errors.New("client is nil")
	}
	c.grpc = client.NewGrpcClientWithTimeout(c.url, timeout)
	err := c.grpc.Start(opts...)
	if err != nil {
		return fmt.Errorf("can't start grpc client: %v", err)
	}
	return nil
}

// keepConnect keeps the connection and if it fails - reconnect
func (c *TronClient) keepConnect() error {
	_, err := c.grpc.GetNodeInfo()
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			return c.grpc.Reconnect(c.url)
		}
		return fmt.Errorf("can't connect to the node: %v", err)
	}
	return nil
}