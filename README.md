# tracklet
A CLI portfolio tracker intended to keep track of investments.

### Supported exchanges :
- Binance (WIP)
- Coinbase (Todo)

# Installation
`make install`\
It will build the project (*golang 1.18* required), create the config file in your homedir, and copy the binary into
your path.

# Setup
Modify your config file under `$HOME/.tracklet/tracklet.yaml` with the necessary information.

# Usage
`tracklet [exchange] process` : Gather data from binance account and save to file to allow wallet calculation.\
`tracklet [exchange] wallet` : Perform calculation to build wallet data.

# Uninstall
To remove **tracklet** : `make uninstall`.
> Note that the configuration file is backup under `/tmp/tracklet.yaml`, just in case during the process.
