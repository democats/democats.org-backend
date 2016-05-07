#!/bin/bash

source $HOME/.bash_profile

mv $HOME/to-websockets/blockchains_output/* $HOME/blockchains/

cd $HOME/blockchains
git add -A
git commit -m "Blockchains updated"
git push
