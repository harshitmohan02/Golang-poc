package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	database "projectname_projectmanager/driver"
	models "projectname_projectmanager/model"

	"github.com/gorilla/mux"
)

// UpdateProspectGetData : Get all the Prospects of a particular Program Manager or Project Manager
func (C *Commander) UpdateProspectGetData(w http.ResponseWriter, r *http.Request) {
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
		offsets, ok := r.URL.Query()["pages"]
		if !ok || len(offsets[0]) < 1 {
			fmt.Fprintf(w, "Url Param pages is missing")
			return
		}
		pages := offsets[0]
		var prospects []models.Prospects
		var Pag models.Pagination
		Pag.Limit = 10

		result, err1 := db.Query("call getallupdateprospect_Project(?,?)", UserName, pages)
		catch(err1)
		defer result.Close()

		count, err2 := db.Query("call getallupdateprospect_Projectcount(?)", UserName)
		catch(err2)
		defer count.Close()

		for result.Next() {
			var prospect models.Prospects
			erro := result.Scan(&prospect.ID, &prospect.Project, &prospect.Manager, &prospect.Prospect, &prospect.Status, &prospect.Comments, &prospect.Challenges)
			catch(erro)
			prospects = append(prospects, prospect)
		}
		var co int
		if count.Next() != false {
			count.Scan(&co)
		} else {
			co = 0
		}
		Pag.TotalData = co
		Pag.Data = prospects
		x1 := co / 10
		x := co % 10
		x2 := x1 + 1

		if x == 0 {
			Pag.TotalPages = x1
		} else {
			Pag.TotalPages = x2
		}
		x, _ = strconv.Atoi(pages)
		if Pag.TotalPages != 0 {
			x1 = (x + 1)
		}
		Pag.Page = x1
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(Pag)

	} else if (strings.Contains(strings.ToLower(rol), "program manager")) == true {
		offsets, ok := r.URL.Query()["pages"]
		if !ok || len(offsets[0]) < 1 {
			fmt.Fprintf(w, "Url Param pages is missing")
			return
		}
		pages := offsets[0]
		var prospects []models.Prospects
		var Pag models.Pagination
		result, err1 := db.Query("call getallupdateprospect_Program(?,?)", UserName, pages)
		catch(err1)
		defer result.Close()
		Pag.Limit = 10
		count, err2 := db.Query("call getallupdateprospect_Programcount(?)", UserName)
		catch(err2)
		defer count.Close()

		for result.Next() {
			var prospect models.Prospects
			erro := result.Scan(&prospect.ID, &prospect.Project, &prospect.Manager, &prospect.Prospect, &prospect.Status, &prospect.Comments, &prospect.Challenges)
			catch(erro)
			prospects = append(prospects, prospect)
			Pag.Data = prospects
		}
		var co int
		if count.Next() != false {
			count.Scan(&co)
		} else {
			co = 0
		}
		Pag.TotalData = co
		Pag.Data = prospects
		x1 := co / 10
		x := co % 10
		x2 := x1 + 1

		if x == 0 {
			Pag.TotalPages = x1
		} else {
			Pag.TotalPages = x2
		}
		x, _ = strconv.Atoi(pages)
		if Pag.TotalPages != 0 {
			x1 = (x + 1)
		}
		Pag.Page = x1
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(Pag)
	} else {
		fmt.Println("Not")
		w.WriteHeader(http.StatusNotFound)
	}
}

// UpdateProspectCreateData : to insert the data
func (C *Commander) UpdateProspectCreateData(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println(rol, UserName)
	if (strings.Contains(strings.ToLower(rol), "project manager")) == true {
		var prospect models.Prospects
		var ManagerID int
		json.NewDecoder(r.Body).Decode(&prospect)
		stmt1, err1 := db.Query("call updateprospects_getmanagerID(?,?)", prospect.Project, UserName)
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
				stmt, erro := db.Prepare("INSERT INTO update_prospects(sub_project_manager_id,prospect,status,comments,challenges,created_at,updated_at) VALUES(?,?,?,?,?,now(),now())")
				catch(erro)
				_, erro = stmt.Exec(ManagerID, prospect.Prospect, prospect.Status, prospect.Comments, prospect.Challenges)
				catch(erro)
				w.WriteHeader(http.StatusCreated)

				respondwithJSON(w, http.StatusOK, map[string]string{"message": "Inserted Successfully"})
			} else {
				respondwithJSON(w, http.StatusBadRequest, map[string]string{"message": "Duplicates cannot be created"})
			}
		} else {
			respondwithJSON(w, http.StatusBadRequest, map[string]string{"message": "Project not under you"})
		}
	} else {
		respondwithJSON(w, http.StatusUnauthorized, map[string]string{"message": "Not Authorized"})
	}
}

// UpdateProspectGetDataID : to get the data according to search
func (C *Commander) UpdateProspectGetDataID(w http.ResponseWriter, r *http.Request) { // Get all the Prospects of a particular Program Manager or Project Manager
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
		offsets, ok := r.URL.Query()["pages"]
		if !ok || len(offsets[0]) < 1 {
			fmt.Fprintf(w, "Url Parameter 'pages' is missing")
			return
		}
		pages := offsets[0]
		i, _ := strconv.Atoi(pages)
		Offset := i * 10
		var prospects []models.Prospects
		var Pag models.Pagination
		params := mux.Vars(r)
		key := params["id"]
		var per string = "'" + key + "%'"
		Pag.Limit = 10

		result, err1 := db.Query("SELECT update_prospects.id,projects.project_name,manager.manager_name,prospect,status,comments,challenges FROM update_prospects LEFT JOIN manager_project ON update_prospects.manager_project_id = manager_project.id LEFT JOIN manager ON manager_project.manager_id = manager.id LEFT JOIN projects ON manager_project.project_id = projects.id WHERE projects.program_manager = ? AND update_prospects.is_active = 1 AND ((projects.project_name LIKE "+per+") OR (manager.manager_name LIKE "+per+") OR (prospect LIKE "+per+") OR (status LIKE "+per+") OR (comments LIKE "+per+") OR (challenges LIKE "+per+")) LIMIT ?,10", UserName, Offset)
		catch(err1)
		defer result.Close()
		for result.Next() {
			var prospect models.Prospects
			erro := result.Scan(&prospect.ID, &prospect.Project, &prospect.Manager, &prospect.Prospect, &prospect.Status, &prospect.Comments, &prospect.Challenges)
			catch(erro)
			prospects = append(prospects, prospect)
		}
		count, err2 := db.Query("SELECT count(update_prospects.id) FROM update_prospects LEFT JOIN manager_project ON update_prospects.manager_project_id = manager_project.id LEFT JOIN manager ON manager_project.manager_id = manager.id LEFT JOIN projects ON manager_project.project_id = projects.id WHERE projects.program_manager = ? AND update_prospects.is_active = 1 AND ((projects.project_name LIKE "+per+") OR (manager.manager_name LIKE "+per+") OR (prospect LIKE "+per+") OR (status LIKE "+per+") OR (comments LIKE "+per+") OR (challenges LIKE "+per+"))", UserName)
		catch(err2)
		defer count.Close()
		var co int
		if count.Next() != false {
			count.Scan(&co)
		} else {
			co = 0
		}
		Pag.TotalData = co
		Pag.Data = prospects
		x1 := co / 10
		x := co % 10
		x2 := x1 + 1

		if x == 0 {
			Pag.TotalPages = x1
		} else {
			Pag.TotalPages = x2
		}
		x, _ = strconv.Atoi(pages)
		if Pag.TotalPages != 0 {
			x1 = (x + 1)
		}
		Pag.Page = x1
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(Pag)
	} else if (strings.Contains(strings.ToLower(rol), "project manager")) == true {
		offsets, ok := r.URL.Query()["pages"]
		if !ok || len(offsets[0]) < 1 {
			fmt.Fprintf(w, "Url Parameter 'pages' is Missing")
			return
		}
		pages := offsets[0]
		i, _ := strconv.Atoi(pages)
		Offset := i * 10
		var prospects []models.Prospects
		var Pag models.Pagination
		params := mux.Vars(r)
		key := params["id"]
		var per string = "'" + key + "%'"
		Pag.Limit = 10

		result, err1 := db.Query("SELECT update_prospects.id,projects.project_name,manager.manager_name,prospect,status,comments,challenges FROM update_prospects LEFT JOIN manager_project ON update_prospects.manager_project_id = manager_project.id LEFT JOIN manager ON manager_project.manager_id = manager.id LEFT JOIN projects ON manager_project.project_id = projects.id WHERE manager.manager_email_id = ? AND update_prospects.is_active = 1 AND ((projects.project_name LIKE "+per+") OR (manager.manager_name LIKE "+per+") OR (prospect LIKE "+per+") OR (status LIKE "+per+") OR (comments LIKE "+per+") OR (challenges LIKE "+per+")) LIMIT ?,10", UserName, Offset)
		catch(err1)
		defer result.Close()
		Pag.Limit = 10
		count, err2 := db.Query("SELECT count(update_prospects.id) FROM update_prospects LEFT JOIN manager_project ON update_prospects.manager_project_id = manager_project.id LEFT JOIN manager ON manager_project.manager_id = manager.id LEFT JOIN projects ON manager_project.project_id = projects.id WHERE manager.manager_email_id = ? AND update_prospects.is_active = 1 AND ((projects.project_name LIKE "+per+") OR (manager.manager_name LIKE "+per+") OR (prospect LIKE "+per+") OR (status LIKE "+per+") OR (comments LIKE "+per+") OR (challenges LIKE "+per+"))", UserName)
		catch(err2)
		defer count.Close()
		for count.Next() {
			err3 := count.Scan(&Pag.TotalData)
			catch(err3)
		}
		Pag.TotalPages = Pag.TotalData / Pag.Limit
		for result.Next() {
			var prospect models.Prospects
			erro := result.Scan(&prospect.ID, &prospect.Project, &prospect.Manager, &prospect.Prospect, &prospect.Status, &prospect.Comments, &prospect.Challenges)
			catch(erro)
			prospects = append(prospects, prospect)
			Pag.Data = prospects
		}
		var co int
		if count.Next() != false {
			count.Scan(&co)
		} else {
			co = 0
		}
		Pag.TotalData = co
		Pag.Data = prospects
		x1 := co / 10
		x := co % 10
		x2 := x1 + 1

		if x == 0 {
			Pag.TotalPages = x1
		} else {
			Pag.TotalPages = x2
		}
		x, _ = strconv.Atoi(pages)
		if Pag.TotalPages != 0 {
			x1 = (x + 1)
		}
		Pag.Page = x1
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Pag)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// UpdateProspectUpdateData : to update the data
func (C *Commander) UpdateProspectUpdateData(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println(rol, UserName)
	if (strings.Contains(strings.ToLower(rol), "project manager")) == true {
		var prospect models.Prospects
		var ManagerID int
		json.NewDecoder(r.Body).Decode(&prospect)
		stmt1, err1 := db.Query("call updateprospects_getmanagerID(?,?)", prospect.Project, UserName)
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
				query, err := db.Prepare("UPDATE update_prospects SET prospect = ?, status = ?, comments = ?, challenges = ? WHERE sub_project_manager_id = ? AND is_active = 1")
				catch(err)
				_, er := query.Exec(prospect.Prospect, prospect.Status, prospect.Comments, prospect.Challenges, ManagerID)
				catch(er)
				defer query.Close()
				respondwithJSON(w, http.StatusOK, map[string]string{"message": "Updated Successfully"})
			} else {
				respondwithJSON(w, http.StatusBadRequest, map[string]string{"Message": "Data Not Found"})
			}
		} else {
			respondwithJSON(w, http.StatusBadRequest, map[string]string{"Message": "Project not under you"})
		}
	} else {
		respondwithJSON(w, http.StatusUnauthorized, map[string]string{"Message": "Not Authorized"})
	}
}

// UpdateProspectDeleteData : to delete the data
func (C *Commander) UpdateProspectDeleteData(w http.ResponseWriter, r *http.Request) {
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
		var prospect models.Prospects
		var ManagerID int
		json.NewDecoder(r.Body).Decode(&prospect)
		fmt.Println(prospect.ID)

		stmt1, err1 := db.Query("SELECT sub_project_manager_id FROM update_prospects WHERE id = ?", prospect.ID)
		catch(err1)
		defer stmt1.Close()
		if stmt1.Next() != false {
			err2 := stmt1.Scan(&ManagerID)
			catch(err2)
		}
		fmt.Println(ManagerID)

		if ManagerID != 0 {
			var email string
			stmt4, err7 := db.Query("call deleteupdateprospect(?)", ManagerID)
			catch(err7)
			defer stmt4.Close()
			if stmt4.Next() != false {
				err8 := stmt4.Scan(&email)
				catch(err8)
			}
			if email == UserName {
				stmt, err := db.Prepare("UPDATE update_prospects SET is_active = 0 WHERE sub_project_manager_id = ?")
				catch(err)
				_, err = stmt.Exec(prospect.ID)
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
