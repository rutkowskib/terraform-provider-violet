resource "violet_webhook" "example" {
  event           = "OFFER_UPDATED"
  remote_endpoint = "https://test.com/"
}
