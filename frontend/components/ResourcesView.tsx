import { ResourceCard } from "@/components/ResourceCard";
import { ResourceTabs } from "@/components/ResourceTabs";
import { ScrapeButton } from "@/components/ScrapeButton";
import type { Resource, ResourceTag } from "@/lib/api";

type ResourcesViewProps = {
  activeTag?: string;
  emptyLabel?: string;
  resources: Resource[];
  subtitle: string;
  tags: ResourceTag[];
  title: string;
};

// Shared layout for the all-resources page and the tag-specific pages.
export function ResourcesView({
  activeTag,
  emptyLabel = "resources",
  resources,
  subtitle,
  tags,
  title,
}: ResourcesViewProps) {
  return (
    <main className="mx-auto flex min-h-screen w-full max-w-6xl flex-col gap-8 px-5 py-8 sm:px-8">
      <PageHeader title={title} subtitle={subtitle} />
      <ResourceTabs activeTag={activeTag} tags={tags} />
      <ResourceGrid emptyLabel={emptyLabel} resources={resources} />
    </main>
  );
}

function PageHeader({ title, subtitle }: { title: string; subtitle: string }) {
  return (
    <header className="flex flex-col gap-5 sm:flex-row sm:items-start sm:justify-between">
      <div className="flex flex-col gap-3">
        <p className="text-sm font-semibold uppercase tracking-normal text-primary">
          Team resources
        </p>
        <div className="flex flex-col gap-2">
          <h1 className="text-4xl font-bold tracking-normal text-foreground sm:text-5xl">
            {title}
          </h1>
          <p className="max-w-2xl text-base leading-7 text-muted-foreground">
            {subtitle}
          </p>
        </div>
      </div>
      <ScrapeButton />
    </header>
  );
}

function ResourceGrid({
  emptyLabel,
  resources,
}: {
  emptyLabel: string;
  resources: Resource[];
}) {
  if (resources.length === 0) {
    return (
      <section className="rounded-lg border border-dashed border-border bg-card p-8 text-center text-muted-foreground">
        No {emptyLabel} yet. Use Scrape now to load the latest resources.
      </section>
    );
  }

  return (
    <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
      {resources.map((resource) => (
        <ResourceCard key={resource.id} resource={resource} />
      ))}
    </section>
  );
}
