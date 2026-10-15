# How Tor and the dark web work

## The simple version (explain it like I'm a kid)

Imagine you want to send a secret note to a friend, but you don't want anyone along the way to know it was *you* who sent it, or where it's *finally* going. So you ask three other friends to help.

You wrap your note in three envelopes, one inside another, like Russian nesting dolls. The outermost envelope is addressed to Friend 1. Inside it is another envelope addressed to Friend 2. Inside *that* is one addressed to Friend 3. And inside that final envelope is your real note, addressed to your actual friend.

You hand the whole triple-wrapped package to Friend 1. They open their layer and see only "pass this to Friend 2" — they don't know what's inside or where it ultimately ends up. Friend 2 opens their layer and sees only "pass this to Friend 3" — they don't know it came from you. Friend 3 opens the last layer, finally sees the real note, and delivers it — but has no idea it started with you.

Nobody in the chain ever knows both *who sent it* and *what it says and where it's really going*. That's the whole trick, and it's why it's called "onion routing" — layers, like an onion.

---

## What problem Tor is actually solving

Regular HTTPS encryption hides *what* you're sending — your ISP can't read your messages. But it doesn't hide *who's talking to whom*. Anyone watching the network (your ISP, a government, the website itself) can still see that your IP address connected to that website's IP address, even if the content is scrambled. Tor's target is that metadata — the connection graph — not just the content.

---

## Building the circuit (the technical part)

Your Tor client picks three relays out of the thousands run by volunteers worldwide: a **guard** (entry) relay, a **middle** relay, and an **exit** relay. It builds a "circuit" through them one hop at a time — this is called telescoping. For each hop, the client performs a Diffie-Hellman-style key exchange (Tor uses an optimized version called the **ntor handshake**, built on Curve25519) to negotiate a private symmetric key with that relay alone. So each relay ends up sharing a key only with you — never with each other.

### The layered encryption (the actual "onion")

Once the circuit exists, your client encrypts your data three times, once per relay, in reverse order: the exit relay's key goes on first (innermost), then the middle relay's key, then the guard relay's key (outermost). Each relay can only peel off the layer meant for it, revealing just enough to know where to forward the packet next — nothing more:

```
[ Guard's layer
    [ Middle relay's layer
        [ Exit relay's layer
            { your real message + the real destination address }
        ]
    ]
]
```

### Cells, directory authorities, and the guard relay

- Data isn't sent in arbitrary chunks — it's split into fixed-size units called **cells** (traditionally 512 bytes). Making every unit the same size is deliberate: it stops an eavesdropper from fingerprinting your traffic by its packet sizes.
- Your Tor client doesn't guess which relays exist. A small set of trusted **directory authorities** vote on a signed list every hour — the "consensus" — listing every known relay, its public key, bandwidth, and exit rules. This consensus is cached and mirrored across the network so your client isn't obviously phoning one central server every time it wants a relay list.
- The **guard relay** is chosen and then kept for weeks or months at a time (not re-randomized per session), specifically to reduce the odds that, over time, an attacker running many relays eventually becomes your guard and starts learning things about you.

### Following the path: what each hop knows

| Hop | Knows | Doesn't know |
|---|---|---|
| You (Tor client) | Everything — it built the whole circuit | — |
| Guard relay | Your real IP address | What site you're visiting, or the content |
| Middle relay | The guard's IP and the exit's IP | Your real IP, the destination, or the content |
| Exit relay | The destination site, and unencrypted content if the site itself isn't using HTTPS | Your real IP |
| Destination site | The exit relay's IP address | Your real IP |

One important nuance: once traffic reaches the exit relay, it's decrypted back to plain form before leaving the Tor network toward the actual website (HTTPS, if the site uses it, is a separate layer of encryption that still protects the content here). A malicious exit relay operator can, in principle, read or tamper with unencrypted traffic passing through it — it just still has no idea *who* sent it. That's the reasoning behind "always use HTTPS even over Tor."

---

## Hidden services — how a server makes itself anonymous

This is a different trick from the one above. So far we've covered *you* anonymously reaching an ordinary website. A `.onion` address is a *server* hiding its own location, so that nobody — not even the relays helping it — knows where it physically is.

### Explain it like I'm a kid: the secret clubhouse

Imagine a shop owner who wants people to be able to buy things from them, but doesn't want anyone to know where their shop actually is.

Here's what they do: they pick a few trusted friends scattered around town and tell each one, "If someone comes looking for me, don't tell them where I am — just pass along a message that I'll meet them somewhere." Then they write a note and pin it to a hidden bulletin board that only tells people *which friends* to talk to — not the shop's real location.

Two more pieces make this work, and they split into a public half and a secret half. The shop owner has a personal **signet ring** — the physical object, always in their pocket, never shown to anyone. What people actually get is the **wax imprint** the ring leaves on a letter — a recognizable pattern that anyone can see, copy down, and recognize again later. That imprint pattern is exactly what your address for the shop is made from, and it's not secret at all — it's handed out freely to anyone who wants to visit.

If the shop owner pinned their note to the board under that same fixed imprint every day, board-watchers could track the same shop posting over and over. So instead, each day, both the shop owner *and* anyone who already has the imprint pattern can independently do a bit of public arithmetic — mixing the imprint together with today's date, using a fixed rule anyone can look up — and land on the same fresh-looking **code name**, say "GRIMSBY-19" today, "PELICAN-04" tomorrow. Nobody needs the actual ring for this part, just the imprint. That's the spot on the board where today's note goes.

But knowing the slot isn't the same as being allowed to post there convincingly. To make today's note trustworthy — proof it really came from the shop owner, and not some forger who also worked out that "GRIMSBY-19" was today's slot — the shop owner needs the *actual ring itself*, pressed together with today's date, to leave a fresh authenticating seal on the note. Only the real ring, in the real owner's pocket, can produce a seal that checks out. A stranger who's never seen the imprint can't compute "GRIMSBY-19" in the first place, and even a stranger who somehow guessed the slot couldn't forge a seal that passes inspection, because that step needs the ring itself, not just its imprint.

When you want to buy something, you find that note, walk up to one of the trusted friends, and say "I'm looking for the shop owner." That friend passes your message along a chain toward the shop owner — without ever learning where *you* live either. The shop owner gets the message and says, "Let's meet at the old oak tree at noon." That oak tree isn't the shop owner's home and isn't your home — it's neutral ground picked just for this one meeting.

Both of you then travel there through your own separate chains of trusted friends, and you do your business at the tree. The oak tree never finds out where either of you actually lives. It just knows two people showed up and handed things back and forth.

### The technical version

1. The service generates a long-term keypair (**Ed25519**, in the current version-3 protocol) and selects several relays to act as **introduction points**.
2. It publishes a signed "descriptor" — its public key plus its introduction points — into a distributed hash table maintained by relays flagged as **HSDir** (hidden service directory nodes). It publishes this over its own Tor circuit, so even the act of publishing doesn't reveal the server's real IP.
3. To make that descriptor's storage location unpredictable to outsiders (and to stop anyone from just watching one fixed spot in the DHT), the service doesn't publish under its raw public key. It derives a **blinded key** from its long-term key plus the current time period, and the descriptor is stored under a hash of *that*. The blinded key changes on a schedule, so the storage location shifts over time without ever exposing the underlying long-term identity key.
4. The `.onion` address itself is derived directly from the long-term public key (base32-encoded, with a version byte and checksum, 56 characters long). This makes it "self-authenticating" — there's no DNS to trust or spoof. If you have the correct address, you're cryptographically guaranteed you're talking to the holder of that specific key, not an impersonator.
5. When you want to visit the site, your client fetches the current descriptor from the HSDir network, builds a circuit to one of the listed introduction points, and separately builds a circuit to a randomly chosen **rendezvous point**.
6. Your client tells the hidden service (via the introduction point) to meet you at that rendezvous point. The service then builds its own circuit to that same rendezvous relay.
7. The rendezvous relay now just shuttles encrypted traffic between your circuit and the service's circuit. It never learns either party's real IP, and neither party ever learns the other's.
8. Introduction points themselves are rotated periodically and can be swapped out if one becomes unreliable, so the service isn't permanently dependent on any single relay to be reachable.

Both ends stay anonymous to each other and to every relay involved — that's the full mechanism that makes a `.onion` site reachable without anyone, including Tor relay operators, learning what server is actually behind it.

---

## "Dark web" vs "deep web" vs "surface web"

These get conflated a lot, so it's worth being precise:

- **Surface web** — pages search engines can crawl and index.
- **Deep web** — everything *not* indexed: your email inbox, your bank's login-protected pages, private databases. This is the vast majority of the internet, and it's almost entirely mundane.
- **Dark web** — specifically the small set of sites only reachable through an anonymity network's own addressing scheme, like Tor's `.onion` (or I2P's "eepsites"). It's a tiny sliver of the deep web, not a synonym for it.

## What's actually there

Both mundane and less mundane things exist side by side. On the legitimate end: news organizations like the BBC run official `.onion` mirrors so people in censored or surveilled countries can reach them; SecureDrop instances let journalists receive tips from anonymous whistleblowers; activists and dissidents use it to communicate where the open internet is monitored or blocked. On the illicit end, it's also home to black markets and forums for various illegal activity — worth knowing exists as a factual matter, though not something worth detailing further here.

## Where the anonymity actually breaks down

Tor's protection isn't absolute, and it's worth understanding the real threat model:

- **Traffic correlation** — if an adversary can watch traffic entering the network *and* traffic leaving it at the same time, they can match timing and volume patterns statistically and unmask the connection. This is Tor's main acknowledged weakness: it assumes no single party can watch the whole network at once, which is reasonable against most attackers but not against the most resourced state-level ones.
- **Malicious exit nodes** — as mentioned, anyone can run one, and a bad actor running one can see unencrypted traffic passing through it.
- **Operational mistakes, not cryptographic breaks** — most real-world dark web takedowns (Silk Road being the famous case) came from operators making ordinary mistakes: reusing a username, misconfiguring a server, leaving identifying details in a forum post — not from anyone actually breaking Tor's encryption or routing.

## A note on legality

Using Tor itself is legal in most countries (it's blocked or restricted in a handful of places with heavy internet censorship). What's illegal is whatever illegal activity someone chooses to do with it — the network is just a transport layer, morally neutral in the same way the postal service is.
