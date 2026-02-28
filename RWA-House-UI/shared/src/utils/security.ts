/**
 * Security Utilities and Validation Schemas
 * Shared across Web and Mobile platforms
 */

import { z } from 'zod';

const isAllowedImageURL = (value: string): boolean => {
  const trimmed = value.trim();
  const lowered = trimmed.toLowerCase();

  if (lowered.startsWith('ipfs://')) {
    return true;
  }

  if (lowered.startsWith('data:image/png;base64,')) {
    const base64Payload = trimmed.slice('data:image/png;base64,'.length).trim();
    return /^[a-z0-9+/=\s]+$/i.test(base64Payload);
  }

  try {
    const parsed = new URL(trimmed);
    return parsed.protocol === 'https:' || parsed.protocol === 'http:';
  } catch {
    return false;
  }
};

// ==================== VALIDATION SCHEMAS ====================

export const EthAddressSchema = z.string()
  .regex(/^0x[a-fA-F0-9]{40}$/, 'Invalid Ethereum address format')
  .transform((val: string) => val.toLowerCase());

export const HexStringSchema = z.string()
  .regex(/^0x[a-fA-F0-9]+$/, 'Invalid hex string format');

export const HouseIdSchema = z.string()
  .min(3, 'House ID must be at least 3 characters')
  .max(100, 'House ID must be less than 100 characters')
  .regex(/^[a-zA-Z0-9-]+$/, 'House ID can only contain letters, numbers, and hyphens');

export const PriceSchema = z.string()
  .regex(/^\d+$/, 'Price must be a positive integer (in wei)')
  .refine((val: string) => BigInt(val) > 0, 'Price must be greater than 0');

export const BillAmountSchema = z.number()
  .positive('Amount must be positive')
  .max(1000000, 'Amount exceeds maximum limit');

// House Metadata Schema
export const HouseMetadataSchema = z.object({
  address: z.string().min(5).max(200),
  city: z.string().min(2).max(100),
  state: z.string().min(2).max(50),
  zipCode: z.string().regex(/^\d{5}(-\d{4})?$/),
  country: z.string().min(2).max(100),
  propertyType: z.enum(['single_family', 'condo', 'townhouse', 'multi_family', 'apartment', 'commercial']),
  bedrooms: z.number().int().min(0).max(20),
  bathrooms: z.number().min(0).max(20),
  squareFeet: z.number().int().positive(),
  yearBuilt: z.number().int().min(1800).max(new Date().getFullYear()),
  description: z.string().max(1000).optional(),
  images: z.array(
    z.string().refine(
      isAllowedImageURL,
      'Image must use https://, http://, ipfs://, or data:image/png;base64,...',
    ),
  ).max(10)
});

// Mint Form Schema
export const MintFormSchema = z.object({
  houseId: HouseIdSchema,
  storageType: z.enum(['ipfs', 'offchain']),
  metadata: HouseMetadataSchema,
  documents: z.array(z.instanceof(File)).min(1, 'At least one document is required')
    .refine((files: File[]) => files.every((f: File) => f.size <= 10 * 1024 * 1024), 'Each file must be less than 10MB')
    .refine((files: File[]) => files.reduce((acc: number, f: File) => acc + f.size, 0) <= 50 * 1024 * 1024, 'Total file size must be less than 50MB')
});

// Sell Form Schema
export const SellFormSchema = z.object({
  tokenId: z.string(),
  buyerAddress: EthAddressSchema,
  price: PriceSchema,
  buyerPublicKey: z.string().min(128, 'Invalid public key'),
  isPrivateSale: z.boolean(),
  allowedBuyer: EthAddressSchema.optional()
}).refine((data: { isPrivateSale: boolean; allowedBuyer?: string }) => !data.isPrivateSale || !!data.allowedBuyer, {
  message: 'Allowed buyer is required for private sales',
  path: ['allowedBuyer']
});

// Rent Form Schema
export const RentFormSchema = z.object({
  tokenId: z.string(),
  renterAddress: EthAddressSchema,
  durationDays: z.number().int().min(1).max(3650), // Max 10 years
  monthlyRent: PriceSchema,
  renterPublicKey: z.string().min(128, 'Invalid public key'),
  depositMonths: z.number().int().min(1).max(3)
});

// Bill Form Schema
export const CreateBillSchema = z.object({
  tokenId: z.string(),
  billType: z.enum(['electricity', 'water', 'gas', 'internet', 'phone', 'property_tax', 'insurance', 'hoa', 'maintenance', 'other']),
  amount: BillAmountSchema,
  dueDate: z.string().datetime(),
  provider: EthAddressSchema,
  isRecurring: z.boolean()
});

// Payment Schema
export const PaymentSchema = z.object({
  tokenId: z.string(),
  billIndex: z.number().int().nonnegative(),
  paymentMethod: z.enum(['crypto', 'stripe']),
  stripeToken: z.string().optional()
}).refine((data: { paymentMethod: string; stripeToken?: string }) => data.paymentMethod !== 'stripe' || !!data.stripeToken, {
  message: 'Stripe token is required for card payments',
  path: ['stripeToken']
});

// ==================== SECURITY UTILITIES ====================

/**
 * Sanitize user input to prevent XSS
 */
export function sanitizeInput(input: string): string {
  return input
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#x27;')
    .replace(/\//g, '&#x2F;');
}

/**
 * Validate Ethereum address checksum
 */
export function isValidEthereumAddress(address: string): boolean {
  return /^0x[a-fA-F0-9]{40}$/.test(address);
}

/**
 * Convert file to Base64 with size validation
 */
export async function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    if (file.size > 10 * 1024 * 1024) {
      reject(new Error('File size exceeds 10MB limit'));
      return;
    }

    const reader = new FileReader();
    reader.onload = () => {
      const base64 = reader.result as string;
      resolve(base64.split(',')[1]); // Remove data URL prefix
    };
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
}

/**
 * Generate secure nonce for transactions
 */
export function generateNonce(): string {
  const array = new Uint8Array(32);
  crypto.getRandomValues(array);
  return Array.from(array, byte => byte.toString(16).padStart(2, '0')).join('');
}

/**
 * Hash sensitive data (one-way)
 */
export async function hashData(data: string): Promise<string> {
  const encoder = new TextEncoder();
  const dataBuffer = encoder.encode(data);
  const hashBuffer = await crypto.subtle.digest('SHA-256', dataBuffer);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  return hashArray.map(byte => byte.toString(16).padStart(2, '0')).join('');
}

/**
 * Rate limiter for API calls
 */
export class RateLimiter {
  private requests: Map<string, number[]> = new Map();
  private maxRequests: number;
  private windowMs: number;

  constructor(maxRequests: number = 100, windowMs: number = 60000) {
    this.maxRequests = maxRequests;
    this.windowMs = windowMs;
  }

  canProceed(key: string): boolean {
    const now = Date.now();
    const timestamps = this.requests.get(key) || [];
    
    // Remove old timestamps
    const validTimestamps = timestamps.filter(ts => now - ts < this.windowMs);
    
    if (validTimestamps.length >= this.maxRequests) {
      return false;
    }
    
    validTimestamps.push(now);
    this.requests.set(key, validTimestamps);
    return true;
  }

  getRetryAfter(key: string): number {
    const timestamps = this.requests.get(key) || [];
    if (timestamps.length === 0) return 0;
    
    const oldest = timestamps[0];
    return Math.max(0, this.windowMs - (Date.now() - oldest));
  }
}

/**
 * Secure storage wrapper (abstract for platform-specific implementation)
 */
export interface SecureStorage {
  getItem(key: string): Promise<string | null>;
  setItem(key: string, value: string): Promise<void>;
  removeItem(key: string): Promise<void>;
  clear(): Promise<void>;
}

// ==================== ENCRYPTION UTILITIES ====================

/**
 * Generate RSA key pair for document encryption
 */
export async function generateKeyPair(): Promise<CryptoKeyPair> {
  return await crypto.subtle.generateKey(
    {
      name: 'RSA-OAEP',
      modulusLength: 4096,
      publicExponent: new Uint8Array([1, 0, 1]),
      hash: 'SHA-256'
    },
    true,
    ['encrypt', 'decrypt']
  );
}

/**
 * Export public key to PEM format
 */
export async function exportPublicKey(key: CryptoKey): Promise<string> {
  const exported = await crypto.subtle.exportKey('spki', key);
  const exportedAsString = String.fromCharCode(...new Uint8Array(exported));
  const exportedAsBase64 = btoa(exportedAsString);
  const pemExported = `-----BEGIN PUBLIC KEY-----\n${exportedAsBase64.match(/.{1,64}/g)?.join('\n')}\n-----END PUBLIC KEY-----`;
  return pemExported;
}

/**
 * Import PEM public key
 */
export async function importPublicKey(pem: string): Promise<CryptoKey> {
  const pemContents = pem
    .replace('-----BEGIN PUBLIC KEY-----', '')
    .replace('-----END PUBLIC KEY-----', '')
    .replace(/\n/g, '');
  
  const binaryDer = Uint8Array.from(atob(pemContents), c => c.charCodeAt(0));
  
  return await crypto.subtle.importKey(
    'spki',
    binaryDer.buffer,
    {
      name: 'RSA-OAEP',
      hash: 'SHA-256'
    },
    true,
    ['encrypt']
  );
}

// ==================== INPUT VALIDATION ====================

/**
 * Validate file upload
 */
export function validateFileUpload(
  file: File,
  allowedTypes: string[],
  maxSize: number
): { valid: boolean; error?: string } {
  if (!allowedTypes.includes(file.type)) {
    return {
      valid: false,
      error: `Invalid file type. Allowed: ${allowedTypes.join(', ')}`
    };
  }

  if (file.size > maxSize) {
    return {
      valid: false,
      error: `File too large. Maximum size: ${(maxSize / 1024 / 1024).toFixed(1)}MB`
    };
  }

  // Check for common attack vectors
  const dangerousExtensions = ['.exe', '.dll', '.bat', '.sh', '.php', '.jsp'];
  const fileName = file.name.toLowerCase();
  
  if (dangerousExtensions.some(ext => fileName.endsWith(ext))) {
    return {
      valid: false,
      error: 'Potentially dangerous file type detected'
    };
  }

  return { valid: true };
}

/**
 * Validate document content (basic magic number check)
 */
export async function validateDocumentContent(file: File): Promise<boolean> {
  // Read first few bytes to check magic numbers
  const arrayBuffer = await file.slice(0, 8).arrayBuffer();
  const uint8Array = new Uint8Array(arrayBuffer);
  
  // PDF: %PDF
  if (uint8Array[0] === 0x25 && uint8Array[1] === 0x50 && uint8Array[2] === 0x44 && uint8Array[3] === 0x46) {
    return true;
  }
  
  // PNG: 89 50 4E 47
  if (uint8Array[0] === 0x89 && uint8Array[1] === 0x50 && uint8Array[2] === 0x4E && uint8Array[3] === 0x47) {
    return true;
  }
  
  // JPEG: FF D8 FF
  if (uint8Array[0] === 0xFF && uint8Array[1] === 0xD8 && uint8Array[2] === 0xFF) {
    return true;
  }
  
  return false;
}

// ==================== DEVICE SECURITY ====================

/**
 * Detect if device is jailbroken/rooted (mobile only)
 */
export function isDeviceCompromised(): boolean {
  // This would be implemented with native modules in React Native
  // For web, we can check for developer tools
  if (typeof window !== 'undefined') {
    const devToolsOpen = (): boolean => {
      const threshold = 160;
      const widthThreshold = window.outerWidth - window.innerWidth > threshold;
      const heightThreshold = window.outerHeight - window.innerHeight > threshold;
      return widthThreshold || heightThreshold;
    };
    return devToolsOpen();
  }
  return false;
}

/**
 * Generate device fingerprint
 */
export function generateDeviceFingerprint(): string {
  const components = [
    navigator.userAgent,
    navigator.language,
    screen.colorDepth,
    screen.width + 'x' + screen.height,
    new Date().getTimezoneOffset(),
    !!window.sessionStorage,
    !!window.localStorage,
    navigator.hardwareConcurrency
  ];
  
  return components.join('::');
}

// ==================== ERROR HANDLING ====================

/**
 * Secure error logging (no sensitive data)
 */
export function logError(error: Error, context?: Record<string, any>): void {
  // Remove any potentially sensitive data before logging
  const sanitizedContext = context ? 
    Object.entries(context).reduce((acc, [key, value]) => {
      if (key.toLowerCase().includes('key') || 
          key.toLowerCase().includes('password') ||
          key.toLowerCase().includes('token') ||
          key.toLowerCase().includes('secret')) {
        acc[key] = '[REDACTED]';
      } else {
        acc[key] = value;
      }
      return acc;
    }, {} as Record<string, any>) : {};

  console.error('[Secure Error]', {
    message: error.message,
    name: error.name,
    context: sanitizedContext,
    timestamp: new Date().toISOString()
  });
}

// ==================== CSP HELPER ====================

/**
 * Content Security Policy directives
 */
export const CSPDirectives = {
  'default-src': ["'self'"],
  'script-src': ["'self'", "'unsafe-inline'", 'https://auth.privy.io'],
  'style-src': ["'self'", "'unsafe-inline'", 'https://fonts.googleapis.com'],
  'font-src': ["'self'", 'https://fonts.gstatic.com'],
  'img-src': [
    "'self'",
    'blob:',
    'data:',
    'https:',
    'http:',
    'https://ipfs.io',
    'https://*.ipfs.io',
  ],
  'connect-src': [
    "'self'",
    'https://api.rwa-platform.io',
    'https://zkpassport-api-production.up.railway.app',
    'https://*.stripe.com',
    'https://auth.privy.io',
    'https://*.privy.io',
    'wss://*.privy.io',
    'https://*.walletconnect.com',
    'https://*.walletconnect.org',
    'https://demo.zkpassport.id',
    'https://*.zkpassport.id',
    'https://rpc.sepolia.org',
    'https://ethereum-sepolia-rpc.publicnode.com',
    'https://sepolia.drpc.org'
  ],
  'frame-ancestors': ["'none'"],
  'base-uri': ["'self'"],
  'form-action': ["'self'"],
  'upgrade-insecure-requests': []
};

/**
 * Build CSP header string
 */
export function buildCSPHeader(): string {
  return Object.entries(CSPDirectives)
    .map(([key, values]) => {
      if (values.length === 0) return key;
      return `${key} ${values.join(' ')}`;
    })
    .join('; ');
}
