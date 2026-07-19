import { describe, expect, it } from "vitest";

import { WebSession } from "./web-session";

const secret = "0123456789abcdef0123456789abcdef";
const session: WebSession = {
  user: { id: "user-123", name: "Maya Chen", email: "maya@libry.test" },
  accessToken: "access-token",
  refreshToken: "refresh-token",
  tokenType: "Bearer",
  scope: "library:read library:write",
  expiresAt: 2_000,
};

describe("WebSession", () => {
  it("encrypts and authenticates the server session", () => {
    const sealed = WebSession.seal(session, secret);

    expect(WebSession.open(sealed, secret)).toEqual({ status: "valid", session });
    expect(sealed).not.toContain("access-token");
    expect(sealed).not.toContain("refresh-token");
  });

  it("rejects a tampered session without throwing", () => {
    const sealed = WebSession.seal(session, secret);
    const tampered = sealed.replace("v1.", "v1.A");

    expect(WebSession.open(tampered, secret)).toEqual({ status: "invalid" });
  });

  it("refreshes shortly before access-token expiry", () => {
    expect(WebSession.needsRefresh(session, 1_969)).toBe(false);
    expect(WebSession.needsRefresh(session, 1_970)).toBe(true);
  });
});
