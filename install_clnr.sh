#!/bin/bash
git clone https://github.com/arcxevodov/clnr;
cd clnr;
make;
sudo make install;
cd ..;
sudo rm -r clnr;