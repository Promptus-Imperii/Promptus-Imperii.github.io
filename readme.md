# Promptus Imperii
The new website of the Promptus Imperii study association.

## Technologies
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

## Set-up
- Install Hugo: https://gohugo.io/installation/. For Windows, use [winget](https://gohugo.io/installation/windows/#winget) for the easiest install.
- Install [GoLang](https://go.dev/dl/)
- Install any NPM-like package manager: since it is only used for building Tailwind and launching scripts, you can probably use any package manager (NPM, Bun, Yarn, etc.). Just make sure to use the latest (LTS) version.
  - If you don't have npm installed, I recommend [nvm](https://github.com/nvm-sh/nvm) for Linux/MacOs and [nvm-windows](https://github.com/coreybutler/nvm-windows) for Windows
  - After installing nvm(-windows), run `nvm install --lts` (MacOs/Linux) or `nvm install lts` (windows) and activate the downloaded version with `nvm use --lts` (MacOs/Linux) or `nvm use <version-number>` (Windows)
  - On windows, run `Set-ExecutionPolicy RemoteSigned -Scope CurrentUser` in Powershell to allow the running of PowerShell scripts
- To install the npm packages, run `npm install`
- To run the website locally, run `npm run dev` or equivalent. The first time can take a while.
- To build the website, run `npm run build` or equivalent.

## Developing with the backend
When testing the backend, start the backend by following the instructions on that repo. Change the following variables:
- `serverURL` in [config/_default/params.toml](config/_default/params.toml)
- `URL` in [assets/js/submit.js](assets/js/submit.js) and [assets/js/email.js](assets/js/email.js)

Tip: use `http://www.randomiban.com/?country=Netherlands` to test with IBAN numbers.

**Make sure to not commit these changes!**

## Maintainance
### Frontend
- `npm update` and test the website. TODO: Migrate to masonary layout when it is stable.

### Backend
- `go get -u -all` and test the endpoint by submitting a form with the frontend.

### Both
Check every once in a while if the dependencies are still maintained and switch to updated ones where possible (especially with the backend).

## Adding new activity posts
- Adding a new activity happens in two folders: [content/nl/activities](content/nl/activities) and [assets/images/activities](assets/images/activities).
- To create a new activity post: use the template provided below or copy an existing file in the linked **content** folder.
  - Warning: folder and file names may not contain capital letters.
  - Use the `yyyy-mm-dd-my-activity` template to prevent duplicate activity entries over time and allow for sorting the files by date.
- To add pictures: create a folder with the same name as the corresponding file in the linked **assets** folder. All of the activity's images are placed here.
  - To force a specific image order in the gallery: prefix the filenames with `N__` where **N** indicates the position of an image.
- After adding pictures, run `npm run process-images` (before committing!). This process alters the images in two ways:
  - All pictures larger than 1080p are scaled down to 1080p to keep the size of the repository small.
  - All images are "rotated" in order to mitigate these bugs: https://discourse.gohugo.io/t/why-is-hugo-rotating-my-image/23329/5 and https://github.com/golang/go/issues/4341.
  - If you manually need to scale the images down: do not compress them! The website already takes care of compression.
  - If a photo is still rotated after running the script, the solution is simple: rotate the photo a few times with an editing program (don't think too hard: grab Windows Photos, Paint.NET or something similar, but beware: old school Windows Photo Viewer gives problems). If the rotation metadata is changed and the orientation of the photo is correct, the problem is solved.
- When finished: start the website and check if
  - The gallery works
  - Pictures are not rotated or clipped
  - The category and tag overview are correct and do not contain duplicate entries
  - The dates are correct
  - There are no content mismatches
  - Any links are working
- When making significant changes to the template: make sure to update it everywhere using a text replacement tool.

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