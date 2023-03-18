Kup50-TFS

Kup50-TFS is a CLI tool designed to download and manage changes made by a specific author in an Azure DevOps Git repository, generating a ZIP archive with the modified files.

Table of Contents

Features
Requirements
Installation
Usage
License
Features

Connects to an Azure DevOps Git repository using the provided Personal Access Token (PAT) and organization/project details
Filters changes made by a specific author
Downloads and saves modified files
Generates a ZIP archive containing the changes
Requirements

Go 1.20 or higher
Installation

Clone the repository:
bash
Copy code
git clone https://github.com/dominikus1993/kup50-tfs.git
Build the binary:
bash
Copy code
cd kup50-tfs
go build -o kup50-tfs
Add the binary to your $PATH (optional):
bash
Copy code
mv kup50-tfs /usr/local/bin
Usage

bash
Copy code
kup50-tfs --organization <organization_url> --pat <personal_access_token> --project <project_name> [--author <author_name>]
<organization_url> - The URL of your Azure DevOps organization (e.g., https://dev.azure.com/myorg)
<personal_access_token> - Your Azure DevOps Personal Access Token
<project_name> - The name of the project within your organization
<author_name> (optional) - The name of the author whose changes you want to download. If not provided, defaults to "Dominik Kotecki" it's me xD.
License

This project is open source and available under the MIT License.