import React, { useState } from 'react';
import toast from 'react-hot-toast';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { ethers } from 'ethers';
import { useAuth } from './AuthProvider';
import { useUXMode } from './UXModeProvider';
import { apiClient } from '../../../shared/src/utils/api';

interface ListingFormProps {
  tokenId: string;
}

export const ListingForm: React.FC<ListingFormProps> = ({ tokenId }) => {
  const { chainId, walletAddress } = useAuth();
  const { mode } = useUXMode();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const initialListingType: 'for_sale' | 'for_rent' =
    searchParams.get('type') === 'for_rent' ? 'for_rent' : 'for_sale';
  const [listingType, setListingType] = useState<'for_sale' | 'for_rent'>(
    initialListingType,
  );
  const [price, setPrice] = useState('');
  const [isPrivateSale, setIsPrivateSale] = useState(false);
  const [allowedBuyer, setAllowedBuyer] = useState('');
  const [duration, setDuration] = useState('30');
  const [loading, setLoading] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  const parseChainIdValue = (value: string | null): number | null => {
    if (!value) return null;
    const trimmed = value.trim();
    if (!trimmed) return null;
    const parts = trimmed.split(':');
    const parsed = Number.parseInt(parts[parts.length - 1], 10);
    return Number.isFinite(parsed) ? parsed : null;
  };

  const expectedChainId = Number.parseInt(String(import.meta.env.VITE_EXPECTED_CHAIN_ID || ''), 10);
  const connectedChainId = parseChainIdValue(chainId);
  const wrongChain = Number.isFinite(expectedChainId) && connectedChainId !== null && connectedChainId !== expectedChainId;
  const allowedBuyerInvalid =
    isPrivateSale
    && allowedBuyer.trim().length > 0
    && !/^0x[a-fA-F0-9]{40}$/.test(allowedBuyer.trim());

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError(null);

    const rawPrice = price.trim();
    if (!rawPrice) {
      const message = 'Please enter a valid price';
      setFormError(message);
      toast.error(message);
      return;
    }

    let priceInWei: string;
    try {
      priceInWei = mode === 'easy' ? ethers.parseEther(rawPrice).toString() : rawPrice;
    } catch {
      const message = mode === 'easy' ? 'Use a valid ETH amount (example: 1.25)' : 'Use a valid wei amount';
      setFormError(message);
      toast.error(message);
      return;
    }

    if (BigInt(priceInWei) <= 0n) {
      const message = 'Price must be greater than zero';
      setFormError(message);
      toast.error(message);
      return;
    }

    if (isPrivateSale && !allowedBuyer) {
      const message = 'Please enter allowed buyer address for private sale';
      setFormError(message);
      toast.error(message);
      return;
    }

    if (isPrivateSale && !/^0x[a-fA-F0-9]{40}$/.test(allowedBuyer)) {
      const message = 'Allowed buyer must be a valid Ethereum address';
      setFormError(message);
      toast.error(message);
      return;
    }

    try {
      setLoading(true);

      if (!walletAddress) {
        const message = 'Connect your wallet to create a listing';
        setFormError(message);
        toast.error(message);
        return;
      }

      if (wrongChain) {
        if (Number.isFinite(expectedChainId)) {
          const message = `Wrong network. Switch to chain ${expectedChainId} first.`;
          setFormError(message);
          toast.error(message);
        }
        return;
      }

      const response = await apiClient.createListing({
        action: 'create_listing',
        ownerAddress: walletAddress,
        tokenId,
        listingType,
        price: priceInWei,
        preferredToken: ethers.ZeroAddress,
        isPrivateSale,
        allowedBuyer: isPrivateSale ? allowedBuyer : undefined,
        durationDays: Number(duration || '30'),
      });

      if (!response.success) {
        const message = response.message || 'Failed to create listing';
        setFormError(message);
        toast.error(message);
        return;
      }

      toast.success(
        mode === 'degen'
          ? 'Listing created via CRE workflow!'
          : 'Listing published successfully'
      );
      navigate(`/houses/${tokenId}`);
    } catch (error: any) {
      if (error?.shortMessage) {
        setFormError(error.shortMessage);
        toast.error(error.shortMessage);
      } else {
        const message = error.message || 'An error occurred';
        setFormError(message);
        toast.error(message);
      }
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
              {mode === 'degen' ? 'Create Listing' : 'List Property'}
            </h1>
            <p className="form-subtitle">
              {mode === 'degen'
                ? 'Set sale or rental terms for this tokenized asset.'
                : 'Set sale or rental terms for this property.'}
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

          {/* Listing Type */}
          <div className="form-field">
            <label className="block text-sm text-[var(--text-secondary)] mb-2 font-mono" id="listing-type-label">
              {mode === 'degen' ? 'Listing Type' : 'What do you want to do?'}
            </label>
            <div className="grid grid-cols-1 gap-3 sm:grid-cols-2" role="radiogroup" aria-labelledby="listing-type-label">
              <label className={`form-choice-card rounded-xl border px-4 py-3 ${listingType === 'for_sale' ? 'form-choice-card-active border-[#60a5fa] bg-[rgba(96,165,250,0.14)]' : 'border-[rgba(148,163,184,0.45)] bg-[rgba(15,23,42,0.42)]'}`}>
                <input
                  type="radio"
                  value="for_sale"
                  checked={listingType === 'for_sale'}
                  onChange={(e) => setListingType(e.target.value as 'for_sale')}
                  className="mr-2"
                />
                <span className="text-[var(--text-primary)]">{mode === 'degen' ? 'For Sale' : 'Sell this property'}</span>
              </label>
              <label className={`form-choice-card rounded-xl border px-4 py-3 ${listingType === 'for_rent' ? 'form-choice-card-active border-[#60a5fa] bg-[rgba(96,165,250,0.14)]' : 'border-[rgba(148,163,184,0.45)] bg-[rgba(15,23,42,0.42)]'}`}>
                <input
                  type="radio"
                  value="for_rent"
                  checked={listingType === 'for_rent'}
                  onChange={(e) => setListingType(e.target.value as 'for_rent')}
                  className="mr-2"
                />
                <span className="text-[var(--text-primary)]">{mode === 'degen' ? 'For Rent' : 'Rent this property'}</span>
              </label>
            </div>
          </div>

          {/* Price */}
          <div className="form-field">
            <label htmlFor="listing-price" className="block text-sm text-[var(--text-secondary)] mb-2 font-mono">
              {mode === 'degen'
                ? 'Price (in wei)'
                : listingType === 'for_sale'
                  ? 'Sale Price (ETH)'
                  : 'Monthly Rent (ETH)'}
            </label>
            <input
              id="listing-price"
              type={mode === 'degen' ? 'number' : 'text'}
              value={price}
              onChange={(e) => setPrice(e.target.value)}
              placeholder={mode === 'degen' ? '1000000000000000000' : '1.0'}
              className="cyber-input font-mono"
              required
              aria-describedby="listing-price-help"
            />
            <p id="listing-price-help" className="form-help mt-2 font-mono">
              {mode === 'degen'
                ? 'Enter amount in wei (1 ETH = 1000000000000000000 wei)'
                : 'Enter ETH amount. We convert it to wei automatically.'}
            </p>
          </div>

          {/* Duration (for rent) */}
          {listingType === 'for_rent' && (
            <div className="form-field">
              <label htmlFor="listing-duration" className="block text-sm text-[var(--text-secondary)] mb-2 font-mono">
                {mode === 'degen' ? 'Listing Duration (days)' : 'How many days should this stay listed?'}
              </label>
              <input
                id="listing-duration"
                type="number"
                value={duration}
                onChange={(e) => setDuration(e.target.value)}
                min="1"
                max="365"
                className="cyber-input font-mono"
              />
            </div>
          )}

          {/* Private Sale */}
          <div className="form-field">
            <label className="flex items-center rounded-lg border border-slate-600/55 bg-slate-900/45 px-3 py-2.5">
              <input
                type="checkbox"
                checked={isPrivateSale}
                onChange={(e) => setIsPrivateSale(e.target.checked)}
                className="mr-2"
              />
              <span className="text-sm font-medium text-[var(--text-primary)]">Private Sale</span>
            </label>
            <p className="form-help mt-2 font-mono">
              {mode === 'degen' ? 'Only allow a specific buyer to purchase' : 'Restrict this listing to one specific wallet'}
            </p>
          </div>

          {/* Allowed Buyer */}
          {isPrivateSale && (
            <div className="form-field">
              <label htmlFor="listing-allowed-buyer" className="block text-sm text-[var(--text-secondary)] mb-2 font-mono">
                Allowed Buyer Address
              </label>
              <input
                id="listing-allowed-buyer"
                type="text"
                value={allowedBuyer}
                onChange={(e) => setAllowedBuyer(e.target.value)}
                placeholder="0x..."
                className="cyber-input font-mono"
                pattern="^0x[a-fA-F0-9]{40}$"
                aria-invalid={allowedBuyerInvalid}
                aria-describedby={allowedBuyerInvalid ? 'listing-allowed-buyer-error' : undefined}
              />
              {allowedBuyerInvalid && (
                <p id="listing-allowed-buyer-error" className="form-error mt-1">
                  Enter a valid wallet address (0x + 40 hex characters).
                </p>
              )}
            </div>
          )}

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
              {loading
                ? mode === 'degen' ? 'Creating...' : 'Publishing...'
                : mode === 'degen' ? 'Create Listing' : 'Publish Listing'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
