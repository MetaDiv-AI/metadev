package typescript

import (
	"os"
	"path/filepath"
	"time"

	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Builder struct {
	Typescriptify *typescriptify.TypeScriptify
	ApiInfos      []ApiInfo
}

func NewBuilder() *Builder {
	convertor := typescriptify.New()
	convertor.CreateInterface = true
	convertor.BackupDir = ""
	convertor.ManageType(time.Time{}, typescriptify.TypeOptions{TSType: "Date", TSTransform: "new Date(__VALUE__)"})
	// Handle primitive.ObjectID as string
	convertor.ManageType(primitive.ObjectID{}, typescriptify.TypeOptions{TSType: "string"})
	// Handle map[string]any as indexed signature object
	convertor.ManageType(map[string]any{}, typescriptify.TypeOptions{TSType: "{[key: string]: any}"})
	// Handle []map[string]any as array of indexed signature objects
	convertor.ManageType([]map[string]any{}, typescriptify.TypeOptions{TSType: "{[key: string]: any}[]"})
	return &Builder{
		Typescriptify: convertor,
		ApiInfos:      make([]ApiInfo, 0),
	}
}

func (b *Builder) Build(appName string, folderPath string) {
	if folderPath == "" {
		folderPath = "./apis"
	}

	targetFolder := filepath.Join(folderPath, appName)

	// reset target folder to avoid stale files from previous runs
	os.RemoveAll(targetFolder)
	os.MkdirAll(targetFolder, os.ModePerm)

	os.WriteFile(filepath.Join(targetFolder, "general.ts"), []byte(`
export interface Response<T> {
    success: boolean;
    time: string;
    trace_id: string;
    duration: number;
    page?: Page;
    error?: string;
    data?: T;
}
export interface Page {
    page: number;
    size: number;
    total: number;
}
`), os.ModePerm)
	b.Typescriptify.ConvertToFile(filepath.Join(targetFolder, "models.ts"))
	scriptBuilder := NewApiScriptBuilder()
	script := scriptBuilder.Build(b.ApiInfos)
	os.WriteFile(filepath.Join(targetFolder, "api.ts"), []byte(script), os.ModePerm)
}
