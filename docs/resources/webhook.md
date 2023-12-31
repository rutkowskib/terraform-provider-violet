---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "violet_webhook Resource - terraform-provider-violet"
subcategory: ""
description: |-
  Resource to manage Violet webhook
---

# violet_webhook (Resource)

Resource to manage Violet webhook

## Example Usage

```terraform
resource "violet_webhook" "example" {
  event           = "OFFER_UPDATED"
  remote_endpoint = "https://test.com/"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `event` (String) Event webhook will be subscribed to
- `remote_endpoint` (String) Endpoint that webhook will be publishing to

### Read-Only

- `app_id` (Number) App Id of application this webhook belongs to
- `date_created` (String) Creation date of webhook
- `date_last_modified` (String) Date of last modification of the webhook
- `id` (Number) Webhook id
- `status` (String) Status of webhook

## Import

Import is supported using the following syntax:

```shell
# Webhook can be imported using id
terraform import violet_webhook.example 10198
```
