package config

var ExcludedExtensions = map[string]struct{}{
	".exe": {}, ".dll": {}, ".so": {}, ".dylib": {},
	".zip": {}, ".tar": {}, ".gz": {}, ".bz2": {},
	".mp3": {}, ".wav": {}, ".flac": {}, ".aac": {},
	".mp4": {}, ".mkv": {}, ".avi": {}, ".mov": {},
	".jpg": {}, ".jpeg": {}, ".png": {}, ".gif": {}, ".bmp": {},
	".iso": {}, ".img": {}, ".vmdk": {}, ".qcow2": {},
	".pdf": {}, ".doc": {}, ".docx": {}, ".xls": {}, ".xlsx": {}, ".ppt": {}, ".pptx": {},
	".csv": {}, ".log": {}, ".tmp": {}, ".bak": {},
	".db": {}, ".sqlite": {}, ".mdb": {}, ".accdb": {},
	".psd": {}, ".ai": {}, ".sketch": {}, ".xcf": {},
	".swf": {}, ".flv": {}, ".webm": {},
	".apk": {}, ".ipa": {},
}


var ExcludedDirs = map[string]struct{}{
	".git":          {}, // Exclude all Git-related files
	"node_modules":  {}, // Exclude Node.js dependencies
	"vendor":        {}, // Exclude Go vendor dependencies
	".idea":         {}, // Exclude JetBrains IDE project files
	".vscode":       {}, // Exclude VSCode settings
	"__pycache__":   {}, // Exclude Python cache
	"target":        {}, // Exclude Rust target directory
	"bin":           {}, // Exclude binary build directories
	"obj":           {}, // Exclude object files directory
}


