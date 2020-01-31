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

func (c *Commander) Getallretentions(w http.ResponseWriter, r *http.Request) {
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
		var posts []models.Retention
		var totalretention models.Totalretentions
		var Pag models.Pagination
		offsets, ok := r.URL.Query()["pages"]
		if !ok || len(offsets[0]) < 1 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		pages := offsets[0]
		i, _ := strconv.Atoi(pages)
		offset := i * 10
		result, err := db.Query("select retention.id, projects.project_name, manager.manager_name, retention_initiated, retained from retention left join manager_project on retention.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE projects.program_manager = ? AND retention.is_active = 1 LIMIT ?, 10", UserName, offset)
		if err != nil {
			panic(err.Error())
		}
		for result.Next() {
			var post models.Retention
			err := result.Scan(&post.ID, &post.ProjectName, &post.ProjectManagerName, &post.RetentionInitiated, &post.Retained)
			if err != nil {
				panic(err.Error())
			}

			posts = append(posts, post)

			//t.Data = posts
		}
		defer result.Close()
		count1, err1 := db.Query("SELECT count(id) from retention WHERE is_active = 1 AND manager_project_id in (SELECT id from manager_project where project_id in (select id from projects where program_manager = ?))", UserName)
		if err1 != nil {
			panic(err1.Error())
		}

		defer count1.Close()
		for count1.Next() {
			err2 := count1.Scan(&Pag.TotalData)
			if err2 != nil {
				panic(err2.Error())
			}

		}
		count, err1 := db.Query("select sum(retention_initiated), sum(retained) from retention left join manager_project on retention.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE projects.program_manager = ? AND retention.is_active = 1", UserName)
		if err1 != nil {
			panic(err1.Error())
		}

		defer count.Close()
		for count.Next() {
			err2 := count.Scan(&totalretention.TotalRetentionInitiated, &totalretention.TotalRetained)
			if err2 != nil {
				panic(err2.Error())
			}

		}
		totalretention.Data = posts
		Pag.Data = totalretention
		Pag.Limit = 10
		Pag.TotalPages = (Pag.TotalData / Pag.Limit) + 1
		//Pag.Data = posts
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
		var posts []models.Retention
		var totalretention models.Totalretentions
		var Pag models.Pagination
		offsets, ok := r.URL.Query()["pages"]
		if !ok || len(offsets[0]) < 1 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		pages := offsets[0]
		i, _ := strconv.Atoi(pages)
		offset := i * 10
		result, err := db.Query("select retention.id, projects.project_name, manager.manager_name, retention_initiated, retained from retention left join manager_project on retention.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE manager.manager_email_id = ? AND retention.is_active = 1 LIMIT ?, 10", UserName, offset)
		if err != nil {
			panic(err.Error())
		}
		for result.Next() {
			var post models.Retention
			err := result.Scan(&post.ID, &post.ProjectName, &post.ProjectManagerName, &post.RetentionInitiated, &post.Retained)
			if err != nil {
				panic(err.Error())
			}

			posts = append(posts, post)

			//t.Data = posts
		}

		defer result.Close()

		count, err1 := db.Query("select sum(retention_initiated), sum(retained) from retention left join manager_project on retention.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE manager.manager_email_id = ? AND retention.is_active = 1", UserName)
		if err1 != nil {
			panic(err1.Error())
		}

		defer count.Close()
		for count.Next() {
			err2 := count.Scan(&totalretention.TotalRetentionInitiated, &totalretention.TotalRetained)
			if err2 != nil {
				panic(err2.Error())
			}

		}
		count1, err1 := db.Query("SELECT count(id) from retention WHERE is_active = 1 AND manager_project_id in (SELECT id from manager_project where manager_id in (select id from manager where manager_email_id = ?))", UserName)
		if err1 != nil {
			panic(err1.Error())
		}

		defer count1.Close()
		for count1.Next() {
			err2 := count1.Scan(&Pag.TotalData)
			if err2 != nil {
				panic(err2.Error())
			}

		}
		totalretention.Data = posts
		Pag.Data = totalretention
		Pag.Limit = 10
		Pag.TotalPages = (Pag.TotalData / Pag.Limit) + 1
		//Pag.Data = posts
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
func (c *Commander) Createretentions(w http.ResponseWriter, r *http.Request) {
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

		var post models.Retention
		json.NewDecoder(r.Body).Decode(&post)
		stmt1, err1 := db.Query("select id from manager_project where project_id in ( select id from projects where project_name = ? ) and manager_id in (select id from manager where manager_name = ? )", post.ProjectName, post.ProjectManagerName)
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
		if ManagerID != 0 {
			var dublicateID int
			query := db.QueryRow("SELECT id from retention where manager_project_id = ? AND retention_initiated = ? AND retained = ?", ManagerID, post.RetentionInitiated, post.Retained)
			//catch(err)
			query.Scan(&dublicateID)
			if dublicateID == 0 {

				stmt, err := db.Prepare("INSERT INTO retention(manager_project_id, retention_initiated, retained, created_at, updated_at) VALUES(?, ?, ?, now(), now())")
				catch(err)

				stmt.Exec(ManagerID, post.RetentionInitiated, post.Retained)

				w.WriteHeader(http.StatusCreated)
				fmt.Fprintf(w, "New post was created")
			} else {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Duplicates record found")
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Project not under you")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorised access")
	}

}

func (c *Commander) Deleteretentions(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")

	var post models.Retention
	json.NewDecoder(r.Body).Decode(&post)

	stmt, err := db.Prepare("Update retention set is_active = 0 where id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(post.ID)
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	respondwithJSON(w, http.StatusOK, map[string]string{"message": "deleted successfully"})

}

func (c *Commander) Updateretentions(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")

	var post models.Retention

	json.NewDecoder(r.Body).Decode(&post)

	query, err := db.Prepare("Update retention set retention_initiated = ?, retained = ?, updated_at = ? where id = ?")
	catch(err)
	update := time.Now()
	fmt.Println(update)
	_, er := query.Exec(post.RetentionInitiated, post.Retained, update, post.ID)
	catch(er)

	defer query.Close()

	respondwithJSON(w, http.StatusOK, map[string]string{"message": "update successfully"})

}

func (c *Commander) Getretentionbyprojectname(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")
	var rol string
	Rol, err := db.Query("select role from login where username = ?", UserName)
	catch(err)

	for Rol.Next() {
		Rol.Scan(&rol)
	}
	defer Rol.Close()
	if (strings.Contains(rol, "Program Manager")) || (strings.Contains(rol, "program manager")) == true {
		var posts []models.Retention
		var totalretention models.Totalretentions
		var Pag models.Pagination
		params := mux.Vars(r)
		key := params["projects.project_name"]
		var per string = "'" + key + "%'"
		var Offset int
		Pages := r.FormValue("pages")
		i1, _ := strconv.Atoi(Pages)
		Offset = 10 * i1
		count1, _ := db.Query("SELECT count(retention.id) from retention left join manager_project on retention.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE projects.program_manager = ? AND retention.is_active = 1 AND ((project_name LIKE "+per+") OR (manager.manager_name LIKE "+per+"))", UserName)
		defer count1.Close()

		result, err := db.Query("select retention.id, projects.project_name, manager.manager_name, retention_initiated, retained from retention left join manager_project on retention.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE projects.program_manager = ? AND retention.is_active = 1 AND ((projects.project_name LIKE "+per+") OR (manager.manager_name LIKE "+per+")) LIMIT ?,10", UserName, Offset)

		if err != nil {
			panic(err.Error())
		}
		for result.Next() {
			var post models.Retention
			result.Scan(&post.ID, &post.ProjectName, &post.ProjectManagerName, &post.RetentionInitiated, &post.Retained)

			posts = append(posts, post)

		}
		defer result.Close()

		count, err1 := db.Query("select ifnull(sum(retention_initiated), 0), ifnull(sum(retained), 0) from retention left join manager_project on retention.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE projects.program_manager = ? AND retention.is_active = 1 AND ((projects.project_name LIKE "+per+") OR (manager.manager_name LIKE "+per+"))", UserName)
		if err1 != nil {

			panic(err1.Error())
		}

		for count.Next() {
			err2 := count.Scan(&totalretention.TotalRetentionInitiated, &totalretention.TotalRetained)
			if err2 != nil {
				panic(err2.Error())
			}

		}
		defer count.Close()
		var co int

		if count1.Next() != false {
			count1.Scan(&co)
		} else {
			co = 0
		}
		totalretention.Data = posts
		Pag.Data = totalretention
		Pag.TotalData = co
		Pag.Limit = 10
		//Pag.Data = posts
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
			x1 = (x + 1)
		}
		Pag.Page = x1
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Pag)
	} else if (strings.Contains(rol, "Project Manager")) || (strings.Contains(rol, "project manager")) == true {
		var posts []models.Retention
		var totalretention models.Totalretentions
		var Pag models.Pagination
		params := mux.Vars(r)
		key := params["projects.project_name"]
		var per string = "'" + key + "%'"
		var Offset int
		Pages := r.FormValue("pages")
		i1, _ := strconv.Atoi(Pages)
		Offset = 10 * i1
		count1, _ := db.Query("SELECT count(retention.id) from retention left join manager_project on retention.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE manager.manager_email_id = ? AND retention.is_active = 1 AND ((project_name LIKE "+per+") OR (manager.manager_name LIKE "+per+"))", UserName)
		defer count1.Close()
		result, err := db.Query("select retention.id, projects.project_name, manager.manager_name, retention_initiated, retained from retention left join manager_project on retention.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE manager.manager_email_id = ? AND retention.is_active = 1 AND ((projects.project_name LIKE "+per+") OR (manager.manager_name LIKE "+per+")) LIMIT ?,10", UserName, Offset)

		if err != nil {
			panic(err.Error())
		}
		for result.Next() {
			var post models.Retention
			result.Scan(&post.ID, &post.ProjectName, &post.ProjectManagerName, &post.RetentionInitiated, &post.Retained)

			posts = append(posts, post)
		}
		defer result.Close()

		count, err1 := db.Query("select ifnull(sum(retention_initiated), 0), ifnull(sum(retained), 0) from retention left join manager_project on retention.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE manager.manager_email_id = ? AND retention.is_active = 1 AND ((project_name LIKE "+per+") OR (manager.manager_name LIKE "+per+"))", UserName)
		if err1 != nil {

			panic(err1.Error())
		}

		for count.Next() {
			err2 := count.Scan(&totalretention.TotalRetentionInitiated, &totalretention.TotalRetained)
			if err2 != nil {
				panic(err2.Error())
			}

		}
		defer count.Close()
		var co int

		if count1.Next() != false {
			count1.Scan(&co)
		} else {
			co = 0
		}
		totalretention.Data = posts
		Pag.Data = totalretention
		Pag.TotalData = co
		Pag.Limit = 10
		//Pag.Data = posts
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
			x1 = (x + 1)
		}
		Pag.Page = x1
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Pag)
	} else {

		w.WriteHeader(http.StatusUnauthorized)
	}

}
