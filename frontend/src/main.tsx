import { StrictMode } from "react"
import { createRoot } from "react-dom/client"
import "./index.css"
import { BrowserRouter, Routes, Route } from "react-router-dom"
import App from "./App"
import LoginPage from "./LoginPage"
import NewMovie from "./NewMovie"
createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/login" element={<LoginPage isLogin={true} />} />
        <Route path="/register" element={<LoginPage isLogin={false} />} />
        <Route path="/admin/movies/add" element={<NewMovie />} />
      </Routes>
    </BrowserRouter>
  </StrictMode>,
)

