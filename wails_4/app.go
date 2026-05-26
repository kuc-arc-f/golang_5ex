package main

import (
	"os"
	"context"
	"encoding/json"
	"fmt"

	"wails_4/handler"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

type ActionReq struct {
    Action string `json:"action"`
    Data  string    `json:"data"`
}
type ActionRes struct {
	Ret  int
	Data string
}

// JavaScript から呼ばれる
func (a *App) SendMessage(msg string) string {
    println("JSから受信:", msg)
    var req ActionReq
    json1 := []byte(msg)

    // 第1引数: JSONデータ ([]byte)
    // 第2引数: 格納先の構造体のポインタ (&p)
    err := json.Unmarshal(json1, &req)
    if err != nil {
        fmt.Println("JSON Unmarshalエラー:", err)
        return ""
    }
    var acttion_str = fmt.Sprintf("%s", req.Action)
    fmt.Printf("ret_str=%s\n", acttion_str)
    if acttion_str == "todo_create" {
        fmt.Printf("data=%s\n", req.Data)

        type DataReq struct {
                Title string `json:"title"`
        }
        var data_req DataReq

        json2 := []byte(req.Data)
        err = json.Unmarshal(json2, &data_req)
        if err != nil {
                fmt.Println("JSON Unmarshalエラー:", err)
                return ""
        }
        fmt.Printf("Title=%s\n", data_req.Title)
        err = handler.CmdAdd(data_req.Title)
        if err != nil {
                fmt.Println("handler.CmdAddエラー:", err)
                return ""
        }
    }
    if acttion_str == "todo_list" {
        fmt.Printf("#start-todo_list \n")
        var ar_str = handler.CmdList()
        fmt.Printf("ar_str=%s\n", ar_str)
        res := ActionRes{Data: ar_str, Ret: 200}
        json3, err := json.Marshal(res)
        if err != nil {
            fmt.Println(err)
            return ""
        }	
        // JavaScriptへ返信イベント送信
        runtime.EventsEmit(a.ctx, "go-message", string(json3))
        return string(json3)        
    }
    if acttion_str == "todo_delete" {
        fmt.Printf("data=%s\n", req.Data)

        type DeleteReq struct {
                Id int `json:"id"`
        }
        var data_req DeleteReq

        json2 := []byte(req.Data)
        err = json.Unmarshal(json2, &data_req)
        if err != nil {
                fmt.Println("JSON Unmarshalエラー:", err)
                return ""
        }
        fmt.Printf("id=%d\n", data_req.Id)
        var id = data_req.Id;
        var err = handler.CmdDelete(id)
        res := ActionRes{Data: string(id), Ret: 200}
        json3, err := json.Marshal(res)
        if err != nil {
            fmt.Println(err)
            return ""
        }	
        runtime.EventsEmit(a.ctx, "go-message", string(json3))
        return string(json3)        
    }

    res := ActionRes{Data: "", Ret: 200}
    json3, err := json.Marshal(res)
    if err != nil {
        fmt.Println(err)
        return ""
    }	
    // JavaScriptへ返信イベント送信
    runtime.EventsEmit(a.ctx, "go-message", string(json3))
    return string(json3)
}

// JavaScript から呼ばれる
func (a *App) TestMessage(msg string) string {

  var err error

	err = handler.CmdAdd("test-data-1")
	println("JSから受信:", msg)

	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
	}
	// JavaScriptへ返信イベント送信
	runtime.EventsEmit(a.ctx, "go-message", "GOから返信: "+msg)

	return "JSから受信:" + msg

}
