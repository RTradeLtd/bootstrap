package bootstrap

import (
	"context"
	"testing"

	testutils "github.com/RTradeLtd/go-libp2p-testutils"
	"github.com/multiformats/go-multiaddr"
)

func Test_DefaultBootstrapPeers(t *testing.T) {
	if _, err := DefaultBootstrapPeers(); err != nil {
		t.Fatal(err)
	}
}

func Test_DynamicBootstrap(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ds := testutils.NewDatastore(t)
	ps := testutils.NewPeerstore(t)
	logger := testutils.NewLogger(t)
	pk := testutils.NewPrivateKey(t)
	addrs := []multiaddr.Multiaddr{testutils.NewMultiaddr(t)}
	host, dht := testutils.NewLibp2pHostAndDHT(
		ctx,
		t,
		logger.Desugar(),
		ds,
		ps,
		pk,
		addrs,
		nil,
	)
	closeFunc := func() {
		dht.Close()
		host.Close()
	}
	defer closeFunc()
	// test against no previous peers
	if err := DynamicBootstrap(ctx, logger.Desugar(), dht, host); err != nil {
		t.Fatal(err)
	}
	// get bootstrap peers
	peers, err := DefaultBootstrapPeers()
	if err != nil {
		t.Fatal(err)
	}
	// do a regular bootstrap
	if err := Bootstrap(ctx, logger.Desugar(), dht, host, peers); err != nil {
		t.Fatal(err)
	}
	// test against previous peers
	if err := DynamicBootstrap(ctx, logger.Desugar(), dht, host); err != nil {
		t.Fatal(err)
	}
	cancel()
	closeFunc()
}
func Test_Bootstrap(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ds := testutils.NewDatastore(t)
	ps := testutils.NewPeerstore(t)
	logger := testutils.NewLogger(t)
	pk := testutils.NewPrivateKey(t)
	addrs := []multiaddr.Multiaddr{testutils.NewMultiaddr(t)}
	host, dht := testutils.NewLibp2pHostAndDHT(
		ctx,
		t,
		logger.Desugar(),
		ds,
		ps,
		pk,
		addrs,
		nil,
	)
	closeFunc := func() {
		dht.Close()
		host.Close()
	}
	defer closeFunc()
	peers, err := DefaultBootstrapPeers()
	if err != nil {
		t.Fatal(err)
	}
	if err := Bootstrap(ctx, logger.Desugar(), dht, host, peers); err != nil {
		t.Fatal(err)
	}
	cancel()
	closeFunc()
}
