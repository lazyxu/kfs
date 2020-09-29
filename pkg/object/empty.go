package object

var EmptyFile = &File{content: ""}
var EmptyFileHash = EmptyFile.Hash()

var EmptyDir = &Dir{
	baseObject: baseObject{
		TimeImpl: TimeImpl{},
		name:     "",
		hash:     "",
		size:     0,
		mode:     0,
	},
	items: []Object{},
}
var EmptyDirHash = EmptyFile.Hash()
