package main

import (
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
)

func main() {
	err := filepath.Walk("/Users/satoshi/datumix/optimus-backend/", parseFile)

	//f, err := parser.ParseFile(fset, "README.md", nil, 0)
	if err != nil {
		log.Fatal("Error:", err)
	}
	//fmt.Println(f)

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
	_, err = parser.ParseFile(fset, fileName, nil, 0)
	//fmt.Println(f, err)
	return err
}