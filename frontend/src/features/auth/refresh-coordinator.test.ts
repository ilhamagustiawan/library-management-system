import { describe, expect, it, vi } from "vitest";

import { RefreshCoordinator } from "./refresh-coordinator";
import type { TokenResult } from "./oauth-client";

describe("RefreshCoordinator", () => {
  it("shares one rotation across concurrent requests using the same token", async () => {
    const rotate = vi.fn(
      async (): Promise<TokenResult> => ({
        status: "error",
        error: { kind: "token-rejected" },
      }),
    );

    const [first, second] = await Promise.all([
      RefreshCoordinator.run("shared-refresh-token", rotate),
      RefreshCoordinator.run("shared-refresh-token", rotate),
    ]);

    expect(first).toEqual({ status: "error", error: { kind: "token-rejected" } });
    expect(second).toEqual(first);
    expect(rotate).toHaveBeenCalledOnce();
  });
});
