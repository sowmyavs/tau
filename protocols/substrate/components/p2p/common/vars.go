package common

import "github.com/ipfs/go-log/v2"

//TODO: Move to specs

const (
	ServiceName = "substrate"
	Protocol    = "/substrate/v1"
	MinPeers    = 0
	MaxPeers    = 4
)

var Logger = log.Logger("substrate.service.p2p")
