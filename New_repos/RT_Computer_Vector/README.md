# 🌐 Web Exploitation

Write-ups and notes on **web application vulnerabilities** — how they're found, why they work, and how they get fixed. This repo focuses on methodology and reasoning, backed by real-world findings from bug bounty and pentest engagements.

> ⚠️ **For educational and authorized-engagement use only.** Testing applications you don't own or have explicit permission to test is illegal. Nothing here is a substitute for a signed scope/authorization.

---

## 📖 What's in here

| Area | Description |
|---|---|
| **Injection** | SQLi, command injection, template injection — root cause and exploitation logic |
| **XSS** | Reflected, stored, and DOM-based cross-site scripting, plus filter bypass reasoning |
| **Auth & Session** | Authentication bypasses, JWT abuse, session fixation, broken access control |
| **IDOR / Access Control** | Object reference flaws, privilege escalation via broken authorization |
| **SSRF** | Server-side request forgery — internal pivoting and cloud metadata abuse concepts |
| **Deserialization** | Insecure deserialization patterns across languages/frameworks |
| **Business Logic** | Flaws that don't fit a CVE category but break the application's intended flow |

---

## 🧠 Philosophy

Most write-ups here follow the same shape:

1. **What the app was doing wrong** — the root cause, in plain terms
2. **How it was found** — the actual testing/recon process, not just the payload
3. **Why the payload works** — the mechanism, not just "paste this"
4. **How it's fixed** — the defensive takeaway

The goal is understanding the "why" behind a vulnerability class well enough to spot it in a completely different application, not just replaying known payloads.

---

## 🗂️ Structure

```
web-exploitation/
├── injection/         # SQLi, command injection, template injection
├── xss/                # Reflected, stored, DOM-based
├── auth-bypass/         # Authentication & session flaws
├── idor-access-control/  # Broken object/function level auth
├── ssrf/                   # Server-side request forgery
├── deserialization/         # Insecure deserialization notes
└── business-logic/           # Logic flaws & edge cases
```

---

## 🔗 More of my work

Part of a broader collection — see the full hub: **[shivanshsahajpal22-dev.github.io](https://shivanshsahajpal22-dev.github.io/shivanshsahajpal22-dev/)**

- CTF Write-ups
- Red Team Operations
- Exploit Dev
- Malware Analysis
- Social Engineering

---

> "Discipline is the only path to mastery." — *Inspired by Miyamoto Musashi* 
