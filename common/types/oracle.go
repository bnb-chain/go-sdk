package types

import (
	"encoding/binary"

	ctypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/sidechain"
)

type IbcChannelID = ctypes.ChannelID
type IbcChainID = ctypes.ChainID

const (
	prefixLength         = 1
	destIbcChainIDLength = 2
	channelIDLength      = 1
)

var (
	SideChainStorePrefixByIdKey = sidechain.SideChainStorePrefixByIdKey

	PrefixForSendSequenceKey    = sidechain.PrefixForSendSequenceKey
	PrefixForReceiveSequenceKey = sidechain.PrefixForReceiveSequenceKey
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
