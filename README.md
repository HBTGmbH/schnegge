# SCHNEGGE
**SCHnell Noch Etwas Gutes Geld Erarbeiten**

Schnegge ermöglicht es die wichtigsten Funktionen (Buchungen anlegen und anschauen) 
von Salat über ein CLI aufzurufen. 

Einerseits kann Schnegge direkt vom Terminal bedient werden, andererseits kann 
Schnegge aber auch ohne Kommando aufgerufen werden und verfügt dann über eine interaktive 
Shell die bei der Eingabe hilft.

## Vorbereitung in Salat
Damit Schnegge Salat für dich aufrufen kann musst du zuerst in Salat ein Token generiert werden:
1. Im [Salat](https://salat.hbt.de) einloggen
2. `Verwaltung`-> `Einstellungen`
3. `Neuen Token Erzeugen`
4. Einen Kommentar und eine Gültigkeitsdauer eingeben und `Speichern` anklicken
5. Die TokenID und das TokenSecret muss jetzt schnegge bekannt gemacht werden.
   Diese werden dann unter `~/.schnegge` gespeichert:
   
   ```schnegge set -tokenID #### -tokenSecret ####```
6. Testen mit `schnegge list`

## Bedienung

Wenn `schnegge` ohne Parameter aufgerufen wird startet es mit einer interaktiven Shell.
Diese unterstützt den Nutzer bei der Eingabe. 
**Schnegge** kann aber auch bereits auf der Console direkt mir den Commands und
zugehörigen Parametern aufgerufen werden. **Schnegge** unterstützt aktuell die 
Befehle set, list und add. 

Auch bei einem Aufruf ohne ein Command können bereits Parameter angegeben werden. 
Diese bilden dann für diesen Aufruf eine Art Default-Werte. Ein Beispiel wäre z.B.
der Aufruf `schnegge -server http://localhost`. Hier werden dann alle weiteren Commands
gegen diesen Server aufgerufen, auch wenn dieser nicht per set gesetzt wurde.

### set
Setzt konfigurative Werte in der Datei `~/.schnegge`

Im Grunde können alle diese Parameter auch bei einem einzelnen Aufruf von *add* oder *list*
angegeben werden - dies macht nur meistens wenig Sinn (Ausnahme *-auftrag* für das Command *add*).

#### -tokenID
Die TokenID aus Salat.

#### -tokenSecret
Das Secret des Tokens aus Salat.

#### -auftrag
Der Auftrag, der zukünftig von schnegge als Default-Auftrag angenommen wird. 
Hier reicht ein Substring, dieser muss nur eindeutig sein.

Innerhalb der interaktiven Shell gibt es hier eine Unterstützung.

#### -server 
Ein alternativer Server statt https://salat.hbt.de (für die Entwicklung).

#### -noSplash
(kein Parameter notwendig) 

Der Splashscreen wird nicht mehr angezeigt.

### list
Zeigt deine Buchungen an.

#### -datum
Für welches Datum / Zeitbereich sollen die Buchungen angezeigt werden, default ist "heute".
Typische Werte sind:
- heute (default value)
- gestern
- vorgestern
- morgen
- woche / vorwoche (die aktuelle Woche, die vorherige Woche)
- monat (der aktuelle Monat)
- Montag, Dienstag ... der entspreche Wochentag innerhalb der letzten 7 Tage
- Januar, Februar ... der entsprechende Monat innerhalb der letzten 12 Monate
- 25.06.2022
- 2006-01-02
- -2 (von vor 2 Tagen bis heute)

#### -nachAuftrag
(kein Parameter notwendig) 

Die Ausgabe erfolgt nicht nach Datum sortiert, sondern nach Aufträgen.

### add
Fügt eine Buchung in Salat hinzu. Es ist allerdings nur möglich eine Buchung 
für einen Tag anzulegen, Serienbuchen werden bisher nicht unterstützt.

Alle Strings nach den Parametern werden als Kommentar für die Buchung angenommen.

#### -datum
Siehe `-datum` im Abschnitt `list`, allerdings wird hier der erste Tag innerhalb des Zeitraums für 
die Buchung genommen, wenn der Parameter mehrere Tage beschreibt.

So wird z.B. eine Buchung `schnegge add -datum -2` für Vorgestern vorgenommen.

#### -auftrag
Der Auftrag/Subauftrag auf den Gebucht werden soll.
Hier reicht ein Substring, dieser muss nur eindeutig sein.

Innerhalb der interaktiven Shell gibt es hier eine Unterstützung.

#### -stunden
Die Anzahl Stunden, default ist 0.

#### -minuten
Die Anzahl Minuten, default ist 0.

#### -fortbildung
(kein Parameter notwendig)

Gibt an, ob die Buchung Fortbildungscharakter hat (vergl. Haken in Salat)


## Anwendungsbeispiele

### Den aktuellen Tag anzeigen
```schnegge list```

oder 

```schnegge list -datum heute```

### Die letzten 5 Tage anzeigen
```schnegge list -datum -5```

### Den aktuellen Monat nach Projekten sortiert anzeigen
```schnegger list -datum monat -nachAuftrag```

### Den default Auftrag setzen
```schnegger set -auftrag HBT-Zeit```

### Auf HBT Zeit 2 Stunden 15 Minuten buchen
```schnegge add -auftrag HBT-Zeit -stunden 2 -minuten 15 Ich bin der Kommentar```

### Auf Projekt_X 8 Stunden mit Fortbildungsanteil buchen
```schnegge add -auftrag Projekt_X -stunden 8 -fortbildung Das war ein toller Workshop beim Kunden```

## ToDos
- Retry einbauen, falls es wegen Überlast zu falschen Antworten kommt (z.B. bei `schnegge list -datum Monat`)
- Tests schreiben
- Nach set (in der interaktiven Console) am besten die Config neu laden.

## Feature Ideen in der Zukunft
- Samstag/Sonntag und Feiertage einfärben
  - Salat Rest Schnittstelle für Feiertage bereitstellen
  
- Buchung löschen
  - Salat Rest Schnittstelle bereitstellen 
  - Workflow in schnegge bauen
    - z.B. 
      
      ```remove -datum heute -auftrag HBT-Zeit -buchung "Ich bin der Kommentar dazu" ```
- monatliche Freigabe ermöglichen
  - Salat Rest Schnittstelle bereitstellen
  - Workflow in schnegge bauen
    - z.B. erst den aktuellen Monat darstellen und dann muss noch mal explizit bestätigt werden
      
      ```
      > approv
      Aktuell freigegeben bis 31.12.2014
      > approv -datum Monat
      ...
      Möchtest du bis zum 31.06.2022 freigeben (ja/nein)
      > nein
      > approv -datum Monat -nachAuftrag
      ...
      Möchtest du bis zum 31.06.2022 freigeben (ja/nein)
      > ja
      Aktuell freigegeben bis 31.06.2022
      > ...
      ```