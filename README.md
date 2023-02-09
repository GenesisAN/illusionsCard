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

```go
Koikatsu , err := illusionCard.ReadKK("./Card.png")
if err != nil {
    return
}
```

