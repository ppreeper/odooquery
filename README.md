![Build Status](https://github.com/ppreeper/odooquery/actions/workflows/go.yml/badge.svg)
[![coverage](https://pages.github.com/ppreeper/odooquery/coverage.svg)](https://pages.github.com/ppreeper/odooquery/coverage)

# odooquery

Odoo query CLI based on the OdooJRPC library calls

## Usage

```bash
Usage of odooquery:
  -count
        count records
  -d string
        odoo database
  -fields string
        fields
  -filter string
        filter
  -host string
        odoo host specified in config.yml (default "prod")
  -limit int
        limit
  -model string
        model
  -offset int
        offset
```

```bash
odooquery -host odoo -model "res.company"
```
