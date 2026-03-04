import React, { useCallback, useEffect, useMemo, useState } from "react";
import { Link, NavLink, useNavigate } from "react-router-dom";
import toast from "react-hot-toast";

import { useAuth } from "./AuthProvider";
import { BrandMark } from "./BrandMark";
import { useUXMode } from "./UXModeProvider";
import { apiClient } from "@shared/utils/api";
import type { Notification } from "@shared/types";

const NAV_ITEMS: Array<{ readonly to: string; readonly label: string }> = [
  { to: "/dashboard", label: "Dashboard" },
  { to: "/profile", label: "Profile" },
  { to: "/mint", label: "List Property" },
  { to: "/marketplace", label: "Marketplace" },
];

export const Navigation: React.FC = () => {
  const {
    activeWalletClientType,
    chainId,
    connectExternalWallet,
    createEmbeddedWallet,
    embeddedWalletAddress,
    isAuthenticated,
    logout,
    selectActiveWallet,
    switchNetwork,
    walletAddress,
    wallets,
    user,
  } = useAuth();
  const { mode, setMode } = useUXMode();
  const navigate = useNavigate();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [notificationsOpen, setNotificationsOpen] = useState(false);
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [unreadNotifications, setUnreadNotifications] = useState(0);
  const parseChainIdValue = (value: string | null): number | null => {
    if (!value) return null;
    const trimmed = value.trim();
    if (!trimmed) return null;
    const parts = trimmed.split(":");
    const parsed = Number.parseInt(parts[parts.length - 1], 10);
    return Number.isFinite(parsed) ? parsed : null;
  };
  const expectedChainId = Number.parseInt(
    String(import.meta.env.VITE_EXPECTED_CHAIN_ID || ""),
    10,
  );
  const currentChainId = parseChainIdValue(chainId);
  const [selectedChainId, setSelectedChainId] = useState<number | null>(
    Number.isFinite(expectedChainId)
      ? expectedChainId
      : currentChainId || 11155111,
  );

  useEffect(() => {
    if (currentChainId !== null) {
      setSelectedChainId(currentChainId);
    }
  }, [currentChainId]);

  const loadNotifications = useCallback(async () => {
    if (!isAuthenticated || !apiClient.hasAuthToken()) {
      setNotifications([]);
      setUnreadNotifications(0);
      return;
    }
    try {
      const response = await apiClient.getNotifications();
      if (!response.success || !response.data) {
        return;
      }
      setNotifications(response.data.slice(0, 16));
      const unreadCountFromServer = Number(
        (response as { unreadCount?: number }).unreadCount,
      );
      if (Number.isFinite(unreadCountFromServer)) {
        setUnreadNotifications(unreadCountFromServer);
      } else {
        setUnreadNotifications(
          response.data.filter((entry) => !entry.read).length,
        );
      }
    } catch {
      setNotifications([]);
      setUnreadNotifications(0);
    }
  }, [isAuthenticated]);

  useEffect(() => {
    let cancelled = false;
    const safeLoad = async () => {
      if (cancelled) {
        return;
      }
      await loadNotifications();
    };
    void safeLoad();
    const intervalId = window.setInterval(() => {
      void safeLoad();
    }, 12000);
    return () => {
      cancelled = true;
      window.clearInterval(intervalId);
    };
  }, [loadNotifications]);

  const handleNotificationClick = useCallback(
    async (notification: Notification) => {
      if (!notification.read) {
        try {
          await apiClient.markNotificationRead(notification.id);
          setNotifications((previous) =>
            previous.map((entry) =>
              entry.id === notification.id ? { ...entry, read: true } : entry,
            ),
          );
          setUnreadNotifications((previous) =>
            previous > 0 ? previous - 1 : 0,
          );
        } catch {
          // Ignore read-marking failures in the UI.
        }
      }

      const tokenId = String(notification.data?.tokenId ?? "").trim();
      if (tokenId) {
        const conversationId = String(
          notification.data?.conversationId ?? "",
        ).trim();
        const senderWallet = String(notification.data?.from ?? "")
          .trim()
          .toLowerCase();
        const params = new URLSearchParams();
        if (notification.type === "message_received") {
          params.set("tab", "messages");
          if (conversationId) {
            params.set("conversation", conversationId);
          }
          if (/^0x[a-fA-F0-9]{40}$/.test(senderWallet)) {
            params.set("to", senderWallet);
          }
        }
        const query = params.toString();
        navigate(`/houses/${tokenId}${query ? `?${query}` : ""}`);
        setNotificationsOpen(false);
      }
    },
    [navigate],
  );

  const supportedChainIds = useMemo(() => {
    const fromEnv = String(import.meta.env.VITE_SUPPORTED_CHAIN_IDS || "")
      .split(",")
      .map((value) => Number.parseInt(value.trim(), 10))
      .filter((value) => Number.isFinite(value) && value > 0);
    const defaults = [31337, 11155111, 1, 137];
    const merged = new Set<number>([
      ...defaults,
      ...fromEnv,
      ...(Number.isFinite(expectedChainId) ? [expectedChainId] : []),
      ...(currentChainId ? [currentChainId] : []),
    ]);
    return Array.from(merged).sort((a, b) => a - b);
  }, [currentChainId, expectedChainId]);

  const getNetworkName = (id: number): string => {
    if (id === 1) return "Ethereum";
    if (id === 11155111) return "Sepolia";
    if (id === 137) return "Polygon";
    if (id === 31337) return "Anvil";
    return `Chain ${id}`;
  };
  const shortenWalletAddress = (address: string): string => {
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
  };
  const describeWalletType = (
    walletClientType: string | null | undefined,
  ): string => {
    return walletClientType === "privy" ? "Privy" : "External";
  };

  if (!isAuthenticated) {
    return null;
  }

  const handleLogout = async () => {
    await logout();
    navigate("/");
  };

  const handleCreateEmbeddedWallet = () => {
    createEmbeddedWallet().catch((err) => {
      console.error("Create embedded wallet failed:", err);
    });
  };

  const handleSwitchNetwork = async () => {
    if (!selectedChainId) {
      toast.error("Choose a network first.", { id: "network-switch-error" });
      return;
    }
    try {
      await switchNetwork(selectedChainId);
      toast.success(`Switched to ${getNetworkName(selectedChainId)}`);
    } catch (err: any) {
      toast.error(err?.message || "Network switch failed", {
        id: "network-switch-error",
      });
    }
  };
  const handleActiveWalletChange = (
    event: React.ChangeEvent<HTMLSelectElement>,
  ) => {
    const nextAddress = event.target.value;
    if (!nextAddress) {
      return;
    }
    try {
      selectActiveWallet(nextAddress);
      toast.success(
        `Active wallet set to ${shortenWalletAddress(nextAddress)}`,
      );
    } catch (err: any) {
      toast.error(err?.message || "Failed to switch active wallet");
    }
  };

  const networkPillTone =
    currentChainId !== null &&
    Number.isFinite(expectedChainId) &&
    currentChainId !== expectedChainId
      ? "text-amber-200"
      : "text-emerald-200";
  return (
    <nav className="cyber-nav sticky top-0 z-50 border-b border-slate-700/60">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="site-header-banner">
          <Link
            to="/dashboard"
            className="site-brand-link site-brand-link--header"
            aria-label="Go to dashboard"
          >
            <BrandMark
              size="sm"
              showWordmark={false}
              logoAsset="home"
              className="site-brand-logo site-brand-logo--header"
            />
          </Link>

          <div className="hidden lg:flex items-center gap-2 xl:gap-3">
            {NAV_ITEMS.map((item) => (
              <DesktopNavLink key={item.to} to={item.to}>
                {item.label}
              </DesktopNavLink>
            ))}
          </div>

          <div className="site-header-actions hidden md:flex items-center gap-2 lg:gap-3">
            <ModeToggle mode={mode} setMode={setMode} />
            <button
              type="button"
              onClick={() => navigate("/profile")}
              className="btn btn-secondary header-action-btn !px-3 !py-2 !text-xs"
            >
              Wallet overview
            </button>
            {!embeddedWalletAddress && (
              <button
                type="button"
                onClick={handleCreateEmbeddedWallet}
                className="btn btn-secondary header-action-btn !px-3 !py-2 !text-xs"
              >
                Create Privy wallet
              </button>
            )}
            <button
              type="button"
              onClick={connectExternalWallet}
              className="btn btn-primary header-action-btn !px-3 !py-2 !text-xs"
            >
              Connect wallet
            </button>
            <div className="relative">
              <button
                type="button"
                onClick={() => setNotificationsOpen((value) => !value)}
                className="btn btn-secondary header-action-btn !px-3 !py-2 !text-xs"
                aria-label="Open notifications"
              >
                <span aria-hidden="true">🔔</span>
                Notifications
                {unreadNotifications > 0 && (
                  <span className="number-pill number-pill-xs number-pill-mono">
                    {unreadNotifications}
                  </span>
                )}
              </button>
              {notificationsOpen && (
                <div className="absolute right-0 z-50 mt-2 w-[320px] rounded-xl border border-slate-700/70 bg-slate-950/95 p-2 shadow-2xl">
                  <p className="px-3 py-2 text-xs font-semibold uppercase tracking-[0.08em] text-slate-400">
                    Notifications
                  </p>
                  <div className="max-h-80 overflow-y-auto">
                    {notifications.length === 0 ? (
                      <p className="px-3 py-3 text-xs text-slate-400">
                        No new updates.
                      </p>
                    ) : (
                      notifications.map((notification) => (
                        <button
                          key={notification.id}
                          type="button"
                          onClick={() => handleNotificationClick(notification)}
                          className={`rounded-lg px-3 py-2 ${
                            notification.read
                              ? "bg-transparent"
                              : "bg-blue-500/10"
                          } w-full text-left transition hover:bg-slate-800/55`}
                        >
                          <p className="text-xs font-semibold text-slate-100">
                            {notification.title}
                          </p>
                          <p className="mt-0.5 text-xs text-slate-300">
                            {notification.message}
                          </p>
                          <p className="mt-1 text-[11px] text-slate-400">
                            {new Date(notification.createdAt).toLocaleString()}
                          </p>
                        </button>
                      ))
                    )}
                  </div>
                </div>
              )}
            </div>
            <div className="cyber-card px-3 py-1.5">
              <label className="text-[11px] font-medium uppercase tracking-[0.08em] text-slate-400">
                Network
              </label>
              <div className="mt-1 flex items-center gap-2">
                <select
                  value={selectedChainId ?? ""}
                  onChange={(event) =>
                    setSelectedChainId(Number.parseInt(event.target.value, 10))
                  }
                  className="cyber-input !h-[44px] !min-h-[44px] !w-[130px] !px-2 !py-1 !text-xs"
                >
                  {supportedChainIds.map((networkId) => (
                    <option key={networkId} value={networkId}>
                      {getNetworkName(networkId)}
                    </option>
                  ))}
                </select>
                <button
                  type="button"
                  onClick={handleSwitchNetwork}
                  className="btn btn-secondary header-action-btn !px-2.5 !py-1.5 !text-[11px]"
                >
                  Switch
                </button>
              </div>
            </div>
            {walletAddress && (
              <div className="cyber-card px-3 py-1.5">
                <p className="text-[11px] font-medium uppercase tracking-[0.08em] text-slate-400">
                  Wallet
                </p>
                <p className="text-xs font-medium text-blue-100">
                  {walletAddress.slice(0, 6)}...{walletAddress.slice(-4)}
                </p>
                {user?.kycStatus === "verified" && (
                  <p className="mt-1">
                    <span className="inline-flex items-center gap-1 rounded-full border border-emerald-300/50 bg-emerald-500/15 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-[0.08em] text-emerald-200">
                      ✓ KYC
                    </span>
                  </p>
                )}
                <p className="mt-0.5 text-[11px] text-slate-300">
                  {describeWalletType(activeWalletClientType)}
                </p>
                <p className={`mt-0.5 text-[11px] ${networkPillTone}`}>
                  {currentChainId !== null ? (
                    <>
                      {getNetworkName(currentChainId)}{" "}
                      <span className="number-pill number-pill-xs number-pill-mono">
                        {currentChainId}
                      </span>
                    </>
                  ) : (
                    "Unknown network"
                  )}
                </p>
              </div>
            )}
            {wallets.length > 1 && (
              <div className="cyber-card px-3 py-1.5">
                <label className="text-[11px] font-medium uppercase tracking-[0.08em] text-slate-400">
                  Active wallet
                </label>
                <div className="mt-1">
                  <select
                    value={walletAddress ?? ""}
                    onChange={handleActiveWalletChange}
                    className="cyber-input !h-[44px] !min-h-[44px] !w-[190px] !px-2 !py-1 !text-xs"
                  >
                    {wallets.map((connectedWallet) => (
                      <option
                        key={connectedWallet.address}
                        value={connectedWallet.address}
                      >
                        {`${describeWalletType(connectedWallet.walletClientType)} · ${shortenWalletAddress(connectedWallet.address)}`}
                      </option>
                    ))}
                  </select>
                </div>
              </div>
            )}
            <button
              onClick={handleLogout}
              className="btn btn-secondary header-action-btn !px-3 !py-2 !text-xs"
            >
              Disconnect
            </button>
          </div>

          <div className="md:hidden flex items-center gap-2">
            <button
              type="button"
              onClick={() => {
                setMobileMenuOpen(false);
                setNotificationsOpen((value) => !value);
              }}
              className="btn btn-secondary header-action-btn header-icon-btn !px-2.5 !py-2 !text-xs"
              aria-label="Open notifications"
            >
              <span aria-hidden="true">🔔</span>
              {unreadNotifications > 0 && (
                <span className="number-pill number-pill-xs number-pill-mono">
                  {unreadNotifications}
                </span>
              )}
            </button>
            <button
              type="button"
              onClick={() => {
                setNotificationsOpen(false);
                setMobileMenuOpen((value) => !value);
              }}
              className="btn btn-secondary header-action-btn !px-3 !py-2 !text-xs"
              aria-label={mobileMenuOpen ? "Close menu" : "Open menu"}
            >
              <span>{mobileMenuOpen ? "Close" : "Menu"}</span>
              <svg
                className="h-4 w-4"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                {mobileMenuOpen ? (
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                ) : (
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M4 7h16M4 12h16M4 17h16"
                  />
                )}
              </svg>
            </button>
          </div>
        </div>
      </div>

      {notificationsOpen && (
        <div className="md:hidden border-t border-slate-700/60 bg-slate-950/95 px-4 py-3">
          <div className="mx-auto max-w-6xl rounded-lg border border-slate-700/60 bg-slate-950/70 p-3">
            <p className="text-xs font-semibold uppercase tracking-[0.08em] text-slate-400">
              Notifications
            </p>
            <div className="mt-2 max-h-56 space-y-2 overflow-y-auto">
              {notifications.length === 0 ? (
                <p className="text-xs text-slate-400">No new updates.</p>
              ) : (
                notifications.map((notification) => (
                  <div
                    key={notification.id}
                    className={`rounded-md px-2.5 py-2 ${
                      notification.read ? "bg-slate-900/40" : "bg-blue-500/10"
                    }`}
                  >
                    <p className="text-xs font-semibold text-slate-100">
                      {notification.title}
                    </p>
                    <p className="mt-0.5 text-xs text-slate-300">
                      {notification.message}
                    </p>
                  </div>
                ))
              )}
            </div>
          </div>
        </div>
      )}

      {mobileMenuOpen && (
        <div className="border-t border-slate-700/60 bg-slate-950/95 px-4 py-4 md:hidden">
          <div className="mx-auto max-w-6xl">
            <div className="mb-4 flex justify-center">
              <ModeToggle mode={mode} setMode={setMode} />
            </div>
            <div className="space-y-1">
              {NAV_ITEMS.map((item) => (
                <MobileNavLink
                  key={item.to}
                  to={item.to}
                  onClick={() => setMobileMenuOpen(false)}
                >
                  {item.label}
                </MobileNavLink>
              ))}
            </div>
            <div className="mt-4 flex flex-wrap items-center justify-center gap-2">
              <button
                type="button"
                onClick={() => {
                  navigate("/profile");
                  setMobileMenuOpen(false);
                }}
                className="btn btn-secondary header-action-btn !px-3 !py-2 !text-xs"
              >
                Wallet overview
              </button>
              {!embeddedWalletAddress && (
                <button
                  type="button"
                  onClick={handleCreateEmbeddedWallet}
                  className="btn btn-secondary header-action-btn !px-3 !py-2 !text-xs"
                >
                  Create Privy wallet
                </button>
              )}
              <button
                type="button"
                onClick={connectExternalWallet}
                className="btn btn-primary header-action-btn !px-3 !py-2 !text-xs"
              >
                Connect wallet
              </button>
            </div>
            <div className="mt-3 rounded-lg border border-slate-700/60 bg-slate-950/70 p-3">
              <label className="text-xs font-medium uppercase tracking-[0.08em] text-slate-400">
                Network
              </label>
              <div className="mt-2 flex gap-2">
                <select
                  value={selectedChainId ?? ""}
                  onChange={(event) =>
                    setSelectedChainId(Number.parseInt(event.target.value, 10))
                  }
                  className="cyber-input !h-[44px] !min-h-[44px] !flex-1 !px-2 !py-1 !text-sm"
                >
                  {supportedChainIds.map((networkId) => (
                    <option key={networkId} value={networkId}>
                      {getNetworkName(networkId)}
                    </option>
                  ))}
                </select>
                <button
                  type="button"
                  onClick={handleSwitchNetwork}
                  className="btn btn-secondary header-action-btn !px-3 !py-2 !text-xs"
                >
                  Switch
                </button>
              </div>
              <p className={`mt-2 text-xs ${networkPillTone}`}>
                Current:{" "}
                {currentChainId !== null ? (
                  <>
                    {getNetworkName(currentChainId)}{" "}
                    <span className="number-pill number-pill-xs number-pill-mono">
                      {currentChainId}
                    </span>
                  </>
                ) : (
                  "Unknown"
                )}
              </p>
            </div>
            {wallets.length > 1 && (
              <div className="mt-3 rounded-lg border border-slate-700/60 bg-slate-950/70 p-3">
                <label className="text-xs font-medium uppercase tracking-[0.08em] text-slate-400">
                  Active wallet
                </label>
                <select
                  value={walletAddress ?? ""}
                  onChange={handleActiveWalletChange}
                  className="cyber-input mt-2 !h-[44px] !min-h-[44px] !w-full !px-2 !py-1 !text-sm"
                >
                  {wallets.map((connectedWallet) => (
                    <option
                      key={connectedWallet.address}
                      value={connectedWallet.address}
                    >
                      {`${describeWalletType(connectedWallet.walletClientType)} · ${shortenWalletAddress(connectedWallet.address)}`}
                    </option>
                  ))}
                </select>
              </div>
            )}
            <button
              type="button"
              onClick={handleLogout}
              className="mt-3 btn btn-secondary header-action-btn !px-3 !py-2 !text-xs"
            >
              Disconnect
            </button>
          </div>
        </div>
      )}
    </nav>
  );
};

const DesktopNavLink: React.FC<{
  readonly to: string;
  readonly children: React.ReactNode;
}> = ({ to, children }) => (
  <NavLink
    to={to}
    className={({ isActive }) =>
      `nav-link-btn ${isActive ? "nav-link-btn-active" : ""}`
    }
  >
    {children}
  </NavLink>
);

const MobileNavLink: React.FC<{
  readonly to: string;
  readonly onClick: () => void;
  readonly children: React.ReactNode;
}> = ({ to, onClick, children }) => (
  <NavLink
    to={to}
    onClick={onClick}
    className={({ isActive }) =>
      `nav-link-btn nav-link-btn-mobile justify-center ${
        isActive ? "nav-link-btn-active" : ""
      }`
    }
  >
    {children}
  </NavLink>
);

const ModeToggle: React.FC<{
  readonly mode: "easy" | "degen";
  readonly setMode: (mode: "easy" | "degen") => void;
}> = ({ mode, setMode }) => {
  const isEasy = mode === "easy";

  return (
    <div className="mode-toggle">
      <button
        type="button"
        onClick={() => setMode("easy")}
        className={`mode-btn ${isEasy ? "active-easy" : ""}`}
      >
        Guided
      </button>
      <button
        type="button"
        onClick={() => setMode("degen")}
        className={`mode-btn ${!isEasy ? "active-degen" : ""}`}
      >
        Pro
      </button>
    </div>
  );
};
