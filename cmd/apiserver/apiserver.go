package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"k8s.io/apiserver/pkg/server"
	"k8s.io/apiserver/pkg/util/logs"

	"github.com/lrx0014/log-tools/cmd/apiserver/app"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	command := app.NewAPIServerCommand(server.SetupSignalHandler())

	logs.InitLogs()
	defer logs.FlushLogs()

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
