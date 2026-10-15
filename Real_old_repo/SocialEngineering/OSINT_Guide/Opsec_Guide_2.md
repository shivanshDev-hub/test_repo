# OPSEC Guide 2: Data Reduction (Hardened Field Book), deduplicated

Oct 10, 2026 · @tiger

## Start here: how this book works

This book shrinks the amount of data you leak every day, using free or cheap tools you set up once. It is for privacy and safety (stalkers, data brokers, hostile networks, doxxing), not for hiding crime. Every section follows the same shape: **Why** (the idea, in plain words), **Do** (copy-paste steps), **Check** (how to prove it worked).

**How this edition is organised.** Guide 1 covers the everyday basics and Guide 3 covers Tor and self-hosting. Where a topic belongs to one of them, this guide keeps a one-line pointer instead of repeating the steps, and the appendix lists every cut. Section numbers are unchanged, so existing cross-references still work.

### Pick your threat model first

A threat model is just the answer to "who am I protecting myself from?" The answer decides how much inconvenience is worth it.

| Level | Who you are hiding from | What helps most |
| --- | --- | --- |
| 1. Trackers | Advertisers, data brokers, app analytics | Ad blocking, privacy frontends, alternative apps, GrapheneOS |
| 2. Nosy networks and people | ISPs, public Wi-Fi, stalkers, doxxers, employers | VPS or obfuscated VPN, kill switch, hardware keys, Session, separate identities |
| 3. Well-funded adversaries | Hostile states, targeted attackers | Qubes OS, strict compartments, physical security. Seek expert advice; a book cannot cover this alone |

Most readers need levels 1 and 2. Level 3 tools are included, but they cost real convenience.

### Five rules that apply everywhere

1. **Do one section at a time.** Finish it, test it, then move on. Half-built setups leak.
2. **One identity per compartment.** Never log into your real accounts from your private browser, VM or phone.
3. **Fix the biggest leak first.** A kill switch matters more than a fancy messenger.
4. **Write down recovery keys on paper.** Lost passphrases (Qubes, GrapheneOS, password managers) cannot be reset.
5. **Updates are security.** Every tool below is worthless if it runs old software.

### Where the ethical line sits

All defensive setups (VPS, obfuscation, GrapheneOS, Qubes, encrypted chat, secure erasure) have full step-by-step instructions. Three areas are **theory only**, because step-by-step help there mostly serves hiding illegal money flows rather than protecting privacy: coin mixers and joiners, Bitcoin to Monero to Bitcoin hopping, and routing payments through your own Monero node (all in the Payments section). Those sections explain the idea and the legal risk, not a how-to.

This is general information, not legal advice. Laws on VPNs, encryption and crypto differ by country, so check yours. Software changes fast: this book is current as of 10 October 2026, and every step links to the official docs, which win if they disagree with this text.

## 1. IP protection: VPS and the kill switch

### Why (theory)

Your IP address is the return address on every packet you send. Websites, trackers and your ISP use it to link your visits into one profile. Hiding it means sending your traffic out through someone else's address. The real question is **whom you trust instead**.

Guide 1 covers choosing a commercial VPN, and Guide 3 covers Tor. This section is about the step between them: running your own server.

|  | Commercial VPN (Guide 1) | Your own VPS |
| --- | --- | --- |
| Who sees your traffic | The VPN company | You (and the hosting company's network) |
| Crowd to hide in | Thousands share each IP | Only you use the IP |
| Link to your identity | Account and payment | Account and payment (usually ID or card) |
| Resists blocking | Easy to detect and block | Strong, if obfuscated (next sections) |
| Effort | 5 minutes | 30 to 60 minutes |

**Developer opinion:** a VPS is better when you want full control, cannot trust any VPN company's "no logs" promise, or live where VPNs are blocked. A reputable VPN is better when you want to hide *among a crowd*. A VPS gives you a quiet, fast, private exit, but it is **not anonymous**: the host knows who pays for it, and a lone user on a static IP is easy to recognise over time. Treat a VPS as "my own trusted pipe", not as a cloak.

### Do: build a hardened VPS

1. **Pick a host.** You need KVM virtualisation, 1 vCPU, 1 GB RAM, 1 IPv4 address and a location near you. Examples (not endorsements; check current payment options and ID rules yourself): [Hetzner Cloud](https://www.hetzner.com/cloud), [1984 Hosting](https://1984hosting.com), [Njalla](https://njal.la), [DigitalOcean](https://www.digitalocean.com), [Vultr](https://www.vultr.com). Many hosts demand ID or a card for new accounts. Choose Debian stable or the latest Ubuntu LTS.
2. **Make an SSH key on your own computer** (never use passwords):

```bash
ssh-keygen -t ed25519 -C "vps-key"
ssh-copy-id root@YOUR_SERVER_IP
```

3. **Log in, update, and add a normal user:**

```bash
ssh root@YOUR_SERVER_IP
apt update && apt full-upgrade -y
adduser admin && usermod -aG sudo admin
rsync --archive --chown=admin:admin ~/.ssh /home/admin
```

4. **Turn off password and root login.** Edit `/etc/ssh/sshd_config.d/99-hardening.conf` and add:

```
PasswordAuthentication no
PermitRootLogin no
```

Then run `systemctl restart ssh`. **Keep your current session open and test a second login in a new window** before you close it.

**Step 5. Firewall and automatic security updates:**

```bash
sudo apt install -y ufw unattended-upgrades
sudo ufw allow OpenSSH
sudo ufw enable
sudo dpkg-reconfigure -plow unattended-upgrades
```

You will open the tunnel port (443) in the next sections. **Check:** `sudo ufw status` shows only the ports you chose.

### Do: a kill switch for all traffic

A **kill switch** means that if the tunnel drops, your traffic stops instead of silently leaving through your real IP. Without one, a VPN is a fair-weather friend.

- **VPN apps (all platforms):** turn on the app's built-in switch. Guide 1 shows where it sits in the main apps.
- **Android and GrapheneOS:** *Settings, Network and internet, VPN*, tap the gear next to your VPN, enable **Always-on VPN** and **Block connections without VPN**.
- **iOS:** there is no system-level kill switch. Use your VPN app's own option and treat iOS as less leak-proof.
- **Windows:** use the VPN app's kill switch. Hand-written firewall rules on Windows can lock beginners out of the network.
- **Linux (works with WireGuard, sing-box or Hysteria in TUN mode):** block everything except the tunnel.

```bash
# Replace wg0 with your tunnel name (wg0, tun0, ...) and VPS_IP with your server's address
sudo ufw default deny incoming
sudo ufw default deny outgoing
sudo ufw allow out on lo
sudo ufw allow out on wg0
sudo ufw allow out to VPS_IP port 443     # lets the tunnel itself connect
sudo ufw allow out 67:68/udp              # DHCP, so Wi-Fi keeps working
sudo ufw enable
```

Also stop IPv6 leaks if your tunnel does not carry IPv6:

```bash
echo 'net.ipv6.conf.all.disable_ipv6=1' | sudo tee /etc/sysctl.d/99-noipv6.conf && sudo sysctl --system
```

### Check

1. Connect, then open [Mullvad's connection check](https://mullvad.net/en/check) and [BrowserLeaks](https://browserleaks.com/ip): the IP must be your VPS or VPN, with no DNS or WebRTC leak.
2. **Pull the plug test:** while a page is loading, stop the tunnel (disconnect the app or `sudo systemctl stop xray`). The page must fail. If it still loads, your kill switch is not working.

## 2. VPN obfuscation: VLESS + REALITY (XTLS-Vision + uTLS)

### Why (theory)

Classic VPN protocols (OpenVPN, WireGuard, IPsec) have a recognisable "shape". Networks that use deep packet inspection (DPI) can spot that shape and block or log it. **Obfuscation** makes your tunnel look like ordinary web browsing instead.

Three pieces do the work in this setup:

- **VLESS** is the lightweight protocol that carries your traffic and checks who you are by UUID.
- **REALITY** makes your server answer like a real, popular website (say, microsoft.com). Anyone probing the server without your key is quietly passed through to the real site. You need no domain name and no certificate, which also means fewer identity links to you.
- **uTLS** makes your client's first message (the TLS "Client Hello") look exactly like Chrome's, not like a proxy program's.

**Why TLS-in-TLS is bad from a privacy view.** When you open an HTTPS site through a proxy that is itself wrapped in TLS, you get encryption inside encryption. The two nested handshakes leave a recognisable pattern of packet sizes and timing, and researchers have shown that this pattern can be detected and used to flag proxy users (for example, Xue et al., USENIX Security 2024). It is a fingerprint that tags *you*. **XTLS-Vision** fixes this by padding the handshake and then passing the inner encryption straight through, so the nested pattern never appears. The REALITY project's own advice is not to pair REALITY with other protocols for this reason. Avoid plain TLS-wrapped proxies (Trojan, VMess over TLS, VLESS over TLS without Vision) when privacy is the goal.

### Do: server setup

You need the hardened VPS from section 1, logged in as your admin user.

**Step 1. Install Xray-core** with the [official installer](https://github.com/XTLS/Xray-install):

```bash
sudo bash -c "$(curl -L https://github.com/XTLS/Xray-install/raw/main/install-release.sh)" @ install
xray version
```

The config lives at `/usr/local/etc/xray/config.json`.

**Step 2. Make your three secrets.** Save the output in a password manager.

```bash
xray uuid               # your user ID
xray x25519             # private key (stays on the server) + public key (goes to clients)
openssl rand -hex 8     # a short ID
```

The labels printed by `xray x25519` differ between versions. The public key is sometimes shown as "Password". Never share the private key.

**Step 3. Pick a camouflage site.** It must be a big, popular site that supports TLS 1.3 and HTTP/2, ideally hosted in the same region as your VPS. Test a candidate:

```bash
curl -sI --tlsv1.3 --http2 https://www.microsoft.com | head -n 1
```

A reply that starts with `HTTP/2` means it qualifies. Do not choose a site you cannot reach from the VPS.

**Step 4. Write the config.** Open `sudo nano /usr/local/etc/xray/config.json` and paste this, replacing the three `PASTE_` values:

```json
{
  "log": { "loglevel": "warning" },
  "inbounds": [
    {
      "listen": "0.0.0.0",
      "port": 443,
      "protocol": "vless",
      "settings": {
        "clients": [
          { "id": "PASTE_UUID", "flow": "xtls-rprx-vision" }
        ],
        "decryption": "none"
      },
      "streamSettings": {
        "network": "tcp",
        "security": "reality",
        "realitySettings": {
          "show": false,
          "target": "www.microsoft.com:443",
          "xver": 0,
          "serverNames": ["www.microsoft.com"],
          "privateKey": "PASTE_PRIVATE_KEY",
          "shortIds": ["PASTE_SHORT_ID"]
        }
      }
    }
  ],
  "outbounds": [
    { "protocol": "freedom", "tag": "direct" },
    { "protocol": "blackhole", "tag": "block" }
  ],
  "routing": {
    "rules": [
      { "type": "field", "ip": ["geoip:private"], "outboundTag": "block" }
    ]
  }
}
```

The `routing` block stops anyone who steals your key from using your server to reach private networks. If an older Xray version rejects `target`, rename it to `dest`.

**Step 5. Test, open the port, start.**

```bash
sudo xray run -test -config /usr/local/etc/xray/config.json
sudo ufw allow 443/tcp
sudo systemctl enable --now xray
sudo systemctl status xray
```

Always run the `-test` line after any edit. Optional speed boost (BBR congestion control):

```bash
echo -e "net.core.default_qdisc=fq\nnet.ipv4.tcp_congestion_control=bbr" | sudo tee /etc/sysctl.d/99-bbr.conf && sudo sysctl --system
```

### Do: client setup

Build a share link from your values and paste it into a client app (the apps can import it from the clipboard):

```
vless://PASTE_UUID@YOUR_SERVER_IP:443?encryption=none&flow=xtls-rprx-vision&security=reality&sni=www.microsoft.com&fp=chrome&pbk=PUBLIC_KEY&sid=PASTE_SHORT_ID&type=tcp#my-vps
```

| Platform | Client | Notes |
| --- | --- | --- |
| Android | [v2rayNG](https://github.com/2dust/v2rayNG) or [sing-box (SFA)](https://github.com/SagerNet/sing-box) | Use VPN mode, then enable Android's Always-on VPN (section 1) |
| Windows | [v2rayN](https://github.com/2dust/v2rayN) | Enable TUN mode so all apps go through the tunnel |
| Linux, macOS, Windows | [Hiddify](https://github.com/hiddify/hiddify-app) | Beginner-friendly, imports the link |
| iOS | Streisand (free) or Shadowrocket (paid), from the App Store | Check the developer name before installing |

**Only download clients from the official GitHub releases or app stores linked above.** Fake proxy apps with malware are common.

### Check

1. Connect, then visit [ifconfig.me](https://ifconfig.me): it must show your VPS address.
2. Run the leak tests and the **pull the plug test** from section 1.
3. On the server: `sudo journalctl -u xray -n 50` shows errors if the client cannot connect. Wrong short ID, UUID or public key are the usual causes.

### Habits that keep it safe

- Give every device its own UUID and short ID, so you can revoke one without touching the rest.
- Re-run the installer command monthly to update Xray.
- Avoid one-click web panels (such as 3x-ui) unless you keep the panel off the public internet. They add an attack surface you do not need.
- Some countries restrict circumvention tools. Know your local law.

## 3. QUIC-based setup: Hysteria2 with Salamander

### Why (theory)

**QUIC** is the fast, UDP-based protocol behind HTTP/3 (YouTube, Chrome, Google). **Hysteria2** is a proxy protocol built on QUIC. It is designed to stay fast on bad connections: mobile data, long distances, Wi-Fi with packet loss. By default it pretends to be an HTTP/3 web server.

**Salamander** is Hysteria2's optional obfuscation layer. It scrambles every packet into random-looking bytes using a shared password, so a network cannot recognise the QUIC pattern. The official docs also list an experimental extension called Gecko that additionally chops up the QUIC handshake.

|  | VLESS + REALITY (section 2) | Hysteria2 (this section) |
| --- | --- | --- |
| Transport | TCP, looks like normal HTTPS | UDP, looks like HTTP/3 or random noise |
| Stealth | Very strong | Good, best with Salamander |
| Speed on lossy links | Normal | Much better |
| Fails when | TCP port 443 is blocked or the server IP is flagged | The network blocks or throttles UDP or QUIC |
| Use it for | Daily browsing, hostile networks | Streaming, video calls, mobile data, backup route |

**When and why from a privacy view:** speed is not privacy, but a slow tunnel gets switched off, and a switched-off tunnel protects nothing. Keep both on one VPS: Xray listens on TCP 443 and Hysteria2 on UDP 443, and they do not collide. If one gets blocked, the other still works.

**Salamander trade-offs.** With Salamander on, your server can no longer answer probes like a real website, because everything looks like noise. Some networks drop unidentifiable UDP, and some mobile carriers do too. A wrong obfuscation password gives a silent timeout, exactly as if the server were off. Start without Salamander, and turn it on only if QUIC itself is blocked or throttled on your network.

### Do: server setup

**Step 1. Install** with the [official script](https://v2.hysteria.network/docs/getting-started/Server/):

```bash
bash <(curl -fsSL https://get.hy2.sh/)
```

**Step 2. Get a certificate.** Option A (you own a domain): point its A record at your VPS and use the built-in ACME client (Step 4). Option B (no domain, fewer identity links): make a self-signed certificate and let the client pin its fingerprint:

```bash
sudo openssl req -x509 -nodes -newkey ec -pkeyopt ec_paramgen_curve:prime256v1 \
  -keyout /etc/hysteria/server.key -out /etc/hysteria/server.crt \
  -subj "/CN=www.bing.com" -days 3650
sudo chown hysteria:hysteria /etc/hysteria/server.key /etc/hysteria/server.crt
openssl x509 -noout -fingerprint -sha256 -in /etc/hysteria/server.crt
```

Save the fingerprint printed by the last command.

**Step 3. Make two passwords** (one to log in, one for Salamander):

```bash
openssl rand -hex 24
openssl rand -hex 24
```

**Step 4. Write `/etc/hysteria/config.yaml`** (Option B shown):

```yaml
listen: :443

tls:
  cert: /etc/hysteria/server.crt
  key: /etc/hysteria/server.key

auth:
  type: password
  password: PASTE_AUTH_PASSWORD

obfs:
  type: salamander
  salamander:
    password: PASTE_OBFS_PASSWORD
```

For Option A, replace the `tls:` block with an `acme:` block (open `80/tcp` while it fetches the certificate):

```yaml
acme:
  domains:
    - your.domain.example
  email: you@example.com
```

To run **without** Salamander, delete the `obfs:` block and add a `masquerade:` block that makes probes see a real website:

```yaml
masquerade:
  type: proxy
  proxy:
    url: https://news.ycombinator.com/
    rewriteHost: true
```

**Step 5. Open the port and start.**

```bash
sudo ufw allow 443/udp
sudo systemctl enable --now hysteria-server.service
sudo systemctl status hysteria-server.service
```

Logs: `sudo journalctl --no-pager -e -u hysteria-server.service`.

### Do: client setup

Easiest route: import a share link into a GUI client with Hysteria2 support ([Hiddify](https://github.com/hiddify/hiddify-app), [sing-box](https://github.com/SagerNet/sing-box) apps, or v2rayNG). Turn on its TUN or VPN mode so every app uses the tunnel.

```
hysteria2://PASTE_AUTH_PASSWORD@YOUR_SERVER_IP:443/?obfs=salamander&obfs-password=PASTE_OBFS_PASSWORD&sni=www.bing.com&insecure=1&pinSHA256=PASTE_FINGERPRINT#hy2-vps
```

With a domain and ACME certificate, drop `insecure` and `pinSHA256`, and set `sni` to your domain.

Command-line route: download `hysteria` from the [official releases](https://github.com/apernet/hysteria/releases) and save this as `client.yaml`:

```yaml
server: YOUR_SERVER_IP:443
auth: PASTE_AUTH_PASSWORD
obfs:
  type: salamander
  salamander:
    password: PASTE_OBFS_PASSWORD
tls:
  sni: www.bing.com
  insecure: true
  pinSHA256: PASTE_FINGERPRINT
socks5:
  listen: 127.0.0.1:1080
```

Run `./hysteria -c client.yaml`. Note that SOCKS mode only covers apps you point at it. For a true kill switch use TUN mode in a GUI client and the firewall rules from section 1 (they already allow the server port on UDP).

### Check

1. `curl --socks5-hostname 127.0.0.1:1080 https://ifconfig.me` must print your VPS address.
2. Timeout? Check the obfuscation password first, then whether your network blocks UDP. Try again without Salamander to tell which.
3. Run the leak tests and pull the plug test from section 1.

Port hopping (spreading traffic across many UDP ports to dodge per-port throttling) is covered in the [official docs](https://v2.hysteria.network/docs/advanced/Full-Client-Config/) under advanced options. Start without it.

## 4. Never verify: posting time and writing style

### Why (theory)

Hiding your IP is not enough. People are also identified by **when** they act and **how** they write. If an anonymous account posts at your lunch break and writes like you, the link is easy to draw. "Never verify" means: never give an observer a way to confirm that two accounts are the same person. Guide 1 covers the first layer (no real name, email or phone number); this section covers the time and style layers.

### Part A: posting time

**The leak.** Timestamps show your time zone, your sleep hours and your work routine. They also line up with real-world events: you post five minutes after something happened to you, or only while your real account is offline.

**The defence.** Never post live. Write now, publish later, at an irregular time. Prefer the platform's own scheduler, because a third-party tool needs a login token to your account and may add a visible "via" label.

| Platform | How to schedule |
| --- | --- |
| X, LinkedIn, Threads (web) | Built-in *Schedule* option in the post composer (availability varies by account type) |
| Facebook and Instagram | [Meta Business Suite](https://business.facebook.com) (Instagram needs a professional account) |
| Telegram | Press and hold *Send*, choose *Schedule Message* |
| Mastodon | Built-in scheduled posts in the composer |
| YouTube | YouTube Studio, *Schedule* at upload |
| TikTok | Built-in scheduling on the web uploader |
| Bluesky, Reddit, anything else | Self-hosted [Postiz](https://github.com/gitroomhq/postiz-app) or [Mixpost](https://mixpost.app), or SaaS tools such as [Buffer](https://buffer.com) |

If you self-host a scheduler on your VPS, keep it off the public internet (reach it only through your VPN or an SSH tunnel), because it stores the keys to all your accounts.

**Rules for the timing itself:**

- Publish at random minutes (7:43, not 8:00) and vary the hour and weekday.
- Do not post the instant your real account is active, or right after real events.
- Avoid a steady rhythm. A fixed fake time zone is also a fingerprint.

**Random delay for tinkerers.** Put this at the top of your crontab (`crontab -e`) so the job posts at a random moment in a two-hour window:

```
SHELL=/bin/bash
0 18 * * * sleep $((RANDOM % 7200)) && /home/you/post.sh
```

**Other time leaks to close:**

- **Photo metadata** (EXIF) also holds the capture time, so strip it as Guide 1's metadata-removal section shows.
- **Git commits** record your time zone. Commit from a shell set to UTC: `TZ=UTC git commit -m "message"`.
- **Device clock:** keep the private device or VM set to UTC.

### Part B: writing style (stylometry)

**The leak.** *Stylometry* uses statistics on how you write: favourite words, sentence length, punctuation habits, spelling variants, emoji, capitalisation. Given a few thousand words, software can often match an anonymous text to a known author. Typing rhythm (keystroke timing) can be measured by websites too.

**The defence is a routine, not one tool:**

**Step 1. Write a "tells list".** Read three of your real posts and note your habits: favourite phrases, the Oxford comma, em dashes, "honestly", double spaces, ellipses, British or American spelling, emoji. Ban them for your private persona.

**Step 2. Choose a flatter style.** Short sentences, common words, one spelling convention, no slang, no idioms from your native language. Plain writing is harder to fingerprint.

**Step 3. Rewrite with a local model.** Use a model that runs on your own machine ([Ollama](https://ollama.com), see section 16) and prompt it: *"Rewrite this text in plain, neutral English. Short sentences. No idioms. Keep the meaning."* Never paste private drafts into a cloud chatbot, since the company then holds your text. Caveat: AI output has its own recognisable flavour, so finish with a manual pass.

**Step 4. Paste, do not type.** Compose in a plain text editor inside the private compartment, then paste. This hides typing rhythm.

**Step 5. Say less.** Every extra word is data. Short, infrequent posts are the strongest defence.

**Tools that help you test yourself:**

- [Stylo](https://github.com/computationalstylistics/stylo) (R package) clusters texts by style. Load a few of your real writings and a few persona posts. If they cluster together, your persona is not different enough.
- [JStylo and Anonymouth](https://github.com/psal/anonymouth) are older research tools built for exactly this problem. They can be fiddly to run (Java), so treat them as optional.
- The [Whonix guide on posting and blogging](https://www.whonix.org/wiki/Surfing_Posting_Blogging) explains stylometry risks in plain language.

### Check

1. Do your private posts avoid every item on your tells list?
2. Does a Stylo test put the persona and your real self in different clusters?
3. Are your accounts unlinked (Guide 1's rules on usernames, emails and phone numbers)? A reverse image search also finds reused profile pictures.
4. What you write about can identify you even when the style does not: local events, your job, your family. Leave those out.

## 5. GrapheneOS on a spare phone

### Why (theory)

A stock phone is a tracking device with a phone app attached. [GrapheneOS](https://grapheneos.org) is a hardened, open-source Android that removes Google's built-in services and adds real security: a stricter app sandbox, a hardened memory allocator, per-app **Network** and **Sensors** switches, and Google Play as an optional, sandboxed app instead of part of the system.

**Developer opinion:** a dedicated GrapheneOS phone is the biggest single privacy upgrade in this book for the money. It also ends the daily bargain of "free" Android: no background Google account, no default telemetry, and apps you can cut off from the internet with one toggle. The cost is some convenience. A few banking apps and tap-to-pay wallets refuse to run, and you will install apps from different places than the Play Store.

**Supported phones (as of 10 October 2026):** Pixel 6 and newer, up to the Pixel 10 series. Buy an **unlocked, non-carrier** model, because carrier variants may not allow bootloader unlocking. The Pixel 11 generation was still unsupported in September 2026, so confirm on the [official site](https://grapheneos.org/faq) before buying. A newer Pixel gets more years of updates.

### Motorola and GrapheneOS: why it deserves support

For five years GrapheneOS ran only on Google's own phones, which means a privacy project depended on the company it protects you from. In March 2026, at Mobile World Congress, Motorola announced a partnership with the GrapheneOS Foundation to build devices that meet GrapheneOS's hardware requirements from the start. On 24 September 2026 GrapheneOS confirmed that the **Motorola Signature 27** will be its first officially supported non-Pixel phone, though no release date for GrapheneOS support has been announced.

Why this is worth cheering for:

- **Choice.** One hardware vendor is a single point of failure, in price and in policy.
- **Pressure on the industry.** GrapheneOS asks for long update commitments (seven years) and hardware security features. If those become a selling point, other makers feel pushed to match them.
- **Motorola does a large part of the porting work,** which makes support more durable than a volunteer effort chasing each new phone.

Honest caveats: the first supported models are expensive flagships, 2026 Motorola phones on sale today do **not** meet GrapheneOS requirements, and Motorola is owned by Lenovo, which some privacy readers will weigh in their own judgement. **Do not buy a current Motorola expecting GrapheneOS to arrive.** Ways to help: buy supported hardware when it launches, donate to the [GrapheneOS Foundation](https://grapheneos.org/donate), and report bugs well.

### Do: install GrapheneOS (about 45 minutes)

Installing **erases everything** on the phone. You need a USB-C data cable and a computer with a browser that supports WebUSB (Chrome, Edge, Brave or Vanadium; Firefox does not work). The official [web installer guide](https://grapheneos.org/install/web) is the final authority and may be newer than this text.

**Step 1. Update and unlock the option.** Boot the stock OS, connect to Wi-Fi, run all system updates, then go to *Settings, About phone*, tap *Build number* 7 times. Go to *Settings, System, Developer options* and turn on **OEM unlocking** (it needs internet).

**Step 2. Enter the bootloader.** Power the phone off, then hold **Volume Down + Power** until the bootloader screen appears. Connect it to the computer.

**Step 3. Open the installer.** Visit [grapheneos.org/install/web](https://grapheneos.org/install/web) and click **Unlock bootloader**. Confirm on the phone with the volume and power buttons.

**Step 4. Flash.** Click **Download release**, then **Flash release**. Do not unplug the cable or let the computer sleep.

**Step 5. Lock the bootloader.** Click **Lock bootloader** and confirm on the phone. This turns on verified boot, which is what makes the system tamper-evident. Never skip it.

**Step 6. First boot.** Set up the phone offline first. After setup, optionally switch **OEM unlocking** off again in *Developer options* to protect the device.

**Step 7. Verify it.** Install the **Auditor** app from the GrapheneOS app store and run an attestation check against [attestation.app](https://attestation.app). This proves the phone runs a genuine, locked GrapheneOS.

### Do: first-day hardening

- **Unlock method:** a random passphrase of six or more words is strongest. If you use a PIN, make it 8 or more digits and turn on PIN scrambling.
- **Auto reboot:** check *Settings, Security and privacy, Device unlock, Auto reboot*. A reboot returns the phone to its most protected state.
- **Cut apps off:** for each app that does not need the internet (notes, calculator, camera), switch off its **Network** permission in *Settings, Apps*.
- **Google Play:** skip it if you can. If one app needs it, install **Sandboxed Google Play** from the GrapheneOS app store and keep it in a separate user profile.
- **Profiles:** use *Settings, System, Multiple users* to make separate profiles for separate identities. Each profile has its own apps and can use its own always-on VPN (section 1).
- **Wi-Fi only:** a SIM card links the phone to your identity in many countries and reveals location to the carrier. A Wi-Fi-only spare phone avoids that.
- **Where to get apps:** GrapheneOS app store, [Accrescent](https://accrescent.app), F-Droid and Obtainium (see section 15).

### Check

1. *Settings, About phone* shows GrapheneOS and the latest release.
2. The boot screen shows the verified-boot message, not a yellow "bootloader unlocked" warning.
3. A leak test (section 1) shows your tunnel address and the kill switch holds.
4. Read the [usage guide](https://grapheneos.org/usage) for app compatibility questions.

## 6. Qubes OS on the PC (with Whonix built in)

### Why (theory)

On a normal PC, one infected browser tab can read your whole disk. [Qubes OS](https://www.qubes-os.org) fixes this by running every activity in its own lightweight virtual machine, called a **qube**. A hacked "untrusted" qube cannot see your "vault" qube. Qubes marks each one with a colour so you never confuse them. (Guide 1 covers lighter compartmentalization with ordinary VMs, sandboxes and separate devices, for people who do not need Qubes.)

**Developer opinion:** at this level, extreme compartmentalisation is not paranoia, it is just good engineering. Qubes also ships with [Whonix](https://www.whonix.org/wiki/Qubes) (a Tor-routed workstation), so anonymous browsing comes pre-wired. The price is a learning curve, a decent computer, and no GPU acceleration inside qubes. Qubes limits *damage*; it does not make you anonymous by itself. The current release is Qubes 4.3.

| Core idea | What it means |
| --- | --- |
| dom0 | The admin layer. Never browse or open files here. |
| Template qube | The master copy of an OS (Fedora, Debian). Update it, and every qube based on it updates. |
| App qube | Your work, personal or untrusted workspace, built from a template. |
| Vault qube | No network at all. Keep password databases and keys here. |
| Disposable qube | Single-use. Open a suspicious file, then it vanishes. |
| sys-net, sys-firewall, sys-usb | Isolated qubes that handle network cards, filtering and USB devices. |
| sys-whonix, anon-whonix | Tor gateway and the anonymous workstation behind it. |

### Do: check your hardware

Qubes needs a CPU with hardware virtualisation (Intel VT-x and VT-d, or AMD-V and AMD-Vi/IOMMU) switched on in the BIOS, plenty of RAM (16 GB is comfortable) and an SSD. Intel integrated graphics work best. Search your laptop model on the [Hardware Compatibility List](https://www.qubes-os.org/hcl/) first, and read the official [system requirements](https://www.qubes-os.org/doc/system-requirements/).

### Do: the "dual boot" reality

Qubes' own documentation does not recommend sharing a disk with another OS, because the other system could tamper with the boot files and weaken Qubes' protection. Pick the safest option you can afford:

| Option | How | Risk |
| --- | --- | --- |
| 1. Second internal drive (best) | Install Qubes on its own SSD. Choose which OS to boot with the firmware boot menu (often F12 or Esc). Unplug the Windows or Linux drive during install. | Lowest |
| 2. Fast external USB 3 SSD | Install onto the external drive. Plug it in only when you want Qubes. Qubes says this works well. | Low |
| 3. Same drive, shared | Shrink the old OS's partition and let the Qubes installer use the free space. Back up first. | Highest. Needs care and can damage the old install |

If the other OS is Windows with BitLocker, **save your BitLocker recovery key first**: changing boot settings can lock you out. On Linux, note your LUKS passphrase and back up the data.

### Do: install

**Step 1. Download and verify.** Get the ISO from [qubes-os.org/downloads](https://www.qubes-os.org/downloads/). Then follow the official [signature verification guide](https://www.qubes-os.org/security/verifying-signatures/) and compare the master key fingerprint across independent sources. This proves nobody swapped the file.

**Step 2. Write the USB stick.** On Linux, find your stick with `lsblk` and double-check the device name, because this command erases the target:

```bash
sudo dd if=Qubes-R4.3-x86_64.iso of=/dev/sdX bs=1M status=progress conv=fsync
```

(Use your exact ISO file name.) On Windows use [Rufus](https://rufus.ie) in **DD Image mode**.

**Step 3. Boot and install.** Boot from the stick, choose language, select the target disk, and tick **encrypt disk**. Use a long passphrase (6+ random words) because it protects everything. Finish and reboot.

**Step 4. First-boot configuration.** Accept the defaults and make sure these are ticked: *default system qubes*, *default application qubes*, **Whonix qubes**, and optionally *use Whonix for updates*.

**Step 5. Update everything first.** Open the Qubes menu, then *System Tools, Qubes Update*, select all and run it. Do this before anything else and then weekly.

### Do: daily use

- **Browse anonymously:** *Qubes menu, anon-whonix, Tor Browser* (or a Whonix disposable). Guide 3, Part 1 explains Tor and Whonix in depth.
- **Open unknown files:** right-click the file, choose *Open in disposable qube*.
- **Copy and paste between qubes:** press Ctrl+C in the source, **Ctrl+Shift+C**, click the destination window, **Ctrl+Shift+V**, then Ctrl+V. This stops one qube from silently reading your clipboard.
- **Move files:** right-click, *Copy to other qube*.
- **Passwords and keys:** keep them in the *vault* qube only.
- **One identity per qube:** create a new qube (*Qubes menu, Create New Qube*) for each persona.
- **Your VPN or tunnel:** make a dedicated VPN qube, set it as the network source of your app qubes, and run the Xray or Hysteria client from sections 2 and 3 inside it. The firewall then blocks leaks for every qube behind it. The Qubes community's VPN guides are linked from the [documentation](https://www.qubes-os.org/doc/).
- **Backups:** *Qubes menu, System Tools, Backup Qubes*, to an encrypted external drive.

### Check

1. Qubes Update reports no pending updates.
2. The Tor Browser check page in anon-whonix says you are using Tor.
3. Your vault qube shows no network connection.
4. A disposable qube closes and leaves nothing behind.

## 7. Advanced ad blocking: Pi-hole and NextDNS

### Why (theory)

Ads are the delivery truck for tracking: every ad request tells a data broker who you are and what you read. Before your device loads anything, it asks a **DNS server** to turn a name (like `ads.example.com`) into an address. A **DNS filter** lies to your device about known tracker names, so the request dies before it starts. It protects every app, even on a smart TV, and it is invisible.

It has limits: it cannot block ads served from the same domain as the content (most of YouTube), and it relies on blocklists. So use layers: a **DNS filter** for the whole network (this section), plus the in-browser blocker you set up in Guide 1.

| Layer | Best for | Trade-off |
| --- | --- | --- |
| Pi-hole (at home) | Whole home network, you own the data | Needs a Raspberry Pi or small Linux box; only protects at home unless you add a VPN |
| NextDNS (cloud) | Phones and laptops anywhere | A third party sees your DNS queries (you can turn logging off) |
| Mullvad DNS (free, no account) | Quick, no setup | Fewer controls |

### Do: Pi-hole on a home server

**Step 1. Prepare.** Use a Raspberry Pi or any always-on Linux machine (Debian or Ubuntu). In your router, give it a fixed address (a *DHCP reservation*).

**Step 2. Install** with the [official installer](https://docs.pi-hole.net/main/basic-install/):

```bash
curl -sSL https://install.pi-hole.net | bash
```

Accept the defaults, choose an upstream DNS provider, and note the admin address shown at the end. Set the admin password:

```bash
sudo pihole setpassword
```

**Step 3. Point your network at it.** In the router's DHCP settings set the **DNS server** to the Pi-hole's address, then reconnect your devices. Where the router does not allow this, set DNS on each device by hand.

**Step 4. Add stronger blocklists.** Open the admin page (`http://pi.hole/admin`), go to *Lists*, and add [HaGeZi's lists](https://github.com/hagezi/dns-blocklists) or [OISD](https://oisd.nl). Then run `pihole -g` to update.

**Step 5. Stop sending your queries to a big DNS company.** Install **Unbound**, a local recursive resolver, using Pi-hole's [Unbound guide](https://docs.pi-hole.net/guides/dns/unbound/). Your queries then go straight to the internet's root servers instead of Google or Cloudflare.

**Away from home?** Run a WireGuard server on the same box with [PiVPN](https://pivpn.io) and connect your phone back home. Never expose a Pi-hole's DNS port to the public internet.

### Do: NextDNS (works anywhere)

**Step 1.** Create an account at [nextdns.io](https://nextdns.io). Your **Setup** page shows a configuration ID (for example `abc123`).

**Step 2. Privacy tab:** turn on blocklists (the NextDNS Ad-Blocking list plus HaGeZi Pro), and enable *Block disguised third-party trackers*.

**Step 3. Settings tab, Logs:** switch logging **off** or set the shortest retention. A privacy tool that keeps your history defeats the point.

**Step 4. Connect each device** (replace `abc123` with your ID):

| Device | Where to set it |
| --- | --- |
| Android and GrapheneOS | *Settings, Network and internet, Private DNS*, choose hostname, enter `abc123.dns.nextdns.io` |
| iPhone and iPad | Download the config profile from your Setup page (or use the NextDNS app) |
| Windows 11 | *Settings, Network and internet, your connection, DNS server assignment, Manual*, turn on DNS over HTTPS with template `https://dns.nextdns.io/abc123` |
| Linux and macOS | Install the [NextDNS CLI](https://github.com/nextdns/nextdns) from its Setup page |
| Firefox | *Settings, Privacy and Security, DNS over HTTPS, Max Protection*, custom provider `https://dns.nextdns.io/abc123` |

For a no-account option, set Private DNS or DoH to Mullvad's filtering resolver, `adblock.dns.mullvad.net`.

### Check

1. Visit the [ad block tester](https://d3ward.github.io/toolz/adblock.html). Most tests should show *blocked*.
2. Pi-hole: `nslookup doubleclick.net` should return `0.0.0.0`.
3. If a site breaks, find the blocked domain in the Pi-hole or NextDNS log, and allow only that domain.
4. Open [BrowserLeaks DNS](https://browserleaks.com/dns). It must show your chosen resolver, and not your ISP, once you also run your VPN or tunnel.

## 8. Privacy frontends for social media

### Why (theory)

A **frontend** is a different website or app that shows you a platform's public content without the platform's trackers, scripts or login walls. A proxy fetches the posts for you, so the platform sees the proxy and not you. This is the cleanest fix for "I only want to read this one post".

**Honest status check:** big platforms fight frontends hard. Many projects have died (Bibliogram) or are fragile. Instances (public servers) come and go every week. So the plan has three parts: use **LibRedirect** to switch to whatever works now, keep a **fallback**, and for platforms with no working frontend, **contain** the platform instead (isolated browser profile, no app). Tightening the account settings themselves is covered in Guide 1.

### Do: set up LibRedirect (once, covers most sites)

**Step 1. Install.** Get [LibRedirect](https://github.com/libredirect/browser_extension) from the Firefox add-ons site or the project's releases page. It redirects links on supported sites to a privacy-respecting frontend automatically.

**Step 2. Turn on the services you use** in its settings and pick a frontend for each. Press *Update instances* now and then, because the lists are pulled from each project's official instance page.

**Step 3. Fallback.** [Farside](https://farside.link) gives links that jump to a working instance, for example `https://farside.link/nitter/USERNAME`.

**Step 4. Want reliability? Self-host one.** Redlib (Reddit) is a good first project. On your VPS, bind it to localhost only and reach it through an SSH tunnel:

```bash
docker run -d --name redlib -p 127.0.0.1:8080:8080 quay.io/redlib/redlib:latest
# on your own computer:
ssh -L 8080:127.0.0.1:8080 admin@YOUR_SERVER_IP    # then open http://127.0.0.1:8080
```

### Platform by platform

| Platform | Frontend situation | What to do |
| --- | --- | --- |
| **Twitter / X** | Nitter exists but is unreliable, many instances are dead | Use LibRedirect or Farside; check the [Nitter instances wiki](https://github.com/zedeus/nitter/wiki/Instances). For posting, consider Mastodon |
| **Instagram** | Frontends such as Proxigram are listed but flaky | Try LibRedirect; otherwise use the website (not the app) in an isolated profile and the container steps below |
| **Facebook** | No working frontend | Use the website in a [Facebook Container](https://addons.mozilla.org/firefox/addon/facebook-container/) in Firefox. Do not install the app |
| **Threads** | No frontend | Same as Instagram. It is tied to your Instagram account, so treat both as one identity |
| **LinkedIn** | No maintained frontend | Browser only, never the app. Account settings are in Guide 1 |
| **GitHub** | No mainstream frontend | Read code with `git clone`. Account hardening is in Guide 1. Consider [Codeberg](https://codeberg.org) for new projects |
| **Telegram** | Public channels have an official no-login preview | Read channels at `https://t.me/s/CHANNELNAME`, or follow them in a feed reader through [RSS-Bridge](https://github.com/RSS-Bridge/rss-bridge) |
| **YouTube** | Invidious and Piped instances are unreliable | Use [NewPipe](https://newpipe.net) on Android and [FreeTube](https://freetubeapp.io) on desktop. They run on your device and need no account |
| **Reddit, TikTok** | [Redlib](https://github.com/redlib-org/redlib) and ProxiTok | Use LibRedirect |

### Do: contain platforms with no frontend

**Step 1. Separate browser profile or qube.** Use a dedicated Firefox profile (or a Qubes qube from section 6) only for Meta and LinkedIn. Never log into them elsewhere.

**Step 2. Firefox containers.** Install the Facebook Container add-on. It walls off Facebook, Instagram and Threads from the rest of your browsing, so they cannot follow you across sites.

**Step 3. Website, not the app.** Apps request many permissions and bundle trackers. On GrapheneOS, if you must use an app, switch off its Network permission when you are not using it, and run it in a separate profile.

### Check

1. Click a YouTube or Reddit link: the address bar should change to your chosen frontend.
2. In a Facebook Container tab, open another site: the Meta tracking buttons should be blocked.

## 9. Privacy alternatives: choosing, the hardened browser and maps

### Why (theory)

Every app you use is a trade: you give data, you get convenience. A privacy alternative flips the trade, and usually costs some **inconvenience**. The skill is spending that inconvenience where it buys the most privacy. At the hardened level you accept mid-to-high inconvenience for very high privacy. Sections 10 to 13 cover each category.

**What lives in Guide 1.** Guide 1 holds the tool-selection test (open code, encryption by default, how little the tool needs to know you, how it makes money, audits, whether your contacts use it) and the everyday app stack: ad blocker, search engine, password manager, email aliases and authenticator. This guide assumes those are in place and only adds what costs more effort.

### Where to keep finding better tools

This field moves fast and projects die. Re-check these every few months:

- [Privacy Guides](https://www.privacyguides.org), a community-run recommendation list with its reasoning shown.
- [privacytools.io](https://privacytools.io), the older directory. Its original contributors moved to Privacy Guides in 2021 ([their explanation](https://www.privacyguides.org/articles/2021/09/14/welcome-to-privacy-guides/)) and the site is no longer run by that team, so treat its listings with care.
- [digdeeper.club](https://digdeeper.club), independent, opinionated articles about practical privacy and surveillance. Read it critically, as with any single author.
- [Awesome Privacy](https://github.com/Lissy93/awesome-privacy), a large curated list.

Treat every list as an input, not a verdict. Check who runs it and whether listings are paid for.

### Do: the hardened browser and offline maps

**1. Browser.** Easiest strong choice: [Mullvad Browser](https://mullvad.net/en/browser), built with the Tor Project to resist fingerprinting and ships with a content blocker. For Firefox users who want control, harden Firefox with the [arkenfox user.js](https://github.com/arkenfox/user.js):

```bash
# Firefox: open about:profiles, create a new profile, click "Open Folder" on its root directory.
# Close Firefox. In that folder:
curl -O https://raw.githubusercontent.com/arkenfox/user.js/master/user.js
curl -O https://raw.githubusercontent.com/arkenfox/user.js/master/updater.sh
chmod +x updater.sh && ./updater.sh
```

Put your own tweaks in a file called `user-overrides.js` so updates do not erase them. Expect some sites to break; the arkenfox wiki explains each setting. On Android use **IronFox** or Firefox with an ad blocker (Guide 1).

**2. Maps.** Use offline maps built on OpenStreetMap data: [Organic Maps](https://organicmaps.app), its community fork [CoMaps](https://comaps.app), or [OsmAnd](https://osmand.net). They need no account and can work with no network.

Migrate one app at a time, keep the old one for two weeks, then delete the old account properly. Section 17 gives the order of work.

### Check

1. Can you name, for each of your 10 most-used apps, a replacement you have tested?
2. Run [coveryourtracks.eff.org](https://coveryourtracks.eff.org): your browser should show strong protection against tracking.

## 10. Communication: messages, calls and media

### Why (theory)

Two things can leak when you chat: the **content** (what you said) and the **metadata** (who talked to whom, when, from which IP). End-to-end encryption protects content. Metadata is harder, and it is where the real differences between messengers lie. A phone number is metadata too: it is a permanent, real-world ID.

| Messenger | Needs a phone number? | Hides your IP and contacts | Forward secrecy | Best for | Main caveat |
| --- | --- | --- | --- | --- | --- |
| [Signal](https://signal.org) (setup in Guide 1) | Yes, to register (a username can hide it from contacts) | Partly | Yes | Everyday chats with normal people, strongest crypto | Phone number at sign-up |
| [Session](https://getsession.org) | **No**, random Session ID | **Yes**, messages go through an onion-routed network | **Not yet**; Protocol V2 is planned | An anonymous account that hides your IP | No forward secrecy today; the foundation had a funding scare in April 2026, since resolved |
| [SimpleX Chat](https://simplex.chat) | **No**, and no user ID of any kind | Strongest, optional Tor | Yes | Highest metadata privacy | Younger app, you share invite links |
| [Briar](https://briarproject.org) | No | Tor, or direct Bluetooth and Wi-Fi with no internet | Yes | Crises, journalists, shutdowns | Android only; both people must be online |
| [Element](https://element.io) (Matrix) | No | No | Yes (content) | Group communities | Server sees metadata; setup is complex |

**Why Session is attractive:** it asks for no phone number and no email, so there is nothing to link to you. Its onion routing hides your IP from the people you talk to. **The honest trade-off:** Session dropped forward secrecy, so if your long-term key is ever stolen, old messages could be decrypted. Treat it as an anonymity tool for contacts who must not know who you are, and choose Signal or SimpleX for the most sensitive content.

### Do: Session

**Step 1. Install** from the [official download page](https://getsession.org/download) (it lists the app stores and F-Droid). Open it and tap *Create Account*.

**Step 2. Choose a display name that is not your real name.** You now have a Session ID (a long string). This is your identity.

**Step 3. Save the recovery password.** Open *Settings, Recovery Password* and write the words on paper. Without it a lost phone means a lost account.

**Step 4. Notifications.** Choose the slower, Google-free background mode if offered, especially on GrapheneOS.

**Step 5. Privacy settings** (*Settings, Privacy*): turn **off** read receipts, typing indicators and link previews; turn **on** screen security.

**Step 6. Add people** by sharing your Session ID or QR code in person, or through another trusted channel. Never post it publicly beside your real identity.

Session calls are in beta and connect peer to peer, which can reveal your IP to the other person. Run your tunnel (sections 1 to 3) when calling.

### Do: SimpleX (maximum metadata privacy)

Install from [simplex.chat](https://simplex.chat), create a profile with any name, and invite contacts with a one-time link or QR code, using *Incognito* mode so each contact sees a different random name. In *Network and servers* you can route through Tor.

### Do: private calls

- **Best default:** Signal voice and video, with the call-relay setting that hides your IP from the other person (Guide 1 shows where it is). Group calls are end-to-end encrypted.
- **No phone number:** SimpleX calls, with the relay option switched on in its call settings.
- **Avoid** unencrypted phone calls and ordinary SMS for anything sensitive.

### Do: private media sharing

**Step 1. Clean the file first.** Remove metadata (time, GPS, device) as Guide 1's metadata-removal section shows.

**Step 2. Send small files inside an encrypted messenger.** For big files, use [OnionShare](https://onionshare.org): choose *Share Files*, click *Start sharing*, and send the `.onion` address. The recipient opens it in Tor Browser. Your computer is the server, so nothing sits on a company's disk.

**Step 3. Encrypt a file before sending it by any channel.** Guide 3, Part 5 shows how with age and minisign.

**Step 4. Photo library:** use [Ente Photos](https://ente.io) (end-to-end encrypted backups) instead of the default cloud gallery. Self-hosting your own library with Immich is in Guide 3, Part 8.

### Check

1. Your Session recovery words are on paper, stored away from the phone.
2. New SimpleX contacts were added in Incognito mode, so each sees a different random name.

## 11. Community, video, music, podcasts and games

### Why (theory)

Entertainment apps are among the hungriest data collectors: what you watch, listen to and play reveals your mood, politics, health and habits, and it is sold or used for ad targeting. The fix has two parts. **Swap** the service for a private one, or **self-host** your own media (Guide 3, Part 7 covers hosting). For services you cannot replace (a streaming subscription), **contain** them: browser not app, separate profile, tracker blocking.

(I read "chatmart" in your outline as chat rooms; marketplaces are not covered here.)

### Community: forums, chat rooms, groups

| Instead of | Use | Why |
| --- | --- | --- |
| Discord, WhatsApp groups | [Matrix](https://matrix.org) with [Element](https://element.io) or [FluffyChat](https://fluffychat.im); [SimpleX](https://simplex.chat) or Signal groups for small trusted groups | Open protocol, end-to-end encrypted rooms, no phone number needed |
| Reddit | [Lemmy](https://join-lemmy.org) (read Reddit through Redlib, section 8) | Federated, open source, no ad tracking |
| X, Facebook | [Mastodon](https://joinmastodon.org) | Federated, no ads. Direct messages are **not** end-to-end encrypted |
| Old-school chat | IRC on [Libera.Chat](https://libera.chat), through Tor | Minimal data, long history |

**Setup: Matrix in 5 steps.**

1. Install Element or FluffyChat from the official site or F-Droid.
2. Choose a server that does not demand a phone number. Pick a [public server](https://servers.joinmatrix.org), or self-host later (Guide 3, Part 6).
3. Create an account with a throwaway username (no real name).
4. When asked to *Set up Secure Backup*, save the **security key on paper**. Without it you lose message history on a new device.
5. Create rooms with **encryption on** (the default for private rooms). Verify each friend's device by comparing the emoji shown on both screens.

**Setup: Mastodon privacy.** Pick a server at [joinmastodon.org/servers](https://joinmastodon.org/servers). Then in *Preferences* turn on *Opt-out of search engine indexing*, *Hide your social graph*, and set your default post visibility to *Followers only*.

### Video entertainment: movies, series, broadcasts

| Need | Private option |
| --- | --- |
| Your own films and shows | A self-hosted media server (Jellyfin and others): Guide 3, Part 7 |
| Online video | [PeerTube](https://joinpeertube.org) (search with [SepiaSearch](https://sepiasearch.org)), YouTube through NewPipe or FreeTube (section 8) |
| Free films | [Internet Archive](https://archive.org) for public-domain titles, and your public library's free streaming services |
| Netflix, Prime, Disney+ | No private alternative exists. **Contain** them: browser (not the app), separate profile, an ad blocker (Guide 1), a unique alias email |
| Smart TV | The TV itself is a tracker. Switch off its "viewing data" and ad-personalisation options (names vary by brand), or **never connect it to the internet** and use a separate streaming box instead |

### Music, podcasts and audiobooks

| Need | Private option |
| --- | --- |
| Buy and own music | [Bandcamp](https://bandcamp.com) (DRM-free files), then play offline with [Auxio](https://github.com/OxygenCobalt/Auxio) on Android or VLC |
| Stream your own library | A self-hosted music server (Navidrome): Guide 3, Part 7 |
| Podcasts | [AntennaPod](https://antennapod.org) (Android, from F-Droid), which uses open RSS feeds and needs no account |
| Audiobooks | [LibriVox](https://librivox.org) for free public-domain readings, and your library's borrowing apps. Your own files: Audiobookshelf, in Guide 3, Part 7 |

### Games: from puzzles to MMOs

**Mobile.** Many free games are ad networks wearing a costume. Prefer offline, open-source games from F-Droid such as [Shattered Pixel Dungeon](https://shatteredpixel.com) and Simon Tatham's Puzzles. On GrapheneOS, switch off the **Network** permission of any offline game. DNS filtering (section 7) blocks the ad networks in the rest.

**PC.** Buy from stores that give you offline installers: [GOG](https://www.gog.com) and [itch.io](https://itch.io) are DRM-free. On Steam, set *Profile, Edit Profile, Privacy Settings* to **Private**, especially *Game details*.

**Anti-cheat warning.** Some online games install kernel-level anti-cheat that runs at the deepest level of your system. Play those on a separate gaming PC or a separate drive, never on the machine that holds your private life.

**Voice chat.** Use [Mumble](https://www.mumble.info) (hosting it is in Guide 3, Part 7) or Matrix voice rooms rather than a data-hungry platform.

**Linux gaming:** install launchers through Flatpak, then tighten their file access with [Flatseal](https://flathub.org/apps/com.github.tchx84.Flatseal).

### Check

1. Your Matrix security key is on paper.
2. Your Steam profile and game details show as private to a logged-out viewer.
3. On the smart TV, the viewing-data option is off, or the TV has no network.

## 12. Work tools: code, mail, meetings and storage

### Why (theory)

These four tools hold your most valuable data: your source code, your correspondence, your conversations and your files. Two rules cover most of the risk. **Encrypt on your own device before data leaves it**, so the provider only ever holds scrambled data. And **separate identities**, so a leak in one place does not unlock the rest.

### Code and dev architecture

| Need | Private option |
| --- | --- |
| Hosting | [Codeberg](https://codeberg.org) (non-profit, runs Forgejo) or [SourceHut](https://sourcehut.org). Running your own Forgejo is in Guide 3, Part 7 |
| Editor | [VSCodium](https://vscodium.com), the same editor without Microsoft telemetry |
| AI help | A local model (section 16), not a cloud assistant that receives your code |
| Secrets | Never in Git. Keep them in KeePassXC or encrypt with [SOPS](https://github.com/getsops/sops) or `age` |
| Risky code | Run it in a disposable qube (section 6) or a throwaway container |

**Do: a clean Git identity with signed commits.** Signing proves commits are yours without exposing your real email:

```bash
git config --global user.name "your-pseudonym"
git config --global user.email "ID+name@users.noreply.github.com"   # or your host's noreply address
git config --global gpg.format ssh
git config --global user.signingkey ~/.ssh/id_ed25519.pub
git config --global commit.gpgsign true
```

Split work and personal identities by folder. Add this to `~/.gitconfig`, and put the work name and email in `~/.gitconfig-work`:

```
[includeIf "gitdir:~/work/"]
    path = ~/.gitconfig-work
```

**Architecture habit:** one qube or VM per project, SSH keys on a hardware key (section 14), and no real name anywhere in commit history.

### Professional mail and messaging

Email was not designed to be private. Even with encryption, the sender, recipient, subject and time are visible to servers. So: pick a privacy-friendly provider, and use end-to-end encryption for anything sensitive.

| Need | Private option |
| --- | --- |
| Mail provider | [Proton Mail](https://proton.me/mail), [Tuta](https://tuta.com), [Mailbox.org](https://mailbox.org) or [Posteo](https://posteo.de) |
| Sign-up aliases | SimpleLogin or addy.io (Guide 1) |
| Work chat | [Element](https://element.io) (Matrix), [Zulip](https://zulip.com) or [Mattermost](https://mattermost.com). Hosting your own is in Guide 3, Part 7 |
| Small teams | Signal groups |

**Do: end-to-end encrypted email with OpenPGP in Thunderbird.**

1. Install [Thunderbird](https://www.thunderbird.net) and add your mail account.
2. Open *Account Settings, End-to-End Encryption, Add Key, Create a new OpenPGP Key* (use the default modern key type).
3. Export your **public key** and share it with contacts. Keep the private key and its passphrase safe.
4. Import each contact's public key, and verify its fingerprint by another channel.
5. Write a new message and switch on **Require encryption**.

Proton to Proton messages are encrypted automatically. Remember that the subject line and recipient list are not hidden.

**Work devices:** if your employer enrolls your phone in management (MDM), they can see and wipe it. Use a separate work profile or a separate cheap device.

### Virtual meetings

| Need | Private option |
| --- | --- |
| Small private calls | Signal (relay on) or Element Call |
| Team meetings | [Jitsi Meet](https://jitsi.org), with end-to-end encryption switched on in *Security options* |
| Run your own server | Jitsi or Nextcloud Talk: Guide 3, Part 7 |
| Forced to use Zoom, Teams or Meet | Browser only, no app, no account, separate profile, mic and camera off by default |

**Meeting hygiene:** use a lobby and a password, share a single window not your whole screen, blur your background, and turn off cloud recording and AI note-taker bots.

### Storage: files and media

| Need | Private option |
| --- | --- |
| E2EE cloud drive | [Proton Drive](https://proton.me/drive) or [Filen](https://filen.io) |
| Encrypt files in *any* cloud | [Cryptomator](https://cryptomator.org) |
| No server at all | [Syncthing](https://syncthing.net), peer-to-peer sync between your own devices (steps in Guide 3, Part 8) |
| Backups | [restic](https://restic.net), encrypted and deduplicated |
| Photos | Ente Photos (section 10); self-hosted Immich is in Guide 3, Part 8 |

Encrypting a drive or a whole disk (LUKS, VeraCrypt, BitLocker, FileVault) is covered by Guide 1's full-disk-encryption section, so it is not repeated here.

**Do: encrypted backups to your VPS with restic:**

```bash
restic -r sftp:admin@YOUR_SERVER_IP:/home/admin/backups init     # choose a strong password
restic -r sftp:admin@YOUR_SERVER_IP:/home/admin/backups backup ~/Documents
restic -r sftp:admin@YOUR_SERVER_IP:/home/admin/backups snapshots
```

The VPS only ever sees encrypted blocks. Keep the restic password in your password manager and on paper. **Test a restore** once.

**Do: Cryptomator for an existing cloud.** Install it, click *Add Vault, Create New Vault*, choose a folder inside your cloud-sync folder, set a long password and **save the recovery key**. Files you drop in the unlocked vault are encrypted, including their names, before the cloud sees them.

### Check

1. `git log --show-signature -1` shows a good signature and a noreply email.
2. An encrypted test email arrives and a stranger's mail server would show only scrambled content.
3. Your backup restores a test file onto a different folder.

## 13. Payments: crypto

*Welcome to the core of the internet.* This is the one section with deliberate limits. Everything defensive and lawful has full steps. Three topics are **theory only**: mixers and joiners, Bitcoin to Monero to Bitcoin hopping, and routing payments in between through your own Monero node. They are explained, with their legal risk, but not turned into a how-to. Everyday payments (cash first, then cards and wallets) are covered in Guide 1, and hosting your own Bitcoin, Ethereum and Monero nodes is in Guide 3, Part 9.

### Why (theory)

Money is the strongest identity link there is. A card ties you to what you bought, where and when. Privacy in payments means **sharing the least necessary**, keeping accounts separate and avoiding public payment feeds, while staying legal and keeping records for tax. Crypto is not a magic cloak: most blockchains are public ledgers, and an exchange that verified your ID can connect you to every address you withdraw to.

### Crypto basics

**Pseudonymous is not anonymous.** An address is a random-looking string, but every payment to and from it is public forever. If one address is ever tied to your name (through an exchange, a shop or a post), the whole history can follow. Good practice is therefore about habits: new address per payment, one identity per wallet, and your own node.

Three legitimate ways to get crypto: buy on a **regulated exchange** (simple, requires ID), trade **peer to peer** (for example with [Bisq](https://bisq.network), where legal in your country), or **earn** it by being paid for work. Rules on reporting and tax differ by country, so keep records of every purchase and sale.

### Bitcoin

**Step 1. Buy** on a regulated exchange available in your country, using a bank transfer. This is the simplest and lowest-risk way.

**Step 2. Get a hardware wallet** (for example [Trezor](https://trezor.io), [BitBox02](https://bitbox.swiss) or [Coldcard](https://coldcard.com)) and set it up offline. Write the **recovery seed words on paper or steel**. Never photograph them, type them into a computer or store them in the cloud.

**Step 3. Install [Sparrow Wallet](https://sparrowwallet.com)** (desktop), and verify the download signature as its site explains. Choose *File, New Wallet, Connect Hardware Wallet*, then create the wallet from your device.

**Step 4. Withdraw from the exchange** to a **fresh Sparrow receiving address**. Send a small test amount first.

**Step 5. Test your backup:** restore the seed on a spare device with a small balance before trusting it with more.

**Step 6. Privacy habits that are fully legal:**

- Use a new address for every receipt. Sparrow shows an unused one each time.
- Do not mix coins from different sources in one transaction. Use Sparrow's *UTXOs* tab to label and choose coins.
- Connect Sparrow to **your own node** in *Preferences, Server*. Public servers learn your addresses. Setting one up is in Guide 3, Part 9.
- Route Sparrow through Tor: in *Preferences, Server*, enable the proxy `127.0.0.1:9050` while Tor runs.
- For small everyday payments, a self-custody Lightning wallet such as [Phoenix](https://phoenix.acinq.co) is more private than on-chain.

### Theory only: mixers and joiners

**What they are.** A *custodial mixer* takes your coins and pays out different coins. A *CoinJoin* (a "joiner") is non-custodial: many people build one shared transaction with equal-sized outputs, so an observer cannot tell which input paid which output. Smart-contract pools on Ethereum do something similar.

**Why people want them.** Anyone who learns one of your addresses (an employer, a landlord, a stalker) can read your balance and payment history. Privacy from that is legitimate.

**Where the law stands (as of October 2026, and changing).** There is no single global answer.

- In the US, regulators have treated custodial mixers as money transmitters that must register and verify customers.
- Tornado Cash was sanctioned in 2022. A US appeals court found the sanction overreached, and the Treasury removed it from the list in March 2025. Its developer Roman Storm was convicted in August 2025 of conspiring to run an unlicensed money-transmitting business, while the jury could not decide the laundering and sanctions counts. His case was still unresolved in the latest reporting I found.
- The Samourai Wallet founders pleaded guilty in 2025 to a similar money-transmitting conspiracy and were sentenced that November.
- In March 2026 a US Treasury report acknowledged that mixers have lawful privacy uses. In early October 2026 the proposed US "mixer rule" was reportedly withdrawn, while a federal appeals court upheld the Bitcoin Fog operator's conviction. Operator and developer liability remains unsettled.

**The practical risks, even for lawful money.** Exchanges can freeze, flag or close accounts that receive coins touching a mixer. Mixing does not make coins "clean" in law: funds that came from a crime stay tainted. If you are considering this, talk to a lawyer in your country first. **No steps are given here.**

### Ethereum

**Step 1. Buy** on a regulated exchange, as with Bitcoin.

**Step 2. Self-custody:** use a hardware wallet with a wallet interface such as [Rabby](https://rabby.io). Seed handling is the same as for Bitcoin.

**Step 3. Save on fees** by using a layer 2 network (Base, Arbitrum, Optimism) for everyday transactions. Send a small test amount before big transfers.

**Step 4. Privacy habits:**

- Every address and contract call is public. Keep **one address per identity** and never post an address next to your real name.
- An ENS name (like `name.eth`) publicly links all activity of that address to you.
- Your wallet's **RPC provider** sees your IP and addresses. Set a privacy-respecting RPC, or run your own node (Guide 3, Part 9), and use Tor or your tunnel (sections 1 to 3).

**Step 5. Stay safe:** review and revoke old token approvals at [revoke.cash](https://revoke.cash), never sign a request you do not understand, and use a separate browser profile for crypto.

Newer designs such as Privacy Pools try to let users prove their funds are not from illicit sources while keeping details private. Treat these as theory for now and check their legal status in your country.

### Monero

Monero hides sender, receiver and amount by default. It is lawful to hold and use in most countries, but it is getting harder to buy: some regulated exchanges have dropped it, Japan and South Korea have banned privacy coins on domestic exchanges, and the EU's new anti-money-laundering rules (applying from July 2027) will stop regulated providers there from handling anonymity-enhancing coins and anonymous accounts. The rules target providers, not individuals, but check your local position.

**Step 1. Get it ethically:** from a regulated exchange where your country still allows it, peer to peer, by being paid in it, or by mining with [P2Pool](https://p2pool.io).

**Step 2. Wallet:** [Feather Wallet](https://featherwallet.org) (desktop, has Tor built in), the official [Monero GUI](https://www.getmonero.org/downloads/), or [Cake Wallet](https://cakewallet.com) on mobile. **Verify the download signature** following the instructions on getmonero.org.

**Step 3. Create the wallet** and write the seed words on paper (never digitally).

**Step 4. Connect:** your own node is best (Guide 3, Part 9). Otherwise use a remote node **through Tor**, because a remote node can see your IP.

**Step 5. Receive:** create a new **subaddress** for each person who pays you. Incoming funds unlock after about 10 confirmations (roughly 20 minutes).

**Step 6. Spend:** paste or scan the merchant's address and confirm.

### Theory only: Bitcoin to Monero to Bitcoin

**The idea.** Swap Bitcoin for Monero, then back to Bitcoin, hoping the trail between the first and last coins goes cold. People call it chain-hopping.

**Legal reality.** Swapping is a lawful trade in most places, but each swap can be a **taxable event**, and compliance teams treat BTC to XMR to BTC patterns as a classic red flag for hiding where money came from. Exchanges can close accounts, and when the purpose is to disguise criminal proceeds it is money laundering. With the 2027 EU rules and bans elsewhere, the route is also getting narrower. **No steps are given here.**

### Theory only: payments in between, through your own Monero node

**The idea.** A node checks the blockchain itself. A wallet that uses *someone else's* node reveals your IP and which transactions you are interested in. Your own node removes that third party, and broadcasting through Tor or a VPN hides your address from the network. That part is ordinary privacy hygiene, legal in most countries, and its setup is in Guide 3, Part 9.

**The "payments in between" part.** If the plan is to insert extra payments to break the trail between where money came from and where it ends up, then for *lawful* funds it is merely extra privacy, but for *unlawful* funds it is the "layering" step of money laundering, and it creates patterns that investigators and exchange compliance look for. Running a node and using Tor are not illegal in most countries, though a few restrict VPNs or Tor. The risk comes from the source of the funds, tax duties and any obligations to verify customers. **No steps for this part are given.**

### Check

1. You have restored your wallet from the seed on a spare device, with a small balance.
2. Your seed words exist only on paper or steel, in a safe place.
3. Sparrow shows a fresh address for each receipt and no coins from different sources are merged.
4. You keep a record of every buy, sale and swap for tax.

Sources for the legal summary: [National Law Review on Tornado Cash](https://natlawreview.com/article/tornado-cash-and-legality-mixers), [CoinDesk on the March 2026 Treasury report](https://www.coindesk.com/policy/2026/03/09/u-s-treasury-signals-shift-on-crypto-mixers-acknowledges-legitimate-privacy-uses), and [Cryptostate's status review](https://cryptostate.wiki/en/articles/crypto-mixers-tornado-cash-coinjoin). Not legal advice.

## 14. Hardware security

### Secure data erasure

**Why (theory).** Deleting a file only removes its entry; the data stays on the drive until it is overwritten. "Quick format" is the same. A proper erase is **one-way**: after it, no tool can bring the data back, because the old bits are gone or the key that unscrambles them has been destroyed. Standards such as NIST SP 800-88 describe three levels: *clear*, *purge* and *destroy*.

The method depends on the drive type:

| Drive | What works | What does not |
| --- | --- | --- |
| Hard disk (HDD) | One overwrite pass (random or zero) | Old 35-pass schemes: unnecessary on modern disks |
| SSD and NVMe | The drive's own **secure erase**, **sanitize** or **cryptographic erase** | Overwriting: wear-levelling hides cells from the OS, so some data survives |
| USB sticks, SD cards | Encrypt from day one, then destroy the key. Physically destroy if highly sensitive | File shredders |
| Phones | Factory reset on an encrypted phone (GrapheneOS, iOS, Android) destroys the keys | Deleting apps |

The golden rule: **turn on full-disk encryption from the start** (Guide 1 shows how for LUKS, BitLocker, FileVault and phones). Then erasing means destroying the key, which is instant and cannot be reversed.

**Do: HDD.** Boot a live USB (not the disk you are wiping) and check the device name with `lsblk`:

```bash
sudo shred -v -n 1 -z /dev/sdX          # one random pass, then a zero pass
```

Or use [nwipe](https://github.com/martijnvanbrummelen/nwipe) from the bootable [ShredOS](https://github.com/PartialVolume/shredos.x86_64) image, which gives a menu. Choose *PRNG stream* or *zero fill*, one round.

**Do: SATA SSD:**

```bash
sudo hdparm -I /dev/sdX | grep -i -E "frozen|erase"     # must say "not frozen"
sudo hdparm --user-master u --security-set-pass p /dev/sdX
sudo hdparm --user-master u --security-erase p /dev/sdX
```

If it says *frozen*, suspend and wake the PC, then retry.

**Do: NVMe SSD:**

```bash
sudo nvme id-ctrl /dev/nvme0 -H | grep -i -E "format|sanitize|crypto"   # see what the drive supports
sudo nvme format /dev/nvme0n1 --ses=2     # cryptographic erase, if supported
sudo nvme format /dev/nvme0n1 --ses=1     # otherwise: user data erase
```

**Windows and macOS:** use the maker's tool for SSDs (Samsung Magician, WD Dashboard and similar), or *Reset this PC, Remove everything, Clean the drive*. On a Mac with FileVault, *Erase All Content and Settings* destroys the key. **Individual files:** shredding is unreliable on SSDs. Keep sensitive files inside an encrypted container, and delete the container and its key instead.

**Check:** `sudo hexdump -C /dev/sdX | head` should show only zeros or noise, with no readable text.

### Antivirus and EDR: developer opinion

An antivirus or EDR (endpoint detection and response) agent runs at the deepest level of your computer, can read every file and every web request, and usually uploads data to its vendor. For a privacy-focused, hardened setup that is a bad trade, for four reasons:

- **Huge attack surface at the highest privilege.** Security researchers have repeatedly found serious flaws in mainstream antivirus engines. A bug there is a bug with full control of your machine.
- **It can break everything at once.** In July 2024 a faulty update from one EDR vendor crashed millions of Windows computers worldwide.
- **It is a data collector.** Many suites inspect your encrypted web traffic and report file and URL data back to the vendor.
- **Your design already does its job.** With compartments (Qubes, GrapheneOS profiles), sandboxing and verified boot, an infected app is boxed in. Prevention by design beats cleaning up afterwards.

Honest limits: this opinion fits a hardened, compartmentalised setup run by someone who follows good habits. On a normal Windows PC that opens attachments and runs random downloads, Guide 1's advice stands: keep the **built-in** Microsoft Defender, add no third-party suite on top, and on Linux scan on demand with ClamAV. Organisations that must prove compliance do need EDR; that is a different problem.

**Do: basic hygiene that replaces it.**

1. **Update everything** automatically: system, browser, apps, router.
2. **Use a standard account** for daily work, not an administrator account.
3. **Verify downloads** with the publisher's hash or signature (`sha256sum file`) and download only from official pages.
4. **Check a suspicious file by its hash**, not by uploading it: `sha256sum file`, then search the hash on [VirusTotal](https://www.virustotal.com). Uploading shares the file.
5. **Open untrusted attachments** in a disposable qube or a throwaway VM.
6. **Stop and check** before clicking: who sent it, what is the real address, and why is it urgent? Log in through bookmarks, never through links in messages.
7. **Switch off Office macros**, and never use cracked software, the top malware source.
8. **Never plug in a found USB stick.** Use a USB data blocker for public chargers.
9. **Keep an offline backup** so ransomware cannot take everything (section 12).

### Hardware keys: MFA and physical unlock

**Why (theory).** A hardware security key (FIDO2) is a small USB or NFC device that proves you are you. Because it checks the **real website address** before answering, it cannot be fooled by a fake login page. SMS codes and app codes can be phished, a key cannot. Buy **two** (one main, one backup) directly from makers such as [Yubico](https://www.yubico.com) or [Nitrokey](https://www.nitrokey.com).

**Do: MFA on your accounts.**

1. Start with your **email account**, because it can reset all the others. Open its security settings and choose *Add security key*.
2. Register **both keys** on each account.
3. Keep the backup key and the printed backup codes in a safe place.
4. Remove SMS as a recovery or login method wherever the service allows.
5. For sites that only support one-time codes, store them in the key with the Yubico Authenticator, or in the authenticator app from Guide 1.
6. Test by logging out and logging in with the key on a private window.

**Do: use the key to unlock devices.**

- **Linux sudo and login** (with `pam_u2f`). Keep a root terminal open during testing so you cannot lock yourself out:

```bash
sudo apt install libpam-u2f
mkdir -p ~/.config/Yubico
pamu2fcfg > ~/.config/Yubico/u2f_keys          # touch the first key
pamu2fcfg -n >> ~/.config/Yubico/u2f_keys      # touch the backup key
# then add this as the first "auth" line of /etc/pam.d/sudo:
#   auth required pam_u2f.so
```

- **Linux disk encryption** (advanced): `sudo systemd-cryptenroll --fido2-device=auto /dev/nvme0n1p3` adds your key as an unlock method. Follow the [Arch Wiki on systemd-cryptenroll](https://wiki.archlinux.org/title/Systemd-cryptenroll) for your distribution, and **keep your passphrase slot** as a fallback.
- **SSH and Git signing:** a key that never leaves the device, and needs a touch and PIN:

```bash
ssh-keygen -t ed25519-sk -O resident -O verify-required -C "hwkey-1"
```

Add the new `.pub` file to your VPS's `~/.ssh/authorized_keys`.

- **KeePassXC:** *Database, Database Security, Add additional protection, Add YubiKey Challenge-Response*. The vault then needs your passphrase *and* the key.
- **Windows:** protect the disk with BitLocker and **TPM plus PIN** (*Group Policy, Require additional authentication at startup*), and use [Yubico's Windows login tool](https://www.yubico.com/products/computer-login-tools/) for sign-in.
- **macOS:** use the key for your Apple ID, and enable FileVault.
- **Phones:** use a long passphrase. No key unlocks a GrapheneOS phone.

### Killing unnecessary hardware

**Why (theory).** Every radio and sensor you leave on can be used to track or attack you: Wi-Fi probes, Bluetooth beacons, microphones, cameras, USB ports. A switched-off device cannot leak. The rule: **if you are not using it, turn it off, and prefer a hardware switch.**

**Wi-Fi and Bluetooth with rfkill (Linux).**

```bash
rfkill list                       # show all radios
sudo rfkill block wifi
sudo rfkill block bluetooth
sudo rfkill block wwan            # cellular modem
sudo rfkill unblock all           # turn everything back on
```

To keep Bluetooth off at every boot, make `/etc/systemd/system/rfkill-block.service` (check the path with `which rfkill`):

```
[Unit]
Description=Block Bluetooth at boot
[Service]
Type=oneshot
ExecStart=/usr/sbin/rfkill block bluetooth
[Install]
WantedBy=multi-user.target
```

Then run `sudo systemctl enable rfkill-block.service`.

| Device | How to switch it off |
| --- | --- |
| Webcam | Software: add `blacklist uvcvideo` to a file in `/etc/modprobe.d/`. Physical covers are in Guide 1, and a laptop with a hardware kill switch is better still |
| Microphone | Hardware switch if available; otherwise a plug-in dummy audio jack, or assign the mic to one qube only in Qubes. Software mute is not proof |
| USB ports | Disable unused ports in the BIOS. On Linux, install [USBGuard](https://usbguard.github.io): with only trusted devices plugged in, run `sudo usbguard generate-policy > rules.conf`, install it as `/etc/usbguard/rules.conf`, then `sudo systemctl enable --now usbguard` |
| Thunderbolt and DMA | Turn off Thunderbolt in the BIOS if unused, and enable IOMMU and Kernel DMA Protection |
| Phone USB-C | On GrapheneOS keep USB peripherals on *charging-only when locked* (or off) in *Settings, Security and privacy* |
| NFC, location, UWB | Off by default on a phone, enabled only when needed |
| Cellular modem | Remove the SIM, or the WWAN card in a laptop |
| Smart-home gadgets | Put them on a separate guest network, and mute or unplug voice assistants |
| Firmware | Set a BIOS password, enable Secure Boot, and disable booting from USB or network if you do not need it |

**Hardware you cannot switch off.** Intel Management Engine and AMD's equivalent are hidden processors inside the chip. You cannot remove them in normal use. Choose hardware from the [Qubes certified list](https://www.qubes-os.org/doc/certified-hardware/) or other open-firmware laptops if this matters to you. Guide 3, Part 10 goes further, with Libreboot and me\_cleaner.

**Tracker tags.** Use [AirGuard](https://github.com/seemoo-lab/AirGuard) on Android to detect AirTags and similar devices following you.

### Check

1. `rfkill list` shows *Soft blocked: yes* for the radios you turned off.
2. A test drive has been erased and shows only zeros.
3. Logging into your email account requires touching a key, and the backup key works too.
4. A USB stick you did not authorise does nothing when plugged in (USBGuard).

## 15. Privacy-focused app downloading

### Why (theory)

Where you get an app from decides two things: who learns what you install, and whether the file is really what it claims to be. A big app store needs an account, builds a profile of every app you run, and still hosts malicious clones. Better sources are **open**, check **signatures** (a mathematical proof that the file came from the real developer and was not altered), and need **no account**.

| OS | Best source | Extras | Avoid |
| --- | --- | --- | --- |
| Android and GrapheneOS | GrapheneOS app store, [Accrescent](https://accrescent.app), [F-Droid](https://f-droid.org) | Obtainium, Aurora Store | Random APK sites |
| Linux | Your distro's signed repositories | Flathub with Flatseal | Unreviewed PPAs, blind `curl \| bash` |
| Windows | `winget` and the vendor's own site | Scoop | Download portals that bundle installers |
| iOS | App Store, with strict settings | Web apps instead of apps | Unknown profiles and sideload tools |

### Android and GrapheneOS

**Step 1. Built-in store first.** On GrapheneOS open *Apps* and install what you need from it (Vanadium browser, Camera, Auditor, PDF Viewer).

**Step 2. Install F-Droid.** It carries only open-source apps built from source. Download the APK from [f-droid.org](https://f-droid.org). Allow *Install unknown apps* for your browser, install the file, then **turn that permission off again**. In F-Droid's settings, add the [IzzyOnDroid repository](https://apt.izzysoft.de/fdroid) for more apps.

**Step 3. Add Accrescent**, a privacy-focused store with strict signature checks and a smaller catalogue, from [accrescent.app](https://accrescent.app).

**Step 4. [Obtainium](https://github.com/ImranR98/Obtainium)** installs apps straight from their developers' release pages and updates them automatically. Paste the project's repository link and tap *Add*. You are trusting that developer directly.

**Step 5. Play Store apps without a Google account:** the [Aurora Store](https://auroraoss.com) with anonymous login. It sometimes breaks, because Google resists it. Banking apps that insist on Google services work better inside Sandboxed Google Play on GrapheneOS (section 5).

**Verify an APK** on a computer with Android build tools, and compare the printed certificate fingerprint with the one the developer publishes:

```bash
apksigner verify --print-certs app.apk
```

The [AppVerifier](https://github.com/soupslurpr/AppVerifier) app does the same on the phone.

### Linux

**Step 1. Use the distribution's repositories first.** Packages there are signed and reviewed: `sudo apt install qbittorrent` (or `dnf` or `pacman`).

**Step 2. Flathub for newer apps**, then limit what each app can touch:

```bash
sudo apt install flatpak
flatpak remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo
flatpak install flathub org.qbittorrent.qBittorrent
flatpak install flathub com.github.tchx84.Flatseal     # permission manager
```

**Step 3. Read before you pipe.** Several tools in this book offer a one-line `curl ... | bash` installer. Safer: download, read, then run.

```bash
curl -fsSL https://example.com/install.sh -o install.sh
less install.sh
bash install.sh
```

**Step 4. Verify downloads** with the publisher's hash or signature: `sha256sum -c SHA256SUMS` or `gpg --verify file.asc file`. Be careful with random PPAs, and read the build recipe of anything you take from the AUR.

### Windows

**Step 1. Use `winget`** (built into Windows 11). It installs from checked manifests:

```powershell
winget search qbittorrent
winget install --id qBittorrent.qBittorrent
winget upgrade --all
```

[Scoop](https://scoop.sh) is an alternative for portable apps.

**Step 2. For anything else,** type the vendor's real address yourself. Never use "download" sites that wrap installers. Then check the hash: `Get-FileHash .\file.exe -Algorithm SHA256`, and compare it to the vendor's value.

**Step 3. Tighten Windows itself.** In *Settings, Privacy and security*: set diagnostic data to the minimum, switch off tailored experiences and ad ID, and **keep Recall off** on PCs that offer it. Review the apps list with [O&O ShutUp10++](https://www.oo-software.com/en/shutup10) if you want a one-screen overview.

### iOS

iOS is the most closed system, so your options are fewer. Make an Apple ID with a unique alias email (Guide 1 covers aliases), and then:

1. *Settings, Privacy and Security, Tracking*: switch off *Allow Apps to Request to Track*.
2. *Settings, Apple Advertising*: turn off *Personalised Ads*. In *Analytics and Improvements*, turn everything off.
3. Before installing, read the **App Privacy** label on the store page and skip apps that collect "data linked to you" without need.
4. Prefer the website to an app. A site in Safari with a content blocker (such as AdGuard or Wipr) leaks less than a tracker-filled app.
5. Set permissions to *While Using* or *Ask Next Time* and turn off *Precise Location*.
6. Where available, use iCloud **Advanced Data Protection** for end-to-end encrypted backups.
7. For people at high risk, *Lockdown Mode* shrinks the attack surface.

### Torrent clients: qBittorrent over uTorrent

Torrents are just a file-sharing technology, fine for Linux distributions, open-source software and public-domain media. Use them only for legal content. **qBittorrent** is open source and ad-free, while uTorrent carries ads and bundled offers, and has a history of shipping unwanted software.

1. Install [qBittorrent](https://www.qbittorrent.org) from the site or with `winget` or `apt`.
2. *Tools, Options, Advanced, Network Interface*: choose your **VPN or tunnel interface** (for example `wg0` or `tun0`). Now the client has no network at all if the tunnel drops.
3. *Options, BitTorrent*: enable **Anonymous mode** and set encryption to *Require encryption*.
4. *Options, Connection*: turn off *UPnP / NAT-PMP* port forwarding.
5. Never torrent over Tor. It harms the Tor network and can leak your IP.

### Check

1. Everything you installed came from a source in the table, and you can say which one.
2. The *Install unknown apps* permission is off for your browser and file manager.
3. `winget upgrade` or `sudo apt update && apt list --upgradable` shows nothing pending.
4. With the tunnel stopped, qBittorrent shows no connections.

## 16. Hosting a private AI model with Ollama

### Why (theory)

When you type into a cloud chatbot, a company stores your words. A **local model** runs entirely on your own computer, so your prompts, drafts and documents never leave it. [Ollama](https://ollama.com) is the simplest way to run one: it downloads open models and serves them on your machine.

**Unfiltered models.** Mainstream models are trained to refuse many requests. "Uncensored" models are versions with those refusals reduced. There are two kinds: **fine-tunes** trained to follow instructions with fewer refusals (the Dolphin family from Cognitive Computations, Hermes from Nous Research) and **abliterated** models, where the refusal behaviour is edited out of an existing model by community uploaders. They are legal to run, and useful for fiction, security research on your own systems, and drafts that touch sensitive topics. Know the trade-offs: quality can drop, they will state falsehoods with confidence, and you are fully responsible for what you do with the output. Pick models whose model card explains how they were made, and never use them to harm anyone.

### What your computer needs

Models are shrunk ("quantised") to fit. At the common Q4 setting, a rough rule is **about 0.6 GB of memory per billion parameters, plus 1 to 2 GB for context**. Figures below are approximate.

| Model size | Memory needed (GPU VRAM, or RAM on a Mac) | Typical hardware | Example |
| --- | --- | --- | --- |
| 3 to 4B | 3 to 4 GB | Any recent laptop, CPU is fine | `gemma3:4b` |
| 7 to 8B | 5 to 6 GB | 8 GB GPU, or 16 GB Mac | `qwen3:8b`, `hermes3:8b` |
| 12 to 14B | 9 to 10 GB | 12 GB GPU, or 16 to 24 GB Mac | Hermes and Qwen 14B builds |
| About 20 to 24B | 14 to 16 GB | 16 to 24 GB GPU, or 32 GB Mac | `gpt-oss:20b`, Dolphin 24B, Mistral Small |
| About 32B | 20 to 24 GB | 24 GB GPU, or 32 to 48 GB Mac | 27 to 36B class models |
| 70B | 40 GB or more | Two GPUs, 48 GB GPU, or a 64 GB+ Mac | Hermes 3 70B |

- **No GPU?** CPU-only works but is slow. Expect a few words per second for an 8B model, and you want at least 16 GB of RAM, with fast memory.
- **Apple Silicon Macs** are excellent because the GPU shares your normal memory.
- **Disk:** each model takes roughly 2 to 45 GB. Keep 100 GB free if you want to try several.
- **A small VPS has no GPU,** so a spare home PC is the better host. Reach it through your VPN or SSH tunnel (section 1).

Model names change monthly. Check the [Ollama library](https://ollama.com/library) for current releases, and read each model card.

### Do: install and lock down

**Step 1. Install.** On Windows and macOS use the installer from [ollama.com/download](https://ollama.com/download). On Linux, read the script first (section 15), then run it:

```bash
curl -fsSL https://ollama.com/install.sh -o ollama-install.sh
less ollama-install.sh
sh ollama-install.sh
```

Or use Docker, with the port bound to your own machine only (add `--gpus=all` for NVIDIA with the container toolkit installed):

```bash
docker run -d -v ollama:/root/.ollama -p 127.0.0.1:11434:11434 --name ollama ollama/ollama
```

**Step 2. Turn off Ollama's cloud features.** Ollama can also run "cloud models" (names ending in `-cloud`) and web search on Ollama's servers, and those **do send your prompts out**. For full privacy, switch them off. On Linux run `sudo systemctl edit ollama` and add:

```
[Service]
Environment="OLLAMA_NO_CLOUD=1"
```

Then `sudo systemctl restart ollama`. (Mac and Windows: set `OLLAMA_NO_CLOUD=1` as an environment variable and restart the app.) The [Ollama FAQ](https://docs.ollama.com/faq) says the log then shows *Ollama cloud disabled: true*.

**Step 3. Keep the server private.** By default Ollama listens only on `127.0.0.1:11434`. **Never** set `OLLAMA_HOST=0.0.0.0` on a machine reachable from the internet: Ollama has no login, so anyone could use it. To use it from another device, tunnel in: `ssh -L 11434:127.0.0.1:11434 user@your-home-pc`.

### Do: download and chat

**Step 4. Pull and run a model** that fits your table row:

```bash
ollama pull qwen3:8b
ollama run qwen3:8b            # type /bye to leave
ollama list                    # see what you have
ollama rm qwen3:8b             # delete one
```

**Step 5. Try a less-filtered model.** Hermes and Dolphin are in the library:

```bash
ollama run hermes3:8b
ollama run dolphin3
```

Community abliterated models are published as **GGUF files** on [Hugging Face](https://huggingface.co). Ollama can run any GGUF straight from a repository:

```bash
ollama run hf.co/USERNAME/REPOSITORY:Q4_K_M
```

Choose a trusted uploader with a detailed model card, a clear licence and plenty of downloads. Pick the quantisation (`Q4_K_M` is the usual balance) that fits your memory. If a name above is not found, search the library; model names are renamed often.

**Step 6. Add a chat window.** [Open WebUI](https://openwebui.com) gives you a ChatGPT-style page that talks to your local Ollama:

```bash
docker run -d -p 127.0.0.1:3000:8080 --add-host=host.docker.internal:host-gateway \
  -v open-webui:/app/backend/data --name open-webui ghcr.io/open-webui/open-webui:main
```

Open `http://127.0.0.1:3000`, create the local admin account, and pick your model. Turn off any optional web-search or cloud integrations in its settings.

**Step 7. Use it for privacy tasks in this book:** rewriting a draft to blur your writing style (section 4) with *"Rewrite this in plain, neutral English. Short sentences. No idioms."*, summarising private documents, or helping with code (section 12) without sending anything to a third party.

### Check

1. `ss -tlnp | grep 11434` shows `127.0.0.1:11434`, not `0.0.0.0`.
2. The Ollama log says *Ollama cloud disabled: true* (`journalctl -u ollama -n 50`).
3. **Offline test:** turn off Wi-Fi and chat. A model that still answers is truly local.
4. `ollama ps` shows which model is loaded and whether it runs on the GPU. If it says CPU, your model is probably too big for your VRAM.

## 17. Master checklist, order of operations and legal notes

### Order of operations

Do not try everything at once. Each stage pays off on its own, and an unfinished big project leaks more than a finished small one. Finish Guide 1 first: its basics (password manager, aliases, authenticator, ad blocker, full-disk encryption, metadata removal) are assumed below.

| When | Do this | Sections |
| --- | --- | --- |
| Day 1 (about 2 hours) | Update everything. Install the hardened browser. Order two hardware keys | 9, 14 |
| Week 1 | Rent and harden a VPS, add the kill switch, set up VLESS + REALITY. Move chats to Session or SimpleX. Install LibRedirect and contain social media | 1, 2, 10, 8 |
| Weeks 2 to 3 | Install GrapheneOS on a spare phone. Set up Pi-hole or NextDNS. Start tested backups | 5, 7, 12 |
| Month 2 | Install Qubes OS. Register hardware keys everywhere. Turn off unused hardware. Run a local AI model | 6, 14, 16 |
| Ongoing | Update monthly, rotate keys and aliases, test a restore each quarter, re-check tool lists every few months | 9 |

### Master checklist

**Network and identity**

- [ ] VPS hardened: key-only SSH, firewall on, automatic updates on
- [ ] Kill switch active on every device, and the pull-the-plug test passed
- [ ] VLESS + REALITY works, and Hysteria2 is ready as a backup (if your network allows UDP)
- [ ] DNS filter (Pi-hole or NextDNS) in place, with logging off
- [ ] Separate identities: no shared username, photo, email or phone number across personas (Guide 1 rules)
- [ ] Posts scheduled at irregular times, Git commits in UTC
- [ ] Writing-style routine followed for private personas (tells list, local rewrite, paste not type)

**Devices**

- [ ] GrapheneOS phone installed, bootloader locked, verified with Auditor
- [ ] Qubes OS installed, updated, with vault, work and untrusted qubes
- [ ] Full-disk encryption on every computer (Guide 1), with passphrases written on paper
- [ ] Radios, ports, camera and mic you do not use are switched off
- [ ] Two hardware keys registered on email and other key accounts

**Apps and data**

- [ ] Apps come only from the sources in section 15, each verified
- [ ] Session recovery words on paper
- [ ] Encrypted backups, with a restore tested in the last 90 days
- [ ] Local AI runs with Ollama cloud features disabled

**Money**

- [ ] Crypto seeds only on paper or steel, with a tested restore
- [ ] Records kept of every crypto purchase, sale and swap for tax

### Legal and ethical notes

- **Local law comes first.** Some countries restrict or ban VPNs, circumvention tools, encryption or Tor. Others require logging by providers. Check yours before you rent a server or install a tunnel.
- **Terms of service.** Hosting companies and platforms have rules. Read them before running proxies or scrapers.
- **Crypto.** Tax and reporting duties apply in most places. The three theory-only topics (mixers and joiners, BTC to XMR to BTC hopping, and routing payments in between through your own node) are described so you understand them and their legal risk, not so you can do them. If you ever consider them, speak to a lawyer first.
- **Not for crime.** Nothing here is a way to hide illegal activity. These tools protect ordinary people from stalkers, brokers, snoops and censors. Using them to harm others is both wrong and, in most places, illegal.
- **Not professional advice.** This is general information. For threats beyond level 2 (journalists, activists, people facing state attackers), get help from groups such as the [EFF](https://www.eff.org), [Access Now's helpline](https://www.accessnow.org/help/) or the [Freedom of the Press Foundation](https://freedom.press).

### Facts this book depends on (as of 10 October 2026)

These change fast. Re-check them before you rely on them.

- GrapheneOS officially supports Pixel 6 up to the Pixel 10 series. The Motorola Signature 27 was confirmed in September 2026 as its first supported non-Pixel phone, with no support date announced.
- Qubes OS 4.3 is the current release.
- Session has no forward secrecy yet, and Protocol V2 is planned.
- Ollama cloud models and web search send prompts to Ollama's servers unless `OLLAMA_NO_CLOUD=1` is set.
- US treatment of mixers is in flux: Tornado Cash is off the sanctions list, its developer's case continues, and operators of mixers have been convicted or have pleaded guilty. EU anti-money-laundering rules restricting regulated providers from handling anonymity-enhancing coins apply from July 2027.
- Proxy and tunnel software (Xray, Hysteria2), privacy frontends and public instances change names, flags and availability often. The official documentation linked in each section always wins over this text.

### Keep learning

The [Qubes documentation](https://www.qubes-os.org/doc/), the [GrapheneOS usage guide](https://grapheneos.org/usage) and the [EFF Surveillance Self-Defense guide](https://ssd.eff.org). Tool directories (Privacy Guides and others) are listed once, in section 9.

## Appendix: what was removed from this guide, and where it lives now

Twelve groups of content were cut or shortened because Guide 1 or Guide 3 owns them. Ownership follows your draft outline. Guide 1 has no written text yet, so its side of this table rests on the outline alone.

| Section | Removed or shortened | Now lives in |
| --- | --- | --- |
| 1 | Tor column of the comparison table; Mullvad, Proton VPN and IVPN kill-switch steps | Guide 3, Part 1 (Tor); Guide 1 (VPN apps) |
| 4 | EXIF stripping with `exiftool` and `mat2`; the detailed account-unlinking check | Guide 1 (metadata removal, never-verify rules) |
| 7 | uBlock Origin table row and install step | Guide 1 (anti-tracking toolkit) |
| 8 | The *Tighten each account* step (GitHub, LinkedIn, Facebook, Instagram, Telegram settings), the noreply-email check, the *compare tool lists* check | Guide 1 (maximum privacy settings per platform); section 9 for tool lists |
| 9 | The 6-question tool test, Exodus Privacy and ToS;DR links, the core app stack (search engine, password manager, aliases, authenticator), the three-ring migration plan | Guide 1 (privacy alternatives, password hygiene); section 17 for order of work |
| 9, 17 | Repeated tool-directory lists | Kept once, in section 9 |
| 10 | The Signal setup steps and checks; the `age` file-encryption step; Immich; the `exiftool` step | Guide 1 (Signal and Telegram); Guide 3, Part 5 (age and minisign); Guide 3, Part 8 (Immich); Guide 1 (metadata) |
| 11 | The Docker Compose media stack (Jellyfin, Navidrome, Audiobookshelf) and its check; those three as table rows; self-hosting Mumble | Guide 3, Part 7 |
| 12 | The self-hosted Jitsi block; Forgejo self-hosting; LUKS, VeraCrypt, BitLocker and FileVault drive steps | Guide 3, Part 7; Guide 1 (full-disk encryption) |
| 13 | The whole *Normal payments* table (cash, virtual cards, subscriptions, wallets, peer-to-peer apps, parcel lockers); the opening explanation of what a Monero node is | Guide 1 (cash); Guide 3, Part 10.3 (lockers); Guide 3, Part 9 (nodes) |
| 14 | The Defender and ClamAV detail; the webcam shutter row | Guide 1 (antivirus, camera covers) |
| 17 | Day-1 password manager, alias and ad-blocker steps; the Signal and everyday-payments checklist items; the second tool-directory list | Guide 1; section 9 |

Every cut was replaced by a one-line pointer, so section numbers and cross-references still work.

### Judgement calls to confirm

- **Guide 1 must pick up what I assumed.** Two items are at risk because your Guide 1 outline does not clearly list them: the 6-question tool test, and virtual-card advice (the outline mentions cash only). If Guide 1 skips them, copy them back from your original Guide 2.
- **Overlaps I left in Guide 2, because Guide 2 comes first.** Hardening a VPS (section 1), Cryptomator and restic (section 12) and USBGuard (section 14) also appear in Guide 3, Parts 7, 8 and 10.1. Guide 3 Part 11 also repeats the Qubes install from section 6. These are better trimmed to pointers in Guide 3, and I can do that if you want.
- **Guides 4 and 5 repeat section 4.** Their outlines list posting at off times, schedulers and style-vector tools under *faking human vector*. I left both alone because you asked about Guides 1 and 3, but one guide should own them.
- **Monero node wording.** Your Guide 2 outline marks running your own node over Tor as theory only, while Guide 3, Part 9 gives practical node steps. I kept Guide 2's theory-only label for the *payments in between* part and sent ordinary node setup to Guide 3. Check you are comfortable with that line.
- **privacytools.io.** Section 9 now says its original contributors moved to Privacy Guides in 2021 and the site is no longer run by that team, per Privacy Guides' own announcement. Your outline asks to mention it, so I kept the link with that note. Revert it if you prefer.
- **Left alone on purpose:** the threat-model table and five rules in the introduction (shared series framing) and the repeated *know your local law* lines.
- **Not re-verified:** all technical steps in the retained sections are copied unchanged from your text.
