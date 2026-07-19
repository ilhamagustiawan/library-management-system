local cjson = require "cjson.safe"
local http = require "resty.http"

local kong = kong
local ngx = ngx

local credential_headers = {
  "X-Credential-Client-ID",
  "X-Credential-Scope",
  "X-Credential-Sub",
  "X-Credential-Token-Type",
  "X-Credential-Iat",
  "X-Credential-Exp",
  "X-Credential-Iss",
  "X-Credential-Audience",
  "X-Credential-Role",
}

local Plugin = {
  PRIORITY = 1700,
  VERSION = "1.1.0",
}

local rejection_messages = {
  [401] = "Unauthorized",
  [403] = "Forbidden",
  [503] = "Authentication service unavailable",
}

local function reject(status, challenge)
  local headers = { ["Content-Type"] = "application/json" }
  if challenge then
    headers["WWW-Authenticate"] = challenge
  elseif status == 401 then
    headers["WWW-Authenticate"] = 'Bearer realm="library-api"'
  end
  return kong.response.exit(status, { message = rejection_messages[status] }, headers)
end

local function bearer_token()
  local header = kong.request.get_header("authorization")
  if type(header) ~= "string" then
    return nil
  end
  return header:match("^%s*[Bb][Ee][Aa][Rr][Ee][Rr]%s+([^%s]+)%s*$")
end

local function introspect(config, token)
  local client = http.new()
  client:set_timeout(config.timeout)
  return client:request_uri(config.introspection_url, {
    method = "POST",
    body = ngx.encode_args({ token = token, token_type_hint = "access_token" }),
    headers = {
      ["Accept"] = "application/json",
      ["Authorization"] = "Basic " .. ngx.encode_base64(config.client_id .. ":" .. config.client_secret),
      ["Content-Type"] = "application/x-www-form-urlencoded",
    },
    keepalive = true,
  })
end

local function set_string_header(name, value)
  if type(value) == "string" and value ~= "" then
    kong.service.request.set_header(name, value)
  end
end

local function validate_active_response(config, response)
  if type(response.client_id) ~= "string" or response.client_id == "" then
    return "client_id"
  end
  if type(response.token_type) ~= "string" or string.lower(response.token_type) ~= "bearer" then
    return "token_type"
  end
  if type(response.iss) ~= "string" or response.iss ~= config.expected_issuer then
    return "iss"
  end
  if type(response.iat) ~= "number" or type(response.exp) ~= "number" or response.iat > response.exp then
    return "token timestamps"
  end
  return nil
end

local function has_required_scopes(required, granted)
  if #required == 0 then
    return true
  end
  if type(granted) ~= "string" then
    return false
  end

  local granted_set = {}
  for scope in string.gmatch(granted, "%S+") do
    granted_set[scope] = true
  end
  for _, scope in ipairs(required) do
    if not granted_set[scope] then
      return false
    end
  end
  return true
end

local function has_audience(expected, audiences)
  if type(audiences) ~= "table" then
    return false
  end
  for _, audience in ipairs(audiences) do
    if audience == expected then
      return true
    end
  end
  return false
end

local function forward_identity(config, response)
  kong.service.request.clear_header("authorization")
  for _, name in ipairs(credential_headers) do
    kong.service.request.clear_header(name)
  end

  set_string_header("X-Credential-Client-ID", response.client_id)
  set_string_header("X-Credential-Scope", response.scope)
  set_string_header("X-Credential-Sub", response.sub)
  set_string_header("X-Credential-Token-Type", response.token_type)
  set_string_header("X-Credential-Iss", response.iss)
  set_string_header("X-Credential-Audience", config.expected_audience)
  set_string_header("X-Credential-Role", response.role)
  if type(response.iat) == "number" then
    kong.service.request.set_header("X-Credential-Iat", tostring(response.iat))
  end
  if type(response.exp) == "number" then
    kong.service.request.set_header("X-Credential-Exp", tostring(response.exp))
  end
end

function Plugin:access(config)
  local token = bearer_token()
  if not token then
    return reject(401)
  end

  local result, request_error = introspect(config, token)
  if request_error or not result then
    kong.log.err("OAuth introspection request failed")
    return reject(503)
  end
  if result.status ~= 200 then
    kong.log.err("OAuth introspection returned status ", result.status)
    return reject(503)
  end

  local response, decode_error = cjson.decode(result.body)
  if decode_error or type(response) ~= "table" then
    kong.log.err("OAuth introspection returned invalid JSON")
    return reject(503)
  end
  if response.active ~= true then
    return reject(401)
  end
  local invalid_field = validate_active_response(config, response)
  if invalid_field then
    kong.log.err("OAuth introspection active response has invalid ", invalid_field)
    return reject(503)
  end
  if response.exp <= ngx.time() then
    return reject(401)
  end
  if not has_audience(config.expected_audience, response.aud) then
    return reject(401, 'Bearer realm="library-api", error="invalid_token"')
  end
  if config.require_subject and (type(response.sub) ~= "string" or response.sub == "") then
    return reject(403)
  end
  if config.require_role and response.role ~= "member" and response.role ~= "admin" then
    return reject(403)
  end

  local required_scopes = config.required_scopes or {}
  if not has_required_scopes(required_scopes, response.scope) then
    local scope = table.concat(required_scopes, " ")
    return reject(403, 'Bearer realm="library-api", error="insufficient_scope", scope="' .. scope .. '"')
  end

  forward_identity(config, response)
end

return Plugin
