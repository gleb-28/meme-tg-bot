# Prompt handling instructions

- ## 1. Architecture map
- `cmd/bot/main.go:1-44` wires up configuration, logger, database, FSM, video and forward services, and registers command/message/keyboard handlers before starting the bot.
- `internal/core/config/config.go:27-79` uses `cleanenv` with struct tags, reads `.env` only when `IS_DOCKERIZED` is not `true`, and always validates env vars via `cleanenv.ReadEnv`, ensuring both local and container deployments succeed with minimal custom logic.
- SQLite connection and migrations live in `internal/db/db.go:1-46`, models reside in `models/chat.go:7-20`. Direct any storage-related questions (user or forward-mode persistence) to these files.
- Commands (`/start`, `/key`, `/change_mode`) and authentication live in `internal/handlers/commands` (see `commands.go:1-25`, `key.go:11-25`, `changeMode.go:1-16`), while FSM states (`StateInitial`, `StateAwaitingKey`, `StateAwaitingForwardChat`) are defined in `internal/fsm/fsm.go:1-78`.
- Text and link handling runs through `internal/handlers/message/message.go:1-34`, with specific handlers in `handleLink.go:1-73`, `handleMessage.go:1-17`, `validateActivationKey.go:1-31`, and `validateForwardChat.go:1-52`. The FSM decides when to route to each handler.
- Forward-mode logic lives across keyboard handlers (`internal/handlers/keyboard/keyboard.go:1-32`, `forwardMode.go:14-140`), the service in `internal/feat/forward/service.go:1-71`, repository `internal/repo/forwardMode.go:1-55`, and session store `internal/session/session.go:10-124`.
- Video download/compression is handled in `pkg/video/video.go:1-177`, with helper utilities in `pkg/utils/*.go`. Watch for directory creation, file renaming, salting, compression, and cleanup.
- User-facing replies are centralized in `internal/text/replies.go:1-24`. Update this file whenever text changes are needed.
- Repository commands live in `Makefile:1-31`, covering `run`, `test`, `tidy`, Docker `build/down/rebuild`, and `deploy`.

## 2. Checklist for prompt processing
1. **Clarify the scope:** determine if the prompt targets commands, messages, forward mode, DB, or video service. Refer to the relevant architecture section when you answer.
2. **Map to files:** command/button changes belong around `internal/handlers/commands` or `keyboard`. Forward-mode fixes touch `internal/feat/forward/service.go`, sessions, and keyboard handlers. Video-related questions point to `pkg/video/video.go` and `pkg/utils`.
3. **Respect FSM:** new flows must trigger `AwaitingKeyEvent` or `AwaitingForwardChatEvent` inside `internal/fsm/fsm.go:1-78` so user states stay coherent.
4. **Guard persistence:** database updates go through `internal/repo/*` and the models in `models/chat.go:7-20`.
5. **Adjust replies when needed:** change the text in `internal/text/replies.go:1-24` and update keyboards in `internal/handlers/keyboard/forwardMode.go:14-140`.

## 3. Answer guidance and conventions
- Point out which files are affected and why (e.g., “forward-mode toggling happens in `internal/handlers/keyboard/forwardMode.go:14-140` plus `internal/feat/forward/service.go:1-71`”).
- Describe side effects: mention FSM events, session state, DB updates, and reply text changes when proposing a solution.
- Ask clarifying questions when requirements are vague (“Which mode should be default?” or “Do we need a new button?”).
- Always remind about verification commands (`Makefile:21-31` – `go test ./... -v`, `go run cmd/bot/main.go`, `make tidy`).
- Keep code, comments, files, and commits in English to match the requested style. Responses can follow the user's language.
- Follow the Git commit style guide at https://github.com/slashsbin/styleguide-git-commit-message when crafting commits.
- Apply Go best practices from https://go.dev/ref/spec and https://go.dev/doc/effective_go in any new or modified Go code.
- Use the documented behaviors from the referenced packages when working in their domains:
  * GORM updates: https://gorm.io/docs/update.html
  * Telebot usage: https://pkg.go.dev/gopkg.in/telebot.v4
  * Looplab FSM: https://pkg.go.dev/github.com/looplab/fsm

## 4. Additional reminders
- Check how tests or scripts are expected to run (e.g., `go test ./... -v`) before editing.
- Always suggest a follow-up step (“verify migrations”, “run `make tidy` and tests”).
- Highlight env/config changes with references to `internal/core/config/config.go:27-80`.
