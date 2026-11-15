# Two stages: Hugo builds the site, then nginx serves it. Dependabot keeps the pinned image digests current.
FROM hugomods/hugo:exts-0.148.2@sha256:9af903e2a217b51c80681247a740edf16d51bb2938446ad5013521fedb936693 AS builder
ARG HUGO_BASEURL=http://localhost:8080
ENV HUGO_BASEURL=${HUGO_BASEURL}
ENV HUGO_CACHEDIR=/hugo_cache
WORKDIR /src

# Install the npm packages first, so Docker can reuse this step when only the content changes.
# --engine-strict=false: the npm in this image is older than package.json asks for. npm ci only
# installs what the lockfile says, so the settings that need a newer npm don't matter here.
COPY package.json package-lock.json .npmrc ./
RUN npm ci --ignore-scripts --engine-strict=false

# The cache mounts keep the resized images and the Hugo modules between builds,
# so a deploy only processes new or changed images.
COPY . .
RUN --mount=type=cache,target=/src/resources/_gen \
    --mount=type=cache,target=/hugo_cache \
    npm run build

# nginx.conf comes from the hugomods/hugo:nginx image this site used before, with cache headers added.
# The official alpine-slim image is about 6 MB; hugomods/hugo:nginx was about 60 MB.
FROM nginx:1.31-alpine-slim@sha256:4c734d33ee731135dc9f2349914a34551c920f4451fca8bb05bf5880eb9079e8
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /src/public /site
