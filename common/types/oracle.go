package types

import (
	"encoding/binary"
)

type IbcChannelID uint8
type IbcChainID uint16

const (
	prefixLength         = 1
	destIbcChainIDLength = 2
	channelIDLength      = 1
)

var (
	SideChainStorePrefixByIdKey = []byte{0x01} // prefix for each key to a side chain store prefix, by side chain id

	PrefixForSendSequenceKey    = []byte{0xf0}
	PrefixForReceiveSequenceKey = []byte{0xf1}
)

func GetReceiveSequenceKey(destIbcChainID IbcChainID, channelID IbcChannelID) []byte {
	return buildChannelSequenceKey(destIbcChainID, channelID, PrefixForReceiveSequenceKey)
}

func buildChannelSequenceKey(destIbcChainID IbcChainID, channelID IbcChannelID, prefix []byte) []byte {
	key := make([]byte, prefixLength+destIbcChainIDLength+channelIDLength)

	copy(key[:prefixLength], prefix)
	binary.BigEndian.PutUint16(key[prefixLength:prefixLength+destIbcChainIDLength], uint16(destIbcChainID))
	copy(key[prefixLength+destIbcChainIDLength:], []byte{byte(channelID)})
	return key
}
