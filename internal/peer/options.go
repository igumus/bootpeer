package peer

import (
	"errors"
	"fmt"
	"strings"
)

type PeerOption func(*peerConfigOption)

type peerConfigOption struct {
	host string
	port int
}

func (pco *peerConfigOption) toListenAddr() string {
	return fmt.Sprintf("/ip4/%s/tcp/%d", pco.host, pco.port)
}

func defaultPeerConfigOption() *peerConfigOption {
	return &peerConfigOption{
		host: "0.0.0.0",
		port: 3001,
	}
}

func validate(pco *peerConfigOption) error {
	if len(pco.host) == 0 {
		return errors.New("[bootstrap] peer configuration failed: host addr should be valid")
	}
	if pco.port < 0 || pco.port > 65535 {
		return errors.New("[bootstrap] peer configuration failed: port is invalid")
	}
	return nil
}

func applyOptions(opts ...PeerOption) (*peerConfigOption, error) {
	cfg := defaultPeerConfigOption()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg, validate(cfg)
}

func WithHost(h string) PeerOption {
	return func(pco *peerConfigOption) {
		pco.host = strings.TrimSpace(h)
	}
}

func WithPort(p int) PeerOption {
	return func(pco *peerConfigOption) {
		pco.port = p
	}
}
