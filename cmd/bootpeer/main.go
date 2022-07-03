package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	ipeer "github.com/igumus/bootpeer/internal/peer"
)

const hostAddr = "0.0.0.0"

func main() {
	log.Println("info: starting bootpeer")

	flagRestPort := flag.Int("rest-port", 2001, "Rest Service Port")
	flagPeerPort := flag.Int("peer-port", 3001, "Peer Port")
	flag.Parse()

	rootCtx := context.Background()
	peer, err := ipeer.New(rootCtx, ipeer.WithPort(*flagPeerPort))
	if err != nil {
		log.Fatalf("err: creating peer faile: %s\n", err.Error())
	}

	httpAddr := fmt.Sprintf("%s:%d", hostAddr, *flagRestPort)
	http.ListenAndServe(httpAddr, peer.NewMux())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("info: graceful shutdown bootpeer")
	peer.Close()
}
