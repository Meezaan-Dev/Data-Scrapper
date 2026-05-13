"use client";

import { RefreshCw } from "lucide-react";
import { useRouter } from "next/navigation";
import { useState, useTransition } from "react";

import type { ScrapeResponse } from "@/lib/api";
import { cn } from "@/lib/utils";

export function ScrapeButton() {
  const router = useRouter();
  const [isScraping, setIsScraping] = useState(false);
  const [isPending, startTransition] = useTransition();
  const [message, setMessage] = useState<string | null>(null);

  const disabled = isScraping || isPending;

  async function handleScrape() {
    setIsScraping(true);
    setMessage(null);

    try {
      const res = await fetch("/api/scrape", { method: "POST" });

      if (!res.ok) {
        throw new Error("Scrape request failed");
      }

      const payload = (await res.json()) as ScrapeResponse;
      setMessage(`Updated ${payload.processed_items} items`);
      startTransition(() => {
        router.refresh();
      });
    } catch {
      setMessage("Scrape failed");
    } finally {
      setIsScraping(false);
    }
  }

  return (
    <div className="flex flex-col items-start gap-2 sm:items-end">
      <button
        className={cn(
          "inline-flex h-10 items-center justify-center gap-2 rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-65",
          disabled && "hover:bg-primary",
        )}
        disabled={disabled}
        onClick={handleScrape}
        type="button"
      >
        <RefreshCw
          aria-hidden="true"
          className={cn("size-4", disabled && "animate-spin")}
        />
        {disabled ? "Scraping..." : "Scrape now"}
      </button>
      <p aria-live="polite" className="min-h-5 text-sm text-muted-foreground">
        {message}
      </p>
    </div>
  );
}
