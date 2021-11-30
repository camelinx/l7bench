local _M = { }

local strchar = string.char
local mthrand = math.random

local function get_random_string( length )
    local res = ""

    for i = 1, length do
        res = res .. strchar( mthrand( 97, 122 ) )
    end

    return res
end

local function get_one_json_kv( )
    return "\"" .. get_random_string( 3 ) .. "\":\"" .. get_random_string( 3 ) .. "\"", 11
end

local function get_json_string( length )
    local json_doc = "{"
    local nbytes   = 1

    local kv, kvlen = get_one_json_kv( )
    json_doc        = json_doc .. kv
    nbytes          = nbytes + kvlen

    while( nbytes < length )
    do
        json_doc = json_doc .. ","
        nbytes   = nbytes + 1

        kv, kvlen = get_one_json_kv( )
        json_doc  = json_doc .. kv
        nbytes    = nbytes + kvlen
    end

    json_doc = json_doc .. "}"
    nbytes   = nbytes + 1

    return nbytes, json_doc
end

local function parse_request( )
    local length = 1
    if( nil ~= ngx.var.arg_response_size )
    then
        length = tonumber( ngx.var.arg_response_size )
    elseif( nil ~= ngx.var.http_x_response_size )
	then
        length = tonumber( ngx.var.http_x_response_size )
	end

    if( ( nil == length ) or ( length < 0 ) )
    then
        length = 1
    end

	local code = ngx.HTTP_OK
    if( nil ~= ngx.var.arg_response_code )
	then
        code = tonumber( ngx.var.arg_response_code )
    elseif( nil ~= ngx.var.http_x_response_code )
    then
        code = tonumber( ngx.var.http_x_response_code )
	end

    if( ( nil == code ) or ( code < ngx.HTTP_CONTINUE ) or ( code > 599 ) )
    then
        code = ngx.HTTP_OK
    end

	return length, code
end

function _M.send_plain_response( )
    local length, code = parse_request( )

    ngx.status = code
    ngx.header[ "Content-Type" ]   = "text/plain"
    ngx.header[ "Content-Length" ] = length
    ngx.print( get_random_string( length ) ) 
end

function _M.send_json_response( )
    local length, code = parse_request( )

    ngx.status = code
    ngx.header[ "Content-Type" ] = "application/json"

    local json_length, json_string = get_json_string( length )
    ngx.header[ "Content-Length" ] = json_length
    ngx.print( json_string )
end

return _M
