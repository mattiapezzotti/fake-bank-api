# Endpoint
Endpoint raggiungibili.

## /api/account/

### GET
Restituisce la lista di tutti gli account nel sistema

### POST
Crea un nuovo account con i seguenti campi:
- name
- surname
  
ritorna nel body della risposta il nuovo id dell’account creato.

### DELETE

Elimina l’account con id specificato dal parametro URL id

## /api/account/{accountId}

### GET
Restituisce il nome e cognome del proprietario nonche il saldo con un elenco degli identificativi di tutte le transazioni effettuate da accountId, in ordine cronologico
ascendente (dalla più vecchia alla più recente). 

Inoltre, introduce un header di risposta
con chiave X-Sistema-Bancario. Il valore dell’header esprime il nome e cognome del proprietario in formato nome;cognome.

### POST
Effettua un versamento di denaro con un importo specificato dalla chiave amount nel body della richiesta.

### PUT
Sovrascrive name e surname del proprietario del conto. Nel body devono quindi essere presenti le chiavi name e surname

### PATCH
Sovrascrive name oppure surname del proprietario del conto. Nel body deve quindi essere presente solamente una chiave tra name o surname.

### HEAD 
Restituisce nome e cognome del proprietario in un header di risposta con chiave X-Sistema-Bancario.

## /api/transfer

### POST
Effettua uno spostamento di denaro con amount positivo da un account a un altro. Amount è specificato nel body della richiesta. Il body della richiesta presenta quindi i seguenti campi:
- from
- to
- amount

## /api/divert

### POST
Annulla una transazione con id specificato dalla chiave id nel body della richiesta