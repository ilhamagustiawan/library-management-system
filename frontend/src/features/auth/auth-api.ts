import { z } from "zod";

import type { LoginInput, RegisterInput } from "./auth-schema";

const successSchema = z.object({
  code: z.literal("LMS-200000"),
  data: z.object({
    id: z.string().min(1),
    name: z.string().min(1),
    email: z.email(),
  }),
});

const errorSchema = z.object({
  code: z.string().min(1),
  message: z.string().min(1),
});

export type LoginResult =
  | { status: "success" }
  | {
      status: "error";
      error: { kind: "rejected" | "unavailable" | "invalid-response"; message: string };
    };

export type RegisterResult = LoginResult;

export type LogoutResult =
  | { status: "success" }
  | { status: "error"; message: string };

async function submit(
  endpoint: string,
  body: object,
  fetcher: typeof fetch,
  unavailableMessage: string,
): Promise<LoginResult> {
  let response: Response;
  try {
    response = await fetcher(endpoint, {
      method: "POST",
      headers: { Accept: "application/json", "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify(body),
    });
  } catch {
    return {
      status: "error",
      error: { kind: "unavailable", message: unavailableMessage },
    };
  }

  try {
    const payload: unknown = await response.json();
    if (response.ok && successSchema.safeParse(payload).success) {
      return { status: "success" };
    }
    const error = errorSchema.safeParse(payload);
    if (!response.ok && error.success) {
      return { status: "error", error: { kind: "rejected", message: error.data.message } };
    }
  } catch {
    // Invalid external responses become a stable client error below.
  }

  return {
    status: "error",
    error: {
      kind: "invalid-response",
      message: "Authentication returned an invalid response. No account data changed; try again.",
    },
  };
}

async function login(
  endpoint: string,
  input: LoginInput,
  fetcher: typeof fetch = fetch,
): Promise<LoginResult> {
  return submit(
    endpoint,
    input,
    fetcher,
    "Authentication service unavailable. Your credentials were not stored; try again.",
  );
}

function register(
  endpoint: string,
  input: RegisterInput,
  fetcher: typeof fetch = fetch,
): Promise<RegisterResult> {
  return submit(
    endpoint,
    { name: input.name, email: input.email, password: input.password },
    fetcher,
    "Registration service unavailable. No account was created; try again.",
  );
}

async function logout(endpoint: string, fetcher: typeof fetch = fetch): Promise<LogoutResult> {
  let response: Response;
  try {
    response = await fetcher(endpoint, {
      method: "POST",
      headers: { Accept: "application/json" },
      credentials: "include",
      keepalive: true,
    });
  } catch {
    return {
      status: "error",
      message: "Could not reach authentication service. You remain signed in; try again.",
    };
  }

  if (response.status === 204 || response.status === 401) {
    return { status: "success" };
  }
  return {
    status: "error",
    message: "Authentication service rejected logout. You remain signed in; try again.",
  };
}

export const AuthApi = { login, logout, register } as const;
