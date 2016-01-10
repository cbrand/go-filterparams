package filterparams

// Filter is one allowed filter for the given entry.
type Filter struct {
	// Identification is the representation in the query parameter.
	Identification string
}

var (
	// FilterEq represents the equal filter.
	FilterEq = &Filter{
		Identification: "eq",
	}
	// FilterLt is a filter for lesser comparison.
	FilterLt = &Filter{
		Identification: "lt",
	}
	// FilterLte is a filter for lesser equal comparison.
	FilterLte = &Filter{
		Identification: "lte",
	}
	// FilterGt is a filter for greater comparison.
	FilterGt = &Filter{
		Identification: "gt",
	}
	// FilterGte is a filter for greater equal comparison.
	FilterGte = &Filter{
		Identification: "gte",
	}
	// FilterIn is a filter for the in comparison.
	FilterIn = &Filter{
		Identification: "in",
	}
	// FilterLike is a filter for the SQL-LIKE clause.
	FilterLike = &Filter{
		Identification: "like",
	}
	// FilterILike is a filter for the SQL-LIKE with ignoring cases.
	FilterILike = &Filter{
		Identification: "ilike",
	}
)
