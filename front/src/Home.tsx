import { useEffect } from 'react'

const Home = () => {
    useEffect(() => {
      fetch('/bff/getToken')
      .then(res => res.json())
      .then(res => {
        localStorage.setItem("Access_Token", res)
      })
    },[])
    return (<div>
      <ol>
        <a href='./chat'>Chat</a>
        <a href='./pago'>Pago</a>
      </ol> 
    </div>)
}

export default Home;