defmodule WsWeb.Channel.ChatChannel do
  use Phoenix.Channel

  def join("room:lobby", _payload, socket) do
    {:ok, socket}
  end

  # def handle_in("message:new", %{"user" => _user, "body" => body}, socket) do
  #   broadcast!(socket, "message:new", %{user: Map.get(socket.assigns.current_user, "sub"), body: body})
  #   {:noreply, socket}
  # end

  def handle_in(_event, {:binary, binary}, socket) do
    message =  WsWeb.Events.Message.decode(binary)
    IO.inspect(message)
    broadcast!(socket, "message:new", %{user: Map.get(socket.assigns.current_user, "sub"), body: message.text})
    {:noreply, socket}
  end
end
