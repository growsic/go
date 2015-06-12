package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //オプションを解析します。デフォルトでは解析しません。
	fmt.Println(r.Form) //このデータはサーバのプリント情報に出力されます。
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello 9191!") //ここでwに入るものがクライアントに出力されます。
}

func main() {
	http.HandleFunc("/", sayhelloName)       //アクセスのルーティングを設定します。
	err := http.ListenAndServe(":9191", nil) //監視するポートを設定します。
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
