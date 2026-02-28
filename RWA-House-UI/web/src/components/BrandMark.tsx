import React from "react";

type BrandMarkSize = "xs" | "sm" | "md" | "lg";

interface BrandMarkProps {
  readonly size?: BrandMarkSize;
  readonly showWordmark?: boolean;
  readonly showTagline?: boolean;
  readonly logoAsset?: "default" | "home" | "logo";
  readonly className?: string;
}

const ICON_SIZE: Record<BrandMarkSize, string> = {
  xs: "h-3 w-3",
  sm: "h-4 w-4",
  md: "h-5 w-5",
  lg: "h-7 w-7",
};

const TITLE_SIZE: Record<BrandMarkSize, string> = {
  xs: "text-xs",
  sm: "text-[13px]",
  md: "text-sm",
  lg: "text-base",
};

export const BrandMark: React.FC<BrandMarkProps> = ({
  size = "md",
  showWordmark = true,
  showTagline = false,
  logoAsset = "default",
  className,
}) => {
  const showHomeButtonLabel = logoAsset === "home" && !showWordmark;
  const logoSrc =
    logoAsset === "home"
      ? "/housemark-blue.svg"
      : logoAsset === "logo"
        ? "/Logo.png"
        : "/sovereign_realty_logo_512.png";

  return (
    <div
      className={`inline-flex items-center gap-2 ${
        showHomeButtonLabel ? "brandmark-home-inline" : ""
      } ${className ?? ""}`}
    >
      {showHomeButtonLabel && (
        <span className="brandmark-home-label">Soverign-Properties</span>
      )}
      <img
        src={logoSrc}
        alt=""
        className={`${ICON_SIZE[size]} shrink-0 object-contain`}
        loading="eager"
        decoding="async"
        aria-hidden="true"
      />

      {showWordmark && (
        <div className="leading-tight">
          <div
            className={`${TITLE_SIZE[size]} font-bold tracking-tight text-slate-50`}
          >
            PropMe
            <span className="bg-gradient-to-r from-sky-300 via-blue-400 to-blue-600 bg-clip-text text-transparent">
              CRE
            </span>
          </div>
          {showTagline && (
            <div className="text-[10px] font-medium uppercase tracking-[0.16em] text-slate-400">
              Private real estate
            </div>
          )}
        </div>
      )}
    </div>
  );
};
