# Phase 6: Sieve Filtering (Weeks 18-19)

**Status**: Not Started
**Priority**: Full Feature (Post-MVP)
**Estimated Duration**: 1-2 weeks
**Dependencies**: Phase 1 (Core Mail), Phase 2 (Anti-Spam)

---

## Overview

Implement RFC 5228 Sieve mail filtering language with common extensions, ManageSieve protocol for remote script management, and a visual rule editor in the web interface.

---

## 6.1 Sieve Interpreter [FULL]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| SV-001 | Sieve base implementation (RFC 5228) | [ ] | F-002 |
| SV-002 | Variables extension (RFC 5229) | [ ] | SV-001 |
| SV-003 | Vacation extension (RFC 5230) | [ ] | SV-001 |
| SV-004 | Relational extension (RFC 5231) | [ ] | SV-001 |
| SV-005 | Subaddress extension (RFC 5233) | [ ] | SV-001 |
| SV-006 | Spamtest extension (RFC 5235) | [ ] | SV-001, AS-002 |

### SV-001: Sieve Parser and Executor

```go
// pkg/sieve/parser.go
package sieve

type Parser struct {
    lexer *Lexer
}

type Script struct {
    Requires []string
    Commands []Command
}

type Command interface {
    Execute(ctx *ExecutionContext) error
}

// Control commands
type IfCommand struct {
    Test     Test
    Commands []Command
    ElsIf    []ElsIfBlock
    Else     []Command
}

type ElsIfBlock struct {
    Test     Test
    Commands []Command
}

type RequireCommand struct {
    Capabilities []string
}

type StopCommand struct{}

// Action commands
type KeepCommand struct{}

type DiscardCommand struct{}

type FileIntoCommand struct {
    Mailbox string
    Create  bool
    Copy    bool
}

type RedirectCommand struct {
    Address string
    Copy    bool
}

type RejectCommand struct {
    Message string
}

type FlagCommand struct {
    Action string // setflag, addflag, removeflag
    Flags  []string
}

// Tests
type Test interface {
    Evaluate(ctx *ExecutionContext) bool
}

type TrueTest struct{}
type FalseTest struct{}

type AddressTest struct {
    AddressPart  string // localpart, domain, all
    MatchType    string // is, contains, matches
    Comparator   string
    Headers      []string
    Keys         []string
}

type HeaderTest struct {
    MatchType  string
    Comparator string
    Headers    []string
    Keys       []string
}

type SizeTest struct {
    Over  bool
    Under bool
    Limit int64
}

type ExistsTest struct {
    Headers []string
}

type AllOfTest struct {
    Tests []Test
}

type AnyOfTest struct {
    Tests []Test
}

type NotTest struct {
    Test Test
}

func (p *Parser) Parse(script string) (*Script, error) {
    p.lexer = NewLexer(script)
    return p.parseScript()
}
```

### SV-001: Sieve Executor

```go
// pkg/sieve/executor.go
package sieve

type Executor struct {
    mailboxService *service.MailboxService
    spamService    *service.SpamService
    logger         logger.Logger
}

type ExecutionContext struct {
    Message     *domain.Message
    RawMessage  []byte
    Headers     map[string][]string
    Envelope    *Envelope
    SpamScore   float64
    Variables   map[string]string
    Flags       []string
    ImplicitKeep bool
    Actions     []Action
}

type Envelope struct {
    From string
    To   string
}

type Action struct {
    Type    string
    Params  map[string]interface{}
}

func (e *Executor) Execute(script *Script, ctx *ExecutionContext) (*Result, error) {
    ctx.ImplicitKeep = true
    ctx.Variables = make(map[string]string)
    ctx.Actions = []Action{}

    for _, cmd := range script.Commands {
        if err := cmd.Execute(ctx); err != nil {
            if err == ErrStop {
                break
            }
            return nil, err
        }
    }

    // Apply implicit keep if no explicit action was taken
    if ctx.ImplicitKeep && !hasExplicitAction(ctx.Actions) {
        ctx.Actions = append(ctx.Actions, Action{Type: "keep"})
    }

    return &Result{Actions: ctx.Actions}, nil
}

func (c *IfCommand) Execute(ctx *ExecutionContext) error {
    if c.Test.Evaluate(ctx) {
        for _, cmd := range c.Commands {
            if err := cmd.Execute(ctx); err != nil {
                return err
            }
        }
        return nil
    }

    for _, elsif := range c.ElsIf {
        if elsif.Test.Evaluate(ctx) {
            for _, cmd := range elsif.Commands {
                if err := cmd.Execute(ctx); err != nil {
                    return err
                }
            }
            return nil
        }
    }

    if c.Else != nil {
        for _, cmd := range c.Else {
            if err := cmd.Execute(ctx); err != nil {
                return err
            }
        }
    }

    return nil
}

func (c *FileIntoCommand) Execute(ctx *ExecutionContext) error {
    ctx.Actions = append(ctx.Actions, Action{
        Type: "fileinto",
        Params: map[string]interface{}{
            "mailbox": c.Mailbox,
            "create":  c.Create,
        },
    })

    if !c.Copy {
        ctx.ImplicitKeep = false
    }

    return nil
}

func (c *DiscardCommand) Execute(ctx *ExecutionContext) error {
    ctx.Actions = append(ctx.Actions, Action{Type: "discard"})
    ctx.ImplicitKeep = false
    return nil
}
```

### SV-002: Variables Extension

```go
// pkg/sieve/extensions/variables.go
package extensions

type SetCommand struct {
    Modifiers []string // lower, upper, lowerfirst, upperfirst, quotewildcard, length
    Name      string
    Value     string
}

func (c *SetCommand) Execute(ctx *sieve.ExecutionContext) error {
    value := expandVariables(c.Value, ctx.Variables)

    for _, mod := range c.Modifiers {
        switch mod {
        case "lower":
            value = strings.ToLower(value)
        case "upper":
            value = strings.ToUpper(value)
        case "lowerfirst":
            if len(value) > 0 {
                value = strings.ToLower(string(value[0])) + value[1:]
            }
        case "upperfirst":
            if len(value) > 0 {
                value = strings.ToUpper(string(value[0])) + value[1:]
            }
        case "length":
            value = strconv.Itoa(len(value))
        }
    }

    ctx.Variables[c.Name] = value
    return nil
}

func expandVariables(s string, vars map[string]string) string {
    re := regexp.MustCompile(`\$\{([a-zA-Z_][a-zA-Z0-9_]*)\}`)
    return re.ReplaceAllStringFunc(s, func(match string) string {
        name := match[2 : len(match)-1]
        if val, ok := vars[name]; ok {
            return val
        }
        return ""
    })
}
```

### SV-003: Vacation Extension

```go
// pkg/sieve/extensions/vacation.go
package extensions

type VacationCommand struct {
    Days      int
    Subject   string
    From      string
    Addresses []string
    Mime      bool
    Handle    string
    Reason    string
}

func (c *VacationCommand) Execute(ctx *sieve.ExecutionContext) error {
    // Check if we should send a vacation response
    if !c.shouldRespond(ctx) {
        return nil
    }

    ctx.Actions = append(ctx.Actions, sieve.Action{
        Type: "vacation",
        Params: map[string]interface{}{
            "days":      c.Days,
            "subject":   c.Subject,
            "from":      c.From,
            "addresses": c.Addresses,
            "reason":    c.Reason,
            "mime":      c.Mime,
        },
    })

    return nil
}

func (c *VacationCommand) shouldRespond(ctx *sieve.ExecutionContext) bool {
    // Don't respond to mailing lists
    if hasHeader(ctx, "List-Id") || hasHeader(ctx, "List-Unsubscribe") {
        return false
    }

    // Don't respond to auto-generated messages
    if hasHeader(ctx, "Auto-Submitted") {
        autoSubmitted := getHeader(ctx, "Auto-Submitted")
        if autoSubmitted != "no" {
            return false
        }
    }

    // Check Precedence header
    if precedence := getHeader(ctx, "Precedence"); precedence != "" {
        if precedence == "bulk" || precedence == "junk" || precedence == "list" {
            return false
        }
    }

    return true
}
```

### SV-006: Spamtest Extension

```go
// pkg/sieve/extensions/spamtest.go
package extensions

type SpamtestTest struct {
    MatchType  string // value, count
    Comparator string
    Value      string
}

func (t *SpamtestTest) Evaluate(ctx *sieve.ExecutionContext) bool {
    // SpamAssassin typically uses 0-10 scale
    // Sieve spamtest uses 0-10 scale
    score := int(ctx.SpamScore)
    if score > 10 {
        score = 10
    }
    if score < 0 {
        score = 0
    }

    testValue, _ := strconv.Atoi(t.Value)

    switch t.MatchType {
    case "ge":
        return score >= testValue
    case "le":
        return score <= testValue
    case "gt":
        return score > testValue
    case "lt":
        return score < testValue
    case "eq":
        return score == testValue
    case "ne":
        return score != testValue
    }

    return false
}

// Example usage in Sieve script:
// require ["spamtest", "fileinto"];
// if spamtest :value "ge" "7" {
//     fileinto "Spam";
// }
```

---

## 6.2 ManageSieve Protocol [FULL]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| MSV-001 | ManageSieve server (RFC 5804) | [ ] | SV-001 |
| MSV-002 | Script upload/download | [ ] | MSV-001 |
| MSV-003 | Script activation | [ ] | MSV-001 |
| MSV-004 | Script validation | [ ] | MSV-001 |

### MSV-001: ManageSieve Server

```go
// internal/managesieve/server.go
package managesieve

import (
    "bufio"
    "net"
)

type Server struct {
    listener     net.Listener
    scriptService *service.SieveScriptService
    userService   *service.UserService
    logger        logger.Logger
}

func NewServer(addr string, scriptService *service.SieveScriptService) (*Server, error) {
    listener, err := net.Listen("tcp", addr)
    if err != nil {
        return nil, err
    }

    return &Server{
        listener:      listener,
        scriptService: scriptService,
    }, nil
}

func (s *Server) Start() {
    for {
        conn, err := s.listener.Accept()
        if err != nil {
            continue
        }
        go s.handleConnection(conn)
    }
}

func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()

    session := &Session{
        conn:          conn,
        reader:        bufio.NewReader(conn),
        scriptService: s.scriptService,
        userService:   s.userService,
    }

    // Send capabilities
    session.sendCapabilities()

    // Command loop
    for {
        line, err := session.reader.ReadString('\n')
        if err != nil {
            return
        }

        cmd, args := parseCommand(line)
        session.handleCommand(cmd, args)
    }
}

type Session struct {
    conn          net.Conn
    reader        *bufio.Reader
    user          *domain.User
    authenticated bool
    scriptService *service.SieveScriptService
    userService   *service.UserService
}

func (s *Session) sendCapabilities() {
    capabilities := []string{
        `"IMPLEMENTATION" "gomailserver"`,
        `"SASL" "PLAIN"`,
        `"SIEVE" "fileinto reject envelope vacation variables relational subaddress spamtest"`,
        `"STARTTLS"`,
        `"VERSION" "1.0"`,
    }

    for _, cap := range capabilities {
        fmt.Fprintf(s.conn, "%s\r\n", cap)
    }
    fmt.Fprintf(s.conn, "OK\r\n")
}

func (s *Session) handleCommand(cmd string, args []string) {
    switch strings.ToUpper(cmd) {
    case "CAPABILITY":
        s.sendCapabilities()
    case "AUTHENTICATE":
        s.handleAuthenticate(args)
    case "LISTSCRIPTS":
        s.handleListScripts()
    case "GETSCRIPT":
        s.handleGetScript(args)
    case "PUTSCRIPT":
        s.handlePutScript(args)
    case "DELETESCRIPT":
        s.handleDeleteScript(args)
    case "SETACTIVE":
        s.handleSetActive(args)
    case "CHECKSCRIPT":
        s.handleCheckScript(args)
    case "HAVESPACE":
        s.handleHaveSpace(args)
    case "LOGOUT":
        fmt.Fprintf(s.conn, "OK \"Goodbye\"\r\n")
        s.conn.Close()
    default:
        fmt.Fprintf(s.conn, "NO \"Unknown command\"\r\n")
    }
}
```

### MSV-002: Script Upload/Download

```go
// internal/managesieve/commands.go
package managesieve

func (s *Session) handleGetScript(args []string) {
    if !s.authenticated {
        fmt.Fprintf(s.conn, "NO \"Not authenticated\"\r\n")
        return
    }

    if len(args) < 1 {
        fmt.Fprintf(s.conn, "NO \"Missing script name\"\r\n")
        return
    }

    scriptName := unquote(args[0])

    script, err := s.scriptService.Get(s.user.ID, scriptName)
    if err != nil {
        fmt.Fprintf(s.conn, "NO \"Script not found\"\r\n")
        return
    }

    // Send script content with byte count
    content := script.Content
    fmt.Fprintf(s.conn, "{%d}\r\n%s\r\nOK\r\n", len(content), content)
}

func (s *Session) handlePutScript(args []string) {
    if !s.authenticated {
        fmt.Fprintf(s.conn, "NO \"Not authenticated\"\r\n")
        return
    }

    if len(args) < 2 {
        fmt.Fprintf(s.conn, "NO \"Missing arguments\"\r\n")
        return
    }

    scriptName := unquote(args[0])

    // Parse literal length
    length := parseLiteralLength(args[1])

    // Read script content
    content := make([]byte, length)
    s.reader.Read(content)

    // Validate script
    parser := sieve.NewParser()
    _, err := parser.Parse(string(content))
    if err != nil {
        fmt.Fprintf(s.conn, "NO \"Script error: %s\"\r\n", err)
        return
    }

    // Save script
    script := &domain.SieveScript{
        UserID:  s.user.ID,
        Name:    scriptName,
        Content: string(content),
    }

    if err := s.scriptService.Save(script); err != nil {
        fmt.Fprintf(s.conn, "NO \"Failed to save script\"\r\n")
        return
    }

    fmt.Fprintf(s.conn, "OK\r\n")
}

func (s *Session) handleListScripts() {
    if !s.authenticated {
        fmt.Fprintf(s.conn, "NO \"Not authenticated\"\r\n")
        return
    }

    scripts, err := s.scriptService.List(s.user.ID)
    if err != nil {
        fmt.Fprintf(s.conn, "NO \"Failed to list scripts\"\r\n")
        return
    }

    for _, script := range scripts {
        if script.Active {
            fmt.Fprintf(s.conn, "\"%s\" ACTIVE\r\n", script.Name)
        } else {
            fmt.Fprintf(s.conn, "\"%s\"\r\n", script.Name)
        }
    }

    fmt.Fprintf(s.conn, "OK\r\n")
}
```

### MSV-003: Script Activation

```go
func (s *Session) handleSetActive(args []string) {
    if !s.authenticated {
        fmt.Fprintf(s.conn, "NO \"Not authenticated\"\r\n")
        return
    }

    if len(args) < 1 {
        fmt.Fprintf(s.conn, "NO \"Missing script name\"\r\n")
        return
    }

    scriptName := unquote(args[0])

    // Empty string deactivates all scripts
    if scriptName == "" {
        if err := s.scriptService.DeactivateAll(s.user.ID); err != nil {
            fmt.Fprintf(s.conn, "NO \"Failed to deactivate scripts\"\r\n")
            return
        }
        fmt.Fprintf(s.conn, "OK\r\n")
        return
    }

    if err := s.scriptService.SetActive(s.user.ID, scriptName); err != nil {
        if err == service.ErrScriptNotFound {
            fmt.Fprintf(s.conn, "NO \"Script not found\"\r\n")
        } else {
            fmt.Fprintf(s.conn, "NO \"Failed to activate script\"\r\n")
        }
        return
    }

    fmt.Fprintf(s.conn, "OK\r\n")
}
```

### MSV-004: Script Validation

```go
func (s *Session) handleCheckScript(args []string) {
    if !s.authenticated {
        fmt.Fprintf(s.conn, "NO \"Not authenticated\"\r\n")
        return
    }

    // Parse literal length
    length := parseLiteralLength(args[0])

    // Read script content
    content := make([]byte, length)
    s.reader.Read(content)

    // Validate script
    parser := sieve.NewParser()
    script, err := parser.Parse(string(content))
    if err != nil {
        fmt.Fprintf(s.conn, "NO \"Script error: %s\"\r\n", err)
        return
    }

    // Check required extensions
    for _, req := range script.Requires {
        if !isSupported(req) {
            fmt.Fprintf(s.conn, "NO \"Unsupported extension: %s\"\r\n", req)
            return
        }
    }

    fmt.Fprintf(s.conn, "OK\r\n")
}

func isSupported(extension string) bool {
    supported := map[string]bool{
        "fileinto":    true,
        "reject":      true,
        "envelope":    true,
        "vacation":    true,
        "variables":   true,
        "relational":  true,
        "subaddress":  true,
        "spamtest":    true,
        "imap4flags":  true,
        "copy":        true,
    }
    return supported[extension]
}
```

---

## 6.3 Sieve UI [FULL]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| SUI-001 | Visual rule editor in portal | [ ] | UP-001, SV-001 |
| SUI-002 | Common filter templates | [ ] | SUI-001 |
| SUI-003 | Raw script editor | [ ] | SUI-001 |
| SUI-004 | Rule testing interface | [ ] | SUI-001, SV-001 |

### SUI-001: Visual Rule Editor Component

```vue
<!-- web/portal/src/components/sieve/RuleEditor.vue -->
<template>
  <div class="rule-editor">
    <div class="rules-list">
      <draggable v-model="rules" handle=".drag-handle">
        <template #item="{ element, index }">
          <div class="rule-card">
            <div class="drag-handle">⋮⋮</div>
            <div class="rule-content">
              <div class="rule-header">
                <input v-model="element.name" placeholder="Rule name" />
                <label>
                  <input type="checkbox" v-model="element.enabled" />
                  Enabled
                </label>
              </div>

              <div class="conditions">
                <h4>If...</h4>
                <div v-for="(cond, i) in element.conditions" :key="i" class="condition">
                  <select v-model="cond.field">
                    <option value="from">From</option>
                    <option value="to">To</option>
                    <option value="subject">Subject</option>
                    <option value="header">Header</option>
                    <option value="size">Size</option>
                    <option value="spam">Spam Score</option>
                  </select>

                  <select v-model="cond.operator">
                    <option value="contains">contains</option>
                    <option value="is">is exactly</option>
                    <option value="matches">matches</option>
                    <option value="gt">greater than</option>
                    <option value="lt">less than</option>
                  </select>

                  <input v-model="cond.value" placeholder="Value" />

                  <button @click="removeCondition(index, i)" class="btn-sm">×</button>
                </div>

                <button @click="addCondition(index)" class="btn-link">+ Add condition</button>
              </div>

              <div class="actions">
                <h4>Then...</h4>
                <div v-for="(action, i) in element.actions" :key="i" class="action">
                  <select v-model="action.type">
                    <option value="fileinto">Move to folder</option>
                    <option value="flag">Add flag</option>
                    <option value="discard">Delete</option>
                    <option value="reject">Reject with message</option>
                    <option value="redirect">Forward to</option>
                    <option value="vacation">Auto-reply</option>
                  </select>

                  <template v-if="action.type === 'fileinto'">
                    <select v-model="action.mailbox">
                      <option v-for="mb in mailboxes" :key="mb.name" :value="mb.name">
                        {{ mb.name }}
                      </option>
                    </select>
                  </template>

                  <template v-if="action.type === 'flag'">
                    <select v-model="action.flag">
                      <option value="\\Flagged">Star</option>
                      <option value="\\Seen">Mark as read</option>
                    </select>
                  </template>

                  <template v-if="action.type === 'redirect'">
                    <input v-model="action.address" placeholder="email@example.com" />
                  </template>

                  <button @click="removeAction(index, i)" class="btn-sm">×</button>
                </div>

                <button @click="addAction(index)" class="btn-link">+ Add action</button>
              </div>
            </div>

            <button @click="removeRule(index)" class="btn-danger">Delete</button>
          </div>
        </template>
      </draggable>
    </div>

    <button @click="addRule" class="btn-primary">Add New Rule</button>

    <div class="editor-actions">
      <button @click="saveRules" class="btn-primary">Save Rules</button>
      <button @click="showRawEditor = true" class="btn-secondary">Edit Raw Script</button>
    </div>

    <modal v-if="showRawEditor" @close="showRawEditor = false">
      <raw-script-editor
        :script="generateScript()"
        @save="importScript"
      />
    </modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import draggable from 'vuedraggable'
import { saveScript, getScript } from '@/api/sieve'

const rules = ref<Rule[]>([])
const mailboxes = ref<Mailbox[]>([])
const showRawEditor = ref(false)

function generateScript(): string {
  const requires = new Set<string>()
  let script = ''

  for (const rule of rules.value) {
    if (!rule.enabled) continue

    // Collect requirements
    rule.actions.forEach(a => {
      if (a.type === 'fileinto') requires.add('fileinto')
      if (a.type === 'reject') requires.add('reject')
      if (a.type === 'vacation') requires.add('vacation')
      if (a.type === 'flag') requires.add('imap4flags')
    })

    rule.conditions.forEach(c => {
      if (c.field === 'spam') requires.add('spamtest')
    })
  }

  // Generate require statement
  if (requires.size > 0) {
    script += `require [${[...requires].map(r => `"${r}"`).join(', ')}];\n\n`
  }

  // Generate rules
  for (const rule of rules.value) {
    if (!rule.enabled) continue

    script += `# ${rule.name}\n`
    script += generateConditions(rule.conditions)
    script += ` {\n`
    script += generateActions(rule.actions)
    script += `}\n\n`
  }

  return script
}

async function saveRules() {
  const script = generateScript()
  await saveScript('main', script)
}
</script>
```

### SUI-002: Common Templates

```typescript
// web/portal/src/data/sieveTemplates.ts
export const templates = [
  {
    name: 'Move mailing lists to folder',
    description: 'Move emails from mailing lists to a dedicated folder',
    rule: {
      conditions: [
        { field: 'header', header: 'List-Id', operator: 'exists', value: '' }
      ],
      actions: [
        { type: 'fileinto', mailbox: 'Lists' }
      ]
    }
  },
  {
    name: 'Star emails from boss',
    description: 'Flag important emails from your manager',
    rule: {
      conditions: [
        { field: 'from', operator: 'contains', value: 'boss@' }
      ],
      actions: [
        { type: 'flag', flag: '\\Flagged' }
      ]
    }
  },
  {
    name: 'Move spam to Junk',
    description: 'Automatically file high-scoring spam',
    rule: {
      conditions: [
        { field: 'spam', operator: 'gt', value: '5' }
      ],
      actions: [
        { type: 'fileinto', mailbox: 'Spam' }
      ]
    }
  },
  {
    name: 'Vacation auto-reply',
    description: 'Send automatic replies when on vacation',
    rule: {
      conditions: [],
      actions: [
        {
          type: 'vacation',
          days: 7,
          subject: 'Out of Office',
          message: 'I am currently out of the office and will respond when I return.'
        }
      ]
    }
  },
  {
    name: 'Forward work emails',
    description: 'Forward certain emails to another address',
    rule: {
      conditions: [
        { field: 'to', operator: 'contains', value: 'work@' }
      ],
      actions: [
        { type: 'redirect', address: 'personal@example.com', copy: true }
      ]
    }
  }
]
```

### SUI-004: Rule Testing Interface

```vue
<!-- web/portal/src/components/sieve/RuleTester.vue -->
<template>
  <div class="rule-tester">
    <h3>Test Your Rules</h3>

    <div class="test-input">
      <h4>Test Message</h4>
      <div class="form-group">
        <label>From:</label>
        <input v-model="testMessage.from" placeholder="sender@example.com" />
      </div>
      <div class="form-group">
        <label>To:</label>
        <input v-model="testMessage.to" placeholder="you@yourdomain.com" />
      </div>
      <div class="form-group">
        <label>Subject:</label>
        <input v-model="testMessage.subject" placeholder="Test subject" />
      </div>
      <div class="form-group">
        <label>Headers (JSON):</label>
        <textarea v-model="testMessage.headers" placeholder='{"List-Id": "test"}' />
      </div>
      <div class="form-group">
        <label>Spam Score (0-10):</label>
        <input v-model.number="testMessage.spamScore" type="number" min="0" max="10" />
      </div>
    </div>

    <button @click="runTest" class="btn-primary">Test Rules</button>

    <div v-if="testResult" class="test-result" :class="{ success: testResult.success }">
      <h4>Result</h4>
      <div v-if="testResult.matchedRules.length">
        <p>Matched Rules:</p>
        <ul>
          <li v-for="rule in testResult.matchedRules" :key="rule">{{ rule }}</li>
        </ul>
      </div>
      <div v-else>
        <p>No rules matched. Message would be kept in INBOX.</p>
      </div>

      <div v-if="testResult.actions.length">
        <p>Actions that would be taken:</p>
        <ul>
          <li v-for="action in testResult.actions" :key="action">{{ action }}</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { testSieveScript } from '@/api/sieve'

const testMessage = ref({
  from: '',
  to: '',
  subject: '',
  headers: '{}',
  spamScore: 0
})

const testResult = ref<TestResult | null>(null)

async function runTest() {
  testResult.value = await testSieveScript(testMessage.value)
}
</script>
```

---

## Acceptance Criteria

### Sieve Interpreter
- [ ] Base RFC 5228 commands work (if, stop, keep, discard, fileinto, redirect, reject)
- [ ] Variables extension works
- [ ] Vacation extension works
- [ ] Relational tests work
- [ ] Subaddress extension works
- [ ] Spamtest extension integrates with SpamAssassin

### ManageSieve Protocol
- [ ] Clients can connect and authenticate
- [ ] Scripts can be uploaded/downloaded
- [ ] Scripts can be activated/deactivated
- [ ] Script validation works

### UI
- [ ] Visual rule editor creates valid Sieve scripts
- [ ] Templates can be applied
- [ ] Raw script editing works
- [ ] Rule testing provides accurate results

---

## Next Phase

After completing Phase 6, proceed to [TASKS7.md](TASKS7.md) - Webmail Client.
