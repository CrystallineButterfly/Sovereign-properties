import React from "react";
import { Link } from "react-router-dom";

import { useAuth } from "../components/AuthProvider";
import { BrandMark } from "../components/BrandMark";

const START_OPTIONS: Array<{
  readonly title: string;
  readonly description: string;
}> = [
  {
    title: "Create your secure account",
    description:
      "Sign in with Privy and provision your embedded wallet to unlock protected workflows.",
  },
  {
    title: "Connect documents and ownership",
    description:
      "Upload and encrypt property records so buyers and renters can receive verified access.",
  },
  {
    title: "Run sale, rental, and bill flows",
    description:
      "Manage listing status, payments, and settlement from one operational workspace.",
  },
];

export const GettingStartedPage: React.FC = () => {
  const { isLoading, loginWithEmbeddedWallet, connectExternalWallet } =
    useAuth();

  return (
    <div className="relative min-h-screen overflow-hidden bg-[#060b14]">
      <div className="absolute inset-0 -z-10">
        <div className="orb orb-1" />
        <div className="orb orb-2" />
        <div className="orb orb-3" />
        <div className="grid-pattern" />
      </div>

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
              <Link
                to="/"
                className="btn btn-secondary !px-4 !py-2 !text-xs sm:!text-sm"
              >
                Back to landing
              </Link>
            </div>
          </div>
        </div>
      </header>

      <main className="page-shell page-shell-form relative z-10 !pb-14 !pt-10">
        <section className="glass-card hero-panel">
          <p className="meta-chip">Launch guide</p>
          <h1 className="mt-4 text-3xl font-bold text-slate-50 md:text-4xl">
            Start your private real-estate workflow with confidence
          </h1>
          <p className="text-panel mt-4 max-w-3xl text-base text-slate-300 md:text-lg">
            This page explains exactly what happens when you begin, so your team
            can choose the right first action and move into production-ready
            operations quickly.
          </p>

          <div className="mt-7 flex flex-col gap-3 sm:flex-row">
            <button
              type="button"
              onClick={loginWithEmbeddedWallet}
              disabled={isLoading}
              className="btn btn-primary !px-8 !py-3"
            >
              {isLoading ? "Connecting…" : "Get started"}
            </button>
            <button
              type="button"
              onClick={connectExternalWallet}
              className="btn btn-secondary !px-8 !py-3"
            >
              Connect external wallet
            </button>
            <Link to="/marketplace" className="btn btn-secondary !px-8 !py-3">
              Browse market
            </Link>
          </div>
        </section>

        <section className="mt-7 grid gap-4 md:grid-cols-3">
          {START_OPTIONS.map((option) => (
            <article key={option.title} className="feature-card !p-6">
              <h2 className="mb-3 text-lg font-semibold text-slate-50">
                {option.title}
              </h2>
              <p className="text-sm text-slate-300">{option.description}</p>
            </article>
          ))}
        </section>
      </main>

      <footer className="relative z-10 border-t border-slate-700/50 py-8">
        <div className="container flex flex-col items-center justify-center gap-4 text-center">
          <Link to="/" className="site-brand-link" aria-label="Go to home">
            <BrandMark
              size="xs"
              showWordmark={false}
              className="site-brand-logo site-brand-logo--footer"
            />
          </Link>
          <p className="text-sm text-slate-400">
            Powered by Chainlink CRE workflows and secure wallet authentication.
          </p>
          <div className="flex items-center gap-4 text-sm">
            <Link
              to="/terms"
              className="text-slate-400 transition hover:text-slate-200"
            >
              Terms
            </Link>
            <Link
              to="/privacy"
              className="text-slate-400 transition hover:text-slate-200"
            >
              Privacy
            </Link>
          </div>
        </div>
      </footer>
    </div>
  );
};
