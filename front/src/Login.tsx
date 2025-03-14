import { MouseEvent } from 'react'

const Login = () => {
    const onClickLogin = (e:MouseEvent<HTMLButtonElement, globalThis.MouseEvent>) => {
        e.preventDefault()

        fetch('https://ejemplo.com/api', {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json'
            },
            body: JSON.stringify({
              nombre: 'Juan',
              edad: 30
            })
          })
            .then(response => response.json())
            .then(data => console.log(data))
            .catch(error => console.error('Error:', error));


    }
    return (<div>
                <label>user</label>
                <input type="text" ></input>
                <label>password</label>
                <input type="password"></input>
                <button onClick={(e:MouseEvent<HTMLButtonElement, globalThis.MouseEvent>) => onClickLogin(e)}>Enviar</button>
            </div>)
}

export default Login;