#!/bin/sh
set -eu

mkdir -p /root/data /var/log

if [ "$#" -gt 0 ]; then
  exec /root/scraper-api "$@"
fi

echo "0 9 * * 1 cd /root && /root/scraper-api scrape >> /var/log/cron.log 2>&1" | crontab -

crond

if [ "${INITIAL_SCRAPE:-true}" = "true" ]; then
  count="$(sqlite3 /root/data/resources.db 'SELECT COUNT(*) FROM resources;' 2>/dev/null || echo 0)"
  if [ "$count" = "0" ]; then
    echo "Running initial scrape..."
    /root/scraper-api scrape || echo "Initial scrape failed; API will still start."
  fi
fi

exec /root/scraper-api serve
