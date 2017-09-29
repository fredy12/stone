#!/usr/bin/env bash

rm -rf /opt/zanecloud/stone
rm -rf /usr/bin/docker-disk
mkdir -p /opt/zanecloud/stone/bin
mkdir -p /etc/zanecloud
cp bin/docker-disk /usr/bin/
cp bin/stone /opt/zanecloud/stone/bin/
cp -f systemd/stone.service /etc/systemd/system/
cp -f systemd/stone.conf /etc/zanecloud/

systemctl daemon-reload
systemctl restart stone
