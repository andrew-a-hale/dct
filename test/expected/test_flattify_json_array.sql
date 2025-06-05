select
	json[0]::decimal
	, json[1]::decimal
from (select '[ 1, 2 ]'::json as json)
