import React, { useEffect, useMemo, useState } from "react";
import toast from "react-hot-toast";
import { useSearchParams } from "react-router-dom";

import { apiClient } from "@shared/utils/api";
import { useAuth } from "../components/AuthProvider";
import { useUXMode } from "../components/UXModeProvider";
import {
  readLatestClaimKeyHash,
  saveLatestClaimKeyHash,
} from "../utils/claimKeyStorage";

const isBytes32 = (value: string) => /^0x[a-fA-F0-9]{64}$/.test(value);
const needsAuthRetry = (message: string): boolean => {
  const normalized = message.toLowerCase();
  return (
    normalized.includes("authentication is required") ||
    normalized.includes("unauthorized") ||
    normalized.includes("http error: 401")
  );
};

export const ClaimKeyPage: React.FC = () => {
  const { walletAddress, signInWithEthereum } = useAuth();
  const { mode } = useUXMode();
  const [searchParams] = useSearchParams();
  const [keyHash, setKeyHash] = useState("");
  const [loading, setLoading] = useState(false);
  const [encryptedKey, setEncryptedKey] = useState<string | null>(null);

  const keyHashValid = useMemo(() => isBytes32(keyHash.trim()), [keyHash]);
  const showKeyHashError = keyHash.trim().length > 0 && !keyHashValid;

  useEffect(() => {
    const fromQuery = String(searchParams.get("keyHash") || "").trim();
    if (isBytes32(fromQuery)) {
      setKeyHash(fromQuery);
      saveLatestClaimKeyHash(fromQuery);
      return;
    }

    const fromStorage = readLatestClaimKeyHash();
    if (fromStorage) {
      setKeyHash(fromStorage);
    }
  }, [searchParams]);

  const handleClaim = async () => {
    const trimmed = keyHash.trim();
    if (!isBytes32(trimmed)) {
      toast.error("Invalid key hash. Expected bytes32: 0x + 64 hex chars");
      return;
    }

    try {
      setLoading(true);
      setEncryptedKey(null);
      saveLatestClaimKeyHash(trimmed);

      let resp = await apiClient.claimKey(
        trimmed,
        walletAddress || undefined,
      );
      if (!resp.success && needsAuthRetry(String(resp.message || ""))) {
        await signInWithEthereum();
        resp = await apiClient.claimKey(trimmed, walletAddress || undefined);
      }
      if (resp.success && resp.data?.encryptedKey) {
        setEncryptedKey(resp.data.encryptedKey);
        toast.success("Encrypted key retrieved");
      } else {
        toast.error(resp.message || "Failed to claim key");
      }
    } catch (e: any) {
      toast.error(e?.message || "Failed to claim key");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="page-shell page-shell-form">
      <div className="cyber-card overflow-hidden">
        <div className="form-header">
          <div>
            <h1 className="form-title">
              {mode === "degen"
                ? "Claim Encrypted Document Key"
                : "Get Your Private Access Key"}
            </h1>
            <p className="form-subtitle">
              {mode === "degen"
                ? "Paste a `keyHash` from a private sale or rental. The mediator returns the encrypted document key intended for you."
                : "Paste the key hash you received after buying or renting. We return your encrypted access key."}
            </p>
          </div>
          <div className="flex items-center">
            <span className="meta-chip">
              {walletAddress
                ? `${walletAddress.slice(0, 6)}...${walletAddress.slice(-4)}`
                : "Not connected"}
            </span>
          </div>
        </div>

        <div className="grid grid-cols-1 gap-5 p-7 md:p-8">
          <div className="form-field">
            <label htmlFor="claim-key-hash" className="block text-sm text-[var(--text-secondary)] mb-2 font-mono">
              Key Hash (bytes32)
            </label>
            <input
              id="claim-key-hash"
              value={keyHash}
              onChange={(e) => setKeyHash(e.target.value)}
              className="cyber-input font-mono"
              placeholder="0x..."
              spellCheck={false}
              aria-invalid={showKeyHashError}
              aria-describedby={`${keyHash.length === 0 ? "claim-key-hash-help" : ""}${showKeyHashError ? " claim-key-hash-error" : ""}`.trim() || undefined}
            />
            <div className="mt-2 text-xs font-mono">
              {keyHash.length === 0 ? (
                <span id="claim-key-hash-help" className="text-[var(--text-secondary)]">
                  Format: `0x` + 64 hex characters
                </span>
              ) : keyHashValid ? (
                <span className="status-active">VALID</span>
              ) : (
                <span className="status-error">INVALID</span>
              )}
            </div>
            {showKeyHashError && (
              <p id="claim-key-hash-error" className="form-error mt-1">
                Enter a valid key hash (`0x` + 64 hexadecimal characters).
              </p>
            )}
          </div>

          <div className="flex flex-col sm:flex-row gap-3">
            <button
              onClick={handleClaim}
              disabled={loading || !keyHashValid}
              className="cyber-btn cyber-btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading
                ? "Retrieving..."
                : mode === "degen"
                  ? "Retrieve Encrypted Key"
                  : "Get Access Key"}
            </button>
            <button
              onClick={() => {
                setKeyHash("");
                setEncryptedKey(null);
              }}
              className="cyber-btn"
              type="button"
            >
              Clear
            </button>
          </div>

          {encryptedKey && (
            <div className="mt-4 cyber-card p-6 border-[rgba(0,255,136,0.35)]">
              <div className="text-xs text-[var(--text-secondary)] font-mono">
                Encrypted Key
              </div>
              <pre className="mt-2 text-xs text-[var(--text-primary)] whitespace-pre-wrap break-all font-mono">
                {encryptedKey}
              </pre>
              <div className="mt-4 flex gap-2">
                <button
                  className="cyber-btn text-sm"
                  onClick={async () => {
                    await navigator.clipboard.writeText(encryptedKey);
                    toast.success("Copied");
                  }}
                >
                  Copy
                </button>
              </div>
              <p className="mt-4 text-sm text-[var(--text-secondary)]">
                {mode === "degen"
                  ? "Decrypt this key locally using your private key, then decrypt the documents from IPFS/offchain storage."
                  : "Keep this key private. Use your wallet/private key tools to decrypt and open your private documents."}
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
