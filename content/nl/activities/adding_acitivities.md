---
# Do not change this to true! This file is meant to be read in the coding editor.
draft: true
---

# Foto's toevoegen
_Lees eerst deze stappen door, maar kijk daarna vooral naar hoe eerdere activiteiten zijn gedaan. Vergeet bij het aanpassen van het template of bij het toevoegen van comments niet om dit overal te doen (het zijn textbestanden, dus gebruik vooral text replace.)_
1. Maak een fotobestand aan op de drive.
2. Maak een bestand aan met een samenvattende bestandsnaam. Voeg aan het begin van deze bestandsnaam het jaar toe: vele activiteiten herhalen zich namelijk en zo voorkom je conflicten.
3. Maak een map met een naam gelijk aan het bestand aan in assets/images/activities
4. Zet een selectie aan foto's in de map. Als je een bepaalde volgorde wilt forceren, gebruik dan [nummer]__[fotonaam], om de link met het fotobestand te bewaren. Vergeet niet formaten anders dan jpg/jpeg of png naar een van de twee om te zetten. Voorkeur gaat dan uit naar png. Als foto's te groot zijn voor Git zonder Git Large File Storage, zet ze dan om naar 1080p (alleen voor de website)
5. Verander de `image_folder` variabele naar de correcte map.
6. Pas de datum aan naar de correcte datum
7. Vul alle andere frontformattervariabelen in
8. Schrijf de post
9. Start de website op en kijk of:
   1.  De galerij werkt en alle foto's niet gedraait zijn*.
   2.  Er geen dubelle tags zijn onstaan
   3.  De datum klopt
   4.  Alle links werken

*Als een foto gedraait is heeft dat te maken met [deze](https://discourse.gohugo.io/t/why-is-hugo-rotating-my-image/23329/5) en [deze](https://github.com/golang/go/issues/4341) bug. Oplossing is simpel: draai een paar keer de foto een een foto-edit programma (denk niet te moeilijk: pak windows photos of iets dergelijks) en draai de foto een paar keer. Als de rotation metadata verandert is en de foto klopt, is het gefixed.