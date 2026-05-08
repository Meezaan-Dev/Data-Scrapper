"use client";

export default function ResourcesError() {
  return (
    <main className="mx-auto flex min-h-screen w-full max-w-3xl flex-col justify-center gap-4 px-5 py-8 text-center">
      <p className="text-sm font-semibold uppercase tracking-normal text-primary">
        Team resources
      </p>
      <h1 className="text-3xl font-bold tracking-normal text-foreground">
        Could not load resources
      </h1>
      <p className="text-muted-foreground">
        Start the scraper API on port 8080, then refresh this page.
      </p>
    </main>
  );
}
