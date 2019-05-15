/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package rpc

import (
	"github.com/seeleteam/scan-api/log"
	"time"
)

// SeeleRPC json_rpc client
type SeeleRPC struct {
	url    string
	scheme string
	conn   *Client
}

// NewRPC create new json_rpc client with given url
func NewRPC(url string, options ...func(rpc *SeeleRPC)) *SeeleRPC {
	rpc := &SeeleRPC{
		url:    url,
		scheme: "tcp",
	}
	for _, option := range options {
		option(rpc)
	}
	return rpc
}

//Connect Create tcp connect
func (rpc *SeeleRPC) Connect() error {
	if rpc.conn == nil {
		conn, err := Dial(rpc.scheme, rpc.url)
		if err != nil {
			return err
		}
		rpc.conn = conn
	}
	return nil
}

//Release release current rpc
func (rpc *SeeleRPC) Release() {
	if rpc != nil && rpc.conn != nil {
		rpc.conn.Close()
		rpc.conn = nil
	}
}

func (rpc *SeeleRPC) call(serviceMethod string, args interface{}, reply interface{}) error {
	if rpc.conn == nil {
		log.Error("rpc_client conn is nil, try to reconnect")
		for rpc.Connect() !=nil {
			time.Sleep(10*time.Second)
			log.Error("rpc_client conn is nil, try to reconnect")
			continue;
		}
	}
	err := rpc.conn.Call(serviceMethod, args, &reply)
	if err != nil {
		return err
	}

	return nil
}
