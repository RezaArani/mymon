https://www.buymeacoffee.com/RezaArani

 
# mymon, A very simple website monitoring tool

A simple and efficient web monitoring platform developed in Go. This platform is designed for low resource consumption and high performance, supporting multiple operating systems. It uses HTTP headers to check the health of target websites and can trigger HTTP webhooks in case of any failures. The platform also supports custom headers and cookies.

## Features

- **Low Resource Consumption**: Minimal overhead and optimized for performance.
- **High Performance**: Fast checks to ensure real-time monitoring.
- **Multi Operating System Support**: Compatible with various operating systems.
- **HTTP Health Checks**: Uses HTTP headers to determine website health.
- **Webhook Integration**: Calls an HTTP webhook in case of failures.
- **Custom Headers and Cookies**: Supports custom HTTP headers and cookies for checks.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Contributing](#contributing)
- [License](#license)
- [Buy Me a Coffee](https://www.buymeacoffee.com/RezaArani)

## Installation

### Prerequisites

- [Go](https://golang.org/dl/) 1.21 or later

### Steps

1. Clone the repository:

```bash
git clone https://github.com/RezaArani/mymon
cd mymon
```

2. Build the application:

```bash
go build mymon.go
```

3. Run the application:

```bash
./mymon
```

### Notes

- The platform does not require any DBMS or database system.
- It can work standalone without any additional dependencies.
- After compilation, the application does not require any installation steps and can run independently.

## Usage

The mymon platform can be configured using a JSON configuration file. This file should specify the targets to be monitored and the settings for each target.

### Example Configuration File

Create a `config.json` file with the following structure:

```json
     {
        "HTTPBinding":":800",
        "WebhookHttpMethod":"GET",
        "FollowRedirects":true,
        "ClientCertFile":"",
        "FourOnly":false,
        "SixOnly":false,
        "MaxRedirects":10,
        "Debug":false,
        "HttpHeaders":[
            "User-Agent:Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 Safari/537.36",
            "Accept:text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
            "Sec-Fetch-Site:none",
            "Sec-Fetch-Mode:navigate",
            "Sec-Fetch-User:?1",
            "sec-ch-ua:\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"96\", \"Google Chrome\";v=\"96\"",
            "sec-ch-ua-platform:\"Windows\"",
            "Accept-Encoding:gzip",
            "Accept-Language:en-US,en;q=0.9",
            "Sec-Fetch-Dest:document",
            "sec-ch-ua-mobile:?0",
            "Connection:keep-alive"
        ],
        "WebhookURL":"http://webhookurl.cc?TestInfo=~TESTINFO~&Url=~URL~&ErrorMessage=~ERROR~",
        "Websites":[
            {
                "url":"https://semmapas.com",
                "interval":5,
                "alarmfailtestcount":3,
                "httpmethod":"GET",
                "ignorecertificate":true
            } 
            
            
        ]
    } 
 
```

- `Websites`: Array of target website URLs to monitor.
- `Websites.interval`: Time interval (in seconds) between health checks.
- `Websites.alarmfailtestcount`: count of failure detection before calling webhook.
- `headers`: Custom headers to include in the health check requests.
- `cookies`: Custom cookies to include in the health check requests.
- `WebhookURL`: The URL of the webhook to call in case of a failure.

### Running with Configuration

To run the platform with a specific configuration file:

```bash
./main  config.json
```
 
## Contributing

Contributions are welcome! Please submit a pull request or open an issue to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/your-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin feature/your-feature`)
5. Create a new Pull Request

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
