# Intesa San Mattia

## Membri del Progetto
Il Gruppo Cool PPS è formato da:
- Mattia Pezzotti (885965) - m.pezzotti3@campus.unimib.it
- Nome Cognome (MAT) - email
- Nome Cognome (MAT) - email

## Introduzione al Progetto
Intesa San Mattia è una piccola applicazione full-stack creata per esercizio e divertimento.
La parte di *Backend* è stata scritta in **GOLANG** con le seguenti librerie:
- [Gin Web Framework](https://github.com/gin-gonic/gin) framework per web.
- [UUID](https://github.com/gofrs/uuid) generatore di ID unici.
- [simdb](https://github.com/sonyarouje/simdb) piccolo embedded database.

Per la parte di *Frontend* è stato utilizzato **HTML+CSS+JS** (con CSS gentilmente gestito da [Bulma](https://bulma.io/)) senza framework JS.

## Descrizione Pipeline CI/CD
Una volta pushato il progetto alla repository verranno automaticamente avviati dei job, nel dettaglio:

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
Scarica tutte le dipendenze, secondo quello descritto nel pre-generato file *go.mod*
La pipeline viene fermata se vengono rilevati dei problemi nel codice.

### Test
Vengono eseguiti i vari test, sia Unit che Integration per verificare la corretta efficacia sia dei singoli componenti sia del sistema in se.

```
$ go fmt
$ go vet
$ go test unit_test.go
$ go test int_test.go
```

Il source code viene formattato, vengono rilevati costrutti non validi e poi vengono avviati i test veri e propri.

Scarica tutte le dipendenze, secondo quello descritto nel pre-generato file *go.mod*
La pipeline viene fermata se vengono rilevati dei problemi nel codice.

### Running with Docker
- Eseguire direttamente il runnabile da Docker Hub con il comando
  
```
$ docker run -dp 4000:4000 mattiapezzotti/pezzotti-api
``` 

