import { render, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { DashboardNoticeToast } from "./dashboard-notice-toast";

const dependencies = vi.hoisted(() => ({ replace: vi.fn(), success: vi.fn() }));

vi.mock("next/navigation", () => ({
  useRouter: () => ({ replace: dependencies.replace }),
}));
vi.mock("sonner", () => ({
  toast: { success: dependencies.success },
}));

describe("DashboardNoticeToast", () => {
  it("announces a borrowed book once and consumes the URL marker", async () => {
    window.history.replaceState(null, "", "/dashboard?borrowed=1&tab=history");
    render(<DashboardNoticeToast notice={{ kind: "book-borrowed" }} />);

    await waitFor(() => {
      expect(dependencies.success).toHaveBeenCalledTimes(1);
      expect(dependencies.replace).toHaveBeenCalledWith("/dashboard?tab=history", { scroll: false });
    });
  });
});
