# Intesa San Mattia

## Creato da
Il progetto è stato creato da **Mattia Pezzotti**, *MAT 885965, m.pezzotti3@campus.unimib.it*

## Scritto in
Il progetto è stato scritto in **GOLANG** per la parte di *Backend* con il Framework **GinGonic** e l'embedded database creato **Simbd** creato da *Sony Arouje*

Per la parte di *Frontend* è stato utilizzato **HTML+CSS+JS** (con CSS gentilmente gestito da [Bulma](https://bulma.io/)) senza framework JS per evitare che il bundle risultasse troppo pesante 


## Considerazioni
Il database è già popolato con qualche dato di prova, se volete reinizializzarlo basta cancellare entrambi i file nella cartella *data*.  
La cancellazione di un qualsiasi account non comporta la cancellazione di nessun movimento di quest'ultimo.

Una transazione di 0 euro non è considerata valida, vi sfido ad andare in una qualsiasi banca e chiedere di spostare 0 euro.
Una transazione con mittente e destinatario identici non è considerata valida, non vedo ragione logica di fare questo movimento.

Gli endpoint raggiungibili sono quelli descritti nel [progetto.pdf](https://elearning.unimib.it/pluginfile.php/1343307/mod_resource/content/4/progetto.pdf) tramite localhost:4000/*endpoint*