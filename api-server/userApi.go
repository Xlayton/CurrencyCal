package main

import (
	"bytes"
	"context"
	"encoding/json"
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

//User struct to represent a user in the db
type User struct {
	UUID         string `json:"uuid"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PassHash     string `json:"password"`
	ProfileImage string `json:"image"`
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

func createUser(w http.ResponseWriter, r *http.Request) {
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

		//Gets and checks username and passwordhash to check for 0 length
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")
		if isStringEmpty(username) || isStringEmpty(password) || isStringEmpty(email) {
			json.NewEncoder(w).Encode(GeneralResponse{400, "Please include a username, email, and password"})
			return
		}
		passwordHash := hashPass([]byte(password))
		//This chunk writes the image uploaded to machine storage to be used later
		imageID, _ := uuid.NewUUID()
		imageIDString := imageID.String()
		imageFilePath := "./image/" + imageIDString + fileExt
		io.Copy(&buf, file)
		ioutil.WriteFile(imageFilePath, buf.Bytes(), 0644)

		//Inserts user into DB with generated uuid
		client, ctx := getDbConnection()
		defer client.Disconnect(ctx)
		coll := client.Database("budgetbuddy").Collection("users")
		userID, _ := uuid.NewUUID()
		testUser := User{userID.String(), username, email, string(passwordHash), imageFilePath[1:]}
		_, err = coll.InsertOne(context.TODO(), testUser)
		if err != nil {
			log.Fatal(err)
		}
		client.Disconnect(ctx)

		//Writes success back
		w.Write([]byte("Inserted"))
	} else {
		w.Write([]byte("404 Page not found"))
	}
	//TODO Galileo CardHolder Implementation
}

func getUser(w http.ResponseWriter, r *http.Request) {
	//Prepare header for json response
	w.Header().Set("Content-Type", "application/json")
	//Assure method is GET
	if r.Method == "GET" {
		//Parse data from params
		r.ParseForm()
		//Get and check for required fields
		userID := r.Form.Get("userid")
		password := r.Form.Get("password")
		if isStringEmpty(userID) || isStringEmpty(password) {
			errResp, _ := json.Marshal(GeneralResponse{400, "Please provide valid username and password"})
			w.Write(errResp)
			return
		}
		emailRegex := regexp.MustCompile(`^(?:[A-Za-z0-9!#$%&'*+\-/=?^_` + "`" + `{|}~])(?:\.?[A-Za-z0-9!#$%&'*+\-/=?^_` + "`" + `{|}~]+)+\@(?:[A-Za-z0-9!#$%&'*+\-/=?^_` + "`" + `{|}~]+)(?:\.?[A-Za-z0-9!#$%&'*+\-/=?^_` + "`" + `{|}~])+$`)
		client, ctx := getDbConnection()
		defer client.Disconnect(ctx)
		coll := client.Database("budgetbuddy").Collection("users")
		var foundUser User
		var err error
		if emailRegex.Match([]byte(userID)) {
			err = coll.FindOne(context.TODO(), bson.M{"email": userID}).Decode(&foundUser)
		} else {
			err = coll.FindOne(context.TODO(), bson.M{"username": userID}).Decode(&foundUser)
		}
		if err != nil {
			log.Println("Error logging in: " + userID + ":" + password)
			errResp, _ := json.Marshal(GeneralResponse{500, "Internal Server Error."})
			w.Write(errResp)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(foundUser.PassHash), []byte(password))
		if err != nil {
			errResp, _ := json.Marshal(GeneralResponse{400, "Please provide valid username/email & password combination"})
			w.Write(errResp)
			return
		}
		resp, _ := json.Marshal(LoginResponse{200, "OK", foundUser})
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

func main() {
	http.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("./image"))))
	http.HandleFunc("/createuser", createUser)
	http.HandleFunc("/updateuser", updateUser)
	http.HandleFunc("/getuser", getUser)
	http.HandleFunc("/deleteuser", deleteUser)
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
