import type { NextRequest } from "next/server";

function hasSameOrigin(request: NextRequest) {
  const origin = request.headers.get("origin");
  if (origin === null) return false;
  try {
    return new URL(origin).origin === request.nextUrl.origin;
  } catch {
    return false;
  }
}

export const AuthRequest = { hasSameOrigin } as const;
