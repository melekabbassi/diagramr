package graph

type DiagramGraph struct {
	Nodes    map[string]*Node `json:"nodes"` // key: fully-qualified ID
	Edges    []Edge           `json:"edges"`
	Metadata Metadata         `json:"metadata"`
}

type Node struct {
	ID         string   `json:"id"`
	Label      string   `json:"label"`
	Kind       NodeKind `json:"kind"`
	Methods    []Method `json:"methods"`
	Fields     []Field  `json:"fields"`
	Package    string   `json:"package"`
	SourceFile string   `json:"source_file"`
	SourceLine int      `json:"source_line"`
}

type NodeKind string
type Relation string
type Visibility string

const (
	KindStruct    NodeKind = "struct"
	KindInterface NodeKind = "interface"
	KindService   NodeKind = "service" // inferred via naming convention
)

const (
	RelEmbeds     Relation = "embeds"     // Go struct embedding
	RelImplements Relation = "implements" // satisfies interface (inferred)
	RelUses       Relation = "uses"       // field of that type
	RelImports    Relation = "imports"    // package-level import
)

const (
	VisPublic  Visibility = "public"  // exported (capital letter)
	VisPrivate Visibility = "private" // unexported (lower case letter)
)

type Method struct {
	Name       string     `json:"name"`
	Params     []Param    `json:"params"`
	Returns    []string   `json:"returns"`
	Visibility Visibility `json:"visibility"`
	IsPointer  bool       `json:"is_pointer"` // pointer receiver
}

type Field struct {
	Name       string     `json:"name"`
	Type       string     `json:"type"`
	Visibility Visibility `json:"visibility"`
	IsEmbedded bool       `json:"is_embedded"`
	Tag        string     `json:"tag,omitempty"`
}

type Edge struct {
	From     string   `json:"from"`
	To       string   `json:"to"`
	Relation Relation `json:"relation"`
	Label    string   `json:"label,omitempty"`
}

type Param struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Metadata struct {
	Language string `json:"language"`
	RootPath string `json:"root_path"`
	ParsedAt string `json:"parsed_at"`
}
