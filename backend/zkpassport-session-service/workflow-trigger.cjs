const fs = require("node:fs");
const path = require("node:path");
const {
  createCipheriv,
  createDecipheriv,
  createHmac,
  createHash,
  randomBytes,
  randomUUID,
  timingSafeEqual,
} = require("node:crypto");
const { ZKPassport } = require("@zkpassport/sdk");

bootstrapEnv();

const WORKFLOW_CHAIN_ID = parsePositiveInteger(
  process.env.WORKFLOW_CHAIN_ID ?? process.env.CHAIN_ID,
  11155111,
);
const WORKFLOW_RPC_URL =
  process.env.WORKFLOW_RPC_URL ??
  process.env.CRE_RPC_URL ??
  process.env.SEPOLIA_RPC ??
  "";
const WORKFLOW_PRIVATE_KEY =
  process.env.WORKFLOW_PRIVATE_KEY ??
  process.env.CRE_ETH_PRIVATE_KEY ??
  process.env.PRIVATE_KEY ??
  "";
const WORKFLOW_CONTRACT_ADDRESS =
  process.env.HOUSE_RWA_CONTRACT_ADDRESS ??
  process.env.CRE_CONTRACT_ADDRESS ??
  process.env.VITE_HOUSE_RWA_ADDRESS ??
  "";
const WORKFLOW_CONFIRMATIONS = parsePositiveInteger(
  process.env.WORKFLOW_CONFIRMATIONS,
  1,
);
const WORKFLOW_TIMEOUT_MS = parsePositiveInteger(
  process.env.WORKFLOW_TIMEOUT_MS,
  120000,
);
const WORKFLOW_MAX_BODY_BYTES = parsePositiveInteger(
  process.env.WORKFLOW_MAX_BODY_BYTES,
  120 * 1024 * 1024,
);
const KYC_LEVEL = parsePositiveInteger(process.env.WORKFLOW_KYC_LEVEL, 2);
const KYC_EXPIRY_DAYS = parsePositiveInteger(
  process.env.WORKFLOW_KYC_EXPIRY_DAYS,
  180,
);
const ZKPASSPORT_DOMAIN = process.env.ZKPASSPORT_DOMAIN ?? "demo.zkpassport.id";
const ZKPASSPORT_SCOPE = process.env.ZKPASSPORT_SCOPE ?? "rwa-house-kyc";
const ZKPASSPORT_DEV_MODE = parseBoolean(process.env.ZKPASSPORT_DEV_MODE, false);
const ZKPASSPORT_VALIDITY_SECONDS = parsePositiveInteger(
  process.env.ZKPASSPORT_VALIDITY_SECONDS,
  900,
);
const ZKPASSPORT_VERIFY_WRITING_DIRECTORY =
  process.env.ZKPASSPORT_VERIFY_WRITING_DIRECTORY ?? "/tmp";
const WORKFLOW_PRIVATE_STORE_PATH =
  process.env.WORKFLOW_PRIVATE_STORE_PATH
  ?? path.join(__dirname, ".workflow-private-store.json");
const WORKFLOW_PRIVATE_STORE_VERSION = 1;
const WORKFLOW_PRIVATE_STORE_PURPOSE = "cre-house-private-metadata";
const WORKFLOW_ACTIVITY_STORE_PATH =
  process.env.WORKFLOW_ACTIVITY_STORE_PATH
  ?? path.join(__dirname, ".workflow-activity-store.json");
const WORKFLOW_ACTIVITY_STORE_VERSION = 1;
const WORKFLOW_ACTIVITY_MAX_NOTIFICATIONS = parsePositiveInteger(
  process.env.WORKFLOW_ACTIVITY_MAX_NOTIFICATIONS,
  200,
);
const WORKFLOW_ACTIVITY_MAX_MESSAGES = parsePositiveInteger(
  process.env.WORKFLOW_ACTIVITY_MAX_MESSAGES,
  500,
);
const WORKFLOW_ACTIVITY_MAX_CONVERSATIONS = parsePositiveInteger(
  process.env.WORKFLOW_ACTIVITY_MAX_CONVERSATIONS,
  200,
);
const WORKFLOW_AUTH_TOKEN_TTL_SECONDS = parsePositiveInteger(
  process.env.WORKFLOW_AUTH_TOKEN_TTL_SECONDS,
  12 * 60 * 60,
);
const WORKFLOW_ALLOW_INSECURE_BEARER = parseBoolean(
  process.env.WORKFLOW_ALLOW_INSECURE_BEARER,
  false,
);
const WORKFLOW_AUTH_SECRET =
  process.env.WORKFLOW_AUTH_SECRET ?? process.env.AUTH_SECRET ?? "";
const WORKFLOW_AUTH_AUDIENCE = "workflow-trigger-api";
const AUTH_REQUIRED_WORKFLOW_ACTIONS = new Set([
  "mint",
  "set_kyc",
  "create_listing",
  "sell",
  "rent",
  "create_bill",
  "pay_bill",
  "update_house_images",
  "claim_key",
]);

if (!String(WORKFLOW_AUTH_SECRET).trim()) {
  throw new Error(
    "WORKFLOW_AUTH_SECRET is required. Refusing to start without a configured auth secret.",
  );
}

const HOUSE_RWA_ABI = [
  "function setKYCVerification(address user,uint8 level,bytes32 verificationHash,uint48 expiryDate)",
  "function hasKYC(address user) view returns (bool)",
  "function mint(address to,string houseId,bytes32 documentHash,string documentURI,uint8 storageType,string verificationData) returns (uint256)",
  "function createListingFromWorkflow(uint256 tokenId,address owner,uint8 listingType,uint96 price,address preferredToken,bool isPrivateSale,address allowedBuyer,uint48 durationDays)",
  "function completeSale(uint256 tokenId,address buyer,bytes32 keyHash,bytes encryptedKey)",
  "function startRental(uint256 tokenId,address renter,uint48 durationDays,uint96 depositAmount,uint96 monthlyRent,bytes encryptedAccessKey)",
  "function createBill(uint256 tokenId,string billType,uint96 amount,uint48 dueDate,address provider,bool isRecurring,uint8 recurrenceInterval) returns (uint256 billIndex)",
  "function recordBillPayment(uint256 tokenId,uint256 billIndex,string paymentMethod,bytes32 paymentReference)",
  "function balanceOf(address owner) view returns (uint256)",
  "function nextTokenId() view returns (uint256)",
  "function ownerOf(uint256 tokenId) view returns (address)",
  "function getHouseDetails(uint256 tokenId) view returns ((string houseId,bytes32 documentHash,string documentURI,uint8 storageType,address originalOwner,uint48 mintedAt,bool isVerified,uint8 documentCount))",
  "function getListing(uint256 tokenId) view returns ((uint8 listingType,uint96 price,address preferredToken,bool isPrivateSale,address allowedBuyer,uint48 createdAt,uint48 expiresAt,uint8 platformFee))",
  "function getBills(uint256 tokenId) view returns ((string billType,uint96 amount,uint48 dueDate,uint48 paidAt,uint8 status,bytes32 paymentReference,bool isRecurring,address provider,uint8 recurrenceInterval)[])",
  "function getActiveRental(uint256 tokenId) view returns ((address renter,uint48 startTime,uint48 endTime,uint96 depositAmount,uint96 monthlyRent,bool isActive,bytes32 encryptedAccessKeyHash,uint8 disputeStatus))",
  "function getTotalBillsCount(uint256 tokenId) view returns (uint256)",
  "function keyExchanges(bytes32) view returns (bytes32 keyHash, bytes encryptedKey, address intendedRecipient, uint48 createdAt, uint48 expiresAt, bool isClaimed, uint8 exchangeType)",
];

const PAYMENT_METHODS = new Set(["crypto", "stripe", "bank_transfer"]);

let cachedRuntimePromise = null;

function bootstrapEnv() {
  const candidates = [
    path.join(__dirname, ".env"),
    path.join(__dirname, "../../.env"),
  ];

  for (const candidate of candidates) {
    loadEnvFile(candidate);
  }
}

function loadEnvFile(filePath) {
  if (!fs.existsSync(filePath)) {
    return;
  }

  const content = fs.readFileSync(filePath, "utf8");
  const lines = content.split(/\r?\n/);

  for (const line of lines) {
    const trimmed = line.trim();
    if (!trimmed || trimmed.startsWith("#")) {
      continue;
    }
    const equalsIndex = trimmed.indexOf("=");
    if (equalsIndex <= 0) {
      continue;
    }

    const key = trimmed.slice(0, equalsIndex).trim();
    if (!key || process.env[key] !== undefined) {
      continue;
    }

    let value = trimmed.slice(equalsIndex + 1).trim();
    if (
      (value.startsWith('"') && value.endsWith('"')) ||
      (value.startsWith("'") && value.endsWith("'"))
    ) {
      value = value.slice(1, -1);
    }
    process.env[key] = value;
  }
}

function parsePositiveInteger(value, fallback) {
  const parsed = Number.parseInt(String(value ?? ""), 10);
  if (!Number.isFinite(parsed) || parsed <= 0) {
    return fallback;
  }
  return parsed;
}

function parseBoolean(value, fallback) {
  if (value === undefined || value === null || value === "") {
    return fallback;
  }

  const normalized = String(value).trim().toLowerCase();
  if (normalized === "true" || normalized === "1" || normalized === "yes") {
    return true;
  }
  if (normalized === "false" || normalized === "0" || normalized === "no") {
    return false;
  }

  return fallback;
}

let cachedAuthTokenSecret = null;

function encodeBase64Url(buffer) {
  return buffer
    .toString("base64")
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=+$/g, "");
}

function decodeBase64Url(value) {
  const normalized = String(value ?? "").replace(/-/g, "+").replace(/_/g, "/");
  const padding = normalized.length % 4;
  const padded = padding === 0 ? normalized : `${normalized}${"=".repeat(4 - padding)}`;
  return Buffer.from(padded, "base64");
}

function resolveAuthTokenSecret() {
  if (cachedAuthTokenSecret) {
    return cachedAuthTokenSecret;
  }

  const sourceMaterial = String(WORKFLOW_AUTH_SECRET).trim();
  if (!sourceMaterial) {
    throw new Error("WORKFLOW_AUTH_SECRET is required for bearer token signing.");
  }
  cachedAuthTokenSecret = createHash("sha256")
    .update(`auth-token:${sourceMaterial}`)
    .digest();
  return cachedAuthTokenSecret;
}

function signAuthTokenPayload(payloadSegment) {
  return createHmac("sha256", resolveAuthTokenSecret())
    .update(payloadSegment, "utf8")
    .digest();
}

function createAuthToken(walletAddress) {
  const normalizedAddress = normalizeWalletAddress(walletAddress);
  const issuedAt = Math.floor(Date.now() / 1000);
  const payload = {
    sub: normalizedAddress,
    aud: WORKFLOW_AUTH_AUDIENCE,
    iat: issuedAt,
    exp: issuedAt + WORKFLOW_AUTH_TOKEN_TTL_SECONDS,
  };
  const payloadSegment = encodeBase64Url(
    Buffer.from(JSON.stringify(payload), "utf8"),
  );
  const signatureSegment = encodeBase64Url(signAuthTokenPayload(payloadSegment));
  return `rwa.${payloadSegment}.${signatureSegment}`;
}

function parseSignedAuthToken(token) {
  const rawToken = String(token ?? "").trim();
  if (!rawToken.startsWith("rwa.")) {
    return "";
  }

  const parts = rawToken.split(".");
  if (parts.length !== 3 || parts[0] !== "rwa") {
    return "";
  }
  const payloadSegment = parts[1];
  const providedSignature = parts[2];
  if (!payloadSegment || !providedSignature) {
    return "";
  }

  let payloadBuffer;
  let providedSignatureBuffer;
  try {
    payloadBuffer = decodeBase64Url(payloadSegment);
    providedSignatureBuffer = decodeBase64Url(providedSignature);
  } catch {
    return "";
  }

  const expectedSignatureBuffer = signAuthTokenPayload(payloadSegment);
  if (
    expectedSignatureBuffer.length !== providedSignatureBuffer.length
    || !timingSafeEqual(expectedSignatureBuffer, providedSignatureBuffer)
  ) {
    return "";
  }

  try {
    const payload = JSON.parse(payloadBuffer.toString("utf8"));
    const subject = normalizeWalletAddress(payload?.sub);
    const audience = String(payload?.aud ?? "");
    const expiresAt = Number.parseInt(String(payload?.exp ?? ""), 10);
    const now = Math.floor(Date.now() / 1000);

    if (!isHexAddress(subject)) {
      return "";
    }
    if (audience !== WORKFLOW_AUTH_AUDIENCE) {
      return "";
    }
    if (!Number.isFinite(expiresAt) || expiresAt <= now) {
      return "";
    }
    return subject;
  } catch {
    return "";
  }
}

function mapListingType(value) {
  if (value === 1) return "for_sale";
  if (value === 2) return "for_rent";
  return "none";
}

function mapStorageType(value) {
  if (value === 1) return "offchain";
  return "ipfs";
}

function toBigInt(value) {
  if (typeof value === "bigint") return value;
  if (typeof value === "number") return BigInt(Math.max(0, Math.trunc(value)));
  if (typeof value === "string" && value !== "") return BigInt(value);
  if (value && typeof value === "object") {
    if (typeof value.toBigInt === "function") {
      return value.toBigInt();
    }
    if (typeof value.toString === "function") {
      const stringValue = value.toString();
      if (/^-?\d+$/.test(stringValue)) {
        return BigInt(stringValue);
      }
    }
  }
  return 0n;
}

function toDate(value) {
  const seconds = toBigInt(value);
  if (seconds <= 0n) return new Date(0);
  const ms = seconds * 1000n;
  if (ms > BigInt(Number.MAX_SAFE_INTEGER)) {
    return new Date(Number.MAX_SAFE_INTEGER);
  }
  return new Date(Number(ms));
}

function normalizeBillType(value) {
  const v = String(value || "other").toLowerCase();
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

function isZeroAddress(address) {
  return (
    !address
    || String(address).toLowerCase() === "0x0000000000000000000000000000000000000000"
  );
}

function buildDefaultMetadata(houseId) {
  return {
    address: houseId || "Unknown Property",
    city: "Unknown",
    state: "NA",
    zipCode: "00000",
    country: "USA",
    propertyType: "single_family",
    bedrooms: 0,
    bathrooms: 0,
    squareFeet: 0,
    yearBuilt: new Date().getFullYear(),
    description: "",
    images: [],
  };
}

const PROPERTY_TYPES = new Set([
  "single_family",
  "condo",
  "townhouse",
  "multi_family",
  "apartment",
  "commercial",
]);

function normalizeHouseMetadata(rawMetadata, houseId) {
  const defaults = buildDefaultMetadata(houseId);
  const source =
    rawMetadata && typeof rawMetadata === "object" ? rawMetadata : {};
  const metadata = { ...defaults };

  const textFields = [
    "address",
    "city",
    "state",
    "zipCode",
    "country",
    "description",
  ];
  for (const field of textFields) {
    const value = source[field];
    if (typeof value === "string" && value.trim()) {
      metadata[field] = value.trim();
    }
  }

  const propertyType = String(source.propertyType ?? "").trim().toLowerCase();
  if (PROPERTY_TYPES.has(propertyType)) {
    metadata.propertyType = propertyType;
  }

  const numberFields = ["bedrooms", "bathrooms", "squareFeet", "yearBuilt"];
  for (const field of numberFields) {
    const parsed = Number(source[field]);
    if (Number.isFinite(parsed) && parsed >= 0) {
      metadata[field] = Math.floor(parsed);
    }
  }

  const images = Array.isArray(source.images)
    ? source.images
        .filter((value) => typeof value === "string")
        .map((value) => value.trim())
        .filter(Boolean)
    : [];
  metadata.images = images;

  return metadata;
}

function isAllowedMetadataImage(value) {
  const trimmed = String(value ?? "").trim();
  if (!trimmed) {
    return false;
  }

  const lowered = trimmed.toLowerCase();
  if (lowered.startsWith("ipfs://")) {
    return true;
  }
  if (lowered.startsWith("http://") || lowered.startsWith("https://")) {
    return true;
  }
  if (lowered.startsWith("data:image/png;base64,")) {
    const payload = trimmed.slice("data:image/png;base64,".length).trim();
    return /^[a-z0-9+/=\s]+$/i.test(payload);
  }
  return false;
}

function parseMetadataImages(rawImages) {
  if (!Array.isArray(rawImages)) {
    throw new Error("metadata.images must be an array.");
  }

  const unique = new Set();
  for (const candidate of rawImages) {
    if (typeof candidate !== "string") {
      continue;
    }
    const trimmed = candidate.trim();
    if (!trimmed) {
      continue;
    }
    if (!isAllowedMetadataImage(trimmed)) {
      throw new Error(
        "metadata.images entries must use https://, http://, ipfs://, or data:image/png;base64,...",
      );
    }
    unique.add(trimmed);
  }

  return Array.from(unique).slice(0, 10);
}

function sortObjectKeys(value) {
  if (Array.isArray(value)) {
    return value.map((entry) => sortObjectKeys(entry));
  }
  if (!value || typeof value !== "object") {
    return value;
  }
  const result = {};
  for (const key of Object.keys(value).sort()) {
    result[key] = sortObjectKeys(value[key]);
  }
  return result;
}

function stableStringify(value) {
  return JSON.stringify(sortObjectKeys(value));
}

function isHexCommitment(value) {
  return /^0x[a-fA-F0-9]{64}$/.test(String(value ?? "").trim());
}

function ensureDirectoryForFile(filePath) {
  const directory = path.dirname(filePath);
  if (!fs.existsSync(directory)) {
    fs.mkdirSync(directory, { recursive: true });
  }
}

function resolvePrivateStoreKey() {
  const configured = String(process.env.WORKFLOW_PRIVATE_STORE_KEY ?? "").trim();
  if (configured) {
    const normalized = configured.startsWith("0x") ? configured.slice(2) : configured;
    if (/^[a-fA-F0-9]{64}$/.test(normalized)) {
      return Buffer.from(normalized, "hex");
    }
    try {
      const decoded = Buffer.from(configured, "base64");
      if (decoded.length >= 32) {
        return decoded.subarray(0, 32);
      }
    } catch {
      // Fall through to deterministic fallback key material.
    }
  }

  const material = `${WORKFLOW_PRIVATE_KEY}:${WORKFLOW_CONTRACT_ADDRESS}:${WORKFLOW_PRIVATE_STORE_PURPOSE}`;
  return createHash("sha256").update(material).digest();
}

function encryptPrivatePayload(payload) {
  const key = resolvePrivateStoreKey();
  const iv = randomBytes(12);
  const cipher = createCipheriv("aes-256-gcm", key, iv);
  const data = Buffer.concat([
    cipher.update(Buffer.from(JSON.stringify(payload), "utf8")),
    cipher.final(),
  ]);
  const tag = cipher.getAuthTag();

  return {
    iv: iv.toString("base64"),
    tag: tag.toString("base64"),
    data: data.toString("base64"),
  };
}

function decryptPrivatePayload(payload) {
  if (!payload || typeof payload !== "object") {
    return null;
  }
  const { iv, tag, data } = payload;
  if (typeof iv !== "string" || typeof tag !== "string" || typeof data !== "string") {
    return null;
  }

  try {
    const key = resolvePrivateStoreKey();
    const decipher = createDecipheriv(
      "aes-256-gcm",
      key,
      Buffer.from(iv, "base64"),
    );
    decipher.setAuthTag(Buffer.from(tag, "base64"));
    const decrypted = Buffer.concat([
      decipher.update(Buffer.from(data, "base64")),
      decipher.final(),
    ]);
    return JSON.parse(decrypted.toString("utf8"));
  } catch {
    return null;
  }
}

function readPrivateStore() {
  if (!fs.existsSync(WORKFLOW_PRIVATE_STORE_PATH)) {
    return {
      version: WORKFLOW_PRIVATE_STORE_VERSION,
      records: {},
    };
  }

  try {
    const raw = fs.readFileSync(WORKFLOW_PRIVATE_STORE_PATH, "utf8");
    const parsed = JSON.parse(raw);
    if (!parsed || typeof parsed !== "object") {
      throw new Error("Private store payload must be an object.");
    }
    const records =
      parsed.records && typeof parsed.records === "object" ? parsed.records : {};
    return {
      version: WORKFLOW_PRIVATE_STORE_VERSION,
      records,
    };
  } catch (error) {
    throw new Error(
      `Failed to read private store at ${WORKFLOW_PRIVATE_STORE_PATH}: ${
        error instanceof Error ? error.message : String(error)
      }`,
    );
  }
}

function writePrivateStore(store) {
  ensureDirectoryForFile(WORKFLOW_PRIVATE_STORE_PATH);
  const output = JSON.stringify(
    {
      version: WORKFLOW_PRIVATE_STORE_VERSION,
      records: store.records || {},
    },
    null,
    2,
  );
  const tempPath = `${WORKFLOW_PRIVATE_STORE_PATH}.tmp`;
  fs.writeFileSync(tempPath, output, "utf8");
  fs.renameSync(tempPath, WORKFLOW_PRIVATE_STORE_PATH);
}

function upsertPrivateRecord(tokenId, record) {
  const store = readPrivateStore();
  store.records[String(tokenId)] = record;
  writePrivateStore(store);
}

function getPrivateRecordForToken(tokenId, records) {
  const recordSource =
    records && typeof records === "object" ? records : readPrivateStore().records;
  const storedRecord = recordSource[String(tokenId)];
  if (!storedRecord || typeof storedRecord !== "object") {
    return null;
  }

  const decrypted = decryptPrivatePayload(storedRecord.encryptedPayload);
  if (!decrypted || typeof decrypted !== "object") {
    return null;
  }

  return {
    tokenId: String(tokenId),
    ownerAddress: String(storedRecord.ownerAddress ?? ""),
    storageType: String(storedRecord.storageType ?? ""),
    documentHash: String(storedRecord.documentHash ?? ""),
    metadataCommitment: String(storedRecord.metadataCommitment ?? ""),
    createdAt: String(storedRecord.createdAt ?? ""),
    updatedAt: String(storedRecord.updatedAt ?? ""),
    houseId: String(decrypted.houseId ?? ""),
    documentURI: String(decrypted.documentURI ?? ""),
    documentsB64: String(decrypted.documentsB64 ?? ""),
    metadata:
      decrypted.metadata && typeof decrypted.metadata === "object"
        ? decrypted.metadata
        : null,
  };
}

function getPrivateDocumentBundleForToken(tokenId, records) {
  const privateRecord = getPrivateRecordForToken(tokenId, records);
  if (!privateRecord) {
    return null;
  }

  const documentsB64 = String(privateRecord.documentsB64 ?? "").trim();
  if (!documentsB64) {
    return null;
  }

  return {
    tokenId: String(tokenId),
    documentsB64,
    documentHash: String(privateRecord.documentHash ?? ""),
    documentURI: String(privateRecord.documentURI ?? ""),
    metadata:
      privateRecord.metadata && typeof privateRecord.metadata === "object"
        ? privateRecord.metadata
        : null,
  };
}

function updatePrivateMetadataImages(tokenId, nextImages, ownerAddress) {
  const store = readPrivateStore();
  const record = store.records[String(tokenId)];
  if (!record || typeof record !== "object") {
    throw new Error("private metadata record not found for token.");
  }

  const decrypted = decryptPrivatePayload(record.encryptedPayload);
  if (!decrypted || typeof decrypted !== "object") {
    throw new Error("private metadata payload is unavailable.");
  }

  const houseId = String(decrypted.houseId ?? "").trim();
  const metadata = normalizeHouseMetadata(decrypted.metadata, houseId);
  metadata.images = parseMetadataImages(nextImages);

  const updatedAt = new Date().toISOString();
  const normalizedOwner = normalizeWalletAddress(ownerAddress);
  if (isHexAddress(normalizedOwner)) {
    record.ownerAddress = normalizedOwner;
  }
  record.updatedAt = updatedAt;
  record.encryptedPayload = encryptPrivatePayload({
    houseId,
    documentURI: String(decrypted.documentURI ?? ""),
    documentsB64: String(decrypted.documentsB64 ?? ""),
    metadata,
  });
  store.records[String(tokenId)] = record;
  writePrivateStore(store);

  return {
    tokenId: String(tokenId),
    images: metadata.images,
    updatedAt,
  };
}

function createEmptyActivityStore() {
  return {
    version: WORKFLOW_ACTIVITY_STORE_VERSION,
    notificationsByWallet: {},
    conversationsById: {},
    conversationIdsByTokenId: {},
    messagesByConversationId: {},
  };
}

function readActivityStore() {
  if (!fs.existsSync(WORKFLOW_ACTIVITY_STORE_PATH)) {
    return createEmptyActivityStore();
  }

  try {
    const raw = fs.readFileSync(WORKFLOW_ACTIVITY_STORE_PATH, "utf8");
    const parsed = JSON.parse(raw);
    if (!parsed || typeof parsed !== "object") {
      return createEmptyActivityStore();
    }

    return {
      version: WORKFLOW_ACTIVITY_STORE_VERSION,
      notificationsByWallet:
        parsed.notificationsByWallet && typeof parsed.notificationsByWallet === "object"
          ? parsed.notificationsByWallet
          : {},
      conversationsById:
        parsed.conversationsById && typeof parsed.conversationsById === "object"
          ? parsed.conversationsById
          : {},
      conversationIdsByTokenId:
        parsed.conversationIdsByTokenId
        && typeof parsed.conversationIdsByTokenId === "object"
          ? parsed.conversationIdsByTokenId
          : {},
      messagesByConversationId:
        parsed.messagesByConversationId
        && typeof parsed.messagesByConversationId === "object"
          ? parsed.messagesByConversationId
          : {},
    };
  } catch (error) {
    throw new Error(
      `Failed to read activity store at ${WORKFLOW_ACTIVITY_STORE_PATH}: ${
        error instanceof Error ? error.message : String(error)
      }`,
    );
  }
}

function writeActivityStore(store) {
  ensureDirectoryForFile(WORKFLOW_ACTIVITY_STORE_PATH);
  const output = JSON.stringify(
    {
      version: WORKFLOW_ACTIVITY_STORE_VERSION,
      notificationsByWallet: store.notificationsByWallet || {},
      conversationsById: store.conversationsById || {},
      conversationIdsByTokenId: store.conversationIdsByTokenId || {},
      messagesByConversationId: store.messagesByConversationId || {},
    },
    null,
    2,
  );
  const tempPath = `${WORKFLOW_ACTIVITY_STORE_PATH}.tmp`;
  fs.writeFileSync(tempPath, output, "utf8");
  fs.renameSync(tempPath, WORKFLOW_ACTIVITY_STORE_PATH);
}

function normalizeNotificationType(type) {
  const normalized = String(type ?? "").trim().toLowerCase();
  if (!normalized) {
    return "transaction_confirmed";
  }
  return normalized;
}

function getNotificationBucket(store, walletAddress) {
  const normalized = normalizeWalletAddress(walletAddress);
  if (!store.notificationsByWallet[normalized]) {
    store.notificationsByWallet[normalized] = [];
  }
  return store.notificationsByWallet[normalized];
}

function addNotificationToStore(store, input) {
  const walletAddress = normalizeWalletAddress(input?.walletAddress);
  if (!isHexAddress(walletAddress)) {
    return null;
  }

  const now = new Date().toISOString();
  const record = {
    id: randomUUID(),
    userId: walletAddress,
    type: normalizeNotificationType(input?.type),
    title: String(input?.title ?? "Notification").trim() || "Notification",
    message: String(input?.message ?? "").trim() || "You have a new update.",
    data:
      input?.data && typeof input.data === "object"
        ? input.data
        : {},
    read: false,
    createdAt: now,
  };

  const bucket = getNotificationBucket(store, walletAddress);
  bucket.push(record);
  if (bucket.length > WORKFLOW_ACTIVITY_MAX_NOTIFICATIONS) {
    bucket.splice(0, bucket.length - WORKFLOW_ACTIVITY_MAX_NOTIFICATIONS);
  }
  return record;
}

function addNotificationsForWallets(input, walletAddresses) {
  if (!Array.isArray(walletAddresses) || walletAddresses.length === 0) {
    return [];
  }
  const normalizedWallets = Array.from(
    new Set(
      walletAddresses
        .map((wallet) => normalizeWalletAddress(wallet))
        .filter((wallet) => isHexAddress(wallet)),
    ),
  );
  if (normalizedWallets.length === 0) {
    return [];
  }

  const store = readActivityStore();
  const records = [];
  for (const walletAddress of normalizedWallets) {
    const record = addNotificationToStore(store, {
      ...input,
      walletAddress,
    });
    if (record) {
      records.push(record);
    }
  }
  writeActivityStore(store);
  return records;
}

function sortByCreatedAtDesc(entries) {
  return entries.sort((left, right) => {
    const leftTime = Date.parse(String(left?.createdAt ?? ""));
    const rightTime = Date.parse(String(right?.createdAt ?? ""));
    return rightTime - leftTime;
  });
}

function listNotificationsForWallet(walletAddress, limit = 80) {
  const normalized = normalizeWalletAddress(walletAddress);
  if (!isHexAddress(normalized)) {
    return {
      notifications: [],
      unreadCount: 0,
    };
  }
  const store = readActivityStore();
  const bucket = Array.isArray(store.notificationsByWallet[normalized])
    ? [...store.notificationsByWallet[normalized]]
    : [];
  const notifications = sortByCreatedAtDesc(bucket).slice(0, Math.max(1, limit));
  const unreadCount = notifications.filter((entry) => !entry.read).length;
  return {
    notifications,
    unreadCount,
  };
}

function markNotificationRead(walletAddress, notificationId) {
  const normalized = normalizeWalletAddress(walletAddress);
  if (!isHexAddress(normalized)) {
    return false;
  }
  const targetId = String(notificationId ?? "").trim();
  if (!targetId) {
    return false;
  }

  const store = readActivityStore();
  const bucket = Array.isArray(store.notificationsByWallet[normalized])
    ? store.notificationsByWallet[normalized]
    : [];
  const entry = bucket.find((notification) => notification.id === targetId);
  if (!entry) {
    return false;
  }
  entry.read = true;
  writeActivityStore(store);
  return true;
}

function normalizeConversationRole(role) {
  const normalized = String(role ?? "").trim().toLowerCase();
  if (normalized === "buyer") return "buyer";
  if (normalized === "seller") return "seller";
  if (normalized === "renter") return "renter";
  if (normalized === "landlord") return "landlord";
  return "";
}

function addRoleToMap(roleMap, walletAddress, role) {
  const normalizedWallet = normalizeWalletAddress(walletAddress);
  const normalizedRole = normalizeConversationRole(role);
  if (!isHexAddress(normalizedWallet) || !normalizedRole) {
    return;
  }

  if (!roleMap.has(normalizedWallet)) {
    roleMap.set(normalizedWallet, new Set());
  }
  roleMap.get(normalizedWallet).add(normalizedRole);
}

function deriveRoleMapForHouse(house) {
  const roleMap = new Map();
  if (!house || typeof house !== "object") {
    return roleMap;
  }

  addRoleToMap(roleMap, house.ownerAddress, "seller");
  addRoleToMap(roleMap, house.originalOwner, "seller");
  addRoleToMap(roleMap, house.ownerAddress, "landlord");

  if (house?.listing?.isPrivateSale) {
    addRoleToMap(roleMap, house?.listing?.allowedBuyer, "buyer");
  }

  if (house?.rental?.renterAddress) {
    addRoleToMap(roleMap, house.rental.renterAddress, "renter");
  }

  return roleMap;
}

function resolveConversationRolePair(senderRoles, recipientRoles) {
  const sender = senderRoles || new Set();
  const recipient = recipientRoles || new Set();

  if (sender.has("seller") && recipient.has("buyer")) {
    return { senderRole: "seller", recipientRole: "buyer" };
  }
  if (sender.has("buyer") && recipient.has("seller")) {
    return { senderRole: "buyer", recipientRole: "seller" };
  }
  if (sender.has("landlord") && recipient.has("renter")) {
    return { senderRole: "landlord", recipientRole: "renter" };
  }
  if (sender.has("renter") && recipient.has("landlord")) {
    return { senderRole: "renter", recipientRole: "landlord" };
  }
  return null;
}

function isSellerForHouse(house, walletAddress) {
  const normalized = normalizeWalletAddress(walletAddress);
  return (
    normalized &&
    (normalized === normalizeWalletAddress(house?.ownerAddress) ||
      normalized === normalizeWalletAddress(house?.originalOwner))
  );
}

function isLandlordForHouse(house, walletAddress) {
  const normalized = normalizeWalletAddress(walletAddress);
  return normalized && normalized === normalizeWalletAddress(house?.ownerAddress);
}

function resolveProspectiveRolePair(house, senderWalletAddress, recipientWalletAddress) {
  const listingType = String(house?.listing?.listingType || "none").toLowerCase();
  if (listingType === "for_sale") {
    const senderIsSeller = isSellerForHouse(house, senderWalletAddress);
    const recipientIsSeller = isSellerForHouse(house, recipientWalletAddress);
    if (senderIsSeller && !recipientIsSeller) {
      return { senderRole: "seller", recipientRole: "buyer" };
    }
    if (!senderIsSeller && recipientIsSeller) {
      return { senderRole: "buyer", recipientRole: "seller" };
    }
  }
  if (listingType === "for_rent") {
    const senderIsLandlord = isLandlordForHouse(house, senderWalletAddress);
    const recipientIsLandlord = isLandlordForHouse(house, recipientWalletAddress);
    if (senderIsLandlord && !recipientIsLandlord) {
      return { senderRole: "landlord", recipientRole: "renter" };
    }
    if (!senderIsLandlord && recipientIsLandlord) {
      return { senderRole: "renter", recipientRole: "landlord" };
    }
  }
  return null;
}

function getTokenConversationIds(store, tokenId) {
  const key = String(tokenId);
  const entries = store.conversationIdsByTokenId[key];
  if (!Array.isArray(entries)) {
    store.conversationIdsByTokenId[key] = [];
    return store.conversationIdsByTokenId[key];
  }
  return entries;
}

function getConversationParticipantsKey(walletA, walletB) {
  return [normalizeWalletAddress(walletA), normalizeWalletAddress(walletB)]
    .sort()
    .join(":");
}

function findConversationForParticipants(store, tokenId, walletA, walletB) {
  const ids = getTokenConversationIds(store, tokenId);
  const key = getConversationParticipantsKey(walletA, walletB);
  for (const id of ids) {
    const conversation = store.conversationsById[id];
    if (!conversation || !Array.isArray(conversation.participants)) {
      continue;
    }
    const participantKey = getConversationParticipantsKey(
      conversation.participants[0]?.walletAddress,
      conversation.participants[1]?.walletAddress,
    );
    if (participantKey === key) {
      return conversation;
    }
  }
  return null;
}

function appendMessageToConversation(store, message) {
  const conversationId = String(message.conversationId);
  if (!store.messagesByConversationId[conversationId]) {
    store.messagesByConversationId[conversationId] = [];
  }
  const bucket = store.messagesByConversationId[conversationId];
  bucket.push(message);
  if (bucket.length > WORKFLOW_ACTIVITY_MAX_MESSAGES) {
    bucket.splice(0, bucket.length - WORKFLOW_ACTIVITY_MAX_MESSAGES);
  }
}

function createConversationRecord(input) {
  const now = new Date().toISOString();
  return {
    id: randomUUID(),
    tokenId: String(input.tokenId),
    houseId: String(input.houseId ?? ""),
    participants: [
      {
        walletAddress: normalizeWalletAddress(input.senderWalletAddress),
        role: normalizeConversationRole(input.senderRole),
      },
      {
        walletAddress: normalizeWalletAddress(input.recipientWalletAddress),
        role: normalizeConversationRole(input.recipientRole),
      },
    ],
    createdAt: now,
    updatedAt: now,
    lastMessageAt: now,
    lastMessagePreview: "",
  };
}

async function readHouseByTokenId(tokenId, contractAddress) {
  const houses = await readHousesFromChain({
    contractAddress: contractAddress || undefined,
  });
  return houses.find((house) => String(house.tokenId) === String(tokenId)) || null;
}

function formatHouseReference(house, tokenId) {
  const houseId = String(house?.houseId ?? "").trim();
  if (houseId && !houseId.toLowerCase().startsWith("private-asset-")) {
    return houseId;
  }
  return `Asset #${tokenId}`;
}

function notifyWorkflowParticipants(input) {
  try {
    addNotificationsForWallets(
      {
        type: input.type,
        title: input.title,
        message: input.message,
        data: input.data && typeof input.data === "object" ? input.data : {},
      },
      Array.isArray(input.walletAddresses) ? input.walletAddresses : [],
    );
  } catch (error) {
    console.warn(
      "[workflow-trigger] failed to persist workflow notifications:",
      error instanceof Error ? error.message : String(error),
    );
  }
}

function clampMessageContent(content) {
  const trimmed = String(content ?? "").trim();
  if (!trimmed) {
    throw new Error("message is required.");
  }
  if (trimmed.length > 2000) {
    throw new Error("message must be 2000 characters or fewer.");
  }
  return trimmed;
}

function listMessagesForConversation(store, conversationId) {
  const bucket = store.messagesByConversationId[String(conversationId)];
  if (!Array.isArray(bucket)) {
    return [];
  }
  const decrypted = bucket.map((entry) => {
    const message = { ...entry };
    if (message.encryptedContent && typeof message.encryptedContent === "object") {
      const payload = decryptPrivatePayload(message.encryptedContent);
      if (payload && typeof payload.content === "string") {
        message.content = payload.content;
      }
    }
    delete message.encryptedContent;
    return message;
  });
  return decrypted.sort((left, right) => {
    const leftTime = Date.parse(String(left?.createdAt ?? ""));
    const rightTime = Date.parse(String(right?.createdAt ?? ""));
    return leftTime - rightTime;
  });
}

function markConversationMessagesRead(store, conversationId, walletAddress) {
  const normalizedWallet = normalizeWalletAddress(walletAddress);
  const messages = store.messagesByConversationId[String(conversationId)];
  if (!Array.isArray(messages)) {
    return;
  }
  for (const message of messages) {
    if (!Array.isArray(message.readBy)) {
      message.readBy = [];
    }
    if (!message.readBy.includes(normalizedWallet)) {
      message.readBy.push(normalizedWallet);
    }
  }
}

async function sendRoleGatedMessage(payload) {
  const tokenId = String(payload?.tokenId ?? "").trim();
  if (!tokenId) {
    throw new Error("tokenId is required.");
  }

  const senderWalletAddress = normalizeWalletAddress(payload?.senderWalletAddress);
  const recipientWalletAddress = normalizeWalletAddress(payload?.recipientWalletAddress);
  if (!isHexAddress(senderWalletAddress)) {
    throw new Error("senderWalletAddress must be a valid EVM address.");
  }
  if (!isHexAddress(recipientWalletAddress)) {
    throw new Error("recipientWalletAddress must be a valid EVM address.");
  }
  if (senderWalletAddress === recipientWalletAddress) {
    throw new Error("recipientWalletAddress must differ from sender.");
  }

  const messageContent = clampMessageContent(payload?.message);
  const house = await readHouseByTokenId(tokenId, payload?.contractAddress);
  if (!house) {
    throw new Error("House not found for tokenId.");
  }

  const roleMap = deriveRoleMapForHouse(house);
  const senderRoles = roleMap.get(senderWalletAddress) || new Set();
  const recipientRoles = roleMap.get(recipientWalletAddress) || new Set();
  const rolePair =
    resolveConversationRolePair(senderRoles, recipientRoles) ||
    resolveProspectiveRolePair(house, senderWalletAddress, recipientWalletAddress);
  if (!rolePair) {
    throw new Error(
      "Conversation denied. Role-gated messaging only allows seller↔buyer or landlord↔renter pairs for active listings/agreements.",
    );
  }

  const store = readActivityStore();
  let conversation = findConversationForParticipants(
    store,
    tokenId,
    senderWalletAddress,
    recipientWalletAddress,
  );
  if (!conversation) {
    conversation = createConversationRecord({
      tokenId,
      houseId: house.houseId,
      senderWalletAddress,
      senderRole: rolePair.senderRole,
      recipientWalletAddress,
      recipientRole: rolePair.recipientRole,
    });
    store.conversationsById[conversation.id] = conversation;
    const tokenConversationIds = getTokenConversationIds(store, tokenId);
    tokenConversationIds.push(conversation.id);
    if (tokenConversationIds.length > WORKFLOW_ACTIVITY_MAX_CONVERSATIONS) {
      tokenConversationIds.splice(
        0,
        tokenConversationIds.length - WORKFLOW_ACTIVITY_MAX_CONVERSATIONS,
      );
    }
  }

  const requestedXmtpMessageId = String(payload?.xmtpMessageId ?? "").trim();
  const existingMessages = Array.isArray(
    store.messagesByConversationId[conversation.id],
  )
    ? store.messagesByConversationId[conversation.id]
    : [];
  if (requestedXmtpMessageId) {
    const existingMessage = existingMessages.find(
      (message) => message.xmtpMessageId === requestedXmtpMessageId,
    );
    if (existingMessage) {
      return {
        conversation,
        message: existingMessage,
      };
    }
  }

  const now = new Date().toISOString();
  const messageRecord = {
    id: randomUUID(),
    conversationId: conversation.id,
    tokenId,
    senderWalletAddress,
    recipientWalletAddress,
    content: "",
    encryptedContent: encryptPrivatePayload({ content: messageContent }),
    transport: "xmtp",
    xmtpMessageId: requestedXmtpMessageId || undefined,
    createdAt: now,
    readBy: [senderWalletAddress],
  };
  appendMessageToConversation(store, messageRecord);

  conversation.updatedAt = now;
  conversation.lastMessageAt = now;
  conversation.lastMessagePreview = "Encrypted private message";
  store.conversationsById[conversation.id] = conversation;

  addNotificationToStore(store, {
    walletAddress: recipientWalletAddress,
    type: "message_received",
    title: "New private XMTP message",
    message: `New private message about ${formatHouseReference(house, tokenId)}.`,
    data: {
      conversationId: conversation.id,
      tokenId,
      from: senderWalletAddress,
    },
  });

  writeActivityStore(store);

  return {
    conversation,
    message: messageRecord,
  };
}

function ensureWalletConversationAccess(conversation, walletAddress) {
  const normalized = normalizeWalletAddress(walletAddress);
  if (!conversation || !Array.isArray(conversation.participants)) {
    return false;
  }
  return conversation.participants.some(
    (participant) =>
      normalizeWalletAddress(participant?.walletAddress) === normalized,
  );
}

function summarizeConversation(conversation, messages, walletAddress) {
  const normalized = normalizeWalletAddress(walletAddress);
  const counterpart =
    conversation.participants.find(
      (participant) =>
        normalizeWalletAddress(participant.walletAddress) !== normalized,
    ) || null;
  const unreadCount = messages.filter((message) => {
    const recipient = normalizeWalletAddress(message.recipientWalletAddress);
    const readBy = Array.isArray(message.readBy) ? message.readBy : [];
    return recipient === normalized && !readBy.includes(normalized);
  }).length;

  const lastMessage = messages.length > 0 ? messages[messages.length - 1] : null;
  return {
    ...conversation,
    counterpartWalletAddress: counterpart?.walletAddress || "",
    counterpartRole: counterpart?.role || "",
    unreadCount,
    lastMessage,
  };
}

function listRoleGatedConversations(walletAddress, tokenId) {
  const normalized = normalizeWalletAddress(walletAddress);
  if (!isHexAddress(normalized)) {
    return [];
  }

  const store = readActivityStore();
  const candidates = [];
  const tokenIds = tokenId
    ? [String(tokenId)]
    : Object.keys(store.conversationIdsByTokenId);

  for (const currentTokenId of tokenIds) {
    const ids = getTokenConversationIds(store, currentTokenId);
    for (const id of ids) {
      const conversation = store.conversationsById[id];
      if (!ensureWalletConversationAccess(conversation, normalized)) {
        continue;
      }
      const messages = listMessagesForConversation(store, id);
      candidates.push(summarizeConversation(conversation, messages, normalized));
    }
  }

  return candidates.sort((left, right) => {
    const leftTime = Date.parse(String(left.lastMessageAt ?? left.updatedAt ?? ""));
    const rightTime = Date.parse(String(right.lastMessageAt ?? right.updatedAt ?? ""));
    return rightTime - leftTime;
  });
}

function getRoleGatedConversationDetails(walletAddress, conversationId) {
  const normalized = normalizeWalletAddress(walletAddress);
  if (!isHexAddress(normalized)) {
    throw new Error("walletAddress must be a valid EVM address.");
  }

  const store = readActivityStore();
  const conversation = store.conversationsById[String(conversationId)];
  if (!ensureWalletConversationAccess(conversation, normalized)) {
    throw new Error("Conversation not found or access denied.");
  }
  markConversationMessagesRead(store, conversation.id, normalized);
  const messages = listMessagesForConversation(store, conversation.id);
  writeActivityStore(store);

  return {
    conversation: summarizeConversation(conversation, messages, normalized),
    messages,
  };
}

function timeoutAfter(ms, label) {
  return new Promise((_, reject) => {
    const timeoutId = setTimeout(() => {
      reject(new Error(`${label} timed out after ${ms}ms`));
    }, ms);
    timeoutId.unref?.();
  });
}

async function withTimeout(promise, ms, label) {
  return Promise.race([promise, timeoutAfter(ms, label)]);
}

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function isRetryableRpcError(error) {
  const message =
    error instanceof Error ? error.message.toLowerCase() : String(error || "").toLowerCase();
  return (
    message.includes("rate limit") ||
    message.includes("too many requests") ||
    message.includes("request limit") ||
    message.includes("429") ||
    message.includes("-32007") ||
    message.includes("timeout") ||
    message.includes("temporar") ||
    message.includes("network") ||
    message.includes("econnreset")
  );
}

async function callWithRetry(callable, maxAttempts = 3) {
  let lastError = null;
  for (let attempt = 1; attempt <= maxAttempts; attempt++) {
    try {
      return await callable();
    } catch (error) {
      lastError = error;
      if (!isRetryableRpcError(error) || attempt === maxAttempts) {
        throw error;
      }
      await sleep(200 * attempt);
    }
  }
  throw lastError || new Error("RPC call failed.");
}

function loadEthers() {
  const candidates = [
    "ethers",
    path.join(__dirname, "../../RWA-House-UI/web/node_modules/ethers"),
  ];

  for (const candidate of candidates) {
    try {
      const imported = require(candidate);
      return imported.ethers ?? imported;
    } catch {
      continue;
    }
  }

  return null;
}

function readBody(request) {
  return new Promise((resolve, reject) => {
    let body = "";
    request.on("data", (chunk) => {
      body += chunk.toString("utf8");
      if (body.length > WORKFLOW_MAX_BODY_BYTES) {
        reject(
          new Error(
            `request payload too large (max ${WORKFLOW_MAX_BODY_BYTES} bytes)`,
          ),
        );
      }
    });
    request.on("end", () => resolve(body));
    request.on("error", reject);
  });
}

function parseJsonBody(body) {
  if (!body.trim()) {
    return {};
  }
  return JSON.parse(body);
}

function validateRequiredEnv() {
  if (!WORKFLOW_RPC_URL) {
    throw new Error(
      "WORKFLOW_RPC_URL is required (or set CRE_RPC_URL / SEPOLIA_RPC).",
    );
  }
  if (!WORKFLOW_PRIVATE_KEY) {
    throw new Error(
      "WORKFLOW_PRIVATE_KEY is required (or set CRE_ETH_PRIVATE_KEY / PRIVATE_KEY).",
    );
  }
  if (!WORKFLOW_CONTRACT_ADDRESS) {
    throw new Error(
      "HOUSE_RWA_CONTRACT_ADDRESS is required for /workflow/trigger.",
    );
  }
}

async function buildRuntime() {
  const ethers = loadEthers();
  if (!ethers) {
    throw new Error(
      "Ethers library not found. Install dependencies in RWA-House-UI/web first.",
    );
  }

  validateRequiredEnv();

  const provider = new ethers.JsonRpcProvider(WORKFLOW_RPC_URL);
  const wallet = new ethers.Wallet(WORKFLOW_PRIVATE_KEY, provider);
  const contract = new ethers.Contract(
    WORKFLOW_CONTRACT_ADDRESS,
    HOUSE_RWA_ABI,
    wallet,
  );

  const [network, code] = await Promise.all([
    provider.getNetwork(),
    provider.getCode(WORKFLOW_CONTRACT_ADDRESS),
  ]);

  if (!code || code === "0x") {
    throw new Error(
      `No contract code found at ${WORKFLOW_CONTRACT_ADDRESS} on WORKFLOW_RPC_URL.`,
    );
  }

  if (Number(network.chainId) !== WORKFLOW_CHAIN_ID) {
    throw new Error(
      `RPC chain mismatch. Expected ${WORKFLOW_CHAIN_ID}, got ${network.chainId}.`,
    );
  }

  return { ethers, provider, wallet, contract };
}

async function getRuntime() {
  if (!cachedRuntimePromise) {
    cachedRuntimePromise = buildRuntime();
  }
  return cachedRuntimePromise;
}

async function getReadContract(overrideContractAddress) {
  const runtime = await getRuntime();
  const { contract, provider, wallet, ethers } = runtime;
  const normalizedDefault = WORKFLOW_CONTRACT_ADDRESS.toLowerCase();
  const normalizedOverride = String(overrideContractAddress || "").trim().toLowerCase();

  if (!normalizedOverride || normalizedOverride === normalizedDefault) {
    return { ...runtime, contract };
  }

  if (!isHexAddress(normalizedOverride)) {
    throw new Error("contractAddress override must be a valid EVM address.");
  }

  const code = await provider.getCode(normalizedOverride);
  if (!code || code === "0x") {
    throw new Error(`No contract code found at ${normalizedOverride}.`);
  }

  const overrideContract = new ethers.Contract(
    normalizedOverride,
    HOUSE_RWA_ABI,
    wallet,
  );
  return { ...runtime, contract: overrideContract };
}

async function readTokenFromChain(tokenId, contract, ethers, privateRecords = null) {
  const tokenLookup = typeof tokenId === "bigint" ? tokenId.toString() : tokenId;
  let ownerAddress;
  try {
    ownerAddress = String(await callWithRetry(() => contract.ownerOf(tokenLookup)));
  } catch {
    return null;
  }

  const houseDetailsRaw = await callWithRetry(
    () => contract.getHouseDetails(tokenLookup),
  ).catch(() => null);
  const listingRaw = await callWithRetry(
    () => contract.getListing(tokenLookup),
  ).catch(() => null);
  const billsRaw = await callWithRetry(
    () => contract.getBills(tokenLookup),
  ).catch(() => []);
  const rentalRaw = await callWithRetry(
    () => contract.getActiveRental(tokenLookup),
  ).catch(() => null);

  if (!houseDetailsRaw) {
    return null;
  }

  const houseDetails = houseDetailsRaw;
  const listing = listingRaw || {};
  const rental = rentalRaw || {};
  const bills = Array.isArray(billsRaw) ? billsRaw : [];

  const onChainHouseId = String(
    houseDetails.houseId ?? houseDetails[0] ?? `Token-${tokenId.toString()}`,
  );
  const onChainDocumentHash = String(houseDetails.documentHash ?? houseDetails[1] ?? "0x");
  const onChainDocumentURI = String(houseDetails.documentURI ?? houseDetails[2] ?? "");
  const storageTypeEnum = Number(houseDetails.storageType ?? houseDetails[3] ?? 0);
  const originalOwner = String(
    houseDetails.originalOwner ?? houseDetails[4] ?? ownerAddress,
  );
  const mintedAt = toDate(houseDetails.mintedAt ?? houseDetails[5] ?? 0n);
  const isVerified = Boolean(houseDetails.isVerified ?? houseDetails[6] ?? false);
  const privateRecord = getPrivateRecordForToken(tokenLookup, privateRecords);
  const privateMetadata =
    privateRecord && privateRecord.metadata && typeof privateRecord.metadata === "object"
      ? privateRecord.metadata
      : null;
  const fallbackHouseId = isHexCommitment(onChainHouseId)
    ? `Private-Asset-${tokenLookup}`
    : onChainHouseId;
  const houseId = privateRecord?.houseId || fallbackHouseId;
  const documentHash = privateRecord?.documentHash || onChainDocumentHash;
  const documentURI = privateRecord?.documentURI || onChainDocumentURI;
  const metadata = normalizeHouseMetadata(privateMetadata, houseId);

  const listingTypeEnum = Number(listing.listingType ?? listing[0] ?? 0);
  const listingType = mapListingType(listingTypeEnum);
  const listingPrice = toBigInt(listing.price ?? listing[1] ?? 0n);
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
          isPrivateSale: Boolean(listing.isPrivateSale ?? listing[3] ?? false),
          allowedBuyer: isZeroAddress(allowedBuyer) ? undefined : allowedBuyer,
          createdAt: toDate(listing.createdAt ?? listing[5] ?? 0n),
          expiresAt: toDate(listing.expiresAt ?? listing[6] ?? 0n),
        };

  const renter = String(rental.renter ?? rental[0] ?? ethers.ZeroAddress);
  const accessKeyHash = String(
    rental.encryptedAccessKeyHash ?? rental[6] ?? ethers.ZeroHash,
  );
  const mappedRental = isZeroAddress(renter)
    ? undefined
    : {
        tokenId: tokenId.toString(),
        renterAddress: renter,
        startTime: toDate(rental.startTime ?? rental[1] ?? 0n),
        endTime: toDate(rental.endTime ?? rental[2] ?? 0n),
        depositAmount: toBigInt(rental.depositAmount ?? rental[3] ?? 0n).toString(),
        depositFormatted: `${ethers.formatEther(toBigInt(rental.depositAmount ?? rental[3] ?? 0n))} ETH`,
        monthlyRent: toBigInt(rental.monthlyRent ?? rental[4] ?? 0n).toString(),
        isActive: Boolean(rental.isActive ?? rental[5] ?? false),
        hasAccessKey: accessKeyHash !== ethers.ZeroHash,
      };

  const mappedBills = bills.map((bill, idx) => {
    const amountCents = Number(toBigInt(bill.amount ?? bill[1] ?? 0n));
    const amountDollars = amountCents / 100;
    const dueDate = toDate(bill.dueDate ?? bill[2] ?? 0n);
    const paidAtRaw = toBigInt(bill.paidAt ?? bill[3] ?? 0n);
    const paymentReference = String(
      bill.paymentReference ?? bill[5] ?? ethers.ZeroHash,
    );
    const provider = String(bill.provider ?? bill[7] ?? ethers.ZeroAddress);
    const status = Number(bill.status ?? bill[4] ?? 0);

    return {
      id: `${tokenId.toString()}-${idx}`,
      tokenId: tokenId.toString(),
      billType: normalizeBillType(String(bill.billType ?? bill[0] ?? "other")),
      amount: amountDollars,
      amountFormatted: `$${amountDollars.toFixed(2)}`,
      dueDate,
      isPaid: status === 1,
      paidAt: paidAtRaw > 0n ? toDate(paidAtRaw) : undefined,
      paymentMethod: status === 1 ? "crypto" : undefined,
      paymentReference:
        paymentReference !== ethers.ZeroHash ? paymentReference : undefined,
      isRecurring: Boolean(bill.isRecurring ?? bill[6] ?? false),
      provider,
      providerName: isZeroAddress(provider)
        ? "Unknown Provider"
        : `${provider.slice(0, 6)}...${provider.slice(-4)}`,
      createdAt: dueDate,
    };
  });

  return {
    tokenId: tokenId.toString(),
    houseId,
    ownerAddress,
    originalOwner,
    documentHash,
    documentURI,
    storageType: mapStorageType(storageTypeEnum),
    mintedAt,
    isVerified,
    metadata,
    metadataCommitment: isHexCommitment(onChainHouseId) ? onChainHouseId : undefined,
    listing: mappedListing,
    rental: mappedRental,
    bills: mappedBills,
  };
}

async function readHousesFromChain(options = {}) {
  const ownerAddress = String(options.ownerAddress || "").trim();
  const contractAddress = String(options.contractAddress || "").trim();
  const { contract, ethers } = await getReadContract(contractAddress);
  let privateRecords = {};
  try {
    const privateStore = readPrivateStore();
    privateRecords =
      privateStore.records && typeof privateStore.records === "object"
        ? privateStore.records
        : {};
  } catch (error) {
    console.warn(
      "[workflow-trigger] unable to read private metadata store; continuing with onchain-only fields",
      error instanceof Error ? error.message : String(error),
    );
  }
  const nextTokenId = toBigInt(await contract.nextTokenId());
  const maxScan = BigInt(
    Number.parseInt(String(process.env.WORKFLOW_MAX_SCAN || "500"), 10) || 500,
  );
  const scanCount = nextTokenId > maxScan ? maxScan : nextTokenId;
  const tokenIds = [];
  for (let i = 0n; i < scanCount; i++) {
    tokenIds.push(i);
  }
  const houses = [];
  for (const tokenId of tokenIds) {
    const token = await readTokenFromChain(
      tokenId,
      contract,
      ethers,
      privateRecords,
    );
    if (token) {
      houses.push(token);
    }
  }

  if (ownerAddress) {
    const normalized = ownerAddress.toLowerCase();
    return houses.filter(
      (house) => String(house.ownerAddress).toLowerCase() === normalized,
    );
  }
  return houses;
}

async function readNativeBalance(address) {
  const { provider, ethers } = await getRuntime();
  const value = await provider.getBalance(address);
  return ethers.formatEther(value);
}

async function readKYCStatusForWallet(walletAddress) {
  const normalized = normalizeWalletAddress(walletAddress);
  if (!isHexAddress(normalized)) {
    return "unverified";
  }

  try {
    const { contract } = await getRuntime();
    const hasKYC = await contract.hasKYC(normalized);
    return Boolean(hasKYC) ? "verified" : "unverified";
  } catch {
    return "unverified";
  }
}

async function buildAuthenticatedUser(walletAddress, emailHint) {
  const normalizedWallet = normalizeWalletAddress(walletAddress);
  const now = new Date().toISOString();
  const kycStatus = await readKYCStatusForWallet(normalizedWallet);

  return {
    id: normalizedWallet,
    email:
      String(emailHint ?? "").trim()
      || `${normalizedWallet.slice(2, 10)}@demo.rwa.house`,
    walletAddress: normalizedWallet,
    chainId: WORKFLOW_CHAIN_ID,
    kycStatus,
    createdAt: now,
    lastLoginAt: now,
    mfaEnabled: false,
    preferences: {
      theme: "dark",
      currency: "USD",
      language: "en",
      autoPayEnabled: false,
      autoPayThreshold: 1000,
      notifications: {
        email: true,
        push: false,
        sms: false,
        transactions: true,
        bills: true,
        security: true,
      },
    },
  };
}

function isHexAddress(value) {
  return /^0x[a-fA-F0-9]{40}$/.test(String(value ?? ""));
}

function requireAuthenticatedActor(payload, message) {
  const actorAddress = normalizeWalletAddress(payload?.actorAddress);
  if (!isHexAddress(actorAddress)) {
    throw new Error(message || "Authentication is required for this action.");
  }
  return actorAddress;
}

function extractWalletAddressFromAuthHeader(rawHeader) {
  const header = typeof rawHeader === "string" ? rawHeader.trim() : "";
  if (!header) {
    return "";
  }

  if (!header.toLowerCase().startsWith("bearer ")) {
    return "";
  }
  const token = header.slice("bearer ".length).trim();
  const signedWalletAddress = parseSignedAuthToken(token);
  if (signedWalletAddress) {
    return signedWalletAddress;
  }

  if (!WORKFLOW_ALLOW_INSECURE_BEARER) {
    return "";
  }

  if (isHexAddress(token)) {
    return normalizeWalletAddress(token);
  }
  if (token.toLowerCase().startsWith("demo-")) {
    const candidate = token.slice("demo-".length).trim();
    if (isHexAddress(candidate)) {
      return normalizeWalletAddress(candidate);
    }
  }
  return "";
}

function normalizeWalletAddress(value) {
  return String(value ?? "").trim().toLowerCase();
}

function extractBoundWallet(queryResult) {
  const bindSection =
    queryResult && typeof queryResult === "object" ? queryResult.bind : undefined;
  if (!bindSection || typeof bindSection !== "object") {
    return "";
  }

  const boundWallet = String(bindSection.user_address ?? "").trim();
  return isHexAddress(boundWallet) ? normalizeWalletAddress(boundWallet) : "";
}

function parseZKProofBundle(rawPayload) {
  if (!rawPayload || typeof rawPayload !== "object") {
    throw new Error("zkpassport proof payload must be an object.");
  }

  const proofs = Array.isArray(rawPayload.proofs) ? rawPayload.proofs : [];
  const queryResult =
    rawPayload.queryResult && typeof rawPayload.queryResult === "object"
      ? rawPayload.queryResult
      : rawPayload.result && typeof rawPayload.result === "object"
        ? rawPayload.result
        : null;

  if (proofs.length === 0) {
    throw new Error("zkpassport proof payload requires at least one proof.");
  }
  if (!queryResult) {
    throw new Error("zkpassport proof payload requires queryResult.");
  }

  return {
    proofs,
    queryResult,
    domain: String(rawPayload.domain ?? "").trim(),
    scope: String(rawPayload.scope ?? "").trim(),
    validity: parsePositiveInteger(
      rawPayload.validity ?? rawPayload.validitySeconds,
      ZKPASSPORT_VALIDITY_SECONDS,
    ),
    devModeRequested: parseBoolean(rawPayload.devMode, ZKPASSPORT_DEV_MODE),
  };
}

function parseBirthdateValue(rawValue) {
  if (rawValue instanceof Date) {
    return Number.isNaN(rawValue.getTime()) ? null : rawValue;
  }

  if (typeof rawValue !== "string" && typeof rawValue !== "number") {
    return null;
  }

  const parsedDate = new Date(String(rawValue));
  if (Number.isNaN(parsedDate.getTime())) {
    return null;
  }
  return parsedDate;
}

function isAtLeast18YearsOld(birthdate) {
  const now = new Date();
  let age = now.getUTCFullYear() - birthdate.getUTCFullYear();
  const monthDelta = now.getUTCMonth() - birthdate.getUTCMonth();
  const dayDelta = now.getUTCDate() - birthdate.getUTCDate();
  if (monthDelta < 0 || (monthDelta === 0 && dayDelta < 0)) {
    age -= 1;
  }
  return age >= 18;
}

function requireAdultFromQueryResult(queryResult) {
  const ageCheck = queryResult?.age?.gte;
  if (ageCheck && ageCheck.result === true) {
    return;
  }

  const disclosedBirthdate = parseBirthdateValue(
    queryResult?.birthdate?.disclose?.result,
  );
  if (!disclosedBirthdate) {
    throw new Error(
      "zkpassport proof must include either age.gte >= 18 or disclose birthdate.",
    );
  }
  if (!isAtLeast18YearsOld(disclosedBirthdate)) {
    throw new Error("zkpassport proof indicates user is under 18.");
  }
}

function isBytes32Hex(value) {
  return /^0x[a-fA-F0-9]{64}$/.test(String(value ?? ""));
}

function parseTokenId(payload) {
  const raw = payload?.tokenId ?? payload?.tokenID;
  if (raw === undefined || raw === null || String(raw).trim() === "") {
    throw new Error("tokenId is required.");
  }
  const parsed = BigInt(String(raw).trim());
  if (parsed < 0n) {
    throw new Error("tokenId must be non-negative.");
  }
  return parsed;
}

function parsePositiveBigInt(value, fieldName) {
  if (value === undefined || value === null || String(value).trim() === "") {
    throw new Error(`${fieldName} is required.`);
  }
  const normalized = String(value).trim();
  if (!/^\d+$/.test(normalized)) {
    throw new Error(`${fieldName} must be an integer string in wei.`);
  }
  const parsed = BigInt(normalized);
  if (parsed <= 0n) {
    throw new Error(`${fieldName} must be greater than zero.`);
  }
  return parsed;
}

function parseDurationDays(value) {
  const parsed = Number.parseInt(String(value ?? ""), 10);
  if (!Number.isFinite(parsed) || parsed <= 0) {
    throw new Error("durationDays must be a positive integer.");
  }
  return BigInt(parsed);
}

function parseListingType(value) {
  const normalized = String(value ?? "").trim().toLowerCase();
  if (normalized === "for_sale") {
    return 1;
  }
  if (normalized === "for_rent") {
    return 2;
  }
  throw new Error("listingType must be either `for_sale` or `for_rent`.");
}

function parseListingDurationDays(value) {
  if (value === undefined || value === null || String(value).trim() === "") {
    return 0;
  }
  const parsed = Number.parseInt(String(value), 10);
  if (!Number.isFinite(parsed) || parsed < 0) {
    throw new Error("durationDays must be zero or a positive integer.");
  }
  return parsed;
}

function normalizeCreateListingError(error) {
  const message = error instanceof Error ? error.message : String(error ?? "");
  const normalized = message.toLowerCase();
  const missingListingMethod =
    normalized.includes("createlistingfromworkflow")
    && (
      normalized.includes("no matching fragment")
      || normalized.includes("is not a function")
      || normalized.includes("function selector was not recognized")
      || normalized.includes("execution reverted")
    );

  if (missingListingMethod) {
    return new Error(
      "HouseRWA contract does not expose createListingFromWorkflow. "
      + "Redeploy/upgrade contracts, then retry listing creation.",
    );
  }

  return error;
}

function parseNonNegativeInteger(value, fieldName) {
  const parsed = Number.parseInt(String(value ?? ""), 10);
  if (!Number.isFinite(parsed) || parsed < 0) {
    throw new Error(`${fieldName} must be a non-negative integer.`);
  }
  return parsed;
}

function parseBillAmountInCents(rawAmount) {
  const amountNumber = Number(rawAmount);
  if (!Number.isFinite(amountNumber) || amountNumber <= 0) {
    throw new Error("amount must be greater than zero.");
  }

  const cents = Math.round(amountNumber * 100);
  if (cents <= 0) {
    throw new Error("amount is too small.");
  }
  return BigInt(cents);
}

function parseDueDate(rawDueDate) {
  const dueDate = new Date(String(rawDueDate ?? ""));
  if (Number.isNaN(dueDate.getTime())) {
    throw new Error("dueDate must be a valid ISO timestamp.");
  }
  const dueUnix = Math.floor(dueDate.getTime() / 1000);
  if (dueUnix <= Math.floor(Date.now() / 1000)) {
    throw new Error("dueDate must be in the future.");
  }
  return { dueDate, dueUnix: BigInt(dueUnix) };
}

function normalizeStorageType(value) {
  const normalized = String(value ?? "ipfs").trim().toLowerCase();
  if (normalized !== "ipfs" && normalized !== "offchain") {
    throw new Error("storageType must be either `ipfs` or `offchain`.");
  }
  return normalized === "offchain" ? 1 : 0;
}

async function submitTransaction(txPromise, label) {
  const tx = await withTimeout(txPromise, WORKFLOW_TIMEOUT_MS, `${label} submit`);
  const receipt = await withTimeout(
    tx.wait(WORKFLOW_CONFIRMATIONS),
    WORKFLOW_TIMEOUT_MS,
    `${label} confirmation`,
  );

  if (!receipt || Number(receipt.status) !== 1) {
    throw new Error(`${label} transaction reverted.`);
  }

  return { tx, receipt };
}

async function ensureKYC(runtime, userAddress) {
  const { ethers, contract } = runtime;
  if (!isHexAddress(userAddress)) {
    throw new Error("Invalid KYC address.");
  }

  const verificationHash = ethers.keccak256(
    ethers.toUtf8Bytes(`kyc:${userAddress.toLowerCase()}:${Date.now()}`),
  );
  const expiryUnix =
    Math.floor(Date.now() / 1000) + KYC_EXPIRY_DAYS * 24 * 60 * 60;

  await submitTransaction(
    contract.setKYCVerification(
      userAddress,
      KYC_LEVEL,
      verificationHash,
      BigInt(expiryUnix),
    ),
    "setKYCVerification",
  );
}

async function verifyZKPassportProof(walletAddress, rawProof) {
  const normalizedWallet = normalizeWalletAddress(walletAddress);
  if (!isHexAddress(normalizedWallet)) {
    throw new Error("walletAddress must be a valid EVM address.");
  }

  const parsed = parseZKProofBundle(rawProof);
  const enforcedScope = ZKPASSPORT_SCOPE.trim();
  if (enforcedScope && parsed.scope && parsed.scope !== enforcedScope) {
    throw new Error(
      `zkpassport scope mismatch. Expected ${enforcedScope}, got ${parsed.scope}.`,
    );
  }

  const scope = enforcedScope || parsed.scope || undefined;
  const requestedDomain = parsed.domain.trim();
  const normalizedRequestedDomain = requestedDomain.toLowerCase();
  const normalizedConfiguredDomain = ZKPASSPORT_DOMAIN.trim().toLowerCase();
  const effectiveDomain =
    normalizedRequestedDomain && normalizedRequestedDomain !== normalizedConfiguredDomain
      ? requestedDomain
      : ZKPASSPORT_DOMAIN;
  const effectiveDevMode = ZKPASSPORT_DEV_MODE;
  if (!effectiveDevMode && parsed.devModeRequested) {
    throw new Error("zkpassport devMode proofs are disabled for this deployment.");
  }

  const zkPassport = new ZKPassport(effectiveDomain);
  const verification = await zkPassport.verify({
    proofs: parsed.proofs,
    queryResult: parsed.queryResult,
    scope,
    validity: parsed.validity,
    devMode: effectiveDevMode,
    writingDirectory: ZKPASSPORT_VERIFY_WRITING_DIRECTORY,
  });

  if (!verification.verified) {
    throw new Error("zkpassport proof verification failed.");
  }

  requireAdultFromQueryResult(parsed.queryResult);

  const boundWallet = extractBoundWallet(parsed.queryResult);
  if (boundWallet && boundWallet !== normalizedWallet) {
    throw new Error("zkpassport proof wallet binding does not match walletAddress.");
  }

  const uniqueIdentifier = String(verification.uniqueIdentifier ?? "").trim();
  const verificationHash = `0x${createHash("sha256")
    .update(
      JSON.stringify({
        provider: "zkpassport",
        domain: effectiveDomain,
        scope: scope || "",
        walletAddress: normalizedWallet,
        uniqueIdentifier,
      }),
    )
    .digest("hex")}`;

  return {
    verified: true,
    level: KYC_LEVEL,
    verificationHash,
    expiresAt:
      Math.floor(Date.now() / 1000) + KYC_EXPIRY_DAYS * 24 * 60 * 60,
    uniqueIdentifier,
    nullifierType: verification.uniqueIdentifierType,
    domain: effectiveDomain,
    scope: scope || "",
    proof: {
      provider: "zkpassport",
      proofs: parsed.proofs,
      queryResult: parsed.queryResult,
      domain: effectiveDomain,
      scope: scope || "",
      devMode: effectiveDevMode,
      validity: parsed.validity,
      uniqueIdentifier,
      nullifierType: verification.uniqueIdentifierType,
      verified: true,
      verifiedAt: new Date().toISOString(),
    },
  };
}

async function ensureKYCFromPayload(runtime, payload, actorAddress) {
  const provider = String(payload?.kycProvider ?? "mock").trim().toLowerCase();
  if (provider === "none") {
    return;
  }

  if (provider === "zkpassport") {
    if (!payload?.kycProof || typeof payload.kycProof !== "object") {
      throw new Error(
        "kycProof is required when kycProvider=zkpassport.",
      );
    }

    const verification = await verifyZKPassportProof(
      actorAddress,
      payload.kycProof,
    );

    const { contract } = runtime;
    await submitTransaction(
      contract.setKYCVerification(
        actorAddress,
        verification.level,
        verification.verificationHash,
        BigInt(verification.expiresAt),
      ),
      "setKYCVerification",
    );
    return;
  }

  await ensureKYC(runtime, actorAddress);
}

function encodeBase64(value, ethers) {
  return Buffer.from(ethers.getBytes(value)).toString("base64");
}

async function handleMint(payload) {
  const runtime = await getRuntime();
  const { contract } = runtime;
  const actorAddress = requireAuthenticatedActor(
    payload,
    "Authentication is required to mint properties.",
  );

  const ownerAddress = String(payload?.ownerAddress ?? "").trim();
  if (!isHexAddress(ownerAddress)) {
    throw new Error("ownerAddress must be a valid EVM address.");
  }
  const normalizedOwner = normalizeWalletAddress(ownerAddress);
  if (normalizedOwner !== actorAddress) {
    throw new Error("ownerAddress must match the authenticated wallet address.");
  }

  const houseId = String(payload?.houseID ?? payload?.houseId ?? "").trim();
  if (!houseId) {
    throw new Error("houseId is required.");
  }

  const metadata = normalizeHouseMetadata(payload?.metadata, houseId);
  const documentsB64 = String(payload?.documentsB64 ?? "").trim();
  if (!documentsB64) {
    throw new Error("documentsB64 is required.");
  }

  const decodedDocs = Buffer.from(documentsB64, "base64");
  if (!decodedDocs.length) {
    throw new Error("documentsB64 must decode to a non-empty payload.");
  }

  const storageType = normalizeStorageType(payload?.storageType);
  const privacySaltHex = randomBytes(32).toString("hex");
  const documentHash = `0x${createHash("sha256")
    .update(Buffer.from(privacySaltHex, "hex"))
    .update(decodedDocs)
    .digest("hex")}`;
  const documentURI = `cre://private/${randomBytes(18).toString("hex")}`;
  const metadataCommitment = `0x${createHash("sha256")
    .update(
      stableStringify({
        ownerAddress: ownerAddress.toLowerCase(),
        houseId,
        metadata,
        documentHash,
        documentURI,
        storageType,
        privacySaltHex,
      }),
    )
    .digest("hex")}`;
  const verificationData = JSON.stringify({
    source: "workflow-trigger-api-private",
    createdAt: new Date().toISOString(),
    metadataCommitment,
    mode: "commitment-only-onchain",
  });

  await ensureKYCFromPayload(runtime, payload, actorAddress);
  const { receipt } = await submitTransaction(
    contract.mint(
      actorAddress,
      metadataCommitment,
      documentHash,
      documentURI,
      storageType,
      verificationData,
    ),
    "mint",
  );

  const nextTokenId = await contract.nextTokenId();
  const tokenId = nextTokenId > 0n ? (nextTokenId - 1n).toString() : "0";
  const encryptedKey = randomBytes(32).toString("base64");
  const persistedAt = new Date().toISOString();
  let persistenceWarning = "";

  try {
    upsertPrivateRecord(tokenId, {
      ownerAddress: actorAddress,
      storageType: storageType === 1 ? "offchain" : "ipfs",
      documentHash,
      metadataCommitment,
      createdAt: persistedAt,
      updatedAt: persistedAt,
      encryptedPayload: encryptPrivatePayload({
        houseId,
        documentURI,
        documentsB64,
        metadata,
      }),
    });
  } catch (error) {
    persistenceWarning =
      error instanceof Error
        ? error.message
        : "Unknown private-metadata persistence error.";
  }

  notifyWorkflowParticipants({
    walletAddresses: [actorAddress],
    type: "transaction_confirmed",
    title: "Property minted",
    message: `Your property ${houseId} was minted with private CRE metadata.`,
    data: {
      action: "mint",
      tokenId,
      houseID: houseId,
      txHash: receipt.hash,
    },
  });

  return {
    success: true,
    message: persistenceWarning
      ? "house minted, but private metadata persistence needs attention"
      : "house minted successfully with private onchain commitment",
    txHash: receipt.hash,
    encryptedKey,
    data: {
      tokenId,
      tokenID: tokenId,
      houseID: houseId,
      documentHash,
      documentURI,
      metadataCommitment,
      storageType: storageType === 1 ? "offchain" : "ipfs",
      kycProvider: String(payload?.kycProvider ?? "mock"),
      warning: persistenceWarning || undefined,
    },
  };
}

async function handleUpdateHouseImages(payload) {
  const runtime = await getRuntime();
  const { contract } = runtime;

  const tokenId = parseTokenId(payload);
  const actorAddress = normalizeWalletAddress(payload?.actorAddress);
  if (!isHexAddress(actorAddress)) {
    throw new Error("Authentication is required to update property images.");
  }

  const ownerAddress = normalizeWalletAddress(await contract.ownerOf(tokenId));
  if (!isHexAddress(ownerAddress)) {
    throw new Error("Unable to determine current property owner.");
  }
  if (ownerAddress !== actorAddress) {
    throw new Error("Only the current property owner can update images.");
  }

  const imageInput = Array.isArray(payload?.images)
    ? payload.images
    : payload?.metadata?.images;
  const updated = updatePrivateMetadataImages(
    tokenId.toString(),
    parseMetadataImages(imageInput),
    ownerAddress,
  );

  notifyWorkflowParticipants({
    walletAddresses: [ownerAddress],
    type: "transaction_confirmed",
    title: "Property images updated",
    message: `Private image gallery updated for token #${tokenId.toString()}.`,
    data: {
      action: "update_house_images",
      tokenId: tokenId.toString(),
      imagesCount: updated.images.length,
    },
  });

  return {
    success: true,
    message: "property images updated successfully",
    data: {
      tokenID: tokenId.toString(),
      tokenId: tokenId.toString(),
      images: updated.images,
      updatedAt: updated.updatedAt,
    },
  };
}

async function handleCreateListing(payload) {
  const runtime = await getRuntime();
  const { ethers, contract } = runtime;
  const actorAddress = requireAuthenticatedActor(
    payload,
    "Authentication is required to create listings.",
  );

  const tokenId = parseTokenId(payload);
  const ownerAddress = String(payload?.ownerAddress ?? "").trim();
  if (!isHexAddress(ownerAddress)) {
    throw new Error("ownerAddress must be a valid EVM address.");
  }
  const normalizedOwnerAddress = normalizeWalletAddress(ownerAddress);
  if (normalizedOwnerAddress !== actorAddress) {
    throw new Error("ownerAddress must match the authenticated wallet address.");
  }
  const onchainOwner = normalizeWalletAddress(await contract.ownerOf(tokenId));
  if (onchainOwner !== actorAddress) {
    throw new Error("Only the current onchain property owner can create listings.");
  }

  const listingType = parseListingType(payload?.listingType);
  const price = parsePositiveBigInt(payload?.price, "price");
  const durationDays = parseListingDurationDays(payload?.durationDays);
  const isPrivateSale = Boolean(payload?.isPrivateSale);

  const preferredTokenRaw = String(payload?.preferredToken ?? "").trim();
  const preferredToken = preferredTokenRaw ? preferredTokenRaw : ethers.ZeroAddress;
  if (!isHexAddress(preferredToken)) {
    throw new Error("preferredToken must be a valid EVM address.");
  }

  const allowedBuyerRaw = String(payload?.allowedBuyer ?? "").trim();
  if (isPrivateSale && !isHexAddress(allowedBuyerRaw)) {
    throw new Error("allowedBuyer must be a valid EVM address for private listings.");
  }
  const allowedBuyer = isPrivateSale ? allowedBuyerRaw : ethers.ZeroAddress;

  await ensureKYCFromPayload(runtime, payload, actorAddress);
  let receipt;
  try {
    ({ receipt } = await submitTransaction(
      contract.createListingFromWorkflow(
        tokenId,
        actorAddress,
        listingType,
        price,
        preferredToken,
        isPrivateSale,
        allowedBuyer,
        BigInt(durationDays),
      ),
      "createListingFromWorkflow",
    ));
  } catch (error) {
    throw normalizeCreateListingError(error);
  }

  return {
    success: true,
    message: "listing created successfully",
    txHash: receipt.hash,
    data: {
      tokenID: tokenId.toString(),
      ownerAddress: actorAddress,
      listingType: listingType === 1 ? "for_sale" : "for_rent",
      price: price.toString(),
      preferredToken,
      isPrivateSale,
      allowedBuyer: isPrivateSale ? allowedBuyer : null,
      durationDays,
      kycProvider: String(payload?.kycProvider ?? "mock"),
    },
  };
}

async function handleSell(payload) {
  const runtime = await getRuntime();
  const { ethers, contract } = runtime;
  const actorAddress = requireAuthenticatedActor(
    payload,
    "Authentication is required to execute property purchases.",
  );

  const tokenId = parseTokenId(payload);
  let sellerAddress = String(payload?.sellerAddress ?? "").trim().toLowerCase();
  if (sellerAddress && !isHexAddress(sellerAddress)) {
    throw new Error("sellerAddress must be a valid EVM address.");
  }
  const buyerAddress = String(payload?.buyerAddress ?? "").trim();
  if (!isHexAddress(buyerAddress)) {
    throw new Error("buyerAddress must be a valid EVM address.");
  }
  if (normalizeWalletAddress(buyerAddress) !== actorAddress) {
    throw new Error("buyerAddress must match the authenticated wallet address.");
  }

  parsePositiveBigInt(payload?.price, "price");
  if (!sellerAddress) {
    sellerAddress = normalizeWalletAddress(await contract.ownerOf(tokenId));
  }
  const houseDetails = await contract.getHouseDetails(tokenId).catch(() => null);
  const houseReference = formatHouseReference(
    houseDetails ? { houseId: houseDetails.houseId ?? houseDetails[0] ?? "" } : null,
    tokenId.toString(),
  );

  await ensureKYCFromPayload(runtime, payload, actorAddress);
  const encryptedKey = randomBytes(32);
  const keyHash = ethers.keccak256(encryptedKey);

  const { receipt } = await submitTransaction(
    contract.completeSale(tokenId, buyerAddress, keyHash, encryptedKey),
    "completeSale",
  );

  notifyWorkflowParticipants({
    walletAddresses: [sellerAddress],
    type: "listing_sold",
    title: "Property sold",
    message: `${houseReference} has been sold to ${buyerAddress.slice(0, 6)}...${buyerAddress.slice(-4)}.`,
    data: {
      action: "sell",
      tokenId: tokenId.toString(),
      buyerAddress,
      txHash: receipt.hash,
    },
  });
  notifyWorkflowParticipants({
    walletAddresses: [buyerAddress],
    type: "transaction_confirmed",
    title: "Purchase completed",
    message: `You purchased ${houseReference}.`,
    data: {
      action: "buy",
      tokenId: tokenId.toString(),
      sellerAddress,
      txHash: receipt.hash,
    },
  });

  return {
    success: true,
    message: "sale completed successfully",
    txHash: receipt.hash,
    encryptedKey: encryptedKey.toString("base64"),
    data: {
      tokenID: tokenId.toString(),
      buyer: buyerAddress,
      price: String(payload?.price),
      isPrivateSale: Boolean(payload?.isPrivateSale),
      keyHash,
      encryptedKey: encryptedKey.toString("base64"),
      kycProvider: String(payload?.kycProvider ?? "mock"),
    },
  };
}

async function handleRent(payload) {
  const runtime = await getRuntime();
  const { ethers, contract } = runtime;
  const actorAddress = requireAuthenticatedActor(
    payload,
    "Authentication is required to start rentals.",
  );

  const tokenId = parseTokenId(payload);
  const landlordAddress = normalizeWalletAddress(await contract.ownerOf(tokenId));
  const renterAddress = String(payload?.renterAddress ?? "").trim();
  if (!isHexAddress(renterAddress)) {
    throw new Error("renterAddress must be a valid EVM address.");
  }
  if (normalizeWalletAddress(renterAddress) !== actorAddress) {
    throw new Error("renterAddress must match the authenticated wallet address.");
  }

  const durationDays = parseDurationDays(payload?.durationDays);
  const monthlyRent = parsePositiveBigInt(payload?.monthlyRent, "monthlyRent");
  const depositAmount =
    payload?.depositAmount === undefined || String(payload?.depositAmount).trim() === ""
      ? monthlyRent
      : parsePositiveBigInt(payload?.depositAmount, "depositAmount");

  await ensureKYCFromPayload(runtime, payload, actorAddress);
  const encryptedAccessKey = randomBytes(32);
  const accessKeyHash = ethers.keccak256(encryptedAccessKey);
  const houseDetails = await contract.getHouseDetails(tokenId).catch(() => null);
  const houseReference = formatHouseReference(
    houseDetails ? { houseId: houseDetails.houseId ?? houseDetails[0] ?? "" } : null,
    tokenId.toString(),
  );

  const { receipt } = await submitTransaction(
    contract.startRental(
      tokenId,
      renterAddress,
      durationDays,
      depositAmount,
      monthlyRent,
      encryptedAccessKey,
    ),
    "startRental",
  );

  notifyWorkflowParticipants({
    walletAddresses: [landlordAddress],
    type: "transaction_confirmed",
    title: "Rental started",
    message: `${houseReference} is now rented to ${renterAddress.slice(0, 6)}...${renterAddress.slice(-4)}.`,
    data: {
      action: "rent",
      tokenId: tokenId.toString(),
      renterAddress,
      txHash: receipt.hash,
    },
  });
  notifyWorkflowParticipants({
    walletAddresses: [renterAddress],
    type: "transaction_confirmed",
    title: "Rental confirmed",
    message: `Your rental for ${houseReference} is active.`,
    data: {
      action: "rent",
      tokenId: tokenId.toString(),
      landlordAddress,
      txHash: receipt.hash,
    },
  });

  return {
    success: true,
    message: "rental started successfully",
    txHash: receipt.hash,
    encryptedKey: encryptedAccessKey.toString("base64"),
    data: {
      tokenID: tokenId.toString(),
      renter: renterAddress,
      durationDays: Number(durationDays),
      monthlyRent: monthlyRent.toString(),
      depositAmount: depositAmount.toString(),
      accessKeyHash,
      kycProvider: String(payload?.kycProvider ?? "mock"),
    },
  };
}

async function handleCreateBill(payload) {
  const runtime = await getRuntime();
  const { contract } = runtime;
  const actorAddress = requireAuthenticatedActor(
    payload,
    "Authentication is required to create bills.",
  );

  const tokenId = parseTokenId(payload);
  const billType = String(payload?.billType ?? "").trim().toLowerCase();
  if (!billType) {
    throw new Error("billType is required.");
  }

  const amount = parseBillAmountInCents(payload?.amount);
  const { dueDate, dueUnix } = parseDueDate(payload?.dueDate);

  const provider = String(payload?.provider ?? "").trim();
  if (!isHexAddress(provider)) {
    throw new Error("provider must be a valid EVM address.");
  }
  const onchainOwner = normalizeWalletAddress(await contract.ownerOf(tokenId));
  if (onchainOwner !== actorAddress) {
    throw new Error("Only the current onchain property owner can create bills.");
  }

  const isRecurring = Boolean(payload?.isRecurring);
  const recurrenceInterval = isRecurring
    ? parsePositiveInteger(payload?.recurrenceInterval, 30)
    : 0;

  let billIndex = "-1";
  try {
    const preCount = await contract.getTotalBillsCount(tokenId);
    billIndex = preCount.toString();
  } catch {
    billIndex = "-1";
  }

  const { receipt } = await submitTransaction(
    contract.createBill(
      tokenId,
      billType,
      amount,
      dueUnix,
      provider,
      isRecurring,
      recurrenceInterval,
    ),
    "createBill",
  );

  return {
    success: true,
    message: "bill created successfully",
    txHash: receipt.hash,
    data: {
      tokenID: tokenId.toString(),
      billType,
      amount: Number(amount),
      dueDate: dueDate.toISOString(),
      isRecurring,
      recurrenceInterval,
      billIndex,
    },
  };
}

async function handlePayBill(payload) {
  const runtime = await getRuntime();
  const { ethers, contract } = runtime;

  const tokenId = parseTokenId(payload);
  const payerAddress = requireAuthenticatedActor(
    payload,
    "Authentication is required to pay bills.",
  );
  const billIndex = parseNonNegativeInteger(payload?.billIndex, "billIndex");

  const method = String(payload?.paymentMethod ?? "").trim().toLowerCase();
  if (!PAYMENT_METHODS.has(method)) {
    throw new Error("paymentMethod must be crypto, stripe, or bank_transfer.");
  }

  const paymentReference = ethers.keccak256(
    ethers.toUtf8Bytes(
      `bill:${tokenId.toString()}:${billIndex}:${method}:${Date.now()}`,
    ),
  );
  const ownerAddress = normalizeWalletAddress(await contract.ownerOf(tokenId));
  const houseDetails = await contract.getHouseDetails(tokenId).catch(() => null);
  const houseReference = formatHouseReference(
    houseDetails ? { houseId: houseDetails.houseId ?? houseDetails[0] ?? "" } : null,
    tokenId.toString(),
  );

  const { receipt } = await submitTransaction(
    contract.recordBillPayment(
      tokenId,
      BigInt(billIndex),
      method,
      paymentReference,
    ),
    "recordBillPayment",
  );

  const notificationWallets = [ownerAddress];
  if (payerAddress) {
    notificationWallets.push(payerAddress);
  }
  notifyWorkflowParticipants({
    walletAddresses: notificationWallets,
    type: "bill_paid",
    title: "Bill paid",
    message: `Bill #${billIndex} for ${houseReference} was paid via ${method}.`,
    data: {
      action: "pay_bill",
      tokenId: tokenId.toString(),
      billIndex,
      paymentMethod: method,
      paymentReference,
      txHash: receipt.hash,
      payerAddress: payerAddress || undefined,
    },
  });

  return {
    success: true,
    message: `payment processed via ${method}`,
    txHash: receipt.hash,
    data: {
      tokenID: tokenId.toString(),
      billIndex,
      paymentMethod: method,
      paymentReference,
    },
  };
}

async function handleSetKYC(payload) {
  const runtime = await getRuntime();
  const actorAddress = requireAuthenticatedActor(
    payload,
    "Authentication is required to update KYC status.",
  );
  const provider = String(payload?.kycProvider ?? "mock").trim().toLowerCase();
  const normalizedProvider =
    provider === "none" || provider === "zkpassport" ? provider : "mock";
  const walletAddressRaw = String(
    payload?.walletAddress ?? payload?.userAddress ?? payload?.address ?? actorAddress,
  ).trim();
  if (!isHexAddress(walletAddressRaw)) {
    throw new Error("walletAddress must be a valid EVM address.");
  }
  const walletAddress = normalizeWalletAddress(walletAddressRaw);
  if (walletAddress !== actorAddress) {
    throw new Error("walletAddress must match the authenticated wallet address.");
  }

  await ensureKYCFromPayload(runtime, payload, walletAddress);

  return {
    success: true,
    message:
      normalizedProvider === "none"
        ? "kyc skipped for anonymous mode"
        : "kyc recorded successfully",
    data: {
      walletAddress,
      kycProvider: normalizedProvider,
    },
  };
}

async function handleClaimKey(payload) {
  const runtime = await getRuntime();
  const { ethers, contract } = runtime;

  const keyHash = String(payload?.keyHash ?? "").trim();
  if (!isBytes32Hex(keyHash)) {
    throw new Error("keyHash must be a bytes32 hex value.");
  }

  const claimant = String(payload?.claimant ?? "").trim();
  if (claimant && !isHexAddress(claimant)) {
    throw new Error("claimant must be a valid EVM address.");
  }

  const exchange = await contract.keyExchanges(keyHash);
  const encryptedKeyRaw = exchange?.encryptedKey ?? exchange?.[1];
  const intendedRecipient = exchange?.intendedRecipient ?? exchange?.[2];
  const createdAtRaw = exchange?.createdAt ?? exchange?.[3] ?? 0n;
  const expiresAtRaw = exchange?.expiresAt ?? exchange?.[4] ?? 0n;
  const isClaimed = exchange?.isClaimed ?? exchange?.[5] ?? false;
  const exchangeType = exchange?.exchangeType ?? exchange?.[6] ?? 0;

  const encryptedKeyBytes = ethers.getBytes(encryptedKeyRaw ?? "0x");
  if (encryptedKeyBytes.length === 0) {
    throw new Error("key exchange not found.");
  }

  if (claimant) {
    const normalizedClaimant = ethers.getAddress(claimant);
    const normalizedRecipient = ethers.getAddress(intendedRecipient);
    if (normalizedClaimant !== normalizedRecipient) {
      throw new Error("claimant is not intended recipient.");
    }
  }

  const encryptedKey = encodeBase64(encryptedKeyRaw, ethers);

  return {
    success: true,
    message: "key fetched successfully",
    encryptedKey,
    data: {
      keyHash,
      encryptedKey,
      intendedRecipient,
      isClaimed: Boolean(isClaimed),
      createdAt: Number(createdAtRaw),
      expiresAt: Number(expiresAtRaw),
      exchangeType: Number(exchangeType),
    },
  };
}

async function routeWorkflowAction(payload) {
  const action = String(payload?.action ?? "").trim().toLowerCase();
  if (!action) {
    throw new Error("action is required.");
  }

  switch (action) {
    case "mint":
      return handleMint(payload);
    case "update_house_images":
      return handleUpdateHouseImages(payload);
    case "set_kyc":
      return handleSetKYC(payload);
    case "create_listing":
      return handleCreateListing(payload);
    case "sell":
      return handleSell(payload);
    case "rent":
      return handleRent(payload);
    case "create_bill":
      return handleCreateBill(payload);
    case "pay_bill":
      return handlePayBill(payload);
    case "claim_key":
      return handleClaimKey(payload);
    default:
      throw new Error(`unknown action: ${action}`);
  }
}

async function handleWorkflowTrigger(request, response, origin, writeJson) {
  let body;
  try {
    body = parseJsonBody(await readBody(request));
  } catch (error) {
    writeJson(
      response,
      400,
      { success: false, message: "Invalid JSON payload." },
      origin,
    );
    return;
  }

  try {
    const action = String(body?.action ?? "").trim().toLowerCase();
    const viewerAddress = extractWalletAddressFromAuthHeader(
      request.headers.authorization,
    );

    if (AUTH_REQUIRED_WORKFLOW_ACTIONS.has(action) && !viewerAddress) {
      writeJson(
        response,
        401,
        {
          success: false,
          message: "Authentication is required for this workflow action.",
        },
        origin,
      );
      return;
    }

    if (viewerAddress) {
      const actorAddressRaw = String(body?.actorAddress ?? "").trim();
      if (actorAddressRaw) {
        const normalizedActorAddress = normalizeWalletAddress(actorAddressRaw);
        if (!isHexAddress(normalizedActorAddress)) {
          writeJson(
            response,
            400,
            {
              success: false,
              message: "actorAddress must be a valid EVM address.",
            },
            origin,
          );
          return;
        }

        if (normalizedActorAddress !== viewerAddress) {
          writeJson(
            response,
            403,
            {
              success: false,
              message: "actorAddress must match the authenticated wallet address.",
            },
            origin,
          );
          return;
        }
      }

      body = {
        ...body,
        actorAddress: viewerAddress,
      };
    }

    if (action === "claim_key") {
      const claimedBy = String(body?.claimant ?? "").trim().toLowerCase();
      if (claimedBy && claimedBy !== viewerAddress) {
        writeJson(
          response,
          403,
          {
            success: false,
            message:
              "claimant must match the authenticated wallet address.",
          },
          origin,
        );
        return;
      }

      body = {
        ...body,
        claimant: viewerAddress,
      };
    }

    const result = await routeWorkflowAction(body);
    writeJson(response, 200, result, origin);
  } catch (error) {
    const message = error instanceof Error ? error.message : "Workflow error.";
    const isServerIssue = /WORKFLOW_|HOUSE_RWA_CONTRACT_ADDRESS|Ethers library|RPC chain mismatch|No contract code/.test(
      message,
    );
    writeJson(
      response,
      isServerIssue ? 500 : 400,
      { success: false, message },
      origin,
    );
  }
}

function parseVerifyRequestPayload(body) {
  const walletAddress = String(
    body?.walletAddress ?? body?.address ?? "",
  ).trim();

  const provider = String(body?.provider ?? "zkpassport")
    .trim()
    .toLowerCase();
  if (provider !== "zkpassport") {
    throw new Error(`unsupported provider: ${provider}`);
  }

  let proofBundle;
  if (body?.proof && typeof body.proof === "object") {
    proofBundle = body.proof;
  } else {
    proofBundle = {
      proofs: body?.proofs,
      queryResult: body?.queryResult ?? body?.result,
      scope: body?.scope,
      devMode: body?.devMode,
      validity: body?.validity,
    };
  }

  return { walletAddress, proofBundle };
}

async function handleVerifyKYC(request, response, origin, writeJson) {
  let body;
  try {
    body = parseJsonBody(await readBody(request));
  } catch {
    writeJson(
      response,
      400,
      { success: false, message: "Invalid JSON payload." },
      origin,
    );
    return;
  }

  try {
    const { walletAddress, proofBundle } = parseVerifyRequestPayload(body);
    const verification = await verifyZKPassportProof(walletAddress, proofBundle);
    writeJson(
      response,
      200,
      {
        success: true,
        message: "ZKPassport proof verified.",
        data: verification,
      },
      origin,
    );
  } catch (error) {
    const message =
      error instanceof Error
        ? error.message
        : "Unable to verify ZKPassport proof.";
    writeJson(
      response,
      400,
      { success: false, message },
      origin,
    );
  }
}

async function handleVerifyKYCForCRE(request, response, origin, writeJson) {
  let body;
  try {
    body = parseJsonBody(await readBody(request));
  } catch {
    writeJson(
      response,
      400,
      { verified: false, message: "Invalid JSON payload." },
      origin,
    );
    return;
  }

  try {
    const { walletAddress, proofBundle } = parseVerifyRequestPayload(body);
    const verification = await verifyZKPassportProof(walletAddress, proofBundle);
    writeJson(
      response,
      200,
      {
        verified: true,
        level: verification.level,
        verificationHash: verification.verificationHash,
        expiresAt: verification.expiresAt,
        message: "proof verified",
      },
      origin,
    );
  } catch (error) {
    const message =
      error instanceof Error
        ? error.message
        : "Unable to verify ZKPassport proof.";
    writeJson(
      response,
      422,
      {
        verified: false,
        message,
      },
      origin,
    );
  }
}

async function handleVerifyWallet(request, response, origin, writeJson) {
  let body;
  try {
    body = parseJsonBody(await readBody(request));
  } catch {
    writeJson(
      response,
      400,
      { success: false, message: "Invalid JSON payload." },
      origin,
    );
    return;
  }

  const walletAddress = String(body?.address ?? "").trim();
  if (!isHexAddress(walletAddress)) {
    writeJson(
      response,
      400,
      { success: false, message: "address must be a valid EVM address." },
      origin,
    );
    return;
  }

  const signature = String(body?.signature ?? "").trim();
  const message = String(body?.message ?? "").trim();
  if (!signature || !message) {
    writeJson(
      response,
      400,
      {
        success: false,
        message: "signature and message are required for wallet verification.",
      },
      origin,
    );
    return;
  }

  let recoveredAddress = "";
  try {
    const ethers = loadEthers();
    recoveredAddress = String(ethers.verifyMessage(message, signature));
  } catch (error) {
    writeJson(
      response,
      400,
      {
        success: false,
        message: `Unable to verify signature: ${
          error instanceof Error ? error.message : String(error)
        }`,
      },
      origin,
    );
    return;
  }

  if (
    !isHexAddress(recoveredAddress)
    || normalizeWalletAddress(recoveredAddress) !== normalizeWalletAddress(walletAddress)
  ) {
    writeJson(
      response,
      401,
      { success: false, message: "Signature does not match wallet address." },
      origin,
    );
    return;
  }

  const token = createAuthToken(walletAddress);
  const user = await buildAuthenticatedUser(walletAddress);

  writeJson(
    response,
    200,
    {
      success: true,
      message: "Wallet verified.",
      data: { token, user },
    },
    origin,
  );
}

function handleRefreshToken(request, response, origin, writeJson) {
  const walletAddress = extractWalletAddressFromAuthHeader(
    request.headers.authorization,
  );
  if (!walletAddress) {
    writeJson(
      response,
      401,
      {
        success: false,
        message: "Authentication is required to refresh the session token.",
      },
      origin,
    );
    return;
  }

  const token = createAuthToken(walletAddress);
  writeJson(
    response,
    200,
    { success: true, message: "Token refreshed.", data: { token } },
    origin,
  );
}

function handleLogout(response, origin, writeJson) {
  writeJson(response, 200, { success: true, message: "Logged out." }, origin);
}

async function getAuthenticatedUserProfile(walletAddress) {
  return buildAuthenticatedUser(walletAddress);
}

function getNotificationsForWallet(walletAddress, limit) {
  return listNotificationsForWallet(walletAddress, limit);
}

function markWalletNotificationRead(walletAddress, notificationId) {
  return markNotificationRead(walletAddress, notificationId);
}

function getRoleGatedConversationsForWallet(walletAddress, tokenId) {
  return listRoleGatedConversations(walletAddress, tokenId);
}

function getRoleGatedConversationForWallet(walletAddress, conversationId) {
  return getRoleGatedConversationDetails(walletAddress, conversationId);
}

async function sendRoleGatedWalletMessage(payload) {
  return sendRoleGatedMessage(payload);
}

module.exports = {
  handleWorkflowTrigger,
  handleVerifyWallet,
  handleVerifyKYC,
  handleVerifyKYCForCRE,
  handleRefreshToken,
  handleLogout,
  readHousesFromChain,
  readNativeBalance,
  readKYCStatusForWallet,
  getPrivateDocumentBundleForToken,
  extractWalletAddressFromAuthHeader,
  getAuthenticatedUserProfile,
  getNotificationsForWallet,
  markWalletNotificationRead,
  getRoleGatedConversationsForWallet,
  getRoleGatedConversationForWallet,
  sendRoleGatedWalletMessage,
};
