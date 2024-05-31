import React from 'react'
import ReactDOM from 'react-dom/client'
import "./shared.sass"
import "./transitions.css"
import Routing from "./routing.tsx";



ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <Routing/>
  </React.StrictMode>
)
