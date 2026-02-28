import React from "react";
import { Link } from "react-router-dom";

export const TermsPage: React.FC = () => {
  return (
    <div className="min-h-screen bg-[#060b14] pb-10">
      <main className="page-shell page-shell-form workspace-surface">
        <div className="cyber-card dashboard-panel text-panel mx-auto max-w-4xl p-6 md:p-8">
          <h1 className="dashboard-title">Terms of Use</h1>
          <p className="dashboard-action-body mt-3 text-sm md:text-base">
            This demo app is for hackathon and testnet use only. Do not treat
            any listing, payment flow, or KYC result as legal or financial
            advice. Always validate all data and transaction details before
            signing.
          </p>
          <p className="dashboard-action-body mt-3 text-sm md:text-base">
            By using this interface, you accept responsibility for wallet
            security, private key management, and all transactions sent from
            your address.
          </p>
          <div className="mt-6">
            <Link to="/" className="btn btn-secondary header-action-btn !text-xs">
              Back to home
            </Link>
          </div>
        </div>
      </main>
    </div>
  );
};
