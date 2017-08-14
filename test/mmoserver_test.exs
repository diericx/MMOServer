defmodule MmoserverTest do
  use ExUnit.Case
  doctest Mmoserver

  test "greets the world" do
    assert Mmoserver.hello() == :world
  end
end
