#!/bin/bash          

STR="Hello World!"
echo $STR  

for i in $(ls / -a -h -R| grep openssl); do 
    echo file: $i
done