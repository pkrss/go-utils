package orm

import (
	"time"
)

func ValueRaw(obj interface{})(interface {}){
	
	if obj == nil {
		return nil
	}

	switch v := obj.(type) {
	case bool,int,int8,int16,int32,int64,uint,uint8,uint16,uint32,uint64,string,float32,float64,time.Time,*time.Time:
		return v
	case JsonTime:
		if v.IsNil(){
			return nil
		}
		return v.RawValue()
	case UUID:
		if v.IsNil(){
			return nil
		}
		return v.RawValue()
	case JsonbField:
		if v.IsNil(){
			return nil
		}
		return v.RawValue()
	}
	return obj
}