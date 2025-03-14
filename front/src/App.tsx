import { } from 'react'

function App() {
  const authUrl = `http://${import.meta.env.VITE_API_URL}:${import.meta.env.VITE_PORT}/hydra/oauth2/auth?client_id=657d0bb0-2314-4c8b-b649-1525af797d72&response_type=code&scope=openid offline&state=abc12asdfasfdsaf3&redirect_uri=http://${import.meta.env.VITE_API_URL}:${import.meta.env.VITE_PORT}/bff/callback`
  return (
    <div>
      Hola
      <a href={authUrl}>login</a>
    </div>
  )
}

export default App
