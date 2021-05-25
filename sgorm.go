package sgorm

import (
    "archive/zip"
    "errors"
    "github.com/CarlosHJuniior/sgorm/handle"
    "github.com/CarlosHJuniior/sgorm/handle/v1_2"
    "github.com/CarlosHJuniior/sgorm/handle/v2004_2"
    "github.com/CarlosHJuniior/sgorm/handle/v2004_3"
)

func Unmarshal(sf *zip.ReadCloser) ([]handle.Course, error) {
    if sf == nil {
        return nil, errors.New("nil scorm package")
    }
    
    handler := getScormHandler(sf)
    if handler == nil {
        return nil, errors.New("wrong manifest")
    }
    
    l, err := handler.MapObjects()
    if err != nil {
        return nil, err
    }
    
    return l, nil
}

func getScormHandler(sf *zip.ReadCloser) handle.ScormHandle {
    handler, err := v2004_3.NewScormHandlerv2004o3(sf)
    if err == nil {
        return handler
    }
    
    handler, err = v1_2.NewScormHandlerv1o2(sf)
    if err == nil {
        return handler
    }
    
    handler, err = v2004_2.NewScormHandlerv2004o2(sf)
    if err == nil {
        return handler
    }
    
    return nil
}
