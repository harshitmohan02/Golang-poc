package handler
import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"

    model "projectname_projectmanager/model"
    database "projectname_projectmanager/driver"
    "github.com/gorilla/mux"
)

func setupResponse(w *http.ResponseWriter, req *http.Request) { //To set all the CORS request
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}





func (C *Commander) Putdata(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
        defer db.Close()
	if Role == "Program Manager" {
		var ProjectId int
		var ManagerId int
		var dat model.Project
		var pro model.Project
		Time := time.Now()
		json.NewDecoder(r.Body).Decode(&dat)
		ProjectName := dat.ProjectName
		ManagerName := dat.ManagerName
		Email := dat.ManagerEmailID
		Flag := "1"
		getProject, _ := db.Query("SELECT id FROM projects WHERE project_name = ?", ProjectName)
		defer getProject.Close()
		if getProject.Next() == false {
			InsertProject, _ := db.Prepare("INSERT INTO projects(project_name,created_by)VALUES(?,?)")
			InsertProject.Exec(ProjectName, UserName)
			defer InsertProject.Close()
		}
		getProjectId, _ := db.Query("SELECT id FROM projects WHERE project_name = ?", ProjectName)
		defer getProjectId.Close()
		if getProjectId.Next() != false {
			getProjectId.Scan(&ProjectId)
			fmt.Println(ProjectId)
		}

		getManager, _ := db.Query("SELECT id FROM manager WHERE manager_name = ? AND manager_email_id = ?", ManagerName, Email)
		defer getManager.Close()
		if getManager.Next() == false {
			InsertManager, _ := db.Prepare("INSERT INTO manager(manager_name, manager_email_id)VALUES(?, ?)")
			InsertManager.Exec(ManagerName, Email)
			defer InsertManager.Close()
		}
		getManagerId, _ := db.Query("SELECT id FROM manager WHERE manager_name = ? AND manager_email_id = ?", ManagerName, Email)
		defer getManagerId.Close()
		if getManagerId.Next() != false {
			getManagerId.Scan(&ManagerId)
			fmt.Println(ManagerId)
		}

		rows, err := db.Query("SELECT flag, Id FROM manager_details WHERE project_id = ? AND manager_id=?", ProjectId, ManagerId)
		if err != nil {
			panic(err.Error())
		}
		defer rows.Close()
		if rows.Next() == false {
			insForm, err := db.Prepare("INSERT INTO manager_details(project_id, manager_id, flag, created_at, updated_at)VALUES(?, ?, ?, ?, ?)")
			if err != nil {
				panic(err.Error())
			}
			insForm.Exec(ProjectId, ManagerId, Flag, Time, Time)
			defer insForm.Close()
			setupResponse(&w, r)
			w.WriteHeader(http.StatusCreated)
		} else {
			rows.Scan(&pro.Flag, &pro.Id)
			fmt.Println(pro.Flag)
			if pro.Flag == "0" {
				update, _ := db.Query("UPDATE manager_details SET flag = 1, updated_at = ? WHERE  id = ?", Time, pro.Id)
                                defer update.Close()
				setupResponse(&w, r)
				w.WriteHeader(http.StatusCreated)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (C *Commander) GetdataByManager(w http.ResponseWriter, r *http.Request) { // send all the data with te requested manager name

        db := database.DbConn()
        defer db.Close()
	fmt.Println(Role)
	if Role == "Program Manager" {
		p := mux.Vars(r)
		key := p["id"]
		var per string = "'" + key + "__%'"
		var Offset int
		Pages := r.FormValue("Pages")
		i1, _ := strconv.Atoi(Pages)
		Offset = 10 * i1
		count, _ := db.Query("SELECT COUNT(id) FROM manager_details WHERE flag = 1 AND manager_id in (select id from manager where manager_name LIKE"+per+") AND project_id in (SELECT id from projects WHERE created_by = ?)", UserName)
		defer count.Close()
		rows, err := db.Query("SELECT projects.project_name, manager.manager_name, manager.manager_email_id, manager_details.flag, manager_details.id FROM manager_details LEFT JOIN projects on manager_details.project_id = projects.id LEFT JOIN manager on manager_details.manager_id = manager.id WHERE manager.manager_name LIKE"+per+" AND projects.created_by = ? AND manager_details.flag = 1 LIMIT ?, 10", UserName, Offset)
		if err != nil {
			fmt.Println("error in running query")
			log.Fatal(err)
		}
		defer rows.Close()
		var co int
		var pro model.Project
		var Proj []model.Project
		for rows.Next() {
			rows.Scan(&pro.ProjectName, &pro.ManagerName, &pro.ManagerEmailID, &pro.Flag, &pro.Id)
			Proj = append(Proj, pro)
		}
		if count.Next() != false {
			count.Scan(&co)
		} else {
			co = 0
		}
		var Pag model.Pagenation
		Pag.TotalData = co
		Pag.Limit = 10
		Pag.Data = Proj
		x1 := co / 10
		x := co % 10
		x2 := x1 + 1
		if x == 0 {
			Pag.TotalPages = x1
		} else {
			Pag.TotalPages = x2
		}
		x, _ = strconv.Atoi(Pages)
		if Pag.TotalPages != 0 {
			x1 = x + 1
		}
		Pag.Page = x1
		setupResponse(&w, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Pag)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (C *Commander) GetdataByProject(w http.ResponseWriter, r *http.Request) { // send all the data with the requested project name

	db := database.DbConn()
        defer db.Close()
	if Role == "Program Manager" {
		var co int
		var Offset int
		Pages := r.FormValue("Pages")
		i1, _ := strconv.Atoi(Pages)
		Offset = 10 * i1
		p := mux.Vars(r)
		key := p["id"]
		var per string = "'" + key + "%'"
		fmt.Println(per)
		count, _ := db.Query("SELECT COUNT(id) FROM manager_details WHERE flag =1 AND project_id IN (select id from projects where project_name LIKE"+per+"AND created_by = ?)", UserName)
		defer count.Close()
		if count.Next() != false {
			count.Scan(&co)
		} else {
			co = 0
		}
		rows, err := db.Query("SELECT projects.project_name, manager.manager_name, manager.manager_email_id, flag,manager_details.id FROM manager_details LEFT JOIN projects on manager_details.project_id = projects.id LEFT JOIN manager on manager_details.manager_id = manager.id WHERE projects.project_name LIKE"+per+" AND projects.created_by = ? AND manager_details.flag = 1 LIMIT ?, 10", UserName, Offset)
		if err != nil {
			fmt.Println("error in running query")
			log.Fatal(err)
		}
		defer rows.Close()
		var pro model.Project
		var Proj []model.Project
		for rows.Next() {
			rows.Scan(&pro.ProjectName, &pro.ManagerName, &pro.ManagerEmailID, &pro.Flag, &pro.Id)
			Proj = append(Proj, pro)
		}
		var Pag model.Pagenation
		Pag.TotalData = co
		Pag.Limit = 10
		Pag.Data = Proj
		x1 := co / 10
		x := co % 10
		x2 := x1 + 1
		if x == 0 {
			Pag.TotalPages = x1
		} else {
			Pag.TotalPages = x2
		}
		x, _ = strconv.Atoi(Pages)
		if Pag.TotalPages != 0 {
			x1 = x + 1
		}
		Pag.Page = x1
		setupResponse(&w, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Pag)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (C *Commander) GetProjectName(w http.ResponseWriter, r *http.Request) { // send all the project name
	db := database.DbConn()
        defer db.Close()
	if Role == "Program Manager" {
		rows, err := db.Query("SELECT  project_name FROM projects WHERE created_by = ?", UserName)
		if err != nil {
			fmt.Println("error in running query")
			log.Fatal(err)
		}
		defer rows.Close()
		var Nam []string
		var Name string
		for rows.Next() {
			rows.Scan(&Name)
			Nam = append(Nam, Name)
		}
		setupResponse(&w, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Nam)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (C *Commander) UpdateData(w http.ResponseWriter, r *http.Request) { //update the table with the given Id
	db := database.DbConn()
        defer db.Close()
	if Role == "Program Manager" {
		var dat model.Project
		Time := time.Now()
		var ProjectId int
		var ManagerId int
		var Pid int
		var Mid int
		fmt.Println(UserName)
		json.NewDecoder(r.Body).Decode(&dat)
		ID := dat.Id
		Pn := dat.ProjectName
		Mn := dat.ManagerName
		Email := dat.ManagerEmailID
		//Flag := dat.Flag
		fmt.Println(ID, Pn, Mn, Email)
		getId, _ := db.Query("select project_id, manager_id from manager_details left join projects on manager_details.project_id = projects.id left join manager on manager_details.manager_id = manager.id where manager_details.id = ?", ID)
		defer getId.Close()
		if getId.Next() != false {
			getId.Scan(&ProjectId, &ManagerId)
		}
		getProject, _ := db.Query("SELECT id FROM projects WHERE project_name = ? AND created_by = ?", Pn, UserName)
		defer getProject.Close()
		if getProject.Next() != false {
			getProject.Scan(&Pid)
			fmt.Println(Pid)
			updateProjectId, _ := db.Query("UPDATE manager_details set project_id = ?, updated_at = ? WHERE id = ?", Pid, Time, ID)
                        defer updateProjectId.Close()
		} else {
			updateProject, _ := db.Query("UPDATE projects set project_name = ? WHERE id = ?", Pn, ProjectId)
                        defer updateProject.Close()
			updateProjectTime, _ := db.Query("UPDATE manager_details set updated_at = ? WHERE id = ?", Time, ID)
                        defer updateProjectTime.Close()
		}
		getManager, _ := db.Query("SELECT id FROM manager WHERE manager_name = ? AND manager_email_id = ?", Mn, Email)
		defer getManager.Close()
		if getManager.Next() != false {
			getManager.Scan(&Mid)
			Update, _ := db.Query("UPDATE manager_details set manager_id = ?, updated_at = ? WHERE id = ?", Mid, Time, ID)
                        defer Update.Close()
		} else {
			UpdateManager, _ :=db.Query("UPDATE manager SET manager_name = ?, manager_email_id = ? WHERE id = ?", Mn, Email, ManagerId)
                        defer UpdateManager.Close()
			UpdateManagerDetails, _ := db.Query("UPDATE manager_details set updated_at = ? WHERE id = ?", Time, ID)
                        defer UpdateManagerDetails.Close()
		}
		setupResponse(&w, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (C *Commander) DeleteData(w http.ResponseWriter, r *http.Request) { //delete data
	db := database.DbConn()
        defer db.Close()
	if Role == "Program Manager" {
		var dat model.Project

		json.NewDecoder(r.Body).Decode(&dat)
		del, _ := db.Query("UPDATE manager_details SET flag = 0 WHERE  id = ?", dat.Id)
		defer del.Close()
		setupResponse(&w, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (C *Commander) GetData(w http.ResponseWriter, r *http.Request) { // get all data
	db := database.DbConn()
        defer db.Close()
	if Role == "Program Manager" {
		var Offset int
		Pages := r.FormValue("Pages")
		i1, _ := strconv.Atoi(Pages)
		Offset = 10 * i1
		count, _ := db.Query("SELECT COUNT(Id) FROM manager_details WHERE flag = 1 AND project_id in (SELECT id FROM projects WHERE created_by = ?)", UserName)
		defer count.Close()
		rows, err := db.Query("SELECT projects.project_name, manager.manager_name, manager.manager_email_id, manager_details.flag, manager_details.Id FROM manager_details LEFT JOIN projects ON manager_details.project_id = projects.id LEFT JOIN manager ON manager_details.manager_id = manager.id  WHERE flag = 1 AND created_by = ? LIMIT ?, 10", UserName, Offset)
		if err != nil {
			fmt.Println("error in running query")
			log.Fatal(err)
		}
		defer rows.Close()
		var co int
		var pro model.Project
		var Proj []model.Project
		for rows.Next() {
			rows.Scan(&pro.ProjectName, &pro.ManagerName, &pro.ManagerEmailID, &pro.Flag, &pro.Id)
			Proj = append(Proj, pro)
		}
		if count.Next() != false {
			count.Scan(&co)
		} else {
			co = 0
		}
		var Pag model.Pagenation
		Pag.TotalData = co
		Pag.Limit = 10
		Pag.Data = Proj
		x1 := co / 10
		x := co % 10
		x2 := x1 + 1
		if x == 0 {
			Pag.TotalPages = x1
		} else {
			Pag.TotalPages = x2
		}
		x, _ = strconv.Atoi(Pages)
		if Pag.TotalPages != 0 {
			x1 = x + 1
		}
		Pag.Page = x1
		setupResponse(&w, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Pag)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (C *Commander) GetOpenPositionByAging(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
        defer db.Close()
	if Role == "Program Manager" {
		var j int
		rows, err := db.Query("SELECT  project_name FROM projects WHERE created_by = ?", UserName) //getting all the project name
		if err != nil {
			fmt.Println("error in running query")
			log.Fatal(err)
		}
		defer rows.Close()
		var Nam []model.Position // array of the structure
		var pro model.Position   // instance of the structure
		var names []string // array of all the project names
		var Name string    // instance of the project name
		var Str string
		var position int
		var Pos []int
		var aging int
		var Aging []int
		for rows.Next() {
			rows.Scan(&Name)
			names = append(names, Name)

		}
		for i := 0; i < len(names); i++ {
			N := names[i]
			fmt.Println(names[i], len(names))
			pos, err := db.Query("SELECT created_at, positions FROM open_positions WHERE manager_details_id in (select id from manager_details WHERE project_id IN (SELECT id FROM projects WHERE project_name = ? AND created_by = ?)) AND flag = 1 ", N, UserName)
			if err != nil {
				fmt.Println("error in running query")
				log.Fatal(err)
			}
			defer pos.Close()
			t1 := time.Now()
			t := t1.Format("2006-01-02")
			fmt.Println(t)
			for pos.Next() {
				pos.Scan(&Str, &position)
				fmt.Println(Str, position)
				DataDiff, err := db.Query("SELECT DATEDIFF(?, ?)", t, Str)
				if err != nil {
					fmt.Println("error in running query")
					log.Fatal(err)
				}
				defer DataDiff.Close()
				for DataDiff.Next() {
					DataDiff.Scan(&aging)
				}
				fmt.Println(aging)
				Aging = append(Aging, aging)
				Pos = append(Pos, position)

			}
			count1 := 0
			count2 := 0
			count3 := 0
			count4 := 0
			count5 := 0
			count6 := 0
			for j = 0; j < len(Aging); j++ {

				fmt.Println(Aging[j], Pos[j])

				if Aging[j] < 15 {
					count1 = count1 + Pos[j]

				} else if Aging[j] > 15 && Aging[j] < 30 {
					count2 = count2 + Pos[j]

				} else if Aging[j] > 30 && Aging[j] < 60 {
					count3 = count3 + Pos[j]

				} else if Aging[j] > 60 && Aging[j] < 90 {
					count4 = count4 + Pos[j]

				} else if Aging[j] > 90 && Aging[j] < 120 {
					count5 = count5 + Pos[j]

				} else {
					count6 = count6 + Pos[j]

				}
			}
			Aging = nil
			Pos = nil
			pro.ProjectName = N
			pro.Between0to15 = count1
			pro.Between15to30 = count2
			pro.Between30to60 = count3
			pro.Between60to90 = count4
			pro.Between90to120 = count5
			pro.Greaterthen120 = count6
			Nam = append(Nam, pro)
		}

		setupResponse(&w, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Nam)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
