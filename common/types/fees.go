package types

import "fmt"

const (
	OperateFeeType  = "operate"
	TransferFeeType = "transfer"
	DexFeeType      = "dex"

	FeeForProposer = FeeDistributeType(0x01)
	FeeForAll      = FeeDistributeType(0x02)
	FeeFree        = FeeDistributeType(0x03)
)

type FeeDistributeType int8

type FeeParam interface {
	GetParamType() string
	Check() error
}

// dexFee
type DexFeeParam struct {
	DexFeeFields []DexFeeField `json:"dex_fee_fields"`
}

type DexFeeField struct {
	FeeName  string `json:"fee_name"`
	FeeValue int64  `json:"fee_value"`
}

func (p *DexFeeParam) GetParamType() string {
	return DexFeeType
}

func (p *DexFeeParam) isNil() bool {
	for _, d := range p.DexFeeFields {
		if d.FeeValue < 0 {
			return true
		}
	}
	return false
}

func (p *DexFeeParam) Check() error {
	if p.isNil() {
		return fmt.Errorf("Dex fee param is less than 0 ")
	}
	return nil
}

// fixedFee
type FixedFeeParams struct {
	MsgType string            `json:"msg_type"`
	Fee     int64             `json:"fee"`
	FeeFor  FeeDistributeType `json:"fee_for"`
}

func (p *FixedFeeParams) GetParamType() string {
	return OperateFeeType
}

func (p *FixedFeeParams) Check() error {
	if p.FeeFor != FeeForProposer && p.FeeFor != FeeForAll && p.FeeFor != FeeFree {
		return fmt.Errorf("fee_for %d is invalid", p.FeeFor)
	}

	if p.Fee < 0 {
		return fmt.Errorf("fee(%d) should not be negative", p.Fee)
	}
	return nil
}

type TransferFeeParam struct {
	FixedFeeParams    `json:"fixed_fee_params"`
	MultiTransferFee  int64 `json:"multi_transfer_fee"`
	LowerLimitAsMulti int64 `json:"lower_limit_as_multi"`
}

func (p *TransferFeeParam) GetParamType() string {
	return TransferFeeType
}

func (p *TransferFeeParam) Check() error {
	err := p.FixedFeeParams.Check()
	if err != nil {
		return err
	}
	if p.Fee <= 0 || p.MultiTransferFee <= 0 {
		return fmt.Errorf("both fee(%d) and multi_transfer_fee(%d) should be positive", p.Fee, p.MultiTransferFee)
	}
	if p.MultiTransferFee > p.Fee {
		return fmt.Errorf("multi_transfer_fee(%d) should not be bigger than fee(%d)", p.MultiTransferFee, p.Fee)
	}
	if p.LowerLimitAsMulti <= 1 {
		return fmt.Errorf("lower_limit_as_multi should > 1")
	}
	return nil
}
