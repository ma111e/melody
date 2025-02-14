### TL;DR
#### Release
Get the latest release at `https://github.com/ma111e/melody/releases`.

```bash
make install            # Set default outfacing interface
make cap                # Set network capabilities to start Melody without elevated privileges
make certs              # Make self signed certs for the HTTPS fileserver
make default_rules      # Enable the default rules
make service            # Create a systemd service to restart the program automatically and launch it at startup 
                        # Note that the script expects that you've installed Melody in /opt/melody

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

#### From source

```bash

git clone https://github.com/ma111e/melody /opt/melody
cd /opt/melody
make build
```

Then continue with the steps from the [release](#release) TL;DR.

#### Docker

```bash
make certs                        # Make self signed certs for the HTTPS fileserver
make default_rules                # Enable the default rules
mkdir -p /opt/melody/logs
cd /opt/melody/

docker pull bonjourmalware/melody:latest

MELODY_CLI="" # Put your CLI options here. Example : export MELODY_CLI="-s -i 'lo' -F 'dst port 5555' -o 'server.http.port: 5555'"

docker run \
    --net=host \
    -e "MELODY_CLI=$MELODY_CLI" \
    --mount type=bind,source="$(pwd)/filter.bpf",target=/app/filter.bpf,readonly \
    --mount type=bind,source="$(pwd)/config.yml",target=/app/config.yml,readonly \
    --mount type=bind,source="$(pwd)/var",target=/app/var,readonly \
    --mount type=bind,source="$(pwd)/rules",target=/app/rules,readonly \
    --mount type=bind,source="$(pwd)/logs",target=/app/logs/ \
    bonjourmalware/melody
```

The logs should start to pile up in `/opt/melody/logs/melody.ndjson`.

### Before we start

Hi !

You'll need the following info :

+ The name of the interface on which you want Melody to listen to 
    + wlp3s0, ens3, enp0s25, eth0...
+ All the IP addresses to exclude from monitoring

!!! Tip
    You will most likely want to ban at least the address from which you access the server and the addresses used for monitoring (don't forget your hosting provider's).
    
Clone the repo to get the default rules and configuration files :

```git
git clone https://github.com/ma111e/melody
```

!!! Tip
    You can also write the configuration files later or use CLI options to use it as a standalone binary

### Firewall
Don't forget to check your firewall to ensure you're not blocking packets from reaching the sensor.

!!! Tip
    You might want to disengage your hosting provider DDoS protection.

### Build from source

!!! warning
    You need the libpcap headers before building Melody. Install them with :
    ```bash
    sudo apt update
    sudo apt install libpcap-dev
    ```

Build Melody with :

```
make build
```

or

```
go build -ldflags="-s -w -extldflags=-static" -o melody
sudo setcap cap_net_raw,cap_setpcap=ep ./melody
```

### Grab a release

You can grab the latest release by visiting [https://github.com/ma111e/melody/releases/latest](https://github.com/ma111e/melody/releases/latest).

### Docker
You can also use Docker and pull the image from Docker Hub :

```
docker pull bonjourmalware/melody:latest
```

Don't forget to create the `logs` directory, or the mount will fail :

```
mkdir -p /opt/melody/logs
```

Run it with :

```
make docker_run
```

or

```
export MELODY_CLI="" # Put your CLI options here. Example : export MELODY_CLI="-s -i 'lo' -F 'dst port 5555' -o 'server.http.port: 5555'"

docker run \
    --net=host \
    -e "MELODY_CLI=$MELODY_CLI" \
    --mount type=bind,source="$(pwd)/filter.bpf",target=/app/filter.bpf,readonly \
    --mount type=bind,source="$(pwd)/config.yml",target=/app/config.yml,readonly \
    --mount type=bind,source="$(pwd)/var",target=/app/var,readonly \
    --mount type=bind,source="$(pwd)/rules",target=/app/rules,readonly \
    --mount type=bind,source="$(pwd)/logs",target=/app/logs/ \
    melody
```

## Configuration
### Melody configuration
All the available configuration options are listed with their default values in the `config.yml` file.

You'll want to look at a few things before getting started :

+ Set the `listen.interface` to the one on which you want Melody to be listening on

!!! Tip
    On most recent linux distribution, you can run `route | grep '^default' | grep -o '[^ ]*$'` to find the default WAN card. Note that you'll need the `net-tools` package (`sudo apt install net-tools`) in order to use the `route` command. 

!!! Tip
    On Windows, you'll want an interface name like `\Device\NPF_{4E273621-5161-46C8-895A-48D0E52A0B83}`.
    
    If you find an interface name with `TCP` in place of `NPF`, try swaping both.
    
    See [Find Windows interfaces]() for more details.
     
    Don't forget to wrap your string with `'` to prevent the parsing of the escaping `\` : 
    ```
    interface: '\Device\NPF_{4E273621-5161-46C8-895A-48D0E52A0B83}'
    ```

!!! Note
    Note that Melody listen on `lo` by default. You can override the listening interface with the `-i` switch.

+ The dummy HTTP/S servers are enabled by default. Disable it if you're not interested by this kind of data, or you're putting Melody next to a web application

+ Default rules are disabled by default. You can enable them by creating a symlink for each rule in the active rule directory specified in the configuration file (`$melody/rules/rules-enabled` by default)

!!! Tip
    To create a symlink, use the following command from the root of the projet :
    ```bash
    ln -rs ./rules/rules-available/$rule.yml ./rules/rules-enabled/
    ```
    
    Use a wildcard to enable all the rulesets :
    ```bash
    ln -rs ./rules/rules-available/*.yml ./rules/rules-enabled/
    ```

### HTTP/S server
In order to capture the full HTTP transactions, the client must have a server to connect to. To ease that process, a dummy HTTP/S server is available.

The default configuration is to answer `200 OK` on every routes, along with a `Server: Apache` header.

#### iptables

To capture the HTTP traffic your server receives on every ports, I recommend using `iptables` to redirect the data from every ports to the one Melody is listening on.

!!! Danger
    Be very careful while applying these modifications. You must at least exclude your remote connection port using the `! --dports ` switch, or you will be locked out. 

To achieve this, add a rule to your `PREROUTING` table :

```
sudo iptables -A PREROUTING -t nat -i $INTERFACE -p tcp -m multiport ! --dports $REMOTE_ACCESS_PORT,$ANOTHER_EXCLUDED_PORT -j REDIRECT --to-port $MELODY_HTTP_PORT
```

Example : 

```
sudo iptables -A PREROUTING -t nat -i ens3 -p tcp -m multiport ! --dports 1234,5678 -j REDIRECT --to-port 10800
```

Here the ports `1234` and `5678` have been excluded from the redirection.

!!! Note
    Using the `sudo iptables -t nat -L PREROUTING -n -v` command, you should see something like this :
    ```
    Chain PREROUTING (policy ACCEPT 1226K packets, 57M bytes)
     pkts bytes target     prot opt in     out     source               destination         
      25M 1243M REDIRECT   tcp  --  ens3   *       0.0.0.0/0            0.0.0.0/0            multiport dports  !1234,5678 redir ports 10800
    ``` 
    
!!! Important
    This is only used to virtually connect the HTTP server to all the ports.
    
    As Melody sits on the data link layer, the program will receive the packets before being handled by network layer programs such as `iptables` or `ufw`. 

!!! Note
    I won't cover it here as I'm now knowledgeable enough on this, but you'll want to look at advfirewall and the `portproxy` command to replicate this on Windows.

## Berkley Packet Filter (BPF)

The next step is the customization of the `filter.bpf` file. **This is where you filter the data that reaches Melody.**

By default, only inbound traffic is allowed and all 127.0.0.0/24 subnet are banned. 

!!! Note
    You can use the `-F|--filter` switch to set a filter via CLI.    

!!! Tip
    If you're using a VPS, you might need to filter out the IP addresses your hosting provider uses to check the status of your server. 

### Source IP filtering
Use `[src|dst] net <network>` to filter packets according to their IP.

Example :

```bpf
inbound and not net 127.0.0.1
```

You can specify a range to exclude using the CIDR notation :

```bpf
inbound and not net 127.0.0.0/24
```

### Port filtering
Use `[src|dst] port <port>` to filter packets according to their port.

Example :

```bpf
not port 1234
```

You can specify a range to exclude using the CIDR notation and the `portrange` keyword :

```bpf
not portrange 1234-5678
```

Your `filter.bpf` should look like this :

```
inbound
and not port 1234
and not net 127.0.0.0/24
and not net 192.0.2.1
```

!!! Important
    Your file should always start with the `inbound` keyword. I recommend adding your filter rules below, starting with an `and` keyword.


### Advanced

[Here is all you need to know about the BPF syntax](https://biot.com/capstats/bpf.html) and [here is a great source of examples to get quickly started](https://www.ibm.com/support/knowledgecenter/SS42VS_7.4/com.ibm.qradar.doc/c_forensics_bpf.html).

## Rules

Melody rules are used to apply tags on matching packets. They have multiple use cases, such as monitoring emerging threats, automated droppers, vulnerability scanners, etc.

You can look into the `$melody/rules/rule-available` and `$melody/internal/rules/test_resources` folders to quickly find working examples. 

### Basics

A rule file can contain multiple rules description and constitute a ruleset. 

Here is a rule example that detects CVE-2020-14882 (Oracle Weblogic RCE) scans or exploitation attempts by matching either of the two URI and the HTTP verb :

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

This rule is part of the default `server.yml` ruleset and tag CVE-2020-14882-related packets.

If a packet match, its tags will be appended to the `matches` field in the packet's log.

See the [Rules section](/rules) for details on the rules syntax.

!!! Important
    You must generate a new UUIDv4 [for each rule you create](/rules#meta).

### Output
#### Stdout
You can redirect the output to stdout by using the `-s` switch.

!!! Example
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

#### Dump
A dump mode is also available with the `-d|--dump` switch. Similar to `tcpdump`, it will print raw packets to stdout instead of Melody json lines.

!!! Example
    ```toml
    PACKET: 76 bytes, wire length 76 cap length 76 @ 2020-11-16 15:46:00.927899 +0100 CET
    - Layer 1 (14 bytes) = Ethernet {Contents=[..14..] Payload=[..62..] SrcMAC=00:00:00:00:00:00 DstMAC=00:00:00:00:00:00 EthernetType=IPv4 Length=0}
    - Layer 2 (20 bytes) = IPv4     {Contents=[..20..] Payload=[..42..] Version=4 IHL=5 TOS=0 Length=62 Id=39110 Flags=DF FragOffset=0 TTL=64 Protocol=TCP Checksum=41969 SrcIP=127.0.0.1 DstIP=127.0.0.1 Options=[] Padding=[]}
    - Layer 3 (32 bytes) = TCP      {Contents=[..32..] Payload=[..10..] SrcPort=58766 DstPort=1234(search-agent) Seq=1906765476 Ack=2514263732 DataOffset=8 FIN=false SYN=false RST=false PSH=true ACK=true URG=false ECE=false CWR=false NS=false Window=512 Checksum=65074 Urgent=0 Options=[TCPOption(NOP:), TCPOption(NOP:), TCPOption(Timestamps:1712590417/1712586943 0x66140e51661400bf)] Padding=[]}
    - Layer 4 (10 bytes) = Payload  10 byte(s)
    ```

!!! Note
    All the console messages are printed to stderr in order to allow piping Melody's data into `jq`.
