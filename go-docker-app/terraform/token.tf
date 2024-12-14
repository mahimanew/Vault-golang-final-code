

# Reference the existing Service Account
data "kubernetes_service_account" "existing" {
  metadata {
    name      = "go-app" # Replace with your service account name
    namespace = "vault"  # Adjust the namespace
  }
}


# Create a Secret linked to the existing Service Account
resource "kubernetes_secret" "vault_service_account_token" {
  metadata {
    name      = "vault-auth-token"
    namespace = "vault" # Adjust to match your namespace
    annotations = {
      "kubernetes.io/service-account.name" = data.kubernetes_service_account.existing.metadata[0].name
    }
  }

  type = "kubernetes.io/service-account-token"
}
