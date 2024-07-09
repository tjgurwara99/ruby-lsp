package handlers

import (
	"testing"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/ruby"
)

func TestCompletion(t *testing.T) {
	language := ruby.GetLanguage()
	parser := sitter.NewParser()
	parser.SetLanguage(language)
	source := `require "something"

class Something
	extend T::Sig
	attr_reader :something
	sig { params(data: String).returns(Something) }
	def initialize(data, a, b)
    	@data = data
	end

    class SomethingElse
    end
end

module SomeModule
	include SomethingElse
	def method_here(some_arg)
    	a = some_arg
        some_integer a
        a
    end

    def some_integer()
    	1
    end

    class Data

    end
end`
	data, err := allIdentifiers([]byte(source), language, parser)
	if err != nil {
		t.Fatal(err)
	}
	_ = data
}
