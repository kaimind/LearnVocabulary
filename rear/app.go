package main

// build cmd:
// GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o vocabulary *.go

// npm run build
// 如果编译后不希望看到react源码：
//     项目根目录新建.env.production文件，添加内容如下：
//     GENERATE_SOURCEMAP=false

import (
	"flag"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

var (
	AppPort    string
	RootPath   string
	MysqlUrl   string
	SessionKey string
)

func CacheControlWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=2592000")
		h.ServeHTTP(w, r)
	})
}

func Index(w http.ResponseWriter, r *http.Request) {
	file := path.Join(RootPath, "index.html")
	http.ServeFile(w, r, file)
}

func main() {
	flag.StringVar(&AppPort, "port", "5000", "http listen port")
	flag.StringVar(&RootPath, "path", "/path/LearnVocabulary/front/build", "app root path")
	flag.StringVar(&MysqlUrl, "mysql", "user:password@(127.0.0.1:3306)/learn_vocabulary", "mysql")
	flag.StringVar(&SessionKey, "session", "Paeh9eivEiJuo1Vu", "session key pairs")
	flag.Parse()

	model := newModel()
	controller := newController(model)

	router := mux.NewRouter()

	router.HandleFunc("/index", Index)

	dir := path.Join(RootPath, "static/")
	staticFileHandler := CacheControlWrapper(http.FileServer(http.Dir(dir)))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileHandler))

	router.HandleFunc("/api/user", controller.getSession).Methods("GET")
	router.HandleFunc("/api/user", controller.postSession).Methods("POST")
	router.HandleFunc("/api/user", controller.deleteSession).Methods("DELETE")

	router.HandleFunc("/api/learn", controller.getLearn).Methods("GET")
	router.HandleFunc("/api/learn", controller.putLearn).Methods("PUT")

	router.HandleFunc("/api/review", controller.getReview).Methods("GET")
	router.HandleFunc("/api/review", controller.putReview).Methods("PUT")

	router.HandleFunc("/api/progress", controller.getProgress).Methods("GET")
	router.HandleFunc("/api/learnlist", controller.getLearnList).Methods("GET")

	http.ListenAndServe(":"+AppPort, router)
}
