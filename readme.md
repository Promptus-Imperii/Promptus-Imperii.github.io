# Promptus Imperii

Website and backend of the Promptus Imperii study association.

## Layout

- `site/`: Hugo + Tailwind website. See [site/readme.md](site/readme.md).
- `backend/`: Go API for signup and email. See [backend/README.md](backend/README.md).

Both are deployed separately on Coolify (base directories `/site` and `/backend`).

## Development

```sh
# site
cd site && npm install && npm run dev

# backend
cd backend && cp .env.example .env && go run .
```
