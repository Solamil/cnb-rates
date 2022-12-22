#!/bin/sh
title="Kurz devizového trhu"
domain="czk.michalkukla.xyz"
git_url="https://github.com/Solamil/cnb-rates"
file="denni_kurz.txt"
url="cnb.cz/cs/financni-trhy/devizovy-trh/kurzy-devizoveho-trhu/kurzy-devizoveho-trhu/"$file
dir="rates"
list_file="$dir/list.txt"
date_file="$dir/date.txt"
index="web/index.html"

[ -d $dir ] || mkdir -pv $dir

download_rates() { curl -sLf "$url" -o $file; }

parse_rates() {

	grep "|" $file | cut -d"|" -f4 > $list_file
	head -n 1 $file | cut -d" " -f1 > $date_file
	codes=$(grep -v "^kód" $list_file)
	option_tags=""
	links_code=""
	for i in $codes; do
		option_tags=$option_tags" <option value=\"$i\"></option>"
		links_code=$links_code" <a href=\"/?code=$i\">$i</a>"
	done
	option_tags=$option_tags" <option value=\"list\"></option>"
	option_tags=$option_tags" <option value=\"date\"></option>"

	grep -v "kód" < $list_file | while IFS= read -r code 
	do
		line=$(grep "$code" $file)
		echo "$line" | grep -o "\|[^\|]*$" | tr "," "." > "$dir/$code.txt"
		echo "$line" | cut -d"|" -f3 >> "$dir/$code.txt"
	done
}

render_html() { 
	echo "<!DOCTYPE html>
	<html>
	<head>
	<title>$title</title>
	<link rel=\"shortcut icon\" href=\"./pics/favicon.svg\" type=\"image/svg+xml\" />
	<meta charset=\"utf-8\"/>
	<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">
	</head>
	<body>
		<form action=\"/\" method=\"GET\">
			<input type=\"submit\" id=\"save_btn\" style=\"position: absolute;
		     left: -9999px; width: 1px; height: 1px;\" tabindex=\"-1\" />
			<span>$domain/?code=</span>
			<input type=\"text\" name=\"code\" value=\"\" list=\"currencies\" autocomplete=\"off\">
			<span>&amount=</span>
			<input type=\"text\" name=\"amount\" value=\"1\" list=\"currencies\" autocomplete=\"off\">
			<datalist id=\"currencies\">
	$option_tags
			</datalist>
		</form>	
		<pre>
$(cat $file)		
		</pre>
		<footer style=\"font-size: x-small\">
			<p>Kurz devizového trhu</p>
			<p>Webová aplikace pro zobrazení kurzu devizového trhu v dalších formátech. Původ dat <a href=\"https://$url\">ČNB</a>.</p>
			<a href=\"/json\">JSON</a>
			<a href=\"/list\">list</a>
			<a href=\"/date\">date</a>
			<a href=\"/denni_kurz.txt\">denni_kurz.txt</a>
			<p>$links_code</p>
			<p> <a href=\"$git_url\">Projekt</a></p>
		</footer>
	</body>
	</html>
	" > $index

}

if [ -f $file ]; then
	current_date=$(date +"%d.%m.%Y")
	weekday=$(date +"%a")
	if [ "$weekday" != "Sun" ] && [ "$weekday" != "Sat" ] && \
		[ "$(cat $date_file)" != "$current_date" ]; then
			download_rates
	fi
else
	download_rates
fi
parse_rates
render_html
