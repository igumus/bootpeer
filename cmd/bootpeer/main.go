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

func main() {
	log.Println("info: starting bootpeer")

	flagRestHost := flag.String("rest-host", "0.0.0.0", "Rest Service Host")
	flagRestPort := flag.Int("rest-port", 2001, "Rest Service Port")
	flagPeerHost := flag.String("peer-host", "0.0.0.0", "Peer Host")
	flagPeerPort := flag.Int("peer-port", 3001, "Peer Port")
	flag.Parse()

	rootCtx := context.Background()
	peer, err := ipeer.New(rootCtx, ipeer.WithHost(*flagPeerHost), ipeer.WithPort(*flagPeerPort))
	if err != nil {
		log.Fatalf("err: creating peer faile: %s\n", err.Error())
	}

	httpAddr := fmt.Sprintf("%s:%d", *flagRestHost, *flagRestPort)
	http.ListenAndServe(httpAddr, peer.NewMux())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	log.Println("info: graceful shutdown bootpeer")
	peer.Close()
}
