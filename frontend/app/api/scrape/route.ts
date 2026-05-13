import { NextResponse } from "next/server";

import { triggerScrape } from "@/lib/api";

export async function POST() {
  try {
    const payload = await triggerScrape();
    return NextResponse.json(payload);
  } catch {
    return NextResponse.json(
      { error: "Failed to trigger scrape" },
      { status: 502 },
    );
  }
}
