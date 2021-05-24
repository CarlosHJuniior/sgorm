package sgorm

import (
    "encoding/xml"
    "sgorm/entities"
)

func Unmarshal(b []byte, c *entities.Course) error {
    var m entities.Manifest
    
    err := xml.Unmarshal(b, &m)
    if err != nil {
        return err
    }
    
    convertObject(m, c)
    return nil
}

func convertObject(m entities.Manifest, c *entities.Course) {

}
