import ReactDOM from 'react-dom/client'
import React from 'react'

//import App from './App.tsx'
import App from './client/Home.tsx'

ReactDOM.createRoot(document.getElementById('app')).render(
    <App />
)
console.log('createRoot')
