package typescript

import (
	"sort"
	"strings"
)

func NewApiScriptBuilder() *ApiScriptBuilder {
	return &ApiScriptBuilder{
		Imports: make(map[string]bool),
		Content: "",
	}
}

type ApiScriptBuilder struct {
	Imports map[string]bool
	Content string
}

func (b *ApiScriptBuilder) Build(infos []ApiInfo) string {
	for _, info := range infos {
		// imports
		if info.Request != "" {
			b.Imports[strings.ReplaceAll(info.Request, "[]", "")] = true
		}
		if info.Response != "" {
			b.Imports[strings.ReplaceAll(info.Response, "[]", "")] = true
		}

		apiContent := "export const " + info.Name + " = ("

		// request variables
		for _, uri := range info.Uris {
			apiContent += strings.ReplaceAll(uri, ",omitempty", "") + ": any, "
		}
		for _, form := range info.Forms {
			apiContent += strings.ReplaceAll(form, ",omitempty", "") + ": any, "
		}
		if info.Request != "" {
			apiContent += "req: " + info.Request + ", "
		}
		apiContent = strings.TrimSuffix(apiContent, ", ")

		// response
		if info.Response != "" {
			apiContent += "): Promise<AxiosResponse<Response<" + info.Response + ">>> => {\n"
		} else {
			apiContent += "): Promise<AxiosResponse<Response<void>>> => {\n"
		}

		// request url
		var url string = info.Route
		url = strings.ReplaceAll(url, ",omitempty", "")
		if strings.Contains(url, ":") {
			elements := strings.Split(url, "/")
			for i := range elements {
				if strings.Contains(elements[i], ":") {
					url = strings.ReplaceAll(url, elements[i], "${"+elements[i][1:]+"}")
				}
			}
		}
		var query string
		if len(info.Forms) > 0 {
			query = "?"
			for i := range info.Forms {
				query += info.Forms[i] + "=${" + info.Forms[i] + "}&"
			}
		}
		query = strings.TrimSuffix(query, "&")
		query = strings.ReplaceAll(query, ",omitempty", "")

		apiContent += "    return axios." + strings.ToLower(info.Method) + "(`" + url + query + "`"
		if info.Request != "" {
			apiContent += ", req"
		}
		apiContent += ");\n"
		apiContent += "}\n\n"
		b.Content += apiContent
	}

	page := "import axios, { AxiosResponse } from 'axios';\n"
	page += "import { Response } from './general';\n"

	importKeys := make([]string, 0)
	for key := range b.Imports {
		importKeys = append(importKeys, key)
	}
	sort.Strings(importKeys)
	for _, key := range importKeys {
		page += "import { " + key + " } from './models';\n"
	}
	page += "\n"
	page += b.Content
	return page
}
