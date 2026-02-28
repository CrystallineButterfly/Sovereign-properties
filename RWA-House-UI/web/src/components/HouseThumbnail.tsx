import React, { useEffect, useMemo, useState } from "react";

import type { House } from "../../../shared/src/types";

const toRenderableImageUrl = (rawUrl: string): string => {
  const trimmed = rawUrl.trim();
  if (!trimmed) {
    return "";
  }

  if (!trimmed.toLowerCase().startsWith("ipfs://")) {
    return trimmed;
  }

  const withoutScheme = trimmed.slice("ipfs://".length).replace(/^ipfs\//i, "");
  return `https://ipfs.io/ipfs/${withoutScheme}`;
};

const getPrimaryImageUrl = (house: House): string | null => {
  const images = house.metadata?.images;
  if (!Array.isArray(images)) {
    return null;
  }

  for (const imageUrl of images) {
    if (typeof imageUrl !== "string") {
      continue;
    }

    const normalized = toRenderableImageUrl(imageUrl);
    if (normalized) {
      return normalized;
    }
  }

  return null;
};

interface HouseThumbnailProps {
  house: House;
  className?: string;
}

export const HouseThumbnail: React.FC<HouseThumbnailProps> = ({
  house,
  className = "",
}) => {
  const imageUrl = useMemo(() => getPrimaryImageUrl(house), [house.metadata?.images]);
  const [isBroken, setIsBroken] = useState(false);

  useEffect(() => {
    setIsBroken(false);
  }, [imageUrl]);

  const baseClassName = ["house-thumbnail", className].filter(Boolean).join(" ");
  const altText = house.metadata?.address
    ? `Property at ${house.metadata.address}`
    : `Token #${house.tokenId}`;

  if (imageUrl && !isBroken) {
    return (
      <div className={baseClassName}>
        <img
          src={imageUrl}
          alt={altText}
          className="house-thumbnail-image"
          loading="lazy"
          onError={() => {
            setIsBroken(true);
          }}
        />
      </div>
    );
  }

  return (
    <div className={baseClassName} role="img" aria-label={`${altText} (no image)`}>
      <div className="house-thumbnail-placeholder">
        <svg
          className="house-thumbnail-icon"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          aria-hidden="true"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={1.6}
            d="M3 11.5l9-7 9 7M5 10.5V20h14v-9.5M10 20v-5h4v5"
          />
        </svg>
        <span className="house-thumbnail-label">No image</span>
      </div>
    </div>
  );
};
