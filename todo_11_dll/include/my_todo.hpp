#pragma once
#include <iostream>
#include <fstream>
#include <vector>
#include <string>
#include <nlohmann/json.hpp> 

using json = nlohmann::json;

const std::string FILE_PATH = "todos.json";

struct Todo {
    int id;
    std::string title;
};
NLOHMANN_DEFINE_TYPE_NON_INTRUSIVE(Todo, id, title)

struct TodoData {
    int max_id;
    std::vector<Todo> items;
};

class MyTodo {
private:
    std::string m_name;

    public:
    explicit MyTodo(std::string str){}

    ~MyTodo() {}

    // JSONファイルからデータ読み込み
    TodoData load_data() {
        TodoData data{0, {}};
        std::ifstream file(FILE_PATH);
        if (file.is_open()) {
            try {
                json j = json::parse(file);
                if (j.contains("max_id") && j["max_id"].is_number()) {
                    data.max_id = j["max_id"].get<int>();
                }
                if (j.contains("items") && j["items"].is_array()) {
                    for (const auto& item : j["items"]) {
                        if (item.contains("id") && item.contains("title")) {
                            data.items.push_back({item["id"].get<int>(), item["title"].get<std::string>()});
                        }
                    }
                }
            } catch (const json::parse_error& e) {
                std::cerr << "JSON解析エラー: " << e.what() << "\n";
            }
        }
        return data;
    }

    // データをJSONファイルに保存
    void save_data(const TodoData& data) {
        json j;
        j["max_id"] = data.max_id;
        j["items"] = json::array();
        for (const auto& item : data.items) {
            j["items"].push_back({{"id", item.id}, {"title", item.title}});
        }
        std::ofstream file(FILE_PATH);
        if (file.is_open()) {
            file << j.dump(2); // インデント2で整形出力
            file.close();
        } else {
            std::cerr << "エラー: ファイルに書き込めません。\n";
        }
    }

    // TODO追加
    void add_todo(TodoData& data, const std::string& title) {
        data.max_id++;
        data.items.push_back({data.max_id, title});
        save_data(data);
        std::cout << "add: #" << data.max_id << " " << title << "\n";
    }

    // TODO一覧表示
    std::string list_todos(const TodoData& data) {
        std::string ret = "";
        std::vector<Todo> items;
        if (data.items.empty()) {
            std::cout << "TODO none\n";
            return ret;
        }
        std::string buff = "";
        for (const auto& item : data.items) {
            Todo rowData;
            rowData.id = item.id;
            rowData.title = item.title;
            items.push_back(rowData);
            std::string row = "[" + std::to_string(item.id)  + "] " + item.title + "\n" ;
            //std::cout << "[" << item.id << "] " << item.title << "\n";
            buff += row;
        }
        json j1 = items;
        std::string json_str = j1.dump();
        ret = json_str;
        return ret;
    }   
    // TODO削除
    void delete_todo(TodoData& data, int id) {
        auto it = std::remove_if(data.items.begin(), data.items.end(),
                                [id](const Todo& t) { return t.id == id; });
        if (it == data.items.end()) {
            //std::cout << "ID #" << id << " は存在しません。\n";
            std::cout << "ID #" << id << " none \n";
        } else {
            data.items.erase(it, data.items.end());
            save_data(data);
            std::cout << "delete: #" << id << "\n";
        }
    }    
    
    void todo_add_handler(std::string input_str){
        try{
            TodoData data = load_data(); 
            add_todo(data, input_str);
        } catch (const std::exception& e) {
            std::cout << "Error , main" << std::endl;
        }  
    }

    std::string todo_list_handler(){
        std::string ret = "";
        try{
            TodoData data = load_data(); 
            ret = list_todos(data);
            return ret;
        } catch (const std::exception& e) {
            std::cout << "Error , main" << std::endl;
            return ret;
        }  
    }

    void todo_delete_handler(int id){
        try{
            TodoData data = load_data();
            //int id = std::stoi(id);
            delete_todo(data, id); 
//            delete_todo(data);
        } catch (const std::exception& e) {
            std::cout << "Error , main" << std::endl;
        }  
    }

};
