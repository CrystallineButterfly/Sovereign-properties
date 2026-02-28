// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IReceiver} from "./keystone/IReceiver.sol";
import {IERC165} from "@openzeppelin/contracts/interfaces/IERC165.sol";
import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";
import {Ownable2Step} from "@openzeppelin/contracts/access/Ownable2Step.sol";

/**
 * @title HouseRWAReceiver
 * @notice Keystone/CRE receiver that forwards signed CRE reports to HouseRWA.
 * @dev Security model:
 *      - `onReport` is callable ONLY by a trusted Forwarder (msg.sender).
 *      - Only an allowlisted set of HouseRWA selectors can be executed.
 *      - The receiver must be authorized as a CRE workflow in HouseRWA.
 */
contract HouseRWAReceiver is IReceiver, Ownable2Step {
    /// @notice HouseRWA proxy/implementation address.
    address public immutable houseRWA;

    /// @notice Chainlink Keystone forwarder allowed to call `onReport`.
    address public forwarder;

    /// @notice Allowlist for forwarded HouseRWA function selectors.
    mapping(bytes4 => bool) public allowedSelectors;

    error HouseRWAReceiver_InvalidReport();
    error HouseRWAReceiver_OnlyForwarder(address caller);
    error HouseRWAReceiver_SelectorNotAllowed(bytes4 selector);
    error HouseRWAReceiver_InvalidAddress();

    event ForwarderUpdated(address indexed oldForwarder, address indexed newForwarder);
    event SelectorUpdated(bytes4 indexed selector, bool allowed);
    event ReportForwarded(address indexed houseRWA, bytes4 indexed selector, bytes32 indexed reportHash);

    constructor(address _houseRWA, address _forwarder, address _owner, bytes4[] memory selectors) Ownable(_owner) {
        if (_houseRWA == address(0) || _forwarder == address(0) || _owner == address(0)) {
            revert HouseRWAReceiver_InvalidAddress();
        }
        houseRWA = _houseRWA;
        forwarder = _forwarder;

        for (uint256 i = 0; i < selectors.length; i++) {
            allowedSelectors[selectors[i]] = true;
            emit SelectorUpdated(selectors[i], true);
        }
    }

    /**
     * @notice Update the trusted forwarder.
     */
    function setForwarder(address newForwarder) external onlyOwner {
        if (newForwarder == address(0)) revert HouseRWAReceiver_InvalidAddress();
        address old = forwarder;
        forwarder = newForwarder;
        emit ForwarderUpdated(old, newForwarder);
    }

    /**
     * @notice Allow/deny a selector for forwarding.
     */
    function setAllowedSelector(bytes4 selector, bool allowed) external onlyOwner {
        allowedSelectors[selector] = allowed;
        emit SelectorUpdated(selector, allowed);
    }

    /**
     * @notice Receive a CRE report and forward it to HouseRWA.
     * @dev `report` is expected to be raw calldata targeting HouseRWA.
     */
    function onReport(bytes calldata, bytes calldata report) external override {
        if (msg.sender != forwarder) revert HouseRWAReceiver_OnlyForwarder(msg.sender);
        if (report.length < 4) revert HouseRWAReceiver_InvalidReport();

        bytes4 selector = bytes4(report);
        if (!allowedSelectors[selector]) revert HouseRWAReceiver_SelectorNotAllowed(selector);

        (bool success, bytes memory returndata) = houseRWA.call(report);
        if (!success) {
            // Bubble up revert reason from HouseRWA.
            assembly {
                revert(add(returndata, 0x20), mload(returndata))
            }
        }

        emit ReportForwarded(houseRWA, selector, keccak256(report));
    }

    function supportsInterface(bytes4 interfaceId) public pure override returns (bool) {
        return interfaceId == type(IReceiver).interfaceId || interfaceId == type(IERC165).interfaceId;
    }
}
