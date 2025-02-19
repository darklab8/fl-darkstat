package darkrpc

import (
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
	sock string
}

const DarkstatRpcSock = "/tmp/darkstat/rpc.sock"

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

	return RpcI(cli)
}

func (r *ClientRpc) getClient() (*rpc.Client, error) {
	// client, err := rpc.DialHTTP("tcp", "127.0.0.1+":1234") // if serving over http
	client, err := rpc.Dial("unix", r.sock) // if connecting over cli over sock

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
	*reply = true
	logus.Log.Info("rpc server got health checked")
	return nil
}

func (r *ClientRpc) GetHealth(args Args, reply *bool) error {
	client, err := r.getClient()
	if logus.Log.CheckWarn(err, "dialing:") {
		return err
	}

	divCall := client.Go("ServerRpc.GetHealth", args, &reply, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	return replyCall.Error
}
