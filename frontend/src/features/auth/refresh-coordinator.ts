import { createHash } from "node:crypto";

import type { TokenResult } from "./oauth-client";

const replayWindowMilliseconds = 10_000;
const entries = new Map<string, { result: Promise<TokenResult> }>();

function tokenKey(refreshToken: string) {
  return createHash("sha256").update(refreshToken, "utf8").digest("base64url");
}

function run(refreshToken: string, rotate: () => Promise<TokenResult>): Promise<TokenResult> {
  const key = tokenKey(refreshToken);
  const existing = entries.get(key);
  if (existing !== undefined) return existing.result;

  const result = rotate();
  entries.set(key, { result });
  const cleanup = setTimeout(() => {
    if (entries.get(key)?.result === result) entries.delete(key);
  }, replayWindowMilliseconds);
  cleanup.unref();
  return result;
}

export const RefreshCoordinator = { run } as const;
