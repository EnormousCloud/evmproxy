package evmproxy

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestGetBytecode(t *testing.T) {

	a, err := common.NewMixedcaseAddressFromString("0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed")
	if err != nil {
		t.Fatal(err)
	}
	got := GetBytecode(a.Address())
	assert.Greater(t, len(got), 40)
	fmt.Printf("%d %x\n", len(got), got)
}
