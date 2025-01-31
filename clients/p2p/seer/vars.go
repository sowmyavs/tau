package seer

import (
	"errors"
	"time"

	"github.com/ipfs/go-log/v2"
)

var (
	MinPeers                 = 0
	MaxPeers                 = 2
	DefaultGeoBeaconInterval = 5 * time.Minute
	ErrorGeoBeaconStopped    = errors.New("geoBeacon Stopped")
	logger                   = log.Logger("seer.p2p.client")
)
