defmodule ChatWeb.Channel.ChatChannel do
  use Phoenix.Channel

  def join("room:lobby", _payload, socket) do
    {:ok, socket}
  end

  def handle_in("message:new", %{"user" => _user, "body" => body}, socket) do
    IO.inspect("Entró acá")
    broadcast!(socket, "message:new", %{user: Map.get(socket.assigns.current_user, "sub"), body: body})
    {:noreply, socket}
  end
end
