import React, { useState } from 'react';
import { apiClient } from '../../../shared/src/utils/api';
import toast from 'react-hot-toast';
import { Link, useNavigate } from 'react-router-dom';
import { useUXMode } from './UXModeProvider';

interface PaymentScreenProps {
  tokenId: string;
}

export const PaymentScreen: React.FC<PaymentScreenProps> = ({ tokenId }) => {
  const { mode } = useUXMode();
  const navigate = useNavigate();
  const [billIndex, setBillIndex] = useState('');
  const [paymentMethod, setPaymentMethod] = useState<'crypto' | 'stripe'>('crypto');
  const [loading, setLoading] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError(null);
    
    if (!billIndex || parseInt(billIndex) < 0) {
      const message = 'Please enter a valid bill index';
      setFormError(message);
      toast.error(message);
      return;
    }

    try {
      setLoading(true);
      
      const response = await apiClient.payBill({
        tokenId: tokenId,
        billIndex: parseInt(billIndex),
        paymentMethod
      });

      if (response.success) {
        toast.success('Payment successful!');
        navigate(`/houses/${tokenId}`);
      } else {
        const message = response.message || 'Payment failed';
        setFormError(message);
        toast.error(message);
      }
    } catch (error: any) {
      const message = error.message || 'An error occurred';
      setFormError(message);
      toast.error(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="page-shell page-shell-form">
      <div className="cyber-card overflow-hidden">
        <div className="form-header">
          <div>
            <h1 className="form-title">
              {mode === 'degen' ? 'Pay Bill' : 'Pay Property Bill'}
            </h1>
            <p className="form-subtitle">
              {mode === 'degen'
                ? 'Submit a bill payment request and let CRE process the settlement.'
                : 'Choose a bill and payment method to complete the charge securely.'}
            </p>
          </div>
          <div className="flex gap-2">
            <span className="meta-chip">
              Token{" "}
              <span className="number-pill number-pill-sm number-pill-mono">
                #{tokenId}
              </span>
            </span>
            <Link to={`/houses/${tokenId}`} className="cyber-btn text-sm">Back</Link>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="mx-auto max-w-3xl p-6 md:p-7 space-y-6">
          {formError && (
            <div role="status" aria-live="polite" className="text-panel border border-rose-400/45 bg-rose-500/10 p-3 text-sm text-rose-200">
              {formError}
            </div>
          )}

          {/* Bill Index */}
          <div className="form-field">
              <label htmlFor="payment-bill-index" className="block text-sm text-[var(--text-secondary)] mb-2 font-mono">
                {mode === 'degen' ? 'Bill Index' : 'Bill Number'}
              </label>
            <input
              id="payment-bill-index"
              type="number"
              value={billIndex}
              onChange={(e) => setBillIndex(e.target.value)}
              placeholder="0"
              min="0"
              className="cyber-input font-mono"
              required
              aria-describedby="payment-bill-index-help"
            />
            <p id="payment-bill-index-help" className="form-help mt-2 font-mono">
              {mode === 'degen'
                ? 'Find the bill index on the Bills tab. The mediator records payment onchain.'
                : 'Open the property Payments tab, find the bill number, then enter it here.'}
            </p>
          </div>

          {/* Payment Method */}
          <div className="form-field">
              <label className="block text-sm text-[var(--text-secondary)] mb-2 font-mono" id="payment-method-label">
                {mode === 'degen' ? 'Payment Method' : 'How do you want to pay?'}
              </label>
            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2" role="radiogroup" aria-labelledby="payment-method-label">
              <label className={`form-choice-card rounded-xl border px-4 py-3 ${paymentMethod === 'crypto' ? 'form-choice-card-active border-[#60a5fa] bg-[rgba(96,165,250,0.14)]' : 'border-[rgba(148,163,184,0.45)] bg-[rgba(15,23,42,0.42)]'}`}>
                <input
                  type="radio"
                  value="crypto"
                  checked={paymentMethod === 'crypto'}
                  onChange={(e) => setPaymentMethod(e.target.value as 'crypto')}
                  className="mr-2"
                />
                <span className="text-[var(--text-primary)]">Cryptocurrency (ETH)</span>
              </label>
              <label className={`form-choice-card rounded-xl border px-4 py-3 ${paymentMethod === 'stripe' ? 'form-choice-card-active border-[#60a5fa] bg-[rgba(96,165,250,0.14)]' : 'border-[rgba(148,163,184,0.45)] bg-[rgba(15,23,42,0.42)]'}`}>
                <input
                  type="radio"
                  value="stripe"
                  checked={paymentMethod === 'stripe'}
                  onChange={(e) => setPaymentMethod(e.target.value as 'stripe')}
                  className="mr-2"
                />
                <span className="text-[var(--text-primary)]">Credit Card (Stripe)</span>
              </label>
            </div>
            <div className="mt-2 form-help font-mono">
              {mode === 'degen'
                ? 'Stripe requires backend endpoints (`/payments/stripe/*`). Crypto uses CRE workflow trigger `pay_bill`.'
                : 'Crypto always works with CRE. Card payments need Stripe endpoints enabled on the backend.'}
            </div>
          </div>

          {/* Submit */}
          <div className="flex space-x-4">
            <button
              type="button"
              onClick={() => navigate(-1)}
              className="cyber-btn flex-1"
              disabled={loading}
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading}
              className="cyber-btn cyber-btn-primary flex-1 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? 'Processing...' : mode === 'degen' ? 'Pay Now' : 'Submit Payment'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
