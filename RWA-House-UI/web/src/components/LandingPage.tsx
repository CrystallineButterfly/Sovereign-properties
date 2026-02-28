import React from "react";
import { Link, useNavigate } from "react-router-dom";

import { useAuth } from "./AuthProvider";
import { BrandMark } from "./BrandMark";
import { useUXMode } from "./UXModeProvider";

const STATS: Array<{ readonly value: string; readonly label: string }> = [
  { value: "$2.4M+", label: "Asset value secured" },
  { value: "150+", label: "Properties onboarded" },
  { value: "<15 min", label: "Average listing setup" },
  { value: "24/7", label: "Private CRE workflows" },
];

const CAPABILITY_CHIPS: readonly string[] = [
  "Encrypted legal documents",
  "KYC-ready workflows",
  "CRE automation",
  "Onchain ownership records",
];

const WORKFLOW_STEPS: Array<{ readonly title: string; readonly body: string }> =
  [
    {
      title: "Create a secure profile",
      body: "Sign in with Privy, pass verification, and get an embedded wallet in minutes.",
    },
    {
      title: "List property with encrypted docs",
      body: "Upload ownership docs once. We encrypt and route them through CRE workflows.",
    },
    {
      title: "Close sale, rental, or bill payment",
      body: "Buyers, renters, and owners settle through secure onchain and offchain rails.",
    },
  ];

const WORKSPACE_CAPABILITIES: Array<{
  readonly label: string;
  readonly status: string;
  readonly guidedSentence: string;
  readonly proSentence: string;
  readonly tone: "ready" | "protected" | "enabled" | "live";
}> = [
  {
    label: "Tokenization service",
    status: "Ready",
    guidedSentence:
      "Create and manage tokenized real-estate assets for sale and rental workflows.",
    proSentence:
      "Solidity ERC-721 contracts mint and track property state while CRE workflows coordinate each execution step.",
    tone: "ready",
  },
  {
    label: "Document security channel",
    status: "Protected",
    guidedSentence:
      "Share ownership and legal documents through private, access-controlled delivery.",
    proSentence:
      "Legal documents stay encrypted offchain; CRE manages recipient-scoped key release behind permissioned contract flows.",
    tone: "protected",
  },
  {
    label: "Billing and rent automation",
    status: "Enabled",
    guidedSentence:
      "Run recurring rent and bill collection with workflow-based automation.",
    proSentence:
      "Recurring payment jobs run through CRE automation with contract-verified state transitions and retry-safe execution.",
    tone: "enabled",
  },
  {
    label: "Settlement operations",
    status: "Live",
    guidedSentence:
      "Execute and monitor property payment settlement across connected rails.",
    proSentence:
      "Escrow and settlement events finalize on smart contracts while CRE pipelines sync post-trade status and notifications.",
    tone: "live",
  },
];

const FEATURE_CARDS: Array<{
  readonly icon: string;
  readonly title: string;
  readonly guidedDescription: string;
  readonly proDescription: string;
}> = [
  {
    icon: "🏠",
    title: "RWA-native listings",
    guidedDescription:
      "Create sale and rent listings with ownership records anchored onchain.",
    proDescription:
      "CRE-triggered listing pipeline that writes ownership and lifecycle updates through Solidity contract entrypoints.",
  },
  {
    icon: "🔐",
    title: "Encrypted legal exchange",
    guidedDescription:
      "Transfer documents to counterparties without exposing raw files publicly.",
    proDescription:
      "End-to-end encrypted legal payload delivery with key escrow/release controlled by CRE workflow conditions.",
  },
  {
    icon: "💳",
    title: "Unified payment operations",
    guidedDescription:
      "Handle rents and bills with onchain settlement and fiat rails where needed.",
    proDescription:
      "Programmable rent and sale settlement rail with contract state reconciliation and CRE-driven automation telemetry.",
  },
];

const CAPABILITY_STATUS_STYLES: Record<
  (typeof WORKSPACE_CAPABILITIES)[number]["tone"],
  { readonly badge: string; readonly dot: string }
> = {
  ready: {
    badge: "bg-emerald-400/15 text-emerald-200 border border-emerald-300/35",
    dot: "bg-emerald-300",
  },
  protected: {
    badge: "bg-blue-400/15 text-blue-200 border border-blue-300/35",
    dot: "bg-blue-300",
  },
  enabled: {
    badge: "bg-indigo-400/15 text-indigo-200 border border-indigo-300/35",
    dot: "bg-indigo-300",
  },
  live: {
    badge: "bg-cyan-400/15 text-cyan-200 border border-cyan-300/35",
    dot: "bg-cyan-300",
  },
};

export const LandingPage: React.FC = () => {
  const { isLoading, loginWithEmbeddedWallet, connectExternalWallet } =
    useAuth();
  const { mode, setMode } = useUXMode();
  const navigate = useNavigate();
  const isEasy = mode === "easy";

  return (
    <div className="min-h-screen bg-[#060b14]">
      <Background />

      <header className="cyber-nav relative z-10 border-b border-slate-700/60">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="site-header-banner">
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
            <div className="site-header-actions">
              <ModeToggle isEasy={isEasy} setMode={setMode} />
              <button
                type="button"
                onClick={() => navigate("/marketplace")}
                className="hidden sm:inline-flex btn btn-secondary !px-4 !py-2 !text-xs"
              >
                Browse homes
              </button>
            </div>
          </div>
        </div>
      </header>

      <main className="relative z-10">
        <section className="container py-16 md:py-24">
          <div className="mx-auto grid max-w-5xl gap-10">
            <div className="mx-auto max-w-3xl text-center">
              <p className="meta-chip meta-chip-script mb-4">
                Private CRE real-estate rail
              </p>
              <h1 className="mb-5 text-4xl font-bold leading-tight text-slate-50 md:text-5xl">
                List, rent, and sell property with professional-grade privacy.
              </h1>
              <p className="text-panel mx-auto max-w-3xl text-base leading-7 text-slate-200 md:text-lg">
                {isEasy
                  ? "A modern real-estate experience with secure document exchange and clear workflows."
                  : "CRE workflows orchestrate Solidity smart-contract execution for tokenization, privacy-preserving document access, and settlement."}
              </p>

              <div className="mt-8 flex flex-col justify-center gap-3 sm:flex-row">
                {isEasy ? (
                  <>
                    <button
                      type="button"
                      onClick={loginWithEmbeddedWallet}
                      disabled={isLoading}
                      className="btn btn-primary !px-7 !py-3"
                    >
                      {isLoading ? "Connecting…" : "Sign in with Privy"}
                    </button>
                    <button
                      type="button"
                      onClick={connectExternalWallet}
                      className="btn btn-secondary !px-7 !py-3"
                    >
                      Connect external wallet
                    </button>
                    <button
                      type="button"
                      onClick={() => navigate("/marketplace")}
                      className="btn btn-secondary !px-7 !py-3"
                    >
                      Explore marketplace
                    </button>
                  </>
                ) : (
                  <>
                    <button
                      type="button"
                      onClick={connectExternalWallet}
                      className="btn btn-primary !px-8 !py-3"
                    >
                      Connect wallet
                    </button>
                    <button
                      type="button"
                      onClick={loginWithEmbeddedWallet}
                      disabled={isLoading}
                      className="btn btn-secondary !px-8 !py-3"
                    >
                      {isLoading ? "Connecting…" : "Sign in (Privy)"}
                    </button>
                  </>
                )}
              </div>
            </div>

            <div className="glass-card workspace-capabilities-panel mx-auto w-full max-w-4xl">
              <div className="workspace-capabilities-header">
                <ul
                  className="flex flex-wrap gap-2"
                  aria-label="Platform capabilities"
                >
                  {CAPABILITY_CHIPS.map((chip) => (
                    <li key={chip}>
                      <TrustPill>{chip}</TrustPill>
                    </li>
                  ))}
                </ul>
                <p className="text-panel mt-3 text-[15px] leading-7 text-slate-300">
                  {isEasy
                    ? "See which platform services are ready before you launch a listing, rental, or sale flow."
                    : "Each module below shows how CRE automation, smart-contract logic, and private data controls are wired in production."}
                </p>
              </div>

              <div className="workspace-capabilities-grid">
                {WORKSPACE_CAPABILITIES.map((item) => {
                  const styles = CAPABILITY_STATUS_STYLES[item.tone];
                  return (
                    <article
                      key={item.label}
                      className="workspace-capability-card"
                    >
                      <div className="flex items-start gap-3">
                        <span
                          className={`mt-[0.55rem] h-2.5 w-2.5 flex-shrink-0 rounded-full ${styles.dot}`}
                        />
                        <div className="min-w-0 flex-1">
                          <div className="flex flex-wrap items-center justify-between gap-2">
                            <h3 className="text-[15px] font-semibold text-slate-100">
                              {item.label}
                            </h3>
                            <span
                              className={`rounded-full px-2.5 py-1 text-[11px] font-semibold tracking-wide ${styles.badge}`}
                            >
                              {item.status}
                            </span>
                          </div>
                          <p className="mt-2.5 text-sm leading-6 text-slate-300">
                            {isEasy ? item.guidedSentence : item.proSentence}
                          </p>
                        </div>
                      </div>
                    </article>
                  );
                })}
              </div>

              {isEasy && (
                <div className="mt-8 flex flex-col gap-3 sm:flex-row">
                  <button
                    type="button"
                    onClick={loginWithEmbeddedWallet}
                    disabled={isLoading}
                    className="btn btn-primary !w-full !py-2.5 !text-sm sm:!flex-1"
                  >
                    Open workspace
                  </button>
                  <button
                    type="button"
                    onClick={() => navigate("/marketplace")}
                    className="btn btn-secondary !w-full !py-2.5 !text-sm sm:!flex-1"
                  >
                    Browse marketplace
                  </button>
                </div>
              )}
            </div>
          </div>
        </section>

        {isEasy && (
          <section className="container pb-6 md:pb-10">
            <div className="grid grid-cols-2 gap-3 md:grid-cols-4">
              {STATS.map((stat) => (
                <div key={stat.label} className="stat-card">
                  <div className="stat-value">{stat.value}</div>
                  <div className="stat-label">{stat.label}</div>
                </div>
              ))}
            </div>
          </section>
        )}

        <section className="container py-12 md:py-16">
          <div className="mb-8 text-center md:mb-10">
            <h2 className="text-3xl font-bold text-slate-50 md:text-4xl">
              {isEasy
                ? "Built for professional digital real estate"
                : "Protocol capabilities"}
            </h2>
            <p className="text-panel mx-auto mt-3 max-w-3xl text-slate-300">
              {isEasy
                ? "Keep ownership flows private while maintaining strong traceability for buyers, renters, and partners."
                : "Technical building blocks for private CRE workflows across identity, document exchange, and settlement."}
            </p>
          </div>
          <div className="grid gap-4 md:grid-cols-3">
            {FEATURE_CARDS.map((card) => (
              <FeatureCard
                key={card.title}
                icon={card.icon}
                title={card.title}
                description={
                  isEasy ? card.guidedDescription : card.proDescription
                }
              />
            ))}
          </div>
        </section>

        {isEasy && (
          <>
            <section className="container py-10 md:py-16">
              <div className="glass-card how-it-works-panel">
                <div className="mb-8 text-center md:mb-10">
                  <h2 className="text-2xl font-bold text-slate-50 md:text-3xl">
                    How it works
                  </h2>
                  <p className="mx-auto mt-2 max-w-2xl text-slate-300">
                    Clear, low-friction flow for property owners, renters, and
                    buyers.
                  </p>
                </div>
                <div className="grid gap-4 md:grid-cols-3 md:gap-5">
                  {WORKFLOW_STEPS.map((step, index) => (
                    <div key={step.title} className="step-card">
                      <div className="step-number">
                        {String(index + 1).padStart(2, "0")}
                      </div>
                      <h3 className="mb-2 text-lg font-semibold text-slate-50">
                        {step.title}
                      </h3>
                      <p className="text-sm text-slate-300">{step.body}</p>
                    </div>
                  ))}
                </div>
              </div>
            </section>

            <section className="container pb-16 pt-6 text-center md:pb-20">
              <h2 className="text-3xl font-bold text-slate-50 md:text-4xl">
                Need help choosing how to start?
              </h2>
              <p className="mx-auto mt-3 max-w-2xl text-slate-300">
                Open the launch guide for clear onboarding paths, what each
                option does, and where to begin first.
              </p>
              <div className="mt-7 flex flex-col items-center justify-center gap-3 sm:flex-row">
                <Link
                  to="/getting-started"
                  className="btn btn-primary !px-8 !py-3"
                >
                  Open launch guide
                </Link>
              </div>
            </section>
          </>
        )}
      </main>
    </div>
  );
};

const ModeToggle: React.FC<{
  readonly isEasy: boolean;
  readonly setMode: (mode: "easy" | "degen") => void;
}> = ({ isEasy, setMode }) => (
  <div className="mode-toggle">
    <button
      type="button"
      className={`mode-btn ${isEasy ? "active-easy" : ""}`}
      onClick={() => setMode("easy")}
    >
      Guided
    </button>
    <button
      type="button"
      className={`mode-btn ${!isEasy ? "active-degen" : ""}`}
      onClick={() => setMode("degen")}
    >
      Pro
    </button>
  </div>
);

const FeatureCard: React.FC<{
  readonly icon: string;
  readonly title: string;
  readonly description: string;
}> = ({ icon, title, description }) => (
  <div className="feature-card">
    <div className="feature-icon">{icon}</div>
    <h3 className="mb-2 text-lg font-semibold text-slate-50">{title}</h3>
    <p className="text-sm text-slate-300">{description}</p>
  </div>
);

const TrustPill: React.FC<{ readonly children: React.ReactNode }> = ({
  children,
}) => (
  <span
    className="inline-flex items-center rounded-full border border-slate-600/70 bg-slate-900/60 px-3 py-1
    text-xs font-medium text-slate-200"
  >
    {children}
  </span>
);

const Background: React.FC = () => (
  <div className="fixed inset-0 -z-10 overflow-hidden">
    <div className="orb orb-1" />
    <div className="orb orb-2" />
    <div className="orb orb-3" />
    <div className="grid-pattern" />
  </div>
);
