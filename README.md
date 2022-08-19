# Intesa San Mattia

## Endpoint
Gli endpoint raggiungibili tramite localhost:4000/*endpoint* e sono leggibili [qui](https://github.com/mattiapezzotti/example-api/blob/main/endpoint.md).

## Scritto in
Il progetto è stato scritto in **GOLANG** per la parte di *Backend* con il Framework **GinGonic**, l'embedded database creato **Simbd** creato da *Sony Arouje* e la libreria di supporto **UUID** per la generazione di questi.

Per la parte di *Frontend* è stato utilizzato **HTML+CSS+JS** (con CSS gentilmente gestito da [Bulma](https://bulma.io/)) senza framework JS.


## Considerazioni
Il database è già popolato con qualche dato di prova, se volete reinizializzarlo basta cancellare entrambi i file nella cartella *data*.  

La cancellazione di un qualsiasi account non comporta la cancellazione di nessun movimento di quest'ultimo.

Una transazione di 0 euro non è considerata valida, non penso sia saggio lasciare effettuare transizioni nulle che intaserebbero il database dei movimenti.


Una transazione con mittente e destinatario identici non è considerata valida, non vedo ragione logica di fare questo movimento, e come sopra, intaserebbe i movimenti.

## Istruzioni

### Building from Source
Codice sorgente reperibile in questa [Repo Github](https://github.com/mattiapezzotti/example-api).
#### GO
1. Installare **GO**, seguendo le istruzioni riportate nel [sito ufficiale](https://go.dev/doc/install).

2. Aprire il terminale nella cartella del progetto e scrivere 
``` 
$ go run .
``` 
3. Aprire il proprio browser preferito e scrivere **localhost:4000** nella barra URL

#### Dipendenze
Se per qualche motivo le dipendenze non vengono installate automaticamente, installarle manualmente:
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [UUID](https://github.com/gofrs/uuid)
- [simdb](https://github.com/sonyarouje/simdb)

Generalmente sono scaricabili tramite il comando
```
$ go get [github-link]
```
Se ancora ci sono problemi provare i comendi
```
$ go tidy mod
$ go build .
```

### Running with Docker
1. Scaricare, installare e aprire [Docker Desktop](https://www.docker.com/products/docker-desktop/) sulla propria macchina.
    
2. Scegliere se:
- Scaricare l' immagine docker ed eseguire il procedimento riportato sopra. L'immagine è reperibile da terminale con 
```
$ docker pull mattiapezzotti/pezzotti-api
``` 
- Eseguire direttamente il runnabile da Docker Hub con il comando
  
```
$ docker run -dp 4000:4000 mattiapezzotti/pezzotti-api
```


