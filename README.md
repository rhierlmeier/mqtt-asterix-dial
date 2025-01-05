# mqtt-asterisk-dial

mqtt-asterisk-dial is a Go application that connects to an MQTT broker. It obereves one or more MQTT topics. When the values of such a topic has a defined
values then an Asterisk call file is written.

The content of the call files is rendered via a GO template. This template has 
access to values of MQTT topics.

# Configuration

The mqtt-asterisk-dial is configured via a configuration yaml file that
can be configured via the `-config` command line param. Default location is `./conf.yaml`).

You can find a sample configuration file [here](./conf.yaml).



## Project Structure

```
mqtt-asterisk-dial
├── cmd
│   └── mqtt-asterisk-dial
│       └── main.go        # Entry point of the application
├── internal
│   ├── config
│   │   └── config.go      # Configuration settings
│   ├── mqtt
│   │   └── client.go      # MQTT client management
│   └── dial
│       └── dial.go        # Dialing logic
├── go.mod                  # Module definition
└── README.md               # Project documentation
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd mqtt-asterisk-dial
   ```

2. **Install dependencies:**
   ```
   go mod tidy
   ```

3. **Configure the application:**
   Update the configuration settings in `internal/config/config.go` or set environment variables as needed.

## Usage

To run the application, execute the following command:

```
go run cmd/mqtt-asterisk-dial/main.go
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.