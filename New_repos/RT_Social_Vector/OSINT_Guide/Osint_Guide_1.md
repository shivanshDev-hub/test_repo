
# OSINT GUIDE PART 1

## MY OSINT DASHBOARD 

> This is people-related OSINT, not company-related. If you want that, it's in the web exploitation guide passive recon 

**last updated** - `04 october 2026`

Organized by Investigative Pivot. Each category leads with the tools actually worth reaching for; anything kept only for completeness sits in a **Deprioritized** bucket at the end of that same category, not mixed in above the fold.

*Tags used inside Deprioritized buckets: **(unvetted)** = arrived in a bulk import from the Awesome-OSINT list and has not been tested yet; **(status unconfirmed)** = looks dead, retired or replaced, but the domain was not re-checked.*
 
---
 
### 1. Automation Frameworks & All-in-One Platforms
 
- **Self-hosted / open-source automation:** [SpiderFoot](https://github.com/smicallef/spiderfoot) or [Recon-ng](https://github.com/lanmaster53/recon-ng) — open these first on any investigation
- **Semi-automatic framework with a module registry:** [sn0int](https://github.com/kpcyrd/sn0int)
- **Method map / "where do I even start" directory:** [OSINT Framework](https://osintframework.com/)
- **Hosted all-in-one platforms:** [Intell Base](https://intelbase.is/) — note: [OSINT CAT](https://www.osintcat.net/) also markets itself here, but its actual core function is breach-checking — see §10 (Exposed Data, Leaks & Code Repositories), listed there once to avoid duplicating the same URL under two labels

**Deprioritized:** Maltego (heavily paywalled/restricted free tier), Datasploit (unmaintained); Sintelix (marketed as a Maltego-style link-analysis tool, little independent evidence behind it); Photon (older crawler, status unconfirmed); Outwit (general data grabber, not investigation-specific); OSINT.SH and Find OSINT Tool (thin launcher pages that duplicate OSINT Framework); newer "all-in-one" platforms from the bulk import: PRISM, OSINTNova, OpenGraph Intel (OGI), IntellyWeave, Taranis AI, LinkScope (and LinkScope Client), Pyba, IntelHub, OSINT-Tool browser extension (all unvetted)
 
---
 
### 2. Person Identity — Name-Based Pivots
 
- Address search: [TruePeopleSearch](https://www.truepeoplesearch.com) (US only) or [192.com](https://www.192.com) (UK)
- Username/email generation: [soxoj's username-generation-guide](https://github.com/soxoj/username-generation-guide) + [transform_username](https://github.com/soxoj/transform_username)
- Phone number by name: [Truecaller](https://www.truecaller.com)
- Company/workplace: [OpenCorporates](https://opencorporates.com) — full corporate-records listing under §8, then pivot to country-based registries
- US public records: [Judyrecords](https://www.judyrecords.com/) (400M+ court cases), [Voter Records](https://voterrecords.com/), [Black Book Online](https://www.blackbookonline.info), [Federal BOP Inmate Locator](https://www.bop.gov/inmateloc)
- Sanctions & public-office holders: [OpenSanctions](https://www.opensanctions.org/search/)
- Genealogy & historical records: [FamilySearch](https://familysearch.org) (free, registration required), [FamilyTreeNow](https://familytreenow.com), [The National Archives (UK)](https://www.nationalarchives.gov.uk)
- International / regional directories: [Infobel](https://www.infobel.com/) (73 countries), [Canada411](https://www.canada411.ca)

**Deprioritized:** Pipl (now fully commercial/paid-only), Spokeo (heavy restrictions/paywalls on free lookups); the same "limited free info, premium upsell" pattern: 411, BeenVerified, InfoTracer, PeekYou, Reunion, SearchBug, Social Catfish, White Pages, ZabaSearch, Ancestry, Genealogy Bank, Fold3, Classmates; credit-limited sales-prospecting tools: Apollo.io, ContactOut; ClustrMaps (overlaps TruePeopleSearch); UK Phone Book (duplicates 192.com); Homemetry (thin reverse-address lookup); Mugshots.com (low signal, takedown-fee business model); VineLink (victim-notification service, not a search tool); Genealogy Links (a link directory, not a search tool); FinalNotes Obituary Research Guide (a guide, not a tool); Offendersearch and VerifiedHer (unvetted — prefer the official registries for the former); BuscaPaginasBlancas (Spain-only white-pages scraper, unvetted — promote it if you work Spanish targets)
 
---
 
### 3. Username Pivots
 
- Search engines: [maigret](https://github.com/soxoj/maigret) or [Sherlock](https://github.com/sherlock-project/sherlock)
- Extra coverage: [Blackbird](https://github.com/p1ngul1n0/blackbird) (username and email, 600+ sites), [Snoop](https://github.com/snooppr/snoop/blob/master/README.en.md) (strong on the Russian-speaking web)
- Stylometry (writing-style analysis): [JStylo](https://github.com/psal/jstylo)
- Snowflake ID decoding: [singhajit.com snowflake-decoder](https://singhajit.com/tools/snowflake-decoder)
- Shortcode/ID converters: [Tweeterid.com](https://tweeterid.com) + [Instagram Shortcode Converter](https://techconverter.me/instagram-shortcode-to-date)

**Deprioritized:** WhatsMyName (older local Python wrappers lacking active updates compared to the web version/maigret); NexFil (overlaps maigret/Sherlock); Name Chk, Name Checkr, Name Checkup, NameKetchup, CheckUser (built to check brand-name availability, not to investigate a person); IDCrawl, User Search, User Searcher, Digital Footprint Check, SherlockEye, TraceFind, Trace (manus.space), user-scanner, Seekr, Antisocial (web front-ends or newer wrappers around the same site lists — unvetted); Cupidcr4wl (username and phone search limited to adult-content platforms, too niche for routine use — unvetted)

`These are quite useful but require proper drafting. if you are reading this help me out by placing them in thier category and them making a merge request`

**Some other platforms:** [tomba](https://tomba.io/) , [InstantUsername.com](https://instantusername.com/) , [whatsmyname CheapCopy](https://whatsmynameapp.org/) , [wolver eye](https://wolverbot.online/) , [osint search.org](https://osintsearch.org/) , [social searcher](https://www.social-searcher.com/) , [intel-techniques.com](https://inteltechniques.com/tools/) , [arkm.com](https://arkm.com/) , {This category is particularly from the osint lectures i attended}

---
 
### 4. Image & Visual Pivots
 
- Reverse image search: [Yandex Images](https://yandex.com/images) or [Google Lens](https://lens.google.com) — third opinion: [Bing Images](https://www.bing.com/images); Chinese-web content: [Baidu Images](https://image.baidu.com)
- Face search: [Facecheck.id](https://facecheck.id) or [lenso.ai](https://lenso.ai) — also [PimEyes](https://pimeyes.com) and [Search4faces](https://search4faces.com/) (covers Russian-language networks)
- Metadata extraction: [Exiftool](https://exiftool.org) or [EXIF.tools](https://exif.tools) — also [Jeffrey's Image Metadata Viewer](https://exif.regex.info/)
- Tampering / edit detection: [FotoForensics](https://www.fotoforensics.com) or [Forensically](https://29a.ch/photo-forensics/)
- Photo-hosting search: [Flickr](https://flickr.com/search/) — geotagged and metadata-rich uploads
- Street view coverage: [Google Street View](https://www.google.com/maps), [Yandex Panorama](https://yandex.com/maps), [Mapillary](https://www.mapillary.com), [KartaView](https://kartaview.org), [Panoramax](https://panoramax.fr) — plus [KartaVision](https://kartavision.com/) (natural-language and image search over KartaView imagery)
- Advanced street-view query engine: [Overpass-turbo](https://overpass-turbo.eu)
- Geolocation hunters: [GeoSeeker](https://geoseeker.com), [GeoAxis](https://geoaxis.com), [Picarta.ai](https://picarta.ai), [GeoSpy](https://geospy.web.app/), [GeoInfer](https://geoinfer.com), [ReverseImageLocation](https://reverseimagelocation.com), [TracePoint](https://kluter.github.io/TracePoint/) (ray-intersection geolocation, runs client-side)
- Geolocation hints database: [GeoHints](https://geohints.com)
- Shadow-based location pinpointing: [SunCalc](https://www.suncalc.org), [ShadowMap](https://shadowmap.org), [ShadowFinder](https://github.com/bellingcat/ShadowFinder), [ShadowFinder-Web](https://kluter.github.io/ShadowFinder-Web) (browser port, no Jupyter needed)
- Past world imagery: [wayback world imagery](https://livingatlas.arcgis.com/wayback/) and [google earth](https://earth.google.com/web/)

**Deprioritized:** TinEye (outperformed by modern neural reverse-search engines); Google Image (superseded by Lens); Dupli Checker, Lycos Image Search, Yahoo Image Search, PhotoBucket (weak indexes); Pixsy (copyright-enforcement product, not a search tool); Clarify, Image Identification Project, Betaface (object/face demos, not investigative search); PicTriev, Faceagle, Surfface (smaller face-search engines — unvetted); ImgOps (a hub of links to the tools above); ImpulseAdventure, JPEGsnoop, ExifLooter, EXIFEditor.io, JIMPL, Metadata Viewer (kriztalz) (redundant next to Exiftool); forensics.media, JPEG Audit, ChronoVerify, ProfileImageIntel (newer image-forensics front-ends — unvetted); DiffChecker image diff (generic); GeoSpyer (wrapper that depends on a third-party geolocation service); Barcode Reader (niche)
 
---
 
### 5. Email Pivots
 
- Search engines: [Epieos](https://epieos.com), [Holehe](https://github.com/megadose/holehe), [Mailmeteor](https://mailmeteor.com) — Gmail specifically: [GHunt](https://github.com/mxrch/GHunt)
- Email → accounts, extra coverage: [Blackbird](https://github.com/p1ngul1n0/blackbird)
- Format prediction: [Mailmeteor Email Permutator](https://mailmeteor.com/email-permutator)
- Company email patterns: [Email Format](https://email-format.com), [Hunter](https://hunter.io) (limited free tier)
- Verification: [MyEmailVerifier](https://myemailverifier.com) or [VerifyEmailAddress](https://www.verifyemailaddress.org) — self-hosted option: [Reacher](https://reacher.email) (open-source, Rust)
- Reputation / risk scoring: [EmailRep](https://emailrep.io)
- PGP key → identities and emails: [PGPKeyAnalyser](https://kriztalz.sh/pgp-key-analyser/) (in-browser, no .asc download) — see also the PGP note under Bonus tools
- Breach data: [Have I Been Pwned](https://haveibeenpwned.com) — broader breach-search catalog in §10; local breach hunting: [h8mail](https://github.com/khast3x/h8mail)
- Mail blacklist check: [MXToolbox SuperTool](https://mxtoolbox.com/SuperTool.aspx)

**Deprioritized:** Skymem (drastically throttled free export limits); Blacklist Checker, Multirbl, Spamhaus lookup (redundant next to MXToolbox); Email Address Validator, Verify Email, MailTester, EmailHippo (verification wrappers, redundant with the verifiers above); Email Extractor (99tools), Email Permutator (polished.app) (duplicates Mailmeteor's); Peepmail (older, status unconfirmed); Snov.io, Toofr, VoilaNorbert (paywalled sales-lead email finders); Pipl (commercial); ThatsThem (people-search aggregator, thin signal); SherlockEye, TraceFind, user-scanner, Minerva OSINT, MailAccess, OSINTEye (newer aggregators wrapping Holehe-style checks — unvetted); DeHashed, LeakCheck, LeakRadar, InfoStealers (breach-oriented, filed under §10)
 
---
 
### 6. Phone Number Pivots
 
- [Truecaller](https://www.truecaller.com) or [PhoneInfoga](https://github.com/sundowndev/phoneinfoga)
- Carrier / line type: [FreeCarrierLookup](https://freecarrierlookup.com/) (carrier, wireless vs landline), [Phone Validator](https://www.phonevalidator.com/index.aspx) (good at catching Google Voice numbers)
- India prefix → state/operator mapping: [InMobPrefix](https://github.com/hstsethi/in-mob-prefix)

**Deprioritized:** Twilio Lookup basic carrier APIs (heavily restricted to verified commercial accounts); EmobileTracker (promises live phone-location tracking, which a web lookup cannot deliver); Sync.ME (redundant next to Truecaller); Reverse Phone Check, Reverse Phone Lookup, SearchPeopleFREE, USPhonebook, TraceFind (thin, upsell-heavy reverse-lookup sites); Spy Dialer (US-centric voicemail-greeting and owner-name lookup, hit-and-miss — unvetted)
 
---
 
### 7. Domain & Website Pivots
 
**Before anything manual:** [theHarvester](https://github.com/laramies/theHarvester)
 
- WHOIS/ownership: [WhoisXML](https://www.whoisxmlapi.com) or [ICANN Lookup](https://lookup.icann.org) — reverse and historical WHOIS: [DomainTools Whois](https://whois.domaintools.com), [Whoisology](https://whoisology.com)
- Subdomain enumeration: [crt.sh](https://crt.sh) or [Amass](https://github.com/owasp-amass/amass) — crt.sh is also your certificate-transparency tool, cross-referenced in §9; probe and resolve the results with [dnsx](https://github.com/projectdiscovery/dnsx); CT-log alternative: [Merklemap](https://www.merklemap.com/)
- DNS records & history: [SecurityTrails](https://securitytrails.com/dns-trails), [DNSDumpster](https://dnsdumpster.com), [Validin](https://app.validin.com/) (free current + historical DNS)
- Tech stack & site profile: [BuiltWith](https://builtwith.com), [Netcraft Site Report](https://sitereport.netcraft.com/), [Web-Check](https://web-check.as93.net/) (all-in-one site and server metadata)
- Same-owner pivots (analytics IDs, tracker codes, favicon hash): [PublicWWW](https://publicwww.com/), [AnalyzeID](https://analyzeid.com/), [NerdyData](https://nerdydata.com), [FaviconHash](https://kriztalz.sh/favicon-hash/) (hash → Shodan/Censys/FOFA)
- Microsoft 365 / Azure tenant lookup from a domain: [AADInternals OSINT](https://aadinternals.com/osint), [PingCastle Tenant Resolution](https://tenantresolution.pingcastle.com)
- Archived/deleted content: [Wayback Machine](https://web.archive.org) — every URL it knows for a domain: [waybackurls](https://github.com/tomnomnom/waybackurls); capture and archiving tools in §27
- Domain → IP resolution: [dig](https://linux.die.net/man/1/dig) — reverse DNS / AS context: [Robtex](https://www.robtex.com)
- Exposed document metadata: [FOCA](https://github.com/ElevenPaths/FOCA) or [Metagoofil](https://github.com/laramies/metagoofil)
- URL anatomy: [Unfurl](https://dfir.blog/unfurl/)

**Deprioritized:** Sublist3r (unmaintained, broken API endpoints); SEO and traffic suites built for analysing your own site, weak ownership signal: ahrefs, Majestic, OpenLinkProfiler, WebMeUp, SEMrush, Serpstat, SpyFu, Quantcast, Quick Sprout, SEO Chat Tools, SEOTools for Excel, SmallSEOTools, Accuranker, Bing Webmaster Tools, Similar Web; WHOIS/DNS clones that add nothing over the tools above: Central Ops / Domain Dossier, Domain Crawler, Easy Whois, Mark Monitor WHOIS, Who.is, WhoIsHostingThis, WhoisMind, WhoIsRequest, Website Informer, HypeStat, StatsCrop, w3snoop, WhatIsMyIPAddress, You Get Signal, Kloth, Dedicated or Not, DNSStuff, Remote DNS Lookup, Verisign DNSSEC debugger, DNSViz (DNSSEC-only), PageGlimpse, RedirectDetective, Follow.net, GraphyStories, DNS History (completedns); same-owner pivots that overlap AnalyzeID and NerdyData: SameID (unvetted); newer all-in-one domain tools (all unvetted): digga, Domain Intelligence, DomainRecon, CC.LA, DNSai, IntoDNS.ai, IP.THC.ORG, Columbus Project, Subhunt, SubDomainRadar.io, CrawlGraph, Cerast Intelligence, WebsiteTechMiner.py (script wrapper around the BuiltWith/Wappalyzer APIs); availability and typosquat tools, brand-protection use rather than people-OSINT: Domain Hunter, Vacato, Squatm3gator, openSquat
 
---
 
### 8. Company & Corporate Records
 
- Ownership/registry: [OpenCorporates](https://opencorporates.com) or [Crunchbase](https://www.crunchbase.com)
- Filings & disclosures (US): [SEC EDGAR](https://www.sec.gov/edgar/search/)
- Beneficial ownership & offshore structures: [OpenOwnership Register](https://register.openownership.org/), [ICIJ Offshore Leaks Database](https://offshoreleaks.icij.org)
- Registry finder by country: [Overseas Company Registers (GOV.UK)](https://www.gov.uk/government/publications/overseas-registries/overseas-registries), [List of company registers (Wikipedia)](https://en.wikipedia.org/wiki/List_of_company_registers)
- Trademark/brand records: [WIPO Brand DB](https://www3.wipo.int/branddb/en/)
- Regional gap-filler: [YouControl](https://youcontrol.com.ua/en/) (Ukraine)

**Deprioritized:** Gleif (too narrow/niche for broad enterprise identification); paywalled databases: Bloomberg, Bureau van Dijk / Orbis directory, Factiva, Hoovers, Mergent Intellect, Mergent Online, Business Source, Plunkett Research, Vault; defunct or replaced: Data.com Connect (retired), EDGAR Online (replaced by the SEC link above, status unconfirmed); marketing, review and ranking sites mis-filed as company research: CompeteShark, Comparably, GoodFirms, Inc. 5000, Forbes Global 2000, TheWebCo (SEMrush, Serpstat and SpyFu are covered under §7); thin directories: BBB, Bizeurope, BrownBook, Europages, European Business Register, Ezilon, Kompass, Scoot, Corporate Information, CorporationWiki, Company Registration Round the World, Canadian Business Research, globalEdge, ICAEW international-registration guide, AllStocksLinks, Tracxn, GuideStar (status unconfirmed, now part of Candid), Section 16 Deadline Calculator (an SEC deadline worksheet, not an investigation tool); Glassdoor, LinkedIn, Xing (see §32); newer cross-registry aggregators: WhiteIntel, Fylings (unvetted)
 
---
 
### 9. Infrastructure, Attack Surface & Threat Intelligence
 
*(Merged from what was previously four overlapping sections — IP-address pivots, attack-surface search engines, IP/domain reputation, and general threat intel all answer the same underlying question: "what do we know about this IP/domain/host's exposure and reputation")*
 
- IP geolocation/ownership: [Ipinfo](https://ipinfo.io) and [ViewDNS.info](https://viewdns.info) — also [DB-IP](https://db-ip.com) and [IP2Location](https://www.ip2location.com/demo) (geolocation + ASN + VPN/proxy flags) — authoritative registration data for North American IP blocks: [ARIN Whois](https://whois.arin.net)
- Exposed devices/services (attack-surface search engines): [Shodan](https://www.shodan.io), [Censys](https://search.censys.io), [FOFA](https://en.fofa.info), [ZoomEye](https://www.zoomeye.ai), [ONYPHE](https://search.onyphe.io), [Hunter Search Engine](https://hunter.how), [FullHunt](https://fullhunt.io)
- Tor/VPN detection: [IPQualityScore](https://www.ipqualityscore.com), [GetIPIntel](https://getipintel.net) — also [Focsec](https://focsec.com) (API) and [ExoneraTor](https://exonerator.torproject.org) (was this IP a Tor relay on a given date?)
- IP/domain reputation: [AbuseIPDB](https://www.abuseipdb.com), [Cisco Talos Intelligence](https://talosintelligence.com/reputation_center), [GreyNoise](https://viz.greynoise.io), [Criminal IP](https://www.criminalip.io), [Shadowserver](https://dashboard.shadowserver.org), [Pulsedive](https://pulsedive.com)
- URL/file scanning & sandboxes: [urlscan.io](https://urlscan.io/), [Hybrid Analysis](https://www.hybrid-analysis.com), [URLhaus](https://urlhaus.abuse.ch) (malicious URLs), [MetaDefender](https://metadefender.com)
- Malware intelligence (abuse.ch ecosystem): [Abuse.ch Hunting](https://hunting.abuse.ch), [MalwareBazaar](https://bazaar.abuse.ch/browse/), [YARAify](https://yaraify.abuse.ch/scan/) `to Run and test a malware [anyrun](https://any.run/)`
- General IOC/CVE feeds: [VirusTotal](https://www.virustotal.com), [AlienVault OTX](https://otx.alienvault.com), [CVE Details](https://www.cvedetails.com) — also [CyberGordon](https://cybergordon.com) (30+ sources in one query) and [Cyberbro](https://github.com/stanfrbd/cyberbro) (self-hosted observable lookup)
- Threat-actor profiling: [Mitre.org](https://www.mitre.org) or [Recorded Future](https://www.recordedfuture.com) (paid) — full actor lists in §21
- Network/BGP intelligence: [BGP.tools](https://bgp.tools), [BGP.he.net](https://bgp.he.net), [BGPView](https://bgpview.io)
- Traffic/infrastructure trends: [Cloudflare Radar](https://radar.cloudflare.com)
- **Deprioritized:** MaxMind GeoLite legacy offline formats; older static signature scanners without cloud threat feeds; Netlas.io (redundant given Shodan/Censys coverage); ODIN (10-searches/day cap too limited vs. main alternatives); SikkerAPI (niche, redundant next to AbuseIPDB/GreyNoise); BrightCloud (overly basic interface, limited free telemetry); CertKit Certificate Search (crt.sh already covers this natively); IP-geolocation clones with no extra signal: IP 2 Geolocation, ip2geo.dev, IP Checking, IP Location, IPLocation.io, IPFingerprints, Infosniper, IPVoid, ISP.Tools; URL-scan and reputation wrappers: URLVoid, urlQuery, urlDNA, TinyScan, ScanMalware, Webscore, Islegitsite, isMalicious, IPOK; PhishStats (status unconfirmed); aa419 Fake Sites Database (scam-site list, manual use only); Browserling (sandboxed link opener, niche); newer feeds and platforms (all unvetted): Crypto Scam & Phishing Threat Intel Feed (spmedia), HoneyLabs, Webscout, DFIR Platform, ThreatLens, MalwareIntel, Malwagon, REScure, defend.network; added in the latest merge (unvetted): TracerouteVisualizer (maps mtr/traceroute output, niche), pygreynoise (just a Python wrapper for the GreyNoise API), CrowdSec (open-source IPS for defending your own servers, not an investigation tool)
 
---
 
### 10. Exposed Data, Leaks & Code Repositories
 
*(Merged from three previously separate lists — code/pastebin leaks, exposed cloud storage, and breach-search engines are all the same investigative move: "has this person/org's data already leaked somewhere")*
 
- Source-code leak scanning: [Gitleaks](https://github.com/gitleaks/gitleaks) or GitHub Code Search — regex over public repos: [grep.app](https://grep.app/), [SourceGraph](https://sourcegraph.com/search)
- Exposed cloud storage: [GrayhatWarfare](https://grayhatwarfare.com) — one of the highest-yield tools in this entire list for open S3 buckets
- Breach/leak search engines: [CheckLeaked](https://checkleaked.cc), [OSINT CAT / OsintCat](https://www.osintcat.net) (same site as the §1 cross-reference), [StealSeek](https://stealseek.com), [Venacus](https://venacus.com), CredenShow, HIB Ransomed, HEROIC.NOW, IKnowYour.Dad, Leaker (passive CLI, 10 breach DBs at once), NOX (recursive async breach/identity pivoting)
- Established breach-data providers: [DeHashed](https://dehashed.com/) (paid, large aggregated dataset), [LeakCheck](https://leakcheck.io/) (email, username, keyword or corporate-domain search)
- Infostealer-infection exposure: [Hudson Rock](https://www.hudsonrock.com/threat-intelligence-cybercrime-tools) (free cybercrime-intel toolkit)
- Historical paste-site & dark-web indexing: [Intelligence X](https://intelx.io/tools)

**Deprioritized:** Pastebin basic unauthenticated scrapers (heavily rate-limited by Cloudflare); haveibeenzuckered (covers a single 2021 Facebook dataset); DeHashed CLI (just a wrapper around the DeHashed API); LeakRadar, InfoStealers, Dropbase (unvetted); GitGuardian public-GitHub monitoring (a product page, not a search tool); SearchCode, Sourcebot, Code Finder (overlap grep.app and SourceGraph); pastebin hosting sites are covered in §23
 
---
 
### 11. Dark Web Monitoring
 
- [Ahmia](https://ahmia.fi) or [Dark.fail](https://dark.fail)
- Onion-service existence and metadata: [onion-lookup](https://onion.ail-project.org/) (free service and API from the AIL project)

**Deprioritized:** Torch (unreliable index uptime, dead mirrors — not to be confused with the Torch browser in §18); OnionScan (status unconfirmed, older tool); Aleph Open Search, VoidAccess (unvetted)
 
---
 
### 12. Cryptocurrency & Blockchain Tracing
 
- [Blockchair](https://blockchair.com) or [Arkham Intelligence](https://platform.arkhamintelligence.com)

**Deprioritized:** BitcoinWhosWho (archaic UI, minimal free tracking depth); Orbit (recursive wallet-graph crawler, older and status unconfirmed)
 
---
 
### 13. Social Media & Messaging Platforms
 
- Platform-aggregator search: [Social-Searcher](https://www.social-searcher.com) or [Social Analyzer](https://github.com/qeeqbox/social-analyzer)
- Name / email / phone → accounts: [IDCrawl](https://www.idcrawl.com/) (name), [Predicta Search](https://predictasearch.com) (email and phone), Epieos (see §5)
- Instant messaging (Telegram/WhatsApp/Discord): [Telemint](https://github.com/telemint/telemint) or Telegram dorks via [Google CSE](https://cse.google.com) — per-platform tools in §20
- Non-Western and alternative networks to check once you have a country: [VKontakte](https://vk.com), [Odnoklassniki](https://ok.ru) (Russia), [Weibo](https://weibo.com), [Qzone](https://qzone.qq.com) (China), [Mixi](https://mixi.jp) (Japan), [Xing](https://www.xing.com) (DACH), [Taringa](https://www.taringa.net) (Latin America), [Draugiem](https://www.draugiem.lv) (Latvia), [Bluesky](https://bsky.app), [Threads](https://www.threads.net)
- Channel/account statistics history: [SocialBlade](https://socialblade.com)

**Deprioritized:** Twint (frequently breaks from X/Twitter API structural changes), ChatWatch (outdated WhatsApp status-tracking scripts); marketing-analytics suites (audience and brand monitoring, not person-level investigation): Audiense, Bottlenose, Brandwatch, Buffer, BuzzSumo, Hootsuite, Klear, Kribrum, MustBePresent, Netvibes, OpinionCrawl, Rival IQ, SocialBakers, Tagboard; Geocreepy (discontinued project, status unconfirmed); Mail.Ru Social Network Search (status unconfirmed); Social DownORNot (an outage checker); UVRX, WATools, Oblivion, Castrick, OSINTNova (unvetted); mainstream platforms with native search only (per-platform tools live in §20): Facebook, Instagram, Pinterest, Reddit, Tumblr, X/Twitter, and Tinder (login-walled app with no public search — reverse-image a profile photo via §4 instead)
 
---
 
### 14. Web Application & Mobile App Recon
 
- WordPress-specific scanning: [Wpscan](https://wpscan.com)
- Mobile app asset discovery (subdomains/URLs/parameters embedded in apps): [BeVigil](https://bevigil.com/search)
- Security-header and best-practice audit: [Mozilla HTTP Observatory](https://developer.mozilla.org/en-US/observatory)
- TLS configuration check: [Qualys SSL Labs](https://www.ssllabs.com/ssltest/)

**Deprioritized:** Router Passwords (default-credential reference, relevant only to authorised testing); ReconFTW, SpiderSuite, Sub3 Suite, osint-recon (automation for authorised engagements — belongs in the web exploitation guide, not a people-OSINT board); Pentest-Tools Google-hacking (paid wrapper around dorks, see §16); PhishingSecLists, keyFinder (unvetted); BeVigil-CLI (command-line and Python wrapper for the BeVigil API, redundant next to the web search above)
 
---
 
### 15. Transportation Tracking
 
- Flight (live): [FlightRadar24](https://www.flightradar24.com)
- Flight (history/data): [FlightAware](https://www.flightaware.com)
- Marine traffic: [MarineTraffic](https://www.marinetraffic.com)
- Vessel/container tracking: [VesselFinder](https://www.vesselfinder.com)
- Rail: [OpenRailwayMap](https://www.openrailwaymap.org)
- Plate formats & vehicle history: [World License Plates](https://www.worldlicenseplates.com/) (identify a plate's country or region — pairs with the plate checklist in the Geo-OSINT part), [FaxVIN](https://www.faxvin.com/) (VIN or plate → make, model, age)

**Deprioritized:** PlaneFinder (less comprehensive tracking density than FlightRadar24); EpicVIN (paid vehicle-history reports, upsell-heavy); Logibook (reference pages for ports, airports and terminals, no live data — unvetted); Phantom Tide (combined vessel and airspace dashboard — unvetted)
 
---
 
### 16. Search Engines — General, National & Specialized
 
**General dorking:**
[Dorkgpt](https://www.dorkgpt.com), [pagodo](https://github.com/opsdisk/pagodo), [Google Hacking Database](https://www.exploit-db.com/google-hacking-database), [SearchDorks](https://kriztalz.sh/search-dorks/) (generates dorks for Google, FOFA, Shodan, Censys and ZoomEye)

**Independent & self-hosted search (a different index gives different results):**
[Brave Search](https://search.brave.com), [Mojeek](https://www.mojeek.com/) (own crawler, no tracking), [SearXNG](https://searxng.org/) (self-hosted metasearch), [Perplexity](https://www.perplexity.ai) (AI answer engine with citations — verify every claim at the source)

**Custom/scoped search:**
[Google Custom Search](https://www.google.com/cse) (build a scoped engine for a recurring investigation), [Million Short](https://millionshort.com) (excludes the top N most popular results — surfaces content SEO-dominant sites bury)
 
**National — Tier 1 (Strategic Global Powers):**
[Baidu](https://www.baidu.com) (China), [Yandex](https://yandex.com) (Russia/Eastern Europe, also strongest reverse-image search), [Naver](https://www.naver.com) (South Korea), [Daum](https://www.daum.net) (South Korea), [SoGou](https://www.sogou.com) (China, indexes WeChat official-account articles)
 
**National — Tier 2 (Global):**
[Google](https://www.google.com), [Bing](https://www.bing.com), [Yahoo](https://www.yahoo.com)
 
**National — Tier 3 (Regional/Localized):**
Alleba (Philippines), Eniro (Sweden/Nordics), Gerdoo (Iran), Goo (Japan), Najdi (Slovenia), Onet.pl (Poland), Orange (France), Parseek (Iran), SAPO (Portugal), Search.ch (Switzerland), Seznam (Czech Republic), Walla (Israel), Zarebin (Iran), Coc Coc (Vietnam)
 
**Visual/clustering search:**
[Carrot2](https://search.carrot2.org) (clusters results by topic, actively maintained/open-source)
 
**Similar-site discovery:**
[SimilarSites](https://www.similarsites.com)
 
**Deprioritized:** Zanran (dead — domain now hosts an unrelated B2B product), Zapmeta (defunct since ~2003, domain blocks automated access), Mamont/mmnt.ru (technically live but Russian-only FTP index, too niche for routine use), 2lingual Search (functional but zero investigative utility, just a bilingual Google wrapper), unindexed basic Google-scraping scripts without captcha-solving or proxy pools; legacy portals: Aol, Goodsearch, Instya, Lycos, Search.com; general and AI wrappers: YOU, Phind (developer-focused), Wolfram Alpha (a computational engine, not a web index), Startpage (Google results through a proxy, redundant); metasearch wrappers: All-in-One, AllTheInternet, Etools, FaganFinder, Goofram, iZito, Myallsearch, Qwant (relies on Bing), Swisscows; privacy engines kept only for completeness: Disconnect Search, Gibiru, Presearch; dork generators and scrapers: AtDork, DorkCraft, DorkGenius, DuckDuckGo URL scraper (unvetted), SerpApi (paid SERP API), SerpScan (PHP dorking CLI); SitesLike (redundant next to SimilarSites)
 
---
 
### 17. Investigative Document & Records Platforms
 
- [Google Pinpoint](https://journaliststudio.google.com/pinpoint/) — parses massive PDF/audio/email/handwritten-note troves
- [DocumentCloud](https://www.documentcloud.org) — annotate, publish, search leaked/public-record documents
- [OCCRP Aleph](https://aleph.occrp.org) — global leaked/public corporate records, court filings, sanctions lists
- [ICIJ Offshore Leaks Database](https://offshoreleaks.icij.org) — offshore entities, officers and intermediaries from the ICIJ leak investigations (cross-referenced in §8)
- [RECAP / CourtListener](https://www.courtlistener.com/recap/) — federal/state US court document archive
- [Judyrecords](https://www.judyrecords.com/) — free full-text search across 400M+ US court cases
- [Internet Archive](https://archive.org) — the parent project behind Wayback Machine (full web-snapshot tool cross-referenced in §7); also hosts books, media, software

**Deprioritized:** De-facto (closed/absorbed), WorldWideScience.org (sunset — homepage now speaks of founding partners in the past tense), Biznar (defunct since ~2015, Deep Web Technologies wound down its federated-search line), CiteSeerX (currently broken — domain serves a placeholder instead of the actual index), Harmari Unified Listings Search (pivoted into a paid product rebranded "Neumo," no longer functions as a general search tool); UniCourt (limited free searches, premium upsell), Caselaw Access Project (historical, not up to date); Find-pdf-doc, Free Full PDF, Scribd, SlideShare (low-signal hosts and aggregators); Epstein Exposed (third-party aggregator — verify against the primary documents; unvetted); California Justice Watch (regional, unvetted); ArkhamMirror (local-first AI document analysis with knowledge graphs — unvetted; not the same product as Arkham Intelligence in §12)
 
---
 
### 18. Personal OPSEC, Privacy Tooling & Anonymity
 
**Fingerprint / leak testing (test yourself before you test anyone else):**
[BrowserLeaks](https://browserleaks.com), [packet.guru](https://packet.guru) (WebRTC/DNS leak test, TLS/TCP fingerprinting, VPN/proxy/Tor detection with a trust score, no signup), [Tor Browser](https://www.torproject.org), [Cover Your Tracks](https://coveryourtracks.eff.org) (EFF's fingerprint test, successor to Panopticlick)

**Privacy-respecting daily-driver apps:**
```
VPN              -> Proton VPN or Mullvad VPN
Browser          -> Mullvad Browser or Firefox / Brave
Search engine    -> Kagi, SearXNG, or DuckDuckGo
Email            -> Proton Mail or Tuta Mail
Messaging        -> Threema or Signal
Cloud storage    -> Tresorit or Proton Drive
Antivirus        -> Bitdefender, ClamAV, or Sophos
Password manager -> 1Password or Bitwarden
Note-taking      -> Obsidian or Joplin
```
**You can always find more pricavy based apps on** -> [privacytools](https://privacytools.io/)

**Investigator hygiene (isolation and disposable identities):**
[Tails](https://tails.boum.org) and [Qubes OS](https://www.qubes-os.org) (isolation), [uBlock Origin](https://github.com/gorhill/uBlock) and [Privacy Badger](https://www.eff.org/privacybadger) (tracker blocking), [Guerrilla Mail](https://www.guerrillamail.com) (burner inbox), [Privacy.com](https://privacy.com) (virtual payment cards), [KeePass](https://keepass.info) (offline password manager), [GNU PG](https://www.gnupg.org/download/index.html) (key handling and encryption)

**Threat modelling and protecting case files at rest:**
[EFF Surveillance Self-Defense](https://ssd.eff.org) (plain-language guides to threat modelling and picking tools), [FileVault](https://support.apple.com/en-us/HT204837) (macOS built-in full-disk encryption — keep case notes encrypted at rest)

**Hardened browser forks and VPN comparison:**
[LibreWolf](https://librewolf.net), [Waterfox](https://www.waterfox.net/), [VPN Comparison by That One Privacy Guy](https://thatoneprivacysite.net/), [OffShore.cat VPN list](https://offshore.cat/?page=vpn)
 
**Deprioritized:** Abandoned public proxy lists; browsers kept only for completeness: Atom, Bromite (discontinued), CentBrowser, Comodo Dragon, Coowon, Maxthon, Sleipnir, Slimjet, SRWare Iron, Torch Browser (a different product from the Torch onion index in §11), UCBrowser, Yandex Browser, Epic Privacy Browser (status unconfirmed), GNU IceCat (niche), plus the mainstream Chrome, Edge, Opera, Safari and Vivaldi (fine for daily use, not privacy tooling); retired, replaced or compromised tools: HTTPS Everywhere (retired, browsers now default to HTTPS), Panopticlick (replaced by Cover Your Tracks), Wickr (consumer app retired), Boxcryptor (discontinued), OpenPGP / Enigmail (Thunderbird now has OpenPGP built in), Do Not Track (never enforced), LastPass (2022 breach; Bitwarden or KeePass preferred), CCleaner (2017 supply-chain compromise, telemetry controversy), WOT (2016 browsing-data scandal); overlapping or older utilities (status unconfirmed): Abine, Adium, ChatSecure, Disconnect, EncFS MP, Eraser, Ghostery, GPG Tools, Guardian Project, Hotspot Shield, justdelete.me, Lockbin, Mailvelope, Master Password, Nixory, NoScript, Open DNS, Oscobo, OSSEC, Peerblock, Pidgin, Pixel Block, Privazer, Script Safe, Securesha, Silent Circle, Snort, SpiderOak, Steganography Online Codec, ZMail, MetaClean (unvetted); added in the latest merge: I2P (alternative anonymity network, niche — Tor covers the common case), mailbox.org (privacy-friendly mail host that overlaps Proton Mail and Tuta in the table above), TorrentFreak List of VPNs (2017-era review, status unconfirmed — use the comparison lists above)

---
### 19. Countries Where Google Maps / Street View Has Significant Restrictions

**Fully or largely absent from Street View:**
- **North Korea** — no Street View at all, satellite imagery is deliberately degraded by the regime, Maps shows almost no road detail
- **China** — Google Maps and Street View are blocked entirely inside China (use **[Baidu Maps](https://map.baidu.com/)** or **[Amap/Gaode](https://amap.com/)** instead — far better coverage for China than anything Google has)
- **Saudi Arabia** — Street View only recently became partially available and coverage is still very patchy outside major cities
- **Iran** — Google Maps accessible but Street View coverage is almost nonexistent, satellite imagery is limited in detail
- **Cuba** — very sparse Street View, road data is poor
- **Bhutan** — almost no Street View, government restricts most mapping
- **Turkmenistan** — effectively no Street View coverage

**Partial/restricted Street View for security reasons:**
- **Israel** — by law (Mapping and Databases Law), Google must blur military sites and certain government buildings before publishing imagery; coverage exists but with deliberate gaps
- **India** — Street View was banned for years over national security concerns, only returned in 2022 with restrictions near border regions and military installations
- **Russia** — Street View still exists but Google stopped updating it after 2022 and the existing imagery is aging fast; use **[Yandex Maps](https://yandex.com/maps/)** (specifically **[Yandex Panorama](https://yandex.com/maps/)**) which has far better and more current Russian Street View coverage
- **Ukraine** — Google paused Street View collection in active conflict zones
- **Pakistan** — Street View coverage is extremely sparse outside Islamabad and Lahore

**Countries where Google Maps exists but local alternatives are far better:**
- **China** $\rightarrow$ **[Baidu Maps](https://map.baidu.com/)** / **[Amap (Gaode)](https://amap.com/)** — Google has no usable China data
- **Russia** $\rightarrow$ **[Yandex Maps](https://yandex.com/maps/)** — much fresher and denser Street View than Google
- **Japan** — Google is good but **[Mapion](https://www.mapion.co.jp/)** and **[Yahoo Japan Maps](https://map.yahoo.co.jp/)** have detail Google misses at hyper-local level
- **South Korea** $\rightarrow$ **[Naver Maps](https://map.naver.com/)** / **[Kakao Maps](https://map.kakao.com/)** — Korean law restricts Google from having full-resolution map data, so local apps are far superior
- **Czech Republic** $\rightarrow$ **[Mapy.cz](https://en.mapy.cz/)** — significantly better hiking/rural coverage
- **Germany** — Google Street View coverage is intentionally very sparse because a large number of German residents legally opted out when Street View launched (right to be forgotten requests blurred entire streets); **[OpenStreetMap](https://www.openstreetmap.org/)** + **[Mapillary](https://www.mapillary.com/)** often have better German ground-level coverage

#### The Practical Implications for Geo-OSINT

> If your image looks like it could be China and Google Maps shows you nothing useful $\rightarrow$ switch to **[Baidu Maps](https://map.baidu.com/)** immediately, it is not a fallback, it is the primary tool for China

> If your image looks Russian and Street View imagery looks old or is missing $\rightarrow$ **[Yandex Maps Street View (Panorama)](https://yandex.com/maps/)** is what you actually want, it has coverage Google does not

> If your image looks South Korean $\rightarrow$ **[Naver Maps](https://map.naver.com/)** or **[Kakao Maps](https://map.kakao.com/)** will show you ground-level imagery in areas Google cannot legally provide at full resolution

> If your image is in a military-adjacent area in Israel, India, or Pakistan $\rightarrow$ the deliberate blur / absence on Google is itself a clue that you are near a sensitive site

> If there is simply no Street View anywhere near your candidate location $\rightarrow$ **[Mapillary](https://www.mapillary.com/)** is community-contributed ground-level imagery and covers many areas Google has never sent a car to

#### Quick Reference — What to Use Where

| Location | Primary Map Tool | Street View Equivalent |
| :--- | :--- | :--- |
| **China** | **[Baidu Maps](https://map.baidu.com/)** / **[Amap](https://amap.com/)** | Baidu Street View |
| **Russia** | **[Yandex Maps](https://yandex.com/maps/)** | Yandex Panorama |
| **South Korea** | **[Naver Maps](https://map.naver.com/)** / **[Kakao](https://map.kakao.com/)** | Naver Street View |
| **North Korea** | None useful | Literally nothing |
| **Germany (sparse)** | **[OpenStreetMap](https://www.openstreetmap.org/)** | **[Mapillary](https://www.mapillary.com/)** |
| **Japan (hyper-local)** | **[Yahoo Japan Maps](https://map.yahoo.co.jp/)** | Google still usable |
| **Global gaps** | **[Mapillary](https://www.mapillary.com/)** | Community-contributed imagery |
| **Conflict zones** | **[LiveUAMap](https://liveuamap.com/)** + **[OSM](https://www.openstreetmap.org/)** | Older cached Google imagery |

> **The bottom line for geo-OSINT:** Google Maps is not a single global truth — it is a patchwork of varying coverage, legal restrictions, deliberate blurring, and aging imagery. Knowing which country you are looking at (which your Stages 1–4 already narrow down) tells you immediately which mapping tool is actually authoritative for that region.
---
### 20. Social Media & Video Geolocation Search Tools + advanced social media search engines 
(Focused specifically on discovering user-generated media, location-tagged uploads, and spatial tracking across platforms during an investigation)

`advanced search engines`
* Honorable mention here btw => [Awesome-Telegram-OSINT](https://github.com/ItIsMeCall911/Awesome-Telegram-OSINT)
* **[Getdaytrends.com](https://getdaytrends.com/)** : Check the active social media trends — also **[Trends24](https://trends24.in)** for trends by location

* **Instagram**
  * **[Toutatis](https://github.com/megadose/toutatis)**: Python tool designed to extract hidden profile details from target Instagram accounts, including linked phone numbers, masked emails, and internal user IDs.
  * **[Osintgram](https://github.com/Datalux/Osintgram)**: Interactive CLI tool that parses an Instagram profile to automatically map followers, target captions, utilized hashtags, and geotagged locations.
  * **[insto](https://github.com/subzeroid/insto)**: Interactive OSINT CLI with 35+ slash-commands (profile, media, followers, shared-followers intersection, location search, posting-cadence histogram, Maltego CSV export); token-based, so no Instagram account is needed.
* **Twitter / X**
  * **[Twint](https://github.com/twintproject/twint)**: Open-source intelligence tool that scrapes tweets, lists, and user metadata without requiring the official Twitter/X API or authentication keys, bypassing standard rate limits.
  * **[Twitonomy](https://www.twitonomy.com/)**: Web-analytics platform providing granular visual stats on a user's posting habits, active hours, retweet networks, and geographical context.
* **LinkedIn**
  * **[CrossLinked](https://github.com/m8sec/CrossLinked)**: Reconnaissance tool designed to harvest valid employee names, structural titles, and email permutations from LinkedIn via search engine scraping without needing an authenticated user account.
* **GitHub**
  * **[GitRecon](https://github.com/atiilla/gitrecon)**: Advanced OSINT and reconnaissance tool designed to scan GitHub and GitLab user profiles, repositories, and organizations for exposed personal email addresses, commit patterns, SSH keys, and internal metadata.
  * **GitHub Advanced Search & API Dorking**: Leveraging targeted query operators (`filename:`, `path:`, `extension:`, `org:`, and `repo:`) combined with scripts like `gitrob` or `gitleaks` to unearth hardcoded secrets, personal naming conventions, and developer footprints.
  * **[GHNames](https://ghnames.com/)**: Username-history lookup that also surfaces associated emails and deleted repositories.
  * **[GitLeak](https://gitleak.io)**: Pulls hidden commit emails, commit timezones, and activity patterns out of `.patch` metadata (a different tool from the Gitleaks secret scanner in §10).
* **Reddit**
  * **[ThinkPol](https://think-pol.com/)**: Professional-grade archival search engine providing full-archive historical text search, covering user deletions, moderator removals, and pre-edit versions.
  * **[Reveddit](https://www.reveddit.com/)**: Transparency utility used to track and audit content silently removed or filtered by moderators on target user profiles or subreddits.
  * **[Arctic Shift](https://github.com/ArthurHeitmann/arctic_shift)** and **[PullPush](https://pullpush.io/)**: Pushshift-style archives of historical Reddit data, useful for deleted or removed posts and comments.
  * **[RedditMetis](https://redditmetis.com/)**: Account summary and activity statistics for a Reddit user, including top posts and posting patterns.
* **Discord**
  * **Snowflake ID / Metadata Extractors**: Specialized lookup utilities (such as Discord developer utilities and historical guild mappers) used to extract account creation timelines, public metadata, and server association patterns via unique user IDs.
  * **Server & Message Indexers**: Third-party public chat archives and server directories (such as DiscordGate) used to track public user activity, historical message presence, and mutual guild memberships across open ecosystems.
* **Telegram**
  * **[Tgstat](https://tgstat.com/) / Livegram**: Dedicated platform search engines and indexes used to track channel growth metrics, message history logs, and cross-reference public group mentions.
  * **[Telerecon](https://github.com/sockysec/Telerecon)** and **[Telepathy](https://github.com/proseltd/Telepathy-Community)**: Investigation frameworks that archive Telegram channels and analyze communication patterns, message frequency, and user interactions.
  * **[TOsint](https://github.com/drego85/tosint)** and **[Maltego Telegram](https://github.com/vognik/maltego-telegram)**: Extract information from Telegram bots and channels, and bring Telegram entities into Maltego graphs.
  * **[Telegago](https://cse.google.com/cse?q=+&cx=006368593537057042503:efxu7xprihg#gsc.tab=0&gsc.q=%20&gsc.page=1)**: Google-powered advanced search built for finding public Telegram channels and chatrooms.
  * **Utility bots:** [username_to_id_bot](https://t.me/username_to_id_bot) and [UsInfoBot](https://t.me/usinfobot) (username ↔ ID), [creationdatebot](https://t.me/creationdatebot) (approximate account creation date), [SangMata](https://t.me/SangMata_beta_bot) and [Unamer](https://telegram.me/unamer_bot) (name and username history), [Maigret OSINT bot](https://t.me/osint_maigret_bot) (username search), [WhoisDomBot](https://t.me/WhoisDomBot) (whois/dig), [FindStickerCreator](https://t.me/SPOwnerBot) (owner of a sticker pack).
  * **More bots, picked for the gaps they fill:** [GeoMacFinder](https://t.me/geomacbot) (Wi-Fi access point location by MAC/BSSID — pairs with WiGLE below), [EyeTON](https://telegram.me/istoneyebot) (TON wallet graph with linked profiles) and [SpyGGbot](https://telegram.me/SpyGGbot) (TON balances, NFT owners, Fragment usernames) — crypto tracing for Telegram-linked wallets, see §12, [Discord Sensor](https://telegram.me/discordsensorbot) (Discord account data), [ChatSearchRobot](https://t.me/ChatSearchRobot) and [Searchforchats](https://telegram.me/searchforchatsbot) (find chats by topic — feeds the network-building routine in the note below), [AvtoNomer](https://t.me/avtonomerbot) (vehicle photos by plate via platesmania.com — pairs with §15), [AgentFNS_Bot](https://t.me/AgentFNS_bot) and [egrul_bot](https://t.me/egrul_bot) (free counterparty checks from official Russian registry data — pairs with §8)
`Telegram is and will stay one of the best tools for threat actor related osint, find groups using bots, then join those groups to get links for more groups -> build you entire telegram networks untill you get the deep down part of Telegram that cannot be reached directly , that is where the main information lies` 

* **Facebook**
  * **[Lookup-ID.com](https://lookup-id.com)**: Resolves a Facebook profile, group, or page URL to its numeric ID, which most other Facebook pivots need.
* **VKontakte**
  * **[VK People Search](https://vk.com/people)** and **[VK Community Search](https://vk.com/communities)**: The native filter-rich search is the strongest tool on VK; pair it with **[Search4faces](https://search4faces.com/)** for face search across Russian-language networks.
* **WhatsApp**
  * **[2Chat WhatsApp checker](https://2chat.co/tools/whatsapp-checker)**: Checks whether a phone number is registered on WhatsApp.
  * **[Groupio](https://en.groupio.app/)**: Finds and searches public WhatsApp groups.
  * **[WhatsApp CheckLeaked](https://whatsapp.checkleaked.cc/)**: WhatsApp number search and profile-photo checker from the same project as CheckLeaked in §10 (API available).
* **Steam**
  * **[OSINT-Steam](https://osint-steam.vercel.app/en)**: Returns public profile information, such as the friends list and possible locations, from Steam profiles.
  * **[Steam-OSINT](https://github.com/matiash26/steam-osint)**: Mutual-friends analysis across full friend lists.

Deprioritized: Standard native platform search bars (natively scrub deep historical metadata, restrict advanced boolean queries, and impose aggressive rate-limits on unauthenticated enumeration)
* **Instagram:** Iconosquare (marketing analytics), Dolphin Radar, IGDetective, Osintgraph (viewers and graph tools — unvetted); instagram_monitor (real-time activity and profile-change tracker with email alerts — unvetted)
* **Twitter / X:** Foller.me, TwChat, OneMillionTweetMap, TweetMap (two variants), Sentiment140, Twitter Search (search.twitter.com, retired) (status unconfirmed); MyTweetAlerts, RiteTag, Tagdef, Twitter Audit, Twitter Chat Schedule (marketing utilities); ExportData and Xquik (paid data platforms — unvetted)
* **LinkedIn:** the-endorser (older, endorsement-graph only), LinkedInDumper (needs an authenticated session and carries terms-of-service risk, CrossLinked covers the unauthenticated route)
* **GitHub:** Glit and Zen (overlap GitRecon), GithubRecon (kriztalz), Related Repos, Shotstars (fake-star analysis, niche); github_monitor (real-time profile and repository change tracker — unvetted)
* **Reddit:** Pushshift API (public access restricted since 2023), Mostly Harmless, Reddit Suite (the old enhancement extension), Imgur search, Subreddits.org, Reddit Archive, REDARCS (static archive up to 2023), Reddit Comment Search, Reddit Comment Lookup, Reddit User Analyser, Universal Scammer List (a moderation list, not an investigation tool)
* **Telegram:** Telegram Channel Joiner (account-growth script, not investigative), GroupDa, TgramSearch, tg.world, TeleSearch, Teleteg (10 results on the free plan), TeleTracker, TeleGraphite, Telegram Finder (unvetted); remaining bots from the bulk import (all unvetted): duplicates of the utility bots above — RegDateBot, SangMataInfo_bot, TuriBot; Russian phone-operator checkers that overlap the carrier lookups in §6 — bmi_np_bot, mnp_bot; vehicle-history bots that are paid or regionally limited (AvtoNomer covers the free public-photo route) — AVinfoBot, avtogram_bot, avtocodbot, EasyVIN, NEUROAUTOSEARCH, Фари; wrappers around tools already above — PimEyes bot (use PimEyes, §4), datXpert (IntelX, §10), Sherlock (Getcontact-based caller-ID lookup), SaveYoutubeBot (yt-dlp, §25), InstaBot, InstaAnonym, MotherSearchBot, OkSearch, Surftg_bot; non-investigative utilities — CryptoBot (anonymous wallet), dCallsBot (anonymous calls, eSIM), SMS Activate (virtual numbers for burner accounts — check each platform's terms before using), getChatList, GetSendGifts, Insight; regional lookups — OpenDataUABot (Ukrainian open data), Search_firm_bot (organizations, banks, postal codes), ShtrafKZBot (Kazakhstan fines and taxes); OPENLOAD Bot (semi-automated scanning suite, authorised engagements only); leak-backed people, phone and password lookups, where data provenance is unknown and legality varies by jurisdiction — ClerkBot, Detectiva, HimeraSearch, ITP Infotrack, Leak OSINT, Osintkit, Getairplane, PasswordSearch
* **Facebook:** Wolfram Alpha Facebook Report, Facebook Search (search.fb.com), SearchIsBack, OsintStalker (status unconfirmed); Fanpage Karma (marketing); Find my Facebook ID (randomtools); Facebook Friend List Scraper (friend-list scraping carries terms-of-service risk — unvetted), fb-sleep-stats (infers sleep habits from online status, looks old — status unconfirmed)
* **VKontakte:** VK.watch, VK5, Barkov.net, the Дезертир VK app (status unconfirmed or unvetted)
* **TikTok, Pinterest, Tumblr:** TikTokStalker (thin third-party profile viewer, unstable), Pingroupie, Pinterest Pin Stats (marketing stats), Tumblr Search (native search bar)

`geo location`
* **[YouTube Geofind (MW Geofind)](https://mattw.io/youtube-geofind/location)**: The direct counterpart enabling radius-based searches and bulk metadata exports for geotagged YouTube videos.
* **[Echosec](https://www.echosec.net/)**: Commercial geospatial OSINT platform that aggregates social media posts (X/Twitter, YouTube, VK, and others) by drawing a bounding box or radius over a map interface.
* **Scylla / Social Media Map Search Tools**: Specialized scripts and wrappers that query platform APIs to filter content based on coordinate bounding limits or embedded location tags.
* **Telegram Geo-Scrapers**: Custom scripts and open-source tools designed to scrape geotagged media and channel posts originating from a precise coordinate point or radius.
* **[Snapchat Map](https://map.snapchat.com/)**: Public web interface for viewing localized, public user snaps across the globe in real time, frequently used to track crowd movements or localized events.
* **Twitter/X Advanced Search Operators**: Utilizing manual string parameters (`geocode:lat,long,radius`) via the search bar to surface tweets containing native GPS coordinates or location markers.
* **[WiGLE](https://wigle.net/)**: Crowdsourced Wi-Fi network map; resolves an SSID or BSSID to a location (general mapping tools are in §26).

Deprioritized: Native platform search bars (natively scrub or heavily obscure granular coordinate data without specialized third-party filtering layers)

---
## New categories (added from the Awesome-OSINT merge)

---

### 21. Threat Actor & Cybercrime Profiling

*(Complements the one-line profiling entry in §9 and the Telegram note in §20)*

- Group and actor encyclopedias: [MITRE ATT&CK Groups](https://attack.mitre.org/groups/), [Malpedia](https://malpedia.caad.fkie.fraunhofer.de/actors) (actors plus malware families), [ETDA Threat Group Cards](https://apt.etda.or.th/cgi-bin/listgroups.cgi), [MISP Galaxy](https://www.misp-galaxy.org/360net/) (machine-readable actor clusters)
- Alias and handle pivoting: [Threat Actor Usernames Scrape](https://github.com/spmedia/Threat-Actor-Usernames-Scrape) (900k+ usernames from cybercrime sources and forums), [threatactorusernames.com](https://threatactorusernames.com/) (3M+ records) — confirm any handle match against the original post before trusting it
- Cross-vendor alias mapping: [APT Groups and Operations](https://docs.google.com/spreadsheets/u/0/d/1H9_xaxQHpWaa4O_Son4Gx0YOIzlcBWMsdvePFX68EKU/pubhtml?pli=1#) (public spreadsheet of actors, sponsors, tools and methods)

**Deprioritized:** vendor-run encyclopedias with narrower coverage than the three above: FortiGuard Threat Actor Encyclopedia, SOCRadar Labs, BreachHQ, Dark Web Informer, Bi.Zone (148 groups), KNOWLEDGENOW (Netenrich), Thales graphical attack explorer; lazarusholic, OPENHUNTING.IO (unvetted); STIX Viewer (a viewer for STIX 2.1 files, useful only once you already have STIX data)

---

### 22. Live Cyber Threat Maps

*(Situational awareness only — these show each vendor's own telemetry, not individual incidents, so there is nothing to pivot from)*

- [Check Point Live Cyber Threat Map](https://threatmap.checkpoint.com/), [Kaspersky Cyberthreat Live Map](https://cybermap.kaspersky.com/), [Radware Live Threat Map](https://livethreatmap.radware.com/) — Cloudflare Radar is in §9

**Deprioritized:** Bitdefender Threat Map, BunkerWeb, Cisco Talos (email and spam view), FortiGuard Labs, HCL Threat Map, IBM X-Force Exchange activity map, Imperva, LIONIC, NETSCOUT, Secure Gateway, Thales Cyberthreat Map, ThreatsEye, World Monitor Tech, Zscaler (all single-vendor telemetry; good for visuals and demos, no investigative pivots)

---

### 23. File Search, Open Directories & Pastebins

- Open-directory search: [ODCrawler](https://odcrawler.xyz/), [eyedex](https://www.eyedex.org/), [NAPALM FTP Indexer](https://www.searchftps.net/) (FTP specifically)
- Reference: [Filesec.io](https://filesec.io/) (catalog of malicious file extensions and their risks)
- Pastes people actually use and that get indexed: [GitHub Gist](https://gist.github.com), [Rentry](https://rentry.co/) — for older or deleted pastes use Intelligence X (§10) and Google dorks rather than the hosting sites themselves

**Deprioritized:** file search: de digger (Google Drive files), Meawfy (Mega.nz files), SearchFiles.de, Find Security Contacts (security.txt index) (all unvetted); pastebin hosting sites, which are rarely indexed and mostly ephemeral — BeanPaste, bpaste, CentOS Pastebin, cl1p, commie, Context, ControlC, Cutapaste, dpaste, Etusivu (pastebin.fi), Friendpaste, hastebin, ideone, ivpaste, jsbin, justpaste, Katbin, Linkode, Logpasta, lesma.eu, NoPaste, nopaste.net, Notes.io, nekobin, New Paste, n0paste, Pastebin.cz, paste.debian, paste.in.ua, paste.kde, Paste.Monster, paste.ubuntu, Paste.Quest, Pastery, PasteSite.Net, PasteShr, snippet.host, Spacebin, TextBin, Textbin-Code, TutPaste, Write.as; pastebins that advertise client-side encryption, so their contents cannot be searched by design — Defuse, HashBin, Nachricht, paaster, Paste.Cash, paste.sh, SafeNote, scrt.link, vaultbin, Verybin, ZBin, ZeroBin

---

### 24. Forums, Blogs & Q&A Search

- Forum search: [Boardreader](https://boardreader.com), [Google Groups](https://groups.google.com) (long-running Usenet and mailing-list archive)
- Blog search: [Twingly](https://www.twingly.com)
- Q&A with profile history: [Quora](https://www.quora.com), [StackExchange](https://stackexchange.com), [Ответы Mail.ru](https://otvet.mail.ru/) (Russian)

**Deprioritized:** Yahoo Groups (shut down), Yahoo Answers (shut down), Answers.com, Ask, eHow (SEO-driven answer farms); BlogSearchEngine, Notey (status unconfirmed); Built With Flarum (lists Flarum forums only), Ning (most communities are walled off), Facebook Groups and LinkedIn Groups (login-walled, see §13), 4chan Search (unvetted)

---

### 25. Video & Livestream Research

- YouTube pivots: [Filmot](https://filmot.com/) (search inside subtitles), [YouTube Metadata](https://mattw.io/youtube-metadata/), [YouTube Geofind](https://mattw.io/youtube-geofind/) (see §20 for the radius search), [Find YouTube Video](https://findyoutubevideo.thetechrobo.ca/) (deleted or private videos by ID across archives), [Archivarix Tube Search](https://tube.archivarix.net/)
- Download and preserve: [yt-dlp](https://github.com/yt-dlp/yt-dlp/)
- Frame-level review: [Frame by Frame](https://chrome.google.com/webstore/detail/frame-by-frame/cclnaabdfgnehogonpeddbgejclcjneh/) (browser plugin for YouTube)
- Public tourist webcams: [Earthcam](https://www.earthcam.com)

**Deprioritized:** Bing Videos, Yahoo Video Search, DailyMotion, Vimeo, YouTube (native search only); Metacafe, Metatube, Veoh (status unconfirmed); Deturl, DownloadHelper (redundant next to yt-dlp); Tubeup (archives to Internet Archive via yt-dlp, a wrapper); Geosearch (superseded by YouTube Geofind); AudienceCue (comment-export utility, unvetted); Video Stabilization Methods (a research list); Internet Archive Open Source Videos (rarely relevant); Insecam (directory of unsecured IP cameras, many of them private — ethically and legally risky; public webcams like Earthcam are the safer pivot; unvetted)

---

### 26. Geospatial Research, Mapping & Satellite Imagery

*(Companion to §4 and §19: those cover image-led geolocation and per-country map coverage, this is the general mapping toolbox)*

- Satellite and aerial imagery: [Google Earth Pro](https://www.google.com/intl/en/earth/versions/#earth-pro) (free desktop, historical imagery), [Sentinel Hub Playground](https://www.sentinel-hub.com/explore/sentinelplayground/) (recent Sentinel-2 passes), [Zoom Earth](https://zoom.earth/) (near-real-time satellite and weather), [USGS EarthExplorer](https://earthexplorer.usgs.gov/) (free government imagery archive: Landsat, aerial photography, declassified satellite data)
- Open map data and desktop GIS: [Open Street Map](https://www.openstreetmap.org), [QGIS](https://qgis.org), [SAS Planet](https://www.sasgis.org/sasplaneta/) (view, download and stitch satellite imagery)
- Wi-Fi and signal geolocation: [WiGLE](https://wigle.net/)
- Conflict and event mapping: [Liveuamap](https://liveuamap.com/)
- Conditions check (cloud cover, weather at a given time): [Windy](https://www.windy.com/)
- Where do I look in this country: [Worldwide OSINT Tools Map](https://cipher387.github.io/osintmap/) (databases and OSINT sources by location)
- Measuring: [MapChecking](https://www.mapchecking.com/) (distance and area)

**Deprioritized:** mapping libraries and cartography services, for building maps rather than investigating: ArcGIS Living Atlas, Atlas, Atlasify, Batchgeo, Bing Maps, CartoDB, Colorbrewer, digiKam, Dominoc925, DualMaps, Esri, GeoNames, Google My Maps, GPSVisualizer, GrassGIS, Here, Inspire Geoportal, InstantAtlas, Kartograph, Leaflet, Map Maker, MapAList, MapBox, Mapchart.net, MapHub, MapJam, Mapline, Mapquest, NGA GEOINT, OpenLayers, Perry-Castañeda Library maps, Scribble Maps, StoryMaps, Tableau, ViaMichelin, View in Google Earth, Zeemaps; duplicates of tools above: Instant Google Street View, KartaView, Mapillary (both in §4), Google Maps, Wikimapia, Pic2Map (covered by Exiftool in §4), SOAR, Satellites Pro; Apify's Google Maps Scraper (paid scraping); SatIntel (satellite-tracking CLI, niche); Flash Earth, Hyperlapse, Polymaps, Modest Maps, CrowdMap, Maperitive, GeoGig, QuickMaps, Yahoo Maps, WorldMap Harvard (older, status unconfirmed); conflict and geopolitical dashboards: World Monitor, Pharos AI, Phantom Tide (unvetted)

---

### 27. Web Archiving, Capture & Evidence Preservation

*(Capture first, analyze second: pages get edited and deleted)*

- Case-aware capture: [Hunchly](https://www.hunch.ly/) (paid, built for online investigations), [archive.today](https://archive.is) (permanent snapshot that survives the original page)
- Pull from the Wayback Machine: [waybackpy](https://github.com/akamhy/waybackpy) (Python and CLI), [waybackurls](https://github.com/tomnomnom/waybackurls) (every known URL for a domain), [WayTrace](https://waytrace.org/) (rebuilds a domain's history: emails, subdomains, analytics IDs, with first and last seen dates)
- Mirror a site for offline review: [HTTrack](https://www.httrack.com), [Cyotek WebCopy](https://www.cyotek.com/cyotek-webcopy)

**Deprioritized:** BlackWidow, CachedPages, CachedView, stored.website (thin cache viewers); Wayback Machine Archiver / waybackpack (redundant next to waybackpy); Wayback-Archive (GeiserX), Website-Diff (unvetted); A1 Website Download, Website Ripper Copier, WebAssistant (paid or older site rippers); SiteSucker (macOS only); Offliberty, Resolver, gmapcatcher, Hooey webprint (status unconfirmed)

---

### 28. Monitoring, Alerts & Keyword/Trend Research

- Page-change detection: [ChangeDetection.io](https://github.com/dgtlmoon/changedetection.io) (self-hosted, open source), [Visualping](https://visualping.io)
- Alerts and feeds: [Google Alerts](https://www.google.com/alerts), [Feedly](https://www.feedly.com), [RSS Bridge](https://bridge.suumitsu.eu) (generates feeds for sites that have none)
- Search-interest and trend research: [Google Trends](https://www.google.com/trends), [Yandex Wordstat](https://wordstat.yandex.com) (Russian-language search demand), [Soovle](https://www.soovle.com) (autocomplete across engines)

**Deprioritized:** SEO keyword tools aimed at marketers: Google Adwords (now folded into Google Ads), Keyword Discovery, Keyword Spy, KeywordTool, One Look Reverse Dictionary, Ubersuggest, Word Tracker; paid monitoring: Mention, Talkwalker, Silobreaker, versionista; duplicate or stale RSS readers (status unconfirmed): Alltop, Awasu, Bridge.Leslibres, DiggReader, FeedBooster, Feed Exileed, Feed Filter Maker, FeedReader, FetchRSS, Flipboard, InfoMinder, Netvibes, Newsblur, OmeaReader, Reeder, RSS Feed Reader, RSS Search Engine, RSS Search Hub, RSSOwl, Selfoss, The Old Reader, WebReader, Winds, Deltafeed; duplicate change-detection tools: ChangeDetection (changedetection.com), FollowThatPage, OnWebChange, WebSite Watcher

---

### 29. News, Media Monitoring & Fact-Checking

- Aggregation and discovery: [Google News](https://news.google.com), [ABYZ News Links](https://www.abyznewslinks.com) (directory of newspapers worldwide), [Newspaper Map](https://newspapermap.com)
- Wire services: [Reuters](https://www.reuters.com), [AP](https://hosted.ap.org)
- Fact-checking: [Snopes](https://www.snopes.com), [Full Fact](https://fullfact.org), [FactCheck.org](https://www.factcheck.org), [Check by Meedan](https://meedan.com/check)
- Verification training: [Verification Handbook](https://verificationhandbook.com)

**Deprioritized:** single-outlet sites better read directly: BBC News, CNN, Euronews, France24, DPA International, Agence France-Presse, Itar-Tass (now TASS, a state outlet), Yahoo News, Bing News; thin aggregators and directories: 1st Headlines, AllYouCanRead, DailyEarth, HeadlineSpot, MagPortal, News Map, News Now, NewsLookup, Newspaperindex, OnlineNewspapers, Paperboy, Topix, WorldNews, World-Newspapers, Google News Print Archive; paid: Factiva, Newspapers.com, NewsWhip, Press Reader, Silobreaker, Storyful; press-release wire: PR Newswire; Newseum Today's Front Pages (the museum closed in 2019, status unconfirmed); APITube (unvetted); digest apps: Flipboard, Inshorts, Newsinshorts, Nod, Reeder, Spike, Superdesk (a newsroom CMS), Trooclick (status unconfirmed); fact-checking: Captain Fact, Emergent (status unconfirmed)

---

### 30. Academic & Grey Literature

- Search: [Google Scholar](https://scholar.google.com), [CORE](https://core.ac.uk/search), [BASE](https://www.base-search.net), [PubMed](https://www.ncbi.nlm.nih.gov/pubmed), [OA.mg](https://oa.mg) (240M+ works, PDFs for open-access papers)
- Regional: [African Journals Online](https://www.ajol.info)

**Deprioritized:** retired or stale: Microsoft Academic (retired 2021), OpenGrey, World Digital Library, Zetoc, HighWire free-article list (status unconfirmed); login-walled or noisy: Academia, ResearchGate, JSTOR (mostly paywalled); publisher portals, search their content through Google Scholar instead: Elsevier, ScienceDirect, Springer, Wiley, Taylor & Francis Online, SAGE Journals, Oxford Journals, NRC Research Press, American Society of Civil Engineers; publishers often flagged in predatory-journal discussions, treat with caution: Academic Journals, Science Publications, ScienceDomain, SCIRP, Open Access Scientific Journals (PagePress); guides and directories rather than search tools: Grey Guide, Grey Literature Strategies, Grey Literature List of Gateways, GreyNet International, Journal Guide, Journal Seek, Bibsonomy, The Open Syllabus Project; Lazy Scholar, Leibniz TIB, Quetzal Search (niche)

---

### 31. Open Data & Statistics

- Core portals: [World Bank Data](https://data.worldbank.org), [UN Data](https://data.un.org), [OECD Data](https://data.oecd.org), [Eurostat](https://ec.europa.eu/eurostat), [CIA World Factbook](https://www.cia.gov/the-world-factbook/)
- Trade flows: [UN COMTRADE](https://comtrade.un.org)
- Conflict events: [Uppsala Conflict Data Program](https://www.pcr.uu.se/research/UCDP)
- Country context: [Gapminder](https://www.gapminder.org/data)

**Deprioritized:** discontinued or replaced: Freebase (shut down), Zanran (dead, see §16), World Bank Doing Business (discontinued 2021), Google Public Data Explorer (status unconfirmed); outdated links: ILO WESO 2015, Penn World Table 8.1, Transparency International CPI 2015 page; paywalled: Statista, Trading Economics, HIS Piers; overlapping UN and agency portals: UNCTAD (four entries), UNDP HDI, Human Development Reports, UNECE, UNESCO UIS, UNIDO, UNStats Social Indicators, UN agency duplicates of the portals above, WHO Data, WTO Statistics, WITS, ILOSTAT, IMF WEO Database, IEA Statistics, International Trade Center, International Labour Comparisons, OECD Aid Database, OECD Factbook, SESRIC (two entries), World Bank (four additional entries), CEPII, BIS Statistics, BP Statistical Review, Eurostat duplicates (EU Open Data Portal), Government of Canada Open Data, Data.gov.uk, US Data and Statistics, AGOA Data Center, AidData, CIESIN, LandMatrix, Latinobarometro, Pew Research Center, PRB Data Finder, PRS Risk Indicators, The Atlas of Economic Complexity; aggregators and mirrors: Index Mundi, NationMaster, Knoema, Junar, Open Data Network, Google Finance, DBPedia, AWS Public Datasets, Berkeley Data Lab, University of Michigan Library datasets, The Data and Story Library, Paul Hensel's data page, Vizala; globalEDGE DIBS; regional or newer (unvetted): Crime Brasil, DataInt Databook, UAV Radar, Voidly Censorship Index

---

### 32. Jobs, Professional & Expert Directories

*(Mostly useful for employer and role verification and for org-chart mapping — CrossLinked in §20 builds employee lists from LinkedIn)*

- Professional networks and job boards: [LinkedIn](https://www.linkedin.com), [Indeed](https://www.indeed.com), [Glassdoor](https://www.glassdoor.com), [Xing](https://www.xing.com) (Germany, Austria, Switzerland)
- Regional boards: [Naukri](https://www.naukri.com) (India), [Reed](https://www.reed.co.uk) (UK), [Seek](https://www.seek.com.au) (Australia), [Eluta](https://www.eluta.ca) (Canada)
- Professional rosters: [USPTO Patent Attorney/Agent Search](https://oedci.uspto.gov/OEDCI/)

**Deprioritized:** generic or legacy job boards: Beyond, CampusCareerCenter, CareerBuilder, College Recruiter, CVFox, Dice (tech only), Eurojobs, Fish4Jobs, Jobs.pl, Jobsite, Monster, SimplyHired, ZipRecruiter (status unconfirmed for the older ones); Craigslist (city-by-city only, no global search); RecruitEm (unvetted); expert-network and PR-matchmaking sites with low investigative yield: ExpertiseFinder, ExpertPages, Experts.com, HARO (status unconfirmed), GlobalExperts, Idealist, Innocentive, Internet Experts, Maven, MuckRack, National Speakers Association, Newswise, PRNewswire, ResearcherID, SheSource, Sources, TRExpertWitness, Zintro; Academia (see §30); Library of Congress Ask a Librarian (a help service, still worth knowing about)

---

### 33. Link Analysis, Visualization & Case Management

- Link analysis: [Gephi](https://gephi.org) (open-source graph visualization), [Visual Investigative Scenarios](https://vis.occrp.org) (OCCRP) — Maltego's limits are noted under §1
- Case notes and data cleanup: [Obsidian](https://obsidian.md) (also in the §18 table), [OpenRefine](https://github.com/OpenRefine) (messy data), [csvkit](https://github.com/wireservice/csvkit) (CSV from the command line)
- Charts, timelines and story maps: [Datawrapper](https://datawrapper.de), [Plotly](https://plot.ly), [D3js](https://d3js.org), [Observable](https://observablehq.com/), [Timeline](https://timeline.knightlab.com) and [StoryMap](https://storymap.knightlab.com) (Knight Lab)

**Deprioritized:** paid or niche link-analysis: Linkurious, Sentinel Visualizer, ORA-LITE, Wynyard Group (status unconfirmed); generic charting and infographic tools, not OSINT-specific: Aeon, Arbor.js, Beaker, Befunky, Bizint, Cacoo, Canva, Chart.js, chartblocks, Circos, creately, Crossfilter, Data Visualization Catalogue, Dropmark, dygraphs, easel.ly, Exhibit, Flot, FusionCharts, Google Charts, GraphX, Highcharts, Hohli, Infogr.am, Inkscape, Java Infovis Toolkit, JpGraph, jqPlot, Knoema, Listify, LocalFocus, Lucidchart, Mapline, Nodebox, OpenLayers, Palladio, Perspective, Piktochart, Pixxa, Preceden, QlikView, Quadrigram, Raphael, RAW, Snappa, Tableau Public, Tableau, Tagul, Textures.js, Tik-tok, Tiki-toki, Timeflow, Timeline (SIMILE), Timetoast, Venngage, Vis.js, Visme, Visualize Free, Visualize.me, visually, Vortex, ZingChart; QuickCode (Python/R environment); Text Machine cipher identifier (unvetted)

---

### 34. Learning, Blogs, Newsletters & Reference Lists

- Methodology and toolkits: [Bellingcat](https://www.bellingcat.com/) (plus its [Online Investigation Toolkit](https://bit.ly/bcattools)), [IntelTechniques](https://inteltechniques.com/), [OSINTCurious](https://osintcurio.us/), [Sector035](https://sector035.nl/), [NixIntel](https://nixintel.info/), [Trace Labs Blog](https://www.tracelabs.org/blog/rss.xml) (missing-person OSINT), [OSINT Dojo](https://www.osintdojo.com/resources/), [OSINT Updates](https://osintupdates.com/) (weekly newsletter)
- Practice: [Cipherstick](https://cipherstick.tech) (free OSINT puzzles, no account) — pairs with the GeoGuessr practice noted in the Geo-OSINT part
- Security journalism and research: [Krebs on Security](https://krebsonsecurity.com/feed/), [The Citizen Lab](https://citizenlab.ca/feed/), [SANS Internet Storm Center Diary](https://isc.sans.edu/rssfeed_full.xml)
- Related lists: [Awesome OSINT](https://github.com/jivoi/awesome-osint) (source of this merge), [OSINT-BIBLE](https://github.com/frangelbarrera/OSINT-BIBLE), [awesome-threat-intelligence](https://github.com/hslatman/awesome-threat-intelligence), [awesome-forensics](https://github.com/Cugu/awesome-forensics), [awesome-privacy](https://github.com/Lissy93/awesome-privacy/), [SecLists](https://github.com/danielmiessler/SecLists), [Infosec Reference](https://github.com/rmusser01/Infosec_Reference), [Awesome Translations](https://github.com/mbiesiad/awesome-translations#tools) (for the language-tools gap)
- Video: [SANS OSINT Summit 2021 playlist](https://www.youtube.com/playlist?list=PLs4eo9Tja8bj3jJvv42LxOkhc2_ylpS9y), [Bendobrown](https://www.youtube.com/c/Bendobrown)

**Deprioritized:** videos that are awareness or entertainment rather than technique: "Amazing mind reader reveals his 'gift'", "See how easily freaks can take over your life", Data to Go; vendor blogs: Skopenow, Social Links, eInvestigator; variable-quality publications: OSINT Techniques, OSINT Ambition Publication, OSINT Team; narrow focus: EUvsDisinfo (disinformation narratives), Inforrm's Blog (media law), OSINT Belarus (regional Telegram channel); OSINT Combine Blog; one-off survey results: Bellingcat tool-wishes spreadsheet and the 2022 "tools open source researchers say they need" post; aggregator pages: appsec.fyi OSINT Resources, Aware-online.com; adjacent security lists, not OSINT-specific: awesome-anti-forensic, awesome-ctf, awesome-hacking, awesome-honeypots, awesome-incident-response, awesome-malware-analysis, awesome-pentest, awesome-sec-talks, awesome-security, security-list, personal-security-checklist; awesome-lockpicking (off-topic); LERS Portal Directory (only relevant to authorised law-enforcement data requests — unvetted)

---

### 35. Gaming Platforms

*(Gaming accounts leak real-world habits: online hours, friend lists, time zones. These monitors log presence and playtime, so use them only inside a scoped investigation or on accounts whose owners consented. Steam profile and friend-list pivots are in §20.)*

- Steam: [steam_monitor](https://github.com/misiektoja/steam_monitor) — online/offline detection, games played, email alerts and CSV logging
- Xbox Live: [xbox_monitor](https://github.com/misiektoja/xbox_monitor) — same feature set for Xbox players
- PlayStation Network: [psn_monitor](https://github.com/misiektoja/psn_monitor) — same feature set for PSN players

**Deprioritized:** lol_monitor (League of Legends only, too game-specific for routine use — unvetted)

---

### 36. Music Streaming Services

*(Listening history is another pattern-of-life source, with the same consent caveat as §35.)*

- Public listening history: [lastfm_monitor](https://github.com/misiektoja/lastfm_monitor) — plays, skips, online/offline detection, email alerts and CSV logging
- Spotify profile and playlist changes: [spotify_profile_monitor](https://github.com/misiektoja/spotify_profile_monitor)

**Deprioritized:** spotify_monitor (tracks Spotify friends' listening activity, so it needs an account that already follows the target — unvetted)

---

**Bonus tools**

> Putting my Instagram reels to go use :) 
```
Ip address -> what is downloaded by the user: http://iknowwhatyoudownload.com/ (this was not opening in my pc :[ )
```

> Practitioners Note: if you ever get someone pgp keys you can decode it using https://cirw.in/gpg-decoder to decode it to get the email id  

<!-- The Geo-OSINT section continues below, unchanged from the original write-up -->
## Geo-Osint 

> practice this using geo-gusser

when looking at a picture how do you extract its location:here are clues to look at
```
1. Context
2. Foreground
3. Background
4. Map and marking
5. Trial and error 
```
**CONTEXT**
- Fastest way to get someone location is through the extraction on meta data
- Reverse image search and image based location search are effective methods
- Tools like street views form google or location query engines for location
- you can analyze the shadows in the image to extract more information


`but before you start looking at random clues go in this order`

The rule: work from broadest (eliminates the most of the world per clue)
to narrowest (confirms exact spot). Checking a "quick tell" out of order
wastes time — a road sign font means nothing until you already know
which country's font system you're comparing it against.

```
STAGE 1 — INSTANT WORLD-CUTTING CLUES (first 5 seconds)
├── LANGUAGE / SCRIPT on any visible text
│     Single biggest narrowing clue — cuts the world to a script
│     family / country group before anything else
├── STEERING WHEEL SIDE / TRAFFIC DIRECTION
│     Cuts the world in half instantly
└── ANY UNIQUE / "ONLY" TELL
      Kangaroo, baobab, gondola, longtail boat, torii gate, etc.
      If one appears → country confirmed → skip straight to Stage 6

STAGE 2 — HEMISPHERE / CLIMATE / LATITUDE
├── SHADOW DIRECTION          → Northern vs Southern Hemisphere
├── VEGETATION TYPE            → Tropical / temperate / high latitude band
└── TERRAIN TYPE               → Desert / mountain / coastline narrows further

STAGE 3 — CONTINENT / CULTURAL REGION
├── ARCHITECTURE STYLE         → Materials, roof type, window style
├── CLOTHING / PEOPLE          → Traditional dress, uniforms, religious clothing
└── ROAD MARKINGS + PLATE COLOR/SCRIPT

STAGE 4 — COUNTRY-SPECIFIC CONFIRMATION
├── ROAD SIGNS                 → Exact wording, colors, font
├── UTILITY POLES / POSTBOXES / POLICE CAR LIVERY
└── SHOP / FUEL / FAST FOOD BRANDS
      Only useful once narrowed to 2-3 candidates — checking brand
      signage against the whole world at Stage 1 is inefficient

STAGE 5 — REGION / CITY WITHIN THE COUNTRY
├── DIALECT / SCRIPT VARIANT   → Regional markers on signs
├── LOCAL-ONLY CHAIN BRANDS    → Not national chains
└── REGIONAL ARCHITECTURE VARIATION

STAGE 6 — EXACT PINPOINT
├── REVERSE IMAGE SEARCH       → Any visible landmark or storefront
├── STREET VIEW CROSS-REFERENCE→ Walk the candidate area virtually
└── UNIQUE STREET FURNITURE    → Manhole cover design, bollard style

STAGE 7 — TIME / DATE CONFIRMATION (always last)
├── SHADOW LENGTH + SUNCALC    → Against your Stage 6 candidate location
├── SEASONAL VEGETATION STATE  → Blossom, autumn color, snow
└── STAR POSITION via Stellarium → Clear night sky only
```
**HERE ARE THE MORE BROAD CHECKLISTS**
```
SHADOW ANALYSIS QUICK CHECKLIST

HEMISPHERE IDENTIFICATION
├── SHADOWS POINT NORTH
│     Sun travels south of overhead      → Northern Hemisphere
│     Countries: USA / Europe / China / Russia / Japan
│
└── SHADOWS POINT SOUTH
      Sun travels north of overhead      → Southern Hemisphere
      Countries: Australia / NZ / South Africa / Argentina / Brazil

TIME OF DAY FROM SHADOW LENGTH
├── VERY LONG shadows (several times object height)
│     Low sun angle                      → Early morning OR late afternoon
│     Combined with orange/golden light  → Sunrise OR sunset (golden hour)
│
├── MEDIUM shadows (roughly equal to object height)
│     Mid-morning OR mid-afternoon       → Roughly 9-10am OR 2-3pm
│
├── SHORT shadows (much shorter than object height)
│     Sun near overhead                  → Around midday (11am - 1pm)
│
└── NO shadows visible
      Sun directly overhead              → Near equator at solar noon ONLY
      Overcast sky                       → Cannot determine time from shadows

SHADOW DIRECTION = TIME OF DAY
├── NORTHERN HEMISPHERE
│     Shadows pointing NORTHWEST         → Morning (sun in southeast)
│     Shadows pointing NORTH             → Midday (sun in south)
│     Shadows pointing NORTHEAST         → Afternoon (sun in southwest)
│     Shadows pointing WEST              → Very early morning
│     Shadows pointing EAST              → Very late afternoon / evening
│
└── SOUTHERN HEMISPHERE
      Shadows pointing SOUTHWEST         → Morning (sun in northeast)
      Shadows pointing SOUTH             → Midday (sun in north)
      Shadows pointing SOUTHEAST         → Afternoon (sun in northwest)

LATITUDE CLUES FROM SHADOWS
├── VERY SHORT shadows even at midday
│     Sun very high overhead             → Tropical zone (within 23° of equator)
│     Countries: Indonesia / Thailand / Kenya / Brazil / Colombia
│
├── MODERATE shadows at midday
│     Sun at medium angle                → Temperate zone (23° - 60° latitude)
│     Countries: USA / Europe / China / Japan / Australia
│
└── VERY LONG shadows even at midday
      Sun stays low on horizon           → Polar / subpolar regions
      Countries: Iceland / Norway / Finland / Alaska / Southern Argentina

SEASONAL CLUES FROM SHADOWS
├── NORTHERN HEMISPHERE
│     Very short midday shadows          → June / July (summer solstice)
│     Very long midday shadows           → December / January (winter solstice)
│     Medium shadows                     → March / September (equinoxes)
│
└── SOUTHERN HEMISPHERE
      Very short midday shadows          → December / January (summer solstice)
      Very long midday shadows           → June / July (winter solstice)
      Medium shadows                     → March / September (equinoxes)

SHADOW ANGLE QUICK MATH
├── Shadow length = object height × tan(90° - sun altitude)
├── Sun altitude at midday (Northern Hemisphere)
│     90° - latitude + 23.5° (summer)
│     90° - latitude - 23.5° (winter)
│
└── Example: London (51°N) summer midday
      Sun altitude = 90 - 51 + 23.5 = 62.5°
      Shadow = very short (high sun)

SUNCALC METHOD (most reliable)
├── STEP 1  Go to suncalc.org
├── STEP 2  Drop pin on suspected location
├── STEP 3  Set the date from image context / EXIF
├── STEP 4  Adjust time slider until sun angle matches shadows
├── STEP 5  If shadows match → location AND time confirmed
└── STEP 6  Cross reference with hemisphere / season clues

SHADOW COLOR CLUES
├── BLUE tinted shadows                  → Cold climate / snow reflection / winter
├── ORANGE tinted shadows                → Golden hour (sunrise / sunset)
├── HARSH black shadows, no softness     → Tropical / desert (intense direct sun)
└── SOFT diffused shadows                → Overcast / temperate climate

SPECIAL SHADOW SITUATIONS
├── MULTIPLE shadows from one object
│     Two light sources visible          → Urban area, artificial lighting
│     Only at night                      → Cannot use for sun analysis
│
├── SHADOW ON SNOW
│     Blue shadow on white snow          → Winter / high altitude
│     Very long shadows on snow          → High latitude winter (sun very low)
│
├── SHADOWS INSIDE BUILDING / WINDOW
│     Sun coming from specific direction → Confirms compass orientation of room
│     Shadow of window frame on floor    → Can calculate sun angle precisely
│
└── NO SHADOWS AT ALL
      Completely overcast                → Cannot determine time or location
      Dense fog / rain                   → Cannot use shadow analysis

CROSS REFERENCE CHECKLIST
├── Step 1  Which way do shadows point?  → Confirms hemisphere
├── Step 2  How long are the shadows?    → Estimates time of day
├── Step 3  What direction do they point?→ Narrows hour of day
├── Step 4  How short at midday?         → Estimates latitude / season
├── Step 5  Confirm with SunCalc         → Pin exact time + location
└── Step 6  Cross check with vegetation  → Confirms season / hemisphere

QUICK TELLS
┌──────────────────────────────────────────────────┐
│ Shadows point north    → Northern hemisphere     │
│ Shadows point south    → Southern hemisphere     │
│ No shadow at midday    → Equatorial zone         │
│ Very long all day      → High latitude / winter  │
│ Blue tinted shadows    → Cold / snowy climate    │
│ Orange shadows         → Golden hour (6am/6pm)   │
│ Harsh black shadows    → Tropical / desert       │
│ Short midday shadows   → Summer OR tropics       │
│ Long midday shadows    → Winter OR high latitude │
└──────────────────────────────────────────────────┘
```
- Country specific road sign always exist here is list of the few
```
ROAD SIGNS QUICK CHECKLIST

STOP SIGN
├── Red octagon, STOP in English         → USA / Canada / Australia / NZ
├── Red octagon, STOP in French          → France / Quebec Canada
├── Red octagon, ARRÊT                   → Quebec Canada ONLY
├── Red octagon, ALTO                    → Mexico / Latin America
├── Red octagon, PARE                    → Brazil / Portugal
├── Red octagon, СТОП (Cyrillic)         → Russia / Eastern Europe
└── Red octagon, وقف (Arabic)            → Middle East / Gulf states

GIVE WAY / YIELD SIGN
├── Red/white inverted triangle, YIELD   → USA / Canada
├── Red/white inverted triangle, GIVE WAY→ UK / Australia / NZ / Ireland
├── Red/white inverted triangle, VORFAHRT→ Germany
└── Plain inverted red triangle          → Most of Europe (no text)

HIGHWAY / MOTORWAY SIGNS
├── Green background white text          → USA / Canada / Australia
├── Blue background white text           → Europe / Japan / China
├── Yellow background black text         → Germany (federal roads)
├── Blue with white A + number           → France (Autoroute)
├── Green with white M + number          → Spain (Motorway)
├── Blue with white E + number           → Pan-European E-roads
└── Green with white BR + number         → Brazil

SPEED LIMIT SIGNS
├── White circle red border, black number→ Europe / UK / Australia
├── White rectangle, black number        → USA / Canada
├── Yellow diamond, black number         → Australia (advisory speed)
└── White circle, no border, black number→ Japan / China

WARNING SIGNS
├── Yellow diamond black symbol          → USA / Canada / Australia / NZ
├── Red triangle white background        → Europe / UK / most world
├── Yellow triangle black symbol         → Japan
└── Orange diamond                       → USA (road construction only)

ANIMAL WARNING SIGNS
├── Kangaroo silhouette                  → Australia ONLY
├── Moose / elk silhouette               → Canada / Scandinavia / Russia
├── Reindeer silhouette                  → Finland / Norway / Sweden ONLY
├── Penguin silhouette                   → New Zealand / South Africa
├── Camel silhouette                     → Middle East / Australia outback
├── Elephant silhouette                  → India / Southeast Asia / Africa
└── Wombat silhouette                    → Australia ONLY

ROAD SURFACE MARKINGS
├── Yellow center + yellow edge lines    → USA / Canada
├── White center + white edge lines      → Europe / Australia
├── Yellow center + white edge lines     → Brazil / some Latin America
├── White center only, no edge           → Developing countries / rural Asia
└── Blue road markings                   → France (some older roads)

DISTANCE / MILEAGE
├── Miles on signs                       → USA / UK / Myanmar
├── Kilometers on signs                  → Rest of world
└── Both shown                           → Canada (transition areas)

LANGUAGE ON SIGNS
├── English only                         → USA / UK / Australia
├── English + French                     → Canada / some African countries
├── English + Maori                      → New Zealand
├── English + Welsh                      → Wales ONLY
├── French + local language              → France / West Africa
├── Cyrillic only                        → Russia / Belarus / Bulgaria
├── Arabic + English                     → UAE / Qatar / Jordan
├── Hebrew + Arabic + English            → Israel ONLY
├── Chinese + Pinyin + English           → China
├── Japanese + English                   → Japan (expressways)
└── Dual script (Latin + local)          → India / Malaysia / Sri Lanka

UNIQUE COUNTRY SPECIFIC SIGNS
├── Brown tourist signs                  → UK / Europe (heritage sites)
├── Blue H hospital signs (very common)  → UK
├── Green SOS boxes on highways          → France / Spain
├── Red torii gate on sign               → Japan (shrine nearby)
├── White on green, very minimal design  → Scandinavia
└── Highly decorative / colorful signs   → India / Pakistan

QUICK TELLS
┌──────────────────────────────────────────────────┐
│ ALTO on stop sign      → Mexico / Latin America  │
│ PARE on stop sign      → Brazil / Portugal       │
│ ARRÊT on stop sign     → Quebec Canada ONLY      │
│ Yellow diamond signs   → USA / Canada            │
│ Red triangle signs     → Europe / UK             │
│ Kangaroo sign          → Australia ONLY          │
│ Reindeer sign          → Scandinavia ONLY        │
│ English + Welsh        → Wales ONLY              │
│ Hebrew + Arabic        → Israel ONLY             │
│ Miles on signs         → USA / UK / Myanmar      │
│ Brown tourist signs    → UK / Europe             │
│ Torii gate sign        → Japan                   │
└──────────────────────────────────────────────────┘
```
- vegetation can too determine the area as well the region for example
```
VEGETATION QUICK CHECKLIST

TREES
├── Palm trees (tall, tropical)          → Within 35° of equator
│     Date palms (grey, spiky top)       → Middle East / North Africa
│     Coconut palms (leaning, beach)     → Tropical coastlines worldwide
│     Oil palms (very dense fronds)      → Malaysia / Indonesia / West Africa
│
├── Eucalyptus (peeling bark, grey-green)→ Australia / California / parts Africa
├── Baobab (massive trunk, tiny top)     → Sub-Saharan Africa / Madagascar ONLY
├── Acacia (flat-topped umbrella shape)  → East Africa / South Africa / Australia
├── Araucaria (umbrella pine, spiky)     → South America / Norfolk Island
├── Cherry blossom (pink spring flower)  → Japan / Korea / China (spring only)
├── Maple (bright red/orange leaves)     → Canada / Northeast USA (autumn only)
├── Olive trees (silver-green, twisted)  → Mediterranean ONLY
├── Cork oak (thick stripped bark)       → Portugal / Spain / Morocco
└── Mangrove (roots above water)         → Tropical coastlines worldwide

FORESTS
├── Dense dark conifer (pine/spruce)     → Scandinavia / Russia / Canada
├── Mixed deciduous (green + colour)     → Central Europe / Northeast USA
├── Rainforest (extremely dense green)   → Amazon / Congo / SE Asia
├── Cloud forest (misty, mossy)          → Central America / Colombia / Ethiopia
├── Bamboo forest                        → China / Japan / Southeast Asia
└── Sparse dry scrubland                 → Australia / South Africa / Spain

CROPS / FARMLAND
├── Rice paddies (flooded flat fields)   → Southeast Asia / East Asia
├── Terraced farming (stair steps)       → Southeast Asia / Andes / Mediterranean
├── Tea plantations (neat dark rows)     → India / Sri Lanka / Kenya / China
├── Coffee plantations (hilly, lush)     → Ethiopia / Colombia / Vietnam / Brazil
├── Vineyards (neat low rows)            → France / Italy / Spain / Germany / Chile
├── Lavender fields (purple)             → France (Provence) / Spain
├── Tulip fields (bright color rows)     → Netherlands ONLY
├── Wheat plains (vast flat golden)      → USA / Russia / Ukraine / Australia
├── Sugarcane (tall dense grass)         → Brazil / India / Thailand / Caribbean
└── Coca plants (low bushy rows)         → Colombia / Peru / Bolivia

GROUND / UNDERGROWTH
├── Red / orange clay soil               → Australia outback / parts Africa
├── Black volcanic soil                  → Iceland / Hawaii / Indonesia / Italy
├── White sandy soil                     → Middle East / coastal tropics
├── Heather (purple low shrub)           → UK / Scotland / Ireland / Scandinavia
├── Pampas grass (very tall silver)      → Argentina / Uruguay / Southern Brazil
├── Sagebrush (grey-green low shrub)     → USA Southwest ONLY
├── Fynbos (fine low shrubs)             → South Africa (Cape) ONLY
└── Tundra (flat, no trees, mossy)       → Northern Russia / Canada / Alaska

INSTANT COUNTRY TELLS
┌──────────────────────────────────────────────────┐
│ Baobab tree         → Africa / Madagascar ONLY   │
│ Tulip fields        → Netherlands ONLY           │
│ Eucalyptus          → Australia / California     │
│ Araucaria pine      → South America              │
│ Lavender fields     → France (Provence)          │
│ Sagebrush           → USA Southwest ONLY         │
│ Fynbos              → South Africa ONLY          │
│ Oil palm plantation → Malaysia / Indonesia       │
│ Cherry blossom      → Japan / Korea / China      │
│ Maple autumn color  → Canada / NE USA            │
│ Olive trees         → Mediterranean ONLY         │
│ Pampas grass        → Argentina / Uruguay        │
└──────────────────────────────────────────────────┘
```
- based on the visible terrain 
```
TERRAIN / LANDSCAPE QUICK CHECKLIST

MOUNTAINS
├── PERFECTLY CONE SHAPED volcano
│     Snow-capped, near ocean             → Japan (Fuji) / Indonesia / Philippines
│     Tropical jungle surroundings        → Indonesia / Papua New Guinea
│     Pacific coastline nearby            → Mexico / Central America
│
├── JAGGED SHARP peaks (like teeth)
│     Extreme height, glaciers visible    → Himalayas (Nepal / Tibet / Pakistan)
│     Green valleys below                 → Swiss / Austrian Alps
│     Very sharp dramatic spires          → Patagonia (Argentina / Chile) ONLY
│     Granite dome shapes                 → Yosemite USA / Madagascar
│
├── ROUNDED smooth peaks
│     Green covered, misty                → UK / Ireland / Scandinavia
│     Brown/bare at top                   → Scotland Highlands specifically
│     Very wide flat topped               → South Africa (Table Mountain)
│
├── FLAT TOPPED mountains (mesa / butte)
│     Red/orange rock                     → USA Southwest (Utah / Arizona) ONLY
│     Green flat top                      → Brazil (tepui) / Venezuela
│     Brown flat top                      → South Africa / Namibia
│
└── KARST LIMESTONE towers
      Green jungle covered pillars        → China (Guilin) / Vietnam (Ha Long)
      Surrounded by turquoise water       → Vietnam (Ha Long Bay) ONLY
      Inland dry setting                  → China specifically

DESERTS
├── SAND DUNES (tall, curved)
│     Orange/red sand                     → Namibia / UAE / Saudi Arabia
│     White/pale sand                     → Sahara / Libya / Algeria
│     Pink sand                           → Jordan (Wadi Rum) ONLY
│     Very tall massive dunes             → Namibia (Sossusvlei) ONLY
│
├── ROCKY DESERT (flat, stony)
│     Red rock surface                    → Australia outback ONLY
│     Grey/brown gravel flat              → Gobi (Mongolia / China)
│     Black volcanic rock                 → Iceland / Hawaii / Canary Islands
│     White salt flat                     → Bolivia (Salar de Uyuni) / USA (Bonneville)
│
├── CANYON / GORGE
│     Red/orange layered rock walls       → USA (Grand Canyon / Utah) ONLY
│     Narrow slot canyon                  → USA / Jordan (Petra)
│     Green river at bottom               → USA / Turkey / Ethiopia
│
└── SCRUBLAND / SEMI ARID
      Red soil, sparse dry bush           → Australia outback
      Grey green low scrub                → USA Southwest / Mexico
      Dry yellow grass, flat              → East Africa (savanna)
      Dry yellow grass, hilly             → South Africa / Zimbabwe

COASTLINES
├── FJORDS (steep cliffs straight into water)
│     Very dramatic, near vertical walls  → Norway ONLY
│     Similar but smaller scale           → New Zealand / Chile / Alaska
│
├── CLIFFS
│     White chalk cliffs                  → UK (Dover) / France (Normandy)
│     Black volcanic cliffs               → Iceland / Canary Islands / Hawaii
│     Red sandstone cliffs                → Australia / Portugal (Algarve)
│     Orange limestone cliffs             → Malta / Croatia / Greece
│
├── BEACHES
│     White sand + turquoise water        → Caribbean / Maldives / SE Asia
│     Black sand beach                    → Iceland / Hawaii / New Zealand
│     Pink sand beach                     → Bahamas / Bermuda / Greece (rare)
│     Grey/brown sand, cold looking water → UK / Northern Europe
│     Red sand beach                      → Prince Edward Island Canada ONLY
│     Green sand beach                    → Hawaii (Papakolea) ONLY
│
├── TROPICAL COASTLINE
│     Mangroves at water edge             → Southeast Asia / Caribbean / Africa
│     Coral visible through water         → Great Barrier Reef / SE Asia / Red Sea
│     Palm lined beach                    → Within 35° of equator
│     Rice paddies near coast             → Southeast Asia / Indonesia
│
└── ROCKY COASTLINE
      Granite boulders on beach           → Seychelles ONLY
      Limestone stacks in water           → Australia (Twelve Apostles)
      Basalt hexagonal columns            → Ireland (Giants Causeway) / Iceland
      Small rocky islands covered green   → Greece / Croatia / Norway

PLAINS / FLATLANDS
├── GRASSLAND / SAVANNA
│     Flat, dry yellow grass, acacia trees→ East Africa (Kenya / Tanzania)
│     Flat, dry, red soil visible         → Australia outback
│     Flat, green, very wide open         → Argentina (Pampas) / Uruguay
│     Flat, green, hedgerows dividing     → UK / France / Netherlands
│
├── STEPPE
│     Flat, dry, treeless, windy looking  → Mongolia / Kazakhstan / Russia
│     Similar but colder, snow possible   → Siberia / Northern Kazakhstan
│
├── WETLANDS
│     Flat, reeds, many waterways         → Netherlands / Bangladesh / Louisiana
│     Flooded forest                      → Amazon / Congo / Southeast Asia
│     Salt marshes, coastal               → UK / France Atlantic coast
│
└── AGRICULTURAL PLAINS
      Perfect grid fields                 → USA Midwest / Canada prairies
      Irregular shaped fields             → Europe / UK
      Flooded paddy fields                → Southeast Asia / East Asia
      Circular irrigation patterns        → USA / Saudi Arabia (seen from above)

RIVERS / LAKES
├── RIVERS
│     Brown/red muddy wide river          → Amazon / Congo / Mississippi
│     Bright blue clear river             → New Zealand / Scandinavia / Alps
│     Green river through canyon          → USA Southwest / Turkey
│     Very wide flat slow river           → Russia / Siberia / Bangladesh
│
├── LAKES
│     Bright blue, snow peaks around      → Switzerland / Austria / Canada Rockies
│     Turquoise glacial lake              → Canada (Banff) / New Zealand / Patagonia
│     Pink lake                           → Australia (Lake Hillier) ONLY
│     Salt lake, white edges              → Bolivia / USA (Great Salt Lake)
│     Very large, looks like sea          → Great Lakes USA / Lake Victoria Africa
│
└── GLACIERS
      Glacier descending into sea         → Iceland / Greenland / Alaska
      Glacier on mountain side            → Alps / Himalayas / Andes / Rockies
      Glacier lake at bottom (grey water) → Patagonia / Iceland / New Zealand

INSTANT COUNTRY TELLS
┌──────────────────────────────────────────────────┐
│ Perfect cone volcano      → Japan / Indonesia    │
│ Red rock mesa / butte     → USA Southwest ONLY   │
│ Karst limestone pillars   → China / Vietnam      │
│ Pink desert sand          → Jordan (Wadi Rum)    │
│ Massive orange sand dunes → Namibia              │
│ White chalk cliffs        → UK / France          │
│ Black sand beach          → Iceland / Hawaii     │
│ Granite boulders on beach → Seychelles ONLY      │
│ Basalt hexagonal columns  → Ireland / Iceland    │
│ Fjords                    → Norway / NZ / Chile  │
│ Pink lake                 → Australia ONLY       │
│ Turquoise glacial lake    → Canada / Patagonia   │
│ Sharp spire mountains     → Patagonia ONLY       │
│ Flat topped tepui         → Venezuela / Brazil   │
│ Salt flat white mirror    → Bolivia ONLY         │
└──────────────────────────────────────────────────┘
```
- from people and clothing visible in the image 
```
PEOPLE / CLOTHING QUICK CHECKLIST

TRADITIONAL / NATIONAL DRESS
├── SOUTH ASIA
│     Sari (long draped fabric, women)    → India / Sri Lanka / Bangladesh
│     Salwar Kameez (tunic + trousers)    → India / Pakistan / Bangladesh
│     Lungi (wraparound cloth, men)       → Bangladesh / South India / Myanmar
│     Dhoti (white wraparound, men)       → Rural India ONLY
│     Topi (white cap, men)               → Nepal / Northern India
│
├── SOUTHEAST ASIA
│     Longyi (wraparound skirt, men+women)→ Myanmar ONLY
│     Ao Dai (fitted silk tunic, women)   → Vietnam ONLY
│     Batik patterned clothing            → Indonesia / Malaysia
│     Conical straw hat (non la)          → Vietnam ONLY
│     Sarong common on both genders       → Indonesia / Malaysia / Pacific Islands
│
├── EAST ASIA
│     Hanbok (wide skirt, short jacket)   → Korea ONLY
│     Kimono / Yukata (formal occasions)  → Japan ONLY
│     Qipao / Cheongsam (fitted dress)    → China / Hong Kong
│     Mandarin collar shirts common       → China
│
├── MIDDLE EAST
│     Thobe (long white robe, men)        → Saudi Arabia / Gulf states
│     Keffiyeh (red+white checked scarf)  → Saudi Arabia / Jordan
│     Keffiyeh (black+white checked)      → Palestine / Jordan
│     Abaya (full black robe, women)      → Saudi Arabia / Gulf states
│     Niqab (face veil, eyes only)        → Saudi Arabia / Gulf states
│     Hijab (headscarf only)              → Muslim-majority countries globally
│     Kandura (white robe) + Agal         → UAE / Qatar / Kuwait ONLY
│
├── CENTRAL ASIA
│     Chapan (colorful robe, men)         → Uzbekistan / Tajikistan
│     Kalpak (white felt pointed hat)     → Kyrgyzstan ONLY
│     Tubeteika (small embroidered cap)   → Uzbekistan / Kazakhstan
│     Yurt + traditional felt clothing    → Mongolia / Kazakhstan
│
├── AFRICA
│     Kente cloth (colorful woven strips) → Ghana ONLY
│     Dashiki (colorful loose shirt)      → West Africa broadly
│     Shuka (red checked cloth)           → Maasai Kenya / Tanzania ONLY
│     Boubou (long wide robe)             → West Africa / Senegal
│     Djellaba (long hooded robe)         → Morocco / Algeria / Tunisia
│     Basotho blanket (geometric pattern) → Lesotho ONLY
│
├── EUROPE
│     Lederhosen (leather shorts)         → Bavaria Germany / Austria ONLY
│     Dirndl (dress + apron)              → Bavaria Germany / Austria ONLY
│     Kilt (tartan skirt, men)            → Scotland ONLY
│     Clog shoes visible                  → Netherlands
│     Saami clothing (colorful trim)      → Lapland Norway/Sweden/Finland
│
├── AMERICAS
│     Poncho (woven, colorful)            → Andes (Peru / Bolivia / Ecuador)
│     Sombrero (very wide brim)           → Mexico ONLY
│     Mola (colorful panel clothing)      → Panama (Kuna people) ONLY
│     Indigenous feather headdress        → Amazon Brazil / Central America
│
└── PACIFIC
      Grass skirt / tapa cloth            → Pacific Islands broadly
      Lei (flower garland around neck)    → Hawaii / Pacific Islands
      Lavalava (wraparound cloth)         → Samoa / Tonga / Fiji

RELIGIOUS CLOTHING
├── Turban (large, colorful)              → Sikh India / Punjab ONLY
├── Turban (white, simple)                → Muslim scholars / West Africa
├── Kippah (small round cap)              → Jewish / Israel
├── Buddhist monk orange robes            → Thailand / Myanmar / Sri Lanka
├── Buddhist monk yellow/red robes        → Tibet / Mongolia / Bhutan
├── Orthodox priest black robes + hat     → Russia / Greece / Eastern Europe
├── Catholic nun habit                    → Global but common Latin America / Europe
└── Hare Krishna robes (orange/white)     → India / Global urban areas

HEADWEAR
├── Fez (red flat-topped cylinder)        → Morocco / Turkey / Egypt
├── Keffiyeh (checkered headscarf)        → Middle East
├── Pagri (tied turban)                   → India / Pakistan
├── Deerstalker hat                       → UK / Scotland
├── Akubra (wide brim felt)               → Australia ONLY
├── Gaucho beret (flat, black/grey)       → Argentina / Uruguay ONLY
├── Ushanka (fur ear-flap hat)            → Russia / Eastern Europe / Mongolia
├── Conical bamboo hat                    → Vietnam / China / Southeast Asia
├── Karakul (lamb fur hat)                → Central Asia / Afghanistan
└── Baseball cap worn backwards           → USA / Global youth culture

FOOTWEAR
├── Wooden clogs                          → Netherlands / Belgium
├── Geta (wooden sandals)                 → Japan (traditional)
├── Huarache sandals (woven leather)      → Mexico ONLY
├── Juttis (pointed embroidered)          → India / Pakistan
└── Barefoot common in streets            → Parts of Southeast Asia / Pacific Islands

UNIFORMS
├── POLICE
│     Bobby helmet (tall black dome)      → UK ONLY
│     Blue beret + blue uniform           → France
│     Khaki uniform + beret               → India / Pakistan / many ex-British
│     Grey/olive uniform                  → Russia / Eastern Europe
│     Brown uniform + star badge          → USA (sheriff) ONLY
│
├── MILITARY
│     Red ceremonial uniform + bearskin   → UK (Buckingham Palace) ONLY
│     Olive drab very common              → USA / most NATO
│     Digital camo pattern                → USA / Canada / modern armies
│     Blue/white naval uniform            → Global but prominent Russia/France
│
├── SCHOOL UNIFORMS
│     Blazer + tie + grey trousers        → UK / Australia / NZ / Hong Kong
│     White shirt + dark trousers/skirt   → Japan / South Korea / SE Asia
│     No uniform, casual clothes          → USA / Canada / most of Europe
│     White shirt + blue/red neckerchief  → Cuba / China (Young Pioneers)
│
└── WORKERS
      High-vis yellow/orange vest common  → UK / Australia / Europe
      Conical hat + simple clothes        → Vietnam / China (rural workers)
      White lab coat common               → Global (medical / scientific)
      Salwar kameez work clothing         → India / Pakistan (very common)

BODY FEATURES / MARKINGS
├── Bindi (dot on forehead, women)        → India / Hindu communities globally
├── Mehndi / henna on hands              → India / Middle East / North Africa
├── Lip plate (disc in lower lip)         → Ethiopia (Mursi tribe) ONLY
├── Neck rings (coils on neck)            → Myanmar (Kayan) / Thailand
├── Full body tattoos visible             → Samoa / Polynesia / Maori NZ
└── Scarification patterns               → Sub-Saharan Africa broadly

QUICK TELLS
┌──────────────────────────────────────────────────┐
│ Kilt                      → Scotland ONLY        │
│ Lederhosen / Dirndl       → Bavaria / Austria    │
│ Sombrero                  → Mexico ONLY          │
│ Longyi wraparound         → Myanmar ONLY         │
│ Ao Dai fitted tunic       → Vietnam ONLY         │
│ Hanbok wide skirt         → Korea ONLY           │
│ Maasai red shuka          → Kenya / Tanzania     │
│ Kente cloth strips        → Ghana ONLY           │
│ Akubra wide brim          → Australia ONLY       │
│ Bobby helmet police       → UK ONLY              │
│ Orange monk robes         → Thailand / SE Asia   │
│ Red+white keffiyeh        → Saudi / Jordan       │
│ Conical non la hat        → Vietnam ONLY         │
│ Kalpak pointed hat        → Kyrgyzstan ONLY      │
│ Basotho blanket           → Lesotho ONLY         │
│ Gaucho beret              → Argentina / Uruguay  │
│ Ushanka fur hat           → Russia / Central Asia│
│ Bearskin guard hat        → UK (palace) ONLY     │
│ School blazer + tie       → UK / Australia / NZ  │
│ White + neckerchief       → Cuba / China kids    │
└──────────────────────────────────────────────────┘
```
- based on any animal visible in the photo
```
ANIMALS QUICK CHECKLIST

INSTANTLY UNIQUE (seen = country confirmed)
├── Kangaroo                             → Australia ONLY
├── Koala                                → Australia ONLY
├── Wombat                               → Australia ONLY
├── Platypus                             → Australia ONLY
├── Tasmanian Devil                      → Tasmania Australia ONLY
├── Lemur                                → Madagascar ONLY
├── Kiwi bird                            → New Zealand ONLY
├── Giant Panda                          → China ONLY
├── Komodo Dragon                        → Indonesia (Komodo) ONLY
├── Quokka                               → Australia (Rottnest Island) ONLY
├── Capybara (wild, large rodent)        → South America ONLY
├── Jaguar                               → Central / South America ONLY
├── Llama / Alpaca (wild)                → Peru / Bolivia / Chile ONLY
├── Bison roaming free                   → USA / Canada (national parks)
└── Polar Bear                           → Arctic (Canada / Norway / Russia)

LARGE MAMMALS
├── ELEPHANTS
│     Very large ears, flat back         → Africa (savanna / forest)
│     Smaller ears, domed back           → India / Sri Lanka / SE Asia
│     Used as working animal             → Thailand / Myanmar / Sri Lanka
│
├── RHINO
│     Two horns, grey, square lip        → Africa (Kenya / South Africa)
│     One horn, grey, folded skin        → India / Nepal ONLY
│
├── HIPPO
│     In or near river                   → Sub-Saharan Africa ONLY
│
├── GIRAFFE
│     Open savanna, acacia trees         → East / Southern Africa ONLY
│
├── ZEBRA
│     Open grassland, mixed herds        → East / Southern Africa ONLY
│
├── LION
│     Open savanna, dry grass            → Sub-Saharan Africa
│     Forested / grassland (rare)        → India (Gir Forest) ONLY
│
├── TIGER
│     Dense jungle or grassland          → India / Bangladesh / Russia (Siberia)
│     Managed reserve setting            → India most likely
│
├── GORILLA
│     Dense rainforest                   → Congo / Rwanda / Uganda ONLY
│
├── MOOSE
│     Boreal forest / wetlands           → Canada / Scandinavia / Russia
│     Near road in forest                → Canada / Alaska / Scandinavia
│
├── GRIZZLY / BROWN BEAR
│     Near salmon river / forest         → Alaska / Canada / Russia
│     Mountain setting                   → USA Rockies / European Alps
│
└── REINDEER / CARIBOU
      Semi-domesticated near people      → Finland / Norway / Sweden
      Wild large herds                   → Canada / Alaska / Russia

MEDIUM MAMMALS
├── RED FOX in urban setting             → UK / Western Europe
├── RACCOON near buildings               → USA / Canada ONLY
├── COYOTE in suburban area              → USA / Mexico ONLY
├── DINGO (wild dog, sandy)              → Australia ONLY
├── HYENA near settlement                → Sub-Saharan Africa / North Africa
├── BABOON near roads                    → South Africa / East Africa
├── MACAQUE monkey near temples          → India / Japan / Southeast Asia
├── PROBOSCIS MONKEY (huge nose)         → Borneo (Malaysia) ONLY
├── SNOW MONKEY in hot spring            → Japan ONLY
└── MEERKAT standing upright             → South Africa / Namibia / Botswana

BIRDS
├── EMU (large flightless, grey)         → Australia ONLY
├── OSTRICH (very large, long neck)      → Africa / Australia (farmed)
├── CASSOWARY (blue head, black)         → Australia / Papua New Guinea
├── FLAMINGO (pink, standing in water)   → East Africa / Caribbean / Spain
├── TOUCAN (huge colorful beak)          → Central / South America ONLY
├── MACAW (large colorful parrot)        → South / Central America ONLY
├── PELICAN on coast                     → USA / Australia / Africa / Europe
├── PUFFIN (black/white, colorful beak)  → Iceland / UK / Norway / Canada
├── PENGUIN (wild)
│     Large Emperor / King               → Antarctica / South Georgia
│     Smaller Jackass penguin            → South Africa ONLY
│     Small blue penguin                 → New Zealand / Australia
├── CONDOR (huge black soaring)          → Andes (Peru / Chile / Argentina)
├── BALD EAGLE soaring                   → USA / Canada ONLY
├── RED KITE (rusty red, forked tail)    → UK / Spain / Germany
├── PEACOCK roaming free                 → India / Sri Lanka
├── HORNBILL (large curved beak)         → Southeast Asia / Africa
└── BIRDS OF PARADISE                   → Papua New Guinea ONLY

REPTILES
├── CROCODILE in river                   → Australia / Africa / SE Asia
├── ALLIGATOR near water                 → USA (Florida / Louisiana) ONLY
├── GHARIAL (long thin snout)            → India / Nepal ONLY
├── MONITOR LIZARD near water            → Southeast Asia / Australia / Africa
├── IGUANA roaming freely                → Caribbean / Central America
├── CHAMELEON in tree                    → Madagascar / East Africa
└── COBRA near settlement                → India / Southeast Asia / Africa

MARINE / COASTAL
├── GREAT WHITE SHARK (news / signs)     → South Africa / Australia / USA
├── DUGONG / MANATEE                     → Southeast Asia / Caribbean / Africa
├── SEA LION on rocks                    → USA Pacific / South America / NZ
├── FUR SEAL colony on beach             → South Africa / New Zealand / Alaska
├── WHALE SHARK near boats               → Philippines / Mexico / Australia
└── CLOWNFISH / REEF FISH visible        → Indo-Pacific / Red Sea / Caribbean

FARM / WORKING ANIMALS
├── WATER BUFFALO plowing field          → Southeast Asia / India / China
├── CAMEL as transport                   → Middle East / North Africa / Central Asia
├── YAK carrying loads                   → Tibet / Nepal / Mongolia ONLY
├── ELEPHANT being ridden                → Thailand / India / Sri Lanka
├── DONKEY common on roads               → Middle East / North Africa / Mediterranean
├── ZEBU (humped cattle)                 → India / East Africa / Southeast Asia
└── LLAMA carrying loads                 → Peru / Bolivia ONLY

INSECTS / SMALL CREATURES
├── GIANT CENTIPEDE visible              → Southeast Asia / Australia / South America
├── SCORPION common                      → Middle East / North Africa / SW USA
├── TARANTULA crossing road              → South America / SW USA / Southeast Asia
└── FIREFLIES at dusk                    → USA / Japan / Southeast Asia

INSTANT COUNTRY TELLS
┌──────────────────────────────────────────────────┐
│ Kangaroo                  → Australia ONLY       │
│ Lemur                     → Madagascar ONLY      │
│ Giant Panda               → China ONLY           │
│ Snow monkey in hot spring → Japan ONLY           │
│ One horned rhino          → India / Nepal ONLY   │
│ Proboscis monkey          → Borneo ONLY          │
│ Quokka                    → Australia ONLY       │
│ Capybara wild             → South America ONLY   │
│ Alligator                 → USA Southeast ONLY   │
│ Gharial                   → India / Nepal ONLY   │
│ Bald eagle                → USA / Canada ONLY    │
│ Condor soaring            → Andes ONLY           │
│ Birds of paradise         → Papua New Guinea     │
│ Yak carrying loads        → Tibet / Nepal        │
│ Jackass penguin           → South Africa ONLY    │
│ Dingo                     → Australia ONLY       │
│ Raccoon near buildings    → USA / Canada ONLY    │
│ Llama carrying loads      → Peru / Bolivia       │
└──────────────────────────────────────────────────┘
```
- architectural clues are always there in the nature
```
ARCHITECTURAL QUICK CHECKLIST

BUILDING MATERIAL
├── Timber framing (exposed wood beams)  → UK / Germany / Northern Europe
├── Adobe / mud brick                    → Middle East / North Africa / SW USA
├── Red brick dense rows                 → UK / Netherlands / Eastern USA
├── White painted plaster                → Greece / Mediterranean / Middle East
├── Concrete Soviet blocks               → Former Soviet states / Eastern Europe
└── Bamboo / wood stilts                 → Southeast Asia / Pacific Islands

WINDOWS
├── Very large windows, thin walls       → Netherlands / Scandinavia
├── Small windows, thick walls           → Middle East / North Africa
├── Shutters on every window             → France / Italy / Spain
├── Sash windows (slide up/down)         → UK / USA (older buildings)
└── Floor-to-ceiling glass               → Modern Scandinavia / Japan

WALLS
├── Exposed red/orange brick             → UK / Netherlands / Eastern USA
├── White smooth plaster                 → Greece / Mediterranean
├── Wooden horizontal planks             → Scandinavia / New England USA
├── Wooden vertical planks               → Rural Japan / Korea
├── Stone walls (rough cut)              → UK countryside / Ireland / France
└── Corrugated iron walls                → Rural Africa / Pacific Islands

BALCONIES
├── ornate iron railings                 → France / Italy / Spain
├── Wide wooden balconies                → Scandinavia / Alpine Europe
├── Concrete slab balconies              → Soviet bloc / Eastern Europe
├── No balconies at all                  → UK terraced houses
└── Ornate carved wood balconies         → Rural Russia / Turkey

RELIGIOUS BUILDINGS
├── Gold onion dome                      → Russia / Orthodox countries
├── Blue dome white walls                → Greece (Cyclades islands) ONLY
├── Gothic tall stone spires             → Western Europe
├── Minaret (tall thin tower)            → Muslim-majority country
├── Pagoda tiered tower                  → East / Southeast Asia
├── Wooden stave church                  → Norway ONLY
└── Corrugated iron church               → Pacific Islands / rural Africa

UNIQUE STRUCTURES
├── Windmills (white, rotating)          → Netherlands ONLY
├── Canal houses (tall narrow)           → Netherlands / Belgium
├── Trulli (cone stone roofs)            → Puglia, Italy ONLY
├── Rondavel (round thatched hut)        → South Africa / Zimbabwe
├── Longhouse (very long wooden)         → Malaysia / Borneo / Scandinavia
└── Ger / Yurt (round felt tent)         → Mongolia / Central Asia ONLY

FENCES / WALLS
├── Dry stone walls (no mortar)          → UK / Ireland countryside
├── White picket fence                   → USA / Australia
├── Tall concrete security walls         → South Africa / Latin America
├── Bamboo fencing                       → Japan / Southeast Asia
└── Mud brick compound walls             → West Africa / Middle East

QUICK TELLS
┌─────────────────────────────────────────────────┐
│ Windmill                → Netherlands           │
│ Trulli cone roof        → Italy (Puglia) ONLY   │
│ Yurt / Ger              → Mongolia ONLY         │
│ Blue dome               → Greece (Cyclades)     │
│ Gold onion dome         → Russia                │
│ Wooden stave church     → Norway                │
│ Canal narrow houses     → Netherlands           │
│ Rondavel round hut      → Southern Africa       │
│ Adobe mud brick         → Middle East / SW USA  │
│ Soviet concrete blocks  → Eastern Europe        │
└─────────────────────────────────────────────────┘
```
- Infrastructure clues are always there as well
```
INFRASTRUCTURE QUICK CHECKLIST

POWER POLES
├── Wooden poles          → USA / Canada / Japan
├── Concrete poles        → Europe / Russia / China
└── No visible poles      → Dense urban Europe (underground)

TRAFFIC LIGHTS
├── Horizontal mount      → USA / Canada / some Asia
└── Vertical mount        → Europe / UK / Australia / most world

ROAD CENTER LINE
├── Yellow center line    → USA / Canada / Brazil / Japan
└── White center line     → Europe / Australia / most world

POSTBOXES
├── Red standalone        → UK
├── Yellow               → Germany / France / Spain
├── Blue                 → USA
├── Green                → Ireland / Italy / China
└── Orange               → Netherlands

POLICE CARS
├── Black and white       → USA
├── Yellow/blue checker   → UK
├── Green and yellow      → Australia / New Zealand
└── All white + stripe    → Russia / Eastern Europe

BUS TYPE
├── Red double decker     → UK / Hong Kong
├── Yellow school bus     → USA / Canada ONLY
├── Tuk-tuk / rickshaw    → Thailand / India / Sri Lanka
└── Jeepney               → Philippines ONLY

MANHOLE COVERS
├── Ornate / artistic     → Japan
├── City crest cast iron  → UK / Europe
└── Plain grey concrete   → Developing world

FUEL STATION BRAND
├── Petrobras             → Brazil
├── PEMEX                 → Mexico
├── Lukoil / Rosneft      → Russia
├── Indian Oil / HPCL     → India
├── Sinopec / PetroChina  → China
└── Pertamina             → Indonesia
```
- stars can be used to pin point location and time
```
1. Identify visible stars/constellations
2. Use Stellarium (stellarium.org)
3. Input suspected location
4. Adjust time until star positions match
```
- traffic / wheel direction / license plate / car-info 
```
STEP 1: Which SIDE of road are cars driving on?

├── LEFT SIDE of road (overtaking on right)
│     → Steering wheel on RIGHT side of car
│     → Go to Step 2
│
└── RIGHT SIDE of road (overtaking on left)
      → Steering wheel on LEFT side of car
      → Go to Step 3

STEP 2: LEFT-HAND TRAFFIC countries
├── ASIA
│     Circular script on plates            → Thailand / Sri Lanka
│     Japanese characters on plates        → Japan
│     Hindi / Devanagari on plates         → India / Nepal
│     English only, yellow rear plate      → Malaysia / Singapore
│
├── OCEANIA
│     Plates vary by state, kangaroo       → Australia
│     Similar but smaller country code     → New Zealand
│
├── AFRICA
│     Yellow/white plates, English         → South Africa / Kenya
│     Green/white plates                   → Zimbabwe / Zambia
│
├── EUROPE
│     Blue strip left, GB/UK letters       → United Kingdom
│     Blue strip left, IRL letters         → Ireland
│     Blue strip left, M letters           → Malta
│     Blue strip left, CY letters          → Cyprus
│
└── CARIBBEAN
      English text, island name on plate   → Jamaica / Barbados / Trinidad

STEP 3: RIGHT-HAND TRAFFIC countries
├── NORTH AMERICA
│     Varies wildly by state, no EU strip  → USA
│     Province name on top or bottom       → Canada
│     Green/white, Spanish text            → Mexico
│
├── EUROPE (all right-hand traffic)
│     Blue EU strip on left + country code
│       D   = Germany    F  = France
│       E   = Spain      I  = Italy
│       NL  = Netherlands B = Belgium
│       PL  = Poland     RO = Romania
│       SE  = Sweden     NO = Norway
│       CH  = Switzerland A = Austria
│     White plate, black text, EU strip    → Most of EU
│     Yellow plate front AND rear          → Netherlands specifically
│     Yellow rear only, white front        → UK (older pre-2021 style)
│
├── RUSSIA / EASTERN EUROPE
│     Cyrillic letters on plate            → Russia / Belarus
│     RUS code, white plate black text     → Russia
│     UA code                              → Ukraine
│
├── MIDDLE EAST
│     Arabic script on plate               → Saudi / UAE / Jordan / Egypt
│     Hebrew + English on plate            → Israel
│     Both Arabic and Latin               → Lebanon / Morocco / Tunisia
│
├── EAST ASIA
│     Chinese characters, red background   → China (government vehicles)
│     Chinese characters, white plate      → China (standard)
│     Korean characters on plate           → South Korea
│
└── SOUTH AMERICA
      All right-hand traffic
      MERCOSUR blue strip on left          → Brazil / Argentina / Chile
      Spanish text, country flag strip     → Colombia / Peru / Venezuela

STEP 4: PLATE COLOR QUICK ID
├── All yellow (front + rear)              → Netherlands
├── Red background                         → China government / diplomatic
├── Green background                       → Electric vehicle (many countries)
├── Black plate white text                 → Germany (vintage) / some US states
├── Blue plate white text                  → France (old) / some EU diplomatic
├── White plate blue EU strip left         → Standard EU country
├── No EU strip, highly variable           → USA / Canada / Australia
└── Arabic script only                     → Gulf states (Saudi / UAE / Kuwait)

QUICK CHEAT SHEET
┌─────────────────────────────────────────────────┐
│ LEFT traffic  → UK, Australia, Japan, India      │
│ RIGHT traffic → USA, Europe, China, Russia       │
│ Yellow plates → Netherlands (both) / UK (rear)   │
│ EU blue strip → All European Union countries     │
│ Arabic script → Gulf / Middle East               │
│ Cyrillic      → Russia / Belarus / Ukraine       │
│ No strip,     → USA / Canada / Australia         │
│ wildly varied                                    │
└─────────────────────────────────────────────────┘
```
- Differentiation based on roof
```
STEP 1: Is the roof FLAT?
├── YES → Has water tanks on top           → India / South Asia
│         Soviet-style concrete blocks     → Eastern Europe / Russia
│         Sandy/arid surroundings          → Middle East / North Africa
└── NO  → Step 2

STEP 2: Is it covered in GRASS or TURF?
├── YES → Iceland / Norway / Faroe Islands
└── NO  → Step 3

STEP 3: Do the eaves CURVE UPWARD at corners?
├── YES → Slight curve, dark grey tiles    → Japan
│         Dramatic curve, ornate ridge     → China
│         Very dramatic curve upward       → Thailand / Vietnam
└── NO  → Step 4

STEP 4: What is the MATERIAL?

├── ORANGE / RED clay curved tiles
│     Hot dry surroundings                 → Spain / Italy / Portugal / Greece
│     Tropical surroundings               → Latin America
│     Central European city               → Austria / Hungary / Croatia
│
├── DARK GREY / BLACK slate tiles
│     Chimney stacks visible              → UK / Ireland / Wales
│     Steep narrow houses                 → Northern France
│
├── GRASS / TURF covered                  → Iceland / Norway (Step 2)
│
├── GREEN painted METAL
│     Onion domes nearby                  → Russia
│     Plain urban buildings               → Eastern Europe / Baltic
│
├── CORRUGATED IRON / METAL
│     Rural setting                       → Australia / New Zealand
│     Urban informal housing              → Sub-Saharan Africa
│     Tropical rural                      → Southeast Asia
│
├── ASPHALT SHINGLE (dark granular)
│     Suburban sprawl                     → USA / Canada
│     Similar but steeper                 → Canada specifically
│
├── WOODEN SHINGLES (weathered grey/brown)
│     Coastal or forested setting         → Scandinavia / Canada
│     Siberian wilderness                 → Russia (Siberia)
│
└── VERY STEEP dark wood or slate
      Alpine mountains visible            → Switzerland / Austria / Bavaria
      Forested flat land                  → Scandinavia / Finland

STEP 5: Any UNIQUE ROOF SHAPES?
├── Stepped gable ends (staircase profile) → Netherlands
```
- identifying language based on signs
```
STEP 1: Right-to-left text?
├── YES → Flowing cursive horizontal    → Arabic
│         Diagonal flowing (Nastaliq)   → Urdu
│         Blocky square letters         → Hebrew
│         Arabic + adds پ چ ژ           → Persian
└── NO  → Step 2

STEP 2: Vertical text (top → bottom)?
├── YES → Mongolian (only vertical script)
└── NO  → Step 3

STEP 3: Latin alphabet (A-Z)?
├── NO  → Step 4
└── YES → Step 5

STEP 4: Which non-Latin script?
├── Looks like А Б В (Russian-style)    → Cyrillic family
│     Has ы э ъ                         → Russian
│     Has і ї є                         → Ukrainian
│     Has ђ љ њ                         → Serbian
│     Similar to Russian                → Bulgarian
├── Horizontal bar on top (क ख)         → Hindi / Nepali
├── α β γ δ familiar but not Latin      → Greek
├── Circular loops (ก ข ค)             → Thai
├── Rounded curvy unique (ა ბ გ)        → Georgian
├── Complex syllabic (አ ቡ)             → Amharic
└── Rounded unique (Ա Բ Գ)             → Armenian

STEP 5: Complex square CJK characters?
├── YES → Mixed with rounded あ or angular ア → Japanese
│         Only complex squares               → Chinese
│         Geometric blocks with circles 한   → Korean
└── NO  → Step 6

STEP 6: Any special/accented characters?
├── NO  → English / Malay / Indonesian / Swahili
└── YES → Step 7

STEP 7: Which character?
├── ñ  or  ¿ ¡                          → Spanish
├── ß  or  ä ö ü                        → German
├── é  è  ê  ç  œ                       → French
├── ã  õ  ç                             → Portuguese
├── ı  (dotless i)  ğ  ş               → Turkish
├── ắ  ề  ổ  (stacked accents)          → Vietnamese
├── ą  ę  ł  ź  ż                       → Polish
├── č  ř  ž  š  (háček ˇ)              → Czech / Slovak
├── ő  ű  (double acute)                → Hungarian
├── å  æ  ø                             → Norwegian / Danish
├── å  ä  ö  (no æ/ø)                  → Swedish
└── ș  ț  â  î                          → Romanian
```
- based on shop signs and shop culture
```
FOOD STALLS / SHOP SIGNAGE QUICK CHECKLIST

CONVENIENCE STORE CHAINS
├── 7-Eleven (very dense concentration)  → Thailand (most per capita globally)
├── 7-Eleven + FamilyMart + Lawson       → Japan / Taiwan / China
├── FamilyMart dominant                  → Japan / South Korea / Taiwan
├── Lawson dominant                      → Japan / China
├── CU / GS25 / Emart24                  → South Korea ONLY
├── Spar                                 → Netherlands / UK / South Africa
├── Tesco Express                        → UK / Thailand / Hungary
└── Coles Express / IGA                  → Australia ONLY

SUPERMARKET CHAINS
├── Woolworths / Coles                   → Australia ONLY
├── Pick n Pay / Shoprite / Checkers     → South Africa ONLY
├── Lidl / Aldi green logo               → Europe broadly
├── Carrefour (blue/red/white)           → France / Spain / Belgium / Middle East
├── Tesco (red + blue stripe)            → UK / Ireland / Central Europe
├── Sainsburys / Morrisons / Asda        → UK ONLY
├── Walmart / Target / Kroger            → USA ONLY
├── Loblaws / Metro / Sobeys             → Canada ONLY
├── Migros / Coop (orange/red)           → Switzerland ONLY
├── Rewe / Edeka                         → Germany ONLY
├── Mercadona                            → Spain ONLY
├── Pão de Açúcar / Extra               → Brazil ONLY
└── Reliance Fresh / Big Bazaar          → India ONLY

FAST FOOD REGIONAL TELLS
├── Jollibee (red bee mascot)            → Philippines ONLY
├── Lotteria                             → South Korea / Vietnam / Japan
├── Mos Burger                           → Japan / Taiwan / Southeast Asia
├── Popeyes very dominant                → USA / Canada
├── Whataburger (orange/white striped)   → USA South / Texas ONLY
├── In-N-Out Burger                      → USA West Coast ONLY
├── Nandos (flame logo)                  → South Africa / UK / Australia
├── Wimpy                                → South Africa / UK
├── Steers / Debonairs                   → South Africa ONLY
├── Hungry Jacks                         → Australia ONLY (Burger King rebranded)
├── Max Burgers                          → Sweden ONLY
├── Bob's Burgers                        → Brazil ONLY
├── Hesburger                            → Finland / Baltic states ONLY
└── Hardees dominant                     → USA Southeast / Middle East

STREET FOOD STALLS
├── Pad Thai / Som Tam stall             → Thailand
├── Pho cart / stall                     → Vietnam ONLY
├── Banh Mi stand                        → Vietnam ONLY
├── Satay grill on street               → Malaysia / Indonesia / Singapore
├── Taco / burrito cart                  → Mexico / USA Southwest
├── Arepa stall                          → Colombia / Venezuela ONLY
├── Empanada cart                        → Argentina / Chile / Colombia
├── Kebab / doner stand                  → Turkey / Germany / Middle East
├── Jerk chicken drum barrel             → Jamaica / Caribbean ONLY
├── Bunny chow stall                     → South Africa (Durban) ONLY
├── Roti canai stall                     → Malaysia / Singapore ONLY
├── Dim sum cart                         → China / Hong Kong / SE Asia
├── Takoyaki / yakitori stall            → Japan ONLY
├── Currywurst stand                     → Germany (Berlin) ONLY
└── Borek / simit cart                   → Turkey ONLY

DRINK / BEVERAGE BRANDS
├── Efes beer signs                      → Turkey ONLY
├── Tusker beer                          → Kenya ONLY
├── Castle / Black Label beer            → South Africa ONLY
├── XXXX / VB / Tooheys beer             → Australia ONLY
├── Kingfisher beer signs                → India ONLY
├── Chang / Singha beer                  → Thailand ONLY
├── Tiger beer dominant                  → Singapore / Malaysia / Vietnam
├── Bintang beer                         → Indonesia ONLY
├── Modelo / Corona dominant             → Mexico ONLY
├── Quilmes beer                         → Argentina ONLY
├── Cristal beer                         → Peru / Chile
├── Brahma / Skol / Antarctica           → Brazil ONLY
├── Keo / Leon beer                      → Cyprus ONLY
└── Mythos / Fix beer                    → Greece ONLY

TELECOM / PHONE SHOP SIGNS
├── Jio / Airtel / BSNL / Vi             → India ONLY
├── Turkcell / Vodafone TR / Turk Telekom→ Turkey ONLY
├── MTN / Vodacom / Cell C               → South Africa / Sub-Saharan Africa
├── M-Pesa signage (mobile money)        → Kenya / East Africa ONLY
├── Safaricom green shops                → Kenya ONLY
├── Orange / SFR / Bouygues              → France / French-speaking countries
├── Telstra / Optus / Vodafone AU        → Australia ONLY
├── NTT Docomo / SoftBank / au           → Japan ONLY
├── SK Telecom / KT / LG U+              → South Korea ONLY
├── China Mobile / China Unicom          → China ONLY
└── Digicel shops                        → Caribbean / Pacific Islands ONLY

PHARMACY / DRUGSTORE SIGNS
├── Boots (blue logo)                    → UK / Ireland / Thailand
├── CVS / Walgreens / Rite Aid           → USA ONLY
├── Shoppers Drug Mart                   → Canada ONLY
├── Priceline Pharmacy                   → Australia ONLY
├── Apotheke (green cross)               → Germany / Austria / Switzerland
├── Pharmacie (green cross, France)      → France / Belgium / Luxembourg
├── Cruz Verde / Cruz Blanca             → Chile / Colombia / Mexico
└── Green cross sign broadly             → European pharmacy (universal)

PRODUCT PACKAGING CLUES
├── Cyrillic text on packaging           → Russia / Eastern Europe
├── Arabic text on packaging             → Middle East / North Africa
├── Thai script on packaging             → Thailand
├── Hebrew text on packaging             → Israel
├── Korean script on packaging           → South Korea
├── Chinese + English on packaging       → China / Taiwan / Singapore
├── Metric units only (ml, kg)           → Most of world (non-USA)
├── Imperial units (oz, lbs, fl oz)      → USA primarily
└── Dual metric + imperial               → Canada / UK

MARKET / BAZAAR TYPES
├── Night market (many small stalls lit) → Taiwan / Thailand / Malaysia
├── Wet market (fresh meat/fish/veg)     → Southeast Asia / Hong Kong / China
├── Souq / souk (covered narrow lanes)   → Middle East / North Africa ONLY
├── Grand Bazaar style (arched ceiling)  → Turkey / Iran / Central Asia
├── Floating market (stalls on boats)    → Thailand / Vietnam / Indonesia
├── Sunday car boot / flea market        → UK / Australia / New Zealand
├── Farmers market (organic focus)       → USA / Western Europe / Australia
└── Mercado (covered, colorful stalls)   → Latin America broadly

QUICK TELLS
┌──────────────────────────────────────────────────┐
│ Jollibee red bee           → Philippines ONLY    │
│ Hungry Jacks               → Australia ONLY      │
│ Whataburger orange stripe  → USA South ONLY      │
│ M-Pesa signage             → Kenya / East Africa │
│ Safaricom green            → Kenya ONLY          │
│ Bunny chow stall           → South Africa ONLY   │
│ Currywurst stand           → Germany ONLY        │
│ Pho cart                   → Vietnam ONLY        │
│ Banh Mi stand              → Vietnam ONLY        │
│ Takoyaki stall             → Japan ONLY          │
│ Borek / simit cart         → Turkey ONLY         │
│ Jerk chicken drum barrel   → Jamaica ONLY        │
│ Efes beer signs            → Turkey ONLY         │
│ Tusker beer                → Kenya ONLY          │
│ Kingfisher beer            → India ONLY          │
│ Bintang beer               → Indonesia ONLY      │
│ Hesburger                  → Finland ONLY        │
│ Max Burgers                → Sweden ONLY         │
│ Quilmes beer               → Argentina ONLY      │
│ Roti canai stall           → Malaysia ONLY       │
└──────────────────────────────────────────────────┘
```
- boat and harbors verification 
```
BOATS / HARBORS QUICK CHECKLIST

FISHING BOATS
├── Longtail boat (long engine arm)      → Thailand ONLY
├── Banca boat (outriggers both sides)   → Philippines ONLY
├── Jukung (small outrigger, colorful)   → Indonesia / Bali ONLY
├── Dhow (triangular lateen sail)        → Oman / UAE / East Africa / Yemen
├── Felucca (small wooden sail)          → Egypt / Nile ONLY
├── Sampan (flat bottom, oar steered)    → China / Vietnam / Hong Kong
├── Coracle (round wicker boat)          → Wales / India (Bhutan)
├── Dinghy / trawler, plain white        → UK / Ireland / Northern Europe
├── Colorful wooden hull, Mediterranean  → Greece / Turkey / Croatia / Italy
└── Pirogues (carved wooden canoe)       → West Africa / Pacific Islands

COMMERCIAL / CARGO VESSELS
├── Massive container ships              → Major ports globally (Singapore / Rotterdam)
├── River barges (long flat)             → Europe (Rhine / Danube) / China / Russia
├── Traditional junk (red sails)         → China / Hong Kong ONLY
├── Ferry with cars on deck              → Scandinavia / Greece / UK / Japan
└── Small inter-island ferries           → Philippines / Indonesia / Greece

HARBOR / PORT CLUES
├── HARBOR TYPE
│     Natural fjord harbor               → Norway / New Zealand / Chile
│     Artificial modern container port   → Singapore / Rotterdam / Shanghai
│     Old stone quay walls               → UK / Ireland / Malta / Croatia
│     Floating pontoon docks             → USA / Australia / Scandinavia
│     Wooden stilted dock                → Southeast Asia / Pacific Islands
│
├── HARBOR SURROUNDINGS
│     Whitewashed buildings on cliffs    → Greece (Santorini / Mykonos)
│     Red/orange roofed town behind      → Croatia / Montenegro / Italy
│     Green hills steep into water       → Norway / New Zealand
│     Flat desert behind harbor          → UAE / Saudi / Oman
│     Dense urban skyscrapers            → Hong Kong / Singapore / Dubai
│     Wooden fishing shacks on stilts    → Southeast Asia / Pacific
│
└── HARBOR EQUIPMENT
      Red/green buoys (standard)         → Most of world (IALA system)
      Yellow buoys dominant              → Some Asian waters
      Wooden fish traps in water         → Southeast Asia / Middle East
      Lobster pot stacks on dock         → UK / Ireland / Canada / Maine USA
      Massive cranes, container stacks   → Major global shipping ports

SEA / WATER COLOR
├── Turquoise / crystal clear            → Caribbean / Maldives / Mediterranean
├── Deep blue, very clear                → Open Pacific / Atlantic / Mediterranean
├── Green tinted                         → North Sea / Baltic / Irish Sea
├── Brown / murky                        → River mouths / Bangladesh / Indonesia
├── Grey, choppy, cold looking           → North Atlantic / UK / Norway
└── Red tinted (algae bloom)             → Specific seasonal event anywhere

BOAT ACTIVITY CLUES
├── Rice barge on river                  → Vietnam / Myanmar / Thailand
├── Gondola in canal                     → Venice Italy ONLY
├── Abra (wooden water taxi)             → Dubai Creek ONLY
├── Water taxi, very busy crossing       → Bangladesh / Vietnam / Philippines
├── Speedboat towing skiers              → USA / Australia / Western Europe
└── Reed boat (totora reeds)             → Lake Titicaca Peru / Bolivia ONLY

QUICK TELLS
┌──────────────────────────────────────────────────┐
│ Longtail boat          → Thailand ONLY           │
│ Banca outrigger        → Philippines ONLY        │
│ Gondola                → Venice Italy ONLY       │
│ Dhow triangular sail   → Middle East / E Africa  │
│ Felucca                → Egypt / Nile ONLY       │
│ Red sailed junk        → China / Hong Kong       │
│ Abra water taxi        → Dubai ONLY              │
│ Reed boat              → Lake Titicaca ONLY      │
│ Coracle round boat     → Wales / India           │
│ Jukung outrigger       → Bali / Indonesia        │
│ Lobster pots on dock   → UK / Ireland / Canada   │
│ Whitewashed cliff port → Greece                  │
└──────────────────────────────────────────────────┘
```
- night and artificial lighting quick checklist
```
NIGHT / ARTIFICIAL LIGHTING QUICK CHECKLIST

STREET LIGHT COLOR
├── Warm orange / yellow (sodium)        → UK / older Europe / Middle East
├── Cool white / blue-white (LED)        → Modern cities globally post 2015
├── Warm white LED                       → Japan / South Korea / Modern Europe
├── Green tinted lights                  → Some older US cities / Latin America
└── Inconsistent mix of colors           → Developing world / older infrastructure

NEON SIGNS / SIGNAGE
├── Dense Chinese neon / LED signs       → Hong Kong / China / Taiwan
├── Dense Korean neon signage            → South Korea (Gangnam / Seoul)
├── Japanese neon + vending machines lit → Japan (very dense vending machines)
├── Arabic neon signs                    → Middle East / North Africa
├── Cyrillic neon signs                  → Russia / Eastern Europe
├── Hebrew neon signs                    → Israel ONLY
├── Thai script neon                     → Thailand
└── Latin script only neon               → USA / Europe / Latin America

LIGHT DENSITY FROM ABOVE
├── Extremely dense bright grid          → USA (city grid pattern)
├── Dense but organic pattern            → Europe / older cities
├── Very dense + river visible dark      → Major river cities (London/Paris/Cairo)
├── Sparse with dark gaps                → Developing world / rural areas
├── Almost no light                      → North Korea ONLY (from satellite)
└── Single bright spot in darkness       → Remote mining / oil facility

SPECIFIC LIGHTING CLUES
├── VENDING MACHINES lit on street       → Japan ONLY (extremely dense)
├── Red lanterns hanging lit             → China / Chinese New Year / Vietnam
├── Neon cross on buildings              → South Korea (churches everywhere)
├── Green mosque lights at night         → Muslim-majority countries
├── Orthodox church lit gold at night    → Russia / Greece / Eastern Europe
├── Blue/white lit buildings             → Greece / Israel (national colors)
├── Casino strip lighting                → Las Vegas USA / Macau China / Monaco
└── Stadium floodlights visible          → Urban areas globally

TRAFFIC LIGHT COLOR
├── Standard red/amber/green             → Most of world
├── Countdown timer on lights            → China / South Korea / India
├── Flashing green before red            → Some Eastern European countries
└── Pedestrian light makes sound         → Japan / South Korea / Australia

QUICK TELLS
┌──────────────────────────────────────────────────┐
│ Dense vending machines lit  → Japan ONLY         │
│ Almost no light from above  → North Korea ONLY   │
│ Neon cross on buildings     → South Korea        │
│ Green mosque lights         → Muslim countries   │
│ Red lanterns hanging        → China / Vietnam    │
│ Hebrew neon signs           → Israel ONLY        │
│ Casino strip lighting       → Las Vegas / Macau  │
│ Orange sodium streetlights  → UK / older Europe  │
│ Countdown traffic timers    → China / India      │
└──────────────────────────────────────────────────┘
```
- sports and recreation checklist
```
SPORTS / RECREATION QUICK CHECKLIST

PITCH / FIELD TYPE
├── CRICKET PITCH (flat oval, 22 yard strip center)
│     White flannels, red ball            → UK / Australia / NZ / South Africa
│     Colored kit, white ball             → India / Pakistan / Sri Lanka / UAE
│
├── BASEBALL DIAMOND (90° infield dirt)
│     Large stadiums, very common         → USA / Canada
│     Smaller setups                      → Japan / Cuba / Dominican Republic
│
├── AUSTRALIAN RULES FOOTBALL
│     Oval pitch, tall H + behind posts   → Australia ONLY
│
├── RUGBY PITCH (H posts, rectangular)
│     Very common, village level          → UK / NZ / South Africa / France
│     Less common                         → Australia / Argentina / Japan
│
├── AMERICAN FOOTBALL (yard lines visible)
│     Very wide H posts                   → USA / Canada ONLY
│
├── GAELIC FOOTBALL / HURLING pitch       → Ireland ONLY
│
└── SEPAK TAKRAW court (rattan ball net)  → Southeast Asia ONLY

COURT TYPES
├── Basketball court (very common)        → USA / Global urban areas
├── Petanque / boules gravel court        → France / Southern Europe ONLY
├── Padel court (glass walls)             → Spain / Latin America / Middle East
├── Gateball court                        → Japan / China (elderly sport)
└── Volleyball on beach very common       → Brazil / USA / Europe coastal

WINTER SPORTS
├── Ski lifts on mountain                 → Alps / Rockies / Andes / Japan
├── Outdoor ice rink natural              → Canada / Russia / Scandinavia
├── Curling rink                          → Canada / Scotland / Scandinavia
└── Ski jump structure visible            → Norway / Austria / Finland / Japan

WATER SPORTS
├── Surfing (beach break, consistent)     → Australia / Hawaii / Portugal / South Africa
├── Dragon boat racing                    → China / Hong Kong / Taiwan / Vietnam
├── Waka ama (outrigger canoe racing)     → New Zealand / Pacific Islands ONLY
└── Kite surfing common                   → Netherlands / Brazil / Australia

RECREATIONAL CLUES
├── Pétanque players in park              → France / Southern Europe
├── Chess players in public park          → Russia / Eastern Europe / China
├── Tai chi in park at dawn               → China / Taiwan / SE Asian Chinese communities
├── Sumo wrestling                        → Japan ONLY
├── Kabaddi being played                  → India / Bangladesh / Pakistan
├── Capoeira on street                    → Brazil ONLY
├── Zorbing / bungee common               → New Zealand (adventure capital)
└── Polo being played                     → Argentina / UK / Pakistan / UAE

STADIUM CLUES
├── Very large, modern, air conditioned   → UAE / Qatar / Saudi (indoor stadiums)
├── Old Victorian brick stadium           → UK ONLY
├── Wooden bleachers, small               → USA rural / Latin America
├── Standing terraces still in use        → Eastern Europe / South America
└── Retractable roof stadium              → USA / Australia / Canada / Japan

QUICK TELLS
┌──────────────────────────────────────────────────┐
│ Cricket oval pitch        → British Commonwealth │
│ Baseball diamond          → USA / Japan / Cuba   │
│ AFL oval + behind posts   → Australia ONLY       │
│ Gaelic football           → Ireland ONLY         │
│ Petanque gravel court     → France / S Europe    │
│ Sepak takraw net          → SE Asia ONLY         │
│ Sumo ring                 → Japan ONLY           │
│ Capoeira on street        → Brazil ONLY          │
│ Dragon boat racing        → China / SE Asia      │
│ Waka ama canoe            → NZ / Pacific ONLY    │
│ Kabaddi being played      → India / Pakistan     │
│ Curling rink              → Canada / Scotland    │
│ Polo match                → Argentina / UK       │
│ Victorian brick stadium   → UK ONLY              │
└──────────────────────────────────────────────────┘
```
**Since the location is now identified to the narrowest possible region, let's hop on the map engine and get to the very location !**

### STAGE MANUAL MAP HUNT (THE FINAL SWEEP)

You have a confirmed country, a probable city or region, and a shortlist of 2–4 candidate areas. This is where you stop using clues from the image and start moving through the map manually until reality matches what you see.

---

#### SET UP YOUR WORKING AREA FIRST

```
Open two windows side by side:
  Left  → the image you are locating
  Right → Google Maps (or Apple Maps / Yandex Maps as backup)

Drop a pin on your best candidate area first
Switch to Satellite view
Zoom until individual buildings / blocks are distinguishable
DO NOT open Street View yet — satellite orientation comes first
```

---

#### STEP 1 — ORIENT FROM ABOVE BEFORE GOING TO GROUND LEVEL

```
In satellite view, identify ONE structural feature from the image
that is visible from above:

  Road layout shape     → straight grid / curved / radial / irregular
  Building footprint    → rectangular / courtyard / L-shaped / detached
  Roof color / material → visible from satellite in most cities
  Green space           → park, garden, tree line adjacent to the scene
  Water feature         → river bend, canal, coastline edge
  Parking lot shape     → distinctive in satellite view

Scan your candidate area for that single feature
When you find a match → that block or that road is your entry point
If you find zero matches across the whole candidate area → wrong area,
expand outward or move to your next candidate
```

---

#### STEP 2 — NARROW TO A CANDIDATE BLOCK

```
Once you have an entry point from satellite:

  Look at the road network immediately around it
  Count how many roads intersect your candidate block
  Match that count and angle to what is visible in the image
  (a T-junction looks very different from a 4-way cross or a curve)

From the image extract one more architectural anchor:
  Which direction does the main building face? (shadow tells you this)
  Is the building on a corner, mid-block, or at a road end?
  Is there a visible gap / alley / car park immediately adjacent?

Cross-reference these two things simultaneously on satellite:
  Road geometry + building position together
  When both match → you are on the correct block
```

---

#### STEP 3 — DROP INTO STREET VIEW AND WALK IT

```
Click into Street View at the edge of your candidate block
NOT at the midpoint — start at a corner so you can see
down two roads simultaneously

Move through Street View methodically:
  ONE direction at a time
  Check left and right at every junction before advancing
  Do not skip ahead — a single storefront or pole or fence
  can be the confirmation you need

What to match first in Street View:
  1. Sky / roofline silhouette against the image
     (this is fast and eliminates the wrong side of a block immediately)
  2. Road surface width and marking
  3. Pavement / sidewalk material and width
  4. Any vegetation (single tree, hedge, planter) adjacent to the scene
  5. Storefront / signage if visible in the original image
  6. Building entrance position and door color / style
```

---

#### STEP 4 — THE COMPASS LOCK (CONFIRM YOUR FACING DIRECTION)

```
Once you think you have the right block in Street View:

  Look at the shadow direction in the original image
  Note which direction in Street View the shadow would point
  Rotate your Street View view until shadows in Street View
  would match the original if it were the same time of day

  If the facing direction is consistent with the shadow analysis
  from Stage 7 → you have confirmed orientation
  If it is not → you are either on the wrong block or
  facing the wrong way — rotate 180° and re-check before moving

This is the step most people skip and it is why they end up
one block off or on the wrong side of the street
```

---

#### STEP 5 — FINAL PIN CONFIRMATION

```
You have the right block and the right facing direction.
Now identify the single most specific detail in the image:

  Most specific = least likely to repeat anywhere else nearby

  Good specifics:
    A building with a unique color combination
    A mural or wall art
    A specific shopfront with a readable name
    A distinctive lamp post or bollard style
    A utility box or postbox in an unusual position
    A crack / patch in the road surface (visible in Street View)
    A drain / grate in a specific position relative to the kerb

  Bad specifics (too common to confirm):
    "A red car parked outside"   → not fixed
    "A tree on the left"         → too common
    "A white building"           → too common

Find that one specific detail in Street View
When it matches → place your final pin at that exact position
```

---

#### THE MINDSET FOR THE MANUAL HUNT

```
Do not try to match everything at once
Pick ONE anchor feature, find it, then use it as a base
to find the next feature from there

The guy who narrows 40 blocks to 1 is not doing anything
clever — he is doing this loop very fast:

  Pick a feature → scan → match or discard → move to next block
  repeat until only one block matches two features simultaneously

Speed comes from:
  Committing to a candidate and testing it fully before moving on
  Never trying to hold two candidate blocks in your head at once
  Using satellite to pre-eliminate whole streets before Street View
  (Street View is slow — satellite culls first, Street View confirms)

When you are stuck:
  Go back to the image and find a feature you have not used yet
  The image always has more information than you have extracted
  Change the zoom level on satellite — sometimes a feature is only
  visible at a specific zoom (a narrow alley, a small courtyard)
  Try Yandex Maps Street View or Apple Maps Look Around as a backup
  — they cover some areas where Google Street View is absent or outdated
```
