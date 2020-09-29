package object

import "github.com/billziss-gh/cgofuse/fuse"

const DefaultDirMode = fuse.S_IFDIR | 0755
const DefaultFileMode = fuse.S_IFREG | 0644
