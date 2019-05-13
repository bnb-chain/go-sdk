// +build cgo,ledger

package ledger

import (
	ledger "github.com/binance-chain/ledger-cosmos-go"
)

// If ledger support (build tag) has been enabled, which implies a CGO dependency,
// set the discoverLedger function which is responsible for loading the Ledger
// device at runtime or returning an error.
func init() {
	DiscoverLedger = func() (LedgerSecp256k1, error) {
		device, err := ledger.FindLedgerCosmosUserApp()
		if err != nil {
			return nil, err
		}

		return device, nil
	}
}
