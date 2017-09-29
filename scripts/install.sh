#!/usr/bin/env bash

rm -rf /opt/zanecloud/stone
rm -rf /usr/bin/docker-disk
mkdir -p /opt/zanecloud/stone/bin
cp bin/docker-disk /usr/bin/
cp bin/stone /opt/zanecloud/stone/bin/
cp -r systemd/* /etc/systemd/system/

systemctl restart stone
