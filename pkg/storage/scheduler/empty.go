package scheduler

import "github.com/lazyxu/kfs/object"

var EmptyFile = &object.File{Content: ""}
var EmptyFileHash = EmptyFile.Hash()
