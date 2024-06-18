package vydra

type File struct {
	Path string `xml:"path,attr"`
	Type string `xml:"type,attr"`
}

type Solution struct {
	Tag    string `xml:"tag,attr"`
	Source struct {
		Path string `xml:"path,attr"`
		Type string `xml:"type,attr"`
	} `xml:"source"`
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
}

type ProblemXML struct {
	Revision  int    `xml:"revision,attr"`
	ShortName string `xml:"short-name,attr"`
	Assets    Assets `xml:"assets"`
	Files     Files  `xml:"files"`
}
