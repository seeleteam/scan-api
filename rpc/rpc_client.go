/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package rpc

import (
	"github.com/seeleteam/scan-api/log"
	"runtime"
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
			log.Error(err)
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
		tryCnt := 0
		for rpc.Connect() !=nil {
			if tryCnt >=20 {
				runtime.Goexit()
			}
			tryCnt++
			time.Sleep(10*time.Second)
			log.Error("rpc_client conn is nil, try to reconnect")
			continue;
		}
	}
ErrCont:	err := rpc.conn.Call(serviceMethod, args, &reply)
	if err != nil {
		rpc.conn.Close()
		rpc.conn = nil
		rpc = NewRPC(rpc.url)
		tryCnt := 0
		for rpc.Connect() !=nil {
			if tryCnt >=20 {
				runtime.Goexit()
			}
			tryCnt++

			time.Sleep(10*time.Second)
			log.Error("rpc_client conn err, try to reconnect")
			continue;
		}
		goto ErrCont
		return err
	}

	return nil
}
