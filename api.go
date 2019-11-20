package main
import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"
    "strconv"
    "time"

  _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
)
type Pagenation struct{
    TotalData int `json:"total_data"`
    Limit     int  `json:"limit"`
    TotalPages int `json:"total_pages"`
    Page      int `json:"current_page"`
    Pro    []Project `json:"data"`

}

type Project struct {
    Id             int `json:"Id,omitempty"`
    ProjectName    string `json:"name,omitempty"`
    ManagerName    string `json:"manager_name,omitempty"`
    ManagerEmailID string `json:"manager_email_id,omitempty"`
    Flag           string `json:"flag,omitempty"`
}

type Position struct {
    ProjectName    string `json:"projectname"` 
    Between0to15   int `json:"between0to15"`
    Between15to30  int `json:"between15to30"`
    Between30to60  int `json:"between30to60"`
    Between60to90  int `json:"between60to90"`
    Between90to120 int `json:"between90to120"`
    Greaterthen120 int `json:"greaterthen120"`
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {//To set all the CORS request
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}


func dbConn() (db *sql.DB) {//Database connection
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "root987"
    dbName := "weekly_update"
    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(172.21.234.63:3306)/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}

func authMiddleware(next http.Handler) http.Handler {//Auth Middleware
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  db := dbConn()
  reqToken := r.Header.Get("Authorization")

    splitToken := strings.Split(reqToken, "Bearer ")
    reqToken = splitToken[1]
    user, _ := db.Query("SELECT * FROM token WHERE access_token=?",reqToken)
    if user.Next() == false {
    w.WriteHeader(http.StatusUnauthorized)
    return
    }
    next.ServeHTTP(w, r)
  })
}

func loggingMiddleware(next http.Handler) http.Handler {//Logging Middleware
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        
        log.Println(r.RequestURI)
        
        next.ServeHTTP(w, r)
    })
}


func (c *Commander) Putdata(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    var dat Project
    var pro Project
    Time := time.Now()
    json.NewDecoder(r.Body).Decode(&dat)
    Pn := dat.ProjectName
    Mn := dat.ManagerName
    Email := dat.ManagerEmailID
    Flag := "1"
    rows, _ := db.Query("SELECT flag,Id FROM Project WHERE name = ? AND manager_name=?",Pn,Mn)
    if rows.Next() == false{
	insForm, err := db.Prepare("INSERT INTO Project(name, manager_name,manager_email_id,flag,created_at,updated_at)VALUES(?,?,?,?,?,?)")
	if err != nil {
	    panic(err.Error())
	}
	insForm.Exec(Pn, Mn, Email,Flag,Time,Time)
	defer db.Close()
        setupResponse(&w, r)
        w.WriteHeader(http.StatusCreated)
    }else{
        rows.Scan(&pro.Flag,&pro.Id)
        fmt.Println(pro.Flag)
        if pro.Flag == "0"{
            db.Query("UPDATE Project SET flag = 1,updated_at = ? WHERE  Id = ?",Time,pro.Id)
            setupResponse(&w, r)
            w.WriteHeader(http.StatusCreated)
        }else{
            w.WriteHeader(http.StatusBadRequest)   
        }
    }
}



func (c *Commander) GetdataByManager(w http.ResponseWriter, r *http.Request) {// send all the data with te requested manager name
    db := dbConn()

        Key := "data"
        e := `"` + Key + `"`
        w.Header().Set("Etag", e)
        w.Header().Set("Cache-Control", "max-age=2592000") // 30 days
    
        if match := r.Header.Get("If-None-Match"); match != "" {
            if strings.Contains(match, e) {
                w.WriteHeader(http.StatusNotModified)
                return
            }
        }

   
        p := mux.Vars(r)
	key := p["id"]
        var per string = "'"+key+"%'"
        var Offset int
        Pages := r.FormValue("Pages")
        i1,_ := strconv.Atoi(Pages)
        Offset = 10*i1
        count, _ := db.Query("SELECT COUNT(manager_name) FROM Project WHERE manager_name LIKE"+ per +" AND flag =1")
        defer count.Close()
        rows, err := db.Query("SELECT name,manager_name,manager_email_id,flag,Id FROM Project WHERE manager_name LIKE"+ per +" AND flag = 1 LIMIT ?,10",Offset)
      	if err != nil {
            fmt.Println("error in running query")
      	    log.Fatal(err)
      	}
      	defer rows.Close()
        var co int
        var pro Project
        var Proj []Project
        for rows.Next() {
            rows.Scan(&pro.ProjectName,&pro.ManagerName,&pro.ManagerEmailID,&pro.Flag,&pro.Id)
            Proj = append(Proj,pro)
      	}
        if count.Next()!= false{
            count.Scan(&co)
        }else{
            co = 0
        }
        defer db.Close()
        var Pag Pagenation
        Pag.TotalData=co
        Pag.Limit=10
        Pag.Pro=Proj
        x1 := co/10
        x := co%10
        x2 := x1 + 1
        //str1 := strconv.Itoa(x1)
        //str2 := strconv.Itoa(x2)
        if x == 0{
            Pag.TotalPages = x1
        }else{
            Pag.TotalPages = x2
        }
        x,_= strconv.Atoi(Pages)
        if Pag.TotalPages != 0{
            x1 = x+1
        }
        Pag.Page = x1
        setupResponse(&w, r)
        w.Header().Set("Content-Type","application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(Pag)

  
	

}

func (c *Commander) GetdataByProject(w http.ResponseWriter, r *http.Request) {// send all the data with the requested project name

    db := dbConn()

        Key := "data"
        e := `"` + Key + `"`
        w.Header().Set("Etag", e)
        w.Header().Set("Cache-Control", "max-age=2592000") // 30 days
    
        if match := r.Header.Get("If-None-Match"); match != "" {
            if strings.Contains(match, e) {
                w.WriteHeader(http.StatusNotModified)
                return
            }
        }

        var Offset int
        Pages := r.FormValue("Pages")
        i1,_ := strconv.Atoi(Pages)
        Offset = 10*i1
        p := mux.Vars(r)
	key := p["id"]
        var per string = "'"+key+"%'"
        fmt.Println(per)
        count, _ := db.Query("SELECT COUNT(name) FROM Project WHERE name LIKE"+ per +"AND flag = 1")
        defer count.Close()
        rows, err := db.Query("SELECT name,manager_name,manager_email_id,flag,Id FROM Project WHERE name LIKE"+ per +"AND flag = 1 LIMIT ?,10",Offset)
      	if err != nil {
            fmt.Println("error in running query")
            log.Fatal(err)
      	}
      	defer rows.Close()
        var co int
        var pro Project
        var Proj []Project
        for rows.Next() {
            rows.Scan(&pro.ProjectName,&pro.ManagerName,&pro.ManagerEmailID,&pro.Flag,&pro.Id)
            Proj = append(Proj,pro)
        }

        if count.Next()!= false{
            count.Scan(&co)
        }else{
            co = 0 
        }
        defer db.Close()
        var Pag Pagenation
        Pag.TotalData=co
        Pag.Limit=10
        Pag.Pro=Proj
        x1 := co/10
        x := co%10
        x2 := x1 + 1
        //str1 := strconv.Itoa(x1)
        //str2 := strconv.Itoa(x2)
        if x == 0{
            Pag.TotalPages = x1
        }else{
            Pag.TotalPages = x2
        }
        x,_= strconv.Atoi(Pages)
        if Pag.TotalPages != 0{
            x1 = x+1
        }
        Pag.Page = x1
        setupResponse(&w, r)
        w.Header().Set("Content-Type","application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(Pag)
}


func (c *Commander) GetProjectName(w http.ResponseWriter,r *http.Request) {// send all the project name
    db := dbConn()
      //var Offset int
       /* Pages := r.FormValue("Pages")
        i1,_ := strconv.Atoi(Pages)
        Offset = 10*i1
        fmt.Println(Offset)*/
        rows, err := db.Query("SELECT DISTINCT name FROM Project WHERE flag = 1")
        if err != nil {
            fmt.Println("error in running query")
            log.Fatal(err)
        }
        defer rows.Close()
        defer db.Close()
        var Nam []string
        var Name string
        for rows.Next() {
            rows.Scan(&Name)
            Nam = append(Nam,Name)
        }
        setupResponse(&w, r)
        w.Header().Set("Content-Type","application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(Nam)
}


func (c *Commander) UpdateData(w http.ResponseWriter,r *http.Request) {//update the table with the given Id
    db := dbConn()

        var dat Project
        Time := time.Now()

	json.NewDecoder(r.Body).Decode(&dat)
        ID := dat.Id
	Pn := dat.ProjectName
	Mn := dat.ManagerName
	Email := dat.ManagerEmailID
	//Flag := dat.Flag

        update,_ := db.Query("UPDATE Project SET name = ?, manager_name = ?,manager_email_id = ?,updated_at= ? WHERE Id =?",Pn,Mn,Email,Time,ID)
        defer update.Close()
        defer db.Close()
        setupResponse(&w, r)
        w.Header().Set("Content-Type","application/json")
        w.WriteHeader(http.StatusCreated) 
}


func (c *Commander) DeleteData(w http.ResponseWriter,r *http.Request) {//delete data
    db := dbConn()
    var dat Project
   
        json.NewDecoder(r.Body).Decode(&dat)
        del,_ := db.Query("UPDATE Project SET flag = 0 WHERE  Id = ?",dat.Id)
        defer del.Close()
        defer db.Close()
        setupResponse(&w, r)
        w.Header().Set("Content-Type","application/json")
        w.WriteHeader(http.StatusOK)
}



func (c *Commander) GetData(w http.ResponseWriter, r *http.Request) {// get all data
    db := dbConn()

        key := "data"
        e := `"` + key + `"`
        w.Header().Set("Etag", e)
        w.Header().Set("Cache-Control", "max-age=2592000") // 30 days
    
        if match := r.Header.Get("If-None-Match"); match != "" {
            if strings.Contains(match, e) {
                w.WriteHeader(http.StatusNotModified)
                return
            }
        }


        var Offset int
        Pages := r.FormValue("Pages")
        i1,_ := strconv.Atoi(Pages)
        Offset = 10*i1
        count, _ := db.Query("SELECT COUNT(Id) FROM Project WHERE flag = 1 ")
        defer count.Close()
        rows, err := db.Query("SELECT name,manager_name,manager_email_id,flag,Id FROM Project WHERE flag = 1 LIMIT ?,10",Offset)
      	if err != nil {
            fmt.Println("error in running query")
      	    log.Fatal(err)
      	}
      	defer rows.Close()
        defer db.Close()
        var co int
        var pro Project
        var Proj []Project
        for rows.Next() {
            rows.Scan(&pro.ProjectName,&pro.ManagerName,&pro.ManagerEmailID,&pro.Flag,&pro.Id)
            Proj = append(Proj,pro)
      	}
        if count.Next() != false{
            count.Scan(&co)
        }else{
            co = 0
        }
        var Pag Pagenation
        Pag.TotalData=co
        Pag.Limit=10
        Pag.Pro=Proj
        //x,_ := strconv.Atoi(co)
        x1 := co/10
        x := co%10
        x2 := x1 + 1
        //str1 := strconv.Itoa(x1)
        //str2 := strconv.Itoa(x2)
        if x == 0{
            Pag.TotalPages = x1
        }else{
            Pag.TotalPages = x2
        }
        x,_= strconv.Atoi(Pages)
        if Pag.TotalPages != 0{
            x1 = x+1
        }
        Pag.Page = x1
        setupResponse(&w, r)
        w.Header().Set("Content-Type","application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(Pag)
}


func (c *Commander)GetOpenPositionByAging(w http.ResponseWriter, r *http.Request){
    db := dbConn()
    var j int
    rows, err := db.Query("SELECT DISTINCT project_name FROM open_positions")//getting all the project name
    if err != nil {
        fmt.Println("error in running query")
        log.Fatal(err)
    }
    defer rows.Close()
    var Nam []Position// array of the structure
    var pro Position// instance of the structure
    var names []string// array of all the project names
    var Name string// instance of the project name
    

    var Str string
    var aging int
    var Aging []int
    for rows.Next() {
        rows.Scan(&Name)
        names = append(names,Name)
       
    }
    for i := 0;i<len(names);i++ {
        N := names[i]
        fmt.Println(names[i],len(names))
        pos,err := db.Query("SELECT created_at FROM open_positions WHERE project_name=? AND flag = 1 ",N)
        if err != nil {
        fmt.Println("error in running query")
        log.Fatal(err)
        }
        defer pos.Close()
        t1 := time.Now()
        t := t1.Format("2006-01-02")
        fmt.Println(t)
        for pos.Next() {
            pos.Scan(&Str)
            fmt.Println(Str)
            DataDiff,err := db.Query("SELECT DATEDIFF(?,?)",t,Str)           
            if err != nil {
                fmt.Println("error in running query")
                log.Fatal(err)
            }
            defer DataDiff.Close()
            for DataDiff.Next(){
                DataDiff.Scan(&aging)
            }
            fmt.Println(aging)
            Aging = append(Aging,aging)

      	}
        count1 := 0
        count2 := 0
        count3 := 0
        count4 := 0
        count5 := 0
        count6 := 0
        for j =0 ; j<len(Aging);j++{
           
            fmt.Println(Aging[j])

            if Aging[j] < 15{
                 count1++
       
            }else if Aging[j] > 15&& Aging[j]<30{
                 count2++
       
            }else if Aging[j] > 30&& Aging[j]<60{
                 count3++
       
            }else if Aging[j] > 60&& Aging[j]<90{
                 count4++
       
            }else if Aging[j] > 90&& Aging[j]<120{
                 count5++
       
            }else{
                 count6++
       
            }
        }
        Aging = nil
        pro.ProjectName = N
        pro.Between0to15 = count1
        pro.Between15to30 = count2
        pro.Between30to60 = count3
        pro.Between60to90 = count4
        pro.Between90to120 = count5
        pro.Greaterthen120 = count6
        Nam = append(Nam,pro)
    }

    setupResponse(&w, r)
    w.Header().Set("Content-Type","application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(Nam)
}
