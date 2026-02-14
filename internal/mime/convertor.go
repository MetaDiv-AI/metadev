package mime

import (
	"path/filepath"
	"strings"

	"github.com/MetaDiv-AI/metadev/internal/mime/internal/extmap"
)

func NewConvertor() Convertor {
	return &ConvertorImpl{}
}

// Convertor is an interface that converts file extensions to mime types and vice versa
type Convertor interface {
	// FileExtensionToMime converts a file extension to a mime type
	FileExtensionToMime(ext string) string
	// FileNameToMime converts a file name to a mime type
	FileNameToMime(fileName string) string
	// MimeToFileExtension converts a mime type to a file extension
	MimeToFileExtension(mime string) string
}

type ConvertorImpl struct{}

// FileExtensionToMime converts a file extension to a mime type
func (c *ConvertorImpl) FileExtensionToMime(ext string) string {
	ext = strings.TrimPrefix(ext, ".")
	ext = strings.ToLower(ext)
	ext = "." + ext
	mime, ok := extmap.MapExtToMime[ext]
	if !ok {
		return ""
	}
	return mime
}

// FileNameToMime converts a file name to a mime type
func (c *ConvertorImpl) FileNameToMime(fileName string) string {
	ext := filepath.Ext(fileName)
	return c.FileExtensionToMime(ext)
}

// MimeToFileExtension converts a mime type to a file extension
func (c *ConvertorImpl) MimeToFileExtension(mime string) string {
	ext, ok := extmap.MapMimeToExt[mime]
	if !ok {
		return ""
	}
	return ext
}
