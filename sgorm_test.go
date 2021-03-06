package sgorm

import (
    "archive/zip"
    "log"
    "testing"
)

func TestUnmarshal(t *testing.T) {
    sf, err := zip.OpenReader("test/v1.2.zip")
    if err != nil {
        panic(err)
    }
    
    list, err := Unmarshal(sf)
    if err != nil {
        panic(err)
    }
    
    for _, c := range list {
        log.Println(c.Title)
        
        for _, m := range c.Modules {
            log.Printf("--%v", m.Title)
            
            for _, f := range m.Files {
                log.Printf("----%v [%v : %v]", f.Name, len(f.Data), f.IsRoot)
            }
        }
    }
}
