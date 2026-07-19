# Perpus Digital member web

Next.js app for library members. Includes landing, User Service registration, OAuth login, and a guarded dashboard.

## Run locally

Requires Node.js 20.9+ and npm.

```sh
cp .env.example .env.local
npm install
npm run dev
```

The example values match local Docker defaults. Replace both secrets outside local development. Open `http://localhost:3000`.

Login uses Authorization Code with PKCE. The browser sends credentials through Kong to Auth Service. Next.js exchanges the code server-side, loads user info, and stores access and rotating refresh tokens in an AES-GCM encrypted, HttpOnly cookie. The dashboard refreshes access tokens shortly before expiry. No token enters browser-accessible storage.

Registration uses User Service through Kong. User Service creates credentials through Auth Service, persists the member profile, and publishes `UserRegistered.v1` asynchronously.

Members can return active loans from My Books, History, or Book Details. Every
return is confirmed; overdue returns show the server-calculated fine before the
loan changes.

## Checks

```sh
npm test
npm run lint
npm run build
```

## Auth integration

The auth service must use the same client ID, client secret, callback URI, and scopes. Browser login also requires its `ALLOWED_ORIGIN` to match `http://localhost:3000` and its session cookie policy to permit the configured frontend/auth deployment origins.

`AUTH_ISSUER` and `USER_SERVICE_URL` must point to the public gateway origin. Local default: `http://localhost:8000`.
