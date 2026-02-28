// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {SafeCast} from "@openzeppelin/contracts/utils/math/SafeCast.sol";

library HouseListingLogic {
    using SafeCast for uint256;

    uint8 internal constant LISTING_FOR_SALE = 1;
    uint8 internal constant LISTING_FOR_RENT = 2;

    function requireSalesActive(bool salesPaused) internal pure {
        require(!salesPaused, "HouseRWA: Sales paused");
    }

    function requireRentalsActive(bool rentalsPaused) internal pure {
        require(!rentalsPaused, "HouseRWA: Rentals paused");
    }

    function requirePaymentsActive(bool paymentsPaused) internal pure {
        require(!paymentsPaused, "HouseRWA: Payments paused");
    }

    function requireListingForSale(uint8 listingType, uint48 expiresAt, uint256 currentTimestamp) internal pure {
        require(listingType == LISTING_FOR_SALE, "HouseRWA: Not for sale");
        require(expiresAt > currentTimestamp, "HouseRWA: Listing expired");
    }

    function requireListingForRent(uint8 listingType, uint48 expiresAt, uint256 currentTimestamp) internal pure {
        require(listingType == LISTING_FOR_RENT, "HouseRWA: Not for rent");
        require(expiresAt > currentTimestamp, "HouseRWA: Listing expired");
    }

    function requirePrivateCounterparty(
        bool isPrivateSale,
        address allowedBuyer,
        address candidate,
        string memory errorMessage
    ) internal pure {
        if (isPrivateSale) {
            require(candidate == allowedBuyer, errorMessage);
        }
    }

    function requireNotRented(bool isRented) internal pure {
        require(!isRented, "HouseRWA: Already rented");
    }

    function requireValidListingDuration(uint48 durationDays, uint256 maxRentalDurationDays) internal pure {
        require(durationDays == 0 || durationDays <= maxRentalDurationDays, "HouseRWA: Invalid duration");
    }

    function requireValidRentalDuration(uint48 durationDays, uint256 maxRentalDurationDays) internal pure {
        require(durationDays > 0 && durationDays <= maxRentalDurationDays, "HouseRWA: Invalid duration");
    }

    function computeListingExpiry(uint48 durationDays, uint256 currentTimestamp) internal pure returns (uint48) {
        if (durationDays > 0) {
            return (currentTimestamp + (uint256(durationDays) * 1 days)).toUint48();
        }
        return (currentTimestamp + 30 days).toUint48();
    }

    function computeRentalEnd(uint48 durationDays, uint256 currentTimestamp) internal pure returns (uint48) {
        return (currentTimestamp + (uint256(durationDays) * 1 days)).toUint48();
    }
}
