package c

type Status int 

const (
	Unparsed Status  = iota
	Parsing
	Parsed
)


type FileNoder interface {
	Functions() map[string]any
	Calls() map[string]any
	Classes() map[string]any
	Consts() map[string]any
	Variables() map[string]any
	Includes() map[string]any
	Macros() map[string]any
	// parsing, parsed, unparsed
	Status() int

	// 查找文件实现 在文件对象的本地缓存空间中搜寻
	Find()

}

type File struct {
}

// Function implements FileNoder.
func (f *File) Function(index int) any {
	panic("unimplemented")
}

// Calls implements FileNoder.
func (f *File) Calls() map[string]any {
	panic("unimplemented")
}

// Classes implements FileNoder.
func (f *File) Classes() map[string]any {
	panic("unimplemented")
}

// Consts implements FileNoder.
func (f *File) Consts() map[string]any {
	panic("unimplemented")
}

// Functions implements FileNoder.
func (f *File) Functions() map[string]any {
	panic("unimplemented")
}

// Includes implements FileNoder.
func (f *File) Includes() map[string]any {
	panic("unimplemented")
}

// Macros implements FileNoder.
func (f *File) Macros() map[string]any {
	panic("unimplemented")
}

// Variables implements FileNoder.
func (f *File) Variables() map[string]any {
	panic("unimplemented")
}

var _ FileNoder = (*File)(nil)

type FolderNoder interface {
	// children []FilFileNoder
	RangeFile(func(FileNoder))
	RangeFolder(func(FolderNoder))

}

type Folder struct {
}

// RangeFile implements FolderNoder.
func (f *Folder) RangeFile(func(FileNoder)) {
	panic("unimplemented")
}

// RangeFolder implements FolderNoder.
func (f *Folder) RangeFolder(func(FolderNoder)) {
	panic("unimplemented")
}

var _ FolderNoder = (*Folder)(nil)

func buildTree(rootNode FolderNoder) (FolderNoder, error)

func parseSrcTree(root string) (FolderNoder, error) {
	virtualNode := &Folder{}
	tree, err := buildTree(virtualNode)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

func Save(key any, value any)
func GetCall(key string) Call {

}

// createLocalCache
func createLocalCache(tree FolderNoder) {
	tree.RangeFile(func(file FileNoder) {
		Save("filename", file.Functions())
		Save("filename", file.Classes())
		Save("filename", file.Consts())
		Save("filename", file.Classes())
		Save("filename", file.Macros())
		// ....
	})
	tree.RangeFolder(func(fd FolderNoder) {
		createLocalCache(fd)
	})
}


func handleThirdFunc(call any) {

}

func handleLocalFunc(call any) {

}



func link(fd FolderNoder) {
	fd.RangeFile(func(fn FileNoder) {
		if fn.Status(Parsed) {
			return
		}
		for key, function := range fn.Functions() {
			for call := range function.Calls(){
				calls := GetCall(all.Name)
				switch len(calls) {
				case 0:
					handleThirdFunc(call)
				case 1:
					handleLocalFunc(call)
				default:
					// 1 第三方库 查到的都不是
					// 2 是其中某一个
					// 3 都是（1 比对参数， 2 关联上下文猜）
					// 4 	关联 Include 信息查找
					for include := range fn.Includes() {
						if include.Status(Parsed) && include.Fund(call) {
							// 找到了一个，大概率就是这个了 写入关联关系
						}
						if include.Status(Unparsed) {
							// 开始解析这个include文件 parse 的实现是个当前函数的递归 link(virtual(fileNode))
							include.Parse()
							goto : status = parsed
						}
						if include.Status(Parsing) {
							for !include.Status(Parsed) {
								goto : status = parsed
							}
						}
					}
				}
			}
		}

	})
}

func buildGrap(tree FolderNoder) {
	// 解析完整的类型数据
	createLocalCache(tree)
	// 构建关系
	link()

}
