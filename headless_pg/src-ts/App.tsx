import { Routes, Route } from 'react-router-dom';
import React from 'react';

import Home from './client/Home';
import Login from './client/Login';
import ContetData from './client/ContetData';

export default function App(){
  return(
  <div className="App">
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/login" element={<Login />} />
      <Route path="/data" element={<ContetData />} />
    </Routes>
  </div>
  )
}
/*
<Route path="/about" element={<About />} />
*/