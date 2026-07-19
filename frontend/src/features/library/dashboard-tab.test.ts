import { describe, expect, it } from "vitest";

import { DashboardTab } from "./dashboard-tab";

describe("DashboardTab.fromSearchParam", () => {
  it("accepts the history tab", () => {
    expect(DashboardTab.fromSearchParam("history")).toBe("history");
  });

  it("defaults missing, repeated, and unknown values to books", () => {
    expect(DashboardTab.fromSearchParam(undefined)).toBe("books");
    expect(DashboardTab.fromSearchParam(["history", "books"])).toBe("books");
    expect(DashboardTab.fromSearchParam("unknown")).toBe("books");
  });
});
