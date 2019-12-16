package handler
import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "fmt"
    "reflect"
    "strings"
    "regexp"
  //"github.com/unrolled/logger"

     model"projectname_projectmanager/model"
     database "projectname_projectmanager/driver"
)

type Commander struct{}
var UserName string
var Role string

func authMiddleware(next http.Handler) http.Handler { //Auth Middleware
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var DetailRole string
		db := database.DbConn()
		defer db.Close()
		reqToken := r.Header.Get("Authorization")

		splitToken := strings.Split(reqToken, "Bearer ")
		reqToken = splitToken[1]
		user, _ := db.Query("SELECT username FROM token WHERE access_token=? AND flag = 1", reqToken)
                defer user.Close()
		if user.Next() != false {
			user.Scan(&UserName)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		role, _ := db.Query("SELECT role FROM login WHERE username = ?", UserName)
                defer role.Close()
		if role.Next() != false {
			role.Scan(&DetailRole)
			re := regexp.MustCompile("^Program Manager|^Project Manager")
			Role = re.FindString(DetailRole)
		}
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler { //Logging Middleware
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.RequestURI)

		next.ServeHTTP(w, r)
	})
}

func HandleRequests() {// to handle all the http requests
    var i int
    rout := getconfig()
    c := &Commander{}
    myRouter := mux.NewRouter().StrictSlash(true)
    secure := myRouter.PathPrefix("").Subrouter()
    secure.Use(authMiddleware)
    for i=0 ;i<len(rout.R);i++ {
        fmt.Println(rout.R[i].Path,rout.R[i].Callback,rout.R[i].Method)
        m := reflect.ValueOf(c).MethodByName(rout.R[i].Callback)//calling the method by name(string)
        Call := m.Interface().(func(http.ResponseWriter,*http.Request))//creating an interface of m and assign the function type to m to make it callable
        if rout.R[i].Method == "GET"{
            if rout.R[i].Authorization == "YES" {
                secure.HandleFunc(rout.R[i].Path,Call ).Methods(rout.R[i].Method).Queries("Pages","{Pages}")
            }else{
                myRouter.HandleFunc(rout.R[i].Path,Call ).Methods(rout.R[i].Method).Queries("Pages","{Pages}")
            }
        }else if rout.R[i].Method == ""{
            if rout.R[i].Authorization == "YES" {
                secure.HandleFunc(rout.R[i].Path,Call )
            }else{
                myRouter.HandleFunc(rout.R[i].Path,Call )
            }
        }else{
            if rout.R[i].Authorization == "YES" {
                secure.HandleFunc(rout.R[i].Path,Call ).Methods(rout.R[i].Method)
            }else{
                myRouter.HandleFunc(rout.R[i].Path,Call ).Methods(rout.R[i].Method)
            }
        }
     
    }
    myRouter.Use(loggingMiddleware)//calling logging middleware    
    log.Fatal(http.ListenAndServe(":8000", myRouter))
}



func  getconfig()(c model.Conf){
    yamlFile, err := ioutil.ReadFile("routes.yaml")
    if err != nil {
        log.Printf("yamlFile.Get err   #%v ", err)
    }
    err = yaml.Unmarshal([]byte(yamlFile), &c)
    if err != nil {
        log.Fatalf("Unmarshal: %v", err)
    }
    return c
}



