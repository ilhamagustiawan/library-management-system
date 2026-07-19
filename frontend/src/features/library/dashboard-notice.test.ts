import { describe, expect, it } from "vitest";

import { DashboardNotice } from "./dashboard-notice";

describe("DashboardNotice", () => {
  it("accepts only the exact borrow success marker", () => {
    expect(DashboardNotice.fromSearchParams({ borrowed: "1" })).toEqual({
      kind: "book-borrowed",
    });
    expect(DashboardNotice.fromSearchParams({ borrowed: ["1"] })).toEqual({ kind: "none" });
    expect(DashboardNotice.fromSearchParams({ borrowed: "true" })).toEqual({ kind: "none" });
  });
});
