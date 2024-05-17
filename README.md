# Kessel
Kessel is a command-line tool written in Go that generates an AWS Console URL for a specified AWS CLI profile. This tool simplifies the process of accessing the AWS Management Console with specific IAM roles and credentials configured in your AWS profiles. Kessel uses the Cobra library for command-line interface functionality.

## Features
Reads AWS profile configurations from the `~/.aws/config` file.
Assumes specified IAM roles using AWS SSO.
Generates federated sign-in URLs for the AWS Management Console.
Supports multiple AWS profiles.
Built with Cobra for a powerful and user-friendly CLI experience.

## Prerequisites
- Go 1.16 or later
- AWS CLI configured with your profiles
- AWS SSO configured in your AWS account

## Installation
Clone the repository:

```
git clone https://github.com/your-username/kessel.git
cd kessel
```
## Build the project:
```
go build -o kessel
```
Move the kessel binary to a directory in your PATH:
```
mv kessel /usr/local/bin/
```
## Usage
To generate an AWS Console URL for a specified profile, run the following command:

```
kessel console-url <profile>
```
Replace <profile> with the name of your AWS CLI profile (e.g., app1).

## Example
```
kessel console-url app1
```
This command will output a URL that you can open in your browser to access the AWS Management Console with the specified profile's credentials.

## Configuration
Ensure that your AWS CLI profiles are configured in the ~/.aws/config file with the necessary SSO details. Here is an example configuration for a profile named app1:

```
[profile app1]
sso_account_id = <your-sso-account-id>
sso_start_url = https://your-sso-start-url
sso_region = us-west-2
sso_role_name = AdministratorAccess
region = us-west-2
```
## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any changes or improvements.

## License
This project is licensed under the MIT License. See the LICENSE file for details.

## Acknowledgements

- AWS SDK for Go
- AWS CLI
- Cobra
