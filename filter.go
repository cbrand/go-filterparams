package filterparams

// Filter is one allowed filter for the given entry.
type Filter struct {
	// Identification is the representation in the query parameter.
	Identification string
}

var (
	FilterEq = &Filter{
		Identification: "eq",
	}
	FilterLt = &Filter{
		Identification: "lt",
	}
	FilterLte = &Filter{
		Identification: "lte",
	}
	FilterGt = &Filter{
		Identification: "gt",
	}
	FilterGte = &Filter{
		Identification: "gte",
	}
	FilterIn = &Filter {
		Identification: "in",
	}
	FilterLike = &Filter {
		Identification: "like",
	}
	FilterILike = &Filter {
		Identification: "ilike",
	}
)
