# Frontend RSS Hub Viewer

Next.js viewer for the scraper API.

## Development

```bash
npm install
npm run dev
```

The app reads from `NEXT_PUBLIC_API_URL`, defaulting to `http://localhost:8080`.

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080 npm run dev
```

Routes:

- `/resources`
- `/resources/frontend`
- `/resources/ai-tools`
- `/resources/llm`
- `/resources/aws`
- `/resources/mcp`
- `/resources/tools`

Tabs are loaded from the API's `/api/tags` endpoint so labels and routes stay aligned with `scraper-api/config.json`.
