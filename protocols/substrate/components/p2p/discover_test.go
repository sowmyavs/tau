package p2p_test

import (
	"testing"
	"time"

	commonDreamland "github.com/taubyte/dreamland/core/common"
	dreamland "github.com/taubyte/dreamland/core/services"
	commonIface "github.com/taubyte/go-interfaces/common"
	_ "github.com/taubyte/tau/protocols/substrate"
	"github.com/taubyte/tau/protocols/substrate/components/p2p"
	"gotest.tools/assert"
)

func TestService_Discover(t *testing.T) {
	u := dreamland.Multiverse("TestService_Discover")

	defer u.Stop()

	err := u.StartWithConfig(&commonDreamland.Config{
		Services: map[string]commonIface.ServiceConfig{
			"substrate": {Others: map[string]int{"copies": 4}},
		},
		Simples: map[string]commonDreamland.SimpleConfig{},
	})
	assert.NilError(t, err)

	srv, err := p2p.New(u.Substrate())
	assert.NilError(t, err)

	peers, err := srv.Discover(u.Context(), 5, time.Second*2)
	assert.NilError(t, err)

	assert.Equal(t, len(peers), 3)
}
