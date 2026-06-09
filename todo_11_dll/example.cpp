#include <iostream>
#include <cstring>
#include "include/my_todo.hpp"

extern "C" {

    __declspec(dllexport) void todo_delete(int id) {
        MyTodo todo_helper("");
        todo_helper.todo_delete_handler(id);    
    }

    __declspec(dllexport) int todo_list(char* buffer, int bufferSize) {
        MyTodo todo_helper("");
        std::string ret = todo_helper.todo_list_handler();    
        std::cout << "todo_list: " << ret << std::endl;
        const char* message = ret.c_str();
        //const char* message = ret.;
        int messageLen = std::strlen(message);

        // Go側のバッファサイズが足りているかチェック
        if (bufferSize > messageLen) {
            std::strcpy(buffer, message); // バッファにコピー
            return messageLen;            // コピーした文字数を返す
        }   
        return -1;     
    }

    __declspec(dllexport) void todo_add(const char* str) {
        if (str != nullptr) {
            MyTodo todo_helper("");
            todo_helper.todo_add_handler(std::string(str));    
            std::cout << "[C++] Go from Receive String: " << str << std::endl;
        }
    }

    // 1. GoLang から文字を 【送信】 される関数
    __declspec(dllexport) void SendStringToCpp(const char* str) {
        if (str != nullptr) {
            //std::cout << "[C++] Goから受信した文字列: " << str << std::endl;
            std::cout << "[C++] Go from Receive String: " << str << std::endl;
        }
    }

    // 2. GoLang へ文字を 【受信（返却）】 させる関数
    // Go側が用意したバッファ(buffer)に、C++の文字列を書き込みます
    __declspec(dllexport) int ReceiveStringFromCpp(char* buffer, int bufferSize) {
        const char* message = "Hello from C++ String!";
        int messageLen = std::strlen(message);

        // Go側のバッファサイズが足りているかチェック
        if (bufferSize > messageLen) {
            std::strcpy(buffer, message); // バッファにコピー
            return messageLen;            // コピーした文字数を返す
        }
        
        return -1; // バッファ不足エラー
    }
}