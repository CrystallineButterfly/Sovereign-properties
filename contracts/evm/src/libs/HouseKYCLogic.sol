// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

library HouseKYCLogic {
    function hasValidKyc(
        bool isVerified,
        uint8 verificationLevel,
        uint48 expiryDate,
        uint8 minLevel,
        uint256 currentTimestamp
    ) internal pure returns (bool) {
        return isVerified && verificationLevel >= minLevel && expiryDate > currentTimestamp;
    }

    function requireValidKyc(
        bool isVerified,
        uint8 verificationLevel,
        uint48 expiryDate,
        uint8 minLevel,
        uint256 currentTimestamp,
        string memory errorMessage
    ) internal pure {
        require(hasValidKyc(isVerified, verificationLevel, expiryDate, minLevel, currentTimestamp), errorMessage);
    }
}
