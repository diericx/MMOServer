# MMOSERVER.EX
defmodule Mmoserver do
  use Application

  def start(_type, _args) do
    import Supervisor.Spec, warn: false

    IO.puts "Listening for packets..."

    children = [
      # We will add our children here later
      
      worker(Mmoserver.Main, []),
      worker(Mmoserver.MessageReceiver, [])
    ]

    # Start the main supervisor, and restart failed children individually
    opts = [strategy: :one_for_one, name: AcmeUdpLogger.Supervisor]
    Supervisor.start_link(children, opts)
  end

end
