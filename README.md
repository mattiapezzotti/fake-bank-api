# Intesa San Mattia

## Membri del Progetto
Il Gruppo Cool PPS è formato da:
- Mattia Pezzotti (885965) - m.pezzotti3@campus.unimib.it
- Thomas Howard-Grubb (869248) - t.howardgrubb@campus.unimib.it
- Simone Antonio Basile (826512) - s.basile@campus.unimib.it

## Introduzione al Progetto
Intesa San Mattia è una piccola applicazione full-stack creata per esercizio e divertimento.
La parte di *Backend* è stata scritta in **GOLANG** con le seguenti librerie:
- [Gin Web Framework](https://github.com/gin-gonic/gin) framework per web.
- [UUID](https://github.com/gofrs/uuid) generatore di ID unici.
- [simdb](https://github.com/sonyarouje/simdb) piccolo embedded database.

Per la parte di *Frontend* è stato utilizzato **HTML+CSS+JS** (con CSS gentilmente gestito da [Bulma](https://bulma.io/)) senza framework JS.

I vari endpoint sono descritti [qui](https://gitlab.com/unitestworks/2023_assignment1_ISM/-/blob/develop/endpoint.md).

## Descrizione Pipeline CI/CD
Una volta pushato il progetto alla repository verranno automaticamente avviati dei job, nel dettaglio:

### Cache
Viene utilizzata una cache per evitare di riscaricare ogni volta i moduli go.

### Build
Vengono richiamate tre istruzioni:
```
$ go mod download
```
Scarica tutte le dipendenze, secondo quello descritto nel pre-generato file *go.mod*

```
$ go mod tidy
```
Sistema le dipendenze in caso ce ne fossero di non utilizzate.

```
$ go build .
```
Effettua una build vera e propria del progetto, compilando il codice sorgente.

### Verify
Per verificare che non ci siano problemi con il codice sorgente viene utilizzato un aggregatore di lint: [golangci-lint](https://golangci-lint.run/) pensato appositamente per una pipeline CI/CD.

```
$ allow_failure: false
$ golangci-lint run -v
```

La pipeline viene fermata se vengono rilevati dei problemi nel codice.

### Test
Per prima cosa viene formattato il file e viene controllato per costrutti sospetti.
Vengono eseguiti i vari test, sia Unit che Integration per verificare la corretta efficacia sia dei singoli componenti sia del sistema in se.

```
$ go fmt $(go list ./... | grep -v /vendor/)
$ go vet $(go list ./... | grep -v /vendor/)
$ go test -run "Unit"
$ go test -run "Integration"
```

Il source code viene formattato, vengono rilevati costrutti non validi e poi vengono avviati i test veri e propri.

### Package
Raccoglie le componenti necessari, inclusi il binario, la documentazione (README), i folder data e web, e li impacchetta in un file TAR compresso.
```
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


### Running with Docker
- Eseguire direttamente il runnabile da Docker Hub con il comando
  
```
$ docker run -dp 4000:4000 mattiapezzotti/pezzotti-api
``` 

### Documentazione
-Genera la documentazione e la salva in un file HTML locale chiamato "doc.html" utilizzando il server GoDoc temporaneamente avviato. 

```
$ go install golang.org/x/tools/cmd/godoc@latest
$ godoc -http=:6060 & # Start the GoDoc server in the background
$ sleep 10  # Wait for a moment to ensure the GoDoc server fully starts
$ wget -O doc.html http://localhost:6060/pkg   # Generate the documentation and save it to public/doc.html
$ kill %1  # Stop the GoDoc server
```

