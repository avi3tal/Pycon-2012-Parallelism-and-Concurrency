#!/bin/bash

PYTHON2="python"
PYTHON3="/usr/local/python322/bin/python3"
APPDIR="/local/src/Pycon2012ParallelConcurrent"

rmfile() 
{
	fil=$1
	if [ -f $fil ] ;	then
		rm $fil
	fi
}	

# rename time capture file
# TIMECAPSV should be defined in the caller
rentimef()
{
	timecap=$1
	ext=$(date +%Y-%m-%d-%H-%M-%S)
	TIMECAPCSV="$timecap-$ext.csv" 
	mv $timecap $TIMECAPCSV 
}