import { notFound } from "next/navigation";

import { ResourcesView } from "@/components/ResourcesView";
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
    <ResourcesView
      activeTag={tag}
      emptyLabel={`${labels[tag]} resources`}
      resources={resources}
      subtitle={`Focused updates for ${labels[tag]} from the shared feed list.`}
      title={labels[tag]}
    />
  );
}
