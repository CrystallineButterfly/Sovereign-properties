// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @notice Common Chainlink interface for versioned contracts.
/// @dev Implementations may satisfy this by defining a `public constant string typeAndVersion`.
interface ITypeAndVersion {
    function typeAndVersion() external pure returns (string memory);
}

