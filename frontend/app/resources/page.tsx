import { ResourcesView } from "@/components/ResourcesView";
import { fetchResources, fetchTags } from "@/lib/api";

export default async function ResourcesPage() {
  const [resources, tags] = await Promise.all([fetchResources(), fetchTags()]);

  return (
    <ResourcesView
      emptyLabel="resources"
      resources={resources}
      subtitle="Latest official updates and curated references across frontend engineering, AI tooling, LLMs, AWS, and MCP."
      tags={tags}
      title="Engineering Resources"
    />
  );
}
