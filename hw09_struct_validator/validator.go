package hw09structvalidator

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	// Place your code here.
	rT := reflect.TypeOf(v)
	if rT.Kind() != reflect.Struct {
		return ErrMissType
	}
	if rT.NumField() == 0 {
		return ErrStructWithoutFields
	}
	val := reflect.ValueOf(v)
	var validateErrors ValidationErrors
	for idx := 0; idx <= rT.NumField()-1; idx++ {
		field := rT.Field(idx)
		tags := field.Tag
		validateParams := tags.Get("validate")
		if validateParams == "" {
			continue
		}
		validateErrors = validateStructField(field, val, validateErrors, validateParams)
	}
	return validateErrors
}

func validateStructField(field reflect.StructField,
	val reflect.Value, validateErrors ValidationErrors,
	validateParams string,
) ValidationErrors {
	var valForCheck any
	switch field.Type.Name() {
	case "int", "string":
		if field.Type.Name() == "int" {
			valForCheck = val.FieldByName(field.Name).Int()
		} else {
			valForCheck = val.FieldByName(field.Name).String()
		}
		isValidated, err := ValidateField(valForCheck, validateParams)
		if isValidated && err != nil {
			validateErrors = append(validateErrors, ValidationError{Field: field.Name, Err: err})
		}
		if !isValidated && err != nil {
			validateErrors = append(validateErrors, ValidationError{Field: field.Name, Err: err})
			return validateErrors
		}
	case "[]int", "[]string":
		fmt.Println("Start check the slice")
		sliceVal := val.FieldByName(field.Name)
		var errSlice error
		for i := 0; i < sliceVal.Len(); i++ {
			if field.Type.String() == "[]int" {
				valForCheck = sliceVal.Index(i).Int()
			} else {
				valForCheck = sliceVal.Index(i).String()
			}
			isValidated, err := ValidateField(valForCheck, validateParams)
			if !isValidated && err != nil {
				validateErrors = append(validateErrors, ValidationError{Field: field.Name, Err: errSlice})
				return validateErrors
			}
			if isValidated && err != nil {
				errSlice = errors.Join(errSlice, err)
			}
			if errSlice != nil {
				validateErrors = append(validateErrors, ValidationError{Field: field.Name, Err: errSlice})
				fmt.Println(validateErrors)
			}
		}
	}
	return validateErrors
}

func ValidateField(field any, validateParams string) (bool, error) {
	switch val := field.(type) {
	case int64:
		f, err := getValFuncInt(validateParams)
		if err != nil {
			return false, err
		}
		err = f(int(val))
		if err != nil {
			return true, err
		}

		return true, nil
	case string:
		f, err := getValFuncString(validateParams)
		if err != nil {
			return false, err
		}
		err = f(val)
		if err != nil {
			return true, err
		}
		return true, nil
	default:
		return false, nil
	}
}

func getValFuncInt(validateParams string) (func(int) error, error) {
	validateParamsArr := strings.Split(validateParams, "|")
	var resFunc func(int) error
	var f func(int) error
	for _, v := range validateParamsArr {
		switch {
		case strings.HasPrefix(v, "min:"):
			minParam, err := strconv.Atoi(v[4:])
			if err != nil {
				return nil, ErrInvalidValidateParam
			}
			f = func(val int) error {
				if val > minParam {
					return nil
				}
				return ErrValidateIntMin
			}
		case strings.HasPrefix(v, "max:"):
			maxParam, err := strconv.Atoi(v[4:])
			if err != nil {
				return nil, ErrInvalidValidateParam
			}
			f = func(val int) error {
				if val < maxParam {
					return nil
				}
				return ErrValidateIntMax
			}
		case strings.HasPrefix(v, "in:"):
			options := v[2:]
			f = func(val int) error {
				inOptions := strings.Contains(options, fmt.Sprintf(":%d,", val)) ||
					strings.Contains(options, fmt.Sprintf(",%d,", val)) ||
					strings.HasSuffix(options, fmt.Sprintf("%d", val))
				if inOptions {
					return nil
				}
				return ErrValidateIn
			}
		default:
			f = nil
		}
		if f != nil {
			if resFunc == nil {
				resFunc = f
			} else {
				dubF := resFunc
				curF := f
				resFunc = func(val int) error {
					err := dubF(val)
					errF := curF(val)
					return errors.Join(err, errF)
				}
			}
		}
	}
	return resFunc, nil
}

func getValFuncString(validateParams string) (func(val string) error, error) {
	validateParamsArr := strings.Split(validateParams, "|")
	var resFunc func(val string) error
	var f func(val string) error
	for _, v := range validateParamsArr {
		switch {
		case strings.HasPrefix(v, "len:"):
			limit, err := strconv.Atoi(v[4:])
			if err != nil {
				return nil, ErrInvalidValidateParam
			}
			f = func(val string) error {
				if len(val) == limit {
					return nil
				}
				return ErrValidateStrLen
			}
		case strings.HasPrefix(v, "regexp:"):
			reg := v[7:]
			if len(reg) == 0 {
				return nil, ErrInvalidValidateParam
			}
			f = func(val string) error {
				if FullMatchRegString(val, reg) {
					return nil
				}
				return ErrValidateStrRegexp
			}
		case strings.HasPrefix(v, "in:"):
			options := v[2:]
			f = func(val string) error {
				inOptions := strings.Contains(options, fmt.Sprintf(":%s,", val)) ||
					strings.Contains(options, fmt.Sprintf(",%s,", val)) || strings.HasSuffix(options, val)
				if inOptions {
					return nil
				}
				return ErrValidateIn
			}
		default:
			f = nil
		}
		if f != nil {
			if resFunc == nil {
				resFunc = f
			} else {
				dubF := resFunc
				curF := f
				resFunc = func(val string) error {
					err := dubF(val)
					errF := curF(val)
					return errors.Join(err, errF)
				}
			}
		}
	}
	return resFunc, nil
}

func FullMatchRegString(str string, regParam string) bool {
	reg := regexp.MustCompile(regParam)
	return reg.MatchString(str) && reg.FindString(str) == str
}
