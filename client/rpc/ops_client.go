package rpc

import "github.com/binance-chain/go-sdk/common/types"

func (c *HTTP) GetStakeValidators() ([]types.Validator, error) {
	rawVal,err:= c.ABCIQuery("custom/stake/validators",nil)
	if err!=nil{
		return nil,err
	}
	var validators []types.Validator
	err = c.cdc.UnmarshalJSON(rawVal.Response.GetValue(), &validators)
	return validators,err

}
