#!/bin/sh

for entry in "content/post"/*.md
do
  f=`basename $entry`
  docker run -v `pwd`:/content julianshen/ogpp $f > /tmp/$f
  mv /tmp/$f content/post/$f
done