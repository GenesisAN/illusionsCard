# illusionsCard

Used to obtain the game character card related information



Supported games:

- [x] Koikatsu 

- [x] Koikatsu Studio scenes (including ExtendedSave/KKEx MOD dependencies)

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
import (
	 ic  "github.com/GenesisAN/illusionsCard"
    kk  "github.com/GenesisAN/illusionsCard/KK"
    kks "github.com/GenesisAN/illusionsCard/KKS"
)

```
Use:

```go
card, err := ic.ReadCardFromPath("./Card.png")
if err != nil {
    fmt.Println("读取失败:", err)
    return
}
fmt.Println("卡片路径:", card.GetPath())

// 可选：类型断言以访问具体结构
switch c := card.(type) {
case *kk.KKCharaCard:
    fmt.Println("→ KK角色卡:", c.CharParmeter.Nickname)
    c.PrintZipmodeInfo()

case *kk.KKClothesCard:
    fmt.Println("→ KK服装卡")
    c.PrintZipmodeInfo()

case *kks.SunshineCharaCard:
    fmt.Println("→ KKS角色卡:", c.CharParmeter.Nickname)
    c.PrintZipmodeInfo()

case *kk.KKSceneCard:
	fmt.Println("→ KK工作室卡:", c.Version)
	fmt.Println("内嵌角色数:", len(c.CharaCards))
	c.PrintZipmodeInfo()

default:
    fmt.Printf("未知卡片类型: %T\n", card)
}

```

