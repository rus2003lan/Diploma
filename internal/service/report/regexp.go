package report

var (
	inputRegexps = map[string]string{
		"IsXmlJson":      `<input [^>]*accept=".*(application/xml|\.xml|application/json|\.json)`,
		"IsPlCgi":        `<input [^>]*accept=".*(\.pl|cgi-bin/\*)`,
		"IsImage":        `<input [^>]*accept=".*image/\*`,
		"IsXmlSvg":       `<input [^>]*accept=".*(application/xml|\.xml|image/svg|\.svg)`,
		"IsFile":         `<input [^>]*type="file"`,
		"IsAcceptedFile": `<input [^>]*type="file".*accept="`,
	}

	aRegexps = map[string]string{
		"IsNumber": `^\d+$`,
	}

	scriptRegexps = map[string]string{
		"XMLHttpRequestGET":  `.open\("GET" [^>]*\)`,
		"XMLHttpRequestPOST": `.open\("POST" [^>]*\)`,
		"$.get":              `\$\.get([^>]*)`,
		"$.post":             `\$\.post([^>]*)`,
		"$.getJSON":          `\$\.getJSON([^>]*)`,
		"$.ajax":             `\$\.ajax([^>]*)`,
		"fetch":              ` fetch([^>]*)`,
		"axios.get":          `axios\.get([^>]*)`,
		"superagent.get":     `superagent\.get([^>]*)`,
		"superagent.post":    `superagent\.post([^>]*)`,
		"axios.post":         `axios\.post([^>]*)`,
	}
)