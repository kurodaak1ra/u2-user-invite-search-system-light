package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

type dataTree struct {
	Name     string     `json:"name"`
	Value    string     `json:"value"`
	Children []dataTree `json:"children"`
}
type users struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Master   string `json:"master"`
}
type panelData struct {
	Json     string
	UID      string
	Username string
}
type fun func([]dataTree)

var mapData map[string]users

var store = sessions.NewCookieStore([]byte(os.Getenv("b0090e489f835142")))

func main() {
	loadFile()

	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("html/img"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("html/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("html/js"))))
	http.Handle("/", http.RedirectHandler("/auth", http.StatusFound))
	http.HandleFunc("/auth", auth)
	http.HandleFunc("/chart", chart)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func auth(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   0,
	}
	if auth, ok := session.Values["auth"].(bool); ok && auth {
		http.Redirect(w, r, "/chart", http.StatusFound)
		return
	}

	if r.Method == "GET" {
		t, _ := template.ParseFiles("html/auth.html")
		t.Execute(w, "")
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Printf("[ParseForm Err]: %v", err)
		t, _ := template.ParseFiles("html/auth.html")
		t.Execute(w, "内部解析错误")
		return
	}
	if _, ok := r.Form["uid"]; !ok {
		t, _ := template.ParseFiles("html/auth.html")
		t.Execute(w, "字段异常")
		return
	}
	if _, ok := r.Form["username"]; !ok {
		t, _ := template.ParseFiles("html/auth.html")
		t.Execute(w, "字段异常")
		return
	}
	if r.Form["uid"][0] == "" || r.Form["username"][0] == "" {
		t, _ := template.ParseFiles("html/auth.html")
		t.Execute(w, "请填写完整 UID 和 用户名")
		return
	}
	if mapData[r.Form["uid"][0]].Username == r.Form["username"][0] {
		session.Values["auth"] = true
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/chart", http.StatusFound)
	} else {
		t, _ := template.ParseFiles("html/auth.html")
		t.Execute(w, "UID 与 用户名 不匹配")
	}
}

func chart(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "token")
	if err != nil {
		http.Redirect(w, r, "/auth", http.StatusFound)
		return
	}
	if _, ok := session.Values["auth"]; !ok {
		http.Redirect(w, r, "/auth", http.StatusFound)
		return
	}

	if r.Method == "GET" {
		t, _ := template.ParseFiles("html/panel.html")
		jsonByte, _ := json.Marshal(&mapData)
		t.Execute(w, string(jsonByte))
		//t.Execute(w, panelData{})
		return
	}

	// search
	if err := r.ParseForm(); err != nil {
		log.Printf("[ParseForm Err]: %v", err)
		return
	}
	if _, ok := r.Form["uid"]; !ok {
		return
	}
	if _, ok := r.Form["submit"]; !ok {
		return
	}

	var (
		uid  string
		data dataTree
	)
	t, _ := template.ParseFiles("html/panel.html")
	if r.Form["submit"][0] == "self" {
		if mapData[r.Form["uid"][0]].ID != "" {
			uid = r.Form["uid"][0]
			data = generateJson(uid, mapData)
		}
	} else if r.Form["submit"][0] == "owner" {
		master := searchMaster(r.Form["uid"][0], "owner")
		if master != "" {
			uid = r.Form["uid"][0]
			data = generateJson(master, mapData)
		}
	} else if r.PostFormValue("submit") == "family" {
		master := searchMaster(r.Form["uid"][0], "family")
		if master != "" {
			uid = r.Form["uid"][0]
			data = generateJson(master, mapData)
		}
	}
	jsonByte, _ := json.Marshal(&data)
	t.Execute(w, panelData{
		Json:     string(jsonByte),
		UID:      uid,
		Username: mapData[r.Form["uid"][0]].Username,
	})
}

func loadFile() {
	jsonFile, err := os.Open("html/data.json")
	if err != nil {
		log.Fatalf("[JSON File Load Err]: %v", err)
	}
	defer jsonFile.Close()

	// parse slice
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var resultSlice []users
	json.Unmarshal([]byte(byteValue), &resultSlice)

	// parse map
	resultMap := make(map[string]users)
	for _, res := range resultSlice {
		resultMap[res.ID] = users{
			ID:       res.ID,
			Username: res.Username,
			Master:   res.Master,
		}
	}

	// return
	mapData = resultMap
}

func generateJson(master string, data map[string]users) dataTree {
	var masterTree []dataTree
	masterTree = append(masterTree, dataTree{
		Name:     data[master].Username,
		Value:    data[master].ID,
		Children: []dataTree{},
	})
	for _, xRes := range data {
		var f fun
		f = func(tree []dataTree) {
			for y, yRes := range tree {
				if xRes.Master == yRes.Value {
					tree[y].Children = append(yRes.Children, dataTree{
						Name:     xRes.Username,
						Value:    xRes.ID,
						Children: []dataTree{},
					})
					break
				}
				if len(yRes.Children) != 0 {
					f(yRes.Children)
				}
			}
		}
		f(masterTree)
	}
	return masterTree[0]
}

func searchMaster(uid string, typee string) string {
	if mapData[uid].ID == "" {
		return ""
	}
	if mapData[uid].Master == "0" {
		return uid
	}
	if typee == "owner" {
		if mapData[uid].Master != "0" {
			return mapData[uid].Master
		}
	} else if typee == "family" {
		if mapData[uid].Master != "0" {
			return searchMaster(mapData[uid].Master, "family")
		}
	}
	return ""
}
