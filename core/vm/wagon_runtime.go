package vm

import (
	"crypto/sha256"
	"encoding/binary"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bn256"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bulletproof/tx"

	"github.com/holiman/uint256"

	"golang.org/x/crypto/ripemd160"

	"github.com/PlatONnetwork/PlatON-Go/common"
	imath "github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"github.com/PlatONnetwork/wagon/exec"
	"github.com/PlatONnetwork/wagon/wasm"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/params"

	"math/big"
	"reflect"
)

type VMContext struct {
	evm            *EVM
	contract       *Contract
	config         Config
	gasTable       params.GasTable
	db             StateDB
	Input          []byte
	CallOut        []byte
	Output         []byte
	VariableResult []byte
	readOnly       bool // Whether to throw on stateful modifications
	Revert         bool
	Log            *WasmLogger
}

var (
	ptrSize = uint32(4)
)

func NewVMContext(evm *EVM, contract *Contract, config Config, gasTable params.GasTable, db StateDB) *VMContext {
	return &VMContext{
		evm:      evm,
		contract: contract,
		config:   config,
		gasTable: gasTable,
		db:       db,
	}
}
func addFuncExport(m *wasm.Module, sig wasm.FunctionSig, function wasm.Function, export wasm.ExportEntry) {
	typesLen := len(m.Types.Entries)
	m.Types.Entries = append(m.Types.Entries, sig)
	function.Sig = &m.Types.Entries[typesLen]
	funcLen := len(m.FunctionIndexSpace)
	m.FunctionIndexSpace = append(m.FunctionIndexSpace, function)
	export.Index = uint32(funcLen)
	m.Export.Entries[export.FieldStr] = export
}
func NewHostModule() *wasm.Module {
	m := wasm.NewModule()
	m.Export.Entries = make(map[string]wasm.ExportEntry)

	// uint8_t platon_gas_price(uint8_t gas_price[16])
	// func $platon_gas_price(param $0 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(GasPrice),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_gas_price",
			Kind:     wasm.ExternalFunction,
		},
	)
	// void platon_block_hash(int64_t num, uint8_t hash[32])
	// func $platon_block_hash(param $0 i64) (param $1 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI64, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(BlockHash),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_block_hash",
			Kind:     wasm.ExternalFunction,
		},
	)

	// uint64_t platon_block_number()
	// func $platon_block_number (result i64)
	addFuncExport(m,
		wasm.FunctionSig{
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI64},
		},
		wasm.Function{
			Host: reflect.ValueOf(BlockNumber),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_block_number",
			Kind:     wasm.ExternalFunction,
		},
	)

	// uint64_t platon_gas_limit()
	// func $platon_gas_limit (result i64)
	addFuncExport(m,
		wasm.FunctionSig{
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI64},
		},
		wasm.Function{
			Host: reflect.ValueOf(GasLimit),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_gas_limit",
			Kind:     wasm.ExternalFunction,
		},
	)
	// uint64_t platon_gas()
	// func $platon_gas (result i64)
	addFuncExport(m,
		wasm.FunctionSig{
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI64},
		},
		wasm.Function{
			Host: reflect.ValueOf(Gas),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_gas",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int64_t platon_timestamp()
	// func $timestamp (result i64)
	addFuncExport(m,
		wasm.FunctionSig{
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI64},
		},
		wasm.Function{
			Host: reflect.ValueOf(Timestamp),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_timestamp",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_coinbase(uint8_t addr[20])
	// func $platon_coinbase (param $0 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Coinbase),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_coinbase",
			Kind:     wasm.ExternalFunction,
		},
	)

	// uint8_t platon_balance(const uint8_t addr[20], uint8_t balance[16])
	// func $platon_balance (param $0 i32) (param $1 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Balance),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_balance",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_origin(uint8_t addr[20])
	// func $platon_origin (param $0 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Origin),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_origin",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_caller(uint8_t addr[20])
	// func $platon_caller (param $0 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Caller),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_caller",
			Kind:     wasm.ExternalFunction,
		},
	)

	// uint8_t platon_call_value(uint8_t val[16])
	// func $platon_call_value (param $0 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(CallValue),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_call_value",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_address(uint8_t addr[20])
	// func $platon_address  (param $0 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Address),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_address",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_sha3(const uint8_t *src, size_t srcLen, uint8_t *dest, size_t destLen)
	// func $platon_sha3  (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Sha3),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_sha3",
			Kind:     wasm.ExternalFunction,
		},
	)

	// uint64_t platon_caller_nonce()
	// func $platon_caller_nonce  (result i64)
	addFuncExport(m,
		wasm.FunctionSig{
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI64},
		},
		wasm.Function{
			Host: reflect.ValueOf(CallerNonce),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_caller_nonce",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t platon_transfer(const uint8_t to[20], const uint8_t *amount, size_t len)
	// func $platon_transfer  (param $1 i32) (param $2 i32) (param $3 i32) (result i64)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Transfer),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_transfer",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_set_state(const uint8_t* key, size_t klen, const uint8_t *value, size_t vlen)
	// func $platon_set_state (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(SetState),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_set_state",
			Kind:     wasm.ExternalFunction,
		},
	)

	// size_t platon_get_state_length (const uint8_t* key, size_t klen)
	// func $platon_get_state_length (param $0 i32) (param $1 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{

			Host: reflect.ValueOf(GetStateLength),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_get_state_length",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t platon_get_state(const uint8_t *key, size_t klen, uint8_t *value, size_t vlen)
	// func $platon_get_state (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(GetState),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_get_state",
			Kind:     wasm.ExternalFunction,
		},
	)

	// size_t platon_get_input_length()
	// func $platon_get_input_length  (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(GetInputLength),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_get_input_length",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_get_input(const uint8_t *value)
	// func $platon_get_input (param $0 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(GetInput),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_get_input",
			Kind:     wasm.ExternalFunction,
		},
	)

	// size_t platon_get_call_output_length()
	// func $platon_get_call_output_length  (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(GetCallOutputLength),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_get_call_output_length",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_get_call_output(const uint8_t *value)
	// func $platon_get_call_output (param $0 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(GetCallOutput),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_get_call_output",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_return(const uint8_t *value, const size_t len)
	// func $platon_return(param $0 i32) (param $1 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(ReturnContract),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_return",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_revert()
	// func $platon_return()
	addFuncExport(m,
		wasm.FunctionSig{},
		wasm.Function{
			Host: reflect.ValueOf(Revert),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_revert",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_panic()
	// func $platon_panic()
	addFuncExport(m,
		wasm.FunctionSig{},
		wasm.Function{
			Host: reflect.ValueOf(Panic),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_panic",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_debug(const uint8_t *dst, size_t len)
	// func $platon_debug (param i32 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Debug),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_debug",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t platon_call(const uint8_t to[20], const uint8_t *args, size_t args_len, const uint8_t *value, size_t value_len, const uint8_t *call_cost, size_t call_cost_len);
	// func $platon_call  (param $0 i32) (param $1 i32) (param $2 i32) (param $1 i32) (param $2 i32) (param $1 i32) (param $2 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(CallContract),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_call",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t platon_delegate_call(const uint8_t to[20], const uint8_t *args, size_t args_len, const uint8_t *call_cost, size_t call_cost_len);
	// func $platon_delegate_call (param $0 i32) (param $1 i32) (param $2 i32) (param $1 i32) (param $2 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(DelegateCallContract),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_delegate_call",
			Kind:     wasm.ExternalFunction,
		},
	)

	/*	// int32_t platon_static_call(const uint8_t to[20], const uint8_t* args, size_t argsLen, const uint8_t* callCost, size_t callCostLen);
		// func $platon_static_call (param $0 i32) (param $1 i32) (param $2 i32) (param $1 i32) (param $2 i32) (result i32)
		addFuncExport(m,
			wasm.FunctionSig{
				ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
				ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
			},
			wasm.Function{
				Host: reflect.ValueOf(StaticCallContract),
				Body: &wasm.FunctionBody{},
			},
			wasm.ExportEntry{
				FieldStr: "platon_static_call",
				Kind:     wasm.ExternalFunction,
			},
		)*/

	// int32_t platon_destroy(const uint8_t to[20])
	// func $platon_destroy (param $0 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(DestroyContract),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_destroy",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t platon_migrate(uint8_t new_addr[20], const uint8_t *args, size_t args_len, const uint8_t *value, size_t value_len, const uint8_t *call_cost, size_t call_cost_len);
	// func $platon_migrate  (param $1 i32) (param $2 i32) (param $0 i32) (param $1 i32) (param $2 i32) (param $1 i32) (param $2 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32,
				wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(MigrateContract),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_migrate",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t platon_clone_migrate(const uint8_t old_addr[20], uint8_t newAddr[20], const uint8_t* args, size_t argsLen, const uint8_t* value, size_t valueLen, const uint8_t* callCost, size_t callCostLen);
	// func $platon_clone_migrate (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (param $6 i32) (param $7 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32,
				wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(MigrateCloneContract),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_clone_migrate",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_event(const uint8_t *topic, size_t topic_len, const uint8_t *args, size_t args_len);
	// func $platon_event (param $0 i32) (param $1 i32) (param $0 i32) (param $1 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(EmitEvent),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_event",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t platon_ecrecover(const uint8_t hash[32], const uint8_t* sig, const uint8_t sig_len, uint8_t addr[20])
	// func platon_ecrecover (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Ecrecover),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_ecrecover",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_ripemd160(const uint8_t *input, uint32_t input_len, uint8_t hash[20])
	// func platon_ripemd160 (param $0 i32) (param $1 i32) (param $2 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Ripemd160),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_ripemd160",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_sha256(const uint8_t *input, uint32_t input_len, uint8_t hash[32])
	// func platon_sha256 (param $0 i32) (param $1 i32) (param $2 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Sha256),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_sha256",
			Kind:     wasm.ExternalFunction,
		},
	)

	// size_t rlp_u128_size(uint64_t heigh, uint64_t low);
	// func rlp_u128_size (param $0 i64) (param $1 i64) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI64, wasm.ValueTypeI64},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(RlpU128Size),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "rlp_u128_size",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_rlp_u128(uint64_t heigh, uint64_t low, void * dest);
	// func platon_rlp_u128 (param $0 i64) (param $1 i64) (param $2 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI64, wasm.ValueTypeI64, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(RlpU128),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_rlp_u128",
			Kind:     wasm.ExternalFunction,
		},
	)

	// size_t rlp_bytes_size(const void *data, size_t len);
	// func rlp_bytes_size (param $0 i32) (param $1 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(RlpBytesSize),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "rlp_bytes_size",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_rlp_bytes(const void *data, size_t len, void * dest);
	// func platon_rlp_bytes (param $0 i32) (param $1 i32) (param $2 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(RlpBytes),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_rlp_bytes",
			Kind:     wasm.ExternalFunction,
		},
	)

	// size_t rlp_list_size(size_t len);
	// func rlp_list_size (param $0 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(RlpListSize),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "rlp_list_size",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_rlp_list(const void *data, size_t len, void * dest);
	// func platon_rlp_list (param $0 i32) (param $1 i32) (param $2 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(RlpList),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_rlp_list",
			Kind:     wasm.ExternalFunction,
		},
	)

	// size_t platon_contract_code_length(const uint8_t addr[20]);
	// func platon_contract_code_length (param $0 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(ContractCodeLength),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_contract_code_length",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t platon_contract_code(const uint8_t addr[20], uint8_t *code, size_t code_length);
	// func platon_contract_code (param $0 i32) (param $1 i32) (param $2 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(ContractCode),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_contract_code",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t platon_deploy(uint8_t newAddr[20], const uint8_t* args, size_t argsLen, const uint8_t* value, size_t valueLen, const uint8_t* callCost, size_t callCostLen);
	// func $platon_deploy (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (param $6 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32,
				wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(PlatonDeploy),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_deploy",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t platon_clone(const uint8_t old_addr[20], uint8_t newAddr[20], const uint8_t* args, size_t argsLen, const uint8_t* value, size_t valueLen, const uint8_t* callCost, size_t callCostLen);
	// func $platon_clone (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (param $6 i32) (param $7 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32,
				wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(PlatonClone),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_clone",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t platon_confidential_tx_verify(const uint8_t *tx_data, size_t tx_len);
	// func $platon_confidential_tx_verify (param $0 i32) (param $1 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{

			Host: reflect.ValueOf(ConfidentialTxVerify),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_confidential_tx_verify",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t platon_variable_length_result(uint8_t *result, size_t result_len);
	// func $platon_variable_length_result(param $0 i32) (param $1 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{

			Host: reflect.ValueOf(VariableLengthResult),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_variable_length_result",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int bn256_g1_add(byte x1[32], byte y1[32], byte x2[32], byte y2[32], byte x3[32], byte y3[32]);
	// func $bn256_g1_add(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Bn256G1Add),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "bn256_g1_add",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int bn256_g1_mul(byte x1[32], byte y1[32], byte bigint[32], byte x2[32], byte y2[32]);
	// func $bn256_g1_mul(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Bn256G1Mul),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "bn256_g1_mul",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int bn256_g2_add(byte x11[32], byte y11[32], byte x12[32], byte y12[32], byte x21[32], byte y21[32], byte x22[32], byte y22[32], byte x31[32], byte y31[32], byte x32[32], byte y32[32]);
	// func $bn256_g2_add(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (param $6 i32) (param $7 i32) (param $8 i32) (param $9 i32) (param $10 i32) (param $11 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Bn256G2Add),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "bn256_g2_add",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int bn256_g2_mul(byte x11[32], byte y11[32], byte x12[32], byte y12[32], byte bigint[32] byte x21[32], byte y21[32], byte x22[32], byte y22[32]);
	// func $bn256_g2_mul(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (param $6 i32) (param $7 i32) (param $8 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Bn256G2Mul),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "bn256_g2_mul",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int bn256_pairing(uint8_t* x1[], uint8_t* y1[], uint8_t* x21[], uint8_t* y21[], uint8_t* x22[], uint8_t* y22[], size_t len);
	// func $bn256_pairing(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (param $6 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Bn256Pairing),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "bn256_pairing",
			Kind:     wasm.ExternalFunction,
		},
	)


	// uint32_t bigint_binary_operator(const uint8_t *left, uint8_t left_negative, size_t left_arr_size, const uint8_t *right, uint8_t right_negative, size_t right_arr_size, uint8_t *result, size_t result_arr_size, BinaryOperator binary_operator);
	// func $bigint_binary_operator(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (param $6 i32) (param $7 i32) (param $8 i32) (result  i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(BigintBinaryOperator),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "bigint_binary_operator",
			Kind:     wasm.ExternalFunction,
		},
	)

	// uint32_t bigint_exp_mod(const uint8_t *left, uint8_t left_negative, size_t left_arr_size, const uint8_t *right, uint8_t right_negative, size_t right_arr_size, const uint8_t *mod, uint8_t mod_negative, size_t mod_arr_size, uint8_t *result, size_t result_arr_siz);
	// func $bigint_exp_mod(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (param $6 i32) (param $7 i32) (param $8 i32) (param $9 i32) (param $10 i32) (result  i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(BigintExpMod),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "bigint_exp_mod",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t bigint_cmp(const uint8_t *left, uint8_t left_negative, size_t left_arr_size, const uint8_t *right, uint8_t right_negative, size_t right_arr_size);
	// func $bigint_cmp(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (result  i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(BigintCmp),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "bigint_cmp",
			Kind:     wasm.ExternalFunction,
		},
	)

	// uint32_t bigint_sh(const uint8_t *origin, uint8_t origin_negative, size_t origin_arr_size, uint8_t *result, size_t result_arr_size, uint32_t shift_num, ShiftDirection direction);
	// func $bigint_sh(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (param $6 i32)  (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(BigintSh),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "bigint_sh",
			Kind:     wasm.ExternalFunction,
		},
	)

	// uint32_t string_convert_operator(const uint8_t *str, size_t str_len, uint8_t *result, size_t result_arr_size);
	// func $string_convert_operator(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(StringConvertOperator),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "string_convert_operator",
			Kind:     wasm.ExternalFunction,
		},
	)

	return m
}

func checkGas(ctx *VMContext, gas uint64) {
	if !ctx.contract.UseGas(gas) {
		panic(ErrOutOfGas)
	}
}

func mustReadAt(proc *exec.Process, p []byte, off int64) {
	if _, err := proc.ReadAt(p, off); err != nil {
		panic(err)
	}
}

func mustWriteAt(proc *exec.Process, p []byte, off int64) {
	if _, err := proc.WriteAt(p, off); err != nil {
		panic(err)
	}
}
func GasPrice(proc *exec.Process, gasPrice uint32) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	value := ctx.evm.GasPrice.Bytes()
	_, err := proc.WriteAt(value, int64(gasPrice))
	if err != nil {
		panic(err)
	}

	return uint32(len(value))
}

func BlockHash(proc *exec.Process, num uint64, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasExtStep)
	blockHash := ctx.evm.GetHash(num)
	_, err := proc.WriteAt(blockHash.Bytes(), int64(dst))
	if nil != err {
		panic(err)
	}
}

func BlockNumber(proc *exec.Process) uint64 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	return ctx.evm.BlockNumber.Uint64()
}

func GasLimit(proc *exec.Process) uint64 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	return ctx.evm.GasLimit
}

func Gas(proc *exec.Process) uint64 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	return ctx.contract.Gas
}

func Timestamp(proc *exec.Process) int64 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	return ctx.evm.Time.Int64()
}

func Coinbase(proc *exec.Process, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	coinBase := ctx.evm.Coinbase
	_, err := proc.WriteAt(coinBase.Bytes(), int64(dst))
	if nil != err {
		panic(err)
	}
}

func Balance(proc *exec.Process, dst uint32, balance uint32) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, ctx.gasTable.Balance)
	var addr common.Address
	_, err := proc.ReadAt(addr[:], int64(dst))
	if nil != err {
		panic(err)
	}
	value := ctx.evm.StateDB.GetBalance(addr).Bytes()
	_, err = proc.WriteAt(value, int64(balance))
	if nil != err {
		panic(err)
	}
	return uint32(len(value))
}

func Origin(proc *exec.Process, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	_, err := proc.WriteAt(ctx.evm.Origin.Bytes(), int64(dst))
	if nil != err {
		panic(err)
	}
}

func Caller(proc *exec.Process, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	_, err := proc.WriteAt(ctx.contract.caller.Address().Bytes(), int64(dst))
	if nil != err {
		panic(err)
	}
}

// define: uint8_t callValue();
func CallValue(proc *exec.Process, dst uint32) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	value := ctx.contract.value.Bytes()
	_, err := proc.WriteAt(value, int64(dst))
	if nil != err {
		panic(err)
	}
	return uint32(len(value))
}

// define: void address(char hash[20]);
func Address(proc *exec.Process, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	_, err := proc.WriteAt(ctx.contract.Address().Bytes(), int64(dst))
	if nil != err {
		panic(err)
	}
}

// define: void sha3(char *src, size_t srcLen, char *dest, size_t destLen);
func Sha3(proc *exec.Process, src uint32, srcLen uint32, dst uint32, dstLen uint32) {
	ctx := proc.HostCtx().(*VMContext)
	var (
		gas      uint64
		wordGas  uint64
		overflow bool
	)

	if wordGas, overflow = imath.SafeMul(toWordSize(uint64(srcLen)), params.Sha3WordGas); overflow {
		panic(errGasUintOverflow)
	}
	if gas, overflow = imath.SafeAdd(wordGas, params.Sha3Gas); overflow {
		panic(errGasUintOverflow)
	}

	checkGas(ctx, gas)

	data := make([]byte, srcLen)
	_, err := proc.ReadAt(data, int64(src))
	if nil != err {
		panic(err)
	}
	hash := crypto.Keccak256(data)
	if int(dstLen) < len(hash) {
		panic(ErrWASMSha3DstToShort)
	}
	_, err = proc.WriteAt(hash, int64(dst))
	if nil != err {
		panic(err)
	}
}

func CallerNonce(proc *exec.Process) uint64 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	addr := ctx.contract.Caller()
	return ctx.evm.StateDB.GetNonce(addr)
}

func Transfer(proc *exec.Process, dst uint32, amount uint32, len uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)
	address := make([]byte, common.AddressLength)

	_, err := proc.ReadAt(address, int64(dst))
	if nil != err {
		panic(err)
	}

	value := make([]byte, len)
	_, err = proc.ReadAt(value, int64(amount))
	if nil != err {
		panic(err)
	}
	bValue := new(big.Int)
	// 256 bits
	bValue.SetBytes(value)
	bValue = imath.U256(bValue)
	addr := common.BytesToAddress(address)

	transfersValue := bValue.Sign() != 0
	gas := ctx.gasTable.Calls
	if transfersValue {
		gas += params.CallValueTransferGas
	}
	gasTemp, err := callGas(ctx.contract.Gas, params.TxGas, uint256.NewInt().SetUint64(ctx.contract.Gas))
	if nil != err {
		panic(err)
	}
	ctx.evm.callGasTemp = gasTemp
	gas, overflow := imath.SafeAdd(gas, ctx.evm.callGasTemp)
	if overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)

	gas = ctx.evm.callGasTemp
	if transfersValue {
		if gas, overflow = imath.SafeAdd(gas, params.CallStipend); overflow {
			panic(errGasUintOverflow)
		}
	}

	_, returnGas, err := ctx.evm.Call(ctx.contract, addr, nil, gas, bValue)

	var status int32

	if err != nil {
		status = -1
	} else {
		status = 0
	}
	ctx.contract.Gas += returnGas

	return status
}

// storage external function
func SetState(proc *exec.Process, key uint32, keyLen uint32, val uint32, valLen uint32) {
	ctx := proc.HostCtx().(*VMContext)
	if ctx.readOnly {
		panic(ErrWASMWriteProtection)
	}

	keyBuf := make([]byte, keyLen)
	_, err := proc.ReadAt(keyBuf, int64(key))
	if nil != err {
		panic(err)
	}

	currentValue := ctx.evm.StateDB.GetState(ctx.contract.Address(), keyBuf)
	oldWordSize := toWordSize(uint64(keyLen) + uint64(len(currentValue)))
	newWordSize := toWordSize(uint64(keyLen) + uint64(valLen))

	switch {
	case 0 == len(currentValue) && 0 != valLen:
		checkGas(ctx, newWordSize*params.SstoreSetGas)
	case 0 != len(currentValue) && 0 == valLen:
		ctx.evm.StateDB.AddRefund(oldWordSize * params.SstoreRefundGas)
		checkGas(ctx, oldWordSize*params.SstoreClearGas)
	default:
		var (
			addWordSize    uint64 = 0
			deleteWordSize uint64 = 0
			resetWordSize  uint64 = 0
		)

		if newWordSize >= oldWordSize {
			addWordSize = newWordSize - oldWordSize
			resetWordSize = toWordSize(uint64(len(currentValue)))
		} else {
			deleteWordSize = oldWordSize - newWordSize
			resetWordSize = toWordSize(uint64(valLen))
		}

		if 0 == resetWordSize {
			resetWordSize = 1
		}

		checkGas(ctx, addWordSize*params.SstoreSetGas)
		ctx.evm.StateDB.AddRefund(deleteWordSize * params.SstoreRefundGas)
		checkGas(ctx, deleteWordSize*params.SstoreClearGas)
		checkGas(ctx, resetWordSize*params.SstoreResetGas)
	}

	valBuf := make([]byte, valLen)
	_, err = proc.ReadAt(valBuf, int64(val))
	if nil != err {
		panic(err)
	}
	ctx.evm.StateDB.SetState(ctx.contract.Address(), keyBuf, valBuf)
}

func GetStateLength(proc *exec.Process, key uint32, keyLen uint32) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	keyBuf := make([]byte, keyLen)
	_, err := proc.ReadAt(keyBuf, int64(key))
	if nil != err {
		panic(err)
	}
	val := ctx.evm.StateDB.GetState(ctx.contract.Address(), keyBuf)

	checkGas(ctx, ctx.gasTable.SLoad)

	return uint32(len(val))
}

func GetState(proc *exec.Process, key uint32, keyLen uint32, val uint32, valLen uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, ctx.gasTable.SLoad)

	keyBuf := make([]byte, keyLen)
	_, err := proc.ReadAt(keyBuf, int64(key))
	if nil != err {
		panic(err)
	}
	valBuf := ctx.evm.StateDB.GetState(ctx.contract.Address(), keyBuf)
	vlen := len(valBuf)
	if uint32(vlen) > valLen {
		return -1
	}

	_, err = proc.WriteAt(valBuf, int64(val))
	if nil != err {
		panic(err)
	}
	return int32(vlen)
}

func GetInputLength(proc *exec.Process) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	return uint32(len(ctx.Input))
}

func GetInput(proc *exec.Process, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	_, err := proc.WriteAt(ctx.Input, int64(dst))
	if err != nil {
		panic(err)
	}
}

func GetCallOutputLength(proc *exec.Process) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	return uint32(len(ctx.CallOut))
}

func GetCallOutput(proc *exec.Process, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	_, err := proc.WriteAt(ctx.CallOut, int64(dst))
	if err != nil {
		panic(err)
	}
}

func ReturnContract(proc *exec.Process, dst uint32, len uint32) {
	ctx := proc.HostCtx().(*VMContext)
	var (
		gas      uint64
		overflow bool
	)
	if gas, overflow = imath.SafeMul(params.MemoryGas, uint64(len)); overflow {
		panic(errGasUintOverflow)
	}
	if gas, overflow = imath.SafeAdd(gas, GasQuickStep); overflow {
		panic(errGasUintOverflow)
	}

	checkGas(ctx, gas)
	ctx.Output = make([]byte, len)
	_, err := proc.ReadAt(ctx.Output, int64(dst))
	if err != nil {
		panic(err)
	}
}

func Revert(proc *exec.Process) {
	ctx := proc.HostCtx().(*VMContext)
	ctx.Revert = true
	proc.Terminate()
}

func Panic(proc *exec.Process) {
	panic(ErrWASMPanicOp)
}

func Debug(proc *exec.Process, dst uint32, len uint32) {
	ctx := proc.HostCtx().(*VMContext)
	var (
		gas      uint64
		overflow bool
	)

	if gas, overflow = imath.SafeMul(params.MemoryGas, uint64(len)); overflow {
		panic(errGasUintOverflow)
	}
	if gas, overflow = imath.SafeAdd(gas, GasSlowStep); overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)

	buf := make([]byte, len)
	_, err := proc.ReadAt(buf, int64(dst))
	if nil != err {
		panic(err)
	}
	ctx.Log.Debug("WASM:" + string(buf) + "\n")
	ctx.Log.Flush()
}

func CallContract(proc *exec.Process, addrPtr, args, argsLen, val, valLen, callCost, callCostLen uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	address := make([]byte, common.AddressLength)
	_, err := proc.ReadAt(address, int64(addrPtr))
	if nil != err {
		panic(err)
	}
	addr := common.BytesToAddress(address)

	input := make([]byte, argsLen)
	_, err = proc.ReadAt(input, int64(args))
	if nil != err {
		panic(err)
	}

	value := make([]byte, valLen)
	_, err = proc.ReadAt(value, int64(val))
	if nil != err {
		panic(err)
	}
	bValue := new(big.Int)
	// 256 bits
	bValue.SetBytes(value)
	bValue = imath.U256(bValue)

	cost := make([]byte, callCostLen)
	_, err = proc.ReadAt(cost, int64(callCost))
	if nil != err {
		panic(err)
	}
	bCost := new(big.Int)
	// 256 bits
	bCost.SetBytes(cost)
	bCost = imath.U256(bCost)
	if bCost.Cmp(common.Big0) == 0 {
		bCost = new(big.Int).SetUint64(ctx.contract.Gas)
	}

	gas := ctx.gasTable.Calls
	transfersValue := bValue.Sign() != 0
	if transfersValue && ctx.evm.StateDB.Empty(addr) {
		gas += params.CallNewAccountGas
	}

	if transfersValue {
		gas += params.CallValueTransferGas
	}

	gasTemp, err := callGas(ctx.contract.Gas, gas, uint256.NewInt().SetBytes(bCost.Bytes()))
	if nil != err {
		panic(err)
	}
	ctx.evm.callGasTemp = gasTemp
	gas, overflow := imath.SafeAdd(gas, ctx.evm.callGasTemp)
	if overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)

	gas = ctx.evm.callGasTemp
	if bValue.Sign() != 0 {
		if gas, overflow = imath.SafeAdd(gas, params.CallStipend); overflow {
			panic(errGasUintOverflow)
		}
	}

	ret, returnGas, err := ctx.evm.Call(ctx.contract, addr, input, gas, bValue)

	var status int32

	if err != nil {
		status = -1
	} else {
		status = 0
	}

	if err == nil || err == errExecutionReverted {
		ctx.CallOut = ret
	}

	if nil != err {
		if _, ok := err.(*common.BizError); ok {
			ctx.CallOut = ret
		}
	}

	ctx.contract.Gas += returnGas

	return status
}

func DelegateCallContract(proc *exec.Process, addrPtr, params, paramsLen, callCost, callCostLen uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	address := make([]byte, common.AddressLength)
	_, err := proc.ReadAt(address, int64(addrPtr))
	if nil != err {
		panic(err)
	}
	addr := common.BytesToAddress(address)

	input := make([]byte, paramsLen)
	_, err = proc.ReadAt(input, int64(params))
	if nil != err {
		panic(err)
	}

	cost := make([]byte, callCostLen)
	_, err = proc.ReadAt(cost, int64(callCost))
	if nil != err {
		panic(err)
	}
	bCost := new(big.Int)
	// 256 bits
	bCost.SetBytes(cost)
	bCost = imath.U256(bCost)
	if bCost.Cmp(common.Big0) == 0 {
		bCost = new(big.Int).SetUint64(ctx.contract.Gas)
	}

	gasTemp, err := callGas(ctx.contract.Gas, ctx.gasTable.Calls, uint256.NewInt().SetBytes(bCost.Bytes()))
	if nil != err {
		panic(err)
	}
	ctx.evm.callGasTemp = gasTemp
	gas, overflow := imath.SafeAdd(ctx.gasTable.Calls, ctx.evm.callGasTemp)
	if overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)

	gas = ctx.evm.callGasTemp

	ret, returnGas, err := ctx.evm.DelegateCall(ctx.contract, addr, input, gas)

	var status int32

	if err != nil {
		status = -1
	} else {
		status = 0
	}

	if err == nil || err == errExecutionReverted {
		ctx.CallOut = ret
	}

	if nil != err {
		if _, ok := err.(*common.BizError); ok {
			ctx.CallOut = ret
		}
	}

	ctx.contract.Gas += returnGas

	return status
}

func StaticCallContract(proc *exec.Process, addrPtr, params, paramsLen, callCost, callCostLen uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	address := make([]byte, common.AddressLength)
	_, err := proc.ReadAt(address, int64(addrPtr))
	if nil != err {
		panic(err)
	}
	addr := common.BytesToAddress(address)

	input := make([]byte, paramsLen)
	_, err = proc.ReadAt(input, int64(params))
	if nil != err {
		panic(err)
	}

	cost := make([]byte, callCostLen)
	_, err = proc.ReadAt(cost, int64(callCost))
	if nil != err {
		panic(err)
	}
	bCost := new(big.Int)
	// 256 bits
	bCost.SetBytes(cost)
	bCost = imath.U256(bCost)
	if bCost.Cmp(common.Big0) == 0 {
		bCost = new(big.Int).SetUint64(ctx.contract.Gas)
	}

	gasTemp, err := callGas(ctx.contract.Gas, ctx.gasTable.Calls, uint256.NewInt().SetBytes(bCost.Bytes()))
	if nil != err {
		panic(err)
	}

	ctx.evm.callGasTemp = gasTemp
	gas, overflow := imath.SafeAdd(ctx.gasTable.Calls, ctx.evm.callGasTemp)
	if overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)

	gas = ctx.evm.callGasTemp

	ret, returnGas, err := ctx.evm.StaticCall(ctx.contract, addr, input, gas)

	var status int32

	if err != nil {
		status = -1
	} else {
		status = 0
	}

	if err == nil || err == errExecutionReverted {
		ctx.CallOut = ret
	}

	if nil != err {
		if _, ok := err.(*common.BizError); ok {
			ctx.CallOut = ret
		}
	}

	ctx.contract.Gas += returnGas

	return status
}

func DestroyContract(proc *exec.Process, addrPtr uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	if ctx.readOnly {
		panic(ErrWASMWriteProtection)
	}

	address := make([]byte, common.AddressLength)
	_, err := proc.ReadAt(address, int64(addrPtr))
	if nil != err {
		panic(err)
	}
	addr := common.BytesToAddress(address)

	contractAddr := ctx.contract.Address()

	gas := ctx.gasTable.Suicide
	if ctx.evm.StateDB.Empty(addr) && ctx.evm.StateDB.GetBalance(contractAddr).Sign() != 0 {
		gas += ctx.gasTable.CreateBySuicide
	}

	if !ctx.evm.StateDB.HasSuicided(ctx.contract.Address()) {
		ctx.evm.StateDB.AddRefund(params.SuicideRefundGas)
	}
	checkGas(ctx, gas)

	balance := ctx.evm.StateDB.GetBalance(contractAddr)

	ctx.evm.StateDB.AddBalance(addr, balance)

	ctx.evm.StateDB.Suicide(contractAddr)

	return 0
}

func MigrateContract(proc *exec.Process, newAddr, args, argsLen, val, valLen, callCost, callCostLen uint32) int32 {
	// Cost of gas
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, ctx.gasTable.SLoad)

	// get input
	input := make([]byte, argsLen)
	_, err := proc.ReadAt(input, int64(args))
	if nil != err {
		panic(err)
	}

	return MigrateInnerContract(proc, newAddr, val, valLen, callCost, callCostLen, input)
}

func MigrateInnerContract(proc *exec.Process, newAddr, val, valLen, callCost, callCostLen uint32, input []byte) int32 {
	ctx := proc.HostCtx().(*VMContext)

	if ctx.readOnly {
		panic(ErrWASMWriteProtection)
	}

	// check call depth
	if ctx.evm.depth > int(params.CallCreateDepth) {
		return -1
	}

	value := make([]byte, valLen)
	_, err := proc.ReadAt(value, int64(val))
	if nil != err {
		panic(err)
	}
	bValue := new(big.Int)
	// 256 bits
	bValue.SetBytes(value)
	bValue = imath.U256(bValue)

	cost := make([]byte, callCostLen)
	_, err = proc.ReadAt(cost, int64(callCost))
	if nil != err {
		panic(err)
	}
	bCost := new(big.Int)
	// 256 bits
	bCost.SetBytes(cost)
	bCost = imath.U256(bCost)
	if bCost.Cmp(common.Big0) == 0 {
		bCost = new(big.Int).SetUint64(ctx.contract.Gas)
	}

	gas := MigrateContractGas
	if bValue.Sign() != 0 {
		gas += params.CallNewAccountGas
	}
	gasTemp, err := callGas(ctx.contract.Gas, gas, uint256.NewInt().SetBytes(bCost.Bytes()))
	if nil != err {
		panic(err)
	}

	ctx.evm.callGasTemp = gasTemp
	gas, overflow := imath.SafeAdd(gas, ctx.evm.callGasTemp)
	if overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)
	gas = ctx.evm.callGasTemp

	sender := ctx.contract.caller.Address()
	oldContract := ctx.contract.Address()

	// check code of old contract
	oldCode := ctx.evm.StateDB.GetCode(oldContract)
	if len(oldCode) == 0 {
		panic(ErrWASMOldContractCodeNotExists)
	}

	// check balance of sender
	if !ctx.evm.CanTransfer(ctx.evm.StateDB, sender, bValue) {
		return -1
	}

	nonce := ctx.evm.StateDB.GetNonce(oldContract)
	// create new contract address
	newContract := crypto.CreateAddress(oldContract, nonce)
	ctx.evm.StateDB.SetNonce(oldContract, nonce+1)

	// Ensure there's no existing contract already at the designated address
	contractHash := ctx.evm.StateDB.GetCodeHash(newContract)
	if ctx.evm.StateDB.GetNonce(newContract) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
		panic(ErrContractAddressCollision)
	}

	// Create a new account on the state
	snapshotForSnapshotDB, snapshotForStateDB := ctx.evm.DBSnapshot()
	ctx.evm.StateDB.CreateAccount(newContract)
	ctx.evm.StateDB.SetNonce(newContract, 1)

	oldBalance := new(big.Int).Set(ctx.evm.StateDB.GetBalance(oldContract))

	// migrate balance from old contract to new contract
	ctx.evm.Transfer(ctx.evm.StateDB, oldContract, newContract, oldBalance)
	// transfer balance from sender to new contract
	ctx.evm.Transfer(ctx.evm.StateDB, sender, newContract, bValue)

	// migrate stateObject storage from old contract to new contract
	ctx.evm.StateDB.MigrateStorage(oldContract, newContract)

	// suicided the old contract
	ctx.evm.StateDB.Suicide(oldContract)

	balance := new(big.Int).Add(bValue, oldBalance)

	// init new contract context
	contract := NewContract(AccountRef(sender), AccountRef(newContract), balance, gas)
	contract.SetCallCode(&newContract, crypto.Keccak256Hash(input), input)
	contract.DeployContract = true

	// deploy new contract
	ret, err := run(ctx.evm, contract, nil, false)

	// check whether the max code size has been exceeded
	maxCodeSizeExceeded := len(ret) > params.MaxCodeSize
	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil && !maxCodeSizeExceeded {
		createDataGas := uint64(len(ret)) * params.CreateWasmDataGas
		if contract.UseGas(createDataGas) {
			ctx.evm.StateDB.SetCode(newContract, ret)
		} else {
			err = ErrCodeStoreOutOfGas
		}
	}

	// When an error was returned by the VM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if maxCodeSizeExceeded || (err != nil && err != ErrCodeStoreOutOfGas) {
		ctx.evm.RevertToDBSnapshot(snapshotForSnapshotDB, snapshotForStateDB)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}

	// Assign err if contract code size exceeds the max while the err is still empty.
	if maxCodeSizeExceeded && err == nil {
		err = errMaxCodeSizeExceeded
	}
	ctx.contract.Gas += contract.Gas
	if nil != err {
		panic(err)
	}
	_, err = proc.WriteAt(newContract.Bytes(), int64(newAddr))
	if nil != err {
		panic(err)
	}

	return 0
}

func MigrateCloneContract(proc *exec.Process, oldAddr, newAddr, args, argsLen, val, valLen, callCost, callCostLen uint32) int32 {
	// Cost of gas
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, ctx.gasTable.SLoad)

	// get old contract address
	address := make([]byte, common.AddressLength)
	if _, err := proc.ReadAt(address, int64(oldAddr)); nil != err {
		panic(err)
	}
	contractAddress := common.BytesToAddress(address)

	// get contract code
	contractCode := ctx.evm.StateDB.GetCode(contractAddress)
	if 0 == len(contractCode) {
		return -1
	}

	// get init args
	initArgs := make([]byte, argsLen)
	if _, err := proc.ReadAt(initArgs, int64(args)); nil != err {
		panic(err)
	}

	// rlp encode
	createData := struct {
		Code     []byte
		InitArgs []byte
	}{
		Code:     contractCode,
		InitArgs: initArgs,
	}
	input, err := rlp.EncodeToBytes(createData)
	if nil != err {
		panic(err)
	}

	// add magic number
	realInput := make([]byte, len(WasmInterp), len(WasmInterp)+len(input))
	copy(realInput, WasmInterp[0:])
	realInput = append(realInput, input...)

	return MigrateInnerContract(proc, newAddr, val, valLen, callCost, callCostLen, realInput)
}

func EmitEvent(proc *exec.Process, indexesPtr, indexesLen, args, argsLen uint32) {
	ctx := proc.HostCtx().(*VMContext)

	if ctx.readOnly {
		panic(ErrWASMWriteProtection)
	}

	topics := make([]common.Hash, 0)

	if indexesLen != 0 {

		indexes := make([]byte, indexesLen)
		_, err := proc.ReadAt(indexes, int64(indexesPtr))
		if nil != err {
			panic(err)
		}

		content, _, err := rlp.SplitList(indexes)
		if nil != err {
			panic(err)
		}

		topicCount, err := rlp.CountValues(content)
		if nil != err {
			panic(err)
		}
		if topicCount > WasmTopicNum {
			panic(ErrWASMEventCountToLarge)
		}

		decodeTopics := func(b []byte) ([]byte, []byte, error) {
			member, rest, err := rlp.SplitString(b)
			if nil != err {
				panic(err)
			}
			return member, rest, nil
		}

		for len(content) > 0 {
			mem, tail, err := decodeTopics(content)
			if nil != err {
				panic(err)
			}

			if len(mem) > common.HashLength {
				panic(ErrWASMEventContentToLong)
			}

			topics = append(topics, common.BytesToHash(mem))
			content = tail
		}

	}

	input := make([]byte, argsLen)
	_, err := proc.ReadAt(input, int64(args))
	if nil != err {
		panic(err)
	}

	gas, err := logGas(uint64(len(topics)), uint64(argsLen))
	if nil != err {
		panic(err)
	}
	checkGas(ctx, gas)

	bn := ctx.evm.BlockNumber.Uint64()

	addLog(ctx.evm.StateDB, ctx.contract.Address(), topics, input, bn)
}

func Ecrecover(proc *exec.Process, hashPtr, sigPtr, sigLen, addrPtr uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	checkGas(ctx, params.EcrecoverGas)
	hash := make([]byte, 32)
	_, err := proc.ReadAt(hash, int64(hashPtr))
	if err != nil {
		panic(err)
	}

	sig := make([]byte, sigLen)
	_, err = proc.ReadAt(sig, int64(sigPtr))
	if err != nil {
		panic(err)
	}

	pubKey, err := crypto.Ecrecover(hash, sig)
	if err != nil {
		return -1
	}

	if _, err = proc.WriteAt(crypto.Keccak256(pubKey[1:])[12:], int64(addrPtr)); err != nil {
		return -1
	}
	return 0
}

func Ripemd160(proc *exec.Process, inputPtr, inputLen uint32, outputPtr uint32) {
	ctx := proc.HostCtx().(*VMContext)
	var (
		gas      uint64
		overflow bool
	)
	if gas, overflow = imath.SafeMul(toWordSize(uint64(inputLen)), params.Ripemd160PerWordGas); overflow {
		panic(errGasUintOverflow)
	}
	if gas, overflow = imath.SafeAdd(gas, params.Ripemd160BaseGas); overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)

	input := make([]byte, inputLen)
	_, err := proc.ReadAt(input, int64(inputPtr))
	if err != nil {
		panic(err)
	}
	ripemd := ripemd160.New()
	ripemd.Write(input)
	output := ripemd.Sum(nil)
	proc.WriteAt(output, int64(outputPtr))
}

func Sha256(proc *exec.Process, inputPtr, inputLen uint32, outputPtr uint32) {
	ctx := proc.HostCtx().(*VMContext)
	var (
		gas      uint64
		overflow bool
	)

	if gas, overflow = imath.SafeMul(toWordSize(uint64(inputLen)), params.Sha256PerWordGas); overflow {
		panic(errGasUintOverflow)
	}
	if gas, overflow = imath.SafeAdd(gas, params.Sha256BaseGas); overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)

	input := make([]byte, inputLen)
	_, err := proc.ReadAt(input, int64(inputPtr))
	if err != nil {
		panic(err)
	}
	h := sha256.Sum256(input)

	proc.WriteAt(h[:], int64(outputPtr))
}

func addLog(state StateDB, address common.Address, topics []common.Hash, data []byte, bn uint64) {
	log := &types.Log{
		Address:     address,
		Topics:      topics,
		Data:        data,
		BlockNumber: bn,
	}
	state.AddLog(log)
}

func logGas(topicNum, dataSize uint64) (uint64, error) {
	gas := params.LogGas
	var overflow bool
	if gas, overflow = imath.SafeAdd(gas, topicNum*params.LogTopicGas); overflow {
		return 0, errGasUintOverflow
	}

	var logSizeGas uint64
	if logSizeGas, overflow = imath.SafeMul(dataSize, params.LogDataGas); overflow {
		return 0, errGasUintOverflow
	}
	if gas, overflow = imath.SafeAdd(gas, logSizeGas); overflow {
		return 0, errGasUintOverflow
	}
	return gas, nil
}

// rlp encode
const (
	rlpMaxLengthBytes  = 8
	rlpDataImmLenStart = 0x80
	rlpListStart       = 0xc0
	rlpDataImmLenCount = rlpListStart - rlpDataImmLenStart - rlpMaxLengthBytes
	rlpDataIndLenZero  = rlpDataImmLenStart + rlpDataImmLenCount - 1
	rlpListImmLenCount = 256 - rlpListStart - rlpMaxLengthBytes
	rlpListIndLenZero  = rlpListStart + rlpListImmLenCount - 1
)

func bigEndian(num uint64) []byte {
	var (
		temp   []byte
		result []byte
	)

	for {
		if 0 == num {
			break
		}
		temp = append(temp, byte(num))
		num = num >> 8
	}

	for i := len(temp); i != 0; i = i - 1 {
		result = append(result, temp[i-1])
	}

	return result
}

// size_t rlp_u128_size(uint64_t heigh, uint64_t low);
func RlpU128Size(proc *exec.Process, heigh uint64, low uint64) uint32 {
	// Cost of gas
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)

	var size uint32 = 0
	if (0 == heigh && 0 == low) || (0 == heigh && low < rlpDataImmLenStart) {
		size = 1
	} else {
		dataLen := len(bigEndian(heigh)) + len(bigEndian(low))
		size = uint32(dataLen) + 1
	}

	return size
}

// void platon_rlp_u128(uint64_t heigh, uint64_t low, void * dest);
func RlpU128(proc *exec.Process, heigh uint64, low uint64, dest uint32) {
	ctx := proc.HostCtx().(*VMContext)
	var (
		gas      uint64
		overflow bool
	)

	// rlp result
	var data []byte
	if 0 == heigh && 0 == low {
		data = append(data, rlpDataImmLenStart)
	} else if 0 == heigh && low < rlpDataImmLenStart {
		data = append(data, byte(low))
	} else {
		temp := bigEndian(heigh)
		temp = append(temp, bigEndian(low)...)
		data = append(data, byte(rlpDataImmLenStart+len(temp)))
		data = append(data, temp...)
	}

	// Cost of gas
	if gas, overflow = imath.SafeMul(params.MemoryGas, uint64(len(data))); overflow {
		panic(errGasUintOverflow)
	}
	if gas, overflow = imath.SafeAdd(gas, GasQuickStep); overflow {
		panic(errGasUintOverflow)
	}

	checkGas(ctx, gas)

	// write data
	_, err := proc.WriteAt(data, int64(dest))
	if nil != err {
		panic(err)
	}
}

// size_t rlp_bytes_size(const void *data, size_t len);
func RlpBytesSize(proc *exec.Process, src uint32, length uint32) uint32 {
	// Cost of gas
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)

	// read data
	data := make([]byte, 1)
	_, err := proc.ReadAt(data, int64(src))
	if nil != err {
		panic(err)
	}

	if 1 == length && data[0] < rlpDataImmLenStart {
		return 1
	}

	if length < rlpDataImmLenCount {
		return length + 1
	}

	return uint32(len(bigEndian(uint64(length)))) + length + 1
}

// void platon_rlp_bytes(const void *data, size_t len, void * dest);
func RlpBytes(proc *exec.Process, src uint32, length uint32, dest uint32) {
	// Cost of gas
	ctx := proc.HostCtx().(*VMContext)
	var (
		gas      uint64
		overflow bool
	)

	// read data
	data := make([]byte, length)
	_, err := proc.ReadAt(data, int64(src))
	if nil != err {
		panic(err)
	}

	// get prefixData
	var prefixData []byte
	if 1 == length && data[0] < rlpDataImmLenStart {
		prefixData = []byte{}
	} else if length < rlpDataImmLenCount {
		prefixData = append(prefixData, byte(rlpDataImmLenStart+length))
	} else {
		lengthBytes := bigEndian(uint64(length))
		if len(lengthBytes)+rlpDataIndLenZero > 0xff {
			panic(ErrWASMRlpItemCountTooLarge)
		}
		prefixData = append(prefixData, byte(len(lengthBytes)+rlpDataIndLenZero))
		prefixData = append(prefixData, lengthBytes...)
	}

	if gas, overflow = imath.SafeMul(params.MemoryGas, uint64(len(prefixData)+len(data))); overflow {
		panic(errGasUintOverflow)
	}
	if gas, overflow = imath.SafeAdd(gas, GasQuickStep); overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)

	// write prefix
	_, err = proc.WriteAt(prefixData, int64(dest))
	if nil != err {
		panic(err)
	}

	// write data
	_, err = proc.WriteAt(data, int64(dest)+int64(len(prefixData)))
	if nil != err {
		panic(err)
	}
}

// size_t rlp_list_size(const void *data, size_t len);
func RlpListSize(proc *exec.Process, length uint32) uint32 {
	// Cost of gas
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)

	if length < rlpListImmLenCount {
		return length + 1
	}

	return uint32(len(bigEndian(uint64(length)))) + length + 1
}

// void platon_rlp_list(const void *data, size_t len, void * dest);
func RlpList(proc *exec.Process, src uint32, length uint32, dest uint32) {
	// Cost of gas
	ctx := proc.HostCtx().(*VMContext)
	var (
		gas      uint64
		overflow bool
	)

	// read data
	data := make([]byte, length)
	_, err := proc.ReadAt(data, int64(src))
	if nil != err {
		panic(err)
	}

	// get result
	var prefixData []byte
	if length < rlpListImmLenCount {
		prefixData = append(prefixData, byte(rlpListStart+length))
	} else {
		lengthBytes := bigEndian(uint64(length))
		if len(lengthBytes)+rlpDataIndLenZero > 0xff {
			panic(ErrWASMRlpItemCountTooLarge)
		}
		prefixData = append(prefixData, byte(len(lengthBytes)+rlpListIndLenZero))
		prefixData = append(prefixData, lengthBytes...)
	}

	if gas, overflow = imath.SafeMul(params.MemoryGas, uint64(len(prefixData)+len(data))); overflow {
		panic(errGasUintOverflow)
	}
	if gas, overflow = imath.SafeAdd(gas, GasQuickStep); overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)

	// write prefix
	_, err = proc.WriteAt(prefixData, int64(dest))
	if nil != err {
		panic(err)
	}

	// write data
	_, err = proc.WriteAt(data, int64(dest)+int64(len(prefixData)))
	if nil != err {
		panic(err)
	}
}

// size_t platon_contract_code_length(const uint8_t addr[20]);
func ContractCodeLength(proc *exec.Process, addrPtr uint32) uint32 {
	// Cost of gas
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, ctx.gasTable.SLoad)

	// get contract address
	address := make([]byte, common.AddressLength)
	_, err := proc.ReadAt(address, int64(addrPtr))
	if nil != err {
		panic(err)
	}
	contractAddress := common.BytesToAddress(address)

	// get contract code
	contractCode := ctx.evm.StateDB.GetCode(contractAddress)
	return uint32(len(contractCode))
}

// int32_t platon_contract_code(const uint8_t addr[20], uint8_t *code, size_t code_length);
func ContractCode(proc *exec.Process, addrPtr uint32, code uint32, codeLen uint32) int32 {
	// Cost of gas
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, ctx.gasTable.SLoad)

	// get contract address
	address := make([]byte, common.AddressLength)
	_, err := proc.ReadAt(address, int64(addrPtr))
	if nil != err {
		panic(err)
	}
	contractAddress := common.BytesToAddress(address)

	// get contract code
	contractCode := ctx.evm.StateDB.GetCode(contractAddress)
	if 0 == len(contractCode) || uint32(len(contractCode)) > codeLen {
		return 0
	}

	// write data
	_, err = proc.WriteAt(contractCode, int64(code))
	if nil != err {
		panic(err)
	}

	return int32(len(contractCode))
}

// int32_t platon_deploy(uint8_t newAddr[20], const uint8_t* args, size_t argsLen, const uint8_t* value, size_t valueLen, const uint8_t* callCost, size_t callCostLen);
func PlatonDeploy(proc *exec.Process, newAddr, args, argsLen, val, valLen, callCost, callCostLen uint32) int32 {
	// Cost of gas
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, ctx.gasTable.SLoad)

	// get input
	input := make([]byte, argsLen)
	_, err := proc.ReadAt(input, int64(args))
	if nil != err {
		panic(err)
	}

	return CreateContract(proc, newAddr, val, valLen, callCost, callCostLen, input)
}

func CreateContract(proc *exec.Process, newAddr, val, valLen, callCost, callCostLen uint32, input []byte) int32 {
	ctx := proc.HostCtx().(*VMContext)

	if ctx.readOnly {
		panic(ErrWASMWriteProtection)
	}

	// check call depth
	if ctx.evm.depth > int(params.CallCreateDepth) {
		return -1
	}

	// get transfer value
	valueBytes := make([]byte, valLen)
	if _, err := proc.ReadAt(valueBytes, int64(val)); nil != err {
		panic(err)
	}
	bigValue := new(big.Int)
	bigValue.SetBytes(valueBytes)

	// check balance of sender
	sender := ctx.contract.caller.Address()
	if !ctx.evm.CanTransfer(ctx.evm.StateDB, sender, bigValue) {
		return -1
	}

	// get gas limit
	costBytes := make([]byte, callCostLen)
	if _, err := proc.ReadAt(costBytes, int64(callCost)); nil != err {
		panic(err)
	}
	costValue := new(big.Int)
	costValue.SetBytes(costBytes)
	if costValue.Cmp(common.Big0) == 0 {
		costValue.SetUint64(ctx.contract.Gas)
	}

	gas := params.CallNewAccountGas
	gasTemp, err := callGas(ctx.contract.Gas, gas, uint256.NewInt().SetBytes(costValue.Bytes()))
	if nil != err {
		panic(err)
	}

	ctx.evm.callGasTemp = gasTemp
	gas, overflow := imath.SafeAdd(gas, ctx.evm.callGasTemp)
	if overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)
	gas = ctx.evm.callGasTemp

	// generate new address
	oldContract := ctx.contract.Address()
	nonce := ctx.evm.StateDB.GetNonce(oldContract)
	newContract := crypto.CreateAddress(oldContract, nonce)
	ctx.evm.StateDB.SetNonce(oldContract, nonce+1)
	contractHash := ctx.evm.StateDB.GetCodeHash(newContract)
	if ctx.evm.StateDB.GetNonce(newContract) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
		panic(ErrContractAddressCollision)
	}

	// create new account
	snapshotForSnapshotDB, snapshotForStateDB := ctx.evm.DBSnapshot()
	ctx.evm.StateDB.CreateAccount(newContract)
	ctx.evm.StateDB.SetNonce(newContract, 1)

	// transfer value to new account
	ctx.evm.Transfer(ctx.evm.StateDB, sender, newContract, bigValue)

	// init new contract context
	contract := NewContract(AccountRef(sender), AccountRef(newContract), ctx.evm.StateDB.GetBalance(sender), gas)
	contract.SetCallCode(&newContract, crypto.Keccak256Hash(input), input)
	contract.DeployContract = true

	// deploy new contract
	ret, err := run(ctx.evm, contract, nil, false)

	// check whether the max code size has been exceeded
	maxCodeSizeExceeded := len(ret) > params.MaxCodeSize
	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil && !maxCodeSizeExceeded {
		createDataGas := uint64(len(ret)) * params.CreateWasmDataGas
		if contract.UseGas(createDataGas) {
			ctx.evm.StateDB.SetCode(newContract, ret)
		} else {
			err = ErrCodeStoreOutOfGas
		}
	}

	// When an error was returned by the VM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if maxCodeSizeExceeded || (err != nil && err != ErrCodeStoreOutOfGas) {
		ctx.evm.RevertToDBSnapshot(snapshotForSnapshotDB, snapshotForStateDB)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}

	// Assign err if contract code size exceeds the max while the err is still empty.
	if maxCodeSizeExceeded && err == nil {
		err = errMaxCodeSizeExceeded
	}
	ctx.contract.Gas += contract.Gas
	if nil != err {
		panic(err)
	}
	_, err = proc.WriteAt(newContract.Bytes(), int64(newAddr))
	if nil != err {
		panic(err)
	}

	return 0
}

// int32_t platon_clone(const uint8_t old_addr[20], uint8_t newAddr[20], const uint8_t* args, size_t argsLen, const uint8_t* value, size_t valueLen, const uint8_t* callCost, size_t callCostLen);
func PlatonClone(proc *exec.Process, oldAddr, newAddr, args, argsLen, val, valLen, callCost, callCostLen uint32) int32 {
	// Cost of gas
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, ctx.gasTable.SLoad)

	// get old contract address
	address := make([]byte, common.AddressLength)
	if _, err := proc.ReadAt(address, int64(oldAddr)); nil != err {
		panic(err)
	}
	contractAddress := common.BytesToAddress(address)

	// get contract code
	contractCode := ctx.evm.StateDB.GetCode(contractAddress)
	if 0 == len(contractCode) {
		return -1
	}

	// get init args
	initArgs := make([]byte, argsLen)
	if _, err := proc.ReadAt(initArgs, int64(args)); nil != err {
		panic(err)
	}

	// rlp encode
	createData := struct {
		Code     []byte
		InitArgs []byte
	}{
		Code:     contractCode,
		InitArgs: initArgs,
	}
	input, err := rlp.EncodeToBytes(createData)
	if nil != err {
		panic(err)
	}

	// add magic number
	realInput := make([]byte, len(WasmInterp), len(WasmInterp)+len(input))
	copy(realInput, WasmInterp[0:])
	realInput = append(realInput, input...)

	// create contract
	return CreateContract(proc, newAddr, val, valLen, callCost, callCostLen, realInput)
}

// size_t platon_confidential_tx_verify(const uint8_t *tx_data, size_t tx_len);
// func $platon_confidential_tx_verify (param $0 i32) (param $1 i32) (result i32)
func ConfidentialTxVerify(proc *exec.Process, txData uint32, txLen uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	checkGas(ctx, params.ConfidentialTxVerifyBaseGas+uint64(txLen/256)*params.ConfidentialTxPerNoteGas)

	// read data
	dataBuf := make([]byte, txLen)
	_, err := proc.ReadAt(dataBuf, int64(txData))
	if nil != err {
		panic(err)
	}

	if result, err := tx.VerifyConfidentialTx(dataBuf); nil != err {
		return -1
	} else {
		ctx.VariableResult = result
		return int32(len(result))
	}

}

// int32_t platon_variable_length_result(uint8_t *result, size_t result_len);
// func $platon_variable_length_result(param $0 i32) (param $1 i32) (result i32)
func VariableLengthResult(proc *exec.Process, result uint32, resultLen uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	checkGas(ctx, ctx.gasTable.SLoad)

	if nil == ctx.VariableResult || len(ctx.VariableResult) != int(resultLen) {
		return -1
	}

	_, err := proc.WriteAt(ctx.VariableResult, int64(result))
	if nil != err {
		panic(err)
	}

	return int32(len(ctx.VariableResult))
}

// int bn256_g1_add(byte x1[32], byte y1[32], byte x2[32], byte y2[32], byte x3[32], byte y3[32]);
// func $bn256_g1_add(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (result i32)
func Bn256G1Add(proc *exec.Process, x1, y1, x2, y2, x3, y3 uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	checkGas(ctx, params.Bn256G1AddGas)
	var x1Bytes [32]byte
	var y1Bytes [32]byte
	var x2Bytes [32]byte
	var y2Bytes [32]byte

	mustReadAt(proc, x1Bytes[:], int64(x1))
	mustReadAt(proc, y1Bytes[:], int64(y1))
	mustReadAt(proc, x2Bytes[:], int64(x2))
	mustReadAt(proc, y2Bytes[:], int64(y2))

	var gx1, gx2 bn256.G1
	if _, err := gx1.Unmarshal(append(x1Bytes[:], y1Bytes[:]...)); err != nil {
		return -1
	}
	if _, err := gx2.Unmarshal(append(x2Bytes[:], y2Bytes[:]...)); err != nil {
		return -1
	}

	gx3 := new(bn256.G1)

	gx3.Add(&gx1, &gx2)
	res := gx3.Marshal()

	mustWriteAt(proc, res[:32], int64(x3))
	mustWriteAt(proc, res[32:], int64(y3))
	return 0
}

// int bn256_g1_mul(byte x1[32], byte y1[32], byte bigint[32], byte x2[32], byte y2[32]);
// func $bn256_g1_mul(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (result i32)
func Bn256G1Mul(proc *exec.Process, x1, y1, bigint, x2, y2 uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, params.Bn256G1ScalarMulGas)

	var x1Bytes [32]byte
	var y1Bytes [32]byte
	var bigIntBytes [32]byte
	mustReadAt(proc, x1Bytes[:], int64(x1))
	mustReadAt(proc, y1Bytes[:], int64(y1))
	mustReadAt(proc, bigIntBytes[:], int64(bigint))

	scalar := new(big.Int).SetBytes(bigIntBytes[:])
	var gx1 bn256.G1
	if _, err := gx1.Unmarshal(append(x1Bytes[:], y1Bytes[:]...)); err != nil {
		return -1
	}

	gx3 := new(bn256.G1)

	gx3.ScalarMult(&gx1, scalar)
	res := gx3.Marshal()

	mustWriteAt(proc, res[:32], int64(x2))
	mustWriteAt(proc, res[32:], int64(y2))
	return 0
}

// int bn256_g2_add(byte x11[32], byte x12[32], byte y11[32], byte y12[32], byte x21[32], byte x22[32], byte y21[32], byte y22[32], byte x31[32], byte x32[32], byte y31[32], byte y32[32]);
// func $bn256_g2_add(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (result i32)
func Bn256G2Add(proc *exec.Process, x11, x12, y11, y12, x21, x22, y21, y22, x31, x32, y31, y32 uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, params.Bn256G2AddGas)

	var x11Bytes [32]byte
	var y11Bytes [32]byte
	var x12Bytes [32]byte
	var y12Bytes [32]byte

	var x21Bytes [32]byte
	var y21Bytes [32]byte
	var x22Bytes [32]byte
	var y22Bytes [32]byte

	mustReadAt(proc, x11Bytes[:], int64(x11))
	mustReadAt(proc, y11Bytes[:], int64(y11))
	mustReadAt(proc, x12Bytes[:], int64(x12))
	mustReadAt(proc, y12Bytes[:], int64(y12))

	mustReadAt(proc, x21Bytes[:], int64(x21))
	mustReadAt(proc, y21Bytes[:], int64(y21))
	mustReadAt(proc, x22Bytes[:], int64(x22))
	mustReadAt(proc, y22Bytes[:], int64(y22))

	var gx1, gx2 bn256.G2
	if _, err := gx1.Unmarshal(append(x11Bytes[:], append(x12Bytes[:], append(y11Bytes[:], y12Bytes[:]...)...)...)); err != nil {
		return -1
	}
	if _, err := gx2.Unmarshal(append(x21Bytes[:], append(x22Bytes[:], append(y21Bytes[:], y22Bytes[:]...)...)...)); err != nil {
		return -1
	}

	gx3 := new(bn256.G2)

	gx3.Add(&gx1, &gx2)
	res := gx3.Marshal()

	mustWriteAt(proc, res[:32], int64(x31))
	mustWriteAt(proc, res[32:64], int64(x32))
	mustWriteAt(proc, res[64:96], int64(y31))
	mustWriteAt(proc, res[96:128], int64(y32))

	return 0
}

// int bn256_g2_mul(byte x11[32], byte x12[32], byte y11[32], byte y12[32], byte bigint[32] byte x21[32], byte x22[32], byte y21[32], byte y22[32]);
// func $bn256_g2_mul(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (param $6 i32) (param $7 i32) (param $8 i32) (result i32)
func Bn256G2Mul(proc *exec.Process, x11, x12, y11, y12, bigint, x21, x22, y21, y22 uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, params.Bn256G2ScalarMulGas)

	var x11Bytes [32]byte
	var y11Bytes [32]byte
	var x12Bytes [32]byte
	var y12Bytes [32]byte
	var bigintBytes [32]byte

	mustReadAt(proc, x11Bytes[:], int64(x11))
	mustReadAt(proc, y11Bytes[:], int64(y11))
	mustReadAt(proc, x12Bytes[:], int64(x12))
	mustReadAt(proc, y12Bytes[:], int64(y12))
	mustReadAt(proc, bigintBytes[:], int64(bigint))

	var gx1 bn256.G2
	if _, err := gx1.Unmarshal(append(x11Bytes[:], append(x12Bytes[:], append(y11Bytes[:], y12Bytes[:]...)...)...)); err != nil {
		return -1
	}

	gx3 := new(bn256.G2)

	gx3.ScalarMult(&gx1, new(big.Int).SetBytes(bigintBytes[:]))
	res := gx3.Marshal()

	mustWriteAt(proc, res[:32], int64(x21))
	mustWriteAt(proc, res[32:64], int64(x22))
	mustWriteAt(proc, res[64:96], int64(y21))
	mustWriteAt(proc, res[96:128], int64(y22))

	return 0
}

// int bn256_pairing(uint8_t* x1[], uint8_t* y1[], uint8_t* x21[], uint8_t* x22[], uint8_t* y21[], uint8_t* y22[], size_t len);
// func $bn256_pairing(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (param $6 i32) (result i32)
func Bn256Pairing(proc *exec.Process, x1, y1, x21, x22, y21, y22, len uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, params.Bn256PairingCheckBaseGas+params.Bn256PairingCheckPerPairGas*uint64(len))

	x1Array := make([]byte, len*ptrSize)
	y1Array := make([]byte, len*ptrSize)

	x21Array := make([]byte, len*ptrSize)
	y21Array := make([]byte, len*ptrSize)
	x22Array := make([]byte, len*ptrSize)
	y22Array := make([]byte, len*ptrSize)

	mustReadAt(proc, x1Array[:], int64(x1))
	mustReadAt(proc, y1Array[:], int64(y1))
	mustReadAt(proc, x21Array[:], int64(x21))
	mustReadAt(proc, y21Array[:], int64(y21))
	mustReadAt(proc, x22Array[:], int64(x22))
	mustReadAt(proc, y22Array[:], int64(y22))

	g1s := make([]*bn256.G1, 0, len)
	g2s := make([]*bn256.G2, 0, len)

	for i := uint32(0); i < len; i++ {
		var x1Bytes [32]byte
		var y1Bytes [32]byte

		var x21Bytes [32]byte
		var y21Bytes [32]byte
		var x22Bytes [32]byte
		var y22Bytes [32]byte

		mustReadAt(proc, x1Bytes[:], int64(binary.LittleEndian.Uint32(x1Array[i*4:(i+1)*4])))
		mustReadAt(proc, y1Bytes[:], int64(binary.LittleEndian.Uint32(y1Array[i*4:(i+1)*4])))
		mustReadAt(proc, x21Bytes[:], int64(binary.LittleEndian.Uint32(x21Array[i*4:(i+1)*4])))
		mustReadAt(proc, y21Bytes[:], int64(binary.LittleEndian.Uint32(y21Array[i*4:(i+1)*4])))
		mustReadAt(proc, x22Bytes[:], int64(binary.LittleEndian.Uint32(x22Array[i*4:(i+1)*4])))
		mustReadAt(proc, y22Bytes[:], int64(binary.LittleEndian.Uint32(y22Array[i*4:(i+1)*4])))

		var gx1 bn256.G1
		if _, err := gx1.Unmarshal(append(x1Bytes[:], y1Bytes[:]...)); err != nil {
			return -2
		}
		g1s = append(g1s, &gx1)

		var gx2 bn256.G2
		if _, err := gx2.Unmarshal(append(x21Bytes[:], append(x22Bytes[:], append(y21Bytes[:], y22Bytes[:]...)...)...)); err != nil {
			return -2
		}
		g2s = append(g2s, &gx2)
	}

	if !bn256.PairingCheck(g1s, g2s) {
		return -1
	}

	return 0
}


type BinaryOperatorType uint32

const (
	BIGINTADD BinaryOperatorType = 0x01
	BIGINTSUB BinaryOperatorType = 0x02
	BIGINTMUL BinaryOperatorType = 0x04
	BIGINTDIV BinaryOperatorType = 0x08
	BIGINTMOD BinaryOperatorType = 0x10
	BIGINTAND BinaryOperatorType = 0x20
	BIGINTOR  BinaryOperatorType = 0x40
	BIGINTXOR BinaryOperatorType = 0x80
)

// uint32_t bigint_binary_operator(const uint8_t *left, uint8_t left_negative, size_t left_arr_size, const uint8_t *right, uint8_t right_negative, size_t right_arr_size, uint8_t *result, size_t result_arr_size, BinaryOperator binary_operator);
// func $bigint_binary_operator(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (param $6 i32) (param $7 i32) (param $8 i32) (result  i32)
func BigintBinaryOperator(proc *exec.Process, left uint32, leftNegative uint32, leftArrSize uint32, right uint32, rightNegative uint32, rightArrSize uint32,
	result uint32, resultArrSize uint32, binaryOperator uint32) uint32 {
	// read left data
	leftDataBuf := make([]byte, leftArrSize)
	_, leftErr := proc.ReadAt(leftDataBuf, int64(left))
	if nil != leftErr {
		panic(leftErr)
	}

	leftOperand := new(big.Int)
	leftOperand.SetBytes(leftDataBuf)

	if leftNegative > 0 {
		bigValueTemp := new(big.Int)
		bigValueTemp.Neg(leftOperand)
		leftOperand = bigValueTemp
	}

	// read right data
	rightDataBuf := make([]byte, rightArrSize)
	_, rightErr := proc.ReadAt(rightDataBuf, int64(right))
	if nil != rightErr {
		panic(rightErr)
	}

	rightOperand := new(big.Int)
	rightOperand.SetBytes(rightDataBuf)

	if rightNegative > 0 {
		bigValueTemp := new(big.Int)
		bigValueTemp.Neg(rightOperand)
		rightOperand = bigValueTemp
	}

	// binary operation
	operationResult := new(big.Int)
	switch BinaryOperatorType(binaryOperator) {
	case BIGINTADD:
		operationResult.Add(leftOperand, rightOperand)
	case BIGINTSUB:
		operationResult.Sub(leftOperand, rightOperand)
	case BIGINTMUL:
		operationResult.Mul(leftOperand, rightOperand)
	case BIGINTDIV:
		operationResult.Div(leftOperand, rightOperand)
	case BIGINTMOD:
		operationResult.Mod(leftOperand, rightOperand)
	case BIGINTAND:
		operationResult.And(leftOperand, rightOperand)
	case BIGINTOR:
		operationResult.Or(leftOperand, rightOperand)
	case BIGINTXOR:
		operationResult.Xor(leftOperand, rightOperand)
	default:
		panic("invalid parameter")
	}

	// write result
	bytesResult := operationResult.Bytes()

	// Clear to zero
	zeroResult := make([]byte, resultArrSize)
	_, clearErr := proc.WriteAt(zeroResult, int64(result))
	if nil != clearErr {
		panic(clearErr)
	}

	// Set to actual value
	bytesRealResult := bytesResult
	if len(bytesRealResult) > int(resultArrSize) {
		begin := len(bytesRealResult) - int(resultArrSize)
		bytesRealResult = bytesResult[begin:]
	}
	result = result + resultArrSize - uint32(len(bytesRealResult))

	_, setErr := proc.WriteAt(bytesRealResult, int64(result))
	if nil != setErr {
		panic(setErr)
	}

	var returnResult uint32 = 0

	if -1 == operationResult.Sign() {
		returnResult |= 0x01
	}

	if len(bytesResult) > int(resultArrSize) {
		returnResult |= 0x02
	}

	return returnResult
}

// uint32_t bigint_exp_mod(const uint8_t *left, uint8_t left_negative, size_t left_arr_size, const uint8_t *right, uint8_t right_negative, size_t right_arr_size, const uint8_t *mod, uint8_t mod_negative, size_t mod_arr_size, uint8_t *result, size_t result_arr_siz);
// func $bigint_exp_mod(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (param $6 i32) (param $7 i32) (param $8 i32) (param $9 i32) (param $10 i32) (result  i32)
func BigintExpMod(proc *exec.Process, left uint32, leftNegative uint32, leftArrSize uint32, right uint32, rightNegative uint32, rightArrSize uint32, mod uint32, modNegative uint32, modArrSize uint32,
	result uint32, resultArrSize uint32) uint32 {
	// read left data
	leftDataBuf := make([]byte, leftArrSize)
	_, leftErr := proc.ReadAt(leftDataBuf, int64(left))
	if nil != leftErr {
		panic(leftErr)
	}

	leftOperand := new(big.Int)
	leftOperand.SetBytes(leftDataBuf)

	if leftNegative > 0 {
		bigValueTemp := new(big.Int)
		bigValueTemp.Neg(leftOperand)
		leftOperand = bigValueTemp
	}

	// read right data
	rightDataBuf := make([]byte, rightArrSize)
	_, rightErr := proc.ReadAt(rightDataBuf, int64(right))
	if nil != rightErr {
		panic(rightErr)
	}

	rightOperand := new(big.Int)
	rightOperand.SetBytes(rightDataBuf)

	if rightNegative > 0 {
		bigValueTemp := new(big.Int)
		bigValueTemp.Neg(rightOperand)
		rightOperand = bigValueTemp
	}

	// read mod data
	modDataBuf := make([]byte, modArrSize)
	_, modErr := proc.ReadAt(modDataBuf, int64(mod))
	if nil != modErr {
		panic(modErr)
	}

	modOperand := new(big.Int)
	modOperand.SetBytes(modDataBuf)

	if modNegative > 0 {
		bigValueTemp := new(big.Int)
		bigValueTemp.Neg(rightOperand)
		modOperand = bigValueTemp
	}

	// binary operation
	operationResult := new(big.Int)
	operationResult.Exp(leftOperand, rightOperand, modOperand)

	// write result
	bytesResult := operationResult.Bytes()

	// Clear to zero
	zeroResult := make([]byte, resultArrSize)
	_, clearErr := proc.WriteAt(zeroResult, int64(result))
	if nil != clearErr {
		panic(clearErr)
	}

	//Set to actual value
	// Set to actual value
	bytesRealResult := bytesResult
	if len(bytesRealResult) > int(resultArrSize) {
		begin := len(bytesRealResult) - int(resultArrSize)
		bytesRealResult = bytesResult[begin:]
	}
	result = result + resultArrSize - uint32(len(bytesRealResult))

	_, setErr := proc.WriteAt(bytesRealResult, int64(result))
	if nil != setErr {
		panic(setErr)
	}

	var returnResult uint32 = 0

	if -1 == operationResult.Sign() {
		returnResult |= 0x01
	}

	if len(bytesResult) > int(resultArrSize) {
		returnResult |= 0x02
	}

	return returnResult
}

// int32_t bigint_cmp(const uint8_t *left, uint8_t left_negative, size_t left_arr_size, const uint8_t *right, uint8_t right_negative, size_t right_arr_size);
// func $bigint_cmp(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (result  i32)
func BigintCmp(proc *exec.Process, left uint32, leftNegative uint32, leftArrSize uint32, right uint32, rightNegative uint32, rightArrSize uint32) int32 {
	// read left data
	leftDataBuf := make([]byte, leftArrSize)
	_, leftErr := proc.ReadAt(leftDataBuf, int64(left))
	if nil != leftErr {
		panic(leftErr)
	}

	leftOperand := new(big.Int)
	leftOperand.SetBytes(leftDataBuf)

	if leftNegative > 0 {
		bigValueTemp := new(big.Int)
		bigValueTemp.Neg(leftOperand)
		leftOperand = bigValueTemp
	}

	// read right data
	rightDataBuf := make([]byte, rightArrSize)
	_, rightErr := proc.ReadAt(rightDataBuf, int64(right))
	if nil != rightErr {
		panic(rightErr)
	}

	rightOperand := new(big.Int)
	rightOperand.SetBytes(rightDataBuf)

	if rightNegative > 0 {
		bigValueTemp := new(big.Int)
		bigValueTemp.Neg(rightOperand)
		rightOperand = bigValueTemp
	}

	// compare
	result := int32(leftOperand.Cmp(rightOperand))
	return result
}

type ShiftDirectionType uint32

const (
	LEFT  ShiftDirectionType = 0x01
	RIGHT ShiftDirectionType = 0x02
)

// uint32_t bigint_sh(const uint8_t *origin, uint8_t origin_negative, size_t origin_arr_size, uint8_t *result, size_t result_arr_size, uint32_t shift_num, ShiftDirection direction);
// func $bigint_sh(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (param $4 i32) (param $5 i32) (param $6 i32)  (result i32)
func BigintSh(proc *exec.Process, origin uint32, originNegative uint32, originArrSize uint32, result uint32, resultArrSize uint32, shiftNum uint32, direction uint32) uint32 {
	// read origin data
	originDataBuf := make([]byte, originArrSize)
	_, originErr := proc.ReadAt(originDataBuf, int64(origin))
	if nil != originErr {
		panic(originErr)
	}

	originOperand := new(big.Int)
	originOperand.SetBytes(originDataBuf)

	if originNegative > 0 {
		bigValueTemp := new(big.Int)
		bigValueTemp.Neg(originOperand)
		originOperand = bigValueTemp
	}

	// shift
	operationResult := new(big.Int)
	switch ShiftDirectionType(direction) {
	case LEFT:
		operationResult.Lsh(originOperand, uint(shiftNum))
	case RIGHT:
		operationResult.Rsh(originOperand, uint(shiftNum))
	}

	// write result
	bytesResult := operationResult.Bytes()

	// Clear to zero
	zeroResult := make([]byte, resultArrSize)
	_, clearErr := proc.WriteAt(zeroResult, int64(result))
	if nil != clearErr {
		panic(clearErr)
	}

	//Set to actual value
	// Set to actual value
	bytesRealResult := bytesResult
	if len(bytesRealResult) > int(resultArrSize) {
		begin := len(bytesRealResult) - int(resultArrSize)
		bytesRealResult = bytesResult[begin:]
	}
	result = result + resultArrSize - uint32(len(bytesRealResult))

	_, setErr := proc.WriteAt(bytesRealResult, int64(result))
	if nil != setErr {
		panic(setErr)
	}

	var returnResult uint32 = 0

	if -1 == operationResult.Sign() {
		returnResult |= 0x01
	}

	if len(bytesResult) > int(resultArrSize) {
		returnResult |= 0x02
	}

	return returnResult
}

// uint32_t string_convert_operator(const uint8_t *str, size_t str_len, uint8_t *result, size_t result_arr_size);
// func $string_convert_operator(param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (result i32)
func StringConvertOperator(proc *exec.Process, str uint32, strLen uint32, result uint32, resultArrSize uint32) uint32 {

	// read string data
	dataBuf := make([]byte, strLen)
	_, err := proc.ReadAt(dataBuf, int64(str))
	if nil != err {
		panic(err)
	}

	strValue := string(dataBuf[:])
	value := new(big.Int)
	value.SetString(strValue, 10)

	// write result
	bytesResult := value.Bytes()

	// Clear to zero
	zeroResult := make([]byte, resultArrSize)
	_, clearErr := proc.WriteAt(zeroResult, int64(result))
	if nil != clearErr {
		panic(clearErr)
	}

	//Set to actual value
	bytesRealResult := bytesResult
	if len(bytesRealResult) > int(resultArrSize) {
		begin := len(bytesRealResult) - int(resultArrSize)
		bytesRealResult = bytesResult[begin:]
	}
	result = result + resultArrSize - uint32(len(bytesRealResult))

	_, setErr := proc.WriteAt(bytesRealResult, int64(result))
	if nil != setErr {
		panic(setErr)
	}

	var returnResult uint32 = 0

	if -1 == value.Sign() {
		returnResult |= 0x01
	}

	if len(bytesResult) > int(resultArrSize) {
		returnResult |= 0x02
	}

	return returnResult
}
