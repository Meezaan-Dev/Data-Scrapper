import { ExternalLink } from "lucide-react";
import Image from "next/image";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import type { Resource } from "@/lib/api";

const formatter = new Intl.DateTimeFormat("en", {
  month: "short",
  day: "numeric",
  year: "numeric",
});

// Feed entries can arrive as RSS, Atom, or GitHub release HTML. The card keeps
// display logic small so the backend can store source content without polishing it.
export function ResourceCard({ resource }: { resource: Resource }) {
  const publishedDate = formatter.format(new Date(resource.published_at));
  const summary = stripHTML(resource.summary);
  const iconUrl = faviconURL(resource.link);

  return (
    <Card className="flex h-full flex-col">
      <CardHeader>
        <div className="flex items-start gap-3">
          <Image
            alt=""
            aria-hidden="true"
            className="mt-0.5 size-9 rounded-md border border-border bg-muted"
            height={36}
            src={iconUrl}
            unoptimized
            width={36}
          />
          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center gap-2 text-sm text-muted-foreground">
              <span>{resource.source_name}</span>
              <span aria-hidden="true">/</span>
              <time dateTime={resource.published_at}>{publishedDate}</time>
              <Badge className="ml-auto">{resource.tag}</Badge>
            </div>
            <CardTitle className="mt-2">{resource.title}</CardTitle>
          </div>
        </div>
      </CardHeader>
      <CardContent className="flex flex-1 flex-col gap-4">
        <p className="line-clamp-4 text-sm leading-6 text-muted-foreground">
          {summary || "No summary available."}
        </p>
        <Button
          className="mt-auto w-fit"
          href={resource.link}
          target="_blank"
          rel="noreferrer"
          aria-label={`Open ${resource.title}`}
        >
          Read source
          <ExternalLink aria-hidden="true" size={16} />
        </Button>
      </CardContent>
    </Card>
  );
}

function stripHTML(value: string) {
  // Most summaries are short HTML fragments; plain text is easier to scan in cards.
  return value.replace(/<[^>]*>/g, " ").replace(/\s+/g, " ").trim();
}

function faviconURL(link: string) {
  try {
    const url = new URL(link);
    // Favicons give each source a quick visual signal without maintaining assets.
    return `https://www.google.com/s2/favicons?domain=${url.hostname}&sz=64`;
  } catch {
    return "https://www.google.com/s2/favicons?domain=localhost&sz=64";
  }
}
