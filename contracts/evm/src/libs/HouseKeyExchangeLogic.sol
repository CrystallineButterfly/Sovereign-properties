// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {SafeCast} from "@openzeppelin/contracts/utils/math/SafeCast.sol";

library HouseKeyExchangeLogic {
    using SafeCast for uint256;

    function requireNonZeroHash(bytes32 keyHash) internal pure {
        require(keyHash != bytes32(0), "HouseRWA: Invalid key hash");
    }

    function requireNonEmptyKey(bytes memory encryptedKey, string memory errorMessage) internal pure {
        require(encryptedKey.length > 0, errorMessage);
    }

    function requireUnusedHash(uint48 createdAt) internal pure {
        require(createdAt == 0, "HouseRWA: Key hash already used");
    }

    function requireClaimable(
        address intendedRecipient,
        bool isClaimed,
        uint48 createdAt,
        uint48 expiresAt,
        address claimant,
        uint256 currentTimestamp
    ) internal pure {
        require(intendedRecipient == claimant, "HouseRWA: Not recipient");
        require(!isClaimed, "HouseRWA: Already claimed");
        require(currentTimestamp <= expiresAt, "HouseRWA: Key expired");
        require(currentTimestamp >= createdAt, "HouseRWA: Invalid timestamp");
    }

    function hashKey(bytes memory encryptedKey) internal pure returns (bytes32) {
        return keccak256(encryptedKey);
    }

    function keyExpiry(uint256 currentTimestamp, uint256 ttlSeconds) internal pure returns (uint48) {
        return (currentTimestamp + ttlSeconds).toUint48();
    }
}
