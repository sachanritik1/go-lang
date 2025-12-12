# Frontend (Next.js + shadcn/ui)

This is the Next.js (App Router) frontend for the Go API in the repo root.

## Setup

Create a local env file:

```bash
cp .env.local.example .env.local
```

Default expected backend:

- Go API: `http://localhost:8080`
- Frontend: `http://localhost:3000`

## Run

```bash
npm install
npm run dev
```

## Notes

- The Go API currently does not set CORS headers. This frontend uses Next.js Route Handlers under `src/app/api/**` as a same-origin proxy.
- Auth uses an HttpOnly cookie `auth_token` set by `POST /api/auth/login`.

## UI

shadcn/ui is initialized via its CLI and components live under `src/components/ui/*`.
