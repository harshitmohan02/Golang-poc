package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	database "projectname_projectmanager/driver"
	models "projectname_projectmanager/model"
	"strings"
	"time"
)

func (c *Commander) Getcompleteresignations(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")

	var post models.Activeresignationforyear
	var posts []models.Activeresignationforyear
	var inf []models.Info
	var infs models.Info
	var count [6]int
	var mname string

	var testname string = "abc"
	var ca int = 0
	var date string
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

		wek, err := db.Query("select MONTHNAME(date_of_resignation), date_format(date_of_resignation, '%Y-%m-%d')  from active_resignations left join manager_project on active_resignations.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE projects.program_manager = ? AND active_resignations.is_active = 1 AND date_of_resignation>=date_sub(now(), interval 06 MONTH) order by MONTH(date_of_resignation)", UserName)
		if err != nil {
			panic(err.Error())
		}
		defer wek.Close()
		for wek.Next() {
			wek.Scan(&mname, &date)

			layout := "2006-01-02"
			str, _ := time.Parse(layout, date)
			fmt.Println(mname, date)
			weekNo := NumberOfTheWeekInMonth(str)
			if testname == mname || testname == "abc" {
				count[weekNo]++
				testname = mname
				ca++
			} else {
				for j := 1; j < 6; j++ {
					post.CountNo = count[j]
					post.Week = j
					posts = append(posts, post)

				}
				infs.Month = testname
				infs.Total = ca
				infs.Data = posts
				inf = append(inf, infs)
				ca = 1
				posts = nil
				count[1] = 0
				count[2] = 0
				count[3] = 0
				count[4] = 0
				count[5] = 0
				count[weekNo]++
				testname = mname
			}
			//}
		}
		for j := 1; j < 6; j++ {
			post.CountNo = count[j]
			post.Week = j
			posts = append(posts, post)

		}
		infs.Month = testname
		infs.Total = ca
		ca = 1
		infs.Data = posts
		inf = append(inf, infs)
		posts = nil

		if infs.Month != "abc" {

			json.NewEncoder(w).Encode(inf)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else if (strings.Contains(rol, "project manager")) || (strings.Contains(rol, "Project Manager")) == true {
		wek, err := db.Query("select MONTHNAME(date_of_resignation), date_format(date_of_resignation, '%Y-%m-%d')  from active_resignations left join manager_project on active_resignations.manager_project_id = manager_project.id left join manager on  manager_project.manager_id = manager.id left join projects on manager_project.project_id = projects.id WHERE manager.manager_email_id = ? AND active_resignations.is_active = 1 AND date_of_resignation>=date_sub(now(), interval 06 MONTH) order by MONTH(date_of_resignation)", UserName)
		if err != nil {
			panic(err.Error())
		}
		defer wek.Close()
		for wek.Next() {
			wek.Scan(&mname, &date)

			layout := "2006-01-02"
			str, _ := time.Parse(layout, date)
			fmt.Println(mname, date)
			weekNo := NumberOfTheWeekInMonth(str)
			if testname == mname || testname == "abc" {
				count[weekNo]++
				testname = mname
				ca++
			} else {
				for j := 1; j < 6; j++ {
					post.CountNo = count[j]
					post.Week = j
					posts = append(posts, post)

				}
				infs.Month = testname
				infs.Total = ca
				infs.Data = posts
				inf = append(inf, infs)
				ca = 1
				posts = nil
				count[1] = 0
				count[2] = 0
				count[3] = 0
				count[4] = 0
				count[5] = 0
				count[weekNo]++
				testname = mname
			}
			//}
		}
		for j := 1; j < 6; j++ {
			post.CountNo = count[j]
			post.Week = j
			posts = append(posts, post)

		}
		infs.Month = testname
		infs.Total = ca
		ca = 1
		infs.Data = posts
		inf = append(inf, infs)
		posts = nil

		if infs.Month != "abc" {

			json.NewEncoder(w).Encode(inf)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

}
func NumberOfTheWeekInMonth(now time.Time) int {
	beginningOfTheMonth := time.Date(now.Year(), now.Month(), 1, 1, 1, 1, 1, time.UTC)
	_, thisWeek := now.ISOWeek()
	_, beginningWeek := beginningOfTheMonth.ISOWeek()
	return 1 + thisWeek - beginningWeek
}
