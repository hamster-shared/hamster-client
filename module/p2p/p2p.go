package p2p

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	ipfsp2p "github.com/ipfs/go-ipfs/p2p"
	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	pstore "github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/pnet"
	"github.com/libp2p/go-libp2p-core/protocol"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	ma "github.com/multiformats/go-multiaddr"
	madns "github.com/multiformats/go-multiaddr-dns"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hamster-client/config"

	"log"
	"strings"
	"time"
)

var resolveTimeout = 10 * time.Second

// MakeRoutedHost create a p2p routing client
func MakeRoutedHost(listenPort int, privstr string, peers []string) (host.Host, *dht.IpfsDHT, error) {

	skbytes, err := base64.StdEncoding.DecodeString(privstr)
	if err != nil {
		logrus.Error(err)
		return nil, nil, err
	}
	priv, err := crypto.UnmarshalPrivateKey(skbytes)
	if err != nil {
		logrus.Error(err)
		return nil, nil, err
	}
	bootstrapPeers := convertPeers(peers)

	// load private key swarm.key
	swarmkey := []byte(config.SwarmKey)

	psk, err := pnet.DecodeV1PSK(bytes.NewReader(swarmkey))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to configure private network: %s", err)
	}

	// Generate a key pair for this host. We will use it at least
	// to obtain a valid host ID.
	opts := []libp2p.Option{
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
		libp2p.PrivateNetwork(psk),
	}

	ctx := context.Background()

	basicHost, err := libp2p.New(ctx, opts...)
	if err != nil {
		return nil, nil, err
	}

	// Construct a datastore (needed by the DHT). This is just a simple, in-memory thread-safe datastore.
	dstore := dsync.MutexWrap(ds.NewMapDatastore())

	// Make the DHT
	newDht := dht.NewDHT(ctx, basicHost, dstore)

	// Make the routed host
	routedHost := rhost.Wrap(basicHost, newDht)

	// connect to the chosen ipfs nodes
	cfg := DefaultBootstrapConfig
	cfg.BootstrapPeers = func() []peer.AddrInfo {
		return bootstrapPeers
	}

	id, err := peer.IDFromPrivateKey(priv)
	_, err = Bootstrap(id, routedHost, newDht, cfg)

	if err != nil {
		return nil, nil, err
	}

	// Bootstrap the host
	err = newDht.Bootstrap(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", routedHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	// addr := routedHost.Addrs()[0]
	addrs := routedHost.Addrs()
	log.Println("I can be reached at:")
	for _, addr := range addrs {
		log.Println(addr.Encapsulate(hostAddr))
	}

	return routedHost, newDht, nil
}

// MakeIpfsP2p create ipfs p2p object
func MakeIpfsP2p(h *host.Host) *ipfsp2p.P2P {
	return ipfsp2p.New((*h).ID(), *h, (*h).Peerstore())
}

func CreateIdentity() (Identity, error) {
	ident := Identity{}

	priv, pub, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		return ident, err
	}

	// currently storing key unencrypted. in the future we need to encrypt it.
	// TODO(security)
	skbytes, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return ident, err
	}
	ident.PrivKey = base64.StdEncoding.EncodeToString(skbytes)

	id, err := peer.IDFromPublicKey(pub)
	if err != nil {
		return ident, err
	}
	ident.PeerID = id.Pretty()
	return ident, nil
}

// Identity p2p identity tag structure
type Identity struct {
	PeerID  string
	PrivKey string `json:",omitempty"`
}

// P2pClient p2p client
type P2pClient struct {
	Host  host.Host
	P2P   *ipfsp2p.P2P
	DHT   *dht.IpfsDHT
	db    *gorm.DB
	Peers []string
}

// P2PListenerInfoOutput  p2p monitoring or mapping information
type P2PListenerInfoOutput struct {
	Protocol      string
	ListenAddress string
	TargetAddress string
}

// P2PLsOutput p2p monitor or map information output
type P2PLsOutput struct {
	Listeners []P2PListenerInfoOutput
}

// List p2p monitor message list
func (c *P2pClient) List() *P2PLsOutput {
	output := &P2PLsOutput{}

	c.P2P.ListenersLocal.Lock()
	for _, listener := range c.P2P.ListenersLocal.Listeners {
		output.Listeners = append(output.Listeners, P2PListenerInfoOutput{
			Protocol:      string(listener.Protocol()),
			ListenAddress: listener.ListenAddress().String(),
			TargetAddress: listener.TargetAddress().String(),
		})
	}
	c.P2P.ListenersLocal.Unlock()

	c.P2P.ListenersP2P.Lock()
	for _, listener := range c.P2P.ListenersP2P.Listeners {
		output.Listeners = append(output.Listeners, P2PListenerInfoOutput{
			Protocol:      string(listener.Protocol()),
			ListenAddress: listener.ListenAddress().String(),
			TargetAddress: listener.TargetAddress().String(),
		})
	}
	c.P2P.ListenersP2P.Unlock()

	fmt.Println("p2p link list: ", output)

	return output
}

// Listen map local ports to p2p networks
func (c *P2pClient) Listen(protoOpt string, port int) error {
	log.Println("listening for connections")

	targetOpt := fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)
	proto := protocol.ID(protoOpt)

	target, err := ma.NewMultiaddr(targetOpt)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = c.P2P.ForwardRemote(context.Background(), proto, target, false)
	return err
}

// Forward connect p2p network to remote nodes / map to local port
func (c *P2pClient) Forward(protoOpt string, port int, peerId string) error {

	if peerId == "" {
		return fmt.Errorf("peer id cannot be empty")
	}

	if err := c.CheckForwardHealth(protoOpt, peerId); err != nil {
		// recover
		defer func(){
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
			}
		}()
		fmt.Println("CheckForwardHealth:", peerId)
		fmt.Println("c.Peers:", c.Peers)
		bootstrapPeers := randomSubsetOfPeers(convertPeers(c.Peers), 1)
		if len(bootstrapPeers) == 0 {
			return errors.New("not enough bootstrap peers")
		}
		circuitPeerId := bootstrapPeers[0].ID.Pretty()
		err = c.ConnectCircuit(circuitPeerId, peerId)
		if err != nil {
			return err
		}
	}

	listenOpt := fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)
	targetOpt := fmt.Sprintf("/p2p/%s", peerId)
	listen, err := ma.NewMultiaddr(listenOpt)

	if err != nil {
		log.Println(err)
		return err
	}

	targetAddrInfo, err := parseIpfsAddr(targetOpt)
	protoId := protocol.ID(protoOpt)

	c.P2P.ListenersP2P.Lock()
	defer c.P2P.ListenersP2P.Unlock()

	target, err := ma.NewMultiaddr(targetOpt)

	listeners := c.filterListener(c.P2P.ListenersLocal, func(listener ipfsp2p.Listener) bool {
		return listener.Protocol() == protoId && listener.ListenAddress().String() == listen.String() && listener.TargetAddress().String() == target.String()
	})

	if len(listeners) > 0 {
		return nil
	}
	err = forwardLocal(context.Background(), c.P2P, c.Host.Peerstore(), protoId, listen, targetAddrInfo)
	if err != nil {
		log.Println(err)
		return err
	}
	fmt.Println("======================")
	fmt.Println("forward : protoOpt: ", protoOpt)
	fmt.Println("forward : port: ", port)
	fmt.Println("forward : peerId: ", peerId)
	fmt.Println("======================")
	fmt.Println("remote_node" + peerId + ",forward to" + listenOpt + "success")
	return err
}

func (c *P2pClient) filterListener(listeners *ipfsp2p.Listeners, matchFunc func(listener ipfsp2p.Listener) bool) []ipfsp2p.Listener {
	todo := make([]ipfsp2p.Listener, 0)
	for _, l := range listeners.Listeners {
		if matchFunc(l) {
			todo = append(todo, l)
		}
	}
	return todo

}

func (c *P2pClient) ConnectCircuit(circuitPeer, targetPeer string) error {
	maddr := ma.StringCast(fmt.Sprintf("/p2p/%s/p2p-circuit/p2p/%s", circuitPeer, targetPeer))
	pi, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return err
	}
	err = c.Host.Connect(context.Background(), *pi)
	if err != nil {
		return err
	}
	return nil
}

// CheckForwardHealth check if the remote node is connected
func (c *P2pClient) CheckForwardHealth(protoOpt, target string) error {
	targets, err := parseIpfsAddr(target)
	proto := protocol.ID(protoOpt)
	if err != nil {
		return err
	}
	cctx, cancel := context.WithTimeout(context.Background(), time.Second*3) //TODO: configurable?
	defer cancel()
	stream, err := c.Host.NewStream(cctx, targets.ID, proto)
	if err != nil {
		return err
	} else {
		_ = stream.Close()
		return nil
	}
}

// Close closeP2P monitor or map
func (c *P2pClient) Close(target string) (int, error) {

	targetAddress, err := ma.NewMultiaddr(target)
	if err != nil {
		return 0, err
	}
	match := func(listener ipfsp2p.Listener) bool {
		if !targetAddress.Equal(listener.TargetAddress()) {
			return false
		}
		return true
	}

	done := c.P2P.ListenersLocal.Close(match)
	done += c.P2P.ListenersP2P.Close(match)

	fmt.Println("close target: ", done)

	return done, nil

}

// Destroy destroy and close the p2p client, including all subordinate listeners, stream objects
func (c *P2pClient) Destroy() error {
	for _, stream := range c.P2P.Streams.Streams {
		c.P2P.Streams.Close(stream)
	}
	match := func(listener ipfsp2p.Listener) bool {
		return true
	}
	c.P2P.ListenersP2P.Close(match)
	c.P2P.ListenersLocal.Close(match)
	err := c.Host.Close()
	c.P2P = nil
	c.Host = nil
	return err
}

// forwardLocal forwards local connections to a libp2p service
func forwardLocal(ctx context.Context, p *ipfsp2p.P2P, ps pstore.Peerstore, proto protocol.ID, bindAddr ma.Multiaddr, addr *peer.AddrInfo) error {
	ps.AddAddrs(addr.ID, addr.Addrs, pstore.TempAddrTTL)
	// TODO: return some info
	_, err := p.ForwardLocal(ctx, addr.ID, proto, bindAddr)
	return err
}

// parseIpfsAddr is a function that takes in addr string and return ipfsAddrs
func parseIpfsAddr(addr string) (*peer.AddrInfo, error) {

	if !strings.HasPrefix(addr, "/p2p/") {
		addr = "/p2p/" + addr
	}

	multiaddr, err := ma.NewMultiaddr(addr)
	if err != nil {
		return nil, err
	}

	pi, err := peer.AddrInfoFromP2pAddr(multiaddr)
	if err == nil {
		return pi, nil
	}

	// resolve multiaddr whose protocol is not ma.P_IPFS
	ctx, cancel := context.WithTimeout(context.Background(), resolveTimeout)
	defer cancel()
	addrs, err := madns.Resolve(ctx, multiaddr)
	if err != nil {
		return nil, err
	}
	if len(addrs) == 0 {
		return nil, errors.New("fail to resolve the multiaddr:" + multiaddr.String())
	}
	var info peer.AddrInfo
	for _, addr := range addrs {
		taddr, id := peer.SplitAddr(addr)
		if id == "" {
			// not an ipfs addr, skipping.
			continue
		}
		switch info.ID {
		case "":
			info.ID = id
		case id:
		default:
			return nil, fmt.Errorf(
				"ambiguous multiaddr %s could refer to %s or %s",
				multiaddr,
				info.ID,
				id,
			)
		}
		info.Addrs = append(info.Addrs, taddr)
	}
	return &info, nil
}
