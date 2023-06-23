package main

// net/httpパッケージを使ってウェブサーバ構築
import (
	"fmt"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	// wパラメータ: HTTPレスポンスを書き込む
	// rパラメータ: 受け取ったHTTPリクエストの情報を持つ
	fmt.Fprintf(w, "Welcome to My HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func aboutPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to My AboutPage!")
	fmt.Println("Endpoint Hit: aboutPage")
}

func main() {
	fmt.Println("AccessURL: http://localhost:8080/")
	// http.HandleFunc("/", homePage)       // URLルート("/")をhomePage関数にマッピング
	// http.HandleFunc("/about", aboutPage) // URLルート("/about")

	fs := http.FileServer(http.Dir(""))
	http.Handle("/", fs)

	http.ListenAndServe(":8080", nil) // Webサーバの起動
}
