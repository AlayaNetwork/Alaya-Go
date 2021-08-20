package platonstats

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/AlayaNetwork/Alaya-Go/core/types"
	"github.com/AlayaNetwork/Alaya-Go/rlp"

	"github.com/AlayaNetwork/Alaya-Go/common"
)

var (
	address = common.MustBech32ToAddress("lax1e8su9veseal8t8eyj0zuw49nfkvtqlun2sy6wj")
)

func Test_statsBlockExt(t *testing.T) {
	blockEnc := common.FromHex("f90264f901fda00000000000000000000000000000000000000000000000000000000000000000948888f1f195afa192cfee860698584c030f4c9db1a0ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080832fefd8825208845506eb0780b8510376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f861f85f800a82c35094095e7baea6a6c7c4c2dfeb977efac326af552d870a8023a09bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094fa08a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b180")
	var block *types.Block
	if err := rlp.DecodeBytes(blockEnc, block); err != nil {
		t.Fatal("decode block data error: ", err)
	}

	brief := collectBrief(block)

	blockJsonMapping, err := jsonBlock(block)
	if err != nil {
		t.Fatal("marshal block to json string error", err)
	}
	statsBlockExt := &StatsBlockExt{
		BlockType:   brief.BlockType,
		Epoch:       brief.Epoch,
		NodeID:      brief.NodeID,
		NodeAddress: brief.NodeAddress,
		//Block:        convertBlock(block),
		Block: blockJsonMapping,
	}

	jsonBytes, err := json.Marshal(statsBlockExt)
	if err != nil {
		t.Fatal("marshal platon stats block message to json string error", err)
	} else {
		t.Log("marshal platon stats block", "blockNumber", block.NumberU64(), "json", string(jsonBytes))
	}

}
func T2estUrl(t *testing.T) {
	re := regexp.MustCompile("([^:@]*)(:([^@]*))?@(.+)")
	url := "center:myPasswordd@ws://localhost:1900"
	parts := re.FindStringSubmatch(url)
	for i := 0; i < len(parts); i++ {
		t.Logf("url parts: [%d]%s", i, parts[i])
	}
}
