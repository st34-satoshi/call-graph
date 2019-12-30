package main

import (
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

func main() {
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

	//f, err := parser.ParseFile(fset, "README.md", nil, 0)
	if err != nil {
		log.Fatal("Error:", err)
	}
	for _, path := range paths{
		log.Println(path)
	}
	//log.Println(paths)

}

func parseFile(fileName string, info os.FileInfo, err error)error{
	if err != nil{
		return nil
	}
	if info.IsDir(){
		return nil
	}
	if filepath.Ext(info.Name()) != ".go"{
		return nil
	}
	log.Println(fileName)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil{
		return err
	}
	for _, importSpec := range f.Imports{
		log.Println(importSpec.Path.Value, f.Name)

	}
	return nil
}