package util

func FixedJsonFloat(v interface{}) interface{} {
	if v == nil {
		return v
	}

	switch v.(type) {
	case float64:
		{
			val, _ := v.(float64)
			if val-float64(int(val)) < 0.0001 {
				return int64(val)
			} else {
				return v
			}
		}
	case float32:
		{
			val, _ := v.(float32)
			if val-float32(int(val)) < 0.0001 {
				return int32(val)
			} else {
				return v
			}
		}
	case map[string]interface{}:
		{
			val, _ := v.(map[string]interface{})
			ret := make(map[string]interface{})
			for _k, _v := range val {
				ret[_k] = FixedJsonFloat(_v)
			}
			return ret
		}
	case []interface{}:
		{
			val, _ := v.([]interface{})
			ret := []interface{}{}
			for _, _v := range val {
				ret = append(ret, FixedJsonFloat(_v))
			}
			return ret
		}
	}
	return v
}

func FixJsonMapFloat(data map[string]interface{}) map[string]interface{} {
	if len(data) == 0 {
		return data
	}

	ret := make(map[string]interface{})
	for k, v := range data {
		ret[k] = FixedJsonFloat(v)
	}
	return ret
}
