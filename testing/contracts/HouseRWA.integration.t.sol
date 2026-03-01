// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test, console} from "forge-std/Test.sol";
import {HouseRWA} from "src/HouseRWA.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";
import {IERC721} from "@openzeppelin/contracts/token/ERC721/IERC721.sol";

/**
 * @title HouseRWAIntegrationTest
 * @notice Comprehensive integration test suite for HouseRWA contract
 * @dev Tests all critical paths including security, access controls, and edge cases
 */
contract HouseRWAIntegrationTest is Test {
    HouseRWA public implementation;
    HouseRWA public houseRWA;
    ERC1967Proxy public proxy;
    
    // Test accounts
    address public owner;
    address public creWorkflow;
    address public feeRecipient;
    address public seller;
    address public buyer;
    address public tenant;
    address public provider;
    address public arbitrator;
    address public validator1;
    address public validator2;
    address public unauthorized;

    // Test constants
    uint256 constant INITIAL_BALANCE = 100 ether;
    uint256 constant LISTING_PRICE = 10 ether;
    uint256 constant RENTAL_DEPOSIT = 1 ether;
    uint256 constant MONTHLY_RENT = 0.5 ether;
    
    function makeAddr(string memory name) internal returns (address addr) {
        addr = vm.addr(uint256(keccak256(bytes(name))));
        vm.label(addr, name);
    }
    
    function setUp() public {
        // Setup test accounts
        owner = makeAddr("owner");
        creWorkflow = makeAddr("creWorkflow");
        feeRecipient = makeAddr("feeRecipient");
        seller = makeAddr("seller");
        buyer = makeAddr("buyer");
        tenant = makeAddr("tenant");
        provider = makeAddr("provider");
        arbitrator = makeAddr("arbitrator");
        validator1 = makeAddr("validator1");
        validator2 = makeAddr("validator2");
        unauthorized = makeAddr("unauthorized");
        
        vm.startPrank(owner);
        
        // Deploy implementation
        implementation = new HouseRWA();
        
        // Deploy proxy with initialization
        bytes memory initData = abi.encodeWithSelector(
            HouseRWA.initialize.selector,
            owner,
            feeRecipient,
            creWorkflow
        );
        
        proxy = new ERC1967Proxy(address(implementation), initData);
        houseRWA = HouseRWA(payable(address(proxy)));
        
        // Set up trusted providers and arbitrators
        houseRWA.setTrustedBillProvider(provider, true);
        houseRWA.setArbitrator(arbitrator, true);
        
        vm.stopPrank();

        // Baseline KYC for seller so listing flows can execute.
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(seller, 1, keccak256("seller-kyc"), uint48(block.timestamp + 365 days));
        
        // Fund test accounts
        vm.deal(seller, INITIAL_BALANCE);
        vm.deal(buyer, INITIAL_BALANCE);
        vm.deal(tenant, INITIAL_BALANCE);
        vm.deal(validator1, INITIAL_BALANCE);
        vm.deal(validator2, INITIAL_BALANCE);
    }

    // ============ MINTING FLOW TESTS ============
    
    /**
     * @notice Test complete minting flow with encrypted documents
     */
    function test_MintHouse_WithEncryptedDocuments() public {
        string memory houseId = "house-001";
        bytes32 documentHash = keccak256("encrypted-documents");
        string memory documentURI = "ipfs://QmTestHash123";
        
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            seller,
            houseId,
            documentHash,
            documentURI,
            HouseRWA.DocumentStorageType.IPFS,
            "CRE-verification-proof"
        );
        
        // Verify minting succeeded
        assertEq(tokenId, 0, "First token should have ID 0");
        assertEq(houseRWA.ownerOf(tokenId), seller, "Seller should own the token");
        
        // Verify house data
        HouseRWA.House memory house = houseRWA.getHouseDetails(tokenId);
        assertEq(house.houseId, houseId, "House ID mismatch");
        assertEq(house.documentHash, documentHash, "Document hash mismatch");
        assertEq(house.documentURI, documentURI, "Document URI mismatch");
        assertEq(uint256(house.storageType), uint256(HouseRWA.DocumentStorageType.IPFS), "Storage type mismatch");
        assertEq(house.originalOwner, seller, "Original owner mismatch");
        assertTrue(house.isVerified, "House should be verified");
        
        // Verify enumerable interface
        assertEq(houseRWA.totalSupply(), 1, "Total supply should be 1");
        assertEq(houseRWA.balanceOf(seller), 1, "Seller balance should be 1");
    }
    
    /**
     * @notice Test minting is rejected without proper authorization
     */
    function test_Revert_MintUnauthorized() public {
        vm.prank(unauthorized);
        vm.expectRevert("HouseRWA: Unauthorized CRE workflow");
        houseRWA.mint(
            seller,
            "house-001",
            keccak256("docs"),
            "ipfs://test",
            HouseRWA.DocumentStorageType.IPFS,
            "verification"
        );
    }
    
    /**
     * @notice Test minting is rejected when paused
     */
    function test_Revert_MintWhenPaused() public {
        vm.prank(owner);
        houseRWA.setMintingPaused(true, "Emergency maintenance");
        
        vm.prank(creWorkflow);
        vm.expectRevert("HouseRWA: Minting paused");
        houseRWA.mint(
            seller,
            "house-001",
            keccak256("docs"),
            "ipfs://test",
            HouseRWA.DocumentStorageType.IPFS,
            "verification"
        );
    }
    
    /**
     * @notice Test multiple mints with different storage types
     */
    function test_Mint_MultipleHousesWithDifferentStorage() public {
        // IPFS storage
        vm.prank(creWorkflow);
        uint256 token1 = houseRWA.mint(
            seller,
            "house-ipfs",
            keccak256("ipfs-docs"),
            "ipfs://QmHash1",
            HouseRWA.DocumentStorageType.IPFS,
            "CRE-verified"
        );
        
        // Off-chain storage
        vm.prank(creWorkflow);
        uint256 token2 = houseRWA.mint(
            seller,
            "house-offchain",
            keccak256("offchain-docs"),
            "https://api.example.com/docs/2",
            HouseRWA.DocumentStorageType.OFF_CHAIN,
            "CRE-verified"
        );
        
        // Arweave storage
        vm.prank(creWorkflow);
        uint256 token3 = houseRWA.mint(
            seller,
            "house-arweave",
            keccak256("arweave-docs"),
            "arweave://TxId123",
            HouseRWA.DocumentStorageType.ARWEAVE,
            "CRE-verified"
        );
        
        assertEq(token1, 0);
        assertEq(token2, 1);
        assertEq(token3, 2);
        
        HouseRWA.House memory house1 = houseRWA.getHouseDetails(token1);
        HouseRWA.House memory house2 = houseRWA.getHouseDetails(token2);
        HouseRWA.House memory house3 = houseRWA.getHouseDetails(token3);
        
        assertEq(uint256(house1.storageType), uint256(HouseRWA.DocumentStorageType.IPFS));
        assertEq(uint256(house2.storageType), uint256(HouseRWA.DocumentStorageType.OFF_CHAIN));
        assertEq(uint256(house3.storageType), uint256(HouseRWA.DocumentStorageType.ARWEAVE));
    }

    // ============ PRIVATE SALE FLOW TESTS ============
    
    /**
     * @notice Test complete private sale flow with document transfer
     */
    function test_PrivateSale_CompleteFlow() public {
        // Setup: Mint house and verify KYC for buyer
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            seller,
            "house-sale",
            keccak256("docs"),
            "ipfs://test",
            HouseRWA.DocumentStorageType.IPFS,
            "verified"
        );
        
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(buyer, 1, keccak256("kyc-data"), uint48(block.timestamp + 365 days));
        
        // Step 1: Create private listing
        vm.prank(seller);
        houseRWA.createListing(
            tokenId,
            HouseRWA.ListingType.FOR_SALE,
            uint96(LISTING_PRICE),
            address(0), // ETH preferred
            true, // Private sale
            buyer, // Only this buyer can purchase
            30 // 30 days expiry
        );
        
        HouseRWA.Listing memory listing = houseRWA.getListing(tokenId);
        assertTrue(listing.isPrivateSale, "Should be private sale");
        assertEq(listing.allowedBuyer, buyer, "Allowed buyer mismatch");
        
        // Step 2: Complete sale with key exchange
        bytes32 keyHash = keccak256("decryption-key");
        bytes memory encryptedKey = abi.encodePacked("encrypted-key-data-for-buyer");
        
        vm.prank(creWorkflow);
        houseRWA.completeSale(tokenId, buyer, keyHash, encryptedKey);
        
        // Step 3: Verify ownership transfer
        assertEq(houseRWA.ownerOf(tokenId), buyer, "Buyer should now own the token");
        
        // Step 4: Buyer claims decryption key
        vm.prank(buyer);
        bytes memory claimedKey = houseRWA.claimKey(keyHash);
        assertEq(claimedKey, encryptedKey, "Claimed key should match");
        
        // Step 5: Verify listing is cleared
        listing = houseRWA.getListing(tokenId);
        assertEq(uint256(listing.listingType), uint256(HouseRWA.ListingType.NONE), "Listing should be cleared");
    }
    
    /**
     * @notice Test private sale prevents unauthorized purchases
     */
    function test_Revert_PrivateSaleUnauthorizedBuyer() public {
        address unauthorizedBuyer = makeAddr("unauthorizedBuyer");
        
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            seller,
            "house-private",
            keccak256("docs"),
            "ipfs://test",
            HouseRWA.DocumentStorageType.IPFS,
            "verified"
        );
        
        vm.prank(seller);
        houseRWA.createListing(
            tokenId,
            HouseRWA.ListingType.FOR_SALE,
            uint96(LISTING_PRICE),
            address(0),
            true,
            buyer, // Only specific buyer
            30
        );
        
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(unauthorizedBuyer, 1, keccak256("kyc"), uint48(block.timestamp + 365 days));
        
        vm.prank(creWorkflow);
        vm.expectRevert("HouseRWA: Not allowed buyer");
        houseRWA.completeSale(tokenId, unauthorizedBuyer, keccak256("key"), abi.encodePacked("data"));
    }
    
    /**
     * @notice Test slashing mechanism for disputes
     */
    function test_SaleDisputeAndSlashing() public {
        // Setup validator staking
        vm.prank(validator1);
        houseRWA.stakeAsValidator{value: 2 ether}();
        
        (uint96 stakedAmount,,,,,) = houseRWA.validators(validator1);
        assertEq(stakedAmount, 2 ether, "Validator should have staked 2 ETH");
        
        // Slash validator for misconduct
        uint256 slashAmount = 1 ether;
        uint256 recipientBalanceBefore = feeRecipient.balance;
        
        vm.prank(owner);
        houseRWA.slashValidator(validator1, slashAmount, "Fraudulent sale verification");
        
        bool isSlashed;
        (stakedAmount,,, isSlashed,,) = houseRWA.validators(validator1);
        assertEq(stakedAmount, 1 ether, "Stake should be reduced");
        assertTrue(isSlashed, "Validator should be marked as slashed");
        assertEq(feeRecipient.balance, recipientBalanceBefore + slashAmount, "Fee recipient should receive slashed amount");
    }

    // ============ RENTAL FLOW TESTS ============
    
    /**
     * @notice Test complete rental agreement creation
     */
    function test_Rental_CreateAgreement() public {
        // Setup
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            seller,
            "house-rental",
            keccak256("docs"),
            "ipfs://test",
            HouseRWA.DocumentStorageType.IPFS,
            "verified"
        );
        
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(tenant, 1, keccak256("kyc"), uint48(block.timestamp + 365 days));

        vm.prank(seller);
        houseRWA.createListing(
            tokenId,
            HouseRWA.ListingType.FOR_RENT,
            uint96(MONTHLY_RENT),
            address(0),
            false,
            address(0),
            30
        );

        vm.prank(tenant);
        houseRWA.depositForRental{value: RENTAL_DEPOSIT}(tokenId);
        
        bytes memory encryptedAccessKey = abi.encodePacked("access-key-for-tenant");
        
        // Create rental
        vm.prank(creWorkflow);
        houseRWA.startRental(
            tokenId,
            tenant,
            30, // 30 days
            uint96(RENTAL_DEPOSIT),
            uint96(MONTHLY_RENT),
            encryptedAccessKey
        );
        
        // Verify rental agreement
        HouseRWA.RentalAgreement memory rental = houseRWA.getActiveRental(tokenId);
        assertTrue(rental.isActive, "Rental should be active");
        assertEq(rental.renter, tenant, "Renter mismatch");
        assertEq(rental.depositAmount, RENTAL_DEPOSIT, "Deposit mismatch");
        assertEq(rental.monthlyRent, MONTHLY_RENT, "Monthly rent mismatch");
        assertEq(rental.encryptedAccessKeyHash, keccak256(encryptedAccessKey), "Access key hash mismatch");
    }
    
    /**
     * @notice Test rental ending and deposit return
     */
    function test_Rental_EndAndReturnDeposit() public {
        // Setup rental
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            seller,
            "house-rental",
            keccak256("docs"),
            "ipfs://test",
            HouseRWA.DocumentStorageType.IPFS,
            "verified"
        );
        
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(tenant, 1, keccak256("kyc"), uint48(block.timestamp + 365 days));

        vm.prank(seller);
        houseRWA.createListing(
            tokenId,
            HouseRWA.ListingType.FOR_RENT,
            uint96(MONTHLY_RENT),
            address(0),
            false,
            address(0),
            30
        );

        vm.prank(tenant);
        houseRWA.depositForRental{value: RENTAL_DEPOSIT}(tokenId);
        
        vm.prank(creWorkflow);
        houseRWA.startRental(tokenId, tenant, 30, uint96(RENTAL_DEPOSIT), uint96(MONTHLY_RENT), abi.encodePacked("key"));
        
        uint256 tenantBalanceBefore = tenant.balance;
        
        // Warp past rental end
        vm.warp(block.timestamp + 31 days);
        
        // End rental
        vm.prank(seller);
        houseRWA.endRental(tokenId);
        
        // Verify rental ended and deposit returned
        HouseRWA.RentalAgreement memory rental = houseRWA.getActiveRental(tokenId);
        assertFalse(rental.isActive, "Rental should not be active");
        assertEq(tenant.balance, tenantBalanceBefore + RENTAL_DEPOSIT, "Tenant should receive deposit back");
    }
    
    /**
     * @notice Test rental dispute and resolution
     */
    function test_Rental_DisputeAndResolution() public {
        // Setup
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            seller,
            "house-dispute",
            keccak256("docs"),
            "ipfs://test",
            HouseRWA.DocumentStorageType.IPFS,
            "verified"
        );
        
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(tenant, 1, keccak256("kyc"), uint48(block.timestamp + 365 days));

        vm.prank(seller);
        houseRWA.createListing(
            tokenId,
            HouseRWA.ListingType.FOR_RENT,
            uint96(MONTHLY_RENT),
            address(0),
            false,
            address(0),
            30
        );

        vm.prank(tenant);
        houseRWA.depositForRental{value: RENTAL_DEPOSIT}(tokenId);
        
        vm.prank(creWorkflow);
        houseRWA.startRental(tokenId, tenant, 30, uint96(RENTAL_DEPOSIT), uint96(MONTHLY_RENT), abi.encodePacked("key"));
        
        // Open dispute
        vm.prank(tenant);
        houseRWA.openRentalDispute(tokenId, "Property damage dispute");
        
        HouseRWA.RentalAgreement memory rental = houseRWA.getActiveRental(tokenId);
        assertEq(uint256(rental.disputeStatus), uint256(HouseRWA.DisputeStatus.OPEN), "Dispute should be open");
        
        // Resolve dispute (partial deposit to each party)
        uint256 toOwner = RENTAL_DEPOSIT * 60 / 100;
        uint256 toRenter = RENTAL_DEPOSIT * 40 / 100;
        
        uint256 ownerBalanceBefore = seller.balance;
        uint256 tenantBalanceBefore = tenant.balance;
        
        vm.prank(arbitrator);
        houseRWA.resolveRentalDispute(tokenId, toOwner, toRenter);
        
        rental = houseRWA.getActiveRental(tokenId);
        assertFalse(rental.isActive, "Rental should be resolved");
        assertEq(uint256(rental.disputeStatus), uint256(HouseRWA.DisputeStatus.RESOLVED), "Dispute should be resolved");
        assertEq(seller.balance, ownerBalanceBefore + toOwner, "Owner should receive partial deposit");
        assertEq(tenant.balance, tenantBalanceBefore + toRenter, "Tenant should receive partial deposit");
    }

    // ============ BILL PAYMENT FLOW TESTS ============
    
    /**
     * @notice Test bill creation and payment recording
     */
    function test_Bill_CreateAndPay() public {
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            seller,
            "house-bills",
            keccak256("docs"),
            "ipfs://test",
            HouseRWA.DocumentStorageType.IPFS,
            "verified"
        );
        
        // Create bill
        vm.prank(provider);
        uint256 billIndex = houseRWA.createBill(
            tokenId,
            "electricity",
            15000, // $150.00 in cents
            uint48(block.timestamp + 30 days),
            provider,
            false,
            0
        );
        
        HouseRWA.Bill[] memory bills = houseRWA.getBills(tokenId);
        assertEq(bills.length, 1, "Should have 1 bill");
        assertEq(bills[0].billType, "electricity", "Bill type mismatch");
        assertEq(bills[0].amount, 15000, "Bill amount mismatch");
        assertEq(uint256(bills[0].status), uint256(HouseRWA.BillStatus.PENDING), "Bill should be pending");
        
        // Record payment
        bytes32 paymentRef = keccak256("stripe-payment-id");
        
        vm.prank(creWorkflow);
        houseRWA.recordBillPayment(tokenId, billIndex, "stripe", paymentRef);
        
        bills = houseRWA.getBills(tokenId);
        assertEq(uint256(bills[0].status), uint256(HouseRWA.BillStatus.PAID), "Bill should be paid");
        assertEq(bills[0].paymentReference, paymentRef, "Payment reference mismatch");
        assertGt(bills[0].paidAt, 0, "Paid timestamp should be set");
    }
    
    /**
     * @notice Test recurring bill creation
     */
    function test_Bill_Recurring() public {
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            seller,
            "house-recurring",
            keccak256("docs"),
            "ipfs://test",
            HouseRWA.DocumentStorageType.IPFS,
            "verified"
        );
        
        // Create recurring bill (monthly)
        vm.prank(provider);
        uint256 billIndex = houseRWA.createBill(
            tokenId,
            "internet",
            5000, // $50.00
            uint48(block.timestamp + 30 days),
            provider,
            true,
            30 // 30 days interval
        );
        
        // Pay first bill
        vm.prank(creWorkflow);
        houseRWA.recordBillPayment(tokenId, billIndex, "stripe", keccak256("payment-1"));
        
        // Verify next bill was created
        HouseRWA.Bill[] memory bills = houseRWA.getBills(tokenId);
        assertEq(bills.length, 2, "Should have 2 bills (paid + next)");
        assertEq(uint256(bills[0].status), uint256(HouseRWA.BillStatus.PAID), "First bill should be paid");
        assertEq(uint256(bills[1].status), uint256(HouseRWA.BillStatus.PENDING), "Next bill should be pending");
        assertEq(bills[1].dueDate, bills[0].dueDate + 30 days, "Next due date should be 30 days later");
    }

    // ============ KYC/AML FLOW TESTS ============
    
    /**
     * @notice Test KYC verification flow
     */
    function test_KYC_VerificationFlow() public {
        bytes32 kycHash = keccak256("kyc-verification-data");
        uint48 expiry = uint48(block.timestamp + 365 days);
        
        // Set KYC level 1
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(buyer, 1, kycHash, expiry);
        
        assertTrue(houseRWA.hasKYC(buyer), "Buyer should have KYC");
        
        // Upgrade to level 2
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(buyer, 2, kycHash, expiry);
        
        // Test expiry
        vm.warp(block.timestamp + 366 days);
        assertFalse(houseRWA.hasKYC(buyer), "KYC should expire after 366 days");
    }
    
    /**
     * @notice Test high-value transaction requires higher KYC
     */
    function test_KYC_HighValueTransaction() public {
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            seller,
            "house-highvalue",
            keccak256("docs"),
            "ipfs://test",
            HouseRWA.DocumentStorageType.IPFS,
            "verified"
        );
        
        // Set KYC level 1 for buyer
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(buyer, 1, keccak256("kyc"), uint48(block.timestamp + 365 days));
        
        // Try to list at high value (requires level 2 KYC)
        uint96 highPrice = uint96(500 ether); // Well above $100k threshold
        
        vm.prank(seller);
        vm.expectRevert("HouseRWA: High value KYC required");
        houseRWA.createListing(
            tokenId,
            HouseRWA.ListingType.FOR_SALE,
            highPrice,
            address(0),
            false,
            address(0),
            30
        );
        
        // Upgrade to level 2 KYC
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(buyer, 2, keccak256("kyc"), uint48(block.timestamp + 365 days));
        
        // Also need to KYC the seller for high value
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(seller, 2, keccak256("kyc-seller"), uint48(block.timestamp + 365 days));
        
        // Now listing should succeed
        vm.prank(seller);
        houseRWA.createListing(
            tokenId,
            HouseRWA.ListingType.FOR_SALE,
            highPrice,
            address(0),
            false,
            address(0),
            30
        );
    }

    // ============ SECURITY TESTS ============
    
    /**
     * @notice Test reentrancy protection on stake function
     */
    function test_Security_ReentrancyProtection() public {
        // Attempt reentrancy through stake/unstake
        vm.prank(validator1);
        houseRWA.stakeAsValidator{value: 2 ether}();
        
        // Try to unstake immediately (should fail due to timelock)
        vm.prank(validator1);
        vm.expectRevert("HouseRWA: Stake timelock active");
        houseRWA.unstake(1 ether);
        
        // Warp past timelock
        vm.warp(block.timestamp + 31 days);
        
        // Now unstake should succeed
        vm.prank(validator1);
        houseRWA.unstake(1 ether);
        
        (uint96 stakedAmount,,,,,) = houseRWA.validators(validator1);
        assertEq(stakedAmount, 1 ether, "Should have 1 ETH remaining");
    }
    
    /**
     * @notice Test circuit breaker functionality
     */
    function test_Security_CircuitBreakers() public {
        // Test individual circuit breakers
        vm.startPrank(owner);
        houseRWA.setMintingPaused(true, "Minting paused");
        houseRWA.setSalesPaused(true, "Sales paused");
        houseRWA.setRentalsPaused(true, "Rentals paused");
        houseRWA.setPaymentsPaused(true, "Payments paused");
        vm.stopPrank();
        
        assertTrue(houseRWA.mintingPaused(), "Minting should be paused");
        assertTrue(houseRWA.salesPaused(), "Sales should be paused");
        assertTrue(houseRWA.rentalsPaused(), "Rentals should be paused");
        assertTrue(houseRWA.paymentsPaused(), "Payments should be paused");

        vm.prank(owner);
        houseRWA.emergencyPause();
        assertTrue(houseRWA.paused(), "Contract should be paused");
        
        // Test emergency unpause
        vm.prank(owner);
        houseRWA.emergencyUnpause();
        
        assertFalse(houseRWA.mintingPaused(), "Minting should be unpaused");
        assertFalse(houseRWA.paused(), "Contract should be unpaused");
    }
    
    /**
     * @notice Test access control for admin functions
     */
    function test_Security_AccessControl() public {
        // Try to authorize workflow as non-owner
        address newWorkflow = makeAddr("newWorkflow");
        
        vm.prank(unauthorized);
        vm.expectRevert();
        houseRWA.authorizeCREWorkflow(newWorkflow);
        
        // Owner can authorize
        vm.prank(owner);
        houseRWA.authorizeCREWorkflow(newWorkflow);
        assertTrue(houseRWA.authorizedCREWorkflows(newWorkflow), "New workflow should be authorized");
        
        // Try to set KYC as unauthorized
        vm.prank(unauthorized);
        vm.expectRevert("HouseRWA: Unauthorized CRE workflow");
        houseRWA.setKYCVerification(buyer, 1, keccak256("kyc"), uint48(block.timestamp + 365 days));
    }
    
    /**
     * @notice Test key exchange expiry
     */
    function test_Security_KeyExpiry() public {
        // Setup sale
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(buyer, 1, keccak256("kyc"), uint48(block.timestamp + 365 days));
        
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(seller, "house-key", keccak256("docs"), "ipfs://test", HouseRWA.DocumentStorageType.IPFS, "verified");
        
        vm.prank(seller);
        houseRWA.createListing(tokenId, HouseRWA.ListingType.FOR_SALE, 1 ether, address(0), false, address(0), 30);
        
        bytes32 keyHash = keccak256("key");
        vm.prank(creWorkflow);
        houseRWA.completeSale(tokenId, buyer, keyHash, abi.encodePacked("encrypted"));
        
        // Warp past key expiry (7 days)
        vm.warp(block.timestamp + 8 days);
        
        // Try to claim expired key
        vm.prank(buyer);
        vm.expectRevert("HouseRWA: Key expired");
        houseRWA.claimKey(keyHash);
    }
    
    /**
     * @notice Test double-claim prevention
     */
    function test_Security_DoubleClaimPrevention() public {
        // Setup sale
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(buyer, 1, keccak256("kyc"), uint48(block.timestamp + 365 days));
        
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(seller, "house-double", keccak256("docs"), "ipfs://test", HouseRWA.DocumentStorageType.IPFS, "verified");
        
        vm.prank(seller);
        houseRWA.createListing(tokenId, HouseRWA.ListingType.FOR_SALE, 1 ether, address(0), false, address(0), 30);
        
        bytes32 keyHash = keccak256("key");
        vm.prank(creWorkflow);
        houseRWA.completeSale(tokenId, buyer, keyHash, abi.encodePacked("encrypted"));
        
        // First claim succeeds
        vm.prank(buyer);
        houseRWA.claimKey(keyHash);
        
        // Second claim fails
        vm.prank(buyer);
        vm.expectRevert("HouseRWA: Already claimed");
        houseRWA.claimKey(keyHash);
    }

    // ============ FUZZ TESTS ============
    
    /**
     * @notice Fuzz test minting with valid parameters
     */
    function testFuzz_MintWithValidData(
        address to,
        bytes32 docHash,
        uint8 storageType
    ) public {
        vm.assume(to != address(0));
        vm.assume(to != address(houseRWA));
        vm.assume(docHash != bytes32(0));
        vm.assume(storageType <= 3);
        
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            to,
            "fuzz-test",
            docHash,
            "ipfs://fuzz",
            HouseRWA.DocumentStorageType(storageType),
            "fuzz-verification"
        );
        
        assertEq(houseRWA.ownerOf(tokenId), to);
    }
    
    /**
     * @notice Fuzz test price boundaries
     */
    function testFuzz_ListingPrice(uint96 price) public {
        vm.assume(price > 0);

        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(seller, 2, keccak256("seller-kyc-high"), uint48(block.timestamp + 365 days));
        
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(seller, "fuzz-price", keccak256("docs"), "ipfs://test", HouseRWA.DocumentStorageType.IPFS, "verified");
        
        vm.prank(seller);
        houseRWA.createListing(
            tokenId,
            HouseRWA.ListingType.FOR_SALE,
            price,
            address(0),
            false,
            address(0),
            30
        );
        
        HouseRWA.Listing memory listing = houseRWA.getListing(tokenId);
        assertEq(listing.price, price);
    }

    // ============ EVENT EMISSION TESTS ============
    
    /**
     * @notice Test event emission for CRE integration
     */
    function test_Events_Minting() public {
        vm.prank(creWorkflow);
        
        vm.expectEmit(true, true, false, true);
        emit HouseRWA.HouseMinted(
            0,
            seller,
            "event-test",
            keccak256("docs"),
            HouseRWA.DocumentStorageType.IPFS
        );
        
        houseRWA.mint(
            seller,
            "event-test",
            keccak256("docs"),
            "ipfs://event",
            HouseRWA.DocumentStorageType.IPFS,
            "verified"
        );
    }
    
    /**
     * @notice Test sale event emission
     */
    function test_Events_Sale() public {
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(buyer, 1, keccak256("kyc"), uint48(block.timestamp + 365 days));
        
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(seller, "event-sale", keccak256("docs"), "ipfs://test", HouseRWA.DocumentStorageType.IPFS, "verified");
        
        vm.prank(seller);
        houseRWA.createListing(tokenId, HouseRWA.ListingType.FOR_SALE, 1 ether, address(0), false, address(0), 30);
        
        vm.expectEmit(true, true, true, false);
        emit HouseRWA.HouseSold(tokenId, seller, buyer, 1 ether, 0.025 ether);
        
        vm.prank(creWorkflow);
        houseRWA.completeSale(tokenId, buyer, keccak256("key"), abi.encodePacked("encrypted"));
    }

    // ============ GAS OPTIMIZATION TESTS ============
    
    /**
     * @notice Test gas usage for critical operations
     */
    function test_Gas_Minting() public {
        uint256 gasBefore = gasleft();
        
        vm.prank(creWorkflow);
        houseRWA.mint(
            seller,
            "gas-test",
            keccak256("docs"),
            "ipfs://gas",
            HouseRWA.DocumentStorageType.IPFS,
            "verified"
        );
        
        uint256 gasUsed = gasBefore - gasleft();
        console.log("Gas used for minting:", gasUsed);
        assertLt(gasUsed, 300000, "Minting gas should be under 300k");
    }
    
    function test_Gas_SaleCompletion() public {
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(buyer, 1, keccak256("kyc"), uint48(block.timestamp + 365 days));
        
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(seller, "gas-sale", keccak256("docs"), "ipfs://test", HouseRWA.DocumentStorageType.IPFS, "verified");
        
        vm.prank(seller);
        houseRWA.createListing(tokenId, HouseRWA.ListingType.FOR_SALE, 1 ether, address(0), false, address(0), 30);
        
        uint256 gasBefore = gasleft();
        
        vm.prank(creWorkflow);
        houseRWA.completeSale(tokenId, buyer, keccak256("key"), abi.encodePacked("encrypted"));
        
        uint256 gasUsed = gasBefore - gasleft();
        console.log("Gas used for sale completion:", gasUsed);
        assertLt(gasUsed, 200000, "Sale gas should be under 200k");
    }

    // ============ EDGE CASE TESTS ============
    
    /**
     * @notice Test zero address handling
     */
    function test_Edge_ZeroAddress() public {
        vm.prank(creWorkflow);
        vm.expectRevert("HouseRWA: Invalid recipient");
        houseRWA.mint(address(0), "test", keccak256("docs"), "ipfs://test", HouseRWA.DocumentStorageType.IPFS, "verified");
    }
    
    /**
     * @notice Test empty string handling
     */
    function test_Edge_EmptyStrings() public {
        vm.prank(creWorkflow);
        vm.expectRevert("HouseRWA: Invalid house ID");
        houseRWA.mint(seller, "", keccak256("docs"), "ipfs://test", HouseRWA.DocumentStorageType.IPFS, "verified");
    }
    
    /**
     * @notice Test extremely long house ID
     */
    function test_Edge_LongHouseId() public {
        string memory longId = "a";
        for (uint i = 0; i < 100; i++) {
            longId = string(abi.encodePacked(longId, "a"));
        }
        
        vm.prank(creWorkflow);
        vm.expectRevert("HouseRWA: House ID too long");
        houseRWA.mint(seller, longId, keccak256("docs"), "ipfs://test", HouseRWA.DocumentStorageType.IPFS, "verified");
    }
    
    /**
     * @notice Test same address as seller and buyer
     */
    function test_Edge_SameBuyerSeller() public {
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(seller, 1, keccak256("kyc"), uint48(block.timestamp + 365 days));
        
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(seller, "same", keccak256("docs"), "ipfs://test", HouseRWA.DocumentStorageType.IPFS, "verified");
        
        vm.prank(seller);
        houseRWA.createListing(tokenId, HouseRWA.ListingType.FOR_SALE, 1 ether, address(0), false, address(0), 30);
        
        vm.prank(creWorkflow);
        vm.expectRevert("HouseRWA: Cannot buy own house");
        houseRWA.completeSale(tokenId, seller, keccak256("key"), abi.encodePacked("data"));
    }
    
    receive() external payable {}
}
