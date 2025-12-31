# gomailserver Mail Server (`github.com/btafoya/gomailserver`)

## Autonomous Work Mode

**Claude Code is authorized to work autonomously on this project.** When given a task:

1. **Proceed without asking for confirmation** - Execute the full task from start to finish
2. **Make reasonable decisions** - Use best judgment for implementation details
3. **Follow established patterns** - Match existing code style and project conventions
4. **Complete the work** - Don't stop mid-task or leave partial implementations
5. **Report results** - Summarize what was done when complete
6. **Compact conversation** - When you are using compact, please focus on test output and code changes

### When to Ask Questions
- Requirements are genuinely ambiguous with multiple valid interpretations
- Security implications require explicit user approval
- Destructive operations (deleting data, force push) need confirmation
- The task fundamentally contradicts project requirements

### When NOT to Ask Questions
- Implementation details (which pattern to use, naming conventions)
- File organization decisions that follow existing patterns
- Code style choices that match the codebase
- Standard software engineering decisions

## Guidelines

### ❌ Do NOT Include:
- "Generated with Claude Code" in commit messages
- "Co-Authored-By: Claude Sonnet" in commits
- AI attribution in code comments
- References to Claude in documentation footer/header

### ✅ DO Include:
- Your name and email as the commit author
- Professional commit messages describing WHAT changed
- Standard documentation without AI tool references
- Human authorship for all contributions

## Commit Message Standards

**Good Commit Messages:**
```
Add payment gateway integration with Stripe and PayPal
Update RBAC schema to support multi-agency isolation
Implement webhook handlers for automatic payment confirmation
```

**Bad Commit Messages:**
```
Generated with Claude Code
Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
Add feature (created with AI assistance)
```

## Rationale

- **Professionalism**: Code should reflect human authorship
- **Clarity**: Commit history should describe changes, not tools used
- **Standards**: Follow industry-standard Git practices
- **Ownership**: Maintain clear project ownership and responsibility

## Tool Usage

Claude Code is a development assistant tool, like an IDE or linter. You wouldn't attribute your IDE in commits, and the same applies to AI coding assistants.

Use Claude Code to:
- Generate code snippets and boilerplate
- Review and improve code quality
- Write documentation and specs
- Debug and troubleshoot issues

But always commit and sign work as **btafoya**.

## MCP Tools to use

- Shadcn MCP for Vue frontend UI design
- MemmoryGraph MCP for memory
- Serena MCP for memory and tools
- Context7 MCP for library and implementation reference
- Playwright MCP for Browseer testing

### Operating Principles
- Work **autonomously**: do not ask for human confirmation unless the issue is ambiguous or lacks required information to proceed safely.
- Handle **one issue at a time** from start to finish before moving to the next.
- Prefer **minimal, correct changes** that align with gomailserver’s architecture and style.
- If a proposed fix affects public behavior or compatibility, document the impact explicitly in the issue file and ensure tests cover it.

### Safety & Scope
- Do not introduce insecure defaults. If an issue touches authentication, TLS, DKIM/SPF/DMARC, or mail routing, include a brief security impact note in the issue file.
- Avoid configuration-breaking changes unless the issue explicitly requires it; if unavoidable, document migration steps.

## Project Overview

gomailserver is a composable, all-in-one mail server written in Go. It implements all the functionality required to run a modern, secure e-mail server. It is designed to be a single daemon that replaces complex stacks of software like Postfix, Dovecot, OpenDKIM, OpenSPF, and OpenDMARC. This simplifies configuration and maintenance.

Key features include:
- SMTP for sending (MTA) and receiving (MX) messages.
- IMAP for message storage and access.
- Support for modern security standards like DKIM, SPF, DMARC, DANE, and MTA-STS.
- Modular architecture with a clear configuration file format.

The project is structured as a standard Go application, with a main entrypoint in `cmd/gomailserver/main.go`, core logic separated into `internal/` packages, and a reusable `framework/` for the module system and configuration parsing.

## Building and Running

### Building from Source

The project includes a shell script for building the binaries.

- **To build the server:**
  ```sh
  ./build.sh
  ```
  This will create the `gomailserver` executable in the `./build` directory.

- **To create a static, self-contained build:**
  ```sh
  ./build.sh --static
  ```

- **To install the server:**
  ```sh
  ./build.sh install
  ```
  This will install the binary to `/usr/local/bin` and the configuration to `/etc/gomailserver` by default. These paths can be customized with `--prefix` and `--destdir`.

### Running the Server

The server is started using the `run` subcommand.

```sh
gomailserver run --config /path/to/gomailserver.conf
```

The default configuration file path is `/etc/gomailserver/gomailserver.conf`. An example `gomailserver.conf` is available in the root of the repository.

### Running with Docker

A `Dockerfile` is provided for building and running gomailserver in a container.

- **To build the Docker image:**
  ```sh
  docker build -t gomailserver .
  ```

- **To run the container:**
  ```sh
  docker run -p 25:25 -p 143:143 -p 465:465 -p 587:587 -p 993:993 -v /path/to/data:/data gomailserver
  ```
  The container expects a volume at `/data` for persistent state, including configuration (`/data/gomailserver.conf`), certificates, and mail storage.

## Development

### Testing

The project has both unit/module tests and a suite of integration tests.

- **Run unit and module tests:**
  ```sh
  go test ./...
  ```

- **Run integration tests:**
  ```sh
  cd tests/
  ./run.sh
  ```

### Linting

The project uses `golangci-lint` to enforce code quality and style. The configuration is in `.golangci.yml`.

- **To run the linter:**
  ```sh
  golangci-lint run
  ```

### Dependencies

The project uses Go modules for dependency management. Dependencies are declared in the `go.mod` file.


## Doumentation

### PR

Review PR.md to determine the focus on the current tasks to address any open issues in order to provide pull requests to the parent repo.

### Issues

Docuemnt all issues and their status using markdown with them naming format: ISSUE{number}.md

