package handler

import (
	"fmt"
	"github.com/lilhammer111/hammer-cloud/db"
	"github.com/lilhammer111/hammer-cloud/util"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	pwdSalt = "*#890"
)

// SignUpHandler parse username and pwd field from frontend, and encrypt the pwd.
// And then invoke internally UserSignUp to insert user info into DB
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		html, err := os.ReadFile("./static/view/signup.html")
		if err != nil {
			http.Error(w, "internal server", http.StatusNotFound)
			return
		}

		_, err = w.Write(html)
		if err != nil {
			http.Error(w, "internal server", http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "internal server", http.StatusInternalServerError)
			return
		}

		username := r.Form.Get("username")
		passwd := r.Form.Get("passwd")

		if len(username) < 3 || len(passwd) < 5 {
			http.Error(w, "invalid parameter", http.StatusBadRequest)
			return
		}

		encPasswd := util.Sha1([]byte(passwd + pwdSalt))

		if ok := db.UserSignUp(username, encPasswd); ok {
			_, err := w.Write([]byte("signup success"))
			if err != nil {
				log.Println(err)
			}
		} else {
			_, err := w.Write([]byte("signup fail "))
			if err != nil {
				log.Println(err)
			}
		}
		return
	}

	http.Error(w, "wrong request method", http.StatusBadRequest)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}
	username := r.Form.Get("username")
	pwd := r.Form.Get("password")

	// validates username or pwd
	encpwd := util.Sha1([]byte(pwd + pwdSalt))
	checked := db.UserLogin(username, encpwd)
	if !checked {
		http.Error(w, "invalid username or password", http.StatusBadRequest)
		return
	}
	// generate token
	token := GenToken(username)
	// save the new token, and get res
	updSuc := db.UpdateToken(username, token)
	if !updSuc {
		http.Error(w, "generate token failed", http.StatusInternalServerError)
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
			RedirPath string
			Username  string
		}{
			RedirPath: "http://" + r.Host + "/static/view/home.html",
			Username:  username,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", token)
	_, err = w.Write(resp.JSONBytes())
	if err != nil {
		log.Println(err)
		return
	}
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
