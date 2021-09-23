package rawdb

import (
	"testing"

	"github.com/AlayaNetwork/Alaya-Go/crypto"

	"github.com/stretchr/testify/assert"

	"github.com/AlayaNetwork/Alaya-Go/common"
	"github.com/AlayaNetwork/Alaya-Go/params"
)

func TestReadWriteChainConfig(t *testing.T) {

	chainDb := NewMemoryDatabase()
	config := ReadChainConfig(chainDb, common.ZeroHash)
	assert.Nil(t, config, "the chainConfig is not nil")

	WriteChainConfig(chainDb, common.ZeroHash, params.AlayaChainConfig)
	config = ReadChainConfig(chainDb, common.ZeroHash)
	assert.NotNil(t, config, "the chainConfig is nil")

}

func TestReadWritePreimages(t *testing.T) {
	blob := []byte("test")
	hash := crypto.Keccak256Hash(blob)

	chainDb := NewMemoryDatabase()
	preimage := ReadPreimage(chainDb, hash)
	assert.Equal(t, 0, len(preimage), "the preimage is not nil")

	preimages := make(map[common.Hash][]byte)
	preimages[hash] = common.CopyBytes(blob)
	WritePreimages(chainDb, preimages)

	preimage = ReadPreimage(chainDb, hash)
	assert.NotEqual(t, 0, len(preimage), "the preimage is nil")
}
