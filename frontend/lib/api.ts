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

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

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

  const res = await fetch(url, { cache: "no-store" });

  if (!res.ok) {
    throw new Error("Failed to fetch resources");
  }

  return res.json();
}
