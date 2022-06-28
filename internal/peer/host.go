package peer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	routedhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	ma "github.com/multiformats/go-multiaddr"
)

func printListenAddr(ctx context.Context, host host.Host) {
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", host.ID().Pretty()))
	for _, addr := range host.Addrs() {
		log.Printf("info: peer listening on : %s\n", addr.Encapsulate(hostAddr).String())
	}
}

func generateKeyPair(ctx context.Context) (crypto.PrivKey, error) {
	sk, _, err := crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return nil, err
	}
	return sk, nil
}

func makeHost(ctx context.Context, listenAddr string) (host.Host, error) {
	sk, err := generateKeyPair(ctx)
	if err != nil {
		log.Printf("err: generation key pair failed: %s\n", err.Error())
		return nil, err
	}
	connmgr, err := connmgr.NewConnManager(
		100,
		400,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		return nil, err
	}

	host, err := libp2p.New(
		libp2p.Identity(sk),
		libp2p.ListenAddrStrings(listenAddr),
		libp2p.ConnectionManager(connmgr),
		libp2p.DefaultTransports,
	)
	if err != nil {
		return nil, err
	}
	return host, nil
}

func makePeer(ctx context.Context, listenAddr string) (host.Host, *dht.IpfsDHT, error) {
	host, err := makeHost(ctx, listenAddr)
	if err != nil {
		return nil, nil, err
	}

	idht, err := dht.New(ctx, host, dht.Mode(dht.ModeServer))
	if err != nil {
		return host, nil, err
	}

	if err := idht.Bootstrap(ctx); err != nil {
		log.Printf("warn: dht bootstrapping failed: %s\n", err.Error())
	}
	rhost := routedhost.Wrap(host, idht)

	return rhost, idht, nil
}
