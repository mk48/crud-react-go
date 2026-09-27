/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string
  readonly VITE_CASDOOR_ENDPOINT: string
  readonly VITE_CASDOOR_CLIENT_ID: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
