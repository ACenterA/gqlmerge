package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func Merge(paths ...string) *string {
	schemas := make([]Schema, 0, len(paths))

	for _, path := range paths {
		if sc := getSchema(path); sc != nil {
			schemas = append(schemas, *sc)
		}
	}

	if len(schemas) == 0 {
		return nil
	}

	schema := joinSchemas(schemas)
	ms := MergedSchema{}
	ss := ms.StitchSchema(schema)
	ms.GenerateTemplate(schema)
	ms.GenerateNativeTemplate(schema)

	return &ss
}

func getSchema(path string) *Schema {
	abs, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	sc := &Schema{}
	// at this moment, path should be an absolute path
	sc.GetSchema(abs)

	if len(sc.Files) == 0 {
		return nil
	}

	for _, file := range sc.Files {
		l := NewLexer(file)
		sc.ParseSchema(l)
	}

	return sc
}

func joinSchemas(schemas []Schema) *Schema {
	schema := Schema{}

	for _, s := range schemas {
		schema.Files = append(schema.Files, s.Files...)
		schema.Mutations = append(schema.Mutations, s.Mutations...)
		schema.Queries = append(schema.Queries, s.Queries...)
		schema.Subscriptions = append(schema.Subscriptions, s.Subscriptions...)
		schema.TypeNames = append(schema.TypeNames, s.TypeNames...)
		schema.NativeTypeNames = append(schema.NativeTypeNames, s.NativeTypeNames...)
		schema.Scalars = append(schema.Scalars, s.Scalars...)
		schema.Enums = append(schema.Enums, s.Enums...)
		schema.Interfaces = append(schema.Interfaces, s.Interfaces...)
		schema.Unions = append(schema.Unions, s.Unions...)
		schema.Inputs = append(schema.Inputs, s.Inputs...)
	}

	wg := sync.WaitGroup{}
	wg.Add(9)

	go schema.UniqueMutation(&wg)
	go schema.UniqueQuery(&wg)
	go schema.UniqueTypeName(&wg)
	go schema.UniqueNativeTypeName(&wg)
	go schema.UniqueScalar(&wg)
	go schema.UniqueEnum(&wg)
	go schema.UniqueInterface(&wg)
	go schema.UniqueUnion(&wg)
	go schema.UniqueInput(&wg)

	wg.Wait()

	return &schema
}
