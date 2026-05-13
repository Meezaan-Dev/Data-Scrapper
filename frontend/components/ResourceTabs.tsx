import Link from "next/link";

import type { ResourceTag } from "@/lib/api";
import { cn } from "@/lib/utils";

export function ResourceTabs({
  activeTag,
  tags,
}: {
  activeTag?: string;
  tags: ResourceTag[];
}) {
  const tabs = [
    { label: "All", href: "/resources", tag: undefined },
    ...tags.map((tag) => ({
      label: tag.label,
      href: `/resources/${tag.tag}`,
      tag: tag.tag,
    })),
  ];

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
