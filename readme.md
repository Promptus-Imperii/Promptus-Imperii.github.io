# Promptusimperii
The new website of the Promptus Imperii study association.


# Technologies
This website is built on the following libraries/frameworks are used in this project:
- Hugo (uses Golang Templates).
- A heavily edited version of [Hugoplate](https://github.com/zeon-studio/hugoplate), which uses:
    - A few components (see go.mod), mainly [hugo-modules](https://github.com/TeaDrinkingProgrammer/hugo-modules). A fork is used from the main version due to [this](https://github.com/gethugothemes/hugo-modules/pull/56).
    - The pages have been changed to suit the needs of the study association.
    - All tracking has been stripped out.
    - A signup page was added (see the backend project).
- Tailwind with a few plug-ins (see tailwind.config.js)
- Some build tools to build the [Tailwind CSS](https://tailwindcss.com/).
- The hosting uses [Coolify](https://coolify.io/).

# Set-up
- Install any NPM-like package manager: since it is only used for building Tailwind and launching scripts, you can probably use any package manager (NPM, Bun, Yarn, etc.). Just make sure to use the latest (LTS) version.
- To run the website locally, run `npm run dev` or equivalent.
- To build the website, run `npm run build` or equivalent.

# Developing with the backend
When testing the backend, start the backend by following the instructions on that repo. Change the following variables:
- `serverURL` in config/params.toml
- `URL` in assets/js/submit.js and assets/js/email.js

Tip: use `http://www.randomiban.com/?country=Netherlands` to test with IBAN numbers

**Make sure to not commit these changes!**


# Maintainance

## Frontend
- `npm update` and test the website. TODO: Migrate to masonary layout when it is stable

## Backend
- `go get -u -all` and test the endpoint by submitting a form with the frontend

## Both
Check every once in a while if the dependencies are still maintained and switch to update ones where possible (especially with the backend)


# Post instructies
- Om een post te maken, kopieer een bestaande post die lijkt op degene die je wilt maken. Voorbeeld: maak een post voor een bedrijfsbezoek, kopieer een ander bedrijfsbezoek.
- Vul de volgende velden in:
```yaml
# Title shown at the top of the page.
title: "Introductieborrel ATIx-studieverenigingen"
# Please add a good description and title (see the about us page). Google uses it to recommend the website
description: "Introductieborrel van alle studieverenigingen voor de nieuwe eerstejaarsstudenten van ATIx"
# Date must be in this form, set the time to 00:00:00 since it is not shown
date: 2023-09-01T00:00:00
# This path is relative to the assets folder. Create a folder that is assets/images/activities/file-name
# IMPORTANT! The file name/folder name may not contain capital letters!
image_folder: "images/activities/2023-introductieborrel-studieverenigingen"
# Use the original file name and prepend N__, where N is the order in which you want the image to show in the gallery. Make sure to downscale 4k images to 1080p to keep the size of the git repo manageable.
image: "1__20230901_151607.jpg"
# Don't forget capital letters in the categories and tags. Also check all categories and tags by loading the activities page and looking at the list.
categories: ["Ledenactiviteit"]
# Tags are more specific than categories
tags: ["Introductieborrel"]
# If this is set to true, the content is not pulished
draft: false
```
- Maak de map aan in image/activities en plaats de foto's. Pas de naam aan en zet er N__ voor, waar N de volgorde is waarin de foto's verschijnen in de gallerij. Downscale de foto's eerst naar 1080p als ze 4k zijn (gebruik https://imageresizer.com/), om de grootte van de Git repo laag te houden. Comprimeer ze hier niet! Dit doet de website zelf al.
