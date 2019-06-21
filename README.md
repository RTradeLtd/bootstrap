# bootstrap

[![codecov](https://codecov.io/gh/RTradeLtd/bootstrap/branch/master/graph/badge.svg)](https://codecov.io/gh/RTradeLtd/bootstrap)

`bootstrap` provides helpers for bootstrapping libp2p hosts. It supports bootstrapping off the default libp2p bootstrap peers from `ipfs/go-ipfs-config` combined with the Temporal production IPFS nodes. Additionally it supports a `DynamicBootstrap` method to be used in-combination with a persistent peerstore to enable a "decentralized boostrapping" method that isn't reliant on pre-existing hosts.

The `Bootstrap` and `DefaultBootstrapPeers` are modified versions of those contained in `hsanjuan/ipfs-lite`.