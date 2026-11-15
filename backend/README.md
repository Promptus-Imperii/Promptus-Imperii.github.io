# backend

Small Go API for the website. It handles the signup form and the Altcha captcha, and shows the contact address once the captcha is solved.
It uses the official [Altcha library](https://github.com/altcha-org/altcha-lib-go) (v2 protocol, PBKDF2/SHA-256, for widget 3.x) and [go-mail](https://github.com/wneessen/go-mail) for SMTP. Responses are HTML fragments that [htmx 4](https://four.htmx.org) swaps into the page. The fragments are in [templates/](templates).

## Endpoints

All of them work with or without the `/api` prefix.

| Method | Path | Body | Response |
|---|---|---|---|
| GET | `/api/captcha-challenge` | | Altcha challenge (JSON) |
| POST | `/api/signup` | signup form (form-encoded) | 204 with `HX-Redirect: /signup_success/`, or `#errors` with 422 (invalid input) or 500 (mail failed) |
| POST | `/api/email` | `altcha` | `#email-block` with the contact address, or `#errors` with 422 |

Abuse protection comes from Altcha and the Cloudflare proxy. POST bodies are capped at 64 KB.

## Configuration

See [.env.example](.env.example). In production, Coolify sets the variables.

### Mail

Production sends through a [Microsoft 365 SMTP relay connector](https://learn.microsoft.com/en-us/exchange/mail-flow-best-practices/how-to-set-up-a-multifunction-device-or-application-to-send-email-using-microsoft-365-or-office-365), so the backend has no mail credentials:

- `SMTP_HOST` is the tenant's MX host (`<domain>.mail.protection.outlook.com`). `SMTP_PORT` defaults to 25.
- `SERVER_EMAIL_ADDRESS` is the sender. It must be on one of the tenant's domains but doesn't need a mailbox.
- In the Exchange admin center, an inbound connector from "Your organization's email server" trusts the server's static IP.
- The domain's SPF record includes that IP.
- Hetzner blocks outgoing port 25 until you request an unblock (Console > Limits).

If the server's IP changes, update the connector and the SPF record, or mail stops.

## Development

### Go version

The backend needs Go 1.27.4 or newer (see [go.mod](go.mod)). Go can download that version itself. Run this once:

```sh
go env -w GOTOOLCHAIN=auto
```

- Any Go from 1.21 on works as the starting point, so the Go from your package manager is fine.
- In this project, `go` then downloads the version from `go.mod`, checks it against Go's checksum database and uses it. The download goes to the module cache, and your system Go stays as it is.
- Other projects keep using your system Go, unless their `go.mod` asks for something newer.
- Fedora and some other Linux distributions need this, because their Go packages set `GOTOOLCHAIN=local`, which turns the download off. Without it, `go` stops with `go.mod requires go >= 1.27.4 (running go 1.26.x; GOTOOLCHAIN=local)`.
- To undo it: `go env -u GOTOOLCHAIN`.

### Running locally

```sh
# terminal 1: the site
cd site && npm run dev

# terminal 2: the backend, which serves /api/ and passes everything else on to Hugo
cd backend && DEV=1 go run .
```

Open http://localhost:3000. With `DEV=1`, mails are logged instead of sent, the contact address is `info@example.org`, and the captcha key is random. No other settings are needed. To try other values, set them as environment variables; [.env.example](.env.example) lists them all.

## Testing and checks

Run these from `backend/` before pushing. CI runs the same three.

```sh
go test -race ./...
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0 run
go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...
```

- `go test -race` runs the tests with the race detector, which needs a C compiler (`gcc`).
- [golangci-lint](https://golangci-lint.run) runs the linters in [.golangci.yml](.golangci.yml): the standard set (`go vet`, staticcheck, errcheck and more) plus stricter ones such as gosec, revive and gocritic.
  Fix formatting with `golangci-lint fmt` (gofumpt and goimports).
- [govulncheck](https://go.dev/doc/security/vuln/) reports known vulnerabilities in dependencies that the code calls.

`go run tool@version` downloads and builds each tool the first time, which takes about a minute. You don't need to install anything.

## Production

Coolify runs the backend with the Dockerfile build pack, base directory `/backend` and domain `https://svpromptusimperii.nl/api`.
The [Dockerfile](Dockerfile) builds a static binary on `golang:1.27-alpine` and runs it on `distroless/static`, which includes CA certificates and runs as a non-root user.
We use a Dockerfile because Nixpacks only supports Go up to 1.23.
