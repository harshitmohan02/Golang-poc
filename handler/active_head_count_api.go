package handler

import (
	"encoding/json"
	"net/http"
	database "projectname_projectmanager/driver"
	models "projectname_projectmanager/model"
	"strings"
)

func (C *Commander) Getactiveheadcount(w http.ResponseWriter, r *http.Request) {
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
		result, err := db.Query("SELECT sum(net), sum(billables_count), sum(billing_on_hold), sum(vt_count), sum(pi_i_count), sum(others) from head_count left join manager_project on head_count.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE projects.program_manager = ? AND head_count.is_active = 1", UserName)
		if err != nil {
			panic(err.Error())
		}
		defer result.Close()
		var post models.Active
		for result.Next() {
			err := result.Scan(&post.ActiveHeadCount, &post.Billable, &post.BillingOnHold, &post.ValueTrade, &post.ProjectInvestment, &post.Others)
			if err != nil {
				panic(err.Error())
			}
		}
		json.NewEncoder(w).Encode(post)
	} else if (strings.Contains(rol, "project manager")) || (strings.Contains(rol, "Project Manager")) == true {
		result, err := db.Query("SELECT sum(net), sum(billables_count), sum(billing_on_hold), sum(vt_count), sum(pi_i_count), sum(others) from head_count left join manager_project on head_count.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE manager.manager_email_id = ? AND head_count.is_active = 1", UserName)
		if err != nil {
			panic(err.Error())
		}
		defer result.Close()
		var post models.Active
		for result.Next() {
			err := result.Scan(&post.ActiveHeadCount, &post.Billable, &post.BillingOnHold, &post.ValueTrade, &post.ProjectInvestment, &post.Others)
			if err != nil {
				panic(err.Error())
			}
		}
		json.NewEncoder(w).Encode(post)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
