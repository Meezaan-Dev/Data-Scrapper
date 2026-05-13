import { notFound } from "next/navigation";

import { ResourcesView } from "@/components/ResourcesView";
import { fetchResources, fetchTags } from "@/lib/api";

export default async function TaggedResourcesPage({
  params,
}: {
  params: Promise<{ tag: string }>;
}) {
  const { tag } = await params;
  const tagsPromise = fetchTags();
  const resourcesPromise = fetchResources(tag);
  const tags = await tagsPromise;
  const activeTag = tags.find((candidate) => candidate.tag === tag);

  if (!activeTag) {
    notFound();
  }

  const resources = await resourcesPromise;

  return (
    <ResourcesView
      activeTag={tag}
      emptyLabel={`${activeTag.label} resources`}
      resources={resources}
      subtitle={`Focused updates and references for ${activeTag.label} from the shared resource list.`}
      tags={tags}
      title={activeTag.label}
    />
  );
}
