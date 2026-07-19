import { z } from "zod";

import { SealedValue } from "./sealed-value";

const schema = z.object({
  user: z.object({
    id: z.string().min(1),
    name: z.string().min(1),
    email: z.email(),
  }),
  accessToken: z.string().min(1),
  refreshToken: z.string().min(1),
  tokenType: z.literal("Bearer"),
  scope: z.string(),
  expiresAt: z.number().int().positive(),
});

export type WebSession = z.infer<typeof schema>;

export type OpenWebSessionResult =
  | { status: "valid"; session: WebSession }
  | { status: "invalid" };

function seal(session: WebSession, secret: string) {
  return SealedValue.seal(session, schema, secret);
}

function open(value: string, secret: string): OpenWebSessionResult {
  const result = SealedValue.open(value, schema, secret);
  return result.status === "valid"
    ? { status: "valid", session: result.value }
    : { status: "invalid" };
}

function needsRefresh(session: WebSession, now = Math.floor(Date.now() / 1_000)) {
  return session.expiresAt <= now + 30;
}

export const WebSession = { needsRefresh, open, seal } as const;
