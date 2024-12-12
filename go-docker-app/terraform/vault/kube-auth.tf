provider "vault" {
  address         = "https://3.83.68.194:8202/" # Vault address
  token           = ""                          # Vault token, can be stored in a variable or environment variable
  skip_tls_verify = true
}

provider "aws" {
  region = "us-east-1" # Update with your AWS region
}

# Fetch the EKS cluster details
data "aws_eks_cluster" "cluster" {
  name = "cluster-test"
}

data "aws_eks_cluster_auth" "cluster" {
  name = data.aws_eks_cluster.cluster.name
}

resource "vault_policy" "mysecret" {
  name   = "mysecret"
  policy = <<EOF
path "kv-v2/data/fakeapp/mysecret" {
  capabilities = ["read"]
}
EOF
}

# Define the Kubernetes Auth Backend
resource "vault_auth_backend" "kubernetes" {
  type = "kubernetes"
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority[0].data)
  token                  = data.aws_eks_cluster_auth.cluster.token
}

# Fetch the ServiceAccount token from the Kubernetes secret
data "kubernetes_secret" "sa_token" {
  metadata {
    name      = "vault-sa-token"
    namespace = "vault"
  }
}

# Configure the Kubernetes Auth Backend
resource "vault_kubernetes_auth_backend_config" "kubernetes_config" {
  backend            = vault_auth_backend.kubernetes.path
  token_reviewer_jwt = data.kubernetes_secret.sa_token.data["token"]
  kubernetes_host    = data.aws_eks_cluster.cluster.endpoint
  kubernetes_ca_cert = data.kubernetes_secret.sa_token.data["ca.crt"]
}

# Define the Kubernetes Auth Backend Role
resource "vault_kubernetes_auth_backend_role" "vault_role" {
  role_name                        = "vault-role"
  backend                          = vault_auth_backend.kubernetes.path
  bound_service_account_names      = ["vault-sa"]
  bound_service_account_namespaces = ["vault"]
  token_policies                   = ["default", "mysecret"]
  token_max_ttl                    = 3600 # Max TTL in seconds
}
