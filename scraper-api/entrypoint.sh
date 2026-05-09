#!/bin/sh
set -eu

mkdir -p /root/data /var/log

# Allow the container to run one-off commands, for example:
# docker compose run --rm api scrape
if [ "$#" -gt 0 ]; then
  exec /root/scraper-api "$@"
fi

# BusyBox crond is included with Alpine. It runs the same CLI scrape mode used
# by the manual API endpoint, so there is only one scraping implementation.
echo "0 9 * * 1 cd /root && /root/scraper-api scrape >> /var/log/cron.log 2>&1" | crontab -

crond

# First boot should show a useful UI, so seed the database when the volume is empty.
if [ "${INITIAL_SCRAPE:-true}" = "true" ]; then
  count="$(sqlite3 /root/data/resources.db 'SELECT COUNT(*) FROM resources;' 2>/dev/null || echo 0)"
  if [ "$count" = "0" ]; then
    echo "Running initial scrape..."
    /root/scraper-api scrape || echo "Initial scrape failed; API will still start."
  fi
fi

exec /root/scraper-api serve
