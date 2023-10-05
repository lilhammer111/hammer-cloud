package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lilhammer111/hammer-cloud/db"
	"github.com/lilhammer111/hammer-cloud/util"
	"log"
	"net/http"
	"time"
)

const (
	pwdSalt = "*#890"
)

// SignUpHandler parse username and pwd field from frontend, and encrypt the pwd.
// And then invoke internally UserSignUp to insert user info into DB
func DoSignUpHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	pwd := c.Request.FormValue("password")

	if len(username) < 3 || len(pwd) < 5 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "invalid parameter",
			"code": -1,
		})
		return
	}

	encryptPWD := util.Sha1([]byte(pwd + pwdSalt))

	if ok := db.UserSignUp(username, encryptPWD); ok {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "login succeeded",
			"code": 0,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "login failed",
			"code": -2,
		})
	}
}

func SignUpHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

func DoLoginHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	pwd := c.Request.FormValue("password")

	// validates username or pwd
	encryptPWD := util.Sha1([]byte(pwd + pwdSalt))

	pwdChecked := db.UserLogin(username, encryptPWD)
	if !pwdChecked {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "username or password verification failed",
			"code": -1,
		})
		return
	}
	// generate token
	token := GenToken(username)
	// save the new token, and get res
	updateOK := db.UpdateToken(username, token)
	if !updateOK {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg":  "failed to update token",
			"code": -2,
		})
		return
	}
	// redirect to home, however, by the return path to home instead
	//_, err = w.Write([]byte("http://" + rd.Host + "/static/view/home.html"))
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	resp := util.RespBody{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			RedirectPath string
			Username     string
			Token        string
		}{
			RedirectPath: "/static/view/home.html",
			Username:     username,
			Token:        token,
		},
	}
	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
}

func LoginHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

func GenToken(username string) string {
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	username := r.Form.Get("username")

	user, err := db.GetUserInfo(username)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp := util.RespBody{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	_, err = w.Write(resp.JSONBytes())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
