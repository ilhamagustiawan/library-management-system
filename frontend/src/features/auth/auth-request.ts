import type { NextRequest } from "next/server";

function hasOrigin(request: NextRequest, expectedURL: string) {
  const origin = request.headers.get("origin");
  if (origin === null) return false;
  try {
    return new URL(origin).origin === new URL(expectedURL).origin;
  } catch {
    return false;
  }
}

function hasHost(request: NextRequest, expectedURL: string) {
  const forwardedHost = request.headers.get("x-forwarded-host")?.split(",", 1)[0]?.trim();
  const requestHost = forwardedHost || request.headers.get("host") || request.nextUrl.host;
  try {
    return requestHost === new URL(expectedURL).host;
  } catch {
    return false;
  }
}

export const AuthRequest = { hasHost, hasOrigin } as const;
