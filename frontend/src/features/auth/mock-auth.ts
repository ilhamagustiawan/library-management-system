import type { LoginInput, RegisterInput } from "./auth-schema";
import type { AuthSession } from "./auth-session";

export type AuthResult =
  | { status: "success"; session: AuthSession }
  | {
      status: "error";
      error: { kind: "mock-unavailable"; message: string };
    };

function memberNameFromEmail(email: string) {
  const localPart = email.split("@").at(0) ?? "Member";
  const words = localPart.split(/[._-]+/).filter(Boolean);
  const formatted = words
    .map((word) => `${word.charAt(0).toUpperCase()}${word.slice(1)}`)
    .join(" ");

  return formatted || "Member";
}

async function pause() {
  await new Promise((resolve) => window.setTimeout(resolve, 300));
}

// TODO: Replace with a typed /api/auth adapter when backend contracts are available.
export const MockAuth = {
  async login(input: LoginInput): Promise<AuthResult> {
    await pause();
    return {
      status: "success",
      session: {
        id: "mock-member",
        name: memberNameFromEmail(input.email),
        email: input.email,
      },
    };
  },
  async register(input: RegisterInput): Promise<AuthResult> {
    await pause();
    return {
      status: "success",
      session: {
        id: "mock-member",
        name: input.name,
        email: input.email,
      },
    };
  },
} as const;
