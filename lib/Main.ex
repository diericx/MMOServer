defmodule Mmoserver.Main do
  use GenServer

  @tickDelay 33

  def start_link(opts \\ []) do
    GenServer.start_link(__MODULE__, :ok, opts)
  end

  def init (:ok) do

    IO.puts "Main Server Loop started..."

    # start the main loop, parameter is the initial tick value
    mainLoop(0)

    # return
    {:ok, 1}
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