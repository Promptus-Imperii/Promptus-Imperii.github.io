# Promptusimperii
The new website of the Promptus Imperii study association.

## Technologies
If you want to work on this website there is a few technologies to study:
- Hugo
- Tailwind (with the plug-ins)
- Golang Templates

If you can't find some code for partials for example, check the go.mod file to see the dependencies that are loaded in from external sources.

## Post instructies
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
