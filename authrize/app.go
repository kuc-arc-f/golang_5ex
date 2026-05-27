package main

import (
    "bytes"
	"context"
	"encoding/json"
	"fmt"
    "io"
    "net/http"
	"os"

	"authrize/handler"

	"github.com/wailsapp/wails/v2/pkg/runtime"
    "golang.org/x/crypto/bcrypt"
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
const API_URL_BASE   = "http://localhost:8787"

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
    if acttion_str == "user_create" {
        fmt.Printf("data=%s\n", req.Data)

        type DataReq struct {
                UserName string `json:"username"`
                Email string `json:"email"`
                Password string `json:"password"`
        }
        var data_req DataReq

        json2 := []byte(req.Data)
        err = json.Unmarshal(json2, &data_req)
        if err != nil {
                fmt.Println("JSON Unmarshalエラー:", err)
                return ""
        }
        fmt.Printf("UserName=%s\n", data_req.UserName)
        fmt.Printf("Password=%s\n", data_req.Password)
        // hash化
        hash, err := bcrypt.GenerateFromPassword(
            []byte(data_req.Password),
            bcrypt.DefaultCost,
        )
        if err != nil {
            panic(err)
        }
        fmt.Println("hash=" + string(hash))  
        var send_req DataReq
        send_req.UserName = data_req.UserName
        send_req.Email = data_req.Email
        send_req.Password = string(hash)
        jsonData, err := json.Marshal(&send_req)
        if err != nil {
            fmt.Println("JSONマーシャルエラー:", err)
            return ""
        }
        var endpoint = API_URL_BASE + "/api/users/create"
        fmt.Println("endpoint=" + string(endpoint)) 
        resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
            fmt.Println("リクエスト送信エラー:", err)
            return ""
        }
        defer resp.Body.Close()   
        // レスポンスボディの読み取り
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            fmt.Printf("レスポンス読み取りエラー: %v\n", err)
            return ""
        }        
        //fmt.Printf("body: %v\n", len(body))
        fmt.Printf("body: %v\n", string(body))
        if resp.StatusCode >= 200 && resp.StatusCode < 300 {
            fmt.Println("通信成功！")
            res := ActionRes{Data: "", Ret: 200}
            json3, err := json.Marshal(res)
            if err != nil {
                fmt.Println(err)
                return ""
            }	
            // JavaScriptへ返信イベント送信
            runtime.EventsEmit(a.ctx, "go-message", string(json3))
            return string(json3)             
        } else {
            fmt.Printf("通信失敗: ステータスコード %d\n", resp.StatusCode)
            res := ActionRes{Data: "", Ret: 500}
            json3, err := json.Marshal(res)
            if err != nil {
                fmt.Println(err)
                return ""
            }	
            // JavaScriptへ返信イベント送信
            runtime.EventsEmit(a.ctx, "go-message", string(json3))
            return string(json3)             
        }        
	}
    if acttion_str == "user_get" {
        fmt.Printf("data=%s\n", req.Data)
        type DataReq struct {
                Email string `json:"email"`
                Password string `json:"password"`
        }
        var data_req DataReq
        json2 := []byte(req.Data)
        err = json.Unmarshal(json2, &data_req)
        if err != nil {
                fmt.Println("JSON Unmarshalエラー:", err)
                return ""
        }
        fmt.Printf("Email=%s\n", data_req.Email)
        fmt.Printf("Password=%s\n", data_req.Password)
        var send_req DataReq
        send_req.Email = data_req.Email
        send_req.Password = data_req.Password
        jsonData, err := json.Marshal(&send_req)
        if err != nil {
            fmt.Println("JSONマーシャルエラー:", err)
            return ""
        }
        var endpoint = API_URL_BASE + "/api/users/get"
        fmt.Println("endpoint=" + string(endpoint)) 
        resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
            fmt.Println("リクエスト送信エラー:", err)
            return ""
        }
        defer resp.Body.Close()   
        // レスポンスボディの読み取り
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            fmt.Printf("レスポンス読み取りエラー: %v\n", err)
            return ""
        }        
        fmt.Printf("body: %v\n", len(body))
        if resp.StatusCode >= 200 && resp.StatusCode < 300 {
            fmt.Println("通信成功！")

            type LoginRes struct {
                Id int `json:"id"`
                Email string `json:"email"`
                Password string `json:"password"`
                Name string `json:"name"`
                CreatedAt string `json:"createdAt"`
                UpdatedAt string `json:"updatedAt"`
            }
            type LoginData struct {
                Ret string `json:"ret"`
                Data LoginRes `json:"data"`
            }
            var login_data LoginData
            err = json.Unmarshal(body, &login_data)
            if err != nil {
                    fmt.Println("JSON Unmarshalエラー:", err)
                    return ""
            }
            fmt.Printf("Email=%s\n", login_data.Data.Email)
            fmt.Printf("Password=%s\n", login_data.Data.Password)            
            err = bcrypt.CompareHashAndPassword(
                []byte(login_data.Data.Password),
                []byte(data_req.Password),
            )
            if err != nil {
                fmt.Println("NG")
                res := ActionRes{Data: "", Ret: 400}
                json3, err := json.Marshal(res)
                if err != nil {
                    fmt.Println(err)
                    return ""
                }	
                runtime.EventsEmit(a.ctx, "go-message", string(json3))
                return string(json3)             
            }
            // OK-CompareHashAndPassword
            res := ActionRes{Data: string(body), Ret: 200}
            json3, err := json.Marshal(res)
            if err != nil {
                fmt.Println(err)
                return ""
            }	
            runtime.EventsEmit(a.ctx, "go-message", string(json3))
            return string(json3)             
        }        
        //fmt.Printf("body: %v\n", string(body))
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