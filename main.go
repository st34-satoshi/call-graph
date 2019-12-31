package main

import (
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2{
		log.Fatal("Error: no argument. directory name")
		return
	}
	dirName := os.Args[1]

	// find all go files
	var paths []string
	err := filepath.Walk(dirName,
		func(fileName string, info os.FileInfo, err error) error {
			if err != nil{
				log.Fatal(err)
				return err
			}
			if info.IsDir(){
				return nil
			}
			if filepath.Ext(info.Name()) != ".go"{
				return nil
			}
			paths = append(paths, fileName)
			return nil
		})
	if err != nil {
		log.Fatal("Error:", err)
	}

	// create import package map of each package
	packages := map[string]map[string]int{}
	for _, path := range paths{
		filePackage, importPaths, err := parseFile(path)
		if err != nil{
			log.Fatal(err)
			return
		}
		// remove path from this directory
		fileName := strings.Replace(filePackage, dirName, "", 1)
		importPackages, exist := packages[fileName]
		if exist{
			for _, importPath := range *importPaths{
				importPackages[importPath] = 0  // 0 is no meaning. map need key and value but I use only key.
			}
		}else{
			packages[fileName] = map[string]int{}  // initialize the filePackage value(map).
			importPackages = packages[fileName]
			for _, importPath := range *importPaths{
				importPackages[importPath] = 0  // 0 is no meaning. map need key and value but I use only key.
			}
		}
	}

	// TODO: remove or select
	// remove external package

	// output dot file for visualize using graphviz
	text := "digraph G{\n"
	for k, importPackages := range packages{
		// add to graph dot
		for importPackage, _ :=  range  importPackages{
			// remove external package, select only internal package
			// remove first directory name
			topDirectory := strings.Split(importPackage, "/")
			if len(topDirectory) <= 1{
				continue
			}
			isExternal, _ := isExternalPackage(topDirectory[0])
			if isExternal{
				continue
			}

			// add to text
			// remove top directory of importPackage
			// TODO modify: using slash is not good. windows is not used slash but yen mark
			firstSlash := strings.Index(importPackage, "/")
			importName := importPackage[firstSlash+1:]
			text += `  "` + k + `" -> "` + importName + ";\n"
			//text += `  "` + strings.Replace(k, dirName, "", 1) + `" -> ` + importPackage + ";\n"
		}
	}
	text += "}"
	log.Println(text)
	// output to text file
	dotFileName := `test3.txt`
	file, err := os.Create(dotFileName)
	if err != nil {
		log.Fatal("Error", err)
	}
	defer file.Close()
	file.Write(([]byte)(text))

	// save png file
	err = exec.Command("dot", "-T", "png", dotFileName, "-o", "test4.png").Run()
	if err != nil{
		log.Fatal(err)
	}
}

func isExternalPackage(dirName string) (bool, error){
	externalNames := []string{"github", "com", "net"}
	for _, externalName := range externalNames{
		if strings.Contains(dirName, externalName){
			return true, nil
		}
	}
	return false, nil
}

func parseFile(fileName string) (string, *[]string, error) {
	// return this file package, import package list
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil{
		return "", nil, err
	}
	var importPaths []string
	for _, importSpec := range f.Imports{
		importPaths = append(importPaths, importSpec.Path.Value)
		//log.Println(importSpec.Path.Value, f.Name)
	}
	log.Println(f.Name)
	// this directory path
	lastSlash := strings.LastIndex(fileName, "/")
	if strings.Contains(fileName[lastSlash:], "main"){
		// main package has no directory, add main directory
		return fileName[:lastSlash]+ "/main", &importPaths, nil
	}
	return fileName[:lastSlash], &importPaths, nil
}