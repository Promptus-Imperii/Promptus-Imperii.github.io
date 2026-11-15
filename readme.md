# Promptus Imperii

Website and backend of the Promptus Imperii study association.

## Layout

- `site/`: Hugo + Tailwind website. See [site/readme.md](site/readme.md).
- `backend/`: Go API for signup and email. See [backend/README.md](backend/README.md).

Both are deployed separately on Coolify (base directories `/site` and `/backend`).

## Development

Run the site and the backend in two terminals, then open http://localhost:3000:

```sh
# terminal 1: the site (Hugo)
cd site && npm ci && npm run dev

# terminal 2: the backend, which serves /api/ and passes everything else on to Hugo
cd backend && DEV=1 go run .
```

With `DEV=1` the backend needs no other settings and logs mails instead of sending them. On Windows (PowerShell), start it with `$env:DEV=1; go run .`.

The backend needs Go 1.27. Run `go env -w GOTOOLCHAIN=auto` once and Go downloads it itself; see [backend/README.md](backend/README.md). The site's requirements are in [site/readme.md](site/readme.md).
