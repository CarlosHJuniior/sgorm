package sgorm

import (
    "encoding/xml"
    "errors"
    "github.com/CarlosHJuniior/sgorm/entities"
)

func Unmarshal(b []byte, c *entities.Course) error {
    var m entities.Manifest
    
    err := xml.Unmarshal(b, &m)
    if err != nil {
        return errors.New("wrong manifest")
    }
    
    convertObject(m, c)
    return nil
}

func convertObject(m entities.Manifest, c *entities.Course) {

}
