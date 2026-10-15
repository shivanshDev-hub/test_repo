# OPSEC Field Guide 3: Data Reduction (Targeted Edition)

Oct 10, 2026 · @tiger

## Read this first

This field book teaches you to leak less: fewer identifiers, fewer logs, and fewer third parties sitting between you and the people you talk to. It is the third and highest tier of the series, written for people facing a capable, patient adversary: a determined investigator, a well-resourced stalker or employer, or a state. Every part follows one pattern: a short **Why** (theory a beginner can follow), a numbered **Do this** (steps you can copy), and **Mistakes that break it**.

At this tier, habits break anonymity before maths does. Real-world deanonymisation almost always comes from mistakes (one reused username, one login to a personal account, one photo with GPS metadata), not from cracked encryption. That is why every part ends with its own list of mistakes.

**What this guide assumes you have already done**

Guide 1 (Covered) and Guide 2 (Hardened) are treated as complete, and nothing from them is repeated here. In particular, you already have:

- **From Guide 1:** a threat model, no real-identity sign-ups, old accounts deleted, a password manager, Linux, VMs and sandboxes, MAC spoofing, metadata removal, full-disk encryption, privacy alternatives for everyday apps, and cash for normal payments.
- **From Guide 2:** a VPS or obfuscated VPN with a kill switch, Qubes OS, GrapheneOS on a spare phone, DNS-level ad blocking, privacy frontends, crypto purchasing basics, secure erasure, hardware security keys, and unneeded radios switched off.

If any of that is missing, go back first. Advanced tools on top of a leaky base only give you false confidence.

**How to read this guide**

- Level tags: **\[Beginner\]** works with no prior knowledge, **\[Intermediate\]** needs a terminal, **\[Advanced\]** needs hardware or a server.
- Commands sit in code blocks. Replace anything in `<angle brackets>` with your own value.
- **Theory-only zone** marks the places where I explain the idea, how it fails and why it matters, but give no build or operate steps. The steps there would mainly help someone harm others or dodge legal accountability, rather than protect themselves. The four zones are in Parts 3, 10 (two of them) and 11.

**Staying legal.** Privacy tools are legal in most countries, but some restrict VPNs, strong encryption, or unlicensed radio. Check local law before you deploy anything. I am not a lawyer.

## Part 1: IP protection (Tor, Whonix, Tails, bridges)

Your IP address is the first identifier every server sees, so Part 1 shows how to stop handing it over: Tor for the path, Tails or Whonix for the machine, and bridges for the doorstep.

### Why: what Tor does and does not do \[Beginner\]

Tor wraps your traffic in three layers and sends it through three relays: an entry guard, a middle relay and an exit. The guard knows who you are but not where you are going. The exit knows where you are going but not who you are. No single relay knows both.

**Tor hides** your IP from the sites you visit, and hides which sites you visit from your ISP (they only see that you use Tor).

**Tor does not hide** what you type into a site, your login names, a distinctive browser setup, or a watcher who can see both ends of your connection.

**The correlation attack in one paragraph.** If someone can watch your connection into Tor and the server you reach at the same moment, they can match the pattern of timing and volume (for example, four bursts at 12:01:03 on your side and four bursts at 12:01:03.4 at the server). No encryption stops this, because it is statistics on metadata. Your defences are: fewer and shorter sessions, an entry point nobody is watching, and identities that never mix. Bridges (below) help with the second one only.

### My honest opinion on stacking Tor, VPNs and VPSs

- **Stack less.** Every extra hop is one more company that can log you or be compelled to. Add a layer only if you can say in one sentence which threat it removes.
- **Tor then VPN** (your traffic exits Tor into a VPN) is advanced, breaks Tor's per-site circuit separation, and ties your exit traffic to a paid account. Whonix documents it but treats it as an expert option.
- **The VPS use that is worth it** is a private bridge you control (see below). A VPS used as a remote desktop is IP relocation, not anonymity, and your payment trail leads to it.

### Pick your machine: Tails or Whonix

|  | Tails | Whonix |
| --- | --- | --- |
| What it is | A live OS on a USB stick that forgets everything at shutdown | Two virtual machines (Gateway and Workstation) on a host OS |
| Best for | Short, high-risk sessions on a computer you do not want to leave traces on | Daily use and long projects that need saved files and a stable, separate identity |
| Stops | Traces on the computer, non-persistent malware | IP leaks even if the Workstation is hacked, because it has no direct route out |
| Weak spot | Host firmware, a stolen USB with persistence enabled | A compromised host OS or a VM escape |
| Setup time | About 20 minutes | About 45 to 60 minutes |

Tails is currently in its 7.x series (Debian 13 base) with monthly releases plus emergency fixes, for example [7.10.1, a kernel fix](https://news.tuxmachines.org/n/2026/08/06/Relaxing_Weeks_Ahead.shtml). Whonix is on 18.x (Debian 13 base); Whonix 17 has reached end of security support, so do not start on it.

### Do this A: Tails on a USB stick \[Beginner\]

1. Get a USB stick of at least 8 GB and a computer to prepare it on.
2. Open [tails.net/install](https://tails.net/install/), pick your operating system and follow the page. It walks you through download, verification and writing the image. Do not skip the verification step.
3. Restart the target computer from the stick. Hold the boot-menu key while powering on (commonly Esc, F12, F10 or F2; on a Mac, hold Option).
4. On the Welcome Screen keep the defaults. Leave the administration password unset unless you truly need to install software.
5. Connect to Wi-Fi. In the Tor connection screen choose automatic connection; choose the option to hide that you use Tor if your network is hostile (see Bridges below).
6. Open Tor Browser and visit [check.torproject.org](https://check.torproject.org). It should say you are using Tor.
7. Shut down from the system menu when you finish. Tails wipes RAM on shutdown.
8. Update whenever Tails asks. Do not postpone security updates.

**Mistakes that break it:** turning on Persistent Storage out of habit, resizing or maximising the Tor Browser window, logging into accounts tied to your real name, and opening downloaded documents while still online.

### Do this B: Whonix in VirtualBox \[Beginner\]

You need an x86-64 computer with virtualisation enabled in BIOS, 8 GB RAM (16 GB is better) and about 50 GB free disk.

1. Install VirtualBox from [virtualbox.org](https://www.virtualbox.org).
2. Open [whonix.org/wiki/VirtualBox](https://www.whonix.org/wiki/VirtualBox), download the image and verify its signature as that page explains.
3. Double-click the downloaded `.ova` file and accept the defaults. Both machines import together.
4. Start **Whonix-Gateway** first. When the connection wizard appears, choose Connect (or Configure if you need bridges) and wait for Tor to bootstrap.
5. Start **Whonix-Workstation**, open Tor Browser and check [check.torproject.org](https://check.torproject.org).
6. In each machine's terminal run `sudo apt update && sudo apt full-upgrade`.
7. Take a snapshot of the clean Workstation. Use a separate Workstation clone for each identity so identities never share a machine.

For stronger isolation use [Whonix on KVM](https://www.whonix.org/wiki/KVM) or [Qubes-Whonix](https://www.whonix.org/wiki/Qubes), which the Whonix project considers the strongest option. (Qubes itself is set up in Guide 2.)

### Do this C: Live USB with LUKS, no permanent PC \[Intermediate\]

The idea: store nothing on the computer you borrow. Carry the system and your files on encrypted media. How LUKS works and how to choose a passphrase are covered in Guide 1; this section only adds the "no permanent PC" setup.

1. **Easiest:** use Tails. Persistent Storage, if you enable it, is an encrypted volume on the same stick.
2. **Your own encrypted OS:** install Debian onto an external SSD with the installer option *Guided - use entire disk and set up encrypted LVM*. Unplug the internal drive first so you cannot pick the wrong target. Then run Whonix inside it.
3. **Back up the LUKS header** (check the partition name with `lsblk` first), and store the copy offline, never on the same drive:

```
sudo cryptsetup luksHeaderBackup /dev/<your encrypted partition> --header-backup-file luks-header.img
```

**Mistakes that break it:** suspending instead of shutting down (keys stay in RAM), leaving the drive unlocked and plugged in, and forgetting that some countries can legally compel you to reveal a passphrase. Know your local law.

### Do this D: Physically split Whonix across two machines \[Advanced\]

**Why.** In normal Whonix the Gateway and Workstation share one computer, so a hypervisor bug could bridge them. With two computers, a hacked Workstation still has no path to your real IP, because the Gateway is separate hardware. The Whonix project says this [cuts the trusted computing base on the gateway side by more than half](https://www.whonix.org/wiki/Physical_Isolation).

**You need:** Machine G (Gateway), any x86 PC with two network ports (one can be a USB Ethernet adapter); Machine W (Workstation), a laptop or PC; one Ethernet cable between G's second port and W; G's first port goes to your router.

```
Internet --- Router --- [port 1] G: Gateway (runs Tor) [port 2] --- cable --- W: Workstation
                                                                               (Wi-Fi and Bluetooth removed)
```

Traffic leaves W only through G, and G alone touches your router. If you ever see a second route from W, the setup is wrong.

1. Install a minimal Debian (no desktop) on G.
2. Follow the *How-to: Install Whonix-Gateway* section of [whonix.org/wiki/Physical\_Isolation](https://www.whonix.org/wiki/Physical_Isolation) to turn G into the Gateway. That page is the source of truth; commands change between releases.
3. Assign roles: G's external port faces the router, G's internal port faces W.
4. Prepare W: run the Whonix-Workstation on it, and remove or disable W's Wi-Fi and Bluetooth so the cable to G is its only path out.
5. Test: from W, [check.torproject.org](https://check.torproject.org) should say Tor. Then stop Tor on G; W must lose internet completely (fail-closed).

**Mistakes that break it:** plugging W into the router just to update, leaving W's Wi-Fi card active, and ignoring MAC addresses (the Whonix page itself notes that part is under-documented). Qubes-Whonix with physical isolation is unsupported.

### Bridges, and what they really do against correlation \[Intermediate\]

A **bridge** is a Tor entry point that is not in the public relay list. Because your ISP cannot see a connection to a known Tor guard, it learns less about you. Bridges also work in places that block Tor.

**Honest limit:** a bridge does not defeat end-to-end correlation by itself. It removes one watcher (the one at your ISP who sees you connect to a public Tor IP). Whoever runs the bridge's server can still see your IP, so a private bridge on a rented VPS shifts the observer from your ISP to your VPS host. Real correlation defences are short sessions, one identity per session, and not logging into anything linked to you.

**Use a public bridge \[Beginner\].** In Tor Browser open Settings, then Connection, then Bridges, and either pick a built-in transport (obfs4, or Snowflake if obfs4 is blocked) or request one from [bridges.torproject.org](https://bridges.torproject.org).

**Run your own private bridge on a VPS \[Advanced\].**

1. Rent a small Debian or Ubuntu VPS (1 Mbit/s up and down is the Tor Project's stated minimum for a bridge). Remember the host knows who pays.
2. Install Tor and *lyrebird* (the maintained successor to obfs4proxy) by following the official [bridge setup guide](https://community.torproject.org/relay/setup/bridge/). It has a Debian and Ubuntu page with the correct package repository lines.
3. Edit `/etc/tor/torrc`:

```
BridgeRelay 1
ORPort <a high port you choose, not 9001>
ServerTransportPlugin obfs4 exec /usr/bin/lyrebird
ServerTransportListenAddr obfs4 0.0.0.0:<another port you choose>
ExtORPort auto
ContactInfo <your throwaway email>
PublishServerDescriptor 0
BridgeDistribution none
```

4. Open both chosen ports in the VPS firewall, then restart Tor with `sudo systemctl restart tor`.
5. Read your bridge line with `sudo cat /var/lib/tor/pt_state/obfs4_bridge_line.txt`, replace the placeholder IP and port with your real ones, and paste it into Tor Browser under *Add a bridge manually*.

`PublishServerDescriptor 0` keeps your bridge off the public lists, so only people you give the line to can use it. Never turn an existing public relay into a bridge on the same address; it is already blocklisted.

### Tor transparent proxy with pf(4) on OpenBSD \[Advanced\]

**Why.** A transparent proxy forces every device behind a box through Tor without configuring each device. **Warning:** the Tor Project advises against torifying arbitrary apps, because apps leak identity (cookies, fingerprints, login names). Use this for devices you control and trust, and use Tor Browser for browsing.

**Do this** (OpenBSD box with two network cards: `em0` to the internet, `em1` to your devices):

1. `doas pkg_add tor`, and set the correct system clock first (Tor needs accurate time).
2. Edit `/etc/tor/torrc`. The `pf-divert` option is the documented way to run a [TransPort on OpenBSD](https://man.openbsd.org/pf.conf):

```
VirtualAddrNetworkIPv4 10.192.0.0/10
AutomapHostsOnResolve 1
TransPort 127.0.0.1:9040
TransProxyType pf-divert
DNSPort 127.0.0.1:5353
```

3. Turn on forwarding: `echo 'net.inet.ip.forwarding=1' | doas tee -a /etc/sysctl.conf` and `doas sysctl net.inet.ip.forwarding=1`.
4. Write `/etc/pf.conf`:

```
ext_if = "em0"
int_if = "em1"
set skip on lo
block all

# only the Tor daemon may talk to the outside
pass out on $ext_if proto tcp from (self) user _tor

# device TCP goes into Tor's TransPort, device DNS into Tor's DNSPort
pass in on $int_if inet proto tcp from $int_if:network to any divert-to 127.0.0.1 port 9040
pass in on $int_if inet proto udp from $int_if:network to any port 53 rdr-to 127.0.0.1 port 5353
```

5. Check syntax with `doas pfctl -nf /etc/pf.conf`, load with `doas pfctl -f /etc/pf.conf`, then `doas rcctl enable tor && doas rcctl start tor`. (OpenBSD's [Tor bridge guide](https://community.torproject.org/relay/setup/bridge/openbsd) shows the same Tor install on OpenBSD.)
6. Give each device a static address on `em1`'s network, with the router as gateway and DNS.
7. **Leak test:** on a device run `curl https://check.torproject.org/api/ip` and expect `"IsTor":true`. Then run `doas rcctl stop tor` on the router: the device must lose all internet. If anything still gets out, your rules leak.

**Note:** I wrote this configuration from the Tor manual and the pf.conf(5) manual, but I could not run it on real hardware here. Test it on a spare machine before you trust it, and compare against `man pf.conf` on your OpenBSD release.

**Mistakes that break it:** forgetting `block all` (UDP and IPv6 then leak), allowing the router itself to browse, and sharing one Tor circuit between identities that should never mix.

## Part 2: I2P, the garlic network

Use Tor to visit the normal internet anonymously, and use I2P when you want to talk to services, people and torrents that live inside an anonymous network. They solve different problems, so a good field kit has both.

### Why: how I2P differs from Tor \[Beginner\]

- **Everyone is a router.** Your I2P router relays other people's traffic by default, so your own traffic hides in a crowd. In Tor, most users are not relays.
- **Unidirectional tunnels.** Requests leave through one set of tunnels and replies come back through another. An observer sees half a conversation on each path, which makes timing matches harder.
- **Garlic routing.** Several messages are bundled in one encrypted packet (like cloves in a bulb), so an observer cannot tell how many separate messages a packet holds.
- **No central directory authority.** Routers find each other through a distributed database (the netDb) run by volunteer floodfill routers.
- **Built for internal services.** Sites end in `.i2p` (called eepsites). Clearnet exits, called outproxies, exist but are few and not the main design goal.

|  | Tor | I2P |
| --- | --- | --- |
| Best at | Anonymous access to the normal internet; onion services | Hidden services, torrents, chat and mail inside the network |
| Clearnet browsing | Strong, with the hardened Tor Browser | Weak, depends on a few outproxies |
| Browser fingerprint defence | Excellent (uniform Tor Browser) | Basic (you configure a normal browser) |
| Network size and scrutiny | Larger, heavily studied | Smaller, less studied |
| What your ISP sees | Tor connections (hidden with bridges) | A router talking to many peers; there is no Tor-style bridge system |
| Speed | Medium | Slower to start, good once integrated |

**A common myth.** You will sometimes read that I2P completely eliminates the network sniper problem. That claim is too strong. The Tor *sniper attack* abused circuit flow control to exhaust a relay's memory, and Tor patched it. I2P's tunnel design means that exact mechanism does not carry over directly, as far as I know. But I2P has its own attack surface (Sybil routers, floodfill abuse, DoS, a small network where a few hostile routers matter more), and nobody can prove a whole class of attacks is gone. Treat I2P as differently shaped, not immune.

### Do this: set up I2P \[Beginner to Intermediate\]

Pick **one** of these routes.

**Route 1, Java I2P desktop.** Download from [geti2p.net/en/download](https://geti2p.net/en/download), install, and start *I2P Router*. The router console opens at `http://127.0.0.1:7657`.

**Route 2, Docker (any OS with Docker):**

```
docker run -d --name i2p \
  -v i2pconfig:/i2p/.i2p \
  -p 127.0.0.1:4444:4444 -p 127.0.0.1:7657:7657 \
  geti2p/i2p:latest
```

The [Docker Hub page](https://hub.docker.com/r/geti2p/i2p) says to open the console at `http://127.0.0.1:7657` to finish the setup wizard. Binding to `127.0.0.1` keeps the proxy private to your machine.

**Route 3, i2pd (light, for small servers and Raspberry Pi):** `sudo apt install i2pd` on Debian or Ubuntu. Its web console is at `http://127.0.0.1:7070` and its HTTP proxy is on `127.0.0.1:4444`. See [i2pd.readthedocs.io](https://i2pd.readthedocs.io).

Then:

1. **Wait 10 to 15 minutes** after first start. The router needs time to find peers and build tunnels. The console should show the network status as OK.
2. **Make a separate browser profile.** In Firefox run `firefox -P`, create a profile named `i2p`, then open Settings, Network Settings, choose Manual proxy configuration, set HTTP Proxy `127.0.0.1` port `4444`, tick *Also use this proxy for HTTPS*, and tick *Proxy DNS*. Use this profile for I2P only.
3. **Test:** visit `http://i2p-projekt.i2p` and `http://notbob.i2p` in that profile.
4. **Share bandwidth.** In the console go to Configuration, then Bandwidth, and share as much as you comfortably can. A router that stays online and well connected integrates better and routes traffic for others.
5. **Leave tunnel length at the default (3 hops).** Shortening it speeds things up and weakens anonymity.

### Mistakes that break it

- Using the same browser profile for I2P and the normal web.
- Restarting the router every few minutes, so it never integrates.
- Assuming I2P hides that you run it. Your peers and ISP can see I2P-style traffic.
- Using I2P for clearnet anonymity. Use Tor Browser for that.

**Learn more:** [I2P documentation](https://geti2p.net/en/docs) and the [I2P ports list](https://i2p.net/de/docs/overview/ports) if you need to open firewall ports.

## Part 3: Public Wi-Fi and a separate Wi-Fi adapter

A public network keeps your home connection out of the picture, and a dedicated, randomised adapter keeps your hardware identity out of it too. Used with Tor, this is a standard method for journalists and activists; used carelessly it creates a new trail of cameras, receipts and logs.

### Why: what your Wi-Fi gives away \[Beginner\]

- **Your home IP** is tied by your ISP to your billing identity. A public network moves that off your account.
- **Your MAC address** is the hardware ID of your Wi-Fi chip. Networks you join can log it, and it can follow you from place to place if it never changes. Spoofing it is covered in Guide 1.
- **Network names you know.** Many devices shout the names of previously joined networks (like your home SSID) when they scan. That is a fingerprint of your life.
- **The new leaks public Wi-Fi creates:** CCTV, loyalty cards and receipts, captive portals that ask for an email or phone number, and fake hotspots (evil twins) that imitate a real one.

A **separate USB Wi-Fi adapter** matters because your laptop's built-in card may already be tied to you through earlier use. A fresh adapter, with a randomised MAC, starts with no history.

### Do this \[Beginner to Intermediate\]

1. **Choose an adapter with proven Linux support.** Look up the chipset before buying (MediaTek MT76xx and Atheros ath9k\_htc devices are common choices with in-kernel drivers). Check the [Tails system requirements](https://tails.net/doc/about/requirements/) if you plan to use Tails. Buy it in a way that does not tie it to you, for example cash in a shop.
2. **Randomise the MAC.** Tails does this by default, so leave that setting on at the Welcome Screen. On another Linux system (including the host that runs Whonix), apply the MAC-spoofing steps from Guide 1 to the adapter, then run `ip link show` and confirm the address differs from the one printed on the adapter.
3. **Forget old networks.** In your network settings delete saved networks you do not need, so the device stops announcing them.
4. **Before you arrive, switch your phone off** (better, in a Faraday bag; see Part 10), not at the door. Cell towers and the phone's own Wi-Fi and Bluetooth place you at the venue.
5. **Connect, then go straight to Tor.** If a captive portal needs a click-through, let Tails open its captive-portal browser, accept, and close it. If the portal wants a phone number or ID, leave.
6. **Use a bridge** (Part 1) if you also want the Wi-Fi owner not to learn you use Tor.
7. **Vary your habits.** Do not use the same café at the same hour every week, and pay in cash.

### Mistakes that break it

- Logging into any personal account while on this setup.
- Carrying a phone that is switched on.
- Using the adapter on your home network, which links it to you for good.
- Trusting the hotspot name. Anyone can create an evil twin; Tor and HTTPS protect you from what it can read.

**Theory-only zone: Wi-Fi you are not allowed to use.** Cracking a protected network, borrowing a neighbour's connection without permission, or cloning another device's MAC to get in is unauthorised access in most countries. It also pins your activity on the owner and can frame an innocent person. It is not privacy; it is shifting your risk onto someone else. I explain the concept only: a MAC is a label, not a password, so cloning one proves nothing about who you are, and networks that filter by MAC give no real security. For genuine anonymity use networks that are open to the public, with their permission.

## Part 4: Self-sustained network chains

These are networks you can run without trusting an ISP-controlled internet or a central company, and the honest summary is that Reticulum and Yggdrasil are practical today, cjdns and dn42 are mostly for learning, and Netsukuku is an idea more than a tool.

### Why: overlays and meshes \[Beginner\]

An **overlay network** runs on top of the normal internet (Yggdrasil, cjdns, dn42). A **mesh** can run directly over radio or cable with no internet at all (Reticulum, LoRa meshes). In both, your address is usually derived from a cryptographic key, so nobody hands you an address and nobody can take it away.

**Important:** none of these hide who you are by default. Your neighbours in the network see your real IP (overlays) or can locate your radio (meshes). They give you *independence and encryption*, not anonymity.

| Network | Idea in one line | Anonymity | State today | In this guide |
| --- | --- | --- | --- | --- |
| [Reticulum](https://reticulum.network) | Encrypted networking stack that does not need IP; runs over LoRa, WiFi, TCP and more | Partial: packets carry no source address | Active, labelled beta | Full setup |
| [Yggdrasil](https://yggdrasil-network.github.io) | End-to-end encrypted IPv6 overlay | None | Active, labelled early-stage | Full setup |
| [cjdns](https://github.com/cjdelisle/cjdns) / Hyperboria | Encrypted IPv6 mesh where your address comes from your key | None | Maintained but small | Theory plus tryout |
| [dn42](https://dn42.dev) | A hobby 'mini internet' using real BGP over VPN tunnels | None | Active hobby community | Theory |
| Netsukuku | Server-free fractal mesh routing from the mid-2000s | n/a | Abandoned | Theory |

Hyperboria is the main cjdns network. It is unrelated to the Hyperbola Linux distribution, despite the similar name.

### Yggdrasil, step by step \[Intermediate\]

Yggdrasil's own README calls it an *early-stage* implementation, so treat it as experimental.

1. Install from the [installation page](https://yggdrasil-network.github.io/installation.html) (packages exist for Debian, Windows, macOS, OpenWrt and others).
2. Generate a config: `yggdrasil -genconf > yggdrasil.conf`. The key inside it is your identity. Keep the file private (`chmod 600 yggdrasil.conf`).
3. Add peers. Open `yggdrasil.conf`, find `Peers`, and add two or three entries from the [public peers list](https://github.com/yggdrasil-network/public-peers), for example `"tls://<host>:<port>"`.
4. Run it: `sudo yggdrasil -useconffile yggdrasil.conf` (or enable the systemd service your package installed). On Linux it needs permission to create a network interface.
5. Read your address with `sudo yggdrasilctl getSelf`. It starts with `200::` to `3ff::`.
6. Test by pinging another node's address, for example `ping6 <peer address>`.

**Yggdrasil plus your own VPS (a private mesh):**

1. On the VPS, install Yggdrasil and set `Listen: ["tls://0.0.0.0:<port>"]`. Open that port in the firewall.
2. On each of your devices (home PC, laptop, phone) set `Peers: ["tls://<vps-ip>:<port>"]` and no public peers.
3. Every device now reaches every other over its Yggdrasil address, end-to-end encrypted, with no port forwarding at home.
4. Bind private services (for example SSH) to the Yggdrasil address only, and firewall everything else.

To hide your home IP from the VPS or public peers, newer Yggdrasil versions can reach a peer through a SOCKS proxy using a `socks://` peer URI (check your version's docs). If you do not need to run a full TUN interface, [Yggstack](https://github.com/yggdrasil-network/yggstack) runs without admin rights and can expose a local web server on your Yggdrasil address.

**Mistakes that break it:** listening on all interfaces, assuming the address is anonymous (it is stable and tied to your key), and forgetting that peers see your real IP.

### Reticulum, LXMF, Sideband and LoRa \[Intermediate to Advanced\]

**Why.** Reticulum is built for the worst conditions: slow links, no internet, no central servers. Its docs describe it as a complete stack that does not rely on IP, and it is still labelled beta. **LXMF** is its message format (store-and-forward, like email for mesh). **Sideband** is the friendly app, **NomadNet** is the terminal version, and an **RNode** is an open LoRa radio modem for it.

1. **Install:** `pipx install rns` (or `pip install rns`). Start it with `rnsd`, then check with `rnstatus`.
2. **Install Sideband** on your phone or computer from the [Sideband project page](https://github.com/markqvist/Sideband). Open it, let it create your identity, and note your LXMF address. By default it finds other Reticulum users on your local network.
3. **Message someone:** swap LXMF addresses (in person or over a channel you trust), add the contact, and send. Verify addresses out of band, because a swapped address means talking to a stranger.
4. **Link two sites over the internet:** on a VPS add this to `~/.reticulum/config`:

```
[[Public Server]]
  type = TCPServerInterface
  enabled = yes
  listen_ip = 0.0.0.0
  listen_port = 4242
```

and on your client:

```
[[My VPS]]
  type = TCPClientInterface
  enabled = yes
  target_host = <vps-ip>
  target_port = 4242
```

5. **Add LoRa with an RNode.** Get a supported LoRa board and flash it with `rnodeconf --autoinstall` (from `pip install rns`) or the web flasher listed in the [Reticulum manual](https://reticulum.network/manual/). Then add:

```
[[RNode LoRa]]
  type = RNodeInterface
  enabled = yes
  port = /dev/ttyUSB0
  frequency = <Hz, license-free band for your region>
  bandwidth = 125000
  txpower = 7
  spreadingfactor = 8
  codingrate = 5
```

6. **Run a propagation node** (a store-and-forward relay for messages) with the `lxmd` tool if you want your group to receive messages while offline.

**Why LoRa and what it costs.** LoRa reaches kilometres on a few milliwatts with no infrastructure, but carries only a trickle of data, and **radio is locatable**: anyone with direction-finding equipment can find a transmitter. Never use LoRa to hide *where* you are. Check your national regulator's rules for the license-free band, power and duty cycle before transmitting.

**Mistakes that break it:** trusting unverified LXMF addresses, using a band or power level that is illegal where you live, and assuming 'encrypted' means 'untraceable'.

### cjdns and Hyperboria: theory and failure points

cjdns gives every node an IPv6 address (starting `fc`) derived from its public key, and routes with a distributed hash table. Its README admits it is built on a 'bold and unproven' assumption that a non-hierarchical network can scale. **Where it fails:** (1) no anonymity, peers see your IP and your stable address; (2) you traditionally had to find a human peer to join, and the old community peers list is reported stale; (3) it uses its own crypto and routing, which have had less scrutiny than TLS or WireGuard; (4) it is a small network with few maintainers; (5) malicious nodes can attract traffic and attack routing. If you only want to try it, follow the [cjdns README](https://github.com/cjdelisle/cjdns): build it, run `cjdroute --genconf`, and check peers with `cjdnstool peers show`.

### dn42, Netsukuku and friends: theory

**dn42** (decentralised network 42) is a hobby network where people connect by VPN tunnels, register a private AS number and address block, and run real BGP between themselves. It teaches how the internet's routing works. There is no anonymity at all, because the registry is public. Start at [dn42.dev](https://dn42.dev).

**Netsukuku** was a 2005-era design for a server-free mesh. Nodes group into 'gnodes' so each router stores only a tiny routing table even in a huge network, and a distributed system replaced DNS. The project is abandoned, but its ideas live on in community meshes built with BATMAN-adv and Babel, and in LoRa chat meshes such as Meshtastic. Read it as history, not as a toolkit.

**Honest takeaway:** for a private network between people you know, use Yggdrasil over your own VPS or Reticulum with a TCP link. For off-grid messaging, use Reticulum with LoRa. For anonymity, go back to Part 1.

## Part 5: Social media (video, audio, posts)

The safest post is one the platform cannot read, so this part shows how to encrypt content yourself, move private talk to a server you control, and stop your browser and search engine from profiling you.

### Why: three kinds of exposure \[Beginner\]

1. **Content:** what you post. Fix it by encrypting before you upload.
2. **Metadata:** who you talk to, when, from where, how much. Encryption does not remove this; fewer platforms and fewer contacts do.
3. **Identity:** the link between the account and you (phone number, face, writing style, reused username). Fix it with separate personas (Part 11).

### Do this A: encrypt it yourself with age and minisign \[Beginner to Intermediate\]

**Why age, not PGP.** [age](https://github.com/FiloSottile/age) does one job (encrypt a file to a public key) with modern defaults, no keyring and no web of trust, so there are far fewer ways to misuse it. PGP is large, old and easy to get wrong. [minisign](https://jedisct1.github.io/minisign/) does one job too: sign a file and verify the signature.

**Honest limits.** age does not prove who sent a file, so we add minisign. Neither gives forward secrecy: if your private key leaks, every file ever sent to you can be opened. Metadata (who, when, file size) stays visible. For live chat use a messenger with a ratchet (Signal, OMEMO, Matrix); use age for files, posts, voice notes and videos you hand over asynchronously. PGP remains useful where you must interoperate with email or package signing.

1. **Install.** Debian or Ubuntu: `sudo apt install age minisign`. macOS: `brew install age minisign`.
2. **Make your keys:**

```
age-keygen -o ~/.age-key.txt      # prints your public key (age1...). The file is your SECRET key.
chmod 600 ~/.age-key.txt
minisign -G                       # makes minisign.pub and a password-protected secret key
```

3. **Swap public keys** (the `age1...` line and `minisign.pub`) in person or over a channel you can verify, and read the first and last characters aloud on a call to confirm.
4. **Encrypt for one or more people:**

```
age -r age1<friend> -o clip.mp4.age clip.mp4
age -r age1<friend1> -r age1<friend2> -o note.txt.age note.txt
```

5. **Sign the encrypted file:** `minisign -Sm clip.mp4.age` (creates `clip.mp4.age.minisig`).
6. **Post both files** anywhere that takes files or links. For text-only places add `-a` to `age` to get ASCII text.
7. **Your friend verifies, then decrypts:**

```
minisign -Vm clip.mp4.age -p minisign.pub
age -d -i ~/.age-key.txt -o clip.mp4 clip.mp4.age
```

8. **Group with one shared secret:** `age -p -o file.age file` encrypts with a passphrase instead of keys.

Store your key on a hardware token with [age-plugin-yubikey](https://github.com/str4d/age-plugin-yubikey) if you can.

**Mistakes that break it:** emailing your public key with no verification, keeping the secret key in a synced cloud folder, and posting the plain file name (use a neutral name).

### Do this B: leave the big platforms, run your own \[Intermediate to Advanced\]

Federated tools let you host the social layer. Public posts on them are still public, and your server admin sees everything, so use them for things you want published, not secrets.

| Need | Self-hostable option |
| --- | --- |
| Short posts | [GoToSocial](https://gotosocial.org) (light), Mastodon (heavy) |
| Photos | [Pixelfed](https://pixelfed.org) |
| Video | [PeerTube](https://joinpeertube.org) |
| Audio | [Funkwhale](https://funkwhale.audio) |
| Long posts | [WriteFreely](https://writefreely.org) |

The generic hosting steps (server, TLS, onion access, backups) are in **Part 7**. Do those once and reuse them for every app.

### Do this C: private messaging with XMPP, OMEMO and Tor \[Intermediate to Advanced\]

**Why.** XMPP is an old, open, federated chat standard. **OMEMO** adds end-to-end encryption using a double ratchet (per-message keys). Running the server yourself, reachable over a Tor onion address, hides your server's IP and your clients' IPs from the network.

**Easiest path (no Tor):** [Snikket](https://snikket.org) is a ready-made XMPP server with invite links and apps for phones. Follow its setup guide. It needs a domain name and a small VPS.

**Tor-only path:**

1. On a Debian server install the pieces: `sudo apt install prosody tor`.
2. Create the onion service. Add this to `/etc/tor/torrc`, then `sudo systemctl restart tor`:

```
HiddenServiceDir /var/lib/tor/prosody/
HiddenServicePort 5222 127.0.0.1:5222
HiddenServicePort 5269 127.0.0.1:5269
```

3. Read your address: `sudo cat /var/lib/tor/prosody/hostname`.
4. Install Prosody's Tor module (Prosody 0.12 or later): `sudo prosodyctl install --server=https://modules.prosody.im/rocks/ mod_onions`, per [its page](https://modules.prosody.im/mod_onions).
5. In `/etc/prosody/prosody.cfg.lua` add (replace with your onion hostname):

```
modules_enabled = { "onions" }
onions_tor_all = true
VirtualHost "<yourhash>.onion"
```

6. Create users with `sudo prosodyctl adduser you@<yourhash>.onion`, then restart with `sudo systemctl restart prosody`.
7. **Clients:** [Gajim](https://gajim.org) (desktop, OMEMO built in, SOCKS proxy 127.0.0.1:9050 in account settings), [Conversations](https://conversations.im) (Android, can connect through Orbot), [Dino](https://dino.im) (Linux).
8. **Verify OMEMO fingerprints** with each contact before you trust them. Do not use blind trust for anything important.

The module's own README warns that running a hidden service on a machine that also hosts a normal server can expose the hidden one, so use a dedicated server for an onion-only setup.

### Do this D: switch JavaScript off \[Beginner\]

**Why.** JavaScript is how websites run code on your machine. It powers trackers, fingerprinting (reading fonts, canvas, hardware), and most browser attacks.

| You gain | You lose |
| --- | --- |
| Far fewer trackers and ads | Many web apps (maps, banking pages, editors) stop working |
| Less fingerprinting and fewer drive-by exploits | Some login pages and captchas break |
| Faster pages, less battery drain | Embedded video and audio players may not load |

1. **Tor Browser (best for beginners):** click the shield icon, open Settings, Security, and choose **Safest**. This turns JavaScript off by default for websites.
2. **Normal Firefox with NoScript:** install [NoScript](https://noscript.net) and leave the default (block scripts). Allow a site **temporarily** when you must.
3. **uBlock Origin:** open its settings and turn on *Disable JavaScript* for a global rule, then re-enable per site.

**Trade-off to know:** in Tor Browser, do not mix hand-picked per-site settings. Using the standard levels keeps you in a large crowd; custom setups make you stand out.

### Do this E: private search with SearXNG and 4get \[Intermediate\]

**Why.** SearXNG and 4get are *metasearch* tools: they query other engines on your behalf and strip trackers, so Google or Bing see your server, not you. **Limit:** if you are the only user, queries still link back to your server unless you route them through Tor or share the instance with others.

**SearXNG on your own machine (Docker):**

1. Make a config folder and start the container (the port is bound to localhost):

```
mkdir -p ~/searxng
docker run -d --name searxng -p 127.0.0.1:8888:8080 \
  -v ~/searxng:/etc/searxng \
  -e "BASE_URL=http://localhost:8888/" searxng/searxng:latest
```

2. Stop it, edit `~/searxng/settings.yml`, and set these (see the [settings docs](https://docs.searxng.org)):

```
use_default_settings: true
server:
  secret_key: "<paste the output of: openssl rand -hex 32>"
  limiter: false
  image_proxy: true
  method: "POST"
search:
  autocomplete: ""
outgoing:
  proxies:
    all://:
      - socks5h://<address of your Tor SOCKS port>:9050
```

3. Start it again with `docker start searxng`. Keep the *tracker URL remover* plugin on.
4. Set it as the default search engine in your browser. `method: POST` keeps searches out of browser history but can make some browsers' address-bar search awkward, so add the engine from the SearXNG page's OpenSearch entry.

Inside Docker, `127.0.0.1` means the container itself, so point the Tor proxy line at your Tor container or the host's address. Sending everything through Tor often triggers captchas at big engines; that is the price.

**4get on your own machine (Docker):** the project lives at [git.lolcat.ca/lolcat/4get](https://git.lolcat.ca/lolcat/4get).

```
docker run -d --name 4get -p 127.0.0.1:8080:80 \
  -e FOURGET_SERVER_NAME="localhost" luuul/4get:latest
```

Then open `http://127.0.0.1:8080`. 4get's own settings live in `data/config.php`, a self-documenting file; the project's README also shows how to expose it as a Tor onion service by adding `HiddenServiceDir` to Tor's config. For public-instance lists and comparisons, see the project's about page.

## Part 6: Communication (self-hosted messaging, calls, media)

Picking a private messenger for everyday use is covered in Guides 1 and 2. This part is the step beyond: running your own server so no company sits in the middle. Self-hosting hides your content from companies but not your metadata from the server you rent, and for a lighter private server you can use XMPP with OMEMO (Part 5, Do this C).

### Why Matrix \[Beginner\]

[Matrix](https://matrix.org) with the Element app is federated, supports groups, bridges to other networks, and you can host it yourself. The trade-off is metadata: servers can see who is in which room. Choose it when a team or community wants control of its own communication.

### Do this: self-host Matrix \[Advanced\]

The least painful route is the community [matrix-docker-ansible-deploy](https://github.com/spantaleev/matrix-docker-ansible-deploy) playbook. It installs the Synapse server (the default), a database, the Element web client and the supporting services for you. Recent versions can also install the lighter Rust servers Tuwunel and Continuwuity, but Synapse is the full-featured default.

**You need:** a Debian or Ubuntu server (2 GB RAM is a sensible minimum), a domain name you control, and a computer to run Ansible from.

1. **Point DNS:** create an `A` record for `matrix.<yourdomain>` pointing at the server, and (if you want your short domain on user IDs) set up the `.well-known` delegation the docs describe.
2. **Get the playbook:** `git clone https://github.com/spantaleev/matrix-docker-ansible-deploy.git`, then `cd matrix-docker-ansible-deploy`.
3. **Create your config.** Copy the example inventory as the docs show, then in `inventory/host_vars/matrix.<yourdomain>/vars.yml` set `matrix_domain` to your domain and add the two long random secrets the docs ask for (generate them with `openssl rand -hex 32`).
4. **Harden for privacy** in the same file, using the variable names from the playbook docs: turn **open registration off** (the default), turn **federation off** if your group is closed (this stops your server talking to the wider network), and do not enable analytics or telemetry add-ons.
5. **Install:** follow *Installing* in the [playbook docs](https://github.com/spantaleev/matrix-docker-ansible-deploy/blob/master/docs/configuring-playbook.md). It installs the roles and runs the setup in one command.
6. **Create your first user** with the playbook's user-registration tag, as the docs show, then log in at `https://element.<yourdomain>`.
7. **Turn on end-to-end encryption** for every private room (Element does this by default for new private chats), and **verify devices** with each contact by comparing the emoji or QR code.

**Mistakes that break it:** leaving federation on with an open signup page, forgetting that your server logs IP addresses unless you reduce logging, never updating the server, and skipping device verification.

### Calls and media

- **Calls** need a relay server (TURN) so people behind routers can reach each other. A relay can reveal participants' IPs to the server owner, so host your own, which is what the playbook can do.
- **Voice or video for a group:** see the self-hosted meeting tools in Part 7 (Jitsi, Mumble).

## Part 7: Self-hosting playbook (community, video, music, games, code, mail, meetings)

One reusable recipe covers every service in this part: a hardened small server, a private way to reach it, minimal logs and encrypted backups. Learn the recipe once, then pick apps from the table.

### Why: self-hosting moves trust, it does not remove it \[Beginner\]

Hosting yourself removes ad-tech, data resale and a company reading your files. It adds three duties: you must patch, you must back up, and whoever hosts the machine can still see the disk and memory. A rented server also ties a payment identity to it.

| Where it lives | Good for | Watch out for |
| --- | --- | --- |
| A mini-PC or Raspberry Pi at home | Full control of the hardware | Your ISP may block ports or share IPs (CGNAT); use Tor or WireGuard so you open no ports |
| A rented VPS | Easy, good uplink | The host can inspect the machine; there is a payment trail |
| A colocated box | Control plus a good uplink | Cost and effort |

### Do this: the private hosting recipe \[Intermediate\]

**1. Start from Debian stable.** On a home machine use full-disk encryption (Guide 1) from the first install.

**2. Lock down the basics** (on a new Debian server, as root):

```
adduser admin && usermod -aG sudo admin
# on your own computer: ssh-keygen -t ed25519 && ssh-copy-id admin@<server>
# then on the server:
printf 'PasswordAuthentication no\nPermitRootLogin no\n' > /etc/ssh/sshd_config.d/hardening.conf
apt install -y ufw unattended-upgrades fail2ban
ufw default deny incoming && ufw allow 22/tcp && ufw enable
systemctl restart ssh
```

Test a second SSH login with your key before you close the first one.

**3. Pick an app platform.**

- [YunoHost](https://yunohost.org) installs a web admin panel and over a hundred apps with one command: `curl https://install.yunohost.org | bash`, run on a clean Debian. **Check the install page first for the supported Debian version.** When I looked, YunoHost 12 targeted Debian 12 and Debian 13 support was still in beta, so the script fails on a fresh Debian 13.
- **Docker Compose** is the other route: every app below ships a Compose file in its documentation.

**4. Choose how you reach it.**

- **Tor onion service (no open ports, hides the server's IP).** Add to `/etc/tor/torrc`, restart Tor, then read the address from `hostname`:

```
HiddenServiceDir /var/lib/tor/myapp/
HiddenServicePort 80 127.0.0.1:8080
```

- **WireGuard** (your own devices only). [PiVPN](https://pivpn.io) sets it up with a guided script on a Pi or Debian box, or follow the [WireGuard quickstart](https://www.wireguard.com/quickstart/).
- **Public HTTPS with Let's Encrypt** works, but every certificate name is published in public Certificate Transparency logs, so a subdomain like `secret.example.com` becomes searchable. Use a wildcard certificate, or keep private apps off the public web.

**5. Keep fewer logs.** In `/etc/systemd/journald.conf` set `Storage=volatile` (logs live in RAM), and turn off or anonymise web-server access logs. Disable every analytics or telemetry option in each app.

**6. Back up, encrypted.** Use [restic](https://restic.net):

```
restic init --repo /mnt/usb/backup
restic -r /mnt/usb/backup backup /srv /etc
restic -r /mnt/usb/backup snapshots
```

Keep one copy off-site, and test a restore once.

**7. Update monthly** (`apt update && apt full-upgrade`), subscribe to each app's security announcements, and delete anything you stop using.

**Mistakes that break it:** exposing admin panels to the whole internet, reusing passwords, skipping backups, and running the server and your personal browsing from the same machine.

### Pick your apps

| Category | Self-hostable options | Notes |
| --- | --- | --- |
| Community: forums, groups | [Discourse](https://www.discourse.org) (heavy), [Flarum](https://flarum.org) (light), [Lemmy](https://join-lemmy.org), [Zulip](https://zulip.com) | Closed registration plus invites keeps it private |
| Video: movies, series, broadcasts | [Jellyfin](https://jellyfin.org) (your own library), [PeerTube](https://joinpeertube.org), [Owncast](https://owncast.online) (live streams) | Host only media you have the right to |
| Music, podcasts, audiobooks | [Navidrome](https://www.navidrome.org), [Audiobookshelf](https://www.audiobookshelf.org) | Works with your own files and DRM-free purchases |
| Games | [Luanti](https://www.luanti.org), [Veloren](https://veloren.net), OpenTTD, managed with [LinuxGSM](https://linuxgsm.com) | Open-source or properly licensed servers only. Mobile puzzles: install from [F-Droid](https://f-droid.org) or DRM-free [GOG](https://www.gog.com) and play offline |
| Code and dev work | [Forgejo](https://forgejo.org) (Git hosting), [Woodpecker CI](https://woodpecker-ci.org) | Set a neutral commit email per repository |
| Professional mail | See the mail note below | Email is the hardest service to self-host well |
| Messaging for work | [Zulip](https://zulip.com), Matrix (Part 6) |  |
| Virtual meetings | [Jitsi Meet](https://jitsi.org), [Galene](https://galene.org), [Mumble](https://www.mumble.info) (voice only), [BigBlueButton](https://bigbluebutton.org) (classrooms) | Needs a TURN relay and a few open UDP ports |

**Games: a boundary you should know.** You cannot legally self-host a commercial MMO (the server software is not licensed to you). Running unofficial copies of one breaks copyright and the service terms. Stay with open-source game servers like those above.

**Mail, honestly.** Self-hosting email gives you ownership of your mailbox but not of your metadata: the person you write to still holds the message on their server. It is also hard: you need a static IP with reverse DNS, working SPF, DKIM and DMARC records, an open port 25 (many VPS hosts block it), and a clean IP reputation. Two sane routes:

1. **Own your domain, rent the mailbox:** use a paid, privacy-respecting provider with your own domain, so you can move later without changing your address.
2. **Self-host in one box:** [Mail-in-a-Box](https://mailinabox.email), [Mailcow](https://mailcow.email) or [Stalwart](https://stalw.art) automate most of it.

Email is not end-to-end encrypted by default. Use Thunderbird's built-in OpenPGP only where the other side supports it, or send an age-encrypted attachment (Part 5).

**Jitsi in short:** follow the [Jitsi quick-install guide](https://jitsi.github.io/handbook/docs/devops-guide/devops-guide-quickstart/) on a Debian or Ubuntu server with a domain. It lists the ports to open (80, 443 and 10000/udp among them). Turn on end-to-end encryption in the meeting if your Jitsi version offers it for your clients, and check its current limitations.

## Part 8: Storage (files and media)

Your files stay private when the key to read them never leaves your hands, so this part builds storage in layers: sync without a server, a self-hosted cloud, encryption you control, and backups that survive a disaster.

### Why: who holds the key? \[Beginner\]

- **Server-side encryption:** the server encrypts, so the server owner (the company, or a rented host) also holds the key. It protects against a stolen disk, not against the host.
- **Client-side (end-to-end) encryption:** your device encrypts before upload, and the server only ever sees scrambled data. This is the kind that protects you from the host.

Data reduction first: delete files you no longer need, and keep scans of IDs and passports out of any cloud. Stripping metadata before you share files is covered in Guide 1.

| Need | Tool | Privacy notes |
| --- | --- | --- |
| Sync between your own devices, no server | [Syncthing](https://syncthing.net) | Peer to peer and encrypted in transit; no company in the middle |
| Cloud drive, calendar, contacts | [Nextcloud](https://nextcloud.com) | You control the server, but the server can read files unless you encrypt on the client |
| Photos | [Immich](https://immich.app) | Self-hosted photo library with phone backup |
| Encrypt before upload | [Cryptomator](https://cryptomator.org), [gocryptfs](https://nuetzlich.net/gocryptfs/), age | You hold the key; works on top of any cloud |
| Backups | [restic](https://restic.net) or [BorgBackup](https://www.borgbackup.org) | Encrypted and deduplicated |

### Do this A: Syncthing, the easiest private sync \[Beginner\]

1. Install it: `sudo apt install syncthing` (Linux), or download it from [syncthing.net](https://syncthing.net) for Windows, macOS and Android.
2. Start it. The control page opens at `http://127.0.0.1:8384`. First, set a GUI password under Actions, Settings.
3. On both devices, open Actions, Show ID. Add the other device's ID under *Add Remote Device* and accept the request on the other side.
4. Click *Add Folder*, pick a folder, and tick the device to share it with.
5. **For more privacy:** under Actions, Settings, Connections, untick *Global Discovery* and *Enable Relaying*, then connect devices by direct address or on your local network. If you want an always-on copy on a VPS you do not fully trust, share the folder to it as an *untrusted device* with an encryption password, so the VPS only stores ciphertext.

### Do this B: Nextcloud, your own cloud \[Intermediate to Advanced\]

The official all-in-one (AIO) container installs Nextcloud with its database, cache and backup tooling.

1. Install Docker on a Debian or Ubuntu server and point a domain at it.
2. Run the AIO master container. The [official README](https://github.com/nextcloud/all-in-one) is the source of truth for the current image name; at the time of writing it reads:

```
sudo docker run \
  --sig-proxy=false \
  --name nextcloud-aio-mastercontainer \
  --restart always \
  --publish 80:80 --publish 8080:8080 --publish 8443:8443 \
  --volume nextcloud_aio_mastercontainer:/mnt/docker-aio-config \
  --volume /var/run/docker.sock:/var/run/docker.sock:ro \
  ghcr.io/nextcloud-releases/all-in-one:latest
```

3. Open `https://<server-ip>:8080` (it uses a self-signed certificate), save the passphrase it shows, enter your domain, and start the containers. AIO needs ports 80 and 443 for the certificate unless you put it behind your own reverse proxy; for a private install reached only over WireGuard or Tor, follow its *reverse proxy* documentation instead.
4. **Privacy settings:** in Administration settings, open Sharing and turn off publishing to the global address book and lookup server. Disable apps you do not use. Turn on the built-in backup in the AIO interface and store the backup on another disk.
5. **Add client-side encryption** for anything sensitive: create a Cryptomator vault inside your Nextcloud folder. Nextcloud then stores only the scrambled vault.

### Do this C: encrypt a folder with gocryptfs \[Intermediate\]

```
sudo apt install gocryptfs
mkdir cipher plain
gocryptfs -init cipher         # asks for a password; write down the master key it prints
gocryptfs cipher plain         # mounts the readable view at ./plain
# put files in ./plain; upload or sync only the ./cipher folder
fusermount -u plain            # unmount when you finish
```

Sync or upload the `cipher` folder anywhere. The `plain` folder exists only while mounted.

### Do this D: your photos with Immich \[Intermediate\]

1. Install Docker and Docker Compose on your server.
2. Follow Immich's [install guide](https://immich.app/docs/install/docker-compose): download its `docker-compose.yml` and `example.env` from the latest release (rename the second to `.env`), set a strong database password in `.env`, then run `docker compose up -d`.
3. Install the Immich phone app and point it at your server's address over WireGuard or Tor. Do not expose it to the public internet without strong passwords and updates.

### Do this E: backups that actually save you \[Beginner to Intermediate\]

Follow the **3-2-1 rule**: three copies, on two kinds of media, with one copy off-site and ideally one offline. Encrypt every copy.

1. Use the restic commands from Part 7 (`restic init`, `restic backup`).
2. Run it on a schedule with a systemd timer or cron.
3. **Test a restore** once a quarter. An untested backup is a guess.
4. Keep one copy disconnected from the network, because ransomware encrypts anything it can reach.

**Mistakes that break it:** trusting server-side encryption as if it were end-to-end, losing the Cryptomator or gocryptfs master key, backing up only to a disk plugged into the same machine, and exposing the admin interface to the whole internet.

## Part 9: Running your own payment nodes

A node you run yourself verifies your own payments, so no third party learns which addresses you own; it does not make the payments anonymous, and this part stays strictly on the lawful-privacy side of that line. Paying with cash, buying coins, and the legal limits around mixers and cross-chain swaps are covered in Guides 1 and 2, so this part only covers the nodes.

### Why: what a node does and does not give you \[Beginner\]

- **Without your own node,** your wallet asks someone else's server about your balance. That server learns your IP and which addresses are yours.
- **With your own node,** your wallet asks your machine. Nobody else sees the questions.
- **What a node does not do:** hide your transactions. Bitcoin and Ethereum are public ledgers. Monero hides amounts and addresses by design, but your identity can still leak through exchanges, shipping addresses and your own habits.

**Check the law first.** Crypto is restricted or taxed in many countries, and running a node does not change your reporting duties.

### Do this A: a Bitcoin node \[Intermediate\]

**Easiest:** use a node distribution that installs everything with a web interface, such as [StartOS](https://start9.com), [Umbrel](https://umbrel.com) or [RaspiBlitz](https://github.com/raspiblitz/raspiblitz). Plan for a 1 TB SSD and 4 to 8 GB RAM for a full node (check the project page for current blockchain size), or use pruning to save space.

**Manual, on Linux:**

1. Download Bitcoin Core from [bitcoincore.org](https://bitcoincore.org/en/download/) and **verify the signatures** as that page explains.
2. Create `~/.bitcoin/bitcoin.conf`:

```
prune=10000        # keep about 10 GB of blocks; remove this line for a full node
proxy=127.0.0.1:9050   # send outgoing connections through Tor
listen=1
```

3. Start it: `bitcoind -daemon`, and check progress with `bitcoin-cli getblockchaininfo`. The first sync takes many hours or days.
4. Install a wallet such as [Sparrow](https://sparrowwallet.com) and point it at your node (Server settings, Bitcoin Core).
5. Read Bitcoin Core's own [Tor guide](https://github.com/bitcoin/bitcoin/blob/master/doc/tor.md) before you restrict the node to Tor only; doing that without understanding it can leave you with few peers.

**Accepting payments yourself:** [BTCPay Server](https://btcpayserver.org) is a self-hosted payment processor that talks to your node, so no third party processor sits between you and your customers.

### Do this B: a Monero node \[Intermediate\]

Monero's official guides are at [docs.getmonero.org](https://docs.getmonero.org/running-node/monerod-tori2p/).

1. Download the CLI from [getmonero.org/downloads](https://www.getmonero.org/downloads/) and **verify the hashes** shown there.
2. Start a pruned node (it keeps roughly one eighth of the history but still verifies everything):

```
./monerod --data-dir /mnt/ssd/monero --prune-blockchain
```

3. Let it finish syncing (check with `./monerod status`).
4. **Add Tor for sending.** Run Tor, then add these lines to `monerod.conf` so transactions from your wallets are broadcast over Tor:

```
prune-blockchain=1
tx-proxy=tor,127.0.0.1:9050,16,disable_noise
```

5. Point your wallet at it: in the GUI wallet or Feather, set *Node* to your local node (or its onion address, if you published one with `anonymous-inbound`).

Monero's docs note that the blockchain itself cannot sync over onion or I2P services, only through a SOCKS proxy; that is why `tx-proxy` handles *sending* rather than syncing.

### Do this C: an Ethereum node \[Advanced\]

Ethereum needs two programs running together: an **execution client** (for example [Geth](https://geth.ethereum.org)) and a **consensus client** (for example [Lighthouse](https://lighthouse-book.sigmaprime.io) or Nimbus). Plan for a 2 TB NVMe SSD and 16 GB RAM.

1. **Easiest route:** [Eth Docker](https://ethdocker.com) wraps everything: `git clone` its repository, then run `./ethd install`, `./ethd config`, and `./ethd up`.
2. **Manual route:** follow ethereum.org's [run a node guide](https://ethereum.org/en/developers/docs/nodes-and-clients/run-a-node/).
3. In your wallet, add your node as a custom RPC network (usually `http://<node-ip>:8545`) over your WireGuard or Tor connection. **Never expose port 8545 to the internet.**

Why bother: the default wallet setup sends every address query to a company's servers. Your own node removes that.

### Mistakes that break it

- Trusting downloaded binaries without checking signatures.
- Opening RPC ports to the world.
- Reusing addresses, or paying from an exchange account tied to your ID and expecting the node to hide it.
- Forgetting disk space and having the node stop mid-sync.

## Part 10: Hardware security

Software privacy fails if the hardware underneath it is leaky, so this part covers kill switches, signal blocking, location privacy, safer firmware, and the hidden fingerprint your camera leaves on every photo.

### 10.1 Hardware kill switches and tamper defences \[Beginner to Intermediate\]

**Why.** A kill switch is a way to make a device stop working, lock, or forget its secrets *fast* and without relying on software you might not control.

**What to use:**

1. **Physical privacy switches on the device.** Some laptops ship with real hardware switches for the microphone, camera and Wi-Fi, for example Purism's Librem line and the Framework Laptop's camera and microphone switches. Prefer these over software toggles.
2. **Dead-man cable: BusKill.** [BusKill](https://www.buskill.in) is a magnetic cable that links your laptop to your belt or wrist; if the cable is pulled out, the BusKill app locks the screen or shuts the machine down (you choose). Steps: get the cable, install the BusKill app from its site, plug the cable in, pick the trigger (lock or shutdown) in the app, and clip the far end to yourself. Test that it triggers.
3. **Your encryption key is your real kill switch.** With full-disk encryption set up as in Guide 1, a powered-off laptop is unreadable. For an emergency wipe, `sudo cryptsetup luksErase /dev/<your encrypted partition>` destroys every key slot; the data becomes unrecoverable. **This is instant and irreversible. Do not run it as a test.** Debian-family systems can also use the `cryptsetup-nuke-password` package so that a special password at boot erases the keys; check that it is available for your release.
4. **Block unknown USB devices.** Install [USBGuard](https://usbguard.github.io): `sudo apt install usbguard`, then `sudo usbguard generate-policy | sudo tee /etc/usbguard/rules.conf` with only your own devices plugged in, and `sudo systemctl enable --now usbguard`. Unknown USB devices are now blocked.
5. **Tamper evidence.** Seal screws or ports with glitter nail polish and photograph the pattern. A different pattern next time means someone opened the machine. It does not stop an attack; it tells you one may have happened.
6. **Smaller wins:** a USB data-blocker when charging from public ports, and a mic blocker plug on desktops. (Camera covers are in Guide 1.)
7. **Software 'USB kill' tools** such as [usbkill](https://github.com/hephaest0s/usbkill) shut a computer down when USB devices change. Treat them as an extra, and test them first.

**Theory-only zone: USB Killer devices.** A 'USB Killer' charges capacitors from a USB port and fires high-voltage pulses back into it, which physically destroys the host. I do not give build or use steps. It is a way to damage someone else's equipment, which is criminal damage in nearly every legal system, and it protects nothing: encryption plus a key wipe makes your data unreadable *without* destroying hardware. Protect your own gear with the steps above instead.

### 10.2 Blocking tower, Wi-Fi and GPS signals \[Beginner\]

**Why.** GPS only *receives*, so a satellite cannot track you. Location leaks because your phone uses cell towers, Wi-Fi and Bluetooth scans, and apps that upload the position. A **Faraday bag** is a signal-blocking pouch that stops all of those radios.

1. Buy a bag with several metallic layers and a roll-top closure from a reputable maker (cheap ones often fail).
2. **Test before you trust it:** put your phone inside, close it fully (fold the top at least three times), then call that phone from another phone. It must not ring. Repeat with Bluetooth (try to connect it to a speaker) and Wi-Fi.
3. Re-test every few months and replace a bag that is creased, torn or leaking.
4. **Know the trade-off:** the moment a phone vanishes from the network at a location is itself a data point. Switch off before you reach the sensitive place, not at the door, and remember that powering off fully is stronger than any bag.
5. Airplane mode is not a guarantee. Apps and the hardware itself may still act.

### 10.3 Keeping your home location private \[Beginner\]

**Why.** Your address leaks through deliveries, accounts, public records (property, voter and company registries) and data brokers.

1. **List where your address appears.** Search your name, check your accounts, and note every company that ships to you.
2. **Get a different mailing address.** Options: a private mailbox from a commercial mail receiving agency (in the US it has a street-style address; the provider must verify your ID with USPS Form 1583), a post office box, or a mail-forwarding service. Many couriers will not deliver to a P.O. box, so a private mailbox with a street address is more flexible.
3. **Use pickup points.** Ship online orders to parcel lockers, carrier pickup points or a shop that holds packages (Amazon Hub Locker, UPS Access Point and similar). You collect the parcel, and the seller never gets your home address.
4. **Move your address, service by service.** Change delivery addresses first for shopping and subscriptions. Banks, tax, employers and government agencies may legally require your real residential address; give it to them honestly and keep it out of everything else.
5. **Remove yourself from data brokers** using the deletion steps in Guide 1, and repeat every six months because listings come back.
6. **If you face stalking or abuse,** ask whether your region has an address-confidentiality programme. Many do, and they are built for this.
7. **Property.** Some places let you hold a home through a trust or company to keep your name off public records. Rules differ a lot; ask a local lawyer.

**Mistakes that break it:** posting photos that show your street or house number, giving the private mailbox to some services and the home address to others in a way that lets them be matched, and giving false addresses to banks or the authorities (that is fraud, not privacy).

### 10.4 Safer boot: Libreboot \[Advanced\]

**Why.** Firmware (BIOS or UEFI) runs before your OS and can hide malware that a reinstall never removes. [Libreboot](https://libreboot.org) replaces the vendor firmware with open-source firmware (coreboot-based) you can inspect.

**Be realistic:** Libreboot supports a specific list of mostly older machines, and flashing can brick a computer. Its own docs warn it is risky and tell you to have external flashing equipment ready.

1. **Check support** on the [hardware compatibility list](https://libreboot.org/docs/hcl/). If your machine is not on it, stop.
2. **Easiest:** buy a machine with Libreboot pre-installed, or send yours to Minifree, which offers a [pre-install service](https://libreboot.org/docs/install/).
3. **DIY, using an external programmer:** a Raspberry Pi Pico running Libreboot's *serprog* firmware, a clip for your flash chip, and a Linux computer. Follow the [SPI flashing guide](https://libreboot.org/docs/install/spi.html) line by line, because the exact commands depend on your chip and wiring.
4. **Back up first:** dump the original flash at least three times with `flashprog`, then compare them (`sha1sum dump1.rom dump2.rom dump3.rom`). They must be identical. Keep one copy safe.
5. **Flash, verify and boot,** as the guide says.
6. **Alternatives:** machines sold with open firmware and measured boot, such as Heads on Nitrokey or Purism hardware, if you do not want to flash anything yourself.

### 10.5 Hardening the firmware: me\_cleaner \[Advanced\]

**Why.** On most Intel computers from about 2008 on, the **Management Engine (ME)** is a small hidden computer inside the chipset that runs before and beside your OS. It has had serious security flaws and cannot normally be turned off. [me\_cleaner](https://github.com/corna/me_cleaner) removes most of its code from the firmware image and, on newer chips, sets a flag (HAP) that asks it to stay idle after boot.

**Limits:** it does not work on every chip (firmware protections called Boot Guard can stop a modified image from booting), the project is no longer actively developed, AMD's equivalent (PSP) cannot be cleaned this way, and a mistake can brick the machine. Libreboot's installation docs already include ME handling for supported Intel systems, so use that route if you can.

1. Dump your firmware with an external programmer, three times, and compare the checksums, as in 10.4.
2. Download me\_cleaner and read its README for your chip generation.
3. Run it on a *copy*: `python3 me_cleaner.py -S -O cleaned.rom dump1.rom`. Then verify the result with `python3 me_cleaner.py --check cleaned.rom`.
4. Flash the cleaned image with your external programmer and test-boot.
5. Keep your original dump forever, so you can go back.

### 10.6 Cameras have invisible scratches (PRNU) \[Intermediate\]

**Why.** No sensor is perfect. Each pixel reacts to light slightly differently, and that pattern, called **PRNU** (photo-response non-uniformity), is baked into every photo. Researchers can estimate a camera's pattern from many of its photos and later test whether a different photo came from the same camera. Removing EXIF data does not remove it. This is about knowing the concept: it can link an anonymous picture back to the device that took it.

**What it does and does not do.** An investigator needs a reference: photos from your camera, or the camera itself. It matches *devices*, not people. Heavy resizing, cropping and compression weaken it, and no method removes it perfectly without hurting quality.

**Do this (for protest photos, source protection, or any picture you must publish without linking it to your camera).** Strip the visible metadata first, as in Guide 1; PRNU survives that step.

1. **Disturb the pixel-exact noise pattern.** Crop a few pixels off different edges, shrink the image by at least a third, and re-save at moderate JPEG quality, for example with ImageMagick:

```
magick photo.jpg -strip -gravity NorthWest -chop 9x6 -gravity SouthEast -chop 5x8 \
  -resize 65% -quality 85 out.jpg
```

2. **Use a different crop offset each time.** The same crop on many photos can itself become a pattern.
3. **Tools that automate this** include [StripShot](https://github.com/tailsmails/stripshot), which crops randomly, resizes and adds noise for the same goal, and the research tool [DIPPAS](https://github.com/polimi-ispl/dip_prnu_anonymizer). I have not tested their results independently, so treat their claims as unverified.
4. **Best defence: separate devices.** Use a cheap camera that has never been tied to your identity for anything you publish under another persona, and never take personal photos with it. Do not mix the two, ever.

**Other fingerprints survive:** JPEG settings, lens dust and distortion, shadows and scenery, faces and number plates, and the time you post. Check the *content* of the picture too.

**Theory-only zone: planting a chosen camera's fingerprint.** Researchers have shown that a PRNU pattern can be copied from one camera onto a photo from another so that it falsely points at the first camera. I give no steps. The only use is to frame a specific person or device. A random crop and resize that breaks matching is privacy; forging someone else's signature is not.

## Part 11: Beyond the digital world, IRL compartmentalization

Compartmentalization means keeping parts of your life from being linked, so that one leak cannot unravel everything, and it works only if you stop sharing identifiers, devices, places and people between the parts. The digital side (separate email aliases, password databases, VMs, Qubes, separate phones and payment methods) is built in Guides 1 and 2, and writing-style and face matching are handled in Guides 2 and 4. This part covers what those guides cannot: places, time, paper and people.

### Why: how people get linked \[Beginner\]

Someone who wants to connect two parts of your life looks for **shared things**:

- **Identifiers:** names, emails, phone numbers, usernames, faces, voices.
- **Infrastructure:** the same device, Wi-Fi network, payment card, mailing address, or account recovery email.
- **Behaviour:** the same hours, places, phrases and writing style.
- **People:** a mutual friend who knows both sides.

The rule that follows: **nothing is shared between compartments, and when something must cross, it crosses one way, deliberately, and is never reused.**

### Do this: keep your compartments apart in the real world \[Beginner to Intermediate\]

**Step 1: decide how many compartments you really need.** Most people need two or three. More means more mistakes. Build each compartment's digital kit (its own email alias, password database, device or VM, and payment method) with Guides 1 and 2; the steps below keep the compartments apart in the physical world.

| Compartment | Used for | Identity |
| --- | --- | --- |
| A: Legal | Bank, tax, employer, family, government | Your real, legal name and address |
| B: Everyday online | Hobbies, shopping, ordinary social media | Your usual name or a consistent alias |
| C: Sensitive | Activism, journalism sources, support for a personal crisis | A pen name where the law allows |

**Step 2: keep compartments apart in time and space.**

- Do not use C's devices at home or at work, and do not switch them on next to A's devices. Location and timing link them even when the content does not.
- Do not use C at the same hours as your normal routine.
- Never mention one compartment's details in another, not even as a joke.

**Step 3: watch the paper world.**

- Many colour laser printers print tiny tracking dots that encode the printer's serial number and the date. See the [EFF's list of printers that do or do not](https://www.eff.org/pages/list-printers-which-do-or-do-not-print-tracking-dots). Avoid printing anything you need to stay unlinked.
- Shred papers with your address. Take the labels off parcels before you recycle them.

**Step 4: protect people, too.**

- Tell family and friends what not to share about you, and agree a family code word to verify urgent phone calls.
- Verify strangers who call you by phoning back on a number you find yourself, never one they give you.

### When a link happens (and it will)

1. **Assume the compartment is burned.** Do not try to patch it quietly.
2. Stop using the linked accounts and devices. Delete or retire them.
3. Rebuild with fresh everything; reusing anything recreates the link.
4. Tell people who may be affected.

A written checklist before each session helps: *which identity, device, network, place and payment method am I using, and is any of it shared with another compartment?*

### Staying legal

Pen names, separate emails, private mailboxes and separate devices are normal and lawful. Banks, tax offices, employers, and immigration authorities need your true identity, and you should give it. Compartmentalization is about keeping *strangers* from connecting your life, not about lying to the people the law says you must be truthful with.

**Theory-only zone: false identities.** I do not give steps for forging identity documents, impersonating a real person, building a fake 'paper trail', or hiding from a court order, creditors or an investigation. Identity papers are checked against official registries, forgeries fail often, and the penalties are severe. The concept matters only so you understand why the safe alternatives exist. If you need real protection from a violent person or a stalker, the lawful routes are far stronger: address-confidentiality programmes, protective orders, a legal name change where courts allow it, and specialist support services (a domestic-violence or stalking helpline in your region can help you plan safely, including digital safety).

## Appendix: start-here path, link index, and corrections to your notes

If you read only one page, read this one: it gives the order to do things in, every main link in one place, and a plain list of what I changed or could not verify.

### Start here: the order of effort

| Step | Do this | Part |
| --- | --- | --- |
| 1 | Confirm Guides 1 and 2 are done, and write down who you are hiding from | Read this first |
| 2 | Make a Tails USB stick, check it with check.torproject.org, and run one test session | 1 |
| 3 | Add bridges if your network is hostile; set up Whonix for long-running identities | 1 |
| 4 | Keep your compartments apart in the real world: places, times, paper, people | 11 |
| 5 | Encrypt what you post yourself (age and minisign); set the browser to block JavaScript by default | 5 |
| 6 | Self-host one service, messaging first | 5, 6, 7 |
| 7 | Set up encrypted backups and test a restore | 7, 8 |
| 8 | Move parcels to a private mailbox or pickup point | 10 |
| 9 | Only then add the heavy options your threat model needs: I2P, Yggdrasil or Reticulum, Faraday bags, Libreboot, me\_cleaner | 2, 4, 10 |

**Every month:** update everything, review installed apps and saved Wi-Fi networks, and test one backup.

### What I changed or could not verify in your notes

| What you wrote | What I did |
| --- | --- |
| searchNXG | Read as **SearXNG** |
| hyperbola | Read as **Hyperboria** (the main cjdns network) |
| dns42 | Read as **dn42** |
| cRXN | Searched; found no project by that name. Left out; send a link and I will add it |
| PNUwiper | Searched; found no tool by that name. Covered StripShot, DIPPAS and manual steps instead |
| I2P completely eliminates the sniper problem | Corrected: I2P is differently shaped, not immune (Part 2) |
| PGP is outdated, use age and minisign | Covered, with the limits: no sender authentication in age, no forward secrecy (Part 5) |
| Killer USBs | Covered BusKill, USBGuard and key wipes; USB Killer hardware is theory-only (Part 10) |
| Replace the camera fingerprint with a fake one | Covered only breaking the match; planting a chosen camera's fingerprint is theory-only (Part 10) |
| Self-host everything, including games and mail | Added the limits: no unlicensed MMO servers, and email is the hardest service (Part 7) |
| Anything already owned by Guides 1 and 2 | Removed from this guide and replaced with a one-line pointer where the flow needed it |

**Theory-only zones (four):** Wi-Fi you may not use (3), USB Killer devices (10), planting a camera fingerprint (10), false identities (11).

### How far to trust this guide

- **Checked this session by searching current official or project pages:** Tails and Whonix versions and Whonix physical isolation, Tor bridge setup with lyrebird, Tor's `pf-divert` option, I2P ports and install routes, Reticulum, Yggdrasil, cjdns, Prosody's onion module, 4get, Matrix playbook, YunoHost's Debian requirement, Nextcloud AIO, Monero Tor options, Libreboot docs.
- **Written from my own knowledge of the tools, not copied from a page I opened:** the age, minisign, cryptsetup, USBGuard, restic, gocryptfs, ImageMagick, OpenBSD pf rules, and systemd and SSH commands. Run each tool with `--help` the first time and compare with its manual.
- **Not tested on real hardware:** the OpenBSD Tor router, the Libreboot and me\_cleaner flashing steps, and the LoRa setup. Do those on spare equipment first.
- **I did not open every link.** If one fails, search the project name rather than a mirror.

### Link index

| Topic | Official page |
| --- | --- |
| Tor Browser and Tor Project | [torproject.org](https://www.torproject.org) |
| Get Tor bridges | [bridges.torproject.org](https://bridges.torproject.org) |
| Run a Tor bridge | [community.torproject.org/relay/setup/bridge](https://community.torproject.org/relay/setup/bridge/) |
| Tails | [tails.net](https://tails.net) |
| Whonix | [whonix.org](https://www.whonix.org) |
| Whonix physical isolation | [whonix.org/wiki/Physical\_Isolation](https://www.whonix.org/wiki/Physical_Isolation) |
| I2P | [geti2p.net](https://geti2p.net) |
| Reticulum | [reticulum.network](https://reticulum.network) |
| Sideband app | [github.com/markqvist/Sideband](https://github.com/markqvist/Sideband) |
| Yggdrasil | [yggdrasil-network.github.io](https://yggdrasil-network.github.io) |
| cjdns | [github.com/cjdelisle/cjdns](https://github.com/cjdelisle/cjdns) |
| dn42 | [dn42.dev](https://dn42.dev) |
| age | [github.com/FiloSottile/age](https://github.com/FiloSottile/age) |
| minisign | [jedisct1.github.io/minisign](https://jedisct1.github.io/minisign/) |
| Prosody (XMPP) | [prosody.im](https://prosody.im) |
| Snikket (easy XMPP) | [snikket.org](https://snikket.org) |
| SearXNG docs | [docs.searxng.org](https://docs.searxng.org) |
| 4get | [git.lolcat.ca/lolcat/4get](https://git.lolcat.ca/lolcat/4get) |
| NoScript | [noscript.net](https://noscript.net) |
| Matrix server playbook | [github.com/spantaleev/matrix-docker-ansible-deploy](https://github.com/spantaleev/matrix-docker-ansible-deploy) |
| YunoHost | [yunohost.org](https://yunohost.org) |
| Forgejo | [forgejo.org](https://forgejo.org) |
| Jellyfin | [jellyfin.org](https://jellyfin.org) |
| Syncthing | [syncthing.net](https://syncthing.net) |
| Nextcloud All-in-One | [github.com/nextcloud/all-in-one](https://github.com/nextcloud/all-in-one) |
| Immich | [immich.app](https://immich.app) |
| restic | [restic.net](https://restic.net) |
| Bitcoin Core | [bitcoincore.org](https://bitcoincore.org) |
| Monero node docs | [docs.getmonero.org](https://docs.getmonero.org/running-node/monerod-tori2p/) |
| Eth Docker | [ethdocker.com](https://ethdocker.com) |
| BusKill | [buskill.in](https://www.buskill.in) |
| USBGuard | [usbguard.github.io](https://usbguard.github.io) |
| Libreboot | [libreboot.org](https://libreboot.org) |
| me\_cleaner | [github.com/corna/me\_cleaner](https://github.com/corna/me_cleaner) |
