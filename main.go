package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"

	db "github.com/sonyarouje/simdb"
)

type Account struct {
	AccountID string  `json:"id"`
	Name      string  `json:"name"`
	Surname   string  `json:"surname"`
	Balance   float64 `json:"balance"`
}

type Movimento struct {
	MovimentoID uuid.UUID `json:"id"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	Amount      float64   `json:"amount"`
	Time        time.Time `json:"timestamp"`
}

func randSeq(n int) string {
	var letters = []rune("abcdef123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (a Account) ID() (jsonField string, value interface{}) {
	value = a.AccountID
	jsonField = "id"
	return
}

func (m Movimento) ID() (jsonField string, value interface{}) {
	value = m.MovimentoID
	jsonField = "id"
	return
}

func database() *db.Driver {
	driver, err := db.New("data")
	if err != nil {
		log.Print(err)
	}
	return driver
}

func UUIDv4() uuid.UUID {
	u2, err := uuid.NewV4()
	if err != nil {
		log.Printf("failed to generate UUID: %v", err)
	}
	return u2
}

func checkID(s string) bool {
	return len(s) == 20
}

func main(){
	router := setupRouter()

	err := router.Run("0.0.0.0:4000")

	if err != nil {
		fmt.Println(err)
	}
}

func setupRouter() *gin.Engine {
	rand.Seed(time.Now().UnixNano())
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.LoadHTMLGlob("web/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	router.GET("/transfer", func(c *gin.Context) {
		c.HTML(http.StatusOK, "transfer.html", gin.H{})
	})

	router.GET("/newAccount", func(c *gin.Context) {
		c.HTML(http.StatusOK, "newAccount.html", gin.H{})
	})

	router.GET("/list", func(c *gin.Context) {
		c.HTML(http.StatusOK, "list.html", gin.H{})
	})

	router.GET("/addFunds", func(c *gin.Context) {
		c.HTML(http.StatusOK, "addFunds.html", gin.H{})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.GET("/api/account", getAccounts)
	router.POST("/api/account", postAccount)
	router.DELETE("/api/account", deleteAccount)

	router.GET("/api/account/:id", getAccountByID)
	router.POST("/api/account/:id", versamento_prelievo)
	router.PUT("/api/account/:id", updateWholeOwner)
	router.PATCH("/api/account/:id", updateOwner)
	router.HEAD("/api/account/:id", getOwner)

	router.POST("/api/transfer", trasferimentoDenaro)

	router.POST("/api/divert", giraTrasnazione)

	return router;
}

/****** api/account ******/

/*
   GET: restituisce la lista di tutti gli account nel sistema
*/
func getAccounts(c *gin.Context) {
	driver := database()
	var accounts []Account

	err := driver.Open(Account{}).Get().AsEntity(&accounts)
	if len(accounts) == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "No Account in Database"})
		log.Print(err)
		return
	}
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Accounts retreived correctly", "accounts": accounts})
}

/*
   POST: crea un nuovo account con i seguenti campi:
   ∗ name
   ∗ surname
   e ritorna nel body della risposta il nuovo id dell’account creato. L’id di un account `e una
   stringa di 20 caratteri rappresentante una sequenza di bytes, generati randomicamente
   all’occorrenza, codificati in esadecimale (ad esempio un accountId potrebbe essere 1087b347f1a59277eb98 ).
*/
func postAccount(c *gin.Context) {
	var newAccount Account

	if err := c.BindJSON(&newAccount); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	if strings.TrimSpace(newAccount.Surname) == "" || strings.TrimSpace(newAccount.Name) == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing Name or Surname"})
		return
	}

	newAccount.AccountID = randSeq(20)

	driver := database()

	if err := driver.Insert(newAccount); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Account Created", "newAccount": newAccount})
}

/*
   DELETE: elimina l’account con id specificato dal parametro URL id
*/
func deleteAccount(c *gin.Context) {
	var toDelete Account
	var deleted Account

	if err := c.BindJSON(&toDelete); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	if !checkID(toDelete.AccountID) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID not valid"})
		return
	}

	driver := database()

	if err := driver.Open(Account{}).Where("id", "=", toDelete.AccountID).First().AsEntity(&deleted); err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID does not exist"})
		return
	}

	if err := driver.Delete(toDelete); err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Account Deleted", "deletedUser": deleted})
}

/****** api/account/{accountID} ******/

/*
    GET: restituisce il nome e cognome del proprietario nonch ́e il saldo con un elenco
	degli identificativi di tutte le transazioni effettuate da accountId, in ordine cronologico
	ascendente (dalla pi`u vecchia alla più recente). Inoltre, introduce un header di risposta
	con chiave X-Sistema-Bancario. Il valore dell’header deve esprimere il nome e cognome
	del proprietario in formato nome;cognome.
*/
func getAccountByID(c *gin.Context) {
	var foundAccount Account
	var trasazioniOut []Movimento
	var transazioniIn []Movimento

	id := c.Param("id")

	if !checkID(id) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID not valid"})
		return
	}

	driver := database()

	if err := driver.Open(Movimento{}).Where("from", "=", id).Get().AsEntity(&trasazioniOut); err != nil {
		fmt.Println(err)
	}

	if err := driver.Open(Movimento{}).Where("to", "=", id).Get().AsEntity(&transazioniIn); err != nil {
		fmt.Println(err)
	}

	if err := driver.Open(Account{}).Where("id", "=", id).First().AsEntity(&foundAccount); err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID does not exist"})
		return
	}

	var s = foundAccount.Name + ";" + foundAccount.Surname

	c.Header("X-Sistema-Bancario", s)
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Account Found", "account": foundAccount,
		"In-Transaction": transazioniIn, "Out-Transaction": trasazioniOut})
}

/*
   POST: effettua un versamento di denaro con un importo specificato dalla chiave amount
   nel body della richiesta. Se amount è negativo, viene eseguito un prelievo. Nel caso
   di amount negativo, il server deve generare un errore se il saldo del conto non `e sufficiente, informando il client dell’insuccesso. In caso di successo, nel body della risposta
   viene restituito il nuovo saldo del conto ed un identificativo del versamento/prelievo in
   formato UUID v4.
*/
func versamento_prelievo(c *gin.Context) {
	var newMovimento Movimento
	var account Account

	id := c.Param("id")

	if !checkID(id) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID not valid"})
		return
	}

	newMovimento.MovimentoID = UUIDv4()
	newMovimento.From = "Unknown"
	newMovimento.To = id
	newMovimento.Time = time.Now()

	if err := c.BindJSON(&newMovimento); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	driver := database()

	if err := driver.Open(Account{}).Where("id", "=", id).First().AsEntity(&account); err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID does not exist"})
		return
	}

	if newMovimento.Amount+account.Balance < 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: Insufficent Funds"})
		return
	}

	if err := driver.Insert(newMovimento); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	account.Balance += newMovimento.Amount
	if err := driver.Update(account); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Transaction Completed", "id": newMovimento.MovimentoID, "newBalance": account.Balance})
}

/*
   PUT: modifica (sovrascrive) name e surname del proprietario del conto. Nel body
   devono quindi essere presenti le seguenti chiavi:
   ∗ name
   ∗ surname
*/
func updateWholeOwner(c *gin.Context) {
	var accountToUpdate Account

	id := c.Param("id")

	if !checkID(id) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID not valid"})
		return
	}

	if err := c.BindJSON(&accountToUpdate); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	var newName = accountToUpdate.Name
	var newSurname = accountToUpdate.Surname

	if strings.TrimSpace(newSurname) == "" || strings.TrimSpace(newName) == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing Name or Surname"})
		return
	}

	driver := database()

	if err := driver.Open(Account{}).Where("id", "=", id).First().AsEntity(&accountToUpdate); err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID does not exist"})
		return
	}

	accountToUpdate.Name = newName
	accountToUpdate.Surname = newSurname

	if err := driver.Update(accountToUpdate); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Account Updated", "updatedAccount": accountToUpdate})
}

/*
   PATCH: modifica (sovrascrive) name oppure surname del proprietario del conto. Nel
   body deve quindi essere presente solamente una tra le seguenti chiavi:
   ∗ name
   ∗ surname
*/
func updateOwner(c *gin.Context) {
	var accountToUpdate Account

	id := c.Param("id")

	if !checkID(id) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID not valid"})
		return
	}

	if err := c.BindJSON(&accountToUpdate); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	var newName = accountToUpdate.Name
	var newSurname = accountToUpdate.Surname

	if strings.TrimSpace(newSurname) == "" && strings.TrimSpace(newName) == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing Name or Surname"})
		return
	}

	if strings.TrimSpace(newSurname) != "" && strings.TrimSpace(newName) != "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid Input: only one parameter admitted"})
		return
	}

	var updateName = 0

	if newSurname == "" {
		updateName = 1
	}

	driver := database()

	if err := driver.Open(Account{}).Where("id", "=", id).First().AsEntity(&accountToUpdate); err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID does not exist"})
		return
	}

	if updateName == 1 {
		accountToUpdate.Name = newName
	} else {
		accountToUpdate.Surname = newSurname
	}

	if err := driver.Update(accountToUpdate); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Account Updated", "updatedAccount": accountToUpdate})
}

/*
    HEAD: restituisce nome e cognome del proprietario in un header di risposta con chiave
	X-Sistema-Bancario. Il valore dell’header deve essere in formato nome;cognome. Non
	deve essere presente alcun body di risposta.
*/
func getOwner(c *gin.Context) {
	var foundAccount Account
	id := c.Param("id")

	if !checkID(id) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID not valid"})
		return
	}

	driver := database()

	if err := driver.Open(Account{}).Where("id", "=", id).First().AsEntity(&foundAccount); err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID does not exist"})
		return
	}

	var s = foundAccount.Name + ";" + foundAccount.Surname

	c.Header("x-sistema-bancario", s)

}

/********  /api/transfer ********/

/*
   POST: effettua uno spostamento di denaro con amount positivo da un account a un
   altro. amount è specificato nel body della richiesta. Il server deve generare un errore se
   il saldo del conto di partenza non è sufficiente, informando il client dell’insuccesso. In
   caso di successo, nel body della risposta vengono restituiti i nuovi saldi degli account
   coinvolti nella transazione distinti per accountId ed un identificativo della transazione
   in formato UUID v4. Il body della richiesta presenta quindi i seguenti campi:
   ∗ from
   ∗ to
   ∗ amount
*/
func trasferimentoDenaro(c *gin.Context) {
	var newMovimento Movimento
	var fromAccount Account
	var toAccount Account

	newMovimento.MovimentoID = UUIDv4()
	newMovimento.Time = time.Now()

	if err := c.BindJSON(&newMovimento); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	if newMovimento.From == newMovimento.To {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: Cannot self transfer"})
		return
	}

	if newMovimento.Amount <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: Amount must be positive"})
		return
	}

	if !checkID(newMovimento.From) || !checkID(newMovimento.To) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID not valid"})
		return
	}

	driver := database()

	if err := driver.Open(Account{}).Where("id", "=", newMovimento.From).First().AsEntity(&fromAccount); err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID does not exist"})
		return
	}

	if err := driver.Open(Account{}).Where("id", "=", newMovimento.To).First().AsEntity(&toAccount); err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID does not exist"})
		return
	}

	if fromAccount.Balance-newMovimento.Amount < 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: Insufficent Funds"})
		return
	}

	if err := driver.Insert(newMovimento); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	toAccount.Balance += newMovimento.Amount
	fromAccount.Balance -= newMovimento.Amount
	if err := driver.Update(toAccount); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	if err := driver.Update(fromAccount); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Transaction Completed", "id": newMovimento.MovimentoID,
		"senderID": fromAccount.AccountID, "senderNewBalance": fromAccount.Balance,
		"receiverID": toAccount.AccountID, "receiverNewBalance": toAccount.Balance,
	})
}

/********  /api/divert ********/

/*
   POST: annulla una transazione con id specificato dalla chiave id nel body della richiesta
   ovvero crea una nuova transazione con un nuovo UUID v4 che inverte il trasferimento di
   denaro tra gli account interessati dalla transazione con identificativo id dell’ammontare
   che la transazione ha coinvolto, ma solamente se il saldo dell’account del precedente
   beneficiario consente questa operazione. In caso contrario genera un errore opportuno.
*/
func giraTrasnazione(c *gin.Context) {
	var newMovimento Movimento
	var foundMovimento Movimento
	var fromAccount Account
	var toAccount Account

	newMovimento.MovimentoID = UUIDv4()
	newMovimento.Time = time.Now()

	if err := c.BindJSON(&foundMovimento); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	driver := database()

	if err := driver.Open(Movimento{}).Where("id", "=", foundMovimento.MovimentoID.String()).First().AsEntity(&foundMovimento); err != nil {
		//log.Print(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID does not exist"})
		return
	}

	if err := driver.Open(Account{}).Where("id", "=", foundMovimento.From).First().AsEntity(&toAccount); err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID does not exist"})
		return
	}

	if err := driver.Open(Account{}).Where("id", "=", foundMovimento.To).First().AsEntity(&fromAccount); err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: ID does not exist"})
		return
	}

	newMovimento.From = foundMovimento.To
	newMovimento.To = foundMovimento.From
	newMovimento.Amount = foundMovimento.Amount

	if fromAccount.Balance-newMovimento.Amount < 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Error: Insufficent Funds"})
		return
	}

	if err := driver.Insert(newMovimento); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	toAccount.Balance += newMovimento.Amount
	fromAccount.Balance -= newMovimento.Amount
	if err := driver.Update(toAccount); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	if err := driver.Update(fromAccount); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An Error as Occurred"})
		log.Print(err)
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "Transaction Reverted", "id": newMovimento.MovimentoID,
		"senderID": fromAccount.AccountID, "senderNewBalance": fromAccount.Balance,
		"receiverID": toAccount.AccountID, "receiverNewBalance": toAccount.Balance,
	})
}
