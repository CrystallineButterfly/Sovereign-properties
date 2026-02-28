import React from "react";
import { Link } from "react-router-dom";

export const PrivacyPage: React.FC = () => {
  return (
    <div className="min-h-screen bg-[#060b14] pb-10">
      <main className="page-shell page-shell-form workspace-surface">
        <div className="cyber-card dashboard-panel text-panel mx-auto max-w-4xl p-6 md:p-8">
          <h1 className="dashboard-title">Privacy Notice</h1>
          <p className="dashboard-action-body mt-3 text-sm md:text-base">
            This app stores local session preferences in your browser and sends
            workflow payloads to configured backend services for CRE execution.
            Wallet addresses and transaction metadata may be visible on public
            testnet explorers.
          </p>
          <p className="dashboard-action-body mt-3 text-sm md:text-base">
            For production use, run this stack with your own privacy policy,
            encrypted storage controls, and compliance review.
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
