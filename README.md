# Intesa San Mattia

## Scritto in
Il progetto è stato scritto in **GOLANG** per la parte di *Backend* con il Framework **GinGonic**, l'embedded database creato **Simbd** creato da *Sony Arouje* e la libreria di supporto **UUID** per la generazione di questi.

Per la parte di *Frontend* è stato utilizzato **HTML+CSS+JS** (con CSS gentilmente gestito da [Bulma](https://bulma.io/)) senza framework JS per evitare che il bundle risultasse troppo pesante 


## Considerazioni
Il database è già popolato con qualche dato di prova, se volete reinizializzarlo basta cancellare entrambi i file nella cartella *data*.  

La cancellazione di un qualsiasi account non comporta la cancellazione di nessun movimento di quest'ultimo.

Una transazione di 0 euro non è considerata valida, non penso sia saggio lasciare effettuare transizioni nulle che intaserebbero il database dei movimenti.


Una transazione con mittente e destinatario identici non è considerata valida, non vedo ragione logica di fare questo movimento, e come sopra, intaserebbe i movimenti.

Gli endpoint raggiungibili sono quelli descritti nel [progetto.pdf](https://elearning.unimib.it/pluginfile.php/1343307/mod_resource/content/4/progetto.pdf) tramite localhost:4000/*endpoint*
