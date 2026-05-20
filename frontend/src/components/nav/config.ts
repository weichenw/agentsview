export const SCHEDULER_ROUTE = "scheduler";

export interface NavItem {
  id: string;
  label: string;
  icon: string;
  route: string;
}

export const navItems: NavItem[] = [
  { id: "sessions", label: "Sessions", icon: "grid", route: "/" },
  { id: "scheduler", label: "Scheduler", icon: "clock", route: "/scheduler" },
  { id: "usage", label: "Usage", icon: "chart", route: "/usage" },
  { id: "trends", label: "Trends", icon: "trend", route: "/trends" },
  { id: "insights", label: "Insights", icon: "insight", route: "/insights" },
  { id: "pinned", label: "Pinned", icon: "pin", route: "/pinned" },
  { id: "trash", label: "Trash", icon: "trash", route: "/trash" },
  { id: "settings", label: "Settings", icon: "settings", route: "/settings" },
];
