//import { Routes, Route, Link } from 'react-router-dom';
import {Link } from 'react-router-dom';
import React, { useState , useEffect } from 'react';

function Page() {
  return (
  <div>
    <Link to="/" className="font-bold ms-4" >Home</Link>
    <Link to="/login" className="ms-4" >Login</Link>
    <Link to="/signup" className="ms-4" >Signup</Link>
    <hr className="my-2" />
  </div>
  );
}
//    <Link to="/post" className="ms-4" >Post</Link>
export default Page;
