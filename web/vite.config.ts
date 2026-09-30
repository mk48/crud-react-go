import path from "path"
import tailwindcss from "@tailwindcss/vite"
import react from "@vitejs/plugin-react"
import { defineConfig } from "vite"
import { tanstackRouter } from "@tanstack/router-plugin/vite"

// The Go API serves this app from the same origin in production (embedded
// with -tags embedui). In dev, proxy its routes so the app is same-origin
// here too - no CORS, and the API's /config.js supplies the runtime settings.
const apiServer = process.env.KFAMILY_API_URL ?? "http://localhost:8080"

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    // Please make sure that '@tanstack/router-plugin' is passed before '@vitejs/plugin-react'
    tanstackRouter({
      target: "react",
      autoCodeSplitting: true,
    }),
    react(),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  build: {
    // Where the API embeds it from (api/internal/webui/embed_on.go).
    outDir: "../api/internal/webui/dist",
    emptyOutDir: true,
  },
  server: {
    proxy: {
      "/api": apiServer,
      "/config.js": apiServer,
      "/healthcheck": apiServer,
    },
  },
})
