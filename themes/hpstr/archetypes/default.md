---
date: {{ .Date }}
title: "{{ replace .Name "-" " " | title }}"
images: 
- "https://og.jln.co/tt/jlns1/{{ replace .Name "-" " " | title | base64Encode | replaceRE "=+$" "" }}"
draft: true
---
