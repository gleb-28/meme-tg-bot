# Prompt handling instructions

## 1. Architecture map
- `cmd/bot/main.go:1-44` wires up configuration, logger, database, FSM, session store, services, and registers command/message/keyboard handlers before starting the bot.
- `internal/core/config/config.go:10-39` uses `cleanenv` with struct tags, reads `.env` only when `IS_DOCKERIZED` is not `true`, and always validates env vars via `cleanenv.ReadEnv`.
- SQLite connection and migrations live in `internal/db/db.go:1-47`, models live in `models/chat.go:7-20`, and DB access is encapsulated in `internal/repo/chat.go:9-29` and `internal/repo/forwardMode.go:9-52`.
- Commands are registered in `internal/handlers/commands/commands.go:1-25` with handlers in `start.go:9-13`, `key.go:11-26`, and `changeMode.go:10-16`; authentication middleware is in `internal/handlers/auth/auth.go:9-24`.
- FSM states/events (`StateInitial`, `StateAwaitingKey`, `StateAwaitingForwardChat`) are defined in `internal/fsm/fsm.go:12-78`.
- Message routing is in `internal/handlers/message/message.go:11-40`, URL downloads in `handleLink.go:16-73`, generic message handling in `handleMessage.go:10-25`, activation key validation in `validateActivationKey.go:12-31`, and forward-chat validation in `validateForwardChat.go:13-51`.
- Forward-mode UI and callbacks live in `internal/handlers/keyboard/keyboard.go:12-31` and `internal/handlers/keyboard/forwardMode.go:14-140`, with callback constants in `internal/core/constants/constants.go:3-10`.
- Forward-mode behavior is implemented in `internal/feat/forward/service.go:10-71` with session cache and batching in `internal/session/session.go:11-153`.
- Message forwarding (including albums) lives in `internal/bot.go:39-251`.
- Video download/compression is handled in `pkg/video/video.go:25-177`, with helper utilities in `pkg/utils/*.go`.
- User-facing replies are centralized in `internal/text/replies.go:3-62`.
- Repository commands live in `Makefile:1-32`.

## 2. Checklist for prompt processing
1. Clarify the scope: determine if the prompt targets commands, messages, forward mode, DB, or video service.
2. Map to files: command/button changes belong around `internal/handlers/commands` or `internal/handlers/keyboard`. Forward-mode fixes touch `internal/feat/forward/service.go`, sessions, and keyboard handlers. Video-related questions point to `pkg/video/video.go` and `pkg/utils`.
3. Respect FSM: new flows must trigger `AwaitingKeyEvent` or `AwaitingForwardChatEvent` inside `internal/fsm/fsm.go:12-78` so user states stay coherent.
4. Guard persistence: database updates go through `internal/repo/*` and the models in `models/chat.go:7-20`.
5. Adjust replies when needed: change the text in `internal/text/replies.go:3-62` and update keyboards in `internal/handlers/keyboard/forwardMode.go:14-140` (plus `internal/core/constants/constants.go:3-10` if you add callbacks).

## 3. Answer guidance and conventions
- Point out which files are affected and why (for example: forward-mode toggling happens in `internal/handlers/keyboard/forwardMode.go:14-140` plus `internal/feat/forward/service.go:10-71`).
- Describe side effects: mention FSM events, session state, DB updates, and reply text changes when proposing a solution.
- Ask clarifying questions when requirements are vague (for example: "Which mode should be default?" or "Do we need a new button?").
- Always remind about verification commands (`Makefile:23-32` - `go test ./... -v`, `go run cmd/bot/main.go`, `make tidy`).
- Keep code, comments, files, and commits in English to match the requested style. Responses can follow the user's language.
- Follow the Git commit style guide at https://github.com/slashsbin/styleguide-git-commit-message when crafting commits.
- Apply Go best practices from https://go.dev/ref/spec and https://go.dev/doc/effective_go in any new or modified Go code.
- Use the documented behaviors from the referenced packages when working in their domains:
  * GORM updates: https://gorm.io/docs/update.html
  * Telebot usage: https://pkg.go.dev/gopkg.in/telebot.v4
  * Looplab FSM: https://pkg.go.dev/github.com/looplab/fsm

## 4. Additional reminders
- Check how tests or scripts are expected to run (for example, `go test ./... -v`) before editing.
- Always suggest a follow-up step ("verify migrations", "run `make tidy` and tests").
- Highlight env/config changes with references to `internal/core/config/config.go:10-39`.
- Deployment helpers live under `deploy/` (script, compose manifest, Dockerfile, and env templates) and the Makefile docker targets point at `deploy/docker-compose.yml`.
