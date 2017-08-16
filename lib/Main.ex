defmodule Mmoserver.Main do
  use GenServer

  @tickDelay 33

  def start_link(opts \\ []) do
    GenServer.start_link(__MODULE__, [], name: Main)
  end

  def init (state) do

    IO.puts "Main Server Loop started..."

    # start the main loop, parameter is the initial tick value
    mainLoop(0)

    spawn fn -> IO.puts "got here" end

    # return, why 1??
    {:ok, 1}
  end

  def handle_data(data) do
    GenServer.cast(:main, {:handle_data, data})
  end

  def handle_info({:handle_data, data}, state) do
    # my_function(data)
    IO.puts "Got here2"
    IO.puts inspect(data)
    {:noreply, state}
  end

  # calls respective game functions
  def mainLoop(-1) do
    IO.inspect "Server Loop has ended!" # base case, end of loop
  end
 
  def mainLoop(times) do
    # do shit
    IO.inspect(times) # operation, or body of for loop

    # sleep
    :timer.sleep(@tickDelay);

    # continue the loop RECURSIVELY
    mainLoop(times + 1)
  end

end