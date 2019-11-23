package main

import (
	"encoding/json"
	"github.com/maximecaruchet/simple-webapp-back-go/database"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type MetadataArray struct {
	Data   []database.User `json:"data"`
	Status string          `json:"status"`
}

type MetadataSimple struct {
	Data   database.User `json:"data"`
	Status string        `json:"status"`
}

func addUser(db *database.Database, user database.User) (database.User, error) {
	return db.AddUser(user)
}

func listUsers(db *database.Database) ([]database.User, error) {
	return db.ListUsers()
}

func userHandler(db *database.Database, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		log.Info("post user")

		decoder := json.NewDecoder(r.Body)

		var user database.User
		errDecode := decoder.Decode(&user)
		if errDecode != nil {
			log.Error(errDecode)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, errAddUser := addUser(db, user)
		if errAddUser != nil {
			log.Error(errAddUser)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, errJsonMarshal := json.Marshal(MetadataSimple{
			Data:   user,
			Status: "OK",
		})
		if errJsonMarshal != nil {
			log.Error(errJsonMarshal)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, errWrite := w.Write(b)
		if errWrite != nil {
			log.Error(errWrite)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if r.Method == http.MethodGet {
		log.Info("list users")

		users, errListUsers := listUsers(db)
		if errListUsers != nil {
			log.Error(errListUsers)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, errJsonMarshal := json.Marshal(MetadataArray{
			Data:   users,
			Status: "OK",
		})
		if errJsonMarshal != nil {
			log.Error(errJsonMarshal)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, errWrite := w.Write(b)
		if errWrite != nil {
			log.Error(errWrite)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	db := database.Open(database.DefaultDbName)
	defer func() {
		errClose := db.Close()
		if errClose != nil {
			log.Fatal(errClose)
		}
	}()

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		userHandler(db, w, r)
	})

	log.Info("Listen and serve on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
