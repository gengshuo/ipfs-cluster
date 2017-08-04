// Package descendalloc implements an ipfscluster.util.Allocator returns allocations
// based on sorting the metrics in descending order. Thus, peers with smallest
// metrics are first in the list. This allocator can be used with a number
// of informers, as long as they provide a numeric metric value.
package descendalloc

import (
	"sort"

	"github.com/ipfs/ipfs-cluster/allocator/util"
	"github.com/ipfs/ipfs-cluster/api"

	cid "github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log"
	peer "github.com/libp2p/go-libp2p-peer"
)

var logger = logging.Logger("descendalloc")

type DescendAllocator struct {
	util.SimpleAllocator
}

// NewAllocator returns an initialized Allocator
func NewAllocator() DescendAllocator {
	return DescendAllocator{}
}

// Allocate returns where to allocate a pin request based on metrics which
// carry a numeric value such as "used disk". We do not pay attention to
// the metrics of the currently allocated peers and we just sort the candidates
// based on their metric values (from smallest to largest).
func (alloc DescendAllocator) Allocate(c *cid.Cid, current, candidates map[peer.ID]api.Metric) ([]peer.ID, error) {
	// sort our metrics
	sortable := descendSorter{*util.NewMetricsSorter(candidates)}
	sort.Sort(sortable)
	return sortable.Peers, nil
}

type descendSorter struct {
	util.MetricsSorter
}

// Less reports if the element in position i is less than the element in j
func (s descendSorter) Less(i, j int) bool {
	peeri := s.Peers[i]
	peerj := s.Peers[j]

	x := s.M[peeri]
	y := s.M[peerj]

	return x > y
}
