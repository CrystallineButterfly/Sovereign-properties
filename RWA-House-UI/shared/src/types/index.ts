/**
 * Shared Type Definitions for RWA House Platform
 * Used across Web (React) and Mobile (React Native)
 */

// ==================== USER & AUTHENTICATION ====================

export interface User {
  id: string;
  email: string;
  walletAddress: string;
  chainId: number;
  kycStatus: KYCStatus;
  createdAt: Date;
  lastLoginAt: Date;
  mfaEnabled: boolean;
  preferences: UserPreferences;
}

export type KYCStatus = 'unverified' | 'pending' | 'verified' | 'rejected';
export type KYCProvider = 'none' | 'mock' | 'zkpassport';

export interface UserPreferences {
  theme: 'light' | 'dark' | 'system';
  notifications: NotificationPreferences;
  currency: string;
  language: string;
  autoPayEnabled: boolean;
  autoPayThreshold: number;
}

export interface NotificationPreferences {
  email: boolean;
  push: boolean;
  sms: boolean;
  transactions: boolean;
  bills: boolean;
  security: boolean;
}

// ==================== HOUSE & RWA ====================

export interface House {
  tokenId: string;
  houseId: string;
  ownerAddress: string;
  originalOwner: string;
  documentHash: string;
  documentURI: string;
  storageType: StorageType;
  mintedAt: Date;
  isVerified: boolean;
  metadata: HouseMetadata;
  listing?: Listing;
  rental?: RentalAgreement;
  bills: Bill[];
}

export type StorageType = 'ipfs' | 'offchain';

export interface HouseMetadata {
  address: string;
  city: string;
  state: string;
  zipCode: string;
  country: string;
  propertyType: PropertyType;
  bedrooms: number;
  bathrooms: number;
  squareFeet: number;
  yearBuilt: number;
  description: string;
  images: string[];
}

export type PropertyType = 
  | 'single_family' 
  | 'condo' 
  | 'townhouse' 
  | 'multi_family' 
  | 'apartment'
  | 'commercial';

// ==================== LISTINGS ====================

export interface Listing {
  tokenId: string;
  listingType: ListingType;
  price: string; // in wei
  priceFormatted: string;
  preferredToken: string; // token address (0x0 for ETH)
  isPrivateSale: boolean;
  allowedBuyer?: string;
  createdAt: Date;
  expiresAt?: Date;
}

export type ListingType = 'none' | 'for_sale' | 'for_rent';

export interface ListingFormData {
  listingType: ListingType;
  price: string;
  preferredToken: string;
  isPrivateSale: boolean;
  allowedBuyer?: string;
  duration?: number; // days for rental
}

// ==================== RENTALS ====================

export interface RentalAgreement {
  tokenId: string;
  renterAddress: string;
  startTime: Date;
  endTime: Date;
  depositAmount: string; // in wei
  depositFormatted: string;
  monthlyRent: string; // in wei
  isActive: boolean;
  hasAccessKey: boolean;
}

export interface RentalFormData {
  tokenId: string;
  renterAddress: string;
  durationDays: number;
  monthlyRent: string;
  depositMonths: number;
}

// ==================== BILLS ====================

export interface Bill {
  id: string;
  tokenId: string;
  billType: BillType;
  amount: number;
  amountFormatted: string;
  dueDate: Date;
  isPaid: boolean;
  paidAt?: Date;
  paymentMethod?: PaymentMethod;
  paymentReference?: string;
  isRecurring: boolean;
  provider: string;
  providerName: string;
  createdAt: Date;
}

export type BillType = 
  | 'electricity'
  | 'water'
  | 'gas'
  | 'internet'
  | 'phone'
  | 'property_tax'
  | 'insurance'
  | 'hoa'
  | 'maintenance'
  | 'other';

export type PaymentMethod = 'crypto' | 'stripe' | 'bank_transfer';

export interface BillPaymentData {
  tokenId: string;
  billIndex: number;
  paymentMethod: PaymentMethod;
  stripeToken?: string;
}

export interface CreateBillData {
  tokenId: string;
  billType: BillType;
  amount: number;
  dueDate: string; // ISO 8601
  provider: string;
  isRecurring: boolean;
}

// ==================== TRANSACTIONS ====================

export interface Transaction {
  hash: string;
  type: TransactionType;
  status: TransactionStatus;
  from: string;
  to: string;
  value: string;
  gasUsed?: string;
  gasPrice?: string;
  timestamp: Date;
  tokenId?: string;
  metadata?: Record<string, any>;
}

export type TransactionType = 
  | 'mint'
  | 'transfer'
  | 'sale'
  | 'rental_start'
  | 'rental_end'
  | 'bill_payment'
  | 'listing_create'
  | 'listing_cancel'
  | 'key_claim';

export type TransactionStatus = 'pending' | 'confirmed' | 'failed';

// ==================== DOCUMENTS ====================

export interface Document {
  id: string;
  houseId: string;
  name: string;
  type: DocumentType;
  size: number;
  mimeType: string;
  encrypted: boolean;
  uploadedAt: Date;
  hash: string;
  uri: string;
}

export type DocumentType = 
  | 'deed'
  | 'title'
  | 'inspection'
  | 'appraisal'
  | 'survey'
  | 'tax_record'
  | 'insurance'
  | 'hoa_agreement'
  | 'other';

export interface DocumentUploadData {
  houseId: string;
  documents: File[];
  storageType: StorageType;
}

export interface EncryptedKeyPackage {
  keyHash: string;
  encryptedKey: string;
  recipientAddress: string;
  expiresAt: Date;
}

// ==================== API REQUESTS/RESPONSES ====================

export interface APIResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
  txHash?: string;
  error?: APIError;
}

export interface APIError {
  code: string;
  message: string;
  details?: Record<string, any>;
}

export type ZKPassportSessionStatus =
  | 'pending'
  | 'ready'
  | 'verified'
  | 'failed'
  | 'expired';

export interface ZKPassportSession {
  sessionId: string;
  status: ZKPassportSessionStatus;
  domain?: string;
  mode?: "fast" | "compressed" | "compressed-evm";
  qrCodeUrl?: string;
  deepLinkUrl?: string;
  expiresAt?: string;
  proof?: Record<string, any>;
  message?: string;
}

export interface ZKPassportProofBundle {
  provider?: 'zkpassport';
  proofs: Record<string, any>[];
  queryResult: Record<string, any>;
  scope?: string;
  devMode?: boolean;
  validity?: number;
  uniqueIdentifier?: string;
  nullifierType?: string;
  verified?: boolean;
  verifiedAt?: string;
}

export interface ZKPassportVerificationResult {
  verified: boolean;
  level: number;
  verificationHash: string;
  expiresAt: number;
  uniqueIdentifier?: string;
  nullifierType?: string;
  domain?: string;
  scope?: string;
  proof: ZKPassportProofBundle;
}

// Mint Request
export interface MintRequestPayload {
  action: 'mint';
  ownerAddress: string;
  houseId: string;
  documentsB64: string;
  storageType: StorageType;
  ownerPublicKey: string;
  kycProvider?: KYCProvider;
  kycProof?: Record<string, any>;
  metadata: HouseMetadata;
}

export interface MintResponse {
  tokenId: string;
  txHash: string;
  encryptedKey: string;
  documentHash: string;
}

// Sell Request
export interface SellRequestPayload {
  action: 'sell';
  sellerAddress: string;
  buyerAddress: string;
  tokenId: string;
  price: string;
  buyerPublicKey: string;
  isPrivateSale: boolean;
  kycProvider?: KYCProvider;
  kycProof?: Record<string, any>;
}

export interface SellResponse {
  txHash: string;
  keyHash: string;
}

// Rent Request
export interface RentRequestPayload {
  action: 'rent';
  tokenId: string;
  renterAddress: string;
  durationDays: number;
  monthlyRent: string;
  depositAmount?: string;
  renterPublicKey: string;
  kycProvider?: KYCProvider;
  kycProof?: Record<string, any>;
}

export interface RentResponse {
  txHash: string;
  accessKeyHash: string;
}

export interface CreateListingRequestPayload {
  action: 'create_listing';
  ownerAddress: string;
  tokenId: string;
  listingType: 'for_sale' | 'for_rent';
  price: string;
  preferredToken?: string;
  isPrivateSale: boolean;
  allowedBuyer?: string;
  durationDays?: number;
  kycProvider?: KYCProvider;
  kycProof?: Record<string, any>;
}

// ==================== PRIVY & AUTH ====================

export interface PrivyUser {
  id: string;
  email?: string;
  wallet?: {
    address: string;
    chainId: number;
  };
  mfaMethods: MFAMethod[];
  createdAt: Date;
}

export interface MFAMethod {
  type: 'passkey' | 'totp' | 'sms';
  enabled: boolean;
  verifiedAt?: Date;
}

// ==================== STRIPE ====================

export interface StripePaymentIntent {
  clientSecret: string;
  amount: number;
  currency: string;
  status: string;
}

export interface StripeTokenData {
  token: string;
  cardLast4: string;
  cardBrand: string;
}

// ==================== SECURITY ====================

export interface SecurityEvent {
  id: string;
  userId: string;
  eventType: SecurityEventType;
  severity: 'low' | 'medium' | 'high' | 'critical';
  ipAddress: string;
  userAgent: string;
  timestamp: Date;
  metadata?: Record<string, any>;
}

export type SecurityEventType = 
  | 'login_success'
  | 'login_failure'
  | 'logout'
  | 'mfa_enabled'
  | 'mfa_disabled'
  | 'wallet_created'
  | 'transaction_signed'
  | 'suspicious_activity'
  | 'rate_limit_exceeded';

// ==================== NOTIFICATIONS ====================

export interface Notification {
  id: string;
  userId: string;
  type: NotificationType;
  title: string;
  message: string;
  data?: Record<string, any>;
  read: boolean;
  createdAt: Date;
}

export type NotificationType = 
  | 'transaction_confirmed'
  | 'transaction_failed'
  | 'bill_due'
  | 'bill_paid'
  | 'rental_expiring'
  | 'listing_sold'
  | 'key_claimed'
  | 'security_alert'
  | 'message_received';

// ==================== PRIVATE MESSAGING ====================

export type ConversationRole = 'buyer' | 'seller' | 'renter' | 'landlord';

export interface ConversationParticipant {
  walletAddress: string;
  role: ConversationRole;
}

export interface PrivateMessage {
  id: string;
  conversationId: string;
  tokenId: string;
  senderWalletAddress: string;
  recipientWalletAddress: string;
  content: string;
  transport: 'xmtp';
  xmtpMessageId?: string;
  createdAt: Date;
  readBy: string[];
}

export interface ConversationSummary {
  id: string;
  tokenId: string;
  houseId: string;
  participants: ConversationParticipant[];
  counterpartWalletAddress: string;
  counterpartRole: ConversationRole | '';
  unreadCount: number;
  createdAt: Date;
  updatedAt: Date;
  lastMessageAt?: Date;
  lastMessagePreview?: string;
  lastMessage?: PrivateMessage | null;
}

export interface ConversationDetails {
  conversation: ConversationSummary;
  messages: PrivateMessage[];
}

// ==================== FORM VALIDATION ====================

export interface ValidationError {
  field: string;
  message: string;
  code: string;
}

export interface FormState<T> {
  data: T;
  errors: ValidationError[];
  isValid: boolean;
  isSubmitting: boolean;
  isDirty: boolean;
}
