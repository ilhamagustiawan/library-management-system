# Libry member web

Next.js MVP for library members. Includes landing, mock login/register, and a guarded placeholder dashboard.

## Run locally

Requires Node.js 20.9+ and npm.

```sh
npm install
npm run dev
```

Open `http://localhost:3000`. Authentication is intentionally mocked: any valid form submission succeeds. Only the mock member name and email are stored in browser local storage; passwords are never stored.

## Checks

```sh
npm test
npm run lint
npm run build
```

## Backend integration

The future gateway route is `/api/auth/*`. Replace `src/features/auth/mock-auth.ts` after login and registration request/response contracts are defined.
