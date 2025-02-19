package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/darklab8/fl-darkstat/darkrpc"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
)

func main() {
	app_data := appdata.NewAppData()
	srv := darkrpc.NewRpcServer()
	srv.Serve(app_data)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	srv.Close()
	fmt.Println("did graceful shutdown")
}
