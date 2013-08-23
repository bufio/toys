package view

type ViewData map[string]interface{}

func NewViewData(title string) ViewData {
	data := make(map[string]interface{})
	data["Title"] = title
	data["Scripts"] = []string{}
	data["Styles"] = []string{}
	return data
}

func (data ViewData) SetTitle(title string) ViewData {
	data["Title"] = title
	return data
}

func (data ViewData) AddStrings(key string, paths ...string) ViewData {
	val, ok := data[key].([]string)
	if ok {
		data[key] = append(val, paths...)
	} else {
		data[key] = paths
	}
	return data
}

func (data ViewData) AddScripts(paths ...string) ViewData {
	return data.AddStrings("Scripts", paths...)
}

func (data ViewData) AddStyle(paths ...string) ViewData {
	return data.AddStrings("Styles", paths...)
}
