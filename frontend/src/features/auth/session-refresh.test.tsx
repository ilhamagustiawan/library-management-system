import { render, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { SessionRefresh } from "./session-refresh";

describe("SessionRefresh", () => {
  it("refreshes through POST under an exclusive browser lock", async () => {
    const navigate = vi.fn();
    const fetcher: typeof fetch = async (_input, init) => {
      expect(init?.method).toBe("POST");
      return Response.json({ status: "refreshed" });
    };
    const runExclusive = vi.fn(async (operation: () => Promise<void>) => operation());

    render(
      <SessionRefresh fetcher={fetcher} navigate={navigate} runExclusive={runExclusive} />,
    );

    expect(screen.getByText("Refreshing secure session…")).toBeVisible();
    await waitFor(() => expect(navigate).toHaveBeenCalledWith("/dashboard"));
    expect(runExclusive).toHaveBeenCalledOnce();
  });

  it("returns to login when refresh fails", async () => {
    const navigate = vi.fn();
    const fetcher: typeof fetch = async () => Response.json({ error: "expired" }, { status: 401 });

    render(<SessionRefresh fetcher={fetcher} navigate={navigate} />);

    await waitFor(() => expect(navigate).toHaveBeenCalledWith("/login?error=session_expired"));
  });
});
