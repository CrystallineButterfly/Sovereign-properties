import React from "react";
import {
  BrowserRouter,
  Navigate,
  Route,
  Routes,
  Outlet,
  useParams,
} from "react-router-dom";
import { Toaster } from "react-hot-toast";

import {
  AuthProvider,
  ProtectedRoute,
  useAuth,
} from "./components/AuthProvider";
import { Navigation } from "./components/Navigation";
import { ChainWarning } from "./components/ChainWarning";
import { UXModeProvider } from "./components/UXModeProvider";
import { Dashboard } from "./components/Dashboard";
import { MintHouseForm } from "./components/MintHouseForm";
import { HouseDetails } from "./components/HouseDetails";
import { DocumentsView } from "./components/DocumentsView";
import { ListingForm } from "./components/ListingForm";
import { PaymentScreen } from "./components/PaymentScreen";
import { LandingPage } from "./components/LandingPage";

import { MarketplacePage } from "./pages/MarketplacePage";
import { ClaimKeyPage } from "./pages/ClaimKeyPage";
import { CreateBillPage } from "./pages/CreateBillPage";
import { NotFoundPage } from "./pages/NotFoundPage";
import { ProfilePage } from "./pages/ProfilePage";
import { GettingStartedPage } from "./pages/GettingStartedPage";
import { TermsPage } from "./pages/TermsPage";
import { PrivacyPage } from "./pages/PrivacyPage";

const AuthedLayout: React.FC = () => {
  return (
    <div className="min-h-screen bg-[#060b14]">
      <a href="#app-main-content" className="skip-link">
        Skip to main content
      </a>
      <Navigation />
      <ChainWarning />
      <div id="app-main-content" tabIndex={-1}>
        <Outlet />
      </div>
    </div>
  );
};

const HouseDetailsRoute: React.FC = () => {
  const params = useParams();
  const tokenId = params.tokenId;
  if (!tokenId) return <NotFoundPage />;
  return <HouseDetails tokenId={tokenId} />;
};

const HouseDocumentsRoute: React.FC = () => {
  const params = useParams();
  const tokenId = params.tokenId;
  if (!tokenId) return <NotFoundPage />;
  return <DocumentsView tokenId={tokenId} />;
};

const HouseListingRoute: React.FC = () => {
  const params = useParams();
  const tokenId = params.tokenId;
  if (!tokenId) return <NotFoundPage />;
  return <ListingForm tokenId={tokenId} />;
};

const HousePaymentRoute: React.FC = () => {
  const params = useParams();
  const tokenId = params.tokenId;
  if (!tokenId) return <NotFoundPage />;
  return <PaymentScreen tokenId={tokenId} />;
};

// Root route component that shows landing page or redirects to dashboard
const RootRoute: React.FC = () => {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="loading-container bg-[#060b14]">
        <div className="spinner" />
      </div>
    );
  }

  return isAuthenticated ? (
    <Navigate to="/dashboard" replace />
  ) : (
    <LandingPage />
  );
};

function App() {
  return (
    <AuthProvider>
      <UXModeProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<RootRoute />} />
            <Route path="/getting-started" element={<GettingStartedPage />} />
            <Route path="/terms" element={<TermsPage />} />
            <Route path="/privacy" element={<PrivacyPage />} />

            <Route
              element={
                <ProtectedRoute>
                  <AuthedLayout />
                </ProtectedRoute>
              }
            >
              <Route path="/dashboard" element={<Dashboard />} />
              <Route path="/profile" element={<ProfilePage />} />
              <Route path="/mint" element={<MintHouseForm />} />
              <Route path="/marketplace" element={<MarketplacePage />} />
              <Route path="/claim" element={<ClaimKeyPage />} />

              <Route path="/houses/:tokenId" element={<HouseDetailsRoute />} />
              <Route
                path="/houses/:tokenId/documents"
                element={<HouseDocumentsRoute />}
              />
              <Route
                path="/houses/:tokenId/list"
                element={<HouseListingRoute />}
              />
              <Route
                path="/houses/:tokenId/pay"
                element={<HousePaymentRoute />}
              />
              <Route
                path="/houses/:tokenId/bills/create"
                element={<CreateBillPage />}
              />
            </Route>

            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </BrowserRouter>

        <Toaster
          position="top-right"
          toastOptions={{
            duration: 4500,
            style: {
              background: "rgba(10, 16, 27, 0.92)",
              color: "#e5efff",
              border: "1px solid rgba(148, 163, 184, 0.35)",
              backdropFilter: "blur(12px)",
            },
          }}
        />
      </UXModeProvider>
    </AuthProvider>
  );
}

export default App;
