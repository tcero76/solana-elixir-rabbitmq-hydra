import { type MouseEvent, useEffect, useState } from 'react'

const Login = () => {
  const [user, setUser ] = useState<string>("")
  const [ password, setPassword ] = useState<string>("")
  let loginChallenge:string = "";
  useEffect(() => {
    const params = new URLSearchParams(window.location.search);
    loginChallenge = params.get("login_challenge")??"";
  })
  const onClickLogin = (e:MouseEvent<HTMLButtonElement, globalThis.MouseEvent>) => {
      e.preventDefault()
      fetch(`http://${import.meta.env.VITE_API_URL}:${import.meta.env.VITE_PORT}/bff/login?login_challenge=${loginChallenge}`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            user,
            password
          })
        })
          .then(res => {
            window.location.href = res.url;
          })
          .catch(error => console.error('Error:', error));


  }
  return (<div>
              <label>user</label>
              <input type="text" value={user} onChange={e=>setUser(e.target.value)}></input>
              <label>password</label>
              <input type="password" value={password} onChange={e=>setPassword(e.target.value)}></input>
              <button onClick={(e:MouseEvent<HTMLButtonElement, globalThis.MouseEvent>) => onClickLogin(e)}>Enviar</button>
          </div>)
}

export default Login;