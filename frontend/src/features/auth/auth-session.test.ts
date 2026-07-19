import { describe, expect, it } from "vitest";

import { AuthSession } from "./auth-session";

function createMemoryStorage(initialValue: string | null = null) {
  let value = initialValue;

  return {
    getItem: () => value,
    removeItem: () => {
      value = null;
    },
    setItem: (_key: string, nextValue: string) => {
      value = nextValue;
    },
  };
}

describe("AuthSession", () => {
  it("round-trips a valid mock session", () => {
    const storage = createMemoryStorage();
    const session = {
      id: "mock-member",
      name: "Maya Chen",
      email: "maya@libry.test",
    };

    AuthSession.write(storage, session);

    expect(AuthSession.read(storage)).toEqual(session);
  });

  it("treats malformed stored data as no session", () => {
    const storage = createMemoryStorage('{"name":42}');

    expect(AuthSession.read(storage)).toBeNull();
  });

  it("clears a stored session", () => {
    const storage = createMemoryStorage("stored");

    AuthSession.clear(storage);

    expect(AuthSession.read(storage)).toBeNull();
  });
});
