#!/bin/bash

# hard-wraps the "help" command text into the
# codebase to ensure consistency

set -e

target=source/data_help.go

printf "$(cat text/header_license.txt)\n\n" > $target

printf "// this file was generated by tool/embed.sh: don't modify!\n\npackage main\n\nfunc help(arg string) string {\n\tswitch arg {\n" >> $target

for f in text/help*.txt; do
	name=$(basename ${f%%.txt})
	name=${name/*_}
	data=$(fold -w 64 -s $f)

	printf "\t\tcase \"$name\":\n\t\t\treturn \`\n$data\n\`\n" >> $target
done

printf "\t}\n\treturn \"\"\n}" >> $target