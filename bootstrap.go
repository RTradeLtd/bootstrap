package bootstrap

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	config "github.com/ipfs/go-ipfs-config"
	libcore "github.com/libp2p/go-libp2p-core"
	peer "github.com/libp2p/go-libp2p-core/peer"

	host "github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	"go.uber.org/zap"
)

// TemporalPeerAddresses are the multiaddrs for
// Temporal's production nodes.
var TemporalPeerAddresses = []string{
	"/ip4/172.218.49.115/tcp/4002/ipfs/QmPvnFXWAz1eSghXD6JKpHxaGjbVo4VhBXY2wdBxKPbne5",
	"/ip4/172.218.49.115/tcp/4003/ipfs/QmXow5Vu8YXqvabkptQ7HddvNPpbLhXzmmU53yPCM54EQa",
	"/ip4/35.203.44.77/tcp/4001/ipfs/QmUMtzoRfQ6FttA7RygL8jJf7TZJBbdbZqKTmHfU6QC5Jm",
}

// DefaultBootstrapPeers returns the default lsit
// of bootstrap peers used by go-ipfs, updated
// with the Temporal bootstrap nodes
func DefaultBootstrapPeers() ([]libcore.PeerAddrInfo, error) {
	// conversion copied from go-ipfs
	defaults, err := config.DefaultBootstrapPeers()
	if err != nil {
		return nil, err
	}
	tPeers, err := config.ParseBootstrapPeers(TemporalPeerAddresses)
	if err != nil {
		return nil, err
	}
	defaults = append(defaults, tPeers...)
	pinfos := make(map[peer.ID]*libcore.PeerAddrInfo)
	for _, bootstrap := range defaults {
		pinfo, ok := pinfos[bootstrap.ID]
		if !ok {
			pinfo = new(libcore.PeerAddrInfo)
			pinfos[bootstrap.ID] = pinfo
			pinfo.ID = bootstrap.ID
		}

		pinfo.Addrs = append(pinfo.Addrs, bootstrap.Addrs...)
	}
	var peers []libcore.PeerAddrInfo
	for _, pinfo := range pinfos {
		peers = append(peers, *pinfo)
	}
	return peers, nil
}

// DynamicBootstrap is used to bootstrapa a host off a dynamic list of bootstrap peers.
// This must be used in combination with a host using a datastore backed peerstore
// so that we have a persistent set of peers to boot from, otherwise we just do a default bootstrap
// The final list of bootstrap peers is at most 10 randomly selected from the peerstore, combined
// with the deafult libp2p bootstrap peers
func DynamicBootstrap(ctx context.Context, logger *zap.Logger, dt routing.Routing, hst host.Host) error {
	defaultPeers, err := DefaultBootstrapPeers()
	if err != nil {
		return err
	}
	if len(hst.Peerstore().Peers()) == 0 {
		return Bootstrap(ctx, logger, dt, hst, defaultPeers)
	}
	peers := hst.Peerstore().Peers()
	var (
		peerAddrs = defaultPeers
		found     = make(map[peer.ID]bool)
		limit     = len(peers)
	)
	if limit > 10 {
		limit = 10
	}
	for i := 0; i < limit; i++ {
		pid := peers[rand.Intn(len(peers))]
		if found[pid] {
			continue
		}
		found[pid] = true
		peerAddrs = append(peerAddrs, hst.Peerstore().PeerInfo(pid))
	}
	return Bootstrap(ctx, logger, dt, hst, peerAddrs)
}

// Bootstrap is used to connect our libp2p host to the specified set of peers
func Bootstrap(ctx context.Context, logger *zap.Logger, dt routing.Routing, hst host.Host, peers []libcore.PeerAddrInfo) error {
	var (
		connected = make(chan bool, len(peers))
		wg        sync.WaitGroup
	)
	for _, pinfo := range peers {
		wg.Add(1)
		go func(pinfo libcore.PeerAddrInfo) {
			defer wg.Done()
			err := hst.Connect(ctx, pinfo)
			if err != nil {
				logger.Error("failed to connect to host", zap.Error(err), zap.String("peerid", pinfo.ID.String()))
				return
			}
			logger.Info("successfully connected to peer", zap.String("peerid", pinfo.ID.String()))
			connected <- true
		}(pinfo)
	}

	go func() {
		wg.Wait()
		close(connected)
	}()

	i := 0
	for range connected {
		i++
	}
	if nPeers := len(peers); i < nPeers/2 {
		logger.Warn(fmt.Sprintf("only connected to %d peers out of %d", i, nPeers))
	}

	return dt.Bootstrap(ctx)
}
