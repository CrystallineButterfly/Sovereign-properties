// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {HouseRWA} from "../src/HouseRWA.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

contract HouseRWATest is Test {
    HouseRWA public implementation;
    HouseRWA public houseRWA;
    ERC1967Proxy public proxy;

    address public owner;
    address public creWorkflow;
    address public feeRecipient;
    address public user1;
    address public user2;
    address public provider;

    function makeAddr(string memory name) internal returns (address addr) {
        addr = vm.addr(uint256(keccak256(bytes(name))));
        vm.label(addr, name);
    }

    function setUp() public {
        owner = makeAddr("owner");
        creWorkflow = makeAddr("creWorkflow");
        feeRecipient = makeAddr("feeRecipient");
        user1 = makeAddr("user1");
        user2 = makeAddr("user2");
        provider = makeAddr("provider");

        implementation = new HouseRWA();

        bytes memory initData = abi.encodeWithSelector(HouseRWA.initialize.selector, owner, feeRecipient, creWorkflow);

        proxy = new ERC1967Proxy(address(implementation), initData);
        houseRWA = HouseRWA(payable(address(proxy)));

        vm.deal(user1, 100 ether);
        vm.deal(user2, 100 ether);

        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(user1, 1, keccak256("kyc-user1"), uint48(block.timestamp + 365 days));
    }

    function test_Initialize() public {
        assertEq(houseRWA.owner(), owner);
        assertEq(houseRWA.feeRecipient(), feeRecipient);
        assertTrue(houseRWA.authorizedCREWorkflows(creWorkflow));
    }

    function test_RevertInitializeTwice() public {
        vm.expectRevert();
        houseRWA.initialize(owner, feeRecipient, creWorkflow);
    }

    function test_MintHouse() public {
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            user1,
            "house-123",
            keccak256("encrypted-docs"),
            "ipfs://QmHash",
            HouseRWA.DocumentStorageType.IPFS,
            "CRE-verified"
        );

        assertEq(tokenId, 0);
        assertEq(houseRWA.ownerOf(tokenId), user1);
    }

    function test_RevertMintUnauthorized() public {
        vm.prank(user1);
        vm.expectRevert();
        houseRWA.mint(
            user1,
            "house-123",
            keccak256("encrypted-docs"),
            "ipfs://QmHash",
            HouseRWA.DocumentStorageType.IPFS,
            "CRE-verified"
        );
    }

    function test_CreateListingFromWorkflow() public {
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            user1,
            "house-listing",
            keccak256("encrypted-docs"),
            "ipfs://QmHash",
            HouseRWA.DocumentStorageType.IPFS,
            "CRE-verified"
        );

        vm.prank(creWorkflow);
        houseRWA.createListingFromWorkflow(
            tokenId, user1, HouseRWA.ListingType.FOR_SALE, 2 ether, address(0), false, address(0), 30
        );

        HouseRWA.Listing memory listing = houseRWA.getListing(tokenId);
        assertEq(uint8(listing.listingType), uint8(HouseRWA.ListingType.FOR_SALE));
        assertEq(listing.price, 2 ether);
        assertFalse(listing.isPrivateSale);
    }

    function test_RevertCreateListingFromWorkflowUnauthorized() public {
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            user1,
            "house-listing",
            keccak256("encrypted-docs"),
            "ipfs://QmHash",
            HouseRWA.DocumentStorageType.IPFS,
            "CRE-verified"
        );

        vm.prank(user1);
        vm.expectRevert();
        houseRWA.createListingFromWorkflow(
            tokenId, user1, HouseRWA.ListingType.FOR_SALE, 2 ether, address(0), false, address(0), 30
        );
    }

    function test_RevertCreateListingFromWorkflowWrongOwner() public {
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            user1,
            "house-listing",
            keccak256("encrypted-docs"),
            "ipfs://QmHash",
            HouseRWA.DocumentStorageType.IPFS,
            "CRE-verified"
        );

        vm.prank(creWorkflow);
        vm.expectRevert();
        houseRWA.createListingFromWorkflow(
            tokenId, user2, HouseRWA.ListingType.FOR_SALE, 2 ether, address(0), false, address(0), 30
        );
    }

    function test_CreateListingAndCompleteSale() public {
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(user2, 1, keccak256("kyc-buyer"), uint48(block.timestamp + 365 days));

        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            user1,
            "house-123",
            keccak256("encrypted-docs"),
            "ipfs://QmHash",
            HouseRWA.DocumentStorageType.IPFS,
            "CRE-verified"
        );

        vm.prank(user1);
        houseRWA.createListing(tokenId, HouseRWA.ListingType.FOR_SALE, 1 ether, address(0), false, address(0), 30);

        bytes32 keyHash = keccak256("transfer-key");
        bytes memory encryptedKey = abi.encodePacked("encrypted-key-data");

        vm.prank(creWorkflow);
        houseRWA.completeSale(tokenId, user2, keyHash, encryptedKey);

        assertEq(houseRWA.ownerOf(tokenId), user2);
        bytes memory claimed;
        vm.prank(user2);
        claimed = houseRWA.claimKey(keyHash);
        assertEq(claimed, encryptedKey);
    }

    function test_RevertDoubleClaimKey() public {
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(user2, 1, keccak256("kyc-buyer"), uint48(block.timestamp + 365 days));

        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            user1,
            "house-123",
            keccak256("encrypted-docs"),
            "ipfs://QmHash",
            HouseRWA.DocumentStorageType.IPFS,
            "CRE-verified"
        );

        vm.prank(user1);
        houseRWA.createListing(tokenId, HouseRWA.ListingType.FOR_SALE, 1 ether, address(0), false, address(0), 30);

        bytes32 keyHash = keccak256("transfer-key");
        bytes memory encryptedKey = abi.encodePacked("encrypted-key-data");

        vm.prank(creWorkflow);
        houseRWA.completeSale(tokenId, user2, keyHash, encryptedKey);

        vm.prank(user2);
        houseRWA.claimKey(keyHash);

        vm.prank(user2);
        vm.expectRevert();
        houseRWA.claimKey(keyHash);
    }

    function test_StartAndEndRental() public {
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(user2, 1, keccak256("kyc-renter"), uint48(block.timestamp + 365 days));

        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            user1,
            "house-rent",
            keccak256("encrypted-docs"),
            "ipfs://QmHash",
            HouseRWA.DocumentStorageType.IPFS,
            "CRE-verified"
        );

        vm.prank(user1);
        houseRWA.createListing(tokenId, HouseRWA.ListingType.FOR_RENT, 0.1 ether, address(0), false, address(0), 30);

        vm.prank(user2);
        houseRWA.depositForRental{value: 0.5 ether}(tokenId);

        vm.prank(creWorkflow);
        houseRWA.startRental(tokenId, user2, 30, 0.5 ether, 0.1 ether, abi.encodePacked("access-key"));

        HouseRWA.RentalAgreement memory rental = houseRWA.getActiveRental(tokenId);
        assertTrue(rental.isActive);

        vm.warp(block.timestamp + 31 days);

        vm.prank(user1);
        houseRWA.endRental(tokenId);

        rental = houseRWA.getActiveRental(tokenId);
        assertFalse(rental.isActive);
    }

    function test_CreateAndPayBillWithRecurring() public {
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            user1,
            "house-bill",
            keccak256("encrypted-docs"),
            "ipfs://QmHash",
            HouseRWA.DocumentStorageType.IPFS,
            "CRE-verified"
        );

        vm.prank(owner);
        houseRWA.setTrustedBillProvider(provider, true);

        vm.prank(provider);
        uint256 billIndex =
            houseRWA.createBill(tokenId, "electricity", 15000, uint48(block.timestamp + 30 days), provider, true, 30);

        vm.prank(creWorkflow);
        houseRWA.recordBillPayment(tokenId, billIndex, "stripe", keccak256("payment-ref"));

        HouseRWA.Bill[] memory bills = houseRWA.getBills(tokenId);
        assertEq(bills.length, 2);
        assertEq(uint8(bills[0].status), uint8(HouseRWA.BillStatus.PAID));
        assertEq(uint8(bills[1].status), uint8(HouseRWA.BillStatus.PENDING));
    }

    function test_RevertCreateBillUntrustedProvider() public {
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            user1,
            "house-123",
            keccak256("encrypted-docs"),
            "ipfs://QmHash",
            HouseRWA.DocumentStorageType.IPFS,
            "CRE-verified"
        );

        vm.prank(provider);
        vm.expectRevert();
        houseRWA.createBill(tokenId, "electricity", 15000, uint48(block.timestamp + 30 days), provider, false, 0);
    }

    function test_AuthorizeCREWorkflow() public {
        address newWorkflow = makeAddr("newWorkflow");

        vm.prank(owner);
        houseRWA.authorizeCREWorkflow(newWorkflow);

        assertTrue(houseRWA.authorizedCREWorkflows(newWorkflow));
    }

    function test_SetAndHasKYC() public {
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(user2, 2, keccak256("kyc-data"), uint48(block.timestamp + 365 days));

        assertTrue(houseRWA.hasKYC(user2));

        vm.prank(owner);
        houseRWA.revokeKYC(user2);

        assertFalse(houseRWA.hasKYC(user2));
    }
}
