# go-filterparams #

go-filterparams is a library for analyzing URL parameters for filtering purposes in a backend. It provides a syntax
to map SQL-like queries on top of the query parameters and parses it into go structs.

This is a helper library for providing filter collection APIs. The primary use case for developing the library is to  
use it with a REST-API which uses the [JSONAPI](http://jsonapi.org/) standard. Because of this the syntax is completely
compatible with the standard and encapsulates everything in the `filter` query parameter.

## Example ##

Given the URL (non URL escaped for better readability):
```
/users?filter[param][name][like][no_default_name]=doe&filter[param][first_name]=doe%&filter[binding]=(!no_brand_name&first_name)&filter[order]=name&filter[order]=desc(first_name)
```

It can be parsed by the given function:

```golang
import (
  "net/url"

  "github.com/cbrand/go-filterparams"
  "github.com/cbrand/go-filterparams/filter"
)

def GetFilterArguments(toParseURL *url.URL) (*filterparams.QueryData, error) {
  builder := filterparams.NewBuilder()
  builder.EnableFilter(filter.FilterEq)
  builder.EnableFilter(filter.FilterLike)
  query := builder.CreateQuery()
  return query.Parse(toParseURL.Query())
}
```

This library would parse the result and return the following structs:

```golang
&QueryData{
  filter: &And{
    Left: &Negate{
      Negated: &Parameter {
        Identification: "no_default_name",
        Name: "name",
        Filter: &FilterLike{},
        Value: "doe%",
      },
    },
    Right: &Parameter {
        Identification: "first_name",
        Name: "first_name",
        Filter: &FilterEq{},
        Value: "doe",
    },
  },
  orders: []*Order{
    &Order{
      orderBy: "name",
      orderDesc: false,
    },
    &Order{
      orderBy: "first_name",
      orderDesc: true,
    },
  },
}
```

The data which is parsed can then be applied to your flavor of backend.

## Syntax ##

All arguments must be prefixed by "filter". It is possible to query for specific data with filters, apply orders to the
 result and to combine filters through AND, NOT and OR bindings.

The syntax builds under the filter parameter a virtual object. The keys of the object are simulated through specifying
`[{key}]` in the passed query parameter. Thus `filter[param]` would point to the param key in the filter object.

### Filter specification ###

The solution supports to query data through the `param` subkey.

```
filter[param][{parameter_name}][{operation}][{alias}] = {to_query_value}
```

The `operation` and `alias` parameters may be omitted. If no `alias` is provided the given parameter name is used for it.
If no `operation` is given, the default one is used (in the example this would be equal).

Example:
```
filter[param][phone_number][like]=001%
```

This would add a filter to all phone numbers which start with "001".

### Filter binding ###

Per default all filters are combined through AND clauses. 
You can change that by specifying the `filter[binding]` argument.

This is where the aliases which you can define come into place. 
The binding provides means to combine filters with AND and OR. Also you are able to negate filters here.

The filters are addressed by their alias or name, if no alias is provided.

If you have a filter `search_for_name`, `search_for_phone_number` and `search_for_account_number` defined you can say 
`search_for_name OR NOT search_for_number AND search_for_account_number` by specifying the following filter:

```
filter[binding]=search_for_name|(!search_for_phone_number&search_for_account_number)
```

Even though the brackets are useless here, you can use them in more complex filters.

The following table summarizes the possible configuration options:
<table>
  <thead>
    <tr>
      <th>Type</th>
      <th>Symbol</th>
      <th>Example</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>AND</td>
      <td>&</td>
      <td>a&b</td>
    </tr>
    <tr>
      <td>OR</td>
      <td>|</td>
      <td>a|b</td>
    </tr>
    <tr>
      <td>NOT</td>
      <td>!</td>
      <td>!a</td>
    </tr>
    <tr>
      <td>Bracket</td>
      <td>()</td>
      <td>(a|b)&c</td>
    </tr>
  </tbody>
</table>

### Ordering ###

To specify a sort order of the results the `filter[order]` parameter may be used. The value can be specified multiple
times. To add ordering you have to provide the name of the parameter which should be ordered, not its alias!

If you want to order by `name`, `first_name` and in reverse order `balance` you can do so by specifying the following query url parameters:

```
filter[order]=name&filter[order]=first_name&filter[order]=desc(balance)
```

As you can see the `desc()` definition can be used to indicate reverse ordering.

## Filter definition ##

Not every backend does or should support all possible filter mechanisms. This is why
the filters which should be accepted by the backend have to be added before processing the query parameters.

The package provides a `QueryBuilder` object, which allows the definition of the allowed values.

```golang
filterEq := Filter{Identification: "eq"}
queryBuilder := filterparams.NewBuilder()
queryBuilder.EnableFilter(filterEq).SetDefaultOperation("eq")
query := queryBuilder.CreateQuery()
```

You can define arbitrary filters. If a filter is requested which isn't supported by the given query, 
an error will be returned.


## Notes ##

- There do no yet exist any public projects which use this library to provide transparent mapping to an underlying 
backend. I plan long-term to add another library which does use this package and provide a way to map it on gorm models. 
If you are planning to do this or use it for other data mapping please contact me and I'll add a reference to it 
the README.
- This package has been implemented due to no publicly available project which parses filter arguments that I know of at the moment.
- Depending on your backend it might not make sense to support all features (ordering, parameter binding) of the
language. You might still want to use it to parse your basic parameters though and ignore the rest.

## Used Libraries ##

For evaluating the query binding the [pigeon](https://github.com/PuerkitoBio/pigeon) parser generator is used 
(Licensed under [BSD 3-Clause](http://opensource.org/licenses/BSD-3-Clause)).

## Filterparam libraries ##
- [filterparams - Python](https://github.com/cbrand/python-filterparams) - Implementation of the language parser in python.
- [filterparams - Ruby](https://github.com/cbrand/ruby-filterparams) - Implementation of the language parser in ruby.
- [filterparams-client - JavaScript](https://github.com/cbrand/js-filterparams-client) - Implementation of an object oriented client request configuration approach in JavaScript.

## License ##

Licensed under the [MIT license](https://opensource.org/licenses/MIT). 

