package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json" //token.jsonが必要
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json") //credentials.jsonが必要
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, gmail.MailGoogleComScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	fmt.Println("Created Gmail service", srv)

	msgStr := "From: 'me'\r\n" +
		"reply-to: hogehoge@xxx.com\r\n" + //送信元アドレス
		"To: hoge@yyy.com\r\n" + //送信先アドレス
		"Subject:仮登録完了\r\n" +
		"\r\n" + "本登録は以下のURLからお願いします"
	reader := strings.NewReader(msgStr)
	transformer := japanese.ISO2022JP.NewEncoder()
	msgISO2022JP, err := ioutil.ReadAll(transform.NewReader(reader, transformer))
	if err != nil {
		log.Fatalf("Unable to convert to ISO2022JP: %v", err)
	}
	msg := []byte(msgISO2022JP)
	message := gmail.Message{}
	message.Raw = base64.StdEncoding.EncodeToString(msg)
	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		fmt.Printf("%v", err)
	}
}

// package main

// // net/httpパッケージを使ってウェブサーバ構築
// import (
// 	"fmt"
// 	"net/http"
// )

// func homePage(w http.ResponseWriter, r *http.Request) {
// 	// wパラメータ: HTTPレスポンスを書き込む
// 	// rパラメータ: 受け取ったHTTPリクエストの情報を持つ
// 	fmt.Fprintf(w, "Welcome to My HomePage!")
// 	fmt.Println("Endpoint Hit: homePage")
// }

// func aboutPage(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome to My AboutPage!")
// 	fmt.Println("Endpoint Hit: aboutPage")
// }

// func main() {
// 	fmt.Println("AccessURL: http://localhost:8080/")
// 	http.HandleFunc("/", homePage)       // URLルート("/")をhomePage関数にマッピング
// 	http.HandleFunc("/about", aboutPage) // URLルート("/about")

// 	fs := http.FileServer(http.Dir("templetes/")) // 指定したディレクトリ内のファイルをサーバー上で配信
// 	http.Handle("/index", fs)                     // サーバーのURLルート("/")で静的ファイルを提供

// 	http.ListenAndServe(":8080", nil) // Webサーバの起動
// }
