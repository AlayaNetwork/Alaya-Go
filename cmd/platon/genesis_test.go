// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var customGenesisTests = []struct {
	genesis string
	query   string
	result  string
}{
	// Plain genesis file without anything extra
	{
		genesis: `{
    "alloc":{
        "1000000000000000000000000000000000000001":{
            "balance":"0"
        },
        "1000000000000000000000000000000000000002":{
            "balance":"0"
        },
        "1000000000000000000000000000000000000003":{
            "balance":"200000000000000000000000000"
        },
        "1000000000000000000000000000000000000004":{
            "balance":"0"
        },
        "1000000000000000000000000000000000000005":{
            "balance":"0"
        },
        "60ceca9c1290ee56b98d4e160ef0453f7c40d219":{
            "balance":"8050000000000000000000000000"
        },
        "55bfd49472fd41211545b01713a9c3a97af78b05":{
            "balance":"2000000000000000000000000000"
        }
    },
    "EconomicModel":{
        "Common":{
            "ExpectedMinutes":4,
            "ValidatorCount":4,
            "AdditionalCycleTime":16
        },
        "Staking":{
            "StakeThreshold":               5000000000000000000000000,
            "MinimumThreshold":             10000000000000000000,
            "EpochValidatorNum":            24,
            "HesitateRatio":                1,
            "UnStakeFreezeRatio":           2,
            "ActiveUnDelegateFreezeRatio":  0
        },
        "Slashing":{
           "PackAmountAbnormal":   6,
           "PackAmountHighAbnormal":  2,
           "PackAmountLowSlashRate":  10,
           "PackAmountHighSlashRate":  50,
           "DuplicateSignHighSlashing": 100
        },
        "Gov": {
            "VersionProposalVote_ConsensusRounds": 4,
            "VersionProposalActive_ConsensusRounds": 5,
            "VersionProposal_SupportRate": 0.667,
            "TextProposalVote_ConsensusRounds": 4,
            "TextProposal_VoteRate": 0.5,
            "TextProposal_SupportRate": 0.667,          
            "CancelProposal_VoteRate": 0.50,
            "CancelProposal_SupportRate": 0.667
        },
        "Reward":{
            "NewBlockRate": 50,
            "PlatONFoundationYear": 10 
        }
    },
    "coinbase":"0x0000000000000000000000000000000000000000",
    "extraData":"",
    "gasLimit":"0x2fefd8",
    "nonce":"0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23",
    "parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000",
    "timestamp":"0x00",
    "config":{
        "cbft":{
            "initialNodes":[
                {
                    "node":"enode://4fcc251cf6bf3ea53a748971a223f5676225ee4380b65c7889a2b491e1551d45fe9fcc19c6af54dcf0d5323b5aa8ee1d919791695082bae1f86dd282dba4150f@0.0.0.0:16789",
                    "blsPubKey":"d341a0c485c9ec00cecf7ea16323c547900f6a1bacb9daacb00c2b8bacee631f75d5d31b75814b7f1ae3a4e18b71c617bc2f230daa0c893746ed87b08b2df93ca4ddde2816b3ac410b9980bcc048521562a3b2d00e900fd777d3cf88ce678719"
                }
            ],
            "epoch":1,
            "amount":10,
            "validatorMode":"ppos",
            "period":10000
        }
    }
}`,
		query:  "platon.getBlock(0).nonce",
		result: "0xd31d37efa7b9e9d7df775d9a6f9ddb6f5e3d6dd6b87b471b71ed9be9a69b7b871c71cd1d7f46f96b7f5ed76f7bedad9f71ddb7000000000000000000000000000000000000000000000000000000000000",
	},
	//Genesis file with only cbft config
	{
		genesis: `{
    "alloc":{
        "1000000000000000000000000000000000000001":{
            "balance":"0"
        },
        "1000000000000000000000000000000000000002":{
            "balance":"0"
        },
        "1000000000000000000000000000000000000003":{
            "balance":"200000000000000000000000000"
        },
        "1000000000000000000000000000000000000004":{
            "balance":"0"
        },
        "1000000000000000000000000000000000000005":{
            "balance":"0"
        },
        "60ceca9c1290ee56b98d4e160ef0453f7c40d219":{
            "balance":"8050000000000000000000000000"
        },
        "55bfd49472fd41211545b01713a9c3a97af78b05":{
            "balance":"2000000000000000000000000000"
        }
    },
    "EconomicModel":{
        "Common":{
            "ExpectedMinutes":4,
            "ValidatorCount":4,
            "AdditionalCycleTime":16
        },
        "Staking":{
            "StakeThreshold":               5000000000000000000000000,
            "MinimumThreshold":             10000000000000000000,
            "EpochValidatorNum":            24,
            "HesitateRatio":                1,
            "UnStakeFreezeRatio":           2,
            "ActiveUnDelegateFreezeRatio":  0
        },
        "Slashing":{
           "PackAmountAbnormal":   6,
           "PackAmountHighAbnormal":  2,
           "PackAmountLowSlashRate":  10,
           "PackAmountHighSlashRate":  50,
           "DuplicateSignHighSlashing": 100
        },
        "Gov": {
            "VersionProposalVote_ConsensusRounds": 4,
            "VersionProposalActive_ConsensusRounds": 5,
            "VersionProposal_SupportRate": 0.667,
            "TextProposalVote_ConsensusRounds": 4,
            "TextProposal_VoteRate": 0.5,
            "TextProposal_SupportRate": 0.667,          
            "CancelProposal_VoteRate": 0.50,
            "CancelProposal_SupportRate": 0.667
        },
        "Reward":{
            "NewBlockRate": 50,
            "PlatONFoundationYear": 10 
        }
    },
    "coinbase":"0x0000000000000000000000000000000000000000",
    "extraData":"",
    "gasLimit":"0x2fefd8",
    "nonce":"0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23",
    "parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000",
    "timestamp":"0x00",
    "config":{
        "cbft":{
            "initialNodes":[
                {
                    "node":"enode://4fcc251cf6bf3ea53a748971a223f5676225ee4380b65c7889a2b491e1551d45fe9fcc19c6af54dcf0d5323b5aa8ee1d919791695082bae1f86dd282dba4150f@0.0.0.0:16789",
                    "blsPubKey":"d341a0c485c9ec00cecf7ea16323c547900f6a1bacb9daacb00c2b8bacee631f75d5d31b75814b7f1ae3a4e18b71c617bc2f230daa0c893746ed87b08b2df93ca4ddde2816b3ac410b9980bcc048521562a3b2d00e900fd777d3cf88ce678719"
                }
            ],
            "epoch":1,
            "amount":10,
            "validatorMode":"ppos",
            "period":10000
        }
    }
}`,
		query:  "platon.getBlock(0).nonce",
		result: "0xd31d37efa7b9e9d7df775d9a6f9ddb6f5e3d6dd6b87b471b71ed9be9a69b7b871c71cd1d7f46f96b7f5ed76f7bedad9f71ddb7000000000000000000000000000000000000000000000000000000000000",
	},
	//Genesis file with specific chain configurations
	{
		genesis: `{
    "alloc":{
        "1000000000000000000000000000000000000001":{
            "balance":"0"
        },
        "1000000000000000000000000000000000000002":{
            "balance":"0"
        },
        "1000000000000000000000000000000000000003":{
            "balance":"200000000000000000000000000"
        },
        "1000000000000000000000000000000000000004":{
            "balance":"0"
        },
        "1000000000000000000000000000000000000005":{
            "balance":"0"
        },
        "60ceca9c1290ee56b98d4e160ef0453f7c40d219":{
            "balance":"8050000000000000000000000000"
        },
        "55bfd49472fd41211545b01713a9c3a97af78b05":{
            "balance":"2000000000000000000000000000"
        }
    },
    "EconomicModel":{
        "Common":{
            "ExpectedMinutes":4,
            "ValidatorCount":4,
            "AdditionalCycleTime":16
        },
        "Staking":{
            "StakeThreshold":               5000000000000000000000000,
            "MinimumThreshold":             10000000000000000000,
            "EpochValidatorNum":            24,
            "HesitateRatio":                1,
            "UnStakeFreezeRatio":           2,
            "ActiveUnDelegateFreezeRatio":  0
        },
        "Slashing":{
           "PackAmountAbnormal":   6,
           "PackAmountHighAbnormal":  2,
           "PackAmountLowSlashRate":  10,
           "PackAmountHighSlashRate":  50,
           "DuplicateSignHighSlashing": 100
        },
        "Gov": {
            "VersionProposalVote_ConsensusRounds": 4,
            "VersionProposalActive_ConsensusRounds": 5,
            "VersionProposal_SupportRate": 0.667,
            "TextProposalVote_ConsensusRounds": 4,
            "TextProposal_VoteRate": 0.5,
            "TextProposal_SupportRate": 0.667,          
            "CancelProposal_VoteRate": 0.50,
            "CancelProposal_SupportRate": 0.667
        },
        "Reward":{
            "NewBlockRate": 50,
            "PlatONFoundationYear": 10 
        }
    },
    "coinbase":"0x0000000000000000000000000000000000000000",
    "extraData":"",
    "gasLimit":"0x2fefd8",
    "nonce":"0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23",
    "parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000",
    "timestamp":"0x00",
    "config":{
        "chainId":101,
        "eip155Block":0,
        "interpreter":"wasm",
        "cbft":{
            "initialNodes":[
                {
                    "node":"enode://4fcc251cf6bf3ea53a748971a223f5676225ee4380b65c7889a2b491e1551d45fe9fcc19c6af54dcf0d5323b5aa8ee1d919791695082bae1f86dd282dba4150f@0.0.0.0:16789",
                    "blsPubKey":"d341a0c485c9ec00cecf7ea16323c547900f6a1bacb9daacb00c2b8bacee631f75d5d31b75814b7f1ae3a4e18b71c617bc2f230daa0c893746ed87b08b2df93ca4ddde2816b3ac410b9980bcc048521562a3b2d00e900fd777d3cf88ce678719"
                }
            ],
            "epoch":1,
            "amount":10,
            "validatorMode":"ppos",
            "period":10000
        }
    }
}`,
		query:  "platon.getBlock(0).nonce",
		result: "0xd31d37efa7b9e9d7df775d9a6f9ddb6f5e3d6dd6b87b471b71ed9be9a69b7b871c71cd1d7f46f96b7f5ed76f7bedad9f71ddb7000000000000000000000000000000000000000000000000000000000000",
	},
}

// Tests that initializing Geth with a custom genesis block and chain definitions
// work properly.
func TestCustomGenesis(t *testing.T) {
	for i, tt := range customGenesisTests {
		// Create a temporary data directory to use and inspect later
		datadir := tmpdir(t)
		defer os.RemoveAll(datadir)

		// Initialize the data directory with the custom genesis block
		json := filepath.Join(datadir, "genesis.json")
		if err := ioutil.WriteFile(json, []byte(tt.genesis), 0600); err != nil {
			t.Fatalf("test %d: failed to write genesis file: %v", i, err)
		}
		runGeth(t, "--datadir", datadir, "init", json).WaitExit()

		// Query the custom genesis block
		geth := runGeth(t,
			"--datadir", datadir, "--maxpeers", "0", "--port", "0",
			"--nodiscover", "--nat", "none", "--ipcdisable",
			"--exec", tt.query, "console")
		t.Log("testi", i)
		geth.ExpectRegexp(tt.result)
		geth.ExpectExit()
	}
}
