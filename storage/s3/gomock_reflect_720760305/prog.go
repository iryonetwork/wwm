
package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/golang/mock/mockgen/model"

	pkg_ "github.com/iryonetwork/wwm/storage/s3"
)

func main() {
	its := []struct{
		sym string
		typ reflect.Type
	}{
		
		{ "Storage", reflect.TypeOf((*pkg_.Storage)(nil)).Elem()},
		
		{ "KeyProvider", reflect.TypeOf((*pkg_.KeyProvider)(nil)).Elem()},
		
		{ "Minio", reflect.TypeOf((*pkg_.Minio)(nil)).Elem()},
		
	}
	pkg := &model.Package{
		// NOTE: This behaves contrary to documented behaviour if the
		// package name is not the final component of the import path.
		// The reflect package doesn't expose the package name, though.
		Name: path.Base("github.com/iryonetwork/wwm/storage/s3"),
	}

	for _, it := range its {
		intf, err := model.InterfaceFromInterfaceType(it.typ)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Reflection: %v\n", err)
			os.Exit(1)
		}
		intf.Name = it.sym
		pkg.Interfaces = append(pkg.Interfaces, intf)
	}
	if err := gob.NewEncoder(os.Stdout).Encode(pkg); err != nil {
		fmt.Fprintf(os.Stderr, "gob encode: %v\n", err)
		os.Exit(1)
	}
}
