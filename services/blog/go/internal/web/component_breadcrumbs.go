package web

type breadcrumb struct {
	Path   string
	URL    string
	Name   string
	NoLink bool
}

type breadcrumbs []breadcrumb

func (t breadcrumbs) LDJSON() ldJSONData {
	itemListElement := []ldJSONItem{}
	for i, v := range t {
		url := v.URL
		if v.NoLink {
			url = ""
		}
		itemListElement = append(
			itemListElement,
			ldJSONItem{
				Type:     "ListItem",
				Position: i + 1,
				Name:     v.Name,
				Item:     url,
			},
		)
	}
	return ldJSONData{
		Context:         "https://schema.org",
		Type:            "BreadcrumbList",
		ItemListElement: itemListElement,
	}
}
