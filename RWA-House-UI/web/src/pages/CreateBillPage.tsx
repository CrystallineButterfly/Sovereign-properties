import React, { useMemo, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import toast from "react-hot-toast";

import { apiClient } from "@shared/utils/api";
import type { BillType, CreateBillData } from "@shared/types";
import {
  CreateBillSchema,
  sanitizeInput,
} from "@shared/utils/security";
import { useAuth } from "../components/AuthProvider";
import { useUXMode } from "../components/UXModeProvider";

const BILL_TYPES: { value: BillType; label: string }[] = [
  { value: "electricity", label: "Electricity" },
  { value: "water", label: "Water / Sewer" },
  { value: "gas", label: "Gas / Heating" },
  { value: "internet", label: "Internet" },
  { value: "phone", label: "Phone" },
  { value: "property_tax", label: "Property Tax" },
  { value: "insurance", label: "Insurance" },
  { value: "hoa", label: "HOA" },
  { value: "maintenance", label: "Maintenance" },
  { value: "other", label: "Other" },
];

export const CreateBillPage: React.FC = () => {
  const params = useParams();
  const navigate = useNavigate();
  const { walletAddress } = useAuth();
  const { mode } = useUXMode();

  const tokenId = params.tokenId;

  const [billType, setBillType] = useState<BillType>("electricity");
  const [amount, setAmount] = useState("");
  const [provider, setProvider] = useState(walletAddress || "");
  const [dueDate, setDueDate] = useState(() => {
    const d = new Date();
    d.setDate(d.getDate() + 7);
    d.setHours(12, 0, 0, 0);
    return d.toISOString().slice(0, 16);
  });
  const [isRecurring, setIsRecurring] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  const payload: CreateBillData | null = useMemo(() => {
    if (!tokenId) return null;
    const dueIso = new Date(dueDate).toISOString();
    const amt = amount === "" ? NaN : Number(amount);
    return {
      tokenId,
      billType,
      amount: amt,
      dueDate: dueIso,
      provider: sanitizeInput(provider),
      isRecurring,
    };
  }, [tokenId, billType, amount, dueDate, provider, isRecurring]);

  if (!tokenId) {
    return (
      <div className="page-shell page-shell-form">
        <div className="cyber-card p-8">
          <h1 className="form-title">Missing Token ID</h1>
          <p className="form-subtitle">
            This page requires a tokenId route param.
          </p>
          <div className="mt-6">
            <Link to="/dashboard" className="cyber-btn cyber-btn-primary">
              Back
            </Link>
          </div>
        </div>
      </div>
    );
  }

  const submit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!payload) return;
    setFormError(null);

    // Validate with shared schema
    const parsed = CreateBillSchema.safeParse(payload);
    if (!parsed.success) {
      const msg = parsed.error.errors?.[0]?.message || "Invalid bill data";
      setFormError(msg);
      toast.error(msg);
      return;
    }

    try {
      setSubmitting(true);
      const resp = await apiClient.createBill(payload);
      if (resp.success) {
        toast.success("Bill created");
        navigate(`/houses/${tokenId}`);
      } else {
        const message = resp.message || "Failed to create bill";
        setFormError(message);
        toast.error(message);
      }
    } catch (err: any) {
      const message = err?.message || "Failed to create bill";
      setFormError(message);
      toast.error(message);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="page-shell page-shell-form">
      <div className="cyber-card overflow-hidden">
        <div className="form-header">
          <div>
            <h1 className="form-title">
              {mode === "degen" ? "Create Bill" : "Add Property Bill"}
            </h1>
            <p className="form-subtitle">
              {mode === "degen" ? (
                <>
                  Create a bill for Token{" "}
                  <span className="number-pill number-pill-sm number-pill-mono">
                    #{tokenId}
                  </span>
                  . The mediator will record it onchain.
                </>
              ) : (
                <>
                  Create a bill for Property{" "}
                  <span className="number-pill number-pill-sm number-pill-mono">
                    #{tokenId}
                  </span>
                  . The platform records it privately through CRE.
                </>
              )}
            </p>
          </div>
          <div className="flex items-center">
            <span className="meta-chip">
              Token{" "}
              <span className="number-pill number-pill-sm number-pill-mono">
                #{tokenId}
              </span>
            </span>
          </div>
        </div>

        <form
          onSubmit={submit}
          className="mx-auto max-w-3xl p-6 md:p-7 grid grid-cols-1 gap-5"
        >
          {formError && (
            <div
              role="status"
              aria-live="polite"
              className="text-panel border border-rose-400/45 bg-rose-500/10 p-3 text-sm text-rose-200"
            >
              {formError}
            </div>
          )}

          <div>
            <label
              htmlFor="bill-type"
              className="block text-sm text-[var(--text-secondary)] mb-2 font-mono"
            >
              Bill Type
            </label>
            <select
              id="bill-type"
              className="cyber-input"
              value={billType}
              onChange={(e) => setBillType(e.target.value as BillType)}
            >
              {BILL_TYPES.map((t) => (
                <option key={t.value} value={t.value}>
                  {t.label}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label
              htmlFor="bill-amount"
              className="block text-sm text-[var(--text-secondary)] mb-2 font-mono"
            >
              {mode === "degen" ? "Amount (USD)" : "Amount to Charge (USD)"}
            </label>
            <input
              id="bill-amount"
              className="cyber-input font-mono"
              inputMode="decimal"
              placeholder="125.50"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              aria-describedby="bill-amount-help"
            />
            <p id="bill-amount-help" className="form-help mt-2">
              {mode === "degen"
                ? "For demo: amount is an offchain value; payment recording is done via CRE."
                : "This amount is shown to the tenant and recorded through CRE when paid."}
            </p>
          </div>

          <div>
            <label
              htmlFor="bill-due-date"
              className="block text-sm text-[var(--text-secondary)] mb-2 font-mono"
            >
              Due Date
            </label>
            <input
              id="bill-due-date"
              className="cyber-input font-mono"
              type="datetime-local"
              value={dueDate}
              onChange={(e) => setDueDate(e.target.value)}
              required
            />
          </div>

          <div>
            <label
              htmlFor="bill-provider"
              className="block text-sm text-[var(--text-secondary)] mb-2 font-mono"
            >
              {mode === "degen"
                ? "Provider Address"
                : "Provider Wallet Address"}
            </label>
            <input
              id="bill-provider"
              className="cyber-input font-mono"
              placeholder="0x..."
              value={provider}
              onChange={(e) => setProvider(e.target.value)}
              spellCheck={false}
              aria-describedby="bill-provider-help"
            />
            <p id="bill-provider-help" className="form-help mt-2">
              Use a trusted provider address (or your own) depending on your
              policy.
            </p>
          </div>

          <label className="flex items-center gap-3 rounded-lg border border-slate-600/55 bg-slate-900/45 px-3 py-2.5 text-sm text-[var(--text-primary)]">
            <input
              type="checkbox"
              checked={isRecurring}
              onChange={(e) => setIsRecurring(e.target.checked)}
            />
            <span>
              {mode === "degen"
                ? "Recurring bill (mediator may schedule future bills)"
                : "Repeat this bill automatically"}
            </span>
          </label>

          <div className="flex flex-col sm:flex-row gap-3 pt-2">
            <button
              className="cyber-btn cyber-btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
              type="submit"
              disabled={submitting || !payload}
            >
              {submitting
                ? "Submitting..."
                : mode === "degen"
                  ? "Create Bill"
                  : "Add Bill"}
            </button>
            <Link to={`/houses/${tokenId}`} className="cyber-btn text-center">
              Cancel
            </Link>
          </div>
        </form>
      </div>
    </div>
  );
};
