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
}

const DarkstatSock = "/tmp/darkstat/rpc.sock"

func NewClient() RpcI {
	client, err := rpc.Dial("unix", DarkstatSock) // if serving over unix
	if err != nil {
		log.Fatal("dialing:", err)
	}
	return &ClientRpc{
		client: client,
	}
}

//// Methods

type RpcI interface {
	GetBases(args Args, reply *Reply) error
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
