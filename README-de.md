## Kleingebäck - kleinanzeigen.de Backup

![Kleingebaeck Logo](https://codeberg.org/scip/kleingebaeck/raw/branch/main/.github/assets/kleingebaecklogo-small.png)

[![Go Report Card](https://goreportcard.com/badge/codeberg.org/scip/kleingebaeck)](https://goreportcard.com/report/codeberg.org/scip/kleingebaeck) 
[![status-badge](https://ci.codeberg.org/api/badges/15530/status.svg)](https://ci.codeberg.org/repos/15530)
![GitHub License](https://img.shields.io/github/license/tlinden/kleingebaeck)
[![GitHub release](https://img.shields.io/github/v/release/tlinden/kleingebaeck?color=%2300a719)](https://codeberg.org/scip/kleingebaeck/releases)
[![English](https://codeberg.org/scip/kleingebaeck/raw/branch/main/.github/assets/english.png)](https://codeberg.org/scip/kleingebaeck/raw/branch/main/README.md)

Mit diesem Tool kann man seine Anzeigen bei https://kleinanzeigen.de sichern.

Es kann alle Anzeigen eines Users (oder nur eine Ausgewählte)
inklusive der Bilder herunterladen, die in einem Verzeichnis pro
Anzeige gespeichert werden. In dem Verzeichnis wird eine Datei
`Adlisting.txt` erstellt, in der sich die Inhalte der Anzeige wie
Titel, Preis, Text etc befinden. Bilder werden natürlich auch heruntergeladen.

## ACHTUNG - SICHERHEITS-UPDATE

Fertige vorcompilierte Programme älter als Version `v0.3.12` sind von
Schwachstellen in der Behandlung von HTTP und Zertifikaten
betroffen. Falls Du eine ältere Kleingebäck-Version im Einsatz hast,
bitte update auf Version `v0.3.12` oder höher. Bitte lies auch die [Release Notes für
v0.3.12](https://codeberg.org/scip/kleingebaeck/releases/tag/v0.3.12)
für mehr Details.

## Screenshots

Das ist die Hauptseite meines kleinanzeigen.de Accounts:

![Index](https://codeberg.org/scip/kleingebaeck/raw/branch/main/.github/assets/kleinanzeigen-index.png)

Sichern ich meine Anzeigen:

![Download](https://codeberg.org/scip/kleingebaeck/raw/branch/main/.github/assets/kleinanzeigen-download.png)

Backupverzeichnis nach dem Download:

![Download](https://codeberg.org/scip/kleingebaeck/raw/branch/main/.github/assets/kleinanzeigen-backup.png)

Verzeichnis einer Anzeige:

![Download](https://codeberg.org/scip/kleingebaeck/raw/branch/main/.github/assets/kleinanzeigen-ad.png)

**Das gleiche unter Windows:**

Anzeigen Sichern:

![Download](https://codeberg.org/scip/kleingebaeck/raw/branch/main/.github/assets/cmd-windows.jpg)

Backupverzeichnis nach dem Download

![Download](https://codeberg.org/scip/kleingebaeck/raw/branch/main/.github/assets/liste-windows.jpg)

Und eine Anzeige:

![Download](https://codeberg.org/scip/kleingebaeck/raw/branch/main/.github/assets/adlisting-windows.jpg)

## Installation

Das Tool hat keine weiteren Abhängigkeiten und erfordert auch keine
Anmeldung oder ähnliches. Man kädt sich einfach die ausführbare Datei
für seine Plattform herunter und kann direkt loslegen.

### Installation des vorcompilierten Programms

Auf der Seite [des letzten Releases](https://codeberg.org/scip/kleingebaeck/releases) findet man das Program für sein Betriebssystem und die Plattform (z.b. Windows + Intel)
 
Es gibt 2 Varianten:

1. Direkt das fertige Program für seine Plattform+OS herunterladen,
   z.B. `kleingebaeck-linux-amd64-0.0.5`, nach `kleingebaeck`
   umbenennen und in ein Verzeichnis kopieren, das im `PATH` ist,
   (z.B. nach `$HOME/bin` oder als root nach `/usr/local/bin`).

Um sicher zu gehen, dass an dem Program nicht verändert wurde, kann
man die Signatur vergleichen. Für jeden Download gibt es eine dazu
passende Signatur, in unserem Beispiel wäre das
`kleingebaeck-linux-amd64-0.0.5.sha256`.

Zum Verifizieren ausführen:

```shell
cat kleingebaeck-linux-amd64-0.0.5.sha25 && sha256sum kleingebaeck-linux-amd64-0.0.5
```
Man sollte zweimal den gleichen SHA256 Hash sehen.

2. Man kann auch einen Tarball (tgz Dateiendung) herunterladen,
   auspacken und mit GNU Make installieren:

```shell
tar xvfz kleingebaeck-linux-amd64-0.0.5.tar.gz
cd kleingebaeck-linux-amd64-0.0.5
sudo make install
```

### Installation aus dem Sourcecode

Man muss eine funktionierende Go Buildumgebung in der Version 1.21
installiert haben, um das Programm selber zu compilieren. GNU Make ist
hilfreich, aber nicht unbedingt erforderlich.

Um das Programm zu compilieren, muss man folgende Schritte ausführen:

```shell
git clone https://codeberg.org/scip/kleingebaeck.git
cd kleingebaeck
go mod tidy
make # (oder make)
sudo make install
```

### Docker image benutzen

Ein fertiges Dockerimage mit der aktuellen Programmversion ist immer
verfügbar. Man kann damit z.B. das Tool testen, bevor man es dauerhaft
benutzen möchte.

Um das Image herunterzuladen:
```
docker pull ghcr.io/tlinden/kleingebaeck:latest
```

Um kleingebäck im Image auszuführen und Daten ins lokale Filesystem zu
sichern, kann man so vorgehen:

```shell
mkdir anzeigen
docker run -u `id -u $USER` -v ./anzeigen:/backup ghcr.io/tlinden/kleingebaeck:latest -u XXX -v
ls -l anzeigen/ein-buch-mit-leeren-seiten
total 792
drwxr-xr-x 2 scip root   4096 Jan 23 12:58 ./
drwxr-xr-x 3 scip scip   4096 Jan 23 12:58 ../
-rw-r--r-- 1 scip root 131650 Jan 23 12:58 1.jpg
-rw-r--r-- 1 scip root  81832 Jan 23 12:58 2.jpg
-rw-r--r-- 1 scip root 134050 Jan 23 12:58 3.jpg
-rw-r--r-- 1 scip root   1166 Jan 23 12:58 Adlisting.txt
```

Hier wird der aktuelle User auf den User im Image gemappt und das
lokale Verzeichnis `anzeigen` nach `/backup` innerhalb des Images
gemountet.

Die Optionen `-u XXX -v` sind kleingebäck Optionen. Ersetze `XXX`
durch Deine tatsächliche kleinanzeigen.de Userid.

Eine Liste verfügbarer Images findet man [hier](https://codeberg.org/scip/kleingebaeck/pkgs/container/kleingebaeck/versions?filters%5Bversion_type%5D=tagged)

## Kommandozeilen Optionen:

```
Usage: kleingebaeck [-dvVhmoc] [<ad-listing-url>,...]
Options:
-u --user    <uid>      Backup ads from user with uid <uid>.
-d --debug              Enable debug output.
-v --verbose            Enable verbose output.
-o --outdir  <dir>      Set output dir (default: current directory)
-l --limit   <num>      Limit the ads to download to <num>, default: load all.
-c --config  <file>     Use config file <file> (default: ~/.kleingebaeck).
   --ignoreerrors       Ignore HTTP errors, may lead to incomplete ad backup.
-m --manual             Show manual.
-h --help               Show usage.
-V --version            Show program version.

If one  or more <ad-listing-url>'s  are specified, only  backup those,
otherwise backup all ads of the given user.
```

## Konfiguration

Man    kann    anstelle     von    Kommandlineoptionen    auch    eine
Konfigurationsdatei  verwenden.  Sie  befindet  sich standardmäßig  in
`~/.kleingebaeck`  aber man  kann  mit dem  Parameter  `-c` auch  eine
andere Datei angeben.

Das Format (TOML) ist einfach:

```
user = 1010101
loglevel = verbose
outdir = "test"
```

Im Source gibt es eine Beispieldatei `example.conf` mit Kommentaren.

## Umgebungsvariablen

Man kann darüber hinaus auch Umgebungsvariablen verwenden. Sie
entsprechen den Konfigurationsoptionen, aber gross geschrieben mit dem
Präfix `KLEINGEBAECK_`, z.B.

```shell
% KLEINGEBAECK_OUTDIR=/backup kleingebaeck -v
```

## Benutzung

Um das Tool einsetzen zu können, muss man zunächst seine Userid bei
kleinanzeigen.de herausfinden. Dazu ruft man am besten die Liste
seiner Anzeigen auf, während man NICHT eingeloggt ist:

https://www.kleinanzeigen.de/s-bestandsliste.html?userId=XXXXXX

Der `XXXXX` Teil der URL ist die Userid.

Trage diese Userid in der Konfigurationsdatei ein wie oben
beschrieben. Gib ausserdem das Ausgabeverzeichnis an. Dann einfach nur
`kleingebaeck` ausführen.

Innerhalb des Ausgabeverzeichnisses wird sich dann pro Anzeige ein
Unterverzeichnis befinden. Pro Anzeige gibt es eine Datei
`Adlisting.txt`, die etwa so aussieht:

```default
Title: A book I sell
Price: 99 € VB
Id: 1919191919
Category: Sachbücher
Condition: Sehr Gut
Created: 10.12.2023

This is the description text.

Pay with paypal.
```

Sowie alle Bilder.

Das Format kann man mit der Variable `template` in der Konfiguration
ändern. Die `example.conf` enthält ein Beispiel für das Standard Template.

## Verhalten des Tools

Es gibt einige Dinge über das Verhalten von kleingebäck, über die Du
Bescheid wissen solltest:

- alle HTML Seiten und Bilder werden immer heruntergeladen
- es wird ein (konfigurierbarer) Useragent verwendet
- HTTP Cookies werden beachtet
- bei Fehlern wird dreimal mit unterschiedlichem Abstand erneut
  versucht
- Bilder Downloads laufen parallelisiert mit leicht unterschiedlichen
  zeitlichen Abständen ab
- Gleich aussehende Bilder werden nicht überschrieben

Der letzte Punkt muss genauer erläutert werden:

Wenn man bei Kleinanzeigen.de eine Anzeige einstellt und Bilder
postet, werden diese dort in ihrer Grösse reduziert (durch Kompression
und Verkleinerung der Bilder usw.). Diese reduzierten Bilder werden
dann von kleingebäck heruntergeladen. Falls Du Deine original Bilder
behalten hast, kannst Du diese danach in das Backupverzeichnis
kopieren. Bei einem erneuten kleingebäck-Lauf werden diese Bilder dann
nicht überschrieben.

Wir verwenden dafür einen Algorythmus namens [distance
hashing](https://github.com/corona10/goimagehash). Dieser Algorithmus
prüft die Ähnlichkeit von Bildern. Diese können in ihrer Auflösung,
Kompression, Farbtiefe und vielem mehr manipuliert worden sein und
trotzdem als das "gleiche Bild" erkannt werden (wohlgemerkt nicht "das
selbe": die Dateien sind durchaus unterschiedlich!). Bis zu einer
Distance von 5 überschreiben wir keine Bilder, weil wir dann davon
ausgehen, dass das lokal Vorhandene das Original ist.

Bitte beachte aber, dass dies KEIN Cachingmechanismus ist: die Bilder
werden trotzdem immer alle heruntergeladen. Das muss so sein, da wir
uns nicht die Dateinamen anschauen können, da kleinanzeigen.de diese
nämlich zu Zahlen umbenennt. Und die Dateinamen können sich auch
ändern, wenn der User in der Anzeige die Bilder umarrangiert hat.

Du kannst dieses Verhalten mit der Option **--force** ausschalten. Du
kannst ausserdem mit der Option **--ignoreerrors** auch alle Fehler
ignorieren, die beim Bilderdownload auftreten könnten.

## Documentation

Die Dokumentation kann man
[online](https://codeberg.org/scip/kleingebaeck/raw/branch/main/kleingebaeck.pod)
oder lokal lesen mit: `kleingebaeck --manual`. Hat man das Tool mit
dem Tarball installiert, funktioniert auch `man kleingebaeck`.

## Kleingebäck?

Der Name kommt von "kleinanzeigen backup", verkürzt "klein back", das
englisch ausgesprochene "back" (deutsch bäck) führt dann zu "Kleingebäck".

## Wo bekommt man Hilfe

Obwohl ich gerne von kleingebäck Benutzern in privaten Mails höre, ist
das doch der beste Weg, die Anfrage zu übersehen und zu vergessen.

Um einen Fehler, ein unerwartetes Verhalten, eine Feature Request oder
einen Patch zu übermitteln, eröffne daher bitte einen Issue unter:
https://codeberg.org/scip/kleingebaeck/issues. Danke!

Bitte gebe den fehlgeschlagenen Befehl an, rufe es auch mit Debugging
`-d` auf.

## Ähnliche Projekte

Ich konnte kein Projekt finden, das speziell dafür geeignet ist,
Anzeigen bei kleinanzeigen.de zu sichern.

Aber es gibt ein Projekt, mit dem man ebenfalls Backups erstellen
kann:  [kleinanzeigen-bot](https://github.com/Second-Hand-Friends/kleinanzeigen-bot/).
Aber Vorsicht: kleinanzeigen.de bekämpft Bots aktiv, mit diesem hier
gibt es regelmäßige Probleme, z.B.: 
[issue](https://github.com/Second-Hand-Friends/kleinanzeigen-bot/issues/219).
Das Hauptproblem ist, dass diese Art von Bot sich mit Deinem Account
aktiv einloggt und mit der Seite interagiert. Damit kann die Firma die
Aktivitäten recht einfach Deinem User zuordnen und diesen **sperren**!
Also sei bitte vorsichtig!

**Kleingebäck** erfordert keinen Login, es verwendet lediglich die
öffentlich verfügbare Webseite und ruft diese auf, wie ein normaler
Browser. Tatsächlich gibt es meiner Meinung nach keinen Unterschied zu
einem Browserclient: beide laufen auf Anwenderseite auf Initiative
eines Benutzers. Und mit welchen Browser ich eine Webseite aufrufe,
bleibt immer noch mir überlassen und muss mir nicht von irgendwem
vorgeschrieben werden. Das schliesst die Verwendung von Kleingebäck
mit ein.

Hinzu kommt, dass dieses Tool nicht dazu gedacht ist, rund um die Uhr
zu laufen. Man ruft es ab und zu mal auf, wenn man halt neue Anzeigen
eingestellt hat, vielleicht einmal die Woche oder so. Man weiss ja
selber, wann man was geändert hat. Man benötigt trotzdem den Zugriff
mit dem Browser oder der mobilen App um Kleinanzeigen.de verwalten zu
können.

Meiner Ansicht nach ist das Risiko also sehr minimal, es handelt sich
meiner Meinung nach auch nicht um eine Verletzung der AGBs dort. Aber
das ist nur meine persönliche Meinung, bitte beachtet das. Am Ende
müsst Ihr selbst einschätzen und beurteilen wie hoch Ihr das Risiko
seht und ob Ohr es eingehen möchtet. Für eventuell auftretende
Konsequenzen bin ich nicht verantwortlich. Siehe auch [GPL Lizenz](LICENSE).

Es       gibt      noch       ein      weiteres       Tool      namens
[kleinanzeigen-enhanded](https://kleinanzeigen-enhanced.de/).  Das ist
eine    kostenpflichtige     vollständige    Anzeigenverwaltung    für
Profinutzer. Man  muss eine  monatliche Abogebühr bezahlen.   Das Tool
ist  als  Browsererweiterung  für  Google  Chrome  implementiert,  was
erklärt,  warum sie  Anzeigen  erstellen, ändern  und löschen  können,
obwohl es  gar keine  öffentliche API gibt.   Sieht nach  einer netten
ausgereiften Lösung aus. Mit Backups.
 
## Copyright und License

Lizensiert unter der GNU GENERAL PUBLIC LICENSE Version 3.

## Autor

T.v.Dein <tom AT vondein DOT org>

