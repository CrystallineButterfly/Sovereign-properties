const CLAIM_KEY_HASH_STORAGE_KEY = "rwa_house_last_claim_key_hash";

const isBytes32 = (value: string): boolean => /^0x[a-fA-F0-9]{64}$/.test(value);

export const saveLatestClaimKeyHash = (keyHash: string): boolean => {
  const normalized = String(keyHash || "").trim();
  if (!isBytes32(normalized) || typeof window === "undefined") {
    return false;
  }

  window.localStorage.setItem(CLAIM_KEY_HASH_STORAGE_KEY, normalized);
  return true;
};

export const readLatestClaimKeyHash = (): string | null => {
  if (typeof window === "undefined") {
    return null;
  }

  const value = String(
    window.localStorage.getItem(CLAIM_KEY_HASH_STORAGE_KEY) || "",
  ).trim();
  return isBytes32(value) ? value : null;
};

