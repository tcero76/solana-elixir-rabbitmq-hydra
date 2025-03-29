import { useState, useEffect, useRef } from 'react'
import { Socket, Channel } from 'phoenix';
import Messages from './Messages'
import { type Message } from './Types'
import { events } from './transfer.ts'


const Chat = () => {
    const socket = new Socket(`ws://${import.meta.env.VITE_API_URL}:${import.meta.env.VITE_PORT}/api/socket`,
        { params: { token: localStorage.getItem("Access_Token")}});
    const [channel, setChannel] = useState<Channel>(null);
    const [messages, setMessages] = useState<Message[]>([]);
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
            setMessages((prev) => [...prev, {user:payload.user, message:payload.body}]);
          });
    },[])
    const onClickSend = () => {
      const msg:events.Message = new events.Message({user:"Juan", text: "Hola, mundo!"});
      console.log("🚀 ~ onClickSend ~ msg:", msg.serializeBinary().buffer)
      channel?.push("message:new", msg.serializeBinary().buffer)
    }
    return (
        <>
            <Messages messages={messages}/>
            <button
              onClick={(e:React.MouseEvent<HTMLButtonElement, MouseEvent>) => onClickSend()}>
              Send
            </button>
            
        </>)
}

export default Chat;