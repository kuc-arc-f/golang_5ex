package main

import (
    "bytes"
    "context"
    "fmt"
    "io"
    "encoding/json"
    "net/http"
    "os"
    "log"

    "go-rag-2/handler"

    "github.com/joho/godotenv"
    "github.com/wailsapp/wails/v2/pkg/runtime"
    qdrant "github.com/qdrant/go-client/qdrant"
    "google.golang.org/grpc"    
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
const collectionName = "doc-collection"


type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}


type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

func send_chat(query string) string{
    var input = "日本語で、回答して欲しい。\n 要約して欲しい。\n" + query
    fmt.Printf("input: \n%v\n\n", input)

    history := []Message{
        {
            Role:    "system",
            Content: "You are a helpful assistant. 日本語で答えてください。",
        },
    }
    history = append(history, Message{
        Role:    "user",
        Content: input,
    })    
    var serverURL   = "http://localhost:8090/v1/chat/completions"
    var model       = "local-model"
    var temperature = 0.7

    reqBody := ChatRequest{
        Model:       model,
        Messages:    history,
        Temperature: temperature,
    }

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("JSONマーシャルエラー:", err)
		return ""
	}
	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
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
	
	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
        fmt.Errorf("JSONデコードエラー: %w", err)
		return "" 
	}
	if len(chatResp.Choices) == 0 {
        fmt.Errorf("レスポンスにChoicesがありません")
		return "" 
	}

    var outStr string = chatResp.Choices[0].Message.Content;
    //fmt.Printf("\n outStr %s\n\n", outStr)
    return outStr;
}

// JavaScript から呼ばれる
func (a *App) SendMessage(msg string) string {
    err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}    
	// 環境変数からAPIキーを取得
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("エラー: 環境変数 GEMINI_API_KEY が設定されていません")
        return ""
	}    
    println("JSから受信:", msg)
    var req ActionReq
    json1 := []byte(msg)

    // 第1引数: JSONデータ ([]byte)
    // 第2引数: 格納先の構造体のポインタ (&p)
    err = json.Unmarshal(json1, &req)
    if err != nil {
        fmt.Println("JSON Unmarshalエラー:", err)
        return ""
    }
    var acttion_str = fmt.Sprintf("%s", req.Action)
    fmt.Printf("ret_str=%s\n", acttion_str)
    if acttion_str == "rag_search" {
        fmt.Printf("data=%s\n", req.Data)
        type DataReq struct {
            Query string `json:"query"`
        }
        var data_req DataReq

        json2 := []byte(req.Data)
        err = json.Unmarshal(json2, &data_req)
        if err != nil {
                fmt.Println("JSON Unmarshalエラー:", err)
                return ""
        }
        fmt.Printf("Query=%s\n", data_req.Query)
        // 関数呼び出し
        embeddings, err := handler.GetEmbeddings(data_req.Query, apiKey)
        if err != nil {
            fmt.Printf("エラーが発生しました: %v\n", err)
            return ""
        }           
        // 結果の出力
        fmt.Println("\n取得したベクトルデータ:")
        fmt.Printf("次元数: %d\n", len(embeddings))  
        // 4. ベクトル検索 (KNN Query)
        conn, err := grpc.Dial(
            "localhost:6334",
            grpc.WithInsecure(),
        )
        if err != nil {
            log.Fatal(err)
        }
        defer conn.Close()

        client := qdrant.NewPointsClient(conn)

        fmt.Println("Qdrant client connected")
        var queryVector = embeddings

        resp, err := client.Search(context.Background(), &qdrant.SearchPoints{
            CollectionName: collectionName,
            Vector:         queryVector,
            Limit:          1,
            WithPayload: &qdrant.WithPayloadSelector{
                SelectorOptions: &qdrant.WithPayloadSelector_Enable{
                    Enable: true,
                },
            },
        })
        if err != nil {
            log.Fatal(err)
        }
        var matches string = ""
        for _, p := range resp.Result {
            var contentStr = p.Payload["content"].GetStringValue()
            fmt.Printf("ID=%v , score=%.4f \n" , p.Id, p.Score)
            if (p.Score > 0.6) {
                matches += contentStr + "\n"
            }                        
            //fmt.Printf("ID=%v score=%.4f \npayload=%v\n",
            //	p.Id, p.Score, contentStr)
        }   
        //fmt.Printf("matches=%v\n",matches)
        var outText string = ""
        if (len(matches) > 0){
            outText = `context:` + matches + "\n"
            outText += `user query:` + data_req.Query + "\n"
        }else{
            outText =`user query:` + data_req.Query + "\n"
        }    
        var out_str = send_chat(outText)        
        res := ActionRes{Data: out_str, Ret: 200}
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
	println("JSから受信:", msg)
	// JavaScriptへ返信イベント送信
	runtime.EventsEmit(a.ctx, "go-message", "GOから返信: "+msg)

	return "JSから受信:" + msg

}