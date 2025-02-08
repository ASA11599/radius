package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"github.com/ASA11599/radius-server/internal/server"
	"github.com/ASA11599/radius-server/internal/store"
)

func signalChannel() <-chan os.Signal {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	return signals
}

func serverChannel(a server.Server) <-chan error {
	res := make(chan error)
	go func(c chan<- error) {
		c <- a.Start()
	}(res)
	return res
}

func getHost() string {
	if host, ok := os.LookupEnv("HOST"); ok {
		return host
	}
	return "0.0.0.0"
}

func getPort() int {
	if port, ok := os.LookupEnv("PORT"); ok {
		p, err := strconv.ParseInt(port, 10, 16)
		if err == nil {
			return int(p)
		}
	}
	return 80
}

func main() {
	sigs := signalChannel()
	var rs server.Server = server.NewRadiusServer(getHost(), getPort(), store.NewMemoryStore())
	defer func() {
		if err := rs.Stop(); err != nil {
			fmt.Println("Server stopped with error:", err)
		}
	}()
	select {
	case err := <-serverChannel(rs):
		fmt.Println("Server finished with error:", err)
	case s := <-sigs:
		fmt.Println("Received signal:", s)
	}
}
