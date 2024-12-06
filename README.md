# go-secret-santa

A Go-based Secret Santa script that pairs participants and sends them their Secret Santa assignments via email.

## Prerequisites

Before you begin, ensure you have the following:

1. **Mailgun API Key**: You need a Mailgun API key to send emails. You can sign up for a Mailgun account and get your API key from the Mailgun dashboard.
2. **Email Domain**: You need a domain configured in Mailgun to send emails from. This domain will be used to create the sender's email address.

Currently, Mailgun is the only supported email sender, but support for other email services may be added in the future.

## Installation

### Build from source

1. **Clone the repository:**

    ```sh
    git clone https://github.com/dcmcand/go-secret-santa.git
    cd go-secret-santa
    ```

2. **Install dependencies:**

    ```sh
    go get ./...
    ```

3. **Build the project:**

    ```sh
    go build -o go-secret-santa
    ```

### Download the prebuilt binary

Alternatively, you can download the prebuilt binary from the [GitHub Releases](https://github.com/dcmcand/go-secret-santa/releases) page. Choose the appropriate binary for your operating system and architecture, then download and extract it.

For example, on Linux:

```sh
wget https://github.com/dcmcand/go-secret-santa/releases/download/v1.0.0/go-secret-santa-linux-amd64.tar.gz
tar -xzf go-secret-santa-linux-amd64.tar.gz
```

On macOS:

```sh
wget https://github.com/dcmcand/go-secret-santa/releases/download/v1.0.0/go-secret-santa-darwin-amd64.tar.gz
tar -xzf go-secret-santa-darwin-amd64.tar.gz
```

On Windows, download the `.zip` file and extract it using your preferred method.

## Configuration

1. **Generate the configuration file:**

    If the `config.yaml` file does not exist, you can generate a skeleton configuration file:

    ```sh
    ./go-secret-santa --generate-config
    ```

    This will create a [config.yaml](http://_vscodecontentref_/2) file with the following structure:

    ```yaml
    mailgun:
        apikey: "abc123" # This is the Mailgun API key
    email:
        subject: "Secret Santa" # This is the subject of the Secret Santa email
        address: "santa" # This along with the domain is used to create the email address of the sender
        domain: "example.com" # This is the domain of the email
        sender:
            name: "Santa Claus" # This is the name of the sender for use in the email body
    ```

2. **Generate the participants file:**

    If the [participants.csv](http://_vscodecontentref_/3) file does not exist, you can generate a skeleton participants file:

    ```sh
    ./go-secret-santa --generate-participants
    ```

    This will create a [participants.csv](http://_vscodecontentref_/4) file with the following structure:

    ```csv
    Name,Email,Partner,Interests
    Barney,barney@bedrock.com,Betty,"Bowling, Jokes, Movies"
    Fred,fred@bedrock.com,Wilma,"Bowling, Dinosaurs, Golf"
    Wilma,wilma@bedrock.com,Fred,"Cooking, Gardening, Shopping"
    Betty,betty@bedrock.com,Barney,"Reading, Music, Crafts"
    Pebbles,pebbles@bedrock.com,,"Exploring, Drawing, Sports"
    BamBam,bambam@bedrock.com,,"Rock Music, Cave Painting, Athletics"
    ```

## Usage

1. **Run the Secret Santa script:**

    ```sh
    ./go-secret-santa --participants participants.csv --config config.yaml
    ```

    This will pair the participants and send them their Secret Santa assignments via email.

2. **Dry-run mode:**

    If you want to see the pairings without sending emails, you can use the `--dry-run` flag:

    ```sh
    ./go-secret-santa --participants participants.csv --config config.yaml --dry-run
    ```

3. **Custom email template:**

    If you want to use a custom email template, you can specify the template file with the `--email-template` flag:

    ```sh
    ./go-secret-santa --participants participants.csv --config config.yaml --email-template custom_template.txt
    ```

## Testing

To run the tests, use the following command:

```sh
go test ./...
```

## Contributing

Feel free to open issues or submit pull requests if you have any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/dcmcand/go-secret-santa/blob/main/LICENSE) file for details.