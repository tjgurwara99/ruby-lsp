# typed: true
# frozen_string_literal: true

class Foo::Baz
  attr_reader :value

  def initialize(value)
    @value = value
  end

  def data
    @value
  end
end
