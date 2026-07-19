import { describe, expect, it } from "vitest";

import { OAuthFlowCookie } from "./oauth-flow-cookie";

const secret = "0123456789abcdef0123456789abcdef";

describe("OAuthFlowCookie", () => {
  it("round-trips a fresh authorization flow", () => {
    const flow = { state: "state", codeVerifier: "verifier", createdAt: 1_000 };

    expect(OAuthFlowCookie.open(OAuthFlowCookie.seal(flow, secret), secret, 1_100)).toEqual({
      status: "valid",
      flow,
    });
  });

  it("rejects authorization flows older than ten minutes", () => {
    const flow = { state: "state", codeVerifier: "verifier", createdAt: 1_000 };

    expect(OAuthFlowCookie.open(OAuthFlowCookie.seal(flow, secret), secret, 1_601)).toEqual({
      status: "invalid",
    });
  });
});
