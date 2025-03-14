import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from "react-router";
import App from './App.tsx'
import Login from './Login.tsx';
import Home from './Home'
import Pago from './Pago'

const router = createBrowserRouter([
  {
    path: "/home",
    Component: Home,
  },
  {
    path: "/login",
    Component: Login
  },
  {
    path: "/",
    Component: App
  },
  {
    path: "/pago",
    Component: Pago
  }
]);
createRoot(document.getElementById('root')!).render(
    <StrictMode>
      <RouterProvider router={router}/>
    </StrictMode>
  
)
