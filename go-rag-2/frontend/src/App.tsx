import React, { useState , useEffect } from 'react';
import { marked } from 'marked';
//import Head from '../components/Head';
import {SendMessage} from "../wailsjs/go/main/App";

interface Todo {
  id: string;
  title: string;
  completed: boolean;
}

export default function App() {
  const [text, setText] = useState<string>("");
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
  }, []);

  const chatStart = async function(){
    try{    
      setText("");
      setIsLoading(false);
      const elem = document.getElementById("input_text") as HTMLInputElement;
      let inText = "";
      if(elem){
        inText = elem.value;
      };
      console.log("inText=", inText);
      if(!inText){ return; }
      setIsLoading(true);
      const item = {
        input: inText ,
      };
      const inData = {
        query: inText.trim(),
      };
      const target = {
        action: "rag_search",
        data: JSON.stringify(inData)
      }
      const sendJson = JSON.stringify(target)   
      console.log(sendJson)      
      SendMessage(sendJson).then((result) => {
        console.log("result=", result);
        const j1 = JSON.parse(result)
        console.log(j1);
        if(j1.Ret === 200){
          console.log(j1.Data);
          const targetHtml = marked.parse(j1.Data);
          setText(targetHtml)   
          setIsLoading(false)       
        }
      }).catch((err) => { console.error(err);});          


      return;
    } catch(e){
      console.error(e);
    }
  }

  return (
  <div className="min-h-screen bg-slate-50 py-12 px-4 sm:px-6 lg:px-8 font-sans">
    <div className="flex flex-col w-full max-w-3xl py-4 mx-auto gap-4">
      <div className="flex flex-col gap-2 px-4 bg-white">
        <h1 className="text-2xl font-bold">RAG-Chat</h1>
        <input
          id="input_text"
          type="text"
          defaultValue=""
          className="w-full p-2 border border-gray-300 rounded dark:disabled:bg-gray-700"
          placeholder="Type your message..."
        />
        <button
          type="button"
          className="px-4 py-2 text-white bg-gray-600 rounded hover:bg-gray-700 disabled:bg-gray-700"
          onClick={()=>{chatStart()}}
        > GO
        </button>

      </div>
      <div>
        <div dangerouslySetInnerHTML={{ __html: text }} id="get_text_wrap"
          className="mb-8 p-2 bg-gray-100" />
        {isLoading ? (
          <div 
          className="animate-spin rounded-full h-8 w-8 mx-4 border-t-4 border-b-4 border-blue-500">
          </div>
        ): null}
      </div>

    </div>
  </div>

  );
}
