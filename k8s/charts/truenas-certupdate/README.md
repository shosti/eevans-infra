# TrueNAS Certificate Updater

This is a temporary hack until TrueNAS gets built-in support for Cloudlfare DNS
issuers (see https://jira.ixsystems.com/browse/NAS-104912). See
https://github.com/danb35/deploy-freenas for config reference (certificates will
be mounted at /certs). The key in `configSecret` should be `deploy_config`.
