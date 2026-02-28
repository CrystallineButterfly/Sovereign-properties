// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC721Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC721/ERC721Upgradeable.sol";
import {Ownable2StepUpgradeable} from "@openzeppelin/contracts-upgradeable/access/Ownable2StepUpgradeable.sol";
import {ReentrancyGuardUpgradeable} from "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

contract HouseRWA is
    Initializable,
    ERC721Upgradeable,
    Ownable2StepUpgradeable,
    ReentrancyGuardUpgradeable,
    UUPSUpgradeable
{
    uint256 public constant MAX_BILLS_PER_HOUSE = 100;
    uint8 public constant MAX_RECURRENCE_INTERVAL_DAYS = 90;
    uint256 public constant MAX_BILL_DUE_WINDOW = 365 days;
    uint256 public constant MAX_RENTAL_DURATION = 3650 days;
    uint256 public constant KEY_EXCHANGE_EXPIRY = 7 days;
    uint256 public constant PROTOCOL_FEE_BPS = 250;
    uint256 public constant BPS_DENOMINATOR = 10000;

    enum DocumentStorageType {
        IPFS,
        OFF_CHAIN,
        ARWEAVE,
        ENCRYPTED_DB
    }

    enum ListingType {
        NONE,
        FOR_SALE,
        FOR_RENT,
        AUCTION
    }

    enum BillStatus {
        PENDING,
        PAID,
        OVERDUE,
        DISPUTED,
        CANCELLED
    }

    enum DisputeStatus {
        NONE,
        OPEN,
        RESOLVED,
        ESCALATED
    }

    struct House {
        string houseId;
        bytes32 documentHash;
        string documentURI;
        DocumentStorageType storageType;
        address originalOwner;
        uint48 mintedAt;
        bool isVerified;
        uint8 documentCount;
    }

    struct RentalAgreement {
        address renter;
        uint48 startTime;
        uint48 endTime;
        uint96 depositAmount;
        uint96 monthlyRent;
        bool isActive;
        bytes32 encryptedAccessKeyHash;
        DisputeStatus disputeStatus;
    }

    struct Listing {
        ListingType listingType;
        uint96 price;
        address preferredToken;
        bool isPrivateSale;
        address allowedBuyer;
        uint48 createdAt;
        uint48 expiresAt;
        uint8 platformFee;
    }

    struct Bill {
        string billType;
        uint96 amount;
        uint48 dueDate;
        uint48 paidAt;
        BillStatus status;
        bytes32 paymentReference;
        bool isRecurring;
        address provider;
        uint8 recurrenceInterval;
    }

    struct KeyExchange {
        bytes32 keyHash;
        bytes encryptedKey;
        address intendedRecipient;
        uint48 createdAt;
        uint48 expiresAt;
        bool isClaimed;
        uint8 exchangeType;
    }

    struct KYCInfo {
        bool isVerified;
        uint48 verifiedAt;
        bytes32 verificationHash;
        uint8 verificationLevel;
        uint48 expiryDate;
    }

    mapping(uint256 => House) public houses;
    mapping(uint256 => RentalAgreement) public rentals;
    mapping(uint256 => Listing) public listings;
    mapping(uint256 => Bill[]) public houseBills;
    mapping(bytes32 => KeyExchange) public keyExchanges;
    mapping(uint256 => mapping(address => uint256)) public pendingRentalDeposits;
    mapping(address => KYCInfo) public kycInfo;
    mapping(address => bool) public authorizedCREWorkflows;
    mapping(address => bool) public trustedBillProviders;

    address public feeRecipient;
    uint256 public totalFeesCollected;
    uint256 public nextTokenId;

    uint8 public minKYCLevelForMint;
    uint8 public minKYCLevelForHighValue;
    uint256 public highValueThresholdUSD;

    event HouseMinted(
        uint256 indexed tokenId,
        address indexed owner,
        string houseId,
        bytes32 documentHash,
        DocumentStorageType storageType
    );

    event HouseListed(uint256 indexed tokenId, ListingType listingType, uint256 price, bool isPrivateSale);

    event HouseSold(
        uint256 indexed tokenId, address indexed seller, address indexed buyer, uint256 price, uint256 protocolFee
    );

    event RentalStarted(
        uint256 indexed tokenId, address indexed renter, uint256 startTime, uint256 endTime, uint256 deposit
    );

    event RentalEnded(uint256 indexed tokenId, address indexed renter, uint256 depositReturned);

    event RentalDepositReceived(uint256 indexed tokenId, address indexed renter, uint256 amount);
    event RentalDepositWithdrawn(uint256 indexed tokenId, address indexed renter, uint256 amount);

    event KeyExchangeCreated(
        bytes32 indexed keyHash, uint256 indexed tokenId, address indexed recipient, uint256 expiry
    );

    event KeyClaimed(bytes32 indexed keyHash, address indexed claimant, uint256 timestamp);

    event BillCreated(
        uint256 indexed tokenId, uint256 indexed billIndex, string billType, uint256 amount, uint256 dueDate
    );

    event BillPaid(uint256 indexed tokenId, uint256 indexed billIndex, string paymentMethod, bytes32 paymentReference);

    event BillDisputed(uint256 indexed tokenId, uint256 indexed billIndex, address disputer, string reason);

    event KYCVerified(address indexed user, uint8 level, uint256 expiry);

    modifier onlyCRE() {
        require(authorizedCREWorkflows[msg.sender]);
        _;
    }

    modifier onlyHouseOwner(uint256 tokenId) {
        require(ownerOf(tokenId) == msg.sender);
        _;
    }

    modifier validToken(uint256 tokenId) {
        require(_exists(tokenId));
        _;
    }

    constructor() {
        _disableInitializers();
    }

    function initialize(address _owner, address _feeRecipient, address _initialCREWorkflow) public initializer {
        require(_owner != address(0));
        require(_feeRecipient != address(0));
        require(_initialCREWorkflow != address(0));

        __ERC721_init("HouseRWA", "HRWA");
        __Ownable2Step_init();
        __Ownable_init(_owner);
        __ReentrancyGuard_init();
        __UUPSUpgradeable_init();

        feeRecipient = _feeRecipient;
        authorizedCREWorkflows[_initialCREWorkflow] = true;

        minKYCLevelForMint = 1;
        minKYCLevelForHighValue = 2;
        highValueThresholdUSD = 100000 * 100;
    }

    function authorizeCREWorkflow(address workflow) external onlyOwner {
        require(workflow != address(0));
        require(!authorizedCREWorkflows[workflow]);
        authorizedCREWorkflows[workflow] = true;
    }

    function revokeCREWorkflow(address workflow) external onlyOwner {
        require(authorizedCREWorkflows[workflow]);
        authorizedCREWorkflows[workflow] = false;
    }

    function setFeeRecipient(address _feeRecipient) external onlyOwner {
        require(_feeRecipient != address(0));
        feeRecipient = _feeRecipient;
    }

    function setKYCRequirements(uint8 _minLevelForMint, uint8 _minLevelForHighValue, uint256 _highValueThresholdUSD)
        external
        onlyOwner
    {
        minKYCLevelForMint = _minLevelForMint;
        minKYCLevelForHighValue = _minLevelForHighValue;
        highValueThresholdUSD = _highValueThresholdUSD;
    }

    function setTrustedBillProvider(address provider, bool trusted) external onlyOwner {
        require(provider != address(0));
        trustedBillProviders[provider] = trusted;
    }

    function setKYCVerification(address user, uint8 level, bytes32 verificationHash, uint48 expiryDate)
        external
        onlyCRE
    {
        require(user != address(0));
        require(user != address(this));
        require(level <= 2);
        require(expiryDate > block.timestamp);
        require(expiryDate <= block.timestamp + 365 days);
        require(verificationHash != bytes32(0));

        kycInfo[user] = KYCInfo({
            isVerified: true,
            verifiedAt: uint48(block.timestamp),
            verificationHash: verificationHash,
            verificationLevel: level,
            expiryDate: expiryDate
        });

        emit KYCVerified(user, level, expiryDate);
    }

    function revokeKYC(address user) external onlyOwner {
        kycInfo[user].isVerified = false;
    }

    function mint(
        address to,
        string calldata houseId,
        bytes32 documentHash,
        string calldata documentURI,
        DocumentStorageType storageType,
        string calldata verificationData
    ) external onlyCRE nonReentrant returns (uint256) {
        require(to != address(0));
        require(bytes(houseId).length > 0);
        require(bytes(houseId).length <= 100);
        require(documentHash != bytes32(0));
        require(bytes(documentURI).length > 0);
        require(bytes(verificationData).length > 0);

        uint256 tokenId = nextTokenId++;

        houses[tokenId] = House({
            houseId: houseId,
            documentHash: documentHash,
            documentURI: documentURI,
            storageType: storageType,
            originalOwner: to,
            mintedAt: uint48(block.timestamp),
            isVerified: true,
            documentCount: 1
        });

        _safeMint(to, tokenId);

        emit HouseMinted(tokenId, to, houseId, documentHash, storageType);

        return tokenId;
    }

    function createListing(
        uint256 tokenId,
        ListingType listingType,
        uint96 price,
        address preferredToken,
        bool isPrivateSale,
        address allowedBuyer,
        uint48 durationDays
    ) external onlyHouseOwner(tokenId) validToken(tokenId) nonReentrant {
        require(_hasValidKYC(msg.sender, 1));
        _setListing(tokenId, listingType, price, preferredToken, isPrivateSale, allowedBuyer, durationDays);
    }

    function createListingFromWorkflow(
        uint256 tokenId,
        address owner,
        ListingType listingType,
        uint96 price,
        address preferredToken,
        bool isPrivateSale,
        address allowedBuyer,
        uint48 durationDays
    ) external onlyCRE validToken(tokenId) nonReentrant {
        require(owner != address(0));
        require(ownerOf(tokenId) == owner);
        require(_hasValidKYC(owner, 1));
        _setListing(tokenId, listingType, price, preferredToken, isPrivateSale, allowedBuyer, durationDays);
    }

    function cancelListing(uint256 tokenId) external onlyHouseOwner(tokenId) validToken(tokenId) nonReentrant {
        delete listings[tokenId];
        emit HouseListed(tokenId, ListingType.NONE, 0, false);
    }

    function completeSale(uint256 tokenId, address buyer, bytes32 keyHash, bytes calldata encryptedKey)
        external
        onlyCRE
        validToken(tokenId)
        nonReentrant
    {
        Listing memory listing = listings[tokenId];
        require(listing.listingType == ListingType.FOR_SALE);
        require(listing.expiresAt > block.timestamp);
        require(buyer != address(0));
        require(keyHash != bytes32(0));
        require(encryptedKey.length > 0);

        if (listing.isPrivateSale) {
            require(buyer == listing.allowedBuyer);
        }

        address seller = ownerOf(tokenId);
        require(seller != buyer);
        require(_hasValidKYC(seller, 1));
        require(_hasValidKYC(buyer, 1));
        require(keyExchanges[keyHash].createdAt == 0);

        uint256 fee = (uint256(listing.price) * uint256(listing.platformFee)) / BPS_DENOMINATOR;

        keyExchanges[keyHash] = KeyExchange({
            keyHash: keyHash,
            encryptedKey: encryptedKey,
            intendedRecipient: buyer,
            createdAt: uint48(block.timestamp),
            expiresAt: uint48(block.timestamp + KEY_EXCHANGE_EXPIRY),
            isClaimed: false,
            exchangeType: 0
        });

        _transfer(seller, buyer, tokenId);
        delete listings[tokenId];
        totalFeesCollected += fee;

        emit HouseSold(tokenId, seller, buyer, listing.price, fee);
        emit KeyExchangeCreated(keyHash, tokenId, buyer, block.timestamp + KEY_EXCHANGE_EXPIRY);
    }

    function claimKey(bytes32 keyHash) external nonReentrant returns (bytes memory) {
        KeyExchange storage exchange = keyExchanges[keyHash];

        require(exchange.intendedRecipient == msg.sender);
        require(!exchange.isClaimed);
        require(block.timestamp <= exchange.expiresAt);

        exchange.isClaimed = true;

        emit KeyClaimed(keyHash, msg.sender, block.timestamp);

        return exchange.encryptedKey;
    }

    function startRental(
        uint256 tokenId,
        address renter,
        uint48 durationDays,
        uint96 depositAmount,
        uint96 monthlyRent,
        bytes calldata encryptedAccessKey
    ) external onlyCRE validToken(tokenId) nonReentrant {
        require(!rentals[tokenId].isActive);

        Listing memory listing = listings[tokenId];
        require(listing.listingType == ListingType.FOR_RENT);
        require(listing.expiresAt > block.timestamp);
        if (listing.isPrivateSale) {
            require(renter == listing.allowedBuyer);
        }

        require(renter != address(0));
        require(ownerOf(tokenId) != renter);
        require(durationDays > 0 && durationDays <= MAX_RENTAL_DURATION / 1 days);
        require(monthlyRent == listing.price);
        require(encryptedAccessKey.length > 0);
        require(_hasValidKYC(ownerOf(tokenId), 1), "HouseRWA: Owner KYC required");
        require(_hasValidKYC(renter, 1));
        require(pendingRentalDeposits[tokenId][renter] >= uint256(depositAmount));

        uint48 endTime = uint48(block.timestamp + (durationDays * 1 days));

        rentals[tokenId] = RentalAgreement({
            renter: renter,
            startTime: uint48(block.timestamp),
            endTime: endTime,
            depositAmount: depositAmount,
            monthlyRent: monthlyRent,
            isActive: true,
            encryptedAccessKeyHash: keccak256(encryptedAccessKey),
            disputeStatus: DisputeStatus.NONE
        });

        pendingRentalDeposits[tokenId][renter] -= uint256(depositAmount);
        delete listings[tokenId];

        bytes32 keyHash = keccak256(encryptedAccessKey);
        require(keyExchanges[keyHash].createdAt == 0);
        keyExchanges[keyHash] = KeyExchange({
            keyHash: keyHash,
            encryptedKey: encryptedAccessKey,
            intendedRecipient: renter,
            createdAt: uint48(block.timestamp),
            expiresAt: uint48(block.timestamp + KEY_EXCHANGE_EXPIRY),
            isClaimed: false,
            exchangeType: 1
        });

        emit RentalStarted(tokenId, renter, block.timestamp, endTime, depositAmount);
        emit KeyExchangeCreated(keyHash, tokenId, renter, block.timestamp + KEY_EXCHANGE_EXPIRY);
    }

    function depositForRental(uint256 tokenId) external payable validToken(tokenId) nonReentrant {
        require(!rentals[tokenId].isActive);
        require(msg.value > 0);
        require(_hasValidKYC(msg.sender, 1));

        Listing memory listing = listings[tokenId];
        require(listing.listingType == ListingType.FOR_RENT);
        require(listing.expiresAt > block.timestamp);
        if (listing.isPrivateSale) {
            require(msg.sender == listing.allowedBuyer);
        }

        pendingRentalDeposits[tokenId][msg.sender] += msg.value;
        emit RentalDepositReceived(tokenId, msg.sender, msg.value);
    }

    function withdrawRentalDeposit(uint256 tokenId) external nonReentrant {
        require(!rentals[tokenId].isActive);

        uint256 amount = pendingRentalDeposits[tokenId][msg.sender];
        require(amount > 0);

        pendingRentalDeposits[tokenId][msg.sender] = 0;
        (bool success,) = payable(msg.sender).call{value: amount}("");
        require(success);

        emit RentalDepositWithdrawn(tokenId, msg.sender, amount);
    }

    function endRental(uint256 tokenId) external validToken(tokenId) nonReentrant {
        RentalAgreement storage rental = rentals[tokenId];
        require(rental.isActive);

        bool canEnd = msg.sender == ownerOf(tokenId) || msg.sender == rental.renter
            || authorizedCREWorkflows[msg.sender] || block.timestamp >= rental.endTime;

        require(canEnd);

        address renter = rental.renter;
        uint256 deposit = rental.depositAmount;

        rental.isActive = false;

        (bool success,) = payable(renter).call{value: deposit}("");
        require(success);

        emit RentalEnded(tokenId, renter, deposit);
    }

    function createBill(
        uint256 tokenId,
        string calldata billType,
        uint96 amount,
        uint48 dueDate,
        address provider,
        bool isRecurring,
        uint8 recurrenceInterval
    ) external validToken(tokenId) returns (uint256 billIndex) {
        bool isCRE = authorizedCREWorkflows[msg.sender];
        bool isHouseOwner = msg.sender == ownerOf(tokenId);
        bool isTrustedProvider = trustedBillProviders[msg.sender] && msg.sender == provider;
        require(isCRE || isHouseOwner || isTrustedProvider);
        require(provider != address(0));
        require(bytes(billType).length > 0 && bytes(billType).length <= 50);
        require(amount > 0);
        require(dueDate > block.timestamp);
        require(dueDate <= block.timestamp + MAX_BILL_DUE_WINDOW);
        if (isRecurring) {
            require(recurrenceInterval > 0);
            require(recurrenceInterval <= MAX_RECURRENCE_INTERVAL_DAYS);
        } else {
            require(recurrenceInterval == 0);
        }
        require(houseBills[tokenId].length < MAX_BILLS_PER_HOUSE);

        billIndex = houseBills[tokenId].length;

        houseBills[tokenId].push(
            Bill({
                billType: billType,
                amount: amount,
                dueDate: dueDate,
                paidAt: 0,
                status: BillStatus.PENDING,
                paymentReference: bytes32(0),
                isRecurring: isRecurring,
                provider: provider,
                recurrenceInterval: recurrenceInterval
            })
        );

        emit BillCreated(tokenId, billIndex, billType, amount, dueDate);
    }

    function recordBillPayment(
        uint256 tokenId,
        uint256 billIndex,
        string calldata paymentMethod,
        bytes32 paymentReference
    ) external onlyCRE validToken(tokenId) nonReentrant {
        require(billIndex < houseBills[tokenId].length);
        require(bytes(paymentMethod).length > 0);
        require(paymentReference != bytes32(0));

        Bill storage bill = houseBills[tokenId][billIndex];
        require(bill.status == BillStatus.PENDING || bill.status == BillStatus.OVERDUE);

        bill.status = BillStatus.PAID;
        bill.paidAt = uint48(block.timestamp);
        bill.paymentReference = paymentReference;

        if (bill.isRecurring && bill.recurrenceInterval > 0) {
            require(houseBills[tokenId].length < MAX_BILLS_PER_HOUSE);
            uint48 nextDueDate = uint48(bill.dueDate + (bill.recurrenceInterval * 1 days));

            houseBills[tokenId].push(
                Bill({
                    billType: bill.billType,
                    amount: bill.amount,
                    dueDate: nextDueDate,
                    paidAt: 0,
                    status: BillStatus.PENDING,
                    paymentReference: bytes32(0),
                    isRecurring: true,
                    provider: bill.provider,
                    recurrenceInterval: bill.recurrenceInterval
                })
            );

            emit BillCreated(tokenId, houseBills[tokenId].length - 1, bill.billType, bill.amount, nextDueDate);
        }

        emit BillPaid(tokenId, billIndex, paymentMethod, paymentReference);
    }

    function disputeBill(uint256 tokenId, uint256 billIndex, string calldata reason)
        external
        onlyHouseOwner(tokenId)
        validToken(tokenId)
    {
        require(billIndex < houseBills[tokenId].length);

        Bill storage bill = houseBills[tokenId][billIndex];
        require(bill.status == BillStatus.PENDING);

        bill.status = BillStatus.DISPUTED;

        emit BillDisputed(tokenId, billIndex, msg.sender, reason);
    }

    function getHouseDetails(uint256 tokenId) external view validToken(tokenId) returns (House memory) {
        return houses[tokenId];
    }

    function getListing(uint256 tokenId) external view validToken(tokenId) returns (Listing memory) {
        return listings[tokenId];
    }

    function getBills(uint256 tokenId) external view validToken(tokenId) returns (Bill[] memory) {
        return houseBills[tokenId];
    }

    function getActiveRental(uint256 tokenId) external view validToken(tokenId) returns (RentalAgreement memory) {
        return rentals[tokenId];
    }

    function isRented(uint256 tokenId) external view validToken(tokenId) returns (bool) {
        return rentals[tokenId].isActive;
    }

    function getTotalBillsCount(uint256 tokenId) external view validToken(tokenId) returns (uint256) {
        return houseBills[tokenId].length;
    }

    function hasKYC(address user) external view returns (bool) {
        KYCInfo memory info = kycInfo[user];
        return info.isVerified && info.expiryDate > block.timestamp;
    }

    function _exists(uint256 tokenId) internal view returns (bool) {
        return _ownerOf(tokenId) != address(0);
    }

    function _hasValidKYC(address user, uint8 minLevel) internal view returns (bool) {
        KYCInfo memory info = kycInfo[user];
        return info.isVerified && info.verificationLevel >= minLevel && info.expiryDate > block.timestamp;
    }

    function _setListing(
        uint256 tokenId,
        ListingType listingType,
        uint96 price,
        address preferredToken,
        bool isPrivateSale,
        address allowedBuyer,
        uint48 durationDays
    ) internal {
        require(listingType != ListingType.NONE);
        require(price > 0);
        require(durationDays == 0 || durationDays <= uint48(MAX_RENTAL_DURATION / 1 days));
        require(!rentals[tokenId].isActive || rentals[tokenId].renter == address(0));
        if (isPrivateSale) {
            require(allowedBuyer != address(0));
        }

        uint48 expiry =
            durationDays > 0 ? uint48(block.timestamp + (durationDays * 1 days)) : uint48(block.timestamp + 30 days);

        listings[tokenId] = Listing({
            listingType: listingType,
            price: price,
            preferredToken: preferredToken,
            isPrivateSale: isPrivateSale,
            allowedBuyer: isPrivateSale ? allowedBuyer : address(0),
            createdAt: uint48(block.timestamp),
            expiresAt: expiry,
            platformFee: uint8(PROTOCOL_FEE_BPS)
        });

        emit HouseListed(tokenId, listingType, price, isPrivateSale);
    }

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

    receive() external payable {}
}
