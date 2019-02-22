package main

const generic1 = `
// Standard Chidley structs
{{ range $index, $value := . }}
   type {{$value.TypeName}} struct{       
        //Compressed 
   {{if $value.CharDataName}}           {{$value.CharDataName}} string     //CharData   {{ end }}
        //ATTS
   {{ range $value2 := $value.AttributeNames}}            {{$value2}} string //ATT   {{ end }}
        //Normal
   {{ range $value2 := $value.SubElements}}            
        {{$value2.Name}} {{$value2.TypeName}}
   {{ end }}
  }
{{ end }}

/////////////////////////////
// otira
{{ range $index, $value := . }}
   {{$index}} {{$value}}
{{ end }}
`
