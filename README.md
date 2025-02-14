
<p align="center">
  <h1 align="center">Melody</h1>
  <p align="center">Monitor the Internet's background noise</p>
</p>
<p align="center">
    <a href="https://goreportcard.com/badge/github.com/ma111e/melody"><img src="https://goreportcard.com/badge/github.com/ma111e/melody" alt="Go Report Card"/></a>
    <a href="https://coveralls.io/github/ma111e/melody"><img src="https://coveralls.io/repos/github/ma111e/melody/badge.svg" alt="Coverage Status"/></a>
    <a href="https://hub.docker.com/r/ma111e/melody/builds"><img src="https://img.shields.io/docker/cloud/build/ma111e/melody" alt="Docker build status"/></a>
    <a href="https://hub.docker.com/r/ma111e/melody/builds"><img src="https://img.shields.io/docker/image-size/ma111e/melody?sort=date" alt="Docker image size"/></a>
</p>

<p align="center">
       <a href="https://github.com/ma111e/melody/releases/latest"><img src="https://img.shields.io/github/release/ma111e/melody.svg" alt="Latest release"/></a>
    <a href="https://ma111e.github.io/melody/"><img src="https://img.shields.io/badge/%F0%9F%93%9A-Documentation-informational" alt="Documentation"/></a>
    <a href="https://ma111e.github.io/melody/installation"><img src="https://img.shields.io/badge/%F0%9F%93%9A-Installation-informational" alt="Installation"/></a>
    <a href="https://ma111e.github.io/melody/quickstart"><img src="https://img.shields.io/badge/%F0%9F%93%9A-Quickstart-informational" alt="Quickstart"/></a>
    <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-green.svg" alt="Go Report Card"/></a>
</p>

---

Melody is a transparent internet sensor built for threat intelligence and supported by a detection rule framework which allows you to tag packets of interest for further analysis and threat monitoring.

# Table of Contents

   * [Melody](#melody)
   * [Table of contents](#table-of-contents)
   * [Features](#features)
   * [Wishlist](#wishlist)
   * [Use cases](#use-cases)
      * [Internet facing sensor](#internet-facing-sensor)
      * [Stream analysis](#stream-analysis)
   * [Preview](#preview)
   * [Quickstart](#quickstart)
      * [TL;DR](#tldr)
         * [Release](#release)
         * [From source](#from-source)
         * [Docker](#docker)
   * [Rules](#rules)
      * [Rule example](#rule-example)
   * [Logs](#logs)

# Features
Here are some key features of Melody :

+ Transparent capture
+ Write detection rules and tag specific packets to analyze them at scale 
+ Mock vulnerable websites using the builtin HTTP/S server
+ Supports the main internet protocols over IPv4 and IPv6
+ Handles log rotation for you : Melody is designed to run forever on the smallest VPS
+ Minimal configuration required
+ Standalone mode : configure Melody using only the CLI
+ Easily scalable :
    + Statically compiled binary
    + Up-to-date Docker image

# Wishlist
Since I have to focus on other projects right now, I can't put much time in Melody's development.

There is a lot of rom for improvement though, so here are some features that I'd like to implement someday :
+ ~~Dedicated helper program to create, test and manage rules~~ -> Check Meloctl in `cmd/meloctl`
+ Centralized rules management
+ Per port mock application

# Use cases
## Internet facing sensor

+ Extract trends and patterns from Internet's noise 
+ Index malicious activity, exploitation attempts and targeted scanners
+ Monitor emerging threats exploitation
+ Keep an eye on specific threats

## Stream analysis
+ Build a background noise profile to make targeted attacks stand out
+ Replay captures to tag malicious packets in a suspicious stream

# Preview

<p>
<img src="https://raw.githubusercontent.com/ma111e/melody/master/readme/melody_demo.gif" height="600" />
<img src="https://raw.githubusercontent.com/ma111e/melody/master/readme/melody_demo_dash.png" height="600" />
</p>

# Quickstart
[Quickstart details.](https://ma111e.github.io/melody/installation)

## TL;DR
### Release
Get the latest release at `https://github.com/ma111e/melody/releases`.

```bash
make install               # Set default outfacing interface
make cap                   # Set network capabilities to start Melody without elevated privileges
make certs                 # Make self signed certs for the HTTPS fileserver
make enable_all_rules      # Enable the default rules
make service               # Create a systemd service to restart the program automatically and launch it at startup 

sudo systemctl stop melody  # Stop the service while we're configuring it
```

Update the `filter.bpf` file to filter out unwanted packets.

```bash
sudo systemctl start melody     # Start Melody
sudo systemctl status melody    # Check that Melody is running    
```

The logs should start to pile up in `/opt/melody/logs/melody.ndjson`.

```bash
tail -f /opt/melody/logs/melody.ndjson # | jq
```

### From source

```bash
git clone https://github.com/ma111e/melody /opt/melody
cd /opt/melody
make build
```

Then continue with the steps from the [release](#release) TL;DR.

### Docker

```bash
make certs                           # Make self signed certs for the HTTPS fileserver
make enable_all_rules                # Enable the default rules
mkdir -p /opt/melody/logs
cd /opt/melody/

docker pull ma111e/melody:latest

MELODY_CLI="" # Put your CLI options here. Example : export MELODY_CLI="-s -i 'lo' -F 'dst port 5555' -o 'server.http.port: 5555'"

docker run \
    --net=host \
    -e "MELODY_CLI=$MELODY_CLI" \
    --mount type=bind,source="$(pwd)/filter.bpf",target=/app/filter.bpf,readonly \
    --mount type=bind,source="$(pwd)/config.yml",target=/app/config.yml,readonly \
    --mount type=bind,source="$(pwd)/var",target=/app/var,readonly \
    --mount type=bind,source="$(pwd)/rules",target=/app/rules,readonly \
    --mount type=bind,source="$(pwd)/logs",target=/app/logs/ \
    ma111e/melody
```

The logs should start to pile up in `/opt/melody/logs/melody.ndjson`.

# Rules

[Rule syntax details.](https://ma111e.github.io/melody/installation)

## Example

```yaml
CVE-2020-14882 Oracle Weblogic Server RCE:
  layer: http
  meta:
    id: 3e1d86d8-fba6-4e15-8c74-941c3375fd3e
    version: 1.0
    author: BonjourMalware
    status: stable
    created: 2020/11/07
    modified: 2020/20/07
    description: "Checking or trying to exploit CVE-2020-14882"
    references:
      - "https://nvd.nist.gov/vuln/detail/CVE-2020-14882"
  match:
    http.uri:
      startswith|any|nocase:
        - "/console/css/"
        - "/console/images"
      contains|any|nocase:
        - "console.portal"
        - "consolejndi.portal?test_handle="
  tags:
    cve: "cve-2020-14882"
    vendor: "oracle"
    product: "weblogic"
    impact: "rce"
```

# Logs

[Logs content details.](https://ma111e.github.io/melody/layers)

## Example

Netcat TCP packet over IPv4 :

```json
{
  "tcp": {
    "window": 512,
    "seq": 1906765553,
    "ack": 2514263732,
    "data_offset": 8,
    "flags": "PA",
    "urgent": 0,
    "payload": {
      "content": "I made a discovery today. I found a computer.\n",
      "base64": "SSBtYWRlIGEgZGlzY292ZXJ5IHRvZGF5LiAgSSBmb3VuZCBhIGNvbXB1dGVyLgo=",
      "truncated": false
    }
  },
  "ip": {
    "version": 4,
    "ihl": 5,
    "tos": 0,
    "length": 99,
    "id": 39114,
    "fragbits": "DF",
    "frag_offset": 0,
    "ttl": 64,
    "protocol": 6
  },
  "timestamp": "2020-11-16T15:50:01.277828+01:00",
  "session": "bup9368o4skolf20rt8g",
  "type": "tcp",
  "src_ip": "127.0.0.1",
  "dst_port": 1234,
  "matches": {},
  "inline_matches": [],
  "embedded": {}
}
```
