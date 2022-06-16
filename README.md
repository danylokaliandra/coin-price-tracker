# tracklet (WIP)

# Description
A CLI portfolio tracker.

### Supported exchanges :
- Binance :
    - [How to generate an API Key](https://www.binance.com/en/support/faq/360002502072)
- Kucoin : (WIP)
    - [How to generate an API Key](https://www.kucoin.com/support/360015102174-How-to-Create-an-API)

# Installation
`make install`\
It will build the project (*golang 1.18* required), create the config file in your homedir, and copy the binary into
your path.

# Setup
Modify your config file under `$HOME/.tracklet/tracklet.yaml` with the necessary required information ([see example config file](./config/example.yaml) for required fields).

# Usage
`tracklet [exchange] process` : Gather data from binance account and save to file to allow wallet calculation.\
`tracklet [exchange] wallet` : Perform calculation to build wallet data.

# Uninstall
To remove **tracklet** : `make uninstall`.
> Note that the configuration file is backup under `/tmp/tracklet.yaml`, just in case during the process.
