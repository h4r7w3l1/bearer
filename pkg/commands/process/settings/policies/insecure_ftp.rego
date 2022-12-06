package bearer.insecure_ftp

import future.keywords

sensitive_data_group_uuid := "f6a0c071-5908-4420-bac2-bba28d41223e"

has_sensitive_data if {
    some data_type in input.dataflow.data_types
    some data_category in input.data_categories
    data_category.uuid == data_type.category_uuid
    data_category.group_uuid == sensitive_data_group_uuid
}

medium[item] {
    has_sensitive_data == true

    some detector in input.dataflow.risks
    detector.detector_id == input.policy_id

    location = detector.locations[_]
    item := {
         "category_group": "sensitive data",
         "filename": location.filename,
         "line_number": location.line_number,
         "parent_line_number": location.parent.line_number,
         "parent_content": location.parent.content
    }
}