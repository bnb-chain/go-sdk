package msg

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/bsc/rlp"

	"github.com/bnb-chain/go-sdk/common/types"
	"github.com/bnb-chain/node/plugins/tokens/swap"
	cTypes "github.com/cosmos/cosmos-sdk/types"
)

var (
	SortJSON            = cTypes.SortJSON
	MustSortJSON        = cTypes.MustSortJSON
	CalculateRandomHash = swap.CalculateRandomHash
	CalculateSwapID     = swap.CalculateSwapID
	HexAddress          = cTypes.HexAddress
	HexEncode           = cTypes.HexEncode
	HexDecode           = cTypes.HexDecode
	Has0xPrefix         = cTypes.Has0xPrefix
)

func noneExistPackageProto() interface{} {
	panic("should not exist such package")
}

func ParseClaimPayload(payload []byte) ([]CrossChainPackage, error) {
	packages := Packages{}
	err := rlp.DecodeBytes(payload, &packages)
	if err != nil {
		return nil, err
	}
	decodedPackage := make([]CrossChainPackage, 0, len(packages))
	for _, pack := range packages {
		ptype, relayerFee, err := DecodePackageHeader(pack.Payload)
		if err != nil {
			return nil, err
		}
		if _, exist := protoMetrics[pack.ChannelId]; !exist {
			return nil, fmt.Errorf("channnel id do not exist")
		}
		proto, exist := protoMetrics[pack.ChannelId][ptype]
		if !exist || proto == nil {
			return nil, fmt.Errorf("package type do not exist")
		}
		content := proto()
		err = rlp.DecodeBytes(pack.Payload[PackageHeaderLength:], content)
		if err != nil {
			return nil, err
		}
		decodedPackage = append(decodedPackage, CrossChainPackage{
			PackageType: ptype,
			RelayFee:    relayerFee,
			Content:     content,
		})
	}
	return decodedPackage, nil
}

func CreateSendMsg(from types.AccAddress, fromCoins types.Coins, transfers []Transfer) SendMsg {
	input := NewInput(from, fromCoins)

	output := make([]Output, 0, len(transfers))
	for _, t := range transfers {
		t.Coins = t.Coins.Sort()
		output = append(output, NewOutput(t.ToAddr, t.Coins))
	}
	msg := NewMsgSend([]Input{input}, output)
	return msg
}
