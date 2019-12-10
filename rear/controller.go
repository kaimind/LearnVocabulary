package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Controller struct {
	model Model
}

func (c Controller) getSession(w http.ResponseWriter, r *http.Request) {
	userid := getUserId(c.model, r)
	if userid <= 0 {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}

	rst := map[string]interface{}{"code": 0}
	json.NewEncoder(w).Encode(rst)
}

func (c Controller) postSession(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	userid := c.model.checkUser(username, password)
	if userid <= 0 {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}

	session, _ := c.model.session.Get(r, "session")
	session.Values["userid"] = userid
	session.Save(r, w)

	rst := map[string]interface{}{"code": 0}
	json.NewEncoder(w).Encode(rst)
}

func (c Controller) deleteSession(w http.ResponseWriter, r *http.Request) {
	session, _ := c.model.session.Get(r, "session")
	session.Values["userid"] = 0
	session.Save(r, w)

	rst := map[string]interface{}{"code": 0}
	json.NewEncoder(w).Encode(rst)
}

func (c Controller) getProgress(w http.ResponseWriter, r *http.Request) {
	userid := getUserId(c.model, r)
	if userid <= 0 {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}

	rst := c.model.getUserProgress(userid)
	json.NewEncoder(w).Encode(rst)
}

func (c Controller) getLearn(w http.ResponseWriter, r *http.Request) {
	userid := getUserId(c.model, r)
	if userid <= 0 {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}

	words := c.model.getLearnWords(userid, 8)
	json.NewEncoder(w).Encode(words)
}

func (c Controller) putLearn(w http.ResponseWriter, r *http.Request) {
	userid := getUserId(c.model, r)
	if userid <= 0 {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}

	var learn struct {
		Words []int `json:"words"`
	}
	json.NewDecoder(r.Body).Decode(&learn)
	c.model.putLearnWords(userid, learn.Words)

	rst := map[string]interface{}{"code": 0}
	json.NewEncoder(w).Encode(rst)
}

func (c Controller) getReview(w http.ResponseWriter, r *http.Request) {
	userid := getUserId(c.model, r)
	if userid <= 0 {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}

	learn := r.FormValue("learn")
	learnArray := strings.Split(learn, "|")

	total, complish, words := c.model.getReviewWords(userid, learnArray, 8)
	rst := map[string]interface{}{"total": total, "complish": complish, "words": words}
	json.NewEncoder(w).Encode(rst)
}

func (c Controller) putReview(w http.ResponseWriter, r *http.Request) {
	userid := getUserId(c.model, r)
	if userid <= 0 {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}

	var learn struct {
		Words []int `json:"words"`
	}
	json.NewDecoder(r.Body).Decode(&learn)
	c.model.putReviewWords(userid, learn.Words)

	rst := map[string]interface{}{"code": 0}
	json.NewEncoder(w).Encode(rst)
}

func (c Controller) getLearnList(w http.ResponseWriter, r *http.Request) {
	userid := getUserId(c.model, r)
	if userid <= 0 {
		http.Error(w, "BadRequest", http.StatusBadRequest)
		return
	}

	rst := c.model.getLearnList(userid)
	json.NewEncoder(w).Encode(rst)
}

func newController(model Model) Controller {
	return Controller{
		model: model,
	}
}

func getUserId(m Model, r *http.Request) int {
	session, err := m.session.Get(r, "session")
	if err != nil || session.IsNew {
		return 0
	}

	sessionid, ok := session.Values["userid"].(int)
	if ok {
		return sessionid
	}

	return 0
}
