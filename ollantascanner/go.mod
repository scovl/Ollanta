module github.com/user/ollanta/ollantascanner

go 1.21

require (
	github.com/user/ollanta/ollantacore v0.0.0
	github.com/user/ollanta/ollantaparser v0.0.0
	github.com/user/ollanta/ollantarules v0.0.0
)

require github.com/smacker/go-tree-sitter v0.0.0-20240827094217-dd81d9e9be82 // indirect

replace (
	github.com/user/ollanta/ollantacore => ../ollantacore
	github.com/user/ollanta/ollantaparser => ../ollantaparser
	github.com/user/ollanta/ollantarules => ../ollantarules
)
