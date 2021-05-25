package handle

const ManifestName = "imsmanifest.xml"

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

type ScormHandle interface {
    MapObjects() ([]Course, error)
}
