package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "projectname_projectmanager/driver"
	models "projectname_projectmanager/model"
	"strings"
)

func (c *Commander) Getprojectheadcount(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	var rol string
	Rol, err := db.Query("SELECT role from login where username = ?", UserName)
	if err != nil {
		panic(err.Error())
	}
	for Rol.Next() {
		Rol.Scan(&rol)
	}
	defer Rol.Close()
	if (strings.Contains(rol, "program manager")) || (strings.Contains(rol, "Program Manager")) == true {
		var posts []models.Bottom
		result, err := db.Query("select projects.project_name, sum(net) from head_count left join manager_project on head_count.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE projects.program_manager = ? AND head_count.is_active = 1 group by projects.project_name", UserName)
		if err != nil {
			panic(err.Error())
		}
		defer result.Close()
		for result.Next() {
			var post models.Bottom
			err := result.Scan(&post.Project, &post.HeadCount)
			if err != nil {
				fmt.Println("in the error")
				panic(err.Error())
			}
			posts = append(posts, post)
		}
		json.NewEncoder(w).Encode(posts)
	} else if (strings.Contains(rol, "project manager")) || (strings.Contains(rol, "Project Manager")) == true {
		var posts []models.Bottom
		result, err := db.Query("select projects.project_name, sum(net) from head_count left join manager_project on head_count.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE manager.manager_email_id = ? AND head_count.is_active = 1 group by projects.project_name", UserName)
		if err != nil {
			panic(err.Error())
		}
		defer result.Close()
		for result.Next() {
			var post models.Bottom
			err := result.Scan(&post.Project, &post.HeadCount)
			if err != nil {
				fmt.Println("in the error")
				panic(err.Error())
			}
			posts = append(posts, post)
		}
		json.NewEncoder(w).Encode(posts)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
