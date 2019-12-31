package main

import (
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// find all go files
	var paths []string
	err := filepath.Walk("",
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
		importPackages, exist := packages[filePackage]
		if exist{
			for _, importPath := range *importPaths{
				importPackages[importPath] = 0  // 0 is no meaning. map need key and value but I use only key.
			}
		}else{
			packages[filePackage] = map[string]int{}  // initialize the filePackage value(map).
			importPackages = packages[filePackage]
			for _, importPath := range *importPaths{
				importPackages[importPath] = 0  // 0 is no meaning. map need key and value but I use only key.
			}
		}
	}
	//log.Println(paths)
	for k, v := range packages{
		log.Println("file package ",k)
		for i, _ := range v{
			log.Println(i)
		}
	}

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
	return f.Name.Name, &importPaths, nil
}