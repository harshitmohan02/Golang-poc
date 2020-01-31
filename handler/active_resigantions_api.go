package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "projectname_projectmanager/driver"
	models "projectname_projectmanager/model"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func (c *Commander) GetAllresignations(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	var rol string
	Rol, err := db.Query("SELECT role from login where username = ?", UserName)
	catch(err)
	for Rol.Next() {
		Rol.Scan(&rol)
	}
	defer Rol.Close()
	if (strings.Contains(rol, "program manager")) || (strings.Contains(rol, "Program Manager")) == true {
		var posts []models.Resignations
		var Pag models.Pagination
		offsets, ok := r.URL.Query()["pages"]
		if !ok || len(offsets[0]) < 1 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		pages := offsets[0]
		i, _ := strconv.Atoi(pages)
		offset := i * 10
		result, err := db.Query("call active_resign_get_all_program(?, ?)", UserName, offset)
		catch(err)
		defer result.Close()
		Pag.Limit = 10
		count, err1 := db.Query("SELECT count(id) from active_resignations WHERE is_active = 1 AND manager_project_id in (SELECT id from manager_project where project_id in (select id from projects where program_manager = ?))", UserName)
		if err1 != nil {
			panic(err1.Error())
		}
		defer count.Close()
		for count.Next() {
			err2 := count.Scan(&Pag.TotalData)
			if err2 != nil {
				panic(err2.Error())
			}
		}
		Pag.TotalPages = (Pag.TotalData / Pag.Limit) + 1
		for result.Next() {
			var post models.Resignations
			err := result.Scan(&post.ID, &post.Empname, &post.Project, &post.Manager, &post.Backfillrequired, &post.Regrenonregre, &post.Status, &post.Dateofresignation, &post.Dateofleaving)
			if err != nil {
				panic(err.Error())
			}
			posts = append(posts, post)
		}
		Pag.Data = posts
		x1 := Pag.TotalData / 10
		x := Pag.TotalData % 10
		x2 := x1 + 1

		if x == 0 {
			Pag.TotalPages = x1
		} else {
			Pag.TotalPages = x2
		}
		x, _ = strconv.Atoi(pages)
		if Pag.TotalPages != 0 {
			x1 = x + 1
		}
		Pag.Page = x1
		json.NewEncoder(w).Encode(Pag)
	} else if (strings.Contains(rol, "project manager")) || (strings.Contains(rol, "Project Manager")) == true {
		var posts []models.Resignations
		var Pag models.Pagination
		offsets, ok := r.URL.Query()["pages"]
		if !ok || len(offsets[0]) < 1 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		pages := offsets[0]
		i, _ := strconv.Atoi(pages)
		offset := i * 10
		result, err := db.Query("call active_resign_get_all_project(?, ?)", UserName, offset)
		catch(err)
		defer result.Close()
		Pag.Limit = 10
		count, err1 := db.Query("SELECT count(id) from active_resignations WHERE is_active = 1 AND manager_project_id in (SELECT id from manager_project where manager_id in (select id from manager where manager_email_id = ?))", UserName)
		if err1 != nil {
			panic(err1.Error())
		}
		defer count.Close()
		for count.Next() {
			err2 := count.Scan(&Pag.TotalData)
			if err2 != nil {
				panic(err2.Error())
			}
		}
		Pag.TotalPages = (Pag.TotalData / Pag.Limit) + 1
		for result.Next() {
			var post models.Resignations
			err := result.Scan(&post.ID, &post.Empname, &post.Project, &post.Manager, &post.Backfillrequired, &post.Regrenonregre, &post.Status, &post.Dateofresignation, &post.Dateofleaving)
			if err != nil {
				panic(err.Error())
			}
			posts = append(posts, post)
		}
		defer db.Close()
		Pag.Data = posts
		x1 := Pag.TotalData / 10
		x := Pag.TotalData % 10
		x2 := x1 + 1

		if x == 0 {
			Pag.TotalPages = x1
		} else {
			Pag.TotalPages = x2
		}
		x, _ = strconv.Atoi(pages)
		if Pag.TotalPages != 0 {
			x1 = x + 1
		}
		Pag.Page = x1
		json.NewEncoder(w).Encode(Pag)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
func (c *Commander) CreateResignations(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	var rol string
	var ManagerID int
	Rol, err := db.Query("select role from login where username = ?", UserName)
	catch(err)
	for Rol.Next() {
		Rol.Scan(&rol)
	}
	defer Rol.Close()
	if (strings.Contains(rol, "Project Manager")) || (strings.Contains(rol, "project manager")) == true {
		var post models.Resignations
		fmt.Println(rol)
		json.NewDecoder(r.Body).Decode(&post)
		fmt.Println(post.Project, post.Manager)
		stmt1, err1 := db.Query("select id from manager_project where project_id in ( select id from projects where project_name = ? ) and manager_id in (select id from manager where manager_name = ? )", post.Project, post.Manager)

		if err1 != nil {
			panic(err1.Error())
		}
		defer stmt1.Close()
		if stmt1.Next() != false {
			err2 := stmt1.Scan(&ManagerID)

			if err2 != nil {
				panic(err2.Error())
			}
		}
		//fmt.Println(ManagerID)
		if ManagerID != 0 {
			stmt, err := db.Prepare("INSERT INTO active_resignations(emp_name, manager_project_id, backfill_required, regre_non_regre, status, date_of_resignation, date_of_leaving, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, now(), now())")
			if err != nil {
				panic(err.Error())
			}

			_, err = stmt.Exec(post.Empname, ManagerID, post.Backfillrequired, post.Regrenonregre, post.Status, post.Dateofresignation, post.Dateofleaving)
			if err != nil {
				panic(err.Error())
			}
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "New post was created")
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Project not under you")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorised access")
	}
}
func (c *Commander) GetResignationsbyName(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	var rol string
	var ManagerID int
	Rol, err := db.Query("select role from login where username = ?", UserName)
	catch(err)
	for Rol.Next() {
		Rol.Scan(&rol)
	}
	defer Rol.Close()
	if (strings.Contains(rol, "Project Manager")) || (strings.Contains(rol, "project manager")) == true {
		var post models.Resignations
		stmt1, err1 := db.Query("select id from manager_project where project_id in ( select id from projects where project_name = ? ) and manager_id in (select id from manager where manager_name = ? )", post.Project, post.Manager)
		if err1 != nil {
			panic(err1.Error())
		}
		defer stmt1.Close()
		if stmt1.Next() != false {
			err2 := stmt1.Scan(&ManagerID)
			if err2 != nil {
				panic(err2.Error())
			}
		}
		params := mux.Vars(r)
		key := params["emp_name"]
		var per string = key + "%"
		var Offset int
		var co int //for number of records
		Pages := r.FormValue("pages")
		i1, _ := strconv.Atoi(Pages)
		Offset = 10 * i1
		count, _ := db.Query("call active_resign_get_by_name_project_count(?, ?)", UserName, per)
		defer count.Close()
		fmt.Println(per)
		result, err := db.Query("call active_resign_get_by_name_project_result(?, ?, ?)", UserName, per, Offset)
		if err != nil {
			panic(err.Error())
		}
		defer result.Close()

		var posts []models.Resignations
		for result.Next() {
			err := result.Scan(&post.ID, &post.Empname, &post.Project, &post.Manager, &post.Backfillrequired, &post.Regrenonregre, &post.Status, &post.Dateofresignation, &post.Dateofleaving)
			posts = append(posts, post)
			if err != nil {
				panic(err.Error())
			}
		}
		if count.Next() != false {
			count.Scan(&co)
		} else {
			co = 0
		}
		defer db.Close()
		var Pag models.Pagination
		Pag.TotalData = co
		fmt.Println(co)
		Pag.Limit = 10
		Pag.Data = posts
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Pag)
	} else if (strings.Contains(rol, "Program Manager")) || (strings.Contains(rol, "program manager")) == true {
		var post models.Resignations
		stmt1, err1 := db.Query("select id from manager_project where project_id in ( select id from projects where project_name = ? ) and manager_id in (select id from manager where manager_name = ? )", post.Project, post.Manager)
		if err1 != nil {
			panic(err1.Error())
		}
		defer stmt1.Close()
		if stmt1.Next() != false {
			err2 := stmt1.Scan(&ManagerID)
			if err2 != nil {
				panic(err2.Error())
			}
		}
		params := mux.Vars(r)
		key := params["emp_name"]
		var per string = "'" + key + "%'"
		var Offset int
		Pages := r.FormValue("pages")
		i1, _ := strconv.Atoi(Pages)
		Offset = 10 * i1
		count, _ := db.Query("SELECT count(active_resignations.id) from active_resignations left join manager_project on active_resignations.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE projects.program_manager  = ? AND active_resignations.is_active = 1 AND ((emp_name LIKE "+per+") OR (projects.project_name LIKE "+per+") OR (manager.manager_name LIKE "+per+") or (date_of_resignation LIKE "+per+") OR (date_of_leaving LIKE "+per+"))", UserName)
		defer count.Close()
		fmt.Println(per)
		result, err := db.Query("select active_resignations.id, emp_name, projects.project_name, manager.manager_name, backfill_required, regre_non_regre, status, date_of_resignation, date_of_leaving from active_resignations left join manager_project on active_resignations.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE projects.program_manager = ? AND active_resignations.is_active = 1 AND ((emp_name LIKE "+per+") OR (projects.project_name LIKE "+per+") OR (manager.manager_name LIKE "+per+") OR (date_of_resignation LIKE "+per+") OR (date_of_leaving LIKE "+per+")) LIMIT ?,10", UserName, Offset)
		if err != nil {
			panic(err.Error())
		}
		defer result.Close()
		var co int
		var posts []models.Resignations
		for result.Next() {
			err := result.Scan(&post.ID, &post.Empname, &post.Project, &post.Manager, &post.Backfillrequired, &post.Regrenonregre, &post.Status, &post.Dateofresignation, &post.Dateofleaving)
			posts = append(posts, post)
			if err != nil {
				panic(err.Error())
			}
		}
		if count.Next() != false {
			count.Scan(&co)
		} else {
			co = 0
		}
		defer db.Close()
		var Pag models.Pagination
		fmt.Println(co)
		Pag.TotalData = co
		Pag.Limit = 10
		Pag.Data = posts
		x1 := co / 10
		x := int(co) % 10
		x2 := x1 + 1
		if x == 0 {
			Pag.TotalPages = x1
		} else {
			Pag.TotalPages = x2
		}
		x, _ = strconv.Atoi(Pages)
		if Pag.TotalPages != 0 {
			x1 = (x + 1)
		}
		Pag.Page = x1
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Pag)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
func (c *Commander) DeleteResignations(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	var rol string
	Rol, err := db.Query("select role from login where username = ?", UserName)
	catch(err)
	for Rol.Next() {
		Rol.Scan(&rol)
	}
	defer Rol.Close()
	if (strings.Contains(rol, "Project Manager")) || (strings.Contains(rol, "project manager")) == true {
		var post models.Resignations
		json.NewDecoder(r.Body).Decode(&post)
		stmt, err := db.Prepare("Update active_resignations set is_active = 0 where id = ?")
		if err != nil {
			panic(err.Error())
		}
		_, err = stmt.Exec(post.ID)
		if err != nil {
			panic(err.Error())
		}
		defer stmt.Close()
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "deleted successfully"})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
func (c *Commander) UpdateResignations(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	var rol string
	Rol, err := db.Query("select role from login where username = ?", UserName)
	catch(err)
	for Rol.Next() {
		Rol.Scan(&rol)
	}
	defer Rol.Close()
	if (strings.Contains(rol, "Project Manager")) || (strings.Contains(rol, "project manager")) == true {
		var post models.Resignations
		json.NewDecoder(r.Body).Decode(&post)
		query, err := db.Prepare("Update active_resignations set emp_name = ?, backfill_required = ?, regre_non_regre = ?, status = ?, date_of_resignation = ?, date_of_leaving = ?, updated_at = ? where id = ?")
		catch(err)
		update := time.Now()
		_, er := query.Exec(post.Empname, post.Backfillrequired, post.Regrenonregre, post.Status, post.Dateofresignation, post.Dateofleaving, update, post.ID)
		catch(er)
		defer query.Close()
		respondwithJSON(w, http.StatusOK, map[string]string{"message": "update successfully"})
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}
func catch(err error) {
	if err != nil {
		panic(err)
	}
}
func respondwithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	fmt.Println(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
