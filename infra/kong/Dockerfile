FROM kong:3.9.3

USER root
COPY plugins/lms-oauth2-introspection/handler.lua /usr/local/share/lua/5.1/kong/plugins/lms-oauth2-introspection/handler.lua
COPY plugins/lms-oauth2-introspection/schema.lua /usr/local/share/lua/5.1/kong/plugins/lms-oauth2-introspection/schema.lua
USER kong
