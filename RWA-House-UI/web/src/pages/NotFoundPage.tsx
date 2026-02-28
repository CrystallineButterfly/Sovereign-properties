import React from "react";
import { Link } from "react-router-dom";

import { BrandMark } from "../components/BrandMark";

export const NotFoundPage: React.FC = () => {
  return (
    <div className="relative flex min-h-screen flex-col overflow-hidden bg-[#060b14]">
      <div className="absolute inset-0 -z-10">
        <div className="orb orb-1" />
        <div className="orb orb-2" />
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
                to="/dashboard"
                className="btn btn-secondary !px-3 !py-2 !text-xs"
              >
                Dashboard
              </Link>
            </div>
          </div>
        </div>
      </header>

      <main className="relative z-10 flex flex-1 items-center justify-center px-6 py-12">
        <div className="glass-card max-w-xl text-center">
          <p className="meta-chip mx-auto mb-4">404</p>
          <h1 className="text-4xl font-bold text-slate-50 md:text-5xl">
            Page not found
          </h1>
          <p className="mt-4 text-slate-300">
            The page you requested was moved, removed, or never existed.
          </p>

          <div className="mt-8 flex flex-col justify-center gap-3 sm:flex-row">
            <Link to="/dashboard" className="btn btn-primary !px-7 !py-3">
              Go to dashboard
            </Link>
            <Link to="/marketplace" className="btn btn-secondary !px-7 !py-3">
              Browse marketplace
            </Link>
          </div>
        </div>
      </main>
    </div>
  );
};
