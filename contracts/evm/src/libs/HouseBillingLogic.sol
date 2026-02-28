// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {SafeCast} from "@openzeppelin/contracts/utils/math/SafeCast.sol";

library HouseBillingLogic {
    using SafeCast for uint256;

    uint8 internal constant BILL_PENDING = 0;
    uint8 internal constant BILL_OVERDUE = 2;

    function requireAuthorizedCreator(bool isCre, bool isHouseOwner, bool isTrustedProvider) internal pure {
        require(isCre || isHouseOwner || isTrustedProvider, "HouseRWA: Not authorized");
    }

    function validateCreateBillInput(
        string calldata billType,
        uint96 amount,
        uint48 dueDate,
        bool isRecurring,
        uint8 recurrenceInterval,
        uint256 currentTimestamp,
        uint256 maxBillDueWindow,
        uint8 maxRecurrenceIntervalDays
    ) internal pure {
        require(bytes(billType).length > 0 && bytes(billType).length <= 50, "HouseRWA: Invalid bill type");
        require(amount > 0, "HouseRWA: Invalid amount");
        require(dueDate > currentTimestamp, "HouseRWA: Due date in past");
        require(dueDate <= currentTimestamp + maxBillDueWindow, "HouseRWA: Due date too far");

        if (isRecurring) {
            require(recurrenceInterval > 0, "HouseRWA: Invalid recurrence interval");
            require(recurrenceInterval <= maxRecurrenceIntervalDays, "HouseRWA: Recurrence too long");
            return;
        }

        require(recurrenceInterval == 0, "HouseRWA: Non-recurring interval must be zero");
    }

    function requireProvider(address provider) internal pure {
        require(provider != address(0), "HouseRWA: Invalid provider");
    }

    function requireBillCapacity(uint256 currentCount, uint256 maxCount) internal pure {
        require(currentCount < maxCount, "HouseRWA: Max bills reached");
    }

    function requireBillIndexForPayment(uint256 billIndex, uint256 totalCount) internal pure {
        require(billIndex < totalCount, "HouseRWA: Invalid bill index");
    }

    function requireBillIndexForDispute(uint256 billIndex, uint256 totalCount) internal pure {
        require(billIndex < totalCount, "HouseRWA: Invalid bill");
    }

    function requireBillIsPayable(uint8 status) internal pure {
        require(status == BILL_PENDING || status == BILL_OVERDUE, "HouseRWA: Already paid");
    }

    function nextDueDate(uint48 dueDate, uint8 recurrenceInterval) internal pure returns (uint48) {
        return (uint256(dueDate) + (uint256(recurrenceInterval) * 1 days)).toUint48();
    }
}
