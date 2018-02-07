package persistence

import (
	"io/ioutil"
	"errors"
	"encoding/json"
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"strconv"
	"Conus/model"
	"time"
)

var db *sql.DB
var requests []request

type dbsettings struct {
	Host, Port, User, Password, Dbname, Sslmode string
}

type request struct {
	Name, Request string
}

func (d dbsettings) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", d.Host, d.Port, d.User, d.Password, d.Dbname, d.Sslmode)
}

func Init() {
	settings := getSettings()
	var err error
	db, err = sql.Open("postgres", settings)
	if err != nil {
		panic("Could not open sql driver")
	}
	if err = db.Ping(); err != nil {
		panic("Could not ping database")
	}
	initRequests()
}

func initRequests() {
	jsn, err := ioutil.ReadFile("resources/dbReq.json")
	if err != nil {
		panic(errors.New("Could not read dbReq.json file"))
	}
	err = json.Unmarshal(jsn, &requests)
	if err != nil {
		panic(errors.New("Cannot unmarshal requests"))
	}
}

func getSettings() string {
	jsn, err := ioutil.ReadFile("resources/dbSettings.json")
	if err != nil {
		panic(errors.New("Cannot read db settings"))
	}
	var set = new(dbsettings)
	err = json.Unmarshal(jsn, &set)
	if err != nil {
		panic(errors.New("Cannot unmarshal json file of db settings"))
	}
	return set.String()
}
func getRequestByName(name string) (string, error) {
	for _, j := range requests {
		if j.Name == name {
			return j.Request, nil
		}
	}
	return "", errors.New("Wrong request name")
}
func makeUserQuery(query []string) *sql.Rows {
	req, err := getRequestByName(query[0])
	if err != nil {
		panic(err.Error())
	}
	if len(query) == 1 {
		rows, err := db.Query(req)
		if err != nil {
			panic(err.Error())
		}
		return rows
	}
	rows, err := db.Query(req, query[1])
	if err != nil {
		panic(err.Error())
	}
	return rows
}
func GetAllUsers() (*[]model.User, error) {
	rows := makeUserQuery([]string{"getAllUsers"})
	defer rows.Close()
	users := make([]model.User, 0)
	var id float64
	var firstname, lastname, username, email, password string
	for rows.Next() {
		err := rows.Scan(&id, &firstname, &lastname, &username, &email, &password)
		if err != nil {
			return nil, err
		}
		us := model.User{Id: id, FirstName: firstname, LastName: lastname, Username: username, Email: email, Password: ""}
		users = append(users, us)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &users, nil
}

func GetUserByName(name string) (model.User, error) {
	rows := makeUserQuery([]string{"getUserByName", name})
	defer rows.Close()
	var user model.User
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password)
		if err != nil {
			return model.User{}, err
		}
	}
	if err := rows.Err(); err != nil {
		return model.User{}, err
	}
	user.Password = ""
	return user, nil
}

func GetUserById(f float64) (model.User, error) {
	rows := makeUserQuery([]string{"getUserById", strconv.FormatFloat(f, 'f', 0, 64)})
	defer rows.Close()
	var user model.User
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.Password)
		if err != nil {
			return model.User{}, err
		}
	}
	if err := rows.Err(); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func AddUser(m []string) (int,error) {
	req, err := getRequestByName("addUser")
	if err != nil {
		return 0, err
	}
	//_, err = db.Exec(req, m[0], m[1], m[2], m[3], m[4])
	var id int
	err = db.QueryRow(req, m[0], m[1], m[2], m[3], m[4]).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdateUser(m []string) error {
	req, err := getRequestByName("updateUser")
	if err != nil {
		return err
	}
	_, err = db.Exec(req, m[0], m[1], m[2], m[3], m[4], m[5])
	if err != nil {
		return err
	}
	return nil
}

func DeleteUser(f float64) error {
	req, err := getRequestByName("deleteUser")
	if err != nil {
		return err
	}
	_, err = db.Exec(req, f)
	if err != nil {
		return err
	}
	return nil
}

func GetAllSessions() ([]Session, error) {
	rows := makeUserQuery([]string{"getAllSessions"})
	var ssid string
	var timeAcceced time.Time
	var values []byte
	result := make([]Session, 0)
	for rows.Next(){
		err := rows.Scan(&ssid, &timeAcceced, & values)
		if err != nil{
			return []Session{}, err
		}
		valuesMap, err := Unmarshal(values)
		if err != nil{
			return []Session{}, err
		}
		result = append(result, Session{ssid, timeAcceced, valuesMap})
	}
	return result, nil
}

func SessionInit(s *Session) error {
	if _, err := SessionRead(s.SessionId()); err == nil {
		return nil
	}
	req, err := getRequestByName("initSession")
	if err != nil {
		return err
	}
	js, err := json.Marshal(s.values)
	if err != nil{
		return err
	}
	_, err = db.Exec(req, s.sid, s.timeAcceced, js)
	if err != nil {
		return err
	}
	return nil
}

func SessionRead(sid string) (Session, error) {
	rows := makeUserQuery([]string{"readSession", sid})
	var ssid string
	var timeAcceced time.Time
	var values []byte
	for rows.Next() {
		err := rows.Scan(&ssid, &timeAcceced, &values)
		if err != nil {
			return Session{}, err
		}
	}
	valuesMap, err := Unmarshal(values)
	if err != nil {
		return Session{}, err
	}
	return Session{ssid, timeAcceced, valuesMap}, nil
}

func SessionDestroy(sid string) error{
	req, err := getRequestByName("deleteSessions")
	if err != nil{
		return err
	}
	_, err = db.Exec(req, sid)
	if err != nil{
		return err
	}
	return nil
}

func SessionUpdate(s *Session) error{
	req, err := getRequestByName("updateSession")
	if err != nil{
		return err
	}
	jsn, err := json.Marshal(s.values)
	if err != nil{
		return err
	}
	_, err = db.Exec(req, s.sid, s.timeAcceced, jsn)
	if err != nil{
		return err
	}
	return nil
}

func Unmarshal(values []byte) (map[string]interface{}, error) {
	valuesMap := make(map[string]interface{})
	err := json.Unmarshal(values, valuesMap)
	return valuesMap, err
}
