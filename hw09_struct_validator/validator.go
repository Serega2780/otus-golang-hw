package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	VALIDATE        = "validate"
	LEN             = "len"
	MIN             = "min"
	MAX             = "max"
	REGEXP          = "regexp"
	IN              = "in"
	PIPE            = "|"
	COMMA           = ","
	COLON           = ":"
	BACKSLASH       = "\\"
	DOUBLEBACKSLASH = "\\\\"
)

var (
	ErrNotStruct      = errors.New("not a struct error")
	ErrUnknownTag     = errors.New("unknown tag error %s")
	ErrUnsupportedTag = errors.New("unsupported tag %s for type %s error")
	ErrTagViolation   = errors.New("tag violation error")
	supRefTypes       = supportedRefTypes{reflect.String, reflect.Int, reflect.Slice}
	supStrTags        = []string{REGEXP, LEN, IN}
	supIntTags        = []string{MIN, MAX, IN}
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

type Validation struct {
	kind  reflect.Kind
	fname string
	tags  map[string]interface{}
	value interface{}
}

type Validations []Validation

type supportedRefTypes []reflect.Kind

type supportedTypes interface {
	~int | ~string | reflect.Kind
}

func (v ValidationErrors) Error() string {
	var err error
	for _, e := range v {
		er := fmt.Errorf("%w", e.Err)
		err = errors.Join(err, er)
	}
	return err.Error()
}

func (v ValidationErrors) Len() int {
	return len(v)
}

type ThreadSafeErrors struct {
	vErrors ValidationErrors
	mx      sync.Mutex
}

func Validate(v interface{}) error {
	var wg sync.WaitGroup
	tse := &ThreadSafeErrors{vErrors: ValidationErrors{}}
	var validations Validations

	st := reflect.TypeOf(v)
	sk := st.Kind()
	if sk != reflect.Struct {
		return ErrNotStruct
	}
	sv := reflect.ValueOf(v)

	if st.NumField() == 0 {
		return nil
	}

	for i := 0; i < st.NumField(); i++ {
		tag := st.Field(i).Tag
		fieldKind := sv.Field(i).Kind()
		if !contains(supRefTypes, fieldKind) {
			continue
		}
		s, ok := tag.Lookup(VALIDATE)
		if !ok {
			continue
		}
		tm, err := getTagsMap(s, fieldKind)
		if err != nil {
			return err
		}
		if fieldKind == reflect.Slice {
			if sv.Field(i).Len() == 0 {
				continue
			}
			sliceKind := sv.Field(i).Index(0).Kind()
			tmS, err := getTagsMap(s, sliceKind)
			if err != nil {
				return err
			}
			for j := 0; j < sv.Field(i).Len(); j++ {
				if !contains(supRefTypes, sliceKind) {
					break
				}
				validations.append(sliceKind, st.Field(i).Name, tmS, sv.Field(i).Index(j).Interface())
			}
		} else {
			validations.append(sv.Field(i).Kind(), st.Field(i).Name, tm, sv.Field(i).Interface())
		}
	}
	for _, v := range validations {
		wg.Add(1)
		go validate(&wg, v, tse)
	}
	wg.Wait()

	if len(tse.vErrors) != 0 {
		return tse.vErrors
	}
	return nil
}

func (val *Validations) append(kind reflect.Kind, fname string, tags map[string]interface{}, v interface{}) {
	*val = append(*val, Validation{
		kind:  kind,
		fname: fname,
		tags:  tags,
		value: v,
	})
}

func getTagsMap(tagS string, kind reflect.Kind) (map[string]interface{}, error) {
	var err error
	result := make(map[string]interface{})
	var val interface{}
	tags := strings.Split(tagS, PIPE)
	for _, tag := range tags {
		t := strings.Split(tag, COLON)
		if (kind == reflect.Int && !contains(supIntTags, t[0])) ||
			(kind == reflect.String && !contains(supStrTags, t[0])) {
			return nil, fmt.Errorf(ErrUnsupportedTag.Error(), t[0], kind)
		}
		switch t[0] {
		case REGEXP:
			t[1] = strings.ReplaceAll(t[1], BACKSLASH, DOUBLEBACKSLASH)
			rexp, err := regexp.Compile(t[1])
			if err != nil {
				return nil, err
			}
			result[t[0]] = rexp
		case MIN, MAX, LEN:
			val, err = strconv.Atoi(t[1])
			if err != nil {
				return nil, err
			}
			result[t[0]] = val
		case IN:
			s := strings.Split(t[1], COMMA)
			if kind == reflect.String {
				val = s
			} else {
				if val, err = strToInt(s); err != nil {
					return nil, err
				}
			}
			result[t[0]] = val
		default:
			return nil, fmt.Errorf(ErrUnknownTag.Error(), t[0])
		}
	}

	return result, nil
}

func strToInt(s []string) ([]int, error) {
	var err error
	result := make([]int, len(s))
	for i, v := range s {
		if result[i], err = strconv.Atoi(v); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func contains[S ~[]T, T supportedTypes](st S, rk T) bool {
	for i := range st {
		if rk == st[i] {
			return true
		}
	}
	return false
}

func validate(wg *sync.WaitGroup, validation Validation, safeErrors *ThreadSafeErrors) {
	defer wg.Done()
	for k, v := range validation.tags {
		val := validation.value
		fname := validation.fname
		kind := validation.kind
		switch k {
		case IN:
			if (kind == reflect.String && !contains(v.([]string), fmt.Sprint(val))) ||
				(kind == reflect.Int && !contains(v.([]int), val.(int))) {
				createValErr(safeErrors, fname, k, v)
			}
		case MIN:
			if val.(int) < v.(int) {
				createValErr(safeErrors, fname, k, v)
			}
		case MAX:
			if val.(int) > v.(int) {
				createValErr(safeErrors, fname, k, v)
			}
		case LEN:
			if len(val.(string)) != v.(int) {
				createValErr(safeErrors, fname, k, v)
			}
		case REGEXP:
			if v.(*regexp.Regexp).Match([]byte(val.(string))) {
				createValErr(safeErrors, fname, k, v)
			}
		}
	}
}

func createValErr(safeErrors *ThreadSafeErrors, fname string, k string, v interface{}) {
	er := ValidationError{
		Field: fname,
		Err:   fmt.Errorf("field %s not comply with restriction %s %v %w", fname, k, v, ErrTagViolation),
	}
	safeErrors.mx.Lock()
	safeErrors.vErrors = append(safeErrors.vErrors, er)
	safeErrors.mx.Unlock()
}
