import { describe, expect, it } from "vitest";

import { LoginInput, RegisterInput } from "./auth-schema";

describe("LoginInput", () => {
  it("accepts a valid email and password", () => {
    const result = LoginInput.schema.safeParse({
      email: "maya@perpus-digital.test",
      password: "quietreading",
    });

    expect(result.success).toBe(true);
  });

  it("rejects an invalid email", () => {
    const result = LoginInput.schema.safeParse({
      email: "not-an-email",
      password: "quietreading",
    });

    expect(result.success).toBe(false);
  });

  it("accepts an existing account password shorter than registration policy", () => {
    expect(
      LoginInput.schema.safeParse({ email: "maya@perpus-digital.test", password: "legacy8" }).success,
    ).toBe(true);
  });
});

describe("RegisterInput", () => {
  it("accepts complete matching registration details", () => {
    const result = RegisterInput.schema.safeParse({
      name: "Maya Chen",
      email: "maya@perpus-digital.test",
      password: "quietreading",
      confirmPassword: "quietreading",
      acceptsTerms: true,
    });

    expect(result.success).toBe(true);
  });

  it("rejects mismatched passwords", () => {
    const result = RegisterInput.schema.safeParse({
      name: "Maya Chen",
      email: "maya@perpus-digital.test",
      password: "quietreading",
      confirmPassword: "differentpassword",
      acceptsTerms: true,
    });

    expect(result.success).toBe(false);
  });

  it("requires terms acceptance", () => {
    const result = RegisterInput.schema.safeParse({
      name: "Maya Chen",
      email: "maya@perpus-digital.test",
      password: "quietreading",
      confirmPassword: "quietreading",
      acceptsTerms: false,
    });

    expect(result.success).toBe(false);
  });

  it("requires twelve characters for new passwords", () => {
    const result = RegisterInput.schema.safeParse({
      name: "Maya Chen",
      email: "maya@perpus-digital.test",
      password: "legacy8",
      confirmPassword: "legacy8",
      acceptsTerms: true,
    });

    expect(result.success).toBe(false);
  });
});
