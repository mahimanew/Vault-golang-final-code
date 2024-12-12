

provider "helm" {
  kubernetes {
    config_path = "~/.kube/config" # Path to your kubeconfig file
  }
}

resource "helm_release" "vault" {
  name       = "vault"
  namespace  = "vault"
  chart      = "vault"
  repository = "https://helm.releases.hashicorp.com" # Directly specify the repository URL

  create_namespace = true


  values = [
    file("values.yaml")
  ]
}

