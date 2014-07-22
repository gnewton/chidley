#!/bin/bash

FILES=`ls  *.kml xml/*.xml xml/*.xml.bz2 xml/*.xml.gz`

for f in $FILES
do
    echo "=================================================================="
    echo $f
    echo ""
    echo ""
	./chidley -V $f
	cd chidleyVerity
	go build
	./chidleyVerity
	cd ..
done


