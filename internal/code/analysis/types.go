package analysis

type Location struct {
	Line      int
	Character int
	File      string
}

type Range struct {
	Start *Location
	End   *Location
}

type Node interface {
	ParentName() string
	Identifier() string
	Type() string
	Range() *Range
}

type ClassDecl struct {
	Name       string
	Docs       string
	Parent     Node
	Children   []Node
	SuperClass *ClassDecl
	r          *Range
}

// Range implements Node.
func (c *ClassDecl) Range() *Range {
	return c.r
}

var _ Node = (*ClassDecl)(nil)

func (c *ClassDecl) Type() string {
	return "class_definition"
}

func (c *ClassDecl) Identifier() string {
	return c.Name
}

func (c *ClassDecl) ParentName() string {
	if c.Parent == nil {
		return ""
	}
	return c.Parent.Identifier()
}

type ModuleDecl struct {
	Name     string
	Docs     string
	Parent   Node
	Children []Node
	Location *Location
}

// Range implements Node.
func (m *ModuleDecl) Range() *Range {
	return &Range{
		Start: m.Location,
		End:   m.Location,
	}
}

var _ Node = (*ModuleDecl)(nil)

func (m *ModuleDecl) Type() string {
	return "module_definition"
}
func (m *ModuleDecl) ParentName() string {
	if m.Parent == nil {
		return ""
	}
	return m.Parent.Identifier()
}

func (m *ModuleDecl) Identifier() string {
	return m.Name
}

type MethodDecl struct {
	Name   string
	Parent Node
	r      *Range
	Docs   string
}

var _ Node = (*MethodDecl)(nil)

func (m *MethodDecl) Range() *Range {
	return m.r
}

func (m *MethodDecl) Type() string {
	return "method_definition"
}

func (m *MethodDecl) ParentName() string {
	return m.Parent.Identifier()
}

func (m *MethodDecl) Identifier() string {
	return m.Name
}

type SingletonMethodDecl struct {
	Name   string
	Parent Node
	r      *Range
	Docs   string
}

// Range implements Node.
func (m *SingletonMethodDecl) Range() *Range {
	return m.r
}

var _ Node = (*SingletonMethodDecl)(nil)

func (m *SingletonMethodDecl) Type() string {
	return "singleton_method_definition"
}

func (m *SingletonMethodDecl) ParentName() string {
	if m.Parent == nil {
		return ""
	}
	return m.Parent.Identifier()
}

func (m *SingletonMethodDecl) Identifier() string {
	return m.Name
}

type AliasMethodDecl struct {
	Name   string
	Docs   string
	Parent Node
	r      *Range
}

// Range implements Node.
func (m *AliasMethodDecl) Range() *Range {
	return m.r
}

var _ Node = (*AliasMethodDecl)(nil)

func (m *AliasMethodDecl) Type() string {
	return "alias_method_definition"
}

func (m *AliasMethodDecl) ParentName() string {
	if m.Parent == nil {
		return ""
	}
	return m.Parent.Identifier()
}

func (m *AliasMethodDecl) Identifier() string {
	return m.Name
}
