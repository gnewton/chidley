package main


type Annotator interface {
     Make(fieldname string, repeats bool) string     

}

type AnnotatorImpl struct{
     AnnotationId string
     UseNameSpaceTag bool
     UseNameSpaceTagInName bool
}    