
import React, { useState, useEffect } from 'react';
import { Plus, Edit2, Trash2, ExternalLink, X, Bookmark as BookmarkIcon, Search } from 'lucide-react';
import {Link } from 'react-router-dom';
import Head from "../components/Head"
import LibConfig from "./lib/LibConfig"

//
export default function App() {
  
  const fetchTodos = async () => {
    try {
    } catch (error) {
      console.error('Error fetching todos:', error);
    }
  };

  useEffect(() => {
    const uid = localStorage.getItem(LibConfig.STORAGE_KEY_USER_ID);
    if(!uid) {
      location.href = "#/login";
    }    
  }, []);

  return (
    <div className="min-h-screen bg-neutral-50 text-neutral-900 font-sans pb-12">
      <Head />
      home
    </div>
  );
}


