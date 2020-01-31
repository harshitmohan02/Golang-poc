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

func (c *Commander) Getallheadcount(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")

	var rol string
	Rol, err := db.Query("SELECT role from login where username = ?", UserName)
	//catch(err)
	if err != nil {
		panic(err.Error())
	}

	for Rol.Next() {
		Rol.Scan(&rol)
	}
	defer Rol.Close()
	if (strings.Contains(rol, "program manager")) || (strings.Contains(rol, "Program Manager")) == true {
		var posts []models.HeadCount
		var Pag models.Pagination
		offsets, ok := r.URL.Query()["pages"]

		if !ok || len(offsets[0]) < 1 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		pages := offsets[0]
		i, _ := strconv.Atoi(pages)
		offset := i * 10

		result, err := db.Query("select head_count.id, projects.project_name, manager.manager_name, billables_count, billing_on_hold, vt_count, pi_i_count, others, net from head_count left join manager_project on head_count.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE projects.program_manager = ? AND head_count.is_active = 1 LIMIT ?, 10", UserName, offset)
		if err != nil {
			panic(err.Error())
		}
		defer result.Close()
		Pag.Limit = 10

		count, err1 := db.Query("SELECT count(id) from head_count WHERE is_active = 1 AND manager_project_id in (SELECT id from manager_project where project_id in (select id from projects where program_manager = ?))", UserName)
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
			var post models.HeadCount

			err := result.Scan(&post.ID, &post.ProjectName, &post.ManagerName, &post.BillablesCount, &post.BillingOnHold, &post.VtCount, &post.PiICount, &post.Others, &post.Net)
			//fmt.Println(post.ProjectName)
			if err != nil {
				panic(err.Error())
			}
			// if post.is_active == 1 {
			posts = append(posts, post)
			// fmt.Println(post)
			//}
			//s.Data = posts
		}
		//Pag.Limit = 10
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
		var posts []models.HeadCount
		var Pag models.Pagination
		offsets, ok := r.URL.Query()["pages"]

		if !ok || len(offsets[0]) < 1 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		pages := offsets[0]
		i, _ := strconv.Atoi(pages)
		offset := i * 10
		//fmt.Printf("%T", i)

		result, err := db.Query("select head_count.id, projects.project_name, manager.manager_name, billables_count, billing_on_hold, vt_count, pi_i_count, others, net from head_count left join manager_project on head_count.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE manager.manager_email_id = ? AND head_count.is_active = 1 LIMIT ?, 10", UserName, offset)
		catch(err)

		defer result.Close()
		Pag.Limit = 10

		count, err1 := db.Query("SELECT count(id) from head_count WHERE is_active = 1 AND manager_project_id in (SELECT id from manager_project where manager_id in (select id from manager where manager_email_id = ?))", UserName)
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
			var post models.HeadCount
			err := result.Scan(&post.ID, &post.ProjectName, &post.ManagerName, &post.BillablesCount, &post.BillingOnHold, &post.VtCount, &post.PiICount, &post.Others, &post.Net)
			if err != nil {
				panic(err.Error())
			}

			posts = append(posts, post)

		}
		//s.Data = posts
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
func (c *Commander) Createheadcount(w http.ResponseWriter, r *http.Request) {
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

		var post models.HeadCount
		json.NewDecoder(r.Body).Decode(&post)
		stmt1, err1 := db.Query("select id from manager_project where project_id in ( select id from projects where project_name = ? ) and manager_id in (select id from manager where manager_name = ? )", post.ProjectName, post.ManagerName)
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
			query := db.QueryRow("SELECT id from head_count where manager_project_id = ? AND billables_count = ? AND billing_on_hold = ? AND vt_count = ? AND pi_i_count = ? AND others = ?", ManagerID, post.BillablesCount, post.BillingOnHold, post.VtCount, post.PiICount, post.Others)
			//catch(err)
			query.Scan(&dublicateID)
			if dublicateID == 0 {
				stmt, err := db.Prepare("INSERT INTO head_count(manager_project_id, billables_count, billing_on_hold, vt_count, pi_i_count, others, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, now(), now())")
				if err != nil {
					panic(err.Error())
				}
				json.NewDecoder(r.Body).Decode(&post)

				var Total int
				var Sno int
				_, err = stmt.Exec(ManagerID, post.BillablesCount, post.BillingOnHold, post.VtCount, post.PiICount, post.Others)
				if err != nil {
					panic(err.Error())
				}
				rows, err := db.Query("select id, ifnull(billables_count, 0) + ifnull(billing_on_hold, 0) + ifnull(vt_count, 0) + ifnull(pi_i_count, 0) + ifnull(others, 0) as total from head_count")
				defer rows.Close()
				if err != nil {
					panic(err.Error())
				} else {

					for rows.Next() {
						rows.Scan(&Sno, &Total)
					}

					stmt, err := db.Prepare("update head_count set net = ? where id = ?")
					if err != nil {
						panic(err.Error())
					}

					_, er := stmt.Exec(Total, Sno)

					catch(er)

					if err != nil {
						panic(err.Error())
					}
				}

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

func (c *Commander) Getheadcountbyname(w http.ResponseWriter, r *http.Request) {
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
		var post models.HeadCount
		stmt1, err1 := db.Query("select id from manager_project where project_id in ( select id from projects where project_name = ? ) and manager_id in (select id from manager where manager_name = ? )", post.ProjectName, post.ManagerName)
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
		key := params["projects.project_name"]
		var per string = "'" + key + "%'"
		var Offset int
		Pages := r.FormValue("pages")
		i1, _ := strconv.Atoi(Pages)
		Offset = 10 * i1
		count, _ := db.Query("SELECT count(head_count.id) from head_count left join manager_project on head_count.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE manager.manager_email_id = ? AND head_count.is_active = 1 AND ((project_name LIKE "+per+") OR (manager.manager_name LIKE "+per+"))", UserName)
		defer count.Close()
		result, err := db.Query("select head_count.id, projects.project_name, manager.manager_name, billables_count, billing_on_hold, vt_count, pi_i_count, others, net from head_count left join manager_project on head_count.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE manager.manager_email_id = ? AND head_count.is_active = 1 AND ((project_name LIKE "+per+") OR (manager.manager_name LIKE "+per+")) LIMIT ?,10", UserName, Offset)
		if err != nil {
			panic(err.Error())
		}
		defer result.Close()
		var posts []models.HeadCount
		for result.Next() {
			err := result.Scan(&post.ID, &post.ProjectName, &post.ManagerName, &post.BillablesCount, &post.BillingOnHold, &post.VtCount, &post.PiICount, &post.Others, &post.Net)
			posts = append(posts, post)
			if err != nil {
				panic(err.Error())
			}
		}
		var co int

		if count.Next() != false {
			count.Scan(&co)
		} else {
			co = 0
		}
		// if counter != 0 {
		// 	co = counter
		// } else {
		// 	co = 0
		// }
		var Pag models.Pagination
		Pag.TotalData = co
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
			x1 = (x + 1)
		}
		Pag.Page = x1
		//w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Pag)
		//if post.IsActive == 1 {
		//json.NewEncoder(w).Encode(post)
		//}
	} else if (strings.Contains(rol, "Program Manager")) || (strings.Contains(rol, "program manager")) == true {
		var post models.HeadCount
		stmt1, err1 := db.Query("select id from manager_project where project_id in ( select id from projects where project_name = ? ) and manager_id in (select id from manager where manager_name = ? )", post.ProjectName, post.ManagerName)
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
		key := params["projects.project_name"]
		var per string = "'" + key + "%'"
		var Offset int
		Pages := r.FormValue("pages")
		i1, _ := strconv.Atoi(Pages)
		Offset = 10 * i1
		count, _ := db.Query("SELECT count(head_count.id) from head_count left join manager_project on head_count.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE projects.program_manager = ? AND head_count.is_active = 1 AND ((project_name LIKE "+per+") OR (manager.manager_name LIKE "+per+"))", UserName)
		defer count.Close()
		result, err := db.Query("select head_count.id, projects.project_name, manager.manager_name, billables_count, billing_on_hold, vt_count, pi_i_count, others, net from head_count left join manager_project on head_count.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE projects.program_manager = ? AND head_count.is_active = 1 AND ((project_name LIKE "+per+") OR (manager.manager_name LIKE "+per+")) LIMIT ?,10", UserName, Offset)
		if err != nil {
			panic(err.Error())
		}
		defer result.Close()
		var posts []models.HeadCount
		for result.Next() {
			err := result.Scan(&post.ID, &post.ProjectName, &post.ManagerName, &post.BillablesCount, &post.BillingOnHold, &post.VtCount, &post.PiICount, &post.Others, &post.Net)
			posts = append(posts, post)
			if err != nil {
				panic(err.Error())
			}
		}
		var co int

		if count.Next() != false {
			count.Scan(&co)
		} else {
			co = 0
		}
		var Pag models.Pagination
		Pag.TotalData = co
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
			x1 = (x + 1)
		}
		Pag.Page = x1
		//w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Pag)
		//if post.IsActive == 1 {
		// json.NewEncoder(w).Encode(post)

	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorised access")
	}

}

func (c *Commander) Deleteheadcount(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")

	var post models.HeadCount
	json.NewDecoder(r.Body).Decode(&post)

	stmt, err := db.Prepare("Update head_count set is_active = 0 where id = ?")
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

func (c *Commander) Updateheadcount(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")

	var post models.HeadCount

	json.NewDecoder(r.Body).Decode(&post)

	query, err := db.Prepare("Update head_count set billables_count = ?, billing_on_hold = ?, vt_count = ?, pi_i_count = ?, others = ?, updated_at = ? where id = ?")

	catch(err)
	update := time.Now()
	_, er := query.Exec(post.BillablesCount, post.BillingOnHold, post.VtCount, post.PiICount, post.Others, update, post.ID)
	catch(er)
	fmt.Println(post.BillablesCount, post.BillingOnHold, post.VtCount, post.PiICount, post.Others, post.ID)
	defer query.Close()
	fmt.Println(post.ID)
	rows, err := db.Query("select ifnull(billables_count, 0) + ifnull(billing_on_hold, 0) + ifnull(vt_count, 0) + ifnull(pi_i_count, 0) + ifnull(others, 0) as total from head_count where id = ?", post.ID)

	if err != nil {
		panic(err.Error())
	} else {
		var Total int
		for rows.Next() {
			rows.Scan(&Total)
		}
		fmt.Println(Total)

		stmt, err := db.Query("update head_count set net = ? where id = ?", Total, post.ID)

		if err != nil {
			panic(err.Error())
		}
		//post.Net = Total
		//fmt.Println(post.ID, Total)

		//  stmt.Exec(post.Net, post.ID)
		// catch(er)

		// if err != nil {
		// 	panic(err.Error())
		// }
		for stmt.Next() {
			stmt.Scan(&post.Net, &post.ID)
		}
	}

	respondwithJSON(w, http.StatusOK, map[string]string{"message": "update successfully"})

}
