import Link from "next/link";

import { cn } from "@/lib/utils";

const tabs = [
  { label: "All", href: "/resources", tag: undefined },
  { label: "React", href: "/resources/react", tag: "react" },
  { label: "Next.js", href: "/resources/nextjs", tag: "nextjs" },
  { label: "Vercel", href: "/resources/vercel", tag: "vercel" },
  { label: "AI", href: "/resources/ai", tag: "ai" },
  { label: "npm/tools", href: "/resources/tools", tag: "tools" },
];

export function ResourceTabs({ activeTag }: { activeTag?: string }) {
  return (
    <nav
      aria-label="Resource filters"
      className="flex w-full gap-2 overflow-x-auto border-b border-border pb-3"
    >
      {tabs.map((tab) => {
        const active = tab.tag === activeTag || (!tab.tag && !activeTag);

        return (
          <Link
            key={tab.href}
            href={tab.href}
            className={cn(
              "inline-flex h-9 shrink-0 items-center rounded-md border px-3 text-sm font-medium transition-colors",
              active
                ? "border-primary bg-primary text-primary-foreground"
                : "border-border bg-card text-muted-foreground hover:text-foreground",
            )}
          >
            {tab.label}
          </Link>
        );
      })}
    </nav>
  );
}
