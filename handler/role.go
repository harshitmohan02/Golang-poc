package handler

import (
	"encoding/json"
	"net/http"
	database "projectname_projectmanager/driver"
	model "projectname_projectmanager/model"
)

func (c *Commander) GetRole(writer http.ResponseWriter, request *http.Request) {
	var Project model.UserRole
	Project.ProjectManager = 0
	Project.ProgramManager = 0
	type EmailData struct {
		Email string `json:"email"`
	}
	var Data EmailData
	json.NewDecoder(request.Body).Decode(&Data)
	db := database.DbConn()
	getProjectManagerID, err := db.Query("SELECT id from project_manager where project_manager_email = ? ", Data.Email)
	defer getProjectManagerID.Close()
	if err != nil {
		panic(err)
	}
	if getProjectManagerID.Next() == true {
		Project.ProjectManager = 1
	}
	getProgramManagerID, err := db.Query("SELECT id from program_manager where program_manager_email = ? ", Data.Email)
	defer getProgramManagerID.Close()
	if err != nil {
		panic(err)
	}
	if getProgramManagerID.Next() == true {
		Project.ProgramManager = 1
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(Project)

}
