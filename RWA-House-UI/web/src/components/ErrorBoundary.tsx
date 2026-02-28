import { Component, ErrorInfo, ReactNode } from "react";

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
  errorInfo: ErrorInfo | null;
}

export class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false,
    error: null,
    errorInfo: null,
  };

  public static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error, errorInfo: null };
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error("[ErrorBoundary] Uncaught error:", error, errorInfo);
    this.setState({ errorInfo });
  }

  public render() {
    if (this.state.hasError) {
      return (
        <div
          style={{
            minHeight: "100vh",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            background: "#060b14",
            color: "#fecdd3",
            fontFamily: "Inter, sans-serif",
            padding: "24px",
          }}
        >
          <div style={{ maxWidth: "600px", textAlign: "center" }}>
            <h2 style={{ color: "#fecdd3", marginBottom: "16px" }}>
              ⚠️ Application Error
            </h2>
            <p style={{ color: "#dbeafe", marginBottom: "16px" }}>
              Something went wrong while loading the application.
            </p>
            {this.state.error && (
              <div
                style={{
                  background: "rgba(251, 113, 133, 0.12)",
                  border: "1px solid rgba(251, 113, 133, 0.35)",
                  borderRadius: "8px",
                  padding: "16px",
                  marginBottom: "16px",
                  textAlign: "left",
                }}
              >
                <pre
                  style={{
                    color: "#f8fafc",
                    whiteSpace: "pre-wrap",
                    fontSize: "12px",
                    margin: 0,
                  }}
                >
                  {this.state.error.toString()}
                </pre>
              </div>
            )}
            <button
              onClick={() => window.location.reload()}
              style={{
                padding: "12px 24px",
                background: "linear-gradient(135deg, #60a5fa, #2563eb)",
                color: "#fff",
                border: "none",
                borderRadius: "8px",
                cursor: "pointer",
                fontWeight: "bold",
              }}
            >
              Reload Application
            </button>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}
