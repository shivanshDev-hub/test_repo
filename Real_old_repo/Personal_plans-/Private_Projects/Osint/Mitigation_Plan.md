# Identity Migration & Footprint Reduction Guide

**Goal:** build a new online identity that cannot be linked to your real one, and shrink the old footprint. **Reality check:** archives, scrapers, breach dumps and brokers keep copies. You are *breaking linkability and cutting exposure*, not achieving perfect erasure.

---

## Phase 0: Threat model (30 min)

Write down, per adversary, what they'd use against you.

| Adversary | Capability | What you must hide |
| --- | --- | --- |
| Casual searchers / employers | Google, social search | Name, photos, handles |
| Doxxers / hostile actors | Full OSINT tooling, breach data | Cross-platform links, writing style, location |
| Data brokers | Bulk scraping, purchased data | Phone, email, address |
| Platforms / legal process | Logs, IPs, payment data | Not hideable from the platform itself, so don't break ToS or law |

Pick a **separation level**: *Light* (different handle/email) → *Strong* (separate devices, networks, phone) → *Strict* (Qubes/Whonix, never any overlap). This guide assumes **Strong to Strict**.

---

## Phase 1: Audit the old footprint

### 1.1 Account discovery

1. List every handle, email and phone number you've ever used.
2. Run enumeration against yourself:

   ```
   pip install sherlock-project maigret holehe
   sherlock handle1 handle2 --print-found
   maigret handle1 --html
   holehe your@email.com
   ```
3. Search your mailbox for: `welcome`, `verify your email`, `confirm your account`, `password reset`, `receipt`. Mailboxes expose forgotten accounts.
4. Export saved logins from your browsers and password manager.
5. Check OAuth: Google (myaccount → Security → Third-party access), Facebook/Meta (Apps & Websites), Apple (Sign in with Apple), GitHub, Microsoft, X.
6. Check HIBP for each email and phone.

### 1.2 Passive exposure

- Google your name, handles, emails and phone in quotes. Add `site:` operators for key platforms.
- Reverse-image search old avatars and photos (Google Lens, Yandex, TinEye, PimEyes).
- Wayback Machine: search your profile URLs and personal domains.
- People-search and broker sites: Spokeo, Whitepages, BeenVerified, Radaris, etc. (use local equivalents too).
- Technical traces: `whois` history for your domains, Git commit emails (`git log --format='%ae' | sort -u`), PGP keys on keyservers, Gravatar, conference and speaker pages, podcast bios, paper and patent records, GitHub/GitLab/Docker Hub/npm/PyPI profiles.
- Public data: voter rolls, property records, business registrations (check for a registered-agent option).

### 1.3 Inventory

Build an **encrypted** spreadsheet (store in a VeraCrypt or age-encrypted file):

`Account | URL | Handle | Email | Phone | Linked accounts | Content sensitivity | Action (delete / lock / keep) | Status`

Rank by sensitivity. High-risk: anything with real name, face, location, employer, family.

---

## Phase 2: Build compartments

**Rule:** the new identity never touches anything from the old one.

| Element | Old identity | New identity |
| --- | --- | --- |
| Device | Daily driver | Separate laptop, or VM/OS profile (Qubes ideal) |
| Browser | Existing profiles | Fresh hardened profile (Firefox + uBlock, no sync, no old extensions) |
| Network | Home ISP | Dedicated VPN or Tor, consistent exit country |
| Email | Existing | New provider, created on clean setup |
| Phone | Real number | Separate prepaid SIM |
| Payments | Real cards | Virtual or prepaid only |
| Password vault | Existing | Brand new vault + new master password |
| 2FA | Existing | Hardware key bought separately, or TOTP stored in new vault |

### 2.1 Device

1. Preferred: a second device bought without linking it to your real accounts (no Apple ID/Google sign-in, no Wi-Fi auto-join to your home network).
2. Otherwise a dedicated VM (Whonix or a Linux VM) with its own snapshot.
3. Disable telemetry, Bluetooth, location services and cloud backup. Never log into old accounts on it.

### 2.2 Network

1. Pay for the VPN anonymously (Mullvad accepts cash/Monero) or use Tor for registration.
2. **Registration:** some platforms block Tor or flag it. A residential-reputation VPN exit often works better.
3. **Daily use:** pick one stable exit country/region. Erratic IPs trigger security checks and ban risk.
4. Never connect the clean setup from your home IP without the VPN, even once.
5. Don't use it from places that tie to you (home Wi-Fi without VPN, workplace, regular cafés where you pay with your card).

### 2.3 Email

1. Create a new Proton or Tuta account on the clean setup.
2. Add an alias service (SimpleLogin or addy.io). **Use one alias per platform.** If one leaks, you can trace which service did it and kill the alias.
3. Don't set a recovery email or phone that links to the old identity.
4. Don't forward to or from old addresses.

### 2.4 Phone number

The hardest part. Order by cleanliness:

1. **Prepaid SIM bought with cash**, registered per local law, kept off while your real phone is on, and never powered up at home or alongside your real phone.
2. A **non-VoIP** number from a privacy-focused carrier, if available in your country.
3. VoIP numbers: many platforms block them. Use only where accepted.
4. Where possible, skip phone verification: prefer platforms that allow email-only sign-up.
5. **Never** reuse this number for anything tied to your real identity (bank alerts, delivery apps, contacts).

### 2.5 Payments

Virtual cards (e.g. privacy-oriented services in your country) or cash-bought prepaid cards. Never your real card, even once. Billing name/address will correlate.

### 2.6 Credentials

1. New password manager, generated unique passwords, new master passphrase (6+ random words).
2. Enable 2FA (hardware key preferred). Keep backup codes encrypted offline.
3. Never copy passwords or seeds from the old vault.

---

## Phase 3: Design the persona

### 3.1 Define it first

Write a short persona sheet: handle, "voice", topics, timezone, interests. Keep it **minimal and separate** from your real traits.

### 3.2 Anti-correlation checklist

| Vector | What to do |
| --- | --- |
| **Username** | Brand new, unrelated, not derived from old handles, not a variant. Check namecheckers to avoid collisions, but don't enter it in sites linked to your accounts. |
| **Stylometry** | Change sentence length, vocabulary, punctuation, emoji, spelling variants (UK/US), capitalization, humor. Avoid pet phrases. Write drafts and compare against your old posts. |
| **Timing** | Different posting windows. Don't mirror old patterns. Avoid posting right after old-account activity. |
| **Social graph** | Don't follow the same people, communities, subreddits or Discords. Graph overlap is the top deanonymization method. |
| **Media** | No old photos or avatars. Use illustrations or AI-generated avatars. Strip metadata: `exiftool -all= *.jpg`. Check backgrounds, reflections, and window views. |
| **Content** | No unique anecdotes, employer, city, school, pets, niche hobbies or events you've mentioned before. |
| **Technical** | Don't reuse bios, link-in-bio services, profile layouts, domain registrars or hosting accounts. |
| **Language** | If you code-switch or have regional slang, suppress it. |

### 3.3 Professional work

If your OSINT/threat-intel work needs credibility, keep the **professional identity separate from the private one**. A pseudonymous research persona is fine, but don't mix it with a private personal account.

---

## Phase 4: Create new accounts (procedure)

For each platform:

1. Open the clean browser profile over the agreed network.
2. Register with the platform-specific **email alias**.
3. Use unique credentials from the new vault.
4. If phone is required, use the prepaid SIM, then **remove the number afterward** if the platform allows.
5. Disable contact sync, "find me by email/phone", discoverability, and location tagging.
6. Set the profile to private or minimal until it's warmed up.
7. Don't connect to the old network on any platform that suggests friends.
8. Log the account in the *new* compartment map (not the old spreadsheet).

**Warm-up:** new accounts that go straight to heavy activity get flagged. Ramp up gradually over days to weeks.

**Compliance note:** a pseudonymous personal account is permitted on most platforms. Don't impersonate others, evade bans, or run multiple accounts where it violates terms.

---

## Phase 5: Migration rules (no handoff)

1. **Never transfer anything.** Followers, contacts, exported posts and "move account" features all create links.
2. **No announcements** from the old identity. No cross-follows. No reposts.
3. **Don't reupload** old posts, photos, or code.
4. **Cooling-off gap:** wait 2+ weeks between the new identity going active and decommissioning the old one, so there's no tidy timeline.
5. For the few contacts who must reach the new identity: share it on **Signal** using the new number, tell them not to mention your real name, tag you, or connect on any platform.
6. **Don't ask mutuals to follow** the new account in public.
7. Build the audience organically.

---

## Phase 6: Decommission old accounts

### 6.1 Decide: delete or lock

- **Delete** if the handle has no impersonation value.
- **Lock/scrub instead** for high-visibility handles (private, blank, bland, 2FA on). A deleted handle can be registered by someone else and used against you.
- Keep a **small number of accounts alive** that can't be safely deleted (e.g. ones needed for banking, government, professional licensing), cleaned and compartmentalized.

### 6.2 Per-account procedure

1. **Export your data** (store encrypted, offline).
2. **Scrub content:** bulk delete posts, comments, photos. Tools: Redact, TweetDelete, Power Delete Suite (Reddit), or platform bulk tools.
3. **Edit** bio, name, location, website, birthday, work and education history.
4. **Remove** photos, tags, and albums. Untag yourself from others' posts.
5. **Remove linked data:** phone, secondary emails, connected apps, saved payment methods, trusted devices.
6. **Change** the account email to a throwaway alias (cuts the link to your main email).
7. **Delete** the account.
8. **Do not log in** during the grace period (typically 14-30 days). Logging in reactivates it.
9. **Verify after 30+ days:** profile 404s, search shows nothing.

### 6.3 Suggested order

1. Lowest-value, easiest accounts first (practice).
2. Accounts that reference others (Facebook, Instagram, LinkedIn).
3. Content platforms (X, Reddit, YouTube, forums).
4. Professional platforms last (GitHub, LinkedIn, conference/CV sites), once the new identity is established.

### 6.4 Then purge OAuth and sessions

Revoke third-party app access, log out of all sessions, and remove saved logins from every browser.

---

## Phase 7: Clean beyond platforms

### 7.1 Search engines

- Google: **Results about you** (contact info removal), **Remove outdated content** tool, SafeSearch removal for sensitive content.
- Bing: **Content removal / Outdated Content** requests. DuckDuckGo uses Bing, so it follows.
- Verify with quoted searches after the removal windows (days to weeks).

### 7.2 Data brokers

1. Opt out manually on each broker (list from Phase 1). Re-check every 3-6 months, since records reappear.
2. Optional paid services: DeleteMe, Optery, Incogni, Kanary.
3. Legal erasure requests: **GDPR Art. 17**, **CCPA/CPRA**, or your local equivalent. Send written requests and keep copies.
4. Use a masked email/phone on opt-out forms.

### 7.3 Archives

- Wayback Machine: email **[info@archive.org](mailto:info@archive.org)** with URLs and ownership proof to request exclusion.
- Archive.today and similar: use their takedown contact.

### 7.4 Technical traces

- **WHOIS:** enable privacy, or update to a proxy contact. Note history already exists.
- **Domains:** transfer to a privacy-oriented registrar; use separate registrar accounts per identity.
- **Git:** switch to noreply addresses (`git config --global user.email "ID+user@users.noreply.github.com"`). Rewriting history (`git filter-repo --mailmap ...`) breaks forks and clones and **doesn't remove existing copies**, so consider it carefully.
- **PGP:** delete the key from `keys.openpgp.org`. SKS-era keyservers retain keys permanently.
- **Package registries / papers / patents:** update author metadata where possible; accept that some stays permanent.

### 7.5 Third-party mentions

Ask friends, colleagues, event organizers and employers to remove or untag photos, speaker bios and team pages. Offer specific URLs.

### 7.6 Breaches

Breach data can't be removed. **Rotate every credential** tied to exposed emails, and stop using those emails.

---

## Phase 8: Verify

Run Phase 1 again against the **old** and **new** identity:

1. Old identifiers: confirm what's gone, what remains, what's cached.
2. New handles: search them alongside old names, emails and phone numbers. Nothing should connect.
3. Compare the writing style of new posts against old archives (you can use a stylometry tool for this).
4. Check the new accounts' discoverability settings.

---

## Phase 9: Ongoing OPSEC

- **Compartment map** (encrypted): which email, number, device, network and vault belongs to which identity.
- **Quarterly re-audit** of both identities.
- **Rules:**
  - Never log into the wrong identity.
  - Never reuse passwords or email aliases.
  - Never type the old name or handle on the new setup.
  - Don't comment about local events or weather at your real location.
- **If you slip** (wrong browser, wrong network, reused photo): assume linkage. **Burn the new account and rebuild** rather than hope nobody noticed.
- **Life-event risk:** moving, new job, and travel are common causes of leakage. Re-check your setup before and after.

---

## Common failure points

| Mistake | Result |
| --- | --- |
| Logging into the new account from the home IP once | IP linkage in platform logs |
| Using the same avatar or bio style | Reverse image / stylometric match |
| Installing the same apps with contact sync | "People You May Know" exposes you |
| Using the new SIM in your real phone | IMEI/IMSI linkage |
| Posting about the same niche event | Contextual correlation |
| Announcing the move to close friends publicly | Social graph link |
| Deleting old accounts instantly after creating the new ones | Temporal correlation |
| Using a real-name payment | Billing record linkage |

---

## Quick checklist

- [ ] Threat model written
- [ ] Footprint audited and inventoried (encrypted)
- [ ] Clean device, network, email, phone, payments, vault ready
- [ ] Persona sheet + anti-correlation rules defined
- [ ] New accounts created with aliases, warmed up
- [ ] Cooling-off period observed
- [ ] Old accounts exported, scrubbed, deleted or locked
- [ ] OAuth and sessions revoked
- [ ] Search engine and broker removals filed
- [ ] Archive/WHOIS/Git/PGP traces handled
- [ ] Verification audit completed
- [ ] Quarterly re-audit scheduled

---

# Appendix: Tools, Websites & Required Items

> Services change, get acquired, or break when platforms change their APIs. **Verify each tool's current status and pricing before relying on it.** Bulk-delete tools in particular break often.

## A. Audit your old footprint

| Tool / site | Use |
| --- | --- |
| **Sherlock**, **Maigret** (GitHub) | Find accounts by username across hundreds of sites |
| **WhatsMyName** (whatsmyname.app) | Username enumeration with a maintained site list |
| **Holehe** (GitHub) | Checks which sites an email is registered on |
| **Epieos** (epieos.com) | Email/phone lookups; shows what a Google account exposes |
| **SpiderFoot** (spiderfoot.net), **theHarvester** | Automated footprint scans of domains, emails, names |
| **Have I Been Pwned** (haveibeenpwned.com) | Breach exposure; sign up for notifications |
| **IntelX**, **DeHashed** | Breach/leak search to see what attackers see (paid tiers) |
| **PimEyes**, **Yandex Images**, **TinEye**, **Google Lens** | Reverse image search for avatars and photos |
| **Wayback Machine** (web.archive.org) | See archived copies of profiles and sites |
| **OSINT Framework** (osintframework.com) | Directory of OSINT resources by category |
| **SecurityTrails**, **ViewDNS**, **DomainTools** | WHOIS history for your domains |
| **Google dorks** | `"your name" -site:linkedin.com`, `"handle" site:pastebin.com`, `intext:"your@email"` |

## B. Account discovery & deletion

| Tool / site | Use |
| --- | --- |
| **Password manager export** (Bitwarden, 1Password, browser) | Fastest inventory of old accounts |
| **JustDeleteMe** (justdeleteme.xyz) | Directory of direct deletion links, with difficulty ratings |
| **Google Security Checkup** / **Third-party access** | Revoke OAuth apps and sessions |
| **Meta Accounts Center**, **Apple ID → Sign in with Apple** | Same for Facebook/Instagram and Apple logins |
| **Redact** (redact.dev) | Bulk delete posts across many platforms (paid tiers) |
| **TweetDelete**, **Semiphemeral** | Bulk delete on X |
| **Power Delete Suite**, **Shreddit** | Bulk edit/delete Reddit history |
| **Platform data export pages** | Back up content before deletion (store encrypted) |

## C. Clean infrastructure

**OS and devices**

| Tool | Use |
| --- | --- |
| **Qubes OS** (qubes-os.org) | Compartmentalized OS; best for strict separation |
| **Whonix** (whonix.org) | Tor-routed VM pair |
| **Tails** (tails.net) | Amnesic live USB for sensitive sessions |
| **GrapheneOS** (Pixel phones) | Hardened phone OS for the new identity |
| **VirtualBox / KVM** | Plain VMs if Qubes isn't an option |

**Browsers and leak checks**

| Tool | Use |
| --- | --- |
| **Tor Browser**, **Mullvad Browser** | Anti-fingerprinting browsers |
| **Firefox + arkenfox user.js** | Hardened Firefox profile |
| **uBlock Origin** | Content blocking |
| **Cover Your Tracks** (coveryourtracks.eff.org), **AmIUnique**, **BrowserLeaks** | Fingerprint tests |
| **ipleak.net**, **dnsleaktest.com**, **mullvad.net/check** | IP, DNS and WebRTC leak tests |

**Network**

| Tool | Use |
| --- | --- |
| **Mullvad**, **IVPN** | No-logs VPNs with anonymous payment options |
| **Proton VPN** | Free tier and Tor-over-VPN features |
| **Tor** (torproject.org) | Anonymity network; use bridges if Tor is blocked |

**Email and aliasing**

| Tool | Use |
| --- | --- |
| **Proton Mail**, **Tuta** | Privacy-focused email |
| **SimpleLogin**, **addy.io** | One alias per site |
| **Firefox Relay**, **DuckDuckGo Email Protection** | Alternative alias services |

> Don't use Apple "Hide My Email" or Google aliases for the new identity. They're tied to accounts you already own.

**Messaging**

| Tool | Use |
| --- | --- |
| **Signal** (signal.org) | Use usernames so people never see your number |
| **Session**, **SimpleX**, **Element** | No-phone-number alternatives |

**Credentials**

| Tool | Use |
| --- | --- |
| **Bitwarden**, **KeePassXC**, **Proton Pass** | New, separate password vault |
| **YubiKey**, **Nitrokey** (buy two) | Hardware 2FA, with a spare |
| **Aegis** (Android), **Ente Auth**, **KeePassXC TOTP** | TOTP apps; avoid anything tied to your old cloud account |

**Payments**

| Tool | Use |
| --- | --- |
| **Privacy.com** (US only), **Revolut** disposable cards (where available) | Virtual cards that don't carry your real card number |
| **Cash-bought prepaid cards** | Check activation and KYC rules; many require identity |
| **Monero** | Only if you understand the legal and exchange implications |

**Phone numbers**

| Option | Notes |
| --- | --- |
| **Prepaid SIM, bought with cash** | Cleanest; SIM registration laws vary by country |
| **Privacy-oriented eSIM providers** | Can accept crypto; check KYC rules and device linkage |
| **VoIP apps (e.g. MySudo)** | Many platforms block VoIP numbers |

## D. Persona hardening

| Tool | Use |
| --- | --- |
| **ExifTool** (exiftool.org) | `exiftool -all= file.jpg` strips metadata |
| **mat2** | Metadata anonymisation toolkit for many file types |
| **DiceBear** (dicebear.com), **Boring Avatars** | Generated non-photo avatars |
| **Signal's blur tool**, **GIMP** | Redact faces or backgrounds |
| **stylo** (R package), **JStylo/Anonymouth** | Stylometry analysis; compare your new writing with your old |
| **Namechk-type sites** | Check handle availability. Use from the clean setup only |

## E. Cleanup and removal

| Tool / site | Use |
| --- | --- |
| **Google "Results about you"** (myactivity.google.com/results-about-you) | Remove personal contact info from results |
| **Google "Remove outdated content"** tool | Refresh or remove stale cached results |
| **Bing Content Removal** (Bing Webmaster Tools) | Outdated page and personal-info removal |
| **Wayback exclusion** | Email [info@archive.org](mailto:info@archive.org) with URLs and proof of ownership |
| **DeleteMe**, **Optery**, **Incogni**, **Kanary**, **Privacy Bee** | Paid broker opt-out services |
| **Big Ass Data Broker Opt-Out List** (GitHub: yaelwrites) | Free manual opt-out guide |
| **Privacy Rights Clearinghouse** | Data broker registry and guidance |
| **Legal erasure requests** | GDPR Art. 17 (EU/UK), CCPA/CPRA (California), DPDP Act (India), or your local law; send written requests and keep copies |

## F. Technical traces

| Tool | Use |
| --- | --- |
| **git filter-repo**, **BFG Repo-Cleaner** | Rewrite Git history (breaks forks; doesn't remove copies) |
| **GitHub email privacy settings** | "Keep my email private" + "Block pushes that expose my email" |
| **WHOIS privacy** (via your registrar) | Hide registrant details going forward |
| **keys.openpgp.org** | Manage or delete your PGP key |

## G. Secure storage and notes

| Tool | Use |
| --- | --- |
| **VeraCrypt**, **Cryptomator**, **age** | Encrypted containers for exports and the inventory |
| **CryptPad** (cryptpad.org) | Encrypted collaborative docs and sheets |
| **Paper notebook** | Offline compartment map and recovery codes (store separately from your real-ID documents) |

## H. Monitoring

| Tool | Use |
| --- | --- |
| **Google Alerts** | Alerts on your old name, handles and email (set up from the **old** identity only) |
| **HIBP notifications** | Alerts when emails appear in new breaches |
| **Quarterly re-run** of Sherlock/Maigret + search dorks | Catch re-listings and re-scraped data |

## I. Things you'll need to buy or prepare

- [ ] Second laptop or spare device (no old accounts ever signed in)
- [ ] USB drives (Tails installer, encrypted backups)
- [ ] Two hardware security keys
- [ ] Prepaid SIM + a SIM-capable spare phone (cash purchase)
- [ ] Cash for the VPN subscription and prepaid cards
- [ ] Virtual or prepaid payment method
- [ ] Paper notebook for the compartment map and recovery codes
- [ ] Time budget: roughly 2-4 weeks for a careful migration, then regular maintenance

## J. Learning resources

- **Privacy Guides** (privacyguides.org): vetted tool recommendations
- **EFF Surveillance Self-Defense** (ssd.eff.org): practical threat modeling and guides
- **Michael Bazzell, *Extreme Privacy*** and inteltechniques.com: detailed privacy-by-design and identity separation
- **Whonix and Qubes documentation**: compartmentalization patterns

## K. Minimal starter kit

If you want the smallest viable set: **Mullvad VPN + Tor Browser or hardened Firefox + Proton Mail + SimpleLogin + Bitwarden + hardware key + prepaid SIM + Signal (usernames) + ExifTool + Sherlock/Maigret + JustDeleteMe + a broker opt-out service.**
