// Test KV services and clients.
//
// It's called a "system" test because it doesn't test a component (like
// KVService) in isolation; rather, the test harness constructs a complete
// system comprising of a cluster of services and some KVClients to exercise it.

package main

import (
	"testing"
	"time"
)

func sleepMs(n int) {
	time.Sleep(time.Duration(n) * time.Millisecond)
}

func TestElectionLeaderDisconnect(t *testing.T) {
	n := 5
	h := NewHarness(t, n)
	defer h.Shutdown()

	origLeaderId := h.CheckSingleLeader()

	h.DisconnectServiceFromPeers(origLeaderId)
	sleepMs(350)

	newLeaderId := h.CheckSingleLeader()
	if newLeaderId == origLeaderId {
		t.Errorf("want new leader to be different from orig leader")
	}
}

func TestDisconnect2Followers(t *testing.T) {
	n := 5
	h := NewHarness(t, n)
	defer h.Shutdown()

	origLeaderId := h.CheckSingleLeader()

	// send a PUT request to the cluster
	c1 := h.NewClient()
	prev, found := h.CheckPut(c1, "one", "RAFT")
	if found || prev != "" {
		t.Errorf(`got found=%v, prev=%v, want false/""`, found, prev)
	}

	// disconnect 2 followers
	numDisconn := 2
	disconnIds := []int{}
	for i := range n {
		if i != origLeaderId {
			h.DisconnectServiceFromPeers(i)
			disconnIds = append(disconnIds, i)
			if len(disconnIds) == numDisconn {
				break
			}
		}
	}

	prev, found = h.CheckPut(c1, "two", "pBFT")
	if found || prev != "" {
		t.Errorf(`got found=%v, prev=%v, want false/""`, found, prev)
	}

	for _, id := range disconnIds {
		h.ReconnectServiceToPeers(id)
	}

	h.CheckGet(c1, "two", "pBFT")

	sleepMs(3000)
}

func TestDisconnect3Followers(t *testing.T) {
	n := 5
	h := NewHarness(t, n)
	defer h.Shutdown()

	origLeaderId := h.CheckSingleLeader()

	// send a PUT request to the cluster
	c1 := h.NewClient()
	prev, found := h.CheckPut(c1, "one", "RAFT")
	if found || prev != "" {
		t.Errorf(`got found=%v, prev=%v, want false/""`, found, prev)
	}

	// disconnect 2 followers
	numDisconn := 3
	disconnIds := []int{}
	for i := range n {
		if i != origLeaderId {
			h.DisconnectServiceFromPeers(i)
			disconnIds = append(disconnIds, i)
			if len(disconnIds) == numDisconn {
				break
			}
		}
	}

	prev, found = h.CheckPut(c1, "two", "pBFT")
	if found || prev != "" {
		t.Errorf(`got found=%v, prev=%v, want false/""`, found, prev)
	}

	for _, id := range disconnIds {
		h.ReconnectServiceToPeers(id)
	}

	h.CheckGet(c1, "two", "pBFT")

	sleepMs(2000)
}

func Test2Partition(t *testing.T) {
	n := 5
	h := NewHarness(t, n)
	defer h.Shutdown()

	h.CheckSingleLeader()

	// send a PUT request to the cluster
	c1 := h.NewClient()
	prev, found := h.CheckPut(c1, "one", "RAFT")
	if found || prev != "" {
		t.Errorf(`got found=%v, prev=%v, want false/""`, found, prev)
	}

	cluster1, cluster2 := PartitionInto2ClusterBySize(n, 2)
	h.Disconnect2Cluster(cluster1, cluster2)

	prev, found = h.CheckPut(c1, "one", "Bitcoin")
	if !found || prev != "RAFT" {
		t.Errorf(`got found=%v, prev=%v, want true/"RAFT"`, found, prev)
	}
	sleepMs(300)

	h.Reconnect2Cluster(cluster1, cluster2)

	h.CheckGet(c1, "one", "Bitcoin")

	sleepMs(2000)
}
