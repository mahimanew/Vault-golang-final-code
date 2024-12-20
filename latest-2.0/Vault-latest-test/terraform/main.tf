provider "aws" {
  region = "us-east-1" # Update with your AWS region
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority[0].data)
  token                  = data.aws_eks_cluster_auth.cluster.token
}

# Fetch the EKS cluster details
data "aws_eks_cluster" "cluster" {
  name = "cluster-dev-2"
}

data "aws_eks_cluster_auth" "cluster" {
  name = data.aws_eks_cluster.cluster.name
}

# Helm provider to manage releases
provider "helm" {
  kubernetes {
    config_path = "~/.kube/config" # Adjust if needed
  }
}

resource "kubernetes_namespace" "vault" {
  metadata {
    name = "vault"
  }
}

# Create a Helm release for your Go app
resource "helm_release" "go_app" {
  name      = "go-app"
  chart     = "../go-app/" # Placeholder; will override with custom resources
  namespace = "vault"
  values    = [file("../go-app/values.yaml")]

  # Expose the service as a LoadBalancer
  set {
    name  = "service.type"
    value = "LoadBalancer"
  }

  set {
    name  = "service.port"
    value = "8080"
  }

  depends_on = [kubernetes_namespace.vault] # Ensure namespace is created before deploying the app
}

