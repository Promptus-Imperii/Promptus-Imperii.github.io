# Promptus Imperii
The website of the Promptus Imperii study association.

## Technologies
- [Hugo](https://gohugo.io/), which uses Go templates.
- A heavily edited version of the [Hugoplate](https://github.com/zeon-studio/hugoplate) theme:
    - It uses a few Hugo modules (see go.mod), mostly from [hugo-modules](https://github.com/TeaDrinkingProgrammer/hugo-modules). The gallery comes from a fork because of [this pull request](https://github.com/gethugothemes/hugo-modules/pull/56).
    - The pages are adapted to the study association, all tracking is removed, and there is a signup page that talks to the [backend](../backend/README.md).
- [Tailwind CSS](https://tailwindcss.com/), built with npm.
- [htmx](https://htmx.org/) for the signup and contact forms.
- Hosting on [Coolify](https://coolify.io/), with the [Dockerfile](Dockerfile): Hugo builds the site and nginx serves it with [nginx.conf](nginx.conf).

## Set-up
- Install Hugo extended, version 0.148 (the Dockerfile uses 0.148.2): https://gohugo.io/installation/. On Windows, [winget](https://gohugo.io/installation/windows/#winget) is the easiest.
- Install [Go](https://go.dev/dl/). Hugo needs it to download the modules.
- Install npm 11.10 or newer. The npm that comes with the latest Node LTS is new enough; `npm install -g npm@latest` updates an older one. Older versions refuse to run here, because they don't know the supply-chain settings in [.npmrc](.npmrc). Use npm, not Bun or Yarn, for the same reason.
  - Without npm, install it with [nvm](https://github.com/nvm-sh/nvm) (Linux and macOS) or [nvm-windows](https://github.com/coreybutler/nvm-windows).
  - Then run `nvm install --lts` and `nvm use --lts` (Linux and macOS), or `nvm install lts` and `nvm use <version-number>` (Windows).
  - On Windows, run `Set-ExecutionPolicy RemoteSigned -Scope CurrentUser` in PowerShell, so it can run PowerShell scripts.
- Run `npm ci` to install the npm packages. It installs exactly what `package-lock.json` says.

## Troubleshooting

### `npm run dev` fails with "tailwindcss is not whitelisted" error
Hugo requires `tailwindcss` to be explicitly allowed in the security policy. Ensure your `hugo.toml` includes:
```toml
[security.exec]
allow = ['^(dart-)?sass$', '^go$', '^git$', '^node$', '^postcss$', '^tailwindcss$']
```

### Build errors about deprecated `_build` front matter
Hugo v0.145.0+ removed support for the `_build` front matter key. It has been replaced with `build`. Content files should use:
```yaml
---
build:
  render: "never"
---
```
instead of `_build:`.

## Running locally
Start the site and the backend in two terminals, then open http://localhost:3000:

```sh
# terminal 1, in site/
npm run dev

# terminal 2, in backend/
DEV=1 go run .
```

The backend serves `/api/` and passes everything else on to Hugo, so the forms work the same way as in production. With `DEV=1` it logs mails instead of sending them. On Windows (PowerShell), start it with `$env:DEV=1; go run .`. The backend needs Go 1.27; see [backend/README.md](../backend/README.md).

Only changing content? `hugo server` on its own also works: open http://localhost:1313. The forms won't work there.

To build the site, run `npm run build`. The first build takes a while, because Hugo resizes all images.

Tip: http://www.randomiban.com/?country=Netherlands gives IBANs to test the signup form with.

## npm supply-chain hardening
Attacks like the Shai-Hulud worm publish malicious versions of popular packages. Installing one runs its code, which steals tokens and uses them to spread. [.npmrc](.npmrc) limits what npm does:
- `ignore-scripts=true`: npm never runs the install scripts (`preinstall`, `install`, `postinstall`) of dependencies. `npm run <script>` still works. All current dependencies work without them; `sharp` and `@parcel/watcher` ship prebuilt binaries.
- `min-release-age=7`: `npm install` and `npm update` only pick versions that are at least 7 days old. Malicious versions are usually removed within that time.
- `allow-git=none`: dependencies only come from the npm registry, never from a git repository.
- `engine-strict=true` with `"engines": {"npm": ">=11.10.0"}` in package.json: older npm versions stop with `EBADENGINE`, because they would silently ignore the settings above.

The Docker build and CI use `npm ci`, which fails if `package-lock.json` doesn't match `package.json` and checks every package against the hash in the lockfile. CI also runs `npm audit signatures`, which checks that the registry signed every package. Dependabot waits 7 days before it proposes a new version.

If a new dependency really needs its install script, look for a different package first. Otherwise read what the script does and run it once by hand with `npm rebuild <package> --ignore-scripts=false`.

## Environment variables
When deploying, set `HUGO_BASEURL` to the full URL of the site, such as `https://example.com`. In Coolify this is under Configuration > Environment Variables; tick `Build Variable`.

## Maintenance
- Site: run `npm update` and check the website. TODO: switch to a masonry layout once that is stable.
- Backend: in `backend/`, run `go get -u ./... && go mod tidy`, then the checks from [backend/README.md](../backend/README.md), and submit a form on the site.
- Now and then, check that the dependencies are still maintained, and switch where needed (especially for the backend).

## Adding new activity posts
- An activity lives in two folders: [content/nl/activities](content/nl/activities) for the text and [assets/images/activities](assets/images/activities) for the photos.
- To create a post, copy an existing file in the content folder or use the template below.
  - Folder and file names may not contain capital letters.
  - Name it `yyyy-mm-dd-my-activity`. The date prevents duplicates and sorts the files.
- For the photos, create a folder in the assets folder with the same name as the content file, and put all of the activity's photos in it.
  - To set the order in the gallery, start the file names with `N__`, where N is the position.
- After adding photos, run `npm run process-images` before committing. The script:
  - scales photos down to fit within 1920×1080 (landscape) or 1080×1920 (portrait), to keep the repository small. It keeps the original if the scaled-down file would be larger;
  - turns photos the right way up, using their EXIF orientation, because of these bugs: https://discourse.gohugo.io/t/why-is-hugo-rotating-my-image/23329/5 and https://github.com/golang/go/issues/4341.
  - Scaling down by hand? Don't compress the photos: the website does that already.
  - Is a photo still sideways after the script? Turn it with any photo editor (Windows Photos, Paint.NET and so on, but not the old Windows Photo Viewer, which causes problems) until it is the right way up and save it.
- When you're done, start the website and check that:
  - the gallery works;
  - the photos aren't sideways or cut off;
  - the category and tag overviews are correct and have no duplicates;
  - the dates are correct;
  - the text matches the activity;
  - the links work.
- After big changes to the template, update it in every activity, for example with search and replace.

### Activity template
```
---
# Title shown at the top of the page.
title: "My Activity"
# Not directly shown to users, but used by search engines to index the website for example.
description: "A description of my activity."
# Date must be in the form yyyy-MM-ddTHH:mm:ss ('T' separates the date and time), a time of 00:00:00 will only display the date
date: 2023-09-01T00:00:00
# This path is relative to the assets folder. Create a folder that is assets/images/activities/file-name
# IMPORTANT! The file name/folder name may not contain capital letters!
image_folder: "images/activities/yyyy-mm-dd-my-activity"
# Use the original file name and prepend N__, where N is the order in which you want the image to show in the gallery. Make sure to downscale 4k images to 1080p to keep the size of the git repo manageable. Only use the .jp(e)g, .png and .gif file formats.
image: "1__20230901_151607.jpg"
# Don't forget capital letters in the categories and tags. Also check all categories and tags by loading the activities page and looking at the list.
categories: ["Lezing"]
# Tags are more specific than categories, and are currently only used for sponsor/partner names
tags: ["MyCompany"]
# If this is set to true, the content is not published
draft: false
---

Write about your activity here!

Another line with my conclusion.
```