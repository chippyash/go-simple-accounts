# Chippyash Simple Accounts Client for Go
## github.com/chippyash/go-simple-accounts

Go: 1.18

## What
Provides a Go client for the [Simple Accounts](https://github.com/chippyash/simple-accounts-3) system.
Please refer to that repository for full details.

## How
All the notes and remarks from the original PHP client hold true for the Go client except where noted.

#### Create an Accountant
```go
import (
	"database/sql"
	"github.com/chippyash/go-simple-accounts/sa"
	"github.com/go-sql-driver/mysql"
)
config := mysql.Config{
    User:                 os.Getenv("DBUID"),
    Passwd:               os.Getenv("DBPWD"),
    DBName:               os.Getenv("DBNAME"),
    AllowNativePasswords: true,
    ParseTime:            true,
}
dba, err := sql.Open("mysql", config.FormatDSN())
if err != nil {
	panic(err)
}

accountant = sa.NewAccountant(dba, 0, "GBP")
```
The Go Accountant has an additional parameter, a 3 character currency code. This is held in the Accountant struct for reference
only and is not stored in the database.

#### Create a new Chart
```go
def, err := sa.NewChartDefinition("../tests/_data/personal.xml")
if err != nil {
    panic(err)
}
lastId, err := accountant.CreateChart("Test", "GBP", def)
if err != nil {
    panic(err)
}

```

#### Fetch an existing Chart
```go
//You will have previously saved your chart id somewhere for later retrieval
accountant := sa.NewAccountant(db, chartId, "GBP")
chart, err := accountant.FetchChart()
if err != nil {
    panic(err)
}
```

#### Adding an Account ledger to the COA
```go
nom, _ := sa.NewNominal("1111")
prnt, _ := sa.NewNominal("0000")
err := accountant.AddAccount(nom, sa.NewAcType().Asset(), "foo", &prnt)
if err != nil {
    panic(err)
}
```

#### Deleting an Account ledger from the COA
```go
nom, _ := sa.NewNominal("0001")
err := accountant.DelAccount(nom)
if err != nil {
    panic(err)
}
```

#### Operations on a Chart
##### Get an account
```go
account := chart.GetAccount(sa.MustNewNominal("1000"))
account = chart.GetAccountByName("Liability")
```

##### Get account parent
```go
nominal := chart.GetParentId(sa.MustNewNominal("1000"))
account := chart.GetAccount(chart.GetParentId(sa.MustNewNominal("1000")))
```

##### Testing if an account exists
```go
exists := chart.HasAccount(sa.MustNewNominal("1000"))
```

##### Get the name of the Chart
```go
name := chart.Name()
```

##### Get the currency of the Chart
```go
crcy := chart.Crcy()
```

##### Get a ledger account's values
```go
account := chart.GetAccount(sa.MustNewNominal("1000"))
dr := account.Dr()
cr := account.Cr()
balance, err := account.Balance()
name := account.Name()
acType := account.Type()
```

#### The COA as a Tree
Under the covers, the chart is kept as a [Hierarchy Tree](https://github.com/chippyash/go-hierarchy-tree).  You can
retrieve the tree:
```go
chartTree := chart.Tree() //returns tree.NodeIFace
```

#### Transaction Entries
##### Creating Entries
Two transaction builders are provided:
- SplitTransactionBuilder
- SimpleTransactionBuilder

By default both will build the transaction with the date set to now() and default values
for source, reference and note (i.e. empty values). Both builders support `With...` methods
to set the additional information, e.g.:
```go
txn := sa.NewSplitTransactionBuilder(0).
    WithDate(dt).
    WithNote("foo").
    WithReference(1).
    WithSource("src").
    Build()
```

###### Split Transaction
```go
txn := sa.NewSplitTransactionBuilder(0).
	WithEntries(
        sa.NewEntry(sa.MustNewNominal("1000"), 100, *sa.NewAcType().Dr()),
        sa.NewEntry(sa.MustNewNominal("2000"), 100, *sa.NewAcType().Cr()),
        sa.NewEntry(sa.MustNewNominal("3000"), 100, *sa.NewAcType().Dr()),
        sa.NewEntry(sa.MustNewNominal("4000"), 100, *sa.NewAcType().Cr()),		
	).
	Build()

txn := sa.NewSplitTransactionBuilder(0).
	WithEntry(*sa.NewEntry(sa.MustNewNominal("1000"), 100, *sa.NewAcType().Dr())).
	WithEntry(*sa.NewEntry(sa.MustNewNominal("2000"), 100, *sa.NewAcType().Cr())).
	Build()
```

When creating new transactions, set the id == 0

###### Simple Transaction
```go
txn := sa.NewSimpleTransactionBuilder(0, sa.MustNewNominal("1000"), sa.MustNewNominal("2000"), 100).
    Build()
```

NB. A simple transaction is a split transaction with 2 entries

##### Transaction information
```go
amt, err := txn.GetAmount() //sum(dr + cr) / 2
noms := txn.GetDrAc()  //[]*Nominals
noms := txn.GetCrAc()  //[]*Nominals
simple := txn.IsSimple()  //true if simple else false
balanced := txn.CheckBalance() //true if transaction is balanced else false
```

##### Writing transactions
```go
txnId := accountant.WriteTransactionWithDate(txn, dt)
txnId := accountant.WriteTransaction(txn) //default date to now()  
```

##### Fetching transactions
```go
txn, err := accountant.FetchTransaction(txnId)
entries, err := accountant.FetchAccountJournals("0001")
```

### For Development
#### Setup

- clone repository
- cd to project directory
- `go get ./...`

##### Database
See [Simple Accounts Readme](https://github.com/chippyash/simple-accounts-3/blob/master/README.md) for instructions to build your database.

Support for [Golang Migrate](https://github.com/golang-migrate/migrate) is provided with migration files for the MySql Db variant
in `./db/migrations`. You will need to install golang-migrate.

`go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

Full install instruction for [go-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

- create a database `sa_accounts`
- run the migration

`migrate -path ./db/migrations -database "mysql://<dbuid>:<dbpwd>@tcp(localhost:3306)/sa_accounts?multiStatements=true" up`

replacing `<dbuid>` and `<dbpwd>` with your root user name and password.

Whilst for development, you certainly can continue to use the root user credentials, it is considered bad practice and
you should create a user that has access only to sa_accounts.

If you want to destroy the database content, just use the above command with `down` instead of `up`.


#### Testing
_To run unit tests:_

`go test --tags=unit ./...`

_To run integration tests:_

- Integration tests require a MySql/MariaDb database server running and the simple accounts database set up

`DBUID=<uid> DBPWD=<pwd> DBNAME=<dbname> go test --tags=integration ./...`

replacing `<uid>`, `<pwd>` and `<dbname>` with your credentials

Integration tests will run the unit tests as well


#### Before you do a PR

- update the readme if required
- add new tests for your new/changed code
- run the tests

## References

- [Github](https://github.com/chippyash/go-simple-accounts)
- [PHP Version of Simple Accounts](https://github.com/chippyash/simple-accounts-3)
- [Go Hierarchy Tree used by this library](https://github.com/chippyash/go-hierarchy-tree)
