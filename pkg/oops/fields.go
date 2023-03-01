package oops

type field struct {
	key   string
	value any
}

func F(key string, value any) field {
	return field{
		key:   key,
		value: value,
	}
}

// Fields takes a list of key followed by value pairs and keeps them in a fields list. There is no field deduplication.
func (e *Error) Fields(fields ...field) *Error {
	if e == nil {
		return nil
	}

	e.fields = append(e.fields, fields...)

	return e
}

func (e *Error) Field(key string, value any) *Error {
	if e == nil {
		return nil
	}

	e.fields = append(e.fields, F(key, value))

	return e
}

// FieldsMap returns the fields as a map.
func (e *Error) FieldsMap() map[string]any {
	if len(e.fields) == 0 {
		return nil
	}

	fields := make(map[string]any, len(e.fields))
	for _, f := range e.fields {
		v, ok := fields[f.key]
		if !ok {
			fields[f.key] = f.value
			continue
		}

		switch vv := v.(type) {
		case []any:
			fields[f.key] = append(vv, f.value)
		default:
			fields[f.key] = []any{vv, f.value}
		}
	}

	return fields
}

func Field(err error, key string, value any) error {
	e, _, _ := As(err)
	return e.Fields(F(key, value))
}
