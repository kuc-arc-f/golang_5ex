package main

import (
    "context"
    "fmt"
    "encoding/json"
    //"os"
    "syscall"
    "unsafe"
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
    dll := syscall.NewLazyDLL("example.dll")
    todoAddFunc := dll.NewProc("todo_add")
    todoListFunc := dll.NewProc("todo_list")
    todoDeleteFunc := dll.NewProc("todo_delete")

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
        var input = data_req.Title
        fmt.Println("title=", input)
        title := input
        // Goの文字列を、C言語形式（NULL終端のバイト配列）のポインタに変換
        cStrPointer, err := syscall.BytePtrFromString(title)
        if err != nil {
            panic(err)
        }
        todoAddFunc.Call(uintptr(unsafe.Pointer(cStrPointer)))        
    }
    if acttion_str == "todo_list" {
        fmt.Printf("#start-todo_list \n")
        buffer := make([]byte, 5000)    
        // 2. バッファの先頭ポインタと、そのサイズを取得
        bufPtr := uintptr(unsafe.Pointer(&buffer[0]))
        bufSize := uintptr(len(buffer))

        // 3. C++の関数を呼び出し、バッファに書き込んでもらう
        // 戻り値(r1)には、C++側で書き込まれた文字数が返ってきます
        r1, _, _ := todoListFunc.Call(bufPtr, bufSize)
        stringLen := int(r1)

        if stringLen > 0 {
            // 4. 書き込まれたバイト数分だけ切り出して、Goの文字列に変換
            cppMessage := string(buffer[:stringLen])
            fmt.Printf("[Go]  C++から受信した文字列: %s\n", cppMessage)
            res := ActionRes{Data: cppMessage , Ret: 200}
            json3, err := json.Marshal(res)
            if err != nil {
                fmt.Println(err)
                return ""
            }	
            runtime.EventsEmit(a.ctx, "go-message", string(json3))
            return string(json3)        
        } else {
            fmt.Println("[Go]  文字列の受信に失敗しました（バッファサイズ不足など）")
        }        
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
        todoDeleteFunc.Call(uintptr(id))        
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
	// JavaScriptへ返信イベント送信
	runtime.EventsEmit(a.ctx, "go-message", "GOから返信: "+msg)

	return "JSから受信:" + msg

}