package p2p

import (
	"bytes"
	"cmp"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"net"
	"slices"
	"sync"
	"time"
)

type QngServer struct {
	cfg *Config

	peerFeed event.Feed

	peers map[enode.ID]*Peer

	lock *sync.RWMutex

	s *Server
}

func (srv *QngServer) Start(s *Server) error {
	srv.s = s
	srv.cfg = &s.Config
	return nil
}

func (srv *QngServer) Stop() {
	for srv.PeerCount() > 0 {
		srv.lock.Lock()
		for _, p := range srv.peers {
			p.Disconnect(DiscQuitting)
		}
		srv.lock.Unlock()
		select {
		case <-time.After(time.Second):
			srv.s.log.Info("Waiting for all peers closed in QngServer")
		}
	}
	srv.s.log.Info("QngServer stopped")
}

func (srv *QngServer) Connect(fd net.Conn, dialDest *enode.Node) (bool, error) {
	flags := trustedConn
	if dialDest == nil {
		flags |= inboundConn
	} else {
		flags |= dynDialedConn
	}
	return srv.SetupConn(fd, flags, dialDest)
}

func (srv *QngServer) SetupConn(fd net.Conn, flags connFlag, dialDest *enode.Node) (bool, error) {
	c := &conn{fd: fd, flags: flags, cont: make(chan error)}
	if dialDest == nil {
		c.transport = newRLPX(fd, nil)
	} else {
		c.transport = newRLPX(fd, dialDest.Pubkey())
	}

	ret, err := srv.setupConn(c, dialDest)
	if err != nil {
		if !c.is(inboundConn) {
			markDialError(err)
		} else {
			markServeError(err)
		}
		c.close(err)
	}
	return ret, err
}

func (srv *QngServer) setupConn(c *conn, dialDest *enode.Node) (bool, error) {
	// If dialing, figure out the remote public key.
	if dialDest != nil {
		dialPubkey := new(ecdsa.PublicKey)
		if err := dialDest.Load((*enode.Secp256k1)(dialPubkey)); err != nil {
			err = fmt.Errorf("%w: dial destination doesn't have a secp256k1 public key", errEncHandshakeError)
			srv.s.log.Trace("Setting up connection failed", "addr", c.fd.RemoteAddr(), "conn", c.flags, "err", err)
			return false, err
		}
	}

	// Run the RLPx handshake.
	remotePubkey, err := c.doEncHandshake(srv.cfg.PrivateKey)
	if err != nil {
		srv.s.log.Trace("Failed RLPx handshake", "addr", c.fd.RemoteAddr(), "conn", c.flags, "err", err)
		return false, fmt.Errorf("%w: %v", errEncHandshakeError, err)
	}
	if dialDest != nil {
		c.node = dialDest
	} else {
		c.node = nodeFromConn(remotePubkey, c.fd)
	}
	clog := srv.s.log.New("id", c.node.ID(), "addr", c.fd.RemoteAddr(), "conn", c.flags)

	// Run the capability negotiation handshake.
	phs, err := c.doProtoHandshake(srv.s.ourHandshake)
	if err != nil {
		clog.Trace("Failed p2p handshake", "err", err)
		return false, err
	}
	if id := c.node.ID(); !bytes.Equal(crypto.Keccak256(phs.ID), id[:]) {
		clog.Trace("Wrong devp2p handshake identity", "phsid", hex.EncodeToString(phs.ID))
		return false, DiscUnexpectedIdentity
	}
	c.caps, c.name = phs.Caps, phs.Name
	err = srv.addPeerChecks(c)
	if err != nil {
		clog.Trace("Rejected peer", "err", err)
		return false, err
	}
	_, ret := srv.launchPeer(c)
	return ret, nil
}

func (srv *QngServer) addPeerChecks(c *conn) error {
	// Drop connections with no matching protocols.
	if len(srv.cfg.Protocols) > 0 && countMatchingProtocols(srv.cfg.Protocols, c.caps) == 0 {
		return DiscUselessPeer
	}
	return nil
}

func (srv *QngServer) launchPeer(c *conn) (*Peer, bool) {
	p := newPeer(srv.s.log, c, srv.cfg.Protocols)
	if srv.cfg.EnableMsgEvents {
		// If message events are enabled, pass the peerFeed
		// to the peer.
		p.events = &srv.peerFeed
	}
	ret := srv.runPeer(p)
	return p, ret
}

func (srv *QngServer) runPeer(p *Peer) bool {
	start := time.Now()
	srv.addPeer(p)

	srv.peerFeed.Send(&PeerEvent{
		Type:          PeerEventTypeAdd,
		Peer:          p.ID(),
		RemoteAddress: p.RemoteAddr().String(),
		LocalAddress:  p.LocalAddr().String(),
	})

	// Run the per-peer main loop.
	remoteRequested, err := p.run()

	srv.peerFeed.Send(&PeerEvent{
		Type:          PeerEventTypeDrop,
		Peer:          p.ID(),
		Error:         err.Error(),
		RemoteAddress: p.RemoteAddr().String(),
		LocalAddress:  p.LocalAddr().String(),
	})

	srv.delPeer(p, time.Since(start), remoteRequested, err)
	return remoteRequested
}

func (srv *QngServer) addPeer(p *Peer) {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	srv.peers[p.ID()] = p
	srv.s.log.Debug("Adding QNG p2p peer", "peercount", len(srv.peers), "id", p.ID(), "conn", p.rw.flags, "addr", p.RemoteAddr(), "name", p.Name())
}

func (srv *QngServer) delPeer(p *Peer, dur time.Duration, requested bool, err error) {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	delete(srv.peers, p.ID())
	srv.s.log.Debug("Removing QNG p2p peer", "peercount", len(srv.peers), "id", p.ID(), "duration", dur, "req", requested, "err", err)
}

func (srv *QngServer) PeerCount() int {
	srv.lock.RLock()
	defer srv.lock.RUnlock()

	return len(srv.peers)
}

func (srv *QngServer) Peers() []*Peer {
	srv.lock.RLock()
	defer srv.lock.RUnlock()

	var ps []*Peer
	for _, p := range srv.peers {
		ps = append(ps, p)
	}
	return ps
}

func (srv *QngServer) GetPeer(id enode.ID) *Peer {
	srv.lock.RLock()
	defer srv.lock.RUnlock()

	p, ok := srv.peers[id]
	if !ok {
		return nil
	}
	return p
}

func (srv *QngServer) NodeInfo() *NodeInfo {
	node := srv.s.Self()
	info := &NodeInfo{
		Name:       srv.s.Name,
		Enode:      node.URLv4(),
		ID:         node.ID().String(),
		IP:         node.IPAddr().String(),
		ListenAddr: srv.s.ListenAddr,
		Protocols:  make(map[string]interface{}),
	}
	info.Ports.Discovery = node.UDP()
	info.Ports.Listener = node.TCP()
	info.ENR = node.String()

	for _, proto := range srv.s.Protocols {
		if _, ok := info.Protocols[proto.Name]; !ok {
			nodeInfo := interface{}("unknown")
			if query := proto.NodeInfo; query != nil {
				nodeInfo = proto.NodeInfo()
			}
			info.Protocols[proto.Name] = nodeInfo
		}
	}
	return info
}

func (srv *QngServer) PeersInfo() []*PeerInfo {
	infos := make([]*PeerInfo, 0, srv.PeerCount())
	for _, peer := range srv.Peers() {
		if peer != nil {
			infos = append(infos, peer.Info())
		}
	}
	slices.SortFunc(infos, func(a, b *PeerInfo) int {
		return cmp.Compare(a.ID, b.ID)
	})

	return infos
}

func NewQngServer() *QngServer {
	qs := &QngServer{
		peers: make(map[enode.ID]*Peer),
		lock:  &sync.RWMutex{},
	}
	return qs
}
