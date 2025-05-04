#!/usr/bin/env bash

# SPDX-FileCopyrightText: 2025 Sidings Media
# SPDX-License-Identifier: MIT

# A little helper script for debugging to automatically set the
# capability for running raw ICMP sockets

go build .
sudo setcap cap_net_raw=+ep icmp_exporter
./icmp_exporter --collector.config.file=config.yaml --log.level debug
