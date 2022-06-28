package peer

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/libp2p/go-libp2p-core/host"
	cpeer "github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

type BootstrapPeer interface {
	NewMux() *http.ServeMux
	io.Closer
}

type peer struct {
	host host.Host
	dht  *dht.IpfsDHT
}

func New(ctx context.Context, opts ...PeerOption) (BootstrapPeer, error) {
	cfg, cfgErr := applyOptions(opts...)
	if cfgErr != nil {
		return nil, cfgErr
	}
	h, idht, err := makePeer(ctx, cfg.toListenAddr())
	if err != nil {
		return nil, err
	}

	ret := &peer{
		host: h,
		dht:  idht,
	}

	printListenAddr(ctx, h)
	return ret, nil
}

func (p *peer) NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/addr", p.addrHandler)
	return mux
}

func (p *peer) addrHandler(w http.ResponseWriter, r *http.Request) {
	addrInfo := p.host.Peerstore().PeerInfo(p.host.ID())
	data, err := addrInfo.MarshalJSON()

	var addr cpeer.AddrInfo
	if err := json.Unmarshal(data, &addr); err != nil {
		log.Printf("err: unmarshalling addr info failed: %s\n", err.Error())
	} else {
		log.Println(addr.String())
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("can't write identity"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (p *peer) Close() error {
	if err := p.dht.Close(); err != nil {
		return err
	}
	if err := p.host.Close(); err != nil {
		return err
	}
	return nil
}
