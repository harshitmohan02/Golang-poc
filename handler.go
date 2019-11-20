package main
import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "fmt"
    "reflect"
  //"github.com/unrolled/logger"
)
type conf struct{
  
    R []Routes `yaml:"Routes"`
}


type Commander struct{}


type Routes struct{
    Path string `yaml:"Path"`
    Callback string `yaml:"Callback"`
    Method string `yaml:"Method"`
}


func handleRequests() {// to handle all the http requests
    var i int
    rout := getconfig()
    c := &Commander{}
    myRouter := mux.NewRouter().StrictSlash(true)
    for i=0 ;i<len(rout.R);i++ {
        fmt.Println(rout.R[i].Path,rout.R[i].Callback,rout.R[i].Method)
        m := reflect.ValueOf(c).MethodByName(rout.R[i].Callback)//calling the method by name(string)
        Call := m.Interface().(func(http.ResponseWriter,*http.Request))//creating an interface of m and assign the function type to m to make it callable
        if rout.R[i].Method == "GET"{
          //myRouter.Use(authMiddleware)//calling auth middleware
            myRouter.HandleFunc(rout.R[i].Path,Call ).Methods(rout.R[i].Method).Queries("Pages","{Pages}")
        }else if rout.R[i].Method == ""{
         // myRouter.Use(authMiddleware)
            myRouter.HandleFunc(rout.R[i].Path,Call)
        }else{
         // myRouter.Use(authMiddleware)
            myRouter.HandleFunc(rout.R[i].Path,Call ).Methods(rout.R[i].Method)
        }
     
    }
    myRouter.Use(loggingMiddleware)//calling logging middleware    
    log.Fatal(http.ListenAndServe(":8000", myRouter))
}



func  getconfig()(c conf){
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



