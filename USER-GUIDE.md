# gomailserver - User Guide

**Version**: 1.0  
**Date**: January 11, 2026  
**gomailserver Version**: 0.10.0+

---

## Table of Contents

1. [Getting Started](#getting-started)
2. [Webmail Interface](#webmail-interface)
3. [Managing Email](#managing-email)
4. [Contacts and Calendar](#contacts-and-calendar)
5. [Account Settings](#account-settings)
6. [Security Best Practices](#security-best-practices)
7. [Troubleshooting](#troubleshooting)

---

## Getting Started

### Accessing gomailserver

gomailserver provides a modern webmail interface accessible at:

```
https://your-mail-server.com/webmail/
```

### First Login

1. Open your web browser and navigate to your webmail URL
2. Enter your email address and password
3. Click "Login"
4. If 2FA is enabled, enter your TOTP code from your authenticator app

### Initial Setup

After first login:

1. **Welcome Guide**: You'll be presented with a quick tour of the interface
2. **Profile Setup**: Complete your profile information (name, signature)
3. **Settings Review**: Configure your preferences

---

## Webmail Interface

### Layout Overview

The webmail interface consists of several sections:

- **Sidebar**: Folders, labels, and navigation
- **Message List**: Main email viewing area
- **Reading Pane**: Full message display
- **Composer**: Compose new emails
- **Contacts**: Address book integration
- **Calendar**: Upcoming events and calendar management

### Folder Structure

**Default Folders**:
- **INBOX**: Incoming emails
- **Drafts**: Saved drafts (auto-created when you first send an email)
- **Sent**: Emails you've sent
- **Spam**: Messages flagged as spam by spam filters
- **Trash**: Deleted emails (30-day retention)

**Folder Actions**:
- Click folder to view contents
- Right-click for folder options (rename, empty, mark read)
- Drag and drop messages between folders

### Reading Emails

**Message List View**:
- Shows sender, subject, date, size
- Unread messages bolded
- Starred messages with star icon
- Attachment indicator (paperclip icon)
- Hover for quick actions (reply, forward, delete)

**Reading Pane**:
- Full message display with headers
- HTML and plain text views
- Attachment list with download options
- Reply/Forward buttons in toolbar
- Next/Previous message navigation

**Keyboard Shortcuts**:
- `C` or `N`: Next/Previous message
- `R`: Reply to message
- `F`: Forward message
- `Delete` or `Backspace`: Delete message
- `Escape`: Return to message list
- `?`: Show keyboard shortcuts help

### Composing Emails

**Accessing Composer**:
- Click "Compose" button in toolbar
- Press `C` from message list
- Right-click any email and select "Reply" or "Forward"

**Composer Features**:

**Rich Text Editor**:
- Formatting toolbar (bold, italic, underline, lists, links)
- Font selection and size
- Text and background colors
- Alignment options

**Recipients**:
- **To**: Main recipients (click to add multiple)
- **Cc**: Carbon copy recipients
- **Bcc**: Blind carbon copy recipients
- Type-ahead autocomplete from your contacts

**Attachments**:
- Click attachment icon (paperclip)
- Select files from your computer
- Maximum attachment size: 25MB (configurable by administrator)
- Supported formats: PDF, images, documents, archives

**Composer Options**:
- **Priority**: Normal (default), High, Low
- **Request Read Receipt**: Get confirmation when recipient reads email
- **Schedule Send**: Send at specific date/time
- **Save Draft**: Save as draft to send later
- **Spell Check**: Automatic spell checking

**Drafts Management**:
- **Auto-save**: Composing emails auto-save every 30 seconds
- **Drafts Folder**: All drafts stored in Drafts folder
- **Manual Save**: Press `Ctrl+S` or click "Save Draft" button
- **Discard**: Click "X" or press `Esc` to close without saving

---

## Managing Email

### Message Actions

**Multiple Selection**:
- Click checkbox next to messages
- Use `Shift+Click` for range selection
- `Ctrl+A` to select all messages in current view

**Bulk Actions**:
- **Delete**: Move to trash
- **Mark Read/Unread**: Update read status
- **Move to Folder**: Organize emails into folders
- **Add Label**: Apply custom labels for organization
- **Archive**: Archive emails out of inbox

**Individual Message Actions**:
- **Reply**: Reply to sender
- **Reply All**: Reply to all recipients
- **Forward**: Forward to new recipient
- **Print**: Print email (uses browser print dialog)
- **Download**: Save email as .eml file
- **View Source**: View raw email headers and body

### Searching Emails

**Search Features**:
- **Quick Search**: Search from anywhere with `Ctrl+K` or `Cmd+F`
- **Search By**:
  - From: Filter by sender
  - To: Filter by recipient
  - Subject: Search subject lines
  - Date: Search by date range
  - Attachment: Filter by attachments
  - Body: Search message content

**Advanced Search**:
- `has:attachment` - Find emails with attachments
- `is:unread` - Find only unread emails
- `is:starred` - Find only starred emails
- `in:folder` - Search specific folder
- Boolean operators: `AND`, `OR`, `NOT` (e.g., `from:john@example.com is:unread`)

### Filters and Labels

**Creating Filters**:
- Click "Filter" button in toolbar
- Define filter conditions (sender, subject, date, etc.)
- Save filter for quick access
- Apply multiple filters for complex searches

**Using Labels**:
- Color-coded labels for organization
- Create custom labels (e.g., "Work", "Personal", "Finance")
- Apply labels to messages via right-click menu
- View messages by label in sidebar

### Spam Management

**Reporting Spam**:
- Select message(s)
- Click "Report Spam" in toolbar
- Choose spam reason (if required by admin)
- Message moved to Spam folder

**Whitelist Senders**:
- Open spam message
- Click "Not Spam" button
- Sender added to whitelist

**Spam Folder**:
- Review spam folder regularly
- False positives: Move to inbox and mark as not spam
- False negatives: Report to admin for spam filter improvement

---

## Contacts and Calendar

### Contacts (CardDAV Integration)

Contacts are synchronized with your CardDAV address book, accessible across devices and email clients.

**Accessing Contacts**:
- Click "Contacts" in sidebar
- Search contacts by name or email
- View contact details with all information

**Creating Contacts**:
- Click "New Contact" button
- Fill in contact details:
  - **Name**: Full name
  - **Email**: Email address
  - **Phone**: Phone number (optional)
  - **Organization**: Company or organization
  - **Notes**: Additional information
- Click "Save"

**Contact Fields**:
- **Name**: Required
- **Email**: Required, must be unique
- **Phone**: Optional
- **Mobile**: Optional mobile number
- **Address**: Work, home, other addresses
- **Organization**: Company or department
- **Title**: Job title or position
- **Notes**: Free-form notes
- **Photo**: Profile picture (optional)

**Contact Groups**:
- Organize contacts into groups (e.g., "Work", "Family", "Friends")
- Send group emails quickly
- Share groups with other users (if permissions allow)

**Contact Autocomplete in Composer**:
- Start typing in To/Cc/Bcc fields
- Contacts appear as suggestions
- Arrow keys to navigate suggestions
- Enter to select contact
- Auto-fills email address

### Calendar (CalDAV Integration)

Calendar events are synchronized with your CalDAV calendar, accessible across devices and email clients.

**Accessing Calendar**:
- Click "Calendar" in sidebar or "Calendar" icon in toolbar
- View month/week/day views
- Click events to see details

**Creating Events**:
- Click "New Event" button
- Fill in event details:
  - **Title**: Event title
  - **Start Date/Time**: When event starts
  - **End Date/Time**: When event ends
  - **Location**: Event location
  - **Description**: Event details
  - **Attendees**: Add participants (autocomplete from contacts)
  - **Reminder**: Set reminder (5 min, 15 min, 1 hour, etc.)
- Click "Save"

**Event Types**:
- **Meeting**: Attendees expected (default)
- **Appointment**: One-on-one meeting
- **Reminder**: Personal reminder
- **All Day**: Full-day event
- **Recurring**: Daily, weekly, monthly events

**Calendar Views**:
- **Month View**: Full month overview
- **Week View**: Week by week
- **Day View**: Daily schedule
- **Agenda View**: List of upcoming events

**Meeting Invitations**:
- Receive calendar invitations via email
- Accept, tentatively accept, or decline
- Invitations automatically added to your calendar
- Your responses sent to organizer

---

## Account Settings

### Accessing Settings

Click your profile picture or name in the top-right corner, then select "Settings" from the dropdown menu.

### General Settings

**Personal Information**:
- **Full Name**: Display name
- **Email Address**: Cannot be changed (contact admin)
- **Time Zone**: Select your timezone for correct timestamps
- **Language**: Interface language (if multiple available)
- **Date/Time Format**: 12-hour vs 24-hour, date format

**Display Preferences**:
- **Messages Per Page**: 10, 25, 50, 100
- **Preview Pane**: Show/Hide (on right side)
- **Message List Density**: Comfortable, Compact, Cozy
- **Theme**: Light or Dark mode
- **Font Size**: Adjust text size

### Email Settings

**Signature**:
- **Plain Text Signature**: Simple text signature
- **Rich Text Signature**: Formatted signature with links, images
- **Use Signature Automatically**: Add to all outgoing emails
- **Signature Position**: Above or below reply

**Composing Options**:
- **Default Font**: Font for new emails
- **Default Font Size**: Size for new emails
- **Reply Behavior**: Reply to sender or reply to all
- **Forward Behavior**: Include original message as attachment or inline

**Forwarding**:
- **Enable Forwarding**: Forward all incoming emails to another address
- **Forward To**: Target email address
- **Keep Copy**: Keep copy in inbox or delete after forwarding

### Security Settings

**Change Password**:
- Enter current password
- Enter new password (8+ characters, recommended)
- Confirm new password
- Click "Change Password"

**Password Requirements**:
- Minimum 8 characters
- Recommended 12+ characters
- Mix of letters, numbers, and special characters

**2FA (Two-Factor Authentication)**:
- **Enable/Disable TOTP**: Toggle 2FA for your account
- **Setup TOTP**:
  1. Download authenticator app (Google Authenticator, Authy, etc.)
  2. Scan QR code displayed in webmail
  3. Enter 6-digit code to verify
- **Recovery Codes**: Save backup codes for account recovery
- **Recovery Codes**: Save backup codes for account recovery
  4. Store in secure location

### Vacation Auto-Reply

**Enable Auto-Reply**:
- Toggle on vacation auto-reply
- **Subject**: Auto-reply subject line
- **Message**: Auto-reply message content
- **Start Date**: When auto-reply starts
- **End Date**: When auto-reply ends
- **Auto-Reply To**: Specific addresses or all incoming
- **Save Drafts**: Don't send auto-replies to mailing lists

---

## Security Best Practices

### Password Security

**Create Strong Passwords**:
- Use 12+ characters with mix of upper/lower case, numbers, symbols
- Don't reuse passwords across services
- Don't use personal information (birthdays, names)
- Use password manager for storage
- Change passwords regularly (every 90 days recommended)

### Protecting Your Account

**2FA Best Practices**:
- Always enable 2FA when available
- Keep backup recovery codes in secure location
- Don't share QR codes or secret keys
- Use reputable authenticator apps
- Immediately revoke codes if device is lost/stolen

**Session Security**:
- Always log out when using shared/public computers
- Don't save passwords in browser
- Use private/incognito mode on public devices
- Clear browser cache regularly

### Email Safety

**Recognizing Phishing**:
- Verify sender before clicking links
- Check URL domain (e.g., paypa1.com instead of paypal.com)
- Be suspicious of urgent/pressure tactics
- Verify unexpected attachments before downloading
- Hover over links to see actual URL before clicking

**Handling Attachments Safely**:
- Scan attachments for viruses (admin-provided ClamAV)
- Don't open unexpected attachments from unknown senders
- Be cautious of executable files (.exe, .zip, .rar)
- Download and scan before opening if in doubt

**Reporting Suspicious Activity**:
- Report phishing attempts to admin
- Forward suspicious emails to security team
- Mark as spam to help train filters

---

## Troubleshooting

### Common Issues

**Can't Login**:
- Check email address for typos
- Verify password is correct
- Clear browser cache and cookies
- Try different browser
- Check if account is locked (contact admin)
- Verify 2FA code is correct

**Emails Not Sending**:
- Check sent folder for errors
- Check queue status (if you have admin access)
- Verify recipient email address
- Check attachment size limits
- Check for rate limiting (if you send many emails quickly)

**Not Receiving Emails**:
- Check inbox, spam folder, and trash
- Verify email address is working (ask someone to send test)
- Check your domain's MX records (if self-managed)
- Check your account quota limits
- Verify email filters aren't blocking messages

**Emails Going to Spam**:
- Check spam folder regularly
- Review spam filters in settings
- Whitelist legitimate senders
- Mark false positives as "Not Spam" to train filters
- Report spam to admin for filter improvement

**Slow Performance**:
- Check internet connection speed
- Close unused browser tabs
- Clear browser cache
- Try different browser
- Disable browser extensions
- Check if you have large emails with many attachments

**Sync Issues with Contacts/Calendar**:
- Check your device's internet connection
- Refresh webmail page (F5 or Ctrl+R)
- Check if CalDAV/CardDAV is enabled (contact admin)
- Verify contact/calendar app settings
- Check for error messages

### Getting Help

**Built-in Help**:
- Click "?" icon in webmail interface
- View keyboard shortcuts
- Read user guide (this document)

**Contacting Support**:
- Email your administrator for assistance
- Check organization's help desk or knowledge base
- Provide detailed error messages when reporting issues
- Include screenshots when reporting UI issues

### Webmail vs Email Client

**Advantages of Webmail**:
- Access from anywhere with web browser
- No software installation required
- Contacts and calendar sync across devices
- Full-text search of all emails
- Automatic spam filtering
- Consistent interface across devices

**Advantages of Email Client**:
- Offline access to emails
- Faster response time for heavy email users
- Advanced filtering and rules
- Integration with other applications
- No storage quota limitations

**When to Use Each**:
- **Webmail**: Occasional access, travel, shared computers
- **Email Client**: Daily use, offline access, high-volume email workflows

---

## Appendix

### Keyboard Shortcuts Reference

**Global Shortcuts**:
- `C` or `/`: Compose new email
- `Ctrl+A`: Select all messages
- `Escape`: Deselect all or close composer
- `Ctrl+F`: Search emails
- `Ctrl+K`: Quick search (alternative to Ctrl+F)
- `?`: Help

**Message List Shortcuts**:
- `J`/`K` or `Up`/`Down`: Navigate messages
- `Enter` or `Space`: Open selected message
- `X` or `Delete`: Delete selected messages
- `Shift+C`: Mark selected as read
- `Shift+U`: Mark selected as unread
- `Shift+S` Star selected messages
- `Shift+I`: Archive selected messages

**Reading Pane Shortcuts**:
- `R`: Reply to message
- `F`: Forward message
- `Ctrl+Enter`: Send email (in composer)
- `Esc`: Return to message list
- `P`: Print message
- `Shift+P`: Print conversation

### Mobile Responsive Design

The webmail interface is fully responsive and works on:
- Desktop computers (1920px+ width)
- Tablets (768px-1024px width)
- Mobile phones (320px-767px width)
- Touch-friendly controls for mobile devices

### Dark Mode

Switch between light and dark themes:
- Click profile picture → "Settings" → "Display"
- Select "Dark" or "Light" theme
- Theme persists across sessions
- Reduces eye strain in low-light environments

---

**End of User Guide**

For additional help, see:
- **README**: https://github.com/btafoya/gomailserver
- **Admin Guide**: https://github.com/btafoya/gomailserver/blob/main/ADMIN-GUIDE.md
- **API Documentation**: See API reference guide
- **Support**: Contact your administrator or IT support team
