// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console} from "forge-std/Script.sol";
import {HouseRWA} from "../src/HouseRWA.sol";
import {HouseRWAReceiver} from "../src/HouseRWAReceiver.sol";
import {ERC1967Proxy} from "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";

/**
 * @title Deploy
 * @notice Deployment script for HouseRWA contract with UUPS proxy
 * @dev Uses Foundry's scripting capabilities for deployment
 *
 * Usage:
 *   forge script script/Deploy.s.sol:Deploy --rpc-url $RPC_URL --broadcast
 *
 * Environment Variables:
 *   - PRIVATE_KEY: Deployer private key
 *   - CRE_WORKFLOW_ADDRESS: Address of CRE workflow (defaults to deployer)
 *   - CRE_FORWARDER_ADDRESS: Keystone forwarder allowed to call the receiver (required for report-based writes)
 *   - FEE_RECIPIENT: Address to receive protocol fees (defaults to deployer)
 *   - CRE_STRICT_MODE: If true, revokes CRE_WORKFLOW_ADDRESS after authorizing the receiver
 */
contract Deploy is Script {
    function run() external {
        // Get deployment parameters from environment
        uint256 deployerPrivateKey = uint256(vm.envBytes32("PRIVATE_KEY"));
        address deployer = vm.addr(deployerPrivateKey);

        address creWorkflow = vm.envOr("CRE_WORKFLOW_ADDRESS", deployer);
        address feeRecipient = vm.envOr("FEE_RECIPIENT", deployer);
        address creForwarder = vm.envOr("CRE_FORWARDER_ADDRESS", address(0));
        bool strictMode = vm.envOr("CRE_STRICT_MODE", false);

        if (creForwarder == address(0)) {
            console.log(
                "WARNING: CRE_FORWARDER_ADDRESS not set; defaulting receiver forwarder to deployer (local-only)"
            );
            creForwarder = deployer;
        }

        console.log("=== HouseRWA Deployment ===");
        console.log("Deployer:", deployer);
        console.log("CRE Workflow:", creWorkflow);
        console.log("CRE Forwarder:", creForwarder);
        console.log("Fee Recipient:", feeRecipient);
        console.log("Strict Mode:", strictMode);

        vm.startBroadcast(deployerPrivateKey);

        // Step 1: Deploy implementation contract
        console.log("\n1. Deploying HouseRWA implementation...");
        HouseRWA implementation = new HouseRWA();
        console.log("   Implementation deployed at:", address(implementation));

        // Step 2: Prepare initialization data
        console.log("\n2. Preparing initialization data...");
        bytes memory initData = abi.encodeWithSelector(
            HouseRWA.initialize.selector,
            deployer, // owner
            feeRecipient, // fee recipient
            creWorkflow // initial CRE workflow
        );

        // Step 3: Deploy proxy contract
        console.log("\n3. Deploying ERC1967Proxy...");
        ERC1967Proxy proxy = new ERC1967Proxy(address(implementation), initData);
        console.log("   Proxy deployed at:", address(proxy));

        HouseRWA houseRWA = HouseRWA(payable(address(proxy)));

        // Step 4: Deploy secured receiver contract for CRE report writes
        console.log("\n4. Deploying HouseRWAReceiver...");
        bytes4[] memory selectors = new bytes4[](7);
        selectors[0] = HouseRWA.setKYCVerification.selector;
        selectors[1] = HouseRWA.mint.selector;
        selectors[2] = HouseRWA.completeSale.selector;
        selectors[3] = HouseRWA.startRental.selector;
        selectors[4] = HouseRWA.createBill.selector;
        selectors[5] = HouseRWA.recordBillPayment.selector;
        selectors[6] = HouseRWA.createListingFromWorkflow.selector;

        HouseRWAReceiver receiver = new HouseRWAReceiver(address(proxy), creForwarder, deployer, selectors);
        console.log("   Receiver deployed at:", address(receiver));

        // Authorize receiver as a CRE workflow so it can call onlyCRE functions
        houseRWA.authorizeCREWorkflow(address(receiver));
        console.log("   Receiver authorized in HouseRWA");

        if (strictMode) {
            // Optionally revoke the initial EOA workflow to force report-only writes.
            if (creWorkflow != address(receiver) && creWorkflow != address(0)) {
                houseRWA.revokeCREWorkflow(creWorkflow);
                console.log("   Revoked initial CRE workflow:", creWorkflow);
            }
        }

        vm.stopBroadcast();

        // Step 5: Verification
        console.log("\n5. Verification...");
        require(houseRWA.owner() == deployer, "Owner mismatch");
        require(houseRWA.feeRecipient() == feeRecipient, "Fee recipient mismatch");
        require(houseRWA.authorizedCREWorkflows(address(receiver)), "Receiver not authorized");
        if (!strictMode) {
            require(houseRWA.authorizedCREWorkflows(creWorkflow), "CRE workflow not authorized");
        }
        console.log("   All verifications passed!");

        // Step 6: Save deployment info
        console.log("\n6. Saving deployment info...");
        _saveDeployment(
            address(implementation),
            address(proxy),
            address(receiver),
            creForwarder,
            deployer,
            creWorkflow,
            feeRecipient,
            strictMode
        );

        console.log("\n=== Deployment Complete ===");
        console.log("HouseRWA Proxy Address:", address(proxy));
        console.log("Implementation Address:", address(implementation));
        console.log("Receiver Address:", address(receiver));
        console.log("\nIMPORTANT: Use the PROXY address for all interactions!");
    }

    function _saveDeployment(
        address implementation,
        address proxy,
        address receiver,
        address creForwarder,
        address deployer,
        address creWorkflow,
        address feeRecipient,
        bool strictMode
    ) internal {
        vm.createDir("deployments", true);

        string memory deploymentInfo = string.concat(
            "{\n",
            '  "contractName": "HouseRWA",\n',
            '  "version": "1.0.0",\n',
            '  "implementationAddress": "',
            vm.toString(implementation),
            '",\n',
            '  "proxyAddress": "',
            vm.toString(proxy),
            '",\n',
            '  "receiverAddress": "',
            vm.toString(receiver),
            '",\n',
            '  "creForwarderAddress": "',
            vm.toString(creForwarder),
            '",\n',
            '  "deployer": "',
            vm.toString(deployer),
            '",\n',
            '  "creWorkflowAddress": "',
            vm.toString(creWorkflow),
            '",\n',
            '  "feeRecipient": "',
            vm.toString(feeRecipient),
            '",\n',
            '  "strictMode": ',
            strictMode ? "true" : "false",
            ",\n",
            '  "timestamp": ',
            vm.toString(block.timestamp),
            ",\n",
            '  "chainId": ',
            vm.toString(block.chainid),
            ",\n",
            '  "network": "',
            _getNetworkName(),
            '"\n',
            "}"
        );

        string memory filename = string.concat(
            "deployments/houserwa_", vm.toString(block.chainid), "_", vm.toString(block.timestamp), ".json"
        );

        vm.writeFile(filename, deploymentInfo);
        console.log("   Deployment info saved to:", filename);
    }

    function _getNetworkName() internal view returns (string memory) {
        uint256 chainId = block.chainid;
        if (chainId == 1) return "mainnet";
        if (chainId == 11155111) return "sepolia";
        if (chainId == 137) return "polygon";
        if (chainId == 80001) return "mumbai";
        if (chainId == 42161) return "arbitrum";
        if (chainId == 10) return "optimism";
        if (chainId == 31337) return "anvil";
        return "unknown";
    }
}

/**
 * @title DeployAndSetup
 * @notice Deployment script with full initial setup
 * @dev Deploys contract and performs initial configuration
 */
contract DeployAndSetup is Script {
    function run() external {
        uint256 deployerPrivateKey = uint256(vm.envBytes32("PRIVATE_KEY"));
        address deployer = vm.addr(deployerPrivateKey);

        address creWorkflow = vm.envOr("CRE_WORKFLOW_ADDRESS", deployer);
        address feeRecipient = vm.envOr("FEE_RECIPIENT", deployer);
        address creForwarder = vm.envOr("CRE_FORWARDER_ADDRESS", address(0));
        bool strictMode = vm.envOr("CRE_STRICT_MODE", false);

        if (creForwarder == address(0)) {
            console.log(
                "WARNING: CRE_FORWARDER_ADDRESS not set; defaulting receiver forwarder to deployer (local-only)"
            );
            creForwarder = deployer;
        }

        console.log("=== HouseRWA Deployment with Setup ===");

        vm.startBroadcast(deployerPrivateKey);

        // Deploy implementation
        HouseRWA implementation = new HouseRWA();
        console.log("Implementation:", address(implementation));

        // Deploy proxy
        bytes memory initData =
            abi.encodeWithSelector(HouseRWA.initialize.selector, deployer, feeRecipient, creWorkflow);

        ERC1967Proxy proxy = new ERC1967Proxy(address(implementation), initData);
        HouseRWA houseRWA = HouseRWA(payable(address(proxy)));
        console.log("Proxy:", address(proxy));

        // Deploy receiver + authorize for report-based writes
        bytes4[] memory selectors = new bytes4[](7);
        selectors[0] = HouseRWA.setKYCVerification.selector;
        selectors[1] = HouseRWA.mint.selector;
        selectors[2] = HouseRWA.completeSale.selector;
        selectors[3] = HouseRWA.startRental.selector;
        selectors[4] = HouseRWA.createBill.selector;
        selectors[5] = HouseRWA.recordBillPayment.selector;
        selectors[6] = HouseRWA.createListingFromWorkflow.selector;

        HouseRWAReceiver receiver = new HouseRWAReceiver(address(proxy), creForwarder, deployer, selectors);
        houseRWA.authorizeCREWorkflow(address(receiver));
        console.log("Receiver:", address(receiver));

        if (strictMode) {
            if (creWorkflow != address(receiver) && creWorkflow != address(0)) {
                houseRWA.revokeCREWorkflow(creWorkflow);
                console.log("Revoked initial CRE workflow:", creWorkflow);
            }
        }

        // Setup configuration
        console.log("\nPerforming initial setup...");

        // Set KYC requirements
        houseRWA.setKYCRequirements(1, 2, 100000 * 100); // $100k threshold
        console.log("- KYC requirements set");

        // Note: Set trusted bill providers and arbitrators manually after deployment
        // or add them as parameters to this script

        vm.stopBroadcast();

        console.log("\n=== Deployment Complete ===");
        console.log("Proxy Address (use this!):", address(proxy));
    }
}

/**
 * @title UpgradeHouseRWA
 * @notice Upgrade script for HouseRWA proxy
 * @dev Deploys new implementation and upgrades proxy
 */
contract UpgradeHouseRWA is Script {
    function run() external {
        uint256 deployerPrivateKey = uint256(vm.envBytes32("PRIVATE_KEY"));
        address proxyAddress = address(uint160(uint256(vm.envBytes32("PROXY_ADDRESS"))));

        console.log("=== HouseRWA Upgrade ===");
        console.log("Proxy:", proxyAddress);

        vm.startBroadcast(deployerPrivateKey);

        // Deploy new implementation
        console.log("Deploying new implementation...");
        HouseRWA newImplementation = new HouseRWA();
        console.log("New implementation:", address(newImplementation));

        // Upgrade proxy
        console.log("Upgrading proxy...");
        HouseRWA proxy = HouseRWA(payable(proxyAddress));
        proxy.upgradeToAndCall(address(newImplementation), "");

        vm.stopBroadcast();

        console.log("\n=== Upgrade Complete ===");
        console.log("New implementation:", address(newImplementation));
    }
}
