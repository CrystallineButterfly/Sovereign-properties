/**
 * Privy Authentication Provider for Web
 * Handles authentication, wallet creation, and SIWE
 */

import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  useMemo,
} from "react";
import {
  type ConnectedWallet,
  type EIP1193Provider,
  PrivyProvider as PrivyLibProvider,
  useConnectWallet,
  useCreateWallet,
  usePrivy,
  useWallets,
} from "@privy-io/react-auth";
import { Link } from "react-router-dom";
import { SiweMessage } from "siwe";
import { ethers } from "ethers";
import { AUTH_EXPIRED_EVENT_NAME, apiClient } from "@shared/utils/api";
import type { User, PrivyUser, KYCStatus } from "@shared/types";
import { BrandMark } from "./BrandMark";

interface AuthContextType {
  user: User | null;
  privyUser: PrivyUser | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: () => Promise<void>;
  loginWithEmbeddedWallet: () => Promise<void>;
  connectExternalWallet: () => void;
  createEmbeddedWallet: () => Promise<void>;
  logout: () => Promise<void>;
  walletAddress: string | null;
  activeWalletClientType: string | null;
  embeddedWalletAddress: string | null;
  externalWalletAddresses: string[];
  wallets: ConnectedWallet[];
  selectActiveWallet: (address: string) => void;
  chainId: string | null;
  getEthereumProvider: () => Promise<EIP1193Provider>;
  switchNetwork: (targetChainId: number) => Promise<void>;
  signMessage: (message: string) => Promise<string>;
  signInWithEthereum: () => Promise<{ message: string; signature: string }>;
  setUserKycStatus: (status: KYCStatus) => void;
  refreshCurrentUser: () => Promise<void>;
  error: Error | null;
  loginStep:
    | "idle"
    | "authenticating"
    | "creating-wallet"
    | "signing"
    | "complete";
}

const AuthContext = createContext<AuthContextType | null>(null);

// Privy App ID from environment
const PRIVY_APP_ID = import.meta.env.VITE_PRIVY_APP_ID || "";
const LOCAL_SESSION_MARKER_KEY = "RWA_LOCAL_BROWSER_SESSION_ACTIVE";
const ACTIVE_WALLET_ADDRESS_KEY = "RWA_ACTIVE_WALLET_ADDRESS";
const API_AUTH_TOKEN_STORAGE_KEY = "RWA_API_AUTH_TOKEN";
const SIWE_NONCE_CHARS =
  "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";

const parseChainIdValue = (value: unknown): number | null => {
  if (typeof value === "number" && Number.isFinite(value)) return value;
  if (typeof value !== "string") return null;
  const trimmed = value.trim();
  if (trimmed.length === 0) return null;
  // Privy (or wallet providers) may return strings like "eip155:11155111".
  const parts = trimmed.split(":");
  const last = parts[parts.length - 1];
  const parsed = Number.parseInt(last, 10);
  return Number.isFinite(parsed) ? parsed : null;
};

const generateSiweNonce = (length = 24): string => {
  const size = Math.max(8, Math.floor(length));
  if (typeof crypto !== "undefined" && typeof crypto.getRandomValues === "function") {
    const bytes = new Uint8Array(size);
    crypto.getRandomValues(bytes);
    return Array.from(bytes, (value) =>
      SIWE_NONCE_CHARS[value % SIWE_NONCE_CHARS.length]
    ).join("");
  }

  let nonce = "";
  for (let index = 0; index < size; index += 1) {
    const randomIndex = Math.floor(Math.random() * SIWE_NONCE_CHARS.length);
    nonce += SIWE_NONCE_CHARS[randomIndex];
  }
  return nonce;
};

const buildChainConfig = (
  targetChainId: number,
): {
  chainId: string;
  chainName: string;
  nativeCurrency: { name: string; symbol: string; decimals: number };
  rpcUrls: string[];
  blockExplorerUrls?: string[];
} | null => {
  if (targetChainId === 1) {
    return {
      chainId: "0x1",
      chainName: "Ethereum Mainnet",
      nativeCurrency: { name: "Ether", symbol: "ETH", decimals: 18 },
      rpcUrls: ["https://cloudflare-eth.com"],
      blockExplorerUrls: ["https://etherscan.io"],
    };
  }

  if (targetChainId === 11155111) {
    return {
      chainId: "0xaa36a7",
      chainName: "Sepolia",
      nativeCurrency: { name: "Sepolia ETH", symbol: "SEP", decimals: 18 },
      rpcUrls: ["https://rpc.sepolia.org"],
      blockExplorerUrls: ["https://sepolia.etherscan.io"],
    };
  }

  if (targetChainId === 137) {
    return {
      chainId: "0x89",
      chainName: "Polygon",
      nativeCurrency: { name: "MATIC", symbol: "POL", decimals: 18 },
      rpcUrls: ["https://polygon-rpc.com"],
      blockExplorerUrls: ["https://polygonscan.com"],
    };
  }

  if (targetChainId === 31337) {
    const configuredRpc = String(import.meta.env.VITE_RPC_URL || "").trim();
    return {
      chainId: "0x7a69",
      chainName: "Anvil Local",
      nativeCurrency: { name: "Ether", symbol: "ETH", decimals: 18 },
      rpcUrls: [configuredRpc || "http://127.0.0.1:8545"],
      blockExplorerUrls: [],
    };
  }

  return null;
};

const clearWalletSessionStorage = () => {
  const clearMatchingKeys = (storage: Storage | undefined) => {
    if (!storage) return;
    const keys = Object.keys(storage);
    keys.forEach((key) => {
      const normalized = key.toLowerCase();
      if (
        normalized.includes("privy") ||
        normalized.includes("walletconnect") ||
        normalized.includes("wagmi")
      ) {
        storage.removeItem(key);
      }
    });
  };

  try {
    clearMatchingKeys(window.localStorage);
  } catch {
    // Ignore storage clear failures in restricted browser contexts.
  }

  try {
    clearMatchingKeys(window.sessionStorage);
  } catch {
    // Ignore storage clear failures in restricted browser contexts.
  }
};

const readStoredActiveWalletAddress = (): string | null => {
  try {
    const raw = window.localStorage.getItem(ACTIVE_WALLET_ADDRESS_KEY);
    if (!raw) return null;
    const trimmed = raw.trim();
    return trimmed || null;
  } catch {
    return null;
  }
};

const persistActiveWalletAddress = (address: string | null) => {
  try {
    if (!address) {
      window.localStorage.removeItem(ACTIVE_WALLET_ADDRESS_KEY);
      return;
    }
    window.localStorage.setItem(ACTIVE_WALLET_ADDRESS_KEY, address);
  } catch {
    // Ignore storage write issues in restricted browser contexts.
  }
};

const readStoredApiAuthToken = (): string | null => {
  try {
    const raw = window.localStorage.getItem(API_AUTH_TOKEN_STORAGE_KEY);
    if (!raw) return null;
    const trimmed = raw.trim();
    return trimmed || null;
  } catch {
    return null;
  }
};

const persistApiAuthToken = (token: string | null) => {
  try {
    if (!token) {
      window.localStorage.removeItem(API_AUTH_TOKEN_STORAGE_KEY);
      return;
    }
    window.localStorage.setItem(API_AUTH_TOKEN_STORAGE_KEY, token);
  } catch {
    // Ignore storage write issues in restricted browser contexts.
  }
};

const shouldForceSignOutOnNewLocalSession = (): boolean => {
  const persistAuthOnLocalhost =
    String(
      import.meta.env.VITE_PERSIST_AUTH_ON_LOCALHOST || "",
    ).toLowerCase() === "true";
  if (persistAuthOnLocalhost) return false;

  if (typeof window === "undefined") return false;
  const hostname = window.location.hostname;
  return hostname === "localhost" || hostname === "127.0.0.1";
};

const buildDemoUser = (
  walletAddress: string,
  chainId: number,
  email?: string,
): User => {
  const now = new Date();
  return {
    id: walletAddress,
    email: email || "operator@rwa.house",
    walletAddress,
    chainId,
    kycStatus: "unverified",
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
};

const parseUserDate = (value: unknown, fallback: Date): Date => {
  const parsed = new Date(String(value ?? ""));
  if (Number.isNaN(parsed.getTime())) {
    return fallback;
  }
  return parsed;
};

const asObjectRecord = (value: unknown): Record<string, unknown> => {
  if (!value || typeof value !== "object") {
    return {};
  }
  return value as Record<string, unknown>;
};

const normalizeUserPayload = (rawUser: unknown): User => {
  const userRecord = asObjectRecord(rawUser);
  const walletAddress = String(userRecord.walletAddress || "");
  const chainId = Number.parseInt(String(userRecord.chainId || 1), 10) || 1;
  const preferenceRecord = asObjectRecord(userRecord.preferences);
  const notificationRecord = asObjectRecord(preferenceRecord.notifications);
  const fallbackDate = new Date();
  const fallbackPreferences = buildDemoUser(
    walletAddress || "0x0000000000000000000000000000000000000000",
    chainId,
  ).preferences;
  return {
    ...userRecord,
    walletAddress,
    id: String(userRecord.id || walletAddress || ""),
    email: String(userRecord.email || "operator@rwa.house"),
    chainId,
    kycStatus: (
      userRecord.kycStatus === "verified" ||
      userRecord.kycStatus === "pending" ||
      userRecord.kycStatus === "rejected"
    )
      ? userRecord.kycStatus
      : "unverified",
    createdAt: parseUserDate(userRecord.createdAt, fallbackDate),
    lastLoginAt: parseUserDate(userRecord.lastLoginAt, fallbackDate),
    mfaEnabled: Boolean(userRecord.mfaEnabled),
    preferences:
      typeof preferenceRecord.theme === "string" &&
      typeof preferenceRecord.currency === "string" &&
      typeof preferenceRecord.language === "string"
        ? {
            theme: (
              preferenceRecord.theme === "light" ||
              preferenceRecord.theme === "system"
            )
              ? preferenceRecord.theme
              : "dark",
            currency: String(preferenceRecord.currency),
            language: String(preferenceRecord.language),
            autoPayEnabled: Boolean(preferenceRecord.autoPayEnabled),
            autoPayThreshold:
              Number(preferenceRecord.autoPayThreshold) || 1000,
            notifications: {
              email: Boolean(notificationRecord.email),
              push: Boolean(notificationRecord.push),
              sms: Boolean(notificationRecord.sms),
              transactions: Boolean(notificationRecord.transactions),
              bills: Boolean(notificationRecord.bills),
              security: Boolean(notificationRecord.security),
            },
          }
        : fallbackPreferences,
  };
};

// Auth Provider wrapper
export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  // Validate Privy App ID
  if (!PRIVY_APP_ID) {
    console.error("[AuthProvider] Privy app is not configured.");
    return (
      <div
        style={{
          minHeight: "100vh",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          background: "#060b14",
          color: "#fecdd3",
          fontFamily: "Inter, sans-serif",
          padding: "24px",
        }}
      >
        <div style={{ textAlign: "center" }}>
          <h2>Configuration Error</h2>
          <p>Privy app configuration is missing.</p>
          <p style={{ color: "#94a3b8", fontSize: "14px", marginTop: "10px" }}>
            Please configure Privy to continue.
          </p>
        </div>
      </div>
    );
  }

  return (
    <PrivyLibProvider
      appId={PRIVY_APP_ID}
      config={{
        loginMethods: ["email", "wallet", "google", "twitter"],
        appearance: {
          theme: "dark",
          accentColor: "#3b82f6",
        },
        embeddedWallets: {
          createOnLogin: "users-without-wallets",
          showWalletUIs: true,
        },
        legal: {
          termsAndConditionsUrl: "/terms",
          privacyPolicyUrl: "/privacy",
        },
      }}
    >
      <AuthContextProvider>{children}</AuthContextProvider>
    </PrivyLibProvider>
  );
};

// Internal context provider
const AuthContextProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const {
    login: privyLogin,
    logout: privyLogout,
    authenticated,
    ready,
    user: privyUser,
  } = usePrivy();

  const { connectWallet } = useConnectWallet();
  const { wallets, ready: walletsReady } = useWallets();
  const { createWallet } = useCreateWallet();
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const [isCreatingWallet, setIsCreatingWallet] = useState(false);
  const [isLoggingOut, setIsLoggingOut] = useState(false);
  const [enforceFreshLocalSession, setEnforceFreshLocalSession] =
    useState(false);
  const [walletProvisionAttempted, setWalletProvisionAttempted] =
    useState(false);
  const [pendingExternalSelection, setPendingExternalSelection] =
    useState(false);
  const [preferredWalletAddress, setPreferredWalletAddress] = useState<
    string | null
  >(null);
  const [loginStep, setLoginStep] = useState<
    "idle" | "authenticating" | "creating-wallet" | "signing" | "complete"
  >("idle");

  const setApiAuthToken = useCallback((token: string | null) => {
    const normalized = String(token || "").trim();
    apiClient.setAuthToken(normalized);
    persistApiAuthToken(normalized || null);
  }, []);

  useEffect(() => {
    const storedToken = readStoredApiAuthToken();
    if (storedToken) {
      apiClient.setAuthToken(storedToken);
    }
  }, []);

  useEffect(() => {
    const handleAuthExpired = () => {
      setUser(null);
      setApiAuthToken(null);
      setLoginStep("idle");
      setError(new Error("Session expired. Please sign in again."));
    };

    window.addEventListener(
      AUTH_EXPIRED_EVENT_NAME,
      handleAuthExpired as EventListener,
    );
    return () => {
      window.removeEventListener(
        AUTH_EXPIRED_EVENT_NAME,
        handleAuthExpired as EventListener,
      );
    };
  }, [setApiAuthToken]);

  const embeddedWallet = wallets.find(
    (candidate) => candidate.walletClientType === "privy",
  );
  const externalWallets = wallets.filter(
    (candidate) => candidate.walletClientType !== "privy",
  );
  const wallet = useMemo(() => {
    const normalizedPreferred = preferredWalletAddress?.toLowerCase() || null;
    if (normalizedPreferred) {
      const preferredMatch = wallets.find(
        (candidate) => candidate.address.toLowerCase() === normalizedPreferred,
      );
      if (preferredMatch) return preferredMatch;
    }

    const latestExternalWallet = externalWallets[externalWallets.length - 1];
    return latestExternalWallet || embeddedWallet || wallets[0];
  }, [preferredWalletAddress, wallets, externalWallets, embeddedWallet]);
  const walletAddress = wallet?.address || null;
  const embeddedWalletAddress = embeddedWallet?.address || null;
  const externalWalletAddresses = externalWallets.map(
    (candidate) => candidate.address,
  );
  const activeWalletClientType = wallet?.walletClientType || null;
  const chainId = wallet?.chainId || null;

  useEffect(() => {
    if (typeof window === "undefined") return;
    setPreferredWalletAddress(readStoredActiveWalletAddress());
  }, []);

  useEffect(() => {
    if (walletAddress) {
      persistActiveWalletAddress(walletAddress);
    }
  }, [walletAddress]);

  useEffect(() => {
    if (wallets.length === 0) {
      return;
    }

    const normalizedPreferred = preferredWalletAddress?.toLowerCase() || null;
    const stillExists = normalizedPreferred
      ? wallets.some(
          (candidate) => candidate.address.toLowerCase() === normalizedPreferred,
        )
      : false;
    if (stillExists) {
      return;
    }

    const fallback =
      externalWallets[externalWallets.length - 1]?.address
      || embeddedWallet?.address
      || wallets[0]?.address
      || null;
    setPreferredWalletAddress(fallback);
  }, [wallets, embeddedWallet, externalWallets, preferredWalletAddress]);

  useEffect(() => {
    if (!pendingExternalSelection) return;
    const latestExternal = externalWallets[externalWallets.length - 1];
    if (!latestExternal) return;
    setPreferredWalletAddress(latestExternal.address);
    setPendingExternalSelection(false);
  }, [pendingExternalSelection, externalWallets]);

  useEffect(() => {
    if (!shouldForceSignOutOnNewLocalSession()) {
      return;
    }

    try {
      const existingMarker = window.sessionStorage.getItem(
        LOCAL_SESSION_MARKER_KEY,
      );
      if (!existingMarker) {
        window.sessionStorage.setItem(LOCAL_SESSION_MARKER_KEY, "1");
        setEnforceFreshLocalSession(true);
      }
    } catch {
      setEnforceFreshLocalSession(true);
    }
  }, []);

  useEffect(() => {
    if (!enforceFreshLocalSession || !ready) {
      return;
    }

    let cancelled = false;
    const enforceSignOut = async () => {
      setIsLoading(true);
      clearWalletSessionStorage();
      setApiAuthToken(null);
      setUser(null);
      setWalletProvisionAttempted(false);
      setLoginStep("idle");

      if (authenticated) {
        try {
          await privyLogout();
        } catch (logoutErr) {
          console.warn(
            "[Auth] Localhost session reset could not complete Privy logout:",
            logoutErr,
          );
        }
      }

      if (!cancelled) {
        setEnforceFreshLocalSession(false);
        setIsLoading(false);
      }
    };

    enforceSignOut().catch((err) => {
      console.warn("[Auth] Failed to enforce fresh localhost session:", err);
      if (!cancelled) {
        setEnforceFreshLocalSession(false);
        setIsLoading(false);
      }
    });

    return () => {
      cancelled = true;
    };
  }, [
    enforceFreshLocalSession,
    ready,
    authenticated,
    privyLogout,
    setApiAuthToken,
  ]);

  // Sign message using wallet
  const signMessage = useCallback(
    async (message: string): Promise<string> => {
      if (!wallet) {
        throw new Error("No wallet available");
      }

      try {
        const provider = await wallet.getEthereumProvider();
        const normalizeAccounts = (value: unknown): string[] => {
          if (!Array.isArray(value)) return [];
          return value
            .map((entry) => String(entry || "").trim())
            .filter((entry) => entry.length > 0);
        };

        let accounts = normalizeAccounts(
          await provider.request({ method: "eth_accounts" }),
        );
        if (accounts.length === 0) {
          accounts = normalizeAccounts(
            await provider.request({ method: "eth_requestAccounts" }),
          );
        }

        const activeAccount = String(
          accounts[0] || wallet.address || "",
        ).trim();
        if (!activeAccount) {
          throw new Error("No wallet account available for message signing.");
        }

        const trySign = async (
          method: string,
          params: unknown[],
        ): Promise<string | null> => {
          try {
            const result = await provider.request({ method, params });
            const signature = String(result || "").trim();
            return signature.length > 0 ? signature : null;
          } catch {
            return null;
          }
        };

        const personalSignPrimary = await trySign("personal_sign", [
          message,
          activeAccount,
        ]);
        if (personalSignPrimary) {
          return personalSignPrimary;
        }

        const personalSignFallback = await trySign("personal_sign", [
          activeAccount,
          message,
        ]);
        if (personalSignFallback) {
          return personalSignFallback;
        }

        const ethSignFallback = await trySign("eth_sign", [
          activeAccount,
          ethers.hexlify(ethers.toUtf8Bytes(message)),
        ]);
        if (ethSignFallback) {
          return ethSignFallback;
        }

        const browserProvider = new ethers.BrowserProvider(provider);
        const signer = await browserProvider.getSigner(activeAccount);
        return await signer.signMessage(message);
      } catch (err) {
        console.error("Signing error:", err);
        throw new Error("Failed to sign message");
      }
    },
    [wallet],
  );

  // Sign-In with Ethereum (SIWE)
  const signInWithEthereum = useCallback(async () => {
    if (!walletAddress || !chainId) {
      throw new Error("Wallet not connected");
    }

    try {
      const chainIdNum = parseChainIdValue(chainId) || 1;

      // Create SIWE message
      const message = new SiweMessage({
        domain: window.location.host,
        address: walletAddress,
        statement: "Sign in to PropMeSovereignty Platform",
        uri: window.location.origin,
        version: "1",
        chainId: chainIdNum,
        nonce: await fetchNonce(),
      });

      const messageToSign = message.prepareMessage();
      const signature = await signMessage(messageToSign);

      try {
        const response = await apiClient.verifyWallet(
          walletAddress,
          signature,
          messageToSign,
        );
        if (!response.success || !response.data?.token || !response.data?.user) {
          throw new Error(response.message || "Wallet verification failed.");
        }

        setApiAuthToken(response.data.token);
        setUser(normalizeUserPayload(response.data.user));
      } catch (error) {
        setApiAuthToken(null);
        setUser(null);
        throw new Error(
          `Failed to establish authenticated session: ${
            error instanceof Error ? error.message : "Unknown verification error."
          }`,
        );
      }

      return { message: messageToSign, signature };
    } catch (err) {
      console.warn(
        "[Auth] SIWE unavailable; secure backend session was not established.",
        err,
      );
      throw err;
    }
  }, [walletAddress, chainId, signMessage, setApiAuthToken]);

  // Fetch nonce from server
  const fetchNonce = async (): Promise<string> => {
    // SIWE nonce must be at least 8 alphanumeric characters (EIP-4361).
    return generateSiweNonce(24);
  };

  // Login handler focused on embedded Privy wallet onboarding.
  const loginWithEmbeddedWallet = useCallback(async () => {
    try {
      setError(null);
      setIsLoading(true);
      setLoginStep("authenticating");
      setWalletProvisionAttempted(false);

      await privyLogin({
        loginMethods: ["email", "google", "twitter"],
      });

      // Note: Wallet creation will be handled by the useEffect below
      // The user will be signed in once wallet is ready
    } catch (err) {
      setError(err as Error);
      console.error("Login error:", err);
      setLoginStep("idle");
    }
  }, [privyLogin]);

  const login = useCallback(async () => {
    await loginWithEmbeddedWallet();
  }, [loginWithEmbeddedWallet]);

  const connectExternalWallet = useCallback(() => {
    setError(null);
    setPendingExternalSelection(true);

    if (!authenticated) {
      setIsLoading(true);
      setLoginStep("authenticating");
      setWalletProvisionAttempted(false);
      privyLogin({ loginMethods: ["wallet"] });
      return;
    }

    connectWallet();
  }, [authenticated, connectWallet, privyLogin]);

  const selectActiveWallet = useCallback(
    (address: string) => {
      const normalizedAddress = address.trim().toLowerCase();
      if (!normalizedAddress) {
        throw new Error("Wallet address is required to select an active wallet.");
      }

      const matchedWallet = wallets.find(
        (candidate) => candidate.address.toLowerCase() === normalizedAddress,
      );
      if (!matchedWallet) {
        throw new Error("Selected wallet is not currently connected.");
      }

      setPreferredWalletAddress(matchedWallet.address);
      persistActiveWalletAddress(matchedWallet.address);
    },
    [wallets],
  );

  const createEmbeddedWallet = useCallback(async () => {
    if (!authenticated) {
      await loginWithEmbeddedWallet();
      return;
    }

    setError(null);
    setLoginStep("creating-wallet");
    try {
      await createWallet();
    } catch (err) {
      setError(err as Error);
      console.error("Create embedded wallet failed:", err);
    }
  }, [authenticated, createWallet, loginWithEmbeddedWallet]);

  const setUserKycStatus = useCallback((status: KYCStatus) => {
    setUser((previous) => {
      if (!previous) {
        return previous;
      }
      return {
        ...previous,
        kycStatus: status,
      };
    });
  }, []);

  const refreshCurrentUser = useCallback(async () => {
    if (!authenticated || !walletAddress || !apiClient.hasAuthToken()) {
      return;
    }
    try {
      const response = await apiClient.getCurrentUser();
      if (response.success && response.data) {
        setUser(normalizeUserPayload(response.data));
      }
    } catch (err) {
      console.warn("[Auth] Unable to refresh current user:", err);
    }
  }, [authenticated, walletAddress]);

  const ensureEmbeddedWallet = useCallback(async () => {
    if (
      !ready ||
      !authenticated ||
      !walletsReady ||
      wallets.length > 0 ||
      isCreatingWallet ||
      walletProvisionAttempted
    ) {
      return;
    }

    setWalletProvisionAttempted(true);
    setIsCreatingWallet(true);
    setLoginStep("creating-wallet");

    try {
      console.log("[Auth] Creating embedded wallet...");
      await createWallet();
      console.log("[Auth] Embedded wallet created successfully");
    } catch (err: any) {
      const message = err?.message?.toLowerCase() || "";
      const expectedAlreadyExists =
        message.includes("already has") ||
        message.includes("exists") ||
        message.includes("conflict");

      if (expectedAlreadyExists) {
        console.log("[Auth] Wallet already exists for this user");
      } else {
        console.error("[Auth] Embedded wallet creation failed:", err);
        setError(
          new Error(
            "Wallet creation failed. This might be due to:\n" +
              "1. Browser blocking popups\n" +
              "2. Network connectivity issues\n" +
              "3. Please try refreshing the page and logging in again",
          ),
        );
      }
    } finally {
      setIsCreatingWallet(false);
    }
  }, [
    ready,
    authenticated,
    walletsReady,
    wallets.length,
    isCreatingWallet,
    walletProvisionAttempted,
    createWallet,
  ]);

  // Logout handler
  const logout = useCallback(async () => {
    setIsLoggingOut(true);
    setIsLoading(true);
    try {
      await apiClient.logout();
    } catch (err) {
      console.warn(
        "[Auth] API logout failed, continuing with local wallet logout:",
        err,
      );
    }

    try {
      wallets.forEach((connectedWallet) => {
        try {
          connectedWallet.disconnect?.();
        } catch (disconnectErr) {
          console.warn("[Auth] Wallet disconnect failed:", disconnectErr);
        }
      });
      await privyLogout();
    } catch (err) {
      console.error("Privy logout error:", err);
      clearWalletSessionStorage();
    } finally {
      setUser(null);
      setApiAuthToken(null);
      setWalletProvisionAttempted(false);
      setLoginStep("idle");
      setPendingExternalSelection(false);
      setPreferredWalletAddress(null);
      persistActiveWalletAddress(null);
      setIsLoading(false);
      setIsLoggingOut(false);
    }
  }, [privyLogout, setApiAuthToken, wallets]);

  useEffect(() => {
    if (ready && !authenticated) {
      setUser(null);
      setApiAuthToken(null);
      setLoginStep("idle");
    }
  }, [ready, authenticated, setApiAuthToken]);

  const switchNetwork = useCallback(
    async (targetChainId: number): Promise<void> => {
      if (!Number.isInteger(targetChainId) || targetChainId <= 0) {
        throw new Error("Target chain ID must be a positive integer.");
      }

      const switchErrors: unknown[] = [];
      const hexChainId = `0x${targetChainId.toString(16)}`;

      if (wallet?.switchChain) {
        try {
          await wallet.switchChain(targetChainId);
          return;
        } catch (err) {
          switchErrors.push(err);
          console.warn(
            "[Auth] Privy switchChain failed, falling back to provider:",
            err,
          );
        }
      }

      const ethereum = (window as Window & { ethereum?: any }).ethereum;
      if (!ethereum?.request) {
        if (wallet?.walletClientType === "privy") {
          if (targetChainId === 31337) {
            throw new Error(
              "Privy embedded wallets cannot switch to local Anvil (31337). Connect an external wallet for local testing.",
            );
          }

          const priorError = switchErrors[0];
          const priorMessage =
            priorError instanceof Error
              ? priorError.message
              : priorError
                ? String(priorError)
                : "";

          throw new Error(
            priorMessage ||
              "Privy wallet could not switch networks. Open the Privy wallet UI to switch manually.",
          );
        }

        throw new Error("No wallet provider available to switch networks.");
      }

      try {
        await ethereum.request({
          method: "wallet_switchEthereumChain",
          params: [{ chainId: hexChainId }],
        });
        return;
      } catch (err: any) {
        if (err?.code !== 4902) {
          switchErrors.push(err);
          throw err;
        }
      }

      const chainConfig = buildChainConfig(targetChainId);
      if (!chainConfig) {
        throw new Error(
          `Wallet is missing chain ${targetChainId} and no add-chain config exists for it.`,
        );
      }

      try {
        await ethereum.request({
          method: "wallet_addEthereumChain",
          params: [chainConfig],
        });
        await ethereum.request({
          method: "wallet_switchEthereumChain",
          params: [{ chainId: hexChainId }],
        });
      } catch (addError) {
        switchErrors.push(addError);
        throw new Error(
          `Failed to add chain ${targetChainId} to wallet. ${String(addError)}`,
        );
      }

      if (switchErrors.length > 0) {
        console.warn("[Auth] Prior switch attempts failed:", switchErrors);
      }
    },
    [wallet],
  );

  const getEthereumProvider =
    useCallback(async (): Promise<EIP1193Provider> => {
      if (wallet) {
        return wallet.getEthereumProvider();
      }

      const externalProvider = (window as Window & { ethereum?: any }).ethereum;
      if (externalProvider?.request) {
        return externalProvider;
      }

      throw new Error(
        "No wallet provider available. Sign in with Privy or connect an external wallet.",
      );
    }, [wallet]);

  useEffect(() => {
    if (!authenticated) {
      setWalletProvisionAttempted(false);
    }
  }, [authenticated]);

  // Effect to handle wallet connection changes and SIWE
  useEffect(() => {
    const performSignIn = async () => {
      const userWalletAddress = user?.walletAddress || "";
      const activeWalletAddress = walletAddress?.toLowerCase() || "";
      const shouldSignIn =
        !user || userWalletAddress.toLowerCase() !== activeWalletAddress;

      if (ready && walletsReady && authenticated && walletAddress && wallet && shouldSignIn) {
        setLoginStep("signing");
        setError(null);
        try {
          await signInWithEthereum();
          setLoginStep("complete");
        } catch (err) {
          console.warn("[Auth] SIWE failed; keeping backend session locked.", err);
          setApiAuthToken(null);
          setUser(null);
          setError(
            new Error(
              "Authentication failed. Please reconnect your wallet and sign the SIWE message.",
            ),
          );
          setLoginStep("idle");
        }
      }
    };

    performSignIn();
  }, [
    ready,
    walletsReady,
    authenticated,
    walletAddress,
    wallet,
    user,
    signInWithEthereum,
    setApiAuthToken,
  ]);

  useEffect(() => {
    refreshCurrentUser().catch((err) => {
      console.warn("[Auth] Current user refresh failed:", err);
    });
  }, [refreshCurrentUser]);

  // Reliability fallback: if authenticated but no wallet is present, explicitly create one.
  useEffect(() => {
    ensureEmbeddedWallet().catch((err) => {
      console.warn("Wallet provisioning fallback failed:", err);
    });
  }, [ensureEmbeddedWallet]);

  // Effect to handle initial loading state
  useEffect(() => {
    // Still loading if Privy isn't ready yet
    if (!ready) {
      setIsLoading(true);
      return;
    }

    // Not authenticated = not loading (show login screen)
    if (!authenticated) {
      setIsLoading(false);
      return;
    }

    // Authenticated but waiting for wallet
    if (authenticated && walletsReady && wallets.length === 0) {
      // If we're actively creating wallet, stay loading
      if (isCreatingWallet) {
        setIsLoading(true);
        return;
      }
      // If we haven't attempted yet, we're still loading
      if (!walletProvisionAttempted) {
        setIsLoading(true);
        return;
      }
      // Attempted but no wallet - this is an error state
      if (walletProvisionAttempted && wallets.length === 0) {
        console.error(
          "[Auth] Wallet provision attempted but no wallet available",
        );
        setError(
          new Error("Wallet creation failed. Please refresh and try again."),
        );
        setIsLoading(false);
        return;
      }
    }

    // All good - wallet is ready
    if (authenticated && wallets.length > 0) {
      setIsLoading(false);
      return;
    }

    // Default: not loading
    setIsLoading(false);
  }, [
    ready,
    authenticated,
    walletsReady,
    wallets.length,
    isCreatingWallet,
    walletProvisionAttempted,
  ]);

  const value: AuthContextType = {
    user,
    privyUser: privyUser as PrivyUser | null,
    isAuthenticated:
      !isLoggingOut &&
      authenticated &&
      !!walletAddress &&
      !!user?.walletAddress &&
      apiClient.hasAuthToken(),
    isLoading,
    login,
    loginWithEmbeddedWallet,
    connectExternalWallet,
    createEmbeddedWallet,
    logout,
    walletAddress,
    activeWalletClientType,
    embeddedWalletAddress,
    externalWalletAddresses,
    wallets,
    selectActiveWallet,
    chainId,
    getEthereumProvider,
    switchNetwork,
    signMessage,
    signInWithEthereum,
    setUserKycStatus,
    refreshCurrentUser,
    error,
    loginStep,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

// Hook to use auth context
export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return context;
};

// Protected route wrapper
export const ProtectedRoute: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const {
    isAuthenticated,
    isLoading,
    loginWithEmbeddedWallet,
    connectExternalWallet,
    error,
    loginStep,
  } = useAuth();

  if (isLoading) {
    const getLoadingMessage = () => {
      switch (loginStep) {
        case "authenticating":
          return "Authenticating...";
        case "creating-wallet":
          return "Creating your wallet...";
        case "signing":
          return "Securing your session...";
        default:
          return "Loading...";
      }
    };

    return (
      <div className="loading-container bg-[#060b14]">
        <div className="spinner" />
        <p className="text-sm font-medium text-slate-200">
          {getLoadingMessage()}
        </p>
        {loginStep === "creating-wallet" && (
          <p className="max-w-xs text-center text-xs text-slate-400">
            This may take a moment. Please keep this tab open while we create
            your secure wallet.
          </p>
        )}
      </div>
    );
  }

  if (!isAuthenticated) {
    return (
      <div className="relative min-h-screen overflow-hidden bg-[#060b14]">
        <div className="absolute inset-0 -z-10">
          <div className="orb orb-1" />
          <div className="orb orb-2" />
          <div className="grid-pattern" />
        </div>

        <header className="cyber-nav relative z-10 border-b border-slate-700/60">
          <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
            <div className="site-header-banner !justify-start">
              <Link
                to="/"
                className="site-brand-link site-brand-link--header"
                aria-label="Go to home"
              >
                <BrandMark
                  size="sm"
                  showWordmark={false}
                  logoAsset="home"
                  className="site-brand-logo site-brand-logo--header"
                />
              </Link>
            </div>
          </div>
        </header>

        <div
          className="mx-auto flex min-h-[calc(100vh-5rem)] max-w-6xl items-center
            justify-center px-4 py-10 sm:px-6 lg:px-8"
        >
          <div className="w-full text-center">
            <h1 className="mx-auto max-w-4xl text-4xl font-bold leading-tight text-slate-50 md:text-5xl">
              Secure sign-in for private real-estate workflows
            </h1>
            <p className="text-panel mx-auto mt-4 max-w-3xl text-base text-slate-300 md:text-lg">
              Choose your preferred path: create an embedded Privy wallet or
              connect an existing external wallet.
            </p>

            <section className="glass-card hero-panel mx-auto mt-8 max-w-2xl text-left">
              <h2 className="text-center text-xl font-semibold text-slate-50">
                Access your workspace
              </h2>
              <p className="text-panel mt-2 text-center text-sm text-slate-300">
                Authenticate to mint properties, manage documents, run KYC, and
                settle sales or rentals.
              </p>
              <div className="mt-5 grid gap-3">
                <button
                  onClick={loginWithEmbeddedWallet}
                  className="btn btn-primary w-full !py-3 !text-base"
                >
                  Sign in + create Privy wallet
                </button>
                <button
                  onClick={connectExternalWallet}
                  className="btn btn-secondary w-full !py-3 !text-base"
                >
                  Connect external wallet
                </button>
              </div>

              {error && (
                <div className="mt-4 rounded-lg border border-rose-400/45 bg-rose-500/10 p-3 text-sm text-rose-200 whitespace-pre-line">
                  {error.message}
                </div>
              )}

              <div className="mt-5 flex flex-wrap justify-center gap-2">
                <span className="meta-chip">Privy authentication</span>
                <span className="meta-chip">Embedded wallet support</span>
                <span className="meta-chip">External wallet support</span>
              </div>
            </section>
          </div>
        </div>
      </div>
    );
  }

  return <>{children}</>;
};
