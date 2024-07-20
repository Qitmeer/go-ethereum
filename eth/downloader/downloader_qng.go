// Copyright (c) 2017-2024 The qitmeer developers

package downloader

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/protocols/eth"
	"github.com/ethereum/go-ethereum/log"
	"time"
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
	headers, metas, err := d.fetchQngHeadersByHash(peer, hash, 1, 0, false)
	if err != nil || len(headers) != 1 {
		log.Warn("Failed to fetch sync target", "headers", len(headers), "err", err)
		return err
	}
	// Head header retrieved, if the hash matches, start the actual sync
	if metas[0] != hash {
		log.Error("Received invalid sync target", "want", hash, "have", metas[0])
		return err
	}
	if d.skeleton.filler.(*beaconBackfiller).filling {
		return fmt.Errorf("backfiller is filling:%s", hash.String())
	}
	return d.BeaconSync(mode, headers[0], headers[0])
}

func (d *Downloader) fetchQngHeadersByHash(p *peerConnection, hash common.Hash, amount int, skip int, reverse bool) ([]*types.Header, []common.Hash, error) {
	// Create the response sink and send the network request
	start := time.Now()
	resCh := make(chan *eth.Response)

	req, err := p.peer.RequestHeadersByHash(hash, amount, skip, reverse, resCh)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		req.Close()
	}()

	// Wait until the response arrives, the request is cancelled or times out
	ttl := d.peers.rates.TargetTimeout()

	timeoutTimer := time.NewTimer(ttl)
	defer timeoutTimer.Stop()

	select {
	case <-timeoutTimer.C:
		// Header retrieval timed out, update the metrics
		p.log.Debug("Header request timed out", "elapsed", ttl)
		headerTimeoutMeter.Mark(1)

		return nil, nil, errTimeout

	case res := <-resCh:
		// Headers successfully retrieved, update the metrics
		headerReqTimer.Update(time.Since(start))
		headerInMeter.Mark(int64(len(*res.Res.(*eth.BlockHeadersRequest))))

		// Don't reject the packet even if it turns out to be bad, downloader will
		// disconnect the peer on its own terms. Simply delivery the headers to
		// be processed by the caller
		res.Done <- nil

		return *res.Res.(*eth.BlockHeadersRequest), res.Meta.([]common.Hash), nil
	}
}
