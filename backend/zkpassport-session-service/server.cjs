const { randomUUID } = require("node:crypto");
const { createServer } = require("node:http");
const { URL } = require("node:url");

const { ZKPassport } = require("@zkpassport/sdk");
const {
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
} = require("./workflow-trigger.cjs");

const PORT = parsePositiveInteger(process.env.PORT, 8787);
const HOST = process.env.HOST ?? "0.0.0.0";
const CORS_ORIGIN = process.env.CORS_ORIGIN ?? "*";
const UPSTREAM_API_URL = process.env.UPSTREAM_API_URL ?? "";
const WORKFLOW_RPC_URL =
  process.env.WORKFLOW_RPC_URL ??
  process.env.CRE_RPC_URL ??
  process.env.SEPOLIA_RPC ??
  "";

const SESSION_TTL_SECONDS = parsePositiveInteger(
  process.env.ZKPASSPORT_SESSION_TTL_SECONDS,
  900,
);
const TERMINAL_RETENTION_SECONDS = parsePositiveInteger(
  process.env.ZKPASSPORT_TERMINAL_RETENTION_SECONDS,
  3600,
);

const ZKPASSPORT_DOMAIN =
  process.env.ZKPASSPORT_DOMAIN ??
  process.env.RAILWAY_PUBLIC_DOMAIN ??
  "demo.zkpassport.id";
const ZKPASSPORT_APP_NAME = process.env.ZKPASSPORT_APP_NAME ?? "PropMeCRE";
const ZKPASSPORT_APP_LOGO =
  process.env.ZKPASSPORT_APP_LOGO ??
  "https://zkpassport.id/logo.png";
const ZKPASSPORT_PURPOSE =
  process.env.ZKPASSPORT_PURPOSE ??
  "Verify KYC for secure real-estate workflow access.";
const ZKPASSPORT_SCOPE = process.env.ZKPASSPORT_SCOPE ?? "rwa-house-kyc";
const ZKPASSPORT_CHAIN = process.env.ZKPASSPORT_CHAIN ?? "ethereum_sepolia";
const ZKPASSPORT_BIND_CHAIN = parseBoolean(
  process.env.ZKPASSPORT_BIND_CHAIN,
  false,
);
const ZKPASSPORT_BIND_USER_ADDRESS = parseBoolean(
  process.env.ZKPASSPORT_BIND_USER_ADDRESS,
  true,
);
const ZKPASSPORT_MODE = process.env.ZKPASSPORT_MODE ?? "fast";
const ZKPASSPORT_PROJECT_ID = process.env.ZKPASSPORT_PROJECT_ID ?? "";
const ZKPASSPORT_VALIDITY_SECONDS = parsePositiveInteger(
  process.env.ZKPASSPORT_VALIDITY_SECONDS,
  604800,
);
const ZKPASSPORT_DEV_MODE = parseBoolean(process.env.ZKPASSPORT_DEV_MODE, false);
const ZKPASSPORT_BRIDGE_URL = process.env.ZKPASSPORT_BRIDGE_URL ?? "";
const ZKPASSPORT_CLOUD_PROVER_URL =
  process.env.ZKPASSPORT_CLOUD_PROVER_URL ?? "";
const REQUIRE_KYC_FOR_PRIVATE_DOCUMENTS = parseBoolean(
  process.env.REQUIRE_KYC_FOR_PRIVATE_DOCUMENTS,
  false,
);

const CLAIMS = parseClaims(process.env.ZKPASSPORT_CLAIMS_JSON);
const sessions = new Map();

const SUPPORTED_PROOF_MODES = new Set(["fast", "compressed", "compressed-evm"]);

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

function isTerminalStatus(status) {
  return status === "verified" || status === "failed" || status === "expired";
}

function isLoopbackOrigin(origin) {
  try {
    const parsed = new URL(origin);
    const host = parsed.hostname.toLowerCase();
    return host === "localhost" || host === "127.0.0.1" || host === "::1";
  } catch {
    return false;
  }
}

function resolveCorsOrigin(request) {
  const requestOrigin =
    typeof request.headers.origin === "string" ? request.headers.origin : "";
  if (CORS_ORIGIN === "*") {
    return requestOrigin || "*";
  }

  const allowedOrigins = CORS_ORIGIN.split(",")
    .map((entry) => entry.trim())
    .filter(Boolean);

  if (allowedOrigins.length === 0) {
    return requestOrigin || "*";
  }

  if (allowedOrigins.includes("*")) {
    return requestOrigin || "*";
  }

  if (!requestOrigin) {
    return allowedOrigins[0] || "*";
  }

  if (allowedOrigins.includes(requestOrigin)) {
    return requestOrigin;
  }

  if (
    isLoopbackOrigin(requestOrigin)
    && allowedOrigins.some((origin) => isLoopbackOrigin(origin))
  ) {
    return requestOrigin;
  }

  return null;
}

function createResponseHeaders(origin, extraHeaders = {}) {
  const allowCredentials = origin !== "*";
  return {
    "Content-Type": "application/json; charset=utf-8",
    "Cache-Control": "no-store, no-cache, must-revalidate, proxy-revalidate",
    Pragma: "no-cache",
    Expires: "0",
    "Access-Control-Allow-Origin": origin,
    "Access-Control-Allow-Methods": "GET,POST,DELETE,OPTIONS",
    "Access-Control-Allow-Headers":
      "Content-Type,Authorization,X-Request-ID,X-Timestamp,X-Signature",
    ...(allowCredentials
      ? {
          "Access-Control-Allow-Credentials": "true",
          Vary: "Origin",
        }
      : {}),
    "Access-Control-Max-Age": "86400",
    ...extraHeaders,
  };
}

function writeJson(response, statusCode, payload, origin, extraHeaders = {}) {
  const body = JSON.stringify(payload);
  response.writeHead(statusCode, createResponseHeaders(origin, extraHeaders));
  response.end(body);
}

function buildQrCodeUrl(targetUrl) {
  const encodedUrl = encodeURIComponent(targetUrl);
  return `https://api.qrserver.com/v1/create-qr-code/?size=320x320&data=${encodedUrl}`;
}

function readBody(request) {
  return new Promise((resolve, reject) => {
    let body = "";
    request.on("data", (chunk) => {
      body += chunk.toString("utf8");
      if (body.length > 1_000_000) {
        reject(new Error("request payload too large"));
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

function isValidWalletAddress(walletAddress) {
  return /^0x[a-fA-F0-9]{40}$/.test(walletAddress ?? "");
}

function normalizeWalletAddress(walletAddress) {
  return String(walletAddress ?? "").trim().toLowerCase();
}

function normalizeDomainCandidate(rawDomain) {
  const candidate = String(rawDomain ?? "").trim();
  if (!candidate) {
    return "";
  }

  const withProtocol = /^[a-zA-Z][a-zA-Z\d+\-.]*:\/\//.test(candidate)
    ? candidate
    : `https://${candidate}`;
  try {
    const parsed = new URL(withProtocol);
    return parsed.hostname.trim().toLowerCase();
  } catch {
    return "";
  }
}

function isLoopbackHostname(hostname) {
  const normalized = String(hostname ?? "").trim().toLowerCase();
  return normalized === "localhost" || normalized === "127.0.0.1" || normalized === "::1";
}

function resolveSessionDomain(rawRequestedDomain) {
  const configuredDomain = normalizeDomainCandidate(ZKPASSPORT_DOMAIN);
  const requestedDomain = normalizeDomainCandidate(rawRequestedDomain);

  if (!requestedDomain) {
    return configuredDomain || ZKPASSPORT_DOMAIN;
  }

  if (!configuredDomain) {
    return requestedDomain;
  }

  if (requestedDomain === configuredDomain) {
    return requestedDomain;
  }

  if (isLoopbackHostname(requestedDomain)) {
    if (isLoopbackHostname(configuredDomain) || ZKPASSPORT_DEV_MODE) {
      return requestedDomain;
    }
    return configuredDomain;
  }

  const allowDynamicDomain = parseBoolean(
    process.env.ZKPASSPORT_ALLOW_DYNAMIC_DOMAIN,
    true,
  );
  if (allowDynamicDomain || ZKPASSPORT_DEV_MODE) {
    return requestedDomain;
  }

  return configuredDomain;
}

function normalizeProofMode(rawMode, fallbackMode = ZKPASSPORT_MODE) {
  const candidate = String(rawMode ?? "").trim().toLowerCase();
  if (SUPPORTED_PROOF_MODES.has(candidate)) {
    return candidate;
  }
  const normalizedFallback = String(fallbackMode ?? "").trim().toLowerCase();
  if (SUPPORTED_PROOF_MODES.has(normalizedFallback)) {
    return normalizedFallback;
  }
  return "fast";
}

function extractViewerWalletAddress(request) {
  return extractWalletAddressFromAuthHeader(request.headers.authorization);
}

function buildRedactedMetadata(tokenId) {
  return {
    address: `Private Asset #${tokenId}`,
    city: "Private",
    state: "Private",
    zipCode: "Private",
    country: "Private",
    propertyType: "single_family",
    bedrooms: 0,
    bathrooms: 0,
    squareFeet: 0,
    yearBuilt: 0,
    description: "Sensitive property details are visible only to authorized parties.",
    images: [],
  };
}

function isViewerAuthorizedForHouse(house, viewerWalletAddress) {
  if (!viewerWalletAddress || !house || typeof house !== "object") {
    return false;
  }

  const normalizedViewer = normalizeWalletAddress(viewerWalletAddress);
  const allowed = new Set();
  const addAddress = (candidate) => {
    if (isValidWalletAddress(candidate)) {
      allowed.add(normalizeWalletAddress(candidate));
    }
  };

  addAddress(house.ownerAddress);
  addAddress(house.originalOwner);
  addAddress(house?.rental?.renterAddress);
  if (house?.rental?.isActive) {
    addAddress(house?.rental?.renterAddress);
  }
  if (house?.listing?.isPrivateSale) {
    addAddress(house?.listing?.allowedBuyer);
  }

  return allowed.has(normalizedViewer);
}

function redactHouseForUnauthorizedViewer(house) {
  const tokenId = String(house?.tokenId ?? "unknown");
  return {
    ...house,
    houseId: `Private-Asset-${tokenId}`,
    documentHash: "",
    documentURI: "",
    metadata: buildRedactedMetadata(tokenId),
    bills: [],
    rental: house?.rental
      ? {
          ...house.rental,
          renterAddress: "",
          hasAccessKey: false,
        }
      : undefined,
  };
}

function projectHouseForViewer(house, viewerWalletAddress) {
  if (isViewerAuthorizedForHouse(house, viewerWalletAddress)) {
    return house;
  }
  return redactHouseForUnauthorizedViewer(house);
}

function isViewerCurrentOwner(house, viewerWalletAddress) {
  if (!house || !viewerWalletAddress) {
    return false;
  }
  const ownerAddress = normalizeWalletAddress(house.ownerAddress);
  const viewer = normalizeWalletAddress(viewerWalletAddress);
  if (!isValidWalletAddress(ownerAddress) || !isValidWalletAddress(viewer)) {
    return false;
  }
  return ownerAddress === viewer;
}

function parsePrivateDocumentBundle(documentsB64) {
  const encoded = String(documentsB64 ?? "").trim();
  if (!encoded) {
    return [];
  }

  const decoded = Buffer.from(encoded, "base64").toString("utf8");
  const parsed = JSON.parse(decoded);
  const files = Array.isArray(parsed?.files) ? parsed.files : [];
  const metadata = Array.isArray(parsed?.metadata) ? parsed.metadata : [];

  const documents = [];
  for (let index = 0; index < files.length; index += 1) {
    const base64 = typeof files[index] === "string" ? files[index].trim() : "";
    if (!base64) {
      continue;
    }

    const meta = metadata[index] && typeof metadata[index] === "object"
      ? metadata[index]
      : {};
    const mimeType =
      typeof meta.type === "string" && meta.type.trim()
        ? meta.type.trim()
        : "application/octet-stream";
    const name =
      typeof meta.name === "string" && meta.name.trim()
        ? meta.name.trim()
        : `Document ${index + 1}`;
    const declaredSize = Number(meta.size);
    const size = Number.isFinite(declaredSize) && declaredSize >= 0
      ? declaredSize
      : 0;

    documents.push({
      index,
      name,
      mimeType,
      size,
      base64,
    });
  }

  return documents;
}

function parseClaims(rawClaims) {
  const fallbackClaims = [
    { type: "disclose", field: "birthdate" },
  ];

  if (!rawClaims) {
    return fallbackClaims;
  }

  try {
    const parsed = JSON.parse(rawClaims);
    if (!Array.isArray(parsed)) {
      throw new Error("ZKPASSPORT_CLAIMS_JSON must be an array");
    }
    return parsed;
  } catch (error) {
    console.warn(
      "[zkpassport-session-service] Invalid ZKPASSPORT_CLAIMS_JSON. Falling back to defaults.",
      error,
    );
    return fallbackClaims;
  }
}

function applyClaims(queryBuilder, claims) {
  let builder = queryBuilder;

  for (const claim of claims) {
    const type = String(claim?.type ?? "").toLowerCase();
    const field = String(claim?.field ?? "").trim();
    if (!field) {
      continue;
    }

    if (type === "disclose" && typeof builder.disclose === "function") {
      builder = builder.disclose(field);
      continue;
    }

    if (type === "gte" && typeof builder.gte === "function") {
      builder = builder.gte(field, Number(claim?.value ?? 0));
      continue;
    }

    if (type === "in" && typeof builder.in === "function") {
      builder = builder.in(field, Array.isArray(claim?.values) ? claim.values : []);
      continue;
    }

    if (type === "out" && typeof builder.out === "function") {
      builder = builder.out(
        field,
        Array.isArray(claim?.values) ? claim.values : [],
      );
    }
  }

  return builder;
}

function cancelSessionRequest(session) {
  if (
    session?.zkPassport
    && typeof session.zkPassport.cancelRequest === "function"
    && session?.sessionId
  ) {
    try {
      session.zkPassport.cancelRequest(session.sessionId);
    } catch {
      // Best effort cleanup only.
    }
  }
}

function cleanupSessions() {
  const now = Date.now();
  const retentionMs = TERMINAL_RETENTION_SECONDS * 1000;

  for (const [sessionId, session] of sessions.entries()) {
    const expiresAt = Date.parse(session.expiresAt);
    if (!Number.isFinite(expiresAt)) {
      continue;
    }

    if (now > expiresAt && !isTerminalStatus(session.status)) {
      session.status = "expired";
      session.message = "Session expired. Start a new ZKPassport verification.";
      sessions.set(sessionId, session);
    }

    if (now > expiresAt + retentionMs && isTerminalStatus(session.status)) {
      cancelSessionRequest(session);
      sessions.delete(sessionId);
    }
  }
}

function closeOtherActiveSessionsForWallet(walletAddress) {
  const normalizedWalletAddress = normalizeWalletAddress(walletAddress);
  for (const [sessionId, session] of sessions.entries()) {
    if (normalizeWalletAddress(session.walletAddress) !== normalizedWalletAddress) {
      continue;
    }
    if (isTerminalStatus(session.status)) {
      continue;
    }
    cancelSessionRequest(session);
    sessions.delete(sessionId);
  }
}

function toPublicSession(session) {
  const bridgeConnected =
    typeof session?.requestHandle?.isBridgeConnected === "function"
      ? Boolean(session.requestHandle.isBridgeConnected())
      : Boolean(session.bridgeConnected);
  const requestReceived =
    typeof session?.requestHandle?.requestReceived === "function"
      ? Boolean(session.requestHandle.requestReceived())
      : Boolean(session.requestReceived);

  return {
    sessionId: session.sessionId,
    status: session.status,
    domain: session.domain,
    mode: session.mode,
    qrCodeUrl: session.qrCodeUrl,
    deepLinkUrl: session.deepLinkUrl,
    bridgeConnected,
    requestReceived,
    devMode: Boolean(session.devMode),
    expiresAt: session.expiresAt,
    proof: session.proof,
    message: session.message,
  };
}

function updateSession(sessionId, updater) {
  const current = sessions.get(sessionId);
  if (!current) {
    return;
  }

  const next = updater(current);
  sessions.set(sessionId, next);
}

function normalizeProofEntry(value) {
  if (!value || typeof value !== "object") {
    return null;
  }

  if (value.proof && typeof value.proof === "object") {
    return value.proof;
  }

  return value;
}

function collectSessionProofs(proofEvents) {
  if (!Array.isArray(proofEvents)) {
    return [];
  }

  return proofEvents
    .map((event) => normalizeProofEntry(event))
    .filter((event) => Boolean(event));
}

function collectResultProofs(resultPayload) {
  const proofs = [];

  if (Array.isArray(resultPayload?.proofs)) {
    proofs.push(...resultPayload.proofs);
  }

  if (Array.isArray(resultPayload?.result?.proofs)) {
    proofs.push(...resultPayload.result.proofs);
  }

  if (resultPayload?.proof && typeof resultPayload.proof === "object") {
    proofs.push(resultPayload.proof);
  }

  return proofs
    .map((proof) => normalizeProofEntry(proof))
    .filter((proof) => Boolean(proof));
}

function resolveResultQuery(resultPayload) {
  if (resultPayload?.result && typeof resultPayload.result === "object") {
    return resultPayload.result;
  }

  if (resultPayload?.queryResult && typeof resultPayload.queryResult === "object") {
    return resultPayload.queryResult;
  }

  return {};
}

function createProofBundle(session, resultPayload) {
  const sessionProofs = collectSessionProofs(session.proofEvents);
  const resultProofs = collectResultProofs(resultPayload);
  const proofs = sessionProofs.length > 0 ? sessionProofs : resultProofs;

  return {
    proofs,
    queryResult: resolveResultQuery(resultPayload),
    verified: Boolean(resultPayload?.verified),
    uniqueIdentifier: resultPayload?.uniqueIdentifier ?? "",
    walletAddress: session.walletAddress,
    domain: session.domain,
    scope: ZKPASSPORT_SCOPE,
  };
}

function registerCallback(requestHandle, callbackName, handler) {
  const candidate = requestHandle?.[callbackName];
  if (typeof candidate === "function") {
    candidate(handler);
  }
}

function toErrorMessage(error, fallback) {
  if (error instanceof Error && error.message) {
    return error.message;
  }
  if (typeof error === "string" && error.trim()) {
    return error;
  }
  return fallback;
}

function logSessionEvent(sessionId, event, details = {}) {
  try {
    console.log(
      `[zkpassport-session] ${event} session=${sessionId} details=${JSON.stringify(details)}`,
    );
  } catch {
    console.log(`[zkpassport-session] ${event} session=${sessionId}`);
  }
}

function bindSessionCallbacks(requestHandle, sessionId) {
  registerCallback(requestHandle, "onBridgeConnect", () => {
    logSessionEvent(sessionId, "bridge_connected");
    updateSession(sessionId, (session) => ({
      ...session,
      bridgeConnected: true,
      message: "Bridge connected. Scan the QR code in the ZKPassport app.",
    }));
  });

  registerCallback(requestHandle, "onRequestReceived", () => {
    logSessionEvent(sessionId, "request_received");
    updateSession(sessionId, (session) => ({
      ...session,
      requestReceived: true,
      status: isTerminalStatus(session.status) ? session.status : "ready",
      message: "QR code scanned. Waiting for user confirmation.",
    }));
  });

  registerCallback(requestHandle, "onGeneratingProof", () => {
    logSessionEvent(sessionId, "generating_proof");
    updateSession(sessionId, (session) => ({
      ...session,
      status: isTerminalStatus(session.status) ? session.status : "ready",
      message: "Generating ZK proof. Please keep the app open.",
    }));
  });

  registerCallback(requestHandle, "onProofGenerated", (proofEvent) => {
    logSessionEvent(sessionId, "proof_generated", {
      name: proofEvent?.name ?? "",
      version: proofEvent?.version ?? "",
      total: Number(proofEvent?.total ?? 0),
    });
    updateSession(sessionId, (session) => ({
      ...session,
      proofEvents: [
        ...session.proofEvents,
        proofEvent,
      ],
      message: "Proof generated. Awaiting final verification result.",
    }));
  });

  registerCallback(requestHandle, "onResult", (resultPayload) => {
    logSessionEvent(sessionId, "result", {
      verified: Boolean(resultPayload?.verified),
      hasResult: Boolean(resultPayload?.result),
      hasProofs: Array.isArray(resultPayload?.proofs),
      uniqueIdentifier: String(resultPayload?.uniqueIdentifier ?? ""),
    });
    updateSession(sessionId, (session) => {
      const verified = Boolean(resultPayload?.verified);
      return {
        ...session,
        status: verified ? "verified" : "failed",
        proof: createProofBundle(session, resultPayload),
        message: verified
          ? "ZKPassport verification completed."
          : "ZKPassport verification failed.",
      };
    });
  });

  registerCallback(requestHandle, "onReject", (reason) => {
    logSessionEvent(sessionId, "rejected", {
      reason: toErrorMessage(reason, "User rejected verification request."),
    });
    updateSession(sessionId, (session) => ({
      ...session,
      status: "failed",
      message: `Verification rejected: ${toErrorMessage(
        reason,
        "User rejected verification request.",
      )}`,
    }));
  });

  registerCallback(requestHandle, "onError", (error) => {
    logSessionEvent(sessionId, "error", {
      error: toErrorMessage(error, "Unknown verification error."),
    });
    updateSession(sessionId, (session) => ({
      ...session,
      status: "failed",
      message: `Verification error: ${toErrorMessage(
        error,
        "Unknown verification error.",
      )}`,
    }));
  });

  registerCallback(requestHandle, "onDisconnect", (event) => {
    const reason = toErrorMessage(
      event?.reason,
      "Bridge disconnected before proof generation completed.",
    );
    logSessionEvent(sessionId, "bridge_disconnected", {
      code: Number(event?.code ?? 0),
      reason,
      byServer: Boolean(event?.byServer),
      reconnecting: Boolean(event?.reconnecting),
    });
    updateSession(sessionId, (session) => {
      if (isTerminalStatus(session.status)) {
        return session;
      }
      return {
        ...session,
        bridgeConnected: false,
        message:
          "Bridge disconnected. Re-open the request in ZKPassport and try again.",
      };
    });
  });
}

async function createZKPassportSession(walletAddress, options = {}) {
  const normalizedWalletAddress = walletAddress.toLowerCase();
  closeOtherActiveSessionsForWallet(normalizedWalletAddress);
  const requestedDomain = normalizeDomainCandidate(options.domain);
  const effectiveDomain = resolveSessionDomain(requestedDomain);
  const effectiveMode = normalizeProofMode(options.mode);
  const zkPassport = new ZKPassport(effectiveDomain);
  const requestConfig = {
    name: ZKPASSPORT_APP_NAME,
    logo: ZKPASSPORT_APP_LOGO,
    purpose: ZKPASSPORT_PURPOSE,
    scope: ZKPASSPORT_SCOPE,
    mode: effectiveMode,
    validity: ZKPASSPORT_VALIDITY_SECONDS,
    devMode: ZKPASSPORT_DEV_MODE,
  };
  if (ZKPASSPORT_PROJECT_ID.trim()) {
    requestConfig.projectID = ZKPASSPORT_PROJECT_ID.trim();
  }

  if (ZKPASSPORT_BRIDGE_URL) {
    requestConfig.bridgeUrl = ZKPASSPORT_BRIDGE_URL;
  }
  if (ZKPASSPORT_CLOUD_PROVER_URL) {
    requestConfig.cloudProverUrl = ZKPASSPORT_CLOUD_PROVER_URL;
  }

  const queryBuilder = await zkPassport.request(requestConfig);

  const claimedBuilder = applyClaims(queryBuilder, CLAIMS);
  let requestBuilder = claimedBuilder;
  if (
    ZKPASSPORT_BIND_CHAIN &&
    typeof requestBuilder?.bind === "function" &&
    ZKPASSPORT_CHAIN.trim()
  ) {
    requestBuilder = requestBuilder.bind("chain", ZKPASSPORT_CHAIN.trim());
  }
  if (ZKPASSPORT_BIND_USER_ADDRESS && typeof requestBuilder?.bind === "function") {
    requestBuilder = requestBuilder.bind("user_address", normalizedWalletAddress);
  }
  const requestHandle = requestBuilder.done();
  const sessionId = requestHandle?.requestId || randomUUID();
  const expiresAt = new Date(Date.now() + SESSION_TTL_SECONDS * 1000).toISOString();
  const deepLinkUrl = requestHandle?.url;

  if (!deepLinkUrl) {
    throw new Error("ZKPassport did not return a request URL.");
  }

  const session = {
    sessionId,
    walletAddress: normalizedWalletAddress,
    domain: effectiveDomain,
    mode: effectiveMode,
    status: "pending",
    qrCodeUrl: buildQrCodeUrl(deepLinkUrl),
    deepLinkUrl,
    bridgeConnected:
      typeof requestHandle?.isBridgeConnected === "function"
        ? Boolean(requestHandle.isBridgeConnected())
        : false,
    requestReceived:
      typeof requestHandle?.requestReceived === "function"
        ? Boolean(requestHandle.requestReceived())
        : false,
    devMode: ZKPASSPORT_DEV_MODE,
    expiresAt,
    proof: null,
    proofEvents: [],
    // Keep strong references so SDK bridge callbacks remain alive.
    zkPassport,
    requestHandle,
    message: ZKPASSPORT_DEV_MODE
      ? "Scan the QR code with the ZKPassport app (dev mode enabled)."
      : "Scan the QR code with the ZKPassport app.",
  };

  logSessionEvent(sessionId, "session_created", {
    walletAddress: normalizedWalletAddress,
    requestedDomain,
    domain: effectiveDomain,
    chain: ZKPASSPORT_BIND_CHAIN ? ZKPASSPORT_CHAIN : "",
    bindUserAddress: ZKPASSPORT_BIND_USER_ADDRESS,
    mode: effectiveMode,
    expiresAt,
    devMode: ZKPASSPORT_DEV_MODE,
  });
  sessions.set(sessionId, session);
  bindSessionCallbacks(requestHandle, sessionId);
  return toPublicSession(session);
}

function copyProxyHeaders(upstreamResponse, origin) {
  const headers = createResponseHeaders(origin, {});
  const contentType = upstreamResponse.headers.get("content-type");
  if (contentType) {
    headers["Content-Type"] = contentType;
  }
  return headers;
}

async function proxyToUpstream(request, response, origin, pathWithQuery) {
  if (!UPSTREAM_API_URL) {
    writeJson(
      response,
      404,
      { success: false, message: "Route not found." },
      origin,
    );
    return;
  }

  const targetUrl = new URL(pathWithQuery, UPSTREAM_API_URL).toString();
  const method = request.method ?? "GET";
  const requestBody =
    method === "GET" || method === "HEAD" ? "" : await readBody(request);

  const headers = {};
  for (const [key, value] of Object.entries(request.headers)) {
    if (key.toLowerCase() === "host") {
      continue;
    }
    if (typeof value === "string") {
      headers[key] = value;
    }
  }

  const upstreamResponse = await fetch(targetUrl, {
    method,
    headers,
    body: requestBody ? requestBody : undefined,
  });

  const upstreamBody = await upstreamResponse.text();
  response.writeHead(
    upstreamResponse.status,
    copyProxyHeaders(upstreamResponse, origin),
  );
  response.end(upstreamBody);
}

async function proxyJsonRpc(request, response, origin) {
  if (!WORKFLOW_RPC_URL) {
    writeJson(
      response,
      500,
      {
        success: false,
        message:
          "WORKFLOW_RPC_URL (or CRE_RPC_URL / SEPOLIA_RPC) must be set to use /rpc.",
      },
      origin,
    );
    return;
  }

  if ((request.method ?? "GET").toUpperCase() !== "POST") {
    writeJson(
      response,
      405,
      {
        success: false,
        message: "Method not allowed. Use POST for JSON-RPC.",
      },
      origin,
      { Allow: "POST,OPTIONS" },
    );
    return;
  }

  const requestBody = await readBody(request);
  const upstreamResponse = await fetch(WORKFLOW_RPC_URL, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: requestBody || "{}",
  });

  const upstreamBody = await upstreamResponse.text();
  const contentType =
    upstreamResponse.headers.get("content-type") ||
    "application/json; charset=utf-8";
  response.writeHead(
    upstreamResponse.status,
    createResponseHeaders(origin, {
      "Content-Type": contentType,
      "Cache-Control": "no-store",
    }),
  );
  response.end(upstreamBody);
}

async function handleCreateSession(request, response, origin) {
  let body = {};
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

  const walletAddress = String(body?.walletAddress ?? "").trim();
  const requestedDomain = normalizeDomainCandidate(body?.domain);
  const requestedMode = normalizeProofMode(body?.mode, ZKPASSPORT_MODE);
  if (!isValidWalletAddress(walletAddress)) {
    writeJson(
      response,
      400,
      {
        success: false,
        message: "walletAddress is required and must be a valid EVM address.",
      },
      origin,
    );
    return;
  }

  try {
    const session = await createZKPassportSession(walletAddress, {
      domain: requestedDomain,
      mode: requestedMode,
    });
    writeJson(
      response,
      201,
      { success: true, message: "Session created.", data: session },
      origin,
    );
  } catch (error) {
    writeJson(
      response,
      502,
      {
        success: false,
        message: `Unable to create ZKPassport session: ${toErrorMessage(
          error,
          "Unknown error.",
        )}`,
      },
      origin,
    );
  }
}

function handleGetSession(response, origin, sessionId) {
  const session = sessions.get(sessionId);
  if (!session) {
    writeJson(
      response,
      404,
      { success: false, message: "Session not found." },
      origin,
    );
    return;
  }

  const expiresAtMs = Date.parse(session.expiresAt);
  if (Date.now() > expiresAtMs && !isTerminalStatus(session.status)) {
    session.status = "expired";
    session.message = "Session expired. Start a new verification request.";
    sessions.set(session.sessionId, session);
  }

  writeJson(
    response,
    200,
    {
      success: true,
      message: "Session loaded.",
      data: toPublicSession(sessions.get(sessionId)),
    },
    origin,
  );
}

async function routeRequest(request, response) {
  const origin = resolveCorsOrigin(request);
  if (origin === null) {
    const deniedOrigin =
      typeof request.headers.origin === "string"
        ? request.headers.origin
        : "unknown";
    writeJson(
      response,
      403,
      {
        success: false,
        message: `Origin ${deniedOrigin} is not allowed by CORS_ORIGIN.`,
      },
      "*",
    );
    return;
  }

  if (request.method === "OPTIONS") {
    response.writeHead(204, createResponseHeaders(origin));
    response.end();
    return;
  }

  const requestUrl = new URL(
    request.url ?? "/",
    `http://${request.headers.host ?? "localhost"}`,
  );
  cleanupSessions();

  if (request.method === "GET" && requestUrl.pathname === "/healthz") {
    writeJson(
      response,
      200,
      {
        success: true,
        message: "ok",
        data: {
          status: "ok",
          zkpassport: {
            domain: ZKPASSPORT_DOMAIN,
            scope: ZKPASSPORT_SCOPE,
            chain: ZKPASSPORT_CHAIN,
            bindChain: ZKPASSPORT_BIND_CHAIN,
            bindUserAddress: ZKPASSPORT_BIND_USER_ADDRESS,
            mode: ZKPASSPORT_MODE,
            devMode: ZKPASSPORT_DEV_MODE,
            projectIdConfigured: Boolean(ZKPASSPORT_PROJECT_ID.trim()),
            claims: CLAIMS,
          },
        },
      },
      origin,
    );
    return;
  }

  if (request.method === "POST" && requestUrl.pathname === "/auth/verify-wallet") {
    await handleVerifyWallet(request, response, origin, writeJson);
    return;
  }

  if (request.method === "POST" && requestUrl.pathname === "/auth/refresh") {
    handleRefreshToken(request, response, origin, writeJson);
    return;
  }

  if (request.method === "POST" && requestUrl.pathname === "/auth/logout") {
    handleLogout(response, origin, writeJson);
    return;
  }

  const contractAddressRaw = String(
    requestUrl.searchParams.get("contractAddress") || "",
  ).trim();
  if (contractAddressRaw && !isValidWalletAddress(contractAddressRaw)) {
    writeJson(
      response,
      400,
      {
        success: false,
        message: "contractAddress query must be a valid EVM address.",
      },
      origin,
    );
    return;
  }

  if (request.method === "GET" && requestUrl.pathname === "/auth/me") {
    const viewerWalletAddress = extractViewerWalletAddress(request);
    if (!viewerWalletAddress) {
      writeJson(
        response,
        401,
        {
          success: false,
          message: "Authentication is required to load user profile.",
        },
        origin,
      );
      return;
    }

    try {
      const user = await getAuthenticatedUserProfile(viewerWalletAddress);
      writeJson(
        response,
        200,
        {
          success: true,
          message: "User profile loaded.",
          data: user,
        },
        origin,
      );
    } catch (error) {
      writeJson(
        response,
        500,
        {
          success: false,
          message: `Failed to load user profile: ${toErrorMessage(
            error,
            "Unknown error.",
          )}`,
        },
        origin,
      );
    }
    return;
  }

  if (request.method === "GET" && requestUrl.pathname === "/notifications") {
    const viewerWalletAddress = extractViewerWalletAddress(request);
    if (!viewerWalletAddress) {
      writeJson(
        response,
        401,
        {
          success: false,
          message: "Authentication is required to load notifications.",
          data: [],
        },
        origin,
      );
      return;
    }

    const limit = parsePositiveInteger(
      requestUrl.searchParams.get("limit"),
      80,
    );
    const result = getNotificationsForWallet(viewerWalletAddress, limit);
    writeJson(
      response,
      200,
      {
        success: true,
        message: "Notifications loaded.",
        data: result.notifications,
        unreadCount: result.unreadCount,
      },
      origin,
    );
    return;
  }

  const notificationReadMatch = requestUrl.pathname.match(
    /^\/notifications\/([^/]+)\/read$/,
  );
  if (request.method === "POST" && notificationReadMatch) {
    const viewerWalletAddress = extractViewerWalletAddress(request);
    if (!viewerWalletAddress) {
      writeJson(
        response,
        401,
        {
          success: false,
          message: "Authentication is required to update notifications.",
        },
        origin,
      );
      return;
    }

    const updated = markWalletNotificationRead(
      viewerWalletAddress,
      decodeURIComponent(notificationReadMatch[1]),
    );
    if (!updated) {
      writeJson(
        response,
        404,
        {
          success: false,
          message: "Notification not found.",
        },
        origin,
      );
      return;
    }

    writeJson(
      response,
      200,
      {
        success: true,
        message: "Notification marked as read.",
      },
      origin,
    );
    return;
  }

  if (request.method === "GET" && requestUrl.pathname === "/messages/conversations") {
    const viewerWalletAddress = extractViewerWalletAddress(request);
    if (!viewerWalletAddress) {
      writeJson(
        response,
        401,
        {
          success: false,
          message: "Authentication is required to load conversations.",
          data: [],
        },
        origin,
      );
      return;
    }
    const tokenId = String(requestUrl.searchParams.get("tokenId") || "").trim();
    const conversations = getRoleGatedConversationsForWallet(
      viewerWalletAddress,
      tokenId || undefined,
    );
    writeJson(
      response,
      200,
      {
        success: true,
        message: "Conversations loaded.",
        data: conversations,
      },
      origin,
    );
    return;
  }

  const conversationMatch = requestUrl.pathname.match(
    /^\/messages\/conversations\/([^/]+)$/,
  );
  if (request.method === "GET" && conversationMatch) {
    const viewerWalletAddress = extractViewerWalletAddress(request);
    if (!viewerWalletAddress) {
      writeJson(
        response,
        401,
        {
          success: false,
          message: "Authentication is required to load conversation messages.",
        },
        origin,
      );
      return;
    }
    try {
      const conversation = getRoleGatedConversationForWallet(
        viewerWalletAddress,
        decodeURIComponent(conversationMatch[1]),
      );
      writeJson(
        response,
        200,
        {
          success: true,
          message: "Conversation loaded.",
          data: conversation,
        },
        origin,
      );
    } catch (error) {
      writeJson(
        response,
        403,
        {
          success: false,
          message: toErrorMessage(error, "Unable to load conversation."),
        },
        origin,
      );
    }
    return;
  }

  if (request.method === "POST" && requestUrl.pathname === "/messages/send") {
    const viewerWalletAddress = extractViewerWalletAddress(request);
    if (!viewerWalletAddress) {
      writeJson(
        response,
        401,
        {
          success: false,
          message: "Authentication is required to send private messages.",
        },
        origin,
      );
      return;
    }

    let body = {};
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
      const result = await sendRoleGatedWalletMessage({
        tokenId: body?.tokenId,
        senderWalletAddress: viewerWalletAddress,
        recipientWalletAddress: body?.recipientWalletAddress ?? body?.toAddress,
        message: body?.message,
        xmtpMessageId: body?.xmtpMessageId,
        contractAddress: contractAddressRaw || undefined,
      });
      writeJson(
        response,
        200,
        {
          success: true,
          message: "Private XMTP message sent.",
          data: result,
        },
        origin,
      );
    } catch (error) {
      writeJson(
        response,
        400,
        {
          success: false,
          message: toErrorMessage(error, "Unable to send private message."),
        },
        origin,
      );
    }
    return;
  }

  const viewerWalletAddress = extractViewerWalletAddress(request);

  const singleHouseMatch = requestUrl.pathname.match(/^\/houses\/(\d+)$/);
  if (request.method === "GET" && singleHouseMatch) {
    const tokenId = singleHouseMatch[1];
    try {
      const houses = await readHousesFromChain({
        contractAddress: contractAddressRaw || undefined,
      });
      const house = houses.find((entry) => String(entry.tokenId) === tokenId);
      if (!house) {
        writeJson(
          response,
          404,
          {
            success: false,
            message: "House not found.",
          },
          origin,
        );
        return;
      }

      const projectedHouse = projectHouseForViewer(house, viewerWalletAddress);
      const isAuthorized = isViewerAuthorizedForHouse(
        house,
        viewerWalletAddress,
      );

      writeJson(
        response,
        200,
        {
          success: true,
          message: isAuthorized
            ? "House loaded from CRE workflow state."
            : "House loaded with redacted private fields.",
          data: projectedHouse,
        },
        origin,
      );
    } catch (error) {
      writeJson(
        response,
        502,
        {
          success: false,
          message: `Failed to load house from onchain state: ${toErrorMessage(
            error,
            "Unknown error.",
          )}`,
        },
        origin,
      );
    }
    return;
  }

  const houseBillsMatch = requestUrl.pathname.match(/^\/houses\/(\d+)\/bills$/);
  if (request.method === "GET" && houseBillsMatch) {
    const tokenId = houseBillsMatch[1];
    try {
      const houses = await readHousesFromChain({
        contractAddress: contractAddressRaw || undefined,
      });
      const house = houses.find((entry) => String(entry.tokenId) === tokenId);
      if (!house) {
        writeJson(
          response,
          404,
          {
            success: false,
            message: "House not found.",
            data: [],
          },
          origin,
        );
        return;
      }

      if (!isViewerAuthorizedForHouse(house, viewerWalletAddress)) {
        writeJson(
          response,
          403,
          {
            success: false,
            message:
              "Bills are private. Authenticate as the minter, current owner, or active renter.",
            data: [],
          },
          origin,
        );
        return;
      }

      writeJson(
        response,
        200,
        {
          success: true,
          message: "Bills loaded from CRE workflow state.",
          data: Array.isArray(house.bills) ? house.bills : [],
        },
        origin,
      );
    } catch (error) {
      writeJson(
        response,
        502,
        {
          success: false,
          message: `Failed to load bills from onchain state: ${toErrorMessage(
            error,
            "Unknown error.",
          )}`,
          data: [],
        },
        origin,
      );
    }
    return;
  }

  const houseDocumentsMatch = requestUrl.pathname.match(
    /^\/houses\/(\d+)\/documents$/,
  );
  const houseDocumentContentMatch = requestUrl.pathname.match(
    /^\/houses\/(\d+)\/documents\/content$/,
  );
  if (request.method === "GET" && houseDocumentContentMatch) {
    const tokenId = houseDocumentContentMatch[1];
    if (!viewerWalletAddress) {
      writeJson(
        response,
        401,
        {
          success: false,
          message: "Authentication is required to view private documents.",
          data: { documents: [] },
        },
        origin,
      );
      return;
    }

    try {
      const houses = await readHousesFromChain({
        contractAddress: contractAddressRaw || undefined,
      });
      const house = houses.find((entry) => String(entry.tokenId) === tokenId);
      if (!house) {
        writeJson(
          response,
          404,
          {
            success: false,
            message: "House not found.",
            data: { documents: [] },
          },
          origin,
        );
        return;
      }

      if (!isViewerCurrentOwner(house, viewerWalletAddress)) {
        writeJson(
          response,
          403,
          {
            success: false,
            message:
              "Only the current owner can view private document contents.",
            data: { documents: [] },
          },
          origin,
        );
        return;
      }

      if (REQUIRE_KYC_FOR_PRIVATE_DOCUMENTS) {
        const kycStatus = await readKYCStatusForWallet(viewerWalletAddress);
        if (kycStatus !== "verified") {
          writeJson(
            response,
            403,
            {
              success: false,
              message:
                "KYC verification is required before private documents can be viewed.",
              data: { documents: [] },
            },
            origin,
          );
          return;
        }
      }

      const bundle = getPrivateDocumentBundleForToken(tokenId);
      if (!bundle || !bundle.documentsB64) {
        writeJson(
          response,
          404,
          {
            success: false,
            message:
              "No private document bundle is available for this property.",
            data: { documents: [] },
          },
          origin,
        );
        return;
      }

      const documents = parsePrivateDocumentBundle(bundle.documentsB64);
      if (documents.length === 0) {
        writeJson(
          response,
          404,
          {
            success: false,
            message:
              "Private document bundle is empty or malformed for this property.",
            data: { documents: [] },
          },
          origin,
        );
        return;
      }

      writeJson(
        response,
        200,
        {
          success: true,
          message: "Private document bundle loaded for verified owner.",
          data: {
            documents,
            documentHash: bundle.documentHash,
            documentURI: bundle.documentURI,
          },
        },
        origin,
      );
    } catch (error) {
      writeJson(
        response,
        502,
        {
          success: false,
          message: `Failed to load private document contents: ${toErrorMessage(
            error,
            "Unknown error.",
          )}`,
          data: { documents: [] },
        },
        origin,
      );
    }
    return;
  }

  if (request.method === "GET" && houseDocumentsMatch) {
    const tokenId = houseDocumentsMatch[1];
    try {
      const houses = await readHousesFromChain({
        contractAddress: contractAddressRaw || undefined,
      });
      const house = houses.find((entry) => String(entry.tokenId) === tokenId);
      if (!house) {
        writeJson(
          response,
          404,
          {
            success: false,
            message: "House not found.",
            data: { documents: [] },
          },
          origin,
        );
        return;
      }

      if (!isViewerAuthorizedForHouse(house, viewerWalletAddress)) {
        writeJson(
          response,
          403,
          {
            success: false,
            message:
              "Documents are private. Authenticate as the minter, current owner, or active renter.",
            data: { documents: [] },
          },
          origin,
        );
        return;
      }

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

      writeJson(
        response,
        200,
        {
          success: true,
          message: "Document metadata loaded from CRE confidential store.",
          data: { documents },
        },
        origin,
      );
    } catch (error) {
      writeJson(
        response,
        502,
        {
          success: false,
          message: `Failed to load documents from onchain state: ${toErrorMessage(
            error,
            "Unknown error.",
          )}`,
          data: { documents: [] },
        },
        origin,
      );
    }
    return;
  }

  if (request.method === "GET" && requestUrl.pathname === "/houses") {
    const ownerRaw = String(requestUrl.searchParams.get("owner") || "").trim();
    if (ownerRaw && !isValidWalletAddress(ownerRaw)) {
      writeJson(
        response,
        400,
        {
          success: false,
          message: "owner query must be a valid EVM address.",
        },
        origin,
      );
      return;
    }

    if (ownerRaw && !viewerWalletAddress) {
      writeJson(
        response,
        401,
        {
          success: false,
          message:
            "Authentication is required to query owner-specific private houses.",
          data: [],
        },
        origin,
      );
      return;
    }

    if (
      ownerRaw &&
      viewerWalletAddress &&
      normalizeWalletAddress(ownerRaw) !== viewerWalletAddress
    ) {
      writeJson(
        response,
        403,
        {
          success: false,
          message:
            "Owner query must match the authenticated wallet for private house data.",
          data: [],
        },
        origin,
      );
      return;
    }

    try {
      const houses = await readHousesFromChain({
        ownerAddress: ownerRaw || undefined,
        contractAddress: contractAddressRaw || undefined,
      });
      const projectedHouses = houses.map((house) =>
        projectHouseForViewer(house, viewerWalletAddress),
      );
      writeJson(
        response,
        200,
        {
          success: true,
          message: viewerWalletAddress
            ? "Houses loaded from CRE workflow state."
            : "Houses loaded with redacted private fields.",
          data: projectedHouses,
        },
        origin,
      );
    } catch (error) {
      writeJson(
        response,
        502,
        {
          success: false,
          message: `Failed to load houses from onchain state: ${toErrorMessage(
            error,
            "Unknown error.",
          )}`,
          data: [],
        },
        origin,
      );
    }
    return;
  }

  const balanceMatch = requestUrl.pathname.match(/^\/balances\/(0x[a-fA-F0-9]{40})$/);
  if (request.method === "GET" && balanceMatch) {
    const address = balanceMatch[1];
    try {
      const balance = await readNativeBalance(address);
      writeJson(
        response,
        200,
        {
          success: true,
          message: "Balance loaded from onchain state.",
          data: balance,
        },
        origin,
      );
    } catch (error) {
      writeJson(
        response,
        502,
        {
          success: false,
          message: `Failed to load balance from onchain state: ${toErrorMessage(
            error,
            "Unknown error.",
          )}`,
        },
        origin,
      );
    }
    return;
  }

  if (request.method === "POST" && requestUrl.pathname === "/workflow/trigger") {
    await handleWorkflowTrigger(request, response, origin, writeJson);
    return;
  }

  if (requestUrl.pathname === "/rpc") {
    await proxyJsonRpc(request, response, origin);
    return;
  }

  if (
    request.method === "POST"
    && requestUrl.pathname === "/kyc/zkpassport/verify"
  ) {
    await handleVerifyKYC(request, response, origin, writeJson);
    return;
  }

  if (request.method === "POST" && requestUrl.pathname === "/kyc/verify") {
    await handleVerifyKYCForCRE(request, response, origin, writeJson);
    return;
  }

  if (
    request.method === "POST" &&
    requestUrl.pathname === "/kyc/zkpassport/session"
  ) {
    await handleCreateSession(request, response, origin);
    return;
  }

  const sessionMatch = requestUrl.pathname.match(
    /^\/kyc\/zkpassport\/session\/([^/]+)$/,
  );
  if (request.method === "GET" && sessionMatch) {
    handleGetSession(response, origin, decodeURIComponent(sessionMatch[1]));
    return;
  }

  await proxyToUpstream(
    request,
    response,
    origin,
    `${requestUrl.pathname}${requestUrl.search}`,
  );
}

const server = createServer((request, response) => {
  routeRequest(request, response).catch((error) => {
    const origin = resolveCorsOrigin(request) || "*";
    writeJson(
      response,
      500,
      {
        success: false,
        message: `Unexpected server error: ${toErrorMessage(
          error,
          "Unknown error.",
        )}`,
      },
      origin,
    );
  });
});

server.listen(PORT, HOST, () => {
  console.log(
    `[zkpassport-session-service] listening on ${HOST}:${PORT} (domain=${ZKPASSPORT_DOMAIN})`,
  );
  console.log(`[zkpassport-session-service] chain=${ZKPASSPORT_CHAIN}`);
  console.log(
    `[zkpassport-session-service] bindChain=${ZKPASSPORT_BIND_CHAIN} bindUserAddress=${ZKPASSPORT_BIND_USER_ADDRESS}`,
  );
  console.log(`[zkpassport-session-service] mode=${ZKPASSPORT_MODE}`);
  console.log(`[zkpassport-session-service] CORS_ORIGIN=${CORS_ORIGIN}`);
  console.log(
    `[zkpassport-session-service] claims=${JSON.stringify(CLAIMS)}`,
  );
  if (UPSTREAM_API_URL) {
    console.log(`[zkpassport-session-service] proxy upstream: ${UPSTREAM_API_URL}`);
  }
});
