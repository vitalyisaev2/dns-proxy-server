package main

import "fmt"

func main4(){

    p:= Person{Age:14}

    fmt.Printf("p=%+v\n", p)

    age := p.getAge()
    *age = 15;
    p.addChild("Ana")
    p.addChildByPointer("Carol"); // nao funciona

    name := p.getChild("Ana")

    *name = "Ana Carolina"


    fmt.Printf("p=%+v, p=%p\n", p, &p.Childs)

}

type Person struct {
    Age int;
    Childs []string
}

func (p *Person) getAge() *int {
    return &p.Age
}

func (p *Person) getChilds() *[]string {
    return &p.Childs
}

func (p *Person) getChild(name string) *string {
    for i := range p.Childs {
        fname := &p.Childs[i]
        if *fname == name {
            return fname
        }
    }
    return &p.Childs[0]
}

func (p *Person) addChild(name string) {
    childs := p.getChilds()
    *childs = append(*childs, name)
}

func (p *Person) addChildByPointer(name string) {
    childs := p.getChilds()
    newChilds := append(*childs, name)
    fmt.Printf("childs=%p, newChilds=%p\n", childs, &newChilds)
    childs = &newChilds
}
