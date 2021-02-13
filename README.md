# currency-exchange
Interview project


build docker image
`docker build -t currency-exchange .`

start it
`docker run -it --rm -p 8080:8080 currency-exchange`

get some results
`http "http://localhost:8080/rates/all/2019-01-11"`

```
{
    "CHF": 0.9829,
    "CNY": 6.7596,
    "JPY": 108.34,
    "KRW": 1117.18,
    "NOK": 8.5245,
    "SEK": 8.9305,
    "THB": 31.93,
    "TWD": 30.79
}
```

get single currency/date
`http "http://localhost:8080/rates/cur/NOK/2019-01-11"`

```
{
    "2019-01-11": 8.5245
}
```

get last currency
`http "http://localhost:8080/rates/cur/NOK"`
```
{
    "2021-01-29": 8.5454
}
```

get currency between dates
`http "http://localhost:8080/rates/cur/NOK/from/2018-12-21/to/2019-01-03"`

```
{
    "2018-12-21": 8.7407,
    "2018-12-26": 8.7725,
    "2018-12-27": 8.7884,
    "2018-12-28": 8.7191,
    "2018-12-31": 8.6519,
    "2019-01-02": 8.6976
}
```

store new value
`curl "http://localhost:8080/rates/cur/NOK/2020-05-12" -d '{"rate":11.2123}'`

that works even for new currency

`curl "http://localhost:8080/rates/cur/AUD/2020-05-12" -d '{"rate":21.3001}'`

check that new currency and value are here
`http "http://localhost:8080/rates/all/2020-05-12"`
```
{
    "AUD": 21.3001,
    "CHF": 0.968,
    "CNY": 7.0816,
    "JPY": 107.33,
    "KRW": 1222,
    "NOK": 11.2123,
    "SEK": 9.7394,
    "THB": 32.08,
    "TWD": 29.88
}
```
