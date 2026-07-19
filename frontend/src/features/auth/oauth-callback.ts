import { timingSafeEqual } from "node:crypto";

import type { OAuthFlow } from "./oauth-client";

function statesMatch(actual: string, expected: string) {
  const actualBytes = Buffer.from(actual);
  const expectedBytes = Buffer.from(expected);
  return actualBytes.length === expectedBytes.length && timingSafeEqual(actualBytes, expectedBytes);
}

function validate(code: string | null, state: string | null, error: string | null, flow: OAuthFlow) {
  if (error !== null) {
    return { status: "error", error: "authorization_denied" } as const;
  }
  if (code === null || state === null || !statesMatch(state, flow.state)) {
    return { status: "error", error: "invalid_callback" } as const;
  }
  return { status: "success", code } as const;
}

export const OAuthCallback = { validate } as const;
