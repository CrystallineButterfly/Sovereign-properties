import { test, expect } from '@playwright/test';

// Test configuration
const BASE_URL = process.env.TEST_BASE_URL || 'http://localhost:5173';
const TEST_TIMEOUT = 60000;

test.describe.configure({ mode: 'serial' });

/**
 * House Minting E2E Tests
 */
test.describe('House Minting Flow', () => {
  test.beforeEach(async ({ page }) => {
    // Login with Privy (mock for testing)
    await page.goto(`${BASE_URL}/login`);
    await page.click('[data-testid="privy-login"]');
    // Wait for authentication
    await page.waitForSelector('[data-testid="user-authenticated"]', { timeout: TEST_TIMEOUT });
  });

  test('completes full minting flow', async ({ page }) => {
    await page.goto(`${BASE_URL}/mint`);
    
    // Fill property form
    await page.fill('[data-testid="property-address"]', '123 Test Street, Test City, TC 12345');
    await page.fill('[data-testid="property-price"]', '1000000');
    await page.fill('[data-testid="property-bedrooms"]', '3');
    await page.fill('[data-testid="property-bathrooms"]', '2');
    await page.fill('[data-testid="property-sqft"]', '2000');
    
    // Upload documents
    const deedInput = await page.locator('[data-testid="document-upload-deed"]');
    await deedInput.setInputFiles({
      name: 'test-deed.pdf',
      mimeType: 'application/pdf',
      buffer: Buffer.from('test deed content'),
    });
    
    // Verify upload success
    await expect(page.locator('[data-testid="document-deed-status"]')).toContainText('Uploaded');
    
    // Submit mint
    await page.click('[data-testid="mint-button"]');
    
    // Wait for transaction confirmation
    await expect(page.locator('[data-testid="success-message"]')).toBeVisible({ timeout: TEST_TIMEOUT });
    await expect(page.locator('[data-testid="token-id"]')).not.toBeEmpty();
    
    // Verify token was created
    const tokenId = await page.locator('[data-testid="token-id"]').textContent();
    expect(tokenId).toMatch(/^token-/);
  });

  test('rejects minting without required documents', async ({ page }) => {
    await page.goto(`${BASE_URL}/mint`);
    
    // Fill form without documents
    await page.fill('[data-testid="property-address"]', '456 No Docs St');
    await page.fill('[data-testid="property-price"]', '500000');
    
    // Try to submit
    await page.click('[data-testid="mint-button"]');
    
    // Should show error
    await expect(page.locator('[data-testid="error-message"]')).toContainText('At least one document is required');
  });

  test('validates KYC before minting high value property', async ({ page }) => {
    await page.goto(`${BASE_URL}/mint`);
    
    // High value property
    await page.fill('[data-testid="property-address"]', '789 Luxury Lane');
    await page.fill('[data-testid="property-price"]', '5000000'); // $5M
    
    // Upload required documents
    const deedInput = await page.locator('[data-testid="document-upload-deed"]');
    await deedInput.setInputFiles({
      name: 'deed.pdf',
      mimeType: 'application/pdf',
      buffer: Buffer.from('deed content'),
    });
    
    await page.click('[data-testid="mint-button"]');
    
    // Should require level 2 KYC
    await expect(page.locator('[data-testid="kyc-required-modal"]')).toBeVisible();
    await expect(page.locator('[data-testid="kyc-level-required"]')).toContainText('Level 2');
  });

  test('handles document encryption progress', async ({ page }) => {
    await page.goto(`${BASE_URL}/mint`);
    
    await page.fill('[data-testid="property-address"]', '321 Encryption Ave');
    await page.fill('[data-testid="property-price"]', '750000');
    
    // Upload large document
    const largeDoc = Buffer.alloc(1024 * 1024); // 1MB
    const deedInput = await page.locator('[data-testid="document-upload-deed"]');
    await deedInput.setInputFiles({
      name: 'large-deed.pdf',
      mimeType: 'application/pdf',
      buffer: largeDoc,
    });
    
    await page.click('[data-testid="mint-button"]');
    
    // Should show encryption progress
    await expect(page.locator('[data-testid="encryption-progress"]')).toBeVisible();
    
    // Wait for completion
    await expect(page.locator('[data-testid="success-message"]')).toBeVisible({ timeout: TEST_TIMEOUT * 2 });
  });
});

/**
 * Private Sale E2E Tests
 */
test.describe('Private Sale Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`${BASE_URL}/login`);
    await page.click('[data-testid="privy-login"]');
    await page.waitForSelector('[data-testid="user-authenticated"]', { timeout: TEST_TIMEOUT });
  });

  test('completes private sale flow', async ({ page }) => {
    // Navigate to property
    await page.goto(`${BASE_URL}/properties`);
    
    // Click on a property
    await page.click('[data-testid="property-card"]:first-child');
    
    // Click sell button
    await page.click('[data-testid="sell-button"]');
    
    // Configure sale
    await page.click('[data-testid="sale-type-private"]');
    await page.fill('[data-testid="sale-price"]', '1200000');
    await page.fill('[data-testid="allowed-buyer"]', '0xBuyerAddress123456789');
    await page.fill('[data-testid="sale-duration"]', '30');
    
    // Create listing
    await page.click('[data-testid="create-listing-button"]');
    
    // Wait for listing creation
    await expect(page.locator('[data-testid="listing-created"]')).toBeVisible({ timeout: TEST_TIMEOUT });
    
    // Complete sale (as CRE would)
    await page.click('[data-testid="complete-sale-button"]');
    
    // Wait for sale completion
    await expect(page.locator('[data-testid="sale-completed"]')).toBeVisible({ timeout: TEST_TIMEOUT });
    
    // Verify ownership transfer
    await expect(page.locator('[data-testid="new-owner-address"]')).toContainText('0xBuyer');
  });

  test('prevents unauthorized sale purchases', async ({ page }) => {
    await page.goto(`${BASE_URL}/marketplace`);
    
    // Try to buy a private sale property
    await page.click('[data-testid="private-sale-property"]:first-child');
    await page.click('[data-testid="buy-button"]');
    
    // Should show unauthorized error
    await expect(page.locator('[data-testid="unauthorized-error"]')).toContainText('Not authorized buyer');
  });

  test('handles sale disputes', async ({ page }) => {
    await page.goto(`${BASE_URL}/properties/my-property`);
    await page.click('[data-testid="open-dispute-button"]');
    
    // Fill dispute reason
    await page.fill('[data-testid="dispute-reason"]', 'Fraudulent sale attempt');
    await page.click('[data-testid="submit-dispute"]');
    
    // Verify dispute opened
    await expect(page.locator('[data-testid="dispute-status"]')).toContainText('Open');
  });
});

/**
 * Rental Flow E2E Tests
 */
test.describe('Rental Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`${BASE_URL}/login`);
    await page.click('[data-testid="privy-login"]');
    await page.waitForSelector('[data-testid="user-authenticated"]', { timeout: TEST_TIMEOUT });
  });

  test('creates rental agreement', async ({ page }) => {
    await page.goto(`${BASE_URL}/properties`);
    await page.click('[data-testid="property-card"]:first-child');
    
    // Start rental flow
    await page.click('[data-testid="rent-button"]');
    
    // Configure rental
    await page.fill('[data-testid="renter-address"]', '0xRenterAddress123');
    await page.fill('[data-testid="rental-duration"]', '365');
    await page.fill('[data-testid="deposit-amount"]', '5000');
    await page.fill('[data-testid="monthly-rent"]', '2500');
    
    // Submit rental
    await page.click('[data-testid="create-rental-button"]');
    
    // Wait for confirmation
    await expect(page.locator('[data-testid="rental-created"]')).toBeVisible({ timeout: TEST_TIMEOUT });
    await expect(page.locator('[data-testid="agreement-id"]')).not.toBeEmpty();
  });

  test('processes rental payment', async ({ page }) => {
    await page.goto(`${BASE_URL}/rentals/my-rental`);
    
    // Make payment
    await page.click('[data-testid="pay-rent-button"]');
    
    // Select payment method
    await page.click('[data-testid="payment-crypto"]');
    await page.click('[data-testid="confirm-payment"]');
    
    // Wait for confirmation
    await expect(page.locator('[data-testid="payment-success"]')).toBeVisible({ timeout: TEST_TIMEOUT });
  });

  test('handles late payment notifications', async ({ page }) => {
    await page.goto(`${BASE_URL}/rentals`);
    
    // Check for late payment indicator
    const lateBadge = page.locator('[data-testid="late-payment-badge"]');
    if (await lateBadge.isVisible()) {
      await expect(lateBadge).toContainText('Late');
      
      // Late fee should be displayed
      await expect(page.locator('[data-testid="late-fee-amount"]')).toBeVisible();
    }
  });

  test('ends rental and returns deposit', async ({ page }) => {
    await page.goto(`${BASE_URL}/rentals/my-rental`);
    
    // End rental
    await page.click('[data-testid="end-rental-button"]');
    
    // Confirm
    await page.click('[data-testid="confirm-end-rental"]');
    
    // Wait for deposit return
    await expect(page.locator('[data-testid="rental-ended"]')).toBeVisible({ timeout: TEST_TIMEOUT });
    await expect(page.locator('[data-testid="deposit-returned"]')).toContainText('Deposit returned');
  });
});

/**
 * Bill Payment E2E Tests
 */
test.describe('Bill Payment Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`${BASE_URL}/login`);
    await page.click('[data-testid="privy-login"]');
    await page.waitForSelector('[data-testid="user-authenticated"]', { timeout: TEST_TIMEOUT });
  });

  test('displays bill information', async ({ page }) => {
    await page.goto(`${BASE_URL}/bills`);
    
    // Verify bills are loaded
    await expect(page.locator('[data-testid="bills-list"]')).toBeVisible();
    
    // Check bill details
    const billCard = page.locator('[data-testid="bill-card"]:first-child');
    await expect(billCard.locator('[data-testid="bill-type"]')).not.toBeEmpty();
    await expect(billCard.locator('[data-testid="bill-amount"]')).not.toBeEmpty();
    await expect(billCard.locator('[data-testid="bill-due-date"]')).not.toBeEmpty();
  });

  test('pays bill with crypto', async ({ page }) => {
    await page.goto(`${BASE_URL}/bills`);
    
    // Click pay on first unpaid bill
    await page.click('[data-testid="bill-status-pending"] [data-testid="pay-button"]:first-child');
    
    // Select crypto payment
    await page.click('[data-testid="payment-method-crypto"]');
    await page.click('[data-testid="confirm-payment"]');
    
    // Wait for transaction
    await expect(page.locator('[data-testid="transaction-pending"]')).toBeVisible();
    await expect(page.locator('[data-testid="payment-success"]')).toBeVisible({ timeout: TEST_TIMEOUT });
    
    // Verify bill status updated
    await expect(page.locator('[data-testid="bill-status-paid"]')).toBeVisible();
  });

  test('pays bill with Stripe', async ({ page }) => {
    await page.goto(`${BASE_URL}/bills`);
    
    await page.click('[data-testid="bill-status-pending"] [data-testid="pay-button"]:first-child');
    
    // Select Stripe
    await page.click('[data-testid="payment-method-stripe"]');
    
    // Fill Stripe test card
    await page.fill('[data-testid="card-number"]', '4242424242424242');
    await page.fill('[data-testid="card-expiry"]', '12/25');
    await page.fill('[data-testid="card-cvc"]', '123');
    
    await page.click('[data-testid="submit-stripe-payment"]');
    
    // Wait for Stripe confirmation
    await expect(page.locator('[data-testid="stripe-success"]')).toBeVisible({ timeout: TEST_TIMEOUT });
  });

  test('disputes a bill', async ({ page }) => {
    await page.goto(`${BASE_URL}/bills`);
    
    // Open dispute
    await page.click('[data-testid="bill-card"]:first-child [data-testid="dispute-button"]');
    
    // Fill dispute reason
    await page.fill('[data-testid="dispute-reason"]', 'Incorrect amount charged');
    await page.click('[data-testid="submit-dispute"]');
    
    // Verify dispute created
    await expect(page.locator('[data-testid="dispute-submitted"]')).toBeVisible();
    await expect(page.locator('[data-testid="bill-status-disputed"]')).toContainText('Disputed');
  });
});

/**
 * KYC/AML E2E Tests
 */
test.describe('KYC Verification Flow', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto(`${BASE_URL}/login`);
    await page.click('[data-testid="privy-login"]');
    await page.waitForSelector('[data-testid="user-authenticated"]', { timeout: TEST_TIMEOUT });
  });

  test('shows KYC status', async ({ page }) => {
    await page.goto(`${BASE_URL}/profile`);
    
    // Check KYC status
    await expect(page.locator('[data-testid="kyc-status"]')).toBeVisible();
    
    const kycLevel = await page.locator('[data-testid="kyc-level"]').textContent();
    expect(['None', 'Basic', 'Full']).toContain(kycLevel);
  });

  test('initiates KYC verification', async ({ page }) => {
    await page.goto(`${BASE_URL}/profile`);
    
    await page.click('[data-testid="start-kyc-button"]');
    
    // Fill KYC form
    await page.fill('[data-testid="full-name"]', 'John Doe');
    await page.fill('[data-testid="date-of-birth"]', '1990-01-01');
    await page.fill('[data-testid="ssn"]', '123-45-6789');
    
    // Upload ID document
    const idInput = await page.locator('[data-testid="id-document-upload"]');
    await idInput.setInputFiles({
      name: 'id-card.pdf',
      mimeType: 'application/pdf',
      buffer: Buffer.from('ID document content'),
    });
    
    await page.click('[data-testid="submit-kyc"]');
    
    // Wait for verification
    await expect(page.locator('[data-testid="kyc-pending"]')).toBeVisible();
  });

  test('upgrades from level 1 to level 2', async ({ page }) => {
    await page.goto(`${BASE_URL}/profile`);
    
    // Assume level 1 already verified
    await page.click('[data-testid="upgrade-kyc-button"]');
    
    // Additional verification steps
    await page.fill('[data-testid="income-source"]', 'Employment');
    await page.fill('[data-testid="annual-income"]', '150000');
    
    // Upload proof of income
    const incomeInput = await page.locator('[data-testid="income-proof-upload"]');
    await incomeInput.setInputFiles({
      name: 'paystub.pdf',
      mimeType: 'application/pdf',
      buffer: Buffer.from('Paystub content'),
    });
    
    await page.click('[data-testid="submit-upgrade"]');
    
    await expect(page.locator('[data-testid="kyc-level-2-pending"]')).toBeVisible();
  });
});

/**
 * Error Handling and Edge Cases
 */
test.describe('Error Handling', () => {
  test('handles network errors gracefully', async ({ page }) => {
    // Simulate offline
    await page.context().setOffline(true);
    
    await page.goto(`${BASE_URL}/properties`);
    
    // Should show error message
    await expect(page.locator('[data-testid="network-error"]')).toBeVisible();
    await expect(page.locator('[data-testid="retry-button"]')).toBeVisible();
    
    // Restore network and retry
    await page.context().setOffline(false);
    await page.click('[data-testid="retry-button"]');
    
    // Should load successfully
    await expect(page.locator('[data-testid="properties-list"]')).toBeVisible();
  });

  test('handles transaction rejection', async ({ page }) => {
    await page.goto(`${BASE_URL}/mint`);
    
    // Fill form
    await page.fill('[data-testid="property-address"]', 'Test Address');
    await page.fill('[data-testid="property-price"]', '100000');
    
    // Upload document
    const deedInput = await page.locator('[data-testid="document-upload-deed"]');
    await deedInput.setInputFiles({
      name: 'deed.pdf',
      mimeType: 'application/pdf',
      buffer: Buffer.from('deed content'),
    });
    
    // Mock transaction rejection
    await page.route('**/api/mint', async route => {
      await route.fulfill({
        status: 400,
        body: JSON.stringify({ error: 'Transaction rejected by user' }),
      });
    });
    
    await page.click('[data-testid="mint-button"]');
    
    // Should show error
    await expect(page.locator('[data-testid="transaction-error"]')).toContainText('rejected');
  });

  test('handles session timeout', async ({ page }) => {
    await page.goto(`${BASE_URL}/properties`);
    
    // Simulate session expiration
    await page.evaluate(() => {
      localStorage.removeItem('privy:token');
    });
    
    // Try to perform action
    await page.click('[data-testid="mint-button"]');
    
    // Should redirect to login
    await expect(page).toHaveURL(`${BASE_URL}/login`);
  });
});
