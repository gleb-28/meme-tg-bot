# TODO

## How we work with this list
- Keep items small and outcome-focused; define what “done” means.
- Priority:
    - **P1** — critical
    - **P2** — soon
    - **P3** — nice-to-have
- Use checkboxes:
    - `[ ]` — todo / in progress
    - `[x]` — done

## Backlog

### P1 — critical

- [ ] **Handle yt-dlp can't process Instagram images**
- [ ] **Add timeouts and cancellation for long-running operations**
- [ ] **Prevent forward-mode loops and self-triggering**
  _Notes_: Ignore messages coming from the forward chat itself, Ignore messages sent by the bot

### P2 — soon

- [ ] **Clear session cache for non-auth users**
- [ ] **Improve yt-dlp error handling and user feedback**
- [ ] **Add rate limiting and abuse protection**

### P3 — nice-to-have

- [ ] **Expand README with “forward mode” how-to**
  _Notes_: Describe `/change_mode`, required admin permissions, and that albums/media batches are forwarded
- [ ] **Add CI to run tests and formatting + pre-commit**
  _Notes_: Run `go test ./...` and `gofmt` on pushes/PRs (GitHub Actions)


## Done

- [x] **Add bot.MustEdit method**
