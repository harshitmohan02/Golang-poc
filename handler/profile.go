package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	database "projectname_projectmanager/driver"
	helper "projectname_projectmanager/helper"
	model "projectname_projectmanager/model"
	"strconv"
	"time"
)

// ChangeProfileImage : uploding the Profile Image to Server.
func (C *Commander) ChangeProfileImage(w http.ResponseWriter, r *http.Request) {
	Time := time.Now()
	db := database.DbConn()
	defer db.Close()
	var user model.Profile
	user.Name = UserName
	user.Role = Role
	//reading the user whose image we want to change from the database

	imageName, err := helper.FileUpload(r)
	//here we call the function we made to get the image and save it
	if err != nil {
		http.Error(w, "Invalid Data", http.StatusBadRequest)
		return
	}
	user.ImagePath = imageName
	User, _ := db.Query("SELECT id FROM profile WHERE username = ?", UserName)
	defer User.Close()
	if User.Next() != false {
		UpdateProfile, _ := db.Query("UPDATE profile set username = ?, role = ?, image_path = ?, updated_at = ?", user.Name, user.Role, user.ImagePath, Time)
		defer UpdateProfile.Close()
	} else {
		InsertProfile, _ := db.Query("INSERT into profile(username, role, image_path, created_at, updated_at)VALUES(?, ?, ?, ?, ?)", user.Name, user.Role, user.ImagePath, Time, Time)
		defer InsertProfile.Close()
	}
	setupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// FileServe : Serving file to client.
func (C *Commander) FileServe(w http.ResponseWriter, r *http.Request) {
	db := database.DbConn()
	defer db.Close()
	var Filename string
	GetFilePath, _ := db.Query("SELECT image_path FROM profile WHERE username = ?", UserName)
	defer GetFilePath.Close()
	if GetFilePath.Next() != false {
		GetFilePath.Scan(&Filename)
	}
	if Filename == "" {
		//Get not set, send a 400 bad request
		http.Error(w, "Get 'file' not specified in url.", 400)
		return
	}
	fmt.Println("Client requests: " + Filename)

	//Check if file exists and open
	Openfile, err := os.Open(Filename)
	defer Openfile.Close() //Close after function return
	if err != nil {
		//File not found, send 404
		http.Error(w, "File not found.", 404)
		return
	}

	//File is found, create and send the correct headers

	FileHeader := make([]byte, 512)
	Openfile.Read(FileHeader)
	FileContentType := http.DetectContentType(FileHeader)
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string
	setupResponse(&w, r)
	w.Header().Set("Content-Disposition", "attachment; filename="+Filename)
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)
	//We read 512 bytes from the file already, so we reset the offset back to 0
	Openfile.Seek(0, 0)
	io.Copy(w, Openfile) //'Copy' the file to the client
	return
}
