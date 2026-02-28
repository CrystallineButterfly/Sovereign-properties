import React, { createContext, useContext, useEffect, useMemo, useState } from 'react';

export type UXMode = 'easy' | 'degen';

interface UXModeContextValue {
  mode: UXMode;
  setMode: (mode: UXMode) => void;
  toggleMode: () => void;
}

const STORAGE_KEY = 'rwa-house-ux-mode';

const UXModeContext = createContext<UXModeContextValue | null>(null);

export const UXModeProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [mode, setMode] = useState<UXMode>(() => {
    const persisted = typeof window !== 'undefined' ? window.localStorage.getItem(STORAGE_KEY) : null;
    return persisted === 'degen' ? 'degen' : 'easy';
  });

  useEffect(() => {
    if (typeof window !== 'undefined') {
      window.localStorage.setItem(STORAGE_KEY, mode);
    }
  }, [mode]);

  const value = useMemo<UXModeContextValue>(
    () => ({
      mode,
      setMode,
      toggleMode: () => setMode((prev) => (prev === 'easy' ? 'degen' : 'easy')),
    }),
    [mode]
  );

  return <UXModeContext.Provider value={value}>{children}</UXModeContext.Provider>;
};

export const useUXMode = (): UXModeContextValue => {
  const context = useContext(UXModeContext);
  if (!context) {
    throw new Error('useUXMode must be used within UXModeProvider');
  }
  return context;
};
