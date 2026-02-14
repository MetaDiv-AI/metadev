package extmap

var MapExtToMime = map[string]string{
	// Text files
	".txt":  "text/plain",
	".html": "text/html",
	".htm":  "text/html",
	".css":  "text/css",
	".csv":  "text/csv",
	".xml":  "text/xml",
	".json": "application/json",

	// Image files
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".bmp":  "image/bmp",
	".webp": "image/webp",
	".svg":  "image/svg+xml",
	".ico":  "image/x-icon",

	// Audio files
	".mp3": "audio/mpeg",
	".wav": "audio/wav",
	".ogg": "audio/ogg",
	".oga": "audio/ogg",
	".m4a": "audio/mp4",

	// Video files
	".mp4":  "video/mp4",
	".webm": "video/webm",
	".avi":  "video/x-msvideo",
	".mov":  "video/mp4",

	// Document files
	".pdf":  "application/pdf",
	".doc":  "application/msword",
	".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".xls":  "application/vnd.ms-excel",
	".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	".ppt":  "application/vnd.ms-powerpoint",
	".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",

	// Archive files
	".zip": "application/zip",
	".gz":  "application/gzip",
	".tar": "application/x-tar",
	".7z":  "application/x-7z-compressed",

	// Application files
	".js":   "application/javascript",
	".mjs":  "text/javascript",
	".wasm": "application/wasm",

	// Font files
	".ttf":   "font/ttf",
	".otf":   "font/otf",
	".woff":  "font/woff",
	".woff2": "font/woff2",

	// Additional Text files
	".rtf":  "application/rtf",
	".md":   "text/markdown",
	".yaml": "text/yaml",
	".yml":  "text/yaml",
	".log":  "text/plain",
	".ini":  "text/plain",
	".conf": "text/plain",
	".ts":   "application/typescript",
	".jsx":  "text/jsx",
	".tsx":  "text/tsx",

	// Additional Image files
	".tiff": "image/tiff",
	".tif":  "image/tiff",
	".raw":  "image/x-raw",
	".psd":  "image/vnd.adobe.photoshop",
	".ai":   "application/postscript",
	".eps":  "application/postscript",

	// Additional Audio files
	".midi": "audio/midi",
	".mid":  "audio/midi",
	".aac":  "audio/aac",
	".wma":  "audio/x-ms-wma",
	".ra":   "audio/x-realaudio",
	".flac": "audio/flac",

	// Additional Video files
	".mkv":  "video/x-matroska",
	".wmv":  "video/x-ms-wmv",
	".flv":  "video/x-flv",
	".3gp":  "video/3gpp",
	".m4v":  "video/x-m4v",
	".mpg":  "video/mpeg",
	".mpeg": "video/mpeg",

	// Additional Document files
	".odt":   "application/vnd.oasis.opendocument.text",
	".ods":   "application/vnd.oasis.opendocument.spreadsheet",
	".odp":   "application/vnd.oasis.opendocument.presentation",
	".epub":  "application/epub+zip",
	".tex":   "application/x-tex",
	".latex": "application/x-latex",

	// Additional Archive files
	".rar": "application/vnd.rar",
	".bz":  "application/x-bzip",
	".bz2": "application/x-bzip2",
	".cab": "application/vnd.ms-cab-compressed",
	".iso": "application/x-iso9660-image",

	// Additional Application files
	".swf":  "application/x-shockwave-flash",
	".jar":  "application/java-archive",
	".war":  "application/java-archive",
	".php":  "application/x-httpd-php",
	".py":   "text/x-python",
	".rb":   "text/x-ruby",
	".sh":   "application/x-sh",
	".bash": "application/x-sh",
	".sql":  "application/sql",
	".apk":  "application/vnd.android.package-archive",
	".exe":  "application/x-msdownload",
	".dll":  "application/x-msdownload",
	".deb":  "application/x-debian-package",
	".rpm":  "application/x-rpm",

	// Additional Font files
	".eot": "application/vnd.ms-fontobject",

	// Email files
	".eml": "message/rfc822",
	".msg": "application/vnd.ms-outlook",

	// CAD and 3D files
	".stl": "model/stl",
	".obj": "model/obj",
	".fbx": "application/octet-stream",

	// Scientific/Technical
	".mat": "application/x-matlab-data",
	".nb":  "application/mathematica",
	".r":   "text/x-r",

	// Additional Document/Text files
	".n3":   "text/n3",
	".yang": "application/yang",
	".yin":  "application/yin+xml",
	".ttl":  "text/turtle",
	".sgml": "text/sgml",
	".etx":  "text/x-setext",
	".fly":  "text/vnd.fly",
	".par":  "text/plain-bas",
	".uri":  "text/uri-list",
	".uris": "text/uri-list",
	".urls": "text/uri-list",
	".vcs":  "text/x-vcalendar",
	".wml":  "text/vnd.wap.wml",
	".wmls": "text/vnd.wap.wmlscript",
	".s":    "text/x-asm",
	".c":    "text/x-c",
	".f":    "text/x-fortran",
	".java": "text/x-java-source,java",
	".p":    "text/x-pascal",
	".uu":   "text/x-uuencode",
	".vcf":  "text/x-vcard",

	// Additional Chemical/Scientific formats
	".xyz":  "chemical/x-xyz",
	".cdx":  "chemical/x-cdx",
	".cif":  "chemical/x-cif",
	".cmdf": "chemical/x-cmdf",
	".cml":  "chemical/x-cml",
	".csml": "chemical/x-csml",

	// Additional Model formats
	".mts":  "model/vnd.mts",
	".vtu":  "model/vnd.vtu",
	".wrl":  "model/vrml",
	".iges": "model/iges",
	".igs":  "model/iges",
	".mesh": "model/mesh",
	".msh":  "model/mesh",
	".silo": "model/mesh",
	".dwf":  "model/vnd.dwf",
	".gdl":  "model/vnd.gdl",
	".gtw":  "model/vnd.gtw",

	// Additional Application formats
	".zaz":      "application/vnd.zzazz.deck+xml",
	".zmm":      "application/vnd.handheld-entertainment+xml",
	".zir":      "application/vnd.zul",
	".xop":      "application/xop+xml",
	".xer":      "application/patch-ops-error+xml",
	".wspolicy": "application/wspolicy+xml",
	".davmount": "application/davmount+xml",
	".ccxml":    "application/ccxml+xml",
	".cdmia":    "application/cdmi-capability",
	".cdmic":    "application/cdmi-container",
	".cdmid":    "application/cdmi-domain",
	".cdmio":    "application/cdmi-object",
	".cdmiq":    "application/cdmi-queue",

	// Additional Image formats
	".ktx":  "image/ktx",
	".sub":  "image/vnd.dvb.subtitle",
	".dwg":  "image/vnd.dwg",
	".dxf":  "image/vnd.dxf",
	".fst":  "image/vnd.fst",
	".fbs":  "image/vnd.fastbidsheet",
	".fpx":  "image/vnd.fpx",
	".npx":  "image/vnd.net-fpx",
	".wbmp": "image/vnd.wap.wbmp",
	".xif":  "image/vnd.xiff",
	".mac":  "image/x-macpaint",
	".pic":  "image/x-pict",
	".pct":  "image/x-pict",
	".pnm":  "image/x-portable-anymap",
	".pbm":  "image/x-portable-bitmap",
	".pgm":  "image/x-portable-graymap",
	".ppm":  "image/x-portable-pixmap",
	".rgb":  "image/x-rgb",
	".xbm":  "image/x-xbitmap",
	".xpm":  "image/x-xpixmap",
	".xwd":  "image/x-xwindowdump",

	// Additional Audio formats
	".adp":   "audio/adpcm",
	".au":    "audio/basic",
	".snd":   "audio/basic",
	".kar":   "audio/midi",
	".rmi":   "audio/midi",
	".mp4a":  "audio/mp4",
	".eol":   "audio/vnd.digital-winds",
	".dra":   "audio/vnd.dra",
	".dts":   "audio/vnd.dts",
	".dtshd": "audio/vnd.dts.hd",
	".lvp":   "audio/vnd.lucent.voice",
	".pya":   "audio/vnd.ms-playready.media.pya",
	".aif":   "audio/x-aiff",
	".aiff":  "audio/x-aiff",
	".aifc":  "audio/x-aiff",
	".m3u":   "audio/x-mpegurl",
	".ram":   "audio/x-pn-realaudio",

	// Additional Video formats
	".3g2":   "video/3gpp2",
	".h261":  "video/h261",
	".h263":  "video/h263",
	".h264":  "video/h264",
	".jpgv":  "video/jpeg",
	".jpm":   "video/jpm",
	".jpgm":  "video/jpm",
	".mj2":   "video/mj2",
	".mjp2":  "video/mj2",
	".mp4v":  "video/mp4",
	".mpg4":  "video/mp4",
	".mpe":   "video/mpeg",
	".m1v":   "video/mpeg",
	".m2v":   "video/mpeg",
	".ogv":   "video/ogg",
	".pyv":   "video/vnd.ms-playready.media.pyv",
	".viv":   "video/vnd.vivo",
	".f4v":   "video/x-f4v",
	".fli":   "video/x-fli",
	".movie": "video/x-sgi-movie",

	// Additional Font formats
	".pfa": "application/x-font-type1",
	".pfb": "application/x-font-type1",
	".pfm": "application/x-font-type1",
	".afm": "application/x-font-type1",
	".psf": "application/x-font-linux-psf",
	".pcf": "application/x-font-pcf",
	".snf": "application/x-font-snf",

	// Additional Message formats
	".mht":   "message/rfc822",
	".mhtml": "message/rfc822",
	".nws":   "message/rfc822",
}
