import { describe, expect, it } from "vitest";

import { OAuthCallback } from "./oauth-callback";

const flow = { state: "expected-state", codeVerifier: "verifier", createdAt: 1_000 };

describe("OAuthCallback", () => {
  it("accepts a matching state and authorization code", () => {
    expect(OAuthCallback.validate("code-123", "expected-state", null, flow)).toEqual({
      status: "success",
      code: "code-123",
    });
  });

  it("rejects mismatched state", () => {
    expect(OAuthCallback.validate("code-123", "wrong-state", null, flow)).toEqual({
      status: "error",
      error: "invalid_callback",
    });
  });
});
