// Code generated — DO NOT EDIT.

package house_rwa_receiver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
	"google.golang.org/protobuf/types/known/emptypb"

	pb2 "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values/pb"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm/bindings"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

var (
	_ = bytes.Equal
	_ = errors.New
	_ = fmt.Sprintf
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
	_ = emptypb.Empty{}
	_ = pb.NewBigIntFromInt
	_ = pb2.AggregationType_AGGREGATION_TYPE_COMMON_PREFIX
	_ = bindings.FilterOptions{}
	_ = evm.FilterLogTriggerRequest{}
	_ = cre.ResponseBufferTooSmall
	_ = rpc.API{}
	_ = json.Unmarshal
	_ = reflect.Bool
)

var HouseRWAReceiverMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_houseRWA\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_forwarder\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"selectors\",\"type\":\"bytes4[]\",\"internalType\":\"bytes4[]\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowedSelectors\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"forwarder\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"houseRWA\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"onReport\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"report\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setAllowedSelector\",\"inputs\":[{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"},{\"name\":\"allowed\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setForwarder\",\"inputs\":[{\"name\":\"newForwarder\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"ForwarderUpdated\",\"inputs\":[{\"name\":\"oldForwarder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newForwarder\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ReportForwarded\",\"inputs\":[{\"name\":\"houseRWA\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"selector\",\"type\":\"bytes4\",\"indexed\":true,\"internalType\":\"bytes4\"},{\"name\":\"reportHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SelectorUpdated\",\"inputs\":[{\"name\":\"selector\",\"type\":\"bytes4\",\"indexed\":true,\"internalType\":\"bytes4\"},{\"name\":\"allowed\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"HouseRWAReceiver_InvalidAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"HouseRWAReceiver_InvalidReport\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"HouseRWAReceiver_OnlyForwarder\",\"inputs\":[{\"name\":\"caller\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"HouseRWAReceiver_SelectorNotAllowed\",\"inputs\":[{\"name\":\"selector\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// Structs

// Contract Method Inputs
type AllowedSelectorsInput struct {
	Arg0 [4]byte
}

type OnReportInput struct {
	Arg0   []byte
	Report []byte
}

type SetAllowedSelectorInput struct {
	Selector [4]byte
	Allowed  bool
}

type SetForwarderInput struct {
	NewForwarder common.Address
}

type SupportsInterfaceInput struct {
	InterfaceId [4]byte
}

type TransferOwnershipInput struct {
	NewOwner common.Address
}

// Contract Method Outputs

// Errors
type HouseRWAReceiverInvalidAddress struct {
}

type HouseRWAReceiverInvalidReport struct {
}

type HouseRWAReceiverOnlyForwarder struct {
	Caller common.Address
}

type HouseRWAReceiverSelectorNotAllowed struct {
	Selector [4]byte
}

type OwnableInvalidOwner struct {
	Owner common.Address
}

type OwnableUnauthorizedAccount struct {
	Account common.Address
}

// Events
// The <Event>Topics struct should be used as a filter (for log triggers).
// Note: It is only possible to filter on indexed fields.
// Indexed (string and bytes) fields will be of type common.Hash.
// They need to he (crypto.Keccak256) hashed and passed in.
// Indexed (tuple/slice/array) fields can be passed in as is, the Encode<Event>Topics function will handle the hashing.
//
// The <Event>Decoded struct will be the result of calling decode (Adapt) on the log trigger result.
// Indexed dynamic type fields will be of type common.Hash.

type ForwarderUpdatedTopics struct {
	OldForwarder common.Address
	NewForwarder common.Address
}

type ForwarderUpdatedDecoded struct {
	OldForwarder common.Address
	NewForwarder common.Address
}

type OwnershipTransferStartedTopics struct {
	PreviousOwner common.Address
	NewOwner      common.Address
}

type OwnershipTransferStartedDecoded struct {
	PreviousOwner common.Address
	NewOwner      common.Address
}

type OwnershipTransferredTopics struct {
	PreviousOwner common.Address
	NewOwner      common.Address
}

type OwnershipTransferredDecoded struct {
	PreviousOwner common.Address
	NewOwner      common.Address
}

type ReportForwardedTopics struct {
	HouseRWA   common.Address
	Selector   [4]byte
	ReportHash [32]byte
}

type ReportForwardedDecoded struct {
	HouseRWA   common.Address
	Selector   [4]byte
	ReportHash [32]byte
}

type SelectorUpdatedTopics struct {
	Selector [4]byte
}

type SelectorUpdatedDecoded struct {
	Selector [4]byte
	Allowed  bool
}

// Main Binding Type for HouseRWAReceiver
type HouseRWAReceiver struct {
	Address common.Address
	Options *bindings.ContractInitOptions
	ABI     *abi.ABI
	client  *evm.Client
	Codec   HouseRWAReceiverCodec
}

type HouseRWAReceiverCodec interface {
	EncodeAcceptOwnershipMethodCall() ([]byte, error)
	EncodeAllowedSelectorsMethodCall(in AllowedSelectorsInput) ([]byte, error)
	DecodeAllowedSelectorsMethodOutput(data []byte) (bool, error)
	EncodeForwarderMethodCall() ([]byte, error)
	DecodeForwarderMethodOutput(data []byte) (common.Address, error)
	EncodeHouseRWAMethodCall() ([]byte, error)
	DecodeHouseRWAMethodOutput(data []byte) (common.Address, error)
	EncodeOnReportMethodCall(in OnReportInput) ([]byte, error)
	EncodeOwnerMethodCall() ([]byte, error)
	DecodeOwnerMethodOutput(data []byte) (common.Address, error)
	EncodePendingOwnerMethodCall() ([]byte, error)
	DecodePendingOwnerMethodOutput(data []byte) (common.Address, error)
	EncodeRenounceOwnershipMethodCall() ([]byte, error)
	EncodeSetAllowedSelectorMethodCall(in SetAllowedSelectorInput) ([]byte, error)
	EncodeSetForwarderMethodCall(in SetForwarderInput) ([]byte, error)
	EncodeSupportsInterfaceMethodCall(in SupportsInterfaceInput) ([]byte, error)
	DecodeSupportsInterfaceMethodOutput(data []byte) (bool, error)
	EncodeTransferOwnershipMethodCall(in TransferOwnershipInput) ([]byte, error)
	ForwarderUpdatedLogHash() []byte
	EncodeForwarderUpdatedTopics(evt abi.Event, values []ForwarderUpdatedTopics) ([]*evm.TopicValues, error)
	DecodeForwarderUpdated(log *evm.Log) (*ForwarderUpdatedDecoded, error)
	OwnershipTransferStartedLogHash() []byte
	EncodeOwnershipTransferStartedTopics(evt abi.Event, values []OwnershipTransferStartedTopics) ([]*evm.TopicValues, error)
	DecodeOwnershipTransferStarted(log *evm.Log) (*OwnershipTransferStartedDecoded, error)
	OwnershipTransferredLogHash() []byte
	EncodeOwnershipTransferredTopics(evt abi.Event, values []OwnershipTransferredTopics) ([]*evm.TopicValues, error)
	DecodeOwnershipTransferred(log *evm.Log) (*OwnershipTransferredDecoded, error)
	ReportForwardedLogHash() []byte
	EncodeReportForwardedTopics(evt abi.Event, values []ReportForwardedTopics) ([]*evm.TopicValues, error)
	DecodeReportForwarded(log *evm.Log) (*ReportForwardedDecoded, error)
	SelectorUpdatedLogHash() []byte
	EncodeSelectorUpdatedTopics(evt abi.Event, values []SelectorUpdatedTopics) ([]*evm.TopicValues, error)
	DecodeSelectorUpdated(log *evm.Log) (*SelectorUpdatedDecoded, error)
}

func NewHouseRWAReceiver(
	client *evm.Client,
	address common.Address,
	options *bindings.ContractInitOptions,
) (*HouseRWAReceiver, error) {
	parsed, err := abi.JSON(strings.NewReader(HouseRWAReceiverMetaData.ABI))
	if err != nil {
		return nil, err
	}
	codec, err := NewCodec()
	if err != nil {
		return nil, err
	}
	return &HouseRWAReceiver{
		Address: address,
		Options: options,
		ABI:     &parsed,
		client:  client,
		Codec:   codec,
	}, nil
}

type Codec struct {
	abi *abi.ABI
}

func NewCodec() (HouseRWAReceiverCodec, error) {
	parsed, err := abi.JSON(strings.NewReader(HouseRWAReceiverMetaData.ABI))
	if err != nil {
		return nil, err
	}
	return &Codec{abi: &parsed}, nil
}

func (c *Codec) EncodeAcceptOwnershipMethodCall() ([]byte, error) {
	return c.abi.Pack("acceptOwnership")
}

func (c *Codec) EncodeAllowedSelectorsMethodCall(in AllowedSelectorsInput) ([]byte, error) {
	return c.abi.Pack("allowedSelectors", in.Arg0)
}

func (c *Codec) DecodeAllowedSelectorsMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["allowedSelectors"].Outputs.Unpack(data)
	if err != nil {
		return *new(bool), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(bool), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result bool
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(bool), fmt.Errorf("failed to unmarshal to bool: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeForwarderMethodCall() ([]byte, error) {
	return c.abi.Pack("forwarder")
}

func (c *Codec) DecodeForwarderMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["forwarder"].Outputs.Unpack(data)
	if err != nil {
		return *new(common.Address), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(common.Address), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result common.Address
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(common.Address), fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeHouseRWAMethodCall() ([]byte, error) {
	return c.abi.Pack("houseRWA")
}

func (c *Codec) DecodeHouseRWAMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["houseRWA"].Outputs.Unpack(data)
	if err != nil {
		return *new(common.Address), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(common.Address), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result common.Address
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(common.Address), fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeOnReportMethodCall(in OnReportInput) ([]byte, error) {
	return c.abi.Pack("onReport", in.Arg0, in.Report)
}

func (c *Codec) EncodeOwnerMethodCall() ([]byte, error) {
	return c.abi.Pack("owner")
}

func (c *Codec) DecodeOwnerMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["owner"].Outputs.Unpack(data)
	if err != nil {
		return *new(common.Address), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(common.Address), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result common.Address
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(common.Address), fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodePendingOwnerMethodCall() ([]byte, error) {
	return c.abi.Pack("pendingOwner")
}

func (c *Codec) DecodePendingOwnerMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["pendingOwner"].Outputs.Unpack(data)
	if err != nil {
		return *new(common.Address), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(common.Address), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result common.Address
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(common.Address), fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeRenounceOwnershipMethodCall() ([]byte, error) {
	return c.abi.Pack("renounceOwnership")
}

func (c *Codec) EncodeSetAllowedSelectorMethodCall(in SetAllowedSelectorInput) ([]byte, error) {
	return c.abi.Pack("setAllowedSelector", in.Selector, in.Allowed)
}

func (c *Codec) EncodeSetForwarderMethodCall(in SetForwarderInput) ([]byte, error) {
	return c.abi.Pack("setForwarder", in.NewForwarder)
}

func (c *Codec) EncodeSupportsInterfaceMethodCall(in SupportsInterfaceInput) ([]byte, error) {
	return c.abi.Pack("supportsInterface", in.InterfaceId)
}

func (c *Codec) DecodeSupportsInterfaceMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["supportsInterface"].Outputs.Unpack(data)
	if err != nil {
		return *new(bool), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(bool), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result bool
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(bool), fmt.Errorf("failed to unmarshal to bool: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeTransferOwnershipMethodCall(in TransferOwnershipInput) ([]byte, error) {
	return c.abi.Pack("transferOwnership", in.NewOwner)
}

func (c *Codec) ForwarderUpdatedLogHash() []byte {
	return c.abi.Events["ForwarderUpdated"].ID.Bytes()
}

func (c *Codec) EncodeForwarderUpdatedTopics(
	evt abi.Event,
	values []ForwarderUpdatedTopics,
) ([]*evm.TopicValues, error) {
	var oldForwarderRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.OldForwarder).IsZero() {
			oldForwarderRule = append(oldForwarderRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.OldForwarder)
		if err != nil {
			return nil, err
		}
		oldForwarderRule = append(oldForwarderRule, fieldVal)
	}
	var newForwarderRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.NewForwarder).IsZero() {
			newForwarderRule = append(newForwarderRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.NewForwarder)
		if err != nil {
			return nil, err
		}
		newForwarderRule = append(newForwarderRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		oldForwarderRule,
		newForwarderRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeForwarderUpdated decodes a log into a ForwarderUpdated struct.
func (c *Codec) DecodeForwarderUpdated(log *evm.Log) (*ForwarderUpdatedDecoded, error) {
	event := new(ForwarderUpdatedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ForwarderUpdated", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ForwarderUpdated"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) OwnershipTransferStartedLogHash() []byte {
	return c.abi.Events["OwnershipTransferStarted"].ID.Bytes()
}

func (c *Codec) EncodeOwnershipTransferStartedTopics(
	evt abi.Event,
	values []OwnershipTransferStartedTopics,
) ([]*evm.TopicValues, error) {
	var previousOwnerRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.PreviousOwner).IsZero() {
			previousOwnerRule = append(previousOwnerRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.PreviousOwner)
		if err != nil {
			return nil, err
		}
		previousOwnerRule = append(previousOwnerRule, fieldVal)
	}
	var newOwnerRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.NewOwner).IsZero() {
			newOwnerRule = append(newOwnerRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.NewOwner)
		if err != nil {
			return nil, err
		}
		newOwnerRule = append(newOwnerRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		previousOwnerRule,
		newOwnerRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeOwnershipTransferStarted decodes a log into a OwnershipTransferStarted struct.
func (c *Codec) DecodeOwnershipTransferStarted(log *evm.Log) (*OwnershipTransferStartedDecoded, error) {
	event := new(OwnershipTransferStartedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "OwnershipTransferStarted", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["OwnershipTransferStarted"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) OwnershipTransferredLogHash() []byte {
	return c.abi.Events["OwnershipTransferred"].ID.Bytes()
}

func (c *Codec) EncodeOwnershipTransferredTopics(
	evt abi.Event,
	values []OwnershipTransferredTopics,
) ([]*evm.TopicValues, error) {
	var previousOwnerRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.PreviousOwner).IsZero() {
			previousOwnerRule = append(previousOwnerRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.PreviousOwner)
		if err != nil {
			return nil, err
		}
		previousOwnerRule = append(previousOwnerRule, fieldVal)
	}
	var newOwnerRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.NewOwner).IsZero() {
			newOwnerRule = append(newOwnerRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.NewOwner)
		if err != nil {
			return nil, err
		}
		newOwnerRule = append(newOwnerRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		previousOwnerRule,
		newOwnerRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeOwnershipTransferred decodes a log into a OwnershipTransferred struct.
func (c *Codec) DecodeOwnershipTransferred(log *evm.Log) (*OwnershipTransferredDecoded, error) {
	event := new(OwnershipTransferredDecoded)
	if err := c.abi.UnpackIntoInterface(event, "OwnershipTransferred", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["OwnershipTransferred"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) ReportForwardedLogHash() []byte {
	return c.abi.Events["ReportForwarded"].ID.Bytes()
}

func (c *Codec) EncodeReportForwardedTopics(
	evt abi.Event,
	values []ReportForwardedTopics,
) ([]*evm.TopicValues, error) {
	var houseRWARule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.HouseRWA).IsZero() {
			houseRWARule = append(houseRWARule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.HouseRWA)
		if err != nil {
			return nil, err
		}
		houseRWARule = append(houseRWARule, fieldVal)
	}
	var selectorRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Selector).IsZero() {
			selectorRule = append(selectorRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Selector)
		if err != nil {
			return nil, err
		}
		selectorRule = append(selectorRule, fieldVal)
	}
	var reportHashRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.ReportHash).IsZero() {
			reportHashRule = append(reportHashRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[2], v.ReportHash)
		if err != nil {
			return nil, err
		}
		reportHashRule = append(reportHashRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		houseRWARule,
		selectorRule,
		reportHashRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeReportForwarded decodes a log into a ReportForwarded struct.
func (c *Codec) DecodeReportForwarded(log *evm.Log) (*ReportForwardedDecoded, error) {
	event := new(ReportForwardedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ReportForwarded", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ReportForwarded"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c *Codec) SelectorUpdatedLogHash() []byte {
	return c.abi.Events["SelectorUpdated"].ID.Bytes()
}

func (c *Codec) EncodeSelectorUpdatedTopics(
	evt abi.Event,
	values []SelectorUpdatedTopics,
) ([]*evm.TopicValues, error) {
	var selectorRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Selector).IsZero() {
			selectorRule = append(selectorRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Selector)
		if err != nil {
			return nil, err
		}
		selectorRule = append(selectorRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		selectorRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeSelectorUpdated decodes a log into a SelectorUpdated struct.
func (c *Codec) DecodeSelectorUpdated(log *evm.Log) (*SelectorUpdatedDecoded, error) {
	event := new(SelectorUpdatedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "SelectorUpdated", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["SelectorUpdated"].Inputs {
		if arg.Indexed {
			if arg.Type.T == abi.TupleTy {
				// abigen throws on tuple, so converting to bytes to
				// receive back the common.Hash as is instead of error
				arg.Type.T = abi.BytesTy
			}
			indexed = append(indexed, arg)
		}
	}
	// Convert [][]byte → []common.Hash
	topics := make([]common.Hash, len(log.Topics))
	for i, t := range log.Topics {
		topics[i] = common.BytesToHash(t)
	}

	if err := abi.ParseTopics(event, indexed, topics[1:]); err != nil {
		return nil, err
	}
	return event, nil
}

func (c HouseRWAReceiver) AllowedSelectors(
	runtime cre.Runtime,
	args AllowedSelectorsInput,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeAllowedSelectorsMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[bool](*new(bool), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (bool, error) {
		return c.Codec.DecodeAllowedSelectorsMethodOutput(response.Data)
	})

}

func (c HouseRWAReceiver) Forwarder(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeForwarderMethodCall()
	if err != nil {
		return cre.PromiseFromResult[common.Address](*new(common.Address), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (common.Address, error) {
		return c.Codec.DecodeForwarderMethodOutput(response.Data)
	})

}

func (c HouseRWAReceiver) HouseRWA(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeHouseRWAMethodCall()
	if err != nil {
		return cre.PromiseFromResult[common.Address](*new(common.Address), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (common.Address, error) {
		return c.Codec.DecodeHouseRWAMethodOutput(response.Data)
	})

}

func (c HouseRWAReceiver) Owner(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeOwnerMethodCall()
	if err != nil {
		return cre.PromiseFromResult[common.Address](*new(common.Address), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (common.Address, error) {
		return c.Codec.DecodeOwnerMethodOutput(response.Data)
	})

}

func (c HouseRWAReceiver) PendingOwner(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodePendingOwnerMethodCall()
	if err != nil {
		return cre.PromiseFromResult[common.Address](*new(common.Address), err)
	}

	var bn cre.Promise[*pb.BigInt]
	if blockNumber == nil {
		promise := c.client.HeaderByNumber(runtime, &evm.HeaderByNumberRequest{
			BlockNumber: bindings.FinalizedBlockNumber,
		})

		bn = cre.Then(promise, func(finalizedBlock *evm.HeaderByNumberReply) (*pb.BigInt, error) {
			if finalizedBlock == nil || finalizedBlock.Header == nil {
				return nil, errors.New("failed to get finalized block header")
			}
			return finalizedBlock.Header.BlockNumber, nil
		})
	} else {
		bn = cre.PromiseFromResult(pb.NewBigIntFromInt(blockNumber), nil)
	}

	promise := cre.ThenPromise(bn, func(bn *pb.BigInt) cre.Promise[*evm.CallContractReply] {
		return c.client.CallContract(runtime, &evm.CallContractRequest{
			Call:        &evm.CallMsg{To: c.Address.Bytes(), Data: calldata},
			BlockNumber: bn,
		})
	})
	return cre.Then(promise, func(response *evm.CallContractReply) (common.Address, error) {
		return c.Codec.DecodePendingOwnerMethodOutput(response.Data)
	})

}

func (c HouseRWAReceiver) WriteReport(
	runtime cre.Runtime,
	report *cre.Report,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	return c.client.WriteReport(runtime, &evm.WriteCreReportRequest{
		Receiver:  c.Address.Bytes(),
		Report:    report,
		GasConfig: gasConfig,
	})
}

// DecodeHouseRWAReceiverInvalidAddressError decodes a HouseRWAReceiver_InvalidAddress error from revert data.
func (c *HouseRWAReceiver) DecodeHouseRWAReceiverInvalidAddressError(data []byte) (*HouseRWAReceiverInvalidAddress, error) {
	args := c.ABI.Errors["HouseRWAReceiver_InvalidAddress"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &HouseRWAReceiverInvalidAddress{}, nil
}

// Error implements the error interface for HouseRWAReceiverInvalidAddress.
func (e *HouseRWAReceiverInvalidAddress) Error() string {
	return fmt.Sprintf("HouseRWAReceiverInvalidAddress error:")
}

// DecodeHouseRWAReceiverInvalidReportError decodes a HouseRWAReceiver_InvalidReport error from revert data.
func (c *HouseRWAReceiver) DecodeHouseRWAReceiverInvalidReportError(data []byte) (*HouseRWAReceiverInvalidReport, error) {
	args := c.ABI.Errors["HouseRWAReceiver_InvalidReport"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &HouseRWAReceiverInvalidReport{}, nil
}

// Error implements the error interface for HouseRWAReceiverInvalidReport.
func (e *HouseRWAReceiverInvalidReport) Error() string {
	return fmt.Sprintf("HouseRWAReceiverInvalidReport error:")
}

// DecodeHouseRWAReceiverOnlyForwarderError decodes a HouseRWAReceiver_OnlyForwarder error from revert data.
func (c *HouseRWAReceiver) DecodeHouseRWAReceiverOnlyForwarderError(data []byte) (*HouseRWAReceiverOnlyForwarder, error) {
	args := c.ABI.Errors["HouseRWAReceiver_OnlyForwarder"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	caller, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for caller in HouseRWAReceiverOnlyForwarder error")
	}

	return &HouseRWAReceiverOnlyForwarder{
		Caller: caller,
	}, nil
}

// Error implements the error interface for HouseRWAReceiverOnlyForwarder.
func (e *HouseRWAReceiverOnlyForwarder) Error() string {
	return fmt.Sprintf("HouseRWAReceiverOnlyForwarder error: caller=%v;", e.Caller)
}

// DecodeHouseRWAReceiverSelectorNotAllowedError decodes a HouseRWAReceiver_SelectorNotAllowed error from revert data.
func (c *HouseRWAReceiver) DecodeHouseRWAReceiverSelectorNotAllowedError(data []byte) (*HouseRWAReceiverSelectorNotAllowed, error) {
	args := c.ABI.Errors["HouseRWAReceiver_SelectorNotAllowed"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	selector, ok0 := values[0].([4]byte)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for selector in HouseRWAReceiverSelectorNotAllowed error")
	}

	return &HouseRWAReceiverSelectorNotAllowed{
		Selector: selector,
	}, nil
}

// Error implements the error interface for HouseRWAReceiverSelectorNotAllowed.
func (e *HouseRWAReceiverSelectorNotAllowed) Error() string {
	return fmt.Sprintf("HouseRWAReceiverSelectorNotAllowed error: selector=%v;", e.Selector)
}

// DecodeOwnableInvalidOwnerError decodes a OwnableInvalidOwner error from revert data.
func (c *HouseRWAReceiver) DecodeOwnableInvalidOwnerError(data []byte) (*OwnableInvalidOwner, error) {
	args := c.ABI.Errors["OwnableInvalidOwner"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	owner, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for owner in OwnableInvalidOwner error")
	}

	return &OwnableInvalidOwner{
		Owner: owner,
	}, nil
}

// Error implements the error interface for OwnableInvalidOwner.
func (e *OwnableInvalidOwner) Error() string {
	return fmt.Sprintf("OwnableInvalidOwner error: owner=%v;", e.Owner)
}

// DecodeOwnableUnauthorizedAccountError decodes a OwnableUnauthorizedAccount error from revert data.
func (c *HouseRWAReceiver) DecodeOwnableUnauthorizedAccountError(data []byte) (*OwnableUnauthorizedAccount, error) {
	args := c.ABI.Errors["OwnableUnauthorizedAccount"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	account, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for account in OwnableUnauthorizedAccount error")
	}

	return &OwnableUnauthorizedAccount{
		Account: account,
	}, nil
}

// Error implements the error interface for OwnableUnauthorizedAccount.
func (e *OwnableUnauthorizedAccount) Error() string {
	return fmt.Sprintf("OwnableUnauthorizedAccount error: account=%v;", e.Account)
}

func (c *HouseRWAReceiver) UnpackError(data []byte) (any, error) {
	switch common.Bytes2Hex(data[:4]) {
	case common.Bytes2Hex(c.ABI.Errors["HouseRWAReceiver_InvalidAddress"].ID.Bytes()[:4]):
		return c.DecodeHouseRWAReceiverInvalidAddressError(data)
	case common.Bytes2Hex(c.ABI.Errors["HouseRWAReceiver_InvalidReport"].ID.Bytes()[:4]):
		return c.DecodeHouseRWAReceiverInvalidReportError(data)
	case common.Bytes2Hex(c.ABI.Errors["HouseRWAReceiver_OnlyForwarder"].ID.Bytes()[:4]):
		return c.DecodeHouseRWAReceiverOnlyForwarderError(data)
	case common.Bytes2Hex(c.ABI.Errors["HouseRWAReceiver_SelectorNotAllowed"].ID.Bytes()[:4]):
		return c.DecodeHouseRWAReceiverSelectorNotAllowedError(data)
	case common.Bytes2Hex(c.ABI.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]):
		return c.DecodeOwnableInvalidOwnerError(data)
	case common.Bytes2Hex(c.ABI.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]):
		return c.DecodeOwnableUnauthorizedAccountError(data)
	default:
		return nil, errors.New("unknown error selector")
	}
}

// ForwarderUpdatedTrigger wraps the raw log trigger and provides decoded ForwarderUpdatedDecoded data
type ForwarderUpdatedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                   // Embed the raw trigger
	contract                        *HouseRWAReceiver // Keep reference for decoding
}

// Adapt method that decodes the log into ForwarderUpdated data
func (t *ForwarderUpdatedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ForwarderUpdatedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeForwarderUpdated(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ForwarderUpdated log: %w", err)
	}

	return &bindings.DecodedLog[ForwarderUpdatedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWAReceiver) LogTriggerForwarderUpdatedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ForwarderUpdatedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ForwarderUpdatedDecoded]], error) {
	event := c.ABI.Events["ForwarderUpdated"]
	topics, err := c.Codec.EncodeForwarderUpdatedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ForwarderUpdated: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ForwarderUpdatedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWAReceiver) FilterLogsForwarderUpdated(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ForwarderUpdatedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// OwnershipTransferStartedTrigger wraps the raw log trigger and provides decoded OwnershipTransferStartedDecoded data
type OwnershipTransferStartedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                   // Embed the raw trigger
	contract                        *HouseRWAReceiver // Keep reference for decoding
}

// Adapt method that decodes the log into OwnershipTransferStarted data
func (t *OwnershipTransferStartedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[OwnershipTransferStartedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeOwnershipTransferStarted(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode OwnershipTransferStarted log: %w", err)
	}

	return &bindings.DecodedLog[OwnershipTransferStartedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWAReceiver) LogTriggerOwnershipTransferStartedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []OwnershipTransferStartedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[OwnershipTransferStartedDecoded]], error) {
	event := c.ABI.Events["OwnershipTransferStarted"]
	topics, err := c.Codec.EncodeOwnershipTransferStartedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for OwnershipTransferStarted: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &OwnershipTransferStartedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWAReceiver) FilterLogsOwnershipTransferStarted(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.OwnershipTransferStartedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// OwnershipTransferredTrigger wraps the raw log trigger and provides decoded OwnershipTransferredDecoded data
type OwnershipTransferredTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                   // Embed the raw trigger
	contract                        *HouseRWAReceiver // Keep reference for decoding
}

// Adapt method that decodes the log into OwnershipTransferred data
func (t *OwnershipTransferredTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[OwnershipTransferredDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeOwnershipTransferred(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode OwnershipTransferred log: %w", err)
	}

	return &bindings.DecodedLog[OwnershipTransferredDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWAReceiver) LogTriggerOwnershipTransferredLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []OwnershipTransferredTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[OwnershipTransferredDecoded]], error) {
	event := c.ABI.Events["OwnershipTransferred"]
	topics, err := c.Codec.EncodeOwnershipTransferredTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for OwnershipTransferred: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &OwnershipTransferredTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWAReceiver) FilterLogsOwnershipTransferred(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.OwnershipTransferredLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// ReportForwardedTrigger wraps the raw log trigger and provides decoded ReportForwardedDecoded data
type ReportForwardedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                   // Embed the raw trigger
	contract                        *HouseRWAReceiver // Keep reference for decoding
}

// Adapt method that decodes the log into ReportForwarded data
func (t *ReportForwardedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ReportForwardedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeReportForwarded(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ReportForwarded log: %w", err)
	}

	return &bindings.DecodedLog[ReportForwardedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWAReceiver) LogTriggerReportForwardedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ReportForwardedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ReportForwardedDecoded]], error) {
	event := c.ABI.Events["ReportForwarded"]
	topics, err := c.Codec.EncodeReportForwardedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ReportForwarded: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ReportForwardedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWAReceiver) FilterLogsReportForwarded(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ReportForwardedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// SelectorUpdatedTrigger wraps the raw log trigger and provides decoded SelectorUpdatedDecoded data
type SelectorUpdatedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]                   // Embed the raw trigger
	contract                        *HouseRWAReceiver // Keep reference for decoding
}

// Adapt method that decodes the log into SelectorUpdated data
func (t *SelectorUpdatedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[SelectorUpdatedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeSelectorUpdated(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode SelectorUpdated log: %w", err)
	}

	return &bindings.DecodedLog[SelectorUpdatedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWAReceiver) LogTriggerSelectorUpdatedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []SelectorUpdatedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[SelectorUpdatedDecoded]], error) {
	event := c.ABI.Events["SelectorUpdated"]
	topics, err := c.Codec.EncodeSelectorUpdatedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for SelectorUpdated: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &SelectorUpdatedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWAReceiver) FilterLogsSelectorUpdated(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.SelectorUpdatedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}
