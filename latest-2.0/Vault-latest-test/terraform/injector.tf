resource "helm_release" "vault" {
  name       = "vault"
  namespace  = "vault"
  chart      = "vault"
  repository = "https://helm.releases.hashicorp.com" # Directly specify the repository URL

  create_namespace = true

  values = [
    file("values.yaml")
  ]

  depends_on = [kubernetes_namespace.vault] # Ensure namespace is created before Vault
}
