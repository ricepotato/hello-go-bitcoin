# hello go bitcoin

## links
https://github.com/btcsuite/btcd

https://hou27.tistory.com/entry/Go%EB%A1%9C-%EB%A7%8C%EB%93%9C%EB%8A%94-%EB%B8%94%EB%A1%9D%EC%B2%B4%EC%9D%B8-part-5-Wallet


https://live.blockcypher.com/

https://github.com/btcsuite/btcd/blob/master/btcutil/address_test.go

### explorer

https://live.blockcypher.com/btc-testnet/

https://blockstream.info/testnet/

## 시작

```
go mod init github.com/ricepotato/hello-go-bitcoin
```

## 의존성 설치

x/crypto module 0.16 버전을 사용해야 함

```
go get -u github.com/btcsuite/btcd

go get golang.org/x/crypto@v0.16.0
```