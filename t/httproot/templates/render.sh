#!/bin/sh -x
base="../kleinanzeigen"
mkdir -p $base

echo "Generating /s-bestandsliste.html"
p2cli -t index.tpl -i vars.yaml > $base/s-bestandsliste.html

for idx in 0 1; do
    slug=$(cat vars.yaml | yq ".ads[$idx].slug")
    id=$(cat vars.yaml | yq ".ads[$idx].id")
    mkdir -p $base/$slug/$id
    cat vars.yaml | yq ".ads[$idx]" | p2cli -t ad.tpl -f yaml > $base/$slug/$id/index.html
done
