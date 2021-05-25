package sgorm

import (
    "archive/zip"
    "encoding/xml"
    "errors"
    "fmt"
    "io/ioutil"
    "path/filepath"
    "sgorm/entities"
)

const (
    manifestName = "imsmanifest.xml"
)

func Unmarshal(sf *zip.ReadCloser) ([]entities.Course, error) {
    if sf == nil {
        return nil, errors.New("nil scorm package")
    }
    
    sh, err := readManifest(sf)
    if err != nil {
        return nil, errors.New("wrong manifest")
    }
    
    l, err := sh.mapObjects()
    if err != nil {
        return nil, err
    }
    
    return l, nil
}

func readManifest(sf *zip.ReadCloser) (*scormHandler, error) {
    var mf *zip.File
    for _, f := range sf.File {
        if f.Name == manifestName {
            mf = f
            break
        }
    }
    
    if mf == nil {
        return nil, errors.New("manifest not found")
    }
    
    rc, err := mf.Open()
    if err != nil {
        return nil, err
    }
    
    defer rc.Close()
    
    b, err := ioutil.ReadAll(rc)
    if err != nil {
        return nil, err
    }
    
    var m entities.Manifest
    err = xml.Unmarshal(b, &m)
    if err != nil {
        return nil, err
    }
    
    return &scormHandler{
        manifest: m,
        zfList:   sf.File,
    }, nil
}

type scormHandler struct {
    manifest entities.Manifest
    zfList   []*zip.File
}

func (it scormHandler) mapObjects() ([]entities.Course, error) {
    var courses []entities.Course
    
    mapa := make(map[string]entities.Resource)
    for _, r := range it.manifest.Resources {
        mapa[r.Identifier] = r
    }
    
    for _, o := range it.manifest.Organizations {
        var course entities.Course
        
        course.ID = o.Identifier
        course.Title = o.Title
        
        for _, ip := range o.ItemParents {
            var module entities.Module
            
            module.ID = ip.Identifier
            module.Title = ip.Title
            
            for _, i := range ip.Items {
                r, ok := mapa[i.IdentifierRef]
                if !ok {
                    return nil, fmt.Errorf("[%v] no one resource found", i.IdentifierRef)
                }
                
                mflist, err := it.findResources(r)
                if err != nil {
                    return nil, fmt.Errorf("[%v] error to find resource", i.IdentifierRef)
                }
                
                for _, d := range r.Dependencies {
                    r2, ok := mapa[d.IdentifierRef]
                    if !ok {
                        return nil, fmt.Errorf("[%v] dependency not found", d.IdentifierRef)
                    }
                    
                    dlist, err := it.findResources(r2)
                    if err != nil {
                        return nil, fmt.Errorf("[%v] error to find resource", i.IdentifierRef)
                    }
                    
                    mflist = append(mflist, dlist...)
                }
                
                module.Files = mflist
            }
            
            course.Modules = append(course.Modules, module)
        }
        
        courses = append(courses, course)
    }
    
    return courses, nil
}

func (it scormHandler) findResources(r entities.Resource) ([]entities.ModuleFile, error) {
    var list []entities.ModuleFile
    
    for _, rFile := range r.Files {
        for _, zFile := range it.zfList {
            if rFile.Path == zFile.Name {
                rc, err := zFile.Open()
                if err != nil {
                    return nil, err
                }
                
                b, err := ioutil.ReadAll(rc)
                rc.Close()
                
                if err != nil {
                    return nil, err
                }
                
                _, name := filepath.Split(zFile.Name)
                
                mf := entities.ModuleFile{
                    Name: name,
                    Data: b,
                }
                
                list = append(list, mf)
            }
        }
    }
    
    return list, nil
}
