# typed: true
# frozen_string_literal: true

module Foo
  class Bar
    attr_reader :value

    def initialize(value)
      @value = value
    end

    def data
      @value
    end
  end
end
