package v1_2

import "encoding/xml"

type Manifest struct {
    XMLName       xml.Name       `xml:"manifest"`
    Version       string         `xml:"metadata>schemaversion"`
    Organizations []Organization `xml:"organizations>organization"`
    Resources     []Resource     `xml:"resources>resource"`
}

type Organization struct {
    XMLName    xml.Name `xml:"organization"`
    Identifier string   `xml:"identifier,attr"`
    Title      string   `xml:"title"`
    Items      []Item   `xml:"item"`
}

type Item struct {
    XMLName       xml.Name `xml:"item"`
    Identifier    string   `xml:"identifier,attr"`
    IdentifierRef string   `xml:"identifierref,attr"`
    Title         string   `xml:"title"`
}

type Resource struct {
    XMLName      xml.Name     `xml:"resource"`
    Identifier   string       `xml:"identifier,attr"`
    Type         string       `xml:"type,attr"`
    Files        []File       `xml:"file,omitempty"`
    Dependencies []Dependency `xml:"dependency,omitempty"`
}

type File struct {
    XMLName xml.Name `xml:"file"`
    Path    string   `xml:"href,attr"`
}

type Dependency struct {
    XMLName       xml.Name `xml:"dependency"`
    IdentifierRef string   `xml:"identifierref,attr"`
}
