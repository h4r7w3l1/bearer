critical:
    - rule_dsrid: DSR-1
      rule_display_id: javascript_third_parties_datadog
      rule_description: Do not send sensitive data to Datadog.
      rule_documentation_url: https://curio.sh/reference/rules/javascript_third_parties_datadog
      line_number: 3
      filename: pkg/commands/process/settings/rules/javascript/third_parties/datadog/testdata/unsecure.js
      category_groups:
        - PII
      parent_line_number: 11
      parent_content: client.event("user", "logged_in", {}, user)


--
