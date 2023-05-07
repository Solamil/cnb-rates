#!/bin/sh
title="Kurz devizového trhu"
git_url="https://github.com/Solamil/cnb-rates"
file="denni_kurz.txt"
url="https://cnb.cz/cs/financni-trhy/devizovy-trh/kurzy-devizoveho-trhu/kurzy-devizoveho-trhu/"$file
dir="rates"
list_file="$dir/list.txt"
date_file="$dir/date.txt"
number_file="$dir/number.txt"
trinity_file="$dir/svata_trojice.txt"
web_dir="web"
index="$web_dir/index.html"

[ -d $web_dir ] || mkdir -pv $web_dir
[ -d $dir ] || mkdir -pv $dir

download_rates() { curl -sLf "$url" -o $file; }

parse_rates() {

	grep "|" $file | cut -d"|" -f4 > $list_file
	head -n 1 $file | cut -d" " -f1 > $date_file
	head -n 1 $file | cut -d"#" -f2 > $number_file
	codes=$(grep -v "^kód" $list_file)

	grep -v "kód" < $list_file | while IFS= read -r code 
	do
		line=$(grep "$code" $file)
		echo "$line" | grep -o "\|[^\|]*$" | tr "," "." > "$dir/$code.txt"
		echo "$line" | cut -d"|" -f3 >> "$dir/$code.txt"
	done

	option_tags=""
	links_code=""
	for i in $codes; do
		value=$(cat "${dir}/${i}.txt")
		option_tags=$option_tags" <option value=\"$i\">$i</option>"
		links_code=$links_code" <a href=\"/?code=$i\"><abbr title=\"$value\">$i</abbr></a>"
	done
	printf "1$ %.2fKč 1€ %.2fKč 1£ %.2fKč" "$(head -n 1 "$dir/USD.txt")" "$(head -n 1 "$dir/EUR.txt")" "$(head -n 1 "$dir/GBP.txt")" > "$trinity_file"

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
	
		<p>$(cat "$dir/svata_trojice.txt")</p>
		<pre>
$(cat $file)		
		</pre>
		<footer style=\"font-size: x-small\">
			<p>Kurz devizového trhu pro obyčejné lidi.</p>
			<p>Webová aplikace pro zobrazení kurzu devizového trhu v dalších formátech. Původ dat <a href=\"$url\">ČNB</a>.</p>
			<a href=\"/json\">JSON</a>
			<a href=\"/list\">list</a>
			<a href=\"/date\"><abbr title=\"$(cat $date_file 2>/dev/null)\">date</abbr></a>
			<a href=\"/number\"><abbr title=\"$(cat $number_file 2>/dev/null)\">number</abbr></a>
			<a href=\"/denni_kurz.txt\">denni_kurz.txt</a>
			<a href=\"/svata_trojice\">Svatá trojice</a>
			<a href=\"/holy_trinity\">Holy Trinity</a>
			<a href=\"/holy_trinity?p\">Pretty</a>
			<p>$links_code</p>
			<form action=\"/\" method=\"GET\">
				<input type=\"number\" name=\"amount\" placeholder=\"množství\" step=\"0.01\">
				<select name=\"code\" id=\"select_currency\">
		$option_tags
				</select>
				<input type=\"submit\" id=\"save_btn\" value=\"OK\" />
			</form>	
			<p> <a href=\"$git_url\">Projekt</a></p>
		</footer>
	</body>
	</html>
	" > $index

}

if [ -f $file ]; then
	current_date=$(date +"%d.%m.%Y")
	if [ "$(cat $date_file)" != "$current_date" ]; then
			download_rates
	fi
else
	download_rates
fi
parse_rates
render_html
