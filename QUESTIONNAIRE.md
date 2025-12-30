# Mail Server Project Questionnaire

Please answer the following questions to help refine the project requirements and implementation approach.

## 1. Deployment Environment

**Q1.1**: What is your primary deployment environment?
- [!] Bare metal Linux server
- [2] Virtual Private Server (VPS)
- [3] Docker containers
- [ ] Kubernetes cluster
- [ ] Cloud platform (AWS/GCP/Azure)
- [ ] Other: ___________

**Q1.2**: What Linux distribution will you be using?
- [X] Ubuntu/Debian
- [X] RedHat/CentOS/Rocky Linux
- [X] Alpine Linux
- [ ] Other: ___________

**Q1.3**: Expected scale at launch?
- Number of domains: __Unlimited_________
- Number of users: __Unlimited_________
- Expected daily email volume: __100000 per day_________

**Q1.4**: Expected scale in 1 year?
- Number of domains: __100_________
- Number of users: __200_________
- Expected daily email volume: __10000_________

## 2. Storage and Performance

**Q2.1**: What is your preferred message storage approach?
- [X] Store full message in Postgres (BLOB/TEXT fields)
- [ ] Store message metadata in Postgres, files on disk
- [ ] Hybrid approach (small messages in DB, large on disk)
- [ ] Preference? ___________

**Q2.2**: Expected average message size?
- [ ] < 100 KB (mostly text)
- [ ] 100 KB - 1 MB (some attachments)
- [ ] 1-10 MB (frequent attachments)
- [X] > 10 MB (large attachments common)

**Q2.3**: How long should messages be retained?
- [ ] Forever (never auto-delete)
- [X] User-configurable retention
- [ ] Specific period: ___________ days/months/years

**Q2.4**: Postgres configuration preference?
- [X] Single Postgres instance
- [ ] Postgres replication (master-slave)
- [ ] Postgres Galera cluster
- [ ] Amazon RDS or managed Postgres
- [ ] Other: ___________

## 3. Security and Compliance

**Q3.1**: What level of password security do you require?
- [X] Bcrypt (standard)
- [ ] Argon2 (more secure, slower)
- [ ] PBKDF2
- [ ] Support multiple for migration purposes

**Q3.2**: Should two-factor authentication (2FA) be supported?
- [X] Yes, TOTP (Google Authenticator, etc.)
- [ ] Yes, SMS-based
- [ ] Yes, both
- [ ] No, not initially
- [ ] Future consideration

**Q3.3**: Are there specific compliance requirements?
- [ ] GDPR (EU data protection)
- [ ] HIPAA (healthcare)
- [ ] SOC 2
- [ ] Other: ___________
- [X] None specific

**Q3.4**: Should the system support end-to-end encryption (like PGP/GPG)?
- [X] Yes, built-in
- [ ] No, but don't interfere with client-side encryption
- [ ] Future consideration

**Q3.5**: Greylisting preference?
- [X] Enable by default
- [ ] Configurable per domain
- [ ] Don't implement initially
- [ ] Never (concerns about delayed delivery)

## 4. Email Features

**Q4.1**: Auto-reply/vacation message priority?
- [ ] Must have in initial version
- [ ] Nice to have
- [X] Can be added later
- [ ] Not needed

**Q4.2**: Email forwarding requirements?
- [X] Simple forwarding (one-to-one)
- [X] Multiple destination forwarding
- [X] Conditional forwarding (based on rules)
- [ ] Keep copy when forwarding option

**Q4.3**: Should Sieve filtering be supported (RFC 5228)?
- [X] Yes, critical feature
- [ ] Nice to have
- [ ] Not needed (rely on client-side filtering)

**Q4.4**: Shared mailboxes/shared folders needed?
- [X] Yes, with full ACL support
- [ ] Yes, simple sharing only
- [ ] Not initially
- [ ] Not needed

**Q4.5**: Mailing list support?
- [ ] Yes, built-in mailing lists
- [ ] No, use external tools (Mailman, etc.)
- [X] Future consideration

## 5. CalDAV/CardDAV Details

**Q5.1**: CalDAV priority level?
- [X] Critical - need from day 1
- [ ] Important - within first month
- [ ] Nice to have - can wait
- [ ] Low priority

**Q5.2**: CardDAV priority level?
- [X] Critical - need from day 1
- [ ] Important - within first month
- [ ] Nice to have - can wait
- [ ] Low priority

**Q5.3**: Which clients must be supported?
- [X] Thunderbird
- [X] Apple Mail/Calendar/Contacts
- [X] iOS devices
- [X] Android devices
- [X] Microsoft Outlook
- [X] Evolution (Linux)
- [ ] Other: ___________

**Q5.4**: Calendar features required?
- [ ] Basic event storage/sync
- [ ] Event invitations and RSVP
- [ ] Recurring events
- [ ] Reminders/alarms
- [ ] Resource booking (rooms, equipment)
- [ ] Calendar sharing
- [X] All of the above

## 6. Management and Administration

**Q6.1**: Preferred management API format?
- [X] REST (JSON)
- [ ] gRPC
- [ ] GraphQL
- [ ] Both REST and gRPC
- [ ] No preference

**Q6.2**: Who will manage the system?
- [ ] Technical administrators (CLI comfortable)
- [ ] Mix of technical and non-technical
- [X] Non-technical users (need web UI)

**Q6.3**: Web-based admin interface needed?
- [X] Critical - need it ASAP
- [ ] Important - but can start with API/CLI
- [ ] Not needed - API is sufficient
- [ ] Eventually, but not initially

**Q6.4**: User self-service portal needed?
- [X] Yes - users should manage their own passwords, aliases, etc.
- [ ] No - all managed by administrators
- [ ] Future consideration

**Q6.5**: Preferred monitoring integration?
- [ ] Prometheus/Grafana
- [ ] ELK stack (Elasticsearch, Logstash, Kibana)
- [ ] Datadog
- [ ] New Relic
- [ ] Generic StatsD/syslog
- [ ] Other: ___________
- [X] None/undecided

## 7. TLS/SSL Certificates

**Q7.1**: How will TLS certificates be managed?
- [ ] Manual configuration (I'll provide certificates)
- [X] Let's Encrypt automatic integration (ACME) with Cloudflare DNS by default
- [ ] Certificate management via API
- [ ] Mix of manual and automatic

**Q7.2**: Will you need SNI support (different certs per domain)?
- [X] Yes, each domain has its own certificate
- [ ] No, single wildcard or multi-domain certificate
- [ ] Mix of both

## 8. Anti-Spam and Anti-Virus

**Q8.1**: ClamAV deployment preference?
- [X] Run ClamAV on same server
- [ ] Remote ClamAV service
- [ ] Multiple ClamAV instances (load balancing)
- [ ] Undecided

**Q8.2**: What should happen when ClamAV finds a virus?
- [ ] Reject the email
- [ ] Quarantine for admin review
- [ ] Deliver with warning header
- [X] Configurable per domain/user

**Q8.3**: Additional anti-spam tools needed?
- [X] SpamAssassin and Rspamd integration (https://github.com/teamwork/spamc)
- [ ] Rspamd integration
- [ ] Custom scoring system
- [ ] Just rely on SPF/DKIM/DMARC
- [ ] Other: ___________

**Q8.4**: Should there be a spam quarantine?
- [X] Yes, users can review quarantined messages
- [ ] Yes, admin-only review
- [ ] No, just mark and deliver or reject

## 9. Backup and Disaster Recovery

**Q9.1**: Backup strategy preference?
- [X] Built-in backup functionality
- [ ] Rely on external tools (Postgresdump, etc.)
- [ ] Both

**Q9.2**: Backup schedule?
- [ ] Continuous/real-time
- [ ] Hourly
- [X] Daily
- [ ] Weekly
- [ ] Manual only

**Q9.3**: Backup retention?
- Days: __30 days_________
- Weeks: ___________
- Months: ___________

## 10. Migration and Import

**Q10.1**: Will you need to migrate from an existing system?
- [X] No, fresh start
- [ ] Yes, from: __Gmail_________

**Q10.2**: If yes, what needs to be migrated?
- [ ] User accounts
- [ ] Existing emails
- [ ] Aliases/forwards
- [ ] Calendar events
- [ ] Contacts
- [ ] Domain configurations

**Q10.3**: Import format needed?
- [ ] mbox format
- [ ] Maildir format
- [ ] IMAP sync (from other server)
- [ ] CSV for users/domains
- [ ] Other: ___________

## 11. Logging and Auditing

**Q11.1**: What level of email logging do you need?
- [ ] Basic (connections, auth, errors)
- [ X Standard (+ send/receive events)
- [ ] Detailed (+ full headers)
- [ ] Complete (+ full message content)

**Q11.2**: Log retention period?
- [ ] 7 days
- [X] 30 days
- [ ] 90 days
- [ ] 1 year
- [ ] Forever
- [ ] Other: ___________

**Q11.3**: Should there be an audit trail for admin actions?
- [X] Yes, full audit log
- [ ] Basic logging only
- [ ] Not needed

## 12. Additional Features

**Q12.1**: Should there be rate limiting?
- [ ] Yes, per user
- [ ] Yes, per IP address
- [ ] Yes, per domain
- [X] All of the above
- [ ] Not needed

**Q12.2**: Quota enforcement approach?
- [ ] Hard limit (reject when over quota)
- [ ] Soft limit (warn but allow)
- [X] Both (warning threshold + hard limit)

**Q12.3**: Should there be a catch-all address option?
- [X] Yes, per domain
- [ ] No
- [ ] Optional feature

**Q12.4**: Webhook support priority?
- [X] Critical - needed for integrations
- [ ] Nice to have
- [ ] Not needed initially

**Q12.5**: Is internationalization (i18n) needed?
- [ ] Yes, support for multiple languages
- [ ] English only is fine
- [X] Future consideration

## 13. Development and Testing

**Q13.1**: Do you have a staging/testing environment?
- [X] Yes
- [ ] No, will test in production
- [ ] Will set one up

**Q13.2**: Integration testing with real email services?
- [X] Yes, need to test with Gmail, Outlook, etc.
- [ ] Just local testing is fine initially

**Q13.3**: Performance testing requirements?
- [ ] Load testing with X concurrent users: ___________
- [ ] Stress testing to find limits
- [X] Not needed initially

## 14. Documentation Preferences

**Q14.1**: What documentation is most important to you?
- [ ] API documentation
- [ ] Administrator guide
- [ ] User guide
- [ ] Developer/contributor guide
- [ ] Architecture documentation
- [X] All of the above

**Q14.2**: Documentation format preference?
- [X] Markdown files in repo
- [ ] Wiki
- [ ] Generated from code (Swagger, etc.)
- [ ] Dedicated docs site
- [ ] No preference

## 15. Open Questions

**Q15.1**: Are there any specific features not mentioned that are critical for your use case?

Webmail like Gmail with Categories___________________________________________________________
___________________________________________________________
___________________________________________________________

**Q15.2**: Are there any existing mail servers or systems you'd like this to be compatible with or similar to?

___________________________________________________________
___________________________________________________________

**Q15.3**: What is your expected timeline for production deployment?

ASAP___________________________________________________________

**Q15.4**: What is your budget for infrastructure (server specs, bandwidth)?

Minimal Open Source Project___________________________________________________________

**Q15.5**: Any other concerns, requirements, or questions?

The key is to keep this application as simple to install and manage as possible without sacrificing features.
Use sqlite as the base database with ALL configuration storage within it.
___________________________________________________________
___________________________________________________________

---

## Submission Instructions

Please fill out this questionnaire and provide your answers. You can:
1. Edit this file directly and save it
2. Copy to a new file with your answers
3. Provide answers in another format

This information will help create a more precise implementation plan and prioritize features appropriately.
