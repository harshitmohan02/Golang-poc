package handler

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	database "projectname_projectmanager/driver"
	models "projectname_projectmanager/model"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Getallopenpositions :Func for fetching all entries
func (c *Commander) Getallopenpositions(w http.ResponseWriter, r *http.Request) {

	db := database.DbConn()
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")

	if strings.Contains(Role, "Project Manager") || strings.Contains(Role, "project manager") {
		statusGet, ok := r.URL.Query()["status"]

		if !ok || len(statusGet[0]) < 1 {
			fmt.Fprintf(w, "URL Param status is missing")
			return
		}
		status := statusGet[0]

		if status == "daily" || status == "Daily" {
			offsetGet, ok := r.URL.Query()["pages"]

			if !ok || len(offsetGet[0]) < 1 {
				fmt.Fprintf(w, "Url Param offset is missing")
				return
			}

			off := offsetGet[0]
			offInt, _ := strconv.Atoi(off)
			offset := offInt * 10

			limitGet, ok := r.URL.Query()["limit"]

			if !ok || len(limitGet[0]) < 1 {
				fmt.Fprintf(w, "Url Param limit is missing")
				return
			}
			limit := limitGet[0]
			limitFloat, _ := strconv.ParseFloat(limit, 64)
			fmt.Println(limitFloat)
			var pageLimit models.Pages_Daily
			var dailyArray []models.Daily

			queryDaily, err := db.Query("select o.id,p.sub_project_name, o.designation,o.type_position,o.position, o.priority, o.additional_comment, o.l1_due,o.l2_due, o.client_due, o.created_at,o.is_active from open_positions o inner join sub_project_manager s on o.sub_project_manager_id = s.id inner join sub_project p on s.sub_project_id = p.id where s.manager_id in (select id from project_manager where project_manager_email = ?) LIMIT ?, ? ", UserName, offset, limitFloat)
			catch(err)
			defer queryDaily.Close()

			count, err := db.Query("select count(id) from open_positions where is_active = 1 and sub_project_manager_id in ( select id from sub_project_manager where manager_id in ( select id from project_manager where project_manager_email = ?))", UserName)
			catch(err)
			defer count.Close()

			for count.Next() {
				err2 := count.Scan(&pageLimit.Total_data)
				catch(err2)

			}
			pageLimit.Limit = limitFloat
			pageLimit.Page = math.Ceil(pageLimit.Total_data / pageLimit.Limit)

			for queryDaily.Next() {

				var daily models.Daily
				err := queryDaily.Scan(&daily.Id, &daily.Project_name, &daily.Designation, &daily.Type_position, &daily.Position, &daily.Priority, &daily.Additonal_comment, &daily.L1_due, &daily.L2_due, &daily.Client_due, &daily.Created_at, &daily.Is_active)
				catch(err)

				location, err := time.LoadLocation("Asia/Kolkata")
				catch(err)

				layout := "2006-01-02 15:04:05"
				t, err := time.Parse(layout, daily.Created_at)
				to := time.Now()
				to = to.In(location)
				days := to.Sub(t) / (24 * time.Hour)
				daily.Ageing = int(days)

				if daily.Is_active == "1" {
					dailyArray = append(dailyArray, daily)
				}
				pageLimit.Data = dailyArray

			}
			json.NewEncoder(w).Encode(pageLimit)

		} else if status == "Weekly" || status == "weekly" {

			offsets, ok := r.URL.Query()["pages"]
			if !ok || len(offsets[0]) < 1 {
				fmt.Fprintf(w, "Url Param offset is missing")
				return
			}

			off := offsets[0]
			offInt, _ := strconv.Atoi(off)
			offset := offInt * 10

			limitGet, ok := r.URL.Query()["limit"]

			if !ok || len(limitGet[0]) < 1 {
				fmt.Fprintf(w, "Url Param limit is missing")
				return
			}

			limit := limitGet[0]
			limitFloat, _ := strconv.ParseFloat(limit, 64)

			var pageLimit models.Pages_Weekly
			var weeklyArray []models.Weekly

			queryWeekly, err := db.Query("select o.id,p.sub_project_name, o.designation,o.type_position,o.position, o.priority, o.additional_comment, o.l1_happened,o.l2_happened, o.client_happened, o.is_active, o.created_at from open_positions o inner join sub_project_manager s on o.sub_project_manager_id = s.id inner join sub_project p on s.sub_project_id = p.id where s.manager_id in (select id from project_manager where project_manager_email = ?) LIMIT ?, ?  ", UserName, offset, limitFloat)
			catch(err)
			defer queryWeekly.Close()

			count, err1 := db.Query("select count(id) from open_positions where is_active = 1 and sub_project_manager_id in ( select id from sub_project_manager where manager_id in ( select id from project_manager where project_manager_email = ?))", UserName)
			catch(err1)
			defer count.Close()

			for count.Next() {
				err2 := count.Scan(&pageLimit.Total_data)
				catch(err2)
			}

			pageLimit.Limit = limitFloat
			pageLimit.Page = math.Ceil(pageLimit.Total_data / pageLimit.Limit)

			for queryWeekly.Next() {
				var weekly models.Weekly
				var createdAt string

				err := queryWeekly.Scan(&weekly.Id, &weekly.Project_name, &weekly.Designation, &weekly.Type_position, &weekly.Position, &weekly.Priority, &weekly.Additonal_comment, &weekly.L1_Happened, &weekly.L2_Happened, &weekly.Client_Happened, &weekly.Is_active, &createdAt)
				catch(err)

				location, err := time.LoadLocation("Asia/Kolkata")
				catch(err)

				layout := "2006-01-02 15:04:05"
				t, err := time.Parse(layout, createdAt)
				to := time.Now()
				to = to.In(location)
				days := to.Sub(t) / (24 * time.Hour)
				weekly.Ageing = int(days)

				if weekly.Is_active == "1" {
					weeklyArray = append(weeklyArray, weekly)
				}
				pageLimit.Data = weeklyArray
			}
			json.NewEncoder(w).Encode(pageLimit)
		} else {
			fmt.Fprintf(w, "No status")
		}

	} else if strings.Contains(Role, "Program Manager") || strings.Contains(Role, "program manager") {
		st, ok1 := r.URL.Query()["status"]

		if !ok1 || len(st[0]) < 1 {
			fmt.Fprintf(w, "status is missing")
			return
		}
		status := st[0]

		if status == "daily" || status == "Daily" {
			offsets, ok := r.URL.Query()["pages"]

			if !ok || len(offsets[0]) < 1 {
				fmt.Fprintf(w, "Url Param offset is missing")
				return
			}

			off := offsets[0]
			i, _ := strconv.Atoi(off)
			offset := i * 10

			limitGet, ok := r.URL.Query()["limit"]

			if !ok || len(limitGet[0]) < 1 {
				fmt.Fprintf(w, "Url Param limit is missing")
				return
			}

			limit := limitGet[0]
			limitFloat, _ := strconv.ParseFloat(limit, 64)

			var s models.Pages_Daily
			var posts []models.Daily
			result, err := db.Query("select o.id,p.project_name, o.designation,o.type_position,o.position, o.priority, o.additional_comment, o.l1_due,o.l2_due, o.client_due, o.created_at,o.is_active from open_positions o inner join sub_project_manager s on o.sub_project_manager_id = s.id inner join sub_project sp on s.sub_project_id = sp.id inner join project p on sp.project_id = p.id where p.program_manager_id in (select id from program_manager where program_manager_email = ?) LIMIT ?, ?", UserName, offset, limitFloat)
			catch(err)
			defer result.Close()
			count, err1 := db.Query("select count(id) from open_positions where is_active = 1 and sub_project_manager_id in (select id from sub_project_manager where sub_project_id in ( select id from sub_project where project_id in ( select id from project where program_manager_id in ( select id from program_manager where program_manager_email = ? ))))", UserName)
			catch(err1)

			defer count.Close()
			for count.Next() {
				err2 := count.Scan(&s.Total_data)
				catch(err2)

			}
			s.Limit = limitFloat
			s.Page = math.Ceil(s.Total_data / s.Limit)
			for result.Next() {
				var post models.Daily

				err := result.Scan(&post.Id, &post.Project_name, &post.Designation, &post.Type_position, &post.Position, &post.Priority, &post.Additonal_comment, &post.L1_due, &post.L2_due, &post.Client_due, &post.Created_at, &post.Is_active)
				catch(err)
				location, err := time.LoadLocation("Asia/Kolkata")
				catch(err)
				layout := "2006-01-02 15:04:05"
				t, err := time.Parse(layout, post.Created_at)
				to := time.Now()
				to = to.In(location)
				days := to.Sub(t) / (24 * time.Hour)
				//fmt.Println("days", int(days))

				post.Ageing = int(days)

				if post.Is_active == "1" {
					posts = append(posts, post)
				}
				s.Data = posts

			}
			json.NewEncoder(w).Encode(s)
		} else if status == "Weekly" || status == "weekly" {
			offsets, ok := r.URL.Query()["pages"]

			if !ok || len(offsets[0]) < 1 {
				fmt.Fprintf(w, "Url Param offset is missing")
				return
			}

			off := offsets[0]
			i, _ := strconv.Atoi(off)
			offset := i * 10

			limitGet, ok := r.URL.Query()["limit"]

			if !ok || len(limitGet[0]) < 1 {
				fmt.Fprintf(w, "Url Param limit is missing")
				return
			}

			limit := limitGet[0]
			limitFloat, _ := strconv.ParseFloat(limit, 64)

			var s models.Pages_Weekly
			var posts []models.Weekly
			result, err := db.Query("select o.id,p.project_name, o.designation,o.type_position,o.position, o.priority, o.additional_comment, o.l1_happened,o.l2_happened, o.client_happened, o.is_active, o.created_at from open_positions o inner join sub_project_manager s on o.sub_project_manager_id = s.id inner join sub_project sp on s.sub_project_id = sp.id inner join project p on sp.project_id = p.id where p.program_manager_id in (select id from program_manager where program_manager_email= ?) LIMIT ?, ? ", UserName, offset, limitFloat)
			catch(err)
			defer result.Close()
			count, err1 := db.Query("select count(id) from open_positions where is_active = 1 and sub_project_manager_id in (select id from sub_project_manager where sub_project_id in ( select id from sub_project where project_id in ( select id from project where program_manager_id in ( select id from program_manager where program_manager_email = ? ))))", UserName)
			catch(err1)

			defer count.Close()
			for count.Next() {
				err2 := count.Scan(&s.Total_data)
				catch(err2)

			}
			s.Limit = limitFloat
			s.Page = math.Ceil(s.Total_data / s.Limit)
			for result.Next() {
				var post models.Weekly
				var c string

				err := result.Scan(&post.Id, &post.Project_name, &post.Designation, &post.Type_position, &post.Position, &post.Priority, &post.Additonal_comment, &post.L1_Happened, &post.L2_Happened, &post.Client_Happened, &post.Is_active, &c)
				catch(err)
				location, err := time.LoadLocation("Asia/Kolkata")
				catch(err)
				layout := "2006-01-02 15:04:05"
				t, err := time.Parse(layout, c)
				to := time.Now()
				to = to.In(location)
				days := to.Sub(t) / (24 * time.Hour)
				//fmt.Println("days", int(days))

				post.Ageing = int(days)

				if post.Is_active == "1" {
					posts = append(posts, post)
				}
				s.Data = posts

			}
			json.NewEncoder(w).Encode(s)
		}
	} else {
		fmt.Fprintf(w, "No authentication")
	}
}

// Createopenpositions :Func for creating new entries
func (c *Commander) Createopenpositions(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")

	if strings.Contains(Role, "Project Manager") || strings.Contains(Role, "project manager") {
		st, ok1 := r.URL.Query()["status"]

		if !ok1 || len(st[0]) < 1 {
			w.WriteHeader(http.StatusCreated)
			return
		}
		status := st[0]

		if status == "daily" || status == "Daily" {
			var post models.Daily
			json.NewDecoder(r.Body).Decode(&post)

			q2 := db.QueryRow("select id from sub_project_manager where sub_project_id in ( select id from sub_project where sub_project_name = ? ) and manager_id in ( select id from project_manager where project_manager_email = ? )", post.Project_name, UserName)
			var managerid int
			q2.Scan(&managerid)
			//fmt.Println(managerid)
			if managerid != 0 {
				var id1 int
				q1 := db.QueryRow("select id from open_positions where sub_project_manager_id = ? and designation = ? and type_position = ?  ", managerid, post.Designation, post.Type_position)
				q1.Scan(&id1)
				fmt.Println(id1)
				fmt.Println(managerid, post.Designation, post.Type_position)

				if id1 == 0 {

					stmt, err := db.Prepare("INSERT INTO open_positions( sub_project_manager_id, designation, type_position, position, priority, additional_comment, l1_due, l2_due, client_due, created_at, updated_at) VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, now(), now())")
					catch(err)
					_, err = stmt.Exec(managerid, post.Designation, post.Type_position, post.Position, post.Priority, post.Additonal_comment, post.L1_due, post.L2_due, post.Client_due)
					catch(err)
					w.WriteHeader(http.StatusCreated)
					fmt.Fprintf(w, "New post was created")
				} else {
					fmt.Fprintf(w, "Duplicates cannot be created")
				}
			} else {
				fmt.Fprintf(w, "Project not under you")
			}
		} else if status == "Weekly" || status == "weekly" {
			var post models.Weekly
			json.NewDecoder(r.Body).Decode(&post)

			q2 := db.QueryRow("select id from sub_project_manager where sub_project_id in ( select id from sub_project where sub_project_name = ? ) and manager_id in ( select id from project_manager where project_manager_email = ? )", post.Project_name, UserName)
			var managerid int
			q2.Scan(&managerid)
			fmt.Println(managerid)
			if managerid != 0 {
				var id1 int
				q1 := db.QueryRow("select id from open_positions where sub_project_manager_id = ? and designation = ? and type_position = ?  ", managerid, post.Designation, post.Type_position)
				q1.Scan(&id1)
				fmt.Println(id1)
				fmt.Println(managerid, post.Designation, post.Type_position)

				if id1 != 0 {
					stmt, err := db.Prepare("UPDATE open_positions set l1_happened = ?, l2_happened=?, client_happened= ?,  updated_at =  now() where id = ? ")
					catch(err)

					_, err = stmt.Exec(post.L1_Happened, post.L2_Happened, post.Client_Happened, id1)
					catch(err)
					w.WriteHeader(http.StatusCreated)
					fmt.Fprintf(w, "Fields are updated")
				} else {
					fmt.Fprintf(w, "Not found the specific data")
				}
			} else {
				fmt.Fprintf(w, "Project not under you")
			}
		} else {
			fmt.Fprintf(w, "No status")
		}
	} else {
		fmt.Fprintf(w, "Authentication Failed")
	}

}

// Getopenpositions :Func for fetching specific entries
func (c *Commander) Getopenpositions(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")

	if strings.Contains(Role, "Project Manager") || strings.Contains(Role, "project manager") {
		params := mux.Vars(r)
		st, ok1 := r.URL.Query()["status"]

		if !ok1 || len(st[0]) < 1 {
			fmt.Fprintf(w, "Url Param status is missing")
			return
		}
		status := st[0]

		if status == "daily" || status == "Daily" {

			offsets, ok := r.URL.Query()["pages"]

			if !ok || len(offsets[0]) < 1 {
				fmt.Fprintf(w, "Url Param offset is missing")
				return
			}

			off := offsets[0]
			i, _ := strconv.Atoi(off)
			fmt.Printf("%T", i)

			offset := i * 10
			var s models.Pages_Daily
			var posts []models.Daily
			var st string
			st = "'" + params["id"] + "%'"
			result, err := db.Query("select o.id,p.project_name, o.designation,o.type_position,o.position, o.priority, o.additional_comment, o.l1_due,o.l2_due, o.client_due, o.created_at,o.is_active from open_positions o inner join manager_project m on o.sub_project_manager_id = m.id inner join projects p on p.id=m.project_id where o.is_active=1 and designation like "+st+" OR type_position like "+st+" OR position like "+st+" OR priority like "+st+" OR additional_comment like "+st+" OR l1_due like "+st+" OR l2_due like "+st+" OR client_due like "+st+" and o. sub_project_manager_id in (select id from manager_project where manager_id in (select id from manager where manager_email_id = ? )) LIMIT ?, 10 ", UserName, offset)
			catch(err)
			defer result.Close()
			count, err1 := db.Query("select count(id) from open_positions where is_active = 1 and sub_project_manager_id in ( select id from manager_project where manager_id in ( select id from manager where manager_email_id = ?))", UserName)
			catch(err1)

			defer count.Close()
			for count.Next() {
				err2 := count.Scan(&s.Total_data)
				catch(err2)

			}
			s.Limit = 10
			s.Page = math.Ceil(s.Total_data / s.Limit)
			for result.Next() {
				var post models.Daily

				err := result.Scan(&post.Id, &post.Project_name, &post.Designation, &post.Type_position, &post.Position, &post.Priority, &post.Additonal_comment, &post.L1_due, &post.L2_due, &post.Client_due, &post.Created_at, &post.Is_active)
				catch(err)
				location, err := time.LoadLocation("Asia/Kolkata")
				catch(err)
				layout := "2006-01-02 15:04:05"
				t, err := time.Parse(layout, post.Created_at)
				to := time.Now()
				to = to.In(location)
				days := to.Sub(t) / (24 * time.Hour)
				fmt.Println("days", int(days))
				fmt.Println(post.Created_at)

				post.Ageing = int(days)

				if post.Is_active == "1" {
					posts = append(posts, post)
				}
				s.Data = posts
			}
			json.NewEncoder(w).Encode(s)
		} else if status == "Weekly" || status == "weekly" {
			offsets, ok := r.URL.Query()["pages"]

			if !ok || len(offsets[0]) < 1 {
				fmt.Fprintf(w, "Url Param offset is missing")
				return
			}

			off := offsets[0]
			i, _ := strconv.Atoi(off)
			fmt.Printf("%T", i)

			offset := i * 10
			var s models.Pages_Weekly
			var posts []models.Weekly
			var st string
			st = "'" + params["id"] + "%'"
			result, err := db.Query("select o.id,p.project_name, o.designation,o.type_position,o.position, o.priority, o.additional_comment, o.l1_happened,o.l2_happened, o.client_happened, o.is_active, o.created_at from open_positions o inner join manager_project m on o.sub_project_manager_id = m.id inner join projects p on p.id=m.project_id where o.is_active=1 and designation like "+st+" OR type_position like "+st+" OR position like "+st+" OR priority like "+st+" OR additional_comment like "+st+" OR l1_happened like "+st+" OR l2_happened like "+st+" OR client_happened like "+st+"  and o. sub_project_manager_id in (select id from manager_project where manager_id in (select id from manager where manager_email_id = ? )) LIMIT ?, 10 ", UserName, offset)
			catch(err)
			defer result.Close()
			count, err1 := db.Query("select count(id) from open_positions where is_active = 1 and sub_project_manager_id in ( select id from manager_project where manager_id in ( select id from manager where manager_email_id = ?))", UserName)
			catch(err1)

			defer count.Close()
			for count.Next() {
				err2 := count.Scan(&s.Total_data)
				catch(err2)
			}
			s.Limit = 10
			s.Page = math.Ceil(s.Total_data / s.Limit)
			for result.Next() {
				var post models.Weekly
				var c string

				err := result.Scan(&post.Id, &post.Project_name, &post.Designation, &post.Type_position, &post.Position, &post.Priority, &post.Additonal_comment, &post.L1_Happened, &post.L2_Happened, &post.Client_Happened, &post.Is_active, &c)
				catch(err)
				location, err := time.LoadLocation("Asia/Kolkata")
				catch(err)
				layout := "2006-01-02 15:04:05"
				t, err := time.Parse(layout, c)
				to := time.Now()
				to = to.In(location)
				days := to.Sub(t) / (24 * time.Hour)
				fmt.Println("days", int(days))
				fmt.Println(c)

				post.Ageing = int(days)

				if post.Is_active == "1" {
					posts = append(posts, post)
				}
				s.Data = posts

			}
			json.NewEncoder(w).Encode(s)
		} else {
			fmt.Fprintf(w, "No status")
		}
	} else if strings.Contains(Role, "Program Manager") || strings.Contains(Role, "program manager") {
		params := mux.Vars(r)
		st, ok1 := r.URL.Query()["status"]

		if !ok1 || len(st[0]) < 1 {
			fmt.Fprintf(w, "status is missing")
			return
		}
		status := st[0]

		if status == "daily" || status == "Daily" {
			offsets, ok := r.URL.Query()["pages"]

			if !ok || len(offsets[0]) < 1 {
				fmt.Fprintf(w, "Url Param offset is missing")
				return
			}

			off := offsets[0]
			i, _ := strconv.Atoi(off)
			fmt.Printf("%T", i)

			offset := i * 10
			var s models.Pages_Daily
			var posts []models.Daily
			var st string
			st = "'" + params["id"] + "%'"
			result, err := db.Query("select o.id,p.project_name, o.designation,o.type_position,o.position, o.priority, o.additional_comment, o.l1_due,o.l2_due, o.client_due, o.created_at,o.is_active from open_positions o inner join manager_project m on o.sub_project_manager_id = m.id inner join projects p on p.id=m.project_id where o.is_active =1 and designation like "+st+" OR type_position like "+st+" OR position like "+st+" OR priority like "+st+" OR additional_comment like "+st+" OR l1_due like "+st+" OR l2_due like "+st+" OR client_due like "+st+" and o. sub_project_manager_id in (select id from manager_project where project_id in (select id from projects where program_manager = ? )) LIMIT ?, 10 ", UserName, offset)
			catch(err)
			defer result.Close()
			count, err1 := db.Query("select count(id) from open_positions where is_active = 1 and sub_project_manager_id in ( select id from manager_project where manager_id in ( select id from projects where program_manager = ?))", UserName)
			catch(err1)
			defer count.Close()
			for count.Next() {
				err2 := count.Scan(&s.Total_data)
				catch(err2)
			}
			s.Limit = 10
			s.Page = math.Ceil(s.Total_data / s.Limit)
			for result.Next() {
				var post models.Daily

				err := result.Scan(&post.Id, &post.Project_name, &post.Designation, &post.Type_position, &post.Position, &post.Priority, &post.Additonal_comment, &post.L1_due, &post.L2_due, &post.Client_due, &post.Created_at, &post.Is_active)
				catch(err)
				location, err := time.LoadLocation("Asia/Kolkata")
				catch(err)
				layout := "2006-01-02 15:04:05"
				t, err := time.Parse(layout, post.Created_at)
				to := time.Now()
				to = to.In(location)
				days := to.Sub(t) / (24 * time.Hour)
				//fmt.Println("days", int(days))

				post.Ageing = int(days)

				if post.Is_active == "1" {
					posts = append(posts, post)
				}
				s.Data = posts

			}
			json.NewEncoder(w).Encode(s)
		} else if status == "Weekly" || status == "weekly" {
			offsets, ok := r.URL.Query()["pages"]

			if !ok || len(offsets[0]) < 1 {
				fmt.Fprintf(w, "Url Param offset is missing")
				return
			}

			off := offsets[0]
			i, _ := strconv.Atoi(off)
			fmt.Printf("%T", i)

			offset := i * 10
			var s models.Pages_Weekly
			var posts []models.Weekly
			var st string
			st = "'" + params["id"] + "%'"
			result, err := db.Query("select o.id,p.project_name, o.designation,o.type_position,o.position, o.priority, o.additional_comment, o.l1_happened,o.l2_happened, o.client_happened, o.is_active, o.created_at from open_positions o inner join manager_project m on o.sub_project_manager_id = m.id inner join projects p on p.id=m.project_id where o.is_active =1 and designation like "+st+" OR type_position like "+st+" OR position like "+st+" OR priority like "+st+" OR additional_comment like "+st+" OR l1_happened like "+st+" OR l2_happened like "+st+" OR client_happened like "+st+" and o. sub_project_manager_id in (select id from manager_project where project_id in (select id from projects where program_manager = ? )) LIMIT ?, 10 ", UserName, offset)
			catch(err)
			defer result.Close()
			count, err1 := db.Query("select count(id) from open_positions where is_active = 1 and sub_project_manager_id in ( select id from manager_project where project_id in ( select id from projects where program_manager = ?))", UserName)
			catch(err1)
			defer count.Close()
			for count.Next() {
				err2 := count.Scan(&s.Total_data)
				catch(err2)
			}
			s.Limit = 10
			s.Page = math.Ceil(s.Total_data / s.Limit)
			for result.Next() {
				var post models.Weekly
				var c string

				err := result.Scan(&post.Id, &post.Project_name, &post.Designation, &post.Type_position, &post.Position, &post.Priority, &post.Additonal_comment, &post.L1_Happened, &post.L2_Happened, &post.Client_Happened, &post.Is_active, &c)
				catch(err)
				location, err := time.LoadLocation("Asia/Kolkata")
				catch(err)
				layout := "2006-01-02 15:04:05"
				t, err := time.Parse(layout, c)
				to := time.Now()
				to = to.In(location)
				days := to.Sub(t) / (24 * time.Hour)
				//fmt.Println("days", int(days))

				post.Ageing = int(days)

				if post.Is_active == "1" {
					posts = append(posts, post)
				}
				s.Data = posts

			}
			json.NewEncoder(w).Encode(s)
		}
	} else {
		fmt.Fprintf(w, "No authentication")
	}

}

// Deleteopenpositions :Func for deleting entries
func (c *Commander) Deleteopenpositions(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")

	if strings.Contains(Role, "Project Manager") || strings.Contains(Role, "project manager") {
		var post models.Daily
		json.NewDecoder(r.Body).Decode(&post)
		q2 := db.QueryRow("select project_manager_email from project_manager where id in ( select manager_id from sub_project_manager where id in ( select sub_project_manager_id from open_positions where id = ? ))", post.Id)
		var managerEmail string
		q2.Scan(&managerEmail)
		//	fmt.Println(manager_email)

		if managerEmail == UserName {
			var id1 int
			q1 := db.QueryRow("select id from open_positions where id = ? and is_active = 1", post.Id)
			q1.Scan(&id1)
			fmt.Println(id1)
			if id1 != 0 {
				stmt, err := db.Prepare("Update open_positions set is_active = 0 where id = ? ")
				catch(err)
				_, err = stmt.Exec(id1)
				catch(err)
				defer stmt.Close()
				fmt.Fprintf(w, "deleted successfully")
			} else {
				fmt.Fprintf(w, "Project doesnot exist")
			}
		} else {
			fmt.Fprintf(w, "project not under you")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorised access")
	}
}

// Updateopenpositions :Func for updating entries
func (c *Commander) Updateopenpositions(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()

	w.Header().Set("Content-Type", "application/json")

	if strings.Contains(Role, "Project Manager") || strings.Contains(Role, "project manager") {
		st, ok1 := r.URL.Query()["status"]

		if !ok1 || len(st[0]) < 1 {
			fmt.Fprintf(w, "status is missing")
			return
		}
		status := st[0]

		if status == "daily" || status == "Daily" {
			var post models.Daily
			json.NewDecoder(r.Body).Decode(&post)
			q2 := db.QueryRow("select id from sub_project_manager where sub_project_id in ( select id from sub_project where sub_project_name = ? ) and manager_id in ( select id from project_manager where project_manager_email = ? )", post.Project_name, UserName)
			var managerid int
			q2.Scan(&managerid)
			fmt.Println(managerid)
			if managerid != 0 {
				var id1 int
				q1 := db.QueryRow("select id from open_positions where sub_project_manager_id = ? and designation = ?  ", managerid, post.Designation)
				q1.Scan(&id1)
				fmt.Println(id1)
				if id1 != 0 {
					query, err := db.Prepare("Update open_positions set type_position = ?, position = ?, priority = ?, additional_comment = ?, l1_due = ?, l2_due = ?, client_due = ?, updated_at = ? where id = ?")
					catch(err)
					location, err := time.LoadLocation("Asia/Kolkata")
					catch(err)
					to := time.Now()
					to = to.In(location)
					_, er := query.Exec(post.Type_position, post.Position, post.Priority, post.Additonal_comment, post.L1_due, post.L2_due, post.Client_due, to, id1)
					catch(er)

					defer query.Close()

					fmt.Fprintf(w, "update successfully")
				} else {
					fmt.Fprintf(w, "Project doesnot exist")
				}
			} else {
				fmt.Fprintf(w, "Project not under you")
			}
		} else if status == "Weekly" || status == "weekly" {
			var post models.Weekly
			json.NewDecoder(r.Body).Decode(&post)
			q2 := db.QueryRow("select id from sub_project_manager where sub_project_id in ( select id from sub_project where sub_project_name = ? ) and manager_id in ( select id from project_manager where project_manager_email = ? )", post.Project_name, UserName)
			var managerid int
			q2.Scan(&managerid)
			fmt.Println(managerid)
			if managerid != 0 {
				var id1 int
				q1 := db.QueryRow("select id from open_positions where sub_project_manager_id = ? and designation = ?  ", managerid, post.Designation)
				q1.Scan(&id1)
				fmt.Println(id1)
				if id1 != 0 {
					query, err := db.Prepare("Update open_positions set type_position = ?, position = ?, priority = ?, additional_comment = ?, l1_happened = ?, l2_happened = ?, client_happened = ?, updated_at = ? where id = ?")
					catch(err)
					location, err := time.LoadLocation("Asia/Kolkata")
					catch(err)
					to := time.Now()
					to = to.In(location)
					_, er := query.Exec(post.Type_position, post.Position, post.Priority, post.Additonal_comment, post.L1_Happened, post.L2_Happened, post.Client_Happened, to, id1)
					catch(er)

					defer query.Close()

					fmt.Fprintf(w, "update successfully")
				} else {
					fmt.Fprintf(w, "Project doesnot exist")
				}
			} else {
				fmt.Fprintf(w, "Project not under you")
			}
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorised access")
	}

}
