package handler

import (
    //"database/sql"
    "bytes"
    //"context"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"

    "github.com/tmc/langchaingo/textsplitter"
)

const DATA_DIR = "./data"
const CHUNK_SIZE_MAX = 1000

type ReadParam struct {
    Content  string    `json:"content"`
    Name     string    `json:"name"`
}
// リクエストの構造体
type EmbeddingRequest struct {
    Model   string  `json:"model"`
    Content EmbeddingContent `json:"content"`
}

type EmbeddingContent struct {
    Parts []EmbeddingPart `json:"parts"`
}

type EmbeddingPart struct {
    Text string `json:"text"`
}
// レスポンスの構造体
type EmbeddingResponse struct {
    Embedding Embedding `json:"embedding"`
}
type Embedding struct {
    Values []float32 `json:"values"`
}
// EmbedRequest: llama-server に送信するリクエスト構造体
type EmbedRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

// --- リクエスト用構造体 ---
type EmbedContentRequest struct {
	Model   string  `json:"model"`
	Content Content `json:"content"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}
// --- レスポンス用構造体 ---
type EmbedContentResponse struct {
	Embedding Embedding `json:"embedding"`
}

// EmbedResponse: llama-server から返ってくるレスポンス構造体
// 実際のレスポンス形式は llama.cpp のバージョンにより若干異なる場合がありますが、
// 標準的な OpenAI 互換フォーマットに基づいています。
type EmbedResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Index     int       `json:"index"`
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

func GetEmbeddings(inputText string, apiKey string) ([]float32, error) {
    //var ret []float32 
	// リクエストボディの作成
	reqBody := EmbedContentRequest{
		Model: "models/gemini-embedding-001",
		Content: Content{
			Parts: []Part{
				{Text: inputText},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatalf("JSONマーシャルエラー: %v", err)
	}

	// HTTPリクエストの作成
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-embedding-001:embedContent"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("リクエスト作成エラー: %v", err)
	}

	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", apiKey)

    
	// リクエスト送信
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("HTTPリクエストエラー: %v", err)
	}
	defer resp.Body.Close()

	// レスポンスボディの読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("レスポンス読み取りエラー: %v", err)
	}

	// ステータスコードチェック
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("APIエラー (status=%d): %s", resp.StatusCode, string(body))
	}

	// JSONパース
	var embedResp EmbedContentResponse
	if err := json.Unmarshal(body, &embedResp); err != nil {
		log.Fatalf("JSONアンマーシャルエラー: %v", err)
	}

	// 結果の出力
	fmt.Printf("✅ 埋め込みベクトル取得成功\n")
	fmt.Printf("次元数: %d\n", len(embedResp.Embedding.Values))
	
	// ベクトルの先頭5要素を表示（確認用）
	displayCount := 5
	if len(embedResp.Embedding.Values) < displayCount {
		displayCount = len(embedResp.Embedding.Values)
	}
	fmt.Printf("値 (先頭%d件): %v\n", displayCount, embedResp.Embedding.Values[:displayCount])
	// 最初の要素の embedding ベクトルを返す
	// 配列が複数ある場合はインデックスを調整してください
	//return embedResp.Data[0].Embedding, nil
	return embedResp.Embedding.Values , nil
}

/**
*
* @param
*
* @return
*/
func readTextData() []ReadParam{
    fileItem := []ReadParam{}

    entries, err := os.ReadDir(DATA_DIR)
    if err != nil {
        fmt.Println("フォルダ読み込みエラー:", err)
        return nil
    }
        // textsplitter Setting
    splitter := textsplitter.NewRecursiveCharacter(
        textsplitter.WithChunkSize(CHUNK_SIZE_MAX),
        textsplitter.WithChunkOverlap(10),
    )        

    var row ReadParam
    for _, entry := range entries {
//        if entry.IsDir() ||  (filepath.Ext(entry.Name()) != ".txt")= {
        if entry.IsDir() {
            continue
        }
        if (filepath.Ext(entry.Name()) == ".txt" || filepath.Ext(entry.Name()) == ".md") {
            path := filepath.Join(DATA_DIR, entry.Name())
                    row.Name = entry.Name()

            data, err := os.ReadFile(path)
            if err != nil {
                fmt.Println("ファイル読み込みエラー:", err)
                continue
            }
            row.Content = string(data)
            // chunks add
            chunks, err := splitter.SplitText(row.Content)
            if err != nil {
                log.Fatal(err)
            }
            for i, chunk := range chunks {
                fmt.Printf("Chunk %d:\n%s\n---\n", i+1, chunk)
                row.Content = chunk
                fileItem = append(fileItem, row)
            }    
        }
        //fmt.Printf("=== %s ===\n%s\n\n", entry.Name(), string(data))
    }
    return fileItem
}


