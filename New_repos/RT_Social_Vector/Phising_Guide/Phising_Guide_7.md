# LEGAL METRICS OF PHISHING ENGAGEMENT 

`THE ART OF NOT GETTING IN TROUBLE`

> **Read this first.** This guide is general professional guidance, not legal advice. Laws change and differ by country, state, and even city. Before you sign anything or touch any target, have a qualified lawyer in the client's jurisdiction (and yours, if different) review your contract, your Rules of Engagement, and your plan. Every "check the current rule" note in this guide means: verify it again before the engagement.

---

## How to use this guide

This guide has six sections. They follow the order of a real engagement:

| # | Section | Question it answers |
|---|---------|---------------------|
| 1 | Authorization and legal protection | "Are we allowed to do this, and can we prove it?" |
| 2 | Scope and Rules of Engagement | "What exactly may we do, to whom, and when do we stop?" |
| 3 | Planning | "What is the goal, how do we measure it, and how do we prepare safely?" |
| 4 | Ethics and employee welfare | "How do we make sure nobody gets hurt or fired because of us?" |
| 5 | Reporting | "How do we turn what happened into clear, fair, useful findings?" |
| 6 | Defensive measures | "How do we help the client actually get safer afterwards?" |
| A | Appendices | Templates, checklists, glossary, references |

Each section is written as "do this, then this." Where something is a trap that has hurt real testers, it is marked **TRAP**. Where something is a shortcut that saves time, it is marked **TIP**.

---

## The 12 golden rules (print this page)

1. **No signature, no testing.** Verbal approval and friendly emails are not authorization.
2. **The person who signs must actually own the thing being tested.** A friendly IT manager does not own the building, the cloud account, or other people's jobs.
3. **Every person, place, and system you touch must be in writing in the scope.** If it is not written, it is out of scope.
4. **Third parties are not covered by your client's signature.** Landlords, outsourced call centers, cloud and email providers, and SMS carriers may need their own permission.
5. **Carry proof and a 24/7 phone number that is answered.** Test the number before the engagement, not during an arrest.
6. **Never improvise outside the approved pretexts.** If the conversation goes somewhere that is not on the approved card, stop or exit gently.
7. **Never capture real passwords.** Prove that a submission happened without keeping what was typed.
8. **Test the process, not the person.** Findings blame missing controls, never named individuals.
9. **Nobody gets punished for our test.** Put this in the contract before you start.
10. **Collect the minimum personal data, protect it, and delete it on schedule.**
11. **Stop immediately if real harm appears:** a real emergency, a real attacker, a distressed employee, or a legal threat.
12. **Write everything down as it happens.** A good log protects you, the client, and the employees.

---

## Engagement map (where each section fits)

```
WEEK -4 to -2        WEEK -2 to 0         TEST WINDOW          WEEK +1 to +2        WEEK +2 onward
-----------------    -----------------    -----------------    -----------------    -----------------
Section 1            Section 2            Execute              Section 5            Section 6
Authorization        Scope + RoE          (Section 3 plan      Report + debrief     Defenses +
Contract             Final sign-off       + Section 4 rules)   Retest plan          Retest
Third parties        Infra build          Live logging
Laws                 Pretext approval     Stop conditions
                     Section 3 + 4 prep
```

---

# SECTION 1: AUTHORIZATION AND LEGAL PROTECTION

**Goal of this section:** by the end, you hold a signed contract from the right person, a letter you can show if challenged, written permission from every third party who matters, and a clear picture of which laws apply.

**The one-sentence rule:** *Authorization is the only thing that separates a security tester from a criminal, so it must be signed, specific, provable, and given by someone with the power to give it.*

---

## 1.1 Who must sign (and who must NOT)

### Step 1: Ask the "Four Ownership Questions"

Before you ask anyone to sign, find out who owns each of these. Write the answers in a table. They are often four different parties.

| Question | Why it matters | Example of a trap |
|----------|----------------|-------------------|
| 1. Who is the legal entity that employs the people we will target? | Only the employer can approve testing its staff | Staff are actually employed by a staffing agency or a sister company |
| 2. Who owns or leases the building(s) we may enter? | Only the owner or controller of the premises can authorize entry | The client rents one floor; the landlord controls the lobby, garage, and elevators |
| 3. Who owns and runs the systems (email, cloud, VPN, phones)? | Systems run by vendors need the vendor's rules respected | Email is a SaaS service with its own testing policy |
| 4. Who runs the processes we will test (help desk, payment desk, reception)? | Outsourced processes belong to the outsourcer | The help desk is a call center in another country run by a vendor |

### Step 2: Choose the signer

**Acceptable signers (any one, as long as they have the authority):**

- The CEO or Managing Director
- The CFO
- The General Counsel / Head of Legal
- The CISO or CIO **only if** you have written proof they have delegated authority to approve security testing and to bind the company to contracts
- Any officer named in a board resolution or delegation-of-authority letter for this purpose

**Not acceptable as the only signer:**

| Person | Why not |
|--------|---------|
| IT manager, security analyst, help desk lead | They can request testing, but usually cannot authorize testing of other departments, buildings, or people |
| The friendly person who called you | Having contacted you does not make them authorized |
| Procurement | They can sign the commercial purchase, but not permission to attack the company |
| A department head | Can only speak for their own department |
| A parent or sister company contact | They may not control the legal entity being tested |
| Someone who only emailed you | Email from an unverified sender is easy to fake (the irony is real) |

### Step 3: Verify the signer is real and has authority

Do these five checks. Write down the date, the result, and who did it.

1. **Registry check.** Look up the company in the official company registry for its country (for example Companies House in the UK, the state Secretary of State in the US, the Ministry of Corporate Affairs in India). Confirm the legal name, the registered address, and the listed directors or officers.
2. **Authority proof.** Ask for one of: a board resolution, a delegation-of-authority letter, or a signed statement from the General Counsel confirming the signer's authority.
3. **Call-back.** Find the company's main phone number on its official website or registry listing (not from an email signature). Call it, ask for the signer by name, and confirm they sent the paperwork.
4. **Domain check.** Confirm the signer's email domain matches the company's real domain. Look for lookalike domains.
5. **Signature trail.** Use an e-signature tool with an audit trail, or collect a wet-ink signature and scan it. Keep the original.

> **TIP:** The call-back check takes five minutes and has saved testers from fake "clients" who were really attackers trying to get a free reconnaissance report.

### Step 4: Add a second signer if the Four Ownership Questions show other owners

- Building owned by someone else: the landlord or property manager signs a short consent (see 1.5).
- Staff employed by another entity: that entity's authorized signer approves.
- Outsourced help desk: the vendor's authorized signer approves, or the client confirms in writing that its vendor contract allows testing and names who at the vendor approved.

---

## 1.2 What the signed contract must include

You will usually have **three layers** of paper. Do not blend them into one messy document.

| Layer | Document | Purpose |
|-------|----------|---------|
| Commercial | Master Services Agreement (MSA) and Statement of Work (SOW) | Price, dates, liability, confidentiality, deliverables |
| Operational | Rules of Engagement (RoE), built in Section 2 | Exactly what you may and may not do |
| Proof | Authorization letter (1.4), third-party consents (1.5) | What you show people if challenged |

### The contract checklist

Tick every row. If a row is missing, ask for it before starting.

| # | Clause | What it must say | Why it matters |
|---|--------|------------------|----------------|
| 1 | Parties | Full legal names, registered addresses of client and tester | Wrong entity = no authority |
| 2 | Authorization statement | The client expressly authorizes named testers to attempt the listed social engineering activities against the listed targets during the listed dates | This sentence is your legal shield |
| 3 | Authority representation | The signer states they have authority and that the client owns or controls the targets, or has listed which targets it does not control | Moves responsibility to the client if they lied |
| 4 | Testing window | Start and end dates and times, with time zone | Anything outside the window is unauthorized |
| 5 | Scope reference | Attached scope and RoE, with version number | Prevents arguments later |
| 6 | Named testers | Names (and ID type) of every person authorized | Unknown people are not covered |
| 7 | Subcontractors | Whether allowed; if yes, they must be named and bound by the same terms | A subcontractor outside the contract is outside the shield |
| 8 | Third-party handling | List of third parties and who is responsible for getting each consent | Avoids "I thought you were doing that" |
| 9 | Prohibited actions | Short list (see Section 2) | Sets hard lines |
| 10 | Emergency contacts | Two 24/7 people on each side, with phone numbers | Needed during an incident or detention |
| 11 | Stop-work right | Either party can pause the test at any time | Safety valve |
| 12 | Real incident clause | If you find a real attacker or real compromise, you notify the sponsor immediately and pause related activity | Avoids contaminating a real incident |
| 13 | Critical finding clause | Serious issues are reported within a set time (for example 24 hours), not at the end | Prevents harm during the engagement |
| 14 | No-discipline clause | The client agrees not to discipline, penalize, or identify individuals based on test results | Protects employees (see Section 4) |
| 15 | Confidentiality | Both sides protect the findings and the existence of the test | Prevents leaks that spoil the test |
| 16 | Data protection terms | Roles (controller/processor), what personal data you will handle, where it is stored, retention period, deletion proof, breach notification | Required by GDPR-style laws |
| 17 | Evidence handling | Which evidence you keep, how long, how it is destroyed | Matches the data protection terms |
| 18 | Liability and indemnity | Caps, and who covers what; client indemnifies you for claims caused by their lack of authority | Protects you from a client's false assurances |
| 19 | Insurance | You carry professional indemnity / cyber and general liability; amounts stated | Clients and courts expect it |
| 20 | Change control | Scope changes only by written, signed amendment | Prevents "can you also quickly try..." |
| 21 | Deliverables | Report format, delivery date, debrief meeting, retest terms | Avoids disputes |
| 22 | Governing law and disputes | Which country's law, which courts or arbitration | Matters if something goes wrong |
| 23 | Termination | How either party can end it and what happens to data | Clean exit |
| 24 | Physical safety | Tester will not enter unsafe areas; client will tell you which areas are hazardous | Protects testers |
| 25 | Law enforcement clause | Client agrees to vouch for testers to police and to withdraw any complaint | Helps if police are called |

### Sample clauses (give these to your lawyer to adapt)

> **Authorization.** The Client authorizes the Tester's named personnel listed in Schedule 2 to perform the social engineering activities described in Schedule 1 (the "Authorized Activities") against the targets described in Schedule 1, between [start date/time, time zone] and [end date/time, time zone]. The Client confirms that it has the authority to grant this authorization and that it owns, leases, or otherwise controls each target listed, except for those third-party targets identified in Schedule 3.

> **No-discipline.** The Client agrees that results of the Authorized Activities will be used only to improve systems, processes, and training. The Client will not use individual results to discipline, demote, dismiss, or publicly identify any employee or contractor, and will not request individual-level results from the Tester.

> **Real incident.** If the Tester believes it has discovered evidence of a real compromise or real criminal activity unrelated to the test, the Tester will stop the related activity and notify the Client's Emergency Contact immediately.

> **Law enforcement.** If any Tester is questioned or detained while carrying out Authorized Activities, the Client will promptly confirm the authorization to the questioning party and will not support any complaint arising from the Authorized Activities.

---

## 1.3 Insurance: the quiet safety net

1. Ask your insurer (or broker) in writing: "Does our policy cover social engineering and physical penetration testing, including on-site work and impersonation?" Some policies exclude physical testing or require a notice.
2. Confirm your policy covers the countries where the client operates.
3. Ask the **client** whether their cyber insurance requires notice of the test, or could be voided by "unauthorized" security activity. Ask them to check.
4. Keep certificates of insurance ready to attach to the contract.

> **TRAP:** Many testers discover an exclusion only after a claim. A five-minute email now avoids that.

---

## 1.4 The authorization letter ("get out of jail free" letter)

### What it is

A one-page signed letter that proves you are allowed to be where you are and do what you are doing. It is not magic. It only works if the signer really has authority over the place (see the Iowa lesson below).

### Step-by-step: build it

1. **Put it on the client's letterhead** (or include the client's logo and registered address).
2. **Say who you are.** Full names of each tester and the tester company.
3. **Say what is authorized.** One or two sentences: "...authorized to conduct a physical and social engineering security assessment at the following locations..."
4. **List the locations** by full address, including floors or areas if limited.
5. **List the dates and time window**, with time zone.
6. **List two emergency contacts:** names, titles, mobile numbers, and a note that they are available 24/7 during the window. Neither contact may be a target of the test.
7. **Add a verification phrase.** A codeword the contact will ask for or you will say, so that police or guards can confirm this is not a forged letter. Share the codeword only with the contacts and the testers.
8. **Add the signer's name, title, signature, and date.**
9. **Add the property owner's counter-signature** if the premises are leased or shared.
10. **Add a line:** "Please contact the above persons before taking further action."

### Step-by-step: prepare the contacts (do this a week before)

1. Call each emergency contact from your own phone. Make sure the number works.
2. Ask them to save your number and put it in their favorites.
3. Tell them exactly what they will be asked: "Someone may call to confirm the test. Ask for the codeword. If they give it, confirm the test is authorized, then ask to speak to the police officer or guard."
4. Ask what time zone and sleeping hours apply. If the test runs at night, confirm they will answer at 3 a.m.
5. Do a **dry run**: call them at a normal time with a fake "I am being questioned" call and see how they handle it.
6. Make sure the contacts know who the building owner and local police liaison are, if any.

### Step-by-step: what each tester carries

| Carry | Why |
|-------|-----|
| Printed, signed original of the letter | Paper works when phones die |
| A copy on the phone (offline PDF) | Backup |
| Photo ID | Police will ask |
| Both emergency contact numbers written on paper | Phones can be taken or dead |
| The RoE summary (one page) | Shows exactly what is allowed |
| Charged phone and power bank | Basic |
| Some cash | Taxi or phone call |
| The tester company's contact for your own lawyer or manager | You may need support |

**Do NOT carry:** anything that could be treated as a weapon; anything illegal in that place. Be especially careful with lock-picking tools and bypass tools. In many places, possessing such tools can be treated as evidence of intent to commit burglary, and the law differs by country and state. Only carry them if the RoE explicitly allows it **and** your lawyer confirms it is legal where you are working.

### Step-by-step: what to do if you are challenged

1. **Stay calm and stay polite.** Raised voices help nobody.
2. **Stay in your pretext only as long as it is safe and plausible.** When the pretext fails, switch.
3. **Use the safe line:** "I am conducting an authorized security assessment for [Client]. I have a signed letter, and here is a number you can call to confirm."
4. **Hand over the letter** and the contact numbers. Do not argue about whether you are allowed.
5. **If the guard or police arrive:** keep your hands visible, do not run, do not hide things, follow lawful instructions, and ask them to call the emergency contact.
6. **If you are detained:** do not resist. Say the safe line again. Ask to call your emergency contact and your own manager or lawyer. Do not volunteer extra information beyond the facts of the engagement.
7. **Write down the time, names, and what was said** as soon as you are able.
8. **Tell the sponsor right away** that you were challenged, even if everything was resolved.
9. **Pause that activity** until the sponsor and you agree how to continue.

### Real-world lesson: the Iowa courthouse case (2019)

Two authorized testers were arrested at a county courthouse in Iowa during a physical test. They had authorization from the state judicial branch, but the county said it owned the building and had not agreed. The charges were eventually dropped, but the testers spent a night in jail, and the incident became a famous warning. **The lesson:** the person who signs must control the premises, or the owner of the premises must also sign.

---

## 1.5 Third parties: who, when, how, and what to say

### What counts as a third party?

A third party is anyone who is not your client but could be affected, or whose property, service, or staff you may touch. Your client's signature **does not** automatically cover them.

### The master table

| Third party | When you need them | Who to contact (through the client) | What to ask for | Typical lead time* |
|-------------|--------------------|-------------------------------------|-----------------|--------------------|
| Landlord / property manager | Any physical test in leased or shared buildings | Property or facilities manager | Written consent for the dates, areas, and activities; confirm their security team's procedure | 2 to 3 weeks |
| Building security company (contracted guards) | Physical tests where guards may challenge you | Client's facilities lead; the guard company's account manager | Notify that a test may occur (without details) and agree escalation steps; decide whether guards will be told | 1 to 2 weeks |
| Co-tenants or shared spaces | Shared lobbies, parking, shared offices | Landlord | Landlord confirms the shared areas are included | 2 weeks |
| Email provider (Microsoft 365, Google Workspace, others) | Phishing tests against the client's tenant | Client's IT admin | Read the provider's current security-testing policy; set up required allow-listing or advanced delivery configuration | 1 week |
| Email security gateway vendor | If you want to deliver past filters | Client's email admin | Allow-listing steps, or an explicit decision NOT to allow-list | 3 to 7 days |
| Cloud hosting provider (for your own infrastructure) | When hosting landing pages or senders | Your account manager or provider's policy page | Check their acceptable use policy for simulated phishing | 1 week |
| Domain registrar / DNS provider | Registering lookalike domains | Your registrar's support or abuse team | Confirm the registration will not be flagged or suspended mid-test | 3 to 5 days |
| SMS or telephony provider | Smishing and vishing | The provider's compliance team | Registration of sender IDs and templates where required, and terms for simulation | 2 to 4 weeks (can be longer in some countries) |
| Outsourced help desk / call center | Vishing or help desk testing | Client's vendor manager | The client confirms in writing that its contract allows testing, and ideally the vendor's authorized signer approves | 2 to 4 weeks |
| MSSP / outsourced SOC | They will see alerts from your test | Client's security lead | Decide: tell them (deconflict) or not; document either choice | 1 to 2 weeks |
| Payroll, HR, or other SaaS vendors | If your pretext refers to them or if staff data may come up | Client's vendor manager | Usually just notify; avoid impersonating them without their approval | 1 week |
| Works council / employee representatives | In countries where they must be consulted (for example Germany, the Netherlands, France and others) | Client's HR or legal | Their consultation or consent as local law requires | 4 to 8 weeks |
| Trade union | Where staff are unionized and agreements cover monitoring | Client's HR | Notice or agreement as applicable | 2 to 6 weeks |
| Regulator | Regulated entities (finance, health, critical infrastructure) where frameworks apply | Client's compliance team | Check whether the test falls under a framework (for example TIBER-style programs) | Varies |
| Local police (non-emergency liaison) | Large physical tests; sensitive sites | Client security lead | Optional pre-notification for high-risk sites | 1 week |
| Cyber insurer | If policy conditions require | Client's risk manager | Confirm no breach of policy | 1 week |
| Social media platforms | Pretexting with fake profiles | N/A (usually forbidden) | Most platforms ban fake accounts, so avoid it unless legal and the RoE explicitly allow it | N/A |

\*Lead times are typical planning guesses. Ask each party how long they take.

### Step-by-step: the third-party process

1. **Ask the client for a list.** Give them a one-page questionnaire: "Who owns the building? Who runs your help desk? Who runs your email? Who provides your guards? Who provides your SMS or phone service?"
2. **Classify each item:** client-owned, client-leased, or third-party operated.
3. **Decide who contacts whom.** In almost every case the **client** should make the first contact, because they have the relationship. You supply the technical and scope details.
4. **Give the client a ready-to-send message** (template below).
5. **Collect written consent** or a written "no consent needed" confirmation from the third party. Keep both.
6. **Decide on deconfliction.** For each third party that may see activity (guards, MSSP, call center), decide whether to tell them. If you tell them, tell as few people as possible.
7. **Record everything** in a Third-Party Register (see Appendix B).
8. **If a third party refuses**, remove that target from scope, or redesign the test. Do not "just try."

### Template: email from the client to a third party (landlord example)

> **Subject:** Notice and request for consent: planned security assessment at [address]
>
> Dear [Name],
>
> [Client] has commissioned an authorized security assessment of its offices at [address]. The assessment will include testing of physical access controls and staff security procedures, performed by [Tester company] between [dates and times].
>
> To avoid any confusion or alarm, we ask for your written consent to include the following areas: [list areas, for example: ground-floor lobby, floor 5 reception, parking level B1].
>
> The testers will carry a signed authorization letter and have 24/7 emergency contacts. They will not damage property, force entry, or access restricted or hazardous areas. If building security stops them, please call [Contact name and number] before involving the police.
>
> Please reply to confirm consent or tell us if any areas, dates, or conditions need to change.
>
> Regards,
> [Client signer name and title]

### Template: phone script for the client's contact calling a vendor

> "Hello, this is [Name] from [Client]. We are running an authorized security assessment of our own staff procedures. We are writing to confirm that our testing of [help desk / email / SMS service] is allowed under our contract. Could you point me to the right person to approve this in writing? I will send a short email with the dates and scope."

### Practical notes by third-party type

**Landlords and building security:**
- Ask them to name the person who will confirm the test if a guard calls.
- Ask what the guards are told to do with someone who has no badge. You are testing their real behavior, but you must not put the guard at risk.
- Provide building-specific safety info to your testers (fire exits, restricted areas).

**Cloud and email providers:**
- Policies change. Read the provider's **current** security-testing or penetration-testing page before each engagement. Look specifically for rules about simulated phishing and about impersonating the provider.
- Do not impersonate the provider's own staff or brand without checking their rules.

**Outsourced help desks and call centers:**
- Their staff are not your client's employees. Their local laws and employment rules may differ.
- Provide a way for the vendor to know that a test is happening (through a manager) without telling the front-line agents.

**SMS and telephony:**
- Many countries require sender registration for bulk or commercial messages. Test messages may still get blocked if you do not follow the carrier's rules.
- Use a reputable provider and tell them it is a security simulation.

**Works councils and unions:**
- Where they exist, they often have legal rights over monitoring and employee data. Start early. Bring a one-page plain-language summary, not a technical report.

---

## 1.6 Laws to consider by country

> **Reminder:** this table is a map, not a ruling. Always confirm with local counsel. Laws are amended often.

### The four legal questions that apply almost everywhere

1. **Is access authorized?** (computer misuse and trespass laws)
2. **Is anyone being deceived in a way that could be fraud or personation?** (fraud and impersonation laws; consent and no intent to gain usually matter)
3. **Is personal data being collected or processed?** (privacy and data protection laws)
4. **Are calls, messages, or conversations being recorded or intercepted?** (wiretap, interception, and recording-consent laws)

### Major-country quick reference

| Country / region | Computer misuse and access | Fraud / personation / trespass | Privacy and data protection | Recording and messaging traps |
|------------------|---------------------------|--------------------------------|-----------------------------|-------------------------------|
| **United States** | Computer Fraud and Abuse Act (18 U.S.C. § 1030) and state computer crime laws | Wire fraud (18 U.S.C. § 1343) as a theoretical risk without consent; impersonating a federal officer (18 U.S.C. § 912) and state laws on impersonating police or officials; state trespass and burglary laws; pretexting provisions in the Gramm-Leach-Bliley Act (15 U.S.C. §§ 6821 to 6827) for financial institutions | State privacy laws (for example California's CCPA/CPRA, which covers employee data); HIPAA if the target is a healthcare entity; sector rules | Federal wiretap law plus state recording laws: some states require **all** parties to consent (for example California, Florida, Illinois, Maryland, Massachusetts, Montana, New Hampshire, Pennsylvania, Washington; verify the current list); TCPA rules for texts and automated calls; some states have investigator-licensing rules that can touch security testing (check the state) |
| **United Kingdom** | Computer Misuse Act 1990 | Fraud Act 2006 (false representation); burglary under the Theft Act 1968 can apply to entering as a trespasser with intent; impersonating police is an offense | UK GDPR and Data Protection Act 2018 | Interception rules (RIPA 2000 and lawful business practice regulations) apply to recording; PECR for electronic messaging |
| **European Union** (general) | Directive 2013/40/EU implemented in each member state's criminal code | National fraud and trespass laws | GDPR (lawful basis, minimization, processor agreements, possibly a DPIA); national employee-data rules | ePrivacy rules; national recording laws; works councils may need consultation or consent |
| **Germany** | Criminal Code §§ 202a to 202c (including the "hacker tools" provision) | Fraud and trespass in the Criminal Code | GDPR plus the Federal Data Protection Act; employee data rules | Works council co-determination on monitoring and technical systems; strict recording rules |
| **France** | Criminal Code arts. 323-1 and following | Fraud and personation offenses | GDPR; CNIL guidance | Employee representative body (CSE) consultation may be required |
| **Netherlands** | Computer crime provisions of the Criminal Code | Fraud and trespass | GDPR; Dutch employee-data rules | Works council consent for certain employee-monitoring and data policies |
| **India** | Information Technology Act 2000 (ss. 43, 66, 66C identity theft, 66D cheating by personation using a computer resource, 72 breach of confidentiality) | Bharatiya Nyaya Sanhita 2023 (which replaced the Indian Penal Code) covers cheating, cheating by personation, and criminal trespass | Digital Personal Data Protection Act 2023 and its rules (being phased in; check status) | CERT-In directions require reporting certain real incidents quickly (for example within six hours) and log retention; telecom regulator rules require registered headers and templates for bulk commercial SMS, which affects smishing tests |
| **Canada** | Criminal Code s. 342.1 (unauthorized use of a computer) and s. 430(1.1) (mischief in relation to data) | Identity fraud and personation offenses; provincial trespass acts | PIPEDA (and provincial laws such as Quebec's Law 25) | Interception rules allow one-party consent in federal law; anti-spam law (CASL) governs commercial electronic messages |
| **Australia** | Criminal Code Act 1995 (Part 10.7, computer offences) and state laws | Fraud and trespass in state laws | Privacy Act 1988 and Australian Privacy Principles | Telecommunications (Interception and Access) Act; state surveillance-device laws differ on recording consent; Spam Act 2003 |
| **Singapore** | Computer Misuse Act 1993 | Cheating and personation offenses | Personal Data Protection Act 2012 | Spam Control Act; financial-sector technology risk rules for banks |
| **UAE** | Federal cybercrime law (Federal Decree-Law No. 34 of 2021) | Personation and fraud offenses | Federal Personal Data Protection Law (Decree-Law 45 of 2021); separate regimes in DIFC and ADGM | Strict rules on recording people; local counsel is essential |
| **Saudi Arabia** | Anti-Cyber Crime Law | Personation and fraud offenses | Personal Data Protection Law | Local counsel is essential |
| **Japan** | Act on the Prohibition of Unauthorized Computer Access | Fraud and trespass offenses | Act on the Protection of Personal Information | Local counsel recommended |
| **Brazil** | Law 12,737/2012 (computer intrusion) | Fraud and personation offenses | LGPD | Local counsel recommended |
| **South Africa** | Cybercrimes Act 19 of 2020 | Fraud and trespass offenses | POPIA | Local counsel recommended |

### Practical rules that keep you safe in almost every country

1. **Do not take anything of value.** If you walk off with a laptop, a document, or cash "to prove a point," you may cross from testing into theft. Use a **flag** (see Section 4) instead.
2. **Do not impersonate police, government officials, emergency services, doctors, or regulators.** This is illegal almost everywhere and causes real fear.
3. **Do not record audio or video of people** unless the RoE allows it **and** your lawyer has confirmed the recording rules for every country involved (including where the other person is located). Take written notes instead.
4. **Treat personal data as dangerous.** Collect the minimum, store it encrypted, set a retention period, and delete it on time.
5. **Cross-border work multiplies the rules.** If you are in one country, the client in a second, and the targets in a third, get advice on all three.
6. **Never test an organization without an agreement just to "show them."** Unsolicited testing is unauthorized access, however good your intentions.
7. **If anyone asks you to do something illegal, refuse in writing** and tell your manager.

---

## 1.7 Section 1 go / no-go checklist

Do not start any activity until every box is ticked.

- [ ] Four Ownership Questions answered and recorded
- [ ] Signer verified (registry check, authority proof, call-back, domain, signature trail)
- [ ] MSA and SOW signed
- [ ] Authorization clause, no-discipline clause, real-incident clause, and law-enforcement clause present
- [ ] Insurance checked and certificate available
- [ ] Authorization letter signed (and counter-signed by the property owner where needed)
- [ ] Two 24/7 emergency contacts called, tested, and briefed; dry run done
- [ ] Codeword set and shared with the right people only
- [ ] Third-Party Register complete, with written consents or written "not needed" notes
- [ ] Works council or employee-representative requirements checked
- [ ] Legal review done for each country involved
- [ ] Each tester carries letter, ID, contacts, and RoE summary
- [ ] Each tester knows the "if challenged" steps

---

# SECTION 2: SCOPE AND RULES OF ENGAGEMENT (RoE)

**Goal of this section:** by the end, you have a signed RoE document that says exactly **who** you may target, **how**, **when**, **what you must never do**, **who knows**, and **when everything stops**.

**Think of the RoE as the rulebook of a game.** If both teams do not read the same rulebook, someone gets hurt and everyone argues. We write the rulebook first and play second.

---

## Step 1: Gather the right people (the "RoE team")

You need decisions from several roles. If one is missing, you will discover the gap at the worst moment.

| Role | What they decide | What goes wrong without them |
|------|------------------|------------------------------|
| **Executive sponsor** (CISO, CIO, COO, or similar) | Goals, scope approval, risk appetite | Nobody can approve exceptions |
| **Engagement owner** (your counterpart who handles day-to-day) | Daily questions, change requests | You wait days for answers |
| **Legal / compliance** | Laws, privacy, contract terms, regulator rules | You break a rule nobody told you about |
| **HR** | Employee rules, no-discipline policy, sensitive groups, works councils | A test hits someone on medical leave or in a disciplinary process |
| **Physical security / facilities** | Buildings, guards, safe areas, badge rules | Guards call police without a plan |
| **IT / email admin** | Allow-listing, mail flow, kill-switch actions | Your emails bounce or get blocked |
| **SOC / incident response lead** | Deconfliction: do they know or not | Real incident response is triggered by your test |
| **Help desk / finance lead** (when relevant) | Process rules for those teams | A test accidentally causes a real payment or reset |
| **Communications / PR** | Messaging if the test leaks | A rumor spreads that the company has been attacked |

**Write down** who is on the RoE team, their phone numbers, and which of them are inside the "need-to-know circle" (Step 6).

---

## Step 2: Hold the kickoff meeting (use these questions)

Book 90 minutes. Bring this list. Write the answers into the RoE as you go.

### About the goal
1. Why are you doing this test now? (Audit, incident, regulator, new CEO, routine?)
2. What decision will the results help you make?
3. What would a "good" result look like, and what would a "bad" result look like?

### About the organization
4. How many employees, how many sites, which countries?
5. Which teams handle money, credentials, personal data, or visitors?
6. Are any parts of the business outsourced (help desk, payroll, reception, guards)?
7. Is the company regulated (financial, health, government)? Which rules apply?
8. Are there staff who should not be targeted (see Step 5)?

### About safety and fairness
9. What happened the last time you ran a phishing test? Any complaints?
10. Are employees under stress right now (layoffs, restructuring, merger)?
11. Are there union or works council rules we must follow?

### About logistics
12. What are the blackout dates (payroll, audits, product launches, holidays, board meetings, peak seasons)?
13. Who is awake and reachable at the hours we will test?
14. Which email filters, endpoint tools, and SOC services are in place?

> **TIP:** Record the meeting (with permission) or assign a note-taker. Send the notes back within 24 hours and ask the sponsor to reply "agreed" or correct them. That reply is your first piece of written scope.

---

## Step 3: Choose the attack types (the menu)

For each row, the client chooses **In scope**, **Out of scope**, or **In scope with conditions**. Explain each row in plain language.

| Attack type | What it means in plain words | Typical risk level | Common conditions |
|-------------|------------------------------|--------------------|-------------------|
| **Phishing (email)** | Sending test emails to see whether people click, enter details, or report | Low to medium | Approved pretexts; no real credential storage; allow-listing decision |
| **Spear phishing** | Emails tailored to a team or role using public information | Medium | Role-level information only; no family or personal details |
| **Vishing (phone)** | Phoning staff or the help desk and trying to get information or actions | Medium to high | Approved scripts; no recording unless legal OK; stop rules for distress |
| **Smishing (SMS)** | Sending text messages with test links | Medium | Carrier registration; staff mobile numbers must be company-provided or consented; privacy checks |
| **Physical entry / tailgating** | Trying to walk into a building or area without proper access | High | Authorization letter; landlord consent; guard escalation plan; safe areas only |
| **Pretext visits** | Posing as a delivery person, contractor, or visitor | High | Approved pretexts; no impersonation of police or officials |
| **USB or media drop** | Leaving benign test devices where staff may find them | Medium | Devices contain only a harmless marker, never real malware; client approves locations |
| **QR code lures** | Posters or messages with QR codes to a test page | Low to medium | Landing page is benign; locations approved |
| **Help desk social engineering** | Asking help desk to reset passwords, MFA, or share information | High | Use client-provided test accounts; no real user accounts |
| **Business email compromise simulation** | Fake "executive" payment or invoice requests | High | Payment never actually made; finance lead in the circle; approved wording |
| **Dumpster / paper review** | Looking at discarded paper | Medium | Only the client's own bins on the client's premises; no personal items |
| **Wi-Fi lures** (evil-twin style demonstration) | Rogue network to see if devices connect | High | Often better as a separate technical test; strong conditions |
| **Social media or OSINT pretexting** | Using public information to build believable stories | Medium | Passive only; fake profiles usually forbidden |
| **Insider-assisted tests** | A cooperating employee helps you | Low | Documented; counts as a different type of test |

**How to decide (the baby steps):**

1. Read the row aloud.
2. Ask: "Do you want us to try this?" and "What is the worst thing that could happen if we do?"
3. If the worst case is serious (panic, injury, financial loss, legal trouble), either mark **Out of scope** or add conditions that remove the worst case.
4. Write the decision and the conditions in the RoE table.

---

## Step 4: Write the exclusions (what we will NOT do)

This is as important as the scope. Present this as a default list and let the client add to it.

### Default exclusions

| # | Never do this | Why |
|---|---------------|-----|
| 1 | Impersonate police, courts, tax authorities, or any government | Illegal in many places; causes fear |
| 2 | Impersonate emergency services, doctors, or hospitals | Dangerous and illegal |
| 3 | Use real emergencies (a fire, a threat, a medical emergency) as a pretext | Could trigger real panic or delay real help |
| 4 | Use sensitive personal topics (layoffs, pay cuts, bonuses, disciplinary action, medical issues, family problems, bereavement) | Causes distress and anger |
| 5 | Promise money, gifts, or benefits that will not be delivered | Creates resentment and legal exposure |
| 6 | Target employees' personal email, personal social accounts, or family | Out of the client's authority |
| 7 | Use real, current crises (a pandemic, a natural disaster, a real company incident) | Exploits real fear |
| 8 | Impersonate a specific real person (a named executive) without that person's written approval | Personal rights and reputation |
| 9 | Use or test real malware | Risk of real damage |
| 10 | Force entry, break locks, disable alarms, damage property | Crosses into crime |
| 11 | Enter unsafe areas (labs, clean rooms, operating rooms, ICUs, electrical rooms, construction zones) | Safety |
| 12 | Touch or photograph personal belongings, people's screens, or people's faces | Privacy |
| 13 | Record audio or video of people without RoE approval and legal sign-off | Recording laws |
| 14 | Keep, copy, or use real passwords | See Section 4 |
| 15 | Access real accounts of real employees | Unnecessary risk; use test accounts |
| 16 | Remove property or documents | Theft risk |
| 17 | Test minors or anyone under 18 | Special protections |
| 18 | Continue after a stop signal | Obvious |
| 19 | Contact customers, partners, or vendors of the client unless in scope and approved | Out of authority |
| 20 | Discuss the test with outsiders or on social media | Confidentiality |

### Client-specific exclusions (ask these questions)

- Are there executives or board members who must not be contacted?
- Are there teams under investigation or in a sensitive process?
- Are there buildings or floors that are off limits (for example the CEO's floor, a legal department, a lab)?
- Are there customers or partners who must not be mentioned?
- Are there brands, product names, or logos we must not use?
- Are there times of day when staff are especially vulnerable or busy?

---

## Step 5: Build the target pool (who will be tested)

### 5.1 Decide the pool type

| Pool type | Meaning | When to use |
|-----------|---------|-------------|
| **Whole company** | Every employee with a company email | Measuring overall awareness |
| **Sample** | A random subset (see sample size table in Section 3) | Large companies; cheaper and less disruptive |
| **Department-based** | One or more departments | Targeted risk (finance, HR, help desk) |
| **Role-based** | Roles like executives, new hires, administrators | Testing high-risk roles |
| **Physical area** | Specific doors, floors, buildings | Physical tests |

### 5.2 Build the list (step by step)

1. Ask the client to export a list of in-scope employees from their directory. Prefer **work email addresses and role/department only**. Do not request home addresses, personal phone numbers, or HR data.
2. Ask the client to **remove** people who must be excluded (see 5.3).
3. Assign each person a **random ID** (for example `T-00417`). Keep the key (ID to person) separate, encrypted, and with minimal access.
4. Divide the pool into **groups** big enough to protect privacy (at least 10 people per reportable group; see Section 4).
5. Record the final count, the date of the list, and who approved it.
6. Lock the list. Changes need change control (Step 12).

### 5.3 Who to exclude by default (ask HR)

| Group | Why exclude (or handle specially) |
|-------|-----------------------------------|
| Staff on medical, parental, or bereavement leave | Not appropriate to target; they cannot respond anyway |
| Staff currently in a disciplinary or legal process | Results could be misused |
| Staff known to have relevant health conditions (for example severe anxiety) that HR flags | Welfare |
| Anyone under 18 | Special protections |
| Executives or board, **unless** explicitly included | They need to opt in with a clear sponsor |
| Security, SOC, and incident response staff, **unless** testing them is the goal | They will react correctly and skew results; or are tested separately |
| Staff in countries where the legal review is incomplete | Compliance |
| Contractors and temps, **unless** their employer agrees | Not your client's employees |
| The RoE team and the need-to-know circle | They already know |

### 5.4 Special note on executives and VIPs

Executives are valuable targets because attackers love them. They are also the people most likely to complain loudly. Handle them this way:

1. Ask the sponsor: "Do you want executives included?"
2. If yes, get a written approval from the sponsor and, ideally, from the executives themselves (or the CEO for the whole group).
3. Use the same ethics rules (no personal details, no family, no sensitive topics).
4. Tell the executive group in the final debrief first, privately.

---

## Step 6: Decide who knows (the need-to-know circle)

The test works best when most staff do not know about it. But some people must know.

### The circles

| Circle | Who | What they know | Why |
|--------|-----|----------------|-----|
| **Circle 1: Control team** | Sponsor, engagement owner, legal | Everything | They approve and make decisions |
| **Circle 2: Safety net** | Emergency contacts, physical security lead, HR lead | Dates, scope, escalation steps | They rescue testers and protect staff |
| **Circle 3: Deconfliction** | SOC/IR lead (and sometimes IT admin) | Technical indicators (sender domains, IPs, phone numbers) but not timing or targets | They avoid triggering a real incident response |
| **Circle 4: Everyone else** | Staff | Nothing until the debrief | The test is realistic |

### Rules for the circles

1. Keep Circle 1 and 2 as small as possible, but never zero outside your own team.
2. Record every person in each circle, with their role, and have them sign a short confidentiality acknowledgement.
3. Decide **in advance** whether the SOC will be told. There are two valid designs:
   - **Blind SOC:** the SOC is not told, so you measure real detection. A deconfliction contact (Circle 3) can stop a runaway incident response after a set time (for example if the SOC escalates to external authorities).
   - **Informed SOC:** the SOC is told, so they do not waste effort; used when the goal is to test other things.
4. Never tell a manager "do not tell your team" without a plan for what the manager says if staff ask.

---

## Step 7: Plan the escalation and deconfliction process

This is the plan for the moment when somebody reacts to your test.

### 7.1 Situations to plan for (write a response for each)

| Situation | Who acts | What happens |
|-----------|----------|--------------|
| An employee reports a test email to the security team | SOC | Treated as a real report. After the agreed time, the SOC deconflicts with Circle 3 contact and confirms it is the test |
| The SOC begins a real incident investigation | Circle 3 contact | Contact tells SOC "this is part of an authorized test" after your agreed stage, and shares indicators |
| Guards detain a tester | Emergency contact | Confirms authorization; follow the "if challenged" steps in Section 1 |
| Police are called | Emergency contact and sponsor | Verify testers; the client confirms authorization; the client withdraws complaint |
| An employee becomes upset or distressed | Tester | Stop, disclose gently, tell sponsor (Step 8 stop rules) |
| An employee posts about the test on social media | Sponsor, communications | Decide whether to end the phase early; take down if needed |
| An employee forwards the test to the whole company | Sponsor | Decide whether to end the phase early; use that as a finding (positive reporting) |
| A real attacker is discovered during the test | Tester, sponsor | Pause related activity; hand evidence to incident response |
| A tester realizes they made a mistake (wrong person, wrong building) | Tester | Stop, notify engagement owner, log as a deviation |

### 7.2 The codeword system

Use two words:

- **Stop word:** "Please stop the test now." When the sponsor or emergency contact says it, **all** activity stops within minutes.
- **Verification word:** used by testers to prove to police or guards that the contact is the real contact.

Both words are shared only with named people. Change them for every engagement.

---

## Step 8: Define success criteria and stop conditions

### 8.1 Success criteria

"Success" in a social engineering test does **not** mean "we fooled people." It means "we learned something useful." Define success per goal (more in Section 3).

| Test type | Example of a measurable objective | What "done" looks like |
|-----------|-----------------------------------|------------------------|
| Phishing | Measure click, submit, and report rates in groups | Campaign finished, data collected, sample size met |
| Help desk | Test whether a password reset can be obtained with weak verification | One attempt per scenario, outcomes recorded |
| Physical | Test whether an unauthorized person can reach a specific flag in a restricted area | Flag photographed or attempt blocked and recorded |
| Vishing | Test whether staff disclose a category of information | Calls completed per script, disclosures recorded by category |

### 8.2 Stop conditions (write these exactly)

Stop the activity **immediately** and call the engagement owner if any of these happen:

| # | Stop condition |
|---|----------------|
| 1 | A real emergency is happening (fire, medical event, threat) |
| 2 | Someone is distressed, crying, panicking, or says they feel unsafe |
| 3 | A person mentions a medical issue, a personal crisis, or self-harm |
| 4 | A person thinks a **real** emergency is happening because of your test |
| 5 | A tester is physically challenged or detained (follow the letter steps) |
| 6 | Police or emergency services are called and they are not already aware |
| 7 | A real attacker or real compromise is discovered |
| 8 | You discover that a person or area is out of scope |
| 9 | A tester makes a mistake that touches a non-target |
| 10 | Real credentials or sensitive data were captured by accident |
| 11 | The sponsor, engagement owner, or emergency contact says "stop" |
| 12 | Any legal or regulatory concern is raised |
| 13 | A target says they will complain officially or call a lawyer |
| 14 | A tester feels unsafe |
| 15 | The test is leaked to employees and continuing would be misleading |

### 8.3 How to end gently (script)

When a stop condition happens during a live interaction:

> "I'm really sorry about this. I'm actually part of an authorized security exercise run by [Company], and nothing here is real. You haven't done anything wrong and nobody will be told your name. I'll let the security team know right now. If you have any questions, please contact [name] at [number]."

Then:

1. Thank them.
2. Leave if physical, or end the call/message.
3. Log the incident.
4. Call the engagement owner.

---

## Step 9: Build the schedule and the testing windows

### 9.1 Choose windows

| Test type | Window guidance |
|-----------|-----------------|
| Phishing | Weekdays, business hours in each target time zone; spread over several days to avoid one big spike |
| Vishing | Normal working hours; avoid lunch and end of day |
| Physical | Business hours for entry tests; after-hours only if explicitly in scope and safe |
| Help desk | Normal staffing hours |

### 9.2 Blackout dates (ask and write down)

- Payroll days and month-end closes
- Audit periods
- Product launches or earnings announcements
- Holidays and local festivals
- Known restructurings or layoffs
- Major incidents (do not test while the company is already in crisis)

### 9.3 Build a one-page schedule

| Phase | Dates | Owner | Dependencies |
|-------|-------|-------|--------------|
| Contract signed | | Both | Section 1 complete |
| Third-party consents | | Client | Section 1 table |
| OSINT (passive) | | Tester | Approved sources |
| Infrastructure build and test | | Tester | Approved domains |
| Pretext approval | | Sponsor | Written pretext cards |
| Dry run (internal only) | | Tester | Infrastructure ready |
| Go / no-go meeting | | Sponsor | All checklists complete |
| Live window | | Tester | Go decision |
| Debrief with sponsor | | Both | Draft report |
| Employee debrief message | | Client | Sponsor approval |
| Report delivery | | Tester | Evidence organized |
| Retest plan | | Both | Fixes underway |

---

## Step 10: Data handling and evidence rules

Agree these in writing before collecting anything.

| Question | Decision to record |
|----------|-------------------|
| What personal data will we collect? | Minimum: work email, role, department, test response events |
| Where will it be stored? | Encrypted, access-controlled, named location and country |
| Who can access it? | Named testers only |
| How long will we keep it? | Fixed period (for example 30 to 90 days after the report) |
| How will it be deleted? | Secure deletion with written confirmation |
| What evidence will we collect? | Screenshots of landing pages, logs of events, photos of flags (no faces, no screens) |
| What evidence will we NOT collect? | Passwords, audio of people, photos of personal belongings, contents of documents |
| How will evidence be labeled? | Pseudonymous IDs, UTC timestamps |
| What happens if we capture something by mistake? | Stop, isolate, notify sponsor, delete as agreed, record as a deviation |

---

## Step 11: Set up communications during the test

1. **Daily check-in:** a 10-minute call or message with the engagement owner (what was done, any incidents, tomorrow's plan).
2. **Immediate alerts:** any stop condition triggers a phone call, not an email.
3. **Shared log:** a running activity log (see Section 5) that the engagement owner can read in near real time (without personal data).
4. **Secure channel:** agree a secure messaging or email channel for sensitive details. Do not discuss active tests on public tools.
5. **Out-of-hours contact:** agree who answers at night and on weekends.

---

## Step 12: Set up change control

Scope will want to change. That is fine, but it needs a process.

1. Anyone can request a change, in writing, to the engagement owner.
2. The request says: what changes, why, what risk it adds.
3. The sponsor approves or declines in writing.
4. Legal reviews if the change adds people, places, or channels.
5. The RoE version number goes up (v1.0 to v1.1) and everyone gets the new version.
6. **Nothing changes until the signed version arrives.**

> **TRAP:** "Can you also quickly try the CFO?" on a phone call. Say: "Yes, please send it in writing and we'll add it."

---

## Step 13: Write the RoE document

Use this template. Fill in every section. Do not leave "TBD" in a signed version.

```markdown
# RULES OF ENGAGEMENT (RoE)

**Client:** [legal name]
**Tester:** [legal name]
**Engagement title:** [name]
**RoE version:** v1.0
**Date:** [date]

## 1. Purpose and goals
[Plain-language goals; see Section 3 of this guide]

## 2. Authorization reference
Contract: [MSA/SOW reference and date]
Authorization letter: [reference and date]
Third-party consents: [list and attach]

## 3. Testing window
Start: [date, time, time zone]
End: [date, time, time zone]
Blackout dates: [list]

## 4. Authorized activities (in scope)
| Activity | Conditions | Approved pretext IDs |
|---|---|---|
| | | |

## 5. Excluded activities (out of scope)
[List from Step 4, plus client-specific exclusions]

## 6. Target pool
Type: [whole company / sample / department / role / physical area]
Size: [number]
List date and approver: [details]
Excluded groups: [list]

## 7. Locations (physical tests)
| Address | Areas allowed | Areas forbidden | Owner and consent |
|---|---|---|---|
| | | | |

## 8. Personnel
Tester team: [names, roles]
Subcontractors: [names or "none"]

## 9. Need-to-know circles
Circle 1: [names, roles]
Circle 2: [names, roles]
Circle 3: [names, roles]

## 10. Emergency contacts (24/7)
Client primary: [name, phone]
Client secondary: [name, phone]
Tester primary: [name, phone]
Tester secondary: [name, phone]
Stop word and verification word: [shared separately]

## 11. Escalation and deconfliction
[Table from Step 7]

## 12. Success criteria
[Table from Step 8]

## 13. Stop conditions
[List from Step 8]

## 14. Data handling
[Table from Step 10]

## 15. Reporting
Draft report due: [date]
Final report due: [date]
Debrief meeting: [date]
Aggregation rules: groups of at least [10] people; no individual-level results

## 16. No-discipline commitment
[Client commits not to discipline individuals based on results]

## 17. Change control
[Process from Step 12]

## 18. Signatures
Client authorized signer: [name, title, signature, date]
Tester authorized signer: [name, title, signature, date]
```

---

## Step 14: Final review and sign-off

Walk through this checklist **with the sponsor present**. Do not skip.

1. Read the "Authorized activities" table aloud. Ask: "Is anything missing or too risky?"
2. Read the "Excluded" list aloud. Ask: "Anything to add?"
3. Read the "Stop conditions" aloud. Ask: "Can you say the stop word right now?" (Test it.)
4. Ask HR and legal to confirm they have read the RoE.
5. Ask the sponsor to confirm they understand that the test may upset some people, and that the no-discipline commitment applies.
6. Both sides sign. Both sides keep a copy.
7. Share the **summary page** (not the full RoE) with testers.

---

## Worked example (so you can see it filled in)

**Client:** a mid-sized logistics company, 600 employees, two offices (one owned, one leased), outsourced help desk.

| Decision | Result |
|----------|--------|
| Goal | Measure phishing resilience and test the help desk's verification process |
| Attack types in scope | Phishing (email), vishing against help desk |
| Attack types out of scope | Physical entry, USB, smishing, executives |
| Target pool | Random sample of 250 employees from the 540 non-excluded staff |
| Excluded | 12 on leave, 6 in HR processes, all executives, SOC team |
| Circles | Circle 1: CISO, engagement owner, legal. Circle 2: HR lead, emergency contacts. Circle 3: SOC lead (informed of indicators but not timing) |
| Help desk | Vendor consent obtained; test uses 5 client-provided test accounts |
| Windows | Tue to Thu, 10:00 to 16:00 local time, avoiding month-end |
| Reporting | Group size minimum 10; no per-person results; sponsor agrees to no-discipline clause |
| Stop word | Set and tested |

---

## Common scope mistakes (and the fix)

| Mistake | What happens | Fix |
|---------|--------------|-----|
| Scope says "all employees" | You hit executives, leave-takers, and sensitive staff | List exclusions and use a locked list |
| Scope says "physical security assessment" with no addresses | Dispute about which buildings | List every address and area |
| Verbal "yes" to a new target | Later: "we never approved that" | Written change request |
| No stop word | Confusion when someone calls | Set and test one |
| No plan for guards or police | Detentions | Letter, contacts, landlord consent |
| SOC not considered | Real incident response triggered | Choose blind or informed design |
| No data deletion plan | Personal data sits on laptops forever | Retention table and written deletion confirmation |
| Pretext approved verbally | Tester improvises something offensive | Approved pretext cards (Section 3) |
| Test during a crisis | People already stressed | Blackout dates |
| Vendor staff tested without vendor consent | Contract dispute | Third-party consent |

---

## Section 2 checklist

- [ ] RoE team named, with phone numbers
- [ ] Kickoff meeting held; notes confirmed by sponsor
- [ ] Attack type menu completed (in, out, with conditions)
- [ ] Default and client-specific exclusions written
- [ ] Target pool built, randomized IDs assigned, exclusions applied, list locked
- [ ] Need-to-know circles defined; blind vs informed SOC decided
- [ ] Escalation table complete; stop word and verification word set and tested
- [ ] Success criteria written; stop conditions written
- [ ] Schedule and blackout dates agreed
- [ ] Data handling table agreed
- [ ] Communication channels and daily check-ins set
- [ ] Change control agreed
- [ ] RoE document signed by both sides

---

# SECTION 3: PLANNING (GOALS, METRICS, OSINT, INFRASTRUCTURE, PRETEXTS)

**Goal of this section:** by the end, you know **why** you are testing, **what you will measure**, how you gather public information **safely**, how you build your tools **carefully**, and how you keep every pretext **inside the rules**.

---

## 3.1 How to create and choose the goal of the engagement

### Step 1: Ask the sponsor the one question that matters

> "After this test is over, what decision will you make differently because of the results?"

If the sponsor cannot answer, you do not have a goal yet. Keep asking until you get a sentence like:

- "We will decide whether to buy phishing-resistant MFA."
- "We will decide whether to change the help desk verification process."
- "We will decide whether to redesign visitor entry."
- "We will show the board whether training is working."

### Step 2: Pick the goal type

| Goal type | The question it answers | Best methods | Best primary metric |
|-----------|-------------------------|--------------|---------------------|
| **A. Awareness measurement** | "How do our people respond to suspicious messages?" | Phishing campaigns in groups | Report rate, click rate |
| **B. Technical control testing** | "Do our filters, MFA, and endpoint tools stop social engineering attempts?" | Phishing delivered with and without allow-listing; MFA flow tests with test accounts | Delivery/block rate, MFA bypass success on test accounts |
| **C. Process testing** | "Do our procedures (help desk, payments, visitor entry) hold up under pressure?" | Vishing, scripted help desk tests, payment-change simulations | Verification adherence rate |
| **D. Physical security testing** | "Can an unauthorized person reach sensitive areas?" | Tailgating, pretext visits, flag placement | Entry success rate, challenge rate |
| **E. Objective-based (red team style)** | "Can an attacker achieve a specific business objective using people and process?" | Combination of the above over weeks | Objective achieved or not, detection time |
| **F. Detection and response testing** | "Do we notice and respond?" | Any method, observed by SOC | Time to detect, time to contain |

### Step 3: Write the goal statement using this formula

> "We want to learn **[question]** so that we can decide **[decision]**, measured by **[metrics]** over **[period]** with **[target pool]**."

**Example 1 (awareness):**
"We want to learn whether employees in Finance and HR report suspicious emails, so we can decide whether to add a reporting button and targeted training, measured by report rate and time to first report over two weeks with 120 Finance and HR employees."

**Example 2 (process):**
"We want to learn whether the help desk follows the identity verification procedure when a caller sounds urgent, so we can decide whether to change the procedure, measured by verification adherence rate across 30 calls to five test accounts."

**Example 3 (physical):**
"We want to learn whether an unbadged person can reach the finance floor in our owned building, so we can decide on turnstile and escort policies, measured by entry success and challenge rate over three attempts on three days."

### Step 4: Check the goal against these questions

| Check | Pass condition |
|-------|----------------|
| Does the goal lead to a decision? | Yes |
| Can it be measured? | Yes, with the metrics below |
| Is it achievable within the window and scope? | Yes |
| Is it fair to employees? | The goal is about processes, not catching people |
| Is it safe? | No stop condition is likely to trigger |

### Step 5: Pick **one** primary goal and at most two secondary goals

Too many goals produce shallow results. Write the primary goal first and refuse extras that do not support it.

---

## 3.2 Which metrics to use, and how to measure them

### The principle: measure behavior that leads to safety, not just failure

A low click rate does not mean you are safe. A high report rate means the organization can **detect** attacks. Always measure **both** failure and defensive behavior.

### Phishing metrics (definitions and formulas)

| Metric | Formula | How to measure | Warnings |
|--------|---------|----------------|----------|
| **Delivery rate** | Delivered / Sent | Mail logs and bounce data | Bounces and auto-replies must be excluded from later denominators |
| **Open rate** | Opened / Delivered | Tracking pixel | **Unreliable.** Many email clients preload images, inflating opens. Do not make decisions on this |
| **Click rate** | Unique human clicks / Delivered | Landing page logs with unique token per recipient | Remove clicks from security scanners (see "Cleaning the data") |
| **Submission rate** | Unique human submissions / Delivered | Landing page event that records "form submitted" (never the content) | Count each person once |
| **Click-to-submit ratio** | Submissions / Clicks | Same logs | Shows how convincing the page was after the click |
| **Report rate** | Unique reporters / Delivered | Reporting button logs, SOC ticket data, mailbox rules | The most important metric; ask the client how reports are collected |
| **Time to first report** | Time from first delivery to first report | Timestamps | Shorter is better; shows how quickly the company could react |
| **Median time to report** | Median of report times | Timestamps | Median is better than average (outliers) |
| **Report-before-click ratio** | Reporters who reported before clicking / all reporters | Combined logs | Shows careful behavior |
| **Repeat risk rate** | People who clicked in two or more campaigns / campaign participants | Needs identity across campaigns, so only if the RoE allows; use pseudonyms | Use for training design only, never discipline |
| **Detection time (SOC)** | Time from first click or submission to SOC alert | SOC logs | Needs SOC cooperation |
| **Containment time** | Time from alert to action (mail purge, password reset, session revoke) | Incident logs | Shows response speed |
| **MFA block rate** | Attempts stopped by MFA / attempts using a valid credential on a **test account** | Test-account logs | Only on test accounts |

### Vishing and help desk metrics

| Metric | Formula | How to measure |
|--------|---------|----------------|
| **Contact rate** | Calls answered / calls placed | Call log |
| **Engagement rate** | Calls where the agent engaged with the request / calls answered | Call log notes |
| **Verification adherence rate** | Calls where the agent completed **all** required verification steps / calls that reached a verification point | Scored checklist (below) |
| **Disclosure rate by category** | Calls where information in category X was disclosed / calls reaching that point | Notes by category (for example "employee ID", "manager name"), **not the values** |
| **Unauthorized action rate** | Calls where the agent performed a restricted action (reset, change, unlock) / calls requesting it | Test-account logs |
| **Escalation or report rate** | Calls where the agent escalated or reported the suspicious call / calls answered | Help desk logs, SOC reports |
| **Time to report** | Time from end of call to internal report | Timestamps |

**Verification adherence checklist (example; replace with the client's real procedure):**

| Step required by procedure | Done? (Y/N) | Notes |
|----------------------------|-------------|-------|
| Asked for employee ID | | |
| Verified using a callback to the number on file | | |
| Asked a question not answerable from public sources | | |
| Required manager approval for high-risk request | | |
| Logged the request in the ticketing system | | |
| Refused to proceed when verification failed | | |
| Escalated or reported the suspicious call | | |

### Physical metrics

| Metric | Formula | How to measure |
|--------|---------|----------------|
| **Entry attempt success rate** | Successful entries / attempts | Attempt log |
| **Challenge rate** | Attempts where someone challenged / attempts | Attempt log |
| **Time to challenge** | Time from entry to first challenge | Timestamps |
| **Depth reached** | Which zone (lobby, floor, secure room) was reached | Flag photos |
| **Badge and escort adherence** | Observed compliant people / observed people | Observation notes (counts only, no identification) |
| **Report rate** | Reports to security / attempts | Security logs |
| **Time to response** | Time from report to security response | Security logs |

### Cleaning the data (very important)

Real campaign data is messy. Clean it using these steps:

1. **Remove bounces and auto-replies** from "Delivered."
2. **Identify security scanners.** Email gateways often click every link within seconds. Flag clicks that happen within a very short time of delivery (for example under 5 seconds), come from known scanner IP ranges or user-agent strings, or hit multiple links at once.
3. **Ask the client for the gateway's scanner IP ranges**, and exclude them.
4. **Count each person once** per metric (unique).
5. **Keep a copy of the raw data** and a record of the cleaning rules so a reviewer can reproduce your numbers.
6. **Report both raw and cleaned numbers** if the difference is large.

### Sample size: how many people do you need?

Small groups give noisy numbers. Use this rule of thumb (95 percent confidence, worst case):

| People in group (n) | Rough margin of error |
|---------------------|----------------------|
| 25 | plus or minus 20 percent |
| 50 | plus or minus 14 percent |
| 100 | plus or minus 10 percent |
| 250 | plus or minus 6 percent |
| 400 | plus or minus 5 percent |
| 1,000 | plus or minus 3 percent |

(Rule: margin of error is about 1 divided by the square root of n.)

**How to use it:**

1. If a department has 25 people, do not say "Finance is twice as bad as HR" from small differences. State the margin.
2. When a company is small, test everyone and call it a "census" rather than a sample, and say the numbers are descriptive.
3. If you compare two campaigns, differences smaller than the margin of error are probably noise.

### Make results comparable: the difficulty of each phish

A click rate means nothing without knowing how hard the email was. Rate each phishing email's difficulty before sending:

1. Use a difficulty method such as the **NIST Phish Scale**, which rates a phishing email by the number of detectable cues and by how well the message fits the recipient's work context.
2. Record the rating in the report next to the result.
3. Never compare click rates between campaigns without comparing difficulty.

### Baseline and targets

1. If the client has previous results, use **their own** history as the baseline.
2. If not, treat your test as the baseline and say so.
3. Avoid quoting "industry averages" as truth. Different test designs give wildly different numbers.
4. Suggest targets in terms of **reporting** (for example "median time to first report under 30 minutes") rather than only "click rate under 5 percent."

### Build the measurement plan (fill this in before the test)

| Metric | Data source | Who collects | Cleaning rule | Where recorded |
|--------|-------------|--------------|---------------|----------------|
| | | | | |

---

## 3.3 Passive OSINT: how to do it safely and usefully

### What "passive" means

**Passive** means you look at information that is already public **without interacting** with the client's systems or people. You do not send, call, log in, or scan.

**Active** means you touch the target in any way (for example visiting a restricted page, calling the number, scanning their network). Active reconnaissance needs RoE approval.

### Why do OSINT at all?

1. To make pretexts realistic **at the company and role level** (for example "an IT notice"), not at the personal level.
2. To measure what an attacker could learn (this exposure is itself a finding).
3. To check the target list and the physical setup (building, entrances).

### Allowed sources (usually)

| Source | What you learn | Caution |
|--------|----------------|---------|
| Company website and blog | Services, structure, locations, naming style | None |
| Press releases, news articles | Events, vendors, changes | Do not exploit real crises |
| Job postings | Technology in use, team names | Role-level only |
| Public business registries and filings | Legal names, directors, addresses | None |
| Public professional profiles | Roles, org structure (counts and titles) | Do **not** collect personal details beyond role |
| Certificate transparency logs and public DNS | Domains and subdomains the company uses | Passive lookups only |
| WHOIS and registrar data | Domain ownership | May be redacted |
| Search engine results | Public documents, metadata | Do not open documents that contain personal data you do not need |
| Maps and street imagery | Building layout, entrances, parking | Passive only |
| Public review sites | Culture and complaints | Do not use complaints against people |
| Public social media of the **company** | Branding, events, badge styles | Not personal accounts |

### Forbidden or restricted sources

| Source | Why |
|--------|-----|
| Personal social media of employees (family, hobbies, health, religion, politics, location history) | Out of authority; hurts trust; privacy law |
| Data from breaches or leaked databases | Possession can be unlawful; storing credentials breaks Section 4 rules. If the client wants breach exposure checked, ask **them** to use a domain-owner verified breach lookup service and give you summary counts |
| Fake accounts on social platforms | Platform rules and law |
| Private groups, forums, or direct messages | Not public |
| Anything behind a login | Not passive |
| Information about the client's customers, partners, or families | Out of scope |

### Step-by-step OSINT workflow

1. **Write the questions** the OSINT must answer (for example: "What do emails from IT usually look like? Which departments exist? Where are the entrances?"). If a source does not answer a written question, skip it.
2. **List approved source types** (table above) and get them in the RoE.
3. **Time-box** the work (for example 2 days). OSINT can expand forever.
4. **Log each finding** in the OSINT log (template below) with its source and purpose.
5. **Classify each item** as **Org-level** (company), **Role-level** (job title), or **Personal** (about an individual).
6. **Delete personal items** that are not needed. If a personal item is truly needed for a pretext, ask the sponsor for written approval.
7. **Summarize** the findings in one page for the report (what public information made an attack easier).
8. **Delete or archive** the raw notes at the end of the engagement according to the retention rule.

### OSINT log template

| # | Date (UTC) | Source (URL or name) | What was found | Class (Org / Role / Personal) | Purpose (question it answers) | Keep or delete |
|---|------------|----------------------|----------------|-------------------------------|-------------------------------|----------------|
| | | | | | | |

### Data minimization rules (print and pin)

1. **Role, not person.** Use "someone in Finance" rather than a named person whenever you can.
2. **Work, not life.** Never collect family, health, religion, politics, finances, or location history.
3. **Need, not nice-to-know.** If you cannot say why you need it, delete it.
4. **Encrypt, restrict, expire.** Encrypt notes, limit access, set a deletion date.
5. **The "would they be comfortable" test.** If the employee read your notes, would they feel you stayed professional? If not, remove it.

### What to include in the report from OSINT

- "A person with public internet access could determine X, Y, and Z about the company." (Org-level findings)
- Recommendations (for example: reduce detail in job postings, remove badge photos from public posts, standardize out-of-office messages).

---

## 3.4 Setting up infrastructure carefully

### Principles

1. **One engagement, one set of assets.** Never reuse domains, servers, mailboxes, or numbers across clients.
2. **Everything is logged, labeled, and removable.** You must be able to switch it off in minutes.
3. **No real harm.** Landing pages are benign. Payloads, if any, are harmless markers.
4. **No real credentials stored.** (Section 4.)
5. **Secure your own tools.** If your infrastructure is compromised, employee data could leak.

### Step-by-step setup

**Step 1: Choose domains**

- **Preferred:** use a **client-owned** domain or a neutral lookalike registered and controlled by the client (the client registers it and gives you access). This avoids trademark and legal issues.
- **If you register a lookalike yourself:** get written permission from the brand owner (the client) in the contract or RoE; use a registrar that allows it; record ownership details; set a calendar reminder to release or transfer the domain at the end.
- **Neutral domains** (not imitating any brand) are the safest and are suitable for many tests.

**Step 2: Set up sending infrastructure**

1. Use a reputable email or messaging service that allows security simulations (read their terms).
2. Configure standard email authentication (SPF, DKIM, DMARC) on your sending domain so messages behave like normal mail and are not blocked for the wrong reasons.
3. Create separate sending addresses per scenario.
4. Add a **header or flag** you can identify in logs (for example an internal tag), and give the SOC deconfliction contact the sender indicators (domain, IP) at the agreed stage.

**Step 3: Decide the allow-listing question (write it in the RoE)**

| Option | Meaning | When to choose it |
|--------|---------|-------------------|
| **Allow-list** | You ask IT to let your messages through filters so you test **people** | Goal A (awareness) |
| **No allow-list** | Your messages face filters as a real attacker would | Goal B (technical controls) |
| **Both** | A small first wave without allow-listing, then a second wave with it | When you want both results |

Record which option was used for each wave, because it changes how to interpret the numbers.

**Step 4: Build landing pages**

1. Host over HTTPS with a valid certificate.
2. Include a unique token per recipient in the link so you can count events without identifying anyone to the client.
3. The page records **events** only (visited, form submitted) and **never** stores what the person typed (Section 4).
4. After an interaction, show a short, kind, educational message (see Section 4).
5. No downloads of anything harmful. If you need to test file handling, use a harmless document with a unique marker.

**Step 5: Telephony and SMS**

1. Use a reputable provider; tell them it is a security simulation.
2. Register sender IDs and message templates where required by local rules.
3. Use numbers that are not tied to your personal identity.
4. Prepare a short, approved script for each call (Section 3.5).
5. Keep a call log: date, time, number called (role-level), outcome, no recording unless allowed.

**Step 6: Logging**

1. Synchronize all clocks (use UTC).
2. Log every email sent, every landing page request, every call, every physical attempt.
3. Keep logs in an access-controlled place.
4. **Do not log request bodies** that might contain what people typed.
5. Mask or hash anything that looks like personal data.

**Step 7: Kill switch**

1. Know how to take down every landing page in under five minutes.
2. Know how to stop all scheduled sends and calls in under five minutes.
3. Test the kill switch before launch.
4. Give the engagement owner a way to request a stop (phone call to your lead).

**Step 8: Protect your own infrastructure**

1. Use multi-factor authentication on every account.
2. Use full-disk encryption on every laptop.
3. Give each tester a separate account; no shared logins.
4. Keep infrastructure patched.
5. Use a password manager.
6. Restrict access to the target ID key (the key that links pseudonym to person).

### Pre-launch checklist

- [ ] Test email delivered to a controlled test mailbox; headers checked
- [ ] Authentication (SPF, DKIM, DMARC) passes
- [ ] Landing page loads over HTTPS
- [ ] Landing page records events but **not** typed content (verified with a dummy value, then logs searched for the dummy value)
- [ ] Unique tokens work; counts match
- [ ] Kill switch tested
- [ ] SOC deconfliction contact has indicators
- [ ] Allow-listing configured as agreed (or deliberately not)
- [ ] Telephony numbers and templates registered
- [ ] Time zones and send times checked
- [ ] Backup plan if a provider blocks you
- [ ] Emergency contacts confirmed available

### Teardown checklist (do this within days of the end)

- [ ] Take down all landing pages
- [ ] Stop and delete all mailboxes and sender configurations
- [ ] Remove DNS records pointing to your infrastructure (avoid leaving dangling records anyone could take over)
- [ ] Release, transfer, or park domains as agreed
- [ ] Revoke API keys and provider accounts used for the test
- [ ] Delete or securely archive logs per the retention rule
- [ ] Destroy the pseudonym key when retention ends
- [ ] Write a short teardown confirmation and send it to the sponsor

---

## 3.5 Staying within the agreed rules while pretexting

A **pretext** is the story you use. Every pretext must be written down and approved **before** it is used.

### The pretext card

Fill in one card per pretext. Nothing is used without an approved card.

```markdown
PRETEXT CARD  ID: P-[number]   Version: [n]

Channel:            [email / phone / SMS / physical]
Target group:       [group name, not individuals]
Story (2-3 lines):  [what the message or caller says]
Reason it is believable (work context): [role-level reason]
Desired action:     [click, call back, provide category X, enter building]
What we will NOT say or do (red lines): [list]
Information we may reveal if asked: [list]
Information we must never reveal: [list]
Exit line if challenged: [exact words]
Stop conditions specific to this card: [list]
Risk rating: [low / medium / high]
Approved by (sponsor name, date): [ ]
Reviewed by (legal/HR, date): [ ]
```

### Approval process (step by step)

1. Tester writes the card.
2. A **second tester** reviews it against the exclusions in the RoE.
3. Engagement owner reviews it.
4. HR and legal review it (HR especially for tone and topic).
5. Sponsor approves in writing (email is fine).
6. The approved card is locked (version number).
7. **Only approved cards are used.** Any change creates a new version and needs a new approval.

### Acceptable versus unacceptable pretext themes

| Acceptable (usually) | Not acceptable |
|----------------------|----------------|
| A routine IT notice about a mailbox or system update | Fake layoff, pay cut, or restructuring notice |
| A shared-document notification | Fake bonus or salary increase |
| A delivery or courier notice | Fake legal action or tax notice |
| A generic company policy reminder (not sensitive) | A fake health or safety alert |
| A vendor invoice query (finance test, with finance in the circle) | A message about a real recent incident or real crisis |
| A visitor or contractor arrival (physical test) | A fake police, government, or fire inspector |
| A routine help desk call (non-emotional) | A call that claims a family emergency or medical emergency |
| A meeting or calendar invite | A message that targets a named person's personal life |

### Real-time rules while live

1. **Stay on the card.** If the conversation goes somewhere new, steer back or exit politely.
2. **No threats, no pressure about consequences to the person.** Time pressure phrases should be mild ("when you have a moment today"), not threats ("your account will be deleted and you will be in trouble").
3. **Never invent real names.** Do not use the name of a real employee unless approved in writing.
4. **Never claim authority you do not have** (police, regulator, medical staff).
5. **Never ask for more than the objective needs.** If the goal is a category of information, stop at the category.
6. **If the person is hesitant, do not push harder.** Hesitation is good behavior; record it and move on.
7. **If the person becomes upset, use the gentle-ending script** (Section 2, Step 8).
8. **Two-person rule for vishing:** one person speaks, one listens, takes notes, and has the authority to say "stop." This prevents drift.
9. **Do not record** unless the RoE and legal review allow it and the law of every involved location permits it.
10. **Log everything** immediately after each interaction.

### How to prevent "scope drift"

| Drift | How it starts | The fix |
|-------|---------------|---------|
| "While we are here..." | You notice an unlocked door outside the approved area | Do not enter; log it; request change in writing |
| "Just one more email" | A new idea mid-campaign | Write a new card and get approval |
| "The CEO would be interested" | A senior person is reachable | Only if the CEO is already in scope |
| "This person is chatty" | A conversation goes deep | Return to the script; do not fish for extras |
| "Let's use their real name" | To make it convincing | Not allowed without written approval |

### Vishing script structure (template)

| Part | Content | Example type |
|------|---------|--------------|
| Opening | Neutral greeting and role | "Hi, this is [first name] from [department]." |
| Reason | Plausible and mild | "I'm trying to confirm a routine item for a ticket." |
| Ask | One narrow request tied to the objective | One category only |
| If challenged | Honest exit | "No problem, I will contact them another way." |
| Stop conditions | From the card | Distress, personal matters, real emergency |
| After the call | Log, disclose if required by RoE | Category recorded, not values |

---

## 3.6 Dry run and go / no-go

1. **Internal dry run:** send test messages to your own team; call your own number; walk a safe route. Fix what breaks.
2. **Go/no-go meeting** with the sponsor, engagement owner, and legal:
   - Review the Section 1 checklist (all ticked).
   - Review the Section 2 checklist (all ticked).
   - Review the pre-launch checklist (all ticked).
   - Confirm blackout dates are clear.
   - Confirm emergency contacts are available.
   - Confirm the stop word works.
3. **Decision:** Go, No-go, or Go with changes (write down which).
4. **Record** the decision, the time, and who attended.

---

## Section 3 checklist

- [ ] Goal statement written using the formula; one primary goal
- [ ] Metrics chosen, with formulas, data sources, and cleaning rules
- [ ] Sample size checked against margin of error
- [ ] Phish difficulty rating method chosen
- [ ] OSINT questions written; source types approved; log in use; personal data deleted
- [ ] Domains owned or approved; sending infrastructure authenticated
- [ ] Allow-listing decision recorded
- [ ] Landing pages verified to store events only
- [ ] Kill switch tested
- [ ] Pretext cards written, reviewed, and approved
- [ ] Dry run complete; go/no-go recorded

---

# SECTION 4: ETHICS AND EMPLOYEE WELFARE

**Goal of this section:** by the end, you know how to make sure **no employee is fired, punished, shamed, or exposed because of your test**, how to word findings so they point at **systems and processes**, and how to prove that a credential was submitted **without ever keeping the credential**.

---

## 4.1 The big idea (explained simply)

Imagine a school with a door that has a broken lock. A student walks through it. Do you blame the student, or do you fix the lock?

In security, people are the student, and the **controls** (filters, verification steps, locks, training, reporting tools) are the lock. If an attack works, the question is: **"What should have stopped it, and why did it not?"**

Three ideas support this:

| Idea | Plain meaning | Why it helps |
|------|---------------|--------------|
| **Just culture** | People who make honest mistakes in a flawed system are supported, not punished | Staff report problems when they are not afraid |
| **Swiss cheese model** | Every defense has holes; failures happen when holes line up | A click is one hole; the missing filter, missing MFA, and slow reporting are the others |
| **Test the process, not the person** | The unit of analysis is the control, not the human | Produces fixes that scale |

**Why this is not soft:** if employees fear punishment, they stop reporting real attacks. A company that punishes clickers trains staff to **hide** incidents. That makes the company less safe.

---

## 4.2 Before the test: lock in a no-discipline commitment

You cannot fix this after the report is written. Do it first.

### Step-by-step

1. **Raise it at the kickoff.** Say: "Before we start, we need an agreement about how results will be used."
2. **Explain the risk in one sentence.** "If someone is disciplined because of this test, staff will stop trusting security and stop reporting real attacks."
3. **Ask the sponsor to commit** to the no-discipline clause (Section 1 sample).
4. **Ask HR to confirm in writing** that individual results will not enter personnel records or performance reviews.
5. **Ask legal** whether any local employment rule prevents the commitment.
6. **Put it in the RoE** (Section 17 of the template) and in the contract.
7. **Agree what happens if the commitment is broken:** you may stop the engagement, withhold individual-level data, and withdraw from the project.
8. **Tell the sponsor how you will word the results** (this section) so they know what to expect.

### What to say if the sponsor pushes back

| Pushback | Your response |
|----------|---------------|
| "People who click should be held accountable." | "Attackers will eventually fool someone in any company. Our job is to make sure the company's controls still hold. Punishment stops reporting, and reporting is your best early warning." |
| "We need names to give extra training." | "We can send training invitations automatically through your learning system based on anonymous tokens, so managers do not see names. Or we can design training by department." |
| "Our regulator wants evidence of accountability." | "Accountability is shown by corrective action on controls and training coverage, not by individual blame. We can document that in the report." |
| "We already told staff there would be consequences." | "Then we recommend a communication to reverse that before the test. Otherwise staff will fear reporting." |
| "The CEO insists." | "We can brief the CEO. If the policy cannot be agreed, we should pause or change the test design." |

### Sample policy paragraph the client can adopt

> "Security awareness exercises exist to improve the company's defenses. Results are used only in aggregate. No employee will be disciplined, penalized, or publicly identified because of how they responded to a simulated attack. Employees who report suspicious messages, including simulated ones, are thanked and encouraged."

---

## 4.3 Protect identities: anonymize and aggregate

### Step-by-step

1. **Replace names with random IDs** at the start (Section 2, Step 5).
2. **Keep the key** (ID to person) separate, encrypted, and accessible only to named testers. If you do not need it, destroy it early.
3. **Avoid unique details** in reports: timestamps to the second, unusual job titles, office desk numbers, or free-text that could identify someone.
4. **Report by group, not by person.**
5. **Set a minimum group size** (for example at least 10 people per reported group). If a department is smaller, merge it with similar departments or only report it in a larger bucket.
6. **Suppress small cells.** If a result affects fewer than 5 people in a group (for example 2 people clicked in a group of 14), say "fewer than 5" or "a small number" instead of the exact figure, because exact numbers in a small group can identify people.
7. **Do not publish rankings** of departments, managers, or teams. Rankings create blame.
8. **Do not produce "repeat offender" lists.** If you need repeat-risk data for training design, provide it only to the training system as anonymous tokens.
9. **Strip personal data from screenshots** (names, emails, avatars).
10. **Never include free-text quotes from employees** in the report without removing identifying details.

### Reporting table: allowed versus not allowed

| Allowed | Not allowed |
|---------|-------------|
| "Of 250 recipients, 14 percent clicked." | "Alex in Finance clicked within 2 minutes." |
| "In departments with at least 10 recipients, click rates ranged from X to Y percent." | "Marketing was the worst department." |
| "A small number of recipients (fewer than 5) in one group submitted the form." | "Two people in Legal submitted the form." |
| "Median time to first report was 7 minutes." | "Priya was the first to report." (Unless the reporter consents to being thanked publicly) |
| "Several recipients forwarded the email to security." | "Names of reporters: ..." |

> **TIP:** Reporters can be thanked privately by the security team through the normal reporting channel. Do not name them in your report.

---

## 4.4 How to frame findings around process and control failures

### The five-step reframe (use this for every finding)

| Step | Question | Example |
|------|----------|---------|
| 1. **Observe** | What happened? | "Employees in Group 3 entered credentials on the test page." |
| 2. **List the controls** | Which controls had a chance to stop it? | Mail filter, link scanning, browser warning, MFA, reporting tool, SOC |
| 3. **Find the holes** | Which controls failed or were absent, and why? | "The filter did not flag the new lookalike domain; no external sender banner; MFA could be relayed; no visible report button." |
| 4. **Ask 'why was it easy?'** | What made the attack plausible? | "The message matched a real process (document sharing); the domain looked like the real one." |
| 5. **Recommend** | What would make it hard next time? | "Add lookalike-domain detection, external banner, phishing-resistant MFA, a one-click report button." |

### The "5 Whys" example

**Finding:** 12 percent of recipients submitted credentials.

1. **Why did they submit?** The page looked like the real login page.
2. **Why did it look real?** The fake domain was visually close to the real domain and the company's login page is simple.
3. **Why was the domain not blocked?** There is no monitoring or blocking of newly registered lookalike domains.
4. **Why was a credential alone enough?** The sign-in method relies on passwords plus a code that can be entered on the fake page.
5. **Why is it hard for staff to notice?** The email client gives no cue that the sender is external and the link text hides the real address.

**Result:** the report lists five control gaps, none of which says "the employee should have known better."

### Wording translation table

Use the right column in the report, the debrief, and emails.

| Instead of this (blame) | Write this (process) |
|-------------------------|----------------------|
| "Employees fell for the phishing email." | "The message was not stopped by technical controls, and a share of recipients acted on it." |
| "John clicked the link." | "A recipient in Group 2 clicked the link." |
| "Users were careless." | "The message closely matched routine company mail, and no visual or technical cue reliably distinguished it." |
| "The help desk agent failed." | "The verification procedure did not require a step that would have blocked the request." |
| "Staff do not care about security." | "Staff received no clear guidance on how to report the message, and the report channel is not obvious." |
| "Finance is the weakest department." | "Finance recipients encountered a message that closely resembled their normal invoice process; the process lacked an independent check." |
| "Training failed." | "Annual training does not provide timely, practical cues; consider shorter, role-based, in-the-moment guidance." |
| "The guard let us in." | "The entry procedure relied on individual judgment without a physical barrier or escort rule." |
| "Employees are the weakest link." | "Controls that depend on a person getting every decision right are fragile; add layers that do not." |
| "Repeat clickers." | "Participants who encountered multiple messages and responded to more than one, which suggests the cues were hard to spot." |
| "Employee violated policy." | "The policy was not clearly reinforced at the point of action; consider reminders at the moment of risk." |
| "The user entered their password." | "A credential submission event occurred on the simulated page." |
| "The user ignored the warning." | "The warning did not appear, or appeared in a form that was easy to miss." |
| "Stupid mistake." | (Never use this.) |
| "The CFO was fooled." | "A senior finance role was reached by a pretext matching routine approval workflows." (Only if the CFO is in scope and the group is large enough) |

### How to write a finding (structure)

```markdown
### Finding F-[n]: [Control gap, written as what is missing or weak]

**What we observed (aggregate):** [numbers, groups of 10 or more]
**Which controls had a chance to stop it:** [list]
**Why it worked:** [process/design reasons]
**What made it hard for people to detect:** [cue analysis]
**Risk:** [likelihood x impact, with reasoning]
**Recommendation:** [fix at the control level]
**How to verify the fix:** [retest or metric]
```

Notice the title of the finding is about the **control gap**, such as "No monitoring of newly registered lookalike domains," not "Employees clicked."

### Do and don't list for tone

| Do | Don't |
|----|-------|
| Use "recipients," "participants," "group," "process" | Use names or nicknames |
| Show what worked first, then what did not | Open with failures |
| Explain **why** an attack was convincing | Imply it should have been obvious |
| Praise reporters in aggregate | Single out individuals |
| Use neutral colors in charts | Use red and "shame" language |
| Describe fixes that are easy for people | Recommend "more awareness" as the only fix |
| Say "we were able to" about the tester | Say "they failed to" about staff |

---

## 4.5 What to do when the client asks "who clicked?"

This will happen. Prepare the answer.

### Step-by-step

1. **Pause and stay friendly.** Say: "I understand the question. Let me explain what we can offer."
2. **Remind them of the signed commitment.** "Our contract and RoE say results are used in aggregate."
3. **Offer alternatives that meet the real need:**
   - Group-level results for department training plans.
   - Automatic enrolment in a micro-training module by anonymous token, with no manager visibility.
   - A list of **controls** to fix.
   - A second test later to show improvement.
4. **If they still insist, escalate to the sponsor.** Say: "We can't provide individual results without breaching the agreement. If the policy needs to change, we should pause the engagement and discuss."
5. **If the answer is still yes,** decide whether to continue the engagement. You are allowed to stop. Document the decision.
6. **Never hand over individual data informally** (for example in a chat message or a spreadsheet "just this once").

### Script

> "We keep individual results out of the report on purpose. Our agreement with you says the results will not be used for discipline, and in our experience that is what keeps staff reporting real attacks. What we can give you is group results, a list of specific control gaps, and an automatic way to send training without exposing names."

---

## 4.6 The employee debrief: how to close the loop kindly

### Timing

- Send the debrief message **within a few days of the test ending** (as soon as the live window closes and the sponsor agrees).
- Do not wait weeks; people will have heard rumors.

### Step-by-step

1. **Write the message with the client** (their voice, not yours).
2. **Lead with thanks.** Thank people who reported and thank everyone for taking security seriously.
3. **Explain the purpose.** "This was an authorized exercise to test our systems."
4. **Reassure.** "No individual results will be shared with managers or used in performance reviews."
5. **Explain what we learned (briefly)** at the company level.
6. **Give simple tips** that are easy to remember.
7. **Show how to report** (where, how).
8. **Offer a contact** for questions or concerns.
9. **Brief managers** a day before so they can answer questions kindly.

### Template (all-staff message)

> **Subject:** Thank you: our recent security exercise
>
> Hello everyone,
>
> Over the last few days we ran an authorized security exercise to test how well our systems and processes protect us against realistic attacks. Some of you may have received messages or calls that were part of the exercise.
>
> First, thank you. Several colleagues reported suspicious messages quickly, which is exactly what we want. Reporting is the most powerful thing any of us can do.
>
> To be clear: **nobody is being singled out or disciplined.** Results are looked at only in groups, and they will be used to improve our tools and processes.
>
> What we learned: some of the messages were convincing, and that tells us we need better technical protection and clearer guidance. We are already working on [two or three improvements].
>
> If you ever see something that feels wrong: [how to report: button, email, phone]. You will never get in trouble for reporting, even if it turns out to be a test or a mistake.
>
> Questions or concerns? Please contact [name, contact].
>
> Thank you again,
> [Sponsor name and title]

### Manager briefing (one page)

1. The test is over and was authorized.
2. Do not ask team members whether they clicked.
3. If someone mentions it, thank them for raising it and send them to [contact].
4. If someone seems upset, listen and offer the contact.
5. Do not discuss it in public channels.

### If an employee is upset

1. Listen without defending.
2. Say: "Thank you for telling me. You did not do anything wrong, and your result was not identified."
3. Offer a call with the security lead or HR contact.
4. Log the concern (without names) and report it to the sponsor.
5. Consider whether the pretext should be retired or changed.

---

## 4.7 How to prove a credential was submitted without keeping the credential

### The rule

**Never store, log, copy, or use a real password.** A real password is a key to a person's life. Your job is to prove that a person **would have** handed it over, not to keep it.

### Why this matters

- If you keep passwords, a breach of your systems exposes employees.
- If you test the passwords, you may be committing unauthorized access.
- Laws on data protection treat credentials as sensitive.
- Employees will be furious if their real passwords end up in a report.

### Method A: record the event, not the content (the default)

The landing page should only record **that** something happened, not **what**.

**Conceptual logic (pseudo-code, not a real tool):**

```
on page_visit(request):
    token = request.query["t"]            # unique, pseudonymous
    record_event(token, "visited", now_utc())
    show_page()

on form_submit(request):
    token = request.query["t"]
    record_event(token, "submitted", now_utc())
    # IMPORTANT: do NOT read, store, or log the form fields
    show_education_page()
```

Design rules:

1. The server never reads the password field.
2. The form can be built so the browser does not send the field value at all (for example, the submit action only sends the token and an event name).
3. No analytics, error tracking, or third-party scripts on the page (they could capture field values).
4. No request-body logging on the web server, reverse proxy, CDN, or firewall.
5. Do not record the length of what was typed.
6. Do not record browser autofill content.

### Method B: the canary test (prove your setup is clean)

Before launch, prove that nothing stores typed content.

1. Type a unique, fake value (for example `CANARY-NOT-A-REAL-PASSWORD-9271`) into the password field and submit.
2. Search **every place** logs or data could be stored for that string:
   - Web server and application logs
   - Reverse proxy, load balancer, CDN, and web firewall logs
   - Database and cache
   - Error tracker and monitoring tools
   - Browser developer tools and extensions on the test machine
   - Clipboard managers and password managers on the test machine
   - Backups and snapshots
3. The expected result is **zero matches**.
4. If you find any match, fix the cause and repeat.
5. Save a screenshot or note of the test as proof for the report.

### Method C: use client test accounts when you must prove "the credential would have worked"

Sometimes the client wants to know whether a stolen password plus a relayed code would let an attacker in. Do **not** test this with real employees' credentials.

1. Ask the client to create **test accounts** that look like real accounts (same policies, same MFA type, same roles).
2. Ask the client to give the credentials of those accounts to your team, or to have a **trusted agent** at the client run the flow.
3. Use those accounts to test the MFA and sign-in path.
4. Record the results as "MFA withstood / did not withstand this technique on the test accounts."
5. The client's identity team resets these accounts afterward.

### Method D: have the client verify without you seeing the password

1. The page records a hashed marker, not the password. (Only if your lawyer and the client agree; hashing a password is still processing it, so Method A is better.)
2. The preferred approach is: you tell the client "a submission event occurred for token T-xxxx," and the client does **not** need to verify the password. The event itself is the evidence.

### What to do if you accidentally capture real credentials

| Step | Action |
|------|--------|
| 1 | Stop the campaign for that page immediately |
| 2 | Do not read, copy, or use the captured values |
| 3 | Notify the engagement owner and sponsor by phone |
| 4 | Securely delete the captured data (or isolate it as the client directs) and record who did it and when |
| 5 | Ask the client to force a reset for affected accounts as a precaution |
| 6 | Document it as a **deviation** (what, when, how many, what was done) |
| 7 | Fix the cause before restarting |
| 8 | Report it honestly in the report |

### Where passwords accidentally hide (check each)

| Place | Why it is a risk |
|-------|------------------|
| Web server access or error logs | Sometimes log full URLs or bodies |
| Reverse proxy and CDN | May log request bodies for debugging |
| Web application firewall | May store blocked requests |
| Analytics or heatmap scripts | Record keystrokes |
| Error or monitoring tools | Capture request data on failures |
| Database and queue systems | Backups store everything |
| Browser extensions on the tester machine | Autofill or sync |
| Screenshots and screen recordings | May show typed text |
| Shared chat channels | Someone pastes a log |

### Proof for the report

In the report, show:

- The number of submission events (aggregate).
- The timeline of events.
- A screenshot of the landing page (with no real data).
- A statement of method: "The page recorded events only; no credentials were stored. Method verified with a canary test on [date]."

### Proof for other test types without capturing sensitive data

| Test type | Collect this | Do NOT collect this |
|-----------|--------------|---------------------|
| **Vishing** | Category of information disclosed (for example "employee ID provided"), time, outcome | Actual values (the ID number, the password, the answers) |
| **Help desk** | Whether verification steps were followed; whether the test account reset succeeded | Real accounts' details |
| **Physical entry** | Photo of a **flag** (a client-placed item such as a marked card or token in the target area), time, location | Photos of people, screens, documents, or personal belongings |
| **Badge or door tests** | Whether the door opened, time, area | Copies of real badges unless the RoE allows and the client approves |
| **USB drop** | Event that a device connected to the test marker | Contents of the computer |
| **Documents found** | "Documents containing sensitive-looking labels were present" | Photos of the documents' contents |

> **TIP:** The **flag method** is the safest physical proof. Ask the client to place a harmless marked item (for example a sealed envelope with a code) in the target area. You photograph the item, not the surroundings. No theft, no snooping.

---

## 4.8 The tester code of conduct (every tester signs this)

1. I will only act within the signed scope and RoE.
2. I will treat every person with respect, including people who catch me.
3. I will never use fear, threats, or personal crises as a pretext.
4. I will never touch, take, or photograph people's belongings, faces, or screens.
5. I will never keep real passwords or secrets.
6. I will stop immediately when a stop condition happens.
7. I will tell the truth in the report, including my mistakes.
8. I will protect personal data and delete it on schedule.
9. I will not discuss the engagement outside the agreed circles.
10. I will report anything unlawful or unsafe to my manager immediately.

---

## 4.9 Special situations

| Situation | What to do |
|-----------|------------|
| You find evidence of **real crime or misconduct** (for example fraud, harassment) | Stop related activity, notify the sponsor and your legal contact immediately, and do not investigate further; this is outside the test |
| You discover a **real attacker** or an active compromise | Stop, notify the sponsor and incident response, preserve what you saw, follow the real-incident clause |
| A target is **angry** and threatens a complaint | Stop, disclose politely, notify the engagement owner, offer the contact, and log it |
| The test **leaks** (screenshots on social media) | Notify the sponsor; decide together whether to end the phase early; use the leak as evidence of reporting culture if appropriate |
| An employee **guesses** the test and asks you | If the RoE allows, tell the truth gently; otherwise follow the exit script |
| A **minor** or visitor unexpectedly appears | Stop the physical interaction and leave the area |
| A **medical emergency** happens near you | Drop everything; help or call emergency services; notify the sponsor |
| A tester feels **uncomfortable** with a task | The tester may refuse; no penalty; discuss with the lead |

---

## Section 4 checklist

- [ ] No-discipline commitment signed (contract and RoE) and HR confirmation received
- [ ] Pseudonymous IDs in use; key stored separately
- [ ] Minimum group size set (for example 10); small-cell suppression rule set
- [ ] Wording translation table shared with all report writers
- [ ] Findings titled by control gap, not by people
- [ ] Plan ready for "who clicked?" requests
- [ ] Employee debrief message and manager briefing drafted
- [ ] Landing page records events only; canary test passed with zero matches
- [ ] Test accounts arranged for any MFA or credential-path testing
- [ ] Accidental-capture procedure agreed
- [ ] Flag method agreed for physical evidence
- [ ] Tester code of conduct signed by every tester

---

# SECTION 5: REPORTING (TIMELINE, RATES, CONTROLS, RECOMMENDATIONS, POSITIVES)

**Goal of this section:** by the end, you can turn a pile of logs into a clear timeline, fair success and failure rates, an honest picture of which controls worked, a prioritized set of recommendations, and a report that the client actually uses.

---

## 5.1 How to make a timeline of the techniques used

### Step 1: Keep the activity log **during** the test

You cannot build a good timeline afterward from memory. Start the log on day one. Every tester writes entries as things happen.

**Activity log template (one row per action):**

| Event ID | Time (UTC) | Tester | Channel | Technique | Pretext ID | Target group (pseudonym) | Action taken | Outcome | Evidence ref | Notes |
|----------|-----------|--------|---------|-----------|------------|--------------------------|--------------|---------|--------------|-------|
| E-001 | | | Email | Phishing wave 1 | P-01 | Group A | 125 emails sent | Delivered 121 | EV-01 | |
| E-002 | | | Email | Phishing wave 1 | P-01 | Group A | First click | | EV-02 | |
| E-003 | | | Phone | Vishing | P-04 | Help desk | Call 1 placed | Verification completed, request refused | EV-03 | |

**Rules for the log:**

1. Always use **UTC** (and note the local time zone in the header).
2. Write entries within minutes, not hours.
3. Record **facts**, not opinions.
4. Use pseudonymous IDs, never names.
5. Link every entry to evidence (screenshot, log extract, photo of flag).
6. Record deviations and stop events as they happen.

### Step 2: Add the "client-side" events

The timeline is more useful when it includes what the client did:

- When the first report came in
- When the SOC noticed
- When the filter blocked something
- When IT purged messages
- When the stop word was used (if ever)

Ask the client's SOC or IT team for their timestamps after the test and merge them into your log (with permission).

### Step 3: Build the timeline in three lanes

```
TIME (UTC)   TESTER ACTIONS                CONTROLS TRIGGERED              CLIENT RESPONSE
---------    ---------------------------   -----------------------------   --------------------------
Day 1 10:00  Wave 1 emails sent (Group A)  Gateway scanned links           (none)
Day 1 10:02  First click recorded          Browser showed no warning        (none)
Day 1 10:07  First report received         Reporting button used            SOC ticket opened
Day 1 10:31  SOC triage complete           Mail purge started               Messages removed from 61 inboxes
Day 1 10:40  Wave 1 ends                   Sender domain added to blocklist Deconfliction call
Day 2 11:00  Wave 2 (different pretext)    Gateway blocked 22%              (none)
Day 3 14:00  Help desk call 1              Verification script started      Request refused after step 3
Day 3 15:20  Help desk call 2              Verification skipped by agent    Test account reset
```

### Step 4: Present it clearly in the report

1. A **one-page overview timeline** (start, key events, end).
2. A **detailed timeline appendix** with every logged event (pseudonymous).
3. A **summary table by technique** (count of attempts, outcomes).
4. A **"time-to-event" table** (time to first click, time to first report, time to containment).

### Example summary table by technique

| Technique | Attempts | Blocked by technology | Stopped or reported by people | Reached objective | Not completed |
|-----------|----------|-----------------------|-------------------------------|-------------------|---------------|
| Phishing wave 1 | 125 recipients | 4 | 18 reported | 15 submitted | n/a |
| Phishing wave 2 | 125 recipients | 28 | 22 reported | 9 submitted | n/a |
| Vishing (help desk) | 10 calls | 0 | 4 refused and reported | 3 performed a reset | 1 call dropped |
| Physical entry | 6 attempts | 1 (turnstile) | 2 challenged | 3 reached the lobby; 1 reached floor 5 | 0 |

(Numbers above are illustrative only.)

---

## 5.2 How to calculate success and failure rates

### Step 1: Use clear outcome categories

Every attempt gets **one** outcome. Define these in the report so readers see the same definitions.

| Outcome | Meaning |
|---------|---------|
| **Blocked by technology** | A technical control stopped the attempt before a human saw it |
| **Stopped by a person** | A person recognized and refused the attempt |
| **Reported** | A person reported the attempt to the proper channel |
| **Partial** | The attempt got part of the way (for example click but no submission) |
| **Reached objective** | The tester achieved the defined objective |
| **Not completed** | Technical or scheduling problems prevented a fair attempt |

Note: "Stopped by a person" and "Reported" can overlap. Count them clearly (for example "reported" is a subset of "stopped").

### Step 2: Define the denominator for each rate

| Rate | Numerator | Denominator | Note |
|------|-----------|-------------|------|
| Delivery | Delivered messages | Sent | Excludes bounces |
| Click | Unique human clicks | Delivered | After cleaning scanner clicks |
| Submission | Unique human submissions | Delivered | Event only |
| Report | Unique reporters | Delivered | |
| Technology block rate | Messages blocked | Sent | |
| Help desk refusal rate | Calls where the agent refused | Calls where a restricted action was requested | |
| Physical entry rate | Successful entries | Attempts | |
| Challenge rate | Attempts challenged | Attempts | |

### Step 3: Use a funnel to show where attempts end

```
 Sent                      125  (100%)
  |
  v
 Delivered                 121  ( 97%)
  |
  v
 Opened (not reliable)      --  (do not use)
  |
  v
 Clicked                    24  ( 20% of delivered)
  |
  v
 Submitted event            15  ( 12% of delivered)
  |
 Reported                   18  ( 15% of delivered)
```

(Illustrative.) The funnel shows that more people **reported** than **submitted**, which is a positive finding.

### Step 4: Show rates **with** context

1. State the group size and the margin of error (Section 3).
2. State the **difficulty rating** of the phish.
3. State whether allow-listing was used.
4. State the time window.
5. If you ran more than one wave, show waves separately **and** combined.

### Step 5: Avoid misleading language

| Avoid | Use instead |
|-------|-------------|
| "Success rate 12 percent" (sounds like the tester succeeded) | "12 percent of delivered messages resulted in a submission event" |
| "Failure rate 12 percent" | Same as above |
| "The company failed the test" | "The test identified these control gaps" |
| "Only 88 percent passed" | "88 percent did not submit; 15 percent reported" |
| "Industry average is X" | "Compared with the company's earlier results" (if any) |

### Worked example

Campaign: 125 recipients across two groups; allow-listing used; difficulty rated "moderate."

| Measure | Count | Rate | Interpretation |
|---------|-------|------|----------------|
| Delivered | 121 | 97% of sent | Four messages bounced |
| Clicks (after removing scanner clicks) | 24 | 20% of delivered | Margin of error about plus or minus 9 percent |
| Submission events | 15 | 12% of delivered | Page was convincing to a minority |
| Reports | 18 | 15% of delivered | More reported than submitted |
| Median time to first report | 7 minutes | n/a | Fast reporting; containment could start quickly |
| Time from first report to purge | 24 minutes | n/a | Containment gap of about 24 minutes |

---

## 5.3 How to determine what controls worked

### Step 1: List the controls that **should** have acted

For each technique, list every control that could have stopped, detected, or contained it.

| Layer | Examples (phishing) | Examples (help desk) | Examples (physical) |
|-------|--------------------|----------------------|---------------------|
| **Prevent** | Email filtering, domain monitoring, link rewriting, DMARC | Verification procedure, callback rule | Turnstile, badge reader, visitor system |
| **Deter / slow** | External sender banner, MFA | Manager approval for resets | Escort rules, signage |
| **Detect** | User reporting, SOC alerts | Ticket review, anomaly alerts | Guards, cameras, door alarms |
| **Respond** | Mail purge, password reset | Account lock, callback to the real user | Guard response, lockdown |
| **Recover** | Session revocation, awareness note | Reset all affected accounts | Incident log, procedure update |

### Step 2: Score each control for each technique

Use these ratings:

| Rating | Meaning |
|--------|---------|
| **Effective** | The control stopped or substantially limited the attempt |
| **Partially effective** | It helped but let some attempts through, or acted too late |
| **Ineffective** | It was present but did not help |
| **Absent** | The control did not exist |
| **Bypassed** | The tester avoided it (for example by allow-listing as agreed) |
| **Not tested** | No opportunity to observe |

### Step 3: Fill in the Control Effectiveness Matrix

| Control | Technique | Opportunity to act | Rating | Evidence | Comment |
|---------|-----------|--------------------|--------|----------|---------|
| Email gateway | Phishing wave 2 (no allow-list) | Delivery | Partially effective | Blocked 28 of 125 | Missed new lookalike domain |
| DMARC on client domain | Spoofing test | Delivery | Effective | Spoofed messages rejected | Policy at enforcement |
| External sender banner | Phishing | Reading | Absent | No banner seen in screenshots | Add banner |
| Browser warning | Phishing landing | Click | Ineffective | No warning displayed | Domain not yet flagged |
| MFA (push-based) | Test account flow | Sign-in | Partially effective | Approved on second prompt | Prompt fatigue risk |
| Report button | Phishing | Reading | Effective | 18 reports | Fast |
| SOC alerting | Phishing | After clicks | Ineffective | No alert on first submissions | Add detection |
| Help desk verification | Vishing | Call | Partially effective | 4 of 10 enforced all steps | Procedure unclear |
| Turnstile | Physical | Entry | Effective | Blocked tailgating in lobby | Good |
| Guard challenge | Physical | Floor 5 | Partially effective | 2 of 6 challenged | Polite challenge script needed |

### Step 4: Draw the "layers" picture

For each technique, write down how many layers stopped it:

```
Phishing wave 2 (no allow-list)

Layer 1  Gateway            [##########-----] stopped 22%
Layer 2  Browser warning    [---------------] stopped 0%
Layer 3  User judgment      [#####----------] reported 15%
Layer 4  MFA                [-------#-------] not applicable (credentials only)
Layer 5  SOC detection      [---------------] no alert
Result: 9 of 125 recipients reached the submission step
```

The picture shows the "Swiss cheese": layers 2, 4, and 5 had holes.

### Step 5: Identify the **top three** control gaps

Look for controls that:

1. Failed in more than one technique.
2. Are cheap to fix.
3. Would have changed the outcome the most.

These become your top three recommendations.

---

## 5.4 How to decide recommendations from the results

### Step 1: Rate each finding

Use a simple two-factor score.

| Likelihood (how easily could a real attacker repeat this?) | Score |
|------------------------------------------------------------|-------|
| Hard: needs rare skill or luck | 1 |
| Moderate: needs some skill | 2 |
| Easy: common techniques | 3 |

| Impact (what would happen if it succeeded for real?) | Score |
|------------------------------------------------------|-------|
| Low: limited exposure | 1 |
| Medium: account takeover or small financial loss | 2 |
| High: broad access, large loss, or regulatory harm | 3 |

**Risk = Likelihood x Impact**

| Score | Rating |
|-------|--------|
| 1 to 2 | Low |
| 3 to 4 | Medium |
| 6 | High |
| 9 | Critical |

### Step 2: Estimate effort and time

| Effort | Meaning |
|--------|---------|
| **Low** | Configuration change or a short procedure update (days) |
| **Medium** | New tool or process (weeks) |
| **High** | Project, budget, or structural change (months) |

### Step 3: Place each recommendation on the priority grid

| | Low effort | Medium effort | High effort |
|---|-----------|---------------|-------------|
| **Critical / High risk** | **Do now (within 30 days)** | **Plan now, start within 30 days** | **Fund and schedule (within 12 months)** |
| **Medium risk** | **Quick win (within 60 days)** | Schedule (within 6 months) | Review in the roadmap |
| **Low risk** | Do when convenient | Optional | Probably skip |

### Step 4: Apply these decision rules

1. **Fix controls, not people.** If your recommendation list is mostly "more training," you have not looked deeply enough.
2. **Prefer fixes that work even when a person makes a mistake.** (For example phishing-resistant MFA makes a stolen password far less useful.)
3. **Layer defenses.** One fix is not enough; aim for prevent, detect, and respond.
4. **Start with the cheapest high-impact fix.**
5. **Make each recommendation testable.** State how the client will know it worked.
6. **Name an owner.** "IT" is not an owner; "Head of IT Operations" is.
7. **Limit the list.** Give the top 10 to 15 recommendations; put the rest in an appendix.
8. **Show dependencies.** (For example "report button" must precede "SOC feedback loop.")
9. **Be honest about uncertainty.** If a finding came from a small sample, say so.
10. **Match recommendations to the client's size and budget.** A 40-person company will not build a 24/7 SOC; recommend managed options.

### Step 5: Offer options, not just orders

For major recommendations, offer two or three options with trade-offs.

| Option | Cost | Effort | Risk reduction | Notes |
|--------|------|--------|----------------|-------|
| A: Hardware security keys for all staff | High | High | Very high | Strongest; roll out to high-risk roles first |
| B: Passkeys on managed devices | Medium | Medium | High | Good balance |
| C: App-based MFA with number matching | Low | Low | Medium | Quick win; weaker than A and B |

---

## 5.5 The big chart: failure points and recommendations

How to use these tables:

1. Find the row that matches what your test found.
2. Copy the **recommendation** into your report, adapted to the client's tools and size.
3. Always add an **owner** and a **verification** method.
4. Use **Quick win** for the first 30 to 60 days and **Strategic fix** for 3 to 12 months.

### Table 1: Email and phishing

| Failure point observed | Likely root causes | Quick win (30 to 60 days) | Strategic fix (3 to 12 months) | Typical owner | How to verify |
|------------------------|--------------------|---------------------------|--------------------------------|---------------|---------------|
| Phish delivered to inboxes | Weak filtering; no detection of new or lookalike domains | Tighten gateway policies; block newly registered domains; enable link and attachment scanning | Add advanced email security with impersonation protection; domain monitoring | Head of IT / email admin | Repeat test without allow-listing; compare block rate |
| Spoofed messages with client's own domain were accepted | DMARC not at enforcement; SPF/DKIM gaps | Publish SPF and DKIM for all senders; set DMARC to monitor, then quarantine | Move DMARC to reject; inventory all sending services | Email admin | Run spoof test; check DMARC reports |
| Lookalike domain was not detected | No domain monitoring | Subscribe to a lookalike-domain monitoring service; register obvious variants | Automated takedown process | Security team | Test with a new lookalike domain; measure time to detect |
| No cue that the sender is external | No banner or tag | Add an external sender banner | Add VIP impersonation warnings | Email admin | Screenshot check in a retest |
| Links looked legitimate | No link rewriting or display of real destination | Enable safe-link rewriting; turn on URL reputation checks | Add browser isolation for risky links | Security team | Click test with a benign unknown domain |
| Users submitted credentials on a fake page | Passwords alone grant access; page was convincing | Roll out MFA everywhere; enable breached-password and risky-sign-in checks | Move to phishing-resistant MFA (security keys or passkeys) | IAM lead | Test accounts: relay attempt fails |
| Users opened risky attachments | Macros or risky file types allowed | Block macros from internet files; block dangerous attachment types | Application control; content disarm and reconstruction | Endpoint team | Controlled test file blocked |
| QR code lure succeeded | No guidance; mobile devices unprotected | Add QR code guidance in awareness messages; scan QR in email gateway if available | Mobile threat defense; company policy for QR use | Security / mobile team | Repeat with a benign QR code |
| OAuth or app-consent lure worked | Users can approve apps freely | Require admin approval for new app consent; review current consents | Publisher verification rules; periodic review | IAM / cloud admin | Test app is blocked pending approval |

### Table 2: Identity and MFA

| Failure point observed | Likely root causes | Quick win | Strategic fix | Owner | How to verify |
|------------------------|--------------------|-----------|---------------|-------|---------------|
| MFA approved after repeated prompts (fatigue) | Simple approve button | Enable number matching and additional context in prompts; limit prompts | Phishing-resistant MFA | IAM lead | Test account: repeated prompts not approved |
| MFA code relayed on a fake page | Codes can be re-entered elsewhere | Shorten code validity; add risk-based checks | Security keys or passkeys (bound to the real website) | IAM lead | Relay test on test accounts fails |
| MFA enrollment or reset by weak verification | Help desk can reset MFA easily | Require strong verification and manager approval for MFA changes | Self-service reset with strong identity proofing | IAM + help desk | Scripted test fails to reset |
| Session stolen and reused | Long session lifetime; no token binding | Shorten session lifetime for sensitive apps | Token binding or device-bound sessions | IAM / cloud admin | Replay test fails |
| Same password reused | No password manager; no breached-password checks | Provide a password manager; block known breached passwords | Move to passwordless | IAM lead | Audit shows reduced reuse |
| Legacy authentication allowed | Old protocols enabled | Disable legacy authentication | Conditional access policies | IAM / cloud admin | Test sign-in using legacy method is blocked |

### Table 3: Help desk and internal processes

| Failure point observed | Likely root causes | Quick win | Strategic fix | Owner | How to verify |
|------------------------|--------------------|-----------|---------------|-------|---------------|
| Agent reset password after a short call | Procedure unclear or not enforced | Write a one-page verification script with required steps; require a callback to the number on file | Automated identity verification; self-service with strong checks | Head of Service Desk | Repeat scripted calls; measure adherence rate |
| Agent revealed information (employee ID, manager name, etc.) | No list of what must never be shared | Publish a "never share by phone" list | Limit what the help desk can see for unverified callers | Service desk lead | Repeat calls; zero disclosure |
| Agent skipped steps under time pressure | Targets or pressure encourage speed | Remove speed pressure for sensitive requests; allow agents to say no without penalty | Review metrics so they reward correct verification | Service desk lead / HR | Interviews and repeat tests |
| High-risk requests had no second approval | Single person could approve | Add manager approval for resets of privileged accounts | Role-based approval workflow | IAM + service desk | Workflow log shows approvals |
| No logging or review of resets | Weak audit trail | Log every reset request with verification evidence | Automated anomaly alerts | Service desk lead | Audit sample |
| Outsourced help desk follows weaker rules | Vendor contract lacks equivalent controls | Align vendor procedure with client procedure | Contract clauses with audit rights | Vendor manager | Vendor retest |
| VIP requests handled differently | "Executive exception" | Remove exceptions; add extra checks | VIP verification protocol | CISO | Retest with VIP pretext (if approved) |

### Table 4: Vishing and impersonation

| Failure point observed | Likely root causes | Quick win | Strategic fix | Owner | How to verify |
|------------------------|--------------------|-----------|---------------|-------|---------------|
| Staff provided information to an unknown caller | No rule on what to share; caller ID trusted | Publish "we never ask for X by phone" rules; teach callback to a known number | Caller authentication tools; internal directory for callbacks | Security awareness lead | Repeat vishing; measure disclosure rate |
| Staff trusted caller ID | Spoofable | Remove trust in caller ID in policy | Verified calling for internal calls | Telecom / IT | Test fails |
| Staff did not report the suspicious call | No easy reporting channel | Add a phone and chat reporting channel; give simple script | Integrated reporting across channels | Security team | Report rate improves |
| Staff did not know how to verify a caller | No script | One-page "verify the caller" card | Training with role-play | Awareness lead | Retest |

### Table 5: Finance and business email compromise

| Failure point observed | Likely root causes | Quick win | Strategic fix | Owner | How to verify |
|------------------------|--------------------|-----------|---------------|-------|---------------|
| Payment-change request accepted | No independent verification of changes | Require callback to a known number for bank detail changes; never use contact details from the request | Vendor master data control and dual approval | CFO / finance controller | Simulated request is blocked |
| Executive impersonation accepted | Authority bias; no checks for urgent requests | Add rule: urgent payment requests need second confirmation | Payment approval workflow with limits and dual control | CFO | Simulation blocked at approval step |
| New vendor created without check | Weak onboarding | Add vendor verification steps | Vendor risk management | Procurement | Audit |
| Invoice attached to email was trusted | Process relies on email | Verify invoices against purchase orders | Automated three-way matching | Finance operations | Test invoice flagged |

### Table 6: Reporting, detection, and response

| Failure point observed | Likely root causes | Quick win | Strategic fix | Owner | How to verify |
|------------------------|--------------------|-----------|---------------|-------|---------------|
| Low report rate | No visible button; unclear channel; fear of blame | Add a one-click report button; announce a no-blame policy | Reporting integrated with SOC and feedback loop | CISO | Report rate in retest |
| Reports not acted on quickly | SOC overloaded; no SLA | Set a triage SLA; assign an owner for reported phish | Automated triage and purge | SOC lead | Time-to-purge metric |
| No feedback to reporters | Process gap | Send a thank-you and outcome to every reporter | Automated feedback | SOC lead | Survey and reporting trend |
| SOC did not detect submissions | No detection rules | Add alerts for risky sign-ins, new mailbox rules, new OAuth apps | UEBA / identity threat detection | SOC lead | Retest: alerts fire |
| Slow containment | Manual purge | Document and practice a purge procedure | Automated remediation | SOC lead | Time to containment |
| No playbook for vishing or impersonation | Only email is covered | Write playbooks for calls and impersonation | Regular exercises | IR lead | Tabletop exercise |
| Deconfliction not smooth | No contact process | Define a contact for authorized tests | Standard test-notification process | SOC lead | Next test |

### Table 7: Physical security

| Failure point observed | Likely root causes | Quick win | Strategic fix | Owner | How to verify |
|------------------------|--------------------|-----------|---------------|-------|---------------|
| Tailgating succeeded | Culture of holding doors; no barrier | Signage and a polite "badge please" script; guard coaching | Turnstiles or mantraps in high-risk entrances | Facilities | Repeat entry test |
| Visitor entered without escort | Weak visitor process | Enforce sign-in and escort; issue visible visitor badges | Visitor management system | Facilities | Observation and test |
| No challenge to unknown person | People afraid to challenge; no policy | Adopt a no-blame challenge policy and script | Training and drills | Facilities / HR | Challenge rate improves |
| Delivery or contractor pretext worked | No verification of vendors | Verify appointments with the requesting department | Pre-registered visitor and contractor system | Facilities | Repeat test |
| Unlocked workstations and visible information | No auto-lock; no clean-desk habits | Enforce auto-lock; run clean-desk spot checks | Endpoint policy enforcement; badge-activated locking | IT / Facilities | Walkthrough audit |
| Sensitive paper in bins | No shredding policy | Provide shred bins and a policy | Secure disposal service | Facilities | Bin audit |
| Server or network rooms accessible | Weak door controls | Add access logging and locks; limit keys | Badge access with alarms | Facilities / IT | Test entry blocked |
| USB or media plugged in | No device control | Block or restrict USB mass storage on endpoints; policy "give found devices to security" | Endpoint device control and awareness | Endpoint team | Test device blocked |
| Door propped open or held | Convenience | Door alarm for propped doors | Camera analytics | Facilities | Alarm triggers |
| Cameras did not detect | Coverage gaps | Add or reposition cameras at key entrances | Monitored video and incident process | Security | Test is observed |

### Table 8: Information exposure (OSINT) and third parties

| Failure point observed | Likely root causes | Quick win | Strategic fix | Owner | How to verify |
|------------------------|--------------------|-----------|---------------|-------|---------------|
| Technology stack visible in job postings | Detailed postings | Remove specific product names where possible | Review postings with security | HR / security | Re-run OSINT |
| Org chart and roles easy to find | Public profiles and websites | Reduce unnecessary detail on public pages | Social media guidance (not policing) | Communications | Re-run OSINT |
| Badge or building details on public photos | Photos posted | Add guidance on posting workplace photos | Redesign badges to reduce copying risk | Communications / Facilities | Re-run OSINT |
| Out-of-office messages reveal details | Overly informative replies | Provide a standard out-of-office template | Mail policy | IT | Sample check |
| Vendors accepted unverified requests | Vendor procedures weak | Contract clause requiring verification | Vendor audit rights | Vendor manager | Vendor test |
| Mobile and SMS lures worked | No mobile protection | Policy that company never asks for credentials by SMS; staff reporting channel | Mobile threat defense; MDM controls | Mobile / IT | Repeat smishing test |

### Table 9: Culture, awareness, and governance

| Failure point observed | Likely root causes | Quick win | Strategic fix | Owner | How to verify |
|------------------------|--------------------|-----------|---------------|-------|---------------|
| Staff afraid to report | Fear of blame | Written no-blame policy; thank reporters publicly (in aggregate) | Embed in values and manager training | CISO / HR | Anonymous survey; report rate |
| Training is generic and forgotten | Annual long course | Short monthly micro-lessons; role-based examples | Training tied to real events and simulations | Awareness lead | Report rate and behavior metrics |
| Managers punished clickers | Misunderstanding | Brief managers on the no-discipline rule | Policy in HR handbook | HR | Employee survey |
| Executives exempt from controls | Culture | Apply same controls to executives | Executive sponsorship | CEO / CISO | Audit |
| No ongoing testing | One-off exercise | Plan a retest in 3 to 6 months | Annual program with different scenarios | CISO | Calendar and results |
| Security team unaware of exercises | Poor planning | Notification process | Joint planning with SOC | CISO | Next exercise |

---

## 5.6 Positive findings to include (and how to word them)

A report that lists only failures teaches the client to fear the next test. It also hides what is **worth protecting**. Include real, specific positives.

### Rule: positives must be true, specific, and measured

Never add fake praise to soften bad news. Use data.

### Checklist of positives to look for

| Look for | Example of how to write it |
|----------|----------------------------|
| People who reported | "18 recipients (15 percent) reported the message, more than the 12 percent who submitted." |
| Speed of reporting | "The first report arrived 7 minutes after delivery." |
| People who challenged | "Two of six physical attempts were challenged by staff within minutes, using polite and effective language." |
| Correct help desk behavior | "Four of ten calls ended with the agent insisting on callback verification." |
| Filters that worked | "The gateway blocked 22 percent of messages in the unassisted wave." |
| Email authentication | "Spoofed messages using the client's own domain were rejected." |
| MFA that held | "The MFA method withstood relay attempts on test accounts." |
| Fast containment | "The SOC removed messages from 61 mailboxes within 24 minutes of the first report." |
| Good procedures that worked when followed | "Where the callback step was used, no resets occurred." |
| Improvement against baseline | "Report rate rose from 6 percent to 15 percent since last year." |
| Good communication | "Several staff forwarded the message to colleagues with a warning." |
| Culture | "Staff thanked security for the quick response in the follow-up survey." |
| Safe handling of found media | "Of the placed devices, X were handed to security unplugged." |
| Physical barriers | "Turnstiles prevented tailgating at the main entrance." |
| Staff declined | "Most recipients did nothing, and some replied to ask whether the message was real." |

### How to place positives in the report

1. Put **"What worked well"** right after the executive summary.
2. In each finding, include a **"What worked"** line.
3. In the debrief meeting, start with positives (and mean them).
4. Recommend **protecting and measuring** the positives (for example "keep the report button visible; track report rates monthly").

### Wording examples

| Weak or vague | Strong and specific |
|---------------|---------------------|
| "Staff did well." | "Fifteen percent of recipients reported the message within the first hour; the first report took 7 minutes." |
| "Good filters." | "The gateway blocked 28 of 125 messages in the wave without allow-listing, including all messages carrying macros." |
| "Security-aware culture." | "Staff in three offices challenged an unbadged visitor within five minutes using the recommended polite phrase." |

---

## 5.7 Putting the report together

### Report structure

| # | Section | Contents | Length |
|---|---------|----------|--------|
| 1 | **Executive summary** | Goal, headline results (3 to 5 numbers), top 3 control gaps, what worked, top recommendations | 1 page |
| 2 | **Scope and method** | Authorization, scope, dates, techniques, pool, ethics rules, data handling | 1 to 2 pages |
| 3 | **What worked well** | Positive findings with numbers | 1 page |
| 4 | **Results at a glance** | Summary table by technique, funnel, time-to-event table | 1 to 2 pages |
| 5 | **Timeline** | One-page overview | 1 page |
| 6 | **Findings** | Each finding in the structure from Section 4 | 1 page each |
| 7 | **Control effectiveness matrix** | From 5.3 | 1 to 2 pages |
| 8 | **Recommendations roadmap** | Prioritized grid, owners, time horizons | 2 pages |
| 9 | **Deviations and limitations** | Mistakes, accidental captures, data issues, small samples | 1 page |
| 10 | **Appendices** | Detailed timeline, pretext cards, metrics method, glossary, evidence index | As needed |

### Executive summary template

```markdown
## Executive summary

**Why we tested:** [one sentence goal]
**What we tested:** [techniques, pool size, dates]
**Headline results (aggregate):**
- [Metric 1, with context]
- [Metric 2, with context]
- [Metric 3, with context]

**What worked well:**
- [Positive 1]
- [Positive 2]

**The top three control gaps:**
1. [Control gap, plain language]
2. [Control gap, plain language]
3. [Control gap, plain language]

**Top recommendations (next 30 days):**
1. [Action, owner]
2. [Action, owner]
3. [Action, owner]

**Important note:** No individual results are included. Findings concern systems and processes.
```

### Handling the report safely

1. Mark it **Confidential** and limit distribution to named people.
2. Deliver through a secure channel (not a normal email attachment).
3. Do not include real credentials, personal data, or photos of people.
4. Keep an evidence index with pseudonymous references.
5. Agree retention and destruction dates in advance.
6. Deliver the draft first and ask the sponsor to **fact-check** (dates, scope, and wording), not to change conclusions.

### The debrief meeting (agenda)

| Time | Item |
|------|------|
| 5 min | Thank you; restate goals and ground rules (no names, no blame) |
| 10 min | What worked well |
| 15 min | Timeline and results at a glance |
| 20 min | Top control gaps (one slide each: what we saw, why it worked, what would make it hard) |
| 15 min | Recommendations: quick wins, then strategic |
| 10 min | Decisions: who owns what, dates |
| 5 min | Retest plan and next steps |

**Tips for the debrief:**

1. Start with positives.
2. Use the wording table. If someone says "who clicked?", use the Section 4.5 script.
3. Leave time for questions.
4. End with decisions written on the screen.

### The retest plan

1. Agree on a retest date (usually 3 to 6 months after fixes).
2. Use a **different** pretext so you measure resilience, not memory.
3. Keep the metrics and denominators the same.
4. Compare results with the baseline and show progress.
5. Include the positive trend in the next report.

---

## Section 5 checklist

- [ ] Activity log complete, in UTC, with evidence links
- [ ] Client-side events merged (with permission)
- [ ] Timeline prepared (overview and detail)
- [ ] Outcome categories defined; denominators stated; scanner clicks removed
- [ ] Rates shown with group size, margin of error, difficulty, and allow-listing note
- [ ] Control effectiveness matrix filled in
- [ ] Each finding rated and written as a control gap
- [ ] Recommendations prioritized (grid), with owner and verification
- [ ] Positive findings included, specific and measured
- [ ] Deviations and limitations stated honestly
- [ ] Report delivered securely; distribution list agreed
- [ ] Debrief held; decisions recorded; retest scheduled

---

# SECTION 6: SETTING UP DEFENSIVE MEASURES FOR CLIENTS

**Goal of this section:** by the end, you can turn the report into a working defense plan: which controls to put in place, in what order, how to set each one up, and how to prove it works.

---

## 6.0 Before you start: three ground rules

### Rule 1: Keep roles clear (independence)

If your company both **tests** the client and **builds the defenses**, say so openly.

| Situation | What to do |
|-----------|------------|
| You only advise | Deliver the roadmap and offer to review progress |
| You help implement | Use a separate statement of work and tell the client you will not "grade your own homework" without disclosure |
| You do the retest | Use a different team member or a third party for the retest if possible; if not, disclose it in the report |

> **TRAP:** The temptation to find problems you can then fix (and bill for). Keep recommendations based on evidence, not on your product catalog.

### Rule 2: Turn findings into a backlog with owners

1. Copy the recommendation table from the report into a tracker (spreadsheet or ticket system).
2. For each row add: **owner**, **due date**, **status**, **how to verify**.
3. Review the tracker every two weeks with the engagement owner.
4. Report progress to the sponsor monthly.

### Rule 3: Fix the cheap, high-impact things first

Do **quick wins** before **projects**. A configuration change this week beats a perfect solution next year.

---

## 6.1 The 30 / 90 / 365-day roadmap

| Horizon | Theme | Typical actions |
|---------|-------|-----------------|
| **First 30 days** | Stop the bleeding | External sender banner; report button; no-blame announcement; help desk verification script; callback rule for payment changes; block risky attachments; disable legacy authentication; admin approval for app consent |
| **Days 31 to 90** | Build the foundations | DMARC to quarantine; MFA everywhere; number matching; SOC triage playbook and SLA; visitor and escort procedure; USB restrictions; role-based micro-training; lookalike-domain monitoring |
| **Days 91 to 365** | Make it durable | Phishing-resistant MFA for admins and high-risk roles; DMARC reject; automated phish purge; turnstiles or mantraps where justified; vendor contract clauses; regular exercises; scheduled retest |

---

## 6.2 Email defenses (step by step)

### 6.2.1 Email authentication (SPF, DKIM, DMARC)

Think of these as three stamps that prove an email really came from the company.

1. **List every system that sends email as your domain** (mail platform, marketing tools, ticketing, HR systems, payroll, newsletters). Missing one will break real mail later.
2. **SPF:** publish a record that lists the allowed senders. Keep it within the technical limits (a long list can break it).
3. **DKIM:** turn on signing for each sending system.
4. **DMARC at monitor mode:** publish a DMARC policy that only **reports** (no blocking yet). Send reports to a mailbox or a reporting service.
5. **Read the reports for 2 to 4 weeks.** Fix legitimate senders that fail.
6. **Move to quarantine** (start with a small percentage if your platform supports it, then increase).
7. **Move to reject** when the reports are clean.
8. **Protect parked or unused domains** with a "no mail from here" policy.
9. **Repeat the spoof test** to prove it works.

### 6.2.2 Make the email client help people

| Control | How to set it up | Why |
|---------|------------------|-----|
| External sender banner | Mail admin adds a tag to messages from outside | Gives a visible cue |
| Impersonation protection | List executives and finance roles; enable "display name impersonation" detection | Stops "CEO" lookalikes |
| Safe-link rewriting and click-time checks | Turn on in the gateway | Blocks links that turn malicious later |
| Attachment scanning and sandboxing | Turn on in the gateway | Catches risky files |
| Block macros from internet files | Group policy or management tool | Removes a common risk |
| Block or warn on risky file types | Gateway rules | Reduces exposure |
| Disable automatic external forwarding | Mail policy | Stops silent data leakage |
| New mailbox rule alerts | SOC alert when rules forward or delete mail | Detects account takeover |
| QR code scanning | Enable if the gateway supports it | Reduces QR lures |

### 6.2.3 Lookalike-domain monitoring

1. Register obvious variants of the company's main domains (common misspellings and endings).
2. Subscribe to a service (or use a built-in feature) that alerts on new lookalike registrations.
3. Write a **takedown procedure**: who investigates, who contacts the registrar or host, how long it should take.
4. Feed new bad domains to the email gateway blocklist.

### 6.2.4 The report button

1. Install a **one-click report button** in the email client (all platforms: desktop, web, mobile).
2. Route reports to the SOC queue or a shared mailbox with a clear owner.
3. Auto-reply with "Thank you, we are looking at it."
4. Track: reports per month, time to triage, time to purge.
5. Announce it in all-staff communications and in the new-hire pack.

---

## 6.3 Identity and MFA defenses (step by step)

### Step-by-step

1. **Inventory every sign-in path:** email, VPN, cloud apps, admin consoles, remote desktop, partner portals.
2. **Turn on MFA for everything.** Start with administrators, finance, HR, and executives, then everyone.
3. **Upgrade the MFA method in stages:**

| Stage | Method | Notes |
|-------|--------|-------|
| 1 | App-based push with **number matching** and context | Stops blind approvals; still phishable |
| 2 | Security keys or passkeys for **admins and high-risk roles** | Strongest against fake-page relaying |
| 3 | Passkeys or keys for all staff on managed devices | The long-term goal |
| (avoid) | SMS or voice codes for admins and high-risk roles | Easier to intercept or hijack |

4. **Disable legacy authentication** (old protocols that skip MFA).
5. **Set conditional access rules:** require MFA and a compliant device for sensitive apps; block risky sign-ins (impossible travel, unfamiliar locations).
6. **Control app consent:** require admin approval for new third-party app permissions; review existing ones.
7. **Shorten session lifetimes** for sensitive apps; require re-authentication for high-risk actions.
8. **Provide a password manager** to all staff; block known breached passwords.
9. **Create break-glass accounts:** two emergency admin accounts with strong, offline-stored credentials, excluded from conditional access but heavily monitored.
10. **Test with test accounts** to confirm the controls hold (Section 4.7, Method C).

---

## 6.4 Help desk and identity-verification procedure

### Why this matters

Attackers go after the help desk because a single reset can unlock an account and bypass MFA. The help desk must be **helpful and hard to trick at the same time**.

### Step-by-step: design the procedure

1. **List the requests** the help desk gets (password reset, MFA reset, unlock, new device, change of email, access request).
2. **Rate each request by risk:**

| Request | Risk | Required verification |
|---------|------|-----------------------|
| General information | Low | None beyond basic |
| Password reset (standard user) | Medium | Verify identity (see below) |
| MFA reset or new MFA device | High | Strong verification **and** manager or approver confirmation |
| Privileged account (admin, finance, executive) | Very high | In-person or video verification plus approver; no phone-only resets |
| Change of contact details or bank details | High | Callback and written approval |

3. **Choose verification methods that cannot be answered from public sources.** Avoid questions based on public information (birthday, job title, manager name).
4. **Use callback verification:** the help desk calls back **a number already on file** (not a number the caller provides).
5. **For high-risk requests, require a second person** (the caller's manager or the approver) to confirm through a separate channel.
6. **Add a cooling-off notification:** when an MFA reset happens, send a notice to the user's other registered channels (existing device, alternate address).
7. **Write it as a one-page script** with required steps in order. No step may be skipped.
8. **Give agents explicit permission to refuse** and a script for saying no.
9. **Log every request** with the verification steps taken.
10. **Do not allow executives or VIPs to bypass** the steps.

### Sample one-page procedure

```markdown
HELP DESK IDENTITY VERIFICATION (HIGH-RISK REQUEST)

1. Do not act on an inbound call alone.
2. Record the caller's claimed name, request, and time.
3. Hang up. Call back using the number in the HR directory (not the one given).
4. Verify using: [method that is not publicly available].
5. For MFA reset or privileged accounts: send an approval request to the
   user's manager in the ticketing tool and wait for approval.
6. Complete the reset. Send a notification to all registered contact methods.
7. If any step fails: stop, say the refusal script, create a security ticket.
8. Log all steps in the ticket.

Refusal script:
"I'm sorry, I can't complete this without finishing verification. This
protects your account. I'll call you back on the number we have on file."
```

### Handling pressure (script for agents)

| Pressure from the caller | What the agent says |
|--------------------------|---------------------|
| "I'm in a hurry. The CEO is waiting." | "I understand. Verification takes two minutes and protects you. I'll call you right back on the number on file." |
| "You're just making this hard." | "I know it's frustrating. This protects your account, and I'm not allowed to skip it for anyone." |
| "I'll report you to your manager." | "That's fine. My manager will support me following the procedure." |
| "Just this once." | "I can't make exceptions. I can escalate if there is a real emergency." |

### Outsourced help desks

1. Add the same procedure to the vendor contract as a requirement.
2. Include audit rights and periodic authorized testing.
3. Ask for monthly reporting on verification adherence.
4. Make sure the vendor's agents can say no without penalty.

### Metrics

- Verification adherence rate (from Section 3 checklist)
- Number of MFA resets per month (spikes need review)
- Percent of high-risk requests with a second approval
- Time to complete high-risk requests (balance security and service)

---

## 6.5 Finance and payment verification

### Step-by-step

1. **Write the rule:** "Any change to bank details, any new vendor, and any urgent payment request must be verified by calling a **known** number, never a number in the request."
2. **Keep a trusted vendor directory** with verified phone numbers, maintained by finance.
3. **Require two people** (dual control) for payments above a threshold and for any bank detail change.
4. **Add a waiting period** for bank detail changes before the first payment.
5. **Match invoices to purchase orders** and delivery confirmations.
6. **Create an executive-request rule:** requests from executives for urgent payments need a second confirmation through a different channel.
7. **Train finance staff with real scenarios** (a payment-change request, a fake executive message).
8. **Run a controlled simulation** (with the finance lead in the circle) and confirm the process holds.
9. **Define the response plan** if a bad payment is sent (call the bank immediately; the first hours matter).

---

## 6.6 Reporting, SOC, and response

### Step-by-step: triage playbook for a reported suspicious message

1. **Receive** the report (button or mailbox).
2. **Acknowledge** to the reporter (automatic).
3. **Classify:** malicious, suspicious, or harmless. Check the sender, links, attachments, and similar messages.
4. **Search** the mail system for the same message in other mailboxes.
5. **Contain:** purge from all mailboxes, block the sender and the link domain.
6. **Check for clicks or submissions:** review sign-in logs, new mailbox rules, new app consents.
7. **If credentials may be exposed:** reset passwords, revoke sessions, review MFA devices.
8. **Close the loop:** send the reporter an outcome and thanks.
9. **Record** timestamps for metrics.
10. **Learn:** if a new technique appeared, update filters and training.

### SLA suggestions (adapt to size)

| Stage | Target |
|-------|--------|
| Acknowledge | Immediate (automatic) |
| First triage | Within 30 minutes during business hours |
| Purge from all mailboxes | Within 1 hour of confirmation |
| Reporter feedback | Within 24 hours |

### Detections to add

| Detection | Why |
|-----------|-----|
| Sign-in from a new country or impossible travel | Possible stolen credentials |
| Multiple MFA prompts in a short time | Prompt fatigue attack |
| New mailbox forwarding or hiding rules | Account takeover |
| New third-party app consent | Consent phishing |
| Many users visiting the same new domain | Phishing campaign |
| Help desk MFA reset followed by sign-in from a new place | Social engineering of the help desk |

### Playbooks beyond email

Write short playbooks for: suspicious phone call, SMS lure, impersonation of an executive, suspicious visitor, found USB device.

### Exercises

- **Tabletop exercise** every 6 months: a facilitator describes a scenario ("the CFO's assistant gets a call...") and the team talks through the response.
- **Purple-team style review:** testers and defenders sit together, replay the test, and improve detections.

---

## 6.7 Awareness that works

### What to change from "annual boring training"

| Old way | Better way |
|---------|-----------|
| One long annual course | Short, frequent lessons (a few minutes monthly) |
| Same content for everyone | Role-based examples (finance, HR, help desk, executives, new hires) |
| Training after failure as punishment | Training offered as help, never as discipline |
| Measuring click rate only | Measuring **report rate**, time to report, and confidence |
| Generic warnings | Real examples from the company's own environment (with permission) |
| Shaming | Thanking reporters |

### Program design (step by step)

1. **Pick three behaviors** to teach first: pause on urgency, verify through a second channel, report.
2. **Create simple memory aids**: for example "Stop, Check, Report."
3. **Offer new-hire training** in the first week.
4. **Provide role-based modules** for high-risk teams.
5. **Run simulated phishing responsibly**: written ethics policy, no sensitive topics, no names to managers, difficulty rating, and communication in advance (a general notice that simulations happen).
6. **Make reporting easy and rewarded** (thank-you messages, recognition in aggregate).
7. **Track the right numbers:** report rate, time to first report, repeat rates (anonymous), survey of confidence.
8. **Review quarterly** and adjust.

### Ethics policy for future simulations (one page the client can adopt)

1. Simulations do not use pay, layoffs, bonuses, health, or personal crises.
2. Individual results are never shared with managers or used in reviews.
3. Staff are told that simulations occur, in general terms.
4. Anyone who feels distressed can contact [name].
5. Simulations are approved by security, HR, and legal.

---

## 6.8 Physical security defenses

### Step-by-step

1. **Map every entry point:** main entrances, side doors, loading docks, parking, stairwells, roof access.
2. **Classify zones:** public, general staff, restricted (finance, server rooms, executive floor).
3. **Add barriers where risk is highest:** turnstiles or mantraps at main entrances; badge readers on restricted doors.
4. **Set a visitor procedure:** sign-in, ID check, visible visitor badge, escort at all times; pre-register contractors and deliveries.
5. **Make challenging polite and safe:** adopt a no-blame policy and a script:

> "Hi, can I help you? I don't think I've seen your badge. Let me walk you to reception."

6. **Coach guards and receptionists** on the script and on escalation steps.
7. **Badge rules:** badges visible inside, never lent, reported if lost; badge design should not reveal building details to outsiders.
8. **Anti-tailgating:** door alarms for propped doors; signage ("one badge, one person"); culture of polite refusal.
9. **Clean desk and screen lock:** automatic lock after a short idle time; periodic spot checks (friendly, not punitive).
10. **Secure disposal:** shred bins and a pickup schedule; no sensitive paper in normal bins.
11. **Server and network rooms:** locks, access logs, alarms, camera coverage.
12. **Cameras:** cover entrances and restricted doors; make sure someone can watch or review them.
13. **Found-device rule:** any found USB or device goes to security, unplugged and untouched.
14. **USB control:** block or restrict removable storage on endpoints.
15. **Contractors and cleaners:** vetting, badges, supervised access to sensitive areas.

### Physical metrics to track

- Entry attempts blocked, challenged, or reported (from tests)
- Percentage of doors with working alarms
- Visitor sign-in compliance (sample audits)
- Time for security to respond to a report

---

## 6.9 Vishing and impersonation policy

### The "we never" list (publish it)

| The company will never... | So if someone asks, it is suspicious |
|---------------------------|--------------------------------------|
| Ask for your password or MFA code by phone, SMS, or email | Always |
| Ask you to approve an MFA prompt you did not start | Always |
| Ask you to install software from a link in a message | Always |
| Ask you to change bank details based on a single message | Always |
| Ask you to bypass a procedure because "the executive said so" | Always |
| Ask for gift cards or urgent transfers | Always |

### The "call back" card (one line on a sticker)

> "If in doubt, hang up and call back using the number in the company directory."

### Policy points

1. Staff may refuse any request that feels wrong, **without penalty**.
2. Executives agree publicly that their urgent requests may be verified.
3. Everyone knows how to report suspicious calls.
4. Reception and assistants have a script for unknown callers.

---

## 6.10 Information hygiene (reducing what attackers can learn)

1. **Job postings:** describe skills, not the exact internal product stack, where possible.
2. **Public staff pages:** show what customers need, not full org charts.
3. **Social media guidance:** a friendly guide (not a policing policy) about posting badge photos, building layouts, and project details.
4. **Out-of-office messages:** a standard template with no unnecessary details.
5. **Document metadata:** remove author names and internal paths from public documents.
6. **Visitor and delivery information:** do not publish internal procedures online.
7. **Re-run passive OSINT** after changes and compare.

---

## 6.11 Governance and culture

| Practice | How to set it up |
|----------|------------------|
| **Executive sponsor** | One named executive owns the program and removes blockers |
| **No-blame reporting policy** | Written, signed by the CEO, included in onboarding |
| **Metrics dashboard** | Monthly: report rate, time to triage, time to purge, MFA coverage, help desk adherence, physical findings |
| **Regular exercises** | Annual full test, plus smaller quarterly targeted checks with different pretexts |
| **Third-party contracts** | Add verification and testing clauses for vendors that handle identity, payments, or reception |
| **Incident communications** | Prepared templates for staff, customers, and regulators |
| **Review cycle** | Revisit the roadmap every quarter |
| **Celebrate reporters** | In aggregate, thank staff for reports and share what the report prevented |

---

## 6.12 How to prove the defenses work (verification)

### Step-by-step

1. For each defensive measure, write **one testable statement** (for example "A spoofed email using our domain is rejected").
2. Choose the verification method:

| Measure | Verification method |
|---------|---------------------|
| DMARC enforcement | Send controlled spoof tests; check DMARC reports |
| External sender banner | Screenshot check |
| Report button | Test report; measure acknowledgement and purge time |
| MFA upgrade | Test accounts: relay and fatigue attempts fail |
| Help desk procedure | Scripted calls to test accounts; adherence rate |
| Payment verification | Controlled simulation (finance lead aware) |
| Turnstile and escort rules | Authorized physical test |
| USB control | Test device (harmless marker) is blocked |
| Training | Report rate and time to first report in the next simulation |

3. **Schedule the retest** (3 to 6 months), using **new pretexts** and **the same metrics**.
4. **Compare** against the baseline and report progress.
5. **Close findings** only when the verification passes.
6. **Reopen** a finding if a retest shows regression.

---

## 6.13 Who does what (role checklist)

| Role | Key actions |
|------|-------------|
| **CEO / executive sponsor** | Sign the no-blame policy; fund the roadmap; accept verification of urgent requests |
| **CISO / security lead** | Own the backlog; run the metrics; schedule exercises |
| **IT / email admin** | Banner, DMARC, gateway rules, legacy auth, app consent |
| **IAM lead** | MFA upgrade, conditional access, break-glass accounts |
| **SOC lead** | Triage playbook, SLAs, detections, feedback to reporters |
| **Help desk lead** | Verification procedure, scripts, logging, vendor alignment |
| **Finance lead** | Callback rule, dual control, vendor directory |
| **HR** | No-discipline policy, onboarding training, wellbeing contact |
| **Facilities / physical security** | Entry barriers, visitor process, guard coaching, cameras, disposal |
| **Communications** | Staff debrief, templates, social media guidance |
| **Legal / compliance** | Policies, vendor clauses, privacy for simulations |
| **Managers** | Support reporting; never ask "who clicked?" |
| **All staff** | Pause, verify, report |

---

## Section 6 checklist

- [ ] Independence and roles clarified in writing
- [ ] Backlog created with owners, dates, and verification methods
- [ ] 30-day quick wins scheduled
- [ ] Email authentication plan (SPF, DKIM, DMARC stages) in motion
- [ ] External sender banner and impersonation protection enabled
- [ ] Report button live and routed to a named owner
- [ ] MFA upgrade plan and legacy authentication removal scheduled
- [ ] Help desk procedure published; agents trained; vendor aligned
- [ ] Finance callback and dual-control rules in place
- [ ] SOC triage playbook, SLAs, and detections defined
- [ ] Awareness program redesigned around reporting
- [ ] Physical improvements prioritized (visitors, escort, challenge script, USB, disposal)
- [ ] "We never" list published
- [ ] Governance (sponsor, dashboard, exercises) established
- [ ] Retest scheduled with new pretexts and the same metrics

---

# APPENDICES

---

## Appendix A: Master go / no-go checklist (one page for the whole engagement)

### Before signing
- [ ] Four Ownership Questions answered
- [ ] Signer verified (registry, authority proof, call-back, domain, signature trail)
- [ ] Insurance checked
- [ ] Legal review done for every country involved

### Before testing
- [ ] MSA, SOW, and RoE signed (with no-discipline, real-incident, and law-enforcement clauses)
- [ ] Authorization letter signed (property owner counter-signed where needed)
- [ ] Emergency contacts tested; stop word and verification word set
- [ ] Third-Party Register complete; consents filed
- [ ] Works council or employee-representative requirements met
- [ ] Target pool locked; exclusions applied; pseudonymous IDs assigned
- [ ] Need-to-know circles recorded; blind or informed SOC decided
- [ ] Goal statement and metrics plan written
- [ ] OSINT questions and source list approved
- [ ] Infrastructure built, authenticated, and tested; kill switch tested
- [ ] Landing page canary test passed (zero matches)
- [ ] Pretext cards approved in writing
- [ ] Dry run complete
- [ ] Go/no-go meeting held and recorded

### During testing
- [ ] Activity log kept in UTC with evidence links
- [ ] Daily check-in held
- [ ] Stop conditions understood by every tester
- [ ] Deviations recorded immediately

### After testing
- [ ] Teardown completed and confirmed to the sponsor
- [ ] Data cleaned; scanner clicks removed
- [ ] Report written using the wording table (no names, no blame)
- [ ] Positive findings included
- [ ] Report delivered securely; debrief held
- [ ] Employee debrief message sent
- [ ] Data deleted on schedule; deletion confirmed in writing
- [ ] Retest scheduled

---

## Appendix B: Templates

### B.1 Authorization letter

```markdown
[CLIENT LETTERHEAD]

Date: [date]

TO WHOM IT MAY CONCERN

Subject: Authorization for security assessment

[Client legal name] ("the Client") has engaged [Tester company] to carry out
an authorized physical and social engineering security assessment.

The following persons are authorized to carry out this assessment:
  1. [Full name]   (photo ID presented on request)
  2. [Full name]   (photo ID presented on request)

Locations: [full addresses, floors or areas]
Dates and times: [start] to [end], [time zone]

The activities include attempts to enter the premises and interact with staff
as part of a controlled test. They will not damage property, force entry, or
enter unsafe areas.

If these persons are questioned, please contact either of the following
before taking further action. Both are available 24 hours a day during the
testing window:

  Primary contact:   [Name, title, mobile]
  Secondary contact: [Name, title, mobile]

Verification: the contact will ask for a verification word known only to the
testers and the contact.

Signed for the Client:
[Name, title, signature, date]

Premises owner or manager (if different):
[Name, title, organization, signature, date]
```

### B.2 Third-Party Register

| # | Third party | Role | What may be affected | Contact (name, role) | Who contacted them | Date requested | Consent received (Y/N, date) | File reference | Deconfliction decision |
|---|-------------|------|----------------------|----------------------|--------------------|----------------|------------------------------|----------------|------------------------|
| 1 | | | | | | | | | |
| 2 | | | | | | | | | |
| 3 | | | | | | | | | |

### B.3 Deviation report (use for any mistake, accident, or unplanned event)

```markdown
DEVIATION REPORT   ID: D-[number]

Date/time (UTC):
Reported by:
Engagement:

1. What happened (facts only):
2. Which rule, scope item, or stop condition was involved:
3. Who or what was affected (pseudonyms only):
4. Was any personal data or credential captured? (Y/N, describe):
5. Immediate action taken:
6. Who was notified and when (sponsor, engagement owner, legal):
7. Root cause (process, not person):
8. Corrective action and owner:
9. Does this change the RoE or pretext cards? (Y/N):
10. Included in final report? (Y/N):
```

### B.4 Pretext approval log

| Card ID | Version | Channel | Target group | Risk | Written by | Second reviewer | HR/legal review | Sponsor approval (date) | Status |
|---------|---------|---------|--------------|------|------------|-----------------|-----------------|-------------------------|--------|
| P-01 | 1 | | | | | | | | |
| P-02 | 1 | | | | | | | | |

### B.5 Daily check-in note

```markdown
DAILY CHECK-IN   Date: [date]

Activities today (counts only):
Stop conditions triggered: [none / describe]
Deviations: [none / IDs]
Reports or reactions from staff (aggregate): 
Issues for the sponsor:
Plan for tomorrow:
Changes requested (scope): [none / describe]
```

### B.6 Finding write-up

```markdown
### Finding F-[n]: [Control gap, stated as what is missing or weak]

**Risk rating:** [Low / Medium / High / Critical]   (Likelihood [1-3] x Impact [1-3])
**What we observed (aggregate):**
**Controls that had a chance to stop it:**
**Why it worked:**
**What made it hard for people to detect:**
**What worked well (if anything):**
**Recommendation (quick win):**
**Recommendation (strategic):**
**Owner:**
**How to verify the fix:**
```

### B.7 Data deletion confirmation

```markdown
DATA DELETION CONFIRMATION

Engagement: [name]
Date of deletion: [date]
Data deleted: [pseudonym key, OSINT notes, logs, evidence, test mailboxes, domains]
Method: [secure deletion method]
Performed by: [name]
Verified by: [name]
Backups included: [yes / no, explain]

Signed: [name, date]
```

---

## Appendix C: Glossary (plain words)

| Term | Meaning |
|------|---------|
| **Social engineering** | Tricking people (not computers) into doing something unsafe, such as clicking a link or opening a door |
| **Pentest (penetration test)** | An authorized attempt to find weaknesses by acting like an attacker |
| **RoE (Rules of Engagement)** | The written rulebook for the test |
| **SOW / MSA** | Statement of Work / Master Services Agreement; the commercial contracts |
| **Authorization letter** | A signed letter that proves the testers are allowed to be there |
| **Pretext** | The story used to make a request believable |
| **Phishing** | Fake emails that try to get clicks or information |
| **Spear phishing** | Phishing aimed at a specific group or role |
| **Vishing** | Phishing by voice call |
| **Smishing** | Phishing by text message |
| **Tailgating** | Following someone through a secure door without using your own badge |
| **OSINT** | Open-source intelligence: information found in public sources |
| **Passive OSINT** | Looking at public information without touching the target |
| **Deconfliction** | Making sure a real response team does not mistake the test for a real attack (or vice versa) |
| **SOC** | Security Operations Center: the team that watches for attacks |
| **MSSP** | Managed Security Service Provider: an outside company that does SOC work |
| **MFA** | Multi-factor authentication: more than a password to sign in |
| **Phishing-resistant MFA** | MFA (such as security keys or passkeys) that does not work on a fake website |
| **SPF / DKIM / DMARC** | Email authentication methods that help prove an email is genuine |
| **BEC** | Business email compromise: fake requests, usually for money |
| **Canary test** | Typing a fake, unique value to prove your system does not store typed content |
| **Flag** | A harmless marked item used as proof that a location was reached |
| **Pseudonymization** | Replacing names with random IDs |
| **Aggregation** | Reporting groups, not individuals |
| **Small-cell suppression** | Hiding exact numbers when a group is so small that people could be identified |
| **Just culture** | Treating honest mistakes in a flawed system as learning, not misconduct |
| **Allow-listing** | Telling a filter to let specific test messages through |
| **Kill switch** | A fast way to stop all test activity |
| **Stop word** | A word that means "all activity stops now" |
| **Retest** | Repeating the test later to see whether the fixes worked |
| **Control** | Any protection: a filter, a rule, a lock, a procedure, a training |
| **Defense in depth** | Many layers, so one failure is not a disaster |
| **Controller / processor (privacy law)** | The party that decides why data is used / the party that handles it on their behalf |
| **DPIA** | Data Protection Impact Assessment: a privacy risk review |

---

## Appendix D: Things that get testers into trouble (and the fix)

| What happened | Why it went wrong | The fix |
|---------------|-------------------|---------|
| Tester arrested during a physical test | Signer did not control the premises | Property owner consent and counter-signature (Section 1.4 and 1.5) |
| Phishing test caused a real incident response and an external report | SOC not deconflicted | Circle 3 contact and a defined deconfliction time |
| Employee in distress after a test message | Pretext used a sensitive topic | Exclusion list and HR review of every pretext |
| Employee fired after a test | No no-discipline commitment | Contract clause, HR confirmation, group-only reporting |
| Passwords found in server logs | Request logging enabled | Canary test and event-only landing pages |
| Lookalike domain suspended mid-test | Registrar or host policy ignored | Check provider policy; prefer client-owned domains |
| SMS blocked or flagged as spam | Sender and template not registered | Register sender IDs and templates early |
| Vendor staff complained | Vendor not consulted | Third-Party Register and written consent |
| Client disputed scope after the test | Verbal changes | Written change control |
| Tester recorded a call and broke a law | All-party consent rules ignored | No recording unless legal sign-off |
| Report named a department as "worst" | Ranking and blame | Wording table and aggregate reporting |
| Personal social media used for a pretext | Out of scope and invasive | Passive, role-level OSINT only |
| Real attacker found and test contaminated the case | No real-incident clause | Stop, notify, hand over evidence |
| Lock-picking tools found on a tester | Local law treats them as burglary tools | Do not carry unless RoE and counsel allow |
| Test ran during layoffs | Blackout dates not checked | Ask about stressors at kickoff |
| Executive complained loudly | Executives included without preparation | Explicit sponsor approval and private debrief |
| Data kept for years on a laptop | No retention plan | Retention table, deletion confirmation |
| "Unsolicited test" sent to a company to win business | No authorization | Never test without a signed agreement |

---

## Appendix E: Where to learn more

| Resource | Use |
|----------|-----|
| **PTES (Penetration Testing Execution Standard)** | Overall structure of a penetration test, including pre-engagement |
| **OSSTMM** | Methodology that includes human security testing |
| **NIST SP 800-115** | Technical guide to security testing and assessment |
| **NIST SP 800-50** | Building a security awareness and training program |
| **NIST SP 800-63B** | Guidance on authenticators such as MFA |
| **NIST Phish Scale** | Method for rating the difficulty of phishing emails |
| **TIBER-EU, CBEST, and DORA threat-led testing** | Frameworks for regulated financial entities |
| **CISA and partner guidance on phishing** | Practical defensive guidance (for example "stopping the attack cycle at phase one") |
| **UK NCSC guidance on phishing** | Layered defense approach for organizations |
| **Social-Engineer.org** | Ethics and community resources for social engineering testing |
| **Chris Hadnagy, "Social Engineering: The Science of Human Hacking"** | Foundational book |
| **James Reason, "Human Error" (Swiss cheese model)** | Why failures are system failures |
| **Sidney Dekker, "Just Culture"** | How to treat mistakes fairly |
| **Your local bar association or law society** | Finding a lawyer who understands security testing |

---

## Final reminder

The three questions to ask before every action in an engagement:

1. **Is it written down and signed?** (Authorization)
2. **Could it hurt a person?** (Ethics)
3. **Would I be comfortable if the employee, the client's lawyer, and a police officer all watched me do it?** (Judgment)

If any answer is "no" or "I'm not sure," stop and ask.

*End of guide.*
