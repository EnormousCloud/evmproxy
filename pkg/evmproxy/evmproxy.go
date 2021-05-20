package evmproxy

import (
	"github.com/ethereum/go-ethereum/common"
)

var part1 = []byte{
	0x36, // 0x00000000      36             calldatasize          cds
	0x3d, // 0x00000001      3d             returndatasize        0 cds
	0x3d, // 0x00000002      3d             returndatasize        0 0 cds
	0x37, // 0x00000003      37             calldatacopy
	0x3d, // 0x00000004      3d             returndatasize        0
	0x3d, // 0x00000005      3d             returndatasize        0 0
	0x3d, // 0x00000006      3d             returndatasize        0 0 0
	0x36, // 0x00000007      36             calldatasize          cds 0 0 0
	0x3d, // 0x00000008      3d             returndatasize        0 cds 0 0 0
	0x6f, // 0x00000009      6f bebebebebe. push16 0xbebebebe     0xbebe 0 cds 0 0 0
}

var part2 = []byte{
	0x5a, // 0x0000001e      5a             gas                   gas 0xbebe 0 cds 0 0 0
	0xf4, // 0x0000001f      f4             delegatecall          suc 0
	0x3d, // 0x00000020      3d             returndatasize        rds suc 0
	0x82, // 0x00000021      82             dup3                  0 rds suc 0
	0x80, // 0x00000022      80             dup1                  0 0 rds suc 0
	0x3e, // 0x00000023      3e             returndatacopy        suc 0
	0x90, // 0x00000024      90             swap1                 0 suc
	0x3d, // 0x00000025      3d             returndatasize        rds 0 suc
	0x91, // 0x00000026      91             swap2                 suc 0 rds
	0x60, // 0x00000027      6027           push1 0x27            0x27 suc 0 rds
	0x27,
	0x57, // 0x00000029      57             jumpi                 0 rds
	0xfd, // 0x0000002a      fd             revert
	0x5b, // 0x0000002b      5b             jumpdest              0 rds
	0xf3, // 0x0000002c      f3             return
}

func GetBytecode(addr common.Address) []byte {
	return append(append(part1, addr.Bytes()...), part2...)
}
