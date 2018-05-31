/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package rpc

// SeeleRPC json_rpc client
type SeeleRPC struct {
	url    string
	scheme string
	conn   *Client
}

var (
	gSeeleRPC *SeeleRPC
	//RPCURL address of the node which provided rpc service
	RPCURL = "127.0.0.1:55028"
)

//GetSeeleRPC get rpc connection
func GetSeeleRPC() (*SeeleRPC, error) {
	if gSeeleRPC == nil {
		gSeeleRPC = newSeeleRPC(RPCURL)
	}

	if gSeeleRPC.conn == nil {
		conn, err := Dial(gSeeleRPC.scheme, gSeeleRPC.url)
		if err != nil {
			return nil, err
		}
		gSeeleRPC.conn = conn
	}

	return gSeeleRPC, nil
}

//ReleaseSeeleRPC free rpc connection
func ReleaseSeeleRPC() {
	if gSeeleRPC != nil && gSeeleRPC.conn != nil {
		gSeeleRPC.conn.Close()
		gSeeleRPC.conn = nil
	}
}

// New create new json_rpc client with given url
func newRPC(url string, options ...func(rpc *SeeleRPC)) *SeeleRPC {
	rpc := &SeeleRPC{
		url:    url,
		scheme: "tcp",
	}
	for _, option := range options {
		option(rpc)
	}
	return rpc
}

func newSeeleRPC(url string, options ...func(rpc *SeeleRPC)) *SeeleRPC {
	return newRPC(url, options...)
}

func (rpc *SeeleRPC) call(serviceMethod string, args interface{}, reply interface{}) error {
	err := rpc.conn.Call(serviceMethod, args, &reply)
	if err != nil {
		return err
	}

	return nil
}
