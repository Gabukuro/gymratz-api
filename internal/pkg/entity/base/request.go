package base

type (
	ListQueryParams struct {
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
	}
)

func (p *ListQueryParams) ValidateAndSetDefaults() {
	if p.PerPage == 0 || p.PerPage < 0 {
		p.PerPage = 10
	}

	if p.Page == 0 || p.Page < 0 {
		p.Page = 1
	}

	if p.PerPage > 100 || p.PerPage < 0 {
		p.PerPage = 100
	}
}
