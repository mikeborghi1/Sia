package hostdb

import (
	"errors"
	"sync"

	"github.com/NebulousLabs/Sia/consensus"
	"github.com/NebulousLabs/Sia/modules"
	"github.com/NebulousLabs/Sia/network"
)

var (
	ErrNilState            = errors.New("hostdb can't use nil State")
	ErrMissingGenesisBlock = errors.New("state doesn't have a genesis block")
)

// The HostDB is a database of potential hosts. It assigns a weight to each
// host based on their hosting parameters.
type HostDB struct {
	state       *consensus.State
	recentBlock consensus.BlockID

	hostTree    *hostNode
	activeHosts map[string]*hostNode
	allHosts    map[network.Address]*modules.HostEntry

	mu sync.RWMutex
}

// New returns an empty HostDatabase.
func New(s *consensus.State) (hdb *HostDB, err error) {
	if s == nil {
		err = ErrNilState
		return
	}

	genesisBlock, exists := s.BlockAtHeight(0)
	if !exists {
		if consensus.DEBUG {
			panic(ErrMissingGenesisBlock.Error())
		}
		err = ErrMissingGenesisBlock
		return
	}

	hdb = &HostDB{
		state:       s,
		recentBlock: genesisBlock.ID(),
		activeHosts: make(map[string]*hostNode),
		allHosts:    make(map[network.Address]*modules.HostEntry),
	}
	return
}
