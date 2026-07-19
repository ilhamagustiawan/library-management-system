# Libry member web

Next.js MVP for library members. Includes landing, auth-service registration, OAuth login, and a guarded placeholder dashboard.

## Run locally

Requires Node.js 20.9+ and npm.

```sh
cp .env.example .env.local
npm install
npm run dev
```

The example values match local Docker defaults. Replace both secrets outside local development. Open `http://localhost:3000`.

Login uses Authorization Code with PKCE. The browser sends credentials directly to the auth service. Next.js exchanges the code server-side, loads user info, and stores access and rotating refresh tokens in an AES-GCM encrypted, HttpOnly cookie. The dashboard refreshes access tokens shortly before expiry. No token enters browser-accessible storage.

Registration creates an auth-service account, then starts the OAuth login flow.

## Checks

```sh
npm test
npm run lint
npm run build
```

## Auth integration

The auth service must use the same client ID, client secret, callback URI, and scopes. Browser login also requires its `ALLOWED_ORIGIN` to match `http://localhost:3000` and its session cookie policy to permit the configured frontend/auth deployment origins.
