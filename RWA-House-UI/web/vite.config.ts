import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

const zkpProxyTarget = process.env.ZKPASSPORT_PROXY_TARGET || 'http://127.0.0.1:8787'
const apiProxyTarget = process.env.API_PROXY_TARGET || zkpProxyTarget

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
  ],
  resolve: {
    alias: [
      {
        find: 'ethers/utils',
        replacement: path.resolve(__dirname, './src/shims/ethers-utils.ts'),
      },
      {
        find: '@',
        replacement: path.resolve(__dirname, './src'),
      },
      {
        find: '@shared',
        replacement: path.resolve(__dirname, './src/shared'),
      },
      {
        find: 'zod',
        replacement: path.resolve(__dirname, './node_modules/zod'),
      },
      {
        find: 'ethers',
        replacement: path.resolve(__dirname, './node_modules/ethers'),
      },
    ],
  },
  server: {
    port: 3000,
    host: true,
    https: process.env.NODE_ENV === 'production',
    cors: {
      origin: process.env.ALLOWED_ORIGINS?.split(',') || ['http://localhost:3000'],
      credentials: true,
    },
    proxy: {
      '/healthz': {
        target: zkpProxyTarget,
        changeOrigin: true,
      },
      '/kyc/zkpassport': {
        target: zkpProxyTarget,
        changeOrigin: true,
      },
      '/workflow/trigger': {
        target: apiProxyTarget,
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
    rollupOptions: {
      output: {
        manualChunks: {
          'vendor-react': ['react', 'react-dom'],
          'vendor-web3': ['wagmi', 'viem', 'ethers'],
          'vendor-ui': ['@privy-io/react-auth'],
        },
      },
    },
  },
  define: {
    'process.env': {},
  },
  optimizeDeps: {
    exclude: ["@xmtp/wasm-bindings", "@xmtp/browser-sdk"],
    // Ensure commonjs deps get proper interop in dev
    include: ['bn.js', 'buffer', 'process'],
    needsInterop: ['bn.js'],
    esbuildOptions: {
      define: {
        global: 'globalThis',
      },
    },
  },
})
