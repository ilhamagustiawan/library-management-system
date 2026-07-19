function isAuthorizeReturnTo(raw: string, issuer: string) {
  try {
    const expected = new URL(issuer);
    const target = new URL(raw);
    return (
      target.origin === expected.origin &&
      target.pathname === "/oauth/authorize" &&
      target.username === "" &&
      target.password === "" &&
      target.hash === ""
    );
  } catch {
    return false;
  }
}

export const OAuthNavigation = { isAuthorizeReturnTo } as const;
