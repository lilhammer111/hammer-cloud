package db

import (
	"github.com/lilhammer111/hammer-cloud/db/mysql"
	"log"
)

// UserSignUp insert username and pwd into table tbl_user by ignore mode,
// that means there will be no error even thought insert failed such as
// constriction conflict etc. It's also why, I guess, UserSignUp returns bool type
func UserSignUp(username, passwd string) bool {
	stmtIns, err := mysql.DBConn().Prepare("insert ignore into tbl_user (user_name, user_pwd) values (?, ?);")
	if err != nil {
		log.Println(err)
		return false
	}
	defer stmtIns.Close()

	insRes, err := stmtIns.Exec(username, passwd)
	if err != nil {
		log.Println(err)
		return false
	}

	if rowsAffected, err := insRes.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}

	return false
}

func UserLogin(username, encpwd string) bool {
	stmtSel, err := mysql.DBConn().Prepare("select user_pwd from tbl_user where user_name = ? limit 1;")
	if err != nil {
		log.Println(err)
		return false
	}
	defer stmtSel.Close()

	var selPWD string
	err = stmtSel.QueryRow(username).Scan(&selPWD)
	if err != nil {
		log.Println(err)
		return false
	}

	if selPWD == encpwd {
		return true
	}
	return false
}

// UpdateToken refresh user's token
func UpdateToken(username, token string) bool {
	stmtRep, err := mysql.DBConn().Prepare("replace into tbl_user_token (user_name, user_token) values (?,?);")
	if err != nil {
		log.Println(err)
		return false
	}
	defer stmtRep.Close()

	_, err = stmtRep.Exec(username, token)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

func GetUserInfo(username string) (User, error) {
	user := User{}
	stmtSel, err := mysql.DBConn().
		Prepare("select user_name, signup_at from tbl_user where user_name = ? limit 1")
	if err != nil {
		log.Println(err)
		return user, err
	}
	defer stmtSel.Close()

	err = stmtSel.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}

	return user, nil
}
