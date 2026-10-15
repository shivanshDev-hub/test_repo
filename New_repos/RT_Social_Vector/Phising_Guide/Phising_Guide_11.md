# Social Engineering Incident Response Playbook

Find your situation, follow the steps in order, and never skip sections 1 and 2. Speed matters more than perfection.

## Index: Find Your Situation Fast

Start here (always)
- 1. First 10 minutes
- 2. Money lost or details shared

Email (3)
- 3.1 Clicked a link
- 3.2 Entered a password
- 3.3 Opened an attachment
- 3.4 Fake invoice or vendor change
- 3.5 Email hacked
- 3.6 Extortion email

SMS and chat apps (4)
- 4.1 SMS link
- 4.2 App file (APK)
- 4.3 Friend asks for money
- 4.4 Chat account hijacked
- 4.5 Wrong-number stranger

Phone and video calls (5)
- 5.1 Fake bank call
- 5.2 Fake police, tax, customs or courier
- 5.3 Remote access or tech support
- 5.4 SIM or number threat
- 5.5 Cloned-voice emergency
- 5.6 Fake loan, refund, prize or utility

Online scams (6)
- 6.1 Fake shop or customer care
- 6.2 Marketplace
- 6.3 QR or payment request
- 6.4 Investment or crypto
- 6.5 Fake job
- 6.6 Romance
- 6.7 Lottery or prize
- 6.8 Rental or booking
- 6.9 Scholarship or admission
- 6.10 Fake charity
- 6.11 Sextortion or blackmail
- 6.12 ID documents shared
- 6.13 Loan app threats
- 6.14 Impersonated profile
- 6.15 Fake software or extension

Physical and in-person (7)
- 7.1 Tailgating
- 7.2 Fake delivery
- 7.3 Fake IT, bank or utility staff
- 7.4 Found USB
- 7.5 ATM or card
- 7.6 Lost or stolen phone
- 7.7 Charging or public Wi-Fi
- 7.8 Fake staff caller

Account takeover (8)
- 8.1 MFA spam or code request
- 8.2 SIM swap
- 8.3 Account hijacked
- 8.4 Password breached

Tools
- 9. Evidence
- 10. Who to report to
- 11. After-incident checklist
- 12. Team first hour
- 13. Quick reference

## 1. First 10 Minutes (Any Attack)

1. **Stop.** Hang up, close the chat, stop replying. Send nothing more: no money, codes, documents or photos.
2. **Money moved, or card or bank details shared?** Go to section 2 immediately.
3. **Screenshot before you delete:** number or handle, messages, links, amounts, times.
4. **Installed an app or allowed remote access?** Turn on airplane mode now.
5. **From a different, trusted device,** change the password of the affected account and of your email account, then sign out of all sessions.
6. **Tell someone you trust right now.** Ignore any instruction to keep it secret.
7. **Report it** using section 10.
8. **Ignore anyone who offers to recover your money for a fee.** That is a second scam.

## 2. Money Lost or Details Shared (First Hour)

1. Call your bank, card issuer or payment app's fraud line. Use the number on the back of your card or inside the official app, never a number the attacker gave you.
2. Say: 'I am a victim of fraud.' Give the time, amount, receiving account or ID, and transaction reference.
3. Ask them to block the card, account or wallet; freeze or recall the transaction; flag the receiving account; open a dispute; and give you a complaint reference number.
4. Switch off online, international and contactless payments. Reset PINs and online banking passwords. Remove saved cards from apps and browsers.
5. Call the national financial-fraud helpline and file a complaint on the national cybercrime portal. Do it within the same hour if you can.
6. **Crypto sent:** contact the exchange or wallet provider with the wallet address and transaction hash and ask them to freeze it.
7. **Gift cards bought:** contact the card issuer with the card numbers and receipts.
8. **Wire or international transfer:** ask your bank to send a recall request.
9. File a police report with printed evidence. Keep the complaint number.
10. Log every call: date, name, reference number. Follow up daily for a week.
11. Check all your other accounts for unknown transactions and turn on transaction alerts.

## 3. Email Attacks

### 3.1 You clicked a phishing link but entered nothing

1. Close the page. Enter nothing.
2. Run a full security scan on the device.
3. Report the email as phishing in your email app (at work, forward it to IT instead), then delete it.
4. Watch your account alerts for the next week.

### 3.2 You entered a password or details on a fake page

1. From a clean device, change that password and the password of every account that shares it.
2. Sign out of all sessions and remove unknown devices.
3. In the email account, delete unknown forwarding rules, filters, delegated access, connected apps, and recovery emails or numbers.
4. Turn on MFA using an authenticator app, not SMS.
5. Check Sent and Trash for messages you did not send. Warn your contacts if spam went out.
6. Card or bank details entered: section 2. ID documents entered: section 6.12.

### 3.3 You opened a malicious attachment or ran a file

1. Disconnect from Wi-Fi and cables. Do not log in to anything on this device.
2. Work or college device: call IT or security immediately and stop using it.
3. Personal device: boot into safe mode and run a full scan with a trusted anti-malware tool.
4. From another device, change passwords for email, banking, social media and cloud storage.
5. **Ransomware note on screen:** do not pay, do not contact the attacker, power off, and report. Restore from a clean backup.
6. If the scan cannot clean it, or you are unsure: back up personal documents only, wipe the device and reinstall the operating system.

### 3.4 Fake invoice, executive request or changed vendor bank details

1. Stop the payment. If it was already sent, call the bank at once and request a recall.
2. Alert your finance head and your IT or security team.
3. Verify the request by calling a number you already had on file, not one in the email.
4. Save the original email with full headers. Do not just forward it.
5. Check the compromised mailbox for forwarding rules and reset its credentials.
6. Review every other pending payment and recent bank-detail change from the same sender.
7. Warn the vendor or executive through a separate channel.
8. Report to the police and the national cybercrime portal.

### 3.5 Your email account was hacked

1. Recover it through the provider's recovery page on a clean device.
2. Change the password and turn on MFA.
3. Remove unknown forwarding rules, filters, connected apps, devices and recovery details.
4. List the accounts that use this email and reset their passwords, banking first.
5. Tell contacts not to open anything sent from you.
6. Search the inbox for password-reset emails, receipts and bank notices to see what the attacker touched.

### 3.6 Extortion email ('I hacked your webcam', quotes an old password)

1. Do not reply and do not pay.
2. Change the password they quoted, and anywhere you reused it.
3. Mark it as spam and turn on MFA.
4. If it quotes a current password or includes real private photos, follow 3.2 and 6.11.

## 4. SMS and Messaging-App Attacks

### 4.1 Link in an SMS (delivery fee, fine, bank alert, prize, 'account blocked')

1. Do not tap it. Do not reply, not even 'STOP'.
2. Tapped but entered nothing: close the page and run a scan.
3. Entered card or bank details: section 2 now.
4. Entered a password: follow 3.2.
5. Forward the SMS to your carrier's spam-reporting channel and block the sender.

### 4.2 App file (APK) or software received in a chat

1. Do not install it. Delete the file.
2. Already installed: turn on airplane mode, uninstall it (use safe mode if it blocks removal), and remove it from Accessibility and Device Admin settings.
3. From another device, change passwords for banking, email and messaging apps.
4. Call your bank, because these apps can read SMS codes and your screen.
5. Factory reset the phone if the app had accessibility, SMS or screen-sharing permissions.
6. Report the sender and file to the messaging platform.

### 4.3 A friend or relative messages asking for urgent money

1. Send nothing yet.
2. Call them on their saved number, or ask for your family code word.
3. Ask a question only the real person can answer.
4. Already sent money: section 2.
5. Tell the real person their account may be hijacked, using another channel.
6. Report the number or profile to the platform.

### 4.4 Your own messaging account was hijacked

1. Log in again with your number through the app's recovery flow to push the attacker out.
2. Turn on two-step verification and set a PIN.
3. Open linked devices and remove unknown ones.
4. Warn contacts on another channel to ignore money requests from your account.
5. Check whether your SIM was swapped (section 8.2).

### 4.5 Wrong-number or friendly stranger who moves to investing, trading or a 'job'

1. Stop replying. Block and report the account.
2. Money sent: section 2. Never pay 'taxes' or 'fees' to withdraw your money.
3. Screenshot the chat, profile, wallet addresses and app or site links.
4. Leave any investment-tips or job group you were added to, and report it.
5. Warn anyone you invited.

## 5. Phone and Video Call Attacks

### 5.1 Fake bank, card or account call

1. Hang up. Wait a minute, or use a different phone to call out, because some scammers keep the line open.
2. Call the number on your card or inside the official app.
3. Shared a code, PIN, CVV or password: ask the bank to block cards, freeze the account, and reverse or hold pending transactions.
4. Change PINs and online banking passwords.
5. Moved money to a 'safe account': section 2 now.
6. Installed an app at their request: follow 4.2 and 5.3.
7. Note the number and time, and report it.

### 5.2 Fake police, tax, customs, courier or regulator call (including video-call 'investigation')

1. End the call. Do not stay on camera or 'under supervision'. Real authorities do not work this way.
2. Tell a family member or friend now.
3. Transfer nothing. Share no ID documents or banking details.
4. Already paid: section 2 immediately. The first minutes matter most.
5. Already shared ID documents or photos: section 6.12.
6. Save call logs, numbers, screenshots and any recordings.
7. Verify the claim by visiting the real office or calling a number from its official website.
8. Report to the police and the national cybercrime portal.

### 5.3 Tech-support or remote-access scam (call or pop-up)

1. Disconnect from the internet and close the remote session.
2. Uninstall the remote-access tool (AnyDesk, TeamViewer, QuickSupport and similar).
3. Run a full scan. If they had access for long, or you typed passwords while connected, wipe and reinstall.
4. From a clean device, change email, banking and other passwords.
5. Logged in to banking while connected: call the bank at once (section 2).
6. Accepted a 'refund': it is an overpayment trick. Ask the bank to reverse it and never send the difference back.
7. Pop-up that will not close: restart the browser. Never call the number on screen.

### 5.4 Threat to disconnect your number or SIM, or sudden loss of signal

1. Hang up. Contact your carrier from another phone or its official app.
2. No signal for no reason: treat it as a SIM swap, section 8.2.
3. Never press keys, forward calls or send SMS commands on a caller's instruction.

### 5.5 Cloned-voice or emergency call (relative in trouble, kidnapping, accident)

1. Hang up and call the relative's saved number directly. Try another family member too.
2. Ask for your family code word.
3. Send no money to a stranger's account, gift cards or crypto.
4. Kidnapping claim: call the police from another phone while you verify.
5. Already paid: section 2.
6. Save the number and any recording.

### 5.6 Fake loan, insurance, refund, prize or utility-disconnection call

1. Hang up. Press no keys and confirm no personal details.
2. Verify with the company using the number on a past bill or its official website.
3. Paid a 'processing fee' or 'deposit': section 2.
4. Shared documents or card details: section 6.12 and section 2.
5. Block the number and report it to your carrier.

## 6. Online Scams

### 6.1 Fake shopping site, or a fake customer-care number found in search

1. Stop. Do not call the number and do not make further payments.
2. Paid by card: call your bank, dispute it as fraud, and replace the card.
3. Paid by bank transfer or payment app: section 2.
4. Gave login details: change the passwords. Let them see your screen: follow 5.3.
5. Screenshot the site, URL, order page and contact numbers.
6. Report the site to the hosting provider, browser safe-browsing and the national cybercrime portal.
7. Find real support only through the company's official app or the number on your card or bill.

### 6.2 Marketplace buyer or seller scam

1. Stop. A buyer never needs your PIN, a QR scan or a payment approval to pay you.
2. Decline every 'collect' or payment request. Scan no QR code and enter no PIN to receive money.
3. Already lost money: section 2.
4. Fake payment screenshot, or the 'buyer' overpaid: do not hand over the item or refund anything until your own bank shows cleared funds.
5. Asked to pay a courier or advance-delivery fee: refuse.
6. Report the listing and profile to the marketplace and save the chats.

### 6.3 QR code or payment-request scam

1. Scan or approve nothing further. Decline pending requests.
2. Money left your account: section 2.
3. Fake QR sticker pasted over a real one (shop, parking, charger): photograph it, tell the owner and report it.
4. Open your payment app and remove unknown linked devices and recurring-payment authorisations.

### 6.4 Investment, crypto or trading-app scam

1. Stop all deposits, including 'tax', 'fee' or 'unlock' payments to withdraw.
2. Money already sent: section 2. For crypto, contact the exchange with wallet addresses and transaction hashes.
3. Screenshot the app, site, wallet addresses, transaction IDs and chats.
4. Uninstall the app and change passwords. If it was installed from a link or file, follow 4.2.
5. Report to the platform where you met them, the financial regulator, the police and the national portal.
6. Refuse any 'recovery' offer.

### 6.5 Fake job, task-for-pay, work-from-home or overseas job offer

1. Stop paying. Real employers do not charge for hiring, training, equipment or travel.
2. Share no ID documents, banking login or card details for a 'salary account'. Already shared: 6.12.
3. Money sent: section 2.
4. Task scams pay small amounts early, then demand deposits to continue. Stop at the first deposit demand.
5. Verify the employer on its official website using a contact you find yourself.
6. **If you received money and were told to pass it on:** stop, do not forward it, and tell your bank and the police. You may be used as a money mule.
7. Report to the job platform, the police and the national portal. Warn others in the group.

### 6.6 Romance or friendship scam

1. Stop sending money, gifts or images. Block the account.
2. Money sent: section 2.
3. Reverse-image-search their photos and save the evidence.
4. Intimate images shared: 6.11.
5. Report the profile to the platform.
6. Talk to someone you trust.

### 6.7 Lottery, prize, inheritance or fee-to-release scam

1. Stop. A genuine prize never needs you to pay first.
2. Fee already paid: section 2. Do not pay again to 'release' it.
3. Documents shared: 6.12.
4. Block the sender and report it.

### 6.8 Rental, housing-deposit, ticket or travel-booking scam

1. Stop further payments. Pay no deposit for a place you have not verified.
2. Money sent: section 2. Paid by card: dispute it.
3. Confirm the owner or booking directly with the property manager, platform or airline using contacts you find yourself.
4. Report the listing to the platform and the police. Save the ad, chats and receipts.

### 6.9 Fake scholarship, admission, certificate or loan-approval (advance fee)

1. Stop paying 'processing', 'registration' or 'insurance' fees.
2. Verify through the institution's or lender's official website and phone number.
3. Money lost: section 2. Documents shared: 6.12.
4. Report to the institution, the police and the national portal.

### 6.10 Fake charity or emergency fundraiser

1. Stop. Check the charity in the official registry and on its own website.
2. Money sent: section 2. Paid by card: dispute it.
3. Report the fundraiser to the platform hosting it.

### 6.11 Sextortion, intimate-image blackmail or deepfake threat

1. Do not pay. Payment rarely ends it and demands usually grow. Stop communicating.
2. Do not delete anything: screenshot the threats, profile, usernames, wallet addresses and phone numbers.
3. Lock down your accounts: make friends and followers lists private, change passwords, turn on MFA.
4. Report the profile and the content to the platform. Use its intimate-image removal tool and a trusted takedown service.
5. Report to the police and the national cybercrime portal. You are the victim.
6. Tell a trusted person. If a minor is involved, tell a trusted adult and report immediately.
7. Deepfake of you: save every URL, report it on each platform, and report it to the police.
8. Do not meet the blackmailer and do not send further images.

### 6.12 You shared ID documents, selfies or personal data

1. Write down exactly what you shared and with whom.
2. Tell your bank and ask them to flag your accounts for extra verification.
3. Place a fraud alert or credit freeze with the credit bureaus.
4. Ask your mobile carrier to lock your number with a PIN and alert you to new SIMs or port-outs.
5. Check whether new accounts, SIMs, loans or cards were opened in your name. Pull your credit report.
6. Report the misuse to the authority that issued the document and ask what to do next, such as flagging or reissuing it.
7. Expect follow-up scams that use this data. Keep a record of each one.

### 6.13 Fake or predatory loan app, or threats to leak your contacts

1. Uninstall the app. Revoke its contacts, photos and storage permissions.
2. Do not pay extra 'fees' under pressure.
3. Tell your contacts a threatening message may arrive and to ignore it.
4. Screenshot the threats, numbers, app name and store link.
5. Report the app to the app store, the financial regulator, the police and the national portal.
6. Money lost: section 2. Installed from a file: follow 4.2.

### 6.14 Your profile is impersonated or hijacked

1. **Fake profile of you:** report it for impersonation in the app, warn your contacts with a screenshot of the fake URL, and make your friends list private.
2. **It is asking people for money:** tell the platform and the police.
3. **Your real account is taken over:** follow 8.3.

### 6.15 Fake software, cracked app, fake update pop-up or malicious browser extension

1. Close the browser. Do not run the download.
2. Already ran it: follow 3.3.
3. Remove unknown browser extensions, reset browser settings and clear cookies.
4. From a clean device, change passwords and clear passwords saved in the browser.
5. Update the operating system and browser only through their built-in update settings.

## 7. Physical and In-Person Attacks

### 7.1 Unknown person in a secured area (tailgating)

1. Do not hold the door. Politely ask who they are and who they are visiting.
2. Escort them to reception or security. Never leave them alone.
3. Report to security with a description, the time and the location.
4. If they touched a desk or device: lock it, tell IT and change credentials.
5. Ask for camera footage and visitor logs to be saved.

### 7.2 Fake delivery person or courier at the door

1. Give no code, share no screen and pay for no parcel you did not order.
2. Check the order in the official shopping or courier app before opening the door.
3. Refuse. If it is safe, photograph the person and vehicle.
4. Gave a code or paid: section 2.
5. Report to the courier company and the police.

### 7.3 Fake IT, auditor, bank, utility or inspection staff (home or office)

1. Do not let them in or give access. Ask for ID and verify through the company's official number, not one they give you.
2. Anything they touched or connected: unplug it and treat it as compromised.
3. Report to security or the police. Change passwords for anything they saw.
4. Check desks, ports and network equipment for anything new.

### 7.4 Found or handed USB drive or device

1. Do not plug it in. Hand it to IT or security.
2. Already plugged in: disconnect from the network, power off and report to IT. On a personal device, follow 3.3.
3. Note who gave it to you or where you found it.

### 7.5 ATM or card skimmer, 'helpful' stranger, swapped card, shoulder surfing

1. Loose card slot or odd attachment: do not use the machine. Report to the bank and the ATM owner.
2. A stranger offered help or swapped your card: walk away and call your bank at once to block the card.
3. Someone saw your PIN: change it now and check your statements.
4. Report to the police with the location and time so camera footage can be requested.

### 7.6 Lost or stolen phone, wallet or ID

1. From another device, locate, lock and erase the phone using its find-my-device service.
2. Call your carrier to block the SIM and get a replacement.
3. Call your bank to block cards and disable payment apps.
4. Change passwords for email first, then banking, then social media. Sign out of all sessions.
5. Report the loss to the police and keep a copy of the report.
6. Lost ID: report it to the issuing authority and follow 6.12.
7. Do not click any 'we found your phone' link.

### 7.7 Fake charging station, public Wi-Fi or fake hotspot

1. Unplug the device. Do not use public USB ports or unknown cables.
2. Public Wi-Fi: disconnect and forget the network. If you logged in to anything, change those passwords from a trusted network.
3. Scan the device. Check for new apps or profiles and remove them.
4. Tell the venue about the charger or hotspot.

### 7.8 Caller to your reception, helpdesk or support line pretending to be staff

1. Do not reset passwords, grant access, share information or change details based on a call alone.
2. Verify with a call-back to the internal directory number or through the ticket system.
3. Already acted: tell IT now, reverse the change and force a password and MFA reset on that account.
4. Log the number, time and what they asked. Warn colleagues.

## 8. Account Takeover

### 8.1 MFA prompt spam, or someone asks for your code

1. Deny every prompt. Never read out or type a code for anyone who called or messaged you.
2. Change the password immediately and sign out of all sessions.
3. Open the list of registered MFA devices and remove any you do not recognise.
4. Report to IT or the provider.
5. Move to an authenticator app, passkey or hardware key.
6. If you already approved a prompt, treat the account as hacked: follow 8.3.

### 8.2 SIM swap (sudden loss of signal, unexpected SIM or number-change alerts)

1. Use another phone or visit a carrier store immediately. Get the SIM restored and the number flagged.
2. Call your bank to freeze accounts and block payments.
3. Change passwords in this order: email, banking and payment apps, social media. Move MFA off SMS.
4. Check recent logins and transactions, and start recovery on every account that was taken.
5. Ask the carrier for a SIM-lock or port-out PIN.
6. Report to the police and the national cybercrime portal.

### 8.3 Any account hijacked (email, social, bank, cloud, gaming, marketplace)

1. Use the provider's recovery page on a clean device.
2. Change the password, turn on MFA and sign out of all sessions.
3. Remove unknown devices, recovery contacts, connected apps, forwarding rules and API keys.
4. Review recent activity, messages sent, payments and purchases.
5. Warn your contacts through another channel.
6. Recovery failing: contact the provider's support with proof of identity.
7. Reset every other account that shared the same password.

### 8.4 Your password appeared in a breach, or you reused it on a scam site

1. Change it everywhere you used it, starting with email and banking.
2. Switch to a password manager with a unique password per account.
3. Turn on MFA everywhere it is offered.
4. Check your email address on a breach-checking service.

## 9. Evidence to Save

- Screenshots showing the number or handle, URL and timestamps
- Call logs and any recordings
- The original email with full headers (use 'show original', do not just forward)
- Chat exports
- Bank and card statements, transaction IDs and receipt numbers
- Wallet addresses and transaction hashes
- App and site names, and app-store links
- A written timeline: when it began, what was said, what you shared, what you did
- Do not wipe, factory reset or delete files until evidence is saved, unless you must isolate a live compromise
- Do not reopen malicious links to check them

## 10. Who to Report To (In This Order)

1. **Your bank, card issuer or payment app.** Immediately, if any money is involved.
2. **National financial-fraud helpline and national cybercrime portal.**
3. **Police.** Keep the complaint number.
4. **Your mobile carrier.** For SIM swap, scam numbers and spam.
5. **The platform.** Email provider, social network, messaging app, marketplace, exchange or app store.
6. **Your employer or college IT and security team.**
7. **The national computer emergency response team.** For organisations and serious incidents.
8. **Your regulator or data protection authority.** If other people's data was exposed.
9. **Hosting provider, domain registrar and browser safe-browsing.** To take phishing sites down.
10. **Credit bureaus.** If ID documents or personal data were shared.

**Do this today:** save the official helpline numbers and portal addresses in your phone. Take them only from official government and bank websites, never from search ads or from a caller.

## 11. After the Incident

**Within 24 hours**

- [ ] Changed every affected password from a clean device
- [ ] Turned on MFA for email, banking and social media
- [ ] Signed out of all sessions
- [ ] Replaced any exposed card
- [ ] Filed all reports and saved complaint numbers
- [ ] Turned on transaction alerts for bank and payment accounts

**Within 7 days**

- [ ] Followed up with the bank daily
- [ ] Ran a full device scan
- [ ] Removed unknown apps, browser extensions and app permissions
- [ ] Checked email for unknown forwarding rules and filters
- [ ] Checked your credit report
- [ ] Told contacts if your account was used

**For the next 90 days**

- [ ] Review statements every week
- [ ] Ignore recovery offers
- [ ] Expect follow-up scams that use what you shared, and verify every contact

## 12. For Teams: First-Hour Checklist

1. Take the report and record who, what, when and how.
2. Name an incident lead and start a timeline.
3. **Contain:** disable affected accounts, revoke sessions and tokens, isolate affected devices, block the sender, domain or number, and purge the email from every mailbox.
4. Reset credentials and MFA for everyone who interacted with the attacker.
5. **Finance:** freeze or recall payments and call the bank.
6. Preserve logs, email headers and device images. Do not wipe anything.
7. Check for wider impact: other mailboxes, forwarding rules, new connected apps, and other staff who got the same message.
8. Notify leadership. Notify regulators and affected people if personal data was exposed.
9. Report to the police and the national computer emergency response team.
10. Within 5 working days, hold a no-blame review: signal missed, control that failed, owner and date for each fix.

**Standing rules:** no-blame reporting, dual approval for payments, call-back verification for any change of bank details, and nobody ever reads out a code.

## 13. Quick Reference

| If this happened | Do this first |
| --- | --- |
| Money sent or card details shared | Call the bank's fraud line, then the national helpline and portal (section 2) |
| Code, PIN or password shared | Change credentials from a clean device, sign out everywhere, block cards |
| Caller claims to be an official and demands money or secrecy | Hang up, tell someone, transfer nothing, verify independently (5.2) |
| Clicked a link and entered details | Change the password, turn on MFA, scan the device (3.2, 4.1) |
| Installed an app or let someone remote in | Airplane mode, uninstall, change passwords elsewhere, call the bank (4.2, 5.3) |
| Lost mobile signal for no reason | Call the carrier from another phone, then freeze bank accounts (8.2) |
| Friend or relative asks for urgent money | Call their saved number or use the code word before sending anything (4.3, 5.5) |
| Blackmail or sextortion threat | Do not pay, save evidence, lock down accounts, report (6.11) |
| Told to deposit more to withdraw earnings | Stop paying and report (6.4, 6.5) |
| Shared ID documents | Fraud alert, credit freeze, carrier lock, check for new accounts (6.12) |
| Phone or wallet lost or stolen | Lock the phone, block the SIM and cards, change passwords (7.6) |
| Suspicious payment request at work | Stop payment, verify by call-back, alert finance and IT (3.4) |
| Unsure which one applies | Treat it as the worst case, act now, verify later |

**Memory aid:** Stop. Secure. Save. Report. Recover.
