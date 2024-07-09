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

## 읽을거리


### [비트코인 구조] 개인키(Private key), 공개키(Public key), 주소(Address) 생성
https://kwjdnjs.tistory.com/55


### [비트코인 구조] P2PKH
https://kwjdnjs.tistory.com/61

### [비트코인 구조] P2SH
https://kwjdnjs.tistory.com/62

### [비트코인 구조] 세그윗(Segwit), Bech32 주소
https://kwjdnjs.tistory.com/63

### 탭루트 슈노르 서명

https://kwjdnjs.tistory.com/86

https://kwjdnjs.tistory.com/87

https://academy.gopax.co.kr/taebruteuran-mueosimyeo-biteukoine-eoddeohge-doumi-doelggayo/

### 비트코인 주소 종류

https://blog.naver.com/wolsajubu/222937013012

네이티브 세그윗(bc1 ~ )과 세그윗 (3~)
https://www.ledger.com/ko/academy/%EC%84%B8%EA%B7%B8%EC%9C%97%EA%B3%BC-%EB%84%A4%EC%9D%B4%ED%8B%B0%EB%B8%8C-%EC%84%B8%EA%B7%B8%EC%9C%97bech32%EC%9D%80-%EC%96%B4%EB%96%A4-%EC%B0%A8%EC%9D%B4%EC%A0%90%EC%9D%B4-%EC%9E%88%EB%82%98

## token pocket privkey

token pocket 주소지원
https://medium.com/token-pocket-kr/%ED%86%A0%ED%81%B0%ED%8F%AC%EC%BC%93-4%EA%B0%80%EC%A7%80-%EC%9C%A0%ED%98%95%EC%9D%98-%EB%B9%84%ED%8A%B8%EC%BD%94%EC%9D%B8-%EC%A3%BC%EC%86%8C-%ED%98%95%EC%8B%9D-%EC%A7%80%EC%9B%90-3ae2edd362d2



## examples

### wallet generate

https://github.com/miguelmota/bitcoin-development-with-go/blob/master/en/wallet-generate/README.md

### transfer coin

https://github.com/miguelmota/bitcoin-development-with-go/blob/master/en/transfer-coin/README.md


## RPC 

https://www.quicknode.com/docs/bitcoin