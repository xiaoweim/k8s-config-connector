/**
 * Copyright 2022 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

```hcl
resource "google_project" "project" {
  provider = google-beta

  project_id      = "tf-test%{random_suffix}"
  name            = "tf-test%{random_suffix}"
  org_id          = "123456789"
  billing_account = "000000-0000000-0000000-000000"
}

resource "google_project_service" "apigee" {
  provider = google-beta

  project = google_project.project.project_id
  service = "apigee.googleapis.com"
}

resource "google_project_service" "compute" {
  provider = google-beta

  project = google_project.project.project_id
  service = "compute.googleapis.com"
}

resource "google_project_service" "servicenetworking" {
  provider = google-beta

  project = google_project.project.project_id
  service = "servicenetworking.googleapis.com"
}

resource "google_project_service" "kms" {
  provider = google-beta

  project = google_project.project.project_id
  service = "cloudkms.googleapis.com"
}

resource "google_compute_network" "apigee_network" {
  provider = google-beta

  name       = "apigee-network"
  project    = google_project.project.project_id
  depends_on = [google_project_service.compute]
}

resource "google_compute_global_address" "apigee_range" {
  provider = google-beta

  name          = "apigee-range"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.apigee_network.id
  project       = google_project.project.project_id
}

resource "google_service_networking_connection" "apigee_vpc_connection" {
  provider = google-beta

  network                 = google_compute_network.apigee_network.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.apigee_range.name]
  depends_on              = [google_project_service.servicenetworking]
}

resource "google_kms_key_ring" "apigee_keyring" {
  provider = google-beta

  name       = "apigee-keyring"
  location   = "us-central1"
  project    = google_project.project.project_id
  depends_on = [google_project_service.kms]
}

resource "google_kms_crypto_key" "apigee_key" {
  provider = google-beta

  name            = "apigee-key"
  key_ring        = google_kms_key_ring.apigee_keyring.id
}

resource "google_project_service_identity" "apigee_sa" {
  provider = google-beta

  project = google_project.project.project_id
  service = google_project_service.apigee.service
}

resource "google_kms_crypto_key_iam_binding" "apigee_sa_keyuser" {
  provider = google-beta

  crypto_key_id = google_kms_crypto_key.apigee_key.id
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"

  members = [
    "serviceAccount:${google_project_service_identity.apigee_sa.email}",
  ]
}

resource "google_apigee_organization" "apigee_org" {
  provider = google-beta

  analytics_region                     = "us-central1"
  project_id                           = google_project.project.project_id
  authorized_network                   = google_compute_network.apigee_network.id
  billing_type                         = "PAYG"
  runtime_database_encryption_key_name = google_kms_crypto_key.apigee_key.id

  depends_on = [
    google_service_networking_connection.apigee_vpc_connection,
    google_project_service.apigee,
    google_kms_crypto_key_iam_binding.apigee_sa_keyuser,
  ]
}

resource "google_apigee_environment" "apigee_environment" {
  provider = google-beta

  org_id       = google_apigee_organization.apigee_org.id
  name         = "tf-test%{random_suffix}"
  description  = "Apigee Environment"
  display_name = "tf-test%{random_suffix}"
  node_config {
    min_node_count = "3"
    max_node_count = "5"
  }
}
```
