package main

import (
	//"encoding/json"
	"github.com/hoisie/web"
)

type loginResponse struct {
	Success bool       `json:"success"`
	Files   []jsonFile `json:"files"`
}

type loginReq struct {
	User      string
    Password    string
	Response  chan *Response
}

/* func throwErr(arr *[]jsonFile) string {
	jsonResp := &jsonResponse{Success: false, Files: *arr}
	jresp, err := json.Marshal(jsonResp)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(jresp)
}*/

//Handles logins
func handleLogin(ctx *web.Context, login chan<- *loginReq) string {
    ctx.Request.ParseMultipartForm(4096)
    form := ctx.Request.MultipartForm
    user := form.Value["username"][0]
    password := form.Value["password"][0]

	loginResp := make(chan *Response)
	req := &loginReq{User: user, Password: password, Response: loginResp}
	login <- req
	resp := <-loginResp
	close(loginResp)
    return resp.message;
}
