# KFamily Web

React + TypeScript + Vite + shadcn/ui frontend for the KFamily family-tree
app. Auth is handled via Casdoor (see `src/lib/casdoor.ts` and
`src/lib/auth-context.tsx`); the API lives in `../api`.

No configuration needed: the app talks to the API on its own origin (in dev
`pnpm dev` proxies `/api` and `/config.js` to `http://localhost:8080`, or
`KFAMILY_API_URL`), and reads its Casdoor settings at runtime from the API's
`/config.js` (see `src/lib/runtime-config.ts`).

`pnpm build` writes straight into `../api/internal/webui/dist`, where the API
embeds it (`go build -tags embedui`) - see the root README.

Only one page is wired up so far - `/samples`, a placeholder CRUD resource
for exercising the auth/API plumbing end to end. Replace it once the real
family-tree schema exists.

## Adding components

To add components to your app, run the following command:

```bash
npx shadcn@latest add button
```

This will place the ui components in the `src/components` directory.

## Using components

To use the components in your app, import them as follows:

```tsx
import { Button } from "@/components/ui/button"
```
