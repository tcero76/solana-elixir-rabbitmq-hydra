defmodule WsWeb.MQConsumer do
  use GenServer
  use AMQP

  @queue System.get_env("RABBITMQ_QUEUE") || "default_queue"
  @rabbitmq_host System.get_env("RABBITMQ_HOST") || "localhost"

  def start_link(_opts) do
    GenServer.start_link(__MODULE__, %{}, name: __MODULE__)
  end

  def init(state) do
    case Connection.open("amqp://guest:guest@#{@rabbitmq_host}") do
      {:ok, conn} ->
        {:ok, chan} = Channel.open(conn)
        Queue.declare(chan, @queue, durable: true)
        Basic.consume(chan, @queue, nil, no_ack: true)
        IO.puts("AMQP: conectado al canal #{inspect(chan)}")
        {:ok, %{channel: chan}}

      {:error, _} ->
        Process.send_after(self(), :retry_connect, 5000)
        IO.puts("AMQP: error al conectar al canal")
        {:ok, state}
    end
  end

  def handle_info({:basic_deliver, message, _metadata} = _msg, state) do
    IO.puts("El mensaje recibido es: #{inspect(message)}")
    WsWeb.Endpoint.broadcast("room:lobby", "message:new", %{tipo: "notificacion", body: message})
    {:noreply, state}
  end

  def handle_info({:basic_consume_ok, %{consumer_tag: tag}}, state) do
    IO.puts("Confirmación de consumo recibido: #{inspect(tag)}")
    {:noreply, state}
  end

  def handle_info(msg, state) do
    IO.puts("Mensaje inesperado: #{inspect(msg)}")
    {:noreply, state}
  end
end
