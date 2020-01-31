package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	database "projectname_projectmanager/driver"
	models "projectname_projectmanager/model"
)

//ActionItemGetOverview : To get the overview of Action Items
func (C *Commander) ActionItemGetOverview(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	var overview models.Overview
	t1 := time.Now()
	t := t1.Format("2006-01-02")
	count, err1 := db.Query("SELECT COUNT(id) FROM action_items WHERE closed_in_time = 1")
	catch(err1)
	defer count.Close()
	for count.Next() {
		err2 := count.Scan(&overview.ClosedInTime)
		catch(err2)
	}
	count1, err3 := db.Query("SELECT COUNT(id) FROM action_items WHERE closed_in_time = 0")
	catch(err3)
	defer count1.Close()
	for count1.Next() {
		err4 := count1.Scan(&overview.ClosedOutTime)
		catch(err4)
	}
	count2, err5 := db.Query("SELECT COUNT(id) FROM action_items WHERE ? < target_date AND status LIKE 'inprogress'", t)
	catch(err5)
	defer count2.Close()
	for count2.Next() {
		err6 := count2.Scan(&overview.InProgressInTime)
		catch(err6)
	}
	count3, err7 := db.Query("SELECT COUNT(id) FROM action_items WHERE ? > target_date AND status LIKE 'inprogress'", t)
	catch(err7)
	defer count3.Close()
	for count3.Next() {
		err8 := count3.Scan(&overview.InProgressOutTime)
		catch(err8)
	}
	count4, err9 := db.Query("SELECT COUNT(id) FROM action_items WHERE status LIKE 'onhold'")
	catch(err9)
	defer count4.Close()
	for count4.Next() {
		err10 := count4.Scan(&overview.OnHold)
		catch(err10)
	}
	overview.GrandTotal = 0
	overview.GrandTotal = overview.GrandTotal + overview.ClosedInTime + overview.ClosedOutTime + overview.InProgressInTime + overview.InProgressOutTime + overview.OnHold
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(overview)
}

//ActionItemGetSayDo : to get the SayDo Ratio of Action Items
func (C *Commander) ActionItemGetSayDo(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	var total float64
	var closedintime float64
	var overviews []models.Manager
	var overview models.Manager
	count, err1 := db.Query("SELECT project_manager_name, count(project_manager_name) FROM action_items LEFT JOIN sub_project_manager ON action_items.sub_project_manager_id = sub_project_manager.id LEFT JOIN project_manager ON sub_project_manager.manager_id = project_manager.id LEFT JOIN sub_project ON sub_project_manager.sub_project_id = sub_project.id WHERE action_items.is_active = 1 AND WEEK(CURDATE())-WEEK(action_items.updated_at) < 1 GROUP BY project_manager_name")
	catch(err1)
	defer count.Close()
	for count.Next() {
		err2 := count.Scan(&overview.Name, &total)
		catch(err2)
		Owner := overview.Name
		c, er := db.Query("SELECT count(project_manager_name) FROM action_items LEFT JOIN sub_project_manager ON action_items.sub_project_manager_id = sub_project_manager.id LEFT JOIN project_manager ON sub_project_manager.manager_id = project_manager.id LEFT JOIN sub_project ON sub_project_manager.sub_project_id = sub_project.id WHERE closed_in_time = 1 AND action_items.is_active = 1 AND project_manager.project_manager_email = ? GROUP BY project_manager_name", Owner)
		catch(er)
		defer c.Close()
		for c.Next() {
			er1 := c.Scan(&closedintime)
			catch(er1)
		}
		overview.SayDo = 0.0
		overview.SayDo = ((closedintime / total) * 100)
		stmt, error1 := db.Query("SELECT manager FROM saydo WHERE manager = ?", Owner)
		catch(error1)
		if stmt.Next() != false {
			stmt1, error2 := db.Prepare("UPDATE saydo SET saydoratio = ? WHERE manager = ?")
			catch(error2)
			stmt1.Exec(overview.SayDo, Owner)
		} else {
			stmt2, error3 := db.Prepare("INSERT INTO saydo(manager,saydoratio) VALUES(?,?)")
			catch(error3)
			Say := overview.SayDo
			stmt2.Exec(Owner, Say)
		}
		c1, er2 := db.Query("SELECT saydoratio FROM saydo WHERE manager = ?", Owner)
		catch(er2)
		defer c1.Close()
		for c1.Next() {
			er3 := c1.Scan(&overview.SayDo)
			catch(er3)
		}
		overviews = append(overviews, overview)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(overviews)
}

//ActionItemGetTrend : to get the trend of Action Items
func (C *Commander) ActionItemGetTrend(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	w.Header().Set("Content-Type", "application/json")
	var wk, s int
	var ct, saydo float32
	var counts, counts1 [53]float32
	var week, week1 [53]int
	var SayDOTrend models.Saydo
	var SayDoTrends []models.Saydo

	offsets, ok := r.URL.Query()["year"]
	if !ok || len(offsets[0]) < 1 {
		fmt.Fprintf(w, "Url Parameter 'year' is missing")
		return
	}
	year := offsets[0]
	yr, err := strconv.Atoi(year)
	catch(err)

	count, err1 := db.Query("SELECT COUNT(WEEK(updated_at)),WEEK(updated_at) FROM action WHERE YEAR(updated_at) = ? GROUP BY WEEK(updated_at)", yr)
	defer count.Close()
	catch(err1)
	i := 0
	for count.Next() {
		count.Scan(&ct, &wk)
		counts[wk] = ct
		week[i] = wk
		i++
		s++
	}
	wk = 0
	ct = 0
	closed, err2 := db.Query("SELECT COUNT(WEEK(updated_at)),WEEK(updated_at) FROM action WHERE closed_in_time = 1 AND YEAR(updated_at) = ? GROUP BY WEEK(updated_at) ", yr)
	defer closed.Close()
	catch(err2)
	i = 0
	for closed.Next() {
		closed.Scan(&ct, &wk)
		counts1[wk] = ct
		week1[i] = wk
		i++

	}
	for i = 0; i < s; i++ {
		saydo = counts1[week[i]] / counts[week[i]] * 100
		Week := week[i]
		SayDOTrend.Saydo = saydo
		x, y := WeekRange(yr, Week)
		x1 := x.Format("2006-01-02")
		y1 := y.Format("2006-01-02")
		SayDOTrend.WeekStart = x1
		SayDOTrend.WeekEnd = y1
		SayDoTrends = append(SayDoTrends, SayDOTrend)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(SayDoTrends)
}

// WeekRange : To get the range of week using start date
func WeekRange(year, week int) (start, end time.Time) {
	start = WeekStart(year, week)
	end = start.AddDate(0, 0, 6)
	return
}

// WeekStart : To get the start date using week number
func WeekStart(year, week int) time.Time {
	// Start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)
	// Roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)
	return t
}
