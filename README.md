# illusionsCard

Used to obtain the game character card related information



Supported games:

- [x] Koikatsu 

- [x] Koikatsu Sunshine

Require:

```
github.com/tinylib/msgp
```

### Exampe:

Install:

```shell
go get "github.com/GenesisAN/illusionsCard"
```

Import:

```go
import(
    icb "github.com/GenesisAN/illusionsCard/Base"
    ic "github.com/GenesisAN/illusionsCard"
)
```

Use:

```go
pgb, err := ic.CardTypeRead("./Card.png")
if err != nil {
    fmt.Println(err)
}
switch pgb.Type {
case icb.CT_KK:
    kkinfo, err := ic.ReadKK(pgb)
    fmt.Println(kkinfo, err)
case icb.CT_KKS:
    kksinfo, err := ic.ReadKK(pgb)
    fmt.Println(kksinfo, err)
}
```

