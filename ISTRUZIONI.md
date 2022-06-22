# Istruzioni

## Building from Source
Codice sorgente reperibile in questa [Repo Github](https://github.com/mattiapezzotti/example-api).
### GO
1. Installare **GO**, seguendo le istruzioni riportate nel [sito ufficiale](https://go.dev/doc/install).

2. Aprire il terminale nella cartella del progetto e scrivere 
``` 
$ go run .
``` 
3. Aprire il proprio browser preferito e scrivere **localhost:4000** nella barra URL

### Dipendenze
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

## Running with Docker
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

## Sistema Utilizzato per lo Sviluppo
### Architettura
- Sistema operativo: **Windows 11** a 64 bit, processore basato su **x64**
- Processore: **AMD Ryzen 7 3700U** with Radeon Vega Mobile Gfx     2.30 GHz  
- RAM installata: 8,00 GB (6,94 GB utilizzabile)  

### Software 
- Codice scritto con IDE: **VS Code v1.68**  
- Browser Utilizzato: **Google Chrome 102.0.5005.115 (Build ufficiale) (a 64 bit)**
- Testing: **Postman v9.20.0**

