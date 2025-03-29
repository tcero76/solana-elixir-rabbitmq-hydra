import { useEffect } from 'react'
import Pago from './Pago';
import Chat from './Chat/Chat'

const Home = () => {
    useEffect(() => {
        fetch('/bff/getToken')
            .then(res => res.json())
            .then(res => {
                localStorage.setItem("Access_Token", res)
            })
    }, [])
    return (<div>
      <h1>Home</h1>
      <Chat/>
      <br/>
      <Pago/>
    </div>)
}

export default Home;