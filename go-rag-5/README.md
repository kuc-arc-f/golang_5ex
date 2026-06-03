# go-rag-5

 Version: 0.9.1

 Author  :

 date    : 2026/06/02

 update :

***

Golang Wails , RAG Search SQLite

* Windows use
* embedding: gemini-embedding-001
* model: Gemma-4-E2B
* llama-server , llama.cpp

***
### vector data add

https://github.com/kuc-arc-f/golang_4ex/tree/main/go-rag-3sql

***
### Setup

* llama-server start
* port 8090: gemma-4-E2B

```
#gemma-4-E2B

/usr/local/llama-b8642/llama-server -m /var/lm_data/unsloth/gemma-4-E2B-it-Q4_K_S.gguf \
 --chat-template-kwargs '{"enable_thinking": false}' --port 8090 
```

***
### .env
```
GEMINI_API_KEY=
```
***
### related

https://wails.io/ja/

https://huggingface.co/unsloth/gemma-4-E2B-it-GGUF


***
### blog

https://zenn.dev/knaka0209/scraps/e4f3e01d0d4c58

***

