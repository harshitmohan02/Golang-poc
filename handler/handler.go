package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	database "projectname_projectmanager/driver"
	model "projectname_projectmanager/model"
	"reflect"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

// Commander : structure for commander
type Commander struct{}

// User : for user authentication
var UserName, Role string
var expiration int64

func authMiddleware(next http.Handler) http.Handler { //Auth Middleware
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var error model.Error
		db := database.DbConn()
		defer db.Close()
		reqToken := r.Header.Get("Authorization")
		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]
		user, err := db.Query("SELECT username,role, expiration FROM token WHERE access_token=? AND is_active = '1'", reqToken)
		defer user.Close()
		if err != nil {
			WriteLogFile(err)
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
		}
		if user.Next() != false {
			err := user.Scan(&UserName, &Role, &expiration)
			if err != nil {
				WriteLogFile(err)
				fmt.Println(err)
				w.WriteHeader(http.StatusBadRequest)
			}
			Role = strings.ToLower(Role)
			tm := time.Unix(expiration, 0)
			currentTime := time.Now()
			inTime := currentTime.Before(tm)
			if inTime == false {
				updDB, err := db.Prepare("UPDATE token SET is_active=? WHERE username=?, access_token=?")
				if err != nil {
					WriteLogFile(err)
					panic(err.Error())
				}
				updDB.Exec(0, UserName, reqToken)
				defer updDB.Close()
				updDB, err = db.Prepare("UPDATE refresh_token SET is_active=? WHERE username=?, access_token=?")
				if err != nil {
					WriteLogFile(err)
					panic(err.Error())
				}
				updDB.Exec(0, UserName, reqToken)
				defer updDB.Close()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				error.Message = fmt.Sprintf("Error validating access token: Session has expired on %s", tm)
				json.NewEncoder(w).Encode(error)
				return
			}
			// user, err = db.Query("SELECT role FROM login WHERE username=?", UserName)
			// defer user.Close()
			// if err != nil {
			// 	WriteLogFile(err)
			// 	fmt.Println(err)
			// 	w.WriteHeader(http.StatusBadRequest)
			// }
			// if user.Next() != false {
			// 	err := user.Scan(&Role)
			// 	if err != nil {
			// 		WriteLogFile(err)
			// 		fmt.Println(err)
			// 		w.WriteHeader(http.StatusBadRequest)
			// 	}
			// }
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	})
}

func loggingMiddleware(next http.Handler) http.Handler { //Logging Middleware
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

//HandleRequests : handler function
func HandleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	secure := myRouter.PathPrefix("").Subrouter()
	secure.Use(authMiddleware)
	rout := getconfig()
	c := &Commander{}
	for i := 0; i < len(rout.R); i++ {
		fmt.Println(rout.R[i].Path, rout.R[i].Callback, rout.R[i].Method, rout.R[i].Authorization)
		m := reflect.ValueOf(c).MethodByName(rout.R[i].Callback)
		Call := m.Interface().(func(http.ResponseWriter, *http.Request))

		if rout.R[i].Authorization == "YES" {

			secure.HandleFunc(rout.R[i].Path, Call).Methods(rout.R[i].Method)
		} else {
			myRouter.HandleFunc(rout.R[i].Path, Call).Methods(rout.R[i].Method)
		}

	}
	myRouter.Use(loggingMiddleware)
	log.Fatal(http.ListenAndServe(":8008", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(myRouter)))
}

func getconfig() (c model.Conf) {
	yamlFile, err := ioutil.ReadFile("routes.yaml")
	if err != nil {
		WriteLogFile(err)
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal([]byte(yamlFile), &c)
	//fmt.Println("1")
	if err != nil {
		WriteLogFile(err)
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

//SetupResponse : to setup access control on requests
func SetupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

//WriteLogFile : error logging
func WriteLogFile(err error) {
	f, erro := os.OpenFile("logs/output.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if erro != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
	// log.Println(err)
}

//BadRequest : to handle bad requests
func BadRequest(w http.ResponseWriter, err error) {
	if err != nil {
		WriteLogFile(err)
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
