# Snapper
Aimed to be engine for Snapper android app

## Guide
### Setup
- Clone the repo
- Install Go
- Set GOPATH to this repo root
- Set GOBIN to this repo root/bin

### Build
```shell
$ cd {projectroot}/src
```
```shell
$ go install
```

### Usage
```shell
$ cd {projectroot}/bin
```
```shell
$ ./{buildname}
```
It will run on port 8000

## Input Sample Label
endpoint `localhost:8000/label` POST
make sure the header content type is JSON
```console
{
  "webEntities": [
    {
      "entityId": "/m/06rrc",
      "score": 0.85754,
      "description": "Shoe"
    },
    {
      "entityId": "/m/0lwkh",
      "score": 0.49217,
      "description": "Nike"
    },
    {
      "entityId": "/m/019sc",
      "score": 0.3900155,
      "description": "Black"
    },
    {
      "entityId": "/m/019rjn",
      "score": 0.37383,
      "description": "Futsal"
    },
    {
      "entityId": "/g/12dpwwx05",
      "score": 0.31789,
      "description": "Nike Hypervenom"
    },
    {
      "entityId": "/m/027sf8d",
      "score": 0.30877,
      "description": "Nike Mercurial Vapor"
    },
    {
      "entityId": "/m/04lbp",
      "score": 0.25490758,
      "description": "Leather"
    },
    {
      "entityId": "/m/01sdr",
      "score": 0.18925,
      "description": "Color"
    },
    {
      "entityId": "/g/11dylp_1v",
      "score": 0.18442,
      "description": "Sneakers"
    },
    {
      "entityId": "/m/01g5v",
      "score": 0.17937,
      "description": "Blue"
    },
    {
      "entityId": "/m/092sx5",
      "score": 0.1774,
      "description": "Pricing strategies"
    },
    {
      "entityId": "/m/038hg",
      "score": 0.17397,
      "description": "Green"
    },
    {
      "entityId": "/m/05t5gr",
      "score": 0.17268,
      "description": "Promotion"
    },
    {
      "entityId": "/m/06ntj",
      "score": 0.17248,
      "description": "Sports"
    },
    {
      "entityId": "/m/083jv",
      "score": 0.17213,
      "description": "White"
    }
  ],
  "fullMatchingImages": [
    {
      "url": "https://s0.bukalapak.com/img/037091828/m-1000-1000/IMG_20170104_WA0004_scaled.jpg",
      "dumm": 1
    },
    {
      "url": "https://s1.bukalapak.com/img/1583582511/m-1000-1000/14326926_61f6a647_1150_4d16_a4f3_bba7f259b753_444_444.jpg",
      "dumm": 2
    }
  ]
}
```

## Output Sample
```console
{
  "pairs": [
    {
      "keyword": "Nike Hitam Futsal Sepatu",
      "score": 0.52838886
    },
    {
      "keyword": "Nike Futsal Sepatu",
      "score": 0.5745134
    },
    {
      "keyword": "Nike Sepatu",
      "score": 0.674855
    },
    {
      "keyword": "Futsal Sepatu",
      "score": 0.615685
    }
  ]
}
```


## Input Sample Behaviour
endpoint `localhost:8000/behaviour` GET
need 1 variable `weeks` to get how many weeks from today the behaviour will be gathered 
`weeks` is integer
127.0.0.1:8000/behaviour?weeks=2

## Output Sample
```console
[
  {
    "Jenis": "Sepatu",
    "Merk": "Adidas",
    "Waktu": "2017-05-27T10:30:36Z"
  },
  {
    "Jenis": "",
    "Merk": "Adidas",
    "Waktu": "2017-05-27T10:31:42Z"
  },
  {
    "Jenis": "",
    "Merk": "Adidas",
    "Waktu": "2017-05-27T11:56:58Z"
  },
  {
    "Jenis": "Sepatu",
    "Merk": "Nike",
    "Waktu": "2017-05-27T11:57:29Z"
  }
]
```


