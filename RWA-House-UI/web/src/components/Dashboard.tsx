/**
 * Dashboard component
 * Main user interface showing houses, listings, and quick actions.
 */

import React, { useEffect, useRef, useState } from "react";
import { useAuth } from "./AuthProvider";
import { HouseThumbnail } from "./HouseThumbnail";
import { useUXMode } from "./UXModeProvider";
import {
  apiClient,
  KYC_PROOF_STORAGE_KEY,
  KYC_PROVIDER_STORAGE_KEY,
} from "../../../shared/src/utils/api";
import type {
  House,
  KYCProvider,
  ZKPassportSession,
} from "../../../shared/src/types";
import { useNavigate } from "react-router-dom";
import toast from "react-hot-toast";

const getSessionTone = (
  status: ZKPassportSession["status"],
): { badge: string; text: string } => {
  switch (status) {
    case "verified":
      return {
        badge:
          "border border-emerald-300/45 bg-emerald-400/15 text-emerald-200",
        text: "text-emerald-200",
      };
    case "failed":
    case "expired":
      return {
        badge: "border border-rose-300/45 bg-rose-400/15 text-rose-200",
        text: "text-rose-200",
      };
    case "ready":
      return {
        badge: "border border-blue-300/45 bg-blue-400/15 text-blue-200",
        text: "text-blue-200",
      };
    default:
      return {
        badge: "border border-amber-300/45 bg-amber-400/15 text-amber-200",
        text: "text-amber-200",
      };
  }
};

type DashboardZKSession = ZKPassportSession & {
  bridgeConnected?: boolean;
  requestReceived?: boolean;
  source?: "backend";
};

const parseStoredKYCProvider = (value: string | null): KYCProvider => {
  if (value === "none" || value === "zkpassport") {
    return value;
  }
  return "mock";
};

const parseKYCProofInput = (
  value: string,
): { parsedProof: Record<string, any> | null; isValid: boolean } => {
  const trimmed = value.trim();
  if (!trimmed) {
    return { parsedProof: null, isValid: true };
  }

  try {
    return { parsedProof: JSON.parse(trimmed), isValid: true };
  } catch {
    return { parsedProof: null, isValid: false };
  }
};

export const Dashboard: React.FC = () => {
  const {
    embeddedWalletAddress,
    externalWalletAddresses,
    refreshCurrentUser,
    setUserKycStatus,
    user,
    walletAddress,
  } = useAuth();
  const { mode } = useUXMode();
  const navigate = useNavigate();
  const zkPassportApiBase = apiClient.getZKPassportBaseURL();
  const [houses, setHouses] = useState<House[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [kycProvider, setKYCProvider] = useState<KYCProvider>("mock");
  const [kycProofText, setKYCProofText] = useState("");
  const [kycMessage, setKYCMessage] = useState<string | null>(null);
  const [zkPassportSession, setZKPassportSession] =
    useState<DashboardZKSession | null>(null);
  const [isStartingZKPassport, setIsStartingZKPassport] = useState(false);
  const [isQrImageUnavailable, setIsQrImageUnavailable] = useState(false);
  const processedZKPassportSessionRef = useRef<string | null>(null);
  const isVerifyingZKPassportProofRef = useRef(false);
  const walletDisplay =
    walletAddress && /^0x[a-fA-F0-9]{40}$/.test(walletAddress)
      ? `${walletAddress.slice(0, 6)}...${walletAddress.slice(-4)}`
      : user?.email || "Not connected";
  const isWalletKYCVerified = user?.kycStatus === "verified";
  const activeVerificationWallet =
    walletAddress || embeddedWalletAddress || externalWalletAddresses[0] || "";
  const isAnonymousMode = kycProvider === "none";
  const walletKYCBadge = isAnonymousMode
    ? {
        icon: "🕶️",
        text: "Anonymous mode",
        className:
          "border-slate-300/45 bg-slate-500/20 text-slate-200",
      }
    : isWalletKYCVerified
      ? {
          icon: "✅",
          text: "KYC verified",
          className:
            "border-emerald-300/50 bg-emerald-500/15 text-emerald-200",
        }
      : {
          icon: "⚠️",
          text: "KYC verification required",
          className:
            "border-amber-300/50 bg-amber-500/15 text-amber-200",
        };

  const persistKYCDefaults = (
    provider: KYCProvider,
    proofText: string,
    parsedProof: Record<string, any> | null,
  ) => {
    window.localStorage.setItem(KYC_PROVIDER_STORAGE_KEY, provider);
    if (provider === "none") {
      window.localStorage.removeItem(KYC_PROOF_STORAGE_KEY);
      apiClient.setKYCDefaults("none", null);
      return;
    }

    const trimmedProof = proofText.trim();
    if (trimmedProof) {
      window.localStorage.setItem(KYC_PROOF_STORAGE_KEY, trimmedProof);
    } else {
      window.localStorage.removeItem(KYC_PROOF_STORAGE_KEY);
    }

    apiClient.setKYCDefaults(provider, parsedProof);
  };

  const openMintPage = () => {
    navigate("/mint");
  };

  useEffect(() => {
    loadData();
  }, [walletAddress]);

  useEffect(() => {
    const provider = parseStoredKYCProvider(
      window.localStorage.getItem(KYC_PROVIDER_STORAGE_KEY),
    );
    const proofRaw =
      provider === "none"
        ? ""
        : window.localStorage.getItem(KYC_PROOF_STORAGE_KEY) || "";
    setKYCProvider(provider);
    setKYCProofText(proofRaw);

    if (provider === "none") {
      window.localStorage.removeItem(KYC_PROOF_STORAGE_KEY);
      apiClient.setKYCDefaults("none", null);
      return;
    }

    const { parsedProof, isValid } = parseKYCProofInput(proofRaw);
    if (!isValid) {
      setKYCMessage("Saved proof JSON is invalid. Update and save again.");
    }

    apiClient.setKYCDefaults(provider, parsedProof);
  }, []);

  const chooseKYCMode = () => {
    const provider = kycProvider === "none" ? "mock" : kycProvider;
    const { parsedProof, isValid } = parseKYCProofInput(kycProofText);

    setKYCProvider(provider);
    persistKYCDefaults(provider, kycProofText, isValid ? parsedProof : null);
    setKYCMessage(
      isValid
        ? "KYC mode enabled. Choose mock or ZKPassport for workflow payloads."
        : "KYC mode enabled. Proof JSON is invalid and was not applied.",
    );
  };

  const chooseAnonymousMode = () => {
    setKYCProvider("none");
    setKYCProofText("");
    setZKPassportSession(null);
    setIsQrImageUnavailable(false);
    persistKYCDefaults("none", "", null);
    setKYCMessage(
      "Anonymous mode enabled. Workflow payloads now use kycProvider=none.",
    );
  };

  const saveKYCConfig = () => {
    if (kycProvider === "none") {
      setKYCProofText("");
      setZKPassportSession(null);
      setIsQrImageUnavailable(false);
      persistKYCDefaults("none", "", null);
      setKYCMessage(
        "Anonymous mode saved. Workflow payloads now skip KYC verification.",
      );
      return;
    }

    const { parsedProof, isValid } = parseKYCProofInput(kycProofText);
    if (!isValid) {
      setKYCMessage("KYC proof must be valid JSON");
      return;
    }

    persistKYCDefaults(kycProvider, kycProofText, parsedProof);
    setKYCMessage(
      kycProvider === "zkpassport"
        ? "ZKPassport route enabled for mint/list/sell/rent payloads."
        : "Mock KYC route enabled for mint/list/sell/rent payloads.",
    );
  };

  const clearKYCProof = () => {
    setKYCProofText("");
    persistKYCDefaults(kycProvider, "", null);
    setKYCMessage("Proof JSON cleared.");
  };

  const applyVerifiedZKPassportProof = async (
    session: DashboardZKSession,
    verificationWallet: string,
  ) => {
    if (
      isVerifyingZKPassportProofRef.current ||
      processedZKPassportSessionRef.current === session.sessionId
    ) {
      return;
    }
    if (!session.proof) {
      return;
    }

    isVerifyingZKPassportProofRef.current = true;
    try {
      const verifyResponse = await apiClient.verifyZKPassportProof(
        verificationWallet,
        session.proof as any,
      );
      if (!verifyResponse.success || !verifyResponse.data?.verified) {
        throw new Error(
          verifyResponse.message ||
            "Server verification failed for ZKPassport proof.",
        );
      }

      const verifiedProof = verifyResponse.data.proof || (session.proof as any);
      const formattedProof = JSON.stringify(verifiedProof, null, 2);
      setKYCProvider("zkpassport");
      setKYCProofText(formattedProof);
      window.localStorage.setItem(KYC_PROVIDER_STORAGE_KEY, "zkpassport");
      window.localStorage.setItem(KYC_PROOF_STORAGE_KEY, formattedProof);
      apiClient.setKYCDefaults("zkpassport", verifiedProof);
      const ensureKYCResponse = await apiClient.ensureKYC(verificationWallet);
      if (!ensureKYCResponse.success) {
        throw new Error(
          ensureKYCResponse.message ||
            "Unable to write verified KYC status onchain.",
        );
      }
      processedZKPassportSessionRef.current = session.sessionId;
      setUserKycStatus("verified");
      await refreshCurrentUser();
      setKYCMessage("KYCed! ZKPassport verification confirmed and saved.");
      toast.success("KYCed! ZKPassport verification confirmed.");
    } catch (verifyErr) {
      const message =
        verifyErr instanceof Error
          ? verifyErr.message
          : "Unable to verify ZKPassport proof.";
      const isOnchainSyncIssue =
        message.toLowerCase().includes("onchain")
        || message.toLowerCase().includes("setkycverification")
        || message.toLowerCase().includes("kyc recorded");
      setKYCMessage(
        isOnchainSyncIssue
          ? `ZKPassport proof verified, but onchain KYC sync failed: ${message}`
          : `ZKPassport proof verification failed: ${message}`,
      );
      toast.error(
        isOnchainSyncIssue
          ? "KYC proof verified, but onchain sync failed."
          : "ZKPassport proof verification failed.",
      );
    } finally {
      isVerifyingZKPassportProofRef.current = false;
    }
  };

  useEffect(() => {
    if (!zkPassportSession?.sessionId) {
      return;
    }
    if (zkPassportSession.source !== "backend") {
      return;
    }
    if (
      zkPassportSession.status === "verified" ||
      zkPassportSession.status === "failed" ||
      zkPassportSession.status === "expired"
    ) {
      return;
    }

    let cancelled = false;
    const pollSession = async () => {
      const sessionResponse = await apiClient.getZKPassportSession(
        zkPassportSession.sessionId,
      );
      if (cancelled || !sessionResponse.success || !sessionResponse.data) {
        return;
      }

      const normalizedSession = normalizeZKPassportSession(
        sessionResponse.data as DashboardZKSession,
      );
      setZKPassportSession((previous) => {
        if (!previous || previous.sessionId !== normalizedSession.sessionId) {
          return previous;
        }
        return {
          ...previous,
          ...normalizedSession,
          bridgeConnected:
            normalizedSession.bridgeConnected ?? previous.bridgeConnected,
          requestReceived:
            normalizedSession.requestReceived ?? previous.requestReceived,
        };
      });
    };

    void pollSession();
    const intervalId = window.setInterval(() => {
      void pollSession();
    }, 3000);

    return () => {
      cancelled = true;
      window.clearInterval(intervalId);
    };
  }, [zkPassportSession?.sessionId, zkPassportSession?.status]);

  useEffect(() => {
    if (
      !zkPassportSession?.sessionId ||
      zkPassportSession.status !== "verified"
    ) {
      return;
    }
    if (!zkPassportSession.proof) {
      return;
    }
    if (!/^0x[a-fA-F0-9]{40}$/.test(activeVerificationWallet)) {
      return;
    }

    void applyVerifiedZKPassportProof(
      zkPassportSession,
      activeVerificationWallet,
    );
  }, [
    activeVerificationWallet,
    zkPassportSession?.proof,
    zkPassportSession?.sessionId,
    zkPassportSession?.status,
  ]);

  const startZKPassportFlow = async () => {
    if (!zkPassportApiBase.trim()) {
      const missingBaseMessage =
        "ZKPassport API URL is missing. Set VITE_ZKPASSPORT_API_URL and retry.";
      setKYCMessage(missingBaseMessage);
      toast.error(missingBaseMessage);
      return;
    }

    const verificationWallet = activeVerificationWallet;
    if (!verificationWallet) {
      setKYCMessage(
        "Connect a wallet first to launch ZKPassport verification.",
      );
      return;
    }
    if (!/^0x[a-fA-F0-9]{40}$/.test(verificationWallet)) {
      const invalidWalletMessage =
        "Connected wallet address is invalid for ZKPassport session creation.";
      setKYCMessage(invalidWalletMessage);
      toast.error(invalidWalletMessage);
      return;
    }
    try {
      setIsStartingZKPassport(true);
      setKYCProvider("zkpassport");
      setKYCMessage(null);
      setIsQrImageUnavailable(false);
      processedZKPassportSessionRef.current = null;

      let healthCheckWarning: string | null = null;
      try {
        const healthResponse = await apiClient.getZKPassportHealth();
        if (!healthResponse.success) {
          healthCheckWarning =
            healthResponse.message ||
            "ZKPassport session service health check failed.";
        }
      } catch (healthErr) {
        healthCheckWarning =
          healthErr instanceof Error
            ? healthErr.message
            : "Health check failed.";
      }

      const requestedDomain =
        typeof window !== "undefined" && window.location.hostname
          ? window.location.hostname.trim().toLowerCase()
          : "";
      const shouldSendDomain =
        requestedDomain.length > 0 &&
        requestedDomain !== "localhost" &&
        requestedDomain !== "127.0.0.1" &&
        requestedDomain !== "::1";
      const backendSessionResponse = await apiClient.startZKPassportSession(
        verificationWallet,
        shouldSendDomain
          ? {
              domain: requestedDomain,
            }
          : undefined,
      );
      if (!backendSessionResponse.success || !backendSessionResponse.data) {
        const backendErrorMessage = String(
          backendSessionResponse.message || "Unable to start ZKPassport session.",
        ).trim();
        setKYCMessage(
          `Unable to start backend ZKPassport session: ${backendErrorMessage}`,
        );
        toast.error("Unable to start backend ZKPassport session.");
        return;
      }

      const normalizedSession = normalizeZKPassportSession({
        ...(backendSessionResponse.data as DashboardZKSession),
        source: "backend",
      });
      setZKPassportSession(normalizedSession);
      const sessionMessage =
        !normalizedSession.qrCodeUrl && normalizedSession.deepLinkUrl
          ? "Session started. QR image unavailable, use deep link below."
          : "ZKPassport session started. Scan the QR code and complete verification.";
      const warningSuffix = healthCheckWarning
        ? ` (Health check warning: ${healthCheckWarning})`
        : "";
      setKYCMessage(`${sessionMessage}${warningSuffix}`);
    } catch (err) {
      console.error("Start ZKPassport flow failed:", err);
      const frontendOrigin =
        typeof window !== "undefined" ? window.location.origin : "unknown";
      const isNetworkOrCorsIssue =
        err instanceof TypeError ||
        (err instanceof Error &&
          /failed to fetch|network/i.test(err.message.toLowerCase()));
      const details =
        err instanceof Error && err.message.trim().length > 0
          ? err.message
          : "Unknown network error";
      const guidance = isNetworkOrCorsIssue
        ? "Confirm the session service is running on port 8787, start it with " +
          "`HOST=0.0.0.0`, ensure CORS_ORIGIN allows this frontend origin, and " +
          "verify VITE_ZKPASSPORT_API_URL points to that service."
        : "Check zkpassport-session-service logs for the session creation error.";
      const toastMessage = isNetworkOrCorsIssue
        ? `Unable to start ZKPassport (network/CORS): ${details}`
        : `Unable to start ZKPassport: ${details}`;
      toast.error(toastMessage);
      setKYCMessage(
        `Unable to start ZKPassport flow from ${zkPassportApiBase} ` +
          `(frontend origin ${frontendOrigin}). ${details}. ` +
          guidance,
      );
    } finally {
      setIsStartingZKPassport(false);
    }
  };

  const refreshZKPassportSession = () => {
    if (!zkPassportSession) {
      return;
    }
    if (zkPassportSession.source !== "backend") {
      return;
    }

    void (async () => {
      const sessionResponse = await apiClient.getZKPassportSession(
        zkPassportSession.sessionId,
      );
      if (!sessionResponse.success || !sessionResponse.data) {
        setKYCMessage(
          sessionResponse.message || "Unable to refresh ZKPassport session.",
        );
        return;
      }

      const normalizedSession = normalizeZKPassportSession(
        sessionResponse.data as DashboardZKSession,
      );
      setZKPassportSession((previous) => {
        if (!previous || previous.sessionId !== normalizedSession.sessionId) {
          return {
            ...normalizedSession,
            source: "backend",
          };
        }

        return {
          ...previous,
          ...normalizedSession,
          source: previous.source || "backend",
        };
      });

      if (
        normalizedSession.status === "verified" &&
        normalizedSession.proof &&
        /^0x[a-fA-F0-9]{40}$/.test(activeVerificationWallet)
      ) {
        await applyVerifiedZKPassportProof(
          normalizedSession,
          activeVerificationWallet,
        );
      }
    })();
  };

  const loadData = async () => {
    if (!walletAddress) {
      setHouses([]);
      setIsLoading(false);
      return;
    }

    try {
      setIsLoading(true);
      setError(null);

      // Fetch user's houses
      const housesResponse = await apiClient.getHouses(walletAddress);
      if (housesResponse.success && housesResponse.data) {
        setHouses(housesResponse.data);
        return;
      }

      setHouses([]);
      setError(housesResponse.message || "Failed to load dashboard data");
    } catch (err) {
      setError("Failed to load dashboard data");
      console.error("Dashboard load error:", err);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateBillFromDashboard = () => {
    if (houses.length === 0) {
      toast("Mint a property first, then add bills.", { icon: "ℹ️" });
      openMintPage();
      return;
    }

    navigate(`/houses/${houses[0].tokenId}/bills/create`);
  };

  if (isLoading) {
    return (
      <div className="loading-container bg-[#060b14]">
        <div className="spinner" />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[#060b14] pb-14">
      <header className="border-b border-slate-700/50 bg-slate-950/45">
        <div className="page-shell page-shell-tight !pb-0 !pt-7">
          <div
            className={
              "glass-card hero-panel dashboard-hero-card mx-auto !w-full !max-w-[54rem]" +
              " !rounded-2xl !p-6 md:!p-8"
            }
          >
            <div
              className={
                "flex flex-col items-center gap-6 text-center lg:flex-row" +
                " lg:items-center lg:justify-between lg:text-left"
              }
            >
              <div className="space-y-3">
                <div className="flex items-center justify-center lg:justify-start">
                  <h1 className="dashboard-title">
                    {mode === "degen" ? "Operations center" : "Your properties"}
                  </h1>
                </div>
                <p className="dashboard-meta mx-auto px-3 py-2 text-sm lg:mx-0">
                  <span className="h-2 w-2 rounded-full bg-emerald-400" />
                  {mode === "degen" ? "Operator" : "Wallet"}: {walletDisplay}
                  <span
                    className={`ml-2 inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[11px] font-semibold uppercase tracking-[0.08em] ${walletKYCBadge.className}`}
                  >
                    {walletKYCBadge.icon} {walletKYCBadge.text}
                  </span>
                </p>
              </div>
              <div
                className={
                  "flex flex-wrap items-center justify-center gap-3 rounded-2xl" +
                  " border border-slate-600/45 bg-slate-900/35 px-4 py-3" +
                  " lg:justify-end"
                }
              >
                <button
                  type="button"
                  onClick={loadData}
                  className="btn btn-secondary !px-5 !py-2.5 !text-sm"
                >
                  Refresh
                </button>
                <button
                  type="button"
                  onClick={openMintPage}
                  className="btn btn-primary !px-5 !py-2.5 !text-sm"
                >
                  {mode === "degen" ? "+ Mint asset" : "+ Add property"}
                </button>
              </div>
            </div>
          </div>
        </div>
      </header>

      <main
        className="page-shell page-shell-tight workspace-surface !pt-0"
        style={{ marginTop: "-6rem" }}
      >
        {error && (
          <div className="text-panel mx-auto mb-6 w-full max-w-3xl rounded-xl border border-rose-400/45 bg-rose-500/10 px-6 py-4">
            <p className="flex items-center text-rose-200">
              <svg
                className="w-5 h-5 mr-2"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              {error}
            </p>
          </div>
        )}

        <div className="dashboard-layout mx-auto grid w-full max-w-3xl grid-cols-1 items-start justify-items-center gap-12">
          {/* Main Content */}
          <div className="w-full max-w-3xl space-y-10">
            {/* Quick Stats */}
            <div className="dashboard-stats-grid mx-auto grid w-full max-w-[44rem] grid-cols-1 justify-items-center gap-6 sm:grid-cols-3 sm:gap-6">
              <StatCard
                label="TOTAL ASSETS"
                value={houses.length}
                color="#93c5fd"
                icon="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
              />
              <StatCard
                label="FOR SALE"
                value={
                  houses.filter((h) => h.listing?.listingType === "for_sale")
                    .length
                }
                color="#86efac"
                icon="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
              <StatCard
                label="RENTED"
                value={houses.filter((h) => h.rental?.isActive).length}
                color="#a5b4fc"
                icon="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
              />
            </div>

            {/* Houses List */}
            <div className="cyber-card dashboard-panel mx-auto w-full max-w-[44rem] overflow-hidden">
              <div className="flex flex-col items-center justify-center gap-2 border-b border-slate-700/60 px-6 py-5 text-center">
                <h2 className="dashboard-section-title">
                  {mode === "degen" ? "Asset registry" : "Property list"}
                </h2>
                <p className="dashboard-section-note text-sm">
                  {mode === "degen"
                    ? "Track tokenized assets and their execution-ready state."
                    : "Review properties in your workspace before opening details or billing."}
                </p>
                <span className="dashboard-meta text-xs font-semibold">
                  <span className="number-pair">
                    <span className="number-pill number-pill-sm">
                      {houses.length}
                    </span>
                    <span>entries</span>
                  </span>
                </span>
              </div>

              {houses.length === 0 ? (
                <div className="dashboard-empty-state px-6 py-14 text-center">
                  <div className="mb-4 inline-block rounded-full bg-blue-500/10 p-4">
                    <svg
                      className="h-8 w-8 text-blue-300"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={1.5}
                        d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
                      />
                    </svg>
                  </div>
                  <p className="mb-4 text-slate-300">
                    {mode === "degen"
                      ? "No assets in registry yet."
                      : "No properties added yet."}
                  </p>
                  <button
                    type="button"
                    onClick={openMintPage}
                    className="btn btn-primary !mt-3 !px-5 !py-2.5 !text-sm"
                  >
                    {mode === "degen"
                      ? "Initialize first asset"
                      : "Add your first property"}
                  </button>
                </div>
              ) : (
                <div className="space-y-4 px-3 py-4">
                  {houses.map((house) => (
                    <div
                      key={house.tokenId}
                      className="group rounded-xl border border-slate-700/60 bg-slate-900/35 px-6 py-6 transition-colors hover:bg-slate-900/55"
                    >
                      <div className="flex flex-col gap-5 lg:flex-row lg:items-center lg:justify-between">
                        <div className="flex flex-col gap-5 sm:flex-row sm:items-start sm:gap-6 lg:flex-1">
                          <HouseThumbnail
                            house={house}
                            className="h-40 w-full sm:h-28 sm:w-44 sm:min-w-44"
                          />
                          <div className="flex-1">
                            <div className="flex items-center space-x-3">
                              <h3 className="text-lg font-medium text-slate-50 transition-colors group-hover:text-blue-200">
                                {house.metadata.address}
                              </h3>
                              <StatusBadge house={house} />
                            </div>
                            <p className="mt-1 text-sm text-slate-400">
                              {house.metadata.city}, {house.metadata.state} •{" "}
                              {house.metadata.propertyType}
                            </p>
                            <div className="mt-3 flex flex-wrap items-center gap-4 text-sm text-slate-400">
                              <span className="number-pair">
                                <span className="number-pill number-pill-sm">
                                  {house.metadata.bedrooms}
                                </span>
                                <span>beds</span>
                              </span>
                              <span className="number-pair">
                                <span className="number-pill number-pill-sm">
                                  {house.metadata.bathrooms}
                                </span>
                                <span>baths</span>
                              </span>
                              <span className="number-pair">
                                <span className="number-pill number-pill-sm">
                                  {house.metadata.squareFeet.toLocaleString()}
                                </span>
                                <span>sqft</span>
                              </span>
                            </div>
                          </div>
                        </div>
                        <div className="mt-2 flex flex-wrap items-center gap-x-3.5 gap-y-3">
                          <button
                            type="button"
                            onClick={() =>
                              navigate(`/houses/${house.tokenId}/documents`)
                            }
                            className="btn btn-secondary !min-w-[6.25rem] !justify-center !px-3.5 !py-2 !text-xs"
                          >
                            Documents
                          </button>
                          <button
                            type="button"
                            onClick={() =>
                              navigate(`/houses/${house.tokenId}/bills/create`)
                            }
                            className="btn btn-secondary !min-w-[6.25rem] !justify-center !px-3.5 !py-2 !text-xs"
                          >
                            Add bill
                          </button>
                          <button
                            type="button"
                            onClick={() => navigate(`/houses/${house.tokenId}`)}
                            className="btn btn-secondary !min-w-[6.25rem] !justify-center !px-3.5 !py-2 !text-xs"
                          >
                            Open
                          </button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Sidebar */}
          <div className="dashboard-side-stack w-full max-w-[44rem] space-y-10">
            <div className="cyber-card dashboard-panel mx-auto w-full max-w-[44rem] p-8">
              <h2 className="dashboard-section-title mb-2 text-center">
                {mode === "degen" ? "Execution panel" : "Quick actions"}
              </h2>
              <p className="text-panel mb-8 text-center text-sm text-slate-300">
                {mode === "degen"
                  ? "Run the key CRE-backed actions directly from one place."
                  : "Use these shortcuts to start minting, buying, renting, and billing workflows."}
              </p>
              <div className="quick-action-stack mx-auto w-full">
                <ActionButton
                  onClick={openMintPage}
                  title={mode === "degen" ? "Mint new asset" : "Add property"}
                  description={
                    mode === "degen"
                      ? "Tokenize property documents"
                      : "Create a new digital property record"
                  }
                  icon="M12 4v16m8-8H4"
                />
                <ActionButton
                  onClick={() => navigate("/marketplace")}
                  title={mode === "degen" ? "Browse marketplace" : "Find homes"}
                  description={
                    mode === "degen"
                      ? "Find assets for sale/rent"
                      : "Explore properties for sale or rent"
                  }
                  icon="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
                <ActionButton
                  onClick={handleCreateBillFromDashboard}
                  title={
                    mode === "degen" ? "Create bill" : "Create payment bill"
                  }
                  description={
                    mode === "degen"
                      ? "Open billing setup for your first asset"
                      : "Create utility or rent invoices for a property"
                  }
                  icon="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586A1 1 0 0113.293 3.293l3.414 3.414A1 1 0 0117 7.414V21a2 2 0 01-2 2z"
                />
                <ActionButton
                  onClick={() => navigate("/claim")}
                  title={mode === "degen" ? "Claim key" : "Get access key"}
                  description={
                    mode === "degen"
                      ? "Retrieve encrypted keys from completed flows"
                      : "Recover private document access after sale/rent"
                  }
                  icon="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2h-1V7a5 5 0 10-10 0v4H6a2 2 0 00-2 2v6a2 2 0 002 2z"
                />
              </div>
            </div>

            <div className="cyber-card dashboard-panel mx-auto w-full max-w-[44rem] p-8">
              <h2 className="dashboard-section-title text-center">
                {mode === "degen"
                  ? "KYC route settings"
                  : "Identity verification settings"}
              </h2>
              <p className="text-panel mt-2 text-center text-sm text-slate-300">
                {mode === "degen"
                  ? "Choose to KYC or stay anonymous for CRE mint/list/sell/rent calls."
                  : "Choose whether to verify identity or stay anonymous for mint, list, buy, and rent actions."}
              </p>

              <div className="mt-7 flex flex-wrap justify-center gap-4">
                <button
                  className={`cyber-btn text-sm ${
                    !isAnonymousMode ? "cyber-btn-primary" : ""
                  }`}
                  onClick={chooseKYCMode}
                  type="button"
                >
                  Choose to KYC
                </button>
                <button
                  className={`cyber-btn text-sm ${
                    isAnonymousMode ? "cyber-btn-primary" : ""
                  }`}
                  onClick={chooseAnonymousMode}
                  type="button"
                >
                  Choose to be anon
                </button>
              </div>

              {!isAnonymousMode && (
                <>
                  <div className="mt-5 flex flex-wrap justify-center gap-4">
                    <button
                      className={`cyber-btn text-sm ${
                        kycProvider === "mock" ? "cyber-btn-primary" : ""
                      }`}
                      onClick={() => {
                        setKYCProvider("mock");
                        setKYCMessage(null);
                      }}
                      type="button"
                    >
                      Mock
                    </button>
                    <button
                      className={`cyber-btn text-sm ${
                        kycProvider === "zkpassport" ? "cyber-btn-primary" : ""
                      }`}
                      onClick={() => {
                        setKYCProvider("zkpassport");
                        setKYCMessage(null);
                      }}
                      type="button"
                    >
                      ZKPassport
                    </button>
                  </div>

                  <div className="mt-10 flex flex-wrap justify-center gap-4">
                    <button
                      className="btn btn-primary !px-4 !py-2.5 !text-xs"
                      onClick={startZKPassportFlow}
                      type="button"
                      disabled={isStartingZKPassport}
                    >
                      {isStartingZKPassport
                        ? "Starting session..."
                        : "Start ZKPassport QR"}
                    </button>
                    {zkPassportSession?.sessionId &&
                      zkPassportSession.source === "backend" && (
                        <button
                          className="btn btn-secondary !px-4 !py-2.5 !text-xs"
                          onClick={refreshZKPassportSession}
                          type="button"
                        >
                          Refresh session
                        </button>
                    )}
                  </div>
                  <p className="mt-3 text-center text-xs text-slate-400">
                    Session API: {zkPassportApiBase}
                  </p>

                  {kycProvider === "zkpassport" && (
                    <div className="mt-3 rounded-lg border border-slate-600/60 bg-slate-950/65 p-3">
                      {zkPassportSession ? (
                        <>
                          <div className="flex flex-wrap items-center justify-center gap-2 text-center">
                            <p className="text-xs text-slate-400">
                              Session:{" "}
                              <span className="font-mono text-slate-300">
                                {zkPassportSession.sessionId}
                              </span>
                            </p>
                            <span
                              className={`inline-flex rounded px-2 py-0.5 text-[11px] font-medium uppercase tracking-[0.08em] ${getSessionTone(zkPassportSession.status).badge}`}
                            >
                              {zkPassportSession.status}
                            </span>
                          </div>
                          <p className="mt-2 text-xs text-slate-400">
                            Bridge:{" "}
                            {zkPassportSession.bridgeConnected
                              ? "connected"
                              : "waiting"}
                            {" · "}Request:{" "}
                            {zkPassportSession.requestReceived
                              ? "received"
                              : "pending scan"}
                          </p>

                          {zkPassportSession.qrCodeUrl && !isQrImageUnavailable ? (
                            <div className="mt-3 flex flex-col items-center gap-2 rounded-md border border-slate-700/60 bg-slate-900/70 p-3">
                              <img
                                src={zkPassportSession.qrCodeUrl}
                                alt="ZKPassport verification QR code"
                                className="h-36 w-36 rounded-md border border-slate-700/60 bg-white p-2"
                                onError={() => {
                                  setIsQrImageUnavailable(true);
                                }}
                              />
                              <p className="text-center text-xs text-slate-300">
                                Scan to complete KYC verification with ZKPassport.
                              </p>
                              {zkPassportSession.deepLinkUrl && (
                                <a
                                  href={zkPassportSession.deepLinkUrl}
                                  target="_blank"
                                  rel="noreferrer"
                                  className="text-xs font-medium text-blue-200 transition hover:text-blue-100"
                                >
                                  Open deep link
                                </a>
                              )}
                            </div>
                          ) : isQrImageUnavailable ? (
                            <div className="mt-3 rounded-md border border-amber-300/45 bg-amber-500/10 p-3">
                              <p className="text-center text-xs text-amber-100">
                                QR image provider is blocked or unavailable. Use the
                                deep link below to continue.
                              </p>
                              {zkPassportSession.deepLinkUrl && (
                                <a
                                  href={zkPassportSession.deepLinkUrl}
                                  target="_blank"
                                  rel="noreferrer"
                                  className="mt-2 inline-flex text-xs font-medium text-blue-200 transition hover:text-blue-100"
                                >
                                  Open deep link
                                </a>
                              )}
                            </div>
                          ) : (
                            <p className="mt-2 text-xs text-slate-400">
                              Start a session to generate a QR code for passport
                              proof.
                            </p>
                          )}

                          {zkPassportSession.expiresAt && (
                            <p className="mt-2 text-xs text-slate-400">
                              Expires:{" "}
                              {new Date(
                                zkPassportSession.expiresAt,
                              ).toLocaleString()}
                            </p>
                          )}
                          {zkPassportSession.message && (
                            <p
                              className={`mt-2 text-xs ${getSessionTone(zkPassportSession.status).text}`}
                            >
                              {zkPassportSession.message}
                            </p>
                          )}
                        </>
                      ) : (
                        <p className="text-xs text-slate-400">
                          No active session yet. Click "Start ZKPassport QR".
                        </p>
                      )}
                    </div>
                  )}

                  <label className="mt-5 block text-center text-xs text-slate-400">
                    Optional proof JSON
                  </label>
                  <textarea
                    className="cyber-input mt-2 min-h-[116px] text-xs"
                    placeholder='{"proofs":[],"queryResult":{}}'
                    value={kycProofText}
                    onChange={(event) => {
                      setKYCProofText(event.target.value);
                      setKYCMessage(null);
                    }}
                  />

                  <div className="dashboard-kyc-cta mt-6 flex flex-wrap justify-center gap-4">
                    <button
                      className="btn btn-primary !px-4 !py-2.5 !text-xs"
                      onClick={saveKYCConfig}
                      type="button"
                    >
                      Save KYC config
                    </button>
                    <button
                      className="btn btn-secondary !px-4 !py-2.5 !text-xs"
                      onClick={clearKYCProof}
                      type="button"
                    >
                      Clear proof
                    </button>
                  </div>
                </>
              )}

              {isAnonymousMode && (
                <p className="mt-6 text-center text-xs text-slate-300">
                  Anonymous mode is active. ZKPassport session/proof setup is hidden
                  and workflow requests skip KYC writes.
                </p>
              )}

              {kycMessage && (
                <p
                  aria-live="polite"
                  role="status"
                  className={`mt-3 text-xs ${
                    kycMessage.includes("invalid") ||
                    kycMessage.includes("must be valid") ||
                    kycMessage.toLowerCase().includes("unable") ||
                    kycMessage.toLowerCase().includes("failed")
                      ? "text-rose-300"
                      : "text-emerald-300"
                  }`}
                >
                  {kycMessage}
                </p>
              )}
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

const normalizeZKPassportSession = (
  session: DashboardZKSession,
): DashboardZKSession => {
  const raw = session as DashboardZKSession & {
    qrCode?: string;
    qrUrl?: string;
    deepLink?: string;
    url?: string;
    data?: {
      qrCodeUrl?: string;
      qrCode?: string;
      qrUrl?: string;
      deepLinkUrl?: string;
      deepLink?: string;
      url?: string;
    };
  };
  const nested = raw.data || {};
  const deepLinkUrl =
    session.deepLinkUrl ||
    raw.deepLink ||
    raw.url ||
    nested.deepLinkUrl ||
    nested.deepLink ||
    nested.url;
  const qrCodeUrl =
    session.qrCodeUrl ||
    raw.qrCode ||
    raw.qrUrl ||
    nested.qrCodeUrl ||
    nested.qrCode ||
    nested.qrUrl ||
    buildZKPassportQrUrl(deepLinkUrl);

  return {
    ...session,
    qrCodeUrl,
    deepLinkUrl,
  };
};

const buildZKPassportQrUrl = (deepLinkUrl?: string): string | undefined => {
  if (!deepLinkUrl) {
    return undefined;
  }
  const encodedUrl = encodeURIComponent(deepLinkUrl);
  return `https://api.qrserver.com/v1/create-qr-code/?size=320x320&data=${encodedUrl}`;
};

// Stat Card Component
const StatCard: React.FC<{
  label: string;
  value: number;
  color: string;
  icon: string;
}> = ({ label, value, color, icon }) => (
  <div className="cyber-card dashboard-panel dashboard-kpi-card group relative mx-auto w-full max-w-[13.5rem] overflow-hidden p-6 text-center">
    <div className="absolute top-0 right-0 p-3 opacity-20 transition-opacity group-hover:opacity-40">
      <svg
        className="w-9 h-9"
        style={{ color }}
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={1.5}
          d={icon}
        />
      </svg>
    </div>
    <p className="dashboard-kpi-label mb-2">{label}</p>
    <p className="number-pill number-pill-lg mx-auto mt-1" style={{ color }}>
      {value.toString().padStart(2, "0")}
    </p>
  </div>
);

// Action Button Component
const ActionButton: React.FC<{
  onClick: () => void;
  title: string;
  description: string;
  icon: string;
}> = ({ onClick, title, description, icon }) => (
  <button
    type="button"
    onClick={onClick}
    className="btn btn-secondary quick-action-btn !w-full !rounded-lg !px-4 !py-4 text-center"
  >
    <div className="quick-action-content">
      <div className="quick-action-icon">
        <svg
          className="h-4 w-4 text-blue-200"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d={icon}
          />
        </svg>
      </div>
      <div className="quick-action-copy">
        <p className="dashboard-action-title text-sm">{title}</p>
        <p className="dashboard-action-body text-xs">{description}</p>
      </div>
    </div>
  </button>
);

// Status badge component
const StatusBadge: React.FC<{ house: House }> = ({ house }) => {
  if (house.rental?.isActive) {
    return (
      <span className="inline-flex items-center rounded border border-indigo-300/45 bg-indigo-400/15 px-2.5 py-0.5 text-xs font-medium text-indigo-200">
        Rented
      </span>
    );
  }

  if (house.listing?.listingType === "for_sale") {
    return (
      <span className="inline-flex items-center rounded border border-emerald-300/45 bg-emerald-400/15 px-2.5 py-0.5 text-xs font-medium text-emerald-200">
        For sale
      </span>
    );
  }

  if (house.listing?.listingType === "for_rent") {
    return (
      <span className="inline-flex items-center rounded border border-amber-300/45 bg-amber-400/15 px-2.5 py-0.5 text-xs font-medium text-amber-200">
        For rent
      </span>
    );
  }

  return (
    <span className="inline-flex items-center rounded border border-slate-500/55 bg-slate-600/20 px-2.5 py-0.5 text-xs font-medium text-slate-300">
      Owned
    </span>
  );
};
