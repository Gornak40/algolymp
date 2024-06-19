package vydra

type File struct {
	Path string `xml:"path,attr"`
	Type string `xml:"type,attr"`
}

type Source struct {
	Path string `xml:"path,attr"`
	Type string `xml:"type,attr"`
}

type Executable struct {
	Source Source `xml:"source"`
}

type Solution struct {
	Tag    string `xml:"tag,attr"`
	Source Source `xml:"source"`
}

type Statement struct {
	Charset  string `xml:"charset,attr"`
	Language string `xml:"language,attr"`
	Path     string `xml:"path,attr"`
	Type     string `xml:"type,attr"`
}

type Assets struct {
	Solutions struct {
		Solutions []Solution `xml:"solution"`
	} `xml:"solutions"`
}

type Files struct {
	Resources struct {
		Files []File `xml:"file"`
	} `xml:"resources"`
	Executables struct {
		Executables []Executable `xml:"executable"`
	} `xml:"executables"`
}

type ProblemXML struct {
	Revision   int    `xml:"revision,attr"`
	ShortName  string `xml:"short-name,attr"`
	Assets     Assets `xml:"assets"`
	Files      Files  `xml:"files"`
	Statements struct {
		Statements []Statement `xml:"statement"`
	} `xml:"statements"`
}
