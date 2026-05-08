import { ResourceCard } from "@/components/ResourceCard";
import { ResourceTabs } from "@/components/ResourceTabs";
import { fetchResources } from "@/lib/api";

export default async function ResourcesPage() {
  const resources = await fetchResources();

  return (
    <main className="mx-auto flex min-h-screen w-full max-w-6xl flex-col gap-8 px-5 py-8 sm:px-8">
      <Header title="Frontend Reads" subtitle="Latest updates across the frontend ecosystem." />
      <ResourceTabs />
      <ResourceGrid resources={resources} />
    </main>
  );
}

function Header({ title, subtitle }: { title: string; subtitle: string }) {
  return (
    <header className="flex flex-col gap-3">
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
    </header>
  );
}

function ResourceGrid({
  resources,
}: {
  resources: Awaited<ReturnType<typeof fetchResources>>;
}) {
  if (resources.length === 0) {
    return (
      <section className="rounded-lg border border-dashed border-border bg-card p-8 text-center text-muted-foreground">
        No resources yet. Trigger a scrape from the API, then refresh this page.
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
