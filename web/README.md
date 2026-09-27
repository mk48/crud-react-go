# KFamily Web

React + TypeScript + Vite + shadcn/ui frontend for the KFamily family-tree
app. Auth is handled via Casdoor (see `src/lib/casdoor.ts` and
`src/lib/auth-context.tsx`); the API lives in `../api`.

Copy `.env.example` to `.env` and fill in:

| Variable | Purpose |
|---|---|
| `VITE_API_BASE_URL` | Base URL of the API server |
| `VITE_CASDOOR_ENDPOINT` | Casdoor server URL |
| `VITE_CASDOOR_CLIENT_ID` | Casdoor application client ID (public - the secret stays in the API's env) |

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
