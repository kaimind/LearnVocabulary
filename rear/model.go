package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

type WordInfo struct {
	Wid      int        `json:"wid"`
	Word     string     `json:"word"`
	Category int        `json:"category"`
	Detail   WordDetail `json:"detail"`
	Counter  int        `json:"counter"`
}

type WordDetail struct {
	Read   []WordRead `json:"read"`
	Define WordDefine `json:"define"`
	Use    WordUse    `json:"use"`
}

type WordRead struct {
	Spell string `json:"spell"`
	Audio string `json:"audio"`
}

type WordDefine struct {
	Zh   string `json:"zh"`
	Form string `json:"form"`
	Dual string `json:"dual"`
	En   string `json:"en"`
}

type WordUse struct {
	Collocation string           `json:"collocation"`
	Phrase      string           `json:"phrase"`
	Synonym     string           `json:"synonym"`
	Example     []WordUseExample `json:"example"`
}

type WordUseExample struct {
	En string `json:"en"`
	Zh string `json:"zh"`
}

type Model struct {
	db      *sql.DB
	session *sessions.CookieStore
}

func (w *WordDetail) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b, _ := src.([]byte)
	return json.Unmarshal(b, w)
}

func (m Model) getUserProgress(uid int) map[string]interface{} {
	var tody, todayTarget, total, totalTarget int

	query := "SELECT count(*) FROM learn WHERE uid = ? and TO_DAYS(learn)=TO_DAYS(now())"
	m.db.QueryRow(query, uid).Scan(&tody)

	todayTarget = 72

	query = "SELECT count(*) FROM learn WHERE uid = ? and learn is not null"
	m.db.QueryRow(query, uid).Scan(&total)

	query = "SELECT count(*) FROM learn WHERE uid = ?"
	m.db.QueryRow(query, uid).Scan(&totalTarget)

	return map[string]interface{}{
		"tody":        tody,
		"todayTarget": todayTarget,
		"total":       total,
		"totalTarget": totalTarget,
	}
}

func (m Model) getLearnWords(uid, limit int) []WordInfo {
	query := `select learn.wid, words.word, words.category, words.detail from learn 
		left join words on learn.wid = words.id 
		where learn.uid = ? and learn.counter = 0 order by rand() limit ?`
	rows, _ := m.db.Query(query, uid, limit)
	defer rows.Close()

	words := []WordInfo{}
	for rows.Next() {
		info := WordInfo{}
		rows.Scan(&info.Wid, &info.Word, &info.Category, &info.Detail)
		words = append(words, info)
	}
	return words
}

func (m Model) getReviewWords(uid int, learn []string, limit int) (int, int, []WordInfo) {
	total := 0
	complish := 0
	words := []WordInfo{}

	learnStr := ""
	sep := ""
	for i := 0; i < len(learn); i++ {
		learnStr += sep + fmt.Sprintf("TO_DAYS('%s')", learn[i])
		sep = ", "
	}
	query := `select count(*) from learn where uid=? and TO_DAYS(learn.learn) in (%s)`
	query = fmt.Sprintf(query, learnStr)
	fmt.Println(query)
	m.db.QueryRow(query, uid).Scan(&total)
	if total <= 0 {
		return total, complish, words
	}

	query = `select count(*) from learn where uid=? and TO_DAYS(learn.review)=TO_DAYS(now())
		and TO_DAYS(learn.learn) in (%s)`
	query = fmt.Sprintf(query, learnStr)
	fmt.Println(query)
	m.db.QueryRow(query, uid).Scan(&complish)
	if complish >= total {
		return total, complish, words
	}

	query = `select learn.wid, words.word, words.category, words.detail, learn.counter from learn 
		left join words on learn.wid = words.id 
		where learn.uid = ? and TO_DAYS(learn.learn) in (%s)
		and TO_DAYS(learn.review)<TO_DAYS(now()) order by rand() limit ?`
	query = fmt.Sprintf(query, learnStr)
	fmt.Println(query)
	rows, _ := m.db.Query(query, uid, limit)
	defer rows.Close()

	for rows.Next() {
		info := WordInfo{}
		rows.Scan(&info.Wid, &info.Word, &info.Category, &info.Detail, &info.Counter)
		if info.Wid > 0 {
			words = append(words, info)
		}
	}

	return total, complish, words
}

func (m Model) putLearnWords(uid int, words []int) {
	para := []interface{}{uid}
	query := `update learn set counter=counter+1, learn=now(), review=now() where uid=? and learn is null and wid in (`
	sep := ""
	for i := 0; i < len(words); i++ {
		query += sep + "?"
		sep = ", "
		para = append(para, words[i])
	}
	query += ")"
	_, err := m.db.Exec(query, para...)
	if err != nil {
		log.Println(err)
	}

	query = `insert ignore into review (uid, learn, review) values (?, now(), now())`
	_, err1 := m.db.Exec(query, uid)
	if err1 != nil {
		log.Println(err1)
	}
}

func (m Model) getLearnList(uid int) []map[string]interface{} {
	query := `select learn, TO_DAYS(now())-TO_DAYS(review) from review where uid=? order by learn desc limit 120`
	rows, _ := m.db.Query(query, uid)
	defer rows.Close()

	learnList := []map[string]interface{}{}
	for rows.Next() {
		learn := ""
		review := 0

		rows.Scan(&learn, &review)

		if len(learn) > 10 {
			learn = learn[0:10]
		}
		item := map[string]interface{}{
			"learn":  learn,
			"review": review,
		}
		learnList = append(learnList, item)
	}
	return learnList
}

func (m Model) putReviewWords(uid int, words []int) {
	para := []interface{}{uid}
	wordids := ""
	sep := ""
	for i := 0; i < len(words); i++ {
		wordids += sep + "?"
		sep = ", "
		para = append(para, words[i])
	}

	query := `update learn set counter=counter+1, review=now() 
		where uid=? and TO_DAYS(review)<>TO_DAYS(now()) and wid in (` + wordids + `)`
	_, err := m.db.Exec(query, para...)
	if err != nil {
		log.Println(err)
	}

	query = `select distinct date(learn) from learn where uid=? and wid in (` + wordids + `)`
	rows, _ := m.db.Query(query, para...)
	defer rows.Close()

	learn := []string{}
	for rows.Next() {
		r := ""
		rows.Scan(&r)
		if len(r) > 10 {
			learn = append(learn, r[0:10])
		}
	}

	for _, r := range learn {
		count := 0
		query = `select count(*) from learn where uid=? and TO_DAYS(learn)=TO_DAYS(?) and TO_DAYS(review)<>TO_DAYS(now())`
		m.db.QueryRow(query, uid, r).Scan(&count)
		if count <= 0 {
			query = `update review set review=now(), counter=counter+1 where uid=? and learn=?`
			m.db.Exec(query, uid, r)
		}
	}

}

func (m Model) checkUser(username, password string) int {
	userid := 0
	query := "SELECT id FROM user WHERE username = ? and password = ?"
	m.db.QueryRow(query, username, password).Scan(&userid)
	return userid
}

func newModel() Model {
	db := connectDb()

	createTables(db)

	return Model{
		db:      db,
		session: sessions.NewCookieStore([]byte(SessionKey), []byte(SessionKey)),
	}
}

func connectDb() *sql.DB {
	db, err := sql.Open("mysql", MysqlUrl+"?parseTime=true")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
		return nil
	}

	return db
}

func createTables(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS words(
		id int AUTO_INCREMENT,
		word char(32) NOT NULL,
		category int NOT NULL DEFAULT 0,
		detail json DEFAULT NULL,
		PRIMARY KEY ( id ),
		UNIQUE KEY word (word)
	)ENGINE=InnoDB DEFAULT CHARSET=utf8`
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
		return
	}

	query = `
	CREATE TABLE IF NOT EXISTS learn(
		id INT AUTO_INCREMENT,
		uid INT NOT NULL,
		wid INT NOT NULL,
		learn DATETIME,
		review DATETIME,
		counter SMALLINT DEFAULT 0, 
		PRIMARY KEY ( id )
	)ENGINE=InnoDB DEFAULT CHARSET=utf8`

	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
		return
	}

	query = `
	CREATE TABLE IF NOT EXISTS review(
		id INT AUTO_INCREMENT,
		uid INT NOT NULL,
		learn DATE,
		review DATE,
		counter SMALLINT DEFAULT 0, 
		PRIMARY KEY ( id ),
		UNIQUE KEY uid_learn (uid, learn)
	)ENGINE=InnoDB DEFAULT CHARSET=utf8`

	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
		return
	}

	query = `
	CREATE TABLE IF NOT EXISTS user(
		id INT AUTO_INCREMENT,
		username CHAR(64) NOT NULL UNIQUE,
		password CHAR(64) NOT NULL,
		PRIMARY KEY ( id )
	)ENGINE=InnoDB DEFAULT CHARSET=utf8`

	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
		return
	}
}

// 学员（uid=1）开始学习cet6（category&2）的词书
// insert into learn (uid, wid) select 1, id from words where (category&2)>0;
// (category&1)>0: cet4
// (category&2)>0: cet6
// (category&4)>0: collins 5 star
// (category&8)>0: collins 4 star
// (category&16)>0: collins 3 star
// (category&32)>0: collins 2 star
// (category&64)>0: collins 1 star
