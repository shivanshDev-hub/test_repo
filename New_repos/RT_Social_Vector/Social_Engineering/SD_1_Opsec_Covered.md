# OPSEC Guide 1: Data Reduction — The Field Book

_OPSEC Guide 1 · Field Book edition · Draft v3 (OSINT-first) · Tool and price facts last checked October 2026_

## Reality of the today's World 

**Philosophy.** We have official reached a point where our data have become more valuable then our own self . this is not is not a optional guide but rather a way to fight back in the digital world ! 
### Who this protects you from

- **Tier 1 — Advertisers and data brokers.** Everyone is in their database. Easy to fix.
- **Tier 2 — Criminals, scammers, doxxers, stalkers.** They live off leaked emails, phone numbers and public posts.
- **Tier 3 — Breaches.** Companies you trusted get hacked; whatever you gave them leaks.
- **Tier 4 — People near you.** Employers, schools, ex-partners, nosy neighbours.
- **Tier 5 — A well-funded state or targeted attacker.** This book will _raise the cost_ for them but will not make you invisible. For that tier, continue with Guides 2 and 3 (GrapheneOS, Qubes OS, Whonix, Tails and Tor).

### Scope

`this is devided into many sections like covered , hardened and targetted based on the osint threat model`

- Covered is for everyday normal people like our moms and dads , or maybe some normal teenage boy browsing over the internet 
- Hardened is for the people that particularly work in the industry where they can be watched or targeted , or for the journalist and researcher how must hide their identity 
- Targetted is for the people that have geniune threats behide them like state level adverseries and etc that are being constantly watched , with every activity being monitored or hunted down like osint profs  
### Order of work

This book starts by _looking_, not installing. Parts 1 to 5 find out what is already public about you and get it removed or retired. Only then do Parts 6 onward build the habits and tools that stop the pile from growing back. The order matters: tools protect what you create from now on, they do nothing for what is already out there, and what is already out there is what a stranger will actually find. Tiers 1, 2 and 4 above all live off what is already findable.

Parts 1 to 5 need a few small prerequisites (a clean browser profile, a private notes vault and a removal mailbox: `donot worry these are covered ahead`). Part 1 gives you the 20-minute pre-flight.

## 1. Know Your OSINT Surface — Audit Yourself Before Anyone Else Does

**Philosophy.**Lesten to put this simply is that you cannot defend what you donot know even exist so we are going to figure out what exist about us , we are goint to remove that and then and then only are are going start with fresh new setup and accounts

`I will assume you have idea about how the osint work before we even have a deep dive , if you donot just have a look at the guides before this`
### 🧠 A Little about OSINT flow 

`You can alomst skip this section if you know what osint is !`
#### How strangers build a profile (the mosaic effect)

One fact is harmless. Joined facts are not. A typical chain:

**Here is a example of going from the Forum handle to Username** 
1. An old forum handle matches a username on a photo site.
2. A photo there shows a landmark and a timestamp.
3. The same handle on GitHub has a real email in its commits.
4. That email appears in a breach next to a phone number.
5. The phone number sits on a people-search site beside your home address and relatives.

Every link in that chain is a _join key_ (an identifier that appears in more than one place; Part 7 explains them). The audit hunts for the joins, not just the facts, because breaking one join often unravels the whole chain.

#### The seven layers of your OSINT surface

1. **Search results** for your names, handles, emails and phone numbers.
2. **Things you posted:** social profiles, forums, reviews, code, wishlists, fitness maps.
3. **Things others posted about you:** tags, group photos, club and school pages, news, rosters.
4. **Data-broker and people-search profiles** (Part 3).
5. **Public records:** property, company filings, courts, domain registrations, professional licences, electoral rolls where public.
6. **Leaks and breaches:** what is already in the wild (Part 13 covers the password side).
7. **Archives and caches:** copies of all of the above (Part 4).

### 🛠 Pre-flight (20 minutes)

Parts 1 to 5 need three small things, set up before you start:

1. **A clean browser profile.** Logged-in searches are personalised, and you want a stranger's view. Firefox: type `about:profiles` → _Create a New Profile_ → name it Audit. Chrome, Edge or Brave: _Add profile → continue without an account_. Stay logged out of everything in it.
2. **A private place for your exposure register.** The register is a map of where you leak, so treat it as sensitive. Install [KeePassXC](https://keepassxc.org), create a database, and keep the register as a note inside it. The same database will also hold the old passwords you recover along the way.
3. **A dedicated removal mailbox.** A free [Proton Mail](https://proton.me/mail) or [Tuta](https://tuta.com) account used _only_ for removal requests and confirmations (Parts 3 and 4). It is not your final email setup; that comes in Part 7.

### 🛠 Do this: The audit (form 10min to 2 to 3 hours)

**Step 1. Build your seed list.** Save this in the register as _Seeds_:

- Every name you have used: legal, nicknames, maiden or previous names, spelling variants.
- Every username you remember, current and old.
- Every email address, current and old.
- Every phone number, current and old, written the ways people write it.
- Addresses from the last ten years.
- Employers, schools, clubs and any domains you own.
- Names of close relatives (a pivot strangers use).

Memory aids: your browser's saved logins, your password manager, and your old inboxes (search _welcome_, _verify_, _receipt_ and _password reset_).

**Step 2. Search your seeds.** Run each pattern below on Google, Bing, DuckDuckGo and Brave Search (their indexes differ). Check the Images and News tabs, and look past page one.

| Looking for                  | Query pattern                                                                                           |
| ---------------------------- | ------------------------------------------------------------------------------------------------------- |
| Name plus place              | `"Full Name" "City"`                                                                                    |
| Name plus employer or school | `"Full Name" "Employer"`                                                                                |
| Documents carrying your name | `"Full Name" filetype:pdf` (also `docx`, `xlsx`)                                                        |
| One platform at a time       | `site:linkedin.com "Full Name"` (repeat for facebook.com, instagram.com, x.com, reddit.com, github.com) |
| Your email                   | `"me@example.com"`, then just the part before the @                                                     |
| Your phone number            | each format people write it: `"+1 555 010 0199"`, `"(555) 010-0199"`, `"555-010-0199"`                  |
| Your address                 | `"12 Example Street"` and the street with the city                                                      |
| Your handle                  | `"yourhandle"` and `"yourhandle" -site:yourownsite.com` to hide your own pages                          |

Quotes force an exact match and `-site:` excludes a domain. Log every hit that is _about you_, with its URL, in the register.

`if you have a email logged in you can just idetify the logged in services from the google account and the email inbox--> WILL BE THE SINGLE BEST INFO POINT`

**Step 3. Quick tool sweeps** Run these on tools to quickly identify the OSINT attack data surface points 

- usernames -> [whatsmyname](https://whatsmynameapp.org/) {can and will generate noise so look for the popular option} 
- emails ids -> [mail meteor](https://mailmeteor.com/tools/reverse-email-lookup) {Honestly having main and then using google acc mana or inbox is way faster}
- Phone number -> just use true caller at this point bro :/ 

`later you can go to the osint write up to spend 5-10Hrs having detailed osint on yourself,but write now you must remove the most critical hits there are`

**Step 4. Check your emails against breaches.** Enter each email at [Have I Been Pwned](https://haveibeenpwned.com) and read _which data_ each breach leaked (passwords, phone numbers, addresses). 

**Step 5. Reverse-search your images.** Take your avatars and public profile photos and search them with Google Lens, [TinEye](https://tineye.com) and Bing Visual Search. Optional: [PimEyes](https://pimeyes.com) is a face search over the open web (not social platforms). Its free tier shows thumbnails; seeing source links needs a paid plan, so check current terms. Its opt-out asks for a photo of your face plus a redacted ID and only hides you from PimEyes' own results, so decide on purpose whether that trade is worth it. 

**Step 6 (Later-timetaking) . Check the places search engines miss.**

- **Maps and reviews:** Google Maps contributions, Yelp, Tripadvisor.
- **Wishlists and registries:** Amazon and baby registries often show a name and city.
- **Fitness apps:** public activity maps whose start and end points sit at your front door.
- **Payment apps with social feeds** (for example Venmo in the US).
- **Code:** commit emails (`git log --format='%ae' | sort -u` in each repo), your SSH keys at `github.com/USERNAME.keys`, package-registry author fields (npm, PyPI), and PGP keyservers such as keys.openpgp.org.
- **Domains you own:** look them up at [lookup.icann.org](https://lookup.icann.org). If the registrant shows your name or address, that is Critical.
- **Company filings:** directors' home addresses are public in many countries.
- **Property, court and licensing records:** your local portals, plus real-estate listing sites that show interior photos of your home.
- **Street-level imagery** of your house.
- **Old personal sites and profiles:** blogs, Flickr, SoundCloud, Steam, Goodreads, Last.fm, anything public by default.
- **Old shared files:** Drive, Dropbox or Docs links set to _anyone with the link_.

**Step 7 (Later-timetaking) . Look for what others post.** Search your name inside school, club, team and employer pages, and review tags and group photos. Note who posted each item, because the fix is to ask them (Part 2).

**Score each finding and write it down.** Use this register layout:

| #   | Where (URL)             | What it exposes               | Layer | Joins it to        | Severity | Fix (Part) | Status | Re-check |
| --- | ----------------------- | ----------------------------- | ----- | ------------------ | -------- | ---------- | ------ | -------- |
| 1   | example.com/profile/xyz | Real name + old handle + city | 2     | Handle `darkfox92` | High     | 2          | Open   |          |

Severity scale:

- **Critical:** home address, live location or routine, government ID numbers, financial details, your phone number beside your name on a public page, anything about your children's school or schedule, intimate images, a still-used password in a breach or public repo.
- **High:** your real name linked to a pseudonymous handle (a join), employer plus city, family members' names, photos tied to your name or location, a personal email in public commits, old accounts with reused passwords.
- **Medium:** hobbies, reviews, purchase history, old profile photos, public follows and likes.
- **Low:** details you have deliberately made public, such as a professional profile.

Work in this order: **Critical and easy** today; **Critical and slow** (for example property or company records) start now, because they take weeks; then High; then the rest.

### Decision tree: what did I find, where do I go?

- _An account I still own, with content I made_ → Part 2.
- _A forgotten account I can't get into_ → Part 2 (recovery) and Part 5 (the old email).
- _A people-search or data-broker profile_ → Part 3.
- _A page on someone else's site, or a post by someone else_ → Part 2 (asking others) and Part 4.
- _A search result or a cached or archived copy_ → Part 4.
- _An old email, phone number or handle doing the joining_ → Part 5.
- _A breach with a password I still use_ → Part 13.
- _Photos carrying GPS or device data_ → Part 12.

**Breaking a join:** the audit shows the handle `darkfox92` on a forum, Steam and GitHub, and GitHub's commits carry your real email. That one email joins your pseudonym to your name. Fix: delete the old forum account (Part 2), switch GitHub to the noreply address (Part 14), rewrite the old commits (Part 4.1), then re-search the handle in a month.

## 2. Full Account Deletion — Remove the Posts, Then the Account

**Philosophy.** Now actually erasing the mapped out OSINT attack surface gracefully ! 
### 🧠 Delete vs. deactivate

- **Deactivate** hides your profile. The data stays and you can come back.
- **Delete** starts a removal process, often after a grace period (about 30 days on several big platforms). **Do not log in during that window**; doing so can cancel the deletion.
- **Backups and caches** linger for a while after deletion, and third parties (search engines, archives, scrapers, people with screenshots) keep their own copies. Part 4 deals with those.

### Before you delete (every platform)

1. **Start from your exposure register** (Part 1). Add any account it missed: search your inboxes for 'welcome', 'verify' and 'confirm your account', and check your _Sign in with Google/Apple/Facebook_ lists and your password manager.
2. **Save what you want.** Every platform has a _download your data_ option. Use it first.
- Right after this step delete the posts too ! otherwise ghost posts would be left as it is ! 
1. **Preserve evidence if needed.** If you were harassed or are in a dispute, screenshot and save before you delete.
2. **Hand over what others depend on.** Transfer or archive GitHub projects that other people use.
3. **Find the deletion page.** [Just Delete Me](https://justdeleteme.xyz) is a directory of direct deletion links graded easy, medium, hard or impossible. [AccountKiller](https://www.accountkiller.com) is similar.
4. **Make sure you can receive mail** at the address each account was registered with, so confirmation emails reach you. Close old mailboxes last (Part 5).

### Tools for deleting content (your own accounts only)

- **Userscript managers:** [Violentmonkey](https://violentmonkey.github.io) (open source) or Greasemonkey (Firefox) run small scripts inside your browser. Read a script's source before installing it; a malicious script can steal your session.
- **Use the platform's bulk tool first:** Facebook _Manage activity_; Instagram _Your activity_.
- **Data-broker removal is not account deletion.** Opt-outs from people-search sites and data brokers are Part 3.
- Bulk-deleting scripts may break a platform's terms, can fail when the site changes, and belong only on accounts you own.

### GitHub

1. Download data: _Settings → Account → Export account data_.
- I would say its way better to download data and then re push to some another account since export to another account would create a link anyways 

2. Delete or transfer repositories with the [GitHub CLI](https://cli.github.com):
```
gh auth login
gh auth refresh -s delete_repo
gh repo list YOUR-USERNAME --limit 200
gh repo delete YOUR-USERNAME/old-repo --yes
```

The `--yes` flag skips the confirmation prompt, so double-check the repo name.

**Then:**

- Delete gists and packages and **leave organisations.**
- [Settings → Account → Delete your account](https://github.com/settings/admin). The username can later be taken by someone else, so decide whether you would rather keep an empty shell.
- Forks, clones and cached copies of public repos remain (Part 4).

### LinkedIn

1. _Settings & Privacy → Data privacy → Get a copy of your data_.
2. There is no official bulk delete. Remove posts, comments and recommendations by hand (_Me → Posts & activity_). Blank the profile text and photo first so cached copies show less.
3. _Settings & Privacy → Account preferences → Account management → Close account_. There is a short window to reactivate.
4. Afterwards search your name plus LinkedIn and use the Part 4 steps for stale results.

### Instagram, Facebook and Threads

1. **Accounts Center → Your information and permissions → Download your information**.
2. **Bulk delete:** Instagram _Your activity → Posts / Reels / Comments / Likes → Select → Delete_; Facebook _Settings → Activity log → Manage activity_.
3. **Delete the account:** _Accounts Center → Personal details → Account ownership and control → Deactivation or deletion_ → choose the profile → **Delete account**. Meta usually holds the deletion for about 30 days and finishes within about 90 days of the request.
4. **Threads:** the Threads profile can be deleted on its own without deleting Instagram (_Threads → Settings → Account → Delete or deactivate profile_). Turn off Fediverse sharing _before_ deleting posts (Part 14).
5. Revoke 'Log in with Facebook' sites first, or you may lose access to those accounts.

### X (Twitter)

1. _Settings → Your account → Download an archive of your data_.
2. X has no bulk-delete button. Options: a userscript such as [TweetXer](https://github.com/lucahammer/tweetXer) (uses your downloaded archive) or a paid service. Be wary of any tool asking for your password.
3. Remove likes and bookmarks. Deleting a DM removes it only for you.
4. **Deactivate:** _Settings → Your account → Deactivate your account_. X keeps the account for about 30 days before final deletion.

### Telegram

1. Export: Telegram Desktop → _Settings → Advanced → Export Telegram data_.
2. Transfer or delete groups and channels you own.
3. **Delete now:** open [my.telegram.org/delete](https://my.telegram.org/delete), enter your number, confirm with the code sent inside Telegram, then click **Delete My Account**. This is immediate and permanent. The in-app option only sets a timer (_Settings → Privacy and Security → Delete my account if away for_).
4. Messages you sent stay in other people's chats.

### Content other people posted about you

You can't delete another person's post, but you can ask.

1. **Ask first.** Send a short, friendly message with the exact URL, what you want removed (a tag, a photo, an address) and why. Most people say yes.
2. **Use the platform's tools:** untag yourself, or report the post as a privacy violation (look for _Report → Privacy_ or a form for photos of you). Many platforms have a specific route for doxxing.
3. **Site owners:** for a page on someone's own site (a club roster, a school newsletter, an old forum), email the webmaster; the footer or WHOIS usually names a contact.
4. **If they refuse or ignore you:** log the URL in your register, ask search engines to drop it (Part 4) and consider your local privacy law (Part 3 lists the main ones). If the content is harassment, screenshot it _before_ you ask.

### Verify it worked

- After the grace period try to log in (it should fail) and search the web for your username.
- Keep the 'deletion complete' emails.
- Search Google with `site:instagram.com yourusername` (swap in each platform).

**Synthesis.** Content first, account second, verification third. A deleted account with its posts still cached is a half-finished job, so move straight on to Part 4.

## 3. Data Brokers and People-Search Sites — Opt Out of the Resale of You

**Philosophy.** Data brokers compile public records, purchases, app data and leaked data into profiles, then sell them or display them to anyone. People-search sites are the public-facing kind: a name, an age, an address history, phone numbers and relatives, a few clicks from your front door. They are the most direct route from your name to your home, and they re-list you, so this part is a recurring chore, not a one-time fix.

### 🧠 Three kinds of broker

1. **People-search and background-check sites** (such as Spokeo, Whitepages, BeenVerified, Intelius and Radaris). Anyone can search them. Most are US-focused, but similar sites exist elsewhere.
2. **Marketing and ad-tech brokers.** They sell lists of consumers by interest, income and location. You never see these profiles; you opt out through their privacy portals.
3. **Risk and identity aggregators.** They sell data to lenders, landlords and employers. Credit bureaus are regulated separately, and in countries that have them, a credit freeze limits misuse.

### 🧠 Legal levers (verified October 2026)

- **California residents:** the state's **DROP** platform at [privacy.ca.gov/data-brokers](https://privacy.ca.gov/data-brokers/) sends one deletion request to every registered data broker (500+). Consumer requests opened on 1 January 2026, and from 1 August 2026 brokers must check DROP at least every 45 days and process requests. You verify your identity through the state's identity process and get a DROP ID to track the request. It covers _registered brokers only_, so people-search sites outside the registry still need manual opt-outs.
- **Other US states:** Oregon, Texas and Vermont also require broker registration but have no central portal. Check your attorney general's data broker registry and write to brokers one by one.
- **EU and UK:** the GDPR gives you a right to erasure and a right to object. Send a written request to the broker's privacy contact and keep a copy.
- **Elsewhere:** many countries have some deletion right. Search your data protection regulator's site for 'erase my personal data'.
- **No law applies to you?** Use the broker's own opt-out form. Most large ones have one.

### Default pick and decision tree

- _I have time and want to hand over as little as possible_ → **DIY** with the [Big Ass Data Broker Opt-Out List](https://github.com/yaelwrites/Big-Ass-Data-Broker-Opt-Out-List) and Privacy Guides' [data broker removals](https://www.privacyguides.org/en/data-broker-removals/) page. IntelTechniques also publishes a free [data removal workbook](https://inteltechniques.com/workbook.html).
- _I'm in California_ → **DROP first** (free, one request), then DIY the people-search sites it does not cover.
- _I'm short on time and in the US_ → **a paid service.** Privacy Guides lists EasyOptOuts at about $20 a year at last check. [DeleteMe](https://joindeleteme.com), [Optery](https://www.optery.com) and Incogni are the bigger names. Check current prices and terms.
- _I'm in the EU or UK_ → **GDPR letters** to the sites your audit found, largest first.
- _The profile shows a minor, or a threat is active_ → do DIY removals for the Critical profiles today, and the rest later.

Paid services trade your identity data and money for time. They only cover the brokers they know about, they cannot remove public records, and you should check that their terms say your data is used _only_ for opt-outs. Free DIY costs hours but hands over the least.

### 🛠 Do this: The DIY opt-out loop (5 to 10 minutes per site, plus waiting)

1. **Use the removal mailbox** from Part 1's pre-flight, never your main address.
2. **Start from your register.** Take the broker sites your audit found listing you, then work through the opt-out list for the rest, biggest first.
3. **For each site:** search for yourself and copy the profile URL. Find the site's own opt-out page (the list links to it directly; otherwise search '[site] opt out' and use only the site's own domain, because copycat 'we remove it for you' pages exist). Submit the profile URL and confirm through the link emailed to the removal mailbox. Log the date and a screenshot of the confirmation in the register.
4. **Give the minimum.** A name, the profile URL and an email is usually enough. Do not create an account, do not pay to 'unlock' a record, and ignore the upsell animations. If a site demands an ID upload, black out everything except what is required, or skip it and email its privacy contact instead.
5. **Wait two weeks, then search again.** If the profile is still there, resubmit, name the law that covers you, and escalate to your regulator if you are ignored.
6. **Re-check every three months.** Brokers re-list people and new sites appear. Put the date in the register and in the monthly habit (Part 27).
7. **Marketing brokers:** use the 'Do Not Sell or Share My Personal Information' link on their sites, and turn on the Global Privacy Control signal in your browser (Part 11).

### Public records you cannot opt out of

Property, court, voter and company records are public by law, so no broker opt-out touches them. What you can do:

- Put **WHOIS privacy** on every domain you own.
- Ask the registry or court whether it offers suppression. Some jurisdictions run address-confidentiality programmes for people at risk of stalking or abuse.
- Use a different address going forward. Guide 3 covers location privacy (mailboxes, pickup points and similar).

### ⛔ Boundary — theory only

No steps here for using these sites to look up other people, for filing opt-outs _as someone else_ or on a record that is not yours, or for using the removal process to confirm that a person is in a database. Submit requests only for yourself, or for a family member who has asked you to. **Why:** opt-out forms usually verify nothing more than an email address, which is why filing them for others is abuse, and every removal request tells the broker which details are real.

### Worked chain

**Home address on three sites:** the audit finds your address on three people-search pages and a search result → opt out of each (steps 1 to 4) → request removal of the search result (Part 4) → re-check in two weeks → re-check every three months.

**Synthesis.** Brokers earn money by knowing where you live, and they re-list you. Remove the Critical profiles first, use the law where you have it, give each opt-out the least data it needs, and put the re-check in your calendar, because the first removal is only the start.

## 4. Archives — Getting Your Stuff Out of the Internet's Memory

**Philosophy.** Deleting the original is like burning a letter after the photocopiers have run. The web keeps many copies: search-engine caches, the Wayback Machine, other archives, code mirrors, scrapers and screenshots. Total removal cannot be guaranteed, so the goal is to _remove the copies you can reach, make the rest hard to find, and be honest about what remains._

### Order of operations

1. Delete or edit the **original** (Part 2).
2. Ask **search engines** to drop the page and its cached copy.
3. Ask **archives** to remove snapshots.
4. Clean **code history** (git).
5. Re-check in a month, and re-run the self-audit (Part 1, Steps 2 to 5).

### Search engines first

- **Google, outdated content:** once the original is gone or changed, use [Refresh Outdated Content](https://search.google.com/search-console/remove-outdated-content) to clear the cached snippet.
- **Google, Results about you:** [myactivity.google.com/results-about-you](https://myactivity.google.com/results-about-you) monitors Search for your phone number, email and address and lets you request removal in one click. As of February 2026 the US rollout also covers ID numbers such as passport or licence numbers. Availability varies by country; where it isn't offered, search for Google's personal-information removal form. The trade-off: to monitor, it needs you to give Google your contact details. Removal here hides the result from Google only; the page still exists.
- **Bing:** use the content removal tool in Bing Webmaster Tools. DuckDuckGo and other engines draw on those indexes, so removing the page at the source matters most.
- **Your home on Street View or a map:** use _Report a problem_ on the image and ask for it to be blurred. Other map providers have their own forms.

### 4.1 Deleting git history

**🧠 Theory.** Git remembers every commit forever. Deleting a file in a new commit does not remove it from history, and anyone who cloned or forked still has it. **If it was a secret (key, token, password), assume it is already stolen: revoke and rotate it first.** The rotation is the real fix; history rewriting is clean-up, not safety.

**🛠 Remove a file from all history** with [git-filter-repo](https://github.com/newren/git-filter-repo) (BFG Repo-Cleaner is an older alternative):

```
git clone --mirror https://github.com/YOU/REPO.git
cd REPO.git
git filter-repo --invert-paths --path secrets.txt
git remote add origin https://github.com/YOU/REPO.git
git push --force --all
git push --force --tags
```

**Rewrite the real email baked into old commits.** Create `mailmap.txt` containing one line:

```
Your Handle <ID+USERNAME@users.noreply.github.com> <old.real@email.com>
```

Then run `git filter-repo --mailmap mailmap.txt` in a fresh clone and force-push as above. Switch to the noreply address going forward (Part 14).

Afterwards:

1. Contact [GitHub Support](https://support.github.com) to purge cached views and pull-request references to the old commits. GitHub's guide is called _Removing sensitive data from a repository_ (docs.github.com).
2. Tell collaborators to re-clone, and never merge old clones back in.
3. Forks and clones are not affected. If it matters, delete the repo and recreate it fresh, or ask fork owners to remove theirs.
4. Public code is also mirrored by archives such as [Software Heritage](https://www.softwareheritage.org). Check their site for their removal process.

### 4.2 Internet Archive (Wayback Machine)

Verified October 2026 from the Archive's own help page:

1. Email **info@archive.org**.
2. Include the exact URL(s), the date range you want excluded, the period you controlled the site or account (if relevant), and anything that helps them understand the request, such as 'this snapshot exposes my home address'.
3. They review it and **make no guarantee** about the outcome.
4. For non-consensual intimate images there is a separate address, **takeitdown@archive.org**, described on the same help page. Copyright claims go through their Copyright Policy.
5. `robots.txt` stops future crawling but is unreliable for removing what is already archived. If you own the domain, offer proof of ownership (for example a short text file on the site).

Template:

```
Subject: Exclusion request - example.com/page

Hello, I am the owner of [site/account] and request exclusion of the
following URLs from the Wayback Machine for all dates: [URLs].
I controlled this [site/account] from [date] to [date].
Reason: [personal data / safety / no longer want this public].
Thank you. [Name or alias, contact email]
```

### 4.3 The smaller archives: Archive.today, Ghostarchive and Megalodon

**🧠 Theory.** These are separate from the Internet Archive, and I could not verify a published removal process for any of them. Expect slow, uncertain or no replies, and remember that small archives can vanish without notice.

- **Archive.today** (archive.is, archive.ph and mirrors): the operator is anonymous, the service is controversial, and some sites and networks block it.
- **Ghostarchive:** saves pages on demand without a login.
- **Megalodon** (megalodon.jp): a Japanese archive. Use your browser's translator on its contact or FAQ pages.

**🛠 What to try (the same for all three).**

1. Remove the original and ask search engines to drop the page first.
2. Look for a contact or removal link on the archive's own site.
3. Send the exact snapshot URL, the reason (it contains your personal data) and proof the content is yours. The template from 4.2 works with the names swapped.
4. Keep the request short and your expectations low. If a snapshot puts you at real risk, go through a lawyer or a data-protection complaint rather than repeated emails.

> ⛔ **Boundary — theory only.** File requests only about content you own or about yourself. Do not file false copyright or legal claims (for example DMCA abuse) to suppress criticism, reporting or other people's lawful content. That is illegal and harms others. In theory, archives weigh the public-interest value of records such as journalism and court documents, so personal data is far more likely to be removed than coverage of wrongdoing. Legitimate levers are privacy law (such as GDPR erasure rights) and copyright over your own work.

**Synthesis.** You cannot un-publish the internet, but you can make your data far harder to find and cheaper to ignore. Remove the source, cut off the caches, ask the archives politely, and rotate any secret that ever touched public history.

## 5. Retiring Old Verification Vectors — Emails, Usernames and Phone Numbers

**Philosophy.** Old emails, usernames and numbers are doors you stopped using but never locked. Attackers and data brokers love them because they still open accounts, still receive password resets and still link to you. Close each door properly, in the right order, so you neither lock yourself out nor hand the key to a stranger.

**Dependencies.** Migrating accounts needs a new alias (Part 7) and an authenticator app (Part 13). If those aren't set up yet, do the inventory and triage below now and finish the migration after Part 13. Never close an old mailbox or give up a number until everything on the register that depends on it has been moved or deleted.

### 🧠 The three traps

1. **Recycling.** Some email providers and carriers eventually release dormant addresses and numbers. A stranger can register your old address and receive password-reset emails for any account still tied to it.
2. **Orphaned accounts.** Delete an email before the accounts attached to it and you can no longer prove ownership of them.
3. **Squatted names.** Delete a username and someone can take it and impersonate the old you.

### 5.1 Emails — discard procedure

1. **Inventory.** Start from your exposure register (Part 1). In the old mailbox search for _welcome_, _verify_, _confirm_, _receipt_, _password reset_ and _subscription_. Add each service to the register.
2. **Triage each account:** _delete_ it (Part 2), _migrate_ it (change the login email to a new alias), or _abandon_ it (only for unimportant sites).
3. **Priority order:** banks, cloud storage, domain registrars, app stores, social media, shops with saved cards.
4. **Forward** the old address to your alias for 6 to 12 months. Every forgotten account that emails you gets deleted or migrated.
5. **Update recovery addresses** on every important account.
6. **Close the mailbox last.** Check the provider's rules on dormant addresses. Gmail addresses are not reused after deletion; for other providers, check their policy. If you can't confirm the address won't be recycled, keep it alive but empty, with a long random password and 2FA.
7. Export everything you want first (mail, contacts, calendar).

### 5.2 Usernames — discard procedure

1. Decide per platform: **delete and let go**, **keep a locked empty shell** (stops impersonation), or **rename to random characters, then delete**.
2. For a public persona or professional name, prefer the _locked shell_: delete content, remove photo and bio, set a long random password, turn on 2FA, keep the account.
3. For private throwaway handles, rename to random characters and delete.
4. **Never reuse an old handle** for a new compartment (Part 9). Handles are join keys (Part 7).
5. Search your old handles in quotes with a search engine, and on the sites where you used them, to find leftovers; ask those sites to remove stray mentions.

### 5.3 Phone numbers — discard procedure

1. **List every account tied to your number:** banks, payment apps, messengers, social media, delivery apps, loyalty cards, 2FA.
2. **Move 2FA first** to an authenticator app or security key (Part 13), then remove the number or replace it with one you keep for official purposes only.
3. **Hide it from discovery** on Signal and Telegram (Parts 16 and 14), on WhatsApp (_Privacy_), and in contact lookups on Instagram, Facebook and X (Part 14).
4. **Clean up contact uploads** on Meta, X, LinkedIn and Telegram.
5. **Opt out of people-search sites** that list your number (Part 3).
6. **If you are giving the number up,** carriers recycle numbers after a cooling-off period that varies by country and carrier. Make sure nothing still uses it for resets _before_ you cancel or port out.
7. Keep one **official number** for banking, government, health and the people who must reach you. Treat every other use as optional.

> ⛔ **Boundary — theory only.** No steps here for obtaining numbers or SIMs registered under someone else's identity, SIM farms, or exploiting number recycling. In theory: SIM registration is a legal requirement in many countries, and recycled numbers are a real attack route (reset and SIM-swap takeovers), which is exactly why step 6 matters. The lawful way to reduce exposure is _fewer accounts tied to your number_, not hiding who owns it.

**Synthesis.** An old identifier is only dangerous while something still trusts it. Move trust off it, watch it for a while, then retire it, and never let a stranger inherit a door you left unlocked.

## 6. IP Protection — Don't Hand Out Your Return Address

**Philosophy.** Your IP address is the return address printed on every envelope your devices send. A VPN is a mail-forwarding service: the website sees the forwarder's address, and your internet provider sees only that you talk to the forwarder. You have not become invisible; you have _moved the trust_ from your ISP and the café Wi-Fi to one company whose business model is not selling you.

### 🧠 What an IP address gives away

- Your rough location (usually city level) and your ISP.
- A stable tag that links your visits across days, even without cookies.
- With a court order, your ISP can map it to your name and address.

### 🧠 What a VPN does _not_ do

- It does not stop cookies, logins or browser fingerprinting (Parts 9 and 11).
- It does not make you anonymous if you sign into accounts that know your real name.
- It does not protect you from malware or phishing.
- It is not Tor. If you need to hide from the website itself, use Tor Browser instead.

> 🗣 **Author's take.** I treat a paid VPN like a seatbelt: cheap, boring, always on. 'Free' VPNs are the opposite: if you are not paying, your traffic is the product. Spend about five euros a month and move on. The goal is not to be a ghost; it is to stop strangers building a free profile on you.

### Default pick: Mullvad

Verified October 2026:

- **No email needed.** You press a button and get a random account number. That number is your whole login.
- **Anonymous payment.** Cash sent in an envelope, Monero, Bitcoin and other crypto are accepted, plus card and PayPal.
- **Flat price.** About €5 per month, no long-term lock-in.
- **WireGuard only.** Mullvad removed OpenVPN support in January 2026.
- **Independent audit.** Mullvad has had a third-party audit of its account and payment systems.

### Decision tree

- _I want the most anonymous signup, or I want to pay with cash_ → [Mullvad](https://mullvad.net) or [IVPN](https://www.ivpn.net).
- _I want something easy with a free tier, and an email-based account is fine_ → [Proton VPN](https://protonvpn.com).
- _I need to hide from the website, not just my ISP_ → [Tor Browser](https://www.torproject.org/download/) (slower, but built for exactly this; no logins, no downloads).
- _I want my own server or traffic that looks like ordinary web traffic_ → Guide 2.
- _A stranger sold me a 'lifetime VPN'_ → no. Review sites that push the highest-commission VPN are affiliate marketing. Use [Privacy Guides' VPN page](https://www.privacyguides.org/en/vpn/) instead.

### 🛠 Do this: Mullvad in 10 minutes

1. Go to mullvad.net and choose **Generate account number**. Copy the number into your password manager _now_. If you lose it, nobody can recover the account.
2. Add time: card, crypto or **cash by post** (Part 25 covers the cash route; it takes about one to two weeks to be credited).
3. Install the app for your system from mullvad.net/en/download.
4. Open **Settings** and turn on: _Launch app on start-up_, _Auto-connect_, and the app's **kill switch / lockdown mode** (the label varies by version; it is the option that blocks the internet when the VPN drops).
5. Turn on the **DNS content blockers** (ads, trackers, malware) if you want network-level blocking too.
6. Connect, then test at [mullvad.net/en/check](https://mullvad.net/en/check). You should see that you are connected to Mullvad.
7. Run a second check at [browserleaks.com](https://browserleaks.com) (WebRTC and DNS sections). Your real IP must not appear anywhere.

**Phone:** install the Mullvad app, then on Android go to _Settings → Network → VPN → ⚙ → Always-on VPN_ and enable _Block connections without VPN_. On iPhone keep the app on auto-connect; iOS has no system-wide kill switch.

**Private DNS on Android (works even without a VPN):** _Settings → Network → Private DNS_ → hostname `dns.mullvad.net` or `dns.quad9.net`.

### Fallbacks

- **IVPN:** similar anonymity model, also accepts cash.
- **Proton VPN:** needs an email address. A [Proton](https://proton.me) signup has usually offered a CAPTCHA route instead of phone verification, but that can change.
- **Your own VPN on a rented server (a VPS)** is a Guide 2 topic.

### Worked chain

**Maximum cheap anonymity:** cash by post → Mullvad account number → app with kill switch → Mullvad Browser or hardened Firefox (Part 11) → pseudonymous accounts only (Part 7).

**Synthesis.** A VPN is a trust transfer, not a cloak. Pick one provider that is cheap, audited and does not need your identity to sell you the service, keep it on by default, test it for leaks once, and never undo the benefit by logging in as your real self in the same session.

## 7. Never Verifying — Don't Hand Over the Real You

**Philosophy.** Your real name, email and phone number are _join keys_. Databases that never met each other get merged because they all contain the same three values. Give each service the minimum it needs to work, and give different services different values, so one leak cannot be joined to the next.

### 🧠 Pseudonym vs. impersonation (the rule of this part)

A **pseudonym** is a made-up name used for your own privacy, on services that allow it. That is legal almost everywhere, historically normal, and exactly what this part teaches. **Impersonation** (pretending to be a real person), **fraud** (lying to get money, access or benefits) and **evading a legal identity check** are different things, and they are ⛔ theory-only below.

### The four rules

1. **Name.** One pseudonym per compartment (Part 9). Never reuse a pseudonym between your public and private lives.
2. **Email.** Every service gets its own alias.
3. **Phone number.** Avoid giving it out. A number is the most reliable join key you own, and a SIM is usually registered to your ID.
4. **Accounts.** Ask 'do I need an account?' first. Guest checkout, reading while logged out, or using a website instead of the app all create less data than registering. Never use 'Sign in with Google/Apple/Facebook': it links the new account to your main one.

> 🗣 **Author's take.** Most of my accounts exist because a website made me feel I had to register. The best account is the one I never made. If a service demands my phone number just to read content, the service loses a user.

### 🛠 Do this: Build a persona sheet

1. In your password manager (Part 13) create an **Identity** entry for each compartment. Bitwarden, KeePassXC and Proton Pass all support this.
2. Choose a pseudonym. Invent one, or use [Fake Name Generator](https://www.fakenamegenerator.com/) for ideas. **Use the name only.** Ignore its other fields, do not copy numbers that look like government IDs, and never use output that happens to match a real person.
3. Create an email alias for that persona (below).
4. Save the persona's username, alias and what it is _for_.
5. Keep one persona per purpose: _shopping_, _dev work_, _forums_, _throwaway_.

### 🛠 Do this: Email aliases in 10 minutes

Default: [SimpleLogin](https://simplelogin.io) (owned by Proton) or [Addy.io](https://addy.io) (open source). Alternatives: [Firefox Relay](https://relay.firefox.com), Proton Pass aliases, Apple _Hide My Email_ on iPhone.

1. Create a mailbox that is _not_ tied to your real name: [Proton Mail](https://proton.me/mail) or [Tuta](https://tuta.com) are the usual picks.
2. Create an alias account and link it to that mailbox.
3. Install the browser extension or the password-manager integration.
4. At each signup press _create alias_, paste, done. Replies go out from the alias.
5. If an alias starts getting spam, you know who sold you out. Disable that one alias; nothing else breaks.

### Phone numbers: decision tree

- _Does the service work without one?_ → skip the number.
- _It wants a number for 2FA_ → switch 2FA to an authenticator app or a security key (Part 13).
- _It wants a number to sign up_ → first ask whether you need the service at all. For messaging, Signal is testing signup without a phone number on Android (Part 16); treat it as beta.
- _You need a second number for legitimate separation_ → a prepaid or eSIM line bought in your own name from a licensed carrier hides the number from _websites_, though not from the carrier. Many services block VoIP numbers.

### ⛔ Boundary — theory only

Out of scope here, with no steps: using someone else's identity or documents; buying or renting SIMs registered to other people; SMS-verification farms and 'burner code' websites; mass-creating accounts; evading a ban; getting around a bank, SIM or government identity check ('KYC').

**Why.** Platforms use phone numbers as a cheap anti-abuse barrier, and KYC rules exist to stop fraud and money laundering. Public 'receive SMS' sites are readable by _everyone_, so your 'private' code is public; the numbers get burned quickly; and the whole ecosystem is full of scams. The legal alternative is _less verification surface_, not _fake verification_.

### Platform rules to know

- **LinkedIn and Facebook** require your real name in their terms. If you must use them, use your real name but tighten every setting (Part 14) and limit what you post.
- **GitHub, X, Reddit, Mastodon and most forums** allow pseudonyms.
- **Banks, government, healthcare and taxes** always use your legal identity. Never put an alias there.
- **Do not lie about your age** where age limits exist for child safety or law.

### Worked chain

**Shopping:** guest checkout → alias email → pick up in store and pay cash (Part 25).

**Synthesis.** Verification is a tax on your privacy: every extra identifier you hand over is a join key waiting to be exploited. Spend aliases freely, spend your real identifiers only where the law demands them, and when a service insists on more than it needs, treat that as information about the service.

## 8. Windows → Linux — Leaving the Telemetry Mothership

**Philosophy.** Your operating system sees everything: every file, every keystroke, every network connection. Privacy tools sitting _on top of_ an OS you cannot inspect are a lock on a door the landlord has a key to. Moving to an open-source OS puts the whole stack in public view.

> 🗣 **Author's take.** I left Windows because it stopped feeling like _my_ computer: ads in the Start menu, a steady push toward cloud sign-in, a growing list of data-collecting features, and a support end date (Windows 10 stopped getting free security updates in October 2025) used as a nudge to buy new hardware. Linux is not perfect and not magic, but it is an OS that works for me. This is my opinion. Your hardware, games and workflow may disagree, and that is fine.

### 🧠 Honest expectations

- Linux has a _much smaller_ telemetry footprint, but it is not automatically private: Ubuntu ships optional error and usage reporting, and any app can phone home.
- Desktop Linux is less sandboxed than a phone. Updates, disk encryption (Part 26) and sensible habits matter.
- A few things still work badly: some anti-cheat games, Adobe apps, certain Office-compatibility edge cases. Check first (below).

### Default pick and decision tree (verified October 2026)

- **Default: [Linux Mint](https://linuxmint.com/download.php) 22.3 Cinnamon.** Windows-like layout and long support (through April 2029). Mint 23, built on Ubuntu 26.04 LTS, is planned for around Christmas 2026. Don't wait; install 22.3 now.
- _Want the biggest tutorial ecosystem_ → [Ubuntu 26.04 LTS](https://ubuntu.com/download/desktop) (supported until 2031).
- _Want a customisable Windows-style desktop_ → Kubuntu 26.04 LTS or [Zorin OS](https://zorin.com/os/).
- _Want newer software and don't mind a bit more learning_ → [Fedora Workstation](https://fedoraproject.org/workstation/).
- _Gaming first (NVIDIA)_ → look at Pop!_OS or Fedora, and check your games on [ProtonDB](https://www.protondb.com).
- _Journalist-level threat model_ → Tails or Qubes OS (Guides 2 and 3).
- **Do not distro-hop in your first month.** Pick one, learn it, move on.

### 🛠 Do this: Switch safely

**A. Prepare (while still on Windows)**

1. **Back up** documents, photos, browser data and a password-manager export to an external drive. Open a few files from the backup to prove it works.
2. List your must-have apps and check each on [AlternativeTo](https://alternativeto.net). Typical swaps: LibreOffice (Office), GIMP or Krita (Photoshop), Kdenlive (video editing), OBS Studio (recording), Thunderbird (mail), VLC (media).
3. **Turn off BitLocker / Device Encryption** (Settings → Privacy & security → Device encryption) and wait for decryption to finish. You need this if you will resize partitions for dual-booting.
4. **Turn off Fast Startup** (Control Panel → Power Options → _Choose what the power buttons do_ → untick _Turn on fast startup_). It leaves the Windows disk half-asleep, which breaks dual-booting.

**B. Make the installer USB**

1. Download the ISO from the _official_ page only.
2. Verify the checksum. In Windows PowerShell: `Get-FileHash .\linuxmint-22.3-cinnamon-64bit.iso -Algorithm SHA256` (your file name will match your download). On Linux: `sha256sum file.iso`. Compare the result with the hash on the download page.
3. Write it to a USB stick (8 GB or more) with [Ventoy](https://www.ventoy.net), [Rufus](https://rufus.ie) or [balenaEtcher](https://etcher.balena.io). This erases the stick.

**C. Try before you install**

1. Reboot and tap the boot-menu key (usually F12, F10, F9, Esc or F2, depending on the laptop brand).
2. Boot the USB in _live mode_. Test Wi-Fi, sound, webcam, trackpad and sleep. If something fails, search [your laptop model] + Linux before you install.

**D. Install**

1. Decide: **wipe the whole disk** (simplest, most private) or **dual-boot** (keep Windows for a few apps; more complex).
2. Double-click _Install Linux Mint_. For a wipe choose _Erase disk and install_ → _Advanced features_ → **encrypt the new installation** (this is LUKS, Part 26). Pick a strong passphrase.
3. Choose timezone, create your user, wait, restart, remove the USB.
4. Secure Boot can stay on for Mint, Ubuntu and Fedora. With NVIDIA drivers you may be asked to enroll a key on the next boot; follow the blue screen's instructions.

**E. First 20 minutes after install**

1. Open _Update Manager_ and apply everything (or run `sudo apt update && sudo apt full-upgrade`).
2. Open _Driver Manager_ for GPU and Wi-Fi drivers.
3. Turn on the firewall: `sudo ufw enable`.
4. Set up _Timeshift_ (system snapshots) before you change anything else.
5. Install apps from the Software Manager or [Flathub](https://flathub.org). Flatpak apps are sandboxed; use _Flatseal_ if you want to see or tighten their permissions.
6. **Ubuntu only:** opt out of reporting with `ubuntu-report -f send no` and `sudo systemctl disable --now whoopsie`.
7. Install a browser (Part 11) and your VPN (Part 6).

### Fallbacks

- Not ready to leave? **Dual-boot** with an encrypted Linux install and use Windows only for specific apps, with a _local_ account and telemetry at the minimum.
- Need a Windows app? Try [Wine](https://www.winehq.org) or Bottles, or a Windows VM (Part 9), before dual-booting.
- Gaming: Steam has Proton built in; [Lutris](https://lutris.net) and [Heroic](https://heroicgameslauncher.com) cover other stores. Competitive games with kernel-level anti-cheat often _won't_ run, so check before you commit.

### Worked chain

**Old laptop first:** put Mint on a five-year-old laptop for a month, then decide about your main machine.

**Synthesis.** You don't leave Windows to feel superior; you leave it because privacy is easier to build on an OS you can inspect. Start with a live USB, back up before you touch a partition, encrypt the new install, and let your habits settle for a month before you judge the result.

## 9. Compartmentalization — Don't Keep Everything in One Room

**Philosophy.** Ships have watertight compartments so one hole doesn't sink the vessel. Your digital life needs the same: if your gaming account, bank, work and side project all live in one browser, on one device, under one email, a single slip sinks all of them. You do not need Qubes OS for this. Cheap separation you actually use beats perfect separation you never set up.

### 🧠 The compartment map

Make three or four compartments. Give each its own browser profile or VM, its own email alias and pseudonym (Part 7), and ideally its own password-manager folder.

- **Real-name / sensitive:** banking, government, healthcare, taxes. Boring, locked down, almost no extensions.
- **Daily life:** shopping, streaming, social media under an alias.
- **Public persona:** the identity you post under (creator, developer). Real or pseudonymous, but _never_ connected to the others.
- **Throwaway / risky:** unknown downloads, strange links, free trials, one-time signups.

### Default pick and decision tree

- _Lowest effort (10 minutes)_ → **browser containers or profiles.** Firefox [Multi-Account Containers](https://addons.mozilla.org/firefox/addon/multi-account-containers/), or separate browser profiles.
- _Risky downloads or untrusted software_ → **a virtual machine** ([VirtualBox](https://www.virtualbox.org) on any PC, [virt-manager](https://virt-manager.org) on Linux) or **Windows Sandbox** (Windows Pro/Enterprise).
- _Real separation between identities_ → **a separate device** (old laptop, cheap second phone) or at least a separate OS user account.
- _Phone_ → private space, work profile or separate users (below).
- _Very high threat_ → Qubes OS with Whonix (Guide 2). You can live without them for now.

### 🛠 Do this: Browser compartments in 10 minutes (Firefox)

1. Install Firefox, then add _Multi-Account Containers_.
2. Create containers named **Banking**, **Daily**, **Public persona** and **Throwaway**.
3. Right-click the extension icon → _Always open this site in…_ to pin each site to its container.
4. Remember: containers separate cookies and storage, **not** your IP address or hardware fingerprint. For stronger separation use a different browser profile (type `about:profiles`) or a different browser.

### 🛠 Do this: A disposable VM on a PC (VirtualBox)

1. Install VirtualBox and download a Linux ISO (Part 8).
2. Click _New_, name it Throwaway, choose the ISO, give it 4 GB RAM, 2 CPUs and 40 GB disk.
3. Install the guest OS, then **take a snapshot** and call it clean.
4. Use the VM for untrusted downloads and links. For a suspicious PDF or installer, turn networking off first. Afterwards: _Snapshots → Restore clean_.
5. Hardening: _Settings → General → Advanced_ → shared clipboard and drag-and-drop **Disabled**; add no shared folders; network mode _NAT_.
6. **Windows Sandbox (Pro/Enterprise):** _Turn Windows features on or off_ → tick _Windows Sandbox_ → restart → open it from the Start menu. Everything inside is deleted when you close the window.

### 🛠 Do this: Phone compartments

- **Android 15 and newer:** _Private Space_ (Settings → Security & privacy → Private Space) gives a locked, separate area for sensitive apps.
- **Android, any version:** install [Shelter](https://github.com/PeterCxy/Shelter) from [F-Droid](https://f-droid.org) to clone apps into a work profile, and install risky apps only there.
- **Maximum on Android:** GrapheneOS on a Pixel gives multiple user profiles, each with its own encrypted storage. Guide 2 has the full switch-over.
- **iPhone:** lock or hide apps behind Face ID (iOS 18 and newer), keep throwaway apps under a separate Apple Account, and turn on _Lockdown Mode_ only if you face targeted threats.
- **Separate device:** an inexpensive second phone used only for banking apps is the strongest cheap compartment there is.

### Worked chain

**Public persona:** persona account → its own browser profile → its own VPN server location → its own email alias. That profile never opens your personal mail, bank or social sites.

**Synthesis.** Compartments cost minutes and pay out on the one bad day you will eventually have. Start with browser containers today, add a VM when you begin downloading things, and add a separate device only when the stakes justify it.

## 10. MAC Address Spoofing — Stop Wi-Fi Networks from Recognising Your Hardware

**Philosophy.** Your IP address changes between networks, but your network card's MAC address is a serial number that follows the device. Every shop, airport and hotel Wi-Fi that logs it can recognise you on your next visit. A random MAC per network makes each visit look like a stranger.

### 🧠 What a MAC address is (and isn't)

- A 12-digit hardware ID such as `A4:5E:60:12:34:56`, used only on your _local_ network. Websites never see it.
- The router, the hotspot and anyone else on the same Wi-Fi can see it. Retail analytics and 'free Wi-Fi' portals log it.
- It is not your identity by itself, but combined with a portal login (name, email, phone) it becomes one.

### Default pick: use the built-in randomisation first

Modern systems already randomise MAC addresses. Turn it on and you are most of the way there.

1. **Android 10 and newer:** _Settings → Network & internet → Wi-Fi →_ gear icon next to the network → _Privacy_ → **Use randomized MAC**.
2. **iPhone / iPad (iOS 14+):** _Settings → Wi-Fi →_ (i) next to the network → **Private Wi-Fi Address** on.
3. **Windows 10/11:** _Settings → Network & internet → Wi-Fi → Random hardware addresses_ → On (or _Change daily_). It can also be set per network in that network's properties.
4. **macOS (Sonoma or newer):** _System Settings → Wi-Fi → Details… → Private Wi-Fi address_ → _Fixed_ or _Rotating_.
5. **Linux (NetworkManager):** create the file below, then restart the service.

```
sudo nano /etc/NetworkManager/conf.d/00-macrandomize.conf
```

Put this inside:

```
[device-mac-randomization]
wifi.scan-rand-mac-address=yes

[connection-mac-randomization]
ethernet.cloned-mac-address=random
wifi.cloned-mac-address=random
```

Then run `sudo systemctl restart NetworkManager`. For a single network: `nmcli connection modify 'Cafe-WiFi' 802-11-wireless.cloned-mac-address random` followed by `nmcli connection up 'Cafe-WiFi'`. Use `random` for a new address every connection, or `stable` for one consistent random address per network (kinder to captive portals).

### 🛠 Manual spoofing (only when the built-in option isn't enough)

**Linux, macchanger:**

```
sudo apt install macchanger
ip link
sudo ip link set wlan0 down
sudo macchanger -r wlan0
sudo ip link set wlan0 up
```

Replace `wlan0` with your interface name from `ip link`.

**Windows:** Device Manager → Network adapters → your adapter → Properties → Advanced → _Network Address_ → enter 12 hex digits with no separators. The second digit must be 2, 6, A or E (example: `02AB12CD34EF`). The free [Technitium MAC Address Changer](https://technitium.com/tmac/) does the same with a GUI.

**Also change your device name.** Routers see the hostname your device announces. Names like Alex-iPhone leak who you are. On Linux run `hostnamectl set-hostname laptop`; on Windows use _Settings → System → About → Rename this PC_; on iPhone use _Settings → General → About → Name_.

### Decision tree

- _Phone only_ → the built-in toggle is enough.
- _Laptop on public Wi-Fi (café, airport, hotel)_ → built-in randomisation + generic hostname + VPN (Part 6), and never type your real name into a portal.
- _Portal insists on a name or room number_ → give only what it needs, or use your phone's hotspot instead.
- _Need the same MAC at home (router allow-list, static IP)_ → set `stable` or switch randomisation off for the home network only.
- _Captive portal refuses to load_ → try `stable`, or turn randomisation off briefly for that network.

### ⛔ Boundary — theory only

No steps here for spoofing to impersonate another device on a network you do not own, to bypass MAC allow-lists or parental controls, to dodge time or data limits on paid or hotel Wi-Fi, or to get around a ban from a network. That is unauthorised access or theft of service in many places.

**Why.** MAC filtering is weak security because addresses are visible over the air; serious networks use WPA3 or 802.1X authentication instead. Use this technique for your own privacy on networks you are allowed to use.

**Synthesis.** Randomising your MAC is a switch, not a project. Flip it on for every device, rename the device, and the places that count your visits will see a new stranger each time.

## 11. Social Media & the Web — The Anti-Tracking Toolkit

**Philosophy.** Tracking works because browsers volunteer information: cookies, IP address, screen size, fonts, language. You cannot make yourself blank, but you can block the trackers, make your browser look like everyone else's, and stop handing platforms free raw material (photos, audio, video) to mine.

### Default pick: a browser plus a blocker (verified October 2026)

- **Best all-round: Firefox + uBlock Origin.** The full uBlock Origin works on Firefox (Mozilla has said publicly that this support isn't going anywhere) and on Brave (Brave hosts it on its own servers).
- **Chrome, Chromium and Edge:** the full uBlock Origin is gone. Google pulled the last Manifest V2 extensions from the Chrome Web Store on 31 August 2026, and Edge is phasing them out. The official stopgap is **uBlock Origin Lite**, which is less capable. Switching browser is the real fix.

### Decision tree

- _I want strong protection with little tuning_ → [Brave](https://brave.com) (blocking built in) or [Mullvad Browser](https://mullvad.net/en/browser) (designed so every user looks identical; best with a VPN).
- _I want control and the best blocker_ → [Firefox](https://www.mozilla.org/firefox/) with uBlock Origin, hardened with [Betterfox](https://github.com/yokoffing/Betterfox) (easier) or [arkenfox](https://github.com/arkenfox/user.js) (stricter, can break sites).
- _I want hardened Firefox without doing the work_ → [LibreWolf](https://librewolf.net).
- _I must stay on Chrome_ → uBlock Origin Lite plus [Privacy Badger](https://privacybadger.org).
- _Phone_ → Firefox or Brave on Android with uBlock Origin or Shields; on iPhone use Safari with a content-blocker app, or Brave. Add a private DNS (Part 6).
- _I want every device on my home network covered_ → Guide 2 (network-wide blocking).

### 🛠 Do this: Anti-tracking in 15 minutes (Firefox)

1. Install Firefox. _Settings → Privacy & Security →_ set **Strict** protection and tick _Tell websites not to sell or share my data_ (Global Privacy Control).
2. Install **uBlock Origin** from the official Mozilla add-ons page. Check that the author is Raymond Hill (gorhill); fake clones exist. Source: [github.com/gorhill/uBlock](https://github.com/gorhill/uBlock).
3. Install **[ClearURLs](https://github.com/ClearURLs/Addon)** (strips tracking junk from links) and **Privacy Badger**.
4. Add **Multi-Account Containers** (Part 9).
5. _Settings → Search →_ switch the default engine to DuckDuckGo, Brave Search or Startpage.
6. _Settings → Privacy & Security → Cookies and Site Data_ → delete cookies when Firefox closes (or rely on containers).
7. Test at [coveryourtracks.eff.org](https://coveryourtracks.eff.org).

> 🗣 **Author's take: love of Brave and uBlock Origin, and why to support them.** uBlock Origin is a rarity: free, open source, maintained by one lead developer, and better at protecting you than most commercial products. Brave builds the same idea into the browser itself, so protection doesn't depend on an extension surviving Google's rules. Neither is perfect. Brave ships extras (a crypto wallet, rewards, an AI assistant) that you should switch off in Settings, and it has had controversies in the past, so treat it as 'good, not holy'. But the principle matters: tools that block tracking by default only survive if people use them, recommend them and help improve them. Install uBlock Origin on your family's computers, report broken filter rules to the project, and tell friends.

### Making your photos, audio and video harder to mine

**🧠 Theory.** AI companies and face-search engines scrape public images, voices and clips. Once something is public you cannot fully stop that, but you can raise the cost.

Tools:

- **[Glaze](https://glaze.cs.uchicago.edu)** adds small changes to artwork so image generators struggle to copy your _style_.
- **[Nightshade](https://nightshade.cs.uchicago.edu)** 'poisons' images so models trained on them without permission learn wrong associations. Same research group at the University of Chicago.
- **[Have I Been Trained?](https://haveibeentrained.com)** shows whether your images appear in common training datasets and lets you file opt-outs where offered.
- **Face cloaking:** older research tools such as Fawkes exist, but most are no longer maintained and may not work against current face-recognition systems. Treat them as speed bumps, not shields.

**🛠 Do this**

1. Choose which images are worth protecting (artwork, portfolio, face photos).
2. Run the original through Glaze (and Nightshade for the stronger effect), export, and upload _only_ the processed copy. Keep the original offline.
3. Strip metadata first (Part 12).
4. Switch off the 'use my content to train AI' toggles on every platform (Part 14).
5. On your own website, ask AI crawlers to stay away with `robots.txt`:

```
User-agent: GPTBot
Disallow: /

User-agent: CCBot
Disallow: /

User-agent: ClaudeBot
Disallow: /

User-agent: Google-Extended
Disallow: /
```

Polite crawlers obey this; impolite ones ignore it.

6. **Post fewer identifiable faces, voices and locations.** That is the only protection that always works.

**Honest limits.** This is an arms race. Effectiveness changes with every model release, the tools can slightly alter how your images look, and they do nothing for images already scraped.

> ⛔ **Boundary — theory only.** Protecting your _own_ uploads is self-defence. Deliberately poisoning other people's data, public datasets or someone else's training pipeline to cause harm is out of scope. In theory, data poisoning is an attack on the integrity of shared resources; defenders respond with filtering and provenance checks, which is why it stays a speed bump and not a weapon.

**Synthesis.** You will not become untrackable, but you can stop being _effortlessly_ trackable. Pick one browser and one blocker, keep both updated, and stop giving platforms high-quality training material you did not choose to give.

## 12. Metadata Removal — Photos, Audio, Video and Documents

**Philosophy.** Every file is a letter with a hidden second page. A photo may carry your GPS coordinates, phone model and exact time; a Word file your username; a PDF the software you used. Strip the second page before the file leaves your hands.

### 🧠 What's hiding in your files

- **Photos:** GPS location, device make and model, serial number, timestamp, editing software, an embedded thumbnail of the original (sometimes the _uncropped_ image).
- **Audio and video:** encoder, title tags, device, creation time, GPS (on phones).
- **Documents and PDFs:** author, company, revision history, comments, printer, software version.
- **Filenames and folders** such as `IMG_Alex_Home_2026.jpg`.
- **The content itself:** street signs, reflections in windows, shadows, badges and screens. No tool removes what the picture visibly shows.

### Default pick and decision tree

- _One-off share, beginner_ → [Metadata Cleaner](https://flathub.org/apps/fr.romainvigier.MetadataCleaner) (Linux, Flatpak) or the Windows built-in: right-click file → _Properties → Details → Remove Properties and Personal Information_.
- _Lots of files_ → [ExifTool](https://exiftool.org) (Windows, Mac, Linux).
- _Documents and PDFs_ → [mat2](https://0xacab.org/jvoisin/mat2) (command line, Linux; also powers Metadata Cleaner).
- _Phone_ → turn off location tagging in the camera app. Android: Scrambled Exif (on F-Droid) strips data before sharing. iPhone: _Share sheet → Options → Location off_.
- _High stakes (publishing something sensitive)_ → strip metadata _and_ re-render the image (screenshot, crop, save as a new file) _and_ review the visible content.

### 🛠 Commands (copy and paste)

Photos, strip everything (keep your originals somewhere safe first):

```
exiftool -all= -overwrite_original photo.jpg
```

Whole folder:

```
exiftool -all= -r -overwrite_original ./to-share/
```

Check the result with `exiftool photo.jpg`. It should list only basic file information.

Audio and video with ffmpeg:

```
ffmpeg -i in.mp3 -map_metadata -1 -c:a copy out.mp3
ffmpeg -i in.mp4 -map_metadata -1 -c copy out.mp4
```

PDFs and Office files with mat2:

```
mat2 document.pdf
```

This writes `document.cleaned.pdf`. mat2 can flatten or degrade PDFs as the price of cleaning, so check the result. ExifTool alone is not enough for PDFs, because it appends changes and the old data can remain in the file. In Word, Excel and PowerPoint use _File → Info → Check for Issues → Inspect Document_ and remove personal information.

**Messaging apps:** Signal removes location EXIF from images; many other apps and email do not. Strip first, send second.

**Synthesis.** Metadata is the part of the file you didn't write on purpose. Make stripping a habit at the last step before upload, then look at the picture itself, because the content leaks things metadata never will.

## 13. Basic Password Hygiene — Boring Habits That Stop Most Attacks

**Philosophy.** Most account takeovers are not movie hacking; they are reused passwords from an old breach, tried on your other accounts by bots. Boring habits beat clever tricks: unique passwords, a manager and a second factor stop the large majority of everyday attacks.

### 🧠 Why the boring stuff wins

- Breach lists are traded for free, and bots try your leaked email and password on hundreds of sites ('credential stuffing').
- Humans reuse passwords. A manager makes reuse unnecessary.
- A second factor makes a stolen password useless on its own.
- Email is the master key: whoever controls your email can reset every other account. Protect it first.

### Default pick: a password manager

- **[Bitwarden](https://bitwarden.com)** — open source, audited, free for personal use, works everywhere. The default for beginners.
- **[KeePassXC](https://keepassxc.org)** — one local file, no cloud. Best if you want nothing online; you must back the file up yourself.
- **[Proton Pass](https://proton.me/pass)** — privacy-focused, with built-in email aliases.
- Browser-built-in managers are better than reusing passwords but worse than a dedicated vault.

### 🛠 Do this: Set it up in 20 minutes

1. Create the vault. The **master password is the one password you memorise.** Make it a six-word passphrase picked with real dice from the [EFF word lists](https://www.eff.org/dice), and never reuse it anywhere.
2. Install the browser extension and the phone app. Set autofill to fill only when you click.
3. Change the **email account password first**, then banking, then everything important. Use the generator (20+ random characters).
4. Turn on **two-factor authentication** (below) for email, the password manager, banking and cloud storage.
5. Write down the **recovery codes** and master passphrase and keep them somewhere physical and safe (not a photo on your phone).
6. Delete passwords saved in your browser once they are imported.

### Second factors, ranked

1. **Passkeys or hardware security keys** ([YubiKey](https://www.yubico.com)) are phishing-resistant. Buy two and keep one as a backup. (Guide 2 covers key setup in detail.)
2. **Authenticator apps** such as [Aegis](https://getaegis.app) (Android), [Ente Auth](https://ente.io/auth) or [2FAS](https://2fas.com) are fine for most accounts.
3. **SMS codes** are better than nothing but the weakest option (SIM-swap attacks). Avoid them for important accounts.

### A no-hassle way to rotate passwords

Forget 'change everything every 90 days'. Modern guidance, including NIST's, says forced calendar rotation pushes people toward weaker passwords. Rotate by _event_, plus one tiny routine:

- **Rotate immediately** if a service is breached ([Have I Been Pwned](https://haveibeenpwned.com) alert), you typed the password on a shared or unknown device, someone had access to your phone or laptop, or you see a login alert you didn't cause.
- **Rotation Sunday, twice a year (15 minutes):** change only your top ten (email, password-manager master, banking, main cloud, main social, domain and developer accounts) using the generator.
- **Never-rotate list:** low-value accounts with unique passwords and 2FA.
- Subscribe to Have I Been Pwned notifications for your real address and your aliases.

### Extra easy wins

- Give security questions **random answers** and store them in the manager.
- **Lock your phone and laptop** with a PIN of 6 digits or more (not a pattern) and turn on auto-lock.
- **Never read out a code** to anyone, including a 'support agent' who phoned you.
- Keep the OS, browser and password manager **updated**.
- **Check the address bar** before logging in; use bookmarks instead of clicking links in messages.
- Keep an **encrypted local backup** of your vault export and test the restore once.

### Worked chain

**Lost phone:** recovery codes + the spare YubiKey + the manager's emergency access get you back in.

**Synthesis.** Nobody is hacked because they forgot to be clever. A manager, unique passwords and a second factor on the five accounts that matter will stop most of what actually happens to people. Do those three things before anything fancier in this book.

## 14. Maximum Privacy Settings per Platform (That Still Let You Create)

**Philosophy.** The goal is _public work, private person_: your content can be seen by anyone, but nobody can easily find your phone number, location, contacts, old posts or family. Lock down _discovery, identity and data use_, not the posts you want people to see.

### 🧠 Universal rules for creators

1. **Separate the creator account from the personal one.** Different alias email (Part 7), different browser profile (Part 9).
2. **Control how people find you:** by email, by phone number, by uploaded contacts.
3. **Delete uploaded contacts** and switch off contact syncing.
4. **No live location and no geotags.** Post after you have left the place, not while you are there.
5. **Switch off AI-training and personalised-ads toggles.**
6. **Review connected apps** every quarter and revoke what you don't use.
7. **Download your data once** to see what the platform holds.
8. **Menus move.** Use each app's settings search bar and look for the keyword in bold below. The paths were written for October 2026.

For 2FA on every platform below, use an authenticator app or security key (Part 13).

### GitHub

1. [Settings → Emails](https://github.com/settings/emails): tick **Keep my email addresses private** and **Block command line pushes that expose my email**. Copy the `ID+USERNAME@users.noreply.github.com` address shown there.
2. Tell git to use it, plus a pseudonymous name if you like:

```
git config --global user.email 'ID+USERNAME@users.noreply.github.com'
git config --global user.name 'Your Handle'
```

Replace the placeholders. Commits you already made keep the old email; Part 4 shows how to rewrite history.

**Then:**

- [Settings → Public profile](https://github.com/settings/profile): remove real name (optional), location, company, website and social links. Use a non-photo avatar if you want. Leave **Include private contributions** off.
- [Settings → Password and authentication](https://github.com/settings/security): add a passkey or security key plus an authenticator app, and save the recovery codes.
- **Applications:** revoke OAuth apps you don't use. **SSH and GPG keys:** delete old ones.
- **Copilot:** look for **data use** or **training** options and switch off any that let your code or interactions be used for product improvement, if they are offered to you.
- **Organisations:** set your membership to **Private** from the organisation's People page.

Stars, follows and public repos are always public. Unstar or unfollow anything that reveals personal interests.

### LinkedIn

LinkedIn's terms require your real name, so the goal is to minimise what is exposed.

1. **Settings & Privacy → Visibility:** _Profile viewing options_ → Private mode (or Semi-private); _Edit your public profile_ → limit what logged-out visitors and search engines see; _Who can see your email address_ → Only you; _Who can see your connections_ → Only you; _Profile discovery using email address / phone number_ → No one; _Profile discovery and visibility off LinkedIn_ → off.
2. **Manage active status** → No one.
3. **Data privacy → Data for Generative AI Improvement** → off. Search the settings for **generative AI** if you can't find it.
4. **Advertising data:** switch off profile-data ad targeting, third-party data and contact syncing.
5. **Sign in & security:** two-step verification with an authenticator app; review _Where you're signed in_.
6. Remove phone number, birth year and address from the profile; use a city-level location.

### Instagram, Facebook and Threads (Meta's Accounts Center)

**Instagram**

1. A Professional or Creator account is fine for posting, but it can show a contact button. Hide email and phone in _Edit profile → Contact options_.
2. **Accounts Center → Personal details:** remove details you don't need. **Password and security:** two-factor with an authenticator app; log out unknown devices.
3. **Settings → search tags:** tags and mentions → _People you follow_ or _No one_; limit who can message you.
4. **Activity status** → off. **Story → hide from** → add anyone you don't want.
5. **Suggest your account to others / appear in search engines** → off (labels vary).
6. **Contacts:** switch off _Connect contacts_ and delete uploaded contacts.
7. **Ad preferences:** switch off ads based on partner activity and clear off-platform activity.
8. **Posting:** turn off location tagging in the camera app and the composer, delay posts, strip metadata (Part 12).
9. **AI:** Meta offers a formal objection form to AI training only in some regions. Where it exists, file it; otherwise assume public posts may be used.

**Facebook** (terms: use the name you are known by day to day)

1. **Settings & privacy → Privacy:** _Who can see your future posts_ → Friends (or Public if this is a creator Page); _Limit past posts_ → apply.
2. **How people find and contact you:** by email → Only me; by phone → Only me; **search engines outside Facebook** → off.
3. **Profile and tagging:** review tags before they appear; limit who sees posts you're tagged in.
4. **Your activity off Meta technologies** (Accounts Center): clear history and turn future activity off.
5. Turn on authenticator-app 2FA and check _Where you're logged in_.
6. Create from a **Page** and keep the personal profile quiet.

**Threads**

1. **Settings → Privacy:** private profile, or public with mentions and replies limited to people you follow.
2. **Fediverse sharing** → off. When it is on, your posts are copied to other servers and become harder to delete (Part 2).
3. **Hidden words** on; decide whether to show your linked Instagram account.

### X (Twitter)

1. **Settings and privacy → Privacy and safety → Audience, media and tagging:** keep posts public if you create, but set **Photo tagging** → off.
2. **Discoverability and contacts:** _Let others find you by email_ → off; _by phone_ → off; **Manage contacts** → delete anything uploaded.
3. **Direct messages:** allow requests only from people you choose.
4. **Data sharing and personalization:** switch off personalised ads, inferred identity, partner data sharing and location.
5. **Grok & third-party collaborators:** switch off use of your posts for training.
6. **Security and account access:** add a passkey, security key or authenticator, and revoke old apps.

### Telegram

Telegram's ordinary chats and groups are not end-to-end encrypted (Part 16), so lock down discovery and use it for broadcasts and communities, not private conversation.

1. **Settings → Privacy and Security → Phone Number:** _Who can see my phone number_ → Nobody; _Who can find me by my number_ → My Contacts.
2. **Last seen & online, Profile photos, Bio, Forwarded messages, Calls** → Nobody or My Contacts.
3. **Groups & channels** → My Contacts (stops random adds).
4. **Username:** set one so people can reach you without your number.
5. **Two-step verification:** set a password and a recovery email (use an alias).
6. **Active sessions:** end unknown sessions. **Passcode lock** on. **Auto-delete messages** with a default timer.
7. **Data settings:** turn off _Sync contacts_ and delete synced contacts.
8. **People Nearby:** never make yourself visible.
9. **Delete my account if away for:** pick a short timer if you want a dead-man's switch.

**Synthesis.** Public work and private life only coexist if discovery is turned off. Spend one evening per platform on the identity and discovery settings, and your posts can reach everyone while your phone number reaches no one.

## 15. Privacy Alternatives — The General Setup

**Philosophy.** Every mainstream app is tuned for _its_ growth: more data, more engagement, more lock-in. A privacy alternative costs you a small inconvenience (a new app, a changed habit, sometimes a few euros) in exchange for a large drop in data collection. Adopt alternatives where that trade is lopsided in your favour and skip the rest. One switch a week is a better pace than a weekend of burnout. This part covers the low-inconvenience tier; Guide 2 covers mid-to-high-inconvenience alternatives and Guide 3 covers self-hosting.

### 🧠 Judge any tool with five questions

1. Who pays for it: me, or advertisers and data buyers?
2. Is the code open to inspection, and has it been independently audited?
3. What identifiers does it demand at signup (email, phone, ID)?
4. Is my data end-to-end encrypted, or can the company read it?
5. Can I leave easily, with an export and a real delete button?

### Where to find alternatives and keep finding better ones

- **[Privacy Guides](https://www.privacyguides.org)** is a community-reviewed list of tools with published criteria. List it first.
- **[privacytools.io](https://www.privacytools.io)** is the original name. As far as I know the community behind it moved to Privacy Guides in 2021, so treat Privacy Guides as the maintained reference and privacytools.io as history.
- **[digdeeper.club](https://digdeeper.club)** hosts opinionated essays for people who want to go deeper. I could not verify its current state, so check it yourself and read it critically, as with every opinion source.
- **[AlternativeTo](https://alternativeto.net)** lets you filter by open source and by platform.
- **[F-Droid](https://f-droid.org)** is an app store for open-source Android apps.

### The default stack (one line per category)

- **Browser and blocker:** Firefox or Brave with uBlock Origin (Part 11).
- **Search:** Brave Search, DuckDuckGo or Startpage.
- **Email:** Proton Mail or Tuta, plus aliases (Part 7).
- **Passwords:** Bitwarden or KeePassXC (Part 13).
- **Messaging:** Signal (Part 16).
- **Maps:** [Organic Maps](https://organicmaps.app) or [OsmAnd](https://osmand.net), which work offline without an account.
- **Photos:** Ente Photos (Part 24).
- **Notes and docs:** Joplin, Obsidian (local files) or CryptPad (Part 24).
- **Office:** LibreOffice.
- **Phone keyboard:** an open-source keyboard from F-Droid such as HeliBoard, with cloud features off.

### 🛠 Hardening a mainstream account you can't avoid

Parts 18, 19 and 20 point back here. Apply this list once per service:

1. **Identity:** alias email and a pseudonym, never a social login (Part 7).
2. **Payment:** a gift card or cash-bought voucher where the service accepts one (Part 25).
3. **Access:** use the website in a browser container instead of the app (Part 9).
4. **Settings:** switch off personalised ads, activity sharing ('share my viewing/listening activity') and any data-use or AI-training toggles, and disconnect linked social and scrobbling accounts.
5. **Phone:** remove the advertising ID and deny permissions the app does not need. Android: _Settings → Security & privacy → Privacy → Ads → Delete advertising ID_. iPhone: _Settings → Privacy & Security → Tracking_ → off.
6. **Avatar:** never a real photo of you.

**Synthesis.** Privacy tools are chosen by trade-off, not by purity. Run the five questions, move the high-leverage categories first (browser, email, passwords, messaging), and let the rest follow when it is convenient.

## 16. Communication — Messages, Calls and Media

**Philosophy.** A message has two parts: the _content_ and the _metadata_ (who, when, how often, from where). Encryption protects content; good design also minimises metadata. Choose the messenger whose creators can see the least about you.

### Default pick: Signal

Open source, end-to-end encrypted by default for messages, calls and attachments, and run by a non-profit.

**🛠 Set it up in 10 minutes**

1. Install from [signal.org/download](https://signal.org/download/) or the official app store. Don't use third-party APK sites.
2. Register with a phone number (the default path, below), then set a **PIN** to protect re-registration.
3. **Settings → Privacy → Phone number:** _Who can see my number_ → Nobody; _Who can find me by number_ → Nobody.
4. **Settings → profile → Username:** set a username and share that instead of your number.
5. **Settings → Privacy → Disappearing messages:** set a default timer for new chats.
6. **Settings → Privacy → Advanced → Always relay calls:** on, so contacts don't see your IP address.
7. **Screen lock** on. For important contacts, verify the **safety number** in person or on a trusted channel.

**The honest downside.** Signal has always demanded a phone number to register, even though your contacts never need to see it. A number is the join key most tied to your real ID, and that has been Signal's biggest privacy objection for years. In September 2026 Signal began beta-testing an optional signup without a phone number on Android: a one-time fee of about $3 paid through Google Play, producing a 32-character Account ID and a 64-character Recovery Key, with no recovery if you lose either. Signal says it can't link the payment to the account, but Google still sees that a purchase happened. Phone-number signup remains the default. Existing accounts can't drop their number.

### Decision tree

- _Day-to-day private chat with contacts_ → Signal.
- _I want to pay for a no-phone messenger and keep it simple_ → [Threema](https://threema.ch).
- _I want no phone number and no account of any kind_ → Guide 2 (Session, SimpleX, Briar).
- _WhatsApp_ encrypts message content but gives Meta rich metadata. Use it only where your contacts leave you no choice, and harden it (Part 14).
- _Telegram:_ see below.

**Telegram, honestly.** Telegram asks for a phone number even though it has usernames, ordinary chats and all groups are **not** end-to-end encrypted (only one-to-one Secret Chats are), and its servers can read cloud chats. It is a fine broadcast and community tool. It is not a private messenger.

### Calls

- **Signal calls** (one-to-one and small groups) are end-to-end encrypted.
- **Larger or open meetings:** Jitsi (Part 23).
- **SIP calling:** [Linphone](https://www.linphone.org) with a privacy-respecting SIP provider.
- Ordinary phone calls and SMS are not encrypted. Treat them like postcards.

### Media and files

- Small files and photos: send through Signal.
- One-off large files: a password-protected share link from Proton Drive.
- Strip metadata before sending (Part 12).

### Worked chain

**Community channel:** Telegram public channel for announcements, Signal or Matrix (Part 17) for private discussion.

**Synthesis.** Pick a messenger that encrypts by default, asks for the least, and is built by people with no incentive to read your messages. Today that is Signal for most people, and anything else only for the contacts who force your hand.

## 17. Community — Forums, Chat and Groups

**Philosophy.** Communities leak the most because you talk a lot, in public, for years. Treat every post as permanent and tied to your handle, and pick platforms by _who controls the server_.

### Default pick and decision tree

- _Private chat for a small group_ → Signal group (Part 16).
- _Open community chat with history_ → Matrix with [Element](https://element.io). Private rooms and direct messages can be end-to-end encrypted; public rooms are not.
- _Forum-style discussion_ → [Lemmy](https://join-lemmy.org).
- _Microblogging_ → [Mastodon](https://joinmastodon.org). Everything is public and copied between servers.
- _Voice_ → [Mumble](https://www.mumble.info) (open source; you join a community's server).
- _Old-school chat_ → IRC on [Libera.Chat](https://libera.chat).
- _Discord, Facebook Groups, Reddit_ → fine for casual use under an alias (Part 15 checklist), but assume permanent logging, no end-to-end encryption in chat, and heavy data collection. Keep sensitive topics elsewhere.
- _Running a community server of your own_ → Guide 3.

### 🛠 Do this: Matrix and Element in 10 minutes

1. Choose a homeserver. matrix.org is crowded; others are listed at [servers.joinmatrix.org](https://servers.joinmatrix.org).
2. Install Element (or Element X on phones) and sign up with an alias email (Part 7).
3. **Skip** the prompts to find contacts by email or phone number.
4. Direct messages and private rooms are encrypted by default. In room settings check that **Encrypted** is on.
5. In _Settings → Encryption_ create and save the **recovery key**. Without it a new device can't read old encrypted messages.
6. Verify your other devices with the cross-signing prompt.

### Community hygiene rules

- Use a **pseudonym** that you use nowhere else (Part 7).
- **Never mix** personal details across posts: city, employer, school, schedule, family members, pet names, photos with landmarks.
- Before posting ask: _would I be fine seeing this next to my real name in ten years?_
- Assume **screenshots** exist. Deleting a post does not recall the screenshot.
- Don't click links or open files from strangers on your main profile (Part 9).

> ⛔ **Boundary — theory only.** This guide protects _you_. Techniques for unmasking, tracking or exposing other people are out of scope. In theory, communities are re-identified by writing style, timing and cross-posted details, which is why the hygiene rules above matter; applying that to someone else's identity is harassment or doxxing and is illegal in many places.

**Synthesis.** Communities are where people accidentally write their own dossier. Keep a clean handle, share less than feels natural, and put anything sensitive in an encrypted small group rather than a public room.

## 18. Video Entertainment — Movies, Series and Broadcasts

**Philosophy.** Streaming platforms don't just show you films; they watch you watching. Viewing history reveals your mood, health, politics and relationships. Reduce what they can attach to a real identity, and own what you can.

### Default pick and decision tree

- _Videos on YouTube without an account_ → [FreeTube](https://freetubeapp.io) (desktop and Android) or [NewPipe](https://newpipe.net) (Android, also on F-Droid). Both keep subscriptions locally and need no Google account.
- _Movies and shows I own_ → plain VLC with DRM-free files you purchased or ripped yourself where the law allows. (A home media server is a Guide 3 topic.)
- _Free and legal_ → your library's apps (Libby, Kanopy, Hoopla), and the public-domain collections at the [Internet Archive](https://archive.org).
- _I must use Netflix, Prime, Disney+ or similar_ → apply the Part 15 checklist. For streaming, step 4 means turning off viewing-history personalisation and _share my activity_.
- _Smart TV_ → see the TV section; the best answer is to keep the TV offline.

### 🛠 Do this: FreeTube in 5 minutes

1. Download from [freetubeapp.io](https://freetubeapp.io) and install.
2. In YouTube, export your subscriptions through [Google Takeout](https://takeout.google.com) (choose YouTube, subscriptions only).
3. In FreeTube choose _Settings → Data Settings → Import subscriptions_ and select the file.
4. Log out of Google in your browser. Watch from FreeTube from now on.

### 🛠 Smart TV privacy

Smart TVs run _automatic content recognition_ (ACR), which fingerprints whatever is on screen, even from a game console, and sells the result.

1. **Best:** never connect the TV to Wi-Fi. Use a laptop, a Raspberry Pi with Kodi, or an external streaming box instead.
2. **If it must be online:** search the TV's settings for **ACR**, **viewing information**, **live plus** or **interest-based advertising** and switch each off. Names vary by brand, so use the settings search.
3. Turn off voice assistants and microphones, and decline the broad privacy-policy options at setup.
4. Network-level blocking (Guide 2) also helps with a TV that has to stay online.

> ⛔ **Boundary — theory only.** This guide won't point you to piracy sites. Beyond the legal risk, they are among the biggest malware and tracking hotspots on the web, so they also defeat the privacy goal.

**Synthesis.** You can't make Netflix forget you, but you can make 'you' an alias on a gift card, and you can stop the TV from reporting what it sees. Own the media you love and let the rest be disposable.

## 19. Music Entertainment — Music, Podcasts and Audiobooks

**Philosophy.** Music taste is surprisingly revealing: mood, age, politics, sleep times, even health. Streaming services profile all of it. The privacy-friendly model is the old one: own the files, play them locally.

### Default pick and decision tree

- _Albums I love_ → buy DRM-free on [Bandcamp](https://bandcamp.com), which also pays artists better. Play with VLC or [Strawberry Music Player](https://www.strawberrymusicplayer.org) on desktop.
- _Podcasts_ → [AntennaPod](https://antennapod.org) (open source, Android, no account). Podcasts are just RSS feeds, so any RSS-capable player works, and you can import your OPML file from your old app.
- _Audiobooks_ → your library through Libby, public-domain recordings at [LibriVox](https://librivox.org), and purchases at [Libro.fm](https://libro.fm) (supports local bookshops).
- _Streaming my own library to my phone_ → Guide 3 (self-hosting).
- _I must use Spotify, Apple Music or YouTube Music_ → apply the Part 15 checklist; no major streamer is privacy-first. For music, step 4 also means turning off _share listening activity_, using private sessions, and disconnecting Facebook and any scrobbling service.

### 🛠 Build a local library in 15 minutes

1. Buy or rip your music into one folder: `Music/Artist/Album/01 Track.flac`.
2. Clean tags with [MusicBrainz Picard](https://picard.musicbrainz.org).
3. Point your player at the folder.
4. Back the folder up (Part 24).

**Synthesis.** Own what you love, borrow what you don't, and treat any streaming account as a disposable rental rather than a permanent profile.

## 20. Game Entertainment — From Mobile Puzzles to MMOs

**Philosophy.** Games collect an unusual mix: voice chat, friend lists, hardware IDs, real-time behaviour and, for some titles, kernel-level anti-cheat that runs with full OS privileges. Keep game accounts separate from your real identity, and keep the anti-cheat away from your main system.

### Default pick and decision tree

- _Single-player and indie_ → [GOG](https://www.gog.com) (DRM-free) and [itch.io](https://itch.io). The GOG Galaxy launcher is optional.
- _Steam library_ → fine, with hardening (below). On Linux, use Proton (Part 8).
- _Epic, Ubisoft, EA_ → alias accounts only (Part 15 checklist), and install the launchers in a compartment (Part 9).
- _Games with kernel-level anti-cheat_ → play on a dedicated dual-boot partition or a separate machine, with no personal files and no real-name logins, or skip them.
- _Mobile_ → fewer games, fewer permissions, no ad-funded free-to-play where you can avoid it. Look at the games section of F-Droid.

### 🛠 Harden Steam in 10 minutes

1. **Profile → Edit Profile → Privacy Settings:** _My profile_ → Private (or Friends only); _Game details_ → Private; _Friends list_ → Private; _Inventory_ → Private.
2. Use a username that is not your real name and an alias email (Part 7).
3. Turn on **Steam Guard** with the mobile authenticator.
4. Don't link Facebook, Twitch or other social accounts.
5. In _Settings → In Game_ turn the overlay off unless you need it.
6. Voice chat: use a private Signal call with friends instead of public lobbies.
7. Pay with Steam gift cards if you want your card out of the account (Part 25).

### 🛠 Mobile games

1. Remove the advertising ID and deny permissions the game doesn't need: contacts, location, microphone (Part 15, step 5).
2. Install games from F-Droid or the official stores only.
3. Delete games you no longer play. Each one is another account and another SDK.

> ⛔ **Boundary — theory only.** No steps here for cheats, hardware-ID spoofers or ban evasion. In theory, hardware bans exist because cheating ruins games for everyone and because spoofing tools are a common malware route, so avoiding them also protects your machine.

**Synthesis.** Treat games as entertainment, not identity. Private profiles, alias accounts, and quarantined anti-cheat cover almost all the risk without giving up play.

## 21. Code & Dev Architecture — Work and Code

**Philosophy.** Developer tools leak the most valuable data you own: source code, secrets, behaviour, identity (the email baked into every commit) and, now, your prompts to AI assistants. Your workflow should assume anything sent to a server is retained.

### Default pick and decision tree

- _Git hosting for privacy_ → [Codeberg](https://codeberg.org) (non-profit, runs [Forgejo](https://forgejo.org)).
- _Need GitHub's audience_ → keep GitHub, but use a pseudonymous identity and the noreply email (Part 14).
- _Editor_ → [VSCodium](https://vscodium.com) (VS Code without Microsoft telemetry), or VS Code with telemetry switched off.
- _AI coding assistants_ → assume prompts and code leave your machine. Read each tool's data-use settings and follow your employer's policy. Running a model locally is a Guide 2 topic.
- _Secrets_ → never in the repo. Use environment files and a secrets scanner.
- _Hosting your own Git server_ → Guide 3.

### 🛠 Switch off telemetry in 10 minutes

**VS Code:** open _Settings_ (JSON) and add:

```
"telemetry.telemetryLevel": "off"
```

**Shell:** add these lines to `~/.bashrc` or `~/.zshrc`:

```
export DO_NOT_TRACK=1
export HOMEBREW_NO_ANALYTICS=1
export DOTNET_CLI_TELEMETRY_OPTOUT=1
export NEXT_TELEMETRY_DISABLED=1
```

Then for Go (1.23 and newer) run `go telemetry off`. `DO_NOT_TRACK` is a community convention that some tools honour; the specific variables above are honoured by their own tools.

### 🛠 Keep your identity out of commits

**Use a dedicated SSH key per persona, and set its comment yourself.** The default comment is `user@hostname`, which leaks your username and machine name:

```
ssh-keygen -t ed25519 -C 'yourhandle'
```

**Sign commits with that key** so people can verify you without a real name:

```
git config --global gpg.format ssh
git config --global user.signingkey ~/.ssh/id_ed25519.pub
git config --global commit.gpgsign true
```

Then add the same public key to your host as a _signing key_. Use your noreply email and pseudonymous name in `user.email` and `user.name` (Part 14).

### 🛠 Stop leaking secrets

1. Put secrets in `.env` files and add `.env` to `.gitignore`.
2. Install [gitleaks](https://github.com/gitleaks/gitleaks) as a pre-commit hook so a commit containing a key is blocked before it exists.
3. If a secret ever reaches a public repo, **revoke and rotate it immediately**, then clean history (Part 4.1). Also check the provider's access logs for use of the old key.
4. Never paste secrets, customer data or private code into AI chats.
5. Check CI logs. Build output often prints environment variables.

**Synthesis.** In dev work privacy is mostly _hygiene at the edges_: the commit email, the SSH comment, the telemetry flags and the secrets file. Fix those four and you remove most of what a stranger can learn from your code.

## 22. Professional Communication — Mail and Message

**Philosophy.** Work communication is where real identity is usually mandatory, so the goal changes from anonymity to _control_: keep professional mail on its own address, keep sensitive content out of plain email, and keep work and personal life from leaking into each other.

### Default pick and decision tree

- _Professional mailbox_ → [Proton Mail](https://proton.me/mail), [Tuta](https://tuta.com) or [mailbox.org](https://mailbox.org), ideally on **your own domain** so you can move providers later.
- _Mail from a client that is not end-to-end encrypted_ → treat email like a postcard. Send secrets over Signal, or use a password-protected message or link.
- _Team chat_ → Slack and Teams are fine for work, but assume your company and the vendor can read everything. Running a private team server is a Guide 3 topic.
- _One-to-one professional contacts_ → Signal.

### 🛠 Set up a professional mailbox on your own domain

1. Buy a domain from a registrar (for example Porkbun, Namecheap or Gandi) and keep **WHOIS privacy** on so your home address isn't published.
2. Create a mailbox with your chosen provider and follow its _custom domain_ wizard.
3. At the registrar add the DNS records the wizard shows: **MX** (where mail goes), **SPF**, **DKIM** and **DMARC** (these prove your mail is genuine and stop spoofing).
4. Create aliases such as `hello@` and `billing@` so your main address stays out of public view.
5. Turn on two-factor authentication for the mailbox and the registrar (Part 13).

### 🛠 Encrypt when it matters

1. Proton to Proton is end-to-end encrypted automatically. For other recipients, Proton lets you send a **password-protected** message.
2. For other providers, install [Thunderbird](https://www.thunderbird.net) and create an OpenPGP key in _Account Settings → End-to-End Encryption → Add Key_.
3. Send your public key to the contacts you want to write to securely, and import theirs.
4. Turn off remote content in Thunderbird (_Settings → Privacy & Security → Mail Content → Allow remote content in messages_ off) to defeat tracking pixels.

### Hygiene rules

- Don't put a personal mobile number in your email signature. Use a work alias or a work line.
- Keep **personal accounts out of work devices** and **work accounts out of personal devices** (Part 9).
- Assume work chat is logged and discoverable.
- Keep professional-network settings tight (Part 14).

**Synthesis.** Professional privacy is about _separation and control_, not hiding. Own your domain, keep secrets off plain email, and keep work and personal life in separate rooms.

## 23. Virtual Meetings — Software and Habits

**Philosophy.** A meeting exposes your face, your voice, your room and your attendee list, all through a vendor's servers. Reduce what the platform learns, and reduce what _you_ show.

### Default pick and decision tree

- _Open, no-account meetings_ → [Jitsi Meet](https://meet.jit.si) (open source), with optional end-to-end encryption.
- _Small private calls_ → Signal group calls.
- _I'm forced onto Zoom, Meet or Teams_ → harden (below).
- _Self-hosted team calls or classes_ → Guide 3.

### 🛠 Run a private Jitsi meeting in 5 minutes

1. Go to meet.jit.si and type a **long random room name** (six random words works), then click _Start meeting_.
2. Open the toolbar's **Security options** and add a **password** or turn on the **lobby** so you approve each person.
3. Where it is offered, turn on **end-to-end encryption** in Security options. It works in supported browsers and may disable some features.
4. Share the link and password through Signal, not email.
5. Keep the **camera off** unless needed.

### 🛠 If you must use Zoom, Meet or Teams

1. Use the **web client** in a browser container (Part 9) instead of installing the app.
2. Join with an **alias display name** where the context allows.
3. Start with **camera and microphone off**.
4. Turn on **waiting room** and **passcode** for meetings you host.
5. Never record without telling everyone. Recording laws vary.
6. Share a **window**, not your whole screen, and turn on _Do Not Disturb_ first so notifications don't leak.
7. Check what's visible behind you. Virtual backgrounds can leak edges, so use a plain wall when it matters.
8. Cover the camera when you aren't in a call (Part 26).

**Synthesis.** Choose software that doesn't require an account for you or your guests, and treat the camera and screen-share button like a front door: open it deliberately, never by default.

## 24. Storage — Files and Media

**Philosophy.** Cloud storage is someone else's computer. With **end-to-end encryption** (E2EE) they hold the files but cannot read them. Encrypt _before_ upload, back up in more than one place, and test that you can restore.

### Default pick and decision tree

- _Use my existing cloud but keep files private_ → [Cryptomator](https://cryptomator.org) (open-source encrypted vaults on top of any cloud).
- _A cloud that encrypts for me_ → [Proton Drive](https://proton.me/drive) or [Filen](https://filen.io).
- _Photos_ → [Ente Photos](https://ente.io) (E2EE). Turn off Google Photos backup, and delete the old library only after you have verified the import.
- _Sync with no cloud at all_ → [Syncthing](https://syncthing.net) (device-to-device).
- _Notes and docs_ → [Joplin](https://joplinapp.org) with E2EE sync, Obsidian (local files) or [CryptPad](https://cryptpad.org) for collaborative documents.
- _Backups_ → restic or Borg, to an external disk (below).
- _A cloud I run myself_ → Guide 3.

### 🛠 Encrypt files before they touch the cloud (Cryptomator)

1. Download Cryptomator from cryptomator.org and install it.
2. Click **Add vault → Create new vault**, name it, and choose a location inside your cloud folder (Proton Drive, Nextcloud, Dropbox, and so on).
3. Choose a **long passphrase** (Part 13) and, when offered, **save the recovery key** somewhere offline.
4. **Unlock** the vault. A virtual drive appears; put files in it like a normal folder.
5. Lock the vault when done. Only the encrypted blobs sync.

### 🛠 A simple encrypted backup (restic)

Install restic from [restic.net](https://restic.net) or your package manager. Plug in an external drive, then:

```
restic init --repo /mnt/backup/restic-repo
restic -r /mnt/backup/restic-repo backup ~/Documents ~/Pictures
restic -r /mnt/backup/restic-repo snapshots
```

The `init` step asks for a passphrase. Store it in your password manager. **Test a restore** at least once:

```
restic -r /mnt/backup/restic-repo restore latest --target /tmp/restore-test
```

### The 3-2-1 rule

Keep **3** copies of important data, on **2** different kinds of storage, with **1** offsite (an encrypted drive at a relative's home, or an E2EE cloud). Encrypt all of it.

> 🧠 **Honest limit.** With E2EE storage a lost password means lost data. Keep the recovery key somewhere physical and safe.

**Synthesis.** Privacy for stored data comes from three habits: encrypt before upload, keep more than one copy, and prove you can restore. A backup you've never restored is only a hope.

## 25. Payments — Cash Is the Best Everyday Privacy Tool

**Philosophy.** Every card swipe is a permanent record of what you bought, where, when and who you are. Cash has no account, no profile and no database behind it. For normal people cash is still the single most private way to pay, and the habit costs nothing.

### 🧠 What cash does and doesn't hide

- **Hides:** the purchase from your bank, the card network, payment apps and the data brokers that buy transaction data.
- **Doesn't hide:** CCTV, loyalty cards, phone location, licence-plate cameras, ID checks the law requires, and the ATM withdrawal itself (the bank sees you took cash out, though not how you spent it).
- **Doesn't help online by itself:** websites need a payment method, so you need the bridges below.

### Default pick and decision tree

- _In person, everyday_ → pay cash. Say no to loyalty cards and phone-number lookups.
- _Online services that accept vouchers_ → buy a **voucher or gift card with cash** in a shop.
- _VPN_ → Mullvad **cash by post** (below).
- _Physical goods ordered online_ → choose **in-store pickup and pay at the counter**, or **cash on delivery** where it exists.
- _Bills and utilities_ → pay at the provider's counter or an agent location, in cash, where available.
- _Large purchases_ → expect ID checks and reporting rules (see the boundary box).

### 🛠 How to buy (almost) everything with cash

**Groceries, shops, restaurants, transport.** Pay cash. Decline emailed receipts, don't give a phone number, and skip loyalty cards (or take a store-issued anonymous one).

**Electronics and appliances.** Pay cash in store. Keep the paper receipt for warranty and skip optional registration forms.

**Utilities and rent.** Ask the provider for a cash option, an agent location or a counter payment. Keep the receipt.

**Mullvad VPN by cash post (verified October 2026).**

1. On mullvad.net generate an account number.
2. Write the account number on a slip of paper. Nothing else is required.
3. Put the slip and cash in an envelope. Check Mullvad's payment page for the accepted currencies and the postal address, and copy the address only from that page.
4. Send a small amount, because cash in the post can be lost. Expect one to two weeks before the time appears on your account. Cash payments are non-refundable.

**Online services via vouchers.** Gift cards and digital vouchers bought with cash work for app stores, games, streaming and some VPN resellers. They separate your _bank_ from the account, but the account is still an identity, so use an alias (Part 7).

**Prepaid cards.** Open-loop prepaid cards bought with cash can work for ordinary web shops. Fees are high, and some need ID registration, especially reloadable or higher-value cards.

**Digital goods you want to own.** Prefer DRM-free stores (Bandcamp, GOG, itch.io) that accept gift cards or one-off payments.

**Donations.** Cash in a box or an envelope.

> ⛔ **Boundary — theory only.** No steps here for **structuring** (splitting payments to dodge reporting thresholds), laundering money, evading tax or sanctions, using someone else's identity or payment instruments, or hiding illegal transactions. **Why:** cash-reporting rules exist in many countries, and splitting one transaction into pieces to avoid them is itself a crime ('structuring') even when the money is legitimate. Crypto, mixers and privacy coins are a Guide 2 topic, with their legal and exchange risks. Paying cash for _your own privacy_ is legal. Hiding crime is not, and this book will not help with it.

**Synthesis.** Cash is private at the moment of spending and weak at the moment of withdrawal, so use it where the privacy gain is biggest: everyday purchases, voucher bridges and your VPN.

## 26. Hardware Security — Protecting the Device Itself

**Philosophy.** Software privacy is worth nothing if someone can pick up the device and read the disk, or if malware already runs on it. Hardware security has three jobs: make the disk unreadable when the device is lost, catch malware, and make sure the camera can't be turned into a spy.

### 26.1 Full-disk encryption

**🧠 Theory.** Full-disk encryption (FDE) scrambles the whole drive with a key derived from your passphrase. A stolen laptop or a seized drive is then just noise. Without FDE anyone can boot a USB stick and read your files, whatever your login password says.

**Default pick:** **LUKS** on Linux (set in the installer) and **VeraCrypt** for encrypted containers or external drives. Privacy Guides notes that the operating-system tools (BitLocker, FileVault, LUKS) can use the TPM and Secure Boot chain, which makes them the better choice for the system drive in many setups, so use VeraCrypt mostly for containers and external drives.

**Decision tree**

- _Linux install_ → tick **Encrypt the new installation** during setup (Part 8). That is LUKS.
- _Windows, privacy first_ → move to Linux; otherwise BitLocker with the key kept **local** (below), and VeraCrypt for sensitive containers.
- _External drives and USB sticks_ → VeraCrypt volume or LUKS.
- _Mac_ → FileVault. When setting it up choose to create a recovery key rather than storing it in iCloud.
- _Phones_ → encryption is on by default once you set a screen lock. Use a strong PIN, not a pattern.

**🛠 Check that your Linux install is encrypted.** Run:

```
lsblk -f
```

Look for `crypto_LUKS` in the output. Then back up the LUKS header (replace the device name with the one `lsblk -f` showed for your encrypted partition):

```
sudo cryptsetup luksHeaderBackup /dev/nvme0n1p3 --header-backup-file luks-header.img
```

Store `luks-header.img` offline, for example on an encrypted USB stick, and never on the same disk. A damaged header makes the data unrecoverable.

**Passphrase and habits.** Use a six-word dice passphrase (Part 13). **Shut down instead of suspending** when travelling, because a suspended laptop keeps the key in memory. Add a UEFI/BIOS password. Set a short screen-lock timeout.

**Retiring a device.** Encrypting from day one is also what makes disposal safe: without the key the disk is just noise. If you never encrypted, encrypt first, then factory-reset. Guide 2 covers proper secure-erasure tools, because 'secure delete' tools are unreliable on SSDs and phones.

**🛠 Make an encrypted container with VeraCrypt.**

1. Download VeraCrypt from the official links on its [GitHub project](https://github.com/veracrypt/VeraCrypt) and verify the signature if you can.
2. Choose **Create Volume → Create an encrypted file container → Standard VeraCrypt volume**.
3. Pick a file name and location, leave the default encryption, and choose a size.
4. Enter a long passphrase, move the mouse to generate randomness, and click **Format**.
5. To use it: **Select File**, choose a slot, **Mount**, enter the passphrase, work in the new drive, then **Dismount**.
6. Back up the header with _Volumes → Backup Volume Header_.

**🧠 Why avoid Microsoft's default encryption setup (verified October 2026).** BitLocker itself is strong. The problem is **where the recovery key goes**. Windows 11 usually stores it in your Microsoft account when you sign in with one. Microsoft confirmed it gave the FBI recovery keys for three laptops under a warrant in a 2025 fraud case in Guam, and says it receives about twenty such requests a year; the keys it can hand over are the ones kept in its cloud. Nobody cracked the cipher. Microsoft held the key.

**🛠 If you must use BitLocker, keep the key yourself.**

1. Use Windows Pro and a local account where possible. Control Panel → _BitLocker Drive Encryption → Turn on BitLocker_ and choose **Save to a file** or **Print**, never the Microsoft account.
2. If a key is already in your account, create a new local one from an administrator command prompt, then remove the old protector:

```
manage-bde -protectors -get C:
manage-bde -protectors -add C: -RecoveryPassword
manage-bde -protectors -delete C: -id {OLD-KEY-ID}
```

After that, sign in at account.microsoft.com/devices/recoverykey and delete the stored cloud key. Keep the new 48-digit key offline. Windows 11 Home's automatic _Device Encryption_ may store keys in the cloud again, so Linux, Windows Pro with local keys, or VeraCrypt is cleaner.

### 26.2 Antivirus and EDR

**Philosophy.** Antivirus is a seatbelt, not a force field. You want one good, boring layer you don't have to babysit, plus habits that reduce what it needs to catch.

**Default pick:** Windows → built-in **Microsoft Defender**. Linux → **ClamAV** for on-demand scans of downloads and files you receive, plus updates, sandboxing and sensible habits.

**Why Defender (honestly).**

- It is built in and updated through Windows Update.
- It avoids a third-party vendor running deep inside your system with its own data appetite. Avast, for example, sold users' browsing data through its Jumpshot subsidiary and was fined by the US FTC in 2024.
- Independent labs such as AV-TEST and AV-Comparatives have rated it well for years. Check their latest reports before you rely on that.
- The catch: it is Microsoft, so cloud features send file samples and metadata. Configure them (below).

**🛠 Defender settings (10 minutes).**

1. _Windows Security → Virus & threat protection → Manage settings:_ **Real-time protection** On; **Tamper Protection** On; **Automatic sample submission** Off if you want fewer files leaving the PC (cloud-delivered protection is a trade-off between detection and privacy; decide on purpose).
2. _Ransomware protection → Controlled folder access_ On, and add your important folders.
3. _App & browser control → Reputation-based protection_ On.
4. Run a **Quick scan** weekly and a **Full scan** monthly. For a suspicious machine use _Microsoft Defender Offline scan_.

Optional, in an administrator PowerShell, to never send samples:

```
Set-MpPreference -SubmitSamplesConsent 2
Get-MpComputerStatus
```

**🛠 ClamAV on Linux.** Install, update the signatures and scan your downloads:

```
sudo apt install clamav clamav-daemon clamtk
sudo systemctl stop clamav-freshclam
sudo freshclam
sudo systemctl start clamav-freshclam
clamscan -r -i ~/Downloads
```

`freshclam` updates the signatures, `-r` scans folders recursively and `-i` prints only infected files. ClamTk is a simple GUI.

**ClamAV's honest limits.** It is signature-based, spots brand-new malware less well than Defender, and does not scan in real time by default. It is good for downloads, mail attachments, shared drives and mail servers. It is not a complete endpoint suite, so combine it with updates, least privilege, Flatpak sandboxing and not running unknown installers. Don't upload private files to VirusTotal; it shares them with security vendors.

**🧠 EDR in plain English.** _Endpoint Detection and Response_ watches what programs _do_ (starting processes, touching files, opening network connections), keeps a log and lets a defender respond remotely. Home users get a light version inside Defender.

> ⛔ **Boundary — theory only.** Techniques for evading antivirus or EDR are out of scope. In theory, defenders detect behaviour and attackers adapt, which is why layered habits (updates, least privilege, backups, careful downloads) beat any single product.

### 26.3 Camera covers

**🧠 Theory.** Webcam-spying malware exists. A camera light _usually_ shows when the camera is on, but it is not guaranteed on every device, and researchers have shown that on some older laptops the light could be bypassed. A cover is a free, mechanical, unhackable fix for the _visual_ risk.

**When it makes sense.** Cover the laptop or desktop camera whenever you are not on a call, especially at home and in bedrooms. Skip it if you never leave the device in a private space and the light is hard-wired.

**🛠 Do it properly.**

1. Use a **thin sliding cover** made for laptops. Thick stickers and bulky covers can **crack modern laptop screens** when the lid closes; Apple specifically warns about this for MacBooks, so choose a very thin cover or rely on the indicator.
2. For a desktop webcam, **unplug it** when not in use.
3. **Phones:** a cover isn't practical. Use the operating system instead: on Android 12 and newer use the **Camera access** and **Mic access** Quick Settings toggles; on iPhone watch the green and orange indicator dots and give apps _Allow While Using_ only.
4. **Microphones** can't be covered. Unplug headsets; hardware kill switches are a Guide 3 topic.
5. Remember the limit: covering stops _seeing_ you, not _hearing_ you or stealing files. The real defence is not getting infected.

### Worked chain

**Lost laptop:** full-disk encryption means a thief gets noise; change passwords from another device and revoke sessions.

**Synthesis.** Hardware security is about the worst day: the lost laptop, the infected download, the unwanted camera. Encrypt the disk and keep the key yourself, run one boring antivirus layer, and cover the lens when you aren't using it.

## 27. The 30-Day Plan and Closing Notes

### Week 1 — Find it and remove it (Parts 1 to 4; about 5 hours)

1. Pre-flight: a clean browser profile, a KeePassXC vault for the register and a removal mailbox (Part 1).
2. Run the audit and fill the exposure register (Part 1).
3. Fix Critical findings first: delete or lock accounts and content (Part 2), opt out of people-search sites and file DROP or GDPR requests (Part 3), request search-result removals (Part 4).
4. Set a re-check for two weeks from now. Do not close your main email yet (Part 5).

### Week 2 — Stop new leaks (about 3 hours)

1. Password manager (KeePassXC works), a unique password for your email, and 2FA (Part 13).
2. Install uBlock Origin or switch browsers (Part 11).
3. VPN with the kill switch on (Part 6).
4. Randomise the MAC address on phone and laptop and rename the devices (Part 10).
5. Install Signal and hide your number (Part 16).

### Week 3 — Lock the doors

1. Back up, then encrypt your devices (Parts 24 and 26).
2. Use an email alias for every _new_ signup (Part 7).
3. Tighten one platform per evening (Part 14).
4. Fit a camera cover (Part 26).
5. Keep deleting: five old accounts per sitting until the register shows no Critical or High items (Part 2).

### Week 4 — Retire the old doors and shield what remains

1. Retire old verification vectors now that aliases and 2FA exist (Part 5).
2. Set up compartments: browser containers and one VM (Part 9).
3. Switch two services from Parts 15 to 24 (storage, photos, notes).
4. Try a Linux live USB (Part 8).
5. Re-run the audit (Part 1, Steps 2 to 5), check each register row, and re-check your broker profiles (Part 3). Also run the VPN leak check, the tracking test and a metadata check.

### If you only have one evening

Audit Steps 1, 2 and 4 (about 45 minutes) → remove the worst Critical finding → password manager → email 2FA → VPN.

### The monthly 20-minute habit

Update everything, check Have I Been Pwned, re-run a quick name-and-handle search (Part 1, Step 2), re-check broker profiles every third month (Part 3), review app permissions and connected apps, delete one unused account, and test one backup restore.

### What this guide does not do

- It does not make you invisible to a determined state-level adversary (Guides 2 and 3 raise the cost further).
- It is not legal advice, and it does not replace it where records fall under retention rules or legal holds.
- It cannot remove copies other people made of your data.

### Author's checklist — delete before publishing

**Verified in October 2026:** full uBlock Origin works on Firefox and Brave while Chrome has only uBlock Origin Lite after the August 2026 Manifest V2 removal; Mullvad (account number, cash by post, WireGuard only, about €5 a month); Signal's numberless registration beta on Android (about $3, via Google Play); Microsoft's BitLocker key handover in the Guam case; Telegram's deletion page `my.telegram.org/delete`; the Internet Archive removal procedure; Linux Mint 22.3 current with Mint 23 planned for around Christmas 2026; Ubuntu 26.04 LTS; California's DROP (consumer requests since 1 January 2026, broker processing duties from 1 August 2026); Google's Results about you (US expansion to ID numbers, February 2026); PimEyes' opt-out process (photo plus redacted ID, suppresses PimEyes' own index only); Privacy Guides' data-broker-removal page (EasyOptOuts at about $20 a year); Sherlock and Maigret package names on PyPI.

**Not verified (flagged inline):** removal processes at Archive.today, Ghostarchive and Megalodon; exact menu labels inside each social platform; the current state of digdeeper.club and privacytools.io; the output and false-positive rate of Sherlock, Maigret and WhatsMyName (not run here); each broker's current opt-out page (they move, so lean on the Big Ass list); Google's personal-information removal form outside the US, UK and Australia; PimEyes' current free-tier limits.

**Opinion blocks to rewrite in your own voice:** the OSINT-first take (Part 1), the VPN take (Part 6), the Never-verifying take (Part 7), the Windows-to-Linux take (Part 8) and the Brave and uBlock Origin take (Part 11).

**v3 restructure.** The OSINT self-audit (new Part 1) and all removal work now come first. Old Parts 10, 11 and 12 became Parts 2, 4 and 5; data-broker removal was split out of old Part 10 into new Part 3; old Parts 1 to 9 became Parts 6 to 14, and old Parts 13 to 25 became Parts 15 to 27. Cross-references were renumbered to match.

**Handoff to Guides 2 and 3 (v2 deduplication).** These items were removed from Guide 1 and should be picked up in the later guides; Guide 1 keeps a one-line pointer where noted.

- _Guide 2:_ network-wide blocking with Pi-hole or NextDNS (was in Parts 11 and 18); GrapheneOS (Parts 9 and 15); Qubes OS and Whonix (Part 9); a VPS in place of a VPN (Part 6); Session, SimpleX and Briar (Part 16); local AI models with Ollama (Part 21); secure data erasure (Parts 24 and 26); crypto, mixers and privacy coins (Part 25); hardware-key setup (Part 13 keeps the ranking only).
- _Guide 3:_ self-hosting of media servers (Jellyfin, Navidrome, Audiobookshelf), photos (Immich), Git (Forgejo), team chat (Matrix, Zulip, Mattermost), meetings (BigBlueButton, Nextcloud Talk, Element Call) and search (SearXNG); OnionShare and other Tor tooling (Part 16); hardware kill switches for microphones (Part 26).
- _Dropped outright:_ the Wazuh, osquery and Sysmon pointer in Part 26.2, since Guide 2 argues against running endpoint agents.

**Duplicates collapsed into one home:** Signal's no-number beta and Telegram's lack of end-to-end encryption (Part 16); Mullvad cash-by-post steps (Part 25); secret leaks and key rotation (Part 4.1); the mainstream-account hardening checklist and advertising-ID removal (Part 15); the 'Sign in with…' warning (Part 7); the grace-period warning for deleted accounts (Part 2).

### Final synthesis

Privacy is not a wall; it is a thousand small decisions to give away a little less. People who succeed at it aren't the ones with the fanciest tools but the ones who keep doing the boring parts: the manager, the alias, the update, the deletion request. Start with one evening, come back next month, and let the habits compound.