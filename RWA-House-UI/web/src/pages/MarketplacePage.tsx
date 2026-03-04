import React, { useEffect, useMemo, useRef, useState } from "react";
import { Link } from "react-router-dom";
import toast from "react-hot-toast";
import { ethers } from "ethers";

import { apiClient } from "@shared/utils/api";
import type { House } from "@shared/types";
import {
  fileToBase64,
  validateDocumentContent,
  validateFileUpload,
} from "@shared/utils/security";
import { useAuth } from "../components/AuthProvider";
import { HouseThumbnail } from "../components/HouseThumbnail";
import { useUXMode } from "../components/UXModeProvider";
import { saveLatestClaimKeyHash } from "../utils/claimKeyStorage";

type MarketFilter = "all" | "for_sale" | "for_rent";

interface UploadedMarketImage {
  id: string;
  name: string;
  size: number;
  dataUrl: string;
}

const MAX_MARKET_IMAGE_UPLOAD_SIZE_BYTES = 2 * 1024 * 1024;
const MAX_MARKET_IMAGE_COUNT = 10;

const formatListingPrice = (priceWei?: string): string => {
  if (!priceWei) {
    return "—";
  }

  try {
    const value = Number.parseFloat(ethers.formatEther(priceWei));
    return `${value.toLocaleString(undefined, { maximumFractionDigits: 4 })} ETH`;
  } catch {
    return priceWei;
  }
};

export const MarketplacePage: React.FC = () => {
  const { walletAddress, chainId, getEthereumProvider } = useAuth();
  const { mode } = useUXMode();
  const [houses, setHouses] = useState<House[]>([]);
  const [filter, setFilter] = useState<MarketFilter>("all");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [actionLoadingToken, setActionLoadingToken] = useState<string | null>(
    null,
  );
  const [imageEditorHouse, setImageEditorHouse] = useState<House | null>(null);
  const [imageEditorLinks, setImageEditorLinks] = useState("");
  const [imageEditorUploads, setImageEditorUploads] = useState<
    UploadedMarketImage[]
  >([]);
  const [isImageEditorUploading, setIsImageEditorUploading] = useState(false);
  const [isImageEditorSaving, setIsImageEditorSaving] = useState(false);
  const imageEditorInputRef = useRef<HTMLInputElement>(null);

  const expectedChainId = Number.parseInt(
    String(import.meta.env.VITE_EXPECTED_CHAIN_ID || ""),
    10,
  );
  const parseChainId = (value: string | null): number | null => {
    if (!value) return null;
    const trimmed = value.trim();
    if (!trimmed) return null;
    const parts = trimmed.split(":");
    const parsed = Number.parseInt(parts[parts.length - 1], 10);
    return Number.isFinite(parsed) ? parsed : null;
  };
  const connectedChainId = parseChainId(chainId);
  const wrongChain =
    Number.isFinite(expectedChainId) &&
    connectedChainId !== null &&
    connectedChainId !== expectedChainId;

  const notifyWrongNetwork = () => {
    if (!Number.isFinite(expectedChainId)) {
      return;
    }
    toast.error(
      `Wrong network. Switch your wallet to chain ${expectedChainId} first.`,
    );
  };

  const loadMarketplace = async () => {
    try {
      setLoading(true);
      setError(null);
      const resp = await apiClient.getHouses();
      if (resp.success && resp.data) {
        setHouses(resp.data);
      } else {
        setError(resp.message || "Failed to load marketplace");
      }
    } catch (e: any) {
      setError(e?.message || "Failed to load marketplace");
    } finally {
      setLoading(false);
    }
  };

  const parseImageLinks = (rawValue: string): string[] => {
    return Array.from(
      new Set(
        rawValue
          .split("\n")
          .map((entry) => entry.trim())
          .filter((entry) => entry.length > 0),
      ),
    ).slice(0, MAX_MARKET_IMAGE_COUNT);
  };

  const mergeImageEditorImages = (
    rawLinks: string,
    uploads: UploadedMarketImage[],
  ): string[] => {
    const linkedImages = parseImageLinks(rawLinks);
    const uploadedImages = uploads.map((image) => image.dataUrl);
    return Array.from(
      new Set([...uploadedImages, ...linkedImages]),
    ).slice(0, MAX_MARKET_IMAGE_COUNT);
  };

  const openImageEditor = (house: House) => {
    const existingImages = Array.isArray(house.metadata?.images)
      ? house.metadata.images
      : [];
    const existingUploads: UploadedMarketImage[] = [];
    const existingLinks: string[] = [];

    existingImages.forEach((image, index) => {
      const value = String(image || "").trim();
      if (!value) {
        return;
      }
      if (value.toLowerCase().startsWith("data:image/png;base64,")) {
        existingUploads.push({
          id: `existing:${house.tokenId}:${index}`,
          name: `Uploaded image ${index + 1}`,
          size: 0,
          dataUrl: value,
        });
        return;
      }
      existingLinks.push(value);
    });

    setImageEditorHouse(house);
    setImageEditorUploads(existingUploads.slice(0, MAX_MARKET_IMAGE_COUNT));
    setImageEditorLinks(
      existingLinks.slice(0, MAX_MARKET_IMAGE_COUNT).join("\n"),
    );
  };

  const closeImageEditor = (force = false) => {
    if (isImageEditorSaving && !force) {
      return;
    }
    setImageEditorHouse(null);
    setImageEditorLinks("");
    setImageEditorUploads([]);
    setIsImageEditorUploading(false);
  };

  const removeImageEditorUpload = (imageId: string) => {
    setImageEditorUploads((previous) =>
      previous.filter((image) => image.id !== imageId),
    );
  };

  const handleImageEditorUpload = async (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    const files = Array.from(event.target.files || []);
    if (files.length === 0) {
      return;
    }

    setIsImageEditorUploading(true);
    try {
      const nextUploads = [...imageEditorUploads];
      const duplicateGuard = new Set(nextUploads.map((image) => image.id));
      let rejectedCount = 0;

      for (const file of files) {
        if (nextUploads.length >= MAX_MARKET_IMAGE_COUNT) {
          rejectedCount += 1;
          continue;
        }

        const validation = validateFileUpload(
          file,
          ["image/png"],
          MAX_MARKET_IMAGE_UPLOAD_SIZE_BYTES,
        );
        if (!validation.valid) {
          rejectedCount += 1;
          continue;
        }

        const imageId = `${file.name}:${file.size}:${file.lastModified}`;
        if (duplicateGuard.has(imageId)) {
          rejectedCount += 1;
          continue;
        }

        const contentSafe = await validateDocumentContent(file);
        if (!contentSafe) {
          rejectedCount += 1;
          continue;
        }

        const base64Payload = await fileToBase64(file);
        nextUploads.push({
          id: imageId,
          name: file.name,
          size: file.size,
          dataUrl: `data:image/png;base64,${base64Payload}`,
        });
        duplicateGuard.add(imageId);
      }

      setImageEditorUploads(nextUploads.slice(0, MAX_MARKET_IMAGE_COUNT));
      if (rejectedCount > 0) {
        toast.error(
          `${rejectedCount} image${
            rejectedCount === 1 ? "" : "s"
          } rejected. PNG only, max 2 MB each.`,
        );
      }
    } catch (error: any) {
      toast.error(error?.message || "Unable to process uploaded PNG images");
    } finally {
      setIsImageEditorUploading(false);
      event.target.value = "";
    }
  };

  const handleSaveImageUpdates = async () => {
    if (!imageEditorHouse) {
      return;
    }

    const mergedImages = mergeImageEditorImages(
      imageEditorLinks,
      imageEditorUploads,
    );
    if (mergedImages.length === 0) {
      toast.error("Add at least one image link or PNG upload");
      return;
    }

    try {
      setIsImageEditorSaving(true);
      const response = await apiClient.updateHouseImages(
        imageEditorHouse.tokenId,
        mergedImages,
      );
      if (!response.success) {
        toast.error(response.message || "Failed to update property images");
        return;
      }
      toast.success("Property images updated");
      closeImageEditor(true);
      await loadMarketplace();
    } catch (error: any) {
      toast.error(error?.message || "Failed to update property images");
    } finally {
      setIsImageEditorSaving(false);
    }
  };

  useEffect(() => {
    let mounted = true;

    (async () => {
      await loadMarketplace();
      if (!mounted) return;
    })();

    return () => {
      mounted = false;
    };
  }, []);

  const handleQuickBuy = async (house: House) => {
    const listing = house.listing;
    if (!listing || listing.listingType !== "for_sale") return;
    if (!walletAddress) {
      toast.error("Connect your wallet to buy");
      return;
    }
    if (wrongChain) {
      notifyWrongNetwork();
      return;
    }
    if (
      listing.isPrivateSale &&
      listing.allowedBuyer &&
      listing.allowedBuyer.toLowerCase() !== walletAddress.toLowerCase()
    ) {
      toast.error("Private listing: your wallet is not the allowed buyer");
      return;
    }

    try {
      setActionLoadingToken(house.tokenId);
      const resp = await apiClient.sellHouse({
        action: "sell",
        sellerAddress: house.ownerAddress,
        buyerAddress: walletAddress,
        tokenId: house.tokenId,
        price: listing.price,
        buyerPublicKey: "",
        isPrivateSale: listing.isPrivateSale,
      });

      if (resp.success) {
        const keyHash = String(resp.data?.keyHash || "").trim();
        const keyHashSaved = saveLatestClaimKeyHash(keyHash);
        toast.success(
          resp.txHash
            ? `Buy submitted: ${resp.txHash.slice(0, 10)}...${
              keyHashSaved ? " Key hash saved for claim." : ""
            }`
            : `Buy submitted via CRE${
              keyHashSaved ? " (key hash saved for claim)." : ""
            }`,
        );
        await loadMarketplace();
      } else {
        toast.error(resp.message || "Failed to submit buy action");
      }
    } catch (err: any) {
      toast.error(err?.message || "Failed to submit buy action");
    } finally {
      setActionLoadingToken(null);
    }
  };

  const fundRentalDepositIfNeeded = async (
    tokenId: string,
    renter: string,
    requiredDepositWei: string,
  ) => {
    const contractAddress = import.meta.env.VITE_HOUSE_RWA_ADDRESS;
    if (!contractAddress) {
      throw new Error("Missing VITE_HOUSE_RWA_ADDRESS in web env config");
    }

    const ethereumProvider = await getEthereumProvider();
    const provider = new ethers.BrowserProvider(ethereumProvider);
    const signer = await provider.getSigner();
    const contract = new ethers.Contract(
      contractAddress,
      [
        "function depositForRental(uint256 tokenId) payable",
        "function pendingRentalDeposits(uint256 tokenId,address renter) view returns (uint256)",
      ],
      signer,
    );

    const required = BigInt(requiredDepositWei);
    const pendingRaw = await contract.pendingRentalDeposits(
      BigInt(tokenId),
      renter,
    );
    const pending = BigInt(pendingRaw);

    if (pending >= required) return;

    const shortfall = required - pending;
    const tx = await contract.depositForRental(BigInt(tokenId), {
      value: shortfall,
    });
    await tx.wait();
  };

  const handleQuickRent = async (house: House) => {
    const listing = house.listing;
    if (!listing || listing.listingType !== "for_rent") return;
    if (!walletAddress) {
      toast.error("Connect your wallet to rent");
      return;
    }
    if (wrongChain) {
      notifyWrongNetwork();
      return;
    }

    try {
      setActionLoadingToken(house.tokenId);
      toast.loading("Verifying renter KYC...", { id: `rent-${house.tokenId}` });
      const kycResponse = await apiClient.ensureKYC(walletAddress);
      if (!kycResponse.success) {
        toast.error(kycResponse.message || "Unable to verify renter KYC", {
          id: `rent-${house.tokenId}`,
        });
        return;
      }

      toast.loading(
        mode === "degen"
          ? "Funding rental deposit..."
          : "Funding rental deposit...",
        { id: `rent-${house.tokenId}` },
      );

      await fundRentalDepositIfNeeded(
        house.tokenId,
        walletAddress,
        listing.price,
      );

      toast.loading(
        mode === "degen"
          ? "Submitting CRE rent action..."
          : "Starting your rental...",
        { id: `rent-${house.tokenId}` },
      );
      const resp = await apiClient.rentHouse({
        action: "rent",
        tokenId: house.tokenId,
        renterAddress: walletAddress,
        durationDays: 30,
        monthlyRent: listing.price,
        depositAmount: listing.price,
        renterPublicKey: "",
      });

      if (resp.success) {
        const accessKeyHash = String(resp.data?.accessKeyHash || "").trim();
        const keyHashSaved = saveLatestClaimKeyHash(accessKeyHash);
        toast.success(
          resp.txHash
            ? `Rent submitted: ${resp.txHash.slice(0, 10)}...${
              keyHashSaved ? " Key hash saved for claim." : ""
            }`
            : `Rent submitted via CRE${
              keyHashSaved ? " (key hash saved for claim)." : ""
            }`,
          { id: `rent-${house.tokenId}` },
        );
        await loadMarketplace();
      } else {
        toast.error(resp.message || "Failed to submit rent action", {
          id: `rent-${house.tokenId}`,
        });
      }
    } catch (err: any) {
      toast.error(err?.message || "Failed to submit rent action", {
        id: `rent-${house.tokenId}`,
      });
    } finally {
      setActionLoadingToken(null);
    }
  };

  const filtered = useMemo(() => {
    if (filter === "all") return houses;
    return houses.filter((h) => h.listing?.listingType === filter);
  }, [houses, filter]);

  return (
    <div className="min-h-screen bg-[#060b14] pb-12">
      <main className="page-shell page-shell-tight workspace-surface">
        <header className="page-header text-center">
          <div className="cyber-card dashboard-panel market-hero-card text-panel mx-auto max-w-5xl p-5 md:p-6">
            <div className="space-y-4">
              <div>
                <h1 className="dashboard-title">
                  {mode === "degen" ? "Marketplace" : "Homes Marketplace"}
                </h1>
                <p className="dashboard-action-body text-sm mt-1">
                  {mode === "degen"
                    ? "Browse listed RWAs. Private sales and encrypted document delivery powered by CRE."
                    : "Browse homes listed for private sale or rent. The platform handles secure transfers and private key delivery."}
                </p>
              </div>

              <div
                className="market-filter-group flex flex-wrap items-center justify-center gap-2"
                role="tablist"
                aria-label="Marketplace filters"
              >
                <button
                  type="button"
                  aria-pressed={filter === "all"}
                  className={`market-filter-btn ${filter === "all" ? "market-filter-btn-active" : ""}`}
                  onClick={() => setFilter("all")}
                >
                  {mode === "degen" ? "All" : "All Homes"}
                </button>
                <button
                  type="button"
                  aria-pressed={filter === "for_sale"}
                  className={`market-filter-btn ${filter === "for_sale" ? "market-filter-btn-active" : ""}`}
                  onClick={() => setFilter("for_sale")}
                >
                  {mode === "degen" ? "For Sale" : "Buy"}
                </button>
                <button
                  type="button"
                  aria-pressed={filter === "for_rent"}
                  className={`market-filter-btn ${filter === "for_rent" ? "market-filter-btn-active" : ""}`}
                  onClick={() => setFilter("for_rent")}
                >
                  {mode === "degen" ? "For Rent" : "Rent"}
                </button>
              </div>
            </div>
          </div>
        </header>

        <div className="section-shell space-y-6">
          {wrongChain && (
            <div
              role="status"
              aria-live="polite"
              className="text-panel border border-amber-300/45 bg-amber-500/10 p-4"
            >
              <p className="text-sm text-amber-100">
                Wrong network detected. Switch to chain{" "}
                <span className="number-pill number-pill-sm">
                  {expectedChainId}
                </span>{" "}
                to run buy or rent actions.
              </p>
            </div>
          )}

          {mode === "easy" && (
            <div className="text-panel border-cyan-500/40 bg-[rgba(8,145,178,0.12)] p-4">
              <p className="text-sm text-[var(--text-primary)]">
                Renting flow: deposit is funded first, then the rental is
                created, then bills show up in the property Payments tab.
              </p>
            </div>
          )}

          {error && (
            <div className="text-panel border border-rose-400/45 bg-rose-500/10 p-4">
              <p className="text-rose-100 font-mono text-sm">{error}</p>
              <p className="text-[var(--text-secondary)] text-sm mt-2">
                Tip: set `VITE_API_URL` to your backend and `VITE_RPC_URL` to a
                CORS-safe RPC (recommended: your backend `/rpc`).
              </p>
            </div>
          )}

          {loading ? (
            <div
              role="status"
              aria-live="polite"
              className="flex items-center justify-center py-24"
            >
              <div className="relative">
                <div className="w-16 h-16 border-2 border-[#00f3ff] border-t-transparent rounded-full animate-spin"></div>
                <div
                  className="absolute inset-0 w-16 h-16 border-2 border-[#b026ff] border-b-transparent rounded-full animate-spin"
                  style={{
                    animationDirection: "reverse",
                    animationDuration: "1.4s",
                  }}
                ></div>
              </div>
            </div>
          ) : filtered.length === 0 ? (
            <div className="cyber-card dashboard-panel text-panel p-10 text-center">
              <p className="text-[var(--text-secondary)]">
                No listings match the current filter.
              </p>
              <div className="mt-6">
                <Link
                  to="/mint"
                  className="btn btn-primary !px-5 !py-2.5 !text-sm"
                >
                  {mode === "degen" ? "Mint A Property" : "Add A Property"}
                </Link>
              </div>
            </div>
          ) : (
            <div className="cyber-card dashboard-panel text-panel space-y-6 p-6">
              <div className="flex items-center justify-between border-b border-slate-700/60 pb-4">
                <div>
                  <h2 className="dashboard-section-title">
                    {mode === "degen"
                      ? "Marketplace listings"
                      : "Property listings"}
                  </h2>
                  <p className="dashboard-section-note mt-1 text-sm">
                    {mode === "degen"
                      ? "Each card shows onchain listing state, pricing, and direct CRE actions."
                      : "Browse available homes and launch buy or rental flows directly from each card."}
                  </p>
                </div>
                <span className="dashboard-meta text-xs font-mono">
                  <span className="number-pair">
                    <span className="number-pill number-pill-sm">
                      {filtered.length}
                    </span>
                    <span>
                      {filtered.length === 1 ? "listing" : "listings"}
                    </span>
                  </span>
                </span>
              </div>

              <div className="market-listing-grid grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                {filtered.map((house) => {
                  const isOwner =
                    walletAddress &&
                    house.ownerAddress?.toLowerCase() ===
                      walletAddress.toLowerCase();
                  const listingType = house.listing?.listingType || "none";
                  const canMessageOwner =
                    Boolean(walletAddress) &&
                    !isOwner &&
                    listingType !== "none" &&
                    /^0x[a-fA-F0-9]{40}$/.test(String(house.ownerAddress || ""));
                  const houseOpenPath =
                    canMessageOwner
                      ? `/houses/${house.tokenId}?to=${encodeURIComponent(
                        String(house.ownerAddress || ""),
                      )}`
                      : `/houses/${house.tokenId}`;
                  const messageOwnerPath = `/houses/${house.tokenId}?tab=messages&to=${encodeURIComponent(
                    String(house.ownerAddress || ""),
                  )}`;
                  const badge =
                    listingType === "for_sale"
                      ? { text: "FOR SALE", color: "#00ff88" }
                      : listingType === "for_rent"
                        ? { text: "FOR RENT", color: "#00f3ff" }
                        : { text: "UNLISTED", color: "#8899aa" };

                  return (
                    <article
                      key={house.tokenId}
                      className="cyber-card dashboard-panel listing-card market-listing-card flex h-full flex-col p-6"
                    >
                      <HouseThumbnail
                        house={house}
                        className="mb-4 h-44 w-full"
                      />
                      <div className="flex items-start justify-between gap-4">
                        <div>
                          <h3 className="market-listing-title text-lg leading-tight text-[var(--text-primary)] font-semibold tracking-[0.01em]">
                            {house.metadata?.address ? (
                              house.metadata.address
                            ) : (
                              <>
                                Token{" "}
                                <span className="number-pill number-pill-sm number-pill-mono">
                                  #{house.tokenId}
                                </span>
                              </>
                            )}
                          </h3>
                          <p className="market-listing-subtitle text-sm text-[var(--text-secondary)] mt-2">
                            {house.metadata?.city
                              ? `${house.metadata.city}, ${house.metadata.state}`
                              : "Metadata pending"}
                          </p>
                        </div>

                        <div
                          className="rounded-full border px-3 py-1 text-xs font-mono"
                          style={{
                            borderColor: `${badge.color}55`,
                            color: badge.color,
                            background: `${badge.color}0f`,
                          }}
                        >
                          {badge.text}
                        </div>
                      </div>

                      <div className="market-metrics mt-4 grid grid-cols-3 gap-3 text-xs text-[var(--text-secondary)] font-mono">
                        <div>
                          <div className="text-[var(--text-secondary)]/85">
                            Beds
                          </div>
                          <div className="text-[var(--text-primary)]">
                            <span className="number-pill number-pill-sm">
                              {house.metadata?.bedrooms ?? "-"}
                            </span>
                          </div>
                        </div>
                        <div>
                          <div className="text-[var(--text-secondary)]/85">
                            Baths
                          </div>
                          <div className="text-[var(--text-primary)]">
                            <span className="number-pill number-pill-sm">
                              {house.metadata?.bathrooms ?? "-"}
                            </span>
                          </div>
                        </div>
                        <div>
                          <div className="text-[var(--text-secondary)]/85">
                            Sqft
                          </div>
                          <div className="text-[var(--text-primary)]">
                            <span className="number-pill number-pill-sm">
                              {house.metadata?.squareFeet?.toLocaleString?.() ??
                                "-"}
                            </span>
                          </div>
                        </div>
                      </div>

                      <div className="market-price-panel mt-4 rounded-lg border border-blue-300/25 bg-blue-500/10 px-3 py-2">
                        <p className="text-[11px] uppercase tracking-[0.08em] text-blue-100/80">
                          Listing price
                        </p>
                        <p className="text-sm font-semibold text-blue-100">
                          <span className="number-pill number-pill-sm number-pill-mono">
                            {formatListingPrice(house.listing?.price)}
                          </span>
                        </p>
                      </div>

                      <div className="mt-auto flex items-center justify-between gap-3 pt-5">
                        <div className="text-sm">
                          <div className="text-[var(--text-secondary)]/85 font-mono">
                            Token
                          </div>
                          <div className="text-[var(--text-primary)] font-mono">
                            <span className="number-pill number-pill-sm number-pill-mono">
                              #{house.tokenId}
                            </span>
                          </div>
                        </div>

                        <div className="flex flex-wrap justify-end gap-2">
                          <Link
                            to={houseOpenPath}
                            className="btn btn-secondary !px-3 !py-2 !text-xs"
                          >
                            {mode === "degen" ? "View" : "Open"}
                          </Link>
                          {canMessageOwner && (
                            <Link
                              to={messageOwnerPath}
                              className="btn btn-secondary !px-3 !py-2 !text-xs"
                            >
                              {mode === "degen"
                                ? "Message owner"
                                : "Private message owner"}
                            </Link>
                          )}
                          {isOwner && (
                            <Link
                              to={`/houses/${house.tokenId}/list`}
                              className="btn btn-primary !px-3 !py-2 !text-xs"
                            >
                              {mode === "degen" ? "List" : "Sell / Rent"}
                            </Link>
                          )}
                          {isOwner && (
                            <button
                              type="button"
                              className="btn btn-secondary !px-3 !py-2 !text-xs"
                              onClick={() => openImageEditor(house)}
                            >
                              {mode === "degen"
                                ? "Update images"
                                : "Edit images"}
                            </button>
                          )}
                          {!isOwner && listingType === "for_sale" && (
                            <button
                              className="btn btn-primary !px-3 !py-2 !text-xs disabled:cursor-not-allowed disabled:opacity-50"
                              onClick={() => handleQuickBuy(house)}
                              disabled={
                                actionLoadingToken === house.tokenId ||
                                wrongChain ||
                                (house.listing?.isPrivateSale === true &&
                                  !!house.listing.allowedBuyer &&
                                  house.listing.allowedBuyer.toLowerCase() !==
                                    (walletAddress || "").toLowerCase())
                              }
                            >
                              {actionLoadingToken === house.tokenId
                                ? "Submitting..."
                                : mode === "degen"
                                  ? "Buy"
                                  : "Buy Home"}
                            </button>
                          )}
                          {!isOwner && listingType === "for_rent" && (
                            <button
                              className="btn btn-primary !px-3 !py-2 !text-xs disabled:cursor-not-allowed disabled:opacity-50"
                              onClick={() => handleQuickRent(house)}
                              disabled={
                                actionLoadingToken === house.tokenId ||
                                wrongChain
                              }
                            >
                              {actionLoadingToken === house.tokenId
                                ? "Submitting..."
                                : mode === "degen"
                                  ? "Rent"
                                  : "Start Rent"}
                            </button>
                          )}
                        </div>
                      </div>
                    </article>
                  );
                })}
              </div>
            </div>
          )}
        </div>
      </main>

      {imageEditorHouse && (
        <div className="fixed inset-0 z-[70] flex items-center justify-center bg-black/70 px-4 py-8 backdrop-blur-sm">
          <div className="cyber-card dashboard-panel w-full max-w-3xl p-6 md:p-7">
            <div className="flex items-start justify-between gap-4">
              <div>
                <h2 className="dashboard-section-title">
                  {mode === "degen" ? "Update property images" : "Edit images"}
                </h2>
                <p className="dashboard-section-note mt-1 text-sm">
                  Token{" "}
                  <span className="number-pill number-pill-sm number-pill-mono">
                    #{imageEditorHouse.tokenId}
                  </span>{" "}
                  • {imageEditorHouse.metadata?.address || "Private property"}
                </p>
              </div>
              <button
                type="button"
                className="btn btn-secondary !px-3 !py-2 !text-xs"
                onClick={() => closeImageEditor()}
                disabled={isImageEditorSaving}
              >
                Close
              </button>
            </div>

            <div className="mt-6 space-y-4">
              <div>
                <label className="block text-sm text-[var(--text-secondary)] mb-2 font-mono">
                  Image links (one per line)
                </label>
                <textarea
                  value={imageEditorLinks}
                  onChange={(event) => setImageEditorLinks(event.target.value)}
                  placeholder={
                    "https://example.com/front.jpg\nipfs://bafy.../kitchen.png"
                  }
                  className="cyber-input min-h-[7rem] resize-y font-mono text-xs leading-6"
                  disabled={isImageEditorSaving}
                />
                <p className="form-help mt-1">
                  Supports `https://`, `http://`, `ipfs://`, and uploaded PNG
                  images.
                </p>
              </div>

              <div className="flex flex-wrap items-center gap-3">
                <button
                  type="button"
                  className="btn btn-secondary !px-3 !py-2 !text-xs"
                  onClick={() => imageEditorInputRef.current?.click()}
                  disabled={isImageEditorUploading || isImageEditorSaving}
                >
                  {isImageEditorUploading ? "Processing..." : "Upload PNG image"}
                </button>
                <input
                  ref={imageEditorInputRef}
                  type="file"
                  accept="image/png"
                  multiple
                  className="hidden"
                  onChange={handleImageEditorUpload}
                  disabled={isImageEditorUploading || isImageEditorSaving}
                />
                <span className="text-xs text-[var(--text-secondary)] font-mono">
                  {imageEditorUploads.length} PNG
                  {imageEditorUploads.length === 1 ? "" : "s"} selected
                </span>
              </div>

              {imageEditorUploads.length > 0 && (
                <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
                  {imageEditorUploads.map((image) => (
                    <div
                      key={image.id}
                      className="rounded-lg border border-[rgba(0,243,255,0.2)] bg-[rgba(5,10,19,0.85)] p-3"
                    >
                      <img
                        src={image.dataUrl}
                        alt={image.name}
                        className="h-24 w-full rounded object-cover"
                      />
                      <p className="mt-2 text-xs font-mono text-[var(--text-primary)] truncate">
                        {image.name}
                      </p>
                      <button
                        type="button"
                        className="mt-2 btn btn-secondary !px-2.5 !py-1.5 !text-[11px]"
                        onClick={() => removeImageEditorUpload(image.id)}
                        disabled={isImageEditorSaving}
                      >
                        Remove
                      </button>
                    </div>
                  ))}
                </div>
              )}

              <div className="flex flex-wrap justify-end gap-3 pt-2">
                <button
                  type="button"
                  className="btn btn-secondary !px-4 !py-2 !text-xs"
                  onClick={() => closeImageEditor()}
                  disabled={isImageEditorSaving}
                >
                  Cancel
                </button>
                <button
                  type="button"
                  className="btn btn-primary !px-4 !py-2 !text-xs"
                  onClick={handleSaveImageUpdates}
                  disabled={isImageEditorSaving}
                >
                  {isImageEditorSaving ? "Saving..." : "Save images"}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};
