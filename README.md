# mildew

Seed your [DoD VDP](https://hackerone.com/deptofdefense) recon with the latest list of official "dotmil" domains. mildew crawls all the DoD-maintained website directories to scrape unique `.mil` domains.

Based on the work on [dotmil-domains](https://github.com/esonderegger/dotmil-domains/) project by [esonderegger](https://twitter.com/esonderegger), a research project to map out the DoD's public-facing domain listings:
> There currently isn't a publicly available directory of all the domain names registered under the US military's .mil top-level domain. Such a directory would be useful for people looking to get an aggregate view of military websites and how they are hosted.

## Install
```
go get -u github.com/daehee/mildew
```

## Usage
```
mildew [--roots]
```

## Crawl Sources
The official DoD website directories:
* [U.S. Department of Defense](https://www.defense.gov/Resources/Military-Departments/DOD-Websites/)
* [Air Force](http://www.af.mil/AFSites.aspx)
* [Army](http://www.army.mil/info/a-z/)
* [Navy](https://www.navy.mil/Resources/Navy-Directory/)

## Report Vulnerabilities
Read the DoD Vulnerability Disclosure Policy and submit a vulnerability report at [HackerOne](https://hackerone.com/deptofdefense).
