select
	values[1].a::decimal
	, values[1].b.a::decimal
	, values[2].a::decimal
	, values[2].b[1].values[1]::decimal
	, values[2].b[1].values[2]::decimal
	, values[2].b[1].values[3]::decimal
	, values[2].b[2].values[1]::decimal
	, values[2].b[2].values[2]::decimal
	, values[2].b[2].values[3]::decimal
	, values[3].a::decimal
	, values[3].b::decimal
	, values[4].a::decimal
	, values[4].b::decimal
	, values[5].a::decimal
	, values[5].b::decimal
from 'test/resources/flattify.ndjson';
