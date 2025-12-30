# Future efatures to be added

## DNS Management for Cloudflare (with others providers at later dates).

- Provide automated DNS managment for admin approval during congiuration as part of the UI for adding a new domain. SRV, MX, DKIM, DMARC, and SPF records. Using Cloudflare API for management already present in 3.6 Let's Encrypt Integration from TASKS3.md 3.6 Let's Encrypt Integration.

### Example domain records to add/manage ###

_imap._tcp.example.com          SRV  10  20  143  mail.example.com.
_imaps._tcp.example.com         SRV  0   1   993  .
_pop3._tcp.example.com          SRV  0   1   110  .
_pop3s._tcp.example.com         SRV  0   1   995  .
_smtp._tcp.example.com.         SRV  0   1   25   .
_submission._tcp.example.com.   SRV  10  20  587  mail.example.com.
_autodiscover._tcp.example.com.  0   443 service.example-provider.com.


## Mail Client Autoconfiguration

- Build into the webmail interface and API fully automated mailbox configuration from Apple's Mobileconfig, Microsoft's Autodiscover and Mozilla's Autoconfig in one tool using https://github.com/philband/go-autoconfig as the logic reference. All data for this comes from the domain records stored in sqlite and the url domain passed form the mail client.

### Url Examples ###

- https://autoconfig.thunderbird.net/v1.1/example.com
- https://autoconfig.example.com/mail/config-v1.1.xml?emailaddress=alice@example.com

### Reference ###

- go-autoconfig - A go libray implementation that contains the logic to implement - https://github.com/philband/go-autoconfig
- https://rseichter.github.io/automx2/ for implementation references
- context7 mcp.
