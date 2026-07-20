import { describe, expect, it } from "vitest";

import { AuthApi } from "./auth-api";

describe("AuthApi.login", () => {
  it("sends credentials directly to the auth service with cookies enabled", async () => {
    const fetcher: typeof fetch = async (input, init) => {
      expect(String(input)).toBe("http://localhost:8000/api/v1/auth/login");
      expect(init?.credentials).toBe("include");
      expect(init?.body).toBe(JSON.stringify({ email: "maya@perpus-digital.test", password: "quietreading" }));
      return Response.json({
        code: "LMS-200000",
        data: { id: "user-123", name: "Maya Chen", email: "maya@perpus-digital.test" },
      });
    };

    await expect(
      AuthApi.login(
        "http://localhost:8000/api/v1/auth/login",
        { email: "maya@perpus-digital.test", password: "quietreading" },
        fetcher,
      ),
    ).resolves.toEqual({ status: "success" });
  });

  it("returns the auth service recovery message", async () => {
    const fetcher: typeof fetch = async () =>
      Response.json(
        { code: "LMS-401001", message: "invalid email or password" },
        { status: 401 },
      );

    await expect(
      AuthApi.login(
        "http://localhost:8000/api/v1/auth/login",
        { email: "maya@perpus-digital.test", password: "wrong-password" },
        fetcher,
      ),
    ).resolves.toEqual({
      status: "error",
      error: { kind: "rejected", message: "invalid email or password" },
    });
  });
});

describe("AuthApi.register", () => {
  it("creates the account without sending form-only fields", async () => {
    const fetcher: typeof fetch = async (input, init) => {
      expect(String(input)).toBe("http://localhost:8000/api/v1/users");
      expect(init?.credentials).toBe("include");
      expect(init?.body).toBe(
        JSON.stringify({ name: "Maya Chen", email: "maya@perpus-digital.test", password: "quietreading" }),
      );
      return Response.json(
        {
          code: "LMS-200000",
          data: { id: "user-123", name: "Maya Chen", email: "maya@perpus-digital.test" },
        },
        { status: 201 },
      );
    };

    await expect(
      AuthApi.register(
        "http://localhost:8000/api/v1/users",
        {
          name: "Maya Chen",
          email: "maya@perpus-digital.test",
          password: "quietreading",
          confirmPassword: "quietreading",
          acceptsTerms: true,
        },
        fetcher,
      ),
    ).resolves.toEqual({ status: "success" });
  });
});

describe("AuthApi.logout", () => {
  it("terminates the browser auth-service session", async () => {
    const fetcher: typeof fetch = async (input, init) => {
      expect(String(input)).toBe("http://localhost:8000/api/v1/auth/logout");
      expect(init?.method).toBe("POST");
      expect(init?.credentials).toBe("include");
      expect(init?.keepalive).toBe(true);
      return new Response(null, { status: 204 });
    };

    await expect(
      AuthApi.logout("http://localhost:8000/api/v1/auth/logout", fetcher),
    ).resolves.toEqual({ status: "success" });
  });

  it("preserves local session when auth service is unavailable", async () => {
    const fetcher: typeof fetch = async () => {
      throw new Error("unavailable");
    };

    await expect(
      AuthApi.logout("http://localhost:8000/api/v1/auth/logout", fetcher),
    ).resolves.toEqual({
      status: "error",
      message: "Could not reach authentication service. You remain signed in; try again.",
    });
  });
});
