# Webmail Client

Modern, Gmail-like webmail client built with Nuxt 3, Vue 3, and Tailwind CSS.

## Features

- **Authentication**: JWT-based authentication with user sessions
- **Mailbox Management**: Browse folders (Inbox, Sent, Drafts, etc.)
- **Message List**: Gmail-style message list with avatars and previews
- **Message Detail**: Full message view with attachments
- **Rich Text Composer**: TipTap-based WYSIWYG editor
- **Attachments**: Drag-and-drop file uploads
- **Dark Mode**: System-aware dark mode support
- **Responsive**: Mobile-friendly design
- **Search**: Full-text message search (coming soon)
- **Keyboard Shortcuts**: Gmail-style keyboard navigation (coming soon)
- **PWA**: Offline capability (coming soon)

## Development

```bash
# Install dependencies
pnpm install

# Start development server
pnpm dev

# Build for production
pnpm build

# Preview production build
pnpm preview
```

## Project Structure

```
webmail/
├── assets/
│   └── css/              # Global styles
├── components/
│   ├── mailbox/          # Mailbox sidebar
│   ├── message/          # Message list and detail
│   └── composer/         # Email composer
├── composables/
│   └── useAuth.ts        # Authentication composable
├── layouts/
│   ├── default.vue       # Default layout
│   └── mail.vue          # Mail layout
├── pages/
│   ├── index.vue         # Landing page
│   ├── login.vue         # Login page
│   └── mail/             # Mail pages
├── stores/
│   ├── auth.ts           # Authentication store
│   └── mail.ts           # Mail store
├── lib/
│   └── utils.ts          # Utility functions
└── nuxt.config.ts        # Nuxt configuration
```

## API Integration

The webmail client connects to the backend API at `/api/v1/webmail/`:

- `GET /mailboxes` - List mailboxes
- `GET /mailboxes/:id/messages` - List messages in mailbox
- `GET /messages/:id` - Get message details
- `POST /messages` - Send new message
- `DELETE /messages/:id` - Delete message
- `POST /messages/:id/move` - Move message to folder
- `POST /messages/:id/flags` - Update message flags
- `GET /search` - Search messages
- `GET /attachments/:id` - Download attachment

## Technologies

- **Nuxt 3** - Vue framework
- **Vue 3** - UI framework
- **Tailwind CSS 4** - Styling
- **TipTap** - Rich text editor
- **Pinia** - State management
- **Axios** - HTTP client
- **VueUse** - Vue composables
- **Lucide Icons** - Icon library
- **Radix Vue** - Headless UI components
