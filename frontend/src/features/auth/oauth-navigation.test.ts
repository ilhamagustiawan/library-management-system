import { describe, expect, it } from "vitest";

import { OAuthNavigation } from "./oauth-navigation";

describe("OAuthNavigation", () => {
  it("accepts only the configured authorization endpoint as login return target", () => {
    expect(
      OAuthNavigation.isAuthorizeReturnTo(
        "http://localhost:8000/oauth/authorize?client_id=member-nextjs-web",
        "http://localhost:8000",
      ),
    ).toBe(true);
    expect(
      OAuthNavigation.isAuthorizeReturnTo("https://attacker.test/oauth/authorize", "http://localhost:8000"),
    ).toBe(false);
  });
});
