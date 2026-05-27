import React from 'react'
import ReactDOM from 'react-dom/client'

import { HashRouter, Link, Route, Routes } from 'react-router-dom';
//import {createRoot} from 'react-dom/client'
//import './style.css'
import App from './App'
//import App from './client/home.tsx'
/*
const container = document.getElementById('root')
const root = createRoot(container!)
root.render(
    <HashRouter>
        <App/>
    </HashRouter>
)
*/
ReactDOM.createRoot(document.getElementById('root')).render(
  <HashRouter>
    <App />
  </HashRouter>
)
console.log('createRoot')
