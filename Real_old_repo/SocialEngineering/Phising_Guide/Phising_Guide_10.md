`This still requires more work and expansion, so that's still coming up, I guess.` 

## Part 2: People, Identity and Response

### 5. Physical

#### 5.1 Badge and access control (tailgating)

**Recommend**
- Policy: a clear challenge policy. Any member of staff may ask an unfamiliar person in a restricted area to show a badge. Frame it as expected behavior, not confrontation.
- Signage at badge-controlled doors reminding staff not to hold doors for unverified followers. A simple "please do not tailgate" sign is low cost and sets a visible norm.
- Policy: a no-blame escalation path. Staff who are unsure whether to challenge someone get a named person or extension that handles it for them.

**Implement**
- Identify the doors where tailgating was observed during the engagement and categorize them: badge-controlled, staff-held or unmonitored.
- For restricted-floor doors (finance floor, server room, executive area), verify the badge reader requires a swipe from each entrant, not only the first person through.
- Post the signage at each restricted door, including the challenge and tailgating policy.
- Include a tailgating scenario in the annual fire drill or physical security walkthrough so staff practice challenging in a no-stakes setting.
- ★ Publish the escalation extension on the signage and on the intranet, and test it once a quarter with a scripted call.

#### 5.2 Reception and delivery procedures

**Recommend**
- Policy: a one-page reception verification procedure for couriers and contractors, requiring a delivery slip or work order, a call to the department being visited (never to a number the visitor gives), and a signed visitor log.
- Add "a department name is not sufficient to verify a delivery" to reception training. A caller who only confirms a department exists (public information) is not enough to admit a visitor.

**Implement**
- Put the procedure at each reception desk as a laminated one-page checklist:
  - Ask for a delivery slip or work order.
  - Call the department using the internal directory, not a number provided by the visitor.
  - Log the visitor's name, purpose, escort and time in and out.
  - Issue a clearly distinguishable visitor badge.
- Test with a scripted visit once per quarter and report adherence to the Head of Facilities.

#### 5.3 OSINT exposure reduction

**Recommend**
- Policy: review job postings before publication. Product names (email gateway vendor, MFA provider, cloud platform) hand attackers pretext detail. Describe requirements by capability where possible.
- Policy: an out-of-office standard without travel dates, delegate names or internal escalation paths. "I am unavailable until [date], please contact [team email]" is enough.
- Brief communications and marketing on badge-design exposure. Event photos that clearly show the badge format give an attacker a reference for a convincing lookalike.

**Implement**
- ★ Run a quarterly OSINT sweep for the client: search current job postings for named technology products, review staff-facing and public photos for visible badge designs, and search for out-of-office or travel details published by executives and finance staff.
- ★ Give HR a pre-publication checklist for job postings (no vendor or product names, capability wording instead) and give marketing a photo-review checklist (badges, screens, whiteboards, door-access panels).
- ★ Publish the out-of-office template in the email signature and auto-reply guidance.

---

### 6. Sign-in and Authentication

#### 6.1 Phishing-resistant sign-in

**Recommend**
- Deploy security keys or passkeys first for privileged accounts, finance roles and executives. These are the accounts a real attacker targets after a successful phish, and the only sign-in method that cannot be relayed through a look-alike page.

**Implement**
- Identify the accounts to deploy keys to first: privileged accounts, finance roles and executives.
- Verify the identity platform supports FIDO2 hardware keys and that they cannot be bypassed by falling back to SMS or push.
- Run a pilot of 10 to 20 accounts before full rollout. Using test accounts, confirm the relay test fails on all of them, then expand.
- ★ Enforce the requirement with a policy, such as an authentication-strength or conditional access rule that allows only phishing-resistant methods for those roles, rather than leaving keys as an optional extra.

#### 6.2 Number matching and sign-in context

**Recommend**
- Enable number matching and sign-in context (location, device, application) for all accounts still using authenticator push. A prompt showing the number to enter and the location of the attempt makes blind approval meaningfully harder.

**Implement**
- Enable number matching on all push-based MFA approvals where the platform supports it.
- ★ Updated test: run the relay test against a test account with number matching on, and record the result. Expect that it can still succeed when the relay page displays the correct number to the user. Number matching narrows the attack window but does not close it. Only a security key closes the gap fully.

#### 6.3 Retiring SMS one-time codes

**Recommend**
- Policy: commit to a specific timeline for replacing SMS codes for privileged roles. SMS is the weakest deployed MFA method and can be completed through a relay.

**Implement**
- ★ Pull the report of registered authentication methods per account, and list privileged and finance accounts that still have SMS registered.
- ★ Turn off SMS as a permitted method for those roles once their security keys are in place, and track the remaining count monthly until it reaches zero.

#### 6.4 Risk-based sign-in alerting

**Recommend**
- Enable risk-based sign-in policies. Alerts fire for new country or new ASN sign-ins, repeated MFA prompt sequences (prompt bombing), new inbox rules after sign-in, new app consent grants, and new mail forwarding rules.
- Give the SOC a documented runbook for each alert type: what it means, what the analyst checks, and when to escalate.

**Implement**
- Enable these alert conditions in the identity platform:
  - Sign-in from a new country or previously unseen ASN.
  - Multiple MFA prompts to one account within a short window (prompt bombing).
  - New inbox rule created within 60 minutes of sign-in.
  - New OAuth application consent within 60 minutes of sign-in.
  - New mail forwarding rule created.
- Validate each alert by simulating the condition from a test account.
- Write and publish a runbook per alert type so SOC analysts know the expected triage steps.

---

### 7. Reporting and Containment

#### 7.1 Reporting infrastructure

**Recommend**
- Deploy the report button on every mail client (desktop, web and mobile). A staff member on a phone should be able to report in one tap, not by forwarding to a mailbox they have to remember.
- Policy: a signed, no-blame reporting policy. Staff know that reporting a message that later proves legitimate carries no penalty, and that reporting is valued.
- Close the loop: every person who reports a confirmed phishing message receives an acknowledgment. Even a short "thank you, we are investigating" builds the behavior.

**Implement**
- Confirm the report button is deployed on desktop, web and mobile mail clients.
- Test the button from each client type (desktop Outlook, OWA, iOS Mail, Android Gmail) and confirm it routes to the correct triage mailbox or SIEM feed.
- Set up an automated acknowledgment to the reporter within 5 minutes of the report arriving.
- ★ Send the signed no-blame policy from the CISO to all staff, and link it from the report button's confirmation message.

#### 7.2 Removal speed

**Recommend**
- Automate removal. The manual process took 34 to 52 minutes, and roughly 70 percent of clicks occurred in that window. Automated removal triggered by a confirmed report verdict should bring it below 10 minutes.
- Set a triage-to-removal target (recommended 30 minutes or less), measure it on every incident, and report it to the CISO monthly.

**Implement**
- Implement automated removal using the platform's eDiscovery or message-recall API, triggered by a confirmed verdict from the triage queue.
- Test the full flow end to end: report, triage, confirmed malicious, automated removal from all mailboxes.
- Set the 30-minute triage-to-removal SLA and instrument it so it is measurable per incident.
- ★ Include link-verdict flips from section 1.5 as a second trigger for the same removal automation, so a delayed-weaponization event is handled the same way as a user report.

#### 7.3 SOC runbook for phishing incidents

**Recommend**
- Give the SOC a single documented phishing runbook, owned by a named person and reviewed at least annually.

**Implement**
- Create and publish the runbook covering:
  - How a phishing report enters the queue.
  - Who triages it and in what timeframe.
  - The criteria for confirming malicious versus benign.
  - The automated and manual removal steps.
  - Deconfliction: who to call if removal triggers a wider incident.
  - Post-incident: what to tell reporters, when to send an all-staff notice, and how to document the event.
- Run it as a tabletop exercise with the SOC team before the next engagement.

---

### 8. Help Desk and Verification

#### 8.1 Written procedure enforcement

**Recommend**
- Redesign the procedure so the process itself enforces the key steps, not individual agent judgment: the tool requires a ticket with a verified callback number and a completed verification question before a reset option becomes available.
- Add manager approval as a system-enforced step for MFA device changes, not a manual step an agent can skip under time pressure.
- Send an automatic notification to the user's other registered channels (email, phone) whenever a password or MFA device changes. A real user who did not request it sees the notice within minutes.
- Policy: add procedure adherence as a measured service-level metric in the outsourcing contract, reported to the client monthly.

**Implement**
- Build the verification procedure as a step-by-step checklist embedded in the help desk tool, not a separate PDF.
- Enforce these steps in the tool, with no agent bypass:
  - Employee ID entered.
  - Callback to the number on file completed and logged.
  - Verification question asked and answer logged (the question must not be answerable from public sources).
  - For MFA changes: manager approval ticket number required.
  - Post-reset notification sent to the user's other registered channel.
- Run a 30-day shadow-monitoring period after rollout. Measure adherence per call and feed it back to the vendor weekly.
- ★ Amend the vendor contract with the adherence SLA, a reporting format, and a remedy if adherence falls below target.

#### 8.2 Call handle time

**Recommend**
- Review whether handle-time targets or queue pressure create an incentive to skip the longer verification flow. Calls that followed the full procedure were more than twice as long as calls that did not. If handle time is rewarded, verification will be bypassed.

**Implement**
- ★ Export handle time per call type and compare it with adherence per call. Check whether agents with the shortest handle times also have the lowest adherence.
- ★ Change quality and bonus scoring so that a skipped verification step is a critical failure and handle-time targets exclude calls that require verification.

---

### 9. Continuous Testing and Measurement

#### 9.1 Metrics to track between engagements
- **Phishing report rate**: target is rising over time.
- **Report-before-click rate**: maintain above 80 percent.
- **Time from first report to removal**: under 30 minutes.
- **Help desk adherence rate**: 100 percent of full procedure steps.
- **MFA method distribution**: proportion of accounts on security keys versus push versus SMS.
- **Domain monitoring alert rate**: count of new lookalike domains registered each month.

#### 9.2 Recurring tests to run internally
- Monthly `dnstwist` sweep with a new-domain alert to the security team (section 1.3).
- Quarterly test of the help desk procedure using test accounts and the same script used in this engagement (section 8.1).
- Quarterly drill of the SOC triage-to-removal runbook using a controlled test message (sections 7.2 and 7.3).
- Annual tabletop covering a cross-channel attack (phish, MFA relay, help desk call) to test how the three control areas work together (sections 4, 6 and 8).
- ★ Quarterly gateway regression: repeat the section 1 validation tests (age rules, homoglyphs, protected names, click-time rescanning, QR decoding, archive handling), because gateway updates and vendor changes can silently change behavior.

#### 9.3 What to test in the next engagement
- Business email compromise and payment-change simulation (not in scope for this engagement, recommended in the report).
- Regional office physical access (landlord consent was not received; obtain it for the next cycle).
- SMS smishing wave once the DLT sender ID is registered under the client's PE.
- USB or media drop targeting the finance floor.
- Spear-phishing of executives using OSINT-sourced context.
