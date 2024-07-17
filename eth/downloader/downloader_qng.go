// Copyright (c) 2017-2024 The qitmeer developers

package downloader

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

func (d *Downloader) SyncQng(peerid string, mode SyncMode, hash common.Hash) error {
	d.peers.lock.RLock()
	var peer *peerConnection
	for _, peer = range d.peers.peers {
		if peer.id == peerid {
			break
		}
	}
	d.peers.lock.RUnlock()

	if peer == nil {
		return errBadPeer
	}
	log.Info("Attempting to retrieve sync target", "peer", peer.id, "head", hash.String())
	headers, metas, err := d.fetchHeadersByHash(peer, hash, 1, 0, false)
	if err != nil || len(headers) != 1 {
		log.Warn("Failed to fetch sync target", "headers", len(headers), "err", err)
		return err
	}
	// Head header retrieved, if the hash matches, start the actual sync
	if metas[0] != hash {
		log.Error("Received invalid sync target", "want", hash, "have", metas[0])
		return err
	}
	return d.BeaconSync(mode, headers[0], headers[0])
}
