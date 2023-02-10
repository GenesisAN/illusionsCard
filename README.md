# illusionsCard

Used to obtain the game character card related information



Supported games:

- [x] Koikatsu 

- [ ] Koikatsu Sunshine

Require:

```
github.com/tinylib/msgp
```

### Exampe:

Install:

```shell
go get "github.com/GenesisAN/illusionsCard"
```

Use:

```go
Koikatsu , err := illusionsCard.ReadKK("./Card.png")
if err != nil {
    return
}
```

