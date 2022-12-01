# Catalyst 9800 NETCONF Automated AP Provisioning

This repo contains example code to demonstrate automated provisioning of new Catalyst Access Points (AP) that are connected to a 9800 series wireless controller (WLC).

This code will:

- Establish an MQTT subscription & listen for new AP MAC addresses
- Read provided config map of AP MAC to Site, Policy, and RF tags
- Generate new AP tag configurations based on incoming MAC & configured mapping
- Push AP tag config to WLC

**Note:** This is a companion app to the [NETCONF AP monitor](https://github.com/gve-sw/gve_devnet_c9800_netconf_new_ap_monitor) tool. The AP monitoring tool will detect when new APs are brought online and publish them to MQTT, which this app will subscribe to & automatically provision AP tags.

## Contacts

- Matt Schmitz (mattsc@cisco.com)

## Solution Components

- Cisco Catalyst Wireless Access Points & Wireless Lan Controllers
- MQTT Broker

## Installation/Configuration

**1 - Clone repo:**

```bash
git clone <repo_url>
```

**2 - Provide Config file**

This code relies on a JSON configuration file (`config.json`) to provide the required WLC and MQTT targets. This configuration includes a mapping of AP Ethernet MAC addresses to their assigned site, policy, and RF tags.

A sample configuration file has been provided and uses the format below:

```
{
    "wireless-controllers": [
        {
            "name": "",
            "port": 830
        }
    ],
    "ap-tag-map": {
        "<AP-ETHERNET-MAC>": {
            "site-tag": "",
            "policy-tag": "",
            "rf-tag": ""
        },
        "<AP-ETHERNET-MAC>": {
            "site-tag": "",
            "policy-tag": "",
            "rf-tag": ""
        }
    },
    "mqtt": {
        "broker": "",
        "port": 1883,
        "client-id": "go_mqtt_client",
        "topic": "wireless/ap"
    }
}
```

**3 - Provide WLC Credentials:**

WLC login credentials are provided as environment variables:

```
export WLC_USER=
export WLC_PASSWORD=
```

**4 - Build executable:**

```bash
go build -o netconf-ap-provision
```

## Usage

Run the application with the following command:

```
./netconf-ap-provision
```

## Docker

A `Dockerfile` has been provided for easier deployment of this application. The container can be built and deployed using the following steps:

**1 - Clone repo:**

```bash
git clone <repo_url>
```

**2 - Build the container:**

```
docker build --tag netconf-ap-provision .
```

**3 - Run the container:**

```
docker run -e WLC_USER=<user> -e WLC_PASSWORD=<password> -v <path-to-config.json>:/app/config.json -d netconf-ap-provision
```

# Screenshots

**Example of app execution:**

![/IMAGES/ap-provision.png](/IMAGES/ap-provision.png)

### LICENSE

Provided under Cisco Sample Code License, for details see [LICENSE](LICENSE.md)

### CODE_OF_CONDUCT

Our code of conduct is available [here](CODE_OF_CONDUCT.md)

### CONTRIBUTING

See our contributing guidelines [here](CONTRIBUTING.md)

#### DISCLAIMER

<b>Please note:</b> This script is meant for demo purposes only. All tools/ scripts in this repo are released for use "AS IS" without any warranties of any kind, including, but not limited to their installation, use, or performance. Any use of these scripts and tools is at your own risk. There is no guarantee that they have been through thorough testing in a comparable environment and we are not responsible for any damage or data loss incurred with their use.
You are responsible for reviewing and testing any scripts you run thoroughly before use in any non-testing environment.
