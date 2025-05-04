<!-- 
SPDX-FileCopyrightText: 2022-2025 Sidings Media <contact@sidingsmedia.com>
SPDX-License-Identifier: MIT
-->

# ICMP Exporter

A prometheus ICMP exporter specifically designed to be able to measure
the latency across point to point links. 

## Installation and Usage

The `icmp_exporter` listens on HTTP port 10035 by default. See the
--help output for more options. You must provide a configuration file
using the `--collector.config.file` option to start the exporter.

### Configuration

Configuration of ping targets is done using a simple `yaml` file. No
configuration option is required (although the exporter won't be very
useful if you don't declare any targets).

To allow the exporter to create raw ICMP sockets, you must set the
required capability on the exporter binary as below. Alternatively, you
may run the exporter as root (NOT RECOMMENDED).

```bash
setcap cap_net_raw=+ep /path/to/icmp_exporter
```

##### General Configuration Options:
| Option | Description | Required | Default |
|---|---|---|---|
| `default_ttl` | The TTL to use for a target if none is specified. | :x: | 64 |
| `default_size` | The default size of packets sent if the target does not specify a value. The minimum value is 24. | :x: | 24 |
| `default_count` | The default number of packets to send if the target does not specify a value. | :x: | 1 |
| `default_interval` | The default interval between sending packets in milliseconds if the target does not specify a value. | :x: | 1000 |
| `timeout` | Timeout in milliseconds for ICMP requests after which they are assumed lost. | :x: | 5000 |
| `targets` | A list of targets to ping | :x: | |

##### Target Configuration Options:

| Option | Description | Required | Default | 
|---|---|---|---|
| `host` | IP address or hostname of target. | :white_check_mark: | |
| `interface` | Interface to ping target from. | :x: | |
| `ttl`| TTL for ICMP packet. | :x: | Value of `default_ttl` |
| `size` | Size of ICMP packets to send. The minimum value is 24. | :x: | Value of `default_size` |
| `count` | Number of packets to send. | :x: | Value of `default_count` |
| `interval` | Interval between sending packets in milliseconds. | :x: | Value of `default_interval` | 

#### Examples

##### Minimum configuration

```yaml
targets:
  - host: 8.8.8.8
  - host: 192.168.150.150
```

##### All options
```yaml
default_ttl: 1
default_size: 56
default_count: 2
default_interval: 500

targets:
  - host: 8.8.8.8
    ttl: 64
    count: 5
    size: 24
    interval: 1000

  - host: 192.168.150.150
    interface: eth0
```

## Licence
This repo uses the [REUSE](https://reuse.software) standard in order to
communicate the correct licence for the file. For those unfamiliar with
the standard the licence for each file can be found in one of three
places. The licence will either be in a comment block at the top of the
file, in a `.license` file with the same name as the file, or in the
dep5 file located in the `.reuse` directory. If you are unsure of the
licencing terms please contact
[contact@sidingsmedia.com](mailto:contact@sidingsmedia.com?subject=ICMP%20Exporter%20Licence).
All files committed to this repo must contain valid licencing
information or the pull request can not be accepted.
