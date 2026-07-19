export type DashboardNotice =
  | { kind: "none" }
  | { kind: "book-borrowed" };

type DashboardSearchParams = { borrowed?: string | string[] };

function fromSearchParams(searchParams: DashboardSearchParams): DashboardNotice {
  return searchParams.borrowed === "1" ? { kind: "book-borrowed" } : { kind: "none" };
}

export const DashboardNotice = { fromSearchParams } as const;
