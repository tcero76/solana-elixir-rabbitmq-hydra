import { useState, useEffect } from 'react'
import { Socket, Channel } from 'phoenix';
import Messages from './Messages'
import { type Message } from './Types'
import { events } from './transfer.ts'


const Chat = () => {
    const [socket, setSocket ] = useState<Socket>(null);
    const [channel, setChannel] = useState<Channel>(null);
    const [messages, setMessages] = useState<Message[]>([]);
    useEffect(() => {
      channel?.join()
        .receive("ok", () => {
          console.log("Conectado al canal room:lobby");
          setChannel(channel);
        })
        .receive("error", (resp) => {
          console.log("Error al conectar", resp);
        });
        channel?.on("message:new", (payload) => {
          setMessages((prev) => [...prev, {user:payload.user, message:payload.body}]);
        });
        return () => {
          channel?.leave();
        }
    }, [channel])
    useEffect(() => {
      socket?.connect()
      console.log("🚀 ~ useEffect ~ socket:", socket)
      setChannel(socket?.channel("room:lobby", {}));
      return () => {
        console.log("Cerrando WebSocket...");
        socket?.disconnect();
      }
    },[socket])
    useEffect(() => {
      setSocket(new Socket(`ws://${import.meta.env.VITE_API_URL}:${import.meta.env.VITE_PORT}/api/socket`, { params: { token: localStorage.getItem("Access_Token") }}));
    },[])
    const onClickSend = () => {
      const msg:events.Message = new events.Message({user:"Juan", text: "Hola, mundo!"});
      channel?.push("message:new", msg.serializeBinary().buffer)
    }
    return (
      <>
        <h1>Chat</h1>
          <Messages messages={messages}/>
          <button
            onClick={(e:React.MouseEvent<HTMLButtonElement, MouseEvent>) => onClickSend()}>
            Send
          </button>
          
      </>)
}

export default Chat;