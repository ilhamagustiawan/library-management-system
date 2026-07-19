export type DashboardTab = "books" | "history";

function fromSearchParam(value: string | string[] | undefined): DashboardTab {
  return value === "history" ? "history" : "books";
}

export const DashboardTab = { fromSearchParam } as const;
