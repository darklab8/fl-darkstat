package darkrpc

import (
	"fmt"
	"net/rpc"

	"github.com/darklab8/fl-darkstat/darkstat/appdata"
	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/settings/logus"
)

type ServerRpc struct {
	app_data *appdata.AppData
}

func NewRpc(app_data *appdata.AppData) RpcI {
	return &ServerRpc{app_data: app_data}
}

type Args struct {
}
type Reply struct {
	Bases []*configs_export.Base
}

// / CLIENT///////////////////
type ClientRpc struct {
	sock *string
	port *int
}

const DarkstatRpcSock = "/tmp/darkstat/rpc.sock"

type ClientOpt func(r *ClientRpc)

func WithSockCli(sock string) ClientOpt {
	return func(r *ClientRpc) {
		r.sock = &sock
	}
}

func WithPortCli(port int) ClientOpt {
	return func(r *ClientRpc) {
		r.port = &port
	}
}

func NewClient(opts ...ClientOpt) RpcI {
	cli := &ClientRpc{}

	for _, opt := range opts {
		opt(cli)
	}

	return RpcI(cli)
}

func (r *ClientRpc) getClient() (*rpc.Client, error) {
	var client *rpc.Client
	var err error
	if r.sock == nil && r.port == nil {
		logus.Log.Panic("undefined sock and port for rpc client")
	}
	if r.sock != nil && r.port != nil {
		logus.Log.Panic("both sock and port are defined for rpc client")
	}
	if r.sock != nil {
		fmt.Println("initialized unix client")
		client, err = rpc.Dial("unix", *r.sock) // if connecting over cli over sock
	}
	if r.port != nil {
		info := fmt.Sprintf("127.0.0.1:%d", *r.port)
		fmt.Println("initializing tcp client at ", info)
		client, err = rpc.DialHTTP("tcp", info) // if serving over http
		fmt.Println("initialized tcp client")
	}

	if logus.Log.CheckWarn(err, "dialing:") {
		return nil, err
	}

	return client, err
}

//// Methods

type RpcI interface {
	GetBases(args Args, reply *Reply) error
	GetHealth(args Args, reply *bool) error
	GetInfo(args GetInfoArgs, reply *GetInfoReply) error
}

func (t *ServerRpc) GetBases(args Args, reply *Reply) error {
	reply.Bases = t.app_data.Configs.Bases
	return nil
}

func (r *ClientRpc) GetBases(args Args, reply *Reply) error {
	// Synchronous call
	// return r.client.Call("ServerRpc.GetBases", args, &reply)

	// // Asynchronous call
	client, err := r.getClient()
	if logus.Log.CheckWarn(err, "dialing:") {
		return err
	}

	divCall := client.Go("ServerRpc.GetBases", args, &reply, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	return replyCall.Error
}

func (t *ServerRpc) GetHealth(args Args, reply *bool) error {
	fmt.Println("received get health request")
	*reply = true
	logus.Log.Info("rpc server got health checked")
	return nil
}

func (r *ClientRpc) GetHealth(args Args, reply *bool) error {
	client, err := r.getClient()
	if logus.Log.CheckWarn(err, "dialing:") {
		return err
	}

	fmt.Println("querying server rpc for get health")
	divCall := client.Go("ServerRpc.GetHealth", args, &reply, nil)
	fmt.Println("stuck awaiting server")
	replyCall := <-divCall.Done // will be equal to divCall
	return replyCall.Error
}
