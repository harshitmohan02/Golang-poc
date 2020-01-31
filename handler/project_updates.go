package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	database "projectname_projectmanager/driver"
	models "projectname_projectmanager/model"

	"github.com/gorilla/mux"
)

// ProjectUpdatesGetData : To get the project updates
func (C *Commander) ProjectUpdatesGetData(w http.ResponseWriter, r *http.Request) {
	var rol string
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	Rol, err := db.Query("SELECT role FROM token WHERE username = ?", UserName) //Selecting role from the token
	catch(err)
	for Rol.Next() {
		Rol.Scan(&rol)
	}
	defer Rol.Close()

	if (strings.Contains(strings.ToLower(rol), "program manager")) == true { //Converting the role from database to lower case for comparing
		var updates []models.Updates
		result, err1 := db.Query("call getallprojectupdates_Program(?)", UserName)
		catch(err1)
		defer result.Close()
		for result.Next() {
			var update models.Updates
			err2 := result.Scan(&update.ID, &update.ProjectName, &update.Manager, &update.Ups, &update.Downs, &update.ProjectUpdates, &update.GeneralUpdates, &update.Challenges, &update.NeedHelp, &update.ClientVisits, &update.TeamSize, &update.OpenPositions, &update.HighPerformer, &update.LowPerformer)
			catch(err2)
			updates = append(updates, update)
		}
		json.NewEncoder(w).Encode(updates)

	} else if (strings.Contains(strings.ToLower(rol), "project manager")) == true {
		var updates []models.Updates
		result, err1 := db.Query("call getallprojectupdates_Project(?)", UserName)
		catch(err1)
		defer result.Close()
		for result.Next() {
			var update models.Updates
			err2 := result.Scan(&update.ID, &update.ProjectName, &update.Manager, &update.Ups, &update.Downs, &update.ProjectUpdates, &update.GeneralUpdates, &update.Challenges, &update.NeedHelp, &update.ClientVisits, &update.TeamSize, &update.OpenPositions, &update.HighPerformer, &update.LowPerformer)
			catch(err2)
			updates = append(updates, update)
		}
		json.NewEncoder(w).Encode(updates)
	} else {
		respondwithJSON(w, http.StatusNotFound, map[string]string{"Message": "Not Found"})
	}
}

// ProjectUpdatesCreateData : To create project updates
func (C *Commander) ProjectUpdatesCreateData(w http.ResponseWriter, r *http.Request) {
	var rol string
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	Rol, err := db.Query("SELECT role FROM token WHERE username = ?", UserName)
	catch(err)
	for Rol.Next() {
		Rol.Scan(&rol)
	}
	defer Rol.Close()
	if (strings.Contains(strings.ToLower(rol), "project manager")) == true {
		var update models.Updates
		var ManagerID int
		json.NewDecoder(r.Body).Decode(&update)
		stmt1, err1 := db.Query("call projectupdate_ManagerID(?,?)", update.ProjectName, UserName)
		catch(err1)
		defer stmt1.Close()
		if stmt1.Next() != false {
			err2 := stmt1.Scan(&ManagerID)
			catch(err2)
		}
		if ManagerID != 0 {
			var id1 int
			stmt2, err3 := db.Query("SELECT net FROM head_count WHERE sub_project_manager_id = ?", ManagerID) // Getting Team Size from head count table
			catch(err3)
			defer stmt2.Close()
			if stmt2.Next() != false {
				err4 := stmt2.Scan(&update.TeamSize)
				catch(err4)
			}
			stmt3, err5 := db.Query("SELECT sum(position) FROM open_positions WHERE sub_project_manager_id = ?", ManagerID) // Getting open position from open position table
			catch(err5)
			defer stmt3.Close()
			if stmt3.Next() != false {
				err6 := stmt3.Scan(&update.OpenPositions)
				catch(err6)
			}
			stmt4, err7 := db.Query("SELECT id FROM project_updates WHERE sub_project_manager_id = ?", ManagerID)
			catch(err7)
			defer stmt4.Close()
			if stmt4.Next() != false {
				err8 := stmt4.Scan(&id1)
				catch(err8)
			}
			if id1 == 0 {
				stmt, erro := db.Prepare("INSERT INTO project_updates(ups, downs, project_updates, general_updates, challenges, need_help, client_visits, team_size, open_positions, high_performer, low_performer, manager_project_id, created_at, updated_at ) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,now(),now())")
				catch(erro)
				_, erro = stmt.Exec(update.Ups, update.Downs, update.ProjectUpdates, update.GeneralUpdates, update.Challenges, update.NeedHelp, update.ClientVisits, update.TeamSize, update.OpenPositions, update.HighPerformer, update.LowPerformer, ManagerID)
				catch(erro)
				respondwithJSON(w, http.StatusCreated, map[string]string{"Message": "Inserted Successfully"})
			} else {
				respondwithJSON(w, http.StatusConflict, map[string]string{"Message": "Duplicates cannot be Created"})
			}
		} else {
			respondwithJSON(w, http.StatusBadRequest, map[string]string{"Message": "Project not under you"})
		}

	} else {
		respondwithJSON(w, http.StatusUnauthorized, map[string]string{"Message": "Not Authorized"})
	}
}

// ProjectUpdatesGetDataID : to get the data according to search
func (C *Commander) ProjectUpdatesGetDataID(w http.ResponseWriter, r *http.Request) {
	var rol string
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	Rol, err := db.Query("SELECT role FROM token WHERE username = ?", UserName)
	catch(err)
	for Rol.Next() {
		Rol.Scan(&rol)
	}
	defer Rol.Close()

	if (strings.Contains(strings.ToLower(rol), "program manager")) == true {
		var updates []models.Updates
		params := mux.Vars(r)
		key := params["id"]
		var per string = key + "%"
		fmt.Println(per)
		fmt.Println(UserName)
		result, err1 := db.Query("call getallprojectupdates_ProgramID(?,?)", UserName, per)
		catch(err1)
		defer result.Close()
		for result.Next() {
			var update models.Updates
			err2 := result.Scan(&update.ID, &update.ProjectName, &update.Manager, &update.Ups, &update.Downs, &update.ProjectUpdates, &update.GeneralUpdates, &update.Challenges, &update.NeedHelp, &update.ClientVisits, &update.TeamSize, &update.OpenPositions, &update.HighPerformer, &update.LowPerformer)
			catch(err2)
			updates = append(updates, update)
		}
		json.NewEncoder(w).Encode(updates)

	} else if (strings.Contains(strings.ToLower(rol), "project manager")) == true {
		var updates []models.Updates
		params := mux.Vars(r)
		key := params["id"]
		var per string = "'" + key + "%'"
		result, err1 := db.Query("SELECT project_updates.id,projects.project_name,manager.manager_name,ups,downs,project_updates,general_updates,challenges,need_help,client_visits,team_size,open_positions,high_performer,low_performer FROM project_updates LEFT JOIN manager_project ON project_updates.manager_project_id = manager_project.id LEFT JOIN manager ON manager_project.manager_id = manager.id LEFT JOIN projects ON manager_project.project_id = projects.id WHERE manager.manager_email_id = ? AND project_updates.is_active = 1 AND projects.project_name LIKE "+per+" OR manager.manager_name LIKE "+per+" OR ups LIKE "+per+" OR downs LIKE "+per+" OR project_updates LIKE "+per+" OR general_updates LIKE "+per+" OR challenges LIKE "+per+" OR need_help LIKE "+per+" OR client_visits LIKE "+per+" OR team_size LIKE "+per+" OR open_positions LIKE "+per+" OR high_performer LIKE "+per+" OR low_performer LIKE "+per+" ", UserName)
		catch(err1)
		defer result.Close()
		for result.Next() {
			var update models.Updates
			err2 := result.Scan(&update.ID, &update.ProjectName, &update.Manager, &update.Ups, &update.Downs, &update.ProjectUpdates, &update.GeneralUpdates, &update.Challenges, &update.NeedHelp, &update.ClientVisits, &update.TeamSize, &update.OpenPositions, &update.HighPerformer, &update.LowPerformer)
			catch(err2)
			updates = append(updates, update)
		}
		json.NewEncoder(w).Encode(updates)
	} else {
		respondwithJSON(w, http.StatusNotFound, map[string]string{"Message": "Data Not Found"})
	}
}

// ProjectUpdatesUpdateData : To update a particular prospect
func (C *Commander) ProjectUpdatesUpdateData(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")

	var rol string
	Rol, err := db.Query("SELECT role FROM token WHERE username = ?", UserName)
	catch(err)
	for Rol.Next() {
		Rol.Scan(&rol)
	}
	defer Rol.Close()
	if (strings.Contains(strings.ToLower(rol), "project manager")) == true {
		var update models.Updates
		var ManagerID int
		json.NewDecoder(r.Body).Decode(&update)
		stmt1, err1 := db.Query("call projectupdate_ManagerID(?,?)", update.ProjectName, UserName)
		catch(err1)
		defer stmt1.Close()
		if stmt1.Next() != false {
			err2 := stmt1.Scan(&ManagerID)
			catch(err2)
		}

		if ManagerID != 0 {
			var id1 int
			stmt4, err7 := db.Query("SELECT id FROM project_updates WHERE sub_project_manager_id = ?", ManagerID)
			catch(err7)
			defer stmt4.Close()
			if stmt4.Next() != false {
				err8 := stmt4.Scan(&id1)
				catch(err8)
			}
			if id1 != 0 {
				query, err := db.Prepare("UPDATE project_updates SET ups = ?, downs = ?, project_updates = ?, general_updates = ?, challenges = ?, need_help = ?, client_visits = ?, high_performer = ?, low_performer = ? WHERE sub_project_manager_id = ? AND is_active = 1")
				catch(err)
				_, er := query.Exec(update.Ups, update.Downs, update.ProjectUpdates, update.GeneralUpdates, update.Challenges, update.NeedHelp, update.ClientVisits, update.HighPerformer, update.LowPerformer, ManagerID)
				catch(er)
				defer query.Close()
				respondwithJSON(w, http.StatusOK, map[string]string{"message": "Updated Successfully"})
			} else {
				respondwithJSON(w, http.StatusConflict, map[string]string{"Message": "Data Not Found"})
			}
		} else {
			respondwithJSON(w, http.StatusBadRequest, map[string]string{"Message": "Project not under you"})
		}
	} else {
		respondwithJSON(w, http.StatusUnauthorized, map[string]string{"Message": "Not Authorized"})
	}

}

// ProjectUpdatesDeleteData : To delete a particular prospect
func (C *Commander) ProjectUpdatesDeleteData(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	var rol string
	Rol, err := db.Query("SELECT role FROM token WHERE username = ?", UserName)
	catch(err)
	for Rol.Next() {
		Rol.Scan(&rol)
	}
	defer Rol.Close()
	if (strings.Contains(strings.ToLower(rol), "project manager")) == true {
		var update models.Updates
		var ManagerID int
		json.NewDecoder(r.Body).Decode(&update)
		fmt.Println(update.ID)

		stmt1, err1 := db.Query("SELECT sub_project_manager_id FROM project_updates WHERE id = ?", update.ID)
		catch(err1)
		defer stmt1.Close()
		if stmt1.Next() != false {
			err2 := stmt1.Scan(&ManagerID)
			catch(err2)
		}
		fmt.Println(ManagerID)
		if ManagerID != 0 {
			var id1 string
			stmt4, err7 := db.Query("SELECT project_manager_email FROM project_updates LEFT JOIN sub_project_manager ON project_updates.sub_project_manager_id = sub_project_manager.id LEFT JOIN project_manager ON sub_project_manager.manager_id = project_manager.id LEFT JOIN sub_project ON sub_project_manager.sub_project_id = sub_project.id WHERE project_manager.project_manager_email = ? AND project_updates.is_active = 1", UserName)
			catch(err7)
			defer stmt4.Close()
			if stmt4.Next() != false {
				err8 := stmt4.Scan(&id1)
				catch(err8)
			}
			if id1 == UserName {
				stmt, err := db.Prepare("UPDATE update_prospects SET is_active = 0 WHERE sub_project_manager_id = ?")
				catch(err)
				_, err = stmt.Exec(update.ID)
				catch(err)
				defer stmt.Close()
				respondwithJSON(w, http.StatusOK, map[string]string{"message": "Deleted Successfully"})
			} else {
				respondwithJSON(w, http.StatusUnauthorized, map[string]string{"Message": "Not Authorized to Delete"})
			}
		} else {
			respondwithJSON(w, http.StatusBadRequest, map[string]string{"Message": "Error from Front End"})
		}
	} else {
		respondwithJSON(w, http.StatusUnauthorized, map[string]string{"Message": "Not Authorized"})
	}
}
