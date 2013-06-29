package forms

type buildInfo struct {
}

func Decode(dst interface{}, src map[string][]string,
	convFunc map[string]ConvertFunc) error {
	sInfo, err := Prepare(dst)
	if err != nil {
		return err
	}
	_ = sInfo
	for path, vals := range src {
		println(path, vals)
	}
	return nil
}
