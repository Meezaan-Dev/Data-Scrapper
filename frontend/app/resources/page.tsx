import { ResourcesView } from "@/components/ResourcesView";
import { fetchResources } from "@/lib/api";

export default async function ResourcesPage() {
  const resources = await fetchResources();

  return (
    <ResourcesView
      emptyLabel="resources"
      resources={resources}
      subtitle="Latest updates across the frontend ecosystem."
      title="Frontend Reads"
    />
  );
}
