local next_response
local next_error
local request_header
local exit_response
local cleared_headers
local upstream_headers
local request_options

package.loaded["resty.http"] = {
  new = function()
    return {
      set_timeout = function() end,
      request_uri = function(_, _, options)
        request_options = options
        return next_response, next_error
      end,
    }
  end,
}

ngx = {
  encode_args = function() return "encoded-form" end,
  encode_base64 = function() return "encoded-basic" end,
  time = function() return 100 end,
}

rawset(_G, "kong", {
  request = {
    get_header = function() return request_header end,
  },
  response = {
    exit = function(status, body, headers)
      exit_response = { status = status, body = body, headers = headers }
      return exit_response
    end,
  },
  service = {
    request = {
      clear_header = function(name) cleared_headers[name] = true end,
      set_header = function(name, value) upstream_headers[name] = value end,
    },
  },
  log = { err = function() end },
})

local handler = require "kong.plugins.lms-oauth2-introspection.handler"
local config = {
  introspection_url = "http://auth-service:8081/oauth/introspect",
  client_id = "kong-gateway",
  client_secret = "gateway-secret",
  expected_issuer = "http://auth-service",
  required_scopes = { "library:read" },
  require_subject = true,
  timeout = 2000,
}

local function reset()
  next_response = nil
  next_error = nil
  request_header = nil
  exit_response = nil
  cleared_headers = {}
  upstream_headers = {}
  request_options = nil
end

local function assert_equal(want, got, message)
  if want ~= got then
    error((message or "values differ") .. ": want=" .. tostring(want) .. " got=" .. tostring(got))
  end
end

reset()
handler:access(config)
assert_equal(401, exit_response.status, "missing bearer status")

reset()
request_header = "Bearer opaque-token"
next_response = { status = 200, body = [[{"active":false}]] }
handler:access(config)
assert_equal(401, exit_response.status, "inactive token status")

reset()
request_header = "Bearer opaque-token"
next_response = {
  status = 200,
  body = [[{"active":true,"client_id":"nextjs","scope":"library:read library:write","sub":"user-123","token_type":"Bearer","iat":1,"exp":200,"iss":"http://auth-service"}]],
}
handler:access(config)
assert_equal(nil, exit_response, "active token exit")
assert_equal("Basic encoded-basic", request_options.headers["Authorization"], "introspection basic auth")
assert_equal("encoded-form", request_options.body, "introspection form body")
assert_equal(true, cleared_headers["authorization"], "bearer removed")
assert_equal(true, cleared_headers["X-Credential-Client-ID"], "spoofed client header removed")
assert_equal("nextjs", upstream_headers["X-Credential-Client-ID"], "client header")
assert_equal("user-123", upstream_headers["X-Credential-Sub"], "subject header")
assert_equal("library:read library:write", upstream_headers["X-Credential-Scope"], "scope header")

reset()
request_header = "Bearer opaque-token"
next_response = {
  status = 200,
  body = [[{"active":true,"client_id":"nextjs","scope":"library:read","token_type":"Bearer","iat":1,"exp":200,"iss":"http://auth-service"}]],
}
handler:access(config)
assert_equal(403, exit_response.status, "subjectless token status")

reset()
request_header = "Bearer opaque-token"
next_response = {
  status = 200,
  body = [[{"active":true,"client_id":"nextjs","scope":"library:write","sub":"user-123","token_type":"Bearer","iat":1,"exp":200,"iss":"http://auth-service"}]],
}
handler:access(config)
assert_equal(403, exit_response.status, "insufficient scope status")
assert_equal(
  'Bearer realm="library-api", error="insufficient_scope", scope="library:read"',
  exit_response.headers["WWW-Authenticate"],
  "insufficient scope challenge"
)

reset()
request_header = "Bearer opaque-token"
next_response = {
  status = 200,
  body = [[{"active":true,"client_id":"nextjs","scope":"library:read","sub":"user-123","token_type":"Bearer","iat":1,"exp":200,"iss":"http://wrong-issuer"}]],
}
handler:access(config)
assert_equal(503, exit_response.status, "unexpected issuer status")

reset()
request_header = "Bearer opaque-token"
next_response = {
  status = 200,
  body = [[{"active":true,"client_id":"nextjs","scope":"library:read","sub":"user-123","token_type":"Bearer","iat":1,"exp":99,"iss":"http://auth-service"}]],
}
handler:access(config)
assert_equal(401, exit_response.status, "expired token status")

reset()
request_header = "Bearer opaque-token"
next_error = "connection refused"
handler:access(config)
assert_equal(503, exit_response.status, "introspection outage status")

reset()
request_header = "Bearer opaque-token"
next_response = { status = 200, body = "not-json" }
handler:access(config)
assert_equal(503, exit_response.status, "malformed introspection status")

print("handler_spec: ok")
