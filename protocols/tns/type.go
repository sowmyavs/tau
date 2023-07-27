package tns

import (
	"github.com/taubyte/go-interfaces/kvdb"
	iface "github.com/taubyte/go-interfaces/services/tns"
	"github.com/taubyte/p2p/peer"
	streams "github.com/taubyte/p2p/streams/service"
	"github.com/taubyte/tau/config"
	"github.com/taubyte/tau/protocols/tns/engine"
)

var _ iface.Service = &Service{}

type Service struct {
	node      peer.Node
	db        kvdb.KVDB
	dbFactory kvdb.Factory
	stream    *streams.CommandService
	engine    *engine.Engine
}

func (s *Service) Node() peer.Node {
	return s.node
}

func (s *Service) KV() kvdb.KVDB {
	return s.db
}

type Config struct {
	config.Protocol `yaml:"z,omitempty"`
}
