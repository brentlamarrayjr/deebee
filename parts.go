package deebee


type Parts struct {
	m map[string][]string
	names []string
}

 type PartCallback func(name string, part []string)

func NewParts() *Parts {

	return &Parts{make(map[string][]string), make([]string, 0)}
}

func (p *Parts) Add(name string, part string) {
	p.m[name] = append(p.m[name], part)
	p.names = append(p.names, name)

}

func (p *Parts) Remove(name string, index int){

	delete(p.m, name)
	p.names = append(p.names[:index], p.names[index+1:]...)
}

func (p *Parts) Get(name string) []string {
	return p.m[name]
}

func (p *Parts) Has(name string) bool {

	 if _, has := p.m[name]; has {
	 	return true
	 }

	 return false
}

func (p *Parts) Len(name string) int {
	return len(p.m[name])
}

func (p *Parts) Each(pc PartCallback) {

	m := make(map[string][]string)



	for name, part := range p.m {
		m[name] = part
	}


	for _, name := range p.names {
		pc(name, m[name])
	}


}
