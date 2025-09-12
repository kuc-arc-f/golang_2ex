import ReactDOM from 'react-dom/client'
import React from 'react'


import App from './App';
import { BrowserRouter } from 'react-router-dom'

/*
function App() {
  return(
    <div>App</div>
  )
}
*/
ReactDOM.createRoot(document.getElementById('root')).render(
    <BrowserRouter>
        <App />
    </BrowserRouter>
)
console.log('#createRoot')
