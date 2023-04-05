# Kurz devizového trhu

- Webová aplikace pro zobrazení aktuálního kurzu devizového trhu dle [České národní banky](https://www.cnb.cz/cs/financni-trhy/devizovy-trh/kurzy-devizoveho-trhu/kurzy-devizoveho-trhu/denni_kurz.txt) v dalších formátech.

- Displaying current devision rates of market by [Czech National Bank](https://www.cnb.cz/cs/financni-trhy/devizovy-trh/kurzy-devizoveho-trhu/kurzy-devizoveho-trhu/denni_kurz.txt) in other formats.
## Installation

```
git clone https://github.com/Solamil/cnb-rates
cd cnb-rates/
go run main.go
crontab -e
```
Add to crontab following

```
40 14 * * 1-5 bash /your/path/to/cnb-rates/rates.sh
```

## Author

Solamil (https://github.com/Solamil/cnb-rates)
