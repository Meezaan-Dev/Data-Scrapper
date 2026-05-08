#!/bin/sh
set -eu

mkdir -p /root/data /var/log

echo "0 9 * * 1 cd /root && /root/scraper-api scrape >> /var/log/cron.log 2>&1" | crontab -

crond

exec /root/scraper-api serve
