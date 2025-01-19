package sphinx

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
)

// BOLT 4 Test Vectors
var (
	bolt4PubKeys = []string{
		"02eec7245d6b7d2ccb30380bfbe2a3648cd7a942653f5aa340edcea1f283686619",
		"0324653eac434488002cc06bbfb7f10fe18991e35f9fe4302dbea6d2353dc0ab1c",
		"027f31ebc5462c1fdce1b737ecff52d37d75dea43ce11c74d25aa297165faa2007",
		"032c0b7cf95324a07d05398b240174dc0c2be444d96b159aa6c7f7b1e668680991",
		"02edabbd16b41c8371b92ef2f04c1185b4f03b6dcd52ba9b78d9d7c89c8f221145",
	}
	bolt4SessionKey     = bytes.Repeat([]byte{'A'}, 32)
	bolt4AssocData      = bytes.Repeat([]byte{'B'}, 32)
	bolt4FinalPacketHex = "0002eec7245d6b7d2ccb30380bfbe2..."
)

// newTestRoute creates a route with the specified number of hops.
func newTestRoute(numHops int) ([]*Router, *PaymentPath, *[]HopData, *OnionPacket, error) {
	nodes := make([]*Router, numHops)

	for i := range nodes {
		privKey, err := btcec.NewPrivateKey()
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("error generating key for node: %w", err)
		}
		nodes[i] = NewRouter(&PrivKeyECDH{PrivKey: privKey}, NewMemoryReplayLog())
	}

	var route PaymentPath
	for i, node := range nodes {
		hopData := HopData{
			ForwardAmount: uint64(i),
			OutgoingCltv:  uint32(i),
		}
		copy(hopData.NextAddress[:], bytes.Repeat([]byte{byte(i)}, 8))

		hopPayload, err := NewLegacyHopPayload(&hopData)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("error creating hop payload: %w", err)
		}

		route[i] = OnionHop{
			NodePub:    *node.onionKey.PubKey(),
			HopPayload: hopPayload,
		}
	}

	sessionKey, _ := btcec.PrivKeyFromBytes(bolt4SessionKey)
	fwdMsg, err := NewOnionPacket(&route, sessionKey, nil, DeterministicPacketFiller)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("error creating forwarding message: %w", err)
	}

	var hopsData []HopData
	for _, hop := range route {
		hopData, err := hop.HopPayload.HopData()
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("error getting hop data: %w", err)
		}
		hopsData = append(hopsData, *hopData)
	}

	return nodes, &route, &hopsData, fwdMsg, nil
}

// TestBolt4Packet validates the creation of onion packets against BOLT 4 test vectors.
func TestBolt4Packet(t *testing.T) {
	var (
		route    PaymentPath
		hopsData []HopData
	)

	for i, pubKeyHex := range bolt4PubKeys {
		pubKeyBytes, err := hex.DecodeString(pubKeyHex)
		require.NoError(t, err)

		pubKey, err := btcec.ParsePubKey(pubKeyBytes)
		require.NoError(t, err)

		hopData := HopData{
			ForwardAmount: uint64(i),
			OutgoingCltv:  uint32(i),
		}
		copy(hopData.NextAddress[:], bytes.Repeat([]byte{byte(i)}, 8))
		hopsData = append(hopsData, hopData)

		hopPayload, err := NewLegacyHopPayload(&hopData)
		require.NoError(t, err)

		route[i] = OnionHop{
			NodePub:    *pubKey,
			HopPayload: hopPayload,
		}
	}

	finalPacket, err := hex.DecodeString(bolt4FinalPacketHex)
	require.NoError(t, err)

	sessionKey, _ := btcec.PrivKeyFromBytes(bolt4SessionKey)
	pkt, err := NewOnionPacket(&route, sessionKey, bolt4AssocData, DeterministicPacketFiller)
	require.NoError(t, err)

	var b bytes.Buffer
	require.NoError(t, pkt.Encode(&b))
	require.Equal(t, finalPacket, b.Bytes(), "final packet mismatch")
}

// Other tests can follow a similar pattern...

