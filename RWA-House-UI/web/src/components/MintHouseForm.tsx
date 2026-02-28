/**
 * Mint House Form Component
 * Secure form for minting new house RWA with document upload
 */

import React, { useState, useRef, useCallback, useMemo } from "react";
import { useAuth } from "./AuthProvider";
import { useUXMode } from "./UXModeProvider";
import { apiClient } from "../../../shared/src/utils/api";
import {
  fileToBase64,
  validateFileUpload,
  validateDocumentContent,
  sanitizeInput,
  MintFormSchema,
  generateKeyPair,
  exportPublicKey,
} from "../../../shared/src/utils/security";
import type { HouseMetadata, StorageType } from "../../../shared/src/types";
import toast from "react-hot-toast";
import { useNavigate } from "react-router-dom";

interface MintFormData {
  houseId: string;
  storageType: StorageType;
  metadata: HouseMetadata;
  documents: File[];
}

interface UploadedHouseImage {
  id: string;
  name: string;
  size: number;
  dataUrl: string;
}

type NumericMetadataField = "bedrooms" | "bathrooms" | "squareFeet";

const MAX_HOUSE_IMAGE_UPLOAD_SIZE_BYTES = 2 * 1024 * 1024;
const MAX_HOUSE_IMAGE_COUNT = 10;

const FRIENDLY_FIELD_NAME: Record<string, string> = {
  houseId: "House ID",
  "metadata.address": "Street address",
  "metadata.city": "City",
  "metadata.state": "State",
  "metadata.zipCode": "ZIP code",
  "metadata.propertyType": "Property type",
  "metadata.squareFeet": "Square feet",
  "metadata.images": "Property images",
  documents: "Property documents",
};

export const MintHouseForm: React.FC = () => {
  const { walletAddress } = useAuth();
  const { mode } = useUXMode();
  const navigate = useNavigate();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const houseImageInputRef = useRef<HTMLInputElement>(null);

  const [formData, setFormData] = useState<MintFormData>({
    houseId: "",
    storageType: "ipfs",
    metadata: {
      address: "",
      city: "",
      state: "",
      zipCode: "",
      country: "USA",
      propertyType: "single_family",
      bedrooms: 0,
      bathrooms: 0,
      squareFeet: 0,
      yearBuilt: new Date().getFullYear(),
      description: "",
      images: [],
    },
    documents: [],
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [imageUrlInput, setImageUrlInput] = useState("");
  const [uploadedHouseImages, setUploadedHouseImages] = useState<
    UploadedHouseImage[]
  >([]);
  const [success, setSuccess] = useState(false);
  const [result, setResult] = useState<{
    tokenId: string;
    txHash: string;
  } | null>(null);
  const [isValidatingFiles, setIsValidatingFiles] = useState(false);
  const [numericDrafts, setNumericDrafts] = useState<
    Record<NumericMetadataField, string>
  >({
    bedrooms: "0",
    bathrooms: "0",
    squareFeet: "0",
  });

  const totalDocumentSizeBytes = useMemo(() => {
    return formData.documents.reduce((acc, file) => acc + file.size, 0);
  }, [formData.documents]);

  const totalDocumentSizeLabel = useMemo(() => {
    if (totalDocumentSizeBytes <= 0) return "0 MB";
    return `${(totalDocumentSizeBytes / 1024 / 1024).toFixed(2)} MB`;
  }, [totalDocumentSizeBytes]);

  const formatFileSizeLabel = (bytes: number): string => {
    if (bytes <= 0) {
      return "0 KB";
    }
    return `${(bytes / 1024).toFixed(1)} KB`;
  };

  const parseImageLinkEntries = useCallback((rawValue: string): string[] => {
    return Array.from(
      new Set(
        rawValue
          .split("\n")
          .map((entry) => entry.trim())
          .filter((entry) => entry.length > 0),
      ),
    ).slice(0, MAX_HOUSE_IMAGE_COUNT);
  }, []);

  const mergeMetadataImages = useCallback(
    (rawValue: string, uploadedImages: UploadedHouseImage[]): string[] => {
      const linkedImages = parseImageLinkEntries(rawValue);
      const uploadedImageUrls = uploadedImages.map((image) => image.dataUrl);
      return Array.from(
        new Set([...uploadedImageUrls, ...linkedImages]),
      ).slice(0, MAX_HOUSE_IMAGE_COUNT);
    },
    [parseImageLinkEntries],
  );

  const getFieldError = (key: string): string | null => {
    const message = errors[key];
    return message && message.trim().length > 0 ? message : null;
  };

  // Handle input changes with sanitization
  const handleInputChange = useCallback(
    (field: string, value: any) => {
      const sanitizedValue =
        typeof value === "string" ? sanitizeInput(value) : value;

      if (field.startsWith("metadata.")) {
        const metadataField = field.replace("metadata.", "");
        setFormData((prev) => ({
          ...prev,
          metadata: {
            ...prev.metadata,
            [metadataField]: sanitizedValue,
          },
        }));
      } else {
        setFormData((prev) => ({ ...prev, [field]: sanitizedValue }));
      }

      // Clear error for this field
      if (errors[field] || errors.submit) {
        setErrors((prev) => ({ ...prev, [field]: "", submit: "" }));
      }
    },
    [errors],
  );

  const handleImageUrlsChange = useCallback(
    (rawValue: string) => {
      setImageUrlInput(rawValue);
      const mergedImages = mergeMetadataImages(rawValue, uploadedHouseImages);
      handleInputChange("metadata.images", mergedImages);
    },
    [handleInputChange, mergeMetadataImages, uploadedHouseImages],
  );

  const handleHouseImageSelect = useCallback(
    async (event: React.ChangeEvent<HTMLInputElement>) => {
      const files = Array.from(event.target.files || []);
      if (files.length === 0) {
        return;
      }

      setIsValidatingFiles(true);
      try {
        const nextUploadedImages = [...uploadedHouseImages];
        const duplicateGuard = new Set(
          nextUploadedImages.map((image) => image.id),
        );
        let rejectedCount = 0;

        for (const file of files) {
          if (nextUploadedImages.length >= MAX_HOUSE_IMAGE_COUNT) {
            rejectedCount += 1;
            continue;
          }

          const validation = validateFileUpload(
            file,
            ["image/png"],
            MAX_HOUSE_IMAGE_UPLOAD_SIZE_BYTES,
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

          const isSafeImage = await validateDocumentContent(file);
          if (!isSafeImage) {
            rejectedCount += 1;
            continue;
          }

          const base64Payload = await fileToBase64(file);
          const dataUrl = `data:image/png;base64,${base64Payload}`;
          nextUploadedImages.push({
            id: imageId,
            name: file.name,
            size: file.size,
            dataUrl,
          });
          duplicateGuard.add(imageId);
        }

        setUploadedHouseImages(nextUploadedImages);
        const mergedImages = mergeMetadataImages(imageUrlInput, nextUploadedImages);
        handleInputChange("metadata.images", mergedImages);
        setErrors((prev) => ({
          ...prev,
          "metadata.images":
            rejectedCount > 0
              ? `${rejectedCount} image${
                rejectedCount === 1 ? "" : "s"
              } rejected. Upload PNG only (max 2 MB each).`
              : "",
        }));
      } catch (error) {
        setErrors((prev) => ({
          ...prev,
          "metadata.images": "Unable to process uploaded PNG images.",
        }));
      } finally {
        setIsValidatingFiles(false);
        event.target.value = "";
      }
    },
    [
      handleInputChange,
      imageUrlInput,
      mergeMetadataImages,
      uploadedHouseImages,
    ],
  );

  const removeHouseImage = useCallback(
    (imageId: string) => {
      setUploadedHouseImages((previous) => {
        const nextUploadedImages = previous.filter((image) => image.id !== imageId);
        const mergedImages = mergeMetadataImages(imageUrlInput, nextUploadedImages);
        handleInputChange("metadata.images", mergedImages);
        return nextUploadedImages;
      });
    },
    [handleInputChange, imageUrlInput, mergeMetadataImages],
  );

  const handleNumericMetadataChange = useCallback(
    (field: NumericMetadataField, rawValue: string) => {
      const pattern = field === "bathrooms" ? /^\d*(?:\.\d*)?$/ : /^\d*$/;

      if (!pattern.test(rawValue)) {
        return;
      }

      const parsedValue = rawValue === "" ? 0 : Number(rawValue);
      if (!Number.isFinite(parsedValue)) {
        return;
      }

      if ((field === "bedrooms" || field === "bathrooms") && parsedValue > 20) {
        return;
      }

      setNumericDrafts((prev) => ({ ...prev, [field]: rawValue }));
      handleInputChange(`metadata.${field}`, parsedValue);
    },
    [handleInputChange],
  );

  // Handle file selection
  const handleFileSelect = useCallback(
    async (event: React.ChangeEvent<HTMLInputElement>) => {
      const files = Array.from(event.target.files || []);
      const allowedTypes = [
        "application/pdf",
        "image/jpeg",
        "image/png",
        "image/heic",
      ];
      const maxSize = 10 * 1024 * 1024; // 10MB

      if (files.length === 0) {
        return;
      }

      setIsValidatingFiles(true);

      try {
        const validFiles: File[] = [];
        const newErrors: Record<string, string> = {};
        const duplicateGuard = new Set<string>();
        let rejectedCount = 0;

        for (const [index, file] of files.entries()) {
          const dedupeKey = `${file.name}:${file.size}:${file.lastModified}`;
          if (duplicateGuard.has(dedupeKey)) {
            rejectedCount += 1;
            continue;
          }

          duplicateGuard.add(dedupeKey);
          const validation = validateFileUpload(file, allowedTypes, maxSize);
          if (!validation.valid) {
            rejectedCount += 1;
            newErrors[`file_${index}`] = validation.error || "Invalid file";
            continue;
          }

          if (file.type !== "image/heic") {
            try {
              const isSafeDocument = await validateDocumentContent(file);
              if (!isSafeDocument) {
                rejectedCount += 1;
                newErrors[`file_${index}`] =
                  "File content does not match supported document signatures";
                continue;
              }
            } catch {
              rejectedCount += 1;
              newErrors[`file_${index}`] = "Unable to validate file contents";
              continue;
            }
          }

          validFiles.push(file);
        }

        // Check total size
        const totalSize = validFiles.reduce((acc, file) => acc + file.size, 0);
        if (totalSize > 50 * 1024 * 1024) {
          newErrors.documents = "Total file size must be less than 50MB";
        }

        setFormData((prev) => ({ ...prev, documents: validFiles }));
        setErrors((prev) => ({
          ...prev,
          ...newErrors,
          documents:
            newErrors.documents ||
            (rejectedCount > 0
              ? `${rejectedCount} file${rejectedCount === 1 ? "" : "s"} were rejected. Supported files: PDF, JPG, PNG, HEIC.`
              : ""),
          submit: "",
        }));
      } finally {
        setIsValidatingFiles(false);
        event.target.value = "";
      }
    },
    [],
  );

  // Remove file
  const removeFile = useCallback((index: number) => {
    setFormData((prev) => ({
      ...prev,
      documents: prev.documents.filter((_, i) => i !== index),
    }));
  }, []);

  // Validate form
  const validateForm = useCallback((): boolean => {
    try {
      MintFormSchema.parse(formData);
      return true;
    } catch (error: any) {
      const formattedErrors: Record<string, string> = {};
      error.errors.forEach((err: any) => {
        const path = err.path.join(".");
        formattedErrors[path] = err.message;
      });
      const firstThreeErrors = error.errors
        .slice(0, 3)
        .map((err: any) => {
          const path = err.path.join(".");
          const label = FRIENDLY_FIELD_NAME[path] || path || "Form field";
          return `${label}: ${err.message}`;
        })
        .join(" • ");
      const fallbackSummary = Object.values(formattedErrors)
        .filter((message) => message && message.trim().length > 0)
        .slice(0, 3)
        .join(" • ");
      const summaryMessage =
        firstThreeErrors ||
        fallbackSummary ||
        "Please fix validation issues before submitting";
      setErrors({
        ...formattedErrors,
        submit: summaryMessage,
      });
      toast.error(summaryMessage);
      return false;
    }
  }, [formData]);

  // Handle form submission
  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    if (formData.documents.length === 0) {
      const message = "Upload at least one property document before minting";
      setErrors((prev) => ({ ...prev, documents: message }));
      toast.error(message);
      return;
    }

    if (!validateForm()) {
      return;
    }

    if (!walletAddress) {
      setErrors({ submit: "Please connect your wallet first" });
      toast.error("Connect your wallet to mint");
      return;
    }

    setIsSubmitting(true);
    setUploadProgress(0);
    setErrors({});

    try {
      // Generate key pair for document encryption
      const keyPair = await generateKeyPair();
      const publicKeyPEM = await exportPublicKey(keyPair.publicKey);

      // Convert documents to base64
      setUploadProgress(20);
      const documentBlobs = await Promise.all(
        formData.documents.map(async (file, index) => {
          const base64 = await fileToBase64(file);
          setUploadProgress(
            20 + ((index + 1) / formData.documents.length) * 30,
          );
          return base64;
        }),
      );

      // Combine all documents
      const combinedDocs = JSON.stringify({
        files: documentBlobs,
        metadata: formData.documents.map((f) => ({
          name: f.name,
          type: f.type,
          size: f.size,
        })),
      });

      const documentsB64 = btoa(combinedDocs);
      setUploadProgress(60);

      // Prepare request payload
      const payload = {
        action: "mint",
        ownerAddress: walletAddress,
        houseId: formData.houseId,
        documentsB64,
        storageType: formData.storageType,
        ownerPublicKey: publicKeyPEM,
        metadata: formData.metadata,
      };

      setUploadProgress(80);

      // Submit to API
      const response = await apiClient.mintHouse(
        payload as import("../../../shared/src/types").MintRequestPayload,
      );
      setUploadProgress(100);

      if (response.success && response.data) {
        toast.success("Mint request accepted");
        const txHash = response.txHash || (response.data as any).txHash || "";
        setResult({
          tokenId: response.data.tokenId || (response.data as any).tokenID,
          txHash,
        });
        setSuccess(true);
      } else {
        toast.error(response.message || "Minting failed");
        setErrors({ submit: response.message || "Minting failed" });
      }
    } catch (error: any) {
      console.error("Mint error:", error);
      const isNetworkError =
        error instanceof TypeError ||
        (error?.message &&
          /failed to fetch|network/i.test(String(error.message).toLowerCase()));
      const message = isNetworkError
        ? `Cannot reach backend at ${apiClient.getBaseURL()}. Start zkpassport-session-service and retry.`
        : error.message || "An unexpected error occurred";
      toast.error(message);
      setErrors({ submit: message });
    } finally {
      setIsSubmitting(false);
    }
  };

  if (success && result) {
    return (
      <div className="page-shell page-shell-form">
        <div className="cyber-card p-8 text-center">
          <div className="w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4 bg-[rgba(0,255,136,0.12)] border border-[rgba(0,255,136,0.35)]">
            <svg
              className="w-8 h-8 text-[#00ff88]"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M5 13l4 4L19 7"
              />
            </svg>
          </div>
          <h2
            className="text-2xl mb-2"
            style={{ fontFamily: "Orbitron, sans-serif" }}
          >
            MINT COMPLETE
          </h2>
          <p className="text-[var(--text-secondary)] mb-6">
            {mode === "degen"
              ? "Your property is now an onchain RWA record with encrypted docs."
              : "Your property has been added with private encrypted documents."}
          </p>
          <p className="text-sm text-[#c9dcff] mb-6">
            Next: choose whether to list this property for sale or for rent,
            then publish terms onchain.
          </p>

          <div className="cyber-card p-6 text-left max-w-xl mx-auto">
            <div className="text-xs text-[var(--text-secondary)] font-mono">
              Token ID
            </div>
            <div className="mt-1">
              <span className="number-pill number-pill-lg number-pill-mono">
                {result.tokenId}
              </span>
            </div>
            <div className="mt-4 text-xs text-[var(--text-secondary)] font-mono">
              Transaction
            </div>
            {result.txHash ? (
              <>
                <a
                  href={`https://sepolia.etherscan.io/tx/${result.txHash}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="number-pill number-pill-sm number-pill-mono mt-1 text-[#00f3ff] hover:underline break-all"
                >
                  {result.txHash}
                </a>
                <a
                  href={`https://sepolia.etherscan.io/tx/${result.txHash}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="cyber-btn cyber-btn-primary mt-3 inline-flex items-center justify-center text-sm"
                >
                  Onchain confirmation link
                </a>
              </>
            ) : (
              <p className="text-sm text-amber-200">
                Transaction hash unavailable. Check backend logs for
                confirmation.
              </p>
            )}

            <div className="mt-5 flex flex-col sm:flex-row gap-2">
              <button
                className="cyber-btn text-sm"
                onClick={async () => {
                  await navigator.clipboard.writeText(result.tokenId);
                  toast.success("Token ID copied");
                }}
                type="button"
              >
                Copy Token ID
              </button>
              <button
                className="cyber-btn text-sm"
                onClick={async () => {
                  await navigator.clipboard.writeText(result.txHash);
                  toast.success("Tx hash copied");
                }}
                type="button"
                disabled={!result.txHash}
              >
                Copy Tx Hash
              </button>
            </div>
          </div>

          <div className="cyber-card p-6 text-left max-w-xl mx-auto mt-5">
            <h3 className="text-base font-semibold text-[#dbe7ff] mb-3">
              Recommended next steps
            </h3>
            <ol className="list-decimal pl-5 space-y-2 text-sm text-[#b6c8ea]">
              <li>
                Set listing terms via CRE workflow (
                <span className="font-mono">create_listing</span>).
              </li>
              <li>
                Publish for sale or rent so buyers/renters can execute from
                Marketplace.
              </li>
              <li>Track settlement and bills from the property dashboard.</li>
            </ol>
            <div className="mt-5 grid gap-3 sm:grid-cols-2">
              <button
                onClick={() =>
                  navigate(`/houses/${result.tokenId}/list?type=for_sale`)
                }
                className="cyber-btn cyber-btn-primary"
                type="button"
              >
                Set sale price
              </button>
              <button
                onClick={() =>
                  navigate(`/houses/${result.tokenId}/list?type=for_rent`)
                }
                className="cyber-btn cyber-btn-primary"
                type="button"
              >
                Set rent terms
              </button>
              <button
                onClick={() => navigate("/marketplace")}
                className="cyber-btn sm:col-span-2"
                type="button"
              >
                Open marketplace
              </button>
            </div>
          </div>

          <div className="mt-8 flex flex-col sm:flex-row gap-3 justify-center">
            <button
              onClick={() => navigate(`/houses/${result.tokenId}`)}
              className="cyber-btn cyber-btn-primary"
            >
              View Asset
            </button>
            <button
              onClick={() => navigate("/dashboard")}
              className="cyber-btn"
            >
              Back To Dashboard
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="page-shell page-shell-form">
      <div className="cyber-card overflow-hidden">
        <div className="form-header">
          <div>
            <h1 className="form-title">MINT HOUSE RWA</h1>
            <p className="form-subtitle">
              {mode === "degen"
                ? "Upload legal docs, encrypt them, and mint a private RWA token. CRE acts as the mediator for secure transfers."
                : "Upload property documents to create a private digital property record. CRE handles secure processing in the background."}
            </p>
          </div>
        </div>

        <form
          onSubmit={handleSubmit}
          className="mx-auto max-w-4xl p-6 md:p-7 space-y-6"
        >
          {errors.submit && !isSubmitting && (
            <div
              role="status"
              aria-live="polite"
              className="text-panel border border-rose-400/45 bg-rose-500/10 px-4 py-3 text-sm text-rose-200"
            >
              {errors.submit}
            </div>
          )}

          {/* House ID */}
          <div className="form-field">
            <label
              htmlFor="mint-house-id"
              className="block text-sm text-[var(--text-secondary)] mb-2 font-mono"
            >
              House ID <span className="text-[#ff3366]">*</span>
            </label>
            <input
              id="mint-house-id"
              type="text"
              value={formData.houseId}
              onChange={(e) => handleInputChange("houseId", e.target.value)}
              placeholder="e.g., 123-main-st"
              className={`cyber-input font-mono ${errors.houseId ? "border-[#ff3366]" : ""}`}
              disabled={isSubmitting}
              aria-invalid={Boolean(getFieldError("houseId"))}
              aria-describedby={`mint-house-id-help${getFieldError("houseId") ? " mint-house-id-error" : ""}`}
            />
            {errors.houseId && (
              <p id="mint-house-id-error" className="form-error mt-1">
                {errors.houseId}
              </p>
            )}
            <p id="mint-house-id-help" className="form-help mt-1">
              Unique identifier for your property (letters, numbers, hyphens
              only)
            </p>
          </div>

          {/* Property Address */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="md:col-span-2">
              <label
                htmlFor="mint-address"
                className="block text-sm text-[var(--text-secondary)] mb-2 font-mono"
              >
                Street Address <span className="text-[#ff3366]">*</span>
              </label>
              <input
                id="mint-address"
                type="text"
                value={formData.metadata.address}
                onChange={(e) =>
                  handleInputChange("metadata.address", e.target.value)
                }
                placeholder="123 Main Street"
                className="cyber-input"
                disabled={isSubmitting}
                aria-invalid={Boolean(getFieldError("metadata.address"))}
                aria-describedby={
                  getFieldError("metadata.address")
                    ? "mint-address-error"
                    : undefined
                }
              />
              {errors["metadata.address"] && (
                <p id="mint-address-error" className="form-error mt-1">
                  {errors["metadata.address"]}
                </p>
              )}
            </div>

            <div>
              <label
                htmlFor="mint-city"
                className="block text-sm text-[var(--text-secondary)] mb-2 font-mono"
              >
                City *
              </label>
              <input
                id="mint-city"
                type="text"
                value={formData.metadata.city}
                onChange={(e) =>
                  handleInputChange("metadata.city", e.target.value)
                }
                className="cyber-input"
                disabled={isSubmitting}
                aria-invalid={Boolean(getFieldError("metadata.city"))}
                aria-describedby={
                  getFieldError("metadata.city") ? "mint-city-error" : undefined
                }
              />
              {errors["metadata.city"] && (
                <p id="mint-city-error" className="form-error mt-1">
                  {errors["metadata.city"]}
                </p>
              )}
            </div>

            <div>
              <label
                htmlFor="mint-state"
                className="block text-sm text-[var(--text-secondary)] mb-2 font-mono"
              >
                State *
              </label>
              <input
                id="mint-state"
                type="text"
                value={formData.metadata.state}
                onChange={(e) =>
                  handleInputChange("metadata.state", e.target.value)
                }
                className="cyber-input"
                disabled={isSubmitting}
                aria-invalid={Boolean(getFieldError("metadata.state"))}
                aria-describedby={
                  getFieldError("metadata.state")
                    ? "mint-state-error"
                    : undefined
                }
              />
              {errors["metadata.state"] && (
                <p id="mint-state-error" className="form-error mt-1">
                  {errors["metadata.state"]}
                </p>
              )}
            </div>

            <div>
              <label
                htmlFor="mint-zip"
                className="block text-sm text-[var(--text-secondary)] mb-2 font-mono"
              >
                ZIP Code *
              </label>
              <input
                id="mint-zip"
                type="text"
                value={formData.metadata.zipCode}
                onChange={(e) =>
                  handleInputChange("metadata.zipCode", e.target.value)
                }
                placeholder="12345"
                className="cyber-input font-mono"
                disabled={isSubmitting}
                aria-invalid={Boolean(getFieldError("metadata.zipCode"))}
                aria-describedby={
                  getFieldError("metadata.zipCode")
                    ? "mint-zip-error"
                    : undefined
                }
              />
              {errors["metadata.zipCode"] && (
                <p id="mint-zip-error" className="form-error mt-1">
                  {errors["metadata.zipCode"]}
                </p>
              )}
            </div>

            <div>
              <label
                htmlFor="mint-property-type"
                className="block text-sm text-[var(--text-secondary)] mb-2 font-mono"
              >
                Property Type *
              </label>
              <select
                id="mint-property-type"
                value={formData.metadata.propertyType}
                onChange={(e) =>
                  handleInputChange("metadata.propertyType", e.target.value)
                }
                className="cyber-input"
                disabled={isSubmitting}
                aria-invalid={Boolean(getFieldError("metadata.propertyType"))}
                aria-describedby={
                  getFieldError("metadata.propertyType")
                    ? "mint-property-type-error"
                    : undefined
                }
              >
                <option value="single_family">Single Family</option>
                <option value="condo">Condo</option>
                <option value="townhouse">Townhouse</option>
                <option value="multi_family">Multi-Family</option>
                <option value="apartment">Apartment</option>
                <option value="commercial">Commercial</option>
              </select>
              {errors["metadata.propertyType"] && (
                <p id="mint-property-type-error" className="form-error mt-1">
                  {errors["metadata.propertyType"]}
                </p>
              )}
            </div>
          </div>

          {/* Property Details */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <label
                htmlFor="mint-bedrooms"
                className="block text-sm text-[var(--text-secondary)] mb-2 font-mono"
              >
                Bedrooms
              </label>
              <input
                id="mint-bedrooms"
                type="text"
                inputMode="numeric"
                pattern="[0-9]*"
                value={numericDrafts.bedrooms}
                onChange={(e) =>
                  handleNumericMetadataChange("bedrooms", e.target.value)
                }
                className="cyber-input font-mono"
                disabled={isSubmitting}
              />
            </div>

            <div>
              <label
                htmlFor="mint-bathrooms"
                className="block text-sm text-[var(--text-secondary)] mb-2 font-mono"
              >
                Bathrooms
              </label>
              <input
                id="mint-bathrooms"
                type="text"
                inputMode="decimal"
                value={numericDrafts.bathrooms}
                onChange={(e) =>
                  handleNumericMetadataChange("bathrooms", e.target.value)
                }
                className="cyber-input font-mono"
                disabled={isSubmitting}
              />
            </div>

            <div>
              <label
                htmlFor="mint-square-feet"
                className="block text-sm text-[var(--text-secondary)] mb-2 font-mono"
              >
                Sq Ft
              </label>
              <input
                id="mint-square-feet"
                type="text"
                inputMode="numeric"
                pattern="[0-9]*"
                value={numericDrafts.squareFeet}
                onChange={(e) =>
                  handleNumericMetadataChange("squareFeet", e.target.value)
                }
                className="cyber-input font-mono"
                disabled={isSubmitting}
                aria-invalid={Boolean(getFieldError("metadata.squareFeet"))}
                aria-describedby={
                  getFieldError("metadata.squareFeet")
                    ? "mint-square-feet-error"
                    : undefined
                }
              />
              {errors["metadata.squareFeet"] && (
                <p id="mint-square-feet-error" className="form-error mt-1">
                  {errors["metadata.squareFeet"]}
                </p>
              )}
            </div>

            <div>
              <label
                htmlFor="mint-year-built"
                className="block text-sm text-[var(--text-secondary)] mb-2 font-mono"
              >
                Year Built
              </label>
              <input
                id="mint-year-built"
                type="number"
                min="1800"
                max={new Date().getFullYear()}
                value={formData.metadata.yearBuilt}
                onChange={(e) =>
                  handleInputChange(
                    "metadata.yearBuilt",
                    parseInt(e.target.value) || new Date().getFullYear(),
                  )
                }
                className="cyber-input font-mono"
                disabled={isSubmitting}
              />
            </div>
          </div>

          {/* Property Images */}
          <div className="form-field">
            <label
              htmlFor="mint-image-urls"
              className="block text-sm text-[var(--text-secondary)] mb-2 font-mono"
            >
              Property Images (links + PNG upload)
            </label>
            <textarea
              id="mint-image-urls"
              value={imageUrlInput}
              onChange={(e) => handleImageUrlsChange(e.target.value)}
              placeholder={
                "https://example.com/front.jpg\nipfs://bafy.../kitchen.png"
              }
              className="cyber-input min-h-[7rem] resize-y font-mono text-xs leading-6"
              disabled={isSubmitting}
              aria-invalid={Boolean(getFieldError("metadata.images"))}
              aria-describedby={`mint-image-urls-help${getFieldError("metadata.images") ? " mint-image-urls-error" : ""}`}
            />
            {errors["metadata.images"] && (
              <p id="mint-image-urls-error" className="form-error mt-1">
                {errors["metadata.images"]}
              </p>
            )}
            <div className="mt-3 flex flex-wrap items-center gap-3">
              <button
                type="button"
                onClick={() => houseImageInputRef.current?.click()}
                className="btn btn-secondary !px-3 !py-2 !text-xs"
                disabled={isSubmitting || isValidatingFiles}
              >
                {isValidatingFiles ? "Processing..." : "Upload PNG image"}
              </button>
              <input
                ref={houseImageInputRef}
                type="file"
                accept="image/png"
                multiple
                className="hidden"
                onChange={handleHouseImageSelect}
                disabled={isSubmitting || isValidatingFiles}
              />
              <span className="text-xs text-[var(--text-secondary)] font-mono">
                {uploadedHouseImages.length} PNG
                {uploadedHouseImages.length === 1 ? "" : "s"} uploaded
              </span>
            </div>
            {uploadedHouseImages.length > 0 && (
              <div className="mt-3 grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
                {uploadedHouseImages.map((image) => (
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
                    <p className="text-[11px] text-[var(--text-secondary)]">
                      {formatFileSizeLabel(image.size)}
                    </p>
                    <button
                      type="button"
                      onClick={() => removeHouseImage(image.id)}
                      className="mt-2 btn btn-secondary !px-2.5 !py-1.5 !text-[11px]"
                      disabled={isSubmitting}
                    >
                      Remove
                    </button>
                  </div>
                ))}
              </div>
            )}
            <p id="mint-image-urls-help" className="form-help mt-1">
              Optional. One link per line (https://, http://, or ipfs://) and/or
              upload PNG files (max 2 MB each, up to 10 total). The first valid
              image is used as the house thumbnail in dashboard and marketplace.
            </p>
          </div>

          {/* Storage Type */}
          <div className="text-panel border border-blue-300/35 bg-blue-500/5 p-4 md:p-5">
            <label
              className="block text-sm text-[var(--text-secondary)] mb-2 font-mono"
              id="mint-storage-type-label"
            >
              Document Storage Type
            </label>
            <div
              className="grid grid-cols-1 sm:grid-cols-2 gap-3 max-w-2xl"
              role="radiogroup"
              aria-labelledby="mint-storage-type-label"
            >
              <button
                type="button"
                onClick={() => handleInputChange("storageType", "ipfs")}
                disabled={isSubmitting}
                role="radio"
                aria-checked={formData.storageType === "ipfs"}
                className={`storage-option-card ${
                  formData.storageType === "ipfs"
                    ? "storage-option-card-active"
                    : ""
                }`}
              >
                <span
                  className="storage-option-card-indicator"
                  aria-hidden="true"
                />
                <p className="text-sm font-semibold text-[var(--text-primary)] mt-3">
                  IPFS (decentralized)
                </p>
                <p className="text-xs text-[var(--text-secondary)] mt-1 leading-5">
                  Publicly addressable encrypted document pointer.
                </p>
              </button>
              <button
                type="button"
                onClick={() => handleInputChange("storageType", "offchain")}
                disabled={isSubmitting}
                role="radio"
                aria-checked={formData.storageType === "offchain"}
                className={`storage-option-card ${
                  formData.storageType === "offchain"
                    ? "storage-option-card-active"
                    : ""
                }`}
              >
                <span
                  className="storage-option-card-indicator"
                  aria-hidden="true"
                />
                <p className="text-sm font-semibold text-[var(--text-primary)] mt-3">
                  Private cloud
                </p>
                <p className="text-xs text-[var(--text-secondary)] mt-1 leading-5">
                  Offchain encrypted storage reference.
                </p>
              </button>
            </div>
          </div>

          {/* Document Upload */}
          <div className="text-panel border border-blue-300/35 bg-blue-500/5 p-4 md:p-5">
            <label
              className="block text-sm text-[var(--text-secondary)] mb-2 font-mono"
              id="mint-documents-label"
            >
              Property Documents <span className="text-[#ff3366]">*</span>
            </label>
            <div
              className={`document-upload-panel mx-auto max-w-2xl rounded-xl p-4 transition-colors border border-dashed ${
                errors.documents
                  ? "border-[#ff3366] bg-[rgba(255,51,102,0.08)]"
                  : "border-[rgba(0,243,255,0.35)] bg-[rgba(0,243,255,0.04)]"
              }`}
              aria-labelledby="mint-documents-label"
              role="button"
              tabIndex={isSubmitting || isValidatingFiles ? -1 : 0}
              onClick={() => {
                if (isSubmitting || isValidatingFiles) return;
                fileInputRef.current?.click();
              }}
              onKeyDown={(event) => {
                if (isSubmitting || isValidatingFiles) return;
                if (event.key === "Enter" || event.key === " ") {
                  event.preventDefault();
                  fileInputRef.current?.click();
                }
              }}
            >
              <input
                ref={fileInputRef}
                type="file"
                multiple
                accept=".pdf,.jpg,.jpeg,.png,.heic"
                onChange={handleFileSelect}
                className="hidden"
                disabled={isSubmitting}
              />
              <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                <div className="flex items-start gap-3 text-left">
                  <div className="document-upload-icon" aria-hidden="true">
                    <svg
                      className="h-5 w-5 text-[#8ec8ff]"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={1.8}
                        d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a4 4 0 010 8h-1m-4-2v9m0 0l-3-3m3 3l3-3"
                      />
                    </svg>
                  </div>
                  <div className="space-y-1">
                    <p className="text-sm font-semibold text-[var(--text-primary)]">
                      Add supporting documents
                    </p>
                    <p className="text-[11px] text-[var(--text-secondary)] leading-5">
                      PDF, JPG, PNG, HEIC • Max 10MB each • Max 50MB total
                    </p>
                    <p className="text-[11px] text-[var(--text-secondary)]">
                      <span className="number-pair">
                        <span className="number-pill number-pill-sm">
                          {formData.documents.length}
                        </span>
                        <span>
                          file{formData.documents.length === 1 ? "" : "s"}{" "}
                          selected
                        </span>
                      </span>
                      {" • "}
                      <span className="number-pill number-pill-sm number-pill-mono">
                        {totalDocumentSizeLabel}
                      </span>
                    </p>
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  <button
                    type="button"
                    className="cyber-btn text-xs px-3"
                    onClick={(event) => {
                      event.stopPropagation();
                      fileInputRef.current?.click();
                    }}
                    disabled={isSubmitting || isValidatingFiles}
                  >
                    {isValidatingFiles ? "Validating..." : "Choose files"}
                  </button>
                </div>
              </div>
            </div>

            {formData.documents.length > 0 && (
              <div className="mt-4 space-y-2">
                {formData.documents.map((file, index) => (
                  <div
                    key={index}
                    className="flex items-center justify-between bg-[rgba(0,0,0,0.35)] p-3 rounded border border-[rgba(0,243,255,0.15)]"
                  >
                    <div className="flex min-w-0 items-center">
                      <svg
                        className="w-5 h-5 text-[#00f3ff] mr-2"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z"
                          clipRule="evenodd"
                        />
                      </svg>
                      <span className="truncate text-sm text-[var(--text-primary)] font-mono">
                        {file.name}
                      </span>
                      <span className="text-xs text-[var(--text-secondary)] ml-2">
                        <span className="number-pill number-pill-sm number-pill-mono">
                          {(file.size / 1024 / 1024).toFixed(2)} MB
                        </span>
                      </span>
                    </div>
                    <button
                      type="button"
                      onClick={() => removeFile(index)}
                      className="inline-flex min-h-[2.75rem] min-w-[2.75rem] items-center justify-center rounded-lg text-[#ff3366] transition hover:bg-rose-500/10 hover:text-[#ff3366]/80"
                      disabled={isSubmitting}
                    >
                      <svg
                        className="w-5 h-5"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                          clipRule="evenodd"
                        />
                      </svg>
                    </button>
                  </div>
                ))}
              </div>
            )}
            {errors.documents && (
              <p className="form-error mt-1">{errors.documents}</p>
            )}
          </div>

          {/* Progress Bar */}
          {isSubmitting && (
            <div className="cyber-card p-4">
              <p className="text-sm text-[var(--text-secondary)] mb-2 font-mono">
                Processing (encrypting + uploading)...
              </p>
              <div className="w-full bg-[rgba(0,243,255,0.12)] rounded-full h-2">
                <div
                  className="bg-[#00f3ff] h-2 rounded-full transition-all duration-300"
                  style={{ width: `${uploadProgress}%` }}
                ></div>
              </div>
              <p className="text-xs text-[var(--text-secondary)] mt-1 font-mono">
                <span className="number-pill number-pill-sm number-pill-mono">
                  {uploadProgress}%
                </span>{" "}
                complete
              </p>
            </div>
          )}

          {/* Submit Button */}
          <div className="flex space-x-4">
            <button
              type="button"
              onClick={() => navigate(-1)}
              className="cyber-btn flex-1"
              disabled={isSubmitting}
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isSubmitting}
              className="cyber-btn cyber-btn-primary flex-1 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {isSubmitting
                ? "Minting..."
                : mode === "degen"
                  ? "Mint House RWA"
                  : "Add Property"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
