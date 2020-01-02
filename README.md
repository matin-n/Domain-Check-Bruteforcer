# Domain-Check-Bruteforcer

A tool used to determine availability of domains by bruteforcing character combinations. 

The endpoint to determine domain status was found by sniffing HTTP traffic of the [google domain search website](https://domains.google).

## Usage
1. Run software with `tld=` parameter. Options include com, net, org, or any top level domain.

----
## Libraries Used
* [itertools](https://github.com/ernestosuarez/itertools)
* [gjson](https://github.com/tidwall/gjson)
