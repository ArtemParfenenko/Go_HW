package hw1_tree

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if e := dirTree(os.Stdout, "testdata", false); e != nil {
		panic(e)
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	e := writeFormattedNameToOutput(out, path, printFiles, "")
	return e
}

func writeFormattedNameToOutput(out io.Writer, path string, printFiles bool, parentPrefix string) error {
	selfPrefix := parentPrefix + "├───"
	childPrefix := parentPrefix + "│	"

	selfLastPrefix := parentPrefix + "└───"
	childLastPrefix := parentPrefix + "	"

	file, e := os.Open(path)
	if e != nil {
		return e
	}

	fileInfos, e := getFileInfos(file, printFiles)
	if e != nil {
		return e
	}

	for index, fileInfo := range fileInfos {
		var prefix string
		var childParentPrefix string

		if index < len(fileInfos)-1 {
			prefix = selfPrefix
			childParentPrefix = childPrefix
		} else {
			prefix = selfLastPrefix
			childParentPrefix = childLastPrefix
		}

		e = printFileInfo(out, fileInfo, prefix)
		if e != nil {
			return e
		}

		if fileInfo.IsDir() {
			e = writeFormattedNameToOutput(out, path+"/"+fileInfo.Name(), printFiles, childParentPrefix)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func getFileInfos(file *os.File, printFiles bool) ([]os.FileInfo, error) {
	fileInfos, e := file.Readdir(-1)
	if e != nil {
		return nil, e
	}
	if printFiles {
		return fileInfos, nil
	}

	dirInfos := make([]os.FileInfo, 0, len(fileInfos))
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			dirInfos = append(dirInfos, fileInfo)
		}
	}
	return dirInfos, nil
}

func printFileInfo(out io.Writer, fileInfo os.FileInfo, prefix string) error {
	output := prefix + fileInfo.Name() + getPostfix(fileInfo)
	_, e := out.Write([]byte(output))
	return e
}

func getPostfix(fileInfo os.FileInfo) string {
	postfix := ""
	if isDir := fileInfo.IsDir(); !isDir {
		if size := fileInfo.Size(); size > 0 {
			postfix = fmt.Sprintf(" (%db)", size)
		} else {
			postfix = " (empty)"
		}
	}
	return postfix + "\n"
}
