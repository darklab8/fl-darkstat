package darkrpc

import (
	"log"
	"net/rpc"

	"github.com/darklab8/fl-darkstat/darkstat/configs_export"
	"github.com/darklab8/fl-darkstat/darkstat/router"
)

type ServerRpc struct {
	app_data *router.AppData
}

func NewRpc(app_data *router.AppData) RpcI {
	return &ServerRpc{app_data: app_data}
}

type Args struct {
}
type Reply struct {
	Bases []*configs_export.Base
}

// / CLIENT///////////////////
type ClientRpc struct {
	client *rpc.Client
	sock   string
}

const DarkstatSock = "/tmp/darkstat/rpc.sock"

type ClientOpt func(r *ClientRpc)

func WithSockCli(sock string) ClientOpt {
	return func(r *ClientRpc) {
		r.sock = sock
	}
}

func NewClient(opts ...ClientOpt) RpcI {
	cli := &ClientRpc{}

	for _, opt := range opts {
		opt(cli)
	}

	// client, err := rpc.DialHTTP("tcp", "127.0.0.1+":1234") // if serving over http
	client, err := rpc.Dial("unix", cli.sock) // if connecting over cli over sock
	if err != nil {
		log.Fatal("dialing:", err)
	}

	cli.client = client

	return RpcI(cli)
}

//// Methods

type RpcI interface {
	GetBases(args Args, reply *Reply) error
	GetHealth(args Args, reply *bool) error
}

func (t *ServerRpc) GetBases(args Args, reply *Reply) error {
	reply.Bases = t.app_data.Configs.Bases
	return nil
}

func (r *ClientRpc) GetBases(args Args, reply *Reply) error {
	// Synchronous call
	// return r.client.Call("ServerRpc.GetBases", args, &reply)

	// // Asynchronous call
	divCall := r.client.Go("ServerRpc.GetBases", args, &reply, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	return replyCall.Error
}

func (t *ServerRpc) GetHealth(args Args, reply *bool) error {
	*reply = true
	return nil
}

func (r *ClientRpc) GetHealth(args Args, reply *bool) error {
	divCall := r.client.Go("ServerRpc.GetHealth", args, &reply, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	return replyCall.Error
}
