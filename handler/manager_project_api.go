package handler

import (
	"encoding/json"
	"fmt"

	//"log"
	"net/http"
	database "projectname_projectmanager/driver"
	model "projectname_projectmanager/model"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

//setupResponse : To set all the CORS request
func setupResponse(writer *http.ResponseWriter, request *http.Request) {
	(*writer).Header().Set("Access-Control-Allow-Origin", "*")
	(*writer).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*writer).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// Putdata : insert the data into Manager details table
// func (C *Commander) Putdata(writer http.ResponseWriter, request *http.Request) {
// 	db := database.DbConn()
// 	defer func() {
// 		err := db.Close()
// 		if err != nil {
// 			WriteLogFile(err)
// 			return
// 		}
// 	}()
// 	if strings.Contains(Role, "Program Manager") == true {
// 		var ProjectID int
// 		var ManagerID int
// 		var ManagerDetailData model.Project
// 		Time := time.Now()
// 		json.NewDecoder(request.Body).Decode(&ManagerDetailData)
// 		ProjectName := ManagerDetailData.ProjectName
// 		ManagerName := ManagerDetailData.ManagerName
// 		Email := ManagerDetailData.ManagerEmailID
// 		Flag := 1
// 		getProject, err := db.Query("SELECT id FROM projects WHERE project_name = ?", ProjectName)
// 		if err != nil {
// 			WriteLogFile(err)
// 			return
// 		}
// 		defer func() {
// 			err := getProject.Close()
// 			if err != nil {
// 				WriteLogFile(err)
// 				return
// 			}
// 		}()
// 		if getProject.Next() == false {
// 			InsertProject, err := db.Prepare("INSERT INTO projects(project_name,program_manager)VALUES(?,?)")
// 			if err != nil {
// 				WriteLogFile(err)
// 				return
// 			}
// 			_, err = InsertProject.Exec(ProjectName, UserName)
// 			if err != nil {
// 				WriteLogFile(err)
// 				return
// 			}
// 			defer func() {
// 				err := InsertProject.Close()
// 				if err != nil {
// 					WriteLogFile(err)
// 					return
// 				}
// 			}()
// 		} else {
// 			CheckValadity, err := db.Query("SELECT id FROM projects WHERE project_name = ? AND program_manager = ?", ProjectName, UserName)
// 			if err != nil {
// 				WriteLogFile(err)
// 				return
// 			}
// 			defer func() {
// 				err := CheckValadity.Close()
// 				if err != nil {
// 					WriteLogFile(err)
// 					return
// 				}
// 			}()
// 			if CheckValadity.Next() == false {
// 				writer.WriteHeader(http.StatusPreconditionFailed)
// 				return
// 			}
// 		}
// 		getProjectID, err := db.Query("SELECT id FROM projects WHERE project_name = ?", ProjectName)
// 		if err != nil {
// 			WriteLogFile(err)
// 			return
// 		}
// 		defer func() {
// 			err := getProjectID.Close()
// 			if err != nil {
// 				WriteLogFile(err)
// 				return
// 			}
// 		}()
// 		if getProjectID.Next() != false {
// 			getProjectID.Scan(&ProjectID)
// 			fmt.Println(ProjectID)
// 		}

// 		getManager, err := db.Query("SELECT id FROM manager WHERE manager_name = ? AND manager_email_id = ?", ManagerName, Email)
// 		if err != nil {
// 			WriteLogFile(err)
// 			return
// 		}
// 		defer func() {
// 			err := getManager.Close()
// 			if err != nil {
// 				WriteLogFile(err)
// 				return
// 			}
// 		}()
// 		if getManager.Next() == false {
// 			InsertManager, err := db.Prepare("INSERT INTO manager(manager_name, manager_email_id)VALUES(?, ?)")
// 			if err != nil {
// 				WriteLogFile(err)
// 				return
// 			}
// 			_, err = InsertManager.Exec(ManagerName, Email)
// 			if err != nil {
// 				WriteLogFile(err)
// 				return
// 			}
// 			defer func() {
// 				err := InsertManager.Close()
// 				if err != nil {
// 					WriteLogFile(err)
// 					return
// 				}
// 			}()
// 		}
// 		getManagerID, err := db.Query("SELECT id FROM manager WHERE manager_name = ? AND manager_email_id = ?", ManagerName, Email)
// 		if err != nil {
// 			WriteLogFile(err)
// 			return
// 		}
// 		defer func() {
// 			err := getManagerID.Close()
// 			if err != nil {
// 				WriteLogFile(err)
// 				return
// 			}
// 		}()
// 		if getManagerID.Next() != false {
// 			getManagerID.Scan(&ManagerID)
// 			fmt.Println(ManagerID)
// 		}

// 		CheckAvailability, err := db.Query("SELECT is_active, Id FROM manager_project WHERE project_id = ? AND manager_id = ?", ProjectID, ManagerID)
// 		if err != nil {
// 			WriteLogFile(err)
// 			return
// 		}
// 		defer func() {
// 			err := CheckAvailability.Close()
// 			if err != nil {
// 				WriteLogFile(err)
// 				return
// 			}
// 		}()
// 		if CheckAvailability.Next() == false {
// 			insForm, err := db.Prepare("INSERT INTO manager_project(project_id, manager_id, created_at, updated_at)VALUES(?, ?, ?, ?)")
// 			if err != nil {
// 				WriteLogFile(err)
// 				return
// 			}
// 			_, err = insForm.Exec(ProjectID, ManagerID, Time, Time)
// 			if err != nil {
// 				WriteLogFile(err)
// 				return
// 			}
// 			defer func() {
// 				err := insForm.Close()
// 				if err != nil {
// 					WriteLogFile(err)
// 					return
// 				}
// 			}()
// 			setupResponse(&writer, request)
// 			writer.WriteHeader(http.StatusCreated)
// 		} else {
// 			ID := 0
// 			CheckAvailability.Scan(&Flag, &ID)
// 			fmt.Println(Flag, ID)
// 			if Flag == 0 {
// 				update, err := db.Query("UPDATE manager_project SET is_active = 1, updated_at = ? WHERE  id = ?", Time, ID)
// 				if err != nil {
// 					WriteLogFile(err)
// 					return
// 				}
// 				defer func() {
// 					err := update.Close()
// 					if err != nil {
// 						WriteLogFile(err)
// 						return
// 					}
// 				}()
// 				setupResponse(&writer, request)
// 				writer.WriteHeader(http.StatusCreated)
// 			} else {
// 				writer.WriteHeader(http.StatusBadRequest)
// 			}
// 		}
// 	} else {
// 		writer.WriteHeader(http.StatusNotFound)
// 	}
// }

// GetProjectManagerSearchResult : send all the data with te requested manager name
func (C *Commander) GetProjectManagerSearchResult(w http.ResponseWriter, r *http.Request) {

	db := database.DbConn()
	defer db.Close()
	fmt.Println(Role)
	if strings.Contains(Role, "program manager") == true {
		p := mux.Vars(r)
		key1 := p["id"]
		key := strings.TrimSpace(key1)
		SearchString := key + "%"
		Offset := 0
		fmt.Println(SearchString)
		Pages := r.URL.Query()["Pages"]
		i1, _ := strconv.Atoi(Pages[0])
		Offset = 10 * i1
		Total := 0
		//count, _ := db.Query(" SELECT count(manager_project.id) FROM manager_project LEFT JOIN projects on manager_project.project_id = projects.id LEFT JOIN manager on manager_project.manager_id = manager.id WHERE  manager_project.is_active = 1 AND (manager_name LIKE '"+SearchString+"' OR project_name LIKE '"+SearchString+"' OR manager_email_id LIKE '"+SearchString+"') AND projects.program_manager = ?", UserName)
		//defer count.Close()
		//if count.Next() != false {
		//count.Scan(&Total)
		//} else {
		//Total = 0
		//}
		SearchResults, err := db.Query("call GetManagerDetailsByManagerName(?, ?, ?)", SearchString, UserName, Offset)
		if err != nil {
			WriteLogFile(err)
			return
		}
		defer SearchResults.Close()
		var ManagerDetailData model.Project
		var ManagerDetailsData []model.Project
		for SearchResults.Next() {
			SearchResults.Scan(&ManagerDetailData.ProjectName, &ManagerDetailData.SubProjectName, &ManagerDetailData.ManagerName, &ManagerDetailData.ManagerEmailID, &ManagerDetailData.Id)
			ManagerDetailsData = append(ManagerDetailsData, ManagerDetailData)
			Total++
		}

		var PaginationFormat model.Pagination
		PaginationFormat.TotalData = Total
		PaginationFormat.Limit = 10
		PaginationFormat.Data = ManagerDetailsData
		x1 := Total / 10
		x := Total % 10
		x2 := x1 + 1
		if x == 0 {
			PaginationFormat.TotalPages = x1
		} else {
			PaginationFormat.TotalPages = x2
		}
		x, _ = strconv.Atoi(Pages[0])
		if PaginationFormat.TotalPages != 0 {
			x1 = x + 1
		}
		PaginationFormat.Page = x1
		setupResponse(&w, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PaginationFormat)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// GetProjectName : send all the project name
func (C *Commander) GetProjectName(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	if strings.Contains(Role, "program manager") == true {
		GetProjectName, err := db.Query("SELECT  project_name FROM project WHERE program_manager_id in (SELECT id FROM program_manager WHERE program_manager_email = ?)", UserName)
		if err != nil {
			WriteLogFile(err)
			return
		}
		defer GetProjectName.Close()
		var ProjectNames []string
		var ProjectName string
		for GetProjectName.Next() {
			GetProjectName.Scan(&ProjectName)
			ProjectNames = append(ProjectNames, ProjectName)
		}
		setupResponse(&w, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ProjectNames)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// UpdateData : update the table with the given Id
// func (C *Commander) UpdateData(w http.ResponseWriter, r *http.Request) {
// 	db := database.DbConn()
// 	defer db.Close()
// 	if strings.Contains(Role, "Program Manager") == true {
// 		var ManagerDetail model.Project
// 		Time := time.Now()
// 		var ProjectID int
// 		var ManagerID int
// 		var Pid int
// 		var Mid int
// 		fmt.Println(UserName)
// 		json.NewDecoder(r.Body).Decode(&ManagerDetail)
// 		ID := ManagerDetail.Id
// 		ProjectName := ManagerDetail.ProjectName
// 		ManagerName := ManagerDetail.ManagerName
// 		Email := ManagerDetail.ManagerEmailID
// 		getID, err := db.Query("select project_id, manager_id from manager_project as mp left join projects on mp.project_id = projects.id left join manager on mp.manager_id = manager.id where mp.id = ?", ID)
// 		if err != nil {
// 			WriteLogFile(err)
// 			return
// 		}
// 		defer getID.Close()
// 		if getID.Next() != false {
// 			getID.Scan(&ProjectID, &ManagerID)
// 		}
// 		getProject, _ := db.Query("SELECT id FROM projects WHERE project_name = ? AND program_manager = ?", ProjectName, UserName)
// 		defer getProject.Close()
// 		if getProject.Next() != false {
// 			getProject.Scan(&Pid)
// 			fmt.Println(Pid)
// 			updateProjectID, _ := db.Query("UPDATE manager_project set project_id = ?, updated_at = ? WHERE id = ?", Pid, Time, ID)
// 			defer updateProjectID.Close()
// 		} else {
// 			updateProject, _ := db.Query("UPDATE projects set project_name = ? WHERE id = ?", ProjectName, ProjectID)
// 			defer updateProject.Close()
// 			updateProjectTime, _ := db.Query("UPDATE manager_project set updated_at = ? WHERE id = ?", Time, ID)
// 			defer updateProjectTime.Close()
// 		}
// 		getManager, _ := db.Query("SELECT id FROM manager WHERE manager_name = ? AND manager_email_id = ?", ManagerName, Email)
// 		defer getManager.Close()
// 		if getManager.Next() != false {
// 			getManager.Scan(&Mid)
// 			Update, _ := db.Query("UPDATE manager_project set manager_id = ?, updated_at = ? WHERE id = ?", Mid, Time, ID)
// 			defer Update.Close()
// 		} else {
// 			UpdateManager, _ := db.Query("UPDATE manager SET manager_name = ?, manager_email_id = ? WHERE id = ?", ManagerName, Email, ManagerID)
// 			defer UpdateManager.Close()
// 			UpdateManagerDetails, _ := db.Query("UPDATE manager_project set updated_at = ? WHERE id = ?", Time, ID)
// 			defer UpdateManagerDetails.Close()
// 		}
// 		setupResponse(&w, r)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusCreated)
// 	} else {
// 		w.WriteHeader(http.StatusNotFound)
// 	}
// }

// DeleteData : delete data
// func (C *Commander) DeleteData(w http.ResponseWriter, r *http.Request) {
// 	db := database.DbConn()
// 	defer db.Close()
// 	if strings.Contains(Role, "Program Manager") == true {
// 		var ManagerDetails model.Project

// 		json.NewDecoder(r.Body).Decode(&ManagerDetails)
// 		DeleteManagerDetail, _ := db.Query("UPDATE manager_project SET is_active = 0 WHERE  id = ?", ManagerDetails.Id)
// 		defer DeleteManagerDetail.Close()
// 		setupResponse(&w, r)
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusOK)
// 	} else {
// 		w.WriteHeader(http.StatusNotFound)
// 	}
// }

// GetData : get all data
func (C *Commander) GetData(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	if strings.Contains(Role, "program manager") == true {
		var Offset int
		Pages := r.URL.Query()["Pages"]
		i1, _ := strconv.Atoi(Pages[0])
		Offset = 10 * i1
		count, _ := db.Query("SELECT COUNT(Id) FROM sub_project_manager WHERE sub_project_id in (SELECT id FROM sub_project WHERE project_id in (SELECT id FROM project WHERE program_manager_id in (SELECT id FROM program_manager where program_manager_email = ?)))", UserName)
		defer count.Close()
		GetManagerDetails, err := db.Query("call GetAllManagerDetailsData(?, ?)", UserName, Offset)
		if err != nil {
			WriteLogFile(err)
			return
		}
		defer GetManagerDetails.Close()
		var Total int
		var ManagerDetailData model.Project
		var ManagerDetailsData []model.Project
		for GetManagerDetails.Next() {
			GetManagerDetails.Scan(&ManagerDetailData.ProjectName, &ManagerDetailData.SubProjectName, &ManagerDetailData.ManagerName, &ManagerDetailData.ManagerEmailID, &ManagerDetailData.Id)
			ManagerDetailsData = append(ManagerDetailsData, ManagerDetailData)
		}
		if count.Next() != false {
			count.Scan(&Total)
		} else {
			Total = 0
		}
		var PaginationFormat model.Pagination
		PaginationFormat.TotalData = Total
		PaginationFormat.Limit = 10
		PaginationFormat.Data = ManagerDetailsData
		x1 := Total / 10
		x := Total % 10
		x2 := x1 + 1
		if x == 0 {
			PaginationFormat.TotalPages = x1
		} else {
			PaginationFormat.TotalPages = x2
		}
		x, _ = strconv.Atoi(Pages[0])
		if PaginationFormat.TotalPages != 0 {
			x1 = x + 1
		}
		PaginationFormat.Page = x1
		setupResponse(&w, r)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PaginationFormat)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
