local typedefs = require "kong.db.schema.typedefs"

return {
  name = "lms-oauth2-introspection",
  fields = {
    { consumer = typedefs.no_consumer },
    { protocols = typedefs.protocols_http },
    {
      config = {
        type = "record",
        fields = {
          { introspection_url = typedefs.url({ required = true }) },
          { client_id = { type = "string", required = true, len_min = 1 } },
          { client_secret = { type = "string", required = true, referenceable = true } },
          { expected_issuer = typedefs.url({ required = true }) },
          { expected_audience = { type = "string", required = true, len_min = 1 } },
          {
            required_scopes = {
              type = "array",
              default = {},
              elements = { type = "string", len_min = 1 },
            },
          },
          { require_subject = { type = "boolean", default = false } },
          { require_role = { type = "boolean", default = false } },
          { timeout = { type = "integer", default = 2000, between = { 100, 10000 } } },
        },
      },
    },
  },
}
