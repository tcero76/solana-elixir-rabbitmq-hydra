import { useEffect, useState } from 'react'
import { Socket } from 'phoenix'
import Pago from './Pago';

const socket = new Socket(`ws://${import.meta.env.VITE_API_URL}:${import.meta.env.VITE_PORT}/api/socket`,
    { params: { token: localStorage.getItem("Access_Token")}});
const Home = () => {
    const [channel, setChannel] = useState(null);
    const [messages, setMessages] = useState("");
    useEffect(() => {
        socket.connect()
        const chan = socket.channel("room:lobby", {});
        chan
          .join()
          .receive("ok", () => {
            console.log("Conectado al canal room:lobby");
            setChannel(chan);
          })
          .receive("error", (resp) => {
            console.log("Error al conectar", resp);
          });

          chan.on("message:new", (payload) => {
            setMessages(payload.body);
          });

        fetch('/bff/getToken')
            .then(res => res.json())
            .then(res => {
                localStorage.setItem("Access_Token", res)
            })
    }, [])
    const onClickSend = () => {
        channel?.push("message:new", {
            body: "Hola, mundo!",
            user: "Juan"
        })
    }
    return (<div>
      <label>Home</label>
      <button
        onClick={(e:React.MouseEvent<HTMLButtonElement, MouseEvent>) => onClickSend()}>
          Send
      </button>
      <label>{messages}</label>
      <br/>
      <Pago/>
    </div>)
}

export default Home;