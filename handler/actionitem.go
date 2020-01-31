package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	database "projectname_projectmanager/driver"
	model "projectname_projectmanager/model"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" //blank import
	"github.com/gorilla/mux"
)

var field string = "id,action_item,manager_details_id,meeting_date,target_date,status,closed_date,comment,is_active"
var fields string = "action_items.id,action_items.action_item,action_items.meeting_date,action_items.target_date,action_items.status,action_items.closed_date,action_items.comment,action_items.is_active,projects.project_name,manager.manager_name"

//ActionItemPostData : to post the data into the database
func (c *Commander) ActionItemPostData(w http.ResponseWriter, r *http.Request) {
	var error model.Error
	var actionitem model.ActionItemClosed
	var ManagerDetailsID int
	SetupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
		return
	}
	db := database.DbConn()
	defer db.Close()
	err := json.NewDecoder(r.Body).Decode(&actionitem)
	BadRequest(w, err)
	if actionitem.Status != "open" && actionitem.Status != "inprogress" && actionitem.Status != "closed" && actionitem.Status != "onhold" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		//error.Code = "1265"
		error.Message = "Data truncated for column 'status'. Invalid entry for column 'status'"
		json.NewEncoder(w).Encode(error)
		return
	}
	myDateString := "2006-01-02"
	myMeetingDate, err := time.Parse(myDateString, actionitem.MeetingDate)
	if err != nil {
		WriteLogFile(err)
		panic(err)
	}
	myTargetDate, err := time.Parse(myDateString, actionitem.TargetDate)
	if err != nil {
		WriteLogFile(err)
		panic(err)
	}
	isTargetDateBeforeMeetingDate := myTargetDate.Before(myMeetingDate)
	if isTargetDateBeforeMeetingDate == false {
		var ActionItem, MeetingDate, TargetDate, Status, Comment string
		selDB, err := db.Query("SELECT id FROM manager_project WHERE project_id = (SELECT id FROM projects WHERE project_name=? )and manager_id = (SELECT id FROM manager WHERE manager_email_id =?)", actionitem.ProjectName, UserName)
		defer selDB.Close()
		BadRequest(w, err)
		if selDB.Next() != false {
			err := selDB.Scan(&ManagerDetailsID)
			BadRequest(w, err)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		selDB, err = db.Query("SELECT action_item,meeting_date,target_date,status,comment from action_items WHERE manager_project_id=? AND action_item=? AND is_active='1'", ManagerDetailsID, actionitem.ActionItem)
		defer selDB.Close()
		BadRequest(w, err)
		for selDB.Next() {
			err := selDB.Scan(&ActionItem, &MeetingDate, &TargetDate, &Status, &Comment)
			BadRequest(w, err)
			if ActionItem == actionitem.ActionItem &&
				MeetingDate == actionitem.MeetingDate &&
				TargetDate == actionitem.TargetDate &&
				Status == actionitem.Status &&
				Comment == actionitem.Comment {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				//error.Code = "2627"
				error.Message = "Violation of Unique constraint: Duplicate row insertion"
				json.NewEncoder(w).Encode(error)
				return
			}
		}
		BadRequest(w, err)
		CreatedAt := time.Now()
		insForm, err := db.Prepare("INSERT INTO action_items(action_item,manager_project_id,meeting_date,target_date,status,comment,created_at,updated_at)VALUES(?,?,?,?,?,?,?,?)")
		defer insForm.Close()
		if err != nil {
			WriteLogFile(err)
			return
		}
		insForm.Exec(actionitem.ActionItem,
			ManagerDetailsID,
			actionitem.MeetingDate,
			actionitem.TargetDate,
			actionitem.Status,
			actionitem.Comment,
			CreatedAt,
			CreatedAt)
		BadRequest(w, err)
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		//error.Code = "1292"
		if isTargetDateBeforeMeetingDate == true {
			error.Message = "Incorrect Target Date Value"
		}
		json.NewEncoder(w).Encode(error)
	}
}

//ActionItemUpdateData : to update the details stored in database
func (c *Commander) ActionItemUpdateData(w http.ResponseWriter, r *http.Request) {
	var error model.Error
	var actionitem model.ActionItemClosed
	var closedInTime, managerDetailsID int
	SetupResponse(&w, r)
	fmt.Println("11")
	if (*r).Method == "OPTIONS" {
		fmt.Println("Hello how are you?")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
	}
	fmt.Println("222")
	db := database.DbConn()
	defer db.Close()
	err := json.NewDecoder(r.Body).Decode(&actionitem)
	BadRequest(w, err)
	myDateString := "2006-01-02"
	myMeetingDate, err := time.Parse(myDateString, actionitem.MeetingDate)
	if err != nil {
		WriteLogFile(err)
		panic(err)
	}
	myTargetDate, err := time.Parse(myDateString, actionitem.TargetDate)
	if err != nil {
		WriteLogFile(err)
		panic(err)
	}
	isTargetDateBeforeMeetingDate := myTargetDate.Before(myMeetingDate)
	if actionitem.ClosedDate != "" {
		if actionitem.Status == "closed" {
			myClosedDate, err := time.Parse(myDateString, actionitem.ClosedDate)
			if err != nil {
				WriteLogFile(err)
				panic(err)
			}
			isClosedDateBeforeMeetingDate := myClosedDate.Before(myMeetingDate)
			if isClosedDateBeforeMeetingDate == false && isTargetDateBeforeMeetingDate == false {
				isClosedDateAfterTargetDate := myClosedDate.After(myTargetDate)
				if isClosedDateAfterTargetDate == false {
					closedInTime = 1
				} else {
					closedInTime = 0
				}
				selDB, err := db.Query("SELECT id FROM manager_project WHERE project_id = (SELECT id FROM projects WHERE project_name=? )and manager_id = (SELECT id FROM manager WHERE manager_name =?)", actionitem.ProjectName, actionitem.Owner)
				defer selDB.Close()
				BadRequest(w, err)
				if selDB.Next() != false {
					err := selDB.Scan(&managerDetailsID)
					BadRequest(w, err)
				} else {
					w.WriteHeader(http.StatusBadRequest)
				}

				updatedAt := time.Now()
				updForm, err := db.Prepare("UPDATE action_items SET action_item=? , manager_project_id=? , meeting_date=? , target_date=? , status=? , closed_date=? , comment=? , updated_at=? , closed_in_time=? WHERE id=?")
				if err != nil {
					WriteLogFile(err)
					panic(err.Error())
				}
				defer updForm.Close()
				updForm.Exec(actionitem.ActionItem,
					managerDetailsID,
					actionitem.MeetingDate,
					actionitem.TargetDate,
					actionitem.Status,
					actionitem.ClosedDate,
					actionitem.Comment,
					updatedAt,
					closedInTime,
					actionitem.SNo)
				defer updForm.Close()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				//error.Code = "1292"
				if isClosedDateBeforeMeetingDate == true && isTargetDateBeforeMeetingDate == true {
					error.Message = "Incorrect Closed Date and Target Date Values"
				} else if isTargetDateBeforeMeetingDate == true {
					error.Message = "Incorrect Target Date Value"
				} else if isClosedDateBeforeMeetingDate == true {
					error.Message = "Incorrect Closed Date Value"
				}
				json.NewEncoder(w).Encode(error)
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			//error.Code = "1293"
			error.Message = "Closed date cannot be filled if the status is not closed"
			json.NewEncoder(w).Encode(error)
		}
	} else {
		if isTargetDateBeforeMeetingDate == false {
			selDB, err := db.Query("SELECT id FROM manager_project WHERE project_id = (SELECT id FROM projects WHERE project_name=? )and manager_id = (SELECT id FROM manager WHERE manager_name =?)", actionitem.ProjectName, actionitem.Owner)
			defer selDB.Close()
			BadRequest(w, err)
			if selDB.Next() != false {
				err := selDB.Scan(&managerDetailsID)
				BadRequest(w, err)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
			updatedAt := time.Now()
			updForm, err := db.Prepare("UPDATE action_items SET action_item=? , manager_project_id=? , meeting_date=? , target_date=? , status=? , comment=? , updated_at=? , closed_in_time = NULL WHERE id=?")
			if err != nil {
				WriteLogFile(err)
				panic(err.Error())
			}
			updForm.Exec(actionitem.ActionItem,
				managerDetailsID,
				actionitem.MeetingDate,
				actionitem.TargetDate,
				actionitem.Status,
				actionitem.Comment,
				updatedAt,
				actionitem.SNo)
			defer updForm.Close()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			//error.Code = "1292"
			if isTargetDateBeforeMeetingDate == true {
				error.Message = "Incorrect Target Date Value"
				json.NewEncoder(w).Encode(error)
			}
		}
	}
}

//ActionItemDeleteData : to soft delete the details
func (c *Commander) ActionItemDeleteData(w http.ResponseWriter, r *http.Request) {
	var actionitem model.ActionItemClosed
	SetupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
		return
	}
	db := database.DbConn()
	defer db.Close()
	err := json.NewDecoder(r.Body).Decode(&actionitem)
	BadRequest(w, err)
	updatedAt := time.Now()
	updForm, err := db.Prepare("UPDATE action_items SET is_active = '0', updated_at = ? WHERE id=?")
	if err != nil {
		WriteLogFile(err)
		panic(err.Error())
	}
	defer updForm.Close()
	updForm.Exec(updatedAt, actionitem.SNo)
	w.WriteHeader(http.StatusCreated)
}

//ActionItemGetData : to the get the action items
func (c *Commander) ActionItemGetData(w http.ResponseWriter, r *http.Request) {
	var data []model.ActionItemClosed
	var Page model.Pagination
	var countt int
	var error model.Error
	SetupResponse(&w, r)
	fmt.Println("11")
	if (*r).Method == "OPTIONS" {
		fmt.Println("Hello how are you?")
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
		return
	}
	fmt.Println("222")
	cacheKey := "actionclosed"
	e := `"` + cacheKey + `"`
	w.Header().Set("Etag", e)
	w.Header().Set("Cache-Control", "max-age=2592000") // 30 days

	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, e) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}
	db := database.DbConn()
	defer db.Close()
	//pages := r.FormValue("pages")
	pages := r.URL.Query().Get("pages")
	page, err := strconv.Atoi(pages)
	Page.Page = page + 1
	BadRequest(w, err)
	offset := page * 10
	status := r.URL.Query().Get("status")
	if strings.Contains(Role, "project manager") == true {
		if status == "closed" {
			var SNo int
			var ProjectName, ActionItem, Owner, MeetingDate, TargetDate, Status, ClosedDate, Comment, Flag string
			selDB, err := db.Query("call getClosedActionsProject(?,?)", UserName, offset)
			BadRequest(w, err)
			defer selDB.Close()
			for selDB.Next() {
				countt++
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName, &Owner)
				BadRequest(w, err)
				data = append(data, model.ActionItemClosed{SNo: SNo, ProjectName: ProjectName, ActionItem: ActionItem, Owner: Owner, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: ClosedDate, Comment: Comment, Flag: Flag})
			}
		} else if status == "open" {
			var SNo int
			var ProjectName, ActionItem, Owner, MeetingDate, TargetDate, Status, Comment, Flag string
			var ClosedDate sql.NullString
			selDB, err := db.Query("call getAllActionsProject(?,?)", UserName, offset)
			BadRequest(w, err)
			defer selDB.Close()
			for selDB.Next() {
				countt++
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName, &Owner)
				BadRequest(w, err)
				data = append(data, model.ActionItemClosed{SNo: SNo, ProjectName: ProjectName, ActionItem: ActionItem, Owner: Owner, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: "NA", Comment: Comment, Flag: Flag})
			}
		} else {
			error.Code = "405"
			error.Message = "Method not allowed"
			json.NewEncoder(w).Encode(error)
		}
	} else {
		if status == "closed" {
			var SNo int
			var ProjectName, ActionItem, Owner, MeetingDate, TargetDate, Status, ClosedDate, Comment, Flag string
			selDB, err := db.Query("call getClosedActionsProgram(?,?)", UserName, offset)
			defer selDB.Close()
			BadRequest(w, err)
			for selDB.Next() {
				countt++
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName, &Owner)
				BadRequest(w, err)
				data = append(data, model.ActionItemClosed{SNo: SNo, ProjectName: ProjectName, ActionItem: ActionItem, Owner: Owner, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: ClosedDate, Comment: Comment, Flag: Flag})
			}
		} else if status == "open" {
			var SNo int
			var ProjectName, ActionItem, Owner, MeetingDate, TargetDate, Status, Comment, Flag string
			var ClosedDate sql.NullString
			selDB, err := db.Query("call getAllActionsProgram(?,?)", UserName, offset)
			BadRequest(w, err)
			defer selDB.Close()
			for selDB.Next() {
				countt++
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName, &Owner)
				BadRequest(w, err)
				data = append(data, model.ActionItemClosed{SNo: SNo, ProjectName: ProjectName, ActionItem: ActionItem, Owner: Owner, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: "NA", Comment: Comment, Flag: Flag})
			}
		} else {
			error.Code = "405"
			error.Message = "Method not allowed"
			json.NewEncoder(w).Encode(error)
		}

	}
	w.Header().Set("Content-Type", "application/json")
	Page.TotalData = countt
	Page.Limit = 10
	x := countt
	page = x / 10
	x = x % 10
	if x == 0 {
		Page.TotalPages = page
	} else {
		Page.TotalPages = page + 1
	}
	Page.Data = data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Page)
}

//ActionItemGetDataID : to search the action items
func (c *Commander) ActionItemGetDataID(w http.ResponseWriter, r *http.Request) {
	var data []model.ActionItemClosed
	var Page model.Pagination
	var countt int
	var error model.Error
	SetupResponse(&w, r)
	if (*r).Method == "OPTIONS" {
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.WriteHeader(http.StatusOK)
		return
	}
	db := database.DbConn()
	defer db.Close()
	cacheKey := "actionclosedID"
	e := `"` + cacheKey + `"`
	w.Header().Set("Etag", e)
	w.Header().Set("Cache-Control", "max-age=2592000") // 30 days

	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, e) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}
	//pages := r.FormValue("pages")
	pages := r.URL.Query().Get("pages")
	page, err := strconv.Atoi(pages)
	Page.Page = page + 1
	BadRequest(w, err)
	offset := page * 10
	status := r.URL.Query().Get("status")
	p := mux.Vars(r)
	removeWhiteSpace := p["id"]
	removeWhiteSpace = strings.TrimSpace(removeWhiteSpace)
	key := removeWhiteSpace + "%"
	if strings.Contains(Role, "project manager") == true {
		if status == "closed" {
			var SNo int
			var ProjectName, ActionItem, Owner, MeetingDate, TargetDate, Status, ClosedDate, Comment, Flag string
			selDB, err := db.Query("call getClosedActionsProjectQuery(?,?,?)", key, UserName, offset)
			BadRequest(w, err)
			defer selDB.Close()
			for selDB.Next() {
				countt++
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName, &Owner)
				BadRequest(w, err)
				data = append(data, model.ActionItemClosed{SNo: SNo, ProjectName: ProjectName, ActionItem: ActionItem, Owner: Owner, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: ClosedDate, Comment: Comment, Flag: Flag})
			}
		} else if status == "open" {
			var SNo int
			var ProjectName, ActionItem, Owner, MeetingDate, TargetDate, Status, Comment, Flag string
			var ClosedDate sql.NullString
			selDB, err := db.Query("call getAllActionsProjectQuery(?,?,?)", key, UserName, offset)
			BadRequest(w, err)
			defer selDB.Close()
			for selDB.Next() {
				countt++
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName, &Owner)
				BadRequest(w, err)
				data = append(data, model.ActionItemClosed{SNo: SNo, ProjectName: ProjectName, ActionItem: ActionItem, Owner: Owner, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: "NA", Comment: Comment, Flag: Flag})
			}
		} else {
			error.Code = "405"
			error.Message = "Method not allowed"
			json.NewEncoder(w).Encode(error)
		}
	} else {
		if status == "closed" {
			var SNo int
			var ProjectName, ActionItem, Owner, MeetingDate, TargetDate, Status, ClosedDate, Comment, Flag string
			selDB, err := db.Query("call getClosedActionsProgramQuery(?,?,?)", key, UserName, offset)
			BadRequest(w, err)
			defer selDB.Close()
			for selDB.Next() {
				countt++
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName, &Owner)
				BadRequest(w, err)
				data = append(data, model.ActionItemClosed{SNo: SNo, ProjectName: ProjectName, ActionItem: ActionItem, Owner: Owner, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: ClosedDate, Comment: Comment, Flag: Flag})
			}
		} else if status == "open" {
			var SNo int
			var ProjectName, ActionItem, Owner, MeetingDate, TargetDate, Status, Comment, Flag string
			var ClosedDate sql.NullString
			selDB, err := db.Query("call getAllActionsProgramQuery(?,?,?)", key, UserName, offset)
			BadRequest(w, err)
			defer selDB.Close()
			for selDB.Next() {
				countt++
				err := selDB.Scan(&SNo, &ActionItem, &MeetingDate, &TargetDate, &Status, &ClosedDate, &Comment, &Flag, &ProjectName, &Owner)
				BadRequest(w, err)
				data = append(data, model.ActionItemClosed{SNo: SNo, ProjectName: ProjectName, ActionItem: ActionItem, Owner: Owner, MeetingDate: MeetingDate, TargetDate: TargetDate, Status: Status, ClosedDate: "NA", Comment: Comment, Flag: Flag})
			}
		} else {
			error.Code = "405"
			error.Message = "Method not allowed"
			json.NewEncoder(w).Encode(error)
		}
	}
	Page.TotalData = countt
	Page.Limit = 10
	x := countt
	page = x / 10
	x = x % 10
	if x == 0 {
		Page.TotalPages = page
	} else {
		Page.TotalPages = page + 1
	}
	Page.Data = data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Page)
}
