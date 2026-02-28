import React from "react";
import toast from "react-hot-toast";

import { useAuth } from "./AuthProvider";

const parseChainIdValue = (value: string | null): number | null => {
  if (!value) return null;
  const trimmed = value.trim();
  if (!trimmed) return null;
  const parts = trimmed.split(":");
  const parsed = Number.parseInt(parts[parts.length - 1], 10);
  return Number.isFinite(parsed) ? parsed : null;
};

export const ChainWarning: React.FC = () => {
  const { chainId, switchNetwork } = useAuth();
  const expectedRaw = import.meta.env.VITE_EXPECTED_CHAIN_ID;
  const expected = Number.parseInt(String(expectedRaw || ""), 10);

  if (!Number.isFinite(expected)) return null;

  const current = parseChainIdValue(chainId);
  if (current === null || current === expected) return null;

  const getNetworkName = (id: number): string => {
    if (id === 1) return "Ethereum Mainnet";
    if (id === 11155111) return "Sepolia";
    if (id === 137) return "Polygon";
    if (id === 31337) return "Anvil Local";
    if (id === 10) return "Optimism";
    if (id === 42161) return "Arbitrum";
    return `Chain ${id}`;
  };

  const handleSwitchNetwork = async () => {
    try {
      await switchNetwork(expected);
      toast.success(`Switched to ${getNetworkName(expected)}`);
      return;
    } catch (error: any) {
      toast.error(error?.message || "Network switch failed", {
        id: "network-switch-error",
      });
    }
  };

  return (
    <div className="mx-auto max-w-6xl px-4 pt-3 sm:px-6 lg:px-8">
      <div className="rounded-xl border border-amber-300/40 bg-amber-500/10 px-4 py-3">
        <div className="flex flex-col items-start gap-3 md:flex-row md:items-center md:justify-between">
          <p className="text-xs text-amber-100 md:text-sm">
            Wrong network:{" "}
            <span className="font-semibold text-slate-50">
              {getNetworkName(current)}
            </span>{" "}
            (
            <span className="number-pill number-pill-sm number-pill-mono">
              {current}
            </span>
            ) {"->"} expected{" "}
            <span className="font-semibold text-slate-50">
              {getNetworkName(expected)}
            </span>{" "}
            (
            <span className="number-pill number-pill-sm number-pill-mono">
              {expected}
            </span>
            ).
          </p>
          <button
            type="button"
            onClick={handleSwitchNetwork}
            className="btn btn-secondary self-start shrink-0 !w-auto !py-2 !px-3 !text-xs md:self-auto"
            style={{ width: "fit-content" }}
          >
            Switch Network
          </button>
        </div>
        <p className="mt-2 text-[11px] text-amber-200/90" aria-live="polite">
          Trading actions are paused until your wallet is on a supported chain.
          You can still browse listings and manage your account.
        </p>
      </div>
    </div>
  );
};
