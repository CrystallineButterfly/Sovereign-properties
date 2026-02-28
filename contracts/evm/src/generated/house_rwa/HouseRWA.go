// Code generated — DO NOT EDIT.

package house_rwa

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

var HouseRWAMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"BPS_DENOMINATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"KEY_EXCHANGE_EXPIRY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_BILLS_PER_HOUSE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_DOCUMENTS_PER_HOUSE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_RENTAL_DURATION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_VALIDATOR_STAKE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"PROTOCOL_FEE_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"arbitrators\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"authorizeCREWorkflow\",\"inputs\":[{\"name\":\"workflow\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authorizedCREWorkflows\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancelListing\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimKey\",\"inputs\":[{\"name\":\"keyHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"completeSale\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"buyer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"keyHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"encryptedKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createBill\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"billType\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"amount\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"dueDate\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isRecurring\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"recurrenceInterval\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[{\"name\":\"billIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createListing\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"listingType\",\"type\":\"uint8\",\"internalType\":\"enumHouseRWA.ListingType\"},{\"name\":\"price\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"preferredToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isPrivateSale\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"allowedBuyer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"durationDays\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"depositForRental\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"disputeBill\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"billIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emergencyPause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"emergencyUnpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"endRental\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"feeRecipient\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getActiveRental\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structHouseRWA.RentalAgreement\",\"components\":[{\"name\":\"renter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"startTime\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"endTime\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"depositAmount\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"monthlyRent\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"isActive\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"encryptedAccessKeyHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"disputeStatus\",\"type\":\"uint8\",\"internalType\":\"enumHouseRWA.DisputeStatus\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getApproved\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getBills\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structHouseRWA.Bill[]\",\"components\":[{\"name\":\"billType\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"amount\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"dueDate\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"paidAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumHouseRWA.BillStatus\"},{\"name\":\"paymentReference\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"isRecurring\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recurrenceInterval\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getHouseDetails\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structHouseRWA.House\",\"components\":[{\"name\":\"houseId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"documentHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"documentURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"storageType\",\"type\":\"uint8\",\"internalType\":\"enumHouseRWA.DocumentStorageType\"},{\"name\":\"originalOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"mintedAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"isVerified\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"documentCount\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getListing\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structHouseRWA.Listing\",\"components\":[{\"name\":\"listingType\",\"type\":\"uint8\",\"internalType\":\"enumHouseRWA.ListingType\"},{\"name\":\"price\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"preferredToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isPrivateSale\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"allowedBuyer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"createdAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"expiresAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"platformFee\",\"type\":\"uint8\",\"internalType\":\"uint8\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalBillsCount\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasKYC\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"highValueThresholdUSD\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"houseBills\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"billType\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"amount\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"dueDate\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"paidAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumHouseRWA.BillStatus\"},{\"name\":\"paymentReference\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"isRecurring\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"recurrenceInterval\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"houses\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"houseId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"documentHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"documentURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"storageType\",\"type\":\"uint8\",\"internalType\":\"enumHouseRWA.DocumentStorageType\"},{\"name\":\"originalOwner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"mintedAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"isVerified\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"documentCount\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_feeRecipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_initialCREWorkflow\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isApprovedForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isRented\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"keyExchanges\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"keyHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"encryptedKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"intendedRecipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"createdAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"expiresAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"isClaimed\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"exchangeType\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"kycInfo\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"isVerified\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"verifiedAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"verificationHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"verificationLevel\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"expiryDate\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"listings\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"listingType\",\"type\":\"uint8\",\"internalType\":\"enumHouseRWA.ListingType\"},{\"name\":\"price\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"preferredToken\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isPrivateSale\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"allowedBuyer\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"createdAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"expiresAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"platformFee\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minKYCLevelForHighValue\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"minKYCLevelForMint\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"mint\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"houseId\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"documentHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"documentURI\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"storageType\",\"type\":\"uint8\",\"internalType\":\"enumHouseRWA.DocumentStorageType\"},{\"name\":\"verificationData\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"mintingPaused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"nextTokenId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"openRentalDispute\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"ownerOf\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"paymentsPaused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingRentalDeposits\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"priceFeed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recordBillPayment\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"billIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"paymentMethod\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"paymentReference\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rentals\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"renter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"startTime\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"endTime\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"depositAmount\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"monthlyRent\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"isActive\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"encryptedAccessKeyHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"disputeStatus\",\"type\":\"uint8\",\"internalType\":\"enumHouseRWA.DisputeStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rentalsPaused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"resolveRentalDispute\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"depositToOwner\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"depositToRenter\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeCREWorkflow\",\"inputs\":[{\"name\":\"workflow\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeKYC\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"safeTransferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"salesPaused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setApprovalForAll\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setArbitrator\",\"inputs\":[{\"name\":\"arbitrator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isArbitrator\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setFeeRecipient\",\"inputs\":[{\"name\":\"_feeRecipient\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setKYCRequirements\",\"inputs\":[{\"name\":\"_minLevelForMint\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"_minLevelForHighValue\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"_highValueThresholdUSD\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setKYCVerification\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"level\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"verificationHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"expiryDate\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMintingPaused\",\"inputs\":[{\"name\":\"paused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setPaymentsPaused\",\"inputs\":[{\"name\":\"paused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRentalsPaused\",\"inputs\":[{\"name\":\"paused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setSalesPaused\",\"inputs\":[{\"name\":\"paused\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setTrustedBillProvider\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"trusted\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slashValidator\",\"inputs\":[{\"name\":\"validator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"reason\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"stakeAsValidator\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"startRental\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"renter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"durationDays\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"depositAmount\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"monthlyRent\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"encryptedAccessKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenByIndex\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenOfOwnerByIndex\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"tokenURI\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalFeesCollected\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"trustedBillProviders\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"unstake\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"validators\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"stakedAmount\",\"type\":\"uint96\",\"internalType\":\"uint96\"},{\"name\":\"stakedAt\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"reputation\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"isSlashed\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"successfulValidations\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"failedValidations\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"withdrawRentalDeposit\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ApprovalForAll\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"operator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BillCreated\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"billIndex\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"billType\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"dueDate\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BillDisputed\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"billIndex\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"disputer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"reason\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"BillPaid\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"billIndex\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"paymentMethod\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"paymentReference\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CircuitBreakerTriggered\",\"inputs\":[{\"name\":\"component\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"reason\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EmergencyAction\",\"inputs\":[{\"name\":\"action\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"triggeredBy\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"HouseListed\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"listingType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumHouseRWA.ListingType\"},{\"name\":\"price\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"isPrivateSale\",\"type\":\"bool\",\"indexed\":false,\"internalType\":\"bool\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"HouseMinted\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"houseId\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"documentHash\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"storageType\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumHouseRWA.DocumentStorageType\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"HouseSold\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"seller\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"buyer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"price\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"protocolFee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"KYCVerified\",\"inputs\":[{\"name\":\"user\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"level\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"KeyClaimed\",\"inputs\":[{\"name\":\"keyHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"claimant\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"KeyExchangeCreated\",\"inputs\":[{\"name\":\"keyHash\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"expiry\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RentalDepositReceived\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"renter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RentalDepositWithdrawn\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"renter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RentalEnded\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"renter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"depositReturned\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RentalStarted\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"renter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"startTime\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"endTime\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"deposit\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorSlashed\",\"inputs\":[{\"name\":\"validator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"reason\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorStaked\",\"inputs\":[{\"name\":\"validator\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC721EnumerableForbiddenBatchMint\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC721IncorrectOwner\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InsufficientApproval\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidApprover\",\"inputs\":[{\"name\":\"approver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOperator\",\"inputs\":[{\"name\":\"operator\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidReceiver\",\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721InvalidSender\",\"inputs\":[{\"name\":\"sender\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC721NonexistentToken\",\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"ERC721OutOfBoundsIndex\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"EnforcedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExpectedPause\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedInnerCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
}

// Structs
type Bill struct {
	BillType           string
	Amount             *big.Int
	DueDate            *big.Int
	PaidAt             *big.Int
	Status             uint8
	PaymentReference   [32]byte
	IsRecurring        bool
	Provider           common.Address
	RecurrenceInterval uint8
}

type House struct {
	HouseId       string
	DocumentHash  [32]byte
	DocumentURI   string
	StorageType   uint8
	OriginalOwner common.Address
	MintedAt      *big.Int
	IsVerified    bool
	DocumentCount uint8
}

type Listing struct {
	ListingType    uint8
	Price          *big.Int
	PreferredToken common.Address
	IsPrivateSale  bool
	AllowedBuyer   common.Address
	CreatedAt      *big.Int
	ExpiresAt      *big.Int
	PlatformFee    uint8
}

type RentalAgreement struct {
	Renter                 common.Address
	StartTime              *big.Int
	EndTime                *big.Int
	DepositAmount          *big.Int
	MonthlyRent            *big.Int
	IsActive               bool
	EncryptedAccessKeyHash [32]byte
	DisputeStatus          uint8
}

// Contract Method Inputs
type ApproveInput struct {
	To      common.Address
	TokenId *big.Int
}

type ArbitratorsInput struct {
	Arg0 common.Address
}

type AuthorizeCREWorkflowInput struct {
	Workflow common.Address
}

type AuthorizedCREWorkflowsInput struct {
	Arg0 common.Address
}

type BalanceOfInput struct {
	Owner common.Address
}

type CancelListingInput struct {
	TokenId *big.Int
}

type ClaimKeyInput struct {
	KeyHash [32]byte
}

type CompleteSaleInput struct {
	TokenId      *big.Int
	Buyer        common.Address
	KeyHash      [32]byte
	EncryptedKey []byte
}

type CreateBillInput struct {
	TokenId            *big.Int
	BillType           string
	Amount             *big.Int
	DueDate            *big.Int
	Provider           common.Address
	IsRecurring        bool
	RecurrenceInterval uint8
}

type CreateListingInput struct {
	TokenId        *big.Int
	ListingType    uint8
	Price          *big.Int
	PreferredToken common.Address
	IsPrivateSale  bool
	AllowedBuyer   common.Address
	DurationDays   *big.Int
}

type DepositForRentalInput struct {
	TokenId *big.Int
}

type DisputeBillInput struct {
	TokenId   *big.Int
	BillIndex *big.Int
	Reason    string
}

type EndRentalInput struct {
	TokenId *big.Int
}

type GetActiveRentalInput struct {
	TokenId *big.Int
}

type GetApprovedInput struct {
	TokenId *big.Int
}

type GetBillsInput struct {
	TokenId *big.Int
}

type GetHouseDetailsInput struct {
	TokenId *big.Int
}

type GetListingInput struct {
	TokenId *big.Int
}

type GetTotalBillsCountInput struct {
	TokenId *big.Int
}

type HasKYCInput struct {
	User common.Address
}

type HouseBillsInput struct {
	Arg0 *big.Int
	Arg1 *big.Int
}

type HousesInput struct {
	Arg0 *big.Int
}

type InitializeInput struct {
	Owner              common.Address
	FeeRecipient       common.Address
	InitialCREWorkflow common.Address
}

type IsApprovedForAllInput struct {
	Owner    common.Address
	Operator common.Address
}

type IsRentedInput struct {
	TokenId *big.Int
}

type KeyExchangesInput struct {
	Arg0 [32]byte
}

type KycInfoInput struct {
	Arg0 common.Address
}

type ListingsInput struct {
	Arg0 *big.Int
}

type MintInput struct {
	To               common.Address
	HouseId          string
	DocumentHash     [32]byte
	DocumentURI      string
	StorageType      uint8
	VerificationData string
}

type OpenRentalDisputeInput struct {
	TokenId *big.Int
	Reason  string
}

type OwnerOfInput struct {
	TokenId *big.Int
}

type PendingRentalDepositsInput struct {
	Arg0 *big.Int
	Arg1 common.Address
}

type RecordBillPaymentInput struct {
	TokenId          *big.Int
	BillIndex        *big.Int
	PaymentMethod    string
	PaymentReference [32]byte
}

type RentalsInput struct {
	Arg0 *big.Int
}

type ResolveRentalDisputeInput struct {
	TokenId         *big.Int
	DepositToOwner  *big.Int
	DepositToRenter *big.Int
}

type RevokeCREWorkflowInput struct {
	Workflow common.Address
}

type RevokeKYCInput struct {
	User common.Address
}

type SafeTransferFromInput struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
}

type SafeTransferFrom0Input struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Data    []byte
}

type SetApprovalForAllInput struct {
	Operator common.Address
	Approved bool
}

type SetArbitratorInput struct {
	Arbitrator   common.Address
	IsArbitrator bool
}

type SetFeeRecipientInput struct {
	FeeRecipient common.Address
}

type SetKYCRequirementsInput struct {
	MinLevelForMint       uint8
	MinLevelForHighValue  uint8
	HighValueThresholdUSD *big.Int
}

type SetKYCVerificationInput struct {
	User             common.Address
	Level            uint8
	VerificationHash [32]byte
	ExpiryDate       *big.Int
}

type SetMintingPausedInput struct {
	Paused bool
	Reason string
}

type SetPaymentsPausedInput struct {
	Paused bool
	Reason string
}

type SetRentalsPausedInput struct {
	Paused bool
	Reason string
}

type SetSalesPausedInput struct {
	Paused bool
	Reason string
}

type SetTrustedBillProviderInput struct {
	Provider common.Address
	Trusted  bool
}

type SlashValidatorInput struct {
	Validator common.Address
	Amount    *big.Int
	Reason    string
}

type StartRentalInput struct {
	TokenId            *big.Int
	Renter             common.Address
	DurationDays       *big.Int
	DepositAmount      *big.Int
	MonthlyRent        *big.Int
	EncryptedAccessKey []byte
}

type SupportsInterfaceInput struct {
	InterfaceId [4]byte
}

type TokenByIndexInput struct {
	Index *big.Int
}

type TokenOfOwnerByIndexInput struct {
	Owner common.Address
	Index *big.Int
}

type TokenURIInput struct {
	TokenId *big.Int
}

type TransferFromInput struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
}

type TransferOwnershipInput struct {
	NewOwner common.Address
}

type TrustedBillProvidersInput struct {
	Arg0 common.Address
}

type UnstakeInput struct {
	Amount *big.Int
}

type UpgradeToAndCallInput struct {
	NewImplementation common.Address
	Data              []byte
}

type ValidatorsInput struct {
	Arg0 common.Address
}

type WithdrawRentalDepositInput struct {
	TokenId *big.Int
}

// Contract Method Outputs
type HouseBillsOutput struct {
	BillType           string
	Amount             *big.Int
	DueDate            *big.Int
	PaidAt             *big.Int
	Status             uint8
	PaymentReference   [32]byte
	IsRecurring        bool
	Provider           common.Address
	RecurrenceInterval uint8
}

type HousesOutput struct {
	HouseId       string
	DocumentHash  [32]byte
	DocumentURI   string
	StorageType   uint8
	OriginalOwner common.Address
	MintedAt      *big.Int
	IsVerified    bool
	DocumentCount uint8
}

type KeyExchangesOutput struct {
	KeyHash           [32]byte
	EncryptedKey      []byte
	IntendedRecipient common.Address
	CreatedAt         *big.Int
	ExpiresAt         *big.Int
	IsClaimed         bool
	ExchangeType      uint8
}

type KycInfoOutput struct {
	IsVerified        bool
	VerifiedAt        *big.Int
	VerificationHash  [32]byte
	VerificationLevel uint8
	ExpiryDate        *big.Int
}

type ListingsOutput struct {
	ListingType    uint8
	Price          *big.Int
	PreferredToken common.Address
	IsPrivateSale  bool
	AllowedBuyer   common.Address
	CreatedAt      *big.Int
	ExpiresAt      *big.Int
	PlatformFee    uint8
}

type RentalsOutput struct {
	Renter                 common.Address
	StartTime              *big.Int
	EndTime                *big.Int
	DepositAmount          *big.Int
	MonthlyRent            *big.Int
	IsActive               bool
	EncryptedAccessKeyHash [32]byte
	DisputeStatus          uint8
}

type ValidatorsOutput struct {
	StakedAmount          *big.Int
	StakedAt              *big.Int
	Reputation            uint8
	IsSlashed             bool
	SuccessfulValidations uint8
	FailedValidations     uint8
}

// Errors
type AddressEmptyCode struct {
	Target common.Address
}

type ERC1967InvalidImplementation struct {
	Implementation common.Address
}

type ERC1967NonPayable struct {
}

type ERC721EnumerableForbiddenBatchMint struct {
}

type ERC721IncorrectOwner struct {
	Sender  common.Address
	TokenId *big.Int
	Owner   common.Address
}

type ERC721InsufficientApproval struct {
	Operator common.Address
	TokenId  *big.Int
}

type ERC721InvalidApprover struct {
	Approver common.Address
}

type ERC721InvalidOperator struct {
	Operator common.Address
}

type ERC721InvalidOwner struct {
	Owner common.Address
}

type ERC721InvalidReceiver struct {
	Receiver common.Address
}

type ERC721InvalidSender struct {
	Sender common.Address
}

type ERC721NonexistentToken struct {
	TokenId *big.Int
}

type ERC721OutOfBoundsIndex struct {
	Owner common.Address
	Index *big.Int
}

type EnforcedPause struct {
}

type ExpectedPause struct {
}

type FailedInnerCall struct {
}

type InvalidInitialization struct {
}

type NotInitializing struct {
}

type OwnableInvalidOwner struct {
	Owner common.Address
}

type OwnableUnauthorizedAccount struct {
	Account common.Address
}

type ReentrancyGuardReentrantCall struct {
}

type UUPSUnauthorizedCallContext struct {
}

type UUPSUnsupportedProxiableUUID struct {
	Slot [32]byte
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

type ApprovalTopics struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
}

type ApprovalDecoded struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
}

type ApprovalForAllTopics struct {
	Owner    common.Address
	Operator common.Address
}

type ApprovalForAllDecoded struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
}

type BillCreatedTopics struct {
	TokenId   *big.Int
	BillIndex *big.Int
}

type BillCreatedDecoded struct {
	TokenId   *big.Int
	BillIndex *big.Int
	BillType  string
	Amount    *big.Int
	DueDate   *big.Int
}

type BillDisputedTopics struct {
	TokenId   *big.Int
	BillIndex *big.Int
}

type BillDisputedDecoded struct {
	TokenId   *big.Int
	BillIndex *big.Int
	Disputer  common.Address
	Reason    string
}

type BillPaidTopics struct {
	TokenId   *big.Int
	BillIndex *big.Int
}

type BillPaidDecoded struct {
	TokenId          *big.Int
	BillIndex        *big.Int
	PaymentMethod    string
	PaymentReference [32]byte
}

type CircuitBreakerTriggeredTopics struct {
}

type CircuitBreakerTriggeredDecoded struct {
	Component string
	Reason    string
}

type EmergencyActionTopics struct {
}

type EmergencyActionDecoded struct {
	Action      string
	TriggeredBy common.Address
	Timestamp   *big.Int
}

type HouseListedTopics struct {
	TokenId *big.Int
}

type HouseListedDecoded struct {
	TokenId       *big.Int
	ListingType   uint8
	Price         *big.Int
	IsPrivateSale bool
}

type HouseMintedTopics struct {
	TokenId *big.Int
	Owner   common.Address
}

type HouseMintedDecoded struct {
	TokenId      *big.Int
	Owner        common.Address
	HouseId      string
	DocumentHash [32]byte
	StorageType  uint8
}

type HouseSoldTopics struct {
	TokenId *big.Int
	Seller  common.Address
	Buyer   common.Address
}

type HouseSoldDecoded struct {
	TokenId     *big.Int
	Seller      common.Address
	Buyer       common.Address
	Price       *big.Int
	ProtocolFee *big.Int
}

type InitializedTopics struct {
}

type InitializedDecoded struct {
	Version uint64
}

type KYCVerifiedTopics struct {
	User common.Address
}

type KYCVerifiedDecoded struct {
	User   common.Address
	Level  uint8
	Expiry *big.Int
}

type KeyClaimedTopics struct {
	KeyHash  [32]byte
	Claimant common.Address
}

type KeyClaimedDecoded struct {
	KeyHash   [32]byte
	Claimant  common.Address
	Timestamp *big.Int
}

type KeyExchangeCreatedTopics struct {
	KeyHash   [32]byte
	TokenId   *big.Int
	Recipient common.Address
}

type KeyExchangeCreatedDecoded struct {
	KeyHash   [32]byte
	TokenId   *big.Int
	Recipient common.Address
	Expiry    *big.Int
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

type PausedTopics struct {
}

type PausedDecoded struct {
	Account common.Address
}

type RentalDepositReceivedTopics struct {
	TokenId *big.Int
	Renter  common.Address
}

type RentalDepositReceivedDecoded struct {
	TokenId *big.Int
	Renter  common.Address
	Amount  *big.Int
}

type RentalDepositWithdrawnTopics struct {
	TokenId *big.Int
	Renter  common.Address
}

type RentalDepositWithdrawnDecoded struct {
	TokenId *big.Int
	Renter  common.Address
	Amount  *big.Int
}

type RentalEndedTopics struct {
	TokenId *big.Int
	Renter  common.Address
}

type RentalEndedDecoded struct {
	TokenId         *big.Int
	Renter          common.Address
	DepositReturned *big.Int
}

type RentalStartedTopics struct {
	TokenId *big.Int
	Renter  common.Address
}

type RentalStartedDecoded struct {
	TokenId   *big.Int
	Renter    common.Address
	StartTime *big.Int
	EndTime   *big.Int
	Deposit   *big.Int
}

type TransferTopics struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
}

type TransferDecoded struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
}

type UnpausedTopics struct {
}

type UnpausedDecoded struct {
	Account common.Address
}

type UpgradedTopics struct {
	Implementation common.Address
}

type UpgradedDecoded struct {
	Implementation common.Address
}

type ValidatorSlashedTopics struct {
	Validator common.Address
}

type ValidatorSlashedDecoded struct {
	Validator common.Address
	Amount    *big.Int
	Reason    string
}

type ValidatorStakedTopics struct {
	Validator common.Address
}

type ValidatorStakedDecoded struct {
	Validator common.Address
	Amount    *big.Int
}

// Main Binding Type for HouseRWA
type HouseRWA struct {
	Address common.Address
	Options *bindings.ContractInitOptions
	ABI     *abi.ABI
	client  *evm.Client
	Codec   HouseRWACodec
}

type HouseRWACodec interface {
	EncodeBPSDENOMINATORMethodCall() ([]byte, error)
	DecodeBPSDENOMINATORMethodOutput(data []byte) (*big.Int, error)
	EncodeKEYEXCHANGEEXPIRYMethodCall() ([]byte, error)
	DecodeKEYEXCHANGEEXPIRYMethodOutput(data []byte) (*big.Int, error)
	EncodeMAXBILLSPERHOUSEMethodCall() ([]byte, error)
	DecodeMAXBILLSPERHOUSEMethodOutput(data []byte) (*big.Int, error)
	EncodeMAXDOCUMENTSPERHOUSEMethodCall() ([]byte, error)
	DecodeMAXDOCUMENTSPERHOUSEMethodOutput(data []byte) (*big.Int, error)
	EncodeMAXRENTALDURATIONMethodCall() ([]byte, error)
	DecodeMAXRENTALDURATIONMethodOutput(data []byte) (*big.Int, error)
	EncodeMINVALIDATORSTAKEMethodCall() ([]byte, error)
	DecodeMINVALIDATORSTAKEMethodOutput(data []byte) (*big.Int, error)
	EncodePROTOCOLFEEBPSMethodCall() ([]byte, error)
	DecodePROTOCOLFEEBPSMethodOutput(data []byte) (*big.Int, error)
	EncodeUPGRADEINTERFACEVERSIONMethodCall() ([]byte, error)
	DecodeUPGRADEINTERFACEVERSIONMethodOutput(data []byte) (string, error)
	EncodeAcceptOwnershipMethodCall() ([]byte, error)
	EncodeApproveMethodCall(in ApproveInput) ([]byte, error)
	EncodeArbitratorsMethodCall(in ArbitratorsInput) ([]byte, error)
	DecodeArbitratorsMethodOutput(data []byte) (bool, error)
	EncodeAuthorizeCREWorkflowMethodCall(in AuthorizeCREWorkflowInput) ([]byte, error)
	EncodeAuthorizedCREWorkflowsMethodCall(in AuthorizedCREWorkflowsInput) ([]byte, error)
	DecodeAuthorizedCREWorkflowsMethodOutput(data []byte) (bool, error)
	EncodeBalanceOfMethodCall(in BalanceOfInput) ([]byte, error)
	DecodeBalanceOfMethodOutput(data []byte) (*big.Int, error)
	EncodeCancelListingMethodCall(in CancelListingInput) ([]byte, error)
	EncodeClaimKeyMethodCall(in ClaimKeyInput) ([]byte, error)
	DecodeClaimKeyMethodOutput(data []byte) ([]byte, error)
	EncodeCompleteSaleMethodCall(in CompleteSaleInput) ([]byte, error)
	EncodeCreateBillMethodCall(in CreateBillInput) ([]byte, error)
	DecodeCreateBillMethodOutput(data []byte) (*big.Int, error)
	EncodeCreateListingMethodCall(in CreateListingInput) ([]byte, error)
	EncodeDepositForRentalMethodCall(in DepositForRentalInput) ([]byte, error)
	EncodeDisputeBillMethodCall(in DisputeBillInput) ([]byte, error)
	EncodeEmergencyPauseMethodCall() ([]byte, error)
	EncodeEmergencyUnpauseMethodCall() ([]byte, error)
	EncodeEndRentalMethodCall(in EndRentalInput) ([]byte, error)
	EncodeFeeRecipientMethodCall() ([]byte, error)
	DecodeFeeRecipientMethodOutput(data []byte) (common.Address, error)
	EncodeGetActiveRentalMethodCall(in GetActiveRentalInput) ([]byte, error)
	DecodeGetActiveRentalMethodOutput(data []byte) (RentalAgreement, error)
	EncodeGetApprovedMethodCall(in GetApprovedInput) ([]byte, error)
	DecodeGetApprovedMethodOutput(data []byte) (common.Address, error)
	EncodeGetBillsMethodCall(in GetBillsInput) ([]byte, error)
	DecodeGetBillsMethodOutput(data []byte) ([]Bill, error)
	EncodeGetHouseDetailsMethodCall(in GetHouseDetailsInput) ([]byte, error)
	DecodeGetHouseDetailsMethodOutput(data []byte) (House, error)
	EncodeGetListingMethodCall(in GetListingInput) ([]byte, error)
	DecodeGetListingMethodOutput(data []byte) (Listing, error)
	EncodeGetTotalBillsCountMethodCall(in GetTotalBillsCountInput) ([]byte, error)
	DecodeGetTotalBillsCountMethodOutput(data []byte) (*big.Int, error)
	EncodeHasKYCMethodCall(in HasKYCInput) ([]byte, error)
	DecodeHasKYCMethodOutput(data []byte) (bool, error)
	EncodeHighValueThresholdUSDMethodCall() ([]byte, error)
	DecodeHighValueThresholdUSDMethodOutput(data []byte) (*big.Int, error)
	EncodeHouseBillsMethodCall(in HouseBillsInput) ([]byte, error)
	DecodeHouseBillsMethodOutput(data []byte) (HouseBillsOutput, error)
	EncodeHousesMethodCall(in HousesInput) ([]byte, error)
	DecodeHousesMethodOutput(data []byte) (HousesOutput, error)
	EncodeInitializeMethodCall(in InitializeInput) ([]byte, error)
	EncodeIsApprovedForAllMethodCall(in IsApprovedForAllInput) ([]byte, error)
	DecodeIsApprovedForAllMethodOutput(data []byte) (bool, error)
	EncodeIsRentedMethodCall(in IsRentedInput) ([]byte, error)
	DecodeIsRentedMethodOutput(data []byte) (bool, error)
	EncodeKeyExchangesMethodCall(in KeyExchangesInput) ([]byte, error)
	DecodeKeyExchangesMethodOutput(data []byte) (KeyExchangesOutput, error)
	EncodeKycInfoMethodCall(in KycInfoInput) ([]byte, error)
	DecodeKycInfoMethodOutput(data []byte) (KycInfoOutput, error)
	EncodeListingsMethodCall(in ListingsInput) ([]byte, error)
	DecodeListingsMethodOutput(data []byte) (ListingsOutput, error)
	EncodeMinKYCLevelForHighValueMethodCall() ([]byte, error)
	DecodeMinKYCLevelForHighValueMethodOutput(data []byte) (uint8, error)
	EncodeMinKYCLevelForMintMethodCall() ([]byte, error)
	DecodeMinKYCLevelForMintMethodOutput(data []byte) (uint8, error)
	EncodeMintMethodCall(in MintInput) ([]byte, error)
	DecodeMintMethodOutput(data []byte) (*big.Int, error)
	EncodeMintingPausedMethodCall() ([]byte, error)
	DecodeMintingPausedMethodOutput(data []byte) (bool, error)
	EncodeNameMethodCall() ([]byte, error)
	DecodeNameMethodOutput(data []byte) (string, error)
	EncodeNextTokenIdMethodCall() ([]byte, error)
	DecodeNextTokenIdMethodOutput(data []byte) (*big.Int, error)
	EncodeOpenRentalDisputeMethodCall(in OpenRentalDisputeInput) ([]byte, error)
	EncodeOwnerMethodCall() ([]byte, error)
	DecodeOwnerMethodOutput(data []byte) (common.Address, error)
	EncodeOwnerOfMethodCall(in OwnerOfInput) ([]byte, error)
	DecodeOwnerOfMethodOutput(data []byte) (common.Address, error)
	EncodePausedMethodCall() ([]byte, error)
	DecodePausedMethodOutput(data []byte) (bool, error)
	EncodePaymentsPausedMethodCall() ([]byte, error)
	DecodePaymentsPausedMethodOutput(data []byte) (bool, error)
	EncodePendingOwnerMethodCall() ([]byte, error)
	DecodePendingOwnerMethodOutput(data []byte) (common.Address, error)
	EncodePendingRentalDepositsMethodCall(in PendingRentalDepositsInput) ([]byte, error)
	DecodePendingRentalDepositsMethodOutput(data []byte) (*big.Int, error)
	EncodePriceFeedMethodCall() ([]byte, error)
	DecodePriceFeedMethodOutput(data []byte) (common.Address, error)
	EncodeProxiableUUIDMethodCall() ([]byte, error)
	DecodeProxiableUUIDMethodOutput(data []byte) ([32]byte, error)
	EncodeRecordBillPaymentMethodCall(in RecordBillPaymentInput) ([]byte, error)
	EncodeRenounceOwnershipMethodCall() ([]byte, error)
	EncodeRentalsMethodCall(in RentalsInput) ([]byte, error)
	DecodeRentalsMethodOutput(data []byte) (RentalsOutput, error)
	EncodeRentalsPausedMethodCall() ([]byte, error)
	DecodeRentalsPausedMethodOutput(data []byte) (bool, error)
	EncodeResolveRentalDisputeMethodCall(in ResolveRentalDisputeInput) ([]byte, error)
	EncodeRevokeCREWorkflowMethodCall(in RevokeCREWorkflowInput) ([]byte, error)
	EncodeRevokeKYCMethodCall(in RevokeKYCInput) ([]byte, error)
	EncodeSafeTransferFromMethodCall(in SafeTransferFromInput) ([]byte, error)
	EncodeSafeTransferFrom0MethodCall(in SafeTransferFrom0Input) ([]byte, error)
	EncodeSalesPausedMethodCall() ([]byte, error)
	DecodeSalesPausedMethodOutput(data []byte) (bool, error)
	EncodeSetApprovalForAllMethodCall(in SetApprovalForAllInput) ([]byte, error)
	EncodeSetArbitratorMethodCall(in SetArbitratorInput) ([]byte, error)
	EncodeSetFeeRecipientMethodCall(in SetFeeRecipientInput) ([]byte, error)
	EncodeSetKYCRequirementsMethodCall(in SetKYCRequirementsInput) ([]byte, error)
	EncodeSetKYCVerificationMethodCall(in SetKYCVerificationInput) ([]byte, error)
	EncodeSetMintingPausedMethodCall(in SetMintingPausedInput) ([]byte, error)
	EncodeSetPaymentsPausedMethodCall(in SetPaymentsPausedInput) ([]byte, error)
	EncodeSetRentalsPausedMethodCall(in SetRentalsPausedInput) ([]byte, error)
	EncodeSetSalesPausedMethodCall(in SetSalesPausedInput) ([]byte, error)
	EncodeSetTrustedBillProviderMethodCall(in SetTrustedBillProviderInput) ([]byte, error)
	EncodeSlashValidatorMethodCall(in SlashValidatorInput) ([]byte, error)
	EncodeStakeAsValidatorMethodCall() ([]byte, error)
	EncodeStartRentalMethodCall(in StartRentalInput) ([]byte, error)
	EncodeSupportsInterfaceMethodCall(in SupportsInterfaceInput) ([]byte, error)
	DecodeSupportsInterfaceMethodOutput(data []byte) (bool, error)
	EncodeSymbolMethodCall() ([]byte, error)
	DecodeSymbolMethodOutput(data []byte) (string, error)
	EncodeTokenByIndexMethodCall(in TokenByIndexInput) ([]byte, error)
	DecodeTokenByIndexMethodOutput(data []byte) (*big.Int, error)
	EncodeTokenOfOwnerByIndexMethodCall(in TokenOfOwnerByIndexInput) ([]byte, error)
	DecodeTokenOfOwnerByIndexMethodOutput(data []byte) (*big.Int, error)
	EncodeTokenURIMethodCall(in TokenURIInput) ([]byte, error)
	DecodeTokenURIMethodOutput(data []byte) (string, error)
	EncodeTotalFeesCollectedMethodCall() ([]byte, error)
	DecodeTotalFeesCollectedMethodOutput(data []byte) (*big.Int, error)
	EncodeTotalSupplyMethodCall() ([]byte, error)
	DecodeTotalSupplyMethodOutput(data []byte) (*big.Int, error)
	EncodeTransferFromMethodCall(in TransferFromInput) ([]byte, error)
	EncodeTransferOwnershipMethodCall(in TransferOwnershipInput) ([]byte, error)
	EncodeTrustedBillProvidersMethodCall(in TrustedBillProvidersInput) ([]byte, error)
	DecodeTrustedBillProvidersMethodOutput(data []byte) (bool, error)
	EncodeUnstakeMethodCall(in UnstakeInput) ([]byte, error)
	EncodeUpgradeToAndCallMethodCall(in UpgradeToAndCallInput) ([]byte, error)
	EncodeValidatorsMethodCall(in ValidatorsInput) ([]byte, error)
	DecodeValidatorsMethodOutput(data []byte) (ValidatorsOutput, error)
	EncodeWithdrawRentalDepositMethodCall(in WithdrawRentalDepositInput) ([]byte, error)
	EncodeBillStruct(in Bill) ([]byte, error)
	EncodeHouseStruct(in House) ([]byte, error)
	EncodeListingStruct(in Listing) ([]byte, error)
	EncodeRentalAgreementStruct(in RentalAgreement) ([]byte, error)
	ApprovalLogHash() []byte
	EncodeApprovalTopics(evt abi.Event, values []ApprovalTopics) ([]*evm.TopicValues, error)
	DecodeApproval(log *evm.Log) (*ApprovalDecoded, error)
	ApprovalForAllLogHash() []byte
	EncodeApprovalForAllTopics(evt abi.Event, values []ApprovalForAllTopics) ([]*evm.TopicValues, error)
	DecodeApprovalForAll(log *evm.Log) (*ApprovalForAllDecoded, error)
	BillCreatedLogHash() []byte
	EncodeBillCreatedTopics(evt abi.Event, values []BillCreatedTopics) ([]*evm.TopicValues, error)
	DecodeBillCreated(log *evm.Log) (*BillCreatedDecoded, error)
	BillDisputedLogHash() []byte
	EncodeBillDisputedTopics(evt abi.Event, values []BillDisputedTopics) ([]*evm.TopicValues, error)
	DecodeBillDisputed(log *evm.Log) (*BillDisputedDecoded, error)
	BillPaidLogHash() []byte
	EncodeBillPaidTopics(evt abi.Event, values []BillPaidTopics) ([]*evm.TopicValues, error)
	DecodeBillPaid(log *evm.Log) (*BillPaidDecoded, error)
	CircuitBreakerTriggeredLogHash() []byte
	EncodeCircuitBreakerTriggeredTopics(evt abi.Event, values []CircuitBreakerTriggeredTopics) ([]*evm.TopicValues, error)
	DecodeCircuitBreakerTriggered(log *evm.Log) (*CircuitBreakerTriggeredDecoded, error)
	EmergencyActionLogHash() []byte
	EncodeEmergencyActionTopics(evt abi.Event, values []EmergencyActionTopics) ([]*evm.TopicValues, error)
	DecodeEmergencyAction(log *evm.Log) (*EmergencyActionDecoded, error)
	HouseListedLogHash() []byte
	EncodeHouseListedTopics(evt abi.Event, values []HouseListedTopics) ([]*evm.TopicValues, error)
	DecodeHouseListed(log *evm.Log) (*HouseListedDecoded, error)
	HouseMintedLogHash() []byte
	EncodeHouseMintedTopics(evt abi.Event, values []HouseMintedTopics) ([]*evm.TopicValues, error)
	DecodeHouseMinted(log *evm.Log) (*HouseMintedDecoded, error)
	HouseSoldLogHash() []byte
	EncodeHouseSoldTopics(evt abi.Event, values []HouseSoldTopics) ([]*evm.TopicValues, error)
	DecodeHouseSold(log *evm.Log) (*HouseSoldDecoded, error)
	InitializedLogHash() []byte
	EncodeInitializedTopics(evt abi.Event, values []InitializedTopics) ([]*evm.TopicValues, error)
	DecodeInitialized(log *evm.Log) (*InitializedDecoded, error)
	KYCVerifiedLogHash() []byte
	EncodeKYCVerifiedTopics(evt abi.Event, values []KYCVerifiedTopics) ([]*evm.TopicValues, error)
	DecodeKYCVerified(log *evm.Log) (*KYCVerifiedDecoded, error)
	KeyClaimedLogHash() []byte
	EncodeKeyClaimedTopics(evt abi.Event, values []KeyClaimedTopics) ([]*evm.TopicValues, error)
	DecodeKeyClaimed(log *evm.Log) (*KeyClaimedDecoded, error)
	KeyExchangeCreatedLogHash() []byte
	EncodeKeyExchangeCreatedTopics(evt abi.Event, values []KeyExchangeCreatedTopics) ([]*evm.TopicValues, error)
	DecodeKeyExchangeCreated(log *evm.Log) (*KeyExchangeCreatedDecoded, error)
	OwnershipTransferStartedLogHash() []byte
	EncodeOwnershipTransferStartedTopics(evt abi.Event, values []OwnershipTransferStartedTopics) ([]*evm.TopicValues, error)
	DecodeOwnershipTransferStarted(log *evm.Log) (*OwnershipTransferStartedDecoded, error)
	OwnershipTransferredLogHash() []byte
	EncodeOwnershipTransferredTopics(evt abi.Event, values []OwnershipTransferredTopics) ([]*evm.TopicValues, error)
	DecodeOwnershipTransferred(log *evm.Log) (*OwnershipTransferredDecoded, error)
	PausedLogHash() []byte
	EncodePausedTopics(evt abi.Event, values []PausedTopics) ([]*evm.TopicValues, error)
	DecodePaused(log *evm.Log) (*PausedDecoded, error)
	RentalDepositReceivedLogHash() []byte
	EncodeRentalDepositReceivedTopics(evt abi.Event, values []RentalDepositReceivedTopics) ([]*evm.TopicValues, error)
	DecodeRentalDepositReceived(log *evm.Log) (*RentalDepositReceivedDecoded, error)
	RentalDepositWithdrawnLogHash() []byte
	EncodeRentalDepositWithdrawnTopics(evt abi.Event, values []RentalDepositWithdrawnTopics) ([]*evm.TopicValues, error)
	DecodeRentalDepositWithdrawn(log *evm.Log) (*RentalDepositWithdrawnDecoded, error)
	RentalEndedLogHash() []byte
	EncodeRentalEndedTopics(evt abi.Event, values []RentalEndedTopics) ([]*evm.TopicValues, error)
	DecodeRentalEnded(log *evm.Log) (*RentalEndedDecoded, error)
	RentalStartedLogHash() []byte
	EncodeRentalStartedTopics(evt abi.Event, values []RentalStartedTopics) ([]*evm.TopicValues, error)
	DecodeRentalStarted(log *evm.Log) (*RentalStartedDecoded, error)
	TransferLogHash() []byte
	EncodeTransferTopics(evt abi.Event, values []TransferTopics) ([]*evm.TopicValues, error)
	DecodeTransfer(log *evm.Log) (*TransferDecoded, error)
	UnpausedLogHash() []byte
	EncodeUnpausedTopics(evt abi.Event, values []UnpausedTopics) ([]*evm.TopicValues, error)
	DecodeUnpaused(log *evm.Log) (*UnpausedDecoded, error)
	UpgradedLogHash() []byte
	EncodeUpgradedTopics(evt abi.Event, values []UpgradedTopics) ([]*evm.TopicValues, error)
	DecodeUpgraded(log *evm.Log) (*UpgradedDecoded, error)
	ValidatorSlashedLogHash() []byte
	EncodeValidatorSlashedTopics(evt abi.Event, values []ValidatorSlashedTopics) ([]*evm.TopicValues, error)
	DecodeValidatorSlashed(log *evm.Log) (*ValidatorSlashedDecoded, error)
	ValidatorStakedLogHash() []byte
	EncodeValidatorStakedTopics(evt abi.Event, values []ValidatorStakedTopics) ([]*evm.TopicValues, error)
	DecodeValidatorStaked(log *evm.Log) (*ValidatorStakedDecoded, error)
}

func NewHouseRWA(
	client *evm.Client,
	address common.Address,
	options *bindings.ContractInitOptions,
) (*HouseRWA, error) {
	parsed, err := abi.JSON(strings.NewReader(HouseRWAMetaData.ABI))
	if err != nil {
		return nil, err
	}
	codec, err := NewCodec()
	if err != nil {
		return nil, err
	}
	return &HouseRWA{
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

func NewCodec() (HouseRWACodec, error) {
	parsed, err := abi.JSON(strings.NewReader(HouseRWAMetaData.ABI))
	if err != nil {
		return nil, err
	}
	return &Codec{abi: &parsed}, nil
}

func (c *Codec) EncodeBPSDENOMINATORMethodCall() ([]byte, error) {
	return c.abi.Pack("BPS_DENOMINATOR")
}

func (c *Codec) DecodeBPSDENOMINATORMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["BPS_DENOMINATOR"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeKEYEXCHANGEEXPIRYMethodCall() ([]byte, error) {
	return c.abi.Pack("KEY_EXCHANGE_EXPIRY")
}

func (c *Codec) DecodeKEYEXCHANGEEXPIRYMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["KEY_EXCHANGE_EXPIRY"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeMAXBILLSPERHOUSEMethodCall() ([]byte, error) {
	return c.abi.Pack("MAX_BILLS_PER_HOUSE")
}

func (c *Codec) DecodeMAXBILLSPERHOUSEMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["MAX_BILLS_PER_HOUSE"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeMAXDOCUMENTSPERHOUSEMethodCall() ([]byte, error) {
	return c.abi.Pack("MAX_DOCUMENTS_PER_HOUSE")
}

func (c *Codec) DecodeMAXDOCUMENTSPERHOUSEMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["MAX_DOCUMENTS_PER_HOUSE"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeMAXRENTALDURATIONMethodCall() ([]byte, error) {
	return c.abi.Pack("MAX_RENTAL_DURATION")
}

func (c *Codec) DecodeMAXRENTALDURATIONMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["MAX_RENTAL_DURATION"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeMINVALIDATORSTAKEMethodCall() ([]byte, error) {
	return c.abi.Pack("MIN_VALIDATOR_STAKE")
}

func (c *Codec) DecodeMINVALIDATORSTAKEMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["MIN_VALIDATOR_STAKE"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodePROTOCOLFEEBPSMethodCall() ([]byte, error) {
	return c.abi.Pack("PROTOCOL_FEE_BPS")
}

func (c *Codec) DecodePROTOCOLFEEBPSMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["PROTOCOL_FEE_BPS"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeUPGRADEINTERFACEVERSIONMethodCall() ([]byte, error) {
	return c.abi.Pack("UPGRADE_INTERFACE_VERSION")
}

func (c *Codec) DecodeUPGRADEINTERFACEVERSIONMethodOutput(data []byte) (string, error) {
	vals, err := c.abi.Methods["UPGRADE_INTERFACE_VERSION"].Outputs.Unpack(data)
	if err != nil {
		return *new(string), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(string), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result string
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(string), fmt.Errorf("failed to unmarshal to string: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeAcceptOwnershipMethodCall() ([]byte, error) {
	return c.abi.Pack("acceptOwnership")
}

func (c *Codec) EncodeApproveMethodCall(in ApproveInput) ([]byte, error) {
	return c.abi.Pack("approve", in.To, in.TokenId)
}

func (c *Codec) EncodeArbitratorsMethodCall(in ArbitratorsInput) ([]byte, error) {
	return c.abi.Pack("arbitrators", in.Arg0)
}

func (c *Codec) DecodeArbitratorsMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["arbitrators"].Outputs.Unpack(data)
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

func (c *Codec) EncodeAuthorizeCREWorkflowMethodCall(in AuthorizeCREWorkflowInput) ([]byte, error) {
	return c.abi.Pack("authorizeCREWorkflow", in.Workflow)
}

func (c *Codec) EncodeAuthorizedCREWorkflowsMethodCall(in AuthorizedCREWorkflowsInput) ([]byte, error) {
	return c.abi.Pack("authorizedCREWorkflows", in.Arg0)
}

func (c *Codec) DecodeAuthorizedCREWorkflowsMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["authorizedCREWorkflows"].Outputs.Unpack(data)
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

func (c *Codec) EncodeBalanceOfMethodCall(in BalanceOfInput) ([]byte, error) {
	return c.abi.Pack("balanceOf", in.Owner)
}

func (c *Codec) DecodeBalanceOfMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["balanceOf"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeCancelListingMethodCall(in CancelListingInput) ([]byte, error) {
	return c.abi.Pack("cancelListing", in.TokenId)
}

func (c *Codec) EncodeClaimKeyMethodCall(in ClaimKeyInput) ([]byte, error) {
	return c.abi.Pack("claimKey", in.KeyHash)
}

func (c *Codec) DecodeClaimKeyMethodOutput(data []byte) ([]byte, error) {
	vals, err := c.abi.Methods["claimKey"].Outputs.Unpack(data)
	if err != nil {
		return *new([]byte), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new([]byte), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result []byte
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new([]byte), fmt.Errorf("failed to unmarshal to []byte: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeCompleteSaleMethodCall(in CompleteSaleInput) ([]byte, error) {
	return c.abi.Pack("completeSale", in.TokenId, in.Buyer, in.KeyHash, in.EncryptedKey)
}

func (c *Codec) EncodeCreateBillMethodCall(in CreateBillInput) ([]byte, error) {
	return c.abi.Pack("createBill", in.TokenId, in.BillType, in.Amount, in.DueDate, in.Provider, in.IsRecurring, in.RecurrenceInterval)
}

func (c *Codec) DecodeCreateBillMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["createBill"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeCreateListingMethodCall(in CreateListingInput) ([]byte, error) {
	return c.abi.Pack("createListing", in.TokenId, in.ListingType, in.Price, in.PreferredToken, in.IsPrivateSale, in.AllowedBuyer, in.DurationDays)
}

func (c *Codec) EncodeDepositForRentalMethodCall(in DepositForRentalInput) ([]byte, error) {
	return c.abi.Pack("depositForRental", in.TokenId)
}

func (c *Codec) EncodeDisputeBillMethodCall(in DisputeBillInput) ([]byte, error) {
	return c.abi.Pack("disputeBill", in.TokenId, in.BillIndex, in.Reason)
}

func (c *Codec) EncodeEmergencyPauseMethodCall() ([]byte, error) {
	return c.abi.Pack("emergencyPause")
}

func (c *Codec) EncodeEmergencyUnpauseMethodCall() ([]byte, error) {
	return c.abi.Pack("emergencyUnpause")
}

func (c *Codec) EncodeEndRentalMethodCall(in EndRentalInput) ([]byte, error) {
	return c.abi.Pack("endRental", in.TokenId)
}

func (c *Codec) EncodeFeeRecipientMethodCall() ([]byte, error) {
	return c.abi.Pack("feeRecipient")
}

func (c *Codec) DecodeFeeRecipientMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["feeRecipient"].Outputs.Unpack(data)
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

func (c *Codec) EncodeGetActiveRentalMethodCall(in GetActiveRentalInput) ([]byte, error) {
	return c.abi.Pack("getActiveRental", in.TokenId)
}

func (c *Codec) DecodeGetActiveRentalMethodOutput(data []byte) (RentalAgreement, error) {
	vals, err := c.abi.Methods["getActiveRental"].Outputs.Unpack(data)
	if err != nil {
		return *new(RentalAgreement), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(RentalAgreement), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result RentalAgreement
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(RentalAgreement), fmt.Errorf("failed to unmarshal to RentalAgreement: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetApprovedMethodCall(in GetApprovedInput) ([]byte, error) {
	return c.abi.Pack("getApproved", in.TokenId)
}

func (c *Codec) DecodeGetApprovedMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["getApproved"].Outputs.Unpack(data)
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

func (c *Codec) EncodeGetBillsMethodCall(in GetBillsInput) ([]byte, error) {
	return c.abi.Pack("getBills", in.TokenId)
}

func (c *Codec) DecodeGetBillsMethodOutput(data []byte) ([]Bill, error) {
	vals, err := c.abi.Methods["getBills"].Outputs.Unpack(data)
	if err != nil {
		return *new([]Bill), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new([]Bill), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result []Bill
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new([]Bill), fmt.Errorf("failed to unmarshal to []Bill: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetHouseDetailsMethodCall(in GetHouseDetailsInput) ([]byte, error) {
	return c.abi.Pack("getHouseDetails", in.TokenId)
}

func (c *Codec) DecodeGetHouseDetailsMethodOutput(data []byte) (House, error) {
	vals, err := c.abi.Methods["getHouseDetails"].Outputs.Unpack(data)
	if err != nil {
		return *new(House), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(House), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result House
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(House), fmt.Errorf("failed to unmarshal to House: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetListingMethodCall(in GetListingInput) ([]byte, error) {
	return c.abi.Pack("getListing", in.TokenId)
}

func (c *Codec) DecodeGetListingMethodOutput(data []byte) (Listing, error) {
	vals, err := c.abi.Methods["getListing"].Outputs.Unpack(data)
	if err != nil {
		return *new(Listing), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(Listing), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result Listing
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(Listing), fmt.Errorf("failed to unmarshal to Listing: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeGetTotalBillsCountMethodCall(in GetTotalBillsCountInput) ([]byte, error) {
	return c.abi.Pack("getTotalBillsCount", in.TokenId)
}

func (c *Codec) DecodeGetTotalBillsCountMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["getTotalBillsCount"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeHasKYCMethodCall(in HasKYCInput) ([]byte, error) {
	return c.abi.Pack("hasKYC", in.User)
}

func (c *Codec) DecodeHasKYCMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["hasKYC"].Outputs.Unpack(data)
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

func (c *Codec) EncodeHighValueThresholdUSDMethodCall() ([]byte, error) {
	return c.abi.Pack("highValueThresholdUSD")
}

func (c *Codec) DecodeHighValueThresholdUSDMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["highValueThresholdUSD"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeHouseBillsMethodCall(in HouseBillsInput) ([]byte, error) {
	return c.abi.Pack("houseBills", in.Arg0, in.Arg1)
}

func (c *Codec) DecodeHouseBillsMethodOutput(data []byte) (HouseBillsOutput, error) {
	vals, err := c.abi.Methods["houseBills"].Outputs.Unpack(data)
	if err != nil {
		return HouseBillsOutput{}, err
	}
	if len(vals) != 9 {
		return HouseBillsOutput{}, fmt.Errorf("expected 9 values, got %d", len(vals))
	}
	jsonData0, err := json.Marshal(vals[0])
	if err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to marshal ABI result 0: %w", err)
	}

	var result0 string
	if err := json.Unmarshal(jsonData0, &result0); err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to unmarshal to string: %w", err)
	}
	jsonData1, err := json.Marshal(vals[1])
	if err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to marshal ABI result 1: %w", err)
	}

	var result1 *big.Int
	if err := json.Unmarshal(jsonData1, &result1); err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData2, err := json.Marshal(vals[2])
	if err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to marshal ABI result 2: %w", err)
	}

	var result2 *big.Int
	if err := json.Unmarshal(jsonData2, &result2); err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData3, err := json.Marshal(vals[3])
	if err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to marshal ABI result 3: %w", err)
	}

	var result3 *big.Int
	if err := json.Unmarshal(jsonData3, &result3); err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData4, err := json.Marshal(vals[4])
	if err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to marshal ABI result 4: %w", err)
	}

	var result4 uint8
	if err := json.Unmarshal(jsonData4, &result4); err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to unmarshal to uint8: %w", err)
	}
	jsonData5, err := json.Marshal(vals[5])
	if err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to marshal ABI result 5: %w", err)
	}

	var result5 [32]byte
	if err := json.Unmarshal(jsonData5, &result5); err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to unmarshal to [32]byte: %w", err)
	}
	jsonData6, err := json.Marshal(vals[6])
	if err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to marshal ABI result 6: %w", err)
	}

	var result6 bool
	if err := json.Unmarshal(jsonData6, &result6); err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to unmarshal to bool: %w", err)
	}
	jsonData7, err := json.Marshal(vals[7])
	if err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to marshal ABI result 7: %w", err)
	}

	var result7 common.Address
	if err := json.Unmarshal(jsonData7, &result7); err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}
	jsonData8, err := json.Marshal(vals[8])
	if err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to marshal ABI result 8: %w", err)
	}

	var result8 uint8
	if err := json.Unmarshal(jsonData8, &result8); err != nil {
		return HouseBillsOutput{}, fmt.Errorf("failed to unmarshal to uint8: %w", err)
	}

	return HouseBillsOutput{
		BillType:           result0,
		Amount:             result1,
		DueDate:            result2,
		PaidAt:             result3,
		Status:             result4,
		PaymentReference:   result5,
		IsRecurring:        result6,
		Provider:           result7,
		RecurrenceInterval: result8,
	}, nil
}

func (c *Codec) EncodeHousesMethodCall(in HousesInput) ([]byte, error) {
	return c.abi.Pack("houses", in.Arg0)
}

func (c *Codec) DecodeHousesMethodOutput(data []byte) (HousesOutput, error) {
	vals, err := c.abi.Methods["houses"].Outputs.Unpack(data)
	if err != nil {
		return HousesOutput{}, err
	}
	if len(vals) != 8 {
		return HousesOutput{}, fmt.Errorf("expected 8 values, got %d", len(vals))
	}
	jsonData0, err := json.Marshal(vals[0])
	if err != nil {
		return HousesOutput{}, fmt.Errorf("failed to marshal ABI result 0: %w", err)
	}

	var result0 string
	if err := json.Unmarshal(jsonData0, &result0); err != nil {
		return HousesOutput{}, fmt.Errorf("failed to unmarshal to string: %w", err)
	}
	jsonData1, err := json.Marshal(vals[1])
	if err != nil {
		return HousesOutput{}, fmt.Errorf("failed to marshal ABI result 1: %w", err)
	}

	var result1 [32]byte
	if err := json.Unmarshal(jsonData1, &result1); err != nil {
		return HousesOutput{}, fmt.Errorf("failed to unmarshal to [32]byte: %w", err)
	}
	jsonData2, err := json.Marshal(vals[2])
	if err != nil {
		return HousesOutput{}, fmt.Errorf("failed to marshal ABI result 2: %w", err)
	}

	var result2 string
	if err := json.Unmarshal(jsonData2, &result2); err != nil {
		return HousesOutput{}, fmt.Errorf("failed to unmarshal to string: %w", err)
	}
	jsonData3, err := json.Marshal(vals[3])
	if err != nil {
		return HousesOutput{}, fmt.Errorf("failed to marshal ABI result 3: %w", err)
	}

	var result3 uint8
	if err := json.Unmarshal(jsonData3, &result3); err != nil {
		return HousesOutput{}, fmt.Errorf("failed to unmarshal to uint8: %w", err)
	}
	jsonData4, err := json.Marshal(vals[4])
	if err != nil {
		return HousesOutput{}, fmt.Errorf("failed to marshal ABI result 4: %w", err)
	}

	var result4 common.Address
	if err := json.Unmarshal(jsonData4, &result4); err != nil {
		return HousesOutput{}, fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}
	jsonData5, err := json.Marshal(vals[5])
	if err != nil {
		return HousesOutput{}, fmt.Errorf("failed to marshal ABI result 5: %w", err)
	}

	var result5 *big.Int
	if err := json.Unmarshal(jsonData5, &result5); err != nil {
		return HousesOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData6, err := json.Marshal(vals[6])
	if err != nil {
		return HousesOutput{}, fmt.Errorf("failed to marshal ABI result 6: %w", err)
	}

	var result6 bool
	if err := json.Unmarshal(jsonData6, &result6); err != nil {
		return HousesOutput{}, fmt.Errorf("failed to unmarshal to bool: %w", err)
	}
	jsonData7, err := json.Marshal(vals[7])
	if err != nil {
		return HousesOutput{}, fmt.Errorf("failed to marshal ABI result 7: %w", err)
	}

	var result7 uint8
	if err := json.Unmarshal(jsonData7, &result7); err != nil {
		return HousesOutput{}, fmt.Errorf("failed to unmarshal to uint8: %w", err)
	}

	return HousesOutput{
		HouseId:       result0,
		DocumentHash:  result1,
		DocumentURI:   result2,
		StorageType:   result3,
		OriginalOwner: result4,
		MintedAt:      result5,
		IsVerified:    result6,
		DocumentCount: result7,
	}, nil
}

func (c *Codec) EncodeInitializeMethodCall(in InitializeInput) ([]byte, error) {
	return c.abi.Pack("initialize", in.Owner, in.FeeRecipient, in.InitialCREWorkflow)
}

func (c *Codec) EncodeIsApprovedForAllMethodCall(in IsApprovedForAllInput) ([]byte, error) {
	return c.abi.Pack("isApprovedForAll", in.Owner, in.Operator)
}

func (c *Codec) DecodeIsApprovedForAllMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["isApprovedForAll"].Outputs.Unpack(data)
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

func (c *Codec) EncodeIsRentedMethodCall(in IsRentedInput) ([]byte, error) {
	return c.abi.Pack("isRented", in.TokenId)
}

func (c *Codec) DecodeIsRentedMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["isRented"].Outputs.Unpack(data)
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

func (c *Codec) EncodeKeyExchangesMethodCall(in KeyExchangesInput) ([]byte, error) {
	return c.abi.Pack("keyExchanges", in.Arg0)
}

func (c *Codec) DecodeKeyExchangesMethodOutput(data []byte) (KeyExchangesOutput, error) {
	vals, err := c.abi.Methods["keyExchanges"].Outputs.Unpack(data)
	if err != nil {
		return KeyExchangesOutput{}, err
	}
	if len(vals) != 7 {
		return KeyExchangesOutput{}, fmt.Errorf("expected 7 values, got %d", len(vals))
	}
	jsonData0, err := json.Marshal(vals[0])
	if err != nil {
		return KeyExchangesOutput{}, fmt.Errorf("failed to marshal ABI result 0: %w", err)
	}

	var result0 [32]byte
	if err := json.Unmarshal(jsonData0, &result0); err != nil {
		return KeyExchangesOutput{}, fmt.Errorf("failed to unmarshal to [32]byte: %w", err)
	}
	jsonData1, err := json.Marshal(vals[1])
	if err != nil {
		return KeyExchangesOutput{}, fmt.Errorf("failed to marshal ABI result 1: %w", err)
	}

	var result1 []byte
	if err := json.Unmarshal(jsonData1, &result1); err != nil {
		return KeyExchangesOutput{}, fmt.Errorf("failed to unmarshal to []byte: %w", err)
	}
	jsonData2, err := json.Marshal(vals[2])
	if err != nil {
		return KeyExchangesOutput{}, fmt.Errorf("failed to marshal ABI result 2: %w", err)
	}

	var result2 common.Address
	if err := json.Unmarshal(jsonData2, &result2); err != nil {
		return KeyExchangesOutput{}, fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}
	jsonData3, err := json.Marshal(vals[3])
	if err != nil {
		return KeyExchangesOutput{}, fmt.Errorf("failed to marshal ABI result 3: %w", err)
	}

	var result3 *big.Int
	if err := json.Unmarshal(jsonData3, &result3); err != nil {
		return KeyExchangesOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData4, err := json.Marshal(vals[4])
	if err != nil {
		return KeyExchangesOutput{}, fmt.Errorf("failed to marshal ABI result 4: %w", err)
	}

	var result4 *big.Int
	if err := json.Unmarshal(jsonData4, &result4); err != nil {
		return KeyExchangesOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData5, err := json.Marshal(vals[5])
	if err != nil {
		return KeyExchangesOutput{}, fmt.Errorf("failed to marshal ABI result 5: %w", err)
	}

	var result5 bool
	if err := json.Unmarshal(jsonData5, &result5); err != nil {
		return KeyExchangesOutput{}, fmt.Errorf("failed to unmarshal to bool: %w", err)
	}
	jsonData6, err := json.Marshal(vals[6])
	if err != nil {
		return KeyExchangesOutput{}, fmt.Errorf("failed to marshal ABI result 6: %w", err)
	}

	var result6 uint8
	if err := json.Unmarshal(jsonData6, &result6); err != nil {
		return KeyExchangesOutput{}, fmt.Errorf("failed to unmarshal to uint8: %w", err)
	}

	return KeyExchangesOutput{
		KeyHash:           result0,
		EncryptedKey:      result1,
		IntendedRecipient: result2,
		CreatedAt:         result3,
		ExpiresAt:         result4,
		IsClaimed:         result5,
		ExchangeType:      result6,
	}, nil
}

func (c *Codec) EncodeKycInfoMethodCall(in KycInfoInput) ([]byte, error) {
	return c.abi.Pack("kycInfo", in.Arg0)
}

func (c *Codec) DecodeKycInfoMethodOutput(data []byte) (KycInfoOutput, error) {
	vals, err := c.abi.Methods["kycInfo"].Outputs.Unpack(data)
	if err != nil {
		return KycInfoOutput{}, err
	}
	if len(vals) != 5 {
		return KycInfoOutput{}, fmt.Errorf("expected 5 values, got %d", len(vals))
	}
	jsonData0, err := json.Marshal(vals[0])
	if err != nil {
		return KycInfoOutput{}, fmt.Errorf("failed to marshal ABI result 0: %w", err)
	}

	var result0 bool
	if err := json.Unmarshal(jsonData0, &result0); err != nil {
		return KycInfoOutput{}, fmt.Errorf("failed to unmarshal to bool: %w", err)
	}
	jsonData1, err := json.Marshal(vals[1])
	if err != nil {
		return KycInfoOutput{}, fmt.Errorf("failed to marshal ABI result 1: %w", err)
	}

	var result1 *big.Int
	if err := json.Unmarshal(jsonData1, &result1); err != nil {
		return KycInfoOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData2, err := json.Marshal(vals[2])
	if err != nil {
		return KycInfoOutput{}, fmt.Errorf("failed to marshal ABI result 2: %w", err)
	}

	var result2 [32]byte
	if err := json.Unmarshal(jsonData2, &result2); err != nil {
		return KycInfoOutput{}, fmt.Errorf("failed to unmarshal to [32]byte: %w", err)
	}
	jsonData3, err := json.Marshal(vals[3])
	if err != nil {
		return KycInfoOutput{}, fmt.Errorf("failed to marshal ABI result 3: %w", err)
	}

	var result3 uint8
	if err := json.Unmarshal(jsonData3, &result3); err != nil {
		return KycInfoOutput{}, fmt.Errorf("failed to unmarshal to uint8: %w", err)
	}
	jsonData4, err := json.Marshal(vals[4])
	if err != nil {
		return KycInfoOutput{}, fmt.Errorf("failed to marshal ABI result 4: %w", err)
	}

	var result4 *big.Int
	if err := json.Unmarshal(jsonData4, &result4); err != nil {
		return KycInfoOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return KycInfoOutput{
		IsVerified:        result0,
		VerifiedAt:        result1,
		VerificationHash:  result2,
		VerificationLevel: result3,
		ExpiryDate:        result4,
	}, nil
}

func (c *Codec) EncodeListingsMethodCall(in ListingsInput) ([]byte, error) {
	return c.abi.Pack("listings", in.Arg0)
}

func (c *Codec) DecodeListingsMethodOutput(data []byte) (ListingsOutput, error) {
	vals, err := c.abi.Methods["listings"].Outputs.Unpack(data)
	if err != nil {
		return ListingsOutput{}, err
	}
	if len(vals) != 8 {
		return ListingsOutput{}, fmt.Errorf("expected 8 values, got %d", len(vals))
	}
	jsonData0, err := json.Marshal(vals[0])
	if err != nil {
		return ListingsOutput{}, fmt.Errorf("failed to marshal ABI result 0: %w", err)
	}

	var result0 uint8
	if err := json.Unmarshal(jsonData0, &result0); err != nil {
		return ListingsOutput{}, fmt.Errorf("failed to unmarshal to uint8: %w", err)
	}
	jsonData1, err := json.Marshal(vals[1])
	if err != nil {
		return ListingsOutput{}, fmt.Errorf("failed to marshal ABI result 1: %w", err)
	}

	var result1 *big.Int
	if err := json.Unmarshal(jsonData1, &result1); err != nil {
		return ListingsOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData2, err := json.Marshal(vals[2])
	if err != nil {
		return ListingsOutput{}, fmt.Errorf("failed to marshal ABI result 2: %w", err)
	}

	var result2 common.Address
	if err := json.Unmarshal(jsonData2, &result2); err != nil {
		return ListingsOutput{}, fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}
	jsonData3, err := json.Marshal(vals[3])
	if err != nil {
		return ListingsOutput{}, fmt.Errorf("failed to marshal ABI result 3: %w", err)
	}

	var result3 bool
	if err := json.Unmarshal(jsonData3, &result3); err != nil {
		return ListingsOutput{}, fmt.Errorf("failed to unmarshal to bool: %w", err)
	}
	jsonData4, err := json.Marshal(vals[4])
	if err != nil {
		return ListingsOutput{}, fmt.Errorf("failed to marshal ABI result 4: %w", err)
	}

	var result4 common.Address
	if err := json.Unmarshal(jsonData4, &result4); err != nil {
		return ListingsOutput{}, fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}
	jsonData5, err := json.Marshal(vals[5])
	if err != nil {
		return ListingsOutput{}, fmt.Errorf("failed to marshal ABI result 5: %w", err)
	}

	var result5 *big.Int
	if err := json.Unmarshal(jsonData5, &result5); err != nil {
		return ListingsOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData6, err := json.Marshal(vals[6])
	if err != nil {
		return ListingsOutput{}, fmt.Errorf("failed to marshal ABI result 6: %w", err)
	}

	var result6 *big.Int
	if err := json.Unmarshal(jsonData6, &result6); err != nil {
		return ListingsOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData7, err := json.Marshal(vals[7])
	if err != nil {
		return ListingsOutput{}, fmt.Errorf("failed to marshal ABI result 7: %w", err)
	}

	var result7 uint8
	if err := json.Unmarshal(jsonData7, &result7); err != nil {
		return ListingsOutput{}, fmt.Errorf("failed to unmarshal to uint8: %w", err)
	}

	return ListingsOutput{
		ListingType:    result0,
		Price:          result1,
		PreferredToken: result2,
		IsPrivateSale:  result3,
		AllowedBuyer:   result4,
		CreatedAt:      result5,
		ExpiresAt:      result6,
		PlatformFee:    result7,
	}, nil
}

func (c *Codec) EncodeMinKYCLevelForHighValueMethodCall() ([]byte, error) {
	return c.abi.Pack("minKYCLevelForHighValue")
}

func (c *Codec) DecodeMinKYCLevelForHighValueMethodOutput(data []byte) (uint8, error) {
	vals, err := c.abi.Methods["minKYCLevelForHighValue"].Outputs.Unpack(data)
	if err != nil {
		return *new(uint8), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(uint8), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result uint8
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(uint8), fmt.Errorf("failed to unmarshal to uint8: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeMinKYCLevelForMintMethodCall() ([]byte, error) {
	return c.abi.Pack("minKYCLevelForMint")
}

func (c *Codec) DecodeMinKYCLevelForMintMethodOutput(data []byte) (uint8, error) {
	vals, err := c.abi.Methods["minKYCLevelForMint"].Outputs.Unpack(data)
	if err != nil {
		return *new(uint8), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(uint8), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result uint8
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(uint8), fmt.Errorf("failed to unmarshal to uint8: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeMintMethodCall(in MintInput) ([]byte, error) {
	return c.abi.Pack("mint", in.To, in.HouseId, in.DocumentHash, in.DocumentURI, in.StorageType, in.VerificationData)
}

func (c *Codec) DecodeMintMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["mint"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeMintingPausedMethodCall() ([]byte, error) {
	return c.abi.Pack("mintingPaused")
}

func (c *Codec) DecodeMintingPausedMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["mintingPaused"].Outputs.Unpack(data)
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

func (c *Codec) EncodeNameMethodCall() ([]byte, error) {
	return c.abi.Pack("name")
}

func (c *Codec) DecodeNameMethodOutput(data []byte) (string, error) {
	vals, err := c.abi.Methods["name"].Outputs.Unpack(data)
	if err != nil {
		return *new(string), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(string), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result string
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(string), fmt.Errorf("failed to unmarshal to string: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeNextTokenIdMethodCall() ([]byte, error) {
	return c.abi.Pack("nextTokenId")
}

func (c *Codec) DecodeNextTokenIdMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["nextTokenId"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeOpenRentalDisputeMethodCall(in OpenRentalDisputeInput) ([]byte, error) {
	return c.abi.Pack("openRentalDispute", in.TokenId, in.Reason)
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

func (c *Codec) EncodeOwnerOfMethodCall(in OwnerOfInput) ([]byte, error) {
	return c.abi.Pack("ownerOf", in.TokenId)
}

func (c *Codec) DecodeOwnerOfMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["ownerOf"].Outputs.Unpack(data)
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

func (c *Codec) EncodePausedMethodCall() ([]byte, error) {
	return c.abi.Pack("paused")
}

func (c *Codec) DecodePausedMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["paused"].Outputs.Unpack(data)
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

func (c *Codec) EncodePaymentsPausedMethodCall() ([]byte, error) {
	return c.abi.Pack("paymentsPaused")
}

func (c *Codec) DecodePaymentsPausedMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["paymentsPaused"].Outputs.Unpack(data)
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

func (c *Codec) EncodePendingRentalDepositsMethodCall(in PendingRentalDepositsInput) ([]byte, error) {
	return c.abi.Pack("pendingRentalDeposits", in.Arg0, in.Arg1)
}

func (c *Codec) DecodePendingRentalDepositsMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["pendingRentalDeposits"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodePriceFeedMethodCall() ([]byte, error) {
	return c.abi.Pack("priceFeed")
}

func (c *Codec) DecodePriceFeedMethodOutput(data []byte) (common.Address, error) {
	vals, err := c.abi.Methods["priceFeed"].Outputs.Unpack(data)
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

func (c *Codec) EncodeProxiableUUIDMethodCall() ([]byte, error) {
	return c.abi.Pack("proxiableUUID")
}

func (c *Codec) DecodeProxiableUUIDMethodOutput(data []byte) ([32]byte, error) {
	vals, err := c.abi.Methods["proxiableUUID"].Outputs.Unpack(data)
	if err != nil {
		return *new([32]byte), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new([32]byte), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result [32]byte
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new([32]byte), fmt.Errorf("failed to unmarshal to [32]byte: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeRecordBillPaymentMethodCall(in RecordBillPaymentInput) ([]byte, error) {
	return c.abi.Pack("recordBillPayment", in.TokenId, in.BillIndex, in.PaymentMethod, in.PaymentReference)
}

func (c *Codec) EncodeRenounceOwnershipMethodCall() ([]byte, error) {
	return c.abi.Pack("renounceOwnership")
}

func (c *Codec) EncodeRentalsMethodCall(in RentalsInput) ([]byte, error) {
	return c.abi.Pack("rentals", in.Arg0)
}

func (c *Codec) DecodeRentalsMethodOutput(data []byte) (RentalsOutput, error) {
	vals, err := c.abi.Methods["rentals"].Outputs.Unpack(data)
	if err != nil {
		return RentalsOutput{}, err
	}
	if len(vals) != 8 {
		return RentalsOutput{}, fmt.Errorf("expected 8 values, got %d", len(vals))
	}
	jsonData0, err := json.Marshal(vals[0])
	if err != nil {
		return RentalsOutput{}, fmt.Errorf("failed to marshal ABI result 0: %w", err)
	}

	var result0 common.Address
	if err := json.Unmarshal(jsonData0, &result0); err != nil {
		return RentalsOutput{}, fmt.Errorf("failed to unmarshal to common.Address: %w", err)
	}
	jsonData1, err := json.Marshal(vals[1])
	if err != nil {
		return RentalsOutput{}, fmt.Errorf("failed to marshal ABI result 1: %w", err)
	}

	var result1 *big.Int
	if err := json.Unmarshal(jsonData1, &result1); err != nil {
		return RentalsOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData2, err := json.Marshal(vals[2])
	if err != nil {
		return RentalsOutput{}, fmt.Errorf("failed to marshal ABI result 2: %w", err)
	}

	var result2 *big.Int
	if err := json.Unmarshal(jsonData2, &result2); err != nil {
		return RentalsOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData3, err := json.Marshal(vals[3])
	if err != nil {
		return RentalsOutput{}, fmt.Errorf("failed to marshal ABI result 3: %w", err)
	}

	var result3 *big.Int
	if err := json.Unmarshal(jsonData3, &result3); err != nil {
		return RentalsOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData4, err := json.Marshal(vals[4])
	if err != nil {
		return RentalsOutput{}, fmt.Errorf("failed to marshal ABI result 4: %w", err)
	}

	var result4 *big.Int
	if err := json.Unmarshal(jsonData4, &result4); err != nil {
		return RentalsOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData5, err := json.Marshal(vals[5])
	if err != nil {
		return RentalsOutput{}, fmt.Errorf("failed to marshal ABI result 5: %w", err)
	}

	var result5 bool
	if err := json.Unmarshal(jsonData5, &result5); err != nil {
		return RentalsOutput{}, fmt.Errorf("failed to unmarshal to bool: %w", err)
	}
	jsonData6, err := json.Marshal(vals[6])
	if err != nil {
		return RentalsOutput{}, fmt.Errorf("failed to marshal ABI result 6: %w", err)
	}

	var result6 [32]byte
	if err := json.Unmarshal(jsonData6, &result6); err != nil {
		return RentalsOutput{}, fmt.Errorf("failed to unmarshal to [32]byte: %w", err)
	}
	jsonData7, err := json.Marshal(vals[7])
	if err != nil {
		return RentalsOutput{}, fmt.Errorf("failed to marshal ABI result 7: %w", err)
	}

	var result7 uint8
	if err := json.Unmarshal(jsonData7, &result7); err != nil {
		return RentalsOutput{}, fmt.Errorf("failed to unmarshal to uint8: %w", err)
	}

	return RentalsOutput{
		Renter:                 result0,
		StartTime:              result1,
		EndTime:                result2,
		DepositAmount:          result3,
		MonthlyRent:            result4,
		IsActive:               result5,
		EncryptedAccessKeyHash: result6,
		DisputeStatus:          result7,
	}, nil
}

func (c *Codec) EncodeRentalsPausedMethodCall() ([]byte, error) {
	return c.abi.Pack("rentalsPaused")
}

func (c *Codec) DecodeRentalsPausedMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["rentalsPaused"].Outputs.Unpack(data)
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

func (c *Codec) EncodeResolveRentalDisputeMethodCall(in ResolveRentalDisputeInput) ([]byte, error) {
	return c.abi.Pack("resolveRentalDispute", in.TokenId, in.DepositToOwner, in.DepositToRenter)
}

func (c *Codec) EncodeRevokeCREWorkflowMethodCall(in RevokeCREWorkflowInput) ([]byte, error) {
	return c.abi.Pack("revokeCREWorkflow", in.Workflow)
}

func (c *Codec) EncodeRevokeKYCMethodCall(in RevokeKYCInput) ([]byte, error) {
	return c.abi.Pack("revokeKYC", in.User)
}

func (c *Codec) EncodeSafeTransferFromMethodCall(in SafeTransferFromInput) ([]byte, error) {
	return c.abi.Pack("safeTransferFrom", in.From, in.To, in.TokenId)
}

func (c *Codec) EncodeSafeTransferFrom0MethodCall(in SafeTransferFrom0Input) ([]byte, error) {
	return c.abi.Pack("safeTransferFrom0", in.From, in.To, in.TokenId, in.Data)
}

func (c *Codec) EncodeSalesPausedMethodCall() ([]byte, error) {
	return c.abi.Pack("salesPaused")
}

func (c *Codec) DecodeSalesPausedMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["salesPaused"].Outputs.Unpack(data)
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

func (c *Codec) EncodeSetApprovalForAllMethodCall(in SetApprovalForAllInput) ([]byte, error) {
	return c.abi.Pack("setApprovalForAll", in.Operator, in.Approved)
}

func (c *Codec) EncodeSetArbitratorMethodCall(in SetArbitratorInput) ([]byte, error) {
	return c.abi.Pack("setArbitrator", in.Arbitrator, in.IsArbitrator)
}

func (c *Codec) EncodeSetFeeRecipientMethodCall(in SetFeeRecipientInput) ([]byte, error) {
	return c.abi.Pack("setFeeRecipient", in.FeeRecipient)
}

func (c *Codec) EncodeSetKYCRequirementsMethodCall(in SetKYCRequirementsInput) ([]byte, error) {
	return c.abi.Pack("setKYCRequirements", in.MinLevelForMint, in.MinLevelForHighValue, in.HighValueThresholdUSD)
}

func (c *Codec) EncodeSetKYCVerificationMethodCall(in SetKYCVerificationInput) ([]byte, error) {
	return c.abi.Pack("setKYCVerification", in.User, in.Level, in.VerificationHash, in.ExpiryDate)
}

func (c *Codec) EncodeSetMintingPausedMethodCall(in SetMintingPausedInput) ([]byte, error) {
	return c.abi.Pack("setMintingPaused", in.Paused, in.Reason)
}

func (c *Codec) EncodeSetPaymentsPausedMethodCall(in SetPaymentsPausedInput) ([]byte, error) {
	return c.abi.Pack("setPaymentsPaused", in.Paused, in.Reason)
}

func (c *Codec) EncodeSetRentalsPausedMethodCall(in SetRentalsPausedInput) ([]byte, error) {
	return c.abi.Pack("setRentalsPaused", in.Paused, in.Reason)
}

func (c *Codec) EncodeSetSalesPausedMethodCall(in SetSalesPausedInput) ([]byte, error) {
	return c.abi.Pack("setSalesPaused", in.Paused, in.Reason)
}

func (c *Codec) EncodeSetTrustedBillProviderMethodCall(in SetTrustedBillProviderInput) ([]byte, error) {
	return c.abi.Pack("setTrustedBillProvider", in.Provider, in.Trusted)
}

func (c *Codec) EncodeSlashValidatorMethodCall(in SlashValidatorInput) ([]byte, error) {
	return c.abi.Pack("slashValidator", in.Validator, in.Amount, in.Reason)
}

func (c *Codec) EncodeStakeAsValidatorMethodCall() ([]byte, error) {
	return c.abi.Pack("stakeAsValidator")
}

func (c *Codec) EncodeStartRentalMethodCall(in StartRentalInput) ([]byte, error) {
	return c.abi.Pack("startRental", in.TokenId, in.Renter, in.DurationDays, in.DepositAmount, in.MonthlyRent, in.EncryptedAccessKey)
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

func (c *Codec) EncodeSymbolMethodCall() ([]byte, error) {
	return c.abi.Pack("symbol")
}

func (c *Codec) DecodeSymbolMethodOutput(data []byte) (string, error) {
	vals, err := c.abi.Methods["symbol"].Outputs.Unpack(data)
	if err != nil {
		return *new(string), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(string), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result string
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(string), fmt.Errorf("failed to unmarshal to string: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeTokenByIndexMethodCall(in TokenByIndexInput) ([]byte, error) {
	return c.abi.Pack("tokenByIndex", in.Index)
}

func (c *Codec) DecodeTokenByIndexMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["tokenByIndex"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeTokenOfOwnerByIndexMethodCall(in TokenOfOwnerByIndexInput) ([]byte, error) {
	return c.abi.Pack("tokenOfOwnerByIndex", in.Owner, in.Index)
}

func (c *Codec) DecodeTokenOfOwnerByIndexMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["tokenOfOwnerByIndex"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeTokenURIMethodCall(in TokenURIInput) ([]byte, error) {
	return c.abi.Pack("tokenURI", in.TokenId)
}

func (c *Codec) DecodeTokenURIMethodOutput(data []byte) (string, error) {
	vals, err := c.abi.Methods["tokenURI"].Outputs.Unpack(data)
	if err != nil {
		return *new(string), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(string), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result string
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(string), fmt.Errorf("failed to unmarshal to string: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeTotalFeesCollectedMethodCall() ([]byte, error) {
	return c.abi.Pack("totalFeesCollected")
}

func (c *Codec) DecodeTotalFeesCollectedMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["totalFeesCollected"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeTotalSupplyMethodCall() ([]byte, error) {
	return c.abi.Pack("totalSupply")
}

func (c *Codec) DecodeTotalSupplyMethodOutput(data []byte) (*big.Int, error) {
	vals, err := c.abi.Methods["totalSupply"].Outputs.Unpack(data)
	if err != nil {
		return *new(*big.Int), err
	}
	jsonData, err := json.Marshal(vals[0])
	if err != nil {
		return *new(*big.Int), fmt.Errorf("failed to marshal ABI result: %w", err)
	}

	var result *big.Int
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return *new(*big.Int), fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}

	return result, nil
}

func (c *Codec) EncodeTransferFromMethodCall(in TransferFromInput) ([]byte, error) {
	return c.abi.Pack("transferFrom", in.From, in.To, in.TokenId)
}

func (c *Codec) EncodeTransferOwnershipMethodCall(in TransferOwnershipInput) ([]byte, error) {
	return c.abi.Pack("transferOwnership", in.NewOwner)
}

func (c *Codec) EncodeTrustedBillProvidersMethodCall(in TrustedBillProvidersInput) ([]byte, error) {
	return c.abi.Pack("trustedBillProviders", in.Arg0)
}

func (c *Codec) DecodeTrustedBillProvidersMethodOutput(data []byte) (bool, error) {
	vals, err := c.abi.Methods["trustedBillProviders"].Outputs.Unpack(data)
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

func (c *Codec) EncodeUnstakeMethodCall(in UnstakeInput) ([]byte, error) {
	return c.abi.Pack("unstake", in.Amount)
}

func (c *Codec) EncodeUpgradeToAndCallMethodCall(in UpgradeToAndCallInput) ([]byte, error) {
	return c.abi.Pack("upgradeToAndCall", in.NewImplementation, in.Data)
}

func (c *Codec) EncodeValidatorsMethodCall(in ValidatorsInput) ([]byte, error) {
	return c.abi.Pack("validators", in.Arg0)
}

func (c *Codec) DecodeValidatorsMethodOutput(data []byte) (ValidatorsOutput, error) {
	vals, err := c.abi.Methods["validators"].Outputs.Unpack(data)
	if err != nil {
		return ValidatorsOutput{}, err
	}
	if len(vals) != 6 {
		return ValidatorsOutput{}, fmt.Errorf("expected 6 values, got %d", len(vals))
	}
	jsonData0, err := json.Marshal(vals[0])
	if err != nil {
		return ValidatorsOutput{}, fmt.Errorf("failed to marshal ABI result 0: %w", err)
	}

	var result0 *big.Int
	if err := json.Unmarshal(jsonData0, &result0); err != nil {
		return ValidatorsOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData1, err := json.Marshal(vals[1])
	if err != nil {
		return ValidatorsOutput{}, fmt.Errorf("failed to marshal ABI result 1: %w", err)
	}

	var result1 *big.Int
	if err := json.Unmarshal(jsonData1, &result1); err != nil {
		return ValidatorsOutput{}, fmt.Errorf("failed to unmarshal to *big.Int: %w", err)
	}
	jsonData2, err := json.Marshal(vals[2])
	if err != nil {
		return ValidatorsOutput{}, fmt.Errorf("failed to marshal ABI result 2: %w", err)
	}

	var result2 uint8
	if err := json.Unmarshal(jsonData2, &result2); err != nil {
		return ValidatorsOutput{}, fmt.Errorf("failed to unmarshal to uint8: %w", err)
	}
	jsonData3, err := json.Marshal(vals[3])
	if err != nil {
		return ValidatorsOutput{}, fmt.Errorf("failed to marshal ABI result 3: %w", err)
	}

	var result3 bool
	if err := json.Unmarshal(jsonData3, &result3); err != nil {
		return ValidatorsOutput{}, fmt.Errorf("failed to unmarshal to bool: %w", err)
	}
	jsonData4, err := json.Marshal(vals[4])
	if err != nil {
		return ValidatorsOutput{}, fmt.Errorf("failed to marshal ABI result 4: %w", err)
	}

	var result4 uint8
	if err := json.Unmarshal(jsonData4, &result4); err != nil {
		return ValidatorsOutput{}, fmt.Errorf("failed to unmarshal to uint8: %w", err)
	}
	jsonData5, err := json.Marshal(vals[5])
	if err != nil {
		return ValidatorsOutput{}, fmt.Errorf("failed to marshal ABI result 5: %w", err)
	}

	var result5 uint8
	if err := json.Unmarshal(jsonData5, &result5); err != nil {
		return ValidatorsOutput{}, fmt.Errorf("failed to unmarshal to uint8: %w", err)
	}

	return ValidatorsOutput{
		StakedAmount:          result0,
		StakedAt:              result1,
		Reputation:            result2,
		IsSlashed:             result3,
		SuccessfulValidations: result4,
		FailedValidations:     result5,
	}, nil
}

func (c *Codec) EncodeWithdrawRentalDepositMethodCall(in WithdrawRentalDepositInput) ([]byte, error) {
	return c.abi.Pack("withdrawRentalDeposit", in.TokenId)
}

func (c *Codec) EncodeBillStruct(in Bill) ([]byte, error) {
	tupleType, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "billType", Type: "string"},
			{Name: "amount", Type: "uint96"},
			{Name: "dueDate", Type: "uint48"},
			{Name: "paidAt", Type: "uint48"},
			{Name: "status", Type: "uint8"},
			{Name: "paymentReference", Type: "bytes32"},
			{Name: "isRecurring", Type: "bool"},
			{Name: "provider", Type: "address"},
			{Name: "recurrenceInterval", Type: "uint8"},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tuple type for Bill: %w", err)
	}
	args := abi.Arguments{
		{Name: "bill", Type: tupleType},
	}

	return args.Pack(in)
}
func (c *Codec) EncodeHouseStruct(in House) ([]byte, error) {
	tupleType, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "houseId", Type: "string"},
			{Name: "documentHash", Type: "bytes32"},
			{Name: "documentURI", Type: "string"},
			{Name: "storageType", Type: "uint8"},
			{Name: "originalOwner", Type: "address"},
			{Name: "mintedAt", Type: "uint48"},
			{Name: "isVerified", Type: "bool"},
			{Name: "documentCount", Type: "uint8"},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tuple type for House: %w", err)
	}
	args := abi.Arguments{
		{Name: "house", Type: tupleType},
	}

	return args.Pack(in)
}
func (c *Codec) EncodeListingStruct(in Listing) ([]byte, error) {
	tupleType, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "listingType", Type: "uint8"},
			{Name: "price", Type: "uint96"},
			{Name: "preferredToken", Type: "address"},
			{Name: "isPrivateSale", Type: "bool"},
			{Name: "allowedBuyer", Type: "address"},
			{Name: "createdAt", Type: "uint48"},
			{Name: "expiresAt", Type: "uint48"},
			{Name: "platformFee", Type: "uint8"},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tuple type for Listing: %w", err)
	}
	args := abi.Arguments{
		{Name: "listing", Type: tupleType},
	}

	return args.Pack(in)
}
func (c *Codec) EncodeRentalAgreementStruct(in RentalAgreement) ([]byte, error) {
	tupleType, err := abi.NewType(
		"tuple", "",
		[]abi.ArgumentMarshaling{
			{Name: "renter", Type: "address"},
			{Name: "startTime", Type: "uint48"},
			{Name: "endTime", Type: "uint48"},
			{Name: "depositAmount", Type: "uint96"},
			{Name: "monthlyRent", Type: "uint96"},
			{Name: "isActive", Type: "bool"},
			{Name: "encryptedAccessKeyHash", Type: "bytes32"},
			{Name: "disputeStatus", Type: "uint8"},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create tuple type for RentalAgreement: %w", err)
	}
	args := abi.Arguments{
		{Name: "rentalAgreement", Type: tupleType},
	}

	return args.Pack(in)
}

func (c *Codec) ApprovalLogHash() []byte {
	return c.abi.Events["Approval"].ID.Bytes()
}

func (c *Codec) EncodeApprovalTopics(
	evt abi.Event,
	values []ApprovalTopics,
) ([]*evm.TopicValues, error) {
	var ownerRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Owner).IsZero() {
			ownerRule = append(ownerRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Owner)
		if err != nil {
			return nil, err
		}
		ownerRule = append(ownerRule, fieldVal)
	}
	var approvedRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Approved).IsZero() {
			approvedRule = append(approvedRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Approved)
		if err != nil {
			return nil, err
		}
		approvedRule = append(approvedRule, fieldVal)
	}
	var tokenIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.TokenId).IsZero() {
			tokenIdRule = append(tokenIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[2], v.TokenId)
		if err != nil {
			return nil, err
		}
		tokenIdRule = append(tokenIdRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		ownerRule,
		approvedRule,
		tokenIdRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeApproval decodes a log into a Approval struct.
func (c *Codec) DecodeApproval(log *evm.Log) (*ApprovalDecoded, error) {
	event := new(ApprovalDecoded)
	if err := c.abi.UnpackIntoInterface(event, "Approval", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["Approval"].Inputs {
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

func (c *Codec) ApprovalForAllLogHash() []byte {
	return c.abi.Events["ApprovalForAll"].ID.Bytes()
}

func (c *Codec) EncodeApprovalForAllTopics(
	evt abi.Event,
	values []ApprovalForAllTopics,
) ([]*evm.TopicValues, error) {
	var ownerRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Owner).IsZero() {
			ownerRule = append(ownerRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Owner)
		if err != nil {
			return nil, err
		}
		ownerRule = append(ownerRule, fieldVal)
	}
	var operatorRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Operator).IsZero() {
			operatorRule = append(operatorRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Operator)
		if err != nil {
			return nil, err
		}
		operatorRule = append(operatorRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		ownerRule,
		operatorRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeApprovalForAll decodes a log into a ApprovalForAll struct.
func (c *Codec) DecodeApprovalForAll(log *evm.Log) (*ApprovalForAllDecoded, error) {
	event := new(ApprovalForAllDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ApprovalForAll", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ApprovalForAll"].Inputs {
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

func (c *Codec) BillCreatedLogHash() []byte {
	return c.abi.Events["BillCreated"].ID.Bytes()
}

func (c *Codec) EncodeBillCreatedTopics(
	evt abi.Event,
	values []BillCreatedTopics,
) ([]*evm.TopicValues, error) {
	var tokenIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.TokenId).IsZero() {
			tokenIdRule = append(tokenIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.TokenId)
		if err != nil {
			return nil, err
		}
		tokenIdRule = append(tokenIdRule, fieldVal)
	}
	var billIndexRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.BillIndex).IsZero() {
			billIndexRule = append(billIndexRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.BillIndex)
		if err != nil {
			return nil, err
		}
		billIndexRule = append(billIndexRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		tokenIdRule,
		billIndexRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeBillCreated decodes a log into a BillCreated struct.
func (c *Codec) DecodeBillCreated(log *evm.Log) (*BillCreatedDecoded, error) {
	event := new(BillCreatedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "BillCreated", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["BillCreated"].Inputs {
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

func (c *Codec) BillDisputedLogHash() []byte {
	return c.abi.Events["BillDisputed"].ID.Bytes()
}

func (c *Codec) EncodeBillDisputedTopics(
	evt abi.Event,
	values []BillDisputedTopics,
) ([]*evm.TopicValues, error) {
	var tokenIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.TokenId).IsZero() {
			tokenIdRule = append(tokenIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.TokenId)
		if err != nil {
			return nil, err
		}
		tokenIdRule = append(tokenIdRule, fieldVal)
	}
	var billIndexRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.BillIndex).IsZero() {
			billIndexRule = append(billIndexRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.BillIndex)
		if err != nil {
			return nil, err
		}
		billIndexRule = append(billIndexRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		tokenIdRule,
		billIndexRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeBillDisputed decodes a log into a BillDisputed struct.
func (c *Codec) DecodeBillDisputed(log *evm.Log) (*BillDisputedDecoded, error) {
	event := new(BillDisputedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "BillDisputed", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["BillDisputed"].Inputs {
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

func (c *Codec) BillPaidLogHash() []byte {
	return c.abi.Events["BillPaid"].ID.Bytes()
}

func (c *Codec) EncodeBillPaidTopics(
	evt abi.Event,
	values []BillPaidTopics,
) ([]*evm.TopicValues, error) {
	var tokenIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.TokenId).IsZero() {
			tokenIdRule = append(tokenIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.TokenId)
		if err != nil {
			return nil, err
		}
		tokenIdRule = append(tokenIdRule, fieldVal)
	}
	var billIndexRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.BillIndex).IsZero() {
			billIndexRule = append(billIndexRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.BillIndex)
		if err != nil {
			return nil, err
		}
		billIndexRule = append(billIndexRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		tokenIdRule,
		billIndexRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeBillPaid decodes a log into a BillPaid struct.
func (c *Codec) DecodeBillPaid(log *evm.Log) (*BillPaidDecoded, error) {
	event := new(BillPaidDecoded)
	if err := c.abi.UnpackIntoInterface(event, "BillPaid", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["BillPaid"].Inputs {
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

func (c *Codec) CircuitBreakerTriggeredLogHash() []byte {
	return c.abi.Events["CircuitBreakerTriggered"].ID.Bytes()
}

func (c *Codec) EncodeCircuitBreakerTriggeredTopics(
	evt abi.Event,
	values []CircuitBreakerTriggeredTopics,
) ([]*evm.TopicValues, error) {

	rawTopics, err := abi.MakeTopics()
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeCircuitBreakerTriggered decodes a log into a CircuitBreakerTriggered struct.
func (c *Codec) DecodeCircuitBreakerTriggered(log *evm.Log) (*CircuitBreakerTriggeredDecoded, error) {
	event := new(CircuitBreakerTriggeredDecoded)
	if err := c.abi.UnpackIntoInterface(event, "CircuitBreakerTriggered", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["CircuitBreakerTriggered"].Inputs {
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

func (c *Codec) EmergencyActionLogHash() []byte {
	return c.abi.Events["EmergencyAction"].ID.Bytes()
}

func (c *Codec) EncodeEmergencyActionTopics(
	evt abi.Event,
	values []EmergencyActionTopics,
) ([]*evm.TopicValues, error) {

	rawTopics, err := abi.MakeTopics()
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeEmergencyAction decodes a log into a EmergencyAction struct.
func (c *Codec) DecodeEmergencyAction(log *evm.Log) (*EmergencyActionDecoded, error) {
	event := new(EmergencyActionDecoded)
	if err := c.abi.UnpackIntoInterface(event, "EmergencyAction", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["EmergencyAction"].Inputs {
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

func (c *Codec) HouseListedLogHash() []byte {
	return c.abi.Events["HouseListed"].ID.Bytes()
}

func (c *Codec) EncodeHouseListedTopics(
	evt abi.Event,
	values []HouseListedTopics,
) ([]*evm.TopicValues, error) {
	var tokenIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.TokenId).IsZero() {
			tokenIdRule = append(tokenIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.TokenId)
		if err != nil {
			return nil, err
		}
		tokenIdRule = append(tokenIdRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		tokenIdRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeHouseListed decodes a log into a HouseListed struct.
func (c *Codec) DecodeHouseListed(log *evm.Log) (*HouseListedDecoded, error) {
	event := new(HouseListedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "HouseListed", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["HouseListed"].Inputs {
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

func (c *Codec) HouseMintedLogHash() []byte {
	return c.abi.Events["HouseMinted"].ID.Bytes()
}

func (c *Codec) EncodeHouseMintedTopics(
	evt abi.Event,
	values []HouseMintedTopics,
) ([]*evm.TopicValues, error) {
	var tokenIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.TokenId).IsZero() {
			tokenIdRule = append(tokenIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.TokenId)
		if err != nil {
			return nil, err
		}
		tokenIdRule = append(tokenIdRule, fieldVal)
	}
	var ownerRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Owner).IsZero() {
			ownerRule = append(ownerRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Owner)
		if err != nil {
			return nil, err
		}
		ownerRule = append(ownerRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		tokenIdRule,
		ownerRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeHouseMinted decodes a log into a HouseMinted struct.
func (c *Codec) DecodeHouseMinted(log *evm.Log) (*HouseMintedDecoded, error) {
	event := new(HouseMintedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "HouseMinted", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["HouseMinted"].Inputs {
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

func (c *Codec) HouseSoldLogHash() []byte {
	return c.abi.Events["HouseSold"].ID.Bytes()
}

func (c *Codec) EncodeHouseSoldTopics(
	evt abi.Event,
	values []HouseSoldTopics,
) ([]*evm.TopicValues, error) {
	var tokenIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.TokenId).IsZero() {
			tokenIdRule = append(tokenIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.TokenId)
		if err != nil {
			return nil, err
		}
		tokenIdRule = append(tokenIdRule, fieldVal)
	}
	var sellerRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Seller).IsZero() {
			sellerRule = append(sellerRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Seller)
		if err != nil {
			return nil, err
		}
		sellerRule = append(sellerRule, fieldVal)
	}
	var buyerRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Buyer).IsZero() {
			buyerRule = append(buyerRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[2], v.Buyer)
		if err != nil {
			return nil, err
		}
		buyerRule = append(buyerRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		tokenIdRule,
		sellerRule,
		buyerRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeHouseSold decodes a log into a HouseSold struct.
func (c *Codec) DecodeHouseSold(log *evm.Log) (*HouseSoldDecoded, error) {
	event := new(HouseSoldDecoded)
	if err := c.abi.UnpackIntoInterface(event, "HouseSold", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["HouseSold"].Inputs {
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

func (c *Codec) InitializedLogHash() []byte {
	return c.abi.Events["Initialized"].ID.Bytes()
}

func (c *Codec) EncodeInitializedTopics(
	evt abi.Event,
	values []InitializedTopics,
) ([]*evm.TopicValues, error) {

	rawTopics, err := abi.MakeTopics()
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeInitialized decodes a log into a Initialized struct.
func (c *Codec) DecodeInitialized(log *evm.Log) (*InitializedDecoded, error) {
	event := new(InitializedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "Initialized", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["Initialized"].Inputs {
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

func (c *Codec) KYCVerifiedLogHash() []byte {
	return c.abi.Events["KYCVerified"].ID.Bytes()
}

func (c *Codec) EncodeKYCVerifiedTopics(
	evt abi.Event,
	values []KYCVerifiedTopics,
) ([]*evm.TopicValues, error) {
	var userRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.User).IsZero() {
			userRule = append(userRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.User)
		if err != nil {
			return nil, err
		}
		userRule = append(userRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		userRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeKYCVerified decodes a log into a KYCVerified struct.
func (c *Codec) DecodeKYCVerified(log *evm.Log) (*KYCVerifiedDecoded, error) {
	event := new(KYCVerifiedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "KYCVerified", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["KYCVerified"].Inputs {
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

func (c *Codec) KeyClaimedLogHash() []byte {
	return c.abi.Events["KeyClaimed"].ID.Bytes()
}

func (c *Codec) EncodeKeyClaimedTopics(
	evt abi.Event,
	values []KeyClaimedTopics,
) ([]*evm.TopicValues, error) {
	var keyHashRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.KeyHash).IsZero() {
			keyHashRule = append(keyHashRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.KeyHash)
		if err != nil {
			return nil, err
		}
		keyHashRule = append(keyHashRule, fieldVal)
	}
	var claimantRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Claimant).IsZero() {
			claimantRule = append(claimantRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Claimant)
		if err != nil {
			return nil, err
		}
		claimantRule = append(claimantRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		keyHashRule,
		claimantRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeKeyClaimed decodes a log into a KeyClaimed struct.
func (c *Codec) DecodeKeyClaimed(log *evm.Log) (*KeyClaimedDecoded, error) {
	event := new(KeyClaimedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "KeyClaimed", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["KeyClaimed"].Inputs {
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

func (c *Codec) KeyExchangeCreatedLogHash() []byte {
	return c.abi.Events["KeyExchangeCreated"].ID.Bytes()
}

func (c *Codec) EncodeKeyExchangeCreatedTopics(
	evt abi.Event,
	values []KeyExchangeCreatedTopics,
) ([]*evm.TopicValues, error) {
	var keyHashRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.KeyHash).IsZero() {
			keyHashRule = append(keyHashRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.KeyHash)
		if err != nil {
			return nil, err
		}
		keyHashRule = append(keyHashRule, fieldVal)
	}
	var tokenIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.TokenId).IsZero() {
			tokenIdRule = append(tokenIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.TokenId)
		if err != nil {
			return nil, err
		}
		tokenIdRule = append(tokenIdRule, fieldVal)
	}
	var recipientRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Recipient).IsZero() {
			recipientRule = append(recipientRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[2], v.Recipient)
		if err != nil {
			return nil, err
		}
		recipientRule = append(recipientRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		keyHashRule,
		tokenIdRule,
		recipientRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeKeyExchangeCreated decodes a log into a KeyExchangeCreated struct.
func (c *Codec) DecodeKeyExchangeCreated(log *evm.Log) (*KeyExchangeCreatedDecoded, error) {
	event := new(KeyExchangeCreatedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "KeyExchangeCreated", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["KeyExchangeCreated"].Inputs {
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

func (c *Codec) PausedLogHash() []byte {
	return c.abi.Events["Paused"].ID.Bytes()
}

func (c *Codec) EncodePausedTopics(
	evt abi.Event,
	values []PausedTopics,
) ([]*evm.TopicValues, error) {

	rawTopics, err := abi.MakeTopics()
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodePaused decodes a log into a Paused struct.
func (c *Codec) DecodePaused(log *evm.Log) (*PausedDecoded, error) {
	event := new(PausedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "Paused", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["Paused"].Inputs {
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

func (c *Codec) RentalDepositReceivedLogHash() []byte {
	return c.abi.Events["RentalDepositReceived"].ID.Bytes()
}

func (c *Codec) EncodeRentalDepositReceivedTopics(
	evt abi.Event,
	values []RentalDepositReceivedTopics,
) ([]*evm.TopicValues, error) {
	var tokenIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.TokenId).IsZero() {
			tokenIdRule = append(tokenIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.TokenId)
		if err != nil {
			return nil, err
		}
		tokenIdRule = append(tokenIdRule, fieldVal)
	}
	var renterRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Renter).IsZero() {
			renterRule = append(renterRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Renter)
		if err != nil {
			return nil, err
		}
		renterRule = append(renterRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		tokenIdRule,
		renterRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeRentalDepositReceived decodes a log into a RentalDepositReceived struct.
func (c *Codec) DecodeRentalDepositReceived(log *evm.Log) (*RentalDepositReceivedDecoded, error) {
	event := new(RentalDepositReceivedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "RentalDepositReceived", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["RentalDepositReceived"].Inputs {
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

func (c *Codec) RentalDepositWithdrawnLogHash() []byte {
	return c.abi.Events["RentalDepositWithdrawn"].ID.Bytes()
}

func (c *Codec) EncodeRentalDepositWithdrawnTopics(
	evt abi.Event,
	values []RentalDepositWithdrawnTopics,
) ([]*evm.TopicValues, error) {
	var tokenIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.TokenId).IsZero() {
			tokenIdRule = append(tokenIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.TokenId)
		if err != nil {
			return nil, err
		}
		tokenIdRule = append(tokenIdRule, fieldVal)
	}
	var renterRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Renter).IsZero() {
			renterRule = append(renterRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Renter)
		if err != nil {
			return nil, err
		}
		renterRule = append(renterRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		tokenIdRule,
		renterRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeRentalDepositWithdrawn decodes a log into a RentalDepositWithdrawn struct.
func (c *Codec) DecodeRentalDepositWithdrawn(log *evm.Log) (*RentalDepositWithdrawnDecoded, error) {
	event := new(RentalDepositWithdrawnDecoded)
	if err := c.abi.UnpackIntoInterface(event, "RentalDepositWithdrawn", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["RentalDepositWithdrawn"].Inputs {
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

func (c *Codec) RentalEndedLogHash() []byte {
	return c.abi.Events["RentalEnded"].ID.Bytes()
}

func (c *Codec) EncodeRentalEndedTopics(
	evt abi.Event,
	values []RentalEndedTopics,
) ([]*evm.TopicValues, error) {
	var tokenIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.TokenId).IsZero() {
			tokenIdRule = append(tokenIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.TokenId)
		if err != nil {
			return nil, err
		}
		tokenIdRule = append(tokenIdRule, fieldVal)
	}
	var renterRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Renter).IsZero() {
			renterRule = append(renterRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Renter)
		if err != nil {
			return nil, err
		}
		renterRule = append(renterRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		tokenIdRule,
		renterRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeRentalEnded decodes a log into a RentalEnded struct.
func (c *Codec) DecodeRentalEnded(log *evm.Log) (*RentalEndedDecoded, error) {
	event := new(RentalEndedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "RentalEnded", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["RentalEnded"].Inputs {
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

func (c *Codec) RentalStartedLogHash() []byte {
	return c.abi.Events["RentalStarted"].ID.Bytes()
}

func (c *Codec) EncodeRentalStartedTopics(
	evt abi.Event,
	values []RentalStartedTopics,
) ([]*evm.TopicValues, error) {
	var tokenIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.TokenId).IsZero() {
			tokenIdRule = append(tokenIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.TokenId)
		if err != nil {
			return nil, err
		}
		tokenIdRule = append(tokenIdRule, fieldVal)
	}
	var renterRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Renter).IsZero() {
			renterRule = append(renterRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.Renter)
		if err != nil {
			return nil, err
		}
		renterRule = append(renterRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		tokenIdRule,
		renterRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeRentalStarted decodes a log into a RentalStarted struct.
func (c *Codec) DecodeRentalStarted(log *evm.Log) (*RentalStartedDecoded, error) {
	event := new(RentalStartedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "RentalStarted", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["RentalStarted"].Inputs {
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

func (c *Codec) TransferLogHash() []byte {
	return c.abi.Events["Transfer"].ID.Bytes()
}

func (c *Codec) EncodeTransferTopics(
	evt abi.Event,
	values []TransferTopics,
) ([]*evm.TopicValues, error) {
	var fromRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.From).IsZero() {
			fromRule = append(fromRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.From)
		if err != nil {
			return nil, err
		}
		fromRule = append(fromRule, fieldVal)
	}
	var toRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.To).IsZero() {
			toRule = append(toRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[1], v.To)
		if err != nil {
			return nil, err
		}
		toRule = append(toRule, fieldVal)
	}
	var tokenIdRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.TokenId).IsZero() {
			tokenIdRule = append(tokenIdRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[2], v.TokenId)
		if err != nil {
			return nil, err
		}
		tokenIdRule = append(tokenIdRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		fromRule,
		toRule,
		tokenIdRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeTransfer decodes a log into a Transfer struct.
func (c *Codec) DecodeTransfer(log *evm.Log) (*TransferDecoded, error) {
	event := new(TransferDecoded)
	if err := c.abi.UnpackIntoInterface(event, "Transfer", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["Transfer"].Inputs {
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

func (c *Codec) UnpausedLogHash() []byte {
	return c.abi.Events["Unpaused"].ID.Bytes()
}

func (c *Codec) EncodeUnpausedTopics(
	evt abi.Event,
	values []UnpausedTopics,
) ([]*evm.TopicValues, error) {

	rawTopics, err := abi.MakeTopics()
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeUnpaused decodes a log into a Unpaused struct.
func (c *Codec) DecodeUnpaused(log *evm.Log) (*UnpausedDecoded, error) {
	event := new(UnpausedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "Unpaused", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["Unpaused"].Inputs {
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

func (c *Codec) UpgradedLogHash() []byte {
	return c.abi.Events["Upgraded"].ID.Bytes()
}

func (c *Codec) EncodeUpgradedTopics(
	evt abi.Event,
	values []UpgradedTopics,
) ([]*evm.TopicValues, error) {
	var implementationRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Implementation).IsZero() {
			implementationRule = append(implementationRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Implementation)
		if err != nil {
			return nil, err
		}
		implementationRule = append(implementationRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		implementationRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeUpgraded decodes a log into a Upgraded struct.
func (c *Codec) DecodeUpgraded(log *evm.Log) (*UpgradedDecoded, error) {
	event := new(UpgradedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "Upgraded", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["Upgraded"].Inputs {
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

func (c *Codec) ValidatorSlashedLogHash() []byte {
	return c.abi.Events["ValidatorSlashed"].ID.Bytes()
}

func (c *Codec) EncodeValidatorSlashedTopics(
	evt abi.Event,
	values []ValidatorSlashedTopics,
) ([]*evm.TopicValues, error) {
	var validatorRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Validator).IsZero() {
			validatorRule = append(validatorRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Validator)
		if err != nil {
			return nil, err
		}
		validatorRule = append(validatorRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		validatorRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeValidatorSlashed decodes a log into a ValidatorSlashed struct.
func (c *Codec) DecodeValidatorSlashed(log *evm.Log) (*ValidatorSlashedDecoded, error) {
	event := new(ValidatorSlashedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ValidatorSlashed", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ValidatorSlashed"].Inputs {
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

func (c *Codec) ValidatorStakedLogHash() []byte {
	return c.abi.Events["ValidatorStaked"].ID.Bytes()
}

func (c *Codec) EncodeValidatorStakedTopics(
	evt abi.Event,
	values []ValidatorStakedTopics,
) ([]*evm.TopicValues, error) {
	var validatorRule []interface{}
	for _, v := range values {
		if reflect.ValueOf(v.Validator).IsZero() {
			validatorRule = append(validatorRule, common.Hash{})
			continue
		}
		fieldVal, err := bindings.PrepareTopicArg(evt.Inputs[0], v.Validator)
		if err != nil {
			return nil, err
		}
		validatorRule = append(validatorRule, fieldVal)
	}

	rawTopics, err := abi.MakeTopics(
		validatorRule,
	)
	if err != nil {
		return nil, err
	}

	return bindings.PrepareTopics(rawTopics, evt.ID.Bytes()), nil
}

// DecodeValidatorStaked decodes a log into a ValidatorStaked struct.
func (c *Codec) DecodeValidatorStaked(log *evm.Log) (*ValidatorStakedDecoded, error) {
	event := new(ValidatorStakedDecoded)
	if err := c.abi.UnpackIntoInterface(event, "ValidatorStaked", log.Data); err != nil {
		return nil, err
	}
	var indexed abi.Arguments
	for _, arg := range c.abi.Events["ValidatorStaked"].Inputs {
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

func (c HouseRWA) BPSDENOMINATOR(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeBPSDENOMINATORMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeBPSDENOMINATORMethodOutput(response.Data)
	})

}

func (c HouseRWA) KEYEXCHANGEEXPIRY(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeKEYEXCHANGEEXPIRYMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeKEYEXCHANGEEXPIRYMethodOutput(response.Data)
	})

}

func (c HouseRWA) MAXBILLSPERHOUSE(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeMAXBILLSPERHOUSEMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeMAXBILLSPERHOUSEMethodOutput(response.Data)
	})

}

func (c HouseRWA) MAXDOCUMENTSPERHOUSE(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeMAXDOCUMENTSPERHOUSEMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeMAXDOCUMENTSPERHOUSEMethodOutput(response.Data)
	})

}

func (c HouseRWA) MAXRENTALDURATION(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeMAXRENTALDURATIONMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeMAXRENTALDURATIONMethodOutput(response.Data)
	})

}

func (c HouseRWA) MINVALIDATORSTAKE(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeMINVALIDATORSTAKEMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeMINVALIDATORSTAKEMethodOutput(response.Data)
	})

}

func (c HouseRWA) PROTOCOLFEEBPS(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodePROTOCOLFEEBPSMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodePROTOCOLFEEBPSMethodOutput(response.Data)
	})

}

func (c HouseRWA) UPGRADEINTERFACEVERSION(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[string] {
	calldata, err := c.Codec.EncodeUPGRADEINTERFACEVERSIONMethodCall()
	if err != nil {
		return cre.PromiseFromResult[string](*new(string), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (string, error) {
		return c.Codec.DecodeUPGRADEINTERFACEVERSIONMethodOutput(response.Data)
	})

}

func (c HouseRWA) Arbitrators(
	runtime cre.Runtime,
	args ArbitratorsInput,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeArbitratorsMethodCall(args)
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
		return c.Codec.DecodeArbitratorsMethodOutput(response.Data)
	})

}

func (c HouseRWA) AuthorizedCREWorkflows(
	runtime cre.Runtime,
	args AuthorizedCREWorkflowsInput,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeAuthorizedCREWorkflowsMethodCall(args)
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
		return c.Codec.DecodeAuthorizedCREWorkflowsMethodOutput(response.Data)
	})

}

func (c HouseRWA) BalanceOf(
	runtime cre.Runtime,
	args BalanceOfInput,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeBalanceOfMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeBalanceOfMethodOutput(response.Data)
	})

}

func (c HouseRWA) FeeRecipient(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeFeeRecipientMethodCall()
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
		return c.Codec.DecodeFeeRecipientMethodOutput(response.Data)
	})

}

func (c HouseRWA) GetActiveRental(
	runtime cre.Runtime,
	args GetActiveRentalInput,
	blockNumber *big.Int,
) cre.Promise[RentalAgreement] {
	calldata, err := c.Codec.EncodeGetActiveRentalMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[RentalAgreement](*new(RentalAgreement), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (RentalAgreement, error) {
		return c.Codec.DecodeGetActiveRentalMethodOutput(response.Data)
	})

}

func (c HouseRWA) GetApproved(
	runtime cre.Runtime,
	args GetApprovedInput,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeGetApprovedMethodCall(args)
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
		return c.Codec.DecodeGetApprovedMethodOutput(response.Data)
	})

}

func (c HouseRWA) GetBills(
	runtime cre.Runtime,
	args GetBillsInput,
	blockNumber *big.Int,
) cre.Promise[[]Bill] {
	calldata, err := c.Codec.EncodeGetBillsMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[[]Bill](*new([]Bill), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) ([]Bill, error) {
		return c.Codec.DecodeGetBillsMethodOutput(response.Data)
	})

}

func (c HouseRWA) GetHouseDetails(
	runtime cre.Runtime,
	args GetHouseDetailsInput,
	blockNumber *big.Int,
) cre.Promise[House] {
	calldata, err := c.Codec.EncodeGetHouseDetailsMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[House](*new(House), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (House, error) {
		return c.Codec.DecodeGetHouseDetailsMethodOutput(response.Data)
	})

}

func (c HouseRWA) GetListing(
	runtime cre.Runtime,
	args GetListingInput,
	blockNumber *big.Int,
) cre.Promise[Listing] {
	calldata, err := c.Codec.EncodeGetListingMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[Listing](*new(Listing), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (Listing, error) {
		return c.Codec.DecodeGetListingMethodOutput(response.Data)
	})

}

func (c HouseRWA) GetTotalBillsCount(
	runtime cre.Runtime,
	args GetTotalBillsCountInput,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeGetTotalBillsCountMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeGetTotalBillsCountMethodOutput(response.Data)
	})

}

func (c HouseRWA) HasKYC(
	runtime cre.Runtime,
	args HasKYCInput,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeHasKYCMethodCall(args)
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
		return c.Codec.DecodeHasKYCMethodOutput(response.Data)
	})

}

func (c HouseRWA) HighValueThresholdUSD(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeHighValueThresholdUSDMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeHighValueThresholdUSDMethodOutput(response.Data)
	})

}

func (c HouseRWA) HouseBills(
	runtime cre.Runtime,
	args HouseBillsInput,
	blockNumber *big.Int,
) cre.Promise[HouseBillsOutput] {
	calldata, err := c.Codec.EncodeHouseBillsMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[HouseBillsOutput](HouseBillsOutput{}, err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (HouseBillsOutput, error) {
		return c.Codec.DecodeHouseBillsMethodOutput(response.Data)
	})

}

func (c HouseRWA) Houses(
	runtime cre.Runtime,
	args HousesInput,
	blockNumber *big.Int,
) cre.Promise[HousesOutput] {
	calldata, err := c.Codec.EncodeHousesMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[HousesOutput](HousesOutput{}, err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (HousesOutput, error) {
		return c.Codec.DecodeHousesMethodOutput(response.Data)
	})

}

func (c HouseRWA) IsApprovedForAll(
	runtime cre.Runtime,
	args IsApprovedForAllInput,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeIsApprovedForAllMethodCall(args)
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
		return c.Codec.DecodeIsApprovedForAllMethodOutput(response.Data)
	})

}

func (c HouseRWA) IsRented(
	runtime cre.Runtime,
	args IsRentedInput,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeIsRentedMethodCall(args)
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
		return c.Codec.DecodeIsRentedMethodOutput(response.Data)
	})

}

func (c HouseRWA) KeyExchanges(
	runtime cre.Runtime,
	args KeyExchangesInput,
	blockNumber *big.Int,
) cre.Promise[KeyExchangesOutput] {
	calldata, err := c.Codec.EncodeKeyExchangesMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[KeyExchangesOutput](KeyExchangesOutput{}, err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (KeyExchangesOutput, error) {
		return c.Codec.DecodeKeyExchangesMethodOutput(response.Data)
	})

}

func (c HouseRWA) KycInfo(
	runtime cre.Runtime,
	args KycInfoInput,
	blockNumber *big.Int,
) cre.Promise[KycInfoOutput] {
	calldata, err := c.Codec.EncodeKycInfoMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[KycInfoOutput](KycInfoOutput{}, err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (KycInfoOutput, error) {
		return c.Codec.DecodeKycInfoMethodOutput(response.Data)
	})

}

func (c HouseRWA) Listings(
	runtime cre.Runtime,
	args ListingsInput,
	blockNumber *big.Int,
) cre.Promise[ListingsOutput] {
	calldata, err := c.Codec.EncodeListingsMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[ListingsOutput](ListingsOutput{}, err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (ListingsOutput, error) {
		return c.Codec.DecodeListingsMethodOutput(response.Data)
	})

}

func (c HouseRWA) MinKYCLevelForHighValue(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[uint8] {
	calldata, err := c.Codec.EncodeMinKYCLevelForHighValueMethodCall()
	if err != nil {
		return cre.PromiseFromResult[uint8](*new(uint8), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (uint8, error) {
		return c.Codec.DecodeMinKYCLevelForHighValueMethodOutput(response.Data)
	})

}

func (c HouseRWA) MinKYCLevelForMint(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[uint8] {
	calldata, err := c.Codec.EncodeMinKYCLevelForMintMethodCall()
	if err != nil {
		return cre.PromiseFromResult[uint8](*new(uint8), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (uint8, error) {
		return c.Codec.DecodeMinKYCLevelForMintMethodOutput(response.Data)
	})

}

func (c HouseRWA) MintingPaused(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeMintingPausedMethodCall()
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
		return c.Codec.DecodeMintingPausedMethodOutput(response.Data)
	})

}

func (c HouseRWA) Name(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[string] {
	calldata, err := c.Codec.EncodeNameMethodCall()
	if err != nil {
		return cre.PromiseFromResult[string](*new(string), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (string, error) {
		return c.Codec.DecodeNameMethodOutput(response.Data)
	})

}

func (c HouseRWA) NextTokenId(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeNextTokenIdMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeNextTokenIdMethodOutput(response.Data)
	})

}

func (c HouseRWA) Owner(
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

func (c HouseRWA) OwnerOf(
	runtime cre.Runtime,
	args OwnerOfInput,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodeOwnerOfMethodCall(args)
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
		return c.Codec.DecodeOwnerOfMethodOutput(response.Data)
	})

}

func (c HouseRWA) Paused(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodePausedMethodCall()
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
		return c.Codec.DecodePausedMethodOutput(response.Data)
	})

}

func (c HouseRWA) PaymentsPaused(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodePaymentsPausedMethodCall()
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
		return c.Codec.DecodePaymentsPausedMethodOutput(response.Data)
	})

}

func (c HouseRWA) PendingOwner(
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

func (c HouseRWA) PendingRentalDeposits(
	runtime cre.Runtime,
	args PendingRentalDepositsInput,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodePendingRentalDepositsMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodePendingRentalDepositsMethodOutput(response.Data)
	})

}

func (c HouseRWA) PriceFeed(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[common.Address] {
	calldata, err := c.Codec.EncodePriceFeedMethodCall()
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
		return c.Codec.DecodePriceFeedMethodOutput(response.Data)
	})

}

func (c HouseRWA) ProxiableUUID(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[[32]byte] {
	calldata, err := c.Codec.EncodeProxiableUUIDMethodCall()
	if err != nil {
		return cre.PromiseFromResult[[32]byte](*new([32]byte), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) ([32]byte, error) {
		return c.Codec.DecodeProxiableUUIDMethodOutput(response.Data)
	})

}

func (c HouseRWA) Rentals(
	runtime cre.Runtime,
	args RentalsInput,
	blockNumber *big.Int,
) cre.Promise[RentalsOutput] {
	calldata, err := c.Codec.EncodeRentalsMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[RentalsOutput](RentalsOutput{}, err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (RentalsOutput, error) {
		return c.Codec.DecodeRentalsMethodOutput(response.Data)
	})

}

func (c HouseRWA) RentalsPaused(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeRentalsPausedMethodCall()
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
		return c.Codec.DecodeRentalsPausedMethodOutput(response.Data)
	})

}

func (c HouseRWA) SalesPaused(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeSalesPausedMethodCall()
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
		return c.Codec.DecodeSalesPausedMethodOutput(response.Data)
	})

}

func (c HouseRWA) SupportsInterface(
	runtime cre.Runtime,
	args SupportsInterfaceInput,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeSupportsInterfaceMethodCall(args)
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
		return c.Codec.DecodeSupportsInterfaceMethodOutput(response.Data)
	})

}

func (c HouseRWA) Symbol(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[string] {
	calldata, err := c.Codec.EncodeSymbolMethodCall()
	if err != nil {
		return cre.PromiseFromResult[string](*new(string), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (string, error) {
		return c.Codec.DecodeSymbolMethodOutput(response.Data)
	})

}

func (c HouseRWA) TokenByIndex(
	runtime cre.Runtime,
	args TokenByIndexInput,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeTokenByIndexMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeTokenByIndexMethodOutput(response.Data)
	})

}

func (c HouseRWA) TokenOfOwnerByIndex(
	runtime cre.Runtime,
	args TokenOfOwnerByIndexInput,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeTokenOfOwnerByIndexMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeTokenOfOwnerByIndexMethodOutput(response.Data)
	})

}

func (c HouseRWA) TokenURI(
	runtime cre.Runtime,
	args TokenURIInput,
	blockNumber *big.Int,
) cre.Promise[string] {
	calldata, err := c.Codec.EncodeTokenURIMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[string](*new(string), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (string, error) {
		return c.Codec.DecodeTokenURIMethodOutput(response.Data)
	})

}

func (c HouseRWA) TotalFeesCollected(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeTotalFeesCollectedMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeTotalFeesCollectedMethodOutput(response.Data)
	})

}

func (c HouseRWA) TotalSupply(
	runtime cre.Runtime,
	blockNumber *big.Int,
) cre.Promise[*big.Int] {
	calldata, err := c.Codec.EncodeTotalSupplyMethodCall()
	if err != nil {
		return cre.PromiseFromResult[*big.Int](*new(*big.Int), err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (*big.Int, error) {
		return c.Codec.DecodeTotalSupplyMethodOutput(response.Data)
	})

}

func (c HouseRWA) TrustedBillProviders(
	runtime cre.Runtime,
	args TrustedBillProvidersInput,
	blockNumber *big.Int,
) cre.Promise[bool] {
	calldata, err := c.Codec.EncodeTrustedBillProvidersMethodCall(args)
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
		return c.Codec.DecodeTrustedBillProvidersMethodOutput(response.Data)
	})

}

func (c HouseRWA) Validators(
	runtime cre.Runtime,
	args ValidatorsInput,
	blockNumber *big.Int,
) cre.Promise[ValidatorsOutput] {
	calldata, err := c.Codec.EncodeValidatorsMethodCall(args)
	if err != nil {
		return cre.PromiseFromResult[ValidatorsOutput](ValidatorsOutput{}, err)
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
	return cre.Then(promise, func(response *evm.CallContractReply) (ValidatorsOutput, error) {
		return c.Codec.DecodeValidatorsMethodOutput(response.Data)
	})

}

func (c HouseRWA) WriteReportFromBill(
	runtime cre.Runtime,
	input Bill,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	encoded, err := c.Codec.EncodeBillStruct(input)
	if err != nil {
		return cre.PromiseFromResult[*evm.WriteReportReply](nil, err)
	}
	promise := runtime.GenerateReport(&pb2.ReportRequest{
		EncodedPayload: encoded,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	})

	return cre.ThenPromise(promise, func(report *cre.Report) cre.Promise[*evm.WriteReportReply] {
		return c.client.WriteReport(runtime, &evm.WriteCreReportRequest{
			Receiver:  c.Address.Bytes(),
			Report:    report,
			GasConfig: gasConfig,
		})
	})
}

func (c HouseRWA) WriteReportFromHouse(
	runtime cre.Runtime,
	input House,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	encoded, err := c.Codec.EncodeHouseStruct(input)
	if err != nil {
		return cre.PromiseFromResult[*evm.WriteReportReply](nil, err)
	}
	promise := runtime.GenerateReport(&pb2.ReportRequest{
		EncodedPayload: encoded,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	})

	return cre.ThenPromise(promise, func(report *cre.Report) cre.Promise[*evm.WriteReportReply] {
		return c.client.WriteReport(runtime, &evm.WriteCreReportRequest{
			Receiver:  c.Address.Bytes(),
			Report:    report,
			GasConfig: gasConfig,
		})
	})
}

func (c HouseRWA) WriteReportFromListing(
	runtime cre.Runtime,
	input Listing,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	encoded, err := c.Codec.EncodeListingStruct(input)
	if err != nil {
		return cre.PromiseFromResult[*evm.WriteReportReply](nil, err)
	}
	promise := runtime.GenerateReport(&pb2.ReportRequest{
		EncodedPayload: encoded,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	})

	return cre.ThenPromise(promise, func(report *cre.Report) cre.Promise[*evm.WriteReportReply] {
		return c.client.WriteReport(runtime, &evm.WriteCreReportRequest{
			Receiver:  c.Address.Bytes(),
			Report:    report,
			GasConfig: gasConfig,
		})
	})
}

func (c HouseRWA) WriteReportFromRentalAgreement(
	runtime cre.Runtime,
	input RentalAgreement,
	gasConfig *evm.GasConfig,
) cre.Promise[*evm.WriteReportReply] {
	encoded, err := c.Codec.EncodeRentalAgreementStruct(input)
	if err != nil {
		return cre.PromiseFromResult[*evm.WriteReportReply](nil, err)
	}
	promise := runtime.GenerateReport(&pb2.ReportRequest{
		EncodedPayload: encoded,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	})

	return cre.ThenPromise(promise, func(report *cre.Report) cre.Promise[*evm.WriteReportReply] {
		return c.client.WriteReport(runtime, &evm.WriteCreReportRequest{
			Receiver:  c.Address.Bytes(),
			Report:    report,
			GasConfig: gasConfig,
		})
	})
}

func (c HouseRWA) WriteReport(
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

// DecodeAddressEmptyCodeError decodes a AddressEmptyCode error from revert data.
func (c *HouseRWA) DecodeAddressEmptyCodeError(data []byte) (*AddressEmptyCode, error) {
	args := c.ABI.Errors["AddressEmptyCode"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	target, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for target in AddressEmptyCode error")
	}

	return &AddressEmptyCode{
		Target: target,
	}, nil
}

// Error implements the error interface for AddressEmptyCode.
func (e *AddressEmptyCode) Error() string {
	return fmt.Sprintf("AddressEmptyCode error: target=%v;", e.Target)
}

// DecodeERC1967InvalidImplementationError decodes a ERC1967InvalidImplementation error from revert data.
func (c *HouseRWA) DecodeERC1967InvalidImplementationError(data []byte) (*ERC1967InvalidImplementation, error) {
	args := c.ABI.Errors["ERC1967InvalidImplementation"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	implementation, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for implementation in ERC1967InvalidImplementation error")
	}

	return &ERC1967InvalidImplementation{
		Implementation: implementation,
	}, nil
}

// Error implements the error interface for ERC1967InvalidImplementation.
func (e *ERC1967InvalidImplementation) Error() string {
	return fmt.Sprintf("ERC1967InvalidImplementation error: implementation=%v;", e.Implementation)
}

// DecodeERC1967NonPayableError decodes a ERC1967NonPayable error from revert data.
func (c *HouseRWA) DecodeERC1967NonPayableError(data []byte) (*ERC1967NonPayable, error) {
	args := c.ABI.Errors["ERC1967NonPayable"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &ERC1967NonPayable{}, nil
}

// Error implements the error interface for ERC1967NonPayable.
func (e *ERC1967NonPayable) Error() string {
	return fmt.Sprintf("ERC1967NonPayable error:")
}

// DecodeERC721EnumerableForbiddenBatchMintError decodes a ERC721EnumerableForbiddenBatchMint error from revert data.
func (c *HouseRWA) DecodeERC721EnumerableForbiddenBatchMintError(data []byte) (*ERC721EnumerableForbiddenBatchMint, error) {
	args := c.ABI.Errors["ERC721EnumerableForbiddenBatchMint"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &ERC721EnumerableForbiddenBatchMint{}, nil
}

// Error implements the error interface for ERC721EnumerableForbiddenBatchMint.
func (e *ERC721EnumerableForbiddenBatchMint) Error() string {
	return fmt.Sprintf("ERC721EnumerableForbiddenBatchMint error:")
}

// DecodeERC721IncorrectOwnerError decodes a ERC721IncorrectOwner error from revert data.
func (c *HouseRWA) DecodeERC721IncorrectOwnerError(data []byte) (*ERC721IncorrectOwner, error) {
	args := c.ABI.Errors["ERC721IncorrectOwner"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 3 {
		return nil, fmt.Errorf("expected 3 values, got %d", len(values))
	}

	sender, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for sender in ERC721IncorrectOwner error")
	}

	tokenId, ok1 := values[1].(*big.Int)
	if !ok1 {
		return nil, fmt.Errorf("unexpected type for tokenId in ERC721IncorrectOwner error")
	}

	owner, ok2 := values[2].(common.Address)
	if !ok2 {
		return nil, fmt.Errorf("unexpected type for owner in ERC721IncorrectOwner error")
	}

	return &ERC721IncorrectOwner{
		Sender:  sender,
		TokenId: tokenId,
		Owner:   owner,
	}, nil
}

// Error implements the error interface for ERC721IncorrectOwner.
func (e *ERC721IncorrectOwner) Error() string {
	return fmt.Sprintf("ERC721IncorrectOwner error: sender=%v; tokenId=%v; owner=%v;", e.Sender, e.TokenId, e.Owner)
}

// DecodeERC721InsufficientApprovalError decodes a ERC721InsufficientApproval error from revert data.
func (c *HouseRWA) DecodeERC721InsufficientApprovalError(data []byte) (*ERC721InsufficientApproval, error) {
	args := c.ABI.Errors["ERC721InsufficientApproval"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(values))
	}

	operator, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for operator in ERC721InsufficientApproval error")
	}

	tokenId, ok1 := values[1].(*big.Int)
	if !ok1 {
		return nil, fmt.Errorf("unexpected type for tokenId in ERC721InsufficientApproval error")
	}

	return &ERC721InsufficientApproval{
		Operator: operator,
		TokenId:  tokenId,
	}, nil
}

// Error implements the error interface for ERC721InsufficientApproval.
func (e *ERC721InsufficientApproval) Error() string {
	return fmt.Sprintf("ERC721InsufficientApproval error: operator=%v; tokenId=%v;", e.Operator, e.TokenId)
}

// DecodeERC721InvalidApproverError decodes a ERC721InvalidApprover error from revert data.
func (c *HouseRWA) DecodeERC721InvalidApproverError(data []byte) (*ERC721InvalidApprover, error) {
	args := c.ABI.Errors["ERC721InvalidApprover"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	approver, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for approver in ERC721InvalidApprover error")
	}

	return &ERC721InvalidApprover{
		Approver: approver,
	}, nil
}

// Error implements the error interface for ERC721InvalidApprover.
func (e *ERC721InvalidApprover) Error() string {
	return fmt.Sprintf("ERC721InvalidApprover error: approver=%v;", e.Approver)
}

// DecodeERC721InvalidOperatorError decodes a ERC721InvalidOperator error from revert data.
func (c *HouseRWA) DecodeERC721InvalidOperatorError(data []byte) (*ERC721InvalidOperator, error) {
	args := c.ABI.Errors["ERC721InvalidOperator"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	operator, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for operator in ERC721InvalidOperator error")
	}

	return &ERC721InvalidOperator{
		Operator: operator,
	}, nil
}

// Error implements the error interface for ERC721InvalidOperator.
func (e *ERC721InvalidOperator) Error() string {
	return fmt.Sprintf("ERC721InvalidOperator error: operator=%v;", e.Operator)
}

// DecodeERC721InvalidOwnerError decodes a ERC721InvalidOwner error from revert data.
func (c *HouseRWA) DecodeERC721InvalidOwnerError(data []byte) (*ERC721InvalidOwner, error) {
	args := c.ABI.Errors["ERC721InvalidOwner"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	owner, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for owner in ERC721InvalidOwner error")
	}

	return &ERC721InvalidOwner{
		Owner: owner,
	}, nil
}

// Error implements the error interface for ERC721InvalidOwner.
func (e *ERC721InvalidOwner) Error() string {
	return fmt.Sprintf("ERC721InvalidOwner error: owner=%v;", e.Owner)
}

// DecodeERC721InvalidReceiverError decodes a ERC721InvalidReceiver error from revert data.
func (c *HouseRWA) DecodeERC721InvalidReceiverError(data []byte) (*ERC721InvalidReceiver, error) {
	args := c.ABI.Errors["ERC721InvalidReceiver"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	receiver, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for receiver in ERC721InvalidReceiver error")
	}

	return &ERC721InvalidReceiver{
		Receiver: receiver,
	}, nil
}

// Error implements the error interface for ERC721InvalidReceiver.
func (e *ERC721InvalidReceiver) Error() string {
	return fmt.Sprintf("ERC721InvalidReceiver error: receiver=%v;", e.Receiver)
}

// DecodeERC721InvalidSenderError decodes a ERC721InvalidSender error from revert data.
func (c *HouseRWA) DecodeERC721InvalidSenderError(data []byte) (*ERC721InvalidSender, error) {
	args := c.ABI.Errors["ERC721InvalidSender"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	sender, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for sender in ERC721InvalidSender error")
	}

	return &ERC721InvalidSender{
		Sender: sender,
	}, nil
}

// Error implements the error interface for ERC721InvalidSender.
func (e *ERC721InvalidSender) Error() string {
	return fmt.Sprintf("ERC721InvalidSender error: sender=%v;", e.Sender)
}

// DecodeERC721NonexistentTokenError decodes a ERC721NonexistentToken error from revert data.
func (c *HouseRWA) DecodeERC721NonexistentTokenError(data []byte) (*ERC721NonexistentToken, error) {
	args := c.ABI.Errors["ERC721NonexistentToken"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	tokenId, ok0 := values[0].(*big.Int)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for tokenId in ERC721NonexistentToken error")
	}

	return &ERC721NonexistentToken{
		TokenId: tokenId,
	}, nil
}

// Error implements the error interface for ERC721NonexistentToken.
func (e *ERC721NonexistentToken) Error() string {
	return fmt.Sprintf("ERC721NonexistentToken error: tokenId=%v;", e.TokenId)
}

// DecodeERC721OutOfBoundsIndexError decodes a ERC721OutOfBoundsIndex error from revert data.
func (c *HouseRWA) DecodeERC721OutOfBoundsIndexError(data []byte) (*ERC721OutOfBoundsIndex, error) {
	args := c.ABI.Errors["ERC721OutOfBoundsIndex"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 2 {
		return nil, fmt.Errorf("expected 2 values, got %d", len(values))
	}

	owner, ok0 := values[0].(common.Address)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for owner in ERC721OutOfBoundsIndex error")
	}

	index, ok1 := values[1].(*big.Int)
	if !ok1 {
		return nil, fmt.Errorf("unexpected type for index in ERC721OutOfBoundsIndex error")
	}

	return &ERC721OutOfBoundsIndex{
		Owner: owner,
		Index: index,
	}, nil
}

// Error implements the error interface for ERC721OutOfBoundsIndex.
func (e *ERC721OutOfBoundsIndex) Error() string {
	return fmt.Sprintf("ERC721OutOfBoundsIndex error: owner=%v; index=%v;", e.Owner, e.Index)
}

// DecodeEnforcedPauseError decodes a EnforcedPause error from revert data.
func (c *HouseRWA) DecodeEnforcedPauseError(data []byte) (*EnforcedPause, error) {
	args := c.ABI.Errors["EnforcedPause"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &EnforcedPause{}, nil
}

// Error implements the error interface for EnforcedPause.
func (e *EnforcedPause) Error() string {
	return fmt.Sprintf("EnforcedPause error:")
}

// DecodeExpectedPauseError decodes a ExpectedPause error from revert data.
func (c *HouseRWA) DecodeExpectedPauseError(data []byte) (*ExpectedPause, error) {
	args := c.ABI.Errors["ExpectedPause"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &ExpectedPause{}, nil
}

// Error implements the error interface for ExpectedPause.
func (e *ExpectedPause) Error() string {
	return fmt.Sprintf("ExpectedPause error:")
}

// DecodeFailedInnerCallError decodes a FailedInnerCall error from revert data.
func (c *HouseRWA) DecodeFailedInnerCallError(data []byte) (*FailedInnerCall, error) {
	args := c.ABI.Errors["FailedInnerCall"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &FailedInnerCall{}, nil
}

// Error implements the error interface for FailedInnerCall.
func (e *FailedInnerCall) Error() string {
	return fmt.Sprintf("FailedInnerCall error:")
}

// DecodeInvalidInitializationError decodes a InvalidInitialization error from revert data.
func (c *HouseRWA) DecodeInvalidInitializationError(data []byte) (*InvalidInitialization, error) {
	args := c.ABI.Errors["InvalidInitialization"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &InvalidInitialization{}, nil
}

// Error implements the error interface for InvalidInitialization.
func (e *InvalidInitialization) Error() string {
	return fmt.Sprintf("InvalidInitialization error:")
}

// DecodeNotInitializingError decodes a NotInitializing error from revert data.
func (c *HouseRWA) DecodeNotInitializingError(data []byte) (*NotInitializing, error) {
	args := c.ABI.Errors["NotInitializing"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &NotInitializing{}, nil
}

// Error implements the error interface for NotInitializing.
func (e *NotInitializing) Error() string {
	return fmt.Sprintf("NotInitializing error:")
}

// DecodeOwnableInvalidOwnerError decodes a OwnableInvalidOwner error from revert data.
func (c *HouseRWA) DecodeOwnableInvalidOwnerError(data []byte) (*OwnableInvalidOwner, error) {
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
func (c *HouseRWA) DecodeOwnableUnauthorizedAccountError(data []byte) (*OwnableUnauthorizedAccount, error) {
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

// DecodeReentrancyGuardReentrantCallError decodes a ReentrancyGuardReentrantCall error from revert data.
func (c *HouseRWA) DecodeReentrancyGuardReentrantCallError(data []byte) (*ReentrancyGuardReentrantCall, error) {
	args := c.ABI.Errors["ReentrancyGuardReentrantCall"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &ReentrancyGuardReentrantCall{}, nil
}

// Error implements the error interface for ReentrancyGuardReentrantCall.
func (e *ReentrancyGuardReentrantCall) Error() string {
	return fmt.Sprintf("ReentrancyGuardReentrantCall error:")
}

// DecodeUUPSUnauthorizedCallContextError decodes a UUPSUnauthorizedCallContext error from revert data.
func (c *HouseRWA) DecodeUUPSUnauthorizedCallContextError(data []byte) (*UUPSUnauthorizedCallContext, error) {
	args := c.ABI.Errors["UUPSUnauthorizedCallContext"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 0 {
		return nil, fmt.Errorf("expected 0 values, got %d", len(values))
	}

	return &UUPSUnauthorizedCallContext{}, nil
}

// Error implements the error interface for UUPSUnauthorizedCallContext.
func (e *UUPSUnauthorizedCallContext) Error() string {
	return fmt.Sprintf("UUPSUnauthorizedCallContext error:")
}

// DecodeUUPSUnsupportedProxiableUUIDError decodes a UUPSUnsupportedProxiableUUID error from revert data.
func (c *HouseRWA) DecodeUUPSUnsupportedProxiableUUIDError(data []byte) (*UUPSUnsupportedProxiableUUID, error) {
	args := c.ABI.Errors["UUPSUnsupportedProxiableUUID"].Inputs
	values, err := args.Unpack(data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack error: %w", err)
	}
	if len(values) != 1 {
		return nil, fmt.Errorf("expected 1 values, got %d", len(values))
	}

	slot, ok0 := values[0].([32]byte)
	if !ok0 {
		return nil, fmt.Errorf("unexpected type for slot in UUPSUnsupportedProxiableUUID error")
	}

	return &UUPSUnsupportedProxiableUUID{
		Slot: slot,
	}, nil
}

// Error implements the error interface for UUPSUnsupportedProxiableUUID.
func (e *UUPSUnsupportedProxiableUUID) Error() string {
	return fmt.Sprintf("UUPSUnsupportedProxiableUUID error: slot=%v;", e.Slot)
}

func (c *HouseRWA) UnpackError(data []byte) (any, error) {
	switch common.Bytes2Hex(data[:4]) {
	case common.Bytes2Hex(c.ABI.Errors["AddressEmptyCode"].ID.Bytes()[:4]):
		return c.DecodeAddressEmptyCodeError(data)
	case common.Bytes2Hex(c.ABI.Errors["ERC1967InvalidImplementation"].ID.Bytes()[:4]):
		return c.DecodeERC1967InvalidImplementationError(data)
	case common.Bytes2Hex(c.ABI.Errors["ERC1967NonPayable"].ID.Bytes()[:4]):
		return c.DecodeERC1967NonPayableError(data)
	case common.Bytes2Hex(c.ABI.Errors["ERC721EnumerableForbiddenBatchMint"].ID.Bytes()[:4]):
		return c.DecodeERC721EnumerableForbiddenBatchMintError(data)
	case common.Bytes2Hex(c.ABI.Errors["ERC721IncorrectOwner"].ID.Bytes()[:4]):
		return c.DecodeERC721IncorrectOwnerError(data)
	case common.Bytes2Hex(c.ABI.Errors["ERC721InsufficientApproval"].ID.Bytes()[:4]):
		return c.DecodeERC721InsufficientApprovalError(data)
	case common.Bytes2Hex(c.ABI.Errors["ERC721InvalidApprover"].ID.Bytes()[:4]):
		return c.DecodeERC721InvalidApproverError(data)
	case common.Bytes2Hex(c.ABI.Errors["ERC721InvalidOperator"].ID.Bytes()[:4]):
		return c.DecodeERC721InvalidOperatorError(data)
	case common.Bytes2Hex(c.ABI.Errors["ERC721InvalidOwner"].ID.Bytes()[:4]):
		return c.DecodeERC721InvalidOwnerError(data)
	case common.Bytes2Hex(c.ABI.Errors["ERC721InvalidReceiver"].ID.Bytes()[:4]):
		return c.DecodeERC721InvalidReceiverError(data)
	case common.Bytes2Hex(c.ABI.Errors["ERC721InvalidSender"].ID.Bytes()[:4]):
		return c.DecodeERC721InvalidSenderError(data)
	case common.Bytes2Hex(c.ABI.Errors["ERC721NonexistentToken"].ID.Bytes()[:4]):
		return c.DecodeERC721NonexistentTokenError(data)
	case common.Bytes2Hex(c.ABI.Errors["ERC721OutOfBoundsIndex"].ID.Bytes()[:4]):
		return c.DecodeERC721OutOfBoundsIndexError(data)
	case common.Bytes2Hex(c.ABI.Errors["EnforcedPause"].ID.Bytes()[:4]):
		return c.DecodeEnforcedPauseError(data)
	case common.Bytes2Hex(c.ABI.Errors["ExpectedPause"].ID.Bytes()[:4]):
		return c.DecodeExpectedPauseError(data)
	case common.Bytes2Hex(c.ABI.Errors["FailedInnerCall"].ID.Bytes()[:4]):
		return c.DecodeFailedInnerCallError(data)
	case common.Bytes2Hex(c.ABI.Errors["InvalidInitialization"].ID.Bytes()[:4]):
		return c.DecodeInvalidInitializationError(data)
	case common.Bytes2Hex(c.ABI.Errors["NotInitializing"].ID.Bytes()[:4]):
		return c.DecodeNotInitializingError(data)
	case common.Bytes2Hex(c.ABI.Errors["OwnableInvalidOwner"].ID.Bytes()[:4]):
		return c.DecodeOwnableInvalidOwnerError(data)
	case common.Bytes2Hex(c.ABI.Errors["OwnableUnauthorizedAccount"].ID.Bytes()[:4]):
		return c.DecodeOwnableUnauthorizedAccountError(data)
	case common.Bytes2Hex(c.ABI.Errors["ReentrancyGuardReentrantCall"].ID.Bytes()[:4]):
		return c.DecodeReentrancyGuardReentrantCallError(data)
	case common.Bytes2Hex(c.ABI.Errors["UUPSUnauthorizedCallContext"].ID.Bytes()[:4]):
		return c.DecodeUUPSUnauthorizedCallContextError(data)
	case common.Bytes2Hex(c.ABI.Errors["UUPSUnsupportedProxiableUUID"].ID.Bytes()[:4]):
		return c.DecodeUUPSUnsupportedProxiableUUIDError(data)
	default:
		return nil, errors.New("unknown error selector")
	}
}

// ApprovalTrigger wraps the raw log trigger and provides decoded ApprovalDecoded data
type ApprovalTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into Approval data
func (t *ApprovalTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ApprovalDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeApproval(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Approval log: %w", err)
	}

	return &bindings.DecodedLog[ApprovalDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerApprovalLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ApprovalTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ApprovalDecoded]], error) {
	event := c.ABI.Events["Approval"]
	topics, err := c.Codec.EncodeApprovalTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for Approval: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ApprovalTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsApproval(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ApprovalLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// ApprovalForAllTrigger wraps the raw log trigger and provides decoded ApprovalForAllDecoded data
type ApprovalForAllTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into ApprovalForAll data
func (t *ApprovalForAllTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ApprovalForAllDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeApprovalForAll(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ApprovalForAll log: %w", err)
	}

	return &bindings.DecodedLog[ApprovalForAllDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerApprovalForAllLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ApprovalForAllTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ApprovalForAllDecoded]], error) {
	event := c.ABI.Events["ApprovalForAll"]
	topics, err := c.Codec.EncodeApprovalForAllTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ApprovalForAll: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ApprovalForAllTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsApprovalForAll(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ApprovalForAllLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// BillCreatedTrigger wraps the raw log trigger and provides decoded BillCreatedDecoded data
type BillCreatedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into BillCreated data
func (t *BillCreatedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[BillCreatedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeBillCreated(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode BillCreated log: %w", err)
	}

	return &bindings.DecodedLog[BillCreatedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerBillCreatedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []BillCreatedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[BillCreatedDecoded]], error) {
	event := c.ABI.Events["BillCreated"]
	topics, err := c.Codec.EncodeBillCreatedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for BillCreated: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &BillCreatedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsBillCreated(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.BillCreatedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// BillDisputedTrigger wraps the raw log trigger and provides decoded BillDisputedDecoded data
type BillDisputedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into BillDisputed data
func (t *BillDisputedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[BillDisputedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeBillDisputed(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode BillDisputed log: %w", err)
	}

	return &bindings.DecodedLog[BillDisputedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerBillDisputedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []BillDisputedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[BillDisputedDecoded]], error) {
	event := c.ABI.Events["BillDisputed"]
	topics, err := c.Codec.EncodeBillDisputedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for BillDisputed: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &BillDisputedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsBillDisputed(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.BillDisputedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// BillPaidTrigger wraps the raw log trigger and provides decoded BillPaidDecoded data
type BillPaidTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into BillPaid data
func (t *BillPaidTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[BillPaidDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeBillPaid(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode BillPaid log: %w", err)
	}

	return &bindings.DecodedLog[BillPaidDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerBillPaidLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []BillPaidTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[BillPaidDecoded]], error) {
	event := c.ABI.Events["BillPaid"]
	topics, err := c.Codec.EncodeBillPaidTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for BillPaid: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &BillPaidTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsBillPaid(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.BillPaidLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// CircuitBreakerTriggeredTrigger wraps the raw log trigger and provides decoded CircuitBreakerTriggeredDecoded data
type CircuitBreakerTriggeredTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into CircuitBreakerTriggered data
func (t *CircuitBreakerTriggeredTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[CircuitBreakerTriggeredDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeCircuitBreakerTriggered(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode CircuitBreakerTriggered log: %w", err)
	}

	return &bindings.DecodedLog[CircuitBreakerTriggeredDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerCircuitBreakerTriggeredLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []CircuitBreakerTriggeredTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[CircuitBreakerTriggeredDecoded]], error) {
	event := c.ABI.Events["CircuitBreakerTriggered"]
	topics, err := c.Codec.EncodeCircuitBreakerTriggeredTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for CircuitBreakerTriggered: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &CircuitBreakerTriggeredTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsCircuitBreakerTriggered(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.CircuitBreakerTriggeredLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// EmergencyActionTrigger wraps the raw log trigger and provides decoded EmergencyActionDecoded data
type EmergencyActionTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into EmergencyAction data
func (t *EmergencyActionTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[EmergencyActionDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeEmergencyAction(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode EmergencyAction log: %w", err)
	}

	return &bindings.DecodedLog[EmergencyActionDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerEmergencyActionLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []EmergencyActionTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[EmergencyActionDecoded]], error) {
	event := c.ABI.Events["EmergencyAction"]
	topics, err := c.Codec.EncodeEmergencyActionTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for EmergencyAction: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &EmergencyActionTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsEmergencyAction(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.EmergencyActionLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// HouseListedTrigger wraps the raw log trigger and provides decoded HouseListedDecoded data
type HouseListedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into HouseListed data
func (t *HouseListedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[HouseListedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeHouseListed(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode HouseListed log: %w", err)
	}

	return &bindings.DecodedLog[HouseListedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerHouseListedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []HouseListedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[HouseListedDecoded]], error) {
	event := c.ABI.Events["HouseListed"]
	topics, err := c.Codec.EncodeHouseListedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for HouseListed: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &HouseListedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsHouseListed(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.HouseListedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// HouseMintedTrigger wraps the raw log trigger and provides decoded HouseMintedDecoded data
type HouseMintedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into HouseMinted data
func (t *HouseMintedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[HouseMintedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeHouseMinted(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode HouseMinted log: %w", err)
	}

	return &bindings.DecodedLog[HouseMintedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerHouseMintedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []HouseMintedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[HouseMintedDecoded]], error) {
	event := c.ABI.Events["HouseMinted"]
	topics, err := c.Codec.EncodeHouseMintedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for HouseMinted: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &HouseMintedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsHouseMinted(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.HouseMintedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// HouseSoldTrigger wraps the raw log trigger and provides decoded HouseSoldDecoded data
type HouseSoldTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into HouseSold data
func (t *HouseSoldTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[HouseSoldDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeHouseSold(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode HouseSold log: %w", err)
	}

	return &bindings.DecodedLog[HouseSoldDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerHouseSoldLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []HouseSoldTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[HouseSoldDecoded]], error) {
	event := c.ABI.Events["HouseSold"]
	topics, err := c.Codec.EncodeHouseSoldTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for HouseSold: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &HouseSoldTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsHouseSold(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.HouseSoldLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// InitializedTrigger wraps the raw log trigger and provides decoded InitializedDecoded data
type InitializedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into Initialized data
func (t *InitializedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[InitializedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeInitialized(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Initialized log: %w", err)
	}

	return &bindings.DecodedLog[InitializedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerInitializedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []InitializedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[InitializedDecoded]], error) {
	event := c.ABI.Events["Initialized"]
	topics, err := c.Codec.EncodeInitializedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for Initialized: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &InitializedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsInitialized(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.InitializedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// KYCVerifiedTrigger wraps the raw log trigger and provides decoded KYCVerifiedDecoded data
type KYCVerifiedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into KYCVerified data
func (t *KYCVerifiedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[KYCVerifiedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeKYCVerified(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode KYCVerified log: %w", err)
	}

	return &bindings.DecodedLog[KYCVerifiedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerKYCVerifiedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []KYCVerifiedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[KYCVerifiedDecoded]], error) {
	event := c.ABI.Events["KYCVerified"]
	topics, err := c.Codec.EncodeKYCVerifiedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for KYCVerified: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &KYCVerifiedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsKYCVerified(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.KYCVerifiedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// KeyClaimedTrigger wraps the raw log trigger and provides decoded KeyClaimedDecoded data
type KeyClaimedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into KeyClaimed data
func (t *KeyClaimedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[KeyClaimedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeKeyClaimed(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode KeyClaimed log: %w", err)
	}

	return &bindings.DecodedLog[KeyClaimedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerKeyClaimedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []KeyClaimedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[KeyClaimedDecoded]], error) {
	event := c.ABI.Events["KeyClaimed"]
	topics, err := c.Codec.EncodeKeyClaimedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for KeyClaimed: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &KeyClaimedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsKeyClaimed(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.KeyClaimedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// KeyExchangeCreatedTrigger wraps the raw log trigger and provides decoded KeyExchangeCreatedDecoded data
type KeyExchangeCreatedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into KeyExchangeCreated data
func (t *KeyExchangeCreatedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[KeyExchangeCreatedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeKeyExchangeCreated(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode KeyExchangeCreated log: %w", err)
	}

	return &bindings.DecodedLog[KeyExchangeCreatedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerKeyExchangeCreatedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []KeyExchangeCreatedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[KeyExchangeCreatedDecoded]], error) {
	event := c.ABI.Events["KeyExchangeCreated"]
	topics, err := c.Codec.EncodeKeyExchangeCreatedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for KeyExchangeCreated: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &KeyExchangeCreatedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsKeyExchangeCreated(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.KeyExchangeCreatedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// OwnershipTransferStartedTrigger wraps the raw log trigger and provides decoded OwnershipTransferStartedDecoded data
type OwnershipTransferStartedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
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

func (c *HouseRWA) LogTriggerOwnershipTransferStartedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []OwnershipTransferStartedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[OwnershipTransferStartedDecoded]], error) {
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

func (c *HouseRWA) FilterLogsOwnershipTransferStarted(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
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
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
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

func (c *HouseRWA) LogTriggerOwnershipTransferredLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []OwnershipTransferredTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[OwnershipTransferredDecoded]], error) {
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

func (c *HouseRWA) FilterLogsOwnershipTransferred(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
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

// PausedTrigger wraps the raw log trigger and provides decoded PausedDecoded data
type PausedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into Paused data
func (t *PausedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[PausedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodePaused(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Paused log: %w", err)
	}

	return &bindings.DecodedLog[PausedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerPausedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []PausedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[PausedDecoded]], error) {
	event := c.ABI.Events["Paused"]
	topics, err := c.Codec.EncodePausedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for Paused: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &PausedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsPaused(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.PausedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// RentalDepositReceivedTrigger wraps the raw log trigger and provides decoded RentalDepositReceivedDecoded data
type RentalDepositReceivedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into RentalDepositReceived data
func (t *RentalDepositReceivedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[RentalDepositReceivedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeRentalDepositReceived(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode RentalDepositReceived log: %w", err)
	}

	return &bindings.DecodedLog[RentalDepositReceivedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerRentalDepositReceivedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []RentalDepositReceivedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[RentalDepositReceivedDecoded]], error) {
	event := c.ABI.Events["RentalDepositReceived"]
	topics, err := c.Codec.EncodeRentalDepositReceivedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for RentalDepositReceived: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &RentalDepositReceivedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsRentalDepositReceived(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.RentalDepositReceivedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// RentalDepositWithdrawnTrigger wraps the raw log trigger and provides decoded RentalDepositWithdrawnDecoded data
type RentalDepositWithdrawnTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into RentalDepositWithdrawn data
func (t *RentalDepositWithdrawnTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[RentalDepositWithdrawnDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeRentalDepositWithdrawn(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode RentalDepositWithdrawn log: %w", err)
	}

	return &bindings.DecodedLog[RentalDepositWithdrawnDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerRentalDepositWithdrawnLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []RentalDepositWithdrawnTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[RentalDepositWithdrawnDecoded]], error) {
	event := c.ABI.Events["RentalDepositWithdrawn"]
	topics, err := c.Codec.EncodeRentalDepositWithdrawnTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for RentalDepositWithdrawn: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &RentalDepositWithdrawnTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsRentalDepositWithdrawn(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.RentalDepositWithdrawnLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// RentalEndedTrigger wraps the raw log trigger and provides decoded RentalEndedDecoded data
type RentalEndedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into RentalEnded data
func (t *RentalEndedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[RentalEndedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeRentalEnded(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode RentalEnded log: %w", err)
	}

	return &bindings.DecodedLog[RentalEndedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerRentalEndedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []RentalEndedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[RentalEndedDecoded]], error) {
	event := c.ABI.Events["RentalEnded"]
	topics, err := c.Codec.EncodeRentalEndedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for RentalEnded: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &RentalEndedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsRentalEnded(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.RentalEndedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// RentalStartedTrigger wraps the raw log trigger and provides decoded RentalStartedDecoded data
type RentalStartedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into RentalStarted data
func (t *RentalStartedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[RentalStartedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeRentalStarted(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode RentalStarted log: %w", err)
	}

	return &bindings.DecodedLog[RentalStartedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerRentalStartedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []RentalStartedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[RentalStartedDecoded]], error) {
	event := c.ABI.Events["RentalStarted"]
	topics, err := c.Codec.EncodeRentalStartedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for RentalStarted: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &RentalStartedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsRentalStarted(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.RentalStartedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// TransferTrigger wraps the raw log trigger and provides decoded TransferDecoded data
type TransferTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into Transfer data
func (t *TransferTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[TransferDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeTransfer(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Transfer log: %w", err)
	}

	return &bindings.DecodedLog[TransferDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerTransferLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []TransferTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[TransferDecoded]], error) {
	event := c.ABI.Events["Transfer"]
	topics, err := c.Codec.EncodeTransferTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for Transfer: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &TransferTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsTransfer(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.TransferLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// UnpausedTrigger wraps the raw log trigger and provides decoded UnpausedDecoded data
type UnpausedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into Unpaused data
func (t *UnpausedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[UnpausedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeUnpaused(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Unpaused log: %w", err)
	}

	return &bindings.DecodedLog[UnpausedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerUnpausedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []UnpausedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[UnpausedDecoded]], error) {
	event := c.ABI.Events["Unpaused"]
	topics, err := c.Codec.EncodeUnpausedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for Unpaused: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &UnpausedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsUnpaused(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.UnpausedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// UpgradedTrigger wraps the raw log trigger and provides decoded UpgradedDecoded data
type UpgradedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into Upgraded data
func (t *UpgradedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[UpgradedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeUpgraded(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode Upgraded log: %w", err)
	}

	return &bindings.DecodedLog[UpgradedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerUpgradedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []UpgradedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[UpgradedDecoded]], error) {
	event := c.ABI.Events["Upgraded"]
	topics, err := c.Codec.EncodeUpgradedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for Upgraded: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &UpgradedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsUpgraded(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.UpgradedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// ValidatorSlashedTrigger wraps the raw log trigger and provides decoded ValidatorSlashedDecoded data
type ValidatorSlashedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into ValidatorSlashed data
func (t *ValidatorSlashedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ValidatorSlashedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeValidatorSlashed(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ValidatorSlashed log: %w", err)
	}

	return &bindings.DecodedLog[ValidatorSlashedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerValidatorSlashedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ValidatorSlashedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ValidatorSlashedDecoded]], error) {
	event := c.ABI.Events["ValidatorSlashed"]
	topics, err := c.Codec.EncodeValidatorSlashedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ValidatorSlashed: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ValidatorSlashedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsValidatorSlashed(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ValidatorSlashedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}

// ValidatorStakedTrigger wraps the raw log trigger and provides decoded ValidatorStakedDecoded data
type ValidatorStakedTrigger struct {
	cre.Trigger[*evm.Log, *evm.Log]           // Embed the raw trigger
	contract                        *HouseRWA // Keep reference for decoding
}

// Adapt method that decodes the log into ValidatorStaked data
func (t *ValidatorStakedTrigger) Adapt(l *evm.Log) (*bindings.DecodedLog[ValidatorStakedDecoded], error) {
	// Decode the log using the contract's codec
	decoded, err := t.contract.Codec.DecodeValidatorStaked(l)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ValidatorStaked log: %w", err)
	}

	return &bindings.DecodedLog[ValidatorStakedDecoded]{
		Log:  l,        // Original log
		Data: *decoded, // Decoded data
	}, nil
}

func (c *HouseRWA) LogTriggerValidatorStakedLog(chainSelector uint64, confidence evm.ConfidenceLevel, filters []ValidatorStakedTopics) (cre.Trigger[*evm.Log, *bindings.DecodedLog[ValidatorStakedDecoded]], error) {
	event := c.ABI.Events["ValidatorStaked"]
	topics, err := c.Codec.EncodeValidatorStakedTopics(event, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to encode topics for ValidatorStaked: %w", err)
	}

	rawTrigger := evm.LogTrigger(chainSelector, &evm.FilterLogTriggerRequest{
		Addresses:  [][]byte{c.Address.Bytes()},
		Topics:     topics,
		Confidence: confidence,
	})

	return &ValidatorStakedTrigger{
		Trigger:  rawTrigger,
		contract: c,
	}, nil
}

func (c *HouseRWA) FilterLogsValidatorStaked(runtime cre.Runtime, options *bindings.FilterOptions) (cre.Promise[*evm.FilterLogsReply], error) {
	if options == nil {
		return nil, errors.New("FilterLogs options are required.")
	}
	return c.client.FilterLogs(runtime, &evm.FilterLogsRequest{
		FilterQuery: &evm.FilterQuery{
			Addresses: [][]byte{c.Address.Bytes()},
			Topics: []*evm.Topics{
				{Topic: [][]byte{c.Codec.ValidatorStakedLogHash()}},
			},
			BlockHash: options.BlockHash,
			FromBlock: pb.NewBigIntFromInt(options.FromBlock),
			ToBlock:   pb.NewBigIntFromInt(options.ToBlock),
		},
	}), nil
}
