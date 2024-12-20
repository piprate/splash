package splash

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flowkit/v2"
)

// FlowScriptBuilder is a struct to hold information for running a script
type FlowScriptBuilder struct {
	Connector      *Connector
	FileName       string
	Arguments      []cadence.Value
	ScriptAsString string
}

// Script start a script builder with the inline script as body
func (c *Connector) Script(content string) FlowScriptBuilder {
	return FlowScriptBuilder{
		Connector:      c,
		FileName:       "inline",
		Arguments:      []cadence.Value{},
		ScriptAsString: content,
	}
}

// ScriptFromFile will start a flow script builder
func (c *Connector) ScriptFromFile(filename string) FlowScriptBuilder {
	return FlowScriptBuilder{
		Connector:      c,
		FileName:       filename,
		Arguments:      []cadence.Value{},
		ScriptAsString: "",
	}
}

// AccountArgument add an account as an argument
func (t FlowScriptBuilder) AccountArgument(key string) FlowScriptBuilder {
	f := t.Connector

	account := f.Account(key)
	return t.Argument(cadence.BytesToAddress(account.Address.Bytes()))
}

// RawAccountArgument add an account from a string as an argument
func (t FlowScriptBuilder) RawAccountArgument(key string) FlowScriptBuilder {

	account := flow.HexToAddress(key)
	accountArg := cadence.BytesToAddress(account.Bytes())
	return t.Argument(accountArg)
}

// DateStringAsUnixTimestamp sends a dateString parsed in the timezone as a unix timeszone ufix
func (t FlowScriptBuilder) DateStringAsUnixTimestamp(dateString, timezone string) FlowScriptBuilder {
	return t.UFix64Argument(parseTime(dateString, timezone))
}

// Argument add an argument to the transaction
func (t FlowScriptBuilder) Argument(value cadence.Value) FlowScriptBuilder {
	t.Arguments = append(t.Arguments, value)
	return t
}

// ArgumentList add a list argument to the transaction
func (t FlowScriptBuilder) ArgumentList(values []cadence.Value) FlowScriptBuilder {
	t.Arguments = append(t.Arguments, values...)
	return t
}

// StringArgument add a String Argument to the transaction
func (t FlowScriptBuilder) StringArgument(value string) FlowScriptBuilder {
	return t.Argument(cadence.String(value))
}

// BooleanArgument add a Boolean Argument to the transaction
func (t FlowScriptBuilder) BooleanArgument(value bool) FlowScriptBuilder {
	return t.Argument(cadence.NewBool(value))
}

// BytesArgument add a Bytes Argument to the transaction
func (t FlowScriptBuilder) BytesArgument(value []byte) FlowScriptBuilder {
	return t.Argument(cadence.NewBytes(value))
}

// IntArgument add an Int Argument to the transaction
func (t FlowScriptBuilder) IntArgument(value int) FlowScriptBuilder {
	return t.Argument(cadence.NewInt(value))
}

// Int8Argument add an Int8 Argument to the transaction
func (t FlowScriptBuilder) Int8Argument(value int8) FlowScriptBuilder {
	return t.Argument(cadence.NewInt8(value))
}

// Int16Argument add an Int16 Argument to the transaction
func (t FlowScriptBuilder) Int16Argument(value int16) FlowScriptBuilder {
	return t.Argument(cadence.NewInt16(value))
}

// Int32Argument add an Int32 Argument to the transaction
func (t FlowScriptBuilder) Int32Argument(value int32) FlowScriptBuilder {
	return t.Argument(cadence.NewInt32(value))
}

// Int64Argument add an Int64 Argument to the transaction
func (t FlowScriptBuilder) Int64Argument(value int64) FlowScriptBuilder {
	return t.Argument(cadence.NewInt64(value))
}

// Int128Argument add an Int128 Argument to the transaction
func (t FlowScriptBuilder) Int128Argument(value int) FlowScriptBuilder {
	return t.Argument(cadence.NewInt128(value))
}

// Int256Argument add an Int256 Argument to the transaction
func (t FlowScriptBuilder) Int256Argument(value int) FlowScriptBuilder {
	return t.Argument(cadence.NewInt256(value))
}

// UIntArgument add an UInt Argument to the transaction
func (t FlowScriptBuilder) UIntArgument(value uint) FlowScriptBuilder {
	return t.Argument(cadence.NewUInt(value))
}

// UInt8Argument add an UInt8 Argument to the transaction
func (t FlowScriptBuilder) UInt8Argument(value uint8) FlowScriptBuilder {
	return t.Argument(cadence.NewUInt8(value))
}

// UInt16Argument add an UInt16 Argument to the transaction
func (t FlowScriptBuilder) UInt16Argument(value uint16) FlowScriptBuilder {
	return t.Argument(cadence.NewUInt16(value))
}

// UInt32Argument add an UInt32 Argument to the transaction
func (t FlowScriptBuilder) UInt32Argument(value uint32) FlowScriptBuilder {
	return t.Argument(cadence.NewUInt32(value))
}

// UInt64Argument add an UInt64 Argument to the transaction
func (t FlowScriptBuilder) UInt64Argument(value uint64) FlowScriptBuilder {
	return t.Argument(cadence.NewUInt64(value))
}

// UInt128Argument add an UInt128 Argument to the transaction
func (t FlowScriptBuilder) UInt128Argument(value uint) FlowScriptBuilder {
	return t.Argument(cadence.NewUInt128(value))
}

// UInt256Argument add an UInt256 Argument to the transaction
func (t FlowScriptBuilder) UInt256Argument(value uint) FlowScriptBuilder {
	return t.Argument(cadence.NewUInt256(value))
}

// Word8Argument add a Word8 Argument to the transaction
func (t FlowScriptBuilder) Word8Argument(value uint8) FlowScriptBuilder {
	return t.Argument(cadence.NewWord8(value))
}

// Word16Argument add a Word16 Argument to the transaction
func (t FlowScriptBuilder) Word16Argument(value uint16) FlowScriptBuilder {
	return t.Argument(cadence.NewWord16(value))
}

// Word32Argument add a Word32 Argument to the transaction
func (t FlowScriptBuilder) Word32Argument(value uint32) FlowScriptBuilder {
	return t.Argument(cadence.NewWord32(value))
}

// Word64Argument add a Word64 Argument to the transaction
func (t FlowScriptBuilder) Word64Argument(value uint64) FlowScriptBuilder {
	return t.Argument(cadence.NewWord64(value))
}

// Fix64Argument add a Fix64 Argument to the transaction
func (t FlowScriptBuilder) Fix64Argument(value string) FlowScriptBuilder {
	amount, err := cadence.NewFix64(value)
	if err != nil {
		panic(err)
	}
	return t.Argument(amount)
}

// UFix64Argument add a UFix64 Argument to the transaction
func (t FlowScriptBuilder) UFix64Argument(value string) FlowScriptBuilder {
	amount, err := cadence.NewUFix64(value)
	if err != nil {
		panic(err)
	}
	return t.Argument(amount)
}

// Run executes a read only script
func (t FlowScriptBuilder) Run(ctx context.Context) {
	result := t.RunFailOnError(ctx)
	log.Printf("Script run from result: %v\n", CadenceValueToJSONString(result))
}

// RunReturns executes a read only script
func (t FlowScriptBuilder) RunReturns(ctx context.Context) (cadence.Value, error) {

	f := t.Connector
	scriptFilePath := fmt.Sprintf("./scripts/%s.cdc", t.FileName)

	var err error
	script := []byte(t.ScriptAsString)
	if t.ScriptAsString == "" {
		script, err = f.State.ReaderWriter().ReadFile(scriptFilePath)
		if err != nil {
			return nil, err
		}
	}

	result, err := f.Services.ExecuteScript(
		ctx,
		flowkit.Script{
			Code:     script,
			Args:     t.Arguments,
			Location: scriptFilePath,
		},
		flowkit.LatestScriptQuery,
	)
	if err != nil {
		return nil, err
	}

	if t.ScriptAsString == "" {
		f.Logger.Debug(fmt.Sprintf("Script run from path %s\n", scriptFilePath))
	}
	return result, nil
}

func (t FlowScriptBuilder) RunFailOnError(ctx context.Context) cadence.Value {
	result, err := t.RunReturns(ctx)
	if err != nil {
		t.Connector.Logger.Error(fmt.Sprintf("Error executing script: %s output %v", t.FileName, err))
		os.Exit(1)
	}
	return result
}

// RunReturnsJSONString runs the script and returns pretty printed json string
func (t FlowScriptBuilder) RunReturnsJSONString(ctx context.Context) string {
	return CadenceValueToJSONString(t.RunFailOnError(ctx))
}

// RunReturnsInterface runs the script and returns interface{}
func (t FlowScriptBuilder) RunReturnsInterface(ctx context.Context) interface{} {
	return CadenceValueToInterface(t.RunFailOnError(ctx))
}
