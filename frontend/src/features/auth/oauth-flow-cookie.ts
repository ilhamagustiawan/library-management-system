import { z } from "zod";

import type { OAuthFlow } from "./oauth-client";
import { SealedValue } from "./sealed-value";

const maxAgeSeconds = 10 * 60;
const schema: z.ZodType<OAuthFlow> = z.object({
  state: z.string().min(1),
  codeVerifier: z.string().min(1),
  createdAt: z.number().int().positive(),
});

function seal(flow: OAuthFlow, secret: string) {
  return SealedValue.seal(flow, schema, secret);
}

function open(value: string, secret: string, now = Math.floor(Date.now() / 1_000)) {
  const result = SealedValue.open(value, schema, secret);
  if (result.status === "invalid" || now - result.value.createdAt > maxAgeSeconds) {
    return { status: "invalid" } as const;
  }
  return { status: "valid", flow: result.value } as const;
}

export const OAuthFlowCookie = { maxAgeSeconds, open, seal } as const;
