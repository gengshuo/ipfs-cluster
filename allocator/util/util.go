package util

import (
	"strconv"

	"github.com/ipfs/ipfs-cluster/api"

	rpc "github.com/hsanjuan/go-libp2p-gorpc"
	peer "github.com/libp2p/go-libp2p-peer"
)

type Allocator interface {
	SetClient(c *rpc.Client)
	Shutdown() error
}

type SimpleAllocator struct{}

/*// NewAllocator returns an initialized Allocator
func NewAllocator() *Allocator {
	return &Allocator{}
}*/

// SetClient does nothing in this allocator
func (alloc SimpleAllocator) SetClient(c *rpc.Client) {}

// Shutdown does nothing in this allocator
func (alloc SimpleAllocator) Shutdown() error { return nil }

// MetricsSorter attaches sort.Interface methods to our metrics and sorts
// a slice of peers in the way that interest us
type MetricsSorter struct {
	Peers []peer.ID
	M     map[peer.ID]int
}

func NewMetricsSorter(m map[peer.ID]api.Metric) *MetricsSorter {
	vMap := make(map[peer.ID]int)
	peers := make([]peer.ID, 0, len(m))
	for k, v := range m {
		if v.Discard() {
			continue
		}
		val, err := strconv.Atoi(v.Value)
		if err != nil {
			continue
		}
		peers = append(peers, k)
		vMap[k] = val
	}

	sorter := &MetricsSorter{
		M:     vMap,
		Peers: peers,
	}
	return sorter
}

//type descendSorter struct {
//	MetricsSorter
//}

// Len returns the number of metrics
func (s MetricsSorter) Len() int {
	return len(s.Peers)
}

// Swap swaps the elements in positions i and j
func (s MetricsSorter) Swap(i, j int) {
	temp := s.Peers[i]
	s.Peers[i] = s.Peers[j]
	s.Peers[j] = temp
}
