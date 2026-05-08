import { notFound } from "next/navigation";

import { ResourceCard } from "@/components/ResourceCard";
import { ResourceTabs } from "@/components/ResourceTabs";
import { fetchResources } from "@/lib/api";

const labels: Record<string, string> = {
  react: "React",
  nextjs: "Next.js",
  vercel: "Vercel",
  ai: "AI",
  tools: "npm/tools",
};

export default async function TaggedResourcesPage({
  params,
}: {
  params: Promise<{ tag: string }>;
}) {
  const { tag } = await params;

  if (!labels[tag]) {
    notFound();
  }

  const resources = await fetchResources(tag);

  return (
    <main className="mx-auto flex min-h-screen w-full max-w-6xl flex-col gap-8 px-5 py-8 sm:px-8">
      <header className="flex flex-col gap-3">
        <p className="text-sm font-semibold uppercase tracking-normal text-primary">
          Team resources
        </p>
        <div className="flex flex-col gap-2">
          <h1 className="text-4xl font-bold tracking-normal text-foreground sm:text-5xl">
            {labels[tag]}
          </h1>
          <p className="max-w-2xl text-base leading-7 text-muted-foreground">
            Focused updates for {labels[tag]} from the shared feed list.
          </p>
        </div>
      </header>
      <ResourceTabs activeTag={tag} />
      {resources.length === 0 ? (
        <section className="rounded-lg border border-dashed border-border bg-card p-8 text-center text-muted-foreground">
          No {labels[tag]} resources yet. Trigger a scrape from the API, then refresh this page.
        </section>
      ) : (
        <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {resources.map((resource) => (
            <ResourceCard key={resource.id} resource={resource} />
          ))}
        </section>
      )}
    </main>
  );
}
