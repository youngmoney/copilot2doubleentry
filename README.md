# copilot2doubleentry

Convert exported Copilot transactions to double entry CSV files.

## Example

``` bash
copilot2doubleentry --firstDay 2025-01-01 --lastDay 2025-01-31 --config config.json transactions.csv > converted.csv
```

This can be used in programs such as GnuCash or ledger.

## Config

Because Copilot does not keep enough information for double entry, it
must be added for Income and Internal Transfer categories. Additionally,
it supports the ability to override the category of an Expense which is
useful for excluded categories.

``` json
"overrides": {
    "expense": [
        {
            "match": {
                "description": "re2 matching the name/description",
                "account": "re2 matching the account"
                "outgoing": true
            }
            "account": "optional, category or original account override"
            "split_account": "optional, account or source/destination"
            "always_pair": "if true, always pair with the inverse transaction"
        }
    ]
    "income": []
    "transfer": []
}
```

First match wins.

When matching:

-   the `account` is never the category, always the account field.
-   `outgoing` is true when the money is leaving the account

When overriding Expenses:

-   `account` replaces the category
-   `split_account` replaces the account in the generated double entry

When overriding Income:

-   `account` replaces the account
-   `split_account` replaces the placeholder `INCOME INCOME` in the
    generated double entry

When overriding Internal Transfer:

-   `account` replaces the account
-   `split_account` replaces the placeholder `TRANSFER OUT` or
    `TRANSFER IN` in the generated double entry

`always_pair` is only used for Internal Transfers. It will find a
matching amount transaction across all accounts that has the inverse (ie
10 and -10) within 5 days. Otherwise the program will fail.

## Caveats

-   This program assumes (when printing the CSV file) that all decimals
    are 2 digit precision.
-   Tags are converted to `#hash_tag #style` in the `Notes` field.
