`This write-up is not complete; it is still under active construction.`

# PHISHING DEFENCE EVASION GUIDE 2 

`This is a continuation of guide 1. I would ask you to read that before this.`

## 2. SMS

### 2.1 Fresh throwaway domain per message
- **What it is:** Each campaign or message uses a newly registered domain.
- **Defence it blinds:** Blocklists and reputation feeds.
- **Tell:** Domain age measured in days, generated-looking names, many domains on the same host.
- **Fallback signals:** Cluster by infrastructure (host, nameserver, registrar, certificate), treat new domains in SMS as suspicious by default, block by cluster.

`I won't recommend this as this probably won't work either way and will get blocked quite easily.`

### 2.2 Rotating or spoofed sender IDs and grey routes
- **What it is:** Sender names or numbers change constantly, or messages travel through routes that skip sender registration.
- **Defence it blinds:** Sender reputation and registration checks.
- **Tell:** Template and header mismatch, route anomalies, sender not registered for that content.
- **Fallback signals:** Template fuzzy-hash clustering, validation against registered templates (DLT in India), route-consistency checks.

`Here is how to spof Sender ID ethically for engagement purposes.`

#### Sender ID Customization for Authorized Smishing Engagements

*A compliant process for sender ID customization in client-authorized smishing simulations (India-focused)*

This guide covers only the compliant path: registering a custom sender ID through legitimate telecom/regulatory channels. It does not cover grey routes, unregistered international interconnects, or any method whose purpose is evading carrier or regulator visibility — those fall outside what should be used in an authorized engagement, because authorization from the client doesn't extend to the telecom regulator or carriers who are the other parties with a stake in that traffic.

##### 1. Pre-Engagement: Authorization Before Anything Technical

Before any sender ID or gateway work begins, lock down written authorization that specifically covers SMS/smishing testing — a generic pentest RoE (Rules of Engagement) that only mentions network or email testing is not sufficient cover.

**Required in the RoE / Statement of Work:**

- Explicit scope language authorizing simulated SMS phishing ("smishing") as an in-scope activity, not just phishing generally.
- Named signatory with authority to approve the test on behalf of the client organization (typically CISO, Head of Security, or equivalent — not just a project manager).
- Confirmation of **who owns the telecom/DLT relationship**: does the client already have a Principal Entity (PE) registration and an existing aggregator relationship, or does this need to be set up for the engagement?
- Target population definition (which employee list, which numbers) and a clear statement that the client has the right to test those numbers (internal staff numbers under acceptable-use policy, not third parties).
- Sender ID(s) to be used, agreed in advance and in writing — not decided unilaterally by the testing team.
- Message content/templates to be used, reviewed and approved by the client's legal or compliance function before registration.
- A named point of contact at the client who can vouch for the campaign to carriers, the regulator, or law enforcement if the traffic gets flagged mid-test.
- Explicit stop/abort conditions and a kill-switch process.

Keep this signed documentation on file for the duration of the engagement plus whatever retention period your contract and local law require — it's your evidence that the traffic was authorized if a carrier, regulator, or even a suspicious employee escalates it.

##### 2. Regulatory Groundwork — India (TRAI DLT)

India's Distributed Ledger Technology (DLT) platform, run by TRAI-registered telecom operators, is the compliant channel for any commercial or transactional-style SMS, including simulated smishing content. The grey-route alternative exists specifically to dodge this layer — which is exactly why it's off the table here.

**Steps:**

1. **Principal Entity (PE) Registration.** The entity actually sending the SMS — usually the client, sometimes the testing firm operating under the client's umbrella — must be registered as a Principal Entity with a DLT operator (each major telecom operator, e.g., Jio, Airtel, Vi, BSNL, runs its own DLT platform, but registration typically propagates). If the client already sends transactional/promotional SMS, they likely have this already; confirm rather than assume.
2. **Header (Sender ID) Registration.** Register the specific sender ID/header you intend to use (e.g., "ITSUPP", "SECTEAM") under the client's PE. Headers are typically 6 characters, alphanumeric, and must be registered per-telecom-operator. This is the step that makes the sender ID legitimate rather than spoofed — it's tied to the PE, auditable, and traceable back to the authorized client.
3. **Content Template Registration.** The exact message template (with variable fields marked) must be registered and approved before it can be sent at scale. This also means your phishing-sim message text gets a compliance review pass before launch, which is a reasonable secondary control.
4. **Consent / Scrubbing Considerations.** DLT also enforces consent-based filtering for promotional content. Transactional/service-category headers (which a simulated "IT-Support" or "IT-Security" alert would plausibly fall under) have different consent requirements than promotional ones — work with the aggregator to classify the campaign correctly rather than guessing, since misclassification is itself a compliance issue independent of the phishing-sim framing.

This process routes through the client's own telecom relationship, so client cooperation is genuinely required — this isn't something a testing team can stand up unilaterally, by design.

##### 3. Choosing a Licensed Aggregator

If the client doesn't already have an SMS gateway relationship suitable for this, select a licensed bulk-SMS aggregator that:

- Is itself registered and compliant on the DLT platform (ask for their PE/DLT registration details directly, don't take it on faith).
- Has experience with security-awareness/phishing-simulation traffic specifically — several aggregators explicitly support this use case and know how to classify and register it correctly.
- Can provide delivery logs and audit trails tied to the registered header and template, so you have a defensible record of exactly what was sent, to whom, and when.
- Will not onboard the campaign without seeing the signed client authorization — treat an aggregator willing to skip this as a red flag, not a convenience.

##### 4. Picking and Registering the Sender ID

- **Pick a believable but internal sender ID** (e.g., "ITSUPP", "SECTEAM") rather than a real external brand name, even with client sign-off — this avoids trademark and third-party notification issues a brand owner outside the engagement never agreed to.
- **Register per operator, not once** — a header approved under one telecom operator's DLT platform doesn't automatically propagate to the others; confirm registration on each operator the target population uses.
- **Lock the wording before you register** — the content template tied to the sender ID must match what's actually sent, word for word; changing the message after registration means re-approval, so finalize copy with the client first.

##### 5. What the Sender ID Field Actually Is in Practice

The sender ID isn't a separate technical mechanism — it's a plain field in the aggregator's dashboard or API, same as any other form field.

- **Dashboard:** a dropdown that only lists headers already DLT-approved and mapped to the client's account — you select, you don't type.
- **API:** a parameter like senderId, from, or sender in the send request, carrying that same approved header value.
- **Server-side enforcement:** the aggregator checks the value against the registered header (and the message against the registered template) before sending — an unmapped or unapproved sender ID fails outright (commonly surfaced as FAILED_INVALID_SENDER), it doesn't silently go out.

There's no separate field or setting for an unregistered sender ID. Any send using one skips this validation layer entirely, which puts it back in grey-route territory rather than being a feature of the compliant dashboard.

---

*This guide assumes a genuinely client-authorized engagement with legal/compliance sign-off. It is a process reference, not legal advice — confirm current DLT/TRAI requirements and any jurisdiction-specific telecom regulations with the client's legal team or a licensed aggregator before executing, since registration requirements and fee structures do change.*

### 2.3 Low-and-slow personalized messages
- **What it is:** Small numbers of tailored messages sent over time.
- **Defence it blinds:** Volume and burst rules.
- **Tell:** Volume gives nothing, so content and sender novelty are the signal.
- **Fallback signals:** Intent classification, sender not in contacts combined with urgent action, number age and line type, link analysis.

`For detection on their side, make it so that now of (number people you wanna pretext)/(60% of the total wave time) , lower the limit if you feel things are getting flagged , rest 40% for follow ups and retry buffers`

#### SMS Throughput and Rate Planning

There is no universal "safe messages per second". The safe rate comes from the route's cap, the message's segment count, a headroom margin, and a gradual ramp-up.

##### 1. Calculate in segments, not messages

Carriers throttle and bill per **segment**. One message can be several segments.

| Encoding | Used when | Single segment | Per segment when split |
|---|---|---|---|
| GSM-7 | Plain English text and common symbols | 160 chars | 153 chars |
| Unicode (UCS-2) | Any non-GSM character: Hindi, Marathi, most emoji, curly quotes | 70 chars | 67 chars |

- A few GSM extension characters (`^ { } \ [ ] ~ | €`) count as 2 characters each.
- One emoji or one Devanagari character switches the **whole** message to Unicode and can triple the segment count.
- Compute `segments per message` from the final rendered text (after personalization), not the template.

##### 2. Factors that decide the sending rate

| Factor | Where it comes from | Effect on output |
|---|---|---|
| Route cap | Provider contract or dashboard (segments/sec) | Hard ceiling. Exceeding it queues or drops messages |
| Headroom (use %) | Our choice, 60-70% recommended | Leaves room for retries, replies and bursts |
| Sending window | Regulation (India promotional: roughly 9am-9pm) | Limits hours per day we can send |
| Segments per message | Message length and encoding | Multiplies total load |
| Retry allowance | Expected failure rate, typically 2-5% | Adds extra volume to the total |
| Carrier daily cap | Some US 10DLC trust tiers | Lowers the per-day capacity |
| Ramp-up | Our choice | Controls how fast a new or quiet sender scales |

##### 3. Formulas

```text
segments_per_message = from section 1 (minimum 1)
total_messages       = recipients x (1 + retry%)
total_segments       = total_messages x segments_per_message

safe_rate            = route_cap x use%                 (segments/sec)
daily_capacity       = safe_rate x 3600 x window_hours  (segments/day)
                       (use carrier daily cap instead if it is lower)

pure_send_time       = total_segments / safe_rate       (seconds)
fastest_finish_days  = ceil(total_segments / daily_capacity)

batch_interval       = (batch_size x segments_per_message) / safe_rate   (seconds)

ramp_day_N_messages  = min(day1_volume x growth^(N-1), daily_capacity / segments_per_message)
```

##### 4. Typical route caps (illustrative only)

Ballpark figures. **Always confirm the real limit with the provider or aggregator.**

| Route type | Typical cap |
|---|---|
| Unregistered long number | ~1 segment/sec |
| Toll-free | ~3 segments/sec |
| Registered 10DLC, mid trust | ~10 segments/sec |
| Aggregator / SMPP bind | ~30 segments/sec (set per bind) |
| Short code or large DLT route | ~100 segments/sec |

##### 5. Worked example

**Inputs**

| Input | Value |
|---|---|
| Recipients | 50,000 |
| Message | 99 characters, all GSM-7, so 1 segment |
| Route cap | 30 segments/sec |
| Use % | 60% |
| Window | 12 hours/day |
| Retry allowance | 3% |
| Batch size | 500 messages |

**Calculation**

```text
total_messages      = 50,000 x 1.03            = 51,500
total_segments      = 51,500 x 1               = 51,500
safe_rate           = 30 x 0.60                = 18 segments/sec (1,080 messages/min)
daily_capacity      = 18 x 3600 x 12           = 777,600 segments/day
pure_send_time      = 51,500 / 18              = ~2,861 sec (about 48 min)
fastest_finish_days = 1
batch_interval      = (500 x 1) / 18           = ~28 sec between batches
```

The route could clear this in under an hour, but we deliberately stretch it over several days (next section).

**Same list as a Unicode message (3 segments)**

`total_segments` becomes 154,500, pure send time rises to about 2.4 hours, and the cost triples. Check the encoding before launch.

##### 6. Ramp-up schedule

Start with the most engaged contacts, then grow volume by a fixed multiplier. Example: day 1 = 2,000 messages, growth = 1.5x per day.

| Day | Messages | Cumulative | Done |
|---|---|---|---|
| 1 | 2,000 | 2,000 | 4% |
| 2 | 3,000 | 5,000 | 10% |
| 3 | 4,500 | 9,500 | 18% |
| 4 | 6,750 | 16,250 | 32% |
| 5 | 10,125 | 26,375 | 51% |
| 6 | 15,188 | 41,563 | 81% |
| 7 | 9,938 | 51,500 | 100% |

**Guidance**

- Keep growth between 1.3x and 1.6x per day. Above 2x looks like a spike to carrier filters.
- A new, quiet or recently re-registered sender should ramp over days or weeks, not hours.
- Do not move to the next day's volume until the previous day passes the pause gates below.

##### 7. Pause gates

Check after every batch and every day. Example using the day-1 volume of 2,000 messages:

| Metric | Pause if | At 2,000 messages |
|---|---|---|
| Opt-out rate | above 1% | more than 20 opt-outs |
| Failure rate | above 5% | more than 100 failures |
| Delivery rate | below 90% | fewer than 1,800 delivered |
| Carrier filtering errors | rising or spiking | investigate before sending more |

When a gate trips, stop, find the cause (content, list quality, template mismatch, route issue), fix it, and resume at a lower volume. Never rotate numbers, sender IDs or URLs to get around filtering.

### 2.4 Obfuscated text
- **What it is:** Spacing, Unicode swaps and character substitution to dodge keyword matching.
- **Defence it blinds:** Keyword filters and simple classifiers.
- **Tell:** Unusual character mix, odd spacing, mixed scripts.
- **Fallback signals:** Normalize before classification, character-level and embedding models, URL analysis.

`I have included a bunch of these in the guide before this;so I won't be repeating this again. Visit the write-up before this (phishing write-up 5) for guidance.` 

### 2.5 Reply-to-unlock
- **What it is:** The message asks the recipient to reply (for example "reply Y") so the link becomes clickable or the sender looks trusted.
- **Defence it blinds:** Link scanning, because the link arrives only after the reply.
- **Tell:** Unknown sender requests a reply, with a vague or urgent pretext.
- **Fallback signals:** Treat reply solicitation from unknown senders as a risk feature, watch the follow-up message for links, warn before replying.

`Here is how to integrate a reply sms infra structure into an already existing sms infrastructure.` 

#### Reply-Code Pipeline: Setup and Integration

Message 1 asks the recipient to reply with their seminar code. When a valid reply arrives, message 2 sends the discount link. This section adds the inbound half to the existing outbound sender.

##### 1. Flow

```text
Campaign sender --> Message 1 (no link) --> Recipient
                                               |
                                          replies with code
                                               v
Provider --> Inbound webhook --> Queue --> Reply worker
                                               |
                       STOP? -> suppress   valid code? -> Message 2 (link)
                                           invalid?    -> one retry message
                                               |
                          Existing outbound queue + rate limiter --> Provider
```

##### 2. Setup steps

1. **Get a replyable number.** Alphanumeric sender IDs are one-way. Use a long code or short code with inbound SMS enabled, and confirm the DLT setup with the provider.
2. **Register templates on DLT.** Register message 1, message 2, the retry message and your branded link domain. Do not use public URL shorteners.
3. **Add the inbound webhook.** Point the provider's inbound URL at an HTTPS endpoint in our service. Verify the provider's request signature (see their docs).
4. **Add the data tables** (section 3).
5. **Add the reply worker** (section 4) as a consumer of a new inbound queue.
6. **Route replies through the existing outbound queue** so they share the same rate limiter, opt-out check and logging.
7. **Add metrics and alerts** (section 6), then test before launch (section 7).

##### 3. Data model

```sql
-- one row per recipient per campaign
contacts(
  id, phone, campaign_id,
  state,            -- sent_msg1 | code_received | link_sent | invalid_retry_sent | closed | opted_out
  msg1_sent_at, reply_at, link_sent_at,
  retry_count       -- max 1
)

seminar_codes(code, discount, expires_at, active)

redemptions(phone, code, created_at)   -- unique on (phone, code)

inbound_log(provider_msg_id UNIQUE, phone, body, received_at)   -- idempotency + audit
```

Keep the suppression list as a separate table that the outbound sender checks before every send.

##### 4. Reply worker logic

```text
on inbound message (phone, body, provider_msg_id):
    if provider_msg_id already in inbound_log: return      # duplicate delivery
    insert into inbound_log

    text = uppercase(trim(body))

    if text in [STOP, STOPALL, UNSUBSCRIBE, CANCEL, END, QUIT]:
        add phone to suppression list; state = opted_out
        return                                             # STOP always wins

    contact = find contact by phone with state = sent_msg1 or invalid_retry_sent
    if none: return                                        # unknown sender, ignore

    code = lookup seminar_codes where code = text and active and not expired
    if code and no redemption for (phone, code):
        record redemption; state = link_sent
        enqueue message 2 (link) on outbound queue
    elif contact.retry_count == 0:
        retry_count = 1; state = invalid_retry_sent
        enqueue retry message
    else:
        state = closed                                     # stop after one retry
```

**Rules to keep**

- **Normalize input.** Ignore case and extra spaces, and accept the code inside a short sentence if you want to be lenient.
- **Respond to the webhook fast.** Return 200 immediately and do the work from the queue, so provider retries do not create duplicates.
- **Time limit.** Ignore replies older than the campaign's reply window, for example 48 hours.
- **Link domain.** Use your own branded domain, optionally with a unique token per contact for tracking.

##### 5. Capacity impact

Replies come back in waves after each batch. Message 2 volume is roughly `sends x reply rate x valid-code rate`. Message 2 uses the same route cap as message 1, so the 30-40% headroom in the throughput plan covers it. If replies spike, the reply worker should send from the same limiter and never bypass it.

##### 6. Metrics and alerts

| Metric | Why it matters | Alert when |
|---|---|---|
| Reply rate | Confirms the audience is real and engaged | Far below expectation after first batch |
| Valid-code rate | Spots a bad list or wrong code set | Below 50% of replies |
| Reply to link latency | Shows queue health | Over 1 minute |
| Webhook errors | Replies are being lost | Any sustained failures |
| Opt-out rate | Same pause gate as the campaign | Above 1% |
| Message 2 delivery rate | Link filtering by carriers | Below 90% |

##### 7. Test before launch

- [ ] Send message 1 to internal numbers on both iPhone and Android, reply with a valid code, and confirm message 2 arrives with a working link
- [ ] Reply with a wrong code, then a second wrong code, and confirm exactly one retry message is sent
- [ ] Reply STOP and confirm the number is suppressed and receives nothing further
- [ ] Replay the same webhook payload and confirm no duplicate message 2
- [ ] Reply from a number that never received message 1 and confirm it is ignored
- [ ] Confirm replies and message 2 share the rate limiter

##### 8. Rollout

Launch with the day-1 batch from the ramp-up plan, check the pause gates and the reply metrics above, then scale volume as planned.

---

## 3. Calls

### 3.1 Caller ID spoofing, including local-number spoofing
- **What it is:** The displayed number is faked, often to match the victim's area code or a trusted organization.
- **Defence it blinds:** Caller ID display and number reputation.
- **Tell:** Missing or partial attestation; number does not match the route it arrived on.
- **Fallback signals:** Attestation (STIR/SHAKEN where available), route-consistency checks, "unverified" labeling, callback verification.

`This is very hard to replicate in a legal manner, as it was compared to in the sms case, but I will still provide you with the way to make this happen legally in testing.`

#### Carrier-Sanctioned Spoof Testing for Authorized Voice Engagements

##### Why this is the closest-to-real option
- STIR/SHAKEN attestation, route-consistency checks, and "unverified" labeling are carrier-side controls
- A client's authorization only binds the client — it cannot authorize use of the public telephone network on its own
- The only way to legally exercise these specific controls live is to get the **network operator itself** to consent to the test
- Everything else (local-number proxy, internal PBX) is a substitute that doesn't actually touch these controls

##### 1. Confirm this is worth pursuing before you start
- This path depends entirely on the client's leverage with their carrier — a small client with a standard consumer/business line is unlikely to get this approved
- Best candidates: clients with an **enterprise telecom contract**, a **dedicated account manager**, or an existing **SIP trunk / PRI relationship** with the carrier
- If the client has none of this, say so early and default to the documentary + proxy-testing approach instead of promising a live test you likely can't get approved

##### 2. Scope it in writing before approaching anyone
- Add explicit language to the RoE/SOW covering: "client will request carrier-sanctioned spoof testing of inbound caller ID controls as part of this engagement"
- Name who owns the carrier relationship on the client side (usually IT/telecom procurement, not security) — you'll need their cooperation, not just security's
- Define the exact scope: which numbers may be spoofed-from, which numbers/extensions are the test targets, and the time window
- Get a named signatory who can authorize both the security test **and** the telecom-facing request

##### 3. The client makes the ask, not the testing firm
- The request has to come from the client's own account, since the carrier has a billing/contractual relationship with the client, not with you
- Typically routed through the client's **enterprise account manager** or **fraud/security liaison** at the carrier, not general support
- Frame the ask as: "we are conducting an authorized security assessment and want to verify our inbound caller ID and STIR/SHAKEN controls — can you authorize and assist with a controlled test call displaying a spoofed/internal number to our own lines"

##### 4. What to have ready when the carrier asks questions
- Signed authorization letter referencing the specific test
- Exact source and destination numbers involved
- Requested test window (date/time range)
- A single technical point of contact on your side the carrier can coordinate with
- Confirmation that no numbers outside the client's own organization will be used as either source or destination

##### 5. What carriers commonly offer instead of a true open spoof test
- Some carriers won't let an external party spoof freely even with client sign-off, but will instead:
  - Run the test **themselves**, placing the call with a controlled spoofed CLI **on their end** and sharing the result with the client
  - Provide a **lab/sandbox environment** that isn't on the live public network, simulating attestation outcomes
  - Share **existing STIR/SHAKEN attestation logs and policy documentation** in lieu of a live test, which is functionally the "documentary" fallback but now carrier-confirmed rather than carrier-interviewed
- Treat any of these as a legitimate result — the goal is validated information, not literally dialing the spoofed call yourself

##### 6. If the carrier declines or doesn't respond
- This is common — don't block the engagement on it
- Fall back to the documentary + proxy-testing approach already discussed (carrier interview, local-number proxy, internal PBX)
- Note in the report that a live carrier-sanctioned test was requested and declined/unavailable, which is itself useful information for the client (it tells them how responsive their carrier is to this kind of request)

##### 7. Documentation to keep regardless of outcome
- Copy of the request sent to the carrier and any response
- Names and roles of everyone who approved the request on the client side
- If the carrier runs the test themselves: their methodology notes and result, even if brief
- Timestamp everything — this is your record if anyone downstream questions the traffic

### 3.2 Number rotation and burner or VoIP lines
- **What it is:** Calls come from short-lived numbers or low-cost VoIP lines.
- **Defence it blinds:** Per-number blocklists.
- **Tell:** Young numbers, VoIP line type, shared prefix block, same behavior across numbers.
- **Fallback signals:** Reputation at prefix, carrier and behavior level; line-type checks; call-pattern analytics.

`Okay so in India you can have international VoIP without KYC but lesten who is picking up  a call from a international number anways ... the firms uses it own number registered under its very name so .. phone calls in general is a hard attacks to pick anyways,but i will still drop the method here .. you can still start`

#### Getting a VoIP Number in India (KYC)

##### 1. Pick a cloud telephony provider
- Exotel — [exotel.com](https://exotel.com)
- Knowlarity — [knowlarity.com](https://knowlarity.com)
- TeleCMI — [telecmi.com](https://www.telecmi.com)
- MyOperator — [myoperator.com](https://myoperator.com)
- Ozonetel — [ozonetel.com](https://www.ozonetel.com)

##### 2. Register a business entity
- A sole proprietorship with a current account is usually enough
- Get GST or Udyam registration if you don't already have one
- Check each provider's minimum entity requirement before signing up

##### 3. Sign up and verify
- Go to the provider's signup page (e.g. **Exotel:** [my.exotel.com/auth/register](https://my.exotel.com/auth/register))
- Use a **business email**, not Gmail/Yahoo — it speeds up approval
- Verify your email and phone number (OTP or missed-call verification)
- You'll land in a **trial account** with a shared/test number until KYC clears

##### 4. Submit KYC documents
- **PAN** card
- **Aadhaar**-linked mobile for OTP verification
- **GST** number or company registration (CIN/LLPIN) for the business
- Signing authority details for the account
- Upload in the dashboard under **Settings → KYC / Company Details**
- Approval typically takes **1–3 business days**

##### 5. Buy and configure a number
- Once approved, go to **Numbers → Buy Number**
- Choose local, mobile-style, or toll-free
- Confirm **outbound calling** is included — some virtual numbers are inbound-only
- Attach it to a softphone, dialer app, or your calling software

##### 6. Register for DLT if sending SMS
- Separate TRAI requirement, done through your provider's DLT portal
- Needed for Entity ID, Sender ID, and message templates

#### Minimal-KYC International VoIP

##### 1. Use a mainstream global provider
- Twilio — [twilio.com](https://www.twilio.com)
- Vonage — [vonage.com](https://www.vonage.com)
- RingCentral — [ringcentral.com](https://www.ringcentral.com)

##### 2. Sign up with a standard account
- Business or personal email, phone number, and payment method
- Verify your personal phone number during signup (standard for Twilio trial accounts)

##### 3. Create a regulatory bundle (for most non-US/Canada numbers)
- In **Twilio Console → Phone Numbers → Regulatory Compliance → Bundles**
- Country-specific; requirements vary — some countries need only an address, others need ID
- **Toll-free and mobile numbers typically have the fewest requirements**
- US/Canada toll-free numbers use a lighter "toll-free verification" instead of a full bundle

##### 4. Buy the number
- **Phone Numbers → Buy a Number**
- Filter by country, local/mobile/toll-free, and capabilities (voice, SMS)
- Confirm calling/SMS permissions for whichever country you'll actually be contacting

##### 5. Know the limits
- These numbers are still linked to your account identity and payment method — "minimal KYC" means lighter paperwork, not anonymity
- Calls from a foreign number into India show as international on the recipient's end, which lowers pickup rates
- India-to-India VoIP termination is restricted by TRAI regardless of which provider issued the number

### 3.3 Low-volume targeted vishing
- **What it is:** A few carefully researched calls impersonating a bank, employer or authority.
- **Defence it blinds:** Volume analytics.
- **Tell:** Content only, such as requests for OTPs, remote-access installation or urgent payment.
- **Fallback signals:** Real-time intent flags, in-call warnings, friction on sensitive actions, number age and contact-graph signals.

`Phone calls only mainly succeed when they are mainly spear-phished so do your digging and target top individuals`

### 3.4 Voice cloning 
- **What it is:** A synthetic copy of a familiar voice is used to impersonate a person.
- **Defence it blinds:** Human trust in the voice and synthetic-voice detectors as a standalone control.
- **Tell:** Often nothing reliable in the audio itself.
- **Fallback signals:** Out-of-band verification, shared code words, callback to a known number, approval workflows for sensitive requests. Use the voice detector only as one weak signal.

`Here are a few of the best tools for voice cloning in the market.`

* **[W-okada](https://github.com/w-okada/voice-changer)** (allow custom) , **[Voice.ai](https://voice.ai/)**  (allow custom) , and **[voicemod](https://www.voicemod.net/)** (no custom voice tho) --> voice clones agents for real time 
* **[Elven lab](https://elevenlabs.io/)** , **[Fish audio](https://fish.audio/app/voice-cloning)** , free option *require decent gpu* **[chatter box](https://github.com/resemble-ai/chatterbox)** - to try [without dowloading](https://huggingface.co/spaces/resemble-ai/chatterbox) --> voice cloning agents form text 

### 3.5 One-ring (Wangiri) campaigns
- **What it is:** Very short calls to many numbers to provoke a callback to a premium-rate number.
- **Defence it blinds:** Rules built for conversational scam calls.
- **Tell:** Sub-second calls to many numbers, sequential dialing.
- **Fallback signals:** Pattern analytics, premium-rate and destination checks on the callback number, warning on callback.

`Here I will provide the infra implementation of the auto-calling and hang-up system.`

#### Call Placing Worker (Queue to Hangup)

##### **Overview**

A worker pulls **phone numbers from a queue**, places a short outbound call to each one, and hangs up **before it is answered**. It reports a result for every job. This section covers only the call lifecycle: **dequeue, originate, ring, hang up, report**.

```
Queue -> Worker -> Rate limiter -> Originate (SIP trunk) -> Ring -> Hangup -> Result event
```

##### **Scope**

- **In scope:** queue consumption, concurrency control, call origination, hangup logic, retries, cleanup, metrics.
- **Out of scope:** caller ID selection, verification codes, session logic, fallback channels.

##### **Job Contract**

```json
{
  "jobId": "01J9Z...",
  "to": "+919876543210",
  "enqueuedAt": 1760000000000,
  "expiresAt": 1760000090000,
  "attempt": 0
}
```

- `to` must already be **E.164**. Validate again in the worker.
- `jobId` is the **idempotency key**. It must never produce two calls.
- `expiresAt` lets the worker **drop stale jobs** instead of ringing a phone for a session that no longer exists.

##### **Components**

- **Queue:** Redis (BullMQ or Streams), SQS or RabbitMQ.
- **Worker:** stateless process, scaled horizontally.
- **Rate limiter:** shared (Redis) token bucket for **calls per second** plus a counter for **max concurrent calls**.
- **Call engine:** Asterisk (ARI), FreeSWITCH (ESL), or a provider API.
- **Result channel:** event or webhook to the service that enqueued the job.

##### **Worker Flow**

1. Pull a job from the queue.
2. If `now > expiresAt`, report **EXPIRED** and acknowledge.
3. Set `SET call:{jobId} 1 NX EX 60`. If it already exists, skip (**duplicate**).
4. Wait for the rate limiter (CPS and concurrency slot).
5. Originate the call.
6. Run the **hangup controller** (below) until the call is settled.
7. Release the concurrency slot, emit the result, acknowledge the job.

##### **Concurrency and Rate Limiting**

- Set **max CPS** to what your trunk allows. Exceeding it causes `503` or congestion failures.
- Set **max concurrent calls** to your trunk channel limit minus headroom.
- Use a **shared limiter** so many workers together stay under the limit.
- Release the slot in a `finally` block so crashes do not leak slots. Give slots a **TTL** as a backstop.
- Keep queue priority simple: oldest job first, since older jobs are closer to expiry.

##### **Originate and Hangup Controller (Asterisk ARI)**

Hang up **based on call state**, not on a fixed delay from dialing. Setup time varies by carrier, so the clock starts when the far end begins ringing.

```js
const RING_HOLD_MS  = 1500;   // how long to let it ring
const HARD_LIMIT_S  = 20;     // absolute max, also set on Asterisk side

function placeCall(job) {
  return new Promise((resolve) => {
    const channel = ari.Channel();
    let settled = false, holdTimer = null;

    const finish = async (status, detail) => {
      if (settled) return;
      settled = true;
      clearTimeout(holdTimer); clearTimeout(hardTimer);
      await channel.hangup().catch(() => {});   // idempotent, ignore "not found"
      resolve({ status, detail });
    };

    const hardTimer = setTimeout(() => finish('TIMEOUT'), HARD_LIMIT_S * 1000);

    channel.on('ChannelStateChange', (_, ch) => {
      if (ch.state === 'Ringing' && !holdTimer) {
        holdTimer = setTimeout(() => finish('RUNG'), RING_HOLD_MS);
      }
      if (ch.state === 'Up') finish('ANSWERED');       // picked up early: drop now
    });

    channel.on('ChannelDestroyed', (ev) => finish('FAILED', ev.cause_txt));

    channel.originate({
      endpoint: `PJSIP/${job.to}@carrier-trunk`,
      app: 'flashcall',
      appArgs: job.jobId,
      timeout: HARD_LIMIT_S                              // Asterisk-side safety net
    }).catch((err) => finish('FAILED', err.message));
  });
}
```

For **FreeSWITCH**, use `bgapi originate` with `originate_timeout`, watch `CHANNEL_PROGRESS` (ringing), wait the hold time, then `uuid_kill`.

##### **Hangup Rules**

- **Start the hold timer on Ringing** (SIP 180 or 183), not on dial. Test this on each carrier.
- **Hang up with CANCEL** while unanswered. If the call is already answered, it becomes a BYE automatically.
- **If the user answers, hang up immediately** and mark `ANSWERED`. Still count it as delivered, since the caller ID was displayed. Flag it in metrics.
- **Never hang up before ringing starts.** Early hangup can stop the call from showing as missed.
- **Keep one settle path** (`finish`) guarded by a `settled` flag, so timers, events and errors cannot double-fire.
- **Always set a hard limit** in the worker **and** in the engine (`timeout`), so a crashed worker cannot leave a phone ringing.

##### **Result Codes**

| Status | Meaning | Retry |
|---|---|---|
| **RUNG** | Ringing reached, hung up cleanly | No |
| **ANSWERED** | User picked up, hung up immediately | No |
| **BUSY** | Cause 17 | No |
| **NO_ROUTE** | Unallocated or invalid number (cause 1, 28) | No |
| **UNREACHABLE** | Switched off or no coverage (cause 18, 20) | No |
| **CONGESTION** | Trunk or network congestion (cause 34, 38, 42, SIP 503) | Yes, once |
| **TIMEOUT** | Hard limit reached with no ringing | Yes, once |
| **EXPIRED** | Job older than `expiresAt` | No |

Exact cause behavior differs by carrier, so log the **raw SIP code and Q.850 cause** with every result.

##### **Retry and Dead Letter Rules**

- Retry **only** when the call never reached ringing: `CONGESTION`, `TIMEOUT`, trunk error.
- **Never retry** after `RUNG` or `ANSWERED`. The user already saw a missed call, and a second one causes confusion.
- Retry at most **once**, after 500-1000 ms.
- Drop the retry if it would run past `expiresAt`.
- Send jobs that exhaust retries to a **dead letter queue** with the last error.

##### **Reliability and Cleanup**

- Set the queue **visibility timeout** above `HARD_LIMIT_S` plus headroom, so a job is not redelivered mid-call.
- Run an **orphan reaper** every 10 seconds and on startup. It lists active channels for the app, then hangs up any older than the hard limit.
- On `SIGTERM`, **stop consuming**, let in-flight calls settle, and force hangup after about 5 seconds.
- Make the **result emit** idempotent, and retry it with backoff if the receiver is down.
- Reconnect the ARI websocket automatically. On reconnect, run the reaper before consuming again.

##### **Result Event**

```json
{
  "jobId": "01J9Z...",
  "status": "RUNG",
  "sipCode": 180,
  "q850Cause": null,
  "ringMs": 1512,
  "totalMs": 3240,
  "attempt": 0,
  "trunk": "carrier-trunk",
  "finishedAt": 1760000003240
}
```

##### **Monitoring**

- **Ring rate:** share of calls reaching `RUNG`, split by carrier and trunk.
- **Setup latency:** time from dequeue to ringing (p50, p95, p99).
- **Queue age:** oldest job in the queue, and the share dropped as `EXPIRED`.
- **Concurrency and CPS usage** against limits.
- **Answered rate**, since a rising value means users are not following the instruction.
- **Orphan count** from the reaper, which should stay near zero.
- **Alerts:** failure rate spike, queue age above your session expiry, trunk down.

##### **Testing Checklist**

- Test ring hold values of **1, 1.5 and 2 seconds** and confirm a missed call appears on Jio, Airtel, Vi and BSNL phones.
- Answer a call quickly and confirm instant hangup and the `ANSWERED` status.
- Call invalid, switched-off and busy numbers and check the status mapping.
- Kill a worker mid-call and confirm Asterisk times out and the reaper cleans up.
- Enqueue the same `jobId` twice and confirm only one call is placed.
- Push above the CPS limit and confirm the limiter queues instead of failing.
- Run a load test at peak volume, such as a signup campaign.

### 3.6 Cross-channel handoff
- **What it is:** A call is followed by an SMS link, an email, or a request to install a remote-access tool.
- **Defence it blinds:** Any detector that looks at one channel only.
- **Tell:** A call followed within minutes by a link or an installation prompt for the same recipient.
- **Fallback signals:** Entity graph linking numbers, domains, file hashes and identities across channels; propagate a block from one channel to the others.

`I am providing some emails here where it makes perfect sense to follow up with a second channel.`

#### Appointments & Scheduling

##### Booking or rescheduling calls
- **SMS:** confirmation with date, time, and location
- **Email:** calendar invite and prep instructions

##### Reminder or confirmation calls
- **SMS:** short reminder with a one-tap confirm, reschedule, or cancel link
- **Email:** full details and cancellation policy

##### Service visits and technician calls
- **SMS:** arrival window, technician name, live ETA
- **Email:** job summary and what to prepare

#### Customer Support

##### Troubleshooting calls
- **Email:** ticket number, steps discussed, help-article links
- **SMS:** short status update or case ID

##### Complaint or escalation calls
- **Email:** written acknowledgment, resolution timeline, point of contact
- **SMS:** heads-up when the issue is resolved

##### Order, delivery, or shipment calls
- **SMS:** tracking link and delivery updates
- **Email:** order summary and invoice

##### Returns and refund calls
- **Email:** return label, instructions, refund timeline
- **SMS:** confirmation when the refund is processed

#### Sales & Outreach

##### Missed calls and voicemails
- **SMS:** brief "sorry I missed you" with the reason for calling and a booking link
- **Email:** slightly longer intro if the lead is business-focused

##### Discovery or demo calls
- **Email:** recap, next steps, relevant resources
- **SMS:** meeting link for the follow-up call

##### Quote or pricing calls
- **Email:** written quote or proposal with terms
- **SMS:** nudge to review the quote

#### Payments & Collections

##### Payment reminder calls
- **SMS:** payment link and due date
- **Email:** invoice or statement

##### Payment arrangement calls
- **Email:** written summary of the agreed plan and dates
- **SMS:** reminder before each installment

##### Renewal or subscription calls
- **Email:** renewal terms and pricing
- **SMS:** link to renew or update payment method

#### Onboarding & Accounts

##### Welcome or onboarding calls
- **Email:** getting-started guide, login instructions, key contacts
- **SMS:** link to the first step

##### Account change calls
- **Email:** confirmation of what changed and when
- **SMS:** alert if the change wasn't requested by the customer

#### Healthcare & Wellness

##### Appointment calls
- **SMS:** date, time, clinic address, and a reminder
- **Email:** forms to fill out beforehand and visit prep

##### Results or prescription-ready calls
- **SMS:** generic "your update is ready" message only
- **Email:** link to a secure portal, with no sensitive details in the body

#### Recruiting & HR

##### Screening or interview scheduling calls
- **SMS:** interview time and meeting link
- **Email:** job description, agenda, and interviewer names

##### Offer or onboarding calls
- **Email:** offer letter, documents to sign, start date
- **SMS:** reminder to complete pending paperwork

#### Business & Project Calls

##### Agreement or scope calls
- **Email:** "confirming what we discussed" recap as a written record
- **SMS:** only to flag that the email was sent

##### Multi-party or decision-maker calls
- **Email:** forwardable summary for colleagues who weren't on the call
- **SMS:** quick nudge if a response is needed

#### Logistics & Coordination

##### Address or location calls
- **SMS:** address and map pin
- **Email:** directions, parking, and access instructions

##### Access codes and reference numbers
- **SMS:** the code or number, so it isn't dictated and misheard
- **Email:** backup copy with context

#### Rules of Thumb

##### Which channel for what
- **SMS:** short, urgent, link-based, time-sensitive
- **Email:** long, detailed, attachment-heavy, needs a record

##### Cautions
- **Never** send sensitive data (health, financial, ID numbers) over SMS
- **Get consent** before marketing-style follow-ups (TCPA, TRAI/DND, DLT)
- **Avoid over-messaging:** one well-timed follow-up beats three at once

---

## 4. Combined-Signal Switching Table

| Tell detected | Detector degraded | Shift weight to | Action |
|---|---|---|---|
| Auth passes, new domain, first-time sender | Auth, reputation | Domain age and infrastructure, intent classifier, landing-page analysis | Warning banner, click-time scan |
| Scanner and victim divergence, or challenge gate | URL scanning | Infrastructure age, brand-lure mismatch, sender novelty, kit fingerprint | Delay and quarantine, never deliver as clean |
| Delivery-time content differs from click-time content | Link-rewrite/delivery-time scanner | Click-time rescanning, multi-profile fetch comparison | Rescan on every click, not only at delivery |
| First hop on a trusted/whitelisted platform, final hop unrelated | Domain/URL reputation (first-hop scoring) | Full-chain resolution, final-destination scoring, link creation recency | Flag chain crossing from reputable host to unrelated domain, hold for chain resolution |
| Image or QR heavy, little text | Text classifier | OCR and QR output, visual brand similarity, sender behavior | Re-run text pipeline on extracted content, then warn; elevate risk by default for QR from unknown senders |
| Link on a trusted platform | Domain reputation | Page-level render analysis, sender-platform history, recipient context | Interstitial with real destination |
| High wording variety, one intent cluster | Signatures, keywords | Embeddings, infrastructure clustering, recipient overlap | Block at campaign level |
| Pretext split across multiple messages in a thread | Per-message content classifier | Thread/sequence-level aggregation, post-engagement urgency-shift detection | Score the sequence, not single messages |
| Known sender, behavior shift | Sender reputation, auth | Behavioral anomaly, payment-change rules | Hold and verify out-of-band |
| Sandbox-evasion behavior | Sandbox verdict | Static features, file prevalence, sender signals | Treat as suspicious, delay, re-detonate later |
| Password-protected archive or uncommon container | Static AV/sandbox | Body-text password extraction, container-type policy, file prevalence | Unpack with extracted password, enforce depth limits, detonate |
| Throwaway domain per SMS | Blocklist | Domain age, generated-name features, infrastructure clustering | Block by cluster, warn user |
| Low-volume targeted SMS or call | Volume rules | Intent, number age and line type, contact graph, cross-channel links | Warn, add friction on OTP and payment actions |
| Spoofed-CLI indicators | Caller reputation | Attestation, route consistency, callback verification | Label "unverified", prompt callback |
| Probing pattern | ML thresholds | Randomized thresholds, shadow models, ensembles | Uniform responses, rate limit, tarpit |
| Same domain or number across email, SMS and calls | Single-channel scores | Entity graph | Propagate block across channels |

---

## 5. Design Rules for Safe Switching

1. **Unknown is not safe.** Missing coverage must raise risk or force a fallback, never default to benign.
2. **Re-weight smoothly.** Use hysteresis and expiry instead of hard on/off modes, and run new rules in shadow mode first.
3. **Assume attackers can game the switch.** A faked tell can push a legitimate sender into quarantine or flood the analyst queue. Require corroboration from a second signal family before blocking, and cap how far one tell can move a decision.
4. **Prefer expensive-to-fake signals.** Infrastructure age, sending history and behavioral baselines cost attackers more to change than wording does.
5. **Maintain an evasion registry.** Each entry holds technique, tell, fallback, owner and last-tested date.
6. **Replay evasions continuously.** Run an automated suite weekly that tests each technique against the pipeline, with and without adaptive switching, and track detection and false-positive rates over time.
7. **Treat "well-known" as decaying.** Any technique specific enough to be publicly documented (a named container format, a named cloaking method) has a shelf life — vendors patch for it once it's well known. Re-test the registry's coverage, don't assume a technique stays blocked once it's handled.
8. **Report honestly.** Detection is probabilistic. Agree on precision, recall, false-positive rate and time-to-block with the client, and avoid promising complete coverage.
