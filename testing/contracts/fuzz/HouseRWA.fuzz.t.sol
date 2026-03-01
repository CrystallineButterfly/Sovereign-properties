// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {HouseRWA} from "src/HouseRWA.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

/**
 * @title HouseRWAFuzzTest
 * @notice Fuzz testing suite for HouseRWA contract
 * @dev Tests invariant properties and boundary conditions
 */
contract HouseRWAFuzzTest is Test {
    HouseRWA public implementation;
    HouseRWA public houseRWA;
    ERC1967Proxy public proxy;
    
    address public owner;
    address public creWorkflow;
    address public feeRecipient;
    
    function setUp() public {
        owner = makeAddr("owner");
        creWorkflow = makeAddr("creWorkflow");
        feeRecipient = makeAddr("feeRecipient");
        
        vm.startPrank(owner);
        implementation = new HouseRWA();
        
        bytes memory initData = abi.encodeWithSelector(
            HouseRWA.initialize.selector,
            owner,
            feeRecipient,
            creWorkflow
        );
        
        proxy = new ERC1967Proxy(address(implementation), initData);
        houseRWA = HouseRWA(payable(address(proxy)));
        vm.stopPrank();
    }
    
    /**
     * @notice Fuzz test minting with various valid inputs
     */
    function testFuzz_MintHouse(
        address to,
        bytes32 documentHash,
        uint8 storageType
    ) public {
        vm.assume(to != address(0));
        vm.assume(to != address(houseRWA));
        vm.assume(documentHash != bytes32(0));
        vm.assume(storageType <= 3);
        
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            to,
            "fuzz-test",
            documentHash,
            "ipfs://fuzz",
            HouseRWA.DocumentStorageType(storageType),
            "verification"
        );
        
        assertEq(houseRWA.ownerOf(tokenId), to);
    }
    
    /**
     * @notice Fuzz test listing creation with various prices
     */
    function testFuzz_CreateListing(
        address seller,
        uint96 price,
        uint8 listingType,
        uint16 durationDays
    ) public {
        vm.assume(seller != address(0));
        vm.assume(seller != address(houseRWA));
        vm.assume(price > 0);
        vm.assume(listingType > 0 && listingType <= 3);
        vm.assume(durationDays > 0 && durationDays <= 365);

        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(seller, 2, keccak256("seller-kyc"), uint48(block.timestamp + 365 days));
        
        // Mint token first
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            seller,
            "fuzz-listing",
            keccak256("docs"),
            "ipfs://test",
            HouseRWA.DocumentStorageType.IPFS,
            "verified"
        );
        
        vm.prank(seller);
        houseRWA.createListing(
            tokenId,
            HouseRWA.ListingType(listingType),
            price,
            address(0),
            false,
            address(0),
            uint48(durationDays)
        );
        
        HouseRWA.Listing memory listing = houseRWA.getListing(tokenId);
        assertEq(listing.price, price);
        assertEq(uint256(listing.listingType), uint256(listingType));
    }
    
    /**
     * @notice Fuzz test private sale flow
     */
    function testFuzz_PrivateSale(
        address seller,
        address buyer,
        uint96 price,
        bytes32 keyHash
    ) public {
        vm.assume(seller != address(0) && buyer != address(0));
        vm.assume(seller != buyer);
        vm.assume(seller != address(houseRWA));
        vm.assume(price > 0);
        vm.assume(keyHash != bytes32(0));
        
        // Setup KYC for buyer
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(buyer, 1, keccak256("kyc"), uint48(block.timestamp + 365 days));

        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(seller, 2, keccak256("seller-kyc"), uint48(block.timestamp + 365 days));
        
        // Mint and list
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            seller,
            "fuzz-sale",
            keccak256("docs"),
            "ipfs://test",
            HouseRWA.DocumentStorageType.IPFS,
            "verified"
        );
        
        vm.prank(seller);
        houseRWA.createListing(
            tokenId,
            HouseRWA.ListingType.FOR_SALE,
            price,
            address(0),
            true,
            buyer,
            30
        );
        
        // Complete sale
        vm.prank(creWorkflow);
        houseRWA.completeSale(tokenId, buyer, keyHash, abi.encodePacked("encrypted-key"));
        
        assertEq(houseRWA.ownerOf(tokenId), buyer);
    }
    
    /**
     * @notice Fuzz test bill creation
     */
    function testFuzz_CreateBill(
        string calldata billType,
        uint96 amount,
        uint48 dueDateOffset
    ) public {
        vm.assume(bytes(billType).length > 0 && bytes(billType).length <= 50);
        vm.assume(amount > 0);
        vm.assume(dueDateOffset > 1 days && dueDateOffset <= 365 days);
        
        vm.prank(creWorkflow);
        uint256 tokenId = houseRWA.mint(
            makeAddr("owner"),
            "fuzz-bill",
            keccak256("docs"),
            "ipfs://test",
            HouseRWA.DocumentStorageType.IPFS,
            "verified"
        );
        
        vm.prank(creWorkflow);
        uint256 billIndex = houseRWA.createBill(
            tokenId,
            billType,
            amount,
            uint48(block.timestamp + dueDateOffset),
            makeAddr("provider"),
            false,
            0
        );
        
        HouseRWA.Bill[] memory bills = houseRWA.getBills(tokenId);
        assertEq(bills.length, billIndex + 1);
    }
    
    /**
     * @notice Fuzz test KYC verification
     */
    function testFuzz_KYCVerification(
        address user,
        uint8 level,
        uint48 expiryOffset
    ) public {
        vm.assume(user != address(0));
        vm.assume(level <= 2);
        vm.assume(expiryOffset > 1 days && expiryOffset <= 365 days);
        
        vm.prank(creWorkflow);
        houseRWA.setKYCVerification(
            user,
            level,
            keccak256("kyc-data"),
            uint48(block.timestamp + expiryOffset)
        );
        
        assertTrue(houseRWA.hasKYC(user));
    }
    
    /**
     * @notice Fuzz test validator staking
     */
    function testFuzz_ValidatorStake(
        address validator,
        uint256 stakeAmount
    ) public {
        vm.assume(validator != address(0));
        vm.assume(stakeAmount >= 1 ether);
        vm.assume(stakeAmount <= 1000 ether);
        
        vm.deal(validator, stakeAmount);
        
        vm.prank(validator);
        houseRWA.stakeAsValidator{value: stakeAmount}();
        
        (uint96 stakedAmount,,,,,) = houseRWA.validators(validator);
        assertEq(stakedAmount, stakeAmount);
    }
    
    /**
     * @notice Invariant: Total supply should always equal number of minted tokens
     */
    function testFuzz_Invariant_TotalSupply(
        uint8 numMints
    ) public {
        vm.assume(numMints > 0 && numMints <= 10);
        
        for (uint8 i = 0; i < numMints; i++) {
            address recipient = makeAddr(string(abi.encodePacked("user", i)));
            
            vm.prank(creWorkflow);
            houseRWA.mint(
                recipient,
                string(abi.encodePacked("house-", i)),
                keccak256(abi.encodePacked("docs-", i)),
                "ipfs://test",
                HouseRWA.DocumentStorageType.IPFS,
                "verified"
            );
        }
        
        assertEq(houseRWA.totalSupply(), numMints);
    }
    
    /**
     * @notice Invariant: Token ownership should be consistent
     */
    function testFuzz_Invariant_Ownership(
        address owner1,
        address owner2
    ) public {
        vm.assume(owner1 != address(0) && owner2 != address(0));
        vm.assume(owner1 != owner2);
        
        vm.prank(creWorkflow);
        uint256 tokenId1 = houseRWA.mint(owner1, "house1", keccak256("docs1"), "ipfs://1", HouseRWA.DocumentStorageType.IPFS, "verified");
        
        vm.prank(creWorkflow);
        uint256 tokenId2 = houseRWA.mint(owner2, "house2", keccak256("docs2"), "ipfs://2", HouseRWA.DocumentStorageType.IPFS, "verified");
        
        assertEq(houseRWA.ownerOf(tokenId1), owner1);
        assertEq(houseRWA.ownerOf(tokenId2), owner2);
        assertEq(houseRWA.balanceOf(owner1), 1);
        assertEq(houseRWA.balanceOf(owner2), 1);
    }
    
    function makeAddr(string memory name) internal returns (address addr) {
        addr = vm.addr(uint256(keccak256(bytes(name))));
        vm.label(addr, name);
    }
}
