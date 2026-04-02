** This version is not yet released; this document is gathering release notes
for the future release **

*   Special shout-outs to @anhdle-sso, @cheftako, @codebot-robot, @katrielt, @maqiuyujoyce, and @xiaoweim for their contributions to this release.

## New Alpha Resources (Direct Reconciler):

*   [`AlloyDBUser`](https://cloud.google.com/config-connector/docs/reference/resource-docs/alloydb/alloydbuser)
    *   Manage AlloyDB users.

*   [`AccessContextManagerServicePerimeter`](https://cloud.google.com/config-connector/docs/reference/resource-docs/accesscontextmanager/accesscontextmanagerserviceperimeter)
    *   Manage Access Context Manager Service Perimeters.

*   [`ComputeTargetHTTPSProxy`](https://cloud.google.com/config-connector/docs/reference/resource-docs/compute/computetargethttpsproxy)
    *   Manage Compute Target HTTPS Proxies.

*   [`ParameterManagerParameterVersion`](https://cloud.google.com/config-connector/docs/reference/resource-docs/parametermanager/parametermanagerparameterversion)
    *   Manage Parameter Manager Parameter Versions.

## Reconciliation Improvements:

*   [`BigQueryAnalyticsHubDataExchange`](https://cloud.google.com/config-connector/docs/reference/resource-docs/bigqueryanalyticshub/bigqueryanalyticshubdataexchange)
    *   Added structured reporting diff to enhance diff visibility in logs.

## Bug Fixes:

*   [`TagsTagKey`](https://cloud.google.com/config-connector/docs/reference/resource-docs/tags/tagstagkey) & [`TagsTagValue`](https://cloud.google.com/config-connector/docs/reference/resource-docs/tags/tagstagvalue)
    *   Fixed an issue by properly handling `ALREADY_EXISTS` errors during creation.

*   `ConfigConnector` Core
    *   The `preview` tool now supports alternative controller comparison in preview mode, with graceful timeout handling and execution comments.
