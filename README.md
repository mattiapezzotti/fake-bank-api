# Intesa San Mattia

## Membri del Progetto
Il Gruppo Cool PPS è formato da:
- Mattia Pezzotti (885965) - m.pezzotti3@campus.unimib.it
- Thomas Howard-Grubb (869248) - t.howardgrubb@campus.unimib.it
- Simone Antonio Basile (826512) - s.basile@campus.unimib.it

## Introduzione al Progetto
[Intesa San Mattia](https://gitlab.com/unitestworks/2023_assignment1_ISM) è una piccola applicazione full-stack creata per esercizio e divertimento.
La parte di *Backend* è stata scritta in **GOLANG** con le seguenti librerie:
- [Gin Web Framework](https://github.com/gin-gonic/gin) framework per web.
- [UUID](https://github.com/gofrs/uuid) generatore di ID unici.
- [simdb](https://github.com/sonyarouje/simdb) piccolo embedded database.

Per la parte di *Frontend* è stato utilizzato **HTML+CSS+JS** (con CSS gentilmente gestito da [Bulma](https://bulma.io/)) senza framework JS.

I vari endpoint sono descritti [qui](https://gitlab.com/unitestworks/2023_assignment1_ISM/-/blob/develop/endpoint.md).

## Descrizione Pipeline CI/CD
Una volta pushato il progetto alla repository verranno automaticamente avviati sei job, nel dettaglio:

### Build
Vengono richiamate tre istruzioni:
```sh
$ go mod download
```
Scarica tutte le dipendenze, secondo quello descritto nel pre-generato file *go.mod*

```sh
$ go mod tidy
```
Sistema le dipendenze in caso ce ne fossero di non utilizzate.

```sh
$ go build .
```
Effettua una build vera e propria del progetto, compilando il codice sorgente.

### Verify
Per prima cosa viene formattato il file e viene controllato per costrutti sospetti.

```sh
$ go fmt $(go list ./... | grep -v /vendor/)
$ go vet $(go list ./... | grep -v /vendor/)
```

Per verificare che non ci siano problemi con il codice sorgente viene utilizzato un aggregatore di lint: [golangci-lint](https://golangci-lint.run/) pensato appositamente per una pipeline CI/CD. Viene utilizzata una image apposita, sostituendo quella generale.

```sh
$ golangci-lint run -v
```

La pipeline viene fermata se vengono rilevati dei problemi nel codice.

```sh
$ allow_failure: false
```

### Test
Vengono eseguiti i vari test, sia Unit che Integration separati in due job, per verificare la corretta efficacia sia dei singoli componenti sia del sistema in se.

```sh
$ go test -run "Unit"
$ go test -run "Integration"
```

Go permette facilmente di separare i test, permettendo di runnare solo i test con determinate stringhe all'interno della firma della funzione.

Anche in questo caso la pipeline viene fermata se i test falliscono.

### Package
Raccoglie le componenti necessari, inclusi il binario, la documentazione (README), i folder data e web, e li impacchetta in un file TAR compresso.

```sh
$ mkdir -p release
$ cp mybinaries/* release/
$ cp README.md release/
$ cp endpoint.md release/
$ cp -r web release/
$ cp -r data release/
$ tar -czf release.tar.gz release/
```


### Release
Prende il file release.tar.gz prodotto dal package-step e lo copia in un Docker container che è costruito a partire dal Dockerfile contenuto nel root del progetto. Attraverso una push il container viene inserito nella container directory di gitlab.
Il release-step parte solo dopo il completamento del package-step.

```
$ docker login -u gitlab-ci-token -p $CI_JOB_TOKEN $CI_REGISTRY
$ docker pull $CI_REGISTRY_IMAGE:latest || true
$ docker build --cache-from $CI_REGISTRY_IMAGE:latest -t $CI_REGISTRY_IMAGE:latest -f Dockerfile .
$ docker push $CI_REGISTRY_IMAGE:latest
```
### Dockerfile
Il Dockerfile specifica inanzitutto l'immagine di base per il Docker container:
```
FROM golang:1.15-alpine AS builder
```
Imposta la working directory:
```
WORKDIR /srv/app/
```
Copia il file release.tar.gz da locale nella directory /srv/app/ all'interno dell'immagine Docker.
```
COPY release.tar.gz /srv/app/
```
Imposta che nessun comando predefinito venga eseguito alla creazione del container.
```
ENTRYPOINT []
```

### Documentazione
Genera la documentazione e la salva in un file HTML locale chiamato "doc.html" utilizzando il server GoDoc temporaneamente avviato. 

```sh
$ go install golang.org/x/tools/cmd/godoc@latest
$ godoc -http=:6060 &
$ sleep 10 
$ wget -O doc.html http://localhost:6060/pkg
$ kill %1 
```
## Cache
Viene utilizzata una cache per evitare di riscaricare ogni volta i moduli go.
