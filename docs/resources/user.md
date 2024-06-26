---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "border0_user Resource - terraform-provider-border0"
subcategory: ""
description: |-
  The user resource allows you to create and manage Border0 users.
---

# border0_user (Resource)

The user resource allows you to create and manage Border0 users.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) The display name for the user. A friendly name to help distinguish it among other users.
- `email` (String) The email address of the user. User email must be unique per organization.
- `role` (String) The role for the user. Currently valid values include 'admin', 'member', 'read only', and 'client'.

### Optional

- `notify_by_email` (Boolean) Whether to notify the user that they have been added via email. Defaults to true

### Read-Only

- `id` (String) The ID of this resource.
