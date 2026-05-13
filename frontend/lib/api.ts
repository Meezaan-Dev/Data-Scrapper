export type Resource = {
  id: number;
  title: string;
  link: string;
  summary: string;
  published_at: string;
  source_name: string;
  tag: string;
  created_at: string;
};

export type ResourceTag = {
  tag: string;
  label: string;
};

export type ScrapeResponse = {
  processed_items: number;
  status: "ok";
};

export const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export async function fetchResources(tag?: string, limit?: number): Promise<Resource[]> {
  const params = new URLSearchParams();
  if (tag) {
    params.set("tag", tag);
  }
  if (limit) {
    params.set("limit", String(limit));
  }

  const query = params.toString();
  const url = `${API_BASE_URL}/api/resources${query ? `?${query}` : ""}`;

  // Server-rendered pages should always show the latest local SQLite contents.
  const res = await fetch(url, { cache: "no-store" });

  if (!res.ok) {
    throw new Error("Failed to fetch resources");
  }

  return res.json();
}

export async function fetchTags(): Promise<ResourceTag[]> {
  const res = await fetch(`${API_BASE_URL}/api/tags`, { cache: "no-store" });

  if (!res.ok) {
    throw new Error("Failed to fetch resource tags");
  }

  return res.json();
}

export async function triggerScrape(): Promise<ScrapeResponse> {
  const res = await fetch(`${API_BASE_URL}/scrape`, {
    method: "POST",
    cache: "no-store",
  });

  if (!res.ok) {
    throw new Error("Failed to trigger scrape");
  }

  return res.json();
}
