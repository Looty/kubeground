resource "aws_ecrpublic_repository" "app_repo" {
  repository_name = "${var.name}/${var.name}"

  catalog_data {
    about_text = var.name
  }
}

resource "aws_ecrpublic_repository" "operator_repo" {
  repository_name = "${var.name}/quest-operator"

  catalog_data {
    about_text = var.name
  }
}
