/**
 * API Client for RWA House Platform
 * Works with both Web (fetch) and Mobile (axios/react-native)
 */

import {
  APIResponse,
  MintRequestPayload,
  SellRequestPayload,
  RentRequestPayload,
  BillPaymentData,
  CreateBillData,
  User,
  House,
  HouseMetadata,
  Transaction,
  Notification,
  ConversationSummary,
  ConversationDetails,
  KYCProvider,
  ZKPassportSession,
  ZKPassportProofBundle,
  ZKPassportVerificationResult,
  CreateListingRequestPayload,
} from "../types";
import { ethers } from "ethers";

export const KYC_PROVIDER_STORAGE_KEY = "RWA_KYC_PROVIDER";
export const KYC_PROOF_STORAGE_KEY = "RWA_KYC_PROOF";
export const AUTH_EXPIRED_EVENT_NAME = "rwa:auth-expired";
const HOUSE_METADATA_CACHE_KEY = "RWA_HOUSE_METADATA_CACHE_V1";
const HOUSE_METADATA_CACHE_LIMIT = 200;
const DEFAULT_HOSTED_API_URL = "https://api.rwa-platform.io";

type HouseMetadataCacheEntry = {
  metadata: HouseMetadata;
  updatedAt: number;
};

type HouseMetadataCache = {
  byTokenId: Record<string, HouseMetadataCacheEntry>;
  byHouseId: Record<string, HouseMetadataCacheEntry>;
};

const readEnv = (key: string): string | undefined => {
  const candidateKeys = [key];

  if (key.startsWith("VITE_")) {
    candidateKeys.push(`EXPO_PUBLIC_${key.slice("VITE_".length)}`);
  }
  if (key.startsWith("REACT_APP_")) {
    candidateKeys.push(`EXPO_PUBLIC_${key.slice("REACT_APP_".length)}`);
  }

  const viteEnv = (globalThis as any)?.import?.meta?.env;
  if (viteEnv) {
    for (const envKey of candidateKeys) {
      if (typeof viteEnv[envKey] === "string" && viteEnv[envKey]) {
        return String(viteEnv[envKey]);
      }
    }
  }

  if (typeof process !== "undefined") {
    for (const envKey of candidateKeys) {
      const envVal = (process as any)?.env?.[envKey];
      if (typeof envVal === "string" && envVal) {
        return envVal;
      }
    }
  }

  return undefined;
};

const readStorage = (key: string): string | undefined => {
  try {
    const storage = (globalThis as any)?.localStorage;
    if (!storage) return undefined;
    const value = storage.getItem(key);
    return typeof value === "string" && value.trim() ? value : undefined;
  } catch {
    return undefined;
  }
};

const writeStorage = (key: string, value: string): void => {
  try {
    const storage = (globalThis as any)?.localStorage;
    if (!storage) {
      return;
    }
    storage.setItem(key, value);
  } catch {
    // Ignore writes in storage-restricted environments.
  }
};

const trimTrailingSlash = (value: string): string => {
  return value.replace(/\/+$/, "");
};

const isLoopbackHost = (host: string): boolean => {
  const normalized = host.trim().toLowerCase();
  return (
    normalized === "localhost"
    || normalized === "127.0.0.1"
    || normalized === "::1"
    || normalized === "[::1]"
    || normalized === "0.0.0.0"
  );
};

const resolveLocalAdapterURL = (): string | undefined => {
  const location = (globalThis as any)?.location;
  const hostname =
    typeof location?.hostname === "string" ? location.hostname : "";
  if (!hostname) {
    return undefined;
  }

  if (!isLoopbackHost(hostname)) {
    return undefined;
  }

  const protocol = location?.protocol === "https:" ? "https" : "http";
  return `${protocol}://${hostname}:8787`;
};

const resolveBaseURL = (): string => {
  const configured =
    readEnv("VITE_API_URL")
    || readEnv("REACT_APP_API_URL")
    || readEnv("VITE_ZKPASSPORT_API_URL")
    || readEnv("REACT_APP_ZKPASSPORT_API_URL");
  if (configured) {
    return trimTrailingSlash(configured);
  }

  const localAdapter = resolveLocalAdapterURL();
  if (localAdapter) {
    return localAdapter;
  }

  return DEFAULT_HOSTED_API_URL;
};

const isDefaultHostedAPIBase = (): boolean => {
  return !readEnv("VITE_API_URL") && !readEnv("REACT_APP_API_URL");
};

const resolveRpcURL = (baseURL: string): string => {
  const configured = readEnv("VITE_RPC_URL") || readEnv("REACT_APP_RPC_URL");
  if (configured) {
    return trimTrailingSlash(configured);
  }

  if (!isDefaultHostedAPIBase()) {
    return `${trimTrailingSlash(baseURL)}/rpc`;
  }

  return "https://sepolia.drpc.org";
};

const resolveZKPassportBaseURL = (): string => {
  const configured =
    readEnv("VITE_ZKPASSPORT_API_URL") ||
    readEnv("REACT_APP_ZKPASSPORT_API_URL");
  if (configured) {
    return trimTrailingSlash(configured);
  }

  const localAdapter = resolveLocalAdapterURL();
  if (localAdapter) {
    return localAdapter;
  }

  return resolveBaseURL();
};

const parseBooleanEnv = (value: string | undefined, fallback: boolean): boolean => {
  if (!value) {
    return fallback;
  }
  const normalized = value.trim().toLowerCase();
  if (normalized === "true" || normalized === "1" || normalized === "yes") {
    return true;
  }
  if (normalized === "false" || normalized === "0" || normalized === "no") {
    return false;
  }
  return fallback;
};

const parseIntegerEnv = (
  value: string | undefined,
  fallback: number,
  minimum: number,
): number => {
  if (!value) {
    return fallback;
  }
  const parsed = Number.parseInt(value.trim(), 10);
  if (!Number.isFinite(parsed) || parsed < minimum) {
    return fallback;
  }
  return parsed;
};

// API Configuration
const API_CONFIG = {
  // Prefer explicit env. On localhost dev, default to the local adapter service.
  baseURL: resolveBaseURL(),
  zkPassportBaseURL: resolveZKPassportBaseURL(),
  timeout: 30000,
  retries: 3,
  retryDelay: 1000,
};

const HAS_EXPLICIT_API_CONFIG = Boolean(
  readEnv("VITE_API_URL")
  || readEnv("REACT_APP_API_URL")
  || readEnv("VITE_ZKPASSPORT_API_URL")
  || readEnv("REACT_APP_ZKPASSPORT_API_URL"),
);

const CHAIN_RPC_FALLBACK_ENABLED = parseBooleanEnv(
  readEnv("VITE_ENABLE_CHAIN_RPC_FALLBACK")
    || readEnv("REACT_APP_ENABLE_CHAIN_RPC_FALLBACK"),
  false,
);

const CHAIN_RPC_PUBLIC_CANDIDATES_ENABLED = parseBooleanEnv(
  readEnv("VITE_ENABLE_PUBLIC_RPC_CANDIDATES")
    || readEnv("REACT_APP_ENABLE_PUBLIC_RPC_CANDIDATES"),
  false,
);

const CHAIN_RPC_NETWORK = {
  chainId: parseIntegerEnv(
    readEnv("VITE_CHAIN_ID") || readEnv("REACT_APP_CHAIN_ID"),
    11155111,
    1,
  ),
  name: (
    readEnv("VITE_CHAIN_NAME")
    || readEnv("REACT_APP_CHAIN_NAME")
    || "sepolia"
  ).trim(),
};

const BLOCKED_RPC_CANDIDATES = new Set<string>();

let CHAIN_CONFIG = {
  rpcURL: resolveRpcURL(API_CONFIG.baseURL),
  houseRWAAddress:
    readEnv("VITE_HOUSE_RWA_ADDRESS") ||
    readEnv("REACT_APP_HOUSE_RWA_ADDRESS") ||
    "",
  maxHouseScan: Math.max(
    1,
    Number(readEnv("VITE_MAX_HOUSE_SCAN") || "500") || 500,
  ),
};

const parseKYCProvider = (value: string | undefined): KYCProvider => {
  const normalized = value?.trim().toLowerCase();
  if (normalized === "none") {
    return "none";
  }
  return normalized === "zkpassport" ? "zkpassport" : "mock";
};

const getKYCProviderDefault = (): KYCProvider =>
  parseKYCProvider(
    readStorage(KYC_PROVIDER_STORAGE_KEY) ||
      readEnv("VITE_KYC_PROVIDER") ||
      readEnv("REACT_APP_KYC_PROVIDER"),
  );

const getKYCProofDefaultRaw = (): string | undefined =>
  readStorage(KYC_PROOF_STORAGE_KEY) ||
  readEnv("VITE_KYC_PROOF") ||
  readEnv("REACT_APP_KYC_PROOF");

const HOUSE_RWA_READ_ABI = [
  "function nextTokenId() view returns (uint256)",
  "function ownerOf(uint256 tokenId) view returns (address)",
  "function getHouseDetails(uint256 tokenId) view returns ((string houseId,bytes32 documentHash,string documentURI,uint8 storageType,address originalOwner,uint48 mintedAt,bool isVerified,uint8 documentCount))",
  "function getListing(uint256 tokenId) view returns ((uint8 listingType,uint96 price,address preferredToken,bool isPrivateSale,address allowedBuyer,uint48 createdAt,uint48 expiresAt,uint8 platformFee))",
  "function getBills(uint256 tokenId) view returns ((string billType,uint96 amount,uint48 dueDate,uint48 paidAt,uint8 status,bytes32 paymentReference,bool isRecurring,address provider,uint8 recurrenceInterval)[])",
  "function getActiveRental(uint256 tokenId) view returns ((address renter,uint48 startTime,uint48 endTime,uint96 depositAmount,uint96 monthlyRent,bool isActive,bytes32 encryptedAccessKeyHash,uint8 disputeStatus))",
];

// Request signer interface (for HMAC or JWT signing)
interface RequestSigner {
  signRequest(method: string, path: string, body: string): Promise<string>;
}

// API Client class
export class RWAApiClient {
  private baseURL: string;
  private zkPassportBaseURL: string | null =
    API_CONFIG.zkPassportBaseURL || null;
  private explicitApiConfigured: boolean = HAS_EXPLICIT_API_CONFIG;
  private authToken: string | null = null;
  private requestSigner: RequestSigner | null = null;
  private rateLimiter: Map<string, number> = new Map();
  private kycProviderOverride: KYCProvider | null = null;
  private kycProofOverride: Record<string, any> | null = null;

  private notifyAuthExpired(reason: string) {
    this.authToken = null;
    if (typeof window === "undefined" || typeof window.dispatchEvent !== "function") {
      return;
    }
    window.dispatchEvent(
      new CustomEvent(AUTH_EXPIRED_EVENT_NAME, {
        detail: { reason },
      }),
    );
  }

  constructor(baseURL: string = API_CONFIG.baseURL) {
    this.baseURL = baseURL;
    if (baseURL && trimTrailingSlash(baseURL) !== DEFAULT_HOSTED_API_URL) {
      this.explicitApiConfigured = true;
    }
  }

  // Get current base URL
  getBaseURL(): string {
    return this.baseURL;
  }

  // Set base URL (useful for Vite/Expo env configuration)
  setBaseURL(baseURL: string) {
    this.baseURL = baseURL;
    if (baseURL.trim().length > 0) {
      this.explicitApiConfigured = true;
    }
  }

  // Get current ZKPassport session API URL (falls back to main API URL)
  getZKPassportBaseURL(): string {
    return this.zkPassportBaseURL || this.baseURL;
  }

  // Set dedicated ZKPassport session API URL
  setZKPassportBaseURL(baseURL: string) {
    this.zkPassportBaseURL = baseURL || null;
  }

  setChainConfig(config: {
    rpcURL?: string;
    houseRWAAddress?: string;
    maxHouseScan?: number;
  }) {
    const nextConfig = { ...CHAIN_CONFIG };

    if (typeof config.rpcURL === "string" && config.rpcURL.trim()) {
      nextConfig.rpcURL = config.rpcURL.trim();
    }

    if (
      typeof config.houseRWAAddress === "string" &&
      config.houseRWAAddress.trim()
    ) {
      nextConfig.houseRWAAddress = config.houseRWAAddress.trim();
    }

    if (
      typeof config.maxHouseScan === "number" &&
      Number.isFinite(config.maxHouseScan) &&
      config.maxHouseScan > 0
    ) {
      nextConfig.maxHouseScan = Math.max(1, Math.floor(config.maxHouseScan));
    }

    CHAIN_CONFIG = nextConfig;
  }

  // Set authentication token
  setAuthToken(token: string) {
    const normalized = token.trim();
    this.authToken = normalized || null;
  }

  hasAuthToken(): boolean {
    return typeof this.authToken === "string" && this.authToken.length > 0;
  }

  // Set request signer for additional security
  setRequestSigner(signer: RequestSigner) {
    this.requestSigner = signer;
  }

  setKYCDefaults(provider?: KYCProvider, proof?: Record<string, any> | null) {
    this.kycProviderOverride = provider || null;
    if (provider === "none") {
      this.kycProofOverride = null;
      return;
    }
    this.kycProofOverride = proof || null;
  }

  private getRpcCandidates(): string[] {
    const candidates: string[] = [];
    const pushCandidate = (value: string | null | undefined) => {
      if (typeof value !== "string") return;
      const trimmed = value.trim();
      if (!trimmed) return;
      const normalized = trimTrailingSlash(trimmed);
      if (BLOCKED_RPC_CANDIDATES.has(normalized)) {
        return;
      }
      candidates.push(normalized);
    };
    const pushBaseAsRpc = (value: string | null | undefined) => {
      if (typeof value !== "string") return;
      const trimmed = value.trim();
      if (!trimmed) return;
      const normalized = trimTrailingSlash(trimmed);
      const rpcCandidate = normalized.endsWith("/rpc")
        ? normalized
        : `${normalized}/rpc`;
      if (!BLOCKED_RPC_CANDIDATES.has(rpcCandidate)) {
        candidates.push(rpcCandidate);
      }
    };

    pushCandidate(CHAIN_CONFIG.rpcURL);
    pushBaseAsRpc(this.baseURL);
    pushBaseAsRpc(this.zkPassportBaseURL);

    const envFallbackRPCs =
      readEnv("VITE_FALLBACK_RPC_URLS") || readEnv("REACT_APP_FALLBACK_RPC_URLS");
    if (envFallbackRPCs) {
      envFallbackRPCs
        .split(",")
        .map((entry) => entry.trim())
        .forEach((entry) => pushCandidate(entry));
    }

    if (CHAIN_RPC_PUBLIC_CANDIDATES_ENABLED) {
      pushCandidate("https://rpc.sepolia.org");
      pushCandidate("https://ethereum-sepolia-rpc.publicnode.com");
      pushCandidate("https://sepolia.drpc.org");
    }

    return Array.from(new Set(candidates));
  }

  private createRpcProvider(candidate: string): ethers.JsonRpcProvider {
    return new ethers.JsonRpcProvider(
      candidate,
      CHAIN_RPC_NETWORK,
      { staticNetwork: true },
    );
  }

  private isCspBlockedError(error: unknown): boolean {
    const message = error instanceof Error ? error.message : String(error ?? "");
    const normalized = message.toLowerCase();
    return (
      normalized.includes("content security policy") ||
      normalized.includes("violates the document's content security policy") ||
      normalized.includes("refused to connect because it violates")
    );
  }

  private async withRpcFallback<T>(
    operation: (provider: ethers.JsonRpcProvider) => Promise<T>,
    contextMessage: string,
  ): Promise<T> {
    const candidates = this.getRpcCandidates();
    if (candidates.length === 0) {
      throw new Error(
        `${contextMessage}. No usable RPC candidates remained after CSP filtering.`,
      );
    }
    let lastError: unknown = null;
    const attempts: string[] = [];

    for (const candidate of candidates) {
      const provider = this.createRpcProvider(candidate);
      try {
        const result = await operation(provider);
        BLOCKED_RPC_CANDIDATES.delete(candidate);
        CHAIN_CONFIG = { ...CHAIN_CONFIG, rpcURL: candidate };
        return result;
      } catch (error) {
        lastError = error;
        if (this.isCspBlockedError(error)) {
          BLOCKED_RPC_CANDIDATES.add(candidate);
        }
        const detail = error instanceof Error ? error.message : String(error);
        attempts.push(`${candidate} => ${detail}`);
      } finally {
        provider.destroy();
      }
    }

    const fallbackMessage = attempts.length
      ? `${contextMessage}. RPC attempts: ${attempts.join(" | ")}`
      : contextMessage;
    if (lastError instanceof Error) {
      throw new Error(`${fallbackMessage}. Last error: ${lastError.message}`);
    }
    throw new Error(`${fallbackMessage}. Last error: unknown.`);
  }

  // Base request method
  private async request<T>(
    method: string,
    path: string,
    body?: any,
    customHeaders?: Record<string, string>,
    baseURLOverride?: string,
    requestOptions?: {
      credentialsMode?: RequestCredentials;
    },
  ): Promise<APIResponse<T>> {
    const baseURL = baseURLOverride || this.baseURL;
    const url = `${baseURL}${path}`;
    const requestId = this.generateRequestId();

    // Check rate limiting
    if (!this.checkRateLimit(path)) {
      throw new Error("Rate limit exceeded. Please try again later.");
    }

    // Prepare headers
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      "X-Request-ID": requestId,
      "X-Timestamp": Date.now().toString(),
      ...customHeaders,
    };

    // Add auth token if available
    if (this.authToken) {
      headers["Authorization"] = `Bearer ${this.authToken}`;
    }

    // Sign request if signer is configured
    if (this.requestSigner) {
      const signature = await this.requestSigner.signRequest(
        method,
        path,
        body ? JSON.stringify(body) : "",
      );
      headers["X-Signature"] = signature;
    }

    // Retry logic
    let lastError: Error | null = null;

    for (let attempt = 0; attempt < API_CONFIG.retries; attempt++) {
      try {
        const response = await fetch(url, {
          method,
          headers,
          body: body ? JSON.stringify(body) : undefined,
          credentials: requestOptions?.credentialsMode || "include",
        });

        // Handle HTTP errors
        if (!response.ok) {
          const errorData = await response.json().catch(() => ({}));
          if (response.status === 401) {
            this.notifyAuthExpired(
              errorData.message || "Authentication expired. Please sign in again.",
            );
          }
          throw new Error(
            errorData.message || `HTTP Error: ${response.status}`,
          );
        }

        const data: APIResponse<T> = await response.json();

        // Validate response structure
        if (!this.isValidResponse(data)) {
          throw new Error("Invalid response structure from server");
        }

        return data;
      } catch (error) {
        lastError = error as Error;

        // Don't retry on client errors (4xx)
        if (error instanceof Error && error.message.includes("4")) {
          throw error;
        }

        // Wait before retry
        if (attempt < API_CONFIG.retries - 1) {
          await this.delay(API_CONFIG.retryDelay * Math.pow(2, attempt));
        }
      }
    }

    throw lastError || new Error("Request failed after retries");
  }

  private isApiConfigured(): boolean {
    return this.explicitApiConfigured;
  }

  private allowChainRpcFallback(): boolean {
    return CHAIN_RPC_FALLBACK_ENABLED;
  }

  private async requestPublic<T>(
    path: string,
    baseURL: string,
    options?: {
      includeAuthHeader?: boolean;
    },
  ): Promise<APIResponse<T>> {
    const baseCandidates = this.getLoopbackBaseURLCandidates(baseURL);
    let lastError: unknown = null;

    for (const candidate of baseCandidates) {
      try {
        const headers: Record<string, string> = {
          Accept: "application/json",
        };
        if (options?.includeAuthHeader && this.authToken) {
          headers.Authorization = `Bearer ${this.authToken}`;
        }

        const response = await fetch(`${candidate}${path}`, {
          method: "GET",
          headers,
          credentials: "omit",
        });

        let payload: any = null;
        try {
          payload = await response.json();
        } catch {
          payload = null;
        }

        if (!response.ok) {
          if (response.status === 401 && options?.includeAuthHeader) {
            this.notifyAuthExpired(
              payload?.message || "Authentication expired. Please sign in again.",
            );
          }
          return {
            success: false,
            message:
              payload?.message
              || `HTTP Error: ${response.status} (${response.statusText})`,
          };
        }

        if (!this.isValidResponse(payload)) {
          return {
            success: false,
            message: "Invalid response structure from server",
          };
        }

        return payload as APIResponse<T>;
      } catch (error) {
        lastError = error;
      }
    }

    const message =
      lastError instanceof Error
        ? lastError.message
        : "Public API request failed";
    throw new Error(message);
  }

  private isLikelyNetworkError(error: unknown): boolean {
    if (error instanceof TypeError) {
      return true;
    }

    const message = error instanceof Error ? error.message : String(error ?? "");
    const normalized = message.toLowerCase();
    return (
      normalized.includes("failed to fetch") ||
      normalized.includes("networkerror") ||
      normalized.includes("network request failed") ||
      normalized.includes("load failed")
    );
  }

  private getLoopbackBaseURLCandidates(baseURL: string): string[] {
    const candidates = [baseURL];

    try {
      const parsed = new URL(baseURL);
      const host = parsed.hostname.toLowerCase();
      const fallbackHosts: string[] = [];

      if (host === "localhost") {
        fallbackHosts.push("127.0.0.1", "::1");
      } else if (host === "127.0.0.1") {
        fallbackHosts.push("localhost", "::1");
      } else if (host === "::1") {
        fallbackHosts.push("localhost", "127.0.0.1");
      } else if (host === "0.0.0.0") {
        fallbackHosts.push("localhost", "127.0.0.1", "::1");
      }

      const runtimeHostname =
        typeof globalThis !== "undefined" &&
        typeof (globalThis as any)?.location?.hostname === "string"
          ? String((globalThis as any).location.hostname).trim()
          : "";
      if (runtimeHostname && !fallbackHosts.includes(runtimeHostname)) {
        fallbackHosts.push(runtimeHostname);
      }

      const runtimeOrigin =
        typeof globalThis !== "undefined"
        && typeof (globalThis as any)?.location?.origin === "string"
          ? String((globalThis as any).location.origin).trim()
          : "";
      if (runtimeOrigin) {
        try {
          const runtimeURL = new URL(runtimeOrigin);
          if (
            isLoopbackHost(runtimeURL.hostname)
            && !candidates.includes(trimTrailingSlash(runtimeOrigin))
          ) {
            candidates.push(trimTrailingSlash(runtimeOrigin));
          }
        } catch {
          // Ignore malformed runtime origins.
        }
      }

      for (const fallbackHost of fallbackHosts) {
        const next = new URL(baseURL);
        next.hostname = fallbackHost;
        candidates.push(trimTrailingSlash(next.toString()));
      }
    } catch {
      return candidates;
    }

    return Array.from(new Set(candidates));
  }

  private async requestWithLoopbackFallback<T>(
    method: string,
    path: string,
    body: any,
    customHeaders: Record<string, string> | undefined,
    baseURL: string,
    requestOptions?: {
      credentialsMode?: RequestCredentials;
    },
  ): Promise<APIResponse<T>> {
    const baseCandidates = this.getLoopbackBaseURLCandidates(baseURL);
    let lastError: unknown = null;
    const attempted: string[] = [];

    for (const candidate of baseCandidates) {
      try {
        return await this.request(
          method,
          path,
          body,
          customHeaders,
          candidate,
          requestOptions,
        );
      } catch (error) {
        lastError = error;
        const detail = error instanceof Error ? error.message : String(error);
        attempted.push(`${candidate} => ${detail}`);
        if (!this.isLikelyNetworkError(error)) {
          throw error;
        }
      }
    }

    if (lastError instanceof Error) {
      throw new Error(
        `${lastError.message}. Tried: ${attempted.join(" | ")}`,
      );
    }

    throw new Error(
      `Request failed for all loopback candidates. Tried: ${attempted.join(" | ")}`,
    );
  }

  // Rate limiting check
  private checkRateLimit(path: string): boolean {
    const now = Date.now();
    const key = `${path}`;
    const lastRequest = this.rateLimiter.get(key);

    if (lastRequest && now - lastRequest < 100) {
      // 100ms minimum between requests
      return false;
    }

    this.rateLimiter.set(key, now);
    return true;
  }

  // Generate unique request ID
  private generateRequestId(): string {
    return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
  }

  // Delay utility
  private delay(ms: number): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  // Validate API response structure
  private isValidResponse(response: any): boolean {
    return (
      response &&
      typeof response === "object" &&
      "success" in response &&
      typeof response.success === "boolean" &&
      "message" in response &&
      typeof response.message === "string"
    );
  }

  private isZeroAddress(address: string | undefined | null): boolean {
    if (!address) return true;
    return (
      address.toLowerCase() === "0x0000000000000000000000000000000000000000"
    );
  }

  private mapListingType(value: number): "none" | "for_sale" | "for_rent" {
    if (value === 1) return "for_sale";
    if (value === 2) return "for_rent";
    return "none";
  }

  private mapStorageType(value: number): "ipfs" | "offchain" {
    if (value === 1) return "offchain";
    return "ipfs";
  }

  private normalizeBillType(
    value: string,
  ):
    | "electricity"
    | "water"
    | "gas"
    | "internet"
    | "phone"
    | "property_tax"
    | "insurance"
    | "hoa"
    | "maintenance"
    | "other" {
    const v = value.toLowerCase();
    if (v === "electricity") return "electricity";
    if (v === "water" || v === "sewer") return "water";
    if (v === "gas") return "gas";
    if (v === "internet") return "internet";
    if (v === "phone") return "phone";
    if (v === "property_tax" || v === "tax") return "property_tax";
    if (v === "insurance") return "insurance";
    if (v === "hoa") return "hoa";
    if (v === "maintenance" || v === "utilities") return "maintenance";
    return "other";
  }

  private toBigInt(value: any): bigint {
    if (typeof value === "bigint") return value;
    if (typeof value === "number")
      return BigInt(Math.max(0, Math.trunc(value)));
    if (typeof value === "string" && value !== "") return BigInt(value);
    return 0n;
  }

  private toDateFromSeconds(value: any): Date {
    const seconds = this.toBigInt(value);
    if (seconds <= 0n) return new Date(0);
    const ms = seconds * 1000n;
    if (ms > BigInt(Number.MAX_SAFE_INTEGER)) {
      return new Date(Number.MAX_SAFE_INTEGER);
    }
    return new Date(Number(ms));
  }

  private createEmptyHouseMetadataCache(): HouseMetadataCache {
    return {
      byTokenId: {},
      byHouseId: {},
    };
  }

  private parseCachedMetadata(value: unknown): HouseMetadata | null {
    if (!value || typeof value !== "object") {
      return null;
    }

    const candidate = value as Partial<HouseMetadata>;
    if (typeof candidate.address !== "string" || candidate.address.trim().length === 0) {
      return null;
    }
    if (typeof candidate.city !== "string" || typeof candidate.state !== "string") {
      return null;
    }

    const fallback = this.buildDefaultMetadata(candidate.address);
    const rawImages = Array.isArray(candidate.images) ? candidate.images : [];
    return {
      ...fallback,
      ...candidate,
      images: rawImages.filter(
        (image): image is string =>
          typeof image === "string" && image.trim().length > 0,
      ),
    };
  }

  private parseHouseMetadataCacheMap(
    value: unknown,
  ): Record<string, HouseMetadataCacheEntry> {
    if (!value || typeof value !== "object") {
      return {};
    }

    const entries = Object.entries(value as Record<string, unknown>);
    return entries.reduce<Record<string, HouseMetadataCacheEntry>>(
      (accumulator, [key, rawEntry]) => {
        if (!key.trim() || !rawEntry || typeof rawEntry !== "object") {
          return accumulator;
        }

        const maybeEntry = rawEntry as {
          metadata?: unknown;
          updatedAt?: unknown;
        };
        const metadata = this.parseCachedMetadata(
          maybeEntry.metadata ?? rawEntry,
        );
        if (!metadata) {
          return accumulator;
        }

        const updatedAt =
          typeof maybeEntry.updatedAt === "number" &&
          Number.isFinite(maybeEntry.updatedAt)
            ? maybeEntry.updatedAt
            : Date.now();
        accumulator[key] = { metadata, updatedAt };
        return accumulator;
      },
      {},
    );
  }

  private readHouseMetadataCache(): HouseMetadataCache {
    const rawCache = readStorage(HOUSE_METADATA_CACHE_KEY);
    if (!rawCache) {
      return this.createEmptyHouseMetadataCache();
    }

    try {
      const parsed = JSON.parse(rawCache) as Partial<HouseMetadataCache>;
      return {
        byTokenId: this.parseHouseMetadataCacheMap(parsed.byTokenId),
        byHouseId: this.parseHouseMetadataCacheMap(parsed.byHouseId),
      };
    } catch {
      return this.createEmptyHouseMetadataCache();
    }
  }

  private trimCacheMap(
    cacheMap: Record<string, HouseMetadataCacheEntry>,
  ): Record<string, HouseMetadataCacheEntry> {
    const sorted = Object.entries(cacheMap).sort(
      (a, b) => b[1].updatedAt - a[1].updatedAt,
    );
    return Object.fromEntries(
      sorted.slice(0, HOUSE_METADATA_CACHE_LIMIT),
    ) as Record<string, HouseMetadataCacheEntry>;
  }

  private writeHouseMetadataCache(cache: HouseMetadataCache): void {
    const payload: HouseMetadataCache = {
      byTokenId: this.trimCacheMap(cache.byTokenId),
      byHouseId: this.trimCacheMap(cache.byHouseId),
    };
    writeStorage(HOUSE_METADATA_CACHE_KEY, JSON.stringify(payload));
  }

  private cacheMintedHouseMetadata(
    tokenId: string | null,
    houseId: string,
    metadata: HouseMetadata,
  ): void {
    const normalizedTokenId = tokenId ? tokenId.trim() : "";
    const normalizedHouseId = houseId.trim().toLowerCase();
    if (!normalizedTokenId && !normalizedHouseId) {
      return;
    }

    const cache = this.readHouseMetadataCache();
    const entry: HouseMetadataCacheEntry = {
      metadata: {
        ...metadata,
        images: Array.isArray(metadata.images)
          ? metadata.images.filter(
              (image) => typeof image === "string" && image.trim().length > 0,
            )
          : [],
      },
      updatedAt: Date.now(),
    };

    if (normalizedTokenId) {
      cache.byTokenId[normalizedTokenId] = entry;
    }
    if (normalizedHouseId) {
      cache.byHouseId[normalizedHouseId] = entry;
    }

    this.writeHouseMetadataCache(cache);
  }

  private isLikelyFallbackMetadata(house: House): boolean {
    const address = String(house.metadata?.address || "").trim().toLowerCase();
    const normalizedHouseId = String(house.houseId || "").trim().toLowerCase();
    if (
      !address ||
      address === normalizedHouseId ||
      address === "unknown property"
    ) {
      return true;
    }

    const city = String(house.metadata?.city || "").trim().toLowerCase();
    const state = String(house.metadata?.state || "").trim().toLowerCase();
    return city === "unknown" && (state === "na" || state.length === 0);
  }

  private mergeCachedMetadataIntoHouse(
    house: House,
    cache: HouseMetadataCache,
  ): House {
    const byTokenId = cache.byTokenId[String(house.tokenId || "")];
    const byHouseId = cache.byHouseId[String(house.houseId || "").trim().toLowerCase()];
    const entry = byTokenId || byHouseId;
    if (!entry) {
      return house;
    }

    if (this.isLikelyFallbackMetadata(house)) {
      return {
        ...house,
        metadata: entry.metadata,
      };
    }

    const currentImages = Array.isArray(house.metadata?.images)
      ? house.metadata.images.filter(
          (image) => typeof image === "string" && image.trim().length > 0,
        )
      : [];

    if (currentImages.length === 0 && entry.metadata.images.length > 0) {
      return {
        ...house,
        metadata: {
          ...house.metadata,
          images: entry.metadata.images,
        },
      };
    }

    return house;
  }

  private mergeCachedMetadataIntoHouses(houses: House[]): House[] {
    if (!Array.isArray(houses) || houses.length === 0) {
      return houses;
    }

    const cache = this.readHouseMetadataCache();
    return houses.map((house) => this.mergeCachedMetadataIntoHouse(house, cache));
  }

  private buildDefaultMetadata(houseId: string) {
    return {
      address: houseId || "Unknown Property",
      city: "Unknown",
      state: "NA",
      zipCode: "00000",
      country: "USA",
      propertyType: "single_family" as const,
      bedrooms: 0,
      bathrooms: 0,
      squareFeet: 0,
      yearBuilt: new Date().getFullYear(),
      description: "",
      images: [],
    };
  }

  private async readTokenFromChain(
    tokenId: bigint,
    contract: ethers.Contract,
  ): Promise<House | null> {
    let ownerAddress: string;
    try {
      ownerAddress = String(await contract.ownerOf(tokenId));
    } catch {
      return null;
    }

    const [houseDetailsRaw, listingRaw, billsRaw, rentalRaw] =
      await Promise.all([
        contract.getHouseDetails(tokenId).catch(() => null),
        contract.getListing(tokenId).catch(() => null),
        contract.getBills(tokenId).catch(() => []),
        contract.getActiveRental(tokenId).catch(() => null),
      ]);

    if (!houseDetailsRaw) return null;

    const houseDetails: any = houseDetailsRaw;
    const listing: any = listingRaw || {};
    const rental: any = rentalRaw || {};
    const bills: any[] = Array.isArray(billsRaw) ? billsRaw : [];

    const houseId = String(
      houseDetails.houseId ?? houseDetails[0] ?? `Token-${tokenId.toString()}`,
    );
    const documentHash = String(
      houseDetails.documentHash ?? houseDetails[1] ?? "0x",
    );
    const documentURI = String(
      houseDetails.documentURI ?? houseDetails[2] ?? "",
    );
    const storageTypeEnum = Number(
      houseDetails.storageType ?? houseDetails[3] ?? 0,
    );
    const originalOwner = String(
      houseDetails.originalOwner ?? houseDetails[4] ?? ownerAddress,
    );
    const mintedAt = this.toDateFromSeconds(
      houseDetails.mintedAt ?? houseDetails[5] ?? 0n,
    );
    const isVerified = Boolean(
      houseDetails.isVerified ?? houseDetails[6] ?? false,
    );

    const listingTypeEnum = Number(listing.listingType ?? listing[0] ?? 0);
    const listingType = this.mapListingType(listingTypeEnum);
    const listingPrice = this.toBigInt(listing.price ?? listing[1] ?? 0n);
    const preferredToken = String(
      listing.preferredToken ?? listing[2] ?? ethers.ZeroAddress,
    );
    const allowedBuyer = String(
      listing.allowedBuyer ?? listing[4] ?? ethers.ZeroAddress,
    );

    const mappedListing =
      listingType === "none"
        ? undefined
        : {
            tokenId: tokenId.toString(),
            listingType,
            price: listingPrice.toString(),
            priceFormatted: `${ethers.formatEther(listingPrice)} ETH`,
            preferredToken,
            isPrivateSale: Boolean(
              listing.isPrivateSale ?? listing[3] ?? false,
            ),
            allowedBuyer: this.isZeroAddress(allowedBuyer)
              ? undefined
              : allowedBuyer,
            createdAt: this.toDateFromSeconds(
              listing.createdAt ?? listing[5] ?? 0n,
            ),
            expiresAt: this.toDateFromSeconds(
              listing.expiresAt ?? listing[6] ?? 0n,
            ),
          };

    const renter = String(rental.renter ?? rental[0] ?? ethers.ZeroAddress);
    const accessKeyHash = String(
      rental.encryptedAccessKeyHash ?? rental[6] ?? ethers.ZeroHash,
    );
    const mappedRental = this.isZeroAddress(renter)
      ? undefined
      : {
          tokenId: tokenId.toString(),
          renterAddress: renter,
          startTime: this.toDateFromSeconds(
            rental.startTime ?? rental[1] ?? 0n,
          ),
          endTime: this.toDateFromSeconds(rental.endTime ?? rental[2] ?? 0n),
          depositAmount: this.toBigInt(
            rental.depositAmount ?? rental[3] ?? 0n,
          ).toString(),
          depositFormatted: `${ethers.formatEther(this.toBigInt(rental.depositAmount ?? rental[3] ?? 0n))} ETH`,
          monthlyRent: this.toBigInt(
            rental.monthlyRent ?? rental[4] ?? 0n,
          ).toString(),
          isActive: Boolean(rental.isActive ?? rental[5] ?? false),
          hasAccessKey: accessKeyHash !== ethers.ZeroHash,
        };

    const mappedBills = bills.map((bill: any, idx: number) => {
      const amountCents = Number(this.toBigInt(bill.amount ?? bill[1] ?? 0n));
      const amountDollars = amountCents / 100;
      const dueDate = this.toDateFromSeconds(bill.dueDate ?? bill[2] ?? 0n);
      const paidAtRaw = this.toBigInt(bill.paidAt ?? bill[3] ?? 0n);
      const paymentReference = String(
        bill.paymentReference ?? bill[5] ?? ethers.ZeroHash,
      );
      const provider = String(bill.provider ?? bill[7] ?? ethers.ZeroAddress);
      const status = Number(bill.status ?? bill[4] ?? 0);

      return {
        id: `${tokenId.toString()}-${idx}`,
        tokenId: tokenId.toString(),
        billType: this.normalizeBillType(
          String(bill.billType ?? bill[0] ?? "other"),
        ),
        amount: amountDollars,
        amountFormatted: `$${amountDollars.toFixed(2)}`,
        dueDate,
        isPaid: status === 1,
        paidAt: paidAtRaw > 0n ? this.toDateFromSeconds(paidAtRaw) : undefined,
        paymentMethod: status === 1 ? ("crypto" as const) : undefined,
        paymentReference:
          paymentReference !== ethers.ZeroHash ? paymentReference : undefined,
        isRecurring: Boolean(bill.isRecurring ?? bill[6] ?? false),
        provider,
        providerName: this.isZeroAddress(provider)
          ? "Unknown Provider"
          : `${provider.slice(0, 6)}...${provider.slice(-4)}`,
        createdAt: dueDate,
      };
    });

    const house: House = {
      tokenId: tokenId.toString(),
      houseId,
      ownerAddress,
      originalOwner,
      documentHash,
      documentURI,
      storageType: this.mapStorageType(storageTypeEnum),
      mintedAt,
      isVerified,
      metadata: this.buildDefaultMetadata(houseId),
      listing: mappedListing,
      rental: mappedRental,
      bills: mappedBills,
    };

    return house;
  }

  private async getHousesFromChain(
    ownerAddress?: string,
  ): Promise<APIResponse<House[]>> {
    try {
      if (!CHAIN_CONFIG.houseRWAAddress) {
        throw new Error(
          "HouseRWA address is not configured. Set VITE_HOUSE_RWA_ADDRESS.",
        );
      }

      const houses = await this.withRpcFallback(async (provider) => {
        const contract = new ethers.Contract(
          CHAIN_CONFIG.houseRWAAddress,
          HOUSE_RWA_READ_ABI,
          provider,
        );
        const nextTokenId = this.toBigInt(await contract.nextTokenId());
        const scanCount =
          nextTokenId > BigInt(CHAIN_CONFIG.maxHouseScan)
            ? BigInt(CHAIN_CONFIG.maxHouseScan)
            : nextTokenId;

        const tokenIds: bigint[] = [];
        for (let i = 0n; i < scanCount; i++) {
          tokenIds.push(i);
        }

        const allHousesRaw: House[] = [];
        for (const tokenId of tokenIds) {
          const token = await this.readTokenFromChain(tokenId, contract);
          if (token) {
            allHousesRaw.push(token);
          }
        }
        let allHouses = allHousesRaw;

        if (ownerAddress) {
          const owner = ownerAddress.toLowerCase();
          allHouses = allHouses.filter(
            (h) => h.ownerAddress.toLowerCase() === owner,
          );
        }
        return allHouses;
      }, "Failed to load houses from onchain RPC providers");

      return {
        success: true,
        message: "Houses loaded from onchain state",
        data: this.mergeCachedMetadataIntoHouses(houses),
      };
    } catch (error: any) {
      return {
        success: false,
        message: error?.message || "Failed to load houses from chain",
        data: [],
      };
    }
  }

  private async getHouseFromChain(
    tokenId: string,
  ): Promise<APIResponse<House>> {
    try {
      if (!CHAIN_CONFIG.houseRWAAddress) {
        throw new Error(
          "HouseRWA address is not configured. Set VITE_HOUSE_RWA_ADDRESS.",
        );
      }
      const tokenBI = BigInt(tokenId);
      const house = await this.withRpcFallback(async (provider) => {
        const contract = new ethers.Contract(
          CHAIN_CONFIG.houseRWAAddress,
          HOUSE_RWA_READ_ABI,
          provider,
        );
        return this.readTokenFromChain(tokenBI, contract);
      }, `Failed to load house ${tokenId} from onchain RPC providers`);
      if (!house) {
        return { success: false, message: "House not found onchain" };
      }
      const mergedHouse = this.mergeCachedMetadataIntoHouses([house])[0];
      return {
        success: true,
        message: "House loaded from onchain state",
        data: mergedHouse,
      };
    } catch (error: any) {
      return {
        success: false,
        message: error?.message || "Failed to load house from chain",
      };
    }
  }

  private async getBillsFromChain(
    tokenId: string,
  ): Promise<APIResponse<House["bills"]>> {
    const houseResp = await this.getHouseFromChain(tokenId);
    if (!houseResp.success || !houseResp.data) {
      return {
        success: false,
        message: houseResp.message || "Failed to load bills from chain",
        data: [],
      };
    }

    return {
      success: true,
      message: "Bills loaded from onchain state",
      data: houseResp.data.bills || [],
    };
  }

  private async getHouseDocumentsFromChain(
    tokenId: string,
  ): Promise<APIResponse<{ documents: any[]; encryptedKey?: string }>> {
    const houseResp = await this.getHouseFromChain(tokenId);
    if (!houseResp.success || !houseResp.data) {
      return {
        success: false,
        message: houseResp.message || "Failed to load documents from chain",
        data: { documents: [] },
      };
    }

    const house = houseResp.data;
    const documents = house.documentURI
      ? [
          {
            name: "Encrypted Property Bundle",
            type: "application/octet-stream",
            size: 0,
            uri: house.documentURI,
            hash: house.documentHash,
          },
        ]
      : [];

    return {
      success: true,
      message: "Document metadata loaded from onchain state",
      data: { documents },
    };
  }

  private withKYCPayload<T extends Record<string, any>>(payload: T): T {
    const action = String(payload?.action || "").toLowerCase();
    if (
      action !== "mint"
      && action !== "set_kyc"
      && action !== "sell"
      && action !== "rent"
      && action !== "create_listing"
    ) {
      return payload;
    }

    const nextPayload: Record<string, any> = { ...payload };

    const defaultProvider = this.kycProviderOverride || getKYCProviderDefault();
    const defaultProofRaw = getKYCProofDefaultRaw();

    if (!nextPayload.kycProvider) {
      nextPayload.kycProvider = defaultProvider;
    }

    const effectiveProvider = parseKYCProvider(
      typeof nextPayload.kycProvider === "string"
        ? nextPayload.kycProvider
        : undefined,
    );
    nextPayload.kycProvider = effectiveProvider;

    if (effectiveProvider === "none") {
      delete nextPayload.kycProof;
      return nextPayload as T;
    }

    if (!nextPayload.kycProof && this.kycProofOverride) {
      nextPayload.kycProof = this.kycProofOverride;
      return nextPayload as T;
    }

    if (!nextPayload.kycProof && defaultProofRaw) {
      try {
        nextPayload.kycProof = JSON.parse(defaultProofRaw);
      } catch {
        // Ignore invalid env value and rely on explicit payload proof.
      }
    }

    return nextPayload as T;
  }

  private async triggerWorkflow<T>(payload: Record<string, any>): Promise<
    APIResponse<T>
  > {
    const path = "/workflow/trigger";
    const endpoint = `${this.baseURL}${path}`;

    try {
      return await this.requestWithLoopbackFallback(
        "POST",
        path,
        payload,
        undefined,
        this.baseURL,
      );
    } catch (error: unknown) {
      if (this.isLikelyNetworkError(error)) {
        const message =
          `Unable to reach workflow endpoint ${endpoint}. Confirm the backend `
          + "is running and that CORS allows this frontend origin.";
        throw new Error(message);
      }

      if (error instanceof Error) {
        throw new Error(`Workflow trigger failed via ${endpoint}: ${error.message}`);
      }

      throw new Error(`Workflow trigger failed via ${endpoint}.`);
    }
  }

  // ==================== AUTH ENDPOINTS ====================

  async verifyWallet(
    address: string,
    signature: string,
    message: string,
  ): Promise<APIResponse<{ token: string; user: User }>> {
    return this.request("POST", "/auth/verify-wallet", {
      address,
      signature,
      message,
    });
  }

  async refreshToken(): Promise<APIResponse<{ token: string }>> {
    return this.request("POST", "/auth/refresh");
  }

  async getCurrentUser(): Promise<APIResponse<User>> {
    return this.request("GET", "/auth/me");
  }

  async logout(): Promise<APIResponse<void>> {
    return this.request("POST", "/auth/logout");
  }

  // ==================== HOUSE ENDPOINTS ====================

  async mintHouse(payload: MintRequestPayload): Promise<
    APIResponse<{
      tokenId: string;
      txHash: string;
      encryptedKey: string;
      documentHash: string;
    }>
  > {
    const response = await this.triggerWorkflow<{
      tokenId: string;
      txHash: string;
      encryptedKey: string;
      documentHash: string;
      tokenID?: string;
    }>(this.withKYCPayload(payload));

    if (response.success) {
      const tokenId = String(
        response.data?.tokenId || response.data?.tokenID || "",
      ).trim();
      this.cacheMintedHouseMetadata(
        tokenId || null,
        payload.houseId,
        payload.metadata,
      );
    }

    return response;
  }

  async updateHouseImages(
    tokenId: string,
    images: string[],
  ): Promise<
    APIResponse<{
      tokenId: string;
      images: string[];
      updatedAt: string;
    }>
  > {
    const normalizedImages = Array.from(
      new Set(
        (Array.isArray(images) ? images : [])
          .filter((value) => typeof value === "string")
          .map((value) => value.trim())
          .filter((value) => value.length > 0),
      ),
    ).slice(0, 10);

    if (normalizedImages.length === 0) {
      return {
        success: false,
        message: "Provide at least one valid image to update.",
      };
    }

    return this.triggerWorkflow({
      action: "update_house_images",
      tokenId,
      metadata: { images: normalizedImages },
    });
  }

  async ensureKYC(
    walletAddress: string,
  ): Promise<APIResponse<{ walletAddress: string }>> {
    return this.triggerWorkflow(
      this.withKYCPayload({
        action: "set_kyc",
        walletAddress,
      }),
    );
  }

  async sellHouse(payload: SellRequestPayload): Promise<
    APIResponse<{
      txHash: string;
      keyHash: string;
    }>
  > {
    return this.triggerWorkflow(this.withKYCPayload(payload));
  }

  async rentHouse(payload: RentRequestPayload): Promise<
    APIResponse<{
      txHash: string;
      accessKeyHash: string;
    }>
  > {
    return this.triggerWorkflow(this.withKYCPayload(payload));
  }

  async getHouses(ownerAddress?: string): Promise<APIResponse<House[]>> {
    const apiConfigured = this.isApiConfigured();
    const buildQuery = (includeOwner: boolean): string => {
      const params = new URLSearchParams();
      if (includeOwner && ownerAddress) {
        params.set("owner", ownerAddress);
      }
      if (CHAIN_CONFIG.houseRWAAddress) {
        params.set("contractAddress", CHAIN_CONFIG.houseRWAAddress);
      }
      return params.toString() ? `?${params.toString()}` : "";
    };

    const ownerScopedQuery = buildQuery(Boolean(ownerAddress) && this.hasAuthToken());
    const publicQuery = buildQuery(false);
    const isOwnerScopedAuthError = (message: string): boolean => {
      return message
        .toLowerCase()
        .includes("authentication is required to query owner-specific private houses");
    };

    const tryPublicOwnerFallback = async (): Promise<APIResponse<House[]> | null> => {
      if (!ownerAddress) {
        return null;
      }

      const publicResp = await this.requestPublic<House[]>(
        `/houses${publicQuery}`,
        this.baseURL,
        { includeAuthHeader: true },
      );
      if (!publicResp.success) {
        return null;
      }

      const normalizedOwner = ownerAddress.toLowerCase();
      const filtered = (publicResp.data || []).filter((house) => {
        const owner = String(house.ownerAddress || "").toLowerCase();
        return owner === normalizedOwner;
      });

      return {
        success: true,
        message:
          "Loaded public owner houses. Sign in to include owner-private records.",
        data: this.mergeCachedMetadataIntoHouses(filtered),
      };
    };

    if (ownerAddress && !this.hasAuthToken()) {
      try {
        const fallback = await tryPublicOwnerFallback();
        if (fallback) {
          return fallback;
        }
      } catch {
        // Continue to standard error handling.
      }

      if (apiConfigured) {
        return {
          success: false,
          message:
            "Authentication is required to query owner-specific private houses.",
          data: [],
        };
      }
    }

    try {
      const resp = await this.requestPublic<House[]>(
        `/houses${ownerScopedQuery}`,
        this.baseURL,
        { includeAuthHeader: true },
      );
      if (resp.success) {
        return {
          ...resp,
          data: this.mergeCachedMetadataIntoHouses(resp.data || []),
        };
      }

      if (ownerAddress && isOwnerScopedAuthError(resp.message || "")) {
        try {
          const fallback = await tryPublicOwnerFallback();
          if (fallback) {
            return fallback;
          }
        } catch {
          // Continue to standard error handling.
        }
      }

      if (apiConfigured) {
        return resp;
      }
    } catch (error) {
      const message = error instanceof Error ? error.message : "";
      if (ownerAddress && isOwnerScopedAuthError(message)) {
        try {
          const fallback = await tryPublicOwnerFallback();
          if (fallback) {
            return fallback;
          }
        } catch {
          // Continue to standard error handling.
        }
      }

      if (apiConfigured) {
        return {
          success: false,
          message: `Failed to load houses from API ${this.baseURL}`,
          data: [],
        };
      }
    }
    if (!this.allowChainRpcFallback()) {
      return {
        success: false,
        message:
          "Failed to load houses from API and direct RPC fallback is disabled "
          + "for privacy. Check API/CSP settings or set VITE_ENABLE_CHAIN_RPC_FALLBACK=true.",
        data: [],
      };
    }
    return this.getHousesFromChain(ownerAddress);
  }

  async getHouse(tokenId: string): Promise<APIResponse<House>> {
    const apiConfigured = this.isApiConfigured();
    const params = new URLSearchParams();
    if (CHAIN_CONFIG.houseRWAAddress) {
      params.set("contractAddress", CHAIN_CONFIG.houseRWAAddress);
    }
    const query = params.toString() ? `?${params.toString()}` : "";
    try {
      const resp = await this.requestPublic<House>(
        `/houses/${tokenId}${query}`,
        this.baseURL,
        { includeAuthHeader: true },
      );
      if (resp.success) {
        const mergedHouse = resp.data
          ? this.mergeCachedMetadataIntoHouses([resp.data])[0]
          : undefined;
        return {
          ...resp,
          data: mergedHouse,
        };
      }
      if (apiConfigured) {
        return resp;
      }
    } catch {
      if (apiConfigured) {
        return {
          success: false,
          message: `Failed to load house from API ${this.baseURL}`,
        };
      }
    }
    if (!this.allowChainRpcFallback()) {
      return {
        success: false,
        message:
          "Failed to load house from API and direct RPC fallback is disabled "
          + "for privacy. Check API/CSP settings or set VITE_ENABLE_CHAIN_RPC_FALLBACK=true.",
      };
    }
    return this.getHouseFromChain(tokenId);
  }

  async getHouseDocuments(tokenId: string): Promise<
    APIResponse<{
      documents: any[];
      encryptedKey?: string;
    }>
  > {
    const apiConfigured = this.isApiConfigured();
    const params = new URLSearchParams();
    if (CHAIN_CONFIG.houseRWAAddress) {
      params.set("contractAddress", CHAIN_CONFIG.houseRWAAddress);
    }
    const query = params.toString() ? `?${params.toString()}` : "";
    try {
      const resp = await this.requestPublic<{
        documents: any[];
        encryptedKey?: string;
      }>(`/houses/${tokenId}/documents${query}`, this.baseURL, {
        includeAuthHeader: true,
      });
      if (resp.success) return resp;
      if (apiConfigured) {
        return resp;
      }
    } catch {
      if (apiConfigured) {
        return {
          success: false,
          message: `Failed to load documents from API ${this.baseURL}`,
          data: { documents: [] },
        };
      }
    }
    if (!this.allowChainRpcFallback()) {
      return {
        success: false,
        message:
          "Failed to load documents from API and direct RPC fallback is disabled "
          + "for privacy. Check API/CSP settings or set VITE_ENABLE_CHAIN_RPC_FALLBACK=true.",
        data: { documents: [] },
      };
    }
    return this.getHouseDocumentsFromChain(tokenId);
  }

  async getHouseDocumentContents(tokenId: string): Promise<
    APIResponse<{
      documents: Array<{
        index: number;
        name: string;
        mimeType: string;
        size: number;
        base64: string;
      }>;
      documentHash?: string;
      documentURI?: string;
    }>
  > {
    try {
      const response = await this.requestPublic<{
        documents: Array<{
          index: number;
          name: string;
          mimeType: string;
          size: number;
          base64: string;
        }>;
        documentHash?: string;
        documentURI?: string;
      }>(`/houses/${tokenId}/documents/content`, this.baseURL, {
        includeAuthHeader: true,
      });
      return response;
    } catch (error: any) {
      return {
        success: false,
        message: error?.message || "Failed to load private document contents",
        data: { documents: [] },
      };
    }
  }

  // ==================== BILL ENDPOINTS ====================

  async createBill(payload: CreateBillData): Promise<
    APIResponse<{
      txHash: string;
      billIndex: number;
    }>
  > {
    return this.triggerWorkflow({
      action: "create_bill",
      ...payload,
    });
  }

  async payBill(payload: BillPaymentData): Promise<
    APIResponse<{
      txHash: string;
      paymentReference: string;
    }>
  > {
    return this.triggerWorkflow({
      action: "pay_bill",
      ...payload,
    });
  }

  async getBills(tokenId: string): Promise<APIResponse<House["bills"]>> {
    const apiConfigured = this.isApiConfigured();
    const params = new URLSearchParams();
    if (CHAIN_CONFIG.houseRWAAddress) {
      params.set("contractAddress", CHAIN_CONFIG.houseRWAAddress);
    }
    const query = params.toString() ? `?${params.toString()}` : "";
    try {
      const resp = await this.requestPublic<House["bills"]>(
        `/houses/${tokenId}/bills${query}`,
        this.baseURL,
        { includeAuthHeader: true },
      );
      if (resp.success) return resp;
      if (apiConfigured) {
        return resp;
      }
    } catch {
      if (apiConfigured) {
        return {
          success: false,
          message: `Failed to load bills from API ${this.baseURL}`,
          data: [],
        };
      }
    }
    if (!this.allowChainRpcFallback()) {
      return {
        success: false,
        message:
          "Failed to load bills from API and direct RPC fallback is disabled "
          + "for privacy. Check API/CSP settings or set VITE_ENABLE_CHAIN_RPC_FALLBACK=true.",
        data: [],
      };
    }
    return this.getBillsFromChain(tokenId);
  }

  // ==================== LISTING ENDPOINTS ====================

  async createListing(
    payload: CreateListingRequestPayload,
  ): Promise<APIResponse<{ txHash: string }>> {
    return this.triggerWorkflow(this.withKYCPayload(payload));
  }

  async cancelListing(
    tokenId: string,
  ): Promise<APIResponse<{ txHash: string }>> {
    return this.request("DELETE", `/houses/${tokenId}/listings`);
  }

  // ==================== TRANSACTION ENDPOINTS ====================

  async getTransactions(address?: string): Promise<APIResponse<Transaction[]>> {
    const query = address ? `?address=${address}` : "";
    return this.request("GET", `/transactions${query}`);
  }

  async getTransaction(hash: string): Promise<APIResponse<Transaction>> {
    return this.request("GET", `/transactions/${hash}`);
  }

  // ==================== NOTIFICATION ENDPOINTS ====================

  async getNotifications(): Promise<APIResponse<Notification[]>> {
    return this.request("GET", "/notifications");
  }

  async getNativeBalance(address: string): Promise<APIResponse<string>> {
    const apiConfigured = this.isApiConfigured();
    try {
      if (!/^0x[a-fA-F0-9]{40}$/.test(address)) {
        return {
          success: false,
          message: "Invalid wallet address.",
        };
      }

      try {
        const apiResponse = await this.requestPublic<string>(
          `/balances/${address}`,
          this.baseURL,
          { includeAuthHeader: true },
        );
        if (apiResponse.success && apiResponse.data) {
          return apiResponse;
        }
        if (apiConfigured) {
          return apiResponse;
        }
      } catch {
        if (apiConfigured) {
          return {
            success: false,
            message: `Failed to load balance from API ${this.baseURL}`,
          };
        }
      }
      if (!this.allowChainRpcFallback()) {
        return {
          success: false,
          message:
            "Failed to load balance from API and direct RPC fallback is disabled "
            + "for privacy. Check API/CSP settings or set VITE_ENABLE_CHAIN_RPC_FALLBACK=true.",
        };
      }
      const balance = await this.withRpcFallback(async (provider) => {
        return provider.getBalance(address);
      }, `Failed to read native balance for ${address}`);

      return {
        success: true,
        message: "Balance loaded from onchain RPC.",
        data: ethers.formatEther(balance),
      };
    } catch (error: any) {
      return {
        success: false,
        message: error?.message || "Failed to load native balance.",
      };
    }
  }

  async markNotificationRead(
    notificationId: string,
  ): Promise<APIResponse<void>> {
    return this.request("POST", `/notifications/${notificationId}/read`);
  }

  // ==================== XMTP MESSAGING ENDPOINTS ====================

  async getConversations(
    tokenId?: string,
  ): Promise<APIResponse<ConversationSummary[]>> {
    const query = tokenId ? `?tokenId=${encodeURIComponent(tokenId)}` : "";
    return this.request("GET", `/messages/conversations${query}`);
  }

  async getConversation(
    conversationId: string,
  ): Promise<APIResponse<ConversationDetails>> {
    const encoded = encodeURIComponent(conversationId);
    return this.request("GET", `/messages/conversations/${encoded}`);
  }

  async sendPrivateMessage(payload: {
    tokenId: string;
    recipientWalletAddress: string;
    message: string;
    xmtpMessageId?: string;
  }): Promise<
    APIResponse<{
      conversation: ConversationSummary;
      message: ConversationDetails["messages"][number];
    }>
  > {
    return this.request("POST", "/messages/send", payload);
  }

  // ==================== KYC ENDPOINTS ====================

  async getZKPassportHealth(): Promise<
    APIResponse<{
      status: string;
      zkpassport?: {
        domain?: string;
        scope?: string;
        devMode?: boolean;
      };
    }>
  > {
    return this.requestWithLoopbackFallback(
      "GET",
      "/healthz",
      undefined,
      undefined,
      this.getZKPassportBaseURL(),
      { credentialsMode: "omit" },
    );
  }

  async startZKPassportSession(
    walletAddress: string,
    options?: {
      domain?: string;
      mode?: "fast" | "compressed" | "compressed-evm";
    },
  ): Promise<APIResponse<ZKPassportSession>> {
    return this.requestWithLoopbackFallback(
      "POST",
      "/kyc/zkpassport/session",
      {
        walletAddress,
        domain: options?.domain,
        mode: options?.mode,
      },
      undefined,
      this.getZKPassportBaseURL(),
      { credentialsMode: "omit" },
    );
  }

  async getZKPassportSession(
    sessionId: string,
  ): Promise<APIResponse<ZKPassportSession>> {
    const encoded = encodeURIComponent(sessionId);
    const cacheBust = `ts=${Date.now()}`;
    return this.requestWithLoopbackFallback(
      "GET",
      `/kyc/zkpassport/session/${encoded}?${cacheBust}`,
      undefined,
      undefined,
      this.getZKPassportBaseURL(),
      { credentialsMode: "omit" },
    );
  }

  async verifyZKPassportProof(
    walletAddress: string,
    proof: ZKPassportProofBundle,
  ): Promise<APIResponse<ZKPassportVerificationResult>> {
    return this.requestWithLoopbackFallback(
      "POST",
      "/kyc/zkpassport/verify",
      {
        provider: "zkpassport",
        walletAddress,
        proof,
      },
      undefined,
      this.getZKPassportBaseURL(),
      { credentialsMode: "omit" },
    );
  }

  // ==================== KEY MANAGEMENT ====================

  async claimKey(
    keyHash: string,
    claimant?: string,
  ): Promise<
    APIResponse<{
      encryptedKey: string;
    }>
  > {
    const body: Record<string, any> = {
      action: "claim_key",
      keyHash,
    };

    // Optional claimant field for workflows that validate the caller address.
    // NOTE: In the onchain `claimKey(bytes32)` path, the EVM sender is the claimant.
    // This field is used only for HTTP-trigger simulation/backends.
    if (claimant) {
      body.claimant = claimant;
    }

    return this.triggerWorkflow(body);
  }

  // ==================== STRIPE PAYMENT ENDPOINTS ====================

  async createStripePaymentIntent(
    amount: number,
    currency: string = "usd",
    metadata?: Record<string, string>,
  ): Promise<
    APIResponse<{
      clientSecret: string;
      paymentIntentId: string;
    }>
  > {
    return this.request("POST", "/payments/stripe/create-intent", {
      amount,
      currency,
      metadata,
    });
  }

  async confirmStripePayment(
    paymentIntentId: string,
    tokenId: string,
    billIndex: number,
  ): Promise<
    APIResponse<{
      txHash: string;
      paymentReference: string;
    }>
  > {
    return this.request("POST", "/payments/stripe/confirm", {
      paymentIntentId,
      tokenId,
      billIndex,
    });
  }

  // ==================== UPLOAD ENDPOINTS ====================

  async uploadDocuments(
    tokenId: string,
    files: File[],
  ): Promise<APIResponse<{ documentURIs: string[] }>> {
    const formData = new FormData();
    files.forEach((file, index) => {
      formData.append(`file${index}`, file);
    });

    // Use fetch directly for FormData
    const response = await fetch(
      `${this.baseURL}/houses/${tokenId}/documents`,
      {
        method: "POST",
        headers: {
          Authorization: this.authToken ? `Bearer ${this.authToken}` : "",
          "X-Request-ID": this.generateRequestId(),
        },
        body: formData,
        credentials: "include",
      },
    );

    if (!response.ok) {
      throw new Error(`Upload failed: ${response.status}`);
    }

    return response.json();
  }
}

// Singleton instance
export const apiClient = new RWAApiClient();

// Export configuration for customization
export { API_CONFIG };
