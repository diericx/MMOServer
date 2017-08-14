defmodule Mmoserver.MessageReceiver do
  use GenServer
  require Logger

  def start_link(opts \\ []) do
    GenServer.start_link(__MODULE__, :ok, opts)
  end

  def init (:ok) do
    {:ok, _socket} = :gen_udp.open(21337)
  end

  # Handle UDP data
  def handle_info({:udp, _socket, _ip, _port, data}, state) do
    message = parse_packet(data)
    # Logger.info "Received a secret message! " <> inspect(message)
    {:noreply, state}
  end

  # Ignore everything else
  def handle_info({_, _socket}, state) do
    {:noreply, state}
  end

  def parse_packet(data) do
    # Convert data to string, then split all data
    # WARNING - SPLIT MAY BE EXPENSIVE
    dataString = Kernel.inspect(data)
    vars = String.split(dataString, ",")

    # Get variables
    packetID = Enum.at(vars, 0)
    x = Enum.at(vars, 1)

    # Do stuff with them
    IO.puts "Packet ID:"
    IO.puts packetID
    IO.puts x
  end
end