defmodule ChatWeb.ProtectedController do
  use ChatWeb, :controller

  plug ChatWeb.Plugs.AuthHydra when action in [:index]

  def index(conn, _params) do
    json(conn, %{message: "Acceso permitido. Token válido.", user: conn.assigns.user})
  end
end
