import React from 'react'
import ReactDOM from 'react-dom/client'
//import {createRoot} from 'react-dom/client'
import { HashRouter, Link, Route, Routes } from 'react-router-dom';
//import './style.css'
import App from './App'

/*
const container = document.getElementById('root')
const root = createRoot(container!)
root.render(
    <React.StrictMode>
        <App/>
    </React.StrictMode>
)
*/
ReactDOM.createRoot(document.getElementById('root')).render(
  <HashRouter>
    <App />
  </HashRouter>
)
console.log('createRoot')
