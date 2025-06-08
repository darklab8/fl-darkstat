package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/darklab8/fl-darkstat/darkapis/darkrpc"
	"github.com/darklab8/fl-darkstat/darkstat/appdata"
)

func main() {
	ctx := context.Background()
	app_data := appdata.NewAppData(ctx)
	srv := darkrpc.NewRpcServer()
	srv.Serve(app_data)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	srv.Close()
	fmt.Println("did graceful shutdown")
}
