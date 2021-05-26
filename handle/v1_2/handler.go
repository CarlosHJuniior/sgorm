package v1_2

import (
    "archive/zip"
    "encoding/xml"
    "errors"
    "fmt"
    "github.com/CarlosHJuniior/sgorm/handle"
    "io/ioutil"
)

const version = "1.2"

func NewScormHandlerv1o2(sf *zip.ReadCloser) (handle.ScormHandle, error) {
    var mf *zip.File
    for _, f := range sf.File {
        if f.Name == handle.ManifestName {
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
    
    var m Manifest
    err = xml.Unmarshal(b, &m)
    if err != nil {
        return nil, err
    }
    
    if m.Version != version {
        return nil, errors.New("wrong manifest version")
    }
    
    return &scormhandlerv1o2{
        manifest: m,
        zfList:   sf.File,
    }, nil
}

type scormhandlerv1o2 struct {
    manifest Manifest
    zfList   []*zip.File
}

func (it scormhandlerv1o2) MapObjects() ([]handle.Course, error) {
    var courses []handle.Course
    
    mapa := make(map[string]Resource)
    for _, r := range it.manifest.Resources {
        mapa[r.Identifier] = r
    }
    
    for _, o := range it.manifest.Organizations {
        var course handle.Course
        
        course.ID = o.Identifier
        course.Title = o.Title
        
        for _, ip := range o.Items {
            var module handle.Module
            
            r, ok := mapa[ip.IdentifierRef]
            if !ok {
                return nil, fmt.Errorf("[%v] no one resource found", ip.IdentifierRef)
            }
    
            mflist, err := it.findResources(r)
            if err != nil {
                return nil, fmt.Errorf("[%v] error to find resource", ip.IdentifierRef)
            }
    
            for _, d := range r.Dependencies {
                r2, ok := mapa[d.IdentifierRef]
                if !ok {
                    return nil, fmt.Errorf("[%v] dependency not found", d.IdentifierRef)
                }
        
                dlist, err := it.findResources(r2)
                if err != nil {
                    return nil, fmt.Errorf("[%v] error to find resource", ip.IdentifierRef)
                }
        
                mflist = append(mflist, dlist...)
            }
    
            module.ID = ip.Identifier
            module.Title = ip.Title
            module.Files = mflist
            
            course.Modules = append(course.Modules, module)
        }
        
        courses = append(courses, course)
    }
    
    return courses, nil
}

func (it scormhandlerv1o2) findResources(r Resource) ([]handle.ModuleFile, error) {
    var list []handle.ModuleFile
    
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
                
                mf := handle.ModuleFile{
                    Name: zFile.Name,
                    Data: b,
                    IsRoot: r.HRef == rFile.Path,
                }
                
                list = append(list, mf)
            }
        }
    }
    
    return list, nil
}
