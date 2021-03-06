package lib

import (
	"strings"
	"fmt"
	"os"
	"io/ioutil"
	"text/template"
	"bytes"
)

type MergedSchema struct {
	strings.Builder
}

func (ms *MergedSchema) addIndent(n int) {
	for i := 0; i < n; i++ {
		ms.WriteString(" ")
	}
}

func (ms *MergedSchema) stitchArgument(a *Arg, l int, i int) {
	if l > 2 {
		ms.addIndent(4)
	}
	ms.WriteString(a.Param + ": ")

	if a.IsList {
		ms.WriteString("[")
		ms.WriteString(a.Type)

		if !a.Null {
			ms.WriteString("!")
		}
		ms.WriteString("]")
		if !a.IsListNull {
			ms.WriteString("!")
		}
	} else {
		ms.WriteString(a.Type)
		if a.TypeExt != nil {
			ms.WriteString(" = " + *a.TypeExt)
		}
		if !a.Null {
			ms.WriteString("!")
		}
	}

	if l <= 2 && i != l-1 {
		ms.WriteString(", ")
	}
	if l > 2 && i != l-1 {
		ms.WriteString("\n")
	}
}


func (ms *MergedSchema) GenerateTemplate(s *Schema) {
	paths := []string{
		"graphql.tmpl",
	}

	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"Title": strings.Title,
	}

	if _, err := os.Stat("graphql.tmpl"); os.IsNotExist(err) {
		// path/to/whatever does not exist
	} else {

		t := template.Must(template.New("graphql.tmpl").Funcs(funcMap).ParseFiles(paths...))
		
		var tpl bytes.Buffer
		err := t.Execute(&tpl, s)
		if err != nil {
			panic(err)
		}
		
		bs := tpl.Bytes()
		err = ioutil.WriteFile("graphql.sam", bs, 0644)
		if err != nil {
			panic(err)
		}
	}
}

func (ms *MergedSchema) GenerateNativeTemplate(s *Schema) {
	paths := []string{
		"models.tmpl",
	}

	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"Title": strings.Title,
	}

	if _, err := os.Stat("models.tmpl"); os.IsNotExist(err) {
		// path/to/whatever does not exist
	} else {

		t := template.Must(template.New("models.tmpl").Funcs(funcMap).ParseFiles(paths...))
		
		var tpl bytes.Buffer
		err := t.Execute(&tpl, s)
		if err != nil {
			panic(err)
		}
		
		bs := tpl.Bytes()
		err = ioutil.WriteFile("models.sam", bs, 0644)
		if err != nil {
			panic(err)
		}
	}
}


func (ms *MergedSchema) StitchSchema(s *Schema) string {
	numOfQurs := len(s.Queries)
	numOfMuts := len(s.Mutations)
	numOfSubs := len(s.Subscriptions)

	ms.WriteString("schema {\n")
	if numOfQurs > 0 {
		ms.addIndent(2)
		ms.WriteString("query: Query\n")
	}
	if numOfMuts > 0 {
		ms.addIndent(2)
		ms.WriteString("mutation: Mutation\n")
	}
	if numOfSubs > 0 {
		ms.addIndent(2)
		ms.WriteString("subscription: Subscription\n")
	}
	ms.WriteString("}\n\n")

	if numOfQurs > 0 {
		ms.WriteString(`type Query {
`)
		for _, q := range s.Queries {
			ms.addIndent(2)
			ms.WriteString(q.Name)
			if l := len(q.Args); l > 0 {
				ms.WriteString("(")
				if l > 2 {
					ms.WriteString("\n")
				}

				for i, a := range q.Args {
					ms.stitchArgument(a, l, i)
				}

				if l > 2 {
					ms.WriteString("\n")
					ms.addIndent(2)
				}
				ms.WriteString(")")
			}
			ms.WriteString(": ")
			if q.Resp.IsList {
				ms.WriteString("[")
			}
			ms.WriteString(q.Resp.Name)
			if !q.Resp.Null {
				ms.WriteString("!")
			}
			if q.Resp.IsList {
				ms.WriteString("]")
			}
			if q.Resp.IsList && !q.Resp.IsListNull {
				ms.WriteString("!")
			}

			if q.Directive != nil {
				ms.WriteString(" @" + q.Directive.string)
			}

			ms.WriteString("\n")
		}
		ms.WriteString("}\n\n")
	}

	if numOfMuts > 0 {
		ms.WriteString(`type Mutation {
`)
		for _, m := range s.Mutations {
			ms.addIndent(2)
			ms.WriteString(m.Name)
			if l := len(m.Args); l > 0 {
				ms.WriteString("(")
				if l > 2 {
					ms.WriteString("\n")
				}

				for i, a := range m.Args {
					ms.stitchArgument(a, l, i)
				}

				if l > 2 {
					ms.WriteString("\n")
					ms.addIndent(2)
				}
				ms.WriteString(")")
			}
			ms.WriteString(": ")
			if m.Resp.IsList {
				ms.WriteString("[")
			}
			ms.WriteString(m.Resp.Name)
			if !m.Resp.Null {
				ms.WriteString("!")
			}
			if m.Resp.IsList {
				ms.WriteString("]")
			}
			if m.Resp.IsList && !m.Resp.IsListNull {
				ms.WriteString("!")
			}

			if m.Directive != nil {
				ms.WriteString(" @" + m.Directive.string)
			}

			ms.WriteString("\n")
		}
		ms.WriteString("}\n\n")
	}

	if numOfSubs > 0 {
		ms.WriteString(`type Subscription {
`)
		for _, c := range s.Subscriptions {
			ms.addIndent(2)
			ms.WriteString(c.Name)
			if l := len(c.Args); l > 0 {
				ms.WriteString("(")
				if l > 2 {
					ms.WriteString("\n")
				}

				for i, a := range c.Args {
					ms.stitchArgument(a, l, i)
				}

				if l > 2 {
					ms.WriteString("\n")
					ms.addIndent(2)
				}
				ms.WriteString(")")
			}
			ms.WriteString(": ")
			if c.Resp.IsList {
				ms.WriteString("[")
			}
			ms.WriteString(c.Resp.Name)
			if !c.Resp.Null {
				ms.WriteString("!")
			}
			if c.Resp.IsList {
				ms.WriteString("]")
			}
			if c.Resp.IsList && !c.Resp.IsListNull {
				ms.WriteString("!")
			}

			if c.Directive != nil {
				ms.WriteString(" @" + c.Directive.string)
			}

			ms.WriteString("\n")
		}
		ms.WriteString("}\n\n")
	}

	for i, t := range s.TypeNames {
		ms.WriteString("type ")
		ms.WriteString(t.Name)
		if t.Impl {
			ms.WriteString(" implements " )
			for r := range t.ImplType {
				if (r != 0) {
					ms.WriteString(" & ")
				}
				ms.WriteString(*t.ImplType[r])
			}
		}
		ms.WriteString(" {\n")
                // var propsMethods := make(map[string]Prop{},0)
		propsMethods := map[string]*Prop{}

		if (t.Impl) {
			for _, p := range t.Props {
				propsMethods[p.Name] = p
			}
			if (t.Impl == true) {
			for r := range t.ImplType {
				implTypeObj := *t.ImplType[r]
				implObj := s.InterfacesMap[implTypeObj]
				if (implObj == nil) {
					fmt.Printf("😱 WARNING: Interface '%s' not found\n", implObj)
				} else {
					for _, p := range implObj.Props {
						if _, ok := propsMethods[p.Name]; ok {
							// It exists, already do not override the values ...
							// TODO: Check if same Type? otherwise send warning?
						} else {
							// It did nott exists lets add them
							t.Props = append(t.Props, p)
							propsMethods[p.Name] = p
						}
					}
				}
			}
			}
		}
		for _, p := range t.Props {
			ms.addIndent(2)
			ms.WriteString(p.Name)

			if l := len(p.Args); l > 0 {
				ms.WriteString("(")
				if l > 2 {
					ms.WriteString("\n")
				}
				for i, a := range p.Args {
					ms.stitchArgument(a, l, i)
				}
				if l > 2 {
					ms.WriteString("\n")
					ms.addIndent(2)
				}
				ms.WriteString(")")
			}

			ms.WriteString(": ")
			if p.IsList {
				ms.WriteString("[")
			}
			ms.WriteString(p.Type)
			if !p.Null {
				ms.WriteString("!")
			}
			if p.IsList {
				ms.WriteString("]")
			}
			if p.IsList && !p.IsListNull {
				ms.WriteString("!")
			}

			if p.Directive != nil {
				ms.WriteString(" @" + p.Directive.string)
			}

			ms.WriteString("\n")
		}
		ms.WriteString("}\n")
		if i != len(s.TypeNames)-1 {
			ms.WriteString("\n")
		}
	}
	ms.WriteString("\n")

	for i, c := range s.Scalars {
		ms.WriteString("scalar " + c.Name + "\n")
		if i != len(s.Scalars)-1 {
			ms.WriteString("\n")
		}
	}
	ms.WriteString("\n")

	for i, e := range s.Enums {
		ms.WriteString("enum " + e.Name + " {\n")
		for _, n := range e.Fields {
			ms.addIndent(2)
			ms.WriteString(n + "\n")
		}
		ms.WriteString("}\n")
		if i != len(s.Enums)-1 {
			ms.WriteString("\n")
		}
	}
	ms.WriteString("\n")

	for j, i := range s.Interfaces {
		ms.WriteString("interface " + i.Name + " {\n")

		for _, p := range i.Props {
			ms.addIndent(2)
			ms.WriteString(p.Name)

			if l := len(p.Args); l > 0 {
				ms.WriteString("(")
				if l > 2 {
					ms.WriteString("\n")
				}
				for i, a := range p.Args {
					ms.stitchArgument(a, l, i)
				}
				if l > 2 {
					ms.WriteString("\n")
					ms.addIndent(2)
				}
				ms.WriteString(")")
			}

			ms.WriteString(": ")
			if p.IsList {
				ms.WriteString("[")
			}
			ms.WriteString(p.Type)
			if !p.Null {
				ms.WriteString("!")
			}
			if p.IsList {
				ms.WriteString("]")
			}
			if p.IsList && !p.IsListNull {
				ms.WriteString("!")
			}

			if p.Directive != nil {
				ms.WriteString(" @" + p.Directive.string)
			}

			ms.WriteString("\n")
		}
		ms.WriteString("}\n")
		if j < len(s.Interfaces)-1 {
			ms.WriteString("\n")
		}
	}
	ms.WriteString("\n")

	for _, u := range s.Unions {
		ms.WriteString("union " + u.Name + " = ")
		fields := strings.Join(u.Fields, " | ")
		ms.WriteString(fields + "\n\n")
	}

	for j, i := range s.Inputs {
		ms.WriteString("input " + i.Name + " {\n")

		for _, p := range i.Props {
			ms.addIndent(2)
			ms.WriteString(p.Name + ": ")
			if p.IsList {
				ms.WriteString("[")
			}
			ms.WriteString(p.Type)
			if !p.Null {
				ms.WriteString("!")
			}
			if p.IsList {
				ms.WriteString("]")
			}
			if p.IsList && !p.IsListNull {
				ms.WriteString("!")
			}

			if p.Directive != nil {
				ms.WriteString(" @" + p.Directive.string)
			}

			ms.WriteString("\n")
		}

		ms.WriteString("}\n")
		if j < len(s.Inputs)-1 {
			ms.WriteString("\n")
		}
	}

	return ms.String()
}
