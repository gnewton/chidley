#!/bin/bash

FILES=`ls  xml/*.xml xml/*.xml.bz2 xml/*.xml.gz`

for f in $FILES
do
    echo "=================================================================="
    echo $f
    echo ""
    echo ""
    ./chidley -W $f > test/Test.go
    cd test
    go build
    time ./test -j > /dev/null
    time ./test -j -s > /dev/null
    time ./test -x > /dev/null
    time ./test -x -s > /dev/null
    cd ..
done


# From: http://www.ncbi.nlm.nih.gov/books/NBK25500/ and from openstreetmap.org
declare -a URLS=('http://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&term=science[journal]+AND+breast+cancer+AND+2008[pdat]' 'http://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pubmed&term=science[journal]+AND+breast+cancer+AND+2008[pdat]&usehistory=y' 'http://eutils.ncbi.nlm.nih.gov/entrez/eutils/esummary.fcgi?db=protein&id=6678417,9507199,28558982,28558984,28558988,28558990' 'http://eutils.ncbi.nlm.nih.gov/entrez/eutils/elink.fcgi?dbfrom=protein&db=protein&id=15718680&term=rat[orgn]+AND+srcdb+refseq[prop]&cmd=neighbor_history' 'http://eutils.ncbi.nlm.nih.gov/gquery?term=mouse[orgn]&retmode=xml' 'http://api06.dev.openstreetmap.org/api/capabilities' 'http://api.openstreetmap.org/api/0.6/trackpoints?bbox=0,51.5,0.25,51.75&page=0')

# for u in "${URLS[@]}"
# do
#     echo "#=================================================================="
#     echo "# $u "
# 	./chidley -V -u -s "" -p "T_" -a "Att_" "$u"
# 	cd chidleyVerity
# 	go build
# 	./chidleyVerity 
# 	cd ..
#     echo "#"
# done


