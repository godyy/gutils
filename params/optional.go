package params

func Optional[Param any](params ...Param) (p Param, ok bool) {
	if len(params) > 0 {
		p = params[0]
		ok = true
	}
	return
}

func OptionalDefault[Param any](value Param, params ...Param) Param {
	if len(params) > 0 {
		return params[0]
	} else {
		return value
	}
}
