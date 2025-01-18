package sphinx

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"sync"

	"github.com/btcsuite/btcd/btcec/v2"
)

// Constants defining sizes for various components of the onion packet and routing.
const (
	AddressSize         = 8   // Length of serialized address (channel_id).
	RealmByteSize       = 1   // Size of realm byte.
	AmtForwardSize      = 8   // Size of the amount to forward.
	OutgoingCLTVSize    = 4   // Size of the outgoing CLTV value.
	NumPaddingBytes     = 12  // Padding bytes reserved for future use.
	LegacyHopDataSize   = 65  // Fixed size of hop_data per BOLT 04.
	RoutingInfoSize     = 1300
	NumStreamBytes      = RoutingInfoSize * 2
	KeyLength           = 32  // Length of keys used in encryption.
	BaseVersion         = 0   // Current onion packet version.
	MaxPayloadSize      = RoutingInfoSize
	ErrInvalidOnionHMAC = "invalid HMAC for onion packet"
)

// Errors
var (
	ErrMaxRoutingInfoSizeExceeded = fmt.Errorf(
		"max routing info size of %v bytes exceeded", RoutingInfoSize,
	)
)

// OnionPacket encapsulates the onion routing information for a message.
type OnionPacket struct {
	Version      byte                  // Version of the onion packet.
	EphemeralKey *btcec.PublicKey      // Ephemeral key for ECDH.
	RoutingInfo  [RoutingInfoSize]byte // Routing information.
	HeaderMAC    [sha256.Size]byte     // HMAC for routing data.
}

// Router represents an onion router within the Sphinx network.
type Router struct {
	onionKey SingleKeyECDH
	log      ReplayLog
}

// NewRouter initializes a new Sphinx router with the given onion key and replay log.
func NewRouter(nodeKey SingleKeyECDH, log ReplayLog) *Router {
	return &Router{
		onionKey: nodeKey,
		log:      log,
	}
}

// Start initializes the router's replay log.
func (r *Router) Start() error {
	return r.log.Start()
}

// Stop halts the router's replay log.
func (r *Router) Stop() {
	r.log.Stop()
}

// ProcessOnionPacket processes an incoming onion packet and returns the result.
func (r *Router) ProcessOnionPacket(onionPkt *OnionPacket, assocData []byte,
	incomingCltv uint32, opts ...ProcessOnionOpt) (*ProcessedPacket, error) {

	cfg := &processOnionCfg{}
	for _, o := range opts {
		o(cfg)
	}

	// Compute shared secret for the onion packet.
	sharedSecret, err := r.generateSharedSecret(
		onionPkt.EphemeralKey, cfg.blindingPoint,
	)
	if err != nil {
		return nil, err
	}

	// Check replay log for duplicates.
	hashPrefix := hashSharedSecret(&sharedSecret)
	packet, err := processOnionPacket(onionPkt, &sharedSecret, assocData)
	if err != nil {
		return nil, err
	}

	if err := r.log.Put(hashPrefix, incomingCltv); err != nil {
		return nil, err
	}

	return packet, nil
}

// processOnionPacket handles core packet processing logic.
func processOnionPacket(onionPkt *OnionPacket, sharedSecret *Hash256,
	assocData []byte) (*ProcessedPacket, error) {

	innerPkt, hopPayload, err := unwrapPacket(onionPkt, sharedSecret, assocData)
	if err != nil {
		return nil, err
	}

	action := determineAction(hopPayload)

	return &ProcessedPacket{
		Action:                 action,
		ForwardingInstructions: hopPayload.HopData,
		Payload:                *hopPayload,
		NextPacket:             innerPkt,
	}, nil
}

// Helper function to determine the next action based on hop payload.
func determineAction(payload *HopPayload) ProcessCode {
	if bytes.Equal(payload.HMAC[:], make([]byte, sha256.Size)) {
		return ExitNode
	}
	return MoreHops
}

// generateSharedSecret derives a shared secret for the given onion packet.
func (r *Router) generateSharedSecret(ephemPub *btcec.PublicKey,
	blindingPoint *btcec.PublicKey) (Hash256, error) {
	// Perform ECDH to compute shared secret.
	return computeSharedSecret(r.onionKey, ephemPub, blindingPoint)
}
