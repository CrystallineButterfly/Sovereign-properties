import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import "./index.css";
import { apiClient } from "@shared/utils/api";
import { ErrorBoundary } from "./components/ErrorBoundary";

const root = document.getElementById("root");
if (!root) {
  throw new Error("Root element not found");
}

// Bootstrap marker (helps debug blank screen issues)
root.innerHTML =
  '<div style="min-height:100vh;display:flex;align-items:center;justify-content:center;background:#060b14;color:#dbeafe;font-family:Inter,sans-serif;flex-direction:column;"><div style="font-size:24px;margin-bottom:16px;">🏠</div><div>LOADING PROPMECRE...</div><div style="margin-top:10px;font-size:12px;color:#94a3b8;">Check browser console (F12) for errors</div></div>';

// Configure API base URL (default can be production)
const apiUrl =
  import.meta.env.VITE_API_URL || import.meta.env.VITE_ZKPASSPORT_API_URL;
if (typeof apiUrl === "string" && apiUrl.trim().length > 0) {
  apiClient.setBaseURL(apiUrl.trim());
}

const zkPassportApiUrl = import.meta.env.VITE_ZKPASSPORT_API_URL;
if (
  typeof zkPassportApiUrl === "string" &&
  zkPassportApiUrl.trim().length > 0
) {
  apiClient.setZKPassportBaseURL(zkPassportApiUrl.trim());
}

const houseAddress = import.meta.env.VITE_HOUSE_RWA_ADDRESS;
const rpcUrl = import.meta.env.VITE_RPC_URL;
const maxScanRaw = import.meta.env.VITE_MAX_HOUSE_SCAN;
const parsedMaxScan = Number.parseInt(String(maxScanRaw || ""), 10);

apiClient.setChainConfig({
  rpcURL: typeof rpcUrl === "string" ? rpcUrl : undefined,
  houseRWAAddress: typeof houseAddress === "string" ? houseAddress : undefined,
  maxHouseScan: Number.isFinite(parsedMaxScan) ? parsedMaxScan : undefined,
});

try {
  ReactDOM.createRoot(root).render(
    <React.StrictMode>
      <ErrorBoundary>
        <App />
      </ErrorBoundary>
    </React.StrictMode>,
  );
} catch (err) {
  console.error("[main.tsx] CRASH:", err);
  const message = err instanceof Error ? err.message : String(err);
  const stack = err instanceof Error ? err.stack : "";
  root.innerHTML = `<div style="min-height:100vh;display:flex;align-items:center;justify-content:center;background:#060b14;color:#fecdd3;font-family:Inter,sans-serif;padding:24px;">
    <div style="max-width:600px;">
      <h2 style="margin:0 0 12px 0;">UI CRASHED</h2>
      <pre style="margin:0 0 12px 0;white-space:pre-wrap;color:#f8fafc;">${message}</pre>
      <pre style="margin:0;white-space:pre-wrap;font-size:12px;color:#94a3b8;overflow:auto;max-height:300px;">${stack}</pre>
    </div>
  </div>`;
  throw err;
}
