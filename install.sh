#!/bin/bash
go build
mkdir /usr/local/ddisk_exporter/ -p
mv ddisk-exporter /usr/local/ddisk_exporter/
rm -f /lib/systemd/system/ddisk-exporter.service
cp ddisk-exporter.service /lib/systemd/system/
systemctl daemon-reload
systemctl enable ddisk-exporter
systemctl restart ddisk-exporter
systemctl status ddisk-exporter
