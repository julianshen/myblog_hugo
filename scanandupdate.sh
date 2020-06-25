#!/bin/sh

for entry in "content/post"/*.md
do
  f=`basename $entry`
  docker run -v `pwd`:/blog julianshen/ogpp $f > /tmp/$f
  mv /tmp/$f content/post/$f
done