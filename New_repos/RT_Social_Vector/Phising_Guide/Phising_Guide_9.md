# Phishing Defensive Measures: Unified Playbook

`This still requires practical expansion.` 

*Recommendations and hands-on implementation, kept together under each control.*

**How to read this.** Every control has two labels side by side: **Recommend** (policy, process and awareness decisions the client owns) and **Implement** (hands-on steps the testing team sets up, validates or runs). Items marked ★ are implementation steps that were missing from the earlier draft.

- **Part 1: Inbound attack channels.** Email, SMS, voice, and cross-channel correlation.
- **Part 2: People, identity and response.** Physical, sign-in, reporting and containment, help desk, and ongoing measurement.

---

## Part 1: Inbound Attack Channels

### 1. Email

#### 1.0 Authentication baseline (own domain)

**Recommend**
- Policy: DMARC at `p=reject` on every owned domain, including parked domains that never send mail (publish `v=spf1 -all` and a reject policy on those).

**Implement**
```bash
dig TXT _dmarc.clientdomain.com
dig TXT clientdomain.com
dig TXT selector1._domainkey.clientdomain.com   # repeat per DKIM selector
```
- Confirm `p=reject` (not `none` or `quarantine`), and that the `rua` mailbox is monitored and its reports are reviewed.
- Confirm SPF covers all authorized senders, has no `+all`, and ends in `-all` (`~all` is only acceptable when DMARC is at reject).
- Confirm DKIM selectors are active for every sending system.
- ★ Inventory every third-party sender (marketing, HR, ticketing) and confirm each one aligns with DMARC before moving to reject.

#### 1.1 Attacker-owned domains with valid authentication

**Recommend**
- Policy: authentication passing is not trust. Brief management that a DMARC pass only proves the sender controls the domain it sent from.
- Domain-age threshold in the gateway: a first-time sender from a domain registered in the last 30 days gets a warning banner or is quarantined for review.
- Sender-novelty flagging: label a first-ever message from a never-seen domain "first contact" before delivery.
- Awareness: show staff real attacker-owned domains that pass SPF, DKIM and DMARC.

**Implement**
- Configure the age and first-contact rules, then send test messages from a newly registered look-alike domain and confirm the warning or quarantine fires.
- ★ Check a domain's age quickly when validating:
```bash
curl -s https://rdap.org/domain/example.com | jq -r '.events[] | select(.eventAction=="registration") | .eventDate'
```

#### 1.2 Compromised legitimate accounts and thread hijack

**Recommend**
- Policy: out-of-band verification for any sensitive email request (payment changes, credential resets, access grants), even from a known contact.
- Per-sender behavioral baselines: alert when a known sender appears from a new IP, ASN or mail client, or when a new external recipient joins an established thread.
- Payment-detail change rules: hold any external message containing a bank account number, UPI ID or payment link for human review.
- Awareness: show how a reply inside an existing chain can come from a different address with the same display name.

**Implement**
- ★ Build a mail-flow rule for external senders that holds messages matching an IFSC pattern (`[A-Z]{4}0[A-Z0-9]{6}`), phrases such as "account number", "UPI" or "bank details changed", or a payment link. Test it with fake details from an external test mailbox and confirm the hold.
- ★ Validate the baseline alerts: sign in to a test mailbox from a new IP/ASN, and add a new external recipient to a test thread. Confirm both raise alerts.

#### 1.3 Lookalike domains and display-name spoofing

**Recommend**
- Protected-name list covering executives, brands and trusted partners. Flag or reject mail where the display name matches but the sending domain does not.
- Confusable-character normalization in the gateway: treat domains containing homoglyphs (Cyrillic о, Greek ο, fullwidth forms) of your own domain as suspicious.
- Continuous domain monitoring: monthly sweeps for typosquats, combosquats, homoglyphs and suffix variants, reported to the client's security team, plus Certificate Transparency alerts for any new certificate containing the organization's name.
- Awareness: show the homoglyph table and a side-by-side of two domains that look identical on screen but are different strings.

**Implement**
```bash
pip install dnstwist
dnstwist --registered clientdomain.com
dnstwist --registered --mxcheck clientdomain.com
dnstwist --registered --format csv clientdomain.com > results.csv
dnstwist --registered --dictionary words.txt clientdomain.com
```
- Review output with active MX records first: those domains can already send mail and are the highest immediate risk.
- Schedule the sweep monthly and deliver results to the client's security team.
- ★ Certificate Transparency check (use a CT monitor or vendor feed for real alerting):
```bash
curl -s "https://crt.sh/?q=%25clientdomain%25&output=json" | jq -r '.[].name_value' | sort -u
```
- ★ Build the protected-name list (executives, finance, brands, partners) and load it into the gateway.
- Test a homoglyph domain and confirm confusable detection fires. Test a display-name match from an external domain and confirm the protected-name rule fires.

#### 1.4 Scanner-versus-human cloaking

**Recommend**
- Policy: "could not be scanned" or "clean at delivery time" is never a pass verdict. Treat incomplete inspection as elevated risk.
- Require the gateway to fetch URLs from several scan profiles (IP ranges, user agents, locales) and flag divergent responses.
- Rescan at click time, not only at delivery.

**Implement**
- Verify URL rewriting applies to all messages, including allow-listed ones.
- ★ Point a benign test URL at a page that returns different content per user agent, and confirm the multi-profile fetch flags the divergence. A quick manual comparison:
```bash
curl -s -A "Mozilla/5.0 (Windows NT 10.0) Chrome/126.0" https://test.example/page | sha256sum
curl -s -A "Googlebot/2.1" https://test.example/page | sha256sum
```

#### 1.5 Delayed weaponization

**Recommend**
- Retroactive rescanning of delivered links every 4 to 8 hours for the first 48 hours. A verdict change automatically quarantines every message containing the link.
- Click-time URL rewriting so every click passes through a live scanner.
- Policy: a documented runbook for removing all delivered copies and notifying recipients when a verdict flips from clean to malicious.

**Implement**
- Deliver a message with a link to a benign page, then switch the destination to a vendor-supplied test-malicious URL. Confirm an alert or automated removal fires.
- ★ Time the whole cycle (verdict change to removal) and record it as the baseline for the removal-speed metric in Part 2, section 7.

#### 1.6 Hosting on trusted platforms

**Recommend**
- Page-level analysis beyond domain reputation: render the page and compare logo, brand name and content against the host domain. A bank-branded credential form on a Google Sites page should score worse than the same form on the bank's own domain.
- Policy: publish the list of trusted-platform share links the organization routinely uses (Google Drive, Notion, Dropbox). An unexpected share link from an unknown sender is not safe because of where it is hosted.

**Implement**
- ★ Host a benign, client-branded dummy login page (no credential capture) on a trusted platform, send the link, and confirm the scanner scores it above a normal share link.
- ★ Compile the approved-platform list from mail logs (most-shared hosting domains) and publish it to staff.

#### 1.7 CAPTCHA gates, redirect chains and trusted-domain redirects

**Recommend**
- Policy: a link from an unknown sender that shows a CAPTCHA or a long redirect chain is elevated risk. Report it, do not click through.
- Judge on the full chain: resolve every hop and base the verdict on the final destination. Flag chains where the destination is unrelated to the first hop.
- Score the recency of the first-hop share or short link. A Dropbox link created three minutes before delivery is a different risk from one shared six months ago.

**Implement**
- Document what the gateway does with an unresolved redirect chain and confirm it never defaults to pass.
- ★ Inspect chains and test depth handling:
```bash
curl -sIL -o /dev/null -w '%{url_effective} hops=%{num_redirects}\n' https://short.example/abc
```
- ★ Send a message with a harmless 3+ hop redirect ending on an unrelated test page. Confirm it is flagged, and that a link created minutes earlier scores higher than an older one.

#### 1.8 LLM-paraphrased unique text

**Recommend**
- Move from keyword filters to embedding-based intent classifiers so the intent, not the wording, triggers the alert.
- Cluster by sender infrastructure (hosting, nameserver, registrar, certificate), not by content. Many differently worded messages from one infrastructure are one campaign.

**Implement**
- ★ Write 20 paraphrases of one benign test lure and confirm the classifier gives them the same verdict.
- ★ Cluster a message-trace export by infrastructure:
```python
from collections import defaultdict

def cluster(msgs, keys=('sender_ip_24', 'asn', 'registrar', 'nameserver', 'tls_issuer')):
    groups = defaultdict(list)
    for m in msgs:
        for k in keys:
            if m.get(k):
                groups[(k, m[k])].append(m['id'])
    return {g: ids for g, ids in groups.items() if len(ids) > 1}
```

#### 1.9 Image, QR code and PDF-embedded content (quishing)

**Recommend**
- OCR and QR decoding on all inbound mail. A decoded QR URL is scanned like a plain hyperlink, including click-time rewriting.
- Policy: any QR code from a first-time or unknown sender is elevated risk by default, because the destination can be weaponized after delivery.
- Awareness: scanning with a personal phone camera bypasses corporate endpoint and network controls. Treat that gap as a deliberate evasion of monitoring.

**Implement**
- ★ Generate a benign QR code to a test URL:
```bash
pip install "qrcode[pil]"
python -c "import qrcode; qrcode.make('https://test.clientdomain.example/qr-test').save('qr_test.png')"
```
- ★ Send it three ways (inline image, image attachment, embedded in a PDF) from a first-time sender. Confirm the gateway decodes it, rewrites or scans the URL, and flags the first-time sender.

#### 1.10 Hidden text and zero-width character abuse

**Recommend**
- Normalize before any classification: NFKC form, strip zero-width characters (U+200B, U+200C, U+200D, U+FEFF) and the Tags block (U+E0000 to U+E007F).
- Policy: classifier input is always the normalized, rendered view of the message, never the raw byte stream.

**Implement**
- Deploy this sanitizer on all inbound text before it reaches classifiers (★ updated: it now ends in NFKC as recommended, and `allow_joiners=True` really keeps joiners):
```python
import re, unicodedata
from collections import Counter

TAG_BLOCK = range(0xE0000, 0xE0080)
BIDI_CONTROLS = {0x200E, 0x200F, 0x061C} | set(range(0x202A, 0x202F)) | set(range(0x2066, 0x206A))
EXOTIC_SPACES = {0x00A0, 0x0085, 0x1680, 0x2028, 0x2029, 0x202F, 0x205F, 0x3000} | set(range(0x2000, 0x200B))
INVISIBLE_OTHER = {0x034F, 0x115F, 0x1160, 0x17B4, 0x17B5, 0x2800, 0x3164, 0xFFA0, 0x00AD, 0x180E}
JOINERS = {0x200C, 0x200D}

def clean(text: str, allow_joiners: bool = False):
    report = Counter()
    out = []
    for ch in text:  # inspect raw characters first so the report is accurate
        cp = ord(ch)
        if ch in '\n\t':
            out.append(ch)
        elif cp in EXOTIC_SPACES:
            out.append(' ')
            report['exotic_space'] += 1
        elif cp in TAG_BLOCK:
            report['tag_char'] += 1
        elif cp in BIDI_CONTROLS:
            report['bidi_control'] += 1
        elif allow_joiners and cp in JOINERS:
            out.append(ch)
        elif cp in INVISIBLE_OTHER or cp in JOINERS:
            report['invisible'] += 1
        elif unicodedata.category(ch) in ('Cc', 'Cf', 'Co', 'Cs', 'Cn'):
            report['control_or_private'] += 1
        else:
            out.append(ch)
    result = unicodedata.normalize('NFKC', ''.join(out))
    result = re.sub(r'[ ]{2,}', ' ', result).strip()
    return result, dict(report)
```
- Alert on any non-zero `tag_char` or `bidi_control` count from inbound or form-originated data.
- ★ Test: send the same benign lure three ways (plain, zero-width characters inside keywords, hidden Tags-block text). Confirm the classifier verdict is identical for all three.

#### 1.11 Link-free openers, off-channel moves and multi-message pretext splitting

**Recommend**
- Conversation-trajectory scoring: an unknown sender who solicits a reply with no link or attachment is still scored for intent, especially if the follow-up adds urgency, a payment request, a credential ask or a channel-switch request.
- Score threads and sequences, not only single messages. A benign opener followed by a pressure message to the same recipient is one attack.
- Policy and awareness: an innocent first message followed by a sensitive request is a recognized social engineering technique. Reporting the first message is correct and valuable.

**Implement**
- ★ Create a SIEM correlation rule: external first-contact sender, then recipient replies, then a follow-up from the same sender within 7 days containing urgency, payment or credential keywords or a request to move channels. Raise an alert at the second message.
- ★ Test with a scripted two-message exchange from an external test mailbox and confirm the sequence alert fires even though neither message alone is suspicious.

#### 1.12 Password-protected archives

**Recommend**
- Policy: publish an inbound file-type policy for encrypted archives, disk images and shortcut files, so recipients know what to do with a protected archive from an external sender.
- Gateway extracts passwords from message bodies, attempts to unpack protected archives, and enforces nesting-depth limits.
- Review the permitted container-type list quarterly and add newly publicized container exploits. The list has a short shelf life.

**Implement**
- ★ Download the official EICAR test file from eicar.org, then build test archives:
```bash
zip -P infected eicar_protected.zip eicar.com
zip -P outer nested.zip eicar_protected.zip
```
- Send both (password in the message body for one, absent for the other) and confirm the gateway unpacks and detects, or quarantines what it cannot open.
- ★ Add a recurring quarterly calendar item to review the container-type list.

#### 1.13 Sandbox-aware payloads

**Recommend**
- Policy: a sandbox verdict of "no activity observed" is never clean. Evasive behavior is itself a risk signal.
- Require detonation environments that simulate real user interaction (mouse movement, scrolling, realistic timing), not bare static virtual machines.

**Implement**
- ★ Confirm in the vendor configuration that user-interaction simulation is enabled.
- ★ Submit a benign file that sleeps for several minutes and does nothing, and confirm the verdict is inconclusive or suspicious, not clean.

#### 1.14 Filter probing

**Recommend**
- Policy: uniform responses. The gateway must not reveal blocking thresholds to external senders through different bounce codes, timing or delivery confirmations.
- Run shadow models on live traffic without affecting verdicts, so new evasion patterns are captured before the primary model is tuned to them.

**Implement**
- ★ From an external test host, send one message that will be blocked and one that will be delivered, and compare SMTP replies and timing:
```bash
swaks --to user@clientdomain.com --from probe@testdomain.example --server mx.clientdomain.com --header "Subject: probe-1"
```
- ★ Run the candidate classifier in log-only mode on mirrored traffic and review its disagreements with the production verdicts weekly.

### 2. SMS

#### 2.1 Fresh throwaway domains in SMS

**Recommend**
- Resolve and score every link in an inbound SMS at the SMS gateway or mobile threat defense (MTD) tool. Age the domain from its registration date, not from when the gateway first saw it.
- Cluster by infrastructure, not domain. Two messages with differently named domains on the same hosting IP, nameserver or certificate are one campaign and are blocked together.

**Implement**
- ★ Configure link resolution and domain-age scoring in the MTD or SMS gateway, and send a test SMS linking to a freshly registered test domain to confirm it is scored as high risk.
- ★ Reuse the infrastructure clustering script from section 1.8 on SMS link data, and confirm two test domains sharing one IP are blocked together.

#### 2.2 Rotating or spoofed sender IDs and grey routes (India DLT)

**Recommend**
- Policy: sender IDs used in authorized simulations must be registered under the client's own Principal Entity (PE) on DLT. Any test SMS from an unregistered sender ID or a grey route is out of scope for the engagement.
- Carrier-side DLT validation should use fuzzy-hash matching against registered templates, not only exact matching, to close small variations that still look like the registered template to the recipient.

**Implement**
- Confirm the client's PE registration is active and propagated across Jio, Airtel, Vi and BSNL.
- List all registered sender ID headers and confirm each maps to an authorized aggregator.
- Confirm the aggregator enforces template matching server-side and that an unregistered sender ID returns `FAILED_INVALID_SENDER` instead of going out on a grey route.
- Confirm every smishing test sender ID is registered under the client's PE before any test SMS is sent.
- ★ Send a registered template with a trivial variation (extra space, look-alike character) from a test header and confirm the aggregator rejects it or flags the mismatch.

#### 2.3 Low-and-slow personalized SMS

**Recommend**
- Intent classification on inbound SMS instead of volume rules alone. One well-crafted message asking for an OTP or credential is as dangerous as a campaign of thousands.
- Number-age and line-type checks: flag SMS from VoIP numbers, or from lines registered in the last 30 days by an entity the recipient has no relationship with.

**Implement**
- ★ Enable intent classification and line-type lookup in the MTD stack (or integrate a carrier lookup service), then send a single test message from a VoIP test number that asks for a code. Confirm it is flagged.

#### 2.4 Obfuscated text in SMS

**Recommend**
- Normalize all SMS content before classification: strip zero-width characters, map non-GSM characters, apply NFKC. A keyword filter on unnormalized text can be bypassed with one invisible character.

**Implement**
- ★ Run inbound SMS through the `clean()` sanitizer from section 1.10 before the classifier, and confirm a test message with zero-width characters inside a keyword gets the same verdict as the plain version.

#### 2.5 Reply-to-unlock

**Recommend**
- Policy and awareness: a message from an unknown number asking you to reply with a word or code before a link is sent is a recognized pattern. Do not reply; report the number.
- Flag unknown senders that solicit replies, with no existing contact relationship, as a risk feature in the MTD stack, not only messages containing links.

**Implement**
- ★ Enable the reply-solicitation feature in the MTD policy and test it with a scripted "reply YES to continue" message from an unknown test number. Confirm the flag fires and that staff know how to report it.

#### 2.6 Engagement operations: GSM-7, segments and throughput (testing team)

These steps apply to the testing team before any test or real SMS campaign.

**Implement**
- Validate the rendered text after personalization against the GSM-7 character set, and count segments. ★ Updated: segment counts for single-part messages and emoji (which use two UTF-16 units) are now correct.
```python
GSM7 = set('@£$¥èéùìòÇ\nØø\rÅåΔ_ΦΓΛΩΠΨΣΘΞÆæßÉ !"#¤%&\'()*+,-./0123456789:;<=>?'
           '¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà')
GSM7_EXT = set('^{}\\[~]|€')

def is_gsm7(text):
    return all(ch in GSM7 or ch in GSM7_EXT for ch in text)

def count_segments(text):
    if is_gsm7(text):
        length = len(text) + sum(1 for ch in text if ch in GSM7_EXT)
        return 1 if length <= 160 else -(-length // 153)
    units = len(text.encode('utf-16-le')) // 2   # UCS-2 / UTF-16 code units
    return 1 if units <= 70 else -(-units // 67)
```
- A single emoji or Devanagari character in an otherwise GSM-7 message flips the whole message to UCS-2 and can triple the segment count and cost.
- Plan throughput before every test send:
```text
safe_rate        = route_cap x 0.60                       (segments per second)
daily_capacity   = safe_rate x 3600 x window_hours
total_segments   = total_messages x segments_per_message
fastest_finish   = ceil(total_segments / daily_capacity)  days
batch_interval   = (batch_size x segments_per_message) / safe_rate  (seconds)
```
- Apply pause gates after every batch: opt-out rate above 1%, failure rate above 5%, or delivery rate below 90% means stop and investigate before the next batch.

---

### 3. Voice and Call

#### 3.1 Caller ID spoofing

**Recommend**
- Policy: callback verification is mandatory for any phone request involving account access, credential changes, payment authorization or sensitive disclosure. The callback goes to a number already on file, never one the caller supplies.
- Brief staff that "unverified" or missing caller ID labels are a signal, not a nuisance. Publish a short internal guide on what STIR/SHAKEN labels mean in practice.
- For clients with enterprise telecom relationships: request carrier-level STIR/SHAKEN attestation reporting, which shows what share of inbound calls carry full, partial or no attestation without a live spoof test.

**Implement**
- Request attestation statistics from the carrier: the proportion of inbound calls with full (A), partial (B) and no (C or absent) attestation.
- Verify whether the telephony platform shows "unverified" labels to end users on calls with missing or partial attestation.
- Document the current state as the baseline for the finding, whether or not a live spoof test is possible.
- Carrier-sanctioned spoof test, if the client has an enterprise carrier relationship. Prepare for the carrier approach:
  - A signed authorization letter naming the specific source and destination numbers.
  - The requested test window and a single point of contact for the carrier.
  - Confirmation that no numbers outside the client's own organization are involved.
  - The request must come from the client's own account team, not from the testing firm.
- Document the outcome either way. A declined request is itself a finding about the client's carrier relationship.

#### 3.2 Number rotation and VoIP lines

**Recommend**
- Line-type checking on inbound calls where the platform supports it. A VoIP or prepaid number calling the help desk or a sensitive department is a different risk from a mobile or landline call.
- Monitor call patterns at the prefix and carrier level, not only per number. A new number calling the help desk five times in a week from the same VoIP ASN as a known abused number is a campaign.

**Implement**
- ★ Enable line-type lookup in the telephony platform and tag VoIP and prepaid callers in the call record.
- ★ Build a weekly report from call detail records grouped by carrier and number prefix, listing help-desk and finance-line calls from VoIP ranges, with an alert when a prefix appears above a set threshold.

#### 3.3 Low-volume targeted vishing

**Recommend**
- Policy: mandatory documented out-of-band verification for all help desk requests involving MFA changes, password resets and account unlocks, including a callback to the number on file.
- Brief the outsourced help desk vendor on vishing patterns in the annual vendor review, and confirm the vendor's procedure matches the client's written policy with adherence measured and reported.

**Implement**
- The enforcement, tooling and measurement steps for this control sit together in Part 2, section 8 (Help Desk and Verification).
- ★ Add vishing patterns (urgent executive caller, MFA reset request, unfamiliar callback number) to the vendor's annual briefing pack and confirm delivery in writing.

#### 3.4 Voice cloning

**Recommend**
- Policy: shared code words between trusted parties for any call involving a sensitive action. A pre-agreed, randomly generated phrase defeats a voice clone, which cannot know it.
- Policy: tell all staff and the help desk that a recognizable voice is not verification. Any access, payment or sensitive-data request goes through the documented out-of-band procedure however familiar the caller sounds.
- Add a help desk script question that only the real employee could plausibly answer from internal records, not from public sources, and that is not the caller's name or employee ID.

**Implement**
- ★ Generate code phrases with a cryptographic random source (distribute them in person or through a verified channel, rotate them on a schedule):
```python
import secrets
words = [w.strip() for w in open('wordlist.txt') if w.strip()]
print('-'.join(secrets.choice(words) for _ in range(3)))
```
- ★ Assign code phrases to the highest-risk pairs first (executives with finance, executives with their assistants, help desk escalation contacts), then add the internal-records question to the help desk script.

#### 3.5 One-ring (Wangiri) campaigns

**Recommend**
- Awareness notice: if you get a missed call from an unfamiliar international or premium number, do not call back without checking with IT or your manager.
- Configure the telephony platform to warn before completing a call to a premium-rate number and to block outbound calls to known premium-rate ranges.

**Implement**
- ★ In the PBX or carrier portal, add outbound dial-plan deny rules for known premium-rate ranges, require approval for international dialing on non-privileged lines, and ask the carrier to bar premium-rate numbers at network level.
- ★ Test by dialing a test premium-rate number from a normal extension and confirm the call is blocked or warned.

---

### 4. Cross-Channel Handoff

**Recommend**
- Build an entity graph linking phone numbers, domains, email addresses and file hashes. A block or alert on one channel propagates to the others automatically when they share infrastructure.
- Awareness: a call followed immediately by an SMS link or an email asking you to install software is a recognized pattern. A follow-up on a second channel after an unusual call is a reason to report, not to trust the original call.

**Implement**
- ★ Start with a simple graph in your SIEM or a script, then connect it to gateway block lists:
```python
import networkx as nx

G = nx.Graph()

def link(a, b):
    G.add_edge(a, b)

def blast_radius(ioc):
    return nx.node_connected_component(G, ioc)

# link('phone:+91XXXXXXXXXX', 'domain:bad.example')
# link('domain:bad.example', 'ip:203.0.113.7')
# blast_radius('phone:+91XXXXXXXXXX') -> everything to block together
```
- ★ Test it with one test indicator per channel: block the domain from an email report and confirm the SMS gateway and the call-screening list pick up the linked identifiers.
- ★ Include one cross-channel scenario (call, then SMS link, then email) in the annual tabletop in Part 2, section 9.

---
