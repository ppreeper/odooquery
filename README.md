# odooquery

![Build Status](https://github.com/ppreeper/odooquery/actions/workflows/go.yml/badge.svg)

Odoo query CLI based on the OdooRPC library calls

## Usage

```bash
Usage:
  odooquery <system> <model> [flags]

Flags:
  -c, --count           count records
  -f, --fields string   fields field1,field2,...fieldN
  -F, --filter string   filter "[('field', 'op', value), ...]"
  -h, --help            help for odooquery
  -l, --limit int       limit records, 0 for no limit
  -o, --offset int      offset records, 0 for no offset
```

```bash
odooquery odoosystem res.company
```
