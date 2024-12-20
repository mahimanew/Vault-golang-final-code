provider "vault" {
  address         = "https://34.203.204.99:8202/"  # Vault address
  token           = "hvs.AzprBAo5evG9wx7kECCYkUCf" # Vault token, can be stored in a variable or environment variable
  skip_tls_verify = true
}

resource "vault_policy" "mysecret" {
  name   = "mysecret"
  policy = <<EOF
path "techassurance/data/data/dev" {
  capabilities = ["read","list"]
}
EOF
}

# Define the Kubernetes Auth Backend
resource "vault_auth_backend" "kubernetes" {
  type = "kubernetes"
}

# Fetch the ServiceAccount token from the Kubernetes secret
data "kubernetes_secret" "sa_token" {
  metadata {
    name      = "vault-sa-token"
    namespace = "vault"
  }
}

resource "vault_kubernetes_auth_backend_config" "kubernetes_config" {
  backend            = vault_auth_backend.kubernetes.path
  token_reviewer_jwt = try(data.kubernetes_secret.sa_token.data["token"], null)
  kubernetes_host    = data.aws_eks_cluster.cluster.endpoint
  kubernetes_ca_cert = base64decode(data.aws_eks_cluster.cluster.certificate_authority[0].data) # Explicitly pass the CA cert here

  depends_on = [data.kubernetes_secret.sa_token]
}



# Define the Kubernetes Auth Backend Role
resource "vault_kubernetes_auth_backend_role" "vault_role" {
  role_name                        = "vault-role"
  backend                          = vault_auth_backend.kubernetes.path
  bound_service_account_names      = ["vault"] # service account name
  bound_service_account_namespaces = ["vault"]
  token_policies                   = ["default", "mysecret"]
  token_max_ttl                    = 3600 # Max TTL in seconds

  depends_on = [vault_kubernetes_auth_backend_config.kubernetes_config] # Ensure the auth backend config is applied before the role
}
