import { describe, it, expect, beforeAll } from '@jest/globals';
import { ethers } from 'ethers';

/**
 * Complete RWA Platform Integration Tests
 * Tests the full lifecycle: mint → list → sale → rental → payment
 */

const PROVIDER_URL = process.env.RPC_URL || 'http://localhost:8545';
const CONTRACT_ADDRESS = process.env.CONTRACT_ADDRESS || '';

interface TestContext {
  provider: ethers.JsonRpcProvider;
  seller: ethers.Wallet;
  buyer: ethers.Wallet;
  tenant: ethers.Wallet;
  contract: ethers.Contract;
}

describe('Complete RWA Platform Flow', () => {
  let ctx: TestContext;

  beforeAll(async () => {
    const provider = new ethers.JsonRpcProvider(PROVIDER_URL);
    
    // Create test wallets
    const seller = new ethers.Wallet(ethers.Wallet.createRandom().privateKey, provider);
    const buyer = new ethers.Wallet(ethers.Wallet.createRandom().privateKey, provider);
    const tenant = new ethers.Wallet(ethers.Wallet.createRandom().privateKey, provider);
    
    // Fund wallets
    const funded = await Promise.all([
      fundWallet(provider, seller.address, '10'),
      fundWallet(provider, buyer.address, '10'),
      fundWallet(provider, tenant.address, '5'),
    ]);
    
    // Load contract
    const contract = new ethers.Contract(
      CONTRACT_ADDRESS,
      HouseRWAABI,
      seller
    );
    
    ctx = { provider, seller, buyer, tenant, contract };
  });

  it('executes full lifecycle: mint → list → sale → rental → payment', async () => {
    // Step 1: Mint house as seller
    console.log('🏠 Step 1: Minting house...');
    const mintTx = await mintHouse(ctx);
    const tokenId = mintTx.tokenId;
    expect(tokenId).toBeDefined();
    
    // Verify minting
    const owner = await ctx.contract.ownerOf(tokenId);
    expect(owner).toBe(ctx.seller.address);
    
    // Step 2: List for private sale
    console.log('📋 Step 2: Creating listing...');
    const salePrice = ethers.parseEther('1');
    await listForSale(ctx, tokenId, salePrice);
    
    const listing = await ctx.contract.getListing(tokenId);
    expect(listing.price).toBe(salePrice);
    expect(listing.listingType).toBe(1); // FOR_SALE
    
    // Step 3: Buyer purchases
    console.log('💰 Step 3: Completing sale...');
    await purchaseHouse(ctx, tokenId, ctx.buyer);
    
    // Step 4: Verify ownership transfer
    console.log('✅ Step 4: Verifying ownership...');
    const newOwner = await ctx.contract.ownerOf(tokenId);
    expect(newOwner).toBe(ctx.buyer.address);
    
    // Step 5: Create rental
    console.log('🏘️  Step 5: Creating rental...');
    const rental = await createRental(ctx, tokenId, ctx.tenant.address);
    expect(rental.agreementId).toBeDefined();
    
    // Verify rental is active
    const rentalInfo = await ctx.contract.getActiveRental(tokenId);
    expect(rentalInfo.isActive).toBe(true);
    expect(rentalInfo.renter).toBe(ctx.tenant.address);
    
    // Step 6: Process rental payment
    console.log('💳 Step 6: Processing rental payment...');
    await processRentalPayment(ctx, rental.agreementId);
    
    // Step 7: Record bill payment
    console.log('🧾 Step 7: Recording bill payment...');
    await payBill(ctx, tokenId, 'electricity', 15000); // $150
    
    // Step 8: Verify all records
    console.log('📊 Step 8: Verifying payment history...');
    const paymentHistory = await getPaymentHistory(ctx, tokenId);
    expect(paymentHistory.length).toBeGreaterThanOrEqual(2);
    
    console.log('✨ Full lifecycle test completed successfully!');
  });

  it('handles concurrent operations correctly', async () => {
    // Mint multiple houses
    const tokenIds = await Promise.all([
      mintHouse(ctx).then(tx => tx.tokenId),
      mintHouse(ctx).then(tx => tx.tokenId),
      mintHouse(ctx).then(tx => tx.tokenId),
    ]);
    
    // List all simultaneously
    const listingPromises = tokenIds.map(id => 
      listForSale(ctx, id, ethers.parseEther('0.5'))
    );
    
    await Promise.all(listingPromises);
    
    // Verify all listings created
    for (const tokenId of tokenIds) {
      const listing = await ctx.contract.getListing(tokenId);
      expect(listing.price).toBe(ethers.parseEther('0.5'));
    }
  });

  it('maintains data consistency across operations', async () => {
    const { tokenId } = await mintHouse(ctx);
    
    // Record multiple bills
    const bills = [
      { type: 'electricity', amount: 15000 },
      { type: 'water', amount: 8000 },
      { type: 'internet', amount: 5000 },
    ];
    
    for (const bill of bills) {
      await payBill(ctx, tokenId, bill.type, bill.amount);
    }
    
    // Verify all bills recorded
    const allBills = await ctx.contract.getBills(tokenId);
    expect(allBills.length).toBe(bills.length);
    
    // Verify bill types
    const billTypes = allBills.map((b: any) => b.billType);
    for (const bill of bills) {
      expect(billTypes).toContain(bill.type);
    }
  });

  it('recovers from partial failures gracefully', async () => {
    const { tokenId } = await mintHouse(ctx);
    
    // Create listing
    await listForSale(ctx, tokenId, ethers.parseEther('1'));
    
    // Cancel listing
    await cancelListing(ctx, tokenId);
    
    // Verify listing cleared
    const listing = await ctx.contract.getListing(tokenId);
    expect(listing.listingType).toBe(0); // NONE
    
    // Should be able to list again
    await listForSale(ctx, tokenId, ethers.parseEther('1.5'));
    
    const newListing = await ctx.contract.getListing(tokenId);
    expect(newListing.price).toBe(ethers.parseEther('1.5'));
  });
});

// Helper functions
async function fundWallet(provider: ethers.JsonRpcProvider, address: string, amount: string) {
  // Implementation depends on test environment
  // For Hardhat/Anvil, use the default funded account
}

async function mintHouse(ctx: TestContext): Promise<{ tokenId: string }> {
  // Mock implementation - would call actual contract
  return { tokenId: `token-${Date.now()}` };
}

async function listForSale(ctx: TestContext, tokenId: string, price: bigint) {
  // Mock implementation
}

async function purchaseHouse(ctx: TestContext, tokenId: string, buyer: ethers.Wallet) {
  // Mock implementation
}

async function createRental(ctx: TestContext, tokenId: string, tenantAddress: string): Promise<{ agreementId: string }> {
  // Mock implementation
  return { agreementId: `rental-${Date.now()}` };
}

async function processRentalPayment(ctx: TestContext, rentalId: string) {
  // Mock implementation
}

async function payBill(ctx: TestContext, tokenId: string, billType: string, amount: number) {
  // Mock implementation
}

async function getPaymentHistory(ctx: TestContext, tokenId: string): Promise<any[]> {
  // Mock implementation
  return [];
}

async function cancelListing(ctx: TestContext, tokenId: string) {
  // Mock implementation
}

// ABI for HouseRWA contract
const HouseRWAABI = [
  "function mint(address to, string calldata houseId, bytes32 documentHash, string calldata documentURI, uint8 storageType, string calldata verificationData) external returns (uint256)",
  "function createListing(uint256 tokenId, uint8 listingType, uint96 price, address preferredToken, bool isPrivateSale, address allowedBuyer, uint48 durationDays) external",
  "function completeSale(uint256 tokenId, address buyer, bytes32 keyHash, bytes calldata encryptedKey) external",
  "function startRental(uint256 tokenId, address renter, uint48 durationDays, uint96 depositAmount, uint96 monthlyRent, bytes calldata encryptedAccessKey) external",
  "function recordBillPayment(uint256 tokenId, uint256 billIndex, string calldata paymentMethod, bytes32 paymentReference) external",
  "function ownerOf(uint256 tokenId) external view returns (address)",
  "function getListing(uint256 tokenId) external view returns (tuple(uint8 listingType, uint96 price, address preferredToken, bool isPrivateSale, address allowedBuyer, uint48 createdAt, uint48 expiresAt, uint8 platformFee))",
  "function getActiveRental(uint256 tokenId) external view returns (tuple(address renter, uint48 startTime, uint48 endTime, uint96 depositAmount, uint96 monthlyRent, bool isActive, bytes32 encryptedAccessKeyHash, uint8 disputeStatus))",
  "function getBills(uint256 tokenId) external view returns (tuple(string billType, uint96 amount, uint48 dueDate, uint48 paidAt, uint8 status, bytes32 paymentReference, bool isRecurring, address provider, uint8 recurrenceInterval)[])",
];
