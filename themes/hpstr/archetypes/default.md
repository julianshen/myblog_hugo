---
date: {{ .Date }}
title: "{{ replace .Name "-" " " | title }}"
images: 
- "https://og.jln.co/jlns1/{{ replace .Name "-" " " | title | base64Encode | replaceRE "=+$" "" | replaceRE "\\+" "-" | replaceRE "/" "_"}}"
draft: true
---
