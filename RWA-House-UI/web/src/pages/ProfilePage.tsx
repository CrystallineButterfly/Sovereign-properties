import React, { useCallback, useEffect, useMemo, useState } from "react";

import { useAuth } from "../components/AuthProvider";
import { HouseThumbnail } from "../components/HouseThumbnail";
import type { House } from "@shared/types";
import { apiClient } from "@shared/utils/api";

const shortenAddress = (address: string): string =>
  `${address.slice(0, 6)}...${address.slice(-4)}`;

const formatEth = (value: string | null): string =>
  value === null ? "—" : `${Number.parseFloat(value).toFixed(4)} ETH`;

export const ProfilePage: React.FC = () => {
  const {
    chainId,
    connectExternalWallet,
    createEmbeddedWallet,
    embeddedWalletAddress,
    user,
    walletAddress,
    wallets,
  } = useAuth();

  const [walletBalances, setWalletBalances] = useState<Record<string, string>>(
    {},
  );
  const [ownedHouses, setOwnedHouses] = useState<House[]>([]);
  const [ownedNfts, setOwnedNfts] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const allAddresses = useMemo(() => {
    const addresses = new Set<string>();
    if (walletAddress) {
      addresses.add(walletAddress);
    }
    wallets.forEach((wallet) => addresses.add(wallet.address));
    return Array.from(addresses);
  }, [walletAddress, wallets]);

  const primaryBalance = walletAddress
    ? walletBalances[walletAddress] || null
    : null;
  const walletDisplay = walletAddress ? shortenAddress(walletAddress) : "—";
  const isWalletKYCVerified = user?.kycStatus === "verified";
  const walletKYCBadgeText = isWalletKYCVerified
    ? "KYC verified"
    : "KYC verification required";

  const refreshProfile = useCallback(async () => {
    if (!walletAddress) {
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const housesResponse = await apiClient.getHouses(walletAddress);
      const balances = await Promise.all(
        allAddresses.map(async (address) => {
          const response = await apiClient.getNativeBalance(address);
          return [address, response.success && response.data ? response.data : "0"] as const;
        }),
      );
      setWalletBalances(Object.fromEntries(balances));

      if (housesResponse.success && housesResponse.data) {
        setOwnedHouses(housesResponse.data);
        setOwnedNfts(housesResponse.data.length);
      } else {
        setOwnedHouses([]);
        setOwnedNfts(0);
        setError(housesResponse.message || "Failed to load wallet profile data.");
      }
    } catch (err) {
      console.error("Failed to load profile:", err);
      setError("Failed to load wallet profile data.");
    } finally {
      setIsLoading(false);
    }
  }, [allAddresses, walletAddress]);

  useEffect(() => {
    refreshProfile().catch((err) => {
      console.error("Profile refresh error:", err);
    });
  }, [refreshProfile]);

  return (
    <div className="min-h-screen bg-[#060b14] pb-12">
      <main className="page-shell page-shell-tight workspace-surface">
        <header className="page-header">
          <div className="cyber-card dashboard-panel text-panel mx-auto max-w-5xl !rounded-2xl !p-5 md:!p-6 text-center">
            <h1 className="dashboard-title">Wallet Profile</h1>
            <p className="dashboard-action-body mx-auto mt-2 max-w-3xl text-sm md:text-base">
              Review wallet access, balances, and tokenized property holdings.
            </p>
            <p className="dashboard-meta mx-auto mt-3 inline-flex items-center gap-2 rounded-full px-3 py-2 text-xs">
              <span className="h-2 w-2 rounded-full bg-emerald-400" />
              Wallet: {walletDisplay}
              <span
                className={`inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[11px] font-semibold uppercase tracking-[0.08em] ${
                  isWalletKYCVerified
                    ? "border-emerald-300/50 bg-emerald-500/15 text-emerald-200"
                    : "border-amber-300/50 bg-amber-500/15 text-amber-200"
                }`}
              >
                {isWalletKYCVerified ? "✅" : "⚠️"} {walletKYCBadgeText}
              </span>
            </p>
          </div>
        </header>

        <section className="section-shell space-y-6">
          {error && (
            <div
              role="status"
              aria-live="polite"
              className="text-panel border border-rose-400/45 p-3 text-sm text-rose-200"
            >
              {error}
            </div>
          )}

          <div className="cyber-card dashboard-panel text-panel space-y-4 p-5">
            <div className="flex items-center justify-between">
              <p className="dashboard-section-title text-base">
                Wallet snapshot
              </p>
              <div className="flex items-center gap-2">
                <span
                  className={`inline-flex items-center gap-1 rounded-full border px-2 py-0.5 text-[11px] font-semibold uppercase tracking-[0.08em] ${
                    isWalletKYCVerified
                      ? "border-emerald-300/50 bg-emerald-500/15 text-emerald-200"
                      : "border-amber-300/50 bg-amber-500/15 text-amber-200"
                  }`}
                >
                  {isWalletKYCVerified ? "✅" : "⚠️"} {walletKYCBadgeText}
                </span>
                <span className="text-xs text-slate-400">
                  Refreshed on demand
                </span>
              </div>
            </div>

            <div className="grid gap-4 md:grid-cols-4">
              <InfoCard
                label="Primary wallet"
                value={walletDisplay}
              />
              <InfoCard
                label="Primary balance"
                value={isLoading ? "Loading..." : formatEth(primaryBalance)}
              />
              <InfoCard
                label="Property NFTs"
                value={isLoading ? "Loading..." : `${ownedNfts}`}
              />
              <InfoCard
                label="Connected wallets"
                value={`${allAddresses.length}`}
              />
            </div>
          </div>

          <div className="cyber-card dashboard-panel text-panel space-y-5 p-5">
            <div className="flex items-center justify-between gap-3">
              <p className="dashboard-section-title">Wallets</p>
              <span className="dashboard-meta text-xs font-mono">
                Active chain:{" "}
                <span className="number-pill number-pill-sm number-pill-mono">
                  {chainId || "unknown"}
                </span>
              </span>
            </div>

            <ul className="space-y-3">
              {wallets.map((wallet) => (
                <li
                  key={`${wallet.address}-${wallet.walletClientType ?? "unknown"}`}
                  className="profile-wallet-row"
                >
                  <div>
                    <p className="text-sm font-medium text-slate-100">
                      {shortenAddress(wallet.address)}
                      {wallet.address.toLowerCase() ===
                        String(walletAddress || "").toLowerCase() &&
                        user?.kycStatus === "verified" && (
                        <span className="ml-2 inline-flex items-center gap-1 rounded-full border border-emerald-300/50 bg-emerald-500/15 px-2 py-0.5 text-[11px] font-semibold uppercase tracking-[0.08em] text-emerald-200">
                          ✓ KYC
                        </span>
                      )}
                    </p>
                    <p className="text-xs text-slate-400">
                      {wallet.walletClientType === "privy"
                        ? "Privy embedded wallet"
                        : "External wallet"}
                    </p>
                  </div>
                  <p className="text-xs text-slate-300">
                    <span className="number-pill number-pill-sm number-pill-mono">
                      {formatEth(walletBalances[wallet.address] || null)}
                    </span>
                  </p>
                </li>
              ))}
            </ul>

            <div className="grid gap-3 sm:grid-cols-2">
              {!embeddedWalletAddress && (
                <button
                  type="button"
                  onClick={() => {
                    createEmbeddedWallet().catch((err) => {
                      console.error("Create embedded wallet error:", err);
                    });
                  }}
                  className="btn btn-primary !w-full !py-2.5 !text-sm"
                >
                  Create embedded Privy wallet
                </button>
              )}

              <button
                type="button"
                onClick={connectExternalWallet}
                className="btn btn-secondary !w-full !py-2.5 !text-sm"
              >
                Connect another external wallet
              </button>

              <button
                type="button"
                onClick={() => {
                  refreshProfile().catch((err) => {
                    console.error("Refresh profile error:", err);
                  });
                }}
                className="btn btn-secondary sm:col-span-2 !w-full !py-2.5 !text-sm"
              >
                Refresh balances
              </button>
            </div>
          </div>

          <div className="cyber-card dashboard-panel text-panel space-y-4 p-5">
            <div className="flex items-center justify-between">
              <p className="dashboard-section-title">Property NFTs</p>
              <span className="text-xs text-slate-400 number-pair">
                <span className="number-pill number-pill-sm">{ownedNfts}</span>
                <span>total</span>
              </span>
            </div>

            {ownedHouses.length === 0 ? (
              <p className="text-sm text-slate-400">
                No properties found for this wallet yet.
              </p>
            ) : (
              <div className="grid gap-3 md:grid-cols-2">
                {ownedHouses.map((house) => (
                  <article
                    key={house.tokenId}
                    className="profile-property-card"
                  >
                    <HouseThumbnail
                      house={house}
                      className="mb-3 h-36 w-full"
                    />
                    <p className="text-xs uppercase tracking-[0.08em] text-slate-400">
                      Token{" "}
                      <span className="number-pill number-pill-sm number-pill-mono">
                        #{house.tokenId}
                      </span>
                    </p>
                    <p className="mt-1 text-sm font-medium text-slate-100">
                      {house.metadata.address}
                    </p>
                    <p className="mt-1 text-xs text-slate-400">
                      {house.metadata.city}, {house.metadata.state} •{" "}
                      {house.metadata.propertyType}
                    </p>
                  </article>
                ))}
              </div>
            )}
          </div>
        </section>
      </main>
    </div>
  );
};

const InfoCard: React.FC<{ label: string; value: string }> = ({
  label,
  value,
}) => {
  const hasNumericSignal = /\d/.test(value);
  const valueClass = hasNumericSignal
    ? "number-pill number-pill-lg number-pill-mono"
    : "text-lg font-semibold text-slate-100";

  return (
    <div className="profile-info-card text-center">
      <p className="dashboard-kpi-label text-xs uppercase">
        {label}
      </p>
      <p className="mt-2 text-lg font-semibold text-slate-100">
        <span className={valueClass}>{value}</span>
      </p>
    </div>
  );
};
