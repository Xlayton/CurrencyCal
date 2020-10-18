package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

//Address represents user address
type Address struct {
	City    string `json:"city"`
	State   string `json:"state"`
	Street  string `json:"street"`
	Unit    string `json:"unit,omitempty"`
	ZipCode string `json:"zip_code"`
}

//ShippingAddress represents user address
type ShippingAddress struct {
	City    string `json:"city"`
	State   string `json:"state"`
	Street  string `json:"street"`
	Unit    string `json:"unit,omitempty"`
	ZipCode string `json:"zip_code"`
}

//Income users income data
type Income struct {
	Amount     string `json:"amount"`
	Frequency  string `json:"frequency"`
	Occupation string `json:"occupation"`
	Source     string `json:"source"`
}

//Identity users identification
type Identity struct {
	DOB    string `json:"date_of_birth"`
	ID     string `json:"id"`
	IDType string `json:"id_type"`
}

type Contact struct {
	Name      string `json:"name"`
	AccountID int32  `json:"account_id"`
}

//User struct to represent a user in the db
type User struct {
	UUID            string          `json:"uuid,omitempty"`
	FirstName       string          `json:"first_name"`
	LastName        string          `json:"last_name"`
	PhoneNumber     string          `json:"mobile"`
	Username        string          `json:"username"`
	Email           string          `json:"email"`
	PassHash        string          `json:"password"`
	ProfileImage    string          `json:"image"`
	Agreements      []int32         `json:"agreements"`
	Address         Address         `json:"address"`
	Identity        Identity        `json:"identification"`
	Income          Income          `json:"income"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
	CardholderID    int16           `json:"cardholderid,omitempty"`
	AccountID       int32           `json:"accountID,omitempty"`
	Contacts        []Contact       `json:"contacts,omitempty"`
	Transactions    []int           `json:"transactions,omitempty"`
	//TODO Add Galileo data that matters :\
}

//GeneralResponse represents a JSON response back to the client on failure
type GeneralResponse struct {
	Code    int16  `json:"code"`
	Message string `json:"message"`
}

//LoginResponse represents JSON response back to client on login
type LoginResponse struct {
	Code    int16  `json:"code"`
	Message string `json:"message"`
	User    User   `json:"user"`
}

//AccessResp is the response given by Galileo API
type AccessResp struct {
	AccessToken string `json:"access_token"`
}

//Cardholder :\
type Cardholder struct {
	FirstName       string          `json:"first_name"`
	LastName        string          `json:"last_name"`
	PhoneNumber     string          `json:"mobile"`
	Email           string          `json:"email"`
	Agreements      []int32         `json:"agreements"`
	Address         Address         `json:"address"`
	Identity        Identity        `json:"identification"`
	Income          Income          `json:"income"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
}

//CardholderCreation for creating cardholder
type CardholderCreation struct {
	Cardholder Cardholder `json:"cardholder"`
	ProductID  int32      `json:"product_id"`
}

//CardholderResponse word
type CardholderResponse struct {
	CardholderID int16 `json:"cardholder_id"`
}

//Account lmao
type Account struct {
	AccountID     int32  `json:"account_id"`
	AccountNumber string `json:"account_number"`
	AccountType   string `json:"account_type"`
	Balance       int32  `json:"balance"`
	CreationDate  string `json:"creation_date"`
	Name          string `json:"name"`
	ProductID     int32  `json:"product_id"`
	RoutingNumber string `json:"routing_number"`
	Status        string `json:"status"`
}

//AccountCreateResponse lol
type AccountCreateResponse struct {
	Accounts []Account `json:"accounts"`
}

//DoubleUUID lmao
type DoubleUUID struct {
	AccountUUID       string  `json:"account_uuid"`
	AddingAccountUUID string  `json:"adding_account_uuid"`
	Amount            float32 `json:"amount"`
}

//SingleUUID :\
type SingleUUID struct {
	AccountUUID string `json:"account_uuid"`
}

//AssistLogin sadge
type AssistLogin struct {
	CardholderID int32 `json:"cardholder_id"`
}

//GalileoAccountResp hi
type GalileoAccountResp struct {
	AccountID     int       `json:"account_id"`
	AccountNumber string    `json:"account_number"`
	AccountType   string    `json:"account_type"`
	Balance       float64   `json:"balance"`
	CreationDate  time.Time `json:"creation_date"`
	Name          string    `json:"name"`
	ProductID     int       `json:"product_id"`
	RoutingNumber string    `json:"routing_number"`
	Status        string    `json:"status"`
}

//GalileoTransactionResp mom
type GalileoTransactionResp struct {
	DestinationTransactionID int `json:"destination_transaction_id"`
	SourceTransactionID      int `json:"source_transaction_id"`
}

//GalileoTransaction :|
type GalileoTransaction struct {
	Amount             int       `json:"amount"`
	Description        string    `json:"description"`
	Timestamp          time.Time `json:"timestamp"`
	TransactionID      int       `json:"transaction_id"`
	TransactionSubtype string    `json:"transaction_subtype"`
	TransactionType    string    `json:"transaction_type"`
}

//GalileoTransactionList fj
type GalileoTransactionList struct {
	HasMore      bool                 `json:"has_more"`
	PageSize     int                  `json:"page_size"`
	Transactions []GalileoTransaction `json:"transactions"`
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	//Prepare header for json response
	w.Header().Set("Content-Type", "application/json")
	//Checks for POST method, otherwise responds with 404
	if r.Method == "POST" {
		//Parses data given as multipart form data(needed for profile image)
		r.ParseMultipartForm(16 << 20)
		//This section gets the file data uploaded and defers closing the File generated
		var buf bytes.Buffer
		file, header, err := r.FormFile("profileimage")
		if err != nil {
			if err.Error() == "http: no such file" {
				json.NewEncoder(w).Encode(GeneralResponse{400, "Please upload a jpg png or jpeg < 16MB"})
			} else {
				json.NewEncoder(w).Encode(GeneralResponse{500, "Failure uploading image. Please try again in 1 minute"})
			}
			return

		}
		fileExt := filepath.Ext(header.Filename)
		if !checkFileExtension(fileExt, []string{".jpg", ".png", ".jpeg"}) {
			json.NewEncoder(w).Encode(GeneralResponse{400, "Please upload a jpg png or jpeg < 8MB"})
			return
		}
		defer file.Close()
		imageID, _ := uuid.NewUUID()
		imageIDString := imageID.String()
		imageFilePath := "./image/" + imageIDString + fileExt
		io.Copy(&buf, file)
		ioutil.WriteFile(imageFilePath, buf.Bytes(), 0644)
		//TODO update user with image
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	//Prepare header for json response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//Checks for POST method, otherwise responds with 404
	if r.Method == "POST" {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		user.Contacts = []Contact{}
		user.Transactions = []int{}
		uuid, _ := uuid.NewUUID()
		user.UUID = uuid.String()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		godotenv.Load()
		apiuser := os.Getenv("APIUSER")
		apipass := os.Getenv("APIPASSWORD")
		reqBody, _ := json.Marshal(map[string]string{
			"username": apiuser,
			"password": apipass,
		})
		resp, err := http.Post("https://sandbox.galileo-ft.com/instant/v1/login", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var accessResp AccessResp
		json.Unmarshal(body, &accessResp)
		cardholder := Cardholder{user.FirstName, user.LastName, user.PhoneNumber, user.Email, user.Agreements, user.Address, user.Identity, user.Income, user.ShippingAddress}
		cardholderInfo, _ := json.Marshal(CardholderCreation{Cardholder: cardholder, ProductID: 19469})
		req, err := http.NewRequest("POST", "https://sandbox.galileo-ft.com/instant/v1/cardholders", bytes.NewReader(cardholderInfo))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessResp.AccessToken)
		reqClient := &http.Client{}
		res, err := reqClient.Do(req)
		body, _ = ioutil.ReadAll(res.Body)
		log.Println(string(body))
		var cardholderResp CardholderResponse
		json.Unmarshal(body, &cardholderResp)
		user.CardholderID = cardholderResp.CardholderID
		accountURL := fmt.Sprintf("https://sandbox.galileo-ft.com/instant/v1/cardholders/%d/accounts", cardholderResp.CardholderID)
		log.Println(accountURL)
		req, _ = http.NewRequest("GET", accountURL, nil)
		req.Header.Set("Authorization", "Bearer "+accessResp.AccessToken)
		reqClient = &http.Client{}
		res, _ = reqClient.Do(req)
		body, _ = ioutil.ReadAll(res.Body)
		var accCreateResp AccountCreateResponse
		json.Unmarshal(body, &accCreateResp)
		user.AccountID = accCreateResp.Accounts[0].AccountID
		client, ctx := getDbConnection()
		defer client.Disconnect(ctx)
		coll := client.Database("budgetbuddy").Collection("users")
		_, err = coll.InsertOne(context.TODO(), user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		back, _ := json.Marshal(GeneralResponse{200, "OK"})
		w.Write(back)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	//Prepare header for json response
	w.Header().Set("Content-Type", "application/json")
	//Assure method is GET
	if r.Method == "GET" {
		//Parse data from params
		var loginInfo struct {
			UserID   string `json:"user_id"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&loginInfo)
		//Get and check for required fields
		if isStringEmpty(loginInfo.UserID) || isStringEmpty(loginInfo.Password) {
			errResp, _ := json.Marshal(GeneralResponse{400, "Please provide valid username and password"})
			w.Write(errResp)
			return
		}
		emailRegex := regexp.MustCompile(`^(?:[A-Za-z0-9!#$%&'*+\-/=?^_` + "`" + `{|}~])(?:\.?[A-Za-z0-9!#$%&'*+\-/=?^_` + "`" + `{|}~]+)+\@(?:[A-Za-z0-9!#$%&'*+\-/=?^_` + "`" + `{|}~]+)(?:\.?[A-Za-z0-9!#$%&'*+\-/=?^_` + "`" + `{|}~])+$`)
		client, ctx := getDbConnection()
		defer client.Disconnect(ctx)
		coll := client.Database("budgetbuddy").Collection("users")
		var foundUser User
		if emailRegex.Match([]byte(loginInfo.UserID)) {
			err = coll.FindOne(context.TODO(), bson.M{"email": loginInfo.UserID}).Decode(&foundUser)
		} else {
			err = coll.FindOne(context.TODO(), bson.M{"username": loginInfo.UserID}).Decode(&foundUser)
		}
		if err != nil {
			log.Println("Error logging in: " + loginInfo.UserID + ":" + loginInfo.Password)
			errResp, _ := json.Marshal(GeneralResponse{500, "Internal Server Error."})
			w.Write(errResp)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(foundUser.PassHash), []byte(loginInfo.Password))
		if err != nil {
			errResp, _ := json.Marshal(GeneralResponse{400, "Please provide valid username/email & password combination"})
			w.Write(errResp)
			return
		}
		resp, _ := json.Marshal(LoginResponse{200, "OK", foundUser})
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(resp)
	} else {
		w.Write([]byte("404 Page not found"))
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	//Prepare header for json response
	w.Header().Set("Content-Type", "application/json")
	//Assure method is PUT
	if r.Method == "PUT" {
		//Parse data from body
		r.ParseMultipartForm(32 << 20)
		//Get body fields
		reqUUID := r.Form.Get("uuid")
		username := r.Form.Get("username")
		oldPass := r.Form.Get("oldpass")
		newPass := r.Form.Get("newpass")
		email := r.Form.Get("email")
		file, header, err := r.FormFile("profileimage")
		//Check minimum required fields
		if isStringEmpty(reqUUID) || isStringEmpty(username) || isStringEmpty(oldPass) || isStringEmpty(email) {
			errResp, _ := json.Marshal(GeneralResponse{400, "Please provide necessary fields"})
			w.Write(errResp)
		}
		client, ctx := getDbConnection()
		defer client.Disconnect(ctx)
		coll := client.Database("budgetbuddy").Collection("users")
		var oldUser User
		coll.FindOne(context.TODO(), bson.M{"uuid": reqUUID}).Decode(&oldUser)
		isOldPassCorrectErr := bcrypt.CompareHashAndPassword([]byte(oldUser.PassHash), []byte(oldPass))
		if isOldPassCorrectErr != nil {
			log.Println("Unauthorized edit attempt on user: " + oldUser.UUID)
			errorRes, _ := json.Marshal(GeneralResponse{401, "Unauthorized"})
			w.Write(errorRes)
			return
		}
		if err != nil {
			if isStringEmpty(newPass) {
				update := bson.D{{Key: "$set", Value: bson.D{{Key: "username", Value: username}, {Key: "email", Value: email}}}}
				coll.UpdateOne(context.TODO(), bson.M{"uuid": reqUUID}, update)
			} else {
				newHash := hashPass([]byte(newPass))
				update := bson.D{{Key: "$set", Value: bson.D{{Key: "username", Value: username}, {Key: "email", Value: email}, {Key: "passhash", Value: string(newHash)}}}}
				coll.UpdateOne(context.TODO(), bson.M{"uuid": reqUUID}, update)
			}
		} else {
			fileExt := filepath.Ext(header.Filename)
			if !checkFileExtension(fileExt, []string{".jpg", ".png", ".jpeg"}) {
				json.NewEncoder(w).Encode(GeneralResponse{400, "Please upload a jpg png or jpeg < 8MB"})
				return
			}
			var buf bytes.Buffer
			newImageID, _ := uuid.NewUUID()
			newImageIDString := newImageID.String()
			newImageFilePath := "./image/" + newImageIDString + fileExt
			io.Copy(&buf, file)
			ioutil.WriteFile(newImageFilePath, buf.Bytes(), 0644)
			defer file.Close()
			os.Remove("." + oldUser.ProfileImage)
			if isStringEmpty(newPass) {
				update := bson.D{{Key: "$set", Value: bson.D{{Key: "username", Value: username}, {Key: "email", Value: email}, {Key: "profileimage", Value: newImageFilePath[1:]}}}}
				coll.UpdateOne(context.TODO(), bson.M{"uuid": reqUUID}, update)
			} else {
				newHash := hashPass([]byte(newPass))
				update := bson.D{{Key: "$set", Value: bson.D{{Key: "username", Value: username}, {Key: "email", Value: email}, {Key: "profileimage", Value: newImageFilePath[1:]}, {Key: "passhash", Value: string(newHash)}}}}
				coll.UpdateOne(context.TODO(), bson.M{"uuid": reqUUID}, update)
			}
		}
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	//Prepare header for json response
	w.Header().Set("Content-Type", "application/json")
	//Assure method is DELETE
	if r.Method == "DELETE" {
		//Parse data from params
		r.ParseForm()
		//Get and check required params
		uuid := r.Form.Get("uuid")
		if isStringEmpty(uuid) {
			errResp, _ := json.Marshal(GeneralResponse{400, "Please provide valid username and password"})
			w.Write(errResp)
			return
		}
		client, ctx := getDbConnection()
		defer client.Disconnect(ctx)
		coll := client.Database("budgetbuddy").Collection("users")
		var foundUser User
		coll.FindOne(context.TODO(), bson.M{"uuid": uuid}).Decode(&foundUser)
		_, err := coll.DeleteOne(context.TODO(), bson.M{"uuid": uuid})
		if err != nil {
			log.Println("Error deleting: " + uuid)
			errResp, _ := json.Marshal(GeneralResponse{500, "Internal Server Error."})
			w.Write(errResp)
			return
		}
		if !isStringEmpty(foundUser.ProfileImage) {
			err = os.Remove("." + foundUser.ProfileImage)
			if err != nil {
				log.Println(err)
			}
		}
		resp, _ := json.Marshal(GeneralResponse{200, "OK"})
		w.Write(resp)
	}
}

func addContact(w http.ResponseWriter, r *http.Request) {
	//Prepare header for json response
	w.Header().Set("Content-Type", "application/json")
	//Assure method is POST
	if r.Method == "POST" {
		var contactFields DoubleUUID
		err := json.NewDecoder(r.Body).Decode(&contactFields)
		if err != nil {
			http.Error(w, "400", http.StatusBadRequest)
			return
		}
		client, ctx := getDbConnection()
		defer client.Disconnect(ctx)
		coll := client.Database("budgetbuddy").Collection("users")
		var foundUser User
		var addingUser User
		coll.FindOne(context.TODO(), bson.M{"uuid": contactFields.AccountUUID}).Decode(&foundUser)
		coll.FindOne(context.TODO(), bson.M{"uuid": contactFields.AddingAccountUUID}).Decode(&addingUser)
		foundUser.Contacts = append(foundUser.Contacts, Contact{addingUser.FirstName, addingUser.AccountID})
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "contacts", Value: foundUser.Contacts}}}}
		coll.UpdateOne(context.TODO(), bson.M{"uuid": foundUser.UUID}, update)
		resp, _ := json.Marshal(GeneralResponse{200, "OK"})
		w.Write(resp)
	}
}

func getBalance(w http.ResponseWriter, r *http.Request) {
	//Prepare header for json response
	w.Header().Set("Content-Type", "application/json")
	//Assure method is POST
	if r.Method == "POST" {
		var balanceField SingleUUID
		err := json.NewDecoder(r.Body).Decode(&balanceField)
		if err != nil {
			http.Error(w, "400", http.StatusBadRequest)
			return
		}
		client, ctx := getDbConnection()
		defer client.Disconnect(ctx)
		coll := client.Database("budgetbuddy").Collection("users")
		var foundUser User
		coll.FindOne(context.TODO(), bson.M{"uuid": balanceField.AccountUUID}).Decode(&foundUser)
		godotenv.Load()
		apiuser := os.Getenv("APIUSER")
		apipass := os.Getenv("APIPASSWORD")
		reqBody, _ := json.Marshal(map[string]string{
			"username": apiuser,
			"password": apipass,
		})
		resp, err := http.Post("https://sandbox.galileo-ft.com/instant/v1/login", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var accessResp AccessResp
		json.Unmarshal(body, &accessResp)
		accountURL := fmt.Sprintf("https://sandbox.galileo-ft.com/instant/v1/cardholders/%d/accounts/%d", foundUser.CardholderID, foundUser.AccountID)
		log.Println(accountURL)
		req, _ := http.NewRequest("GET", accountURL, nil)
		req.Header.Set("Authorization", "Bearer "+accessResp.AccessToken)
		reqClient := &http.Client{}
		res, _ := reqClient.Do(req)
		body, _ = ioutil.ReadAll(res.Body)
		var accountResp GalileoAccountResp
		json.Unmarshal(body, &accountResp)
		send, _ := json.Marshal(map[string]string{
			"balance": fmt.Sprintf("%f", accountResp.Balance),
			"msg":     fmt.Sprintf("%s, Your balance is $%f", foundUser.FirstName, accountResp.Balance),
		})
		w.Write(send)
	}
}

func doTransaction(w http.ResponseWriter, r *http.Request) {
	//Prepare header for json response
	w.Header().Set("Content-Type", "application/json")
	//Assure method is POST
	if r.Method == "POST" {
		var contactFields DoubleUUID
		err := json.NewDecoder(r.Body).Decode(&contactFields)
		if err != nil || contactFields.Amount <= 0 {
			http.Error(w, "400", http.StatusBadRequest)
			return
		}
		client, ctx := getDbConnection()
		defer client.Disconnect(ctx)
		coll := client.Database("budgetbuddy").Collection("users")
		var spendingUser User
		var addingUser User
		coll.FindOne(context.TODO(), bson.M{"uuid": contactFields.AccountUUID}).Decode(&spendingUser)
		coll.FindOne(context.TODO(), bson.M{"uuid": contactFields.AddingAccountUUID}).Decode(&addingUser)
		godotenv.Load()
		apiuser := os.Getenv("APIUSER")
		apipass := os.Getenv("APIPASSWORD")
		reqBody, _ := json.Marshal(map[string]string{
			"username": apiuser,
			"password": apipass,
		})
		resp, err := http.Post("https://sandbox.galileo-ft.com/instant/v1/login", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var accessResp AccessResp
		json.Unmarshal(body, &accessResp)
		req2Body, _ := json.Marshal(map[string]string{
			"amount":                 fmt.Sprintf("%f", contactFields.Amount),
			"source_account_id":      fmt.Sprintf("%d", spendingUser.AccountID),
			"destination_account_id": fmt.Sprintf("%d", addingUser.AccountID),
		})
		req, err := http.NewRequest("POST", "https://sandbox.galileo-ft.com/instant/v1/transfers", bytes.NewReader(req2Body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessResp.AccessToken)
		reqClient := &http.Client{}
		res, err := reqClient.Do(req)
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var transResp GalileoTransactionResp
		json.Unmarshal(body, &transResp)
		spendingUser.Transactions = append(spendingUser.Transactions, transResp.SourceTransactionID)
		spendUpdate := bson.D{{Key: "$set", Value: bson.D{{Key: "transactions", Value: spendingUser.Transactions}}}}
		addingUser.Transactions = append(addingUser.Transactions, transResp.DestinationTransactionID)
		addingUpdate := bson.D{{Key: "$set", Value: bson.D{{Key: "transactions", Value: addingUser.Transactions}}}}
		coll.UpdateOne(context.TODO(), bson.M{"uuid": spendingUser.UUID}, spendUpdate)
		coll.UpdateOne(context.TODO(), bson.M{"uuid": addingUser.UUID}, addingUpdate)
		send, _ := json.Marshal(map[string]string{
			"msg": fmt.Sprintf("Successfully sent $%f to %s", contactFields.Amount, addingUser.FirstName),
		})
		w.Write(send)
	}
}

func accountInfo(w http.ResponseWriter, r *http.Request) {
	//Prepare header for json response
	w.Header().Set("Content-Type", "application/json")
	//Assure method is POST
	if r.Method == "POST" {
		var balanceField SingleUUID
		err := json.NewDecoder(r.Body).Decode(&balanceField)
		if err != nil {
			http.Error(w, "400", http.StatusBadRequest)
			return
		}
		client, ctx := getDbConnection()
		defer client.Disconnect(ctx)
		coll := client.Database("budgetbuddy").Collection("users")
		var foundUser User
		coll.FindOne(context.TODO(), bson.M{"uuid": balanceField.AccountUUID}).Decode(&foundUser)
		send, _ := json.Marshal(foundUser)
		w.Write(send)
	}
}

func mostRecentTransaction(w http.ResponseWriter, r *http.Request) {
	//Prepare header for json response
	w.Header().Set("Content-Type", "application/json")
	//Assure method is POST
	if r.Method == "POST" {
		var balanceField SingleUUID
		err := json.NewDecoder(r.Body).Decode(&balanceField)
		if err != nil {
			http.Error(w, "400", http.StatusBadRequest)
			return
		}
		log.Println(balanceField)
		client, ctx := getDbConnection()
		defer client.Disconnect(ctx)
		coll := client.Database("budgetbuddy").Collection("users")
		var foundUser User
		coll.FindOne(context.TODO(), bson.M{"uuid": balanceField.AccountUUID}).Decode(&foundUser)
		lastTransID := foundUser.Transactions[len(foundUser.Transactions)-1]
		log.Println(lastTransID)
		godotenv.Load()
		apiuser := os.Getenv("APIUSER")
		apipass := os.Getenv("APIPASSWORD")
		reqBody, _ := json.Marshal(map[string]string{
			"username": apiuser,
			"password": apipass,
		})
		resp, err := http.Post("https://sandbox.galileo-ft.com/instant/v1/login", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var accessResp AccessResp
		json.Unmarshal(body, &accessResp)
		accountURL := fmt.Sprintf("https://sandbox.galileo-ft.com/instant/v1/cardholders/%d/accounts/%d/transactions", foundUser.CardholderID, foundUser.AccountID)
		log.Println(accountURL)
		req, _ := http.NewRequest("GET", accountURL, nil)
		req.Header.Set("Authorization", "Bearer "+accessResp.AccessToken)
		reqClient := &http.Client{}
		res, _ := reqClient.Do(req)
		body, _ = ioutil.ReadAll(res.Body)
		var transactions GalileoTransactionList
		json.Unmarshal(body, &transactions)
		var sendTransaction GalileoTransaction
		for _, trans := range transactions.Transactions {
			if trans.TransactionID == lastTransID {
				sendTransaction = trans
				break
			}
		}
		send, _ := json.Marshal(sendTransaction)
		w.Write(send)
	}
}

func assistantLogin(w http.ResponseWriter, r *http.Request) {
	//Prepare header for json response
	w.Header().Set("Content-Type", "application/json")
	//Assure method is POST
	if r.Method == "POST" {
		var assistLogin AssistLogin
		err := json.NewDecoder(r.Body).Decode(&assistLogin)
		if err != nil {
			http.Error(w, "400", http.StatusBadRequest)
			return
		}
		client, ctx := getDbConnection()
		defer client.Disconnect(ctx)
		coll := client.Database("budgetbuddy").Collection("users")
		var user User
		log.Println(assistLogin.CardholderID)
		coll.FindOne(context.TODO(), bson.M{"cardholderid": assistLogin.CardholderID}).Decode(&user)
		send, _ := json.Marshal(map[string]string{
			"user_uuid": user.UUID,
		})
		w.Write(send)
	}
}

func main() {
	http.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("./image"))))
	http.HandleFunc("/createuser", createUser)
	http.HandleFunc("/updateuser", updateUser)
	http.HandleFunc("/getuser", getUser)
	http.HandleFunc("/deleteuser", deleteUser)
	http.HandleFunc("/addcontact", addContact)
	http.HandleFunc("/getbalance", getBalance)
	http.HandleFunc("/dotransaction", doTransaction)
	http.HandleFunc("/assistlogin", assistantLogin)
	http.HandleFunc("/accinfo", accountInfo)
	http.HandleFunc("/mostrecenttrans", mostRecentTransaction)
	log.Fatal(http.ListenAndServe(":10000", nil))

}

func hashPass(pass []byte) []byte {
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}
	return hash
}

func isStringEmpty(str string) bool {
	return len(str) <= 0
}

func checkFileExtension(extensionToCheck string, validExtensions []string) bool {
	for _, ext := range validExtensions {
		if ext == extensionToCheck {
			return true
		}
	}
	return false
}

func getDbConnection() (*mongo.Client, context.Context) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	dbURI := os.Getenv("DBURI")
	client, err := mongo.NewClient(options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client, ctx
}
