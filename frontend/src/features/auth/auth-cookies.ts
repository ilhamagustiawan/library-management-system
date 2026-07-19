const flowName = "lms_oauth_flow";
const sessionName = "lms_web_session";
const sessionMaxAgeSeconds = 7 * 24 * 60 * 60;

function options(secure: boolean, maxAge: number) {
  return {
    httpOnly: true,
    secure,
    sameSite: "lax" as const,
    path: "/",
    maxAge,
    priority: "high" as const,
  };
}

export const AuthCookies = { flowName, options, sessionMaxAgeSeconds, sessionName } as const;
