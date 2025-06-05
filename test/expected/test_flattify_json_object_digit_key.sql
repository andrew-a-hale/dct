select
	json."0"[0]::decimal
	, json."0"[1]::decimal
from (select '{"0": [ 1, 2 ]}'::json as json)
