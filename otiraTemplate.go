package main

const one2manyTemplateName = "one2manyTemplate"
const many2manyTemplateName = "many2manyTemplate"
const stringFieldTemplateName = "stringFieldTemplate"

const one2manyTemplate = `

        // One2Many
        //   left={{.Left}}  right={{.Right}} counter={{.Counter}}
	func() {
                one2m{{.Counter}} := otira.NewOneToMany()
                one2m{{.Counter}}.LeftTable = {{.Left}}Table
                one2m{{.Counter}}.RightTable = {{.Right}}Table
                leftField{{.Counter}} := new(otira.FieldDefUint64)
                leftField{{.Counter}}.SetName("{{.RightSql}}")
                {{.Right}}Table.Add(leftField{{.Counter}})
                one2m{{.Counter}}.LeftKeyField = leftField{{.Counter}}
                one2m{{.Counter}}.RightKeyField = {{.Right}}Table.PrimaryKey()
                {{.Left}}Table.AddOneToMany(one2m{{.Counter}})
                
                //TODO
                //one2m{{.Counter}}.RightTableUniqueFields = []FieldDef{.......}
	}()

`

const many2manyTemplate = `

        // Many2Many
        //   left={{.Left}}  right={{.Right}} counter={{.Counter}}
	func() {
	        m2m{{.Counter}} := otira.NewManyToMany()
	        m2m{{.Counter}}.LeftTable = {{.Left}}Table
	        m2m{{.Counter}}.RightTable = {{.Right}}Table
	        {{.Left}}Table.AddManyToMany(m2m{{.Counter}})
                
                //TODO
                //m2m{{.Counter}}.RightTableUniqueFields = []FieldDef{.......}
	}()
`

const stringFieldTemplate = `

         // Field {{.FieldName}} string     
         //     {{.Comment}}
	 {{.FieldVariableName}} := new(otira.FieldDefString)
	 {{.FieldVariableName}}.SetName("{{.FieldName}}")
	 {{.FieldVariableName}}.SetLength({{.FieldLength}})
	 {{.TableVariableName}}Table.Add({{.FieldVariableName}})
 `

const uint64FieldTemplate = `

         // Field {{.FieldName}} uint64
	 {{.FieldVariableName}} := new(otira.FieldDefUint64)
	 {{.FieldVariableName}}.SetName("{{.FieldName}}")
	 {{.TableVariableName}}Table.Add({{.FieldVariableName}})
 `

const floatFieldTemplate = `

         // Field {{.FieldName}} Float
	 {{.FieldVariableName}} := new(otira.FieldDefFloat)
	 {{.FieldVariableName}}.SetName("{{.FieldName}}")
	 {{.TableVariableName}}Table.Add({{.FieldVariableName}})
 `

const boolFieldTemplate = `

         // Field {{.FieldName}} bool
	 {{.FieldVariableName}} := new(otira.FieldDefBool)
	 {{.FieldVariableName}}.SetName("{{.FieldName}}")
	 {{.TableVariableName}}Table.Add({{.FieldVariableName}})
 `

const byteFieldTemplate = `

         // Field {{.FieldName}} bool
	 {{.FieldVariableName}} := new(otira.FieldDefByte)
	 {{.FieldVariableName}}.SetName("{{.FieldName}}")
	 {{.FieldVariableName}}.SetLength({{.FieldLength}})
	 {{.TableVariableName}}Table.Add({{.FieldVariableName}})
 `

var typeMap = map[string]string{
	OFloat:  floatFieldTemplate,
	OBool:   boolFieldTemplate,
	OUint64: uint64FieldTemplate,
	OString: stringFieldTemplate,
	OByte:   byteFieldTemplate,
}
