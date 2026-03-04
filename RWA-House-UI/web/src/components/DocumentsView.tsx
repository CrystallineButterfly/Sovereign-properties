import React, { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import toast from "react-hot-toast";

import { apiClient } from "@shared/utils/api";
import { useAuth } from "./AuthProvider";

interface DocumentsViewProps {
  tokenId: string;
}

interface DocumentDescriptor {
  name?: string;
  type?: string;
  size?: number;
  uri?: string;
  hash?: string;
}

interface PrivateDocumentContent {
  index: number;
  name: string;
  mimeType: string;
  size: number;
  base64: string;
}

const decodeBase64ToBlob = (base64Value: string, mimeType: string): Blob => {
  const binary = window.atob(base64Value);
  const bytes = new Uint8Array(binary.length);
  for (let index = 0; index < binary.length; index += 1) {
    bytes[index] = binary.charCodeAt(index);
  }
  return new Blob([bytes], { type: mimeType || "application/octet-stream" });
};

export const DocumentsView: React.FC<DocumentsViewProps> = ({ tokenId }) => {
  const { user } = useAuth();
  const [documents, setDocuments] = useState<DocumentDescriptor[]>([]);
  const [privateDocuments, setPrivateDocuments] = useState<
    PrivateDocumentContent[]
  >([]);
  const [loading, setLoading] = useState(true);
  const [loadingPrivateBundle, setLoadingPrivateBundle] = useState(false);
  const [openingDocumentIndex, setOpeningDocumentIndex] = useState<
    number | null
  >(null);

  useEffect(() => {
    let mounted = true;
    const loadDocuments = async () => {
      try {
        setLoading(true);
        const response = await apiClient.getHouseDocuments(tokenId);
        if (!mounted) {
          return;
        }
        if (response.success && response.data) {
          setDocuments(
            Array.isArray(response.data.documents) ? response.data.documents : [],
          );
          return;
        }
        setDocuments([]);
        if (response.message) {
          toast.error(response.message);
        }
      } catch (error: unknown) {
        const message =
          error instanceof Error ? error.message : "Failed to load documents";
        if (mounted) {
          toast.error(message);
        }
      } finally {
        if (mounted) {
          setLoading(false);
        }
      }
    };

    void loadDocuments();
    return () => {
      mounted = false;
    };
  }, [tokenId]);

  const openPrivateDocument = (documentContent: PrivateDocumentContent) => {
    const blob = decodeBase64ToBlob(
      documentContent.base64,
      documentContent.mimeType,
    );
    const objectUrl = URL.createObjectURL(blob);
    const opened = window.open(objectUrl, "_blank", "noopener,noreferrer");
    if (!opened) {
      const anchor = document.createElement("a");
      anchor.href = objectUrl;
      anchor.download = documentContent.name || `document-${tokenId}.bin`;
      anchor.click();
    }
    setTimeout(() => URL.revokeObjectURL(objectUrl), 60_000);
  };

  const ensurePrivateBundleLoaded = async (): Promise<PrivateDocumentContent[]> => {
    if (privateDocuments.length > 0) {
      return privateDocuments;
    }

    setLoadingPrivateBundle(true);
    try {
      const response = await apiClient.getHouseDocumentContents(tokenId);
      if (!response.success || !response.data) {
        throw new Error(
          response.message || "Unable to load private document contents",
        );
      }
      const loadedDocuments = Array.isArray(response.data.documents)
        ? response.data.documents
        : [];
      if (loadedDocuments.length === 0) {
        throw new Error(
          "Private document bundle is unavailable for this property.",
        );
      }
      setPrivateDocuments(loadedDocuments);
      return loadedDocuments;
    } finally {
      setLoadingPrivateBundle(false);
    }
  };

  const handleOpenDocument = async (index: number) => {
    try {
      setOpeningDocumentIndex(index);
      const loadedDocuments = await ensurePrivateBundleLoaded();
      const documentContent = loadedDocuments[index] || loadedDocuments[0];
      if (!documentContent) {
        toast.error("No private document content is available to open.");
        return;
      }
      openPrivateDocument(documentContent);
      toast.success("Private document opened");
    } catch (error: unknown) {
      const message =
        error instanceof Error ? error.message : "Unable to open document";
      toast.error(message);
    } finally {
      setOpeningDocumentIndex(null);
    }
  };

  const viewerKYCLabel = useMemo(() => {
    if (!user) {
      return "Not authenticated";
    }
    if (user.kycStatus === "verified") {
      return "KYC verified";
    }
    return `KYC: ${user.kycStatus}`;
  }, [user]);
  const isViewerKYCVerified = user?.kycStatus === "verified";

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-[70vh]">
        <div className="relative">
          <div className="w-16 h-16 border-2 border-[#00f3ff] border-t-transparent rounded-full animate-spin"></div>
          <div
            className="absolute inset-0 w-16 h-16 border-2 border-[#b026ff] border-b-transparent rounded-full animate-spin"
            style={{ animationDirection: "reverse", animationDuration: "1.5s" }}
          ></div>
        </div>
      </div>
    );
  }

  return (
    <div className="page-shell page-shell-form">
      <div className="cyber-card overflow-hidden">
        <div className="form-header">
          <div>
            <h1 className="form-title">Property documents</h1>
            <p className="form-subtitle">
              Token ID{" "}
              <span className="number-pill number-pill-sm number-pill-mono">
                #{tokenId}
              </span>
            </p>
            <div
              className={`mt-2 inline-flex items-center gap-2 rounded-full border px-3 py-1 text-xs font-semibold ${
                isViewerKYCVerified
                  ? "border-emerald-400/45 bg-emerald-500/10 text-emerald-300"
                  : "border-amber-300/45 bg-amber-500/10 text-amber-200"
              }`}
            >
              <span aria-hidden>{isViewerKYCVerified ? "✅" : "⚠️"}</span>
              <span>
                {isViewerKYCVerified ? "KYC verified" : "KYC verification required"}
              </span>
            </div>
            <p className="form-help mt-2">
              Viewing private file contents requires authenticated ownership and
              a verified KYC profile. Current status: {viewerKYCLabel}.
            </p>
          </div>
          <div className="flex flex-wrap gap-2">
            <Link to={`/houses/${tokenId}`} className="cyber-btn text-sm">
              Back
            </Link>
            <Link to="/claim" className="cyber-btn cyber-btn-primary text-sm">
              Claim Key
            </Link>
          </div>
        </div>

        <div className="p-7 md:p-8">
          {documents.length === 0 ? (
            <div className="text-center text-[var(--text-secondary)] py-10">
              <p>No document metadata available</p>
              <p className="text-xs mt-2 font-mono">
                If this is a private sale/rental, use Claim Key to retrieve the
                encrypted document key.
              </p>
            </div>
          ) : (
            <div className="space-y-4">
              {documents.map((doc, index) => (
                <div
                  key={index}
                  className="flex items-center justify-between gap-6 rounded-lg border border-[rgba(0,243,255,0.15)] bg-[rgba(0,0,0,0.35)] p-5"
                >
                  <div>
                    <p className="font-medium text-white">
                      {doc.name ? (
                        doc.name
                      ) : (
                        <>
                          Document{" "}
                          <span className="number-pill number-pill-sm number-pill-mono">
                            {index + 1}
                          </span>
                        </>
                      )}
                    </p>
                    <p className="text-sm text-[var(--text-secondary)] font-mono">
                      Type: {doc.type || "unknown"} • Size:{" "}
                      <span className="number-pill number-pill-sm number-pill-mono">
                        {Number(doc.size || 0)} bytes
                      </span>
                    </p>
                    {doc.uri && (
                      <p className="text-xs text-[var(--text-secondary)] font-mono mt-1 break-all">
                        URI: {doc.uri}
                      </p>
                    )}
                  </div>
                  <button
                    onClick={() => handleOpenDocument(index)}
                    disabled={
                      loadingPrivateBundle || openingDocumentIndex === index
                    }
                    className="cyber-btn cyber-btn-primary text-sm disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {openingDocumentIndex === index
                      ? "Opening..."
                      : loadingPrivateBundle
                        ? "Loading..."
                        : "View"}
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
