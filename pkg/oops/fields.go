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
	e.fields = append(e.fields, fields...)

	return e
}

// FieldsMap returns the fields as a map.
func (e *Error) FieldsMap() map[string]any {
	if len(e.fields) == 0 {
		return nil
	}

	fields := make(map[string]any, len(e.fields))
	for _, f := range e.fields {
		fields[f.key] = f.value
	}

	return fields
}
