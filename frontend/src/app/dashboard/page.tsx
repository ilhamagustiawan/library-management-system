import type { Metadata } from "next";

import { DashboardClient } from "@/features/auth/dashboard-client";

export const metadata: Metadata = {
  title: "Member home",
};

export default function DashboardPage() {
  return <DashboardClient />;
}
