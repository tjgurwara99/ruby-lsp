package code

type Module struct {
	Name        string
	Classes     []*Class
	Modules     []*Module
	ParentScope *Module
	Locations   []Range
}

type Location struct {
	Line      int
	Character int
	FileURI   string
}

type Range struct {
	Start Location
	End   Location
}

type Class struct {
	Name        string
	ParentScope any // its either a class or a module
	Range       Range
	Methods     []*Method
	Classes     []*Class
}

type Method struct {
	Name        string
	ParentScope any // its either a class or a module
	Range       Range
	Args        []string
}
