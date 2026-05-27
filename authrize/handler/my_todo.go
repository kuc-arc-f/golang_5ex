package handler

import (
	"encoding/json"
	"fmt"
	"os"
	//"strconv"
	//"text/tabwriter"
)

const dataFile = "/tmp/todo.json"

type Item struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type Data struct {
	Items []Item `json:"items"`
	MaxID int    `json:"max_id"`
}

// ---------- ファイル操作 ----------

func load() (*Data, error) {
	f, err := os.ReadFile(dataFile)
	if os.IsNotExist(err) {
		return &Data{}, nil
	}
	if err != nil {
		return nil, err
	}
	var d Data
	if err := json.Unmarshal(f, &d); err != nil {
		return nil, err
	}
	return &d, nil
}

func save(d *Data) error {
	b, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, b, 0644)
}


func CmdAdd(title string) error {
	d, err := load()
	if err != nil {
		return err
	}
	d.MaxID++
	d.Items = append(d.Items, Item{ID: d.MaxID, Title: title})
	if err := save(d); err != nil {
		return err
	}
	fmt.Printf("追加しました: [%d] %s\n", d.MaxID, title)
	return nil
}

// example構造体を格納する配列の型
type ListItems []Item

func CmdList() string {
	var arr ListItems
	d, err := load()
	if err != nil {
		return "error , load"
	}
	if len(d.Items) == 0 {
		fmt.Println("TODOはありません。")
		return "error , d.Items=0"
	}
	for _, item := range d.Items {
		var row Item
		row.ID = item.ID;
		row.Title = item.Title;
		arr = append(arr, row)
		//fmt.Fprintf(w, "%d\t%s\n", item.ID, item.Title)
	}
	json3, err := json.Marshal(arr)
	if err != nil {
			fmt.Println(err)
			return "error, json3 Marshal"
	}	
	//fmt.Println("json3=%s \n" , string(json3))
	return string(json3)
}


func CmdDelete(id int) error {
	d, err := load()
	if err != nil {
		return err
	}
	idx := -1
	for i, item := range d.Items {
		if item.ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("ID %d が見つかりません", id)
	}
	removed := d.Items[idx]
	d.Items = append(d.Items[:idx], d.Items[idx+1:]...)
	if err := save(d); err != nil {
		return err
	}
	fmt.Printf("削除しました: [%d] %s\n", removed.ID, removed.Title)
	return nil
}