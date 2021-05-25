package entities

type Course struct {
    ID      string
    Title   string
    Modules []Module
}

type Module struct {
    ID    string
    Title string
    Files []ModuleFile
}

type ModuleFile struct {
    Name string
    Data []byte
}
