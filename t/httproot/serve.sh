#!/bin/sh
ehfs -a :/s-anzeige:./kleinanzeigen \
     -a :/api/v1/prod-ads/images/fc:./img \
     -l localhost:8080 -I index.html
